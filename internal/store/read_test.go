package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"testing"
	"testing/fstest"
	"time"

	"pgregory.net/rapid"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
)

const (
	sOtherID = "01JB000000000000000000000B"
	sOtherCS = "01JB00000000000000000000C2"
)

func mustEncode(t testing.TB, e entry.Entry) []byte {
	t.Helper()
	b, err := entry.Encode(e)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// entryPath is where e belongs, relative to .rhiza.
func entryPath(e entry.Entry) string {
	return path.Join(LogDir, e.ChangeSet.String(), e.ID.String()+".json")
}

func file(b []byte) *fstest.MapFile { return &fstest.MapFile{Data: b, Mode: 0o644} }

func dirFile() *fstest.MapFile { return &fstest.MapFile{Mode: fs.ModeDir | 0o755} }

func TestReadLogProblems(t *testing.T) {
	good := testEntry(t)
	goodPath := entryPath(good)

	other := testEntry(t)
	other.ID = mustID(t, sOtherID)

	newerV := []byte(`"v": ` + strconv.Itoa(entry.FormatVersion+1))
	newer := bytes.Replace(mustEncode(t, other), []byte(`"v": 1`), newerV, 1)
	if bytes.Equal(newer, mustEncode(t, other)) {
		t.Fatal(`no "v": 1 in the encoding`)
	}
	nonCanon := bytes.TrimSuffix(mustEncode(t, other), []byte("\n"))

	cs := path.Join(LogDir, sCS)
	tests := []struct {
		name    string
		path    string // the bad file, relative to .rhiza
		data    *fstest.MapFile
		want    string // the Problem's path
		wantErr error
	}{
		{"stray non-ULID folder", "log/not-a-ulid/" + sOtherID + ".json", file(mustEncode(t, other)), "log/not-a-ulid", ErrStray},
		{"stray README in change set", cs + "/README.md", file([]byte("hi\n")), cs + "/README.md", ErrStray},
		{"stray nested folder", cs + "/sub/" + sOtherID + ".json", file(mustEncode(t, other)), cs + "/sub", ErrStray},
		{"stray empty nested folder", cs + "/sub", dirFile(), cs + "/sub", ErrStray},
		{"stray non-ID file stem", cs + "/notes.json", file(mustEncode(t, other)), cs + "/notes.json", ErrStray},
		{"stray plain file in log", "log/README.md", file([]byte("hi\n")), "log/README.md", ErrStray},
		{"mismatch file name and id", cs + "/01JB000000000000000000000Z.json", file(mustEncode(t, other)), cs + "/01JB000000000000000000000Z.json", ErrMismatch},
		{"mismatch wrong change set folder", "log/" + sOtherCS + "/" + sOtherID + ".json", file(mustEncode(t, other)), "log/" + sOtherCS + "/" + sOtherID + ".json", ErrMismatch},
		{"non-canonical file", cs + "/" + sOtherID + ".json", file(nonCanon), cs + "/" + sOtherID + ".json", entry.ErrInvalid},
		{"newer format file", cs + "/" + sOtherID + ".json", file(newer), cs + "/" + sOtherID + ".json", entry.ErrNewerFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := fstest.MapFS{
				RepoFile: file([]byte("{}\n")),
				goodPath: file(mustEncode(t, good)),
				tt.path:  tt.data,
			}
			got := ReadLog(fsys)
			if len(got.Problems) != 1 {
				t.Fatalf("Problems = %v, want exactly one", got.Problems)
			}
			p := got.Problems[0]
			if p.Path != tt.want {
				t.Errorf("Path = %q, want %q", p.Path, tt.want)
			}
			if !errors.Is(p.Err, tt.wantErr) {
				t.Errorf("Err = %v, want %v", p.Err, tt.wantErr)
			}
			if !reflect.DeepEqual(got.Entries, []entry.Entry{good}) {
				t.Errorf("Entries = %+v, want only the good one", got.Entries)
			}
		})
	}

	t.Run("missing log", func(t *testing.T) {
		got := ReadLog(fstest.MapFS{RepoFile: file([]byte("{}\n"))})
		if len(got.Problems) != 1 {
			t.Fatalf("Problems = %v, want exactly one", got.Problems)
		}
		if p := got.Problems[0]; p.Path != LogDir || !errors.Is(p.Err, fs.ErrNotExist) {
			t.Errorf("Problem = %+v, want %q matching fs.ErrNotExist", p, LogDir)
		}
		if got.Entries == nil || len(got.Entries) != 0 {
			t.Errorf("Entries = %#v, want empty and non-nil", got.Entries)
		}
	})
	t.Run("stray log is a file", func(t *testing.T) {
		got := ReadLog(fstest.MapFS{LogDir: file([]byte("x"))})
		if len(got.Problems) != 1 || got.Problems[0].Path != LogDir {
			t.Errorf("Problems = %v, want one at %q", got.Problems, LogDir)
		}
	})
}

func TestReadLogMixed(t *testing.T) {
	a := testEntry(t)
	b := testEntry(t)
	b.ID = mustID(t, sOtherID)
	b.After = []id.ID{a.ID}
	b.Type = entry.TypeNote
	c := testEntry(t)
	c.ChangeSet = mustID(t, sOtherCS)

	fsys := fstest.MapFS{
		// Listed in reverse, to show the order is the reader's own.
		entryPath(c):                file(mustEncode(t, c)),
		entryPath(b):                file(mustEncode(t, b)),
		entryPath(a):                file(mustEncode(t, a)),
		"log/junk":                  file([]byte("x")),
		"log/" + sCS + "/notes.txt": file([]byte("x")),
		"log/" + sOtherCS + "/" + sOtherID + ".json": file([]byte("{")),
	}
	got := ReadLog(fsys)
	// ...C2 sorts before ...CS.
	if want := []entry.Entry{c, a, b}; !reflect.DeepEqual(got.Entries, want) {
		t.Errorf("Entries = %+v\nwant %+v", got.Entries, want)
	}
	var paths []string
	for _, p := range got.Problems {
		if p.Err == nil || p.Err.Error() == "" {
			t.Errorf("Problem %q has no useful error", p.Path)
		}
		paths = append(paths, p.Path)
	}
	want := []string{
		"log/" + sOtherCS + "/" + sOtherID + ".json",
		"log/" + sCS + "/notes.txt",
		"log/junk",
	}
	if !slices.Equal(paths, want) {
		t.Errorf("Problem paths = %v, want %v", paths, want)
	}
}

func TestReadLogSkipsDotFiles(t *testing.T) {
	e := testEntry(t)
	cs := path.Join(LogDir, sCS)
	fsys := fstest.MapFS{
		entryPath(e):                  file(mustEncode(t, e)),
		"log/.x.tmp":                  file([]byte("partial")),
		"log/.DS_Store":               file([]byte{0, 1}),
		"log/.hidden/whatever":        file([]byte("x")),
		cs + "/.x.tmp":                file([]byte("partial")),
		cs + "/." + sID + ".1234.tmp": file(mustEncode(t, e)),
		cs + "/.DS_Store":             file([]byte{0, 1}),
	}
	got := ReadLog(fsys)
	if len(got.Problems) != 0 {
		t.Errorf("Problems = %v, want none", got.Problems)
	}
	if got.Problems == nil {
		t.Error("Problems is nil, want empty")
	}
	if !reflect.DeepEqual(got.Entries, []entry.Entry{e}) {
		t.Errorf("Entries = %+v", got.Entries)
	}
}

// shuffledFS lists every folder in a rapid-chosen order. fs.ReadDir uses
// ReadDir when an FS has one, so the reader sees the shuffled listings.
type shuffledFS struct {
	fstest.MapFS
	t *rapid.T
}

func (s shuffledFS) ReadDir(name string) ([]fs.DirEntry, error) {
	des, err := s.MapFS.ReadDir(name)
	if err != nil {
		return nil, err
	}
	return rapid.Permutation(des).Draw(s.t, "listing "+name), nil
}

func genStoreID() *rapid.Generator[id.ID] {
	return rapid.Custom(func(t *rapid.T) id.ID {
		var i id.ID
		copy(i[:], rapid.SliceOfN(rapid.Byte(), 16, 16).Draw(t, "id bytes"))
		return i
	})
}

// genLog draws a valid log, with some stray files mixed in.
func genLog(t *rapid.T) fstest.MapFS {
	fsys := fstest.MapFS{RepoFile: file([]byte("{}\n"))}
	sets := rapid.SliceOfNDistinct(genStoreID(), 0, 4, func(i id.ID) id.ID { return i }).Draw(t, "change sets")
	for _, cs := range sets {
		ids := rapid.SliceOfNDistinct(genStoreID(), 0, 4, func(i id.ID) id.ID { return i }).Draw(t, "ids")
		for _, eid := range ids {
			after := rapid.SliceOfNDistinct(genStoreID(), 0, 3, func(i id.ID) id.ID { return i }).Draw(t, "after")
			slices.SortFunc(after, id.Compare)
			typ := rapid.SampledFrom(entry.Types()).Draw(t, "type")
			e := entry.Entry{
				V:         1,
				ID:        eid,
				ChangeSet: cs,
				Actor:     genStoreID().Draw(t, "actor"),
				For:       genStoreID().Draw(t, "for"),
				At:        time.Unix(rapid.Int64Range(0, 4102444800).Draw(t, "at"), 0).UTC(),
				After:     after,
				Type:      typ,
				Body:      json.RawMessage("{}"),
			}
			// A reopen entry is spelled only in format 3.
			if typ == entry.TypeReopen {
				e.V = 3
			}
			b, err := entry.Encode(e)
			if err != nil {
				t.Fatal(err)
			}
			fsys[entryPath(e)] = file(b)
		}
	}
	strays := []string{"log/README.md", "log/junk/x.json", "log/.x.tmp", "log/.DS_Store"}
	for _, cs := range sets {
		dir := path.Join(LogDir, cs.String())
		strays = append(strays,
			dir+"/notes.txt", dir+"/sub/x.json", dir+"/bad.json",
			dir+"/"+sID+".json", dir+"/.x.tmp")
	}
	for _, p := range rapid.SliceOfN(rapid.SampledFrom(strays), 0, 5).Draw(t, "strays") {
		if _, ok := fsys[p]; !ok {
			fsys[p] = file([]byte("{}\n"))
		}
	}
	return fsys
}

type textProblem struct{ Path, Err string }

// textLog turns problem errors into their text, since errors built by
// two reads are never the same values.
func textLog(l Log) ([]entry.Entry, []textProblem) {
	ps := []textProblem{}
	for _, p := range l.Problems {
		ps = append(ps, textProblem{p.Path, p.Err.Error()})
	}
	return l.Entries, ps
}

func TestReadLogOrderProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		fsys := genLog(t)
		wantE, wantP := textLog(ReadLog(fsys))
		gotE, gotP := textLog(ReadLog(shuffledFS{fsys, t}))
		if !reflect.DeepEqual(gotE, wantE) {
			t.Fatalf("Entries differ:\n got %+v\nwant %+v", gotE, wantE)
		}
		if !reflect.DeepEqual(gotP, wantP) {
			t.Fatalf("Problems differ:\n got %+v\nwant %+v", gotP, wantP)
		}
		for i := 1; i < len(gotE); i++ {
			a, b := gotE[i-1], gotE[i]
			if c := id.Compare(a.ChangeSet, b.ChangeSet); c > 0 || c == 0 && id.Compare(a.ID, b.ID) >= 0 {
				t.Fatalf("Entries out of order at %d", i)
			}
		}
	})
}

func TestStoreReadLog(t *testing.T) {
	root, s := mustInit(t)
	gen := testGen()
	var want []entry.Entry
	for i := range 6 {
		e := testEntry(t)
		var err error
		if e.ID, err = gen.New(); err != nil {
			t.Fatal(err)
		}
		if i%2 == 1 {
			e.ChangeSet = mustID(t, sOtherCS)
		}
		if err := s.Write(e); err != nil {
			t.Fatal(err)
		}
		want = append(want, e)
	}
	slices.SortFunc(want, func(a, b entry.Entry) int {
		if c := id.Compare(a.ChangeSet, b.ChangeSet); c != 0 {
			return c
		}
		return id.Compare(a.ID, b.ID)
	})

	// A crashed writer's leftover temp file.
	tmp := filepath.Join(root, Dir, LogDir, sCS, "."+sID+".999.tmp")
	if err := os.WriteFile(tmp, []byte("partial"), 0o600); err != nil {
		t.Fatal(err)
	}

	got := s.ReadLog()
	if len(got.Problems) != 0 {
		t.Errorf("Problems = %v, want none", got.Problems)
	}
	if !reflect.DeepEqual(got.Entries, want) {
		t.Errorf("Entries = %+v\nwant %+v", got.Entries, want)
	}
}

func TestStoreReadLogSymlinkedFolder(t *testing.T) {
	root, s := mustInit(t)
	target := filepath.Join(t.TempDir(), sCS)
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, Dir, LogDir, sCS)); err != nil {
		t.Fatal(err)
	}
	got := s.ReadLog()
	want := path.Join(LogDir, sCS)
	if len(got.Problems) != 1 || got.Problems[0].Path != want || !errors.Is(got.Problems[0].Err, ErrStray) {
		t.Fatalf("Problems = %v, want one ErrStray at %q", got.Problems, want)
	}
}

func TestListLog(t *testing.T) {
	good := testEntry(t)
	goodPath := entryPath(good)
	cs := path.Join(LogDir, sCS)
	corrupt := cs + "/" + sOtherID + ".json"

	tests := []struct {
		name      string
		path      string // the extra file, relative to .rhiza
		data      *fstest.MapFile
		wantPaths []string
		wantProb  string // the one Problem's path, or "" for none
		wantErr   error
	}{
		{"stray plain file in log", "log/README.md", file([]byte("hi\n")), []string{goodPath}, "log/README.md", ErrStray},
		{"stray non-ID folder", "log/not-a-ulid/" + sOtherID + ".json", file([]byte("{")), []string{goodPath}, "log/not-a-ulid", ErrStray},
		{"stray nested folder", cs + "/sub/" + sOtherID + ".json", file([]byte("{")), []string{goodPath}, cs + "/sub", ErrStray},
		{"stray non-json file", cs + "/README.md", file([]byte("hi\n")), []string{goodPath}, cs + "/README.md", ErrStray},
		{"stray non-ID file stem", cs + "/notes.json", file([]byte("{")), []string{goodPath}, cs + "/notes.json", ErrStray},
		{"dot file skipped", cs + "/.x.tmp", file([]byte("partial")), []string{goodPath}, "", nil},
		{"corrupt entry listed not read", corrupt, file([]byte("{")), []string{goodPath, corrupt}, "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := fstest.MapFS{
				RepoFile: file([]byte("{}\n")),
				goodPath: file(mustEncode(t, good)),
				tt.path:  tt.data,
			}
			got := ListLog(fsys)
			want := slices.Clone(tt.wantPaths)
			slices.Sort(want)
			if !slices.Equal(got.Paths, want) {
				t.Errorf("Paths = %v, want %v", got.Paths, want)
			}
			if got.Problems == nil {
				t.Error("Problems is nil, want non-nil")
			}
			if tt.wantProb == "" {
				if len(got.Problems) != 0 {
					t.Errorf("Problems = %v, want none", got.Problems)
				}
				return
			}
			if len(got.Problems) != 1 {
				t.Fatalf("Problems = %v, want exactly one", got.Problems)
			}
			if p := got.Problems[0]; p.Path != tt.wantProb || !errors.Is(p.Err, tt.wantErr) {
				t.Errorf("Problem = %+v, want %q matching %v", p, tt.wantProb, tt.wantErr)
			}
		})
	}

	t.Run("missing log", func(t *testing.T) {
		got := ListLog(fstest.MapFS{RepoFile: file([]byte("{}\n"))})
		if got.Paths == nil || len(got.Paths) != 0 {
			t.Errorf("Paths = %#v, want empty and non-nil", got.Paths)
		}
		if len(got.Problems) != 1 || got.Problems[0].Path != LogDir || !errors.Is(got.Problems[0].Err, fs.ErrNotExist) {
			t.Errorf("Problems = %v, want one at %q matching fs.ErrNotExist", got.Problems, LogDir)
		}
	})
}

func TestReadEntriesSubset(t *testing.T) {
	a := testEntry(t)
	b := testEntry(t)
	b.ID = mustID(t, sOtherID)
	c := testEntry(t)
	c.ChangeSet = mustID(t, sOtherCS)
	corrupt := "log/" + sOtherCS + "/" + sOtherID + ".json"
	mismatch := "log/" + sCS + "/01JB000000000000000000000Z.json"
	fsys := fstest.MapFS{
		entryPath(a): file(mustEncode(t, a)),
		entryPath(b): file(mustEncode(t, b)),
		entryPath(c): file(mustEncode(t, c)),
		corrupt:      file([]byte("{")),
		mismatch:     file(mustEncode(t, b)),
	}

	tests := []struct {
		name      string
		paths     []string
		want      []entry.Entry
		wantProbs []string
		wantErrs  []error
	}{
		{"none", nil, []entry.Entry{}, []string{}, nil},
		{"one", []string{entryPath(b)}, []entry.Entry{b}, []string{}, nil},
		{"subset sorted", []string{entryPath(b), entryPath(c)}, []entry.Entry{c, b}, []string{}, nil},
		{"corrupt among them", []string{corrupt, entryPath(a)}, []entry.Entry{a}, []string{corrupt}, []error{entry.ErrInvalid}},
		{"mismatch", []string{mismatch}, []entry.Entry{}, []string{mismatch}, []error{ErrMismatch}},
		{"missing file", []string{"log/" + sCS + "/01JB000000000000000000000Y.json"}, []entry.Entry{}, []string{"log/" + sCS + "/01JB000000000000000000000Y.json"}, []error{fs.ErrNotExist}},
		{"not an entry path", []string{"log/README.md"}, []entry.Entry{}, []string{"log/README.md"}, []error{ErrStray}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadEntries(fsys, tt.paths)
			if !reflect.DeepEqual(got.Entries, tt.want) {
				t.Errorf("Entries = %+v\nwant %+v", got.Entries, tt.want)
			}
			if got.Problems == nil {
				t.Fatal("Problems is nil, want non-nil")
			}
			paths := []string{}
			for i, p := range got.Problems {
				paths = append(paths, p.Path)
				if i < len(tt.wantErrs) && !errors.Is(p.Err, tt.wantErrs[i]) {
					t.Errorf("Problem %q: Err = %v, want %v", p.Path, p.Err, tt.wantErrs[i])
				}
			}
			if !slices.Equal(paths, tt.wantProbs) {
				t.Errorf("Problem paths = %v, want %v", paths, tt.wantProbs)
			}
		})
	}
}

func TestReadLogIsListLogThenReadEntries(t *testing.T) {
	a := testEntry(t)
	b := testEntry(t)
	b.ID = mustID(t, sOtherID)
	c := testEntry(t)
	c.ChangeSet = mustID(t, sOtherCS)
	other := testEntry(t)
	other.ID = mustID(t, sOtherID)
	cs := path.Join(LogDir, sCS)

	fsys := fstest.MapFS{
		RepoFile:                                file([]byte("{}\n")),
		entryPath(a):                            file(mustEncode(t, a)),
		entryPath(b):                            file(mustEncode(t, b)),
		entryPath(c):                            file(mustEncode(t, c)),
		"log/junk":                              file([]byte("x")),
		"log/not-a-ulid/" + sOtherID + ".json":  file([]byte("{")),
		cs + "/notes.txt":                       file([]byte("x")),
		cs + "/sub/x.json":                      file([]byte("x")),
		cs + "/.x.tmp":                          file([]byte("partial")),
		cs + "/01JB000000000000000000000Z.json": file(mustEncode(t, other)),
		"log/" + sOtherCS + "/" + sOtherID + ".json": file([]byte("{")),
	}
	l := ListLog(fsys)
	r := ReadEntries(fsys, l.Paths)
	merged := Log{Entries: r.Entries, Problems: append(slices.Clone(r.Problems), l.Problems...)}
	sortProblems(merged.Problems)

	got := ReadLog(fsys)
	if !reflect.DeepEqual(got, merged) {
		t.Errorf("ReadLog = %+v\nmerge   = %+v", got, merged)
	}
	if len(got.Entries) != 3 || len(l.Problems) != 4 || len(r.Problems) != 2 {
		t.Errorf("entries %d, layout problems %d, read problems %d; want 3, 4, 2",
			len(got.Entries), len(l.Problems), len(r.Problems))
	}
}

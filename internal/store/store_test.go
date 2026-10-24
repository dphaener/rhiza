package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/canon"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
)

const (
	sID    = "01JB000000000000000000000A"
	sCS    = "01JB00000000000000000000CS"
	sActor = "01JB0000000000000000000ACT"
	sFor   = "01JB0000000000000000000F0R"
)

func mustID(t testing.TB, s string) id.ID {
	t.Helper()
	i, err := id.Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	return i
}

func testGen() *id.Generator {
	clock := func() time.Time { return time.Date(2026, 10, 23, 21, 14, 3, 0, time.UTC) }
	return id.NewGenerator(clock, rand.New(rand.NewSource(1)))
}

func testEntry(t testing.TB) entry.Entry {
	return entry.Entry{
		V:         entry.FormatVersion,
		ID:        mustID(t, sID),
		ChangeSet: mustID(t, sCS),
		Actor:     mustID(t, sActor),
		For:       mustID(t, sFor),
		At:        time.Date(2026, 10, 23, 21, 14, 3, 0, time.UTC),
		After:     []id.ID{},
		Type:      entry.TypeOpen,
		Body:      json.RawMessage("{}"),
	}
}

// snapshot records every path under root with its bytes ("" for folders),
// so a test can prove a failed call changed nothing.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if info.IsDir() {
			out[rel+"/"] = ""
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func listNames(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, de := range des {
		names = append(names, de.Name())
	}
	return names
}

func assertNoTemp(t *testing.T, dir string) {
	t.Helper()
	for _, n := range listNames(t, dir) {
		if strings.HasSuffix(n, ".tmp") {
			t.Errorf("temp file %s left in %s", n, dir)
		}
	}
}

func mustInit(t *testing.T) (string, *Store) {
	t.Helper()
	root := t.TempDir()
	s, err := Init(root, testGen())
	if err != nil {
		t.Fatal(err)
	}
	return root, s
}

func TestInitOpenRoundTrip(t *testing.T) {
	root, s := mustInit(t)
	want, err := testGen().New()
	if err != nil {
		t.Fatal(err)
	}
	if s.Root() != root {
		t.Errorf("Root() = %q, want %q", s.Root(), root)
	}
	if r := s.Repo(); r.Format != 1 || r.ID != want || r.Remotes == nil || len(r.Remotes) != 0 {
		t.Errorf("Init Repo() = %+v", r)
	}

	b, err := os.ReadFile(filepath.Join(root, Dir, RepoFile))
	if err != nil {
		t.Fatal(err)
	}
	if err := canon.Check(b); err != nil {
		t.Errorf("repo.json not canonical: %v", err)
	}
	golden := "{\n  \"format\": 1,\n  \"id\": \"" + want.String() + "\",\n  \"remotes\": []\n}\n"
	if string(b) != golden {
		t.Errorf("repo.json =\n%s\nwant\n%s", b, golden)
	}
	if info, err := os.Stat(filepath.Join(root, Dir, LogDir)); err != nil || !info.IsDir() {
		t.Errorf("log dir missing: %v", err)
	}
	assertNoTemp(t, filepath.Join(root, Dir))

	o, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(o.Repo(), s.Repo()) {
		t.Errorf("Open Repo() = %+v, want %+v", o.Repo(), s.Repo())
	}
	if o.Root() != root {
		t.Errorf("Open Root() = %q", o.Root())
	}
}

func TestOpenRemotes(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	b, err := canon.Marshal(map[string]any{"format": 1, "id": sID, "remotes": []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, Dir, RepoFile), b, 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	want := Repo{Format: 1, ID: mustID(t, sID), Remotes: []string{"a", "b"}}
	if !reflect.DeepEqual(s.Repo(), want) {
		t.Errorf("Repo() = %+v, want %+v", s.Repo(), want)
	}
}

func TestInitExists(t *testing.T) {
	t.Run("folder", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, Dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, Dir, "keep"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		before := snapshot(t, root)
		_, err := Init(root, testGen())
		if !errors.Is(err, ErrExists) {
			t.Fatalf("Init err = %v, want ErrExists", err)
		}
		if after := snapshot(t, root); !reflect.DeepEqual(before, after) {
			t.Errorf("Init changed the tree:\nbefore %v\nafter  %v", before, after)
		}
	})
	t.Run("file", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, Dir), []byte("plain"), 0o644); err != nil {
			t.Fatal(err)
		}
		before := snapshot(t, root)
		_, err := Init(root, testGen())
		if !errors.Is(err, ErrExists) {
			t.Fatalf("Init err = %v, want ErrExists", err)
		}
		if after := snapshot(t, root); !reflect.DeepEqual(before, after) {
			t.Errorf("Init changed the tree:\nbefore %v\nafter  %v", before, after)
		}
	})
}

func TestOpenErrors(t *testing.T) {
	valid := "{\n  \"format\": 1,\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n"
	cases := []struct {
		name  string
		data  string
		want  error
		field string // must appear in the message
	}{
		{"newer format", "{\n  \"format\": 2,\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n", ErrNewerFormat, "format"},
		{"newer format with extra key", "{\n  \"extra\": true,\n  \"format\": 2\n}\n", ErrNewerFormat, "format"},
		{"format zero", "{\n  \"format\": 0,\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n", ErrInvalid, "format"},
		{"format string", "{\n  \"format\": \"1\",\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n", ErrInvalid, "format"},
		{"missing format", "{\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n", ErrInvalid, "format"},
		{"extra key", "{\n  \"extra\": 1,\n  \"format\": 1,\n  \"id\": \"" + sID + "\",\n  \"remotes\": []\n}\n", ErrInvalid, "extra"},
		{"non-canonical", strings.TrimSuffix(valid, "\n"), ErrInvalid, ""},
		{"compact", `{"format":1,"id":"` + sID + `","remotes":[]}`, ErrInvalid, ""},
		{"not an object", "[]\n", ErrInvalid, ""},
		{"lower-case id", "{\n  \"format\": 1,\n  \"id\": \"" + strings.ToLower(sID) + "\",\n  \"remotes\": []\n}\n", ErrInvalid, "id"},
		{"id not a string", "{\n  \"format\": 1,\n  \"id\": 5,\n  \"remotes\": []\n}\n", ErrInvalid, "id"},
		{"missing id", "{\n  \"format\": 1,\n  \"remotes\": []\n}\n", ErrInvalid, "id"},
		{"remotes null", "{\n  \"format\": 1,\n  \"id\": \"" + sID + "\",\n  \"remotes\": null\n}\n", ErrInvalid, "remotes"},
		{"remotes not strings", "{\n  \"format\": 1,\n  \"id\": \"" + sID + "\",\n  \"remotes\": [\n    1\n  ]\n}\n", ErrInvalid, "remotes"},
		{"missing remotes", "{\n  \"format\": 1,\n  \"id\": \"" + sID + "\"\n}\n", ErrInvalid, "remotes"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.Mkdir(filepath.Join(root, Dir), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, Dir, RepoFile), []byte(c.data), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := Open(root)
			if !errors.Is(err, c.want) {
				t.Fatalf("Open err = %v, want %v", err, c.want)
			}
			if c.want == ErrNewerFormat && errors.Is(err, ErrInvalid) {
				t.Errorf("newer format also matches ErrInvalid: %v", err)
			}
			if !errors.Is(err, entry.ErrNewerFormat) && c.want == ErrNewerFormat {
				t.Errorf("err does not match entry.ErrNewerFormat: %v", err)
			}
			if c.field != "" && !strings.Contains(err.Error(), `"`+c.field+`"`) {
				t.Errorf("err %q does not name field %q", err, c.field)
			}
		})
	}
	t.Run("valid control", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, Dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, Dir, RepoFile), []byte(valid), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); err != nil {
			t.Fatal(err)
		}
	})
	t.Run(".genesis is a file", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, Dir), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); !errors.Is(err, ErrInvalid) || !strings.Contains(err.Error(), Dir) {
			t.Fatalf("err = %v, want ErrInvalid naming %s", err, Dir)
		}
	})
	t.Run("repo.json is a folder", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, Dir, RepoFile), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); !errors.Is(err, ErrInvalid) || !strings.Contains(err.Error(), RepoFile) {
			t.Fatalf("err = %v, want ErrInvalid naming %s", err, RepoFile)
		}
	})
}

func TestOpenNotFound(t *testing.T) {
	t.Run("no .genesis", func(t *testing.T) {
		if _, err := Open(t.TempDir()); !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want ErrNotFound", err)
		}
	})
	t.Run("no repo.json", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, Dir), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(root); !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestWrite(t *testing.T) {
	root, s := mustInit(t)
	e := testEntry(t)
	if err := s.Write(e); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, Dir, LogDir, sCS)
	p := filepath.Join(dir, sID+".json")
	info, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != 0o644 {
		t.Errorf("mode = %o, want 644", mode)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := entry.Encode(e)
	if !bytes.Equal(b, want) {
		t.Errorf("bytes =\n%s\nwant\n%s", b, want)
	}
	got, err := entry.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, e) {
		t.Errorf("Decode = %+v, want %+v", got, e)
	}
	assertNoTemp(t, dir)
	if names := listNames(t, dir); !slices.Equal(names, []string{sID + ".json"}) {
		t.Errorf("change set folder = %v", names)
	}

	// A second entry in the same, now existing, change set folder.
	e2 := e
	e2.ID = mustID(t, "01JB000000000000000000000B")
	if err := s.Write(e2); err != nil {
		t.Fatal(err)
	}
	assertNoTemp(t, dir)
	if names := listNames(t, dir); len(names) != 2 {
		t.Errorf("change set folder = %v", names)
	}
}

func TestWriteAppendOnlyEntryExists(t *testing.T) {
	root, s := mustInit(t)
	a := testEntry(t)
	if err := s.Write(a); err != nil {
		t.Fatal(err)
	}
	b := a
	b.Body = json.RawMessage("{\n  \"title\": \"different\"\n}")
	b.Type = entry.TypeNote
	err := s.Write(b)
	if !errors.Is(err, ErrEntryExists) {
		t.Fatalf("second Write err = %v, want ErrEntryExists", err)
	}
	p := filepath.Join(root, Dir, LogDir, sCS, sID+".json")
	if !strings.Contains(err.Error(), p) {
		t.Errorf("err %q does not name %s", err, p)
	}
	t.Logf("second Write: %v", err)
	got, rerr := os.ReadFile(p)
	if rerr != nil {
		t.Fatal(rerr)
	}
	want, _ := entry.Encode(a)
	if !bytes.Equal(got, want) {
		t.Errorf("entry file changed:\n%s\nwant\n%s", got, want)
	}
	assertNoTemp(t, filepath.Join(root, Dir, LogDir, sCS))
}

func TestWriteInvalidWritesNothing(t *testing.T) {
	root, s := mustInit(t)
	logDir := filepath.Join(root, Dir, LogDir)
	before := snapshot(t, root)
	e := testEntry(t)
	e.Actor = id.ID{}
	err := s.Write(e)
	if !errors.Is(err, entry.ErrInvalid) {
		t.Fatalf("err = %v, want entry.ErrInvalid", err)
	}
	if names := listNames(t, logDir); len(names) != 0 {
		t.Errorf("log/ = %v, want empty", names)
	}
	if after := snapshot(t, root); !reflect.DeepEqual(before, after) {
		t.Errorf("tree changed:\nbefore %v\nafter  %v", before, after)
	}
}

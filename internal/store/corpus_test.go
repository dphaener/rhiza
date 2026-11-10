package store

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/entry"
)

// corpusDir holds sample logs from every released format; see its README.
const corpusDir = "../../testdata/corpus"

func corpusCases(t *testing.T) []string {
	t.Helper()
	cases, err := filepath.Glob(filepath.Join(corpusDir, "format-*", "*"))
	if err != nil {
		t.Fatal(err)
	}
	cases = slices.DeleteFunc(cases, func(p string) bool {
		info, err := os.Stat(p)
		return err != nil || !info.IsDir()
	})
	if len(cases) == 0 {
		t.Fatalf("no corpus cases under %s", corpusDir)
	}
	return cases
}

// firstFormat names the earliest format that can spell typ: pin came with
// format 5, reopen with format 3, every other type with format 1.
func firstFormat(typ entry.Type) string {
	if typ == entry.TypePin {
		return "format-5"
	}
	if typ == entry.TypeReopen {
		return "format-3"
	}
	return "format-1"
}

func TestCorpus(t *testing.T) {
	// seen records, per format, the entry types its cases use.
	seen := map[string]map[entry.Type]bool{}
	for _, root := range corpusCases(t) {
		format := filepath.Base(filepath.Dir(root))
		t.Run(format+"/"+filepath.Base(root), func(t *testing.T) {
			s, err := Open(root)
			if err != nil {
				t.Fatal(err)
			}
			got := s.ReadLog()
			for _, p := range got.Problems {
				t.Errorf("problem at %s: %v", p.Path, p.Err)
			}
			if len(got.Entries) == 0 {
				t.Error("no entries")
			}
			if seen[format] == nil {
				seen[format] = map[entry.Type]bool{}
			}
			for _, e := range got.Entries {
				seen[format][e.Type] = true
			}

			files := 0
			err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if p != root && strings.HasPrefix(d.Name(), ".") && d.Name() != Dir {
					t.Errorf("dot file in corpus: %s", p)
				}
				if d.IsDir() || filepath.Ext(p) != ".json" {
					return nil
				}
				files++
				b, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				if err := canon.Check(b); err != nil {
					t.Errorf("%s: %v", p, err)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if files == 0 {
				t.Error("no .json files")
			}
		})
	}

	// Every type is used by some case of the earliest format that has it.
	missing := map[string][]entry.Type{}
	for _, typ := range entry.Types() {
		if f := firstFormat(typ); !seen[f][typ] {
			missing[f] = append(missing[f], typ)
		}
	}
	for f, types := range missing {
		t.Errorf("%s corpus never uses types %v", f, types)
	}
}

// copyTree copies the files and folders under src into dst.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCorpusCorrupt(t *testing.T) {
	root := t.TempDir()
	copyTree(t, filepath.Join(corpusDir, "format-1", "single-change-set"), root)
	s, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	clean := s.ReadLog()
	if len(clean.Problems) != 0 || len(clean.Entries) < 2 {
		t.Fatalf("copy does not read cleanly: %d entries, problems %v", len(clean.Entries), clean.Problems)
	}

	victim := clean.Entries[1]
	rel := entryPath(victim)
	f, err := os.OpenFile(filepath.Join(root, Dir, filepath.FromSlash(rel)), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(" "); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	got := s.ReadLog()
	if len(got.Problems) != 1 {
		t.Fatalf("Problems = %v, want exactly one", got.Problems)
	}
	p := got.Problems[0]
	if p.Path != rel {
		t.Errorf("Path = %q, want %q", p.Path, rel)
	}
	if !errors.Is(p.Err, entry.ErrInvalid) {
		t.Errorf("Err = %v, want entry.ErrInvalid", p.Err)
	}
	t.Logf("corrupt entry: %s: %v", p.Path, p.Err)

	want := slices.Delete(slices.Clone(clean.Entries), 1, 2)
	if !reflect.DeepEqual(got.Entries, want) {
		t.Errorf("Entries = %+v\nwant the other %d", got.Entries, len(want))
	}
}

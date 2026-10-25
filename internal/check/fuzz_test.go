package check

import (
	"io/fs"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dphaener/genesis/internal/store"
)

const corpusDir = "../../testdata/corpus/format-1"

// corpus loads each corpus case's .genesis folder into memory, by case name.
func corpus(t testing.TB) map[string]fstest.MapFS {
	t.Helper()
	des, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]fstest.MapFS{}
	for _, de := range des {
		if !de.IsDir() {
			continue
		}
		root := os.DirFS(path.Join(corpusDir, de.Name(), ".genesis"))
		m := fstest.MapFS{}
		err := fs.WalkDir(root, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			b, err := fs.ReadFile(root, p)
			m[p] = &fstest.MapFile{Data: b}
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		out[de.Name()] = m
	}
	if len(out) == 0 {
		t.Fatal("no corpus cases")
	}
	return out
}

// entryFiles returns the log entry paths in m, sorted.
func entryFiles(m fstest.MapFS) []string {
	var out []string
	for p := range m {
		if strings.HasPrefix(p, store.LogDir+"/") && strings.HasSuffix(p, ".json") {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return out
}

// TestCorpus checks every corpus case with a current render: nothing fails.
func TestCorpus(t *testing.T) {
	for name, m := range corpus(t) {
		t.Run(name, func(t *testing.T) {
			l := store.ReadLog(m)
			rep := Run(Input{Log: l, Rendered: rendered(l), Base: "base"})
			if !rep.OK {
				t.Errorf("not OK: %+v", rep.Findings)
			}
		})
	}
}

// FuzzCheck puts arbitrary bytes at one entry path of a valid repository:
// Run never panics, and fails whenever the log does not read cleanly.
func FuzzCheck(f *testing.F) {
	cases := corpus(f)
	for _, name := range slices.Sorted(maps.Keys(cases)) {
		for _, p := range entryFiles(cases[name]) {
			f.Add(cases[name][p].Data)
		}
	}
	base := cases["single-change-set"]
	target := entryFiles(base)[0]
	f.Fuzz(func(t *testing.T, data []byte) {
		m := fstest.MapFS{}
		for p, file := range base {
			m[p] = file
		}
		m[target] = &fstest.MapFile{Data: data}
		l := store.ReadLog(m)
		rep := Run(Input{Log: l, Rendered: rendered(l)})
		if len(l.Problems) > 0 && rep.OK {
			t.Fatalf("log problems %v but OK", l.Problems)
		}
		stale := Run(Input{Log: l, Rendered: base})
		if len(l.Problems) > 0 && stale.OK {
			t.Fatalf("log problems %v but OK", l.Problems)
		}
	})
}

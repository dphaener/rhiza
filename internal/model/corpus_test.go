package model

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/store"
)

// corpusDir holds sample logs from every released format; see its README.
const corpusDir = "../../testdata/corpus"

// corpusHandles lists handles of live nodes each case's as-is must hold, so
// the test cannot pass on an empty model.
var corpusHandles = map[string][]string{
	"format-1/single-change-set": {"cafe-menu"},
	"format-1/revisions":         {"report-export"},
}

// TestCorpus projects every corpus case: as-is, to-be and the what-if of
// every in-flight change set all have no problems.
func TestCorpus(t *testing.T) {
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
	seen := map[string]bool{}
	for _, root := range cases {
		name := filepath.Base(filepath.Dir(root)) + "/" + filepath.Base(root)
		seen[name] = true
		t.Run(name, func(t *testing.T) {
			l := store.ReadLog(os.DirFS(filepath.Join(root, store.Dir)))
			for _, p := range l.Problems {
				t.Errorf("store problem at %s: %v", p.Path, p.Err)
			}
			if len(l.Entries) == 0 {
				t.Fatal("no entries")
			}
			sets, fps := changeset.FoldLog(l)
			for _, p := range fps {
				t.Errorf("fold problem: %v", p)
			}

			asIs, ps := AsIs(sets)
			for _, p := range ps {
				t.Errorf("AsIs problem: %v", p)
			}
			_, ps = ToBe(sets)
			for _, p := range ps {
				t.Errorf("ToBe problem: %v", p)
			}
			for _, s := range sets {
				if !inFlight(s.State) {
					continue
				}
				_, ps = WhatIf(sets, s.ID)
				for _, p := range ps {
					t.Errorf("WhatIf %s problem: %v", s.ID, p)
				}
			}

			var handles []string
			for _, n := range asIs.Nodes() {
				if !n.Retired {
					handles = append(handles, n.Handle)
				}
			}
			for _, h := range corpusHandles[name] {
				if !slices.Contains(handles, h) {
					t.Errorf("as-is has no live node %q; has %q", h, handles)
				}
			}
		})
	}
	for name := range corpusHandles {
		if !seen[name] {
			t.Errorf("corpus case %s is missing", name)
		}
	}
}

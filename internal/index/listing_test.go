package index

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/store"
)

// TestListLogSame checks listLog lists exactly what store.ListLog does,
// layout problems included.
func TestListLogSame(t *testing.T) {
	f := newFixture(t)
	for n := 1; n <= 20; n++ {
		f.add(csID(n), changeset.Open{Kind: changeset.KindRecord, Intent: "x"})
	}
	writeFile(t, f.gdir, "log/stray.json", []byte("{}"))
	writeFile(t, f.gdir, "log/not-an-id/x.json", []byte("{}"))
	writeFile(t, f.gdir, "log/.hidden/x.json", []byte("{}"))
	writeFile(t, f.gdir, "log/"+csID(1).String()+"/notes.txt", []byte("x"))
	writeFile(t, f.gdir, "log/"+csID(2).String()+"/nested/x.json", []byte("{}"))
	writeFile(t, f.gdir, "log/"+csID(3).String()+"/bad-stem.json", []byte("{}"))
	if err := os.Mkdir(filepath.Join(f.gdir, "log", csID(99).String()), 0o755); err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(f.gdir)
	want, got := store.ListLog(fsys), listLog(fsys)
	if !slices.Equal(got.Paths, want.Paths) {
		t.Fatalf("paths:\n got %v\nwant %v", got.Paths, want.Paths)
	}
	if len(want.Problems) < 5 || fingerprint(got) != fingerprint(want) {
		t.Fatalf("problems:\n got %v\nwant %v", got.Problems, want.Problems)
	}

	// No log folder at all.
	empty := os.DirFS(t.TempDir())
	if fingerprint(listLog(empty)) != fingerprint(store.ListLog(empty)) {
		t.Fatal("missing log folder listed differently")
	}
}

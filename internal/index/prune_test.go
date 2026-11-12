package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPruneOrphansRemovesOnlyGoneWorktrees(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "cache")
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", t.TempDir())
	live, gone := t.TempDir(), t.TempDir()
	for _, r := range []string{live, gone} {
		if err := os.MkdirAll(filepath.Join(r, ".rhiza", "log"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	paths := map[string]string{}
	for _, r := range []string{live, gone} {
		p, err := DefaultPath(r)
		if err != nil {
			t.Fatal(err)
		}
		x, err := Open(r, p)
		if err != nil {
			t.Fatal(err)
		}
		x.Close()
		paths[r] = p
	}
	// A database with no recorded root, as older versions made.
	legacy := filepath.Join(filepath.Dir(paths[live]), "legacy.db")
	if err := os.WriteFile(legacy, []byte("not a db"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(gone); err != nil {
		t.Fatal(err)
	}
	dir, err := IndexDir()
	if err != nil {
		t.Fatal(err)
	}
	got := PruneOrphans(dir)
	if got.Files != 1 || got.Bytes <= 0 {
		t.Fatalf("pruned %+v, want 1 file with bytes", got)
	}
	if _, err := os.Stat(paths[gone]); !os.IsNotExist(err) {
		t.Errorf("orphan index still there")
	}
	for _, p := range []string{paths[live], legacy} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s removed: %v", p, err)
		}
	}
}

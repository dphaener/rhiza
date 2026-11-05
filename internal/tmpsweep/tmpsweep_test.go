package tmpsweep

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStaleRemovesOnlyOldPrefixed(t *testing.T) {
	root := t.TempDir()
	old := time.Now().Add(-2 * time.Hour)
	mk := func(name string, mtime time.Time) string {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Join(p, "sub"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(p, "sub", "f"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(p, mtime, mtime); err != nil {
			t.Fatal(err)
		}
		return p
	}
	staleDir := mk("rhiza-test-old", old)
	fresh := mk("rhiza-test-fresh", time.Now())
	other := mk("other-old", old)

	Stale(root, "rhiza-test-")

	if _, err := os.Stat(staleDir); !os.IsNotExist(err) {
		t.Errorf("stale prefixed folder still there: %v", err)
	}
	for _, p := range []string{fresh, other} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s removed: %v", p, err)
		}
	}
}

func TestStaleMissingRoot(t *testing.T) {
	Stale(filepath.Join(t.TempDir(), "none"), "x")
}

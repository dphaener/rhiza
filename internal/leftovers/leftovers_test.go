package leftovers

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScratchRemovesOnlyOldOrphans(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	mk := func(name string, age time.Duration, size int) string {
		d := filepath.Join(root, name)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "f"), make([]byte, size), 0o644); err != nil {
			t.Fatal(err)
		}
		at := now.Add(-age)
		if err := os.Chtimes(d, at, at); err != nil {
			t.Fatal(err)
		}
		return d
	}
	oldOrphan := mk("s-old-ended", 48*time.Hour, 100)
	oldLive := mk("s-old-live", 48*time.Hour, 10)
	fresh := mk("s-fresh-ended", time.Hour, 10)
	noRun := mk("s-old-norun", 72*time.Hour, 50)

	sessions := []Session{
		{ID: "s-old-ended", Run: "r-ended"},
		{ID: "s-old-live", Run: "r-live"},
		{ID: "s-fresh-ended", Run: "r-ended"},
		{ID: "s-old-norun"},
	}
	live := func(r string) bool { return r == "r-live" }
	globs := func(s Session) []string { return []string{filepath.Join(root, s.ID)} }
	f := Scratch(sessions, live, globs, now)
	if f.Folders != 2 || f.Bytes != 150 {
		t.Fatalf("freed %+v, want 2 folders and 150 bytes", f)
	}
	for _, p := range []string{oldOrphan, noRun} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s still there", p)
		}
	}
	for _, p := range []string{oldLive, fresh} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s removed: %v", p, err)
		}
	}
}

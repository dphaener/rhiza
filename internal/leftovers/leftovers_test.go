package leftovers

import (
	"os"
	"path/filepath"
	"strconv"
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
	parents := func(string) []string { return []string{root} }
	f := Scratch(sessions, live, parents, nil, now)
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

func TestScratchReadsEachParentOnce(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	old := now.Add(-48 * time.Hour)
	var sessions []Session
	for p := 0; p < 3; p++ {
		parent := filepath.Join(root, "claude-1", "p"+string(rune('a'+p)))
		for i := 0; i < 500; i++ {
			id := "s-" + string(rune('a'+p)) + "-" + itoa(i)
			d := filepath.Join(parent, id)
			if err := os.MkdirAll(d, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(d, "f"), []byte("12345"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Chtimes(d, old, old); err != nil {
				t.Fatal(err)
			}
			run := "r-ended"
			if i%10 == 0 {
				run = "r-live"
			}
			sessions = append(sessions, Session{ID: id, Agent: "a", Run: run})
		}
	}
	// 1500 more recorded sessions with no folder on disk.
	for i := 0; i < 1500; i++ {
		sessions = append(sessions, Session{ID: "ghost-" + itoa(i), Agent: "a", Run: "r-ended"})
	}
	reads := map[string]int{}
	read := func(name string) ([]os.DirEntry, error) {
		reads[name]++
		return os.ReadDir(name)
	}
	live := func(r string) bool { return r == "r-live" }
	f := Scratch(sessions, live, func(string) []string { return []string{filepath.Join(root, "claude-*", "*")} }, read, now)
	if f.Folders != 1350 || f.Bytes != 1350*5 {
		t.Fatalf("freed %+v, want 1350 folders and %d bytes", f, 1350*5)
	}
	if len(reads) != 3 {
		t.Fatalf("read %d directories, want 3: %v", len(reads), reads)
	}
	for d, n := range reads {
		if n != 1 {
			t.Errorf("%s read %d times", d, n)
		}
	}
	left, _ := os.ReadDir(filepath.Join(root, "claude-1", "pa"))
	if len(left) != 50 {
		t.Errorf("%d folders left in pa, want the 50 of the live run", len(left))
	}
}

func TestRemoveSessionSafeIDOnly(t *testing.T) {
	root := t.TempDir()
	keep := filepath.Join(root, "keep")
	if err := os.MkdirAll(filepath.Join(root, "p", "sid"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(keep, 0o755); err != nil {
		t.Fatal(err)
	}
	RemoveSession([]string{filepath.Join(root, "p")}, "../keep")
	RemoveSession([]string{filepath.Join(root, "p")}, "*")
	if _, err := os.Stat(keep); err != nil {
		t.Fatal(err)
	}
	RemoveSession([]string{filepath.Join(root, "*")}, "sid")
	if _, err := os.Stat(filepath.Join(root, "p", "sid")); !os.IsNotExist(err) {
		t.Fatal("sid still there")
	}
}

func itoa(i int) string { return strconv.Itoa(i) }

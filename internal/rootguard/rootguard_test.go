package rootguard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	t0 = time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	t1 = t0.Add(time.Minute)
)

// newRoot makes a temp root with a state folder holding a limits file with
// one old line, and returns the root and the state folder.
func newRoot(t *testing.T) (root, state string) {
	t.Helper()
	root = t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "rhiza.toml"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	state = filepath.Join(root, StateDir)
	if err := os.Mkdir(state, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(state, limitsFile), limit(t0.Add(-time.Hour)))
	return root, state
}

func limit(at time.Time) string {
	return fmt.Sprintf(`{"hit_at":%q,"reset_at":%q,"agent":"x"}`+"\n", at.Format(time.RFC3339), at.Add(time.Hour).Format(time.RFC3339))
}

func write(t *testing.T, path, s string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func appendTo(t *testing.T, path, s string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		t.Fatal(err)
	}
}

func TestUntouched(t *testing.T) {
	root, _ := newRoot(t)
	snap := Take(root)
	if got := snap.Check(t0, t1); len(got) != 0 {
		t.Fatalf("untouched folder reported %v", got)
	}
}

func TestNoRoot(t *testing.T) {
	snap := Take(t.TempDir())
	if got := snap.Check(t0, t1); len(got) != 0 {
		t.Fatalf("no root reported %v", got)
	}
}

func TestFileAddedAndChanged(t *testing.T) {
	root, state := newRoot(t)
	write(t, filepath.Join(state, "kept.json"), "a")
	snap := Take(root)

	write(t, filepath.Join(state, "new.json"), "x")
	write(t, filepath.Join(state, "kept.json"), "longer")
	got := strings.Join(snap.Check(t0, t1), "\n")
	for _, want := range []string{"new.json: file added", "kept.json: file changed"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q in %q", want, got)
		}
	}
}

func TestDaemonOwnedFiles(t *testing.T) {
	root, state := newRoot(t)
	write(t, filepath.Join(state, "daemon.lock"), "1")
	snap := Take(root)
	for n := range owned {
		write(t, filepath.Join(state, n), "changed by the daemon\n")
	}
	// limits file rewritten with no new line: nothing to attribute.
	write(t, filepath.Join(state, limitsFile), limit(t0.Add(-time.Hour)))
	if got := snap.Check(t0, t1); len(got) != 0 {
		t.Fatalf("daemon-owned files reported %v", got)
	}
}

func TestLimitLineOutsideSpan(t *testing.T) {
	root, state := newRoot(t)
	snap := Take(root)
	appendTo(t, filepath.Join(state, limitsFile), limit(time.Date(2026, 11, 4, 10, 45, 0, 0, time.UTC)))
	got := snap.Check(t0, t1)
	if len(got) != 1 || !strings.Contains(got[0], "usage.limits.jsonl: line 2 appended") {
		t.Fatalf("got %v", got)
	}
}

func TestLimitLineInsideSpan(t *testing.T) {
	root, state := newRoot(t)
	snap := Take(root)
	appendTo(t, filepath.Join(state, limitsFile), limit(t0.Add(30*time.Second)))
	if got := snap.Check(t0, t1); len(got) != 0 {
		t.Fatalf("driver hit reported %v", got)
	}
}

func TestMainFailsOnWrite(t *testing.T) {
	root, state := newRoot(t)
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	if code := Main(func() int { return 0 }); code != 0 {
		t.Fatalf("clean run exited %d", code)
	}
	code := Main(func() int {
		write(t, filepath.Join(state, "stray"), "x")
		return 0
	})
	if code == 0 {
		t.Fatal("write into the state folder did not fail the run")
	}
}

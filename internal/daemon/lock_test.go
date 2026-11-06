package daemon

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

var t0 = time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)

func TestSecondLockIsHeldAndNamesTheFirst(t *testing.T) {
	common := t.TempDir()
	l, err := TryLock(common, t0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = TryLock(common, t0)
	var h *Held
	if !errors.As(err, &h) {
		t.Fatalf("second TryLock = %v, want *Held", err)
	}
	if h.Holder.Pid != os.Getpid() {
		t.Errorf("holder pid = %d, want %d", h.Holder.Pid, os.Getpid())
	}
	if _, ok := Running(common); !ok {
		t.Error("Running = false while locked")
	}
	l.Unlock()
	if _, ok := Running(common); ok {
		t.Error("Running = true after Unlock")
	}
	l2, err := TryLock(common, t0)
	if err != nil {
		t.Fatalf("TryLock after Unlock: %v", err)
	}
	l2.Unlock()
}

func TestStaleLockFileIsTakenOver(t *testing.T) {
	common := t.TempDir()
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		t.Fatal(err)
	}
	stale := `{"host":"elsewhere","pid":999999,"started":"2026-01-01T00:00:00Z"}`
	if err := os.WriteFile(filepath.Join(Dir(common), File), []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	l, err := TryLock(common, t0)
	if err != nil {
		t.Fatalf("TryLock over a stale file: %v", err)
	}
	defer l.Unlock()
	if h := ReadHolder(common); h.Pid != os.Getpid() || !h.Started.Equal(t0) {
		t.Errorf("holder = %+v, want this process started %v", h, t0)
	}
}

func TestRootLockIsOnePerRoot(t *testing.T) {
	root := t.TempDir()
	l, err := TryLockRoot(root, t0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = TryLockRoot(root, t0)
	var h *Held
	if !errors.As(err, &h) || h.Holder.Pid != os.Getpid() {
		t.Fatalf("second TryLockRoot = %v, want *Held naming this process", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".rhiza-root", File)); err != nil {
		t.Errorf("lock file not in the root's state folder: %v", err)
	}
	l.Unlock()
	l2, err := TryLockRoot(root, t0)
	if err != nil {
		t.Fatalf("after Unlock: %v", err)
	}
	l2.Unlock()
}

// The lock is the root's whichever folder under it a daemon starts from:
// one taken from the root holds off one resolved from a project dir, and
// the other way round, both naming one file in <root>/.rhiza-root.
func TestRootLockFromProjectAndRootAlike(t *testing.T) {
	root := filepath.Join(t.TempDir(), "top")
	if _, err := rootcfg.Init(root); err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(root, "app", "sub")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	resolve := func(dir string) string {
		t.Helper()
		r, ok := rootcfg.Find(dir)
		if !ok {
			t.Fatalf("rootcfg.Find(%s) found no root", dir)
		}
		return r
	}
	for _, order := range [][2]string{{root, project}, {project, root}} {
		first, second := resolve(order[0]), resolve(order[1])
		l, err := TryLockRoot(first, t0)
		if err != nil {
			t.Fatalf("lock from %s: %v", order[0], err)
		}
		_, err = TryLockRoot(second, t0)
		var h *Held
		if !errors.As(err, &h) || h.Holder.Pid != os.Getpid() {
			t.Fatalf("second lock from %s = %v, want *Held naming this process", order[1], err)
		}
		if _, err := os.Stat(filepath.Join(root, ".rhiza-root", File)); err != nil {
			t.Errorf("lock file not in the root's state folder: %v", err)
		}
		l.Unlock()
	}
}

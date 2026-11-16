package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/store"
)

// backdateAge is how far backdateLog sets a folder's modification time
// back: past the index's two-second racy window, as a 2.2 s sleep did.
const backdateAge = 2200 * time.Millisecond

// racyWindow mirrors the index's unexported racyWindow
// (internal/index/listing.go): a folder whose modification time is older
// than this is no longer racy.
const racyWindow = 2 * time.Second

// backdateLog sets the modification time of every folder in root's log
// still inside the index's racy window back to backdateAge ago, as if the
// test had slept past the index's racy window (pattern:
// internal/index/freshness_test.go). Older folders keep their times, so
// a stat fingerprint the index stored over them stays the same, as it
// would over a sleep. The racy check reads only modification times; the
// change time Chtimes moves is in the fingerprint alone, and only for a
// folder that changed anyway.
func backdateLog(t *testing.T, root string) {
	t.Helper()
	logDir := filepath.Join(root, store.Dir, store.LogDir)
	des, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	old := now.Add(-backdateAge)
	cutoff := now.Add(-racyWindow)
	for _, de := range des {
		if !de.IsDir() {
			continue
		}
		dir := filepath.Join(logDir, de.Name())
		fi, err := os.Stat(dir)
		if err != nil {
			t.Fatal(err)
		}
		if fi.ModTime().Before(cutoff) {
			continue
		}
		if err := os.Chtimes(dir, old, old); err != nil {
			t.Fatal(err)
		}
	}
}

// TestBackdateLogPassesRacyWindow checks index.LogStamp is racy for a log
// folder just written and not racy once backdateLog has run, with no
// sleep.
func TestBackdateLogPassesRacyWindow(t *testing.T) {
	root := t.TempDir()
	ctx := testCtx(t, root)
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	newChange(t, ctx, `{"intent":"Just written"}`)
	fp, racy := index.LogStamp(root)
	if fp == "" {
		t.Fatal("no stamp for the log")
	}
	if !racy {
		t.Fatal("a folder just written is not racy")
	}
	backdateLog(t, root)
	after, racy := index.LogStamp(root)
	if after == "" {
		t.Fatal("no stamp for the backdated log")
	}
	if racy {
		t.Fatal("still racy after backdateLog")
	}
	if after == fp {
		t.Error("backdating did not change the stamp")
	}
}

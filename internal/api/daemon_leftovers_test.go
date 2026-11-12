package api

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// TestDaemonSweepLeftoversPrunesOnlyOrphans builds index databases and
// session scratch folders, live and orphaned, injects the clock, and runs
// the daemon's sweep: only the orphans go, the stderr line carries the
// freed bytes, and a second call within the hour does nothing.
func TestDaemonSweepLeftoversPrunesOnlyOrphans(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	common := t.TempDir()
	var stderr bytes.Buffer
	ctx := &Ctx{
		Home:   home,
		Root:   t.TempDir(),
		Now:    func() time.Time { return now },
		Git:    func(args ...string) (string, error) { return common + "\n", nil },
		Stderr: &stderr,
	}

	// Index databases: one for a worktree that stays, one for a gone one.
	trees := t.TempDir()
	dbs := map[string]string{}
	for _, n := range []string{"live", "gone"} {
		root := filepath.Join(trees, n)
		if err := os.MkdirAll(filepath.Join(root, ".rhiza", "log"), 0o755); err != nil {
			t.Fatal(err)
		}
		p, err := index.DefaultPath(root)
		if err != nil {
			t.Fatal(err)
		}
		x, err := index.Open(root, p)
		if err != nil {
			t.Fatal(err)
		}
		x.Close()
		dbs[n] = p
	}
	if err := os.RemoveAll(filepath.Join(trees, "gone")); err != nil {
		t.Fatal(err)
	}
	// The bytes an index sweep frees include the WAL and shared-memory
	// files its own read makes, so learn them by pruning a copy first.
	probe := t.TempDir()
	for _, suf := range []string{"", "-wal", "-shm"} {
		b, err := os.ReadFile(dbs["gone"] + suf)
		if err != nil {
			continue
		}
		if err := os.WriteFile(filepath.Join(probe, "copy.db"+suf), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	pr := index.PruneOrphans(probe)
	if pr.Files != 1 || pr.Bytes <= 0 {
		t.Fatalf("probe pruned %+v, want 1 file with bytes", pr)
	}
	wantBytes := pr.Bytes

	// Runs: one live, one finished.
	runs := run.NewStore(run.Dir(common))
	for _, id := range []string{"run-live", "run-ended"} {
		if err := runs.Create(id, run.Entry{Type: run.TypeStarted, At: now.Add(-72 * time.Hour), Started: &run.Started{Workflow: "build"}}); err != nil {
			t.Fatal(err)
		}
	}
	if err := runs.Append("run-ended", run.Entry{Type: run.TypeFinished, At: now.Add(-48 * time.Hour), End: run.Succeeded}); err != nil {
		t.Fatal(err)
	}

	// Sessions and their scratch folders.
	sessions := session.NewStore(session.Dir(common))
	old, fresh := now.Add(-48*time.Hour), now.Add(-time.Hour)
	scratch := map[string]string{}
	mk := func(sid, rn string, mtime time.Time) {
		t.Helper()
		if err := sessions.Write(session.Record{ID: sid, Agent: interviewAgent, Run: rn}); err != nil {
			t.Fatal(err)
		}
		d := filepath.Join(home, "scratch", sid)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		f := filepath.Join(d, "notes.txt")
		if err := os.WriteFile(f, []byte("scratch of "+sid), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, p := range []string{f, d} {
			if err := os.Chtimes(p, mtime, mtime); err != nil {
				t.Fatal(err)
			}
		}
		scratch[sid] = d
	}
	mk("s-live-old", "run-live", old)
	mk("s-ended-old", "run-ended", old)
	mk("s-ended-fresh", "run-ended", fresh)
	mk("s-norun-old", "", old)
	mk("s-unknown-old", "run-unknown", old)
	gone := []string{"s-ended-old", "s-norun-old", "s-unknown-old"}
	for _, sid := range gone {
		wantBytes += int64(len("scratch of " + sid))
	}

	var sched daemon.SweepSchedule
	daemonSweepLeftovers(ctx, &sched)

	for _, sid := range gone {
		if _, err := os.Stat(scratch[sid]); !os.IsNotExist(err) {
			t.Errorf("orphan scratch %s stays: %v", sid, err)
		}
	}
	for _, sid := range []string{"s-live-old", "s-ended-fresh"} {
		if _, err := os.Stat(scratch[sid]); err != nil {
			t.Errorf("scratch %s went: %v", sid, err)
		}
	}
	if _, err := os.Stat(dbs["gone"]); !os.IsNotExist(err) {
		t.Errorf("orphan index stays: %v", err)
	}
	if _, err := os.Stat(dbs["live"]); err != nil {
		t.Errorf("live worktree's index went: %v", err)
	}
	want := "rhiza: daemon: leftovers sweep freed " + daemon.FormatSize(uint64(wantBytes)) +
		": 1 index databases, 3 session scratch folders\n"
	if got := stderr.String(); got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}

	// Within the hour the sweep does nothing, even with a new orphan.
	stderr.Reset()
	mk("s-later-old", "run-ended", old)
	now = now.Add(59 * time.Minute)
	daemonSweepLeftovers(ctx, &sched)
	if _, err := os.Stat(scratch["s-later-old"]); err != nil {
		t.Fatalf("a sweep within the hour removed %s: %v", "s-later-old", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("a sweep within the hour logged %q", stderr.String())
	}

	// An hour after the first, it runs again.
	now = now.Add(time.Minute)
	daemonSweepLeftovers(ctx, &sched)
	if _, err := os.Stat(scratch["s-later-old"]); !os.IsNotExist(err) {
		t.Fatalf("the hourly sweep left %s: %v", "s-later-old", err)
	}
	if got := stderr.String(); !strings.Contains(got, "0 index databases, 1 session scratch folders") {
		t.Fatalf("second sweep logged %q", got)
	}
	if _, err := os.Stat(scratch["s-live-old"]); err != nil {
		t.Fatalf("a live run's scratch went on the second sweep: %v", err)
	}
}

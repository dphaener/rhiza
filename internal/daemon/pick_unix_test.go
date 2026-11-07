//go:build unix

package daemon

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A run another holder has locked through its store, as the starter does
// between journaling a run and recording its process, is not resumed; once
// the lock is let go the same run is.
func TestResumableLeavesAStoreLockedRun(t *testing.T) {
	rs := run.NewStore(t.TempDir())
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	runs := []run.Journal{{ID: "R1", Entries: []run.Entry{
		{Type: run.TypeStarted, At: now.Add(-time.Hour), Started: &run.Started{Change: "C1"}},
	}}}
	dead := func(int) bool { return false }
	unlock, ok, err := rs.TryLock("R1")
	if err != nil || !ok {
		t.Fatalf("TryLock = %v, %v", ok, err)
	}
	if got := Resumable(runs, "h", dead, rs.Live, func(string) string { return "accepted" }, now); len(got) != 0 {
		t.Fatalf("Resumable with the lock held = %v, want none", got)
	}
	unlock()
	if got := Resumable(runs, "h", dead, rs.Live, func(string) string { return "accepted" }, now); len(got) != 1 || got[0] != "R1" {
		t.Fatalf("Resumable with the lock free = %v, want [R1]", got)
	}
}

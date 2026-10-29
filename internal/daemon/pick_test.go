package daemon

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

func TestResumable(t *testing.T) {
	dead := func(int) bool { return false }
	proc := run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}}
	journal := func(id string, es ...run.Entry) run.Journal {
		return run.Journal{ID: id, Entries: append([]run.Entry{{Type: run.TypeStarted, Started: &run.Started{Change: id}}}, es...)}
	}
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	since := func(j run.Journal, d time.Duration) run.Journal {
		j.Entries[0].At = now.Add(-d)
		return j
	}
	runs := []run.Journal{
		since(journal("interrupted", proc), time.Hour),
		since(journal("again", proc, run.Entry{Type: run.TypeDaemonResumed}), time.Hour),
		since(journal("failed", proc, run.Entry{Type: run.TypeFinished, End: run.Failed}), time.Hour),
		// Journaled, its process not yet recorded: starting.
		since(journal("starting"), Unrecorded-time.Second),
		since(journal("unrecorded"), Unrecorded),
		// Another process holds its lock: running whatever the journal says.
		since(journal("locked"), time.Hour),
	}
	locked := func(id string) bool { return id == "locked" }
	accepted := func(string) bool { return true }
	got := Resumable(runs, "h", dead, locked, accepted, now)
	if len(got) != 2 || got[0] != "interrupted" || got[1] != "unrecorded" {
		t.Fatalf("Resumable = %v, want [interrupted unrecorded]", got)
	}
}

func TestResumableOnlyAcceptedChangeSets(t *testing.T) {
	dead := func(int) bool { return false }
	proc := run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}}
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	states := map[string]bool{"accepted": true}
	var runs []run.Journal
	for _, cs := range []string{"accepted", "landed", "rejected", "withdrawn", "proposed", "draft"} {
		runs = append(runs, run.Journal{ID: "run-" + cs, Entries: []run.Entry{
			{Type: run.TypeStarted, At: now.Add(-time.Hour), Started: &run.Started{Change: cs}}, proc}})
	}
	got := Resumable(runs, "h", dead, func(string) bool { return false }, func(cs string) bool { return states[cs] }, now)
	if len(got) != 1 || got[0] != "run-accepted" {
		t.Fatalf("Resumable = %v, want [run-accepted]", got)
	}
}

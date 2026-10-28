package api

import (
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
)

func pickJournal(cs string, entries ...run.Entry) run.Journal {
	at := time.Unix(1700000000, 0)
	j := run.Journal{ID: "R" + cs, Entries: []run.Entry{{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs}}}}
	return run.Journal{ID: j.ID, Entries: append(j.Entries, entries...)}
}

// The daemon picks the accepted, ready change sets that serve a milestone
// and have a leaf, none with a live run or an unattended failed run, in
// roadmap order, under the running runs' cap.
func TestDaemonPicks(t *testing.T) {
	f := newPlanFix()
	m := f.intent("m", id.ID{}, "a")
	f.milestone(m)
	plain := f.intent("plain", id.ID{}, "b")
	mk := func(name string, serves id.ID) id.ID {
		cs := f.cs(name, changeset.StateAccepted)
		f.serve(cs, serves)
		return cs
	}
	first, second := mk("first", m), mk("second", m)
	leafless, backlog := mk("leafless", m), mk("backlog", plain)
	running, interrupted, failed := mk("running", m), mk("interrupted", m), mk("failed", m)
	draft := f.cs("draft", changeset.StateDraft)
	f.serve(draft, m)

	hasLeaves := func(cs id.ID) bool { return cs != leafless }
	cands := daemonCandidates(f.core(), hasLeaves)
	if len(cands) == 0 || slices.Contains(cands, leafless.String()) || slices.Contains(cands, backlog.String()) || slices.Contains(cands, draft.String()) {
		t.Fatalf("candidates = %v", cands)
	}

	host := "h"
	up := func(pid int) bool { return pid == 1 }
	at := time.Unix(1700000000, 0)
	proc := func(pid int) run.Entry {
		return run.Entry{Type: run.TypeProcess, At: at, Process: &run.Process{Host: host, Pid: pid}}
	}
	fail := func(detail string) run.Entry {
		return run.Entry{Type: run.TypeFinished, At: at, End: run.Failed, Detail: detail}
	}
	runs := func(failure string) []run.Journal {
		return []run.Journal{
			pickJournal(running.String(), proc(1)),
			pickJournal(interrupted.String(), proc(2)),
			pickJournal(failed.String(), proc(3), fail(failure)),
		}
	}
	pick := func(max int, failure string) []string {
		return daemon.Pick(cands, runs(failure), nil, host, up, max)
	}
	want := func(max int, failure string, w ...id.ID) {
		t.Helper()
		var ws []string
		for _, i := range w {
			ws = append(ws, i.String())
		}
		if got := pick(max, failure); !slices.Equal(got, ws) {
			t.Errorf("max %d, failure %q: picked %v, want %v", max, failure, got, ws)
		}
	}
	want(2, "boom", first)
	want(1, "boom")
	want(4, "boom", first, second)
	want(2, run.AbandonedPrefix+"me", first)
	want(4, run.AbandonedPrefix+"me", first, second, failed)

	// A change set the daemon just started counts as running until its
	// journal shows.
	if got := daemon.Pick(cands, runs("boom"), []string{first.String()}, host, up, 3); !slices.Equal(got, []string{second.String()}) {
		t.Errorf("with %s pending picked %v", first, got)
	}
}

// A start whose run journals leaves pending. One that never journals is
// held past the grace, taking no slot, until a new journal of its change
// set shows.
func TestDaemonHeld(t *testing.T) {
	at := time.Unix(1700000000, 0)
	old := pickJournal("B")
	pending := map[string]daemonPending{"A": {at: at}, "B": {at: at, runs: 1}}
	unjournaled := map[string]daemonPending{}

	if got := daemonHeld(pending, unjournaled, []run.Journal{old}, at.Add(time.Second)); !slices.Equal(got, []string{"A", "B"}) {
		t.Fatalf("in the grace held %v", got)
	}
	if got := daemonHeld(pending, unjournaled, []run.Journal{old, pickJournal("A")}, at.Add(time.Second)); !slices.Equal(got, []string{"B"}) {
		t.Fatalf("after A journaled held %v", got)
	}
	if got := daemonHeld(pending, unjournaled, []run.Journal{old}, at.Add(daemonStartGrace)); len(got) != 0 || len(pending) != 0 {
		t.Fatalf("past the grace held %v, pending %v", got, pending)
	}
	if _, ok := unjournaled["B"]; !ok || len(unjournaled) != 1 {
		t.Fatalf("unjournaled = %v", unjournaled)
	}
	daemonHeld(pending, unjournaled, []run.Journal{old}, at.Add(time.Hour))
	if _, ok := unjournaled["B"]; !ok {
		t.Fatalf("B left unjournaled with no new journal")
	}
	daemonHeld(pending, unjournaled, []run.Journal{old, pickJournal("B")}, at.Add(time.Hour))
	if len(unjournaled) != 0 {
		t.Fatalf("B still unjournaled after a new journal: %v", unjournaled)
	}
}

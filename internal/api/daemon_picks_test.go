package api

import (
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/policy"
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
	notReady := mk("not-ready", m)
	f.dep(notReady, first)

	hasLeaves := func(cs id.ID) bool { return cs != leafless }
	cands, skips := daemonCandidates(f.core(), hasLeaves, daemonSchedule{Backlog: policy.BacklogSkip})
	reasons := map[string]string{}
	for _, sk := range skips {
		reasons[sk.Change] = sk.Reason
	}
	for cs, want := range map[id.ID]string{backlog: daemon.SkipNoMilestone, leafless: daemon.SkipNoLeaves, notReady: daemon.SkipNotReady} {
		if reasons[cs.String()] != want {
			t.Errorf("skip of %s = %q, want %q (skips %v)", cs, reasons[cs.String()], want, skips)
		}
	}
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
		return daemon.Pick(cands, runs(failure), nil, host, up, max, daemon.Ship{})
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
	if got := daemon.Pick(cands, runs("boom"), []string{first.String()}, host, up, 3, daemon.Ship{}); !slices.Equal(got, []string{second.String()}) {
		t.Errorf("with %s pending picked %v", first, got)
	}
}

// A start whose run journals leaves pending. One that never journals is
// held past the grace, taking no slot, until a new journal of its change
// set shows, 15 minutes pass since the start or it is no longer accepted.
func TestDaemonHeld(t *testing.T) {
	yes := func(string) bool { return true }
	at := time.Unix(1700000000, 0)
	old := pickJournal("B")
	pending := map[string]daemonPending{"A": {at: at}, "B": {at: at, runs: 1}}
	unjournaled := map[string]daemonPending{}

	if got := daemonHeld(pending, unjournaled, []run.Journal{old}, yes, at.Add(time.Second)); !slices.Equal(got, []string{"A", "B"}) {
		t.Fatalf("in the grace held %v", got)
	}
	if got := daemonHeld(pending, unjournaled, []run.Journal{old, pickJournal("A")}, yes, at.Add(time.Second)); !slices.Equal(got, []string{"B"}) {
		t.Fatalf("after A journaled held %v", got)
	}
	if got := daemonHeld(pending, unjournaled, []run.Journal{old}, yes, at.Add(daemonStartGrace)); len(got) != 0 || len(pending) != 0 {
		t.Fatalf("past the grace held %v, pending %v", got, pending)
	}
	if _, ok := unjournaled["B"]; !ok || len(unjournaled) != 1 {
		t.Fatalf("unjournaled = %v", unjournaled)
	}
	daemonHeld(pending, unjournaled, []run.Journal{old}, yes, at.Add(daemonUnjournaledHold-time.Second))
	if _, ok := unjournaled["B"]; !ok {
		t.Fatalf("B released before 15 minutes")
	}
	daemonHeld(pending, unjournaled, []run.Journal{old}, yes, at.Add(daemonUnjournaledHold))
	if len(unjournaled) != 0 {
		t.Fatalf("B still held at 15 minutes: %v", unjournaled)
	}
	// Picked again, it is pending, and a start with no journal is held anew.
	later := at.Add(daemonUnjournaledHold)
	pending["B"] = daemonPending{at: later, runs: 1}
	if got := daemonHeld(pending, unjournaled, []run.Journal{old}, yes, later.Add(time.Second)); !slices.Equal(got, []string{"B"}) {
		t.Fatalf("after the re-pick held %v", got)
	}
	// A new journal releases it early.
	unjournaled["B"] = daemonPending{at: later, runs: 1}
	daemonHeld(pending, unjournaled, []run.Journal{old, pickJournal("B")}, yes, later.Add(time.Minute))
	if len(unjournaled) != 0 {
		t.Fatalf("B still unjournaled after a new journal: %v", unjournaled)
	}
	// A landed or withdrawn change set leaves the hold at once.
	unjournaled["B"] = daemonPending{at: later, runs: 1}
	daemonHeld(pending, unjournaled, []run.Journal{old}, func(string) bool { return false }, later.Add(time.Minute))
	if len(unjournaled) != 0 {
		t.Fatalf("B still held after it landed: %v", unjournaled)
	}
}

// The backlog rule orders the backlog against the milestone work: after
// (the default) behind it, before ahead of it, oldest accepted first
// inside; skip and an unknown value skip it as no-milestone, the unknown
// one with its reason in the detail.
func TestDaemonCandidatesBacklog(t *testing.T) {
	f := newPlanFix()
	m := f.intent("m", id.ID{}, "a")
	f.milestone(m)
	plain := f.intent("plain", id.ID{}, "b")
	mk := func(name string, serves id.ID) id.ID {
		cs := f.cs(name, changeset.StateAccepted)
		f.serve(cs, serves)
		return cs
	}
	m1, m2 := mk("m1", m), mk("m2", m)
	b1, b2 := mk("b1", plain), mk("b2", plain)
	since := map[id.ID]time.Time{b1: time.Unix(2000, 0), b2: time.Unix(1000, 0)}
	yes := func(id.ID) bool { return true }
	strs := func(ids ...id.ID) []string {
		var out []string
		for _, i := range ids {
			out = append(out, i.String())
		}
		return out
	}
	sched := func(backlog string, detail ...string) daemonSchedule {
		return daemonSchedule{Backlog: backlog, Detail: detail, Since: func(cs id.ID) time.Time { return since[cs] }}
	}
	for backlog, want := range map[string][]string{
		policy.BacklogAfter:  strs(m1, m2, b2, b1),
		policy.BacklogBefore: strs(b2, b1, m1, m2),
		policy.BacklogSkip:   strs(m1, m2),
	} {
		cands, skips := daemonCandidates(f.core(), yes, sched(backlog))
		if !slices.Equal(cands, want) {
			t.Errorf("%s: candidates = %v, want %v", backlog, cands, want)
		}
		if backlog == policy.BacklogSkip && len(skips) != 2 {
			t.Errorf("skip: skips = %v", skips)
		}
		if backlog != policy.BacklogSkip && len(skips) != 0 {
			t.Errorf("%s: skips = %v", backlog, skips)
		}
	}
	cands, skips := daemonCandidates(f.core(), yes, sched("", "bad value"))
	if !slices.Equal(cands, strs(m1, m2)) || len(skips) != 2 {
		t.Fatalf("unknown: candidates %v skips %v", cands, skips)
	}
	for _, sk := range skips {
		if sk.Reason != daemon.SkipNoMilestone || !slices.Equal(sk.Detail, []string{"bad value"}) {
			t.Errorf("unknown: skip = %+v", sk)
		}
	}
}

package daemon

import (
	"reflect"
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
	accepted := func(string) string { return "accepted" }
	got := Resumable(runs, "h", dead, locked, accepted, now)
	if len(got) != 2 || got[0] != "interrupted" || got[1] != "unrecorded" {
		t.Fatalf("Resumable = %v, want [interrupted unrecorded]", got)
	}
}

func TestResumableOnlyAcceptedOrLandedPastLand(t *testing.T) {
	dead := func(int) bool { return false }
	proc := run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}}
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	states := map[string]string{}
	var runs []run.Journal
	add := func(cs, state, workflow string, land bool) {
		states[cs] = state
		es := []run.Entry{{Type: run.TypeStarted, At: now.Add(-time.Hour), Started: &run.Started{Change: cs, Workflow: workflow}}, proc}
		if land {
			es = append(es, run.Entry{Type: run.TypeEntered, Step: "land"})
		}
		runs = append(runs, run.Journal{ID: "run-" + cs, Entries: es})
	}
	add("accepted", "accepted", "ship-change", false)
	add("landed-past", "landed", "ship-change", true)
	add("landed-never", "landed", "ship-change", false)
	add("landed-plan", "landed", "plan", false)
	add("withdrawn-plan", "withdrawn", "plan", false)
	add("rejected-plan", "rejected", "plan", false)
	add("superseded-plan", "superseded", "plan", false)
	for _, s := range []string{"rejected", "withdrawn", "proposed", "draft"} {
		add(s, s, "ship-change", true)
	}
	got := Resumable(runs, "h", dead, func(string) bool { return false }, func(cs string) string { return states[cs] }, now)
	if len(got) != 3 || got[0] != "run-accepted" || got[1] != "run-landed-past" || got[2] != "run-landed-plan" {
		t.Fatalf("Resumable = %v, want [run-accepted run-landed-past run-landed-plan]", got)
	}
}

func TestSkips(t *testing.T) {
	const host = "h"
	up := func(pid int) bool { return pid == 1 }
	at := time.Unix(1700000000, 0)
	journal := func(cs string, es ...run.Entry) run.Journal {
		return run.Journal{ID: "R" + cs, Entries: append([]run.Entry{{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs}}}, es...)}
	}
	proc := func(pid int) run.Entry {
		return run.Entry{Type: run.TypeProcess, At: at, Process: &run.Process{Host: host, Pid: pid}}
	}
	fail := func(detail string) run.Entry {
		return run.Entry{Type: run.TypeFinished, At: at, End: run.Failed, Detail: detail}
	}
	retried := run.Entry{Type: run.TypeFinished, At: at, End: run.Failed, Detail: "boom"}
	runs := []run.Journal{
		journal("running", proc(1)),
		journal("interrupted", proc(2)),
		journal("failed", proc(3), fail("boom")),
		journal("abandoned", proc(3), fail(run.AbandonedPrefix+"me")),
		// Failed, then a newer run that is not failed: retried.
		journal("retried", proc(3), retried),
		journal("retried", proc(1)),
	}
	// Ship runs: shipped ended after the accept; reopened ended before a
	// reopen; plan is a succeeded run of another workflow.
	accepted := at.Add(-time.Hour)
	done := func(cs, wf string, end time.Time) run.Journal {
		return run.Journal{ID: "R" + cs, Entries: []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs, Workflow: wf}},
			{Type: run.TypeFinished, At: end, End: run.Succeeded},
		}}
	}
	runs = append(runs,
		done("shipped", "ship", at.Add(time.Minute)),
		done("plan", "plan", at.Add(time.Minute)),
		done("reopened", "ship", at.Add(time.Minute)),
		done("stale", "ship", accepted.Add(-time.Minute)),
	)
	since := map[string]time.Time{"shipped": accepted, "plan": accepted, "reopened": at.Add(2 * time.Minute), "stale": accepted}
	ship := Ship{Workflow: "ship", Since: func(cs string) time.Time { return since[cs] }}
	cands := []string{"running", "interrupted", "failed", "abandoned", "retried", "held", "fresh", "late", "shipped", "plan", "reopened", "stale"}
	cases := []struct {
		name string
		max  int
		want []Skip
	}{
		{"one slot", 1, []Skip{
			{Change: "running", Reason: SkipLive}, {Change: "interrupted", Reason: SkipLive},
			{Change: "failed", Reason: SkipFailed}, {Change: "abandoned", Reason: SkipNoSlot},
			{Change: "retried", Reason: SkipLive}, {Change: "held", Reason: SkipHeld},
			{Change: "fresh", Reason: SkipNoSlot}, {Change: "late", Reason: SkipNoSlot},
			{Change: "shipped", Reason: SkipShipped}, {Change: "plan", Reason: SkipNoSlot},
			{Change: "reopened", Reason: SkipNoSlot}, {Change: "stale", Reason: SkipNoSlot},
		}},
		{"room", 9, []Skip{
			{Change: "running", Reason: SkipLive}, {Change: "interrupted", Reason: SkipLive},
			{Change: "failed", Reason: SkipFailed}, {Change: "retried", Reason: SkipLive},
			{Change: "held", Reason: SkipHeld}, {Change: "shipped", Reason: SkipShipped},
		}},
	}
	for _, c := range cases {
		got := Skips(cands, runs, []string{"held"}, host, up, c.max, ship)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: Skips = %v, want %v", c.name, got, c.want)
		}
		picked := Pick(cands, runs, []string{"held"}, host, up, c.max, ship)
		if len(picked)+len(got) != len(cands) {
			t.Errorf("%s: picked %v and skipped %v do not split %v", c.name, picked, got, cands)
		}
	}
}

func TestPickSkipsShipped(t *testing.T) {
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	accepted := at.Add(-time.Hour)
	done := func(cs, wf string) run.Journal {
		return run.Journal{ID: "R" + cs, Entries: []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs, Workflow: wf}},
			{Type: run.TypeFinished, At: at.Add(time.Minute), End: run.Succeeded},
		}}
	}
	since := map[string]time.Time{
		"shipped": accepted, "plan": accepted,
		// Reopened after the ship run ended.
		"reopened": at.Add(2 * time.Minute),
	}
	ship := Ship{Workflow: "ship", Since: func(cs string) time.Time { return since[cs] }}
	runs := []run.Journal{done("shipped", "ship"), done("plan", "plan"), done("reopened", "ship")}
	cases := []struct {
		name, change string
		picked       bool
	}{
		{"succeeded ship run after the accept is skipped", "shipped", false},
		{"succeeded plan run is picked", "plan", true},
		{"succeeded ship run followed by a reopen is picked", "reopened", true},
	}
	for _, c := range cases {
		got := Pick([]string{c.change}, runs, nil, "h", func(int) bool { return false }, 2, ship)
		if (len(got) == 1) != c.picked {
			t.Errorf("%s: Pick = %v", c.name, got)
		}
	}
}

// A request's plan run is resumed after its change set lands early whatever
// workflow it runs; a non-plan run that never reached land is not.
func TestResumablePlanRunUnderOtherWorkflow(t *testing.T) {
	dead := func(int) bool { return false }
	proc := run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: "h", Pid: 1}}
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	mk := func(id string, plan bool) run.Journal {
		s := &run.Started{Change: id, Workflow: "bugfix", Plan: plan}
		return run.Journal{ID: id, Entries: []run.Entry{{Type: run.TypeStarted, At: now.Add(-time.Hour), Started: s}, proc}}
	}
	runs := []run.Journal{mk("plan", true), mk("other", false)}
	got := Resumable(runs, "h", dead, func(string) bool { return false }, func(string) string { return "landed" }, now)
	if len(got) != 1 || got[0] != "plan" {
		t.Fatalf("Resumable = %v, want [plan]", got)
	}
}

// With For set, a change set has shipped when a succeeded run of the workflow
// dispatch gives it ended, not of Workflow.
func TestPickShippedUnderDispatchedWorkflow(t *testing.T) {
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	done := func(cs, wf string) run.Journal {
		return run.Journal{ID: "R" + cs, Entries: []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs, Workflow: wf}},
			{Type: run.TypeFinished, At: at.Add(time.Minute), End: run.Succeeded},
		}}
	}
	ship := Ship{Workflow: "ship", Since: func(string) time.Time { return at.Add(-time.Hour) },
		For: func(cs string) string {
			if cs == "bug" {
				return "bugfix"
			}
			return "ship"
		}}
	runs := []run.Journal{done("bug", "bugfix"), done("plain", "bugfix")}
	got := Pick([]string{"bug", "plain"}, runs, nil, "h", func(int) bool { return false }, 2, ship)
	if len(got) != 1 || got[0] != "plain" {
		t.Fatalf("Pick = %v, want [plain]", got)
	}
}

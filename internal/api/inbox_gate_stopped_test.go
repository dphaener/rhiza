package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A plan-milestone run stopped at its gate on an open question is an open
// item whether or not a process still holds the run's lock (one lingering
// after the stop), until a process journals itself after the stop (a
// resume), and its question carries the statement even when a change set
// that serves no plan intent (the review's) created it.
func TestInboxListsGateStoppedPlanRun(t *testing.T) {
	root := t.TempDir()
	ctx := testCtx(root)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	call("init", `{}`)
	// The review's change set asks the question and serves nothing.
	call("change.new", `{"intent":"Review round"}`)
	call("add", `{"kind":"question","handle":"review-r3-f1","fields":["statement=Should the pages show plain commands?"]}`)
	own := call("change.new", `{"intent":"Milestone plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"plan-intent"}`)

	at := time.Date(2026, 11, 10, 1, 0, 0, 0, time.UTC)
	stopped := run.OpenStop + "review-r3-f1; answer each with a Decision, then rhiza run resume R1"
	entries := func(extra ...run.Entry) run.Journal {
		es := []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan-milestone", Change: own.String(), Worktree: root, Plan: true, MilestonePlan: true}},
			{Type: run.TypeProcess, At: at.Add(time.Second), Process: &run.Process{Host: "h", Pid: 1}},
			{Type: run.TypeEntered, At: at.Add(2 * time.Second), Step: "gate", Visit: 1},
			{Type: run.TypeGated, At: at.Add(3 * time.Second), Changes: []string{own.String()}},
			{Type: run.TypeStopped, At: at.Add(10 * time.Second), Step: "gate", Detail: stopped},
		}
		return run.Journal{ID: "R1", Entries: append(es, extra...)}
	}
	state := func(string) (string, bool) { return "accepted", true }
	open := func(j run.Journal, live bool) []InboxItem {
		var out []InboxItem
		for _, it := range inboxItems([]run.Journal{j}, "h", nil, func(string) bool { return live }, state, nil, nil, nil) {
			if it.Kind == InboxOpen {
				out = append(out, it)
			}
		}
		return out
	}
	if got := open(entries(), false); len(got) != 1 {
		t.Fatalf("a stopped run lists %d open items, want 1", len(got))
	}
	// A lingering holder of the lock does not hide the stop.
	if got := open(entries(), true); len(got) != 1 {
		t.Fatalf("a stopped run whose lock is held lists %d open items, want 1", len(got))
	}
	// A process journalled after the stop is a resume.
	resumed := entries(run.Entry{Type: run.TypeProcess, At: at.Add(time.Minute), Process: &run.Process{Host: "h", Pid: 2}})
	if got := open(resumed, true); len(got) != 0 {
		t.Fatalf("a run being resumed lists %d open items, want 0", len(got))
	}
	// The question's statement and handle come from the change set that
	// created it.
	qs := inboxGateQuestions(ctx, entries().State())
	if len(qs) != 1 || qs[0].Handle != "review-r3-f1" || qs[0].Statement != "Should the pages show plain commands?" {
		t.Fatalf("questions = %+v", qs)
	}
}

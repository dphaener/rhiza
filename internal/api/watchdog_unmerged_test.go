package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogGreenUnmerged drives the green-but-unmerged pass with a fake
// gh read: an autonomous failed run at its merge step is resumed, a live run
// is left alone, and a succeeded run, or a run that needs the owner before
// its merge, gets an inbox item with no resume.
func TestWatchdogGreenUnmerged(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = `name = "w"
[step.wait]
kind = "wait"
until = "rhiza ship check --settled"
timeout = "1h"
on.pass = "merge"
[step.merge]
kind = "command"
run = "rhiza ship merge"
`
	const gated = `name = "g"
[step.wait]
kind = "wait"
until = "rhiza ship check --settled"
timeout = "1h"
on.pass = "signoff"
[step.review]
kind = "checks"
on.pass = "signoff"
[[step.review.check]]
name = "code"
prompt = "check {check}"
[step.review.critic]
prompt = "criticise {check}"
[step.review.fixer]
prompt = "fix {leaf}"
[step.review.verifier]
prompt = "verify {leaf}"
[step.signoff]
kind = "signoff"
on.pass = "merge"
on.changes = "review"
[step.merge]
kind = "command"
run = "rhiza ship merge"
`
	started := func(d string) run.Entry {
		return run.Entry{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: d}}
	}
	failedAt := func(id, d, step string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			started(d),
			{Type: run.TypeEntered, At: t0, Step: step, Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: step, Visit: 1, Outcome: "fail", Detail: "exit 1"},
			{Type: run.TypeFinished, At: t0, Step: step, End: "failed", Detail: "step " + step + " failed"},
		}}
	}
	js := []run.Journal{
		failedAt("merge", def, "merge"),
		failedAt("live", def, "merge"),
		failedAt("gate", gated, "wait"),
		{ID: "ok", Entries: []run.Entry{
			started(def),
			{Type: run.TypeEntered, At: t0, Step: "merge", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "merge", Visit: 1, Outcome: "pass"},
			{Type: run.TypeFinished, At: t0, Step: "merge", End: "succeeded"},
		}},
	}
	green := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS", RunID: "7"}}}
	pr := green
	var resumed, log []string
	d := watchdogDeps{
		ci:     func(run.State) (watchdogPR, error) { return pr, nil },
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	live := func(id string) bool { return id == "live" }
	mem, holds, err := watchdogGreenUnmerged(pol, js, live, nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(resumed, ",") != "merge" {
		t.Fatalf("resumed %v, want only merge", resumed)
	}
	got := map[string]bool{}
	for _, h := range holds {
		got[h.Run] = true
	}
	if len(holds) != 2 || !got["gate"] || !got["ok"] {
		t.Fatalf("holds %+v, want gate and ok", holds)
	}
	// Budget of one: the next poll resumes nothing and holds the run too.
	resumed = nil
	_, holds, err = watchdogGreenUnmerged(pol, js, live, mem, d)
	if err != nil || len(resumed) != 0 || len(holds) != 3 {
		t.Fatalf("over budget: resumed %v holds %d err %v", resumed, len(holds), err)
	}
	// Red or absent checks, or no open PR: nothing.
	for _, p := range []watchdogPR{{}, {Head: "aaaa111"}, {Head: "aaaa111", Checks: []ShipCheck{{Name: "t", Status: "COMPLETED", Conclusion: "FAILURE"}}}} {
		pr = p
		resumed = nil
		_, holds, _ = watchdogGreenUnmerged(pol, js, live, nil, d)
		if len(resumed) != 0 || len(holds) != 0 {
			t.Errorf("pr %+v: resumed %v holds %v", p, resumed, holds)
		}
	}
}

package api

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
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
	// Budget of one, spent on an attempt: the next poll resumes nothing and
	// leaves the run to the budget hold, unread.
	resumed = nil
	_, holds, err = watchdogGreenUnmerged(pol, js, live, mem, d)
	if err != nil || len(resumed) != 0 || len(holds) != 2 {
		t.Fatalf("over budget: resumed %v holds %d err %v", resumed, len(holds), err)
	}
	// Spent without an attempt on record (a budget of zero): held by the pass.
	_, holds, _ = watchdogGreenUnmerged(WatchdogPolicy{Enabled: true, AutonomousOnly: true}, js, live, nil, d)
	if len(holds) != 3 {
		t.Fatalf("budget zero: holds %d, want 3", len(holds))
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

// TestWatchdogGreenUnmergedReadsOnce checks that the pass's gh reads stay
// bounded: a run found with no open PR is not read again on later polls,
// through watchdogPasses as the daemon runs them, until its journal moves.
func TestWatchdogGreenUnmergedReadsOnce(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = `name = "w"
[step.merge]
kind = "command"
run = "rhiza ship merge"
`
	ok := run.Journal{ID: "ok", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
		{Type: run.TypeEntered, At: t0, Step: "merge", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "merge", Visit: 1, Outcome: "pass"},
		{Type: run.TypeFinished, At: t0, Step: "merge", End: "succeeded"},
	}}
	reads := 0
	d := watchdogDeps{
		ci:     func(run.State) (watchdogPR, error) { reads++; return watchdogPR{}, nil },
		resume: func(string) error { t.Fatal("resumed"); return nil },
		log:    func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	none := func(string) bool { return false }
	js := []run.Journal{ok}
	var mem []daemon.WatchdogRun
	for range 3 {
		var err error
		mem, _, err = watchdogPasses(pol, js, none, t0, mem, d)
		if err != nil {
			t.Fatal(err)
		}
	}
	if reads != 1 {
		t.Fatalf("%d PR reads over three polls, want 1", reads)
	}
	// A journal that moved on (here: a later entry) is read again, once.
	moved := run.Journal{ID: "ok", Entries: append(append([]run.Entry{}, ok.Entries...),
		run.Entry{Type: run.TypeFinished, At: t0.Add(time.Minute), Step: "merge", End: "succeeded"})}
	js = []run.Journal{moved}
	for range 2 {
		mem, _, _ = watchdogPasses(pol, js, none, t0, mem, d)
	}
	if reads != 2 {
		t.Fatalf("%d PR reads after the journal moved, want 2", reads)
	}
}

// TestWatchdogReadPR checks the shared PR read gives a run the PR only when
// the change set's ship record names that run: another run's, or none, is no
// PR, and the owning run still reads it.
func TestWatchdogReadPR(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(greenRollup)
	g.answers["run list --commit abc1234def --json databaseId,attempt"] = `[]`
	cs := strings.TrimSuffix(filepath.Base(state), ".json")
	read := func(runID string) watchdogPR {
		t.Helper()
		pr, err := watchdogReadPR(f.ctx, run.State{ID: runID, Change: cs})
		if err != nil {
			t.Fatal(err)
		}
		return pr
	}
	// No record, then a record naming no run: no run owns the PR.
	if pr := read("SHIP"); pr.Head != "" {
		t.Errorf("no record: head %q, want none", pr.Head)
	}
	if err := writeShipState(state, shipState{PR: 7}); err != nil {
		t.Fatal(err)
	}
	if pr := read("SHIP"); pr.Head != "" {
		t.Errorf("record naming no run: head %q, want none", pr.Head)
	}
	// A record naming the ship run: it owns the PR, the plan run does not.
	if err := writeShipState(state, shipState{PR: 7, Run: "SHIP"}); err != nil {
		t.Fatal(err)
	}
	if pr := read("SHIP"); pr.Head != "abc1234def" {
		t.Errorf("owner: head %q, want abc1234def", pr.Head)
	}
	if pr := read("PLAN"); pr.Head != "" {
		t.Errorf("other run: head %q, want none", pr.Head)
	}
}

// TestWatchdogGreenUnmergedSkipsBusyChange checks a succeeded run of a
// change set whose other run is live, or has a continue queued, gets no
// hold and no resume.
func TestWatchdogGreenUnmergedSkipsBusyChange(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = `name = "w"
[step.merge]
kind = "command"
run = "rhiza ship merge"
`
	mk := func(id, change string, ended bool) run.Journal {
		es := []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def, Change: change}},
			{Type: run.TypeEntered, At: t0, Step: "merge", Visit: 1},
		}
		if ended {
			es = append(es,
				run.Entry{Type: run.TypeEnded, At: t0, Step: "merge", Visit: 1, Outcome: "pass"},
				run.Entry{Type: run.TypeFinished, At: t0, Step: "merge", End: "succeeded"})
		}
		return run.Journal{ID: id, Entries: es}
	}
	js := []run.Journal{mk("plan", "CS", true), mk("ship", "CS", false)}
	green := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS", RunID: "7"}}}
	var resumed []string
	d := watchdogDeps{
		ci:     func(run.State) (watchdogPR, error) { return green, nil },
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	live := func(id string) bool { return id == "ship" }
	if _, holds, err := watchdogGreenUnmerged(pol, js, live, nil, d); err != nil || len(holds) != 0 || len(resumed) != 0 {
		t.Fatalf("live run: holds %v resumed %v err %v", holds, resumed, err)
	}
	d.queued = func(id string) bool { return id == "ship" }
	if _, holds, err := watchdogGreenUnmerged(pol, js, func(string) bool { return false }, nil, d); err != nil || len(holds) != 0 || len(resumed) != 0 {
		t.Fatalf("queued run: holds %v resumed %v err %v", holds, resumed, err)
	}
	// Nothing else running: the succeeded run is held as before.
	d.queued = nil
	if _, holds, _ := watchdogGreenUnmerged(pol, js[:1], func(string) bool { return false }, nil, d); len(holds) != 1 {
		t.Fatalf("idle: holds %v, want 1", holds)
	}
}

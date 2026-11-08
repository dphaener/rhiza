package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogLog fakes a transient failure and a stalled-attended hold over
// several polls: each action is logged as one line, kept in the journal on
// the state newest first with a repeated hold as one action, and health
// lists them.
func TestWatchdogLog(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	js := []run.Journal{{ID: "gh", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\ngh: HTTP 502 Bad Gateway"},
		{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
	}}}
	now := t0.Add(time.Hour)
	var lines []string
	var state daemon.State
	d := watchdogDeps{
		resume: func(string) error { return nil },
		now:    func() time.Time { return now },
	}
	d.log = watchdogJournalLog(func(f string, a ...any) { lines = append(lines, fmt.Sprintf(f, a...)) }, &state, d.now)
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 3, AutonomousOnly: true}
	for range 2 {
		var err error
		if state.Watchdog, state.WatchdogHolds, err = watchdogPasses(pol, js, func(string) bool { return false }, now, state.Watchdog, d); err != nil {
			t.Fatal(err)
		}
		now = now.Add(time.Minute)
	}
	if len(lines) != 2 {
		t.Fatalf("logged %d lines, want one per action: %q", len(lines), lines)
	}
	acts := state.WatchdogActions
	if len(acts) != 2 || !acts[0].At.After(acts[1].At) {
		t.Fatalf("journal %+v, want two actions, newest first", acts)
	}
	for _, a := range acts {
		if strings.HasPrefix(a.Line, "watchdog:") || !strings.Contains(a.Line, "run gh") || !strings.Contains(a.Line, "502") || !strings.Contains(a.Line, "resumed with --retry") {
			t.Errorf("action %q lacks what it saw, did and the run", a.Line)
		}
	}
	// A line logged again replaces its earlier entry, and the journal is bounded.
	state.WatchdogActions = watchdogRecord(state.WatchdogActions, now, acts[1].Line)
	if len(state.WatchdogActions) != 2 || state.WatchdogActions[0].Line != acts[1].Line {
		t.Fatalf("repeat %+v", state.WatchdogActions)
	}
	for i := range 30 {
		state.WatchdogActions = watchdogRecord(state.WatchdogActions, now, strings.Repeat("x", i+1))
	}
	if len(state.WatchdogActions) != watchdogActionsKept {
		t.Fatalf("kept %d, want %d", len(state.WatchdogActions), watchdogActionsKept)
	}
	h := healthWatchdogActions(HealthPart{}, watchdogActionLines(acts))
	if len(h.Detail) != 3 || !strings.Contains(h.Detail[1], "2030-01-01T01:01:00Z") {
		t.Fatalf("health detail %q", h.Detail)
	}
}

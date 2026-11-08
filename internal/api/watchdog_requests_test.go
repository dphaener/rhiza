package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogRequests drives the stuck-request pass: a request whose plan
// run is gone is resumed, one whose run failed on a transient signature is
// retried, one failed for another reason is left for the owner, and the
// budget holds.
func TestWatchdogRequests(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"p\"\n[step.plan]\nkind = \"command\"\nrun = \"true\"\n"
	started := run.Entry{Type: run.TypeStarted, At: t0, Started: &run.Started{Definition: def}}
	entered := run.Entry{Type: run.TypeEntered, At: t0, Step: "plan", Visit: 1}
	failed := func(id, detail string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{started, entered,
			{Type: run.TypeEnded, At: t0, Step: "plan", Visit: 1, Outcome: "fail", Detail: detail},
			{Type: run.TypeFinished, At: t0, Step: "plan", End: "failed", Detail: detail}}}
	}
	js := []run.Journal{
		{ID: "gone", Entries: []run.Entry{started, entered}},
		{ID: "alive", Entries: []run.Entry{started, entered}},
		failed("flake", "fatal: Unable to create '/x/.git/index.lock': File exists."),
		failed("other", "exit 1"),
	}
	reqs := []daemon.Request{
		{ID: "r1", Kind: daemon.KindFeature, Run: "gone"},
		{ID: "r2", Kind: daemon.KindFeature, Run: "alive"},
		{ID: "r3", Kind: daemon.KindBug, Run: "flake"},
		{ID: "r4", Kind: daemon.KindFeature, Run: "other"},
		{ID: "r5", Kind: daemon.KindMilestone, Run: "gone"},
		{ID: "r6", Kind: daemon.KindFeature},
	}
	var resumed, log []string
	d := watchdogDeps{
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	live := func(id string) bool { return id == "alive" }
	mem, err := watchdogRequests(pol, reqs, js, live, nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(resumed, ","); got != "gone,flake" {
		t.Fatalf("resumed %q, want gone,flake", got)
	}
	if len(log) != 2 || !strings.Contains(log[0], "request r1") || !strings.Contains(log[1], "index.lock") {
		t.Fatalf("log %q", log)
	}
	resumed = nil
	if _, err := watchdogRequests(pol, reqs, js, live, mem, d); err != nil || len(resumed) != 0 {
		t.Fatalf("over budget: resumed %v err %v", resumed, err)
	}
}

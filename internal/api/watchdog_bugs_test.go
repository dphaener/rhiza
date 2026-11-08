package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogBugs drives the bug filing: a transient-failure recovery
// queues one bug request with the run and request IDs, the step and the log
// lines; a second recovery of the same kind and step queues none while the
// first request is open; a CI flake queues none; fileBugs = false queues
// none.
func TestWatchdogBugs(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"w\"\n[step.build]\nkind = \"command\"\nrun = \"true\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\n"
	failed := func(id, step, detail string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
			{Type: run.TypeEntered, At: t0, Step: step, Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: step, Visit: 1, Outcome: "fail", Detail: detail},
			{Type: run.TypeFinished, At: t0, Step: step, End: "failed", Detail: "step " + step + " failed"},
		}}
	}
	lock := "exit 1\nfatal: Unable to create '/x/.git/index.lock': File exists."
	js := []run.Journal{failed("one", "build", lock), failed("two", "build", lock)}
	queue := []daemon.Request{{ID: "req1", Kind: daemon.KindFeature, Run: "one"}}
	d := watchdogDeps{
		resume:   func(string) error { return nil },
		requests: func() ([]daemon.Request, error) { return queue, nil },
		file: func(text string) error {
			queue = append(queue, daemon.Request{ID: "bug" + string(rune('a'+len(queue))), Kind: daemon.KindBug, Text: text})
			return nil
		},
		log: func(string, ...any) {},
	}
	nolive := func(string) bool { return false }
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 5, CIReruns: 1, FileBugs: true, AutonomousOnly: true}
	pass := func(pol WatchdogPolicy, js []run.Journal) {
		t.Helper()
		if _, _, err := watchdogPasses(pol, js, nolive, t0.Add(time.Hour), nil, d); err != nil {
			t.Fatal(err)
		}
	}
	bugs := func() []daemon.Request {
		var out []daemon.Request
		for _, r := range queue {
			if r.Kind == daemon.KindBug {
				out = append(out, r)
			}
		}
		return out
	}

	pass(pol, js[:1])
	b := bugs()
	if len(b) != 1 {
		t.Fatalf("want one bug request, got %+v", b)
	}
	for _, want := range []string{"transient failure at step build", "run one", "request req1", "index.lock"} {
		if !strings.Contains(b[0].Text, want) {
			t.Errorf("bug text %q lacks %q", b[0].Text, want)
		}
	}
	// The same kind and step, another run, while the first is open.
	pass(pol, js[1:])
	if n := len(bugs()); n != 1 {
		t.Fatalf("second matching recovery queued a request: %d bugs", n)
	}
	// A CI flake is not a bug.
	queue = nil
	d.ci = func(run.State) (watchdogPR, error) {
		return watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "FAILURE", RunID: "7"}}, Attempts: map[string]int{"7": 1}}, nil
	}
	d.rerun = func(run.State, string) error { return nil }
	pass(pol, []run.Journal{failed("flake", "wait", "failed: test")})
	if n := len(bugs()); n != 0 {
		t.Fatalf("CI flake queued %d bugs", n)
	}
	// fileBugs = false.
	off := pol
	off.FileBugs = false
	pass(off, js[:1])
	if n := len(bugs()); n != 0 {
		t.Fatalf("fileBugs = false queued %d bugs", n)
	}
}

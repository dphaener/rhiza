package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/daemon"
)

// TestFlakyTestBugs drives the flaky-test bug filing: one bug request per
// newly flaky test naming the test, its runs and heads; none for a second
// flake of a test while its request is open; two for two tests; none with
// fileBugs = false.
func TestFlakyTestBugs(t *testing.T) {
	var queue []daemon.Request
	var runs []CIRun
	d := watchdogDeps{
		requests: func() ([]daemon.Request, error) { return queue, nil },
		file: func(text string) error {
			queue = append(queue, daemon.Request{ID: "b" + string(rune('a'+len(queue))), Kind: daemon.KindBug, Text: text})
			return nil
		},
		flakes: func() ([]CIRun, error) { return runs, nil },
		log:    func(string, ...any) {},
	}
	pol := WatchdogPolicy{FileBugs: true}
	poll := func(pol WatchdogPolicy) {
		t.Helper()
		if err := watchdogFlakyTests(pol, d); err != nil {
			t.Fatal(err)
		}
	}

	runs = []CIRun{{Head: "aaa111", PR: 7, Result: ciGreen, Flakes: []string{"TestA"}}}
	poll(pol)
	if len(queue) != 1 {
		t.Fatalf("want one bug request, got %+v", queue)
	}
	for _, want := range []string{"flaky test TestA:", "PR #7", "aaa111", "locally", "quarantine"} {
		if !strings.Contains(queue[0].Text, want) {
			t.Errorf("bug text %q lacks %q", queue[0].Text, want)
		}
	}
	// A repeated poll and a second flake of the same test queue none.
	poll(pol)
	runs = append(runs, CIRun{Head: "bbb222", PR: 8, Result: ciGreen, Flakes: []string{"TestA"}})
	poll(pol)
	if len(queue) != 1 {
		t.Fatalf("same test queued again: %+v", queue)
	}
	// Two flaky tests, one already open: one more request.
	runs = append(runs, CIRun{Head: "ccc333", PR: 9, Result: ciGreen, Flakes: []string{"TestB", "TestAB"}})
	poll(pol)
	if len(queue) != 3 {
		t.Fatalf("want TestB and TestAB queued once each: %+v", queue)
	}
	// fileBugs = false queues none.
	queue = nil
	poll(WatchdogPolicy{})
	if len(queue) != 0 {
		t.Fatalf("fileBugs off queued: %+v", queue)
	}
}

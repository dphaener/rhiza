package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestFlakyTestBugs drives the flaky-test bug filing: one bug request per
// newly flaky test naming the test, its runs and heads; none for a second
// flake of a test while its request is open; two for two tests; none with
// fileBugs = false; and none for the flakes a planned request named.
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
	runs = append(runs, CIRun{Head: "bbb222", PR: 8, Result: ciGreen, Time: "2026-11-08T11:00:00Z", Flakes: []string{"TestA"}})
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
	// Once TestA's request is planned, the flakes it named queue nothing
	// more, poll after poll; only a later flake of TestA files it again.
	// bbb222 flaked while it was open, so it is covered too.
	queue[0].Run = "R1"
	end := time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC)
	d.all = []run.Journal{{ID: "R1", Entries: []run.Entry{{Type: run.TypeFinished, At: end, End: run.Succeeded}}}}
	if state, _ := daemon.RequestState(queue[0], d.all); state != daemon.RequestPlanned {
		t.Fatalf("request state %s, want planned", state)
	}
	poll(pol)
	poll(pol)
	if len(queue) != 3 {
		t.Fatalf("planned request's flakes queued again: %+v", queue)
	}
	runs = append(runs, CIRun{Head: "ddd444", PR: 10, Result: ciGreen, Time: "2026-11-08T13:00:00Z", Flakes: []string{"TestA"}})
	poll(pol)
	poll(pol)
	if len(queue) != 4 || !strings.Contains(queue[3].Text, "ddd444") {
		t.Fatalf("want one new TestA request for the later flake: %+v", queue)
	}
	// fileBugs = false queues none.
	queue = nil
	poll(WatchdogPolicy{})
	if len(queue) != 0 {
		t.Fatalf("fileBugs off queued: %+v", queue)
	}
}

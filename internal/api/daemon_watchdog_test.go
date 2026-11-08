package api

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// TestWatchdogStalled drives the stalled-step recovery with a fake clock,
// fake sessions and a fake run driver.
func TestWatchdogStalled(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	sessions := session.NewStore(t.TempDir())
	report := func(id string, auto bool, at time.Time) {
		rec := session.Record{ID: id, Autonomous: auto}
		rec.Status.Hook = &session.Report{State: session.StateWorking, At: at}
		if err := sessions.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	journal := func(id string, auto bool, sid string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: auto}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeSession, At: t0, Step: "build", Session: sid},
		}}
	}
	var log, sent, stopped, resumed []string
	d := watchdogDeps{
		sessions: sessions,
		step:     func(run.State) (string, time.Duration) { return "agent", 0 },
		send:     func(sid, _ string) error { sent = append(sent, sid); return nil },
		stop:     func(st run.State, _ []string) error { stopped = append(stopped, st.ID); return nil },
		resume:   func(id string) error { resumed = append(resumed, id); return nil },
		log:      func(f string, a ...any) { log = append(log, strings.TrimSpace(fmt.Sprintf(f, a...))) },
	}
	pol := watchdogPolicyOf(nil)
	_ = pol
	pol = WatchdogPolicy{Enabled: true, AgentNudge: 20 * time.Minute, AgentRetry: 40 * time.Minute, CommandLimit: 30 * time.Minute, WaitLimit: 90 * time.Minute, RetryBudget: 1, AutonomousOnly: true}

	report("auto", true, t0)
	report("busy", true, t0)
	report("owner", false, t0)
	js := []run.Journal{journal("A", true, "auto"), journal("B", true, "busy"), journal("O", false, "owner")}
	var mem []daemon.WatchdogRun
	var holds []daemon.WatchdogHold
	poll := func(now time.Time) {
		t.Helper()
		var err error
		if mem, holds, err = watchdogStalled(pol, js, now, mem, d); err != nil {
			t.Fatal(err)
		}
	}

	poll(t0.Add(19 * time.Minute))
	if len(sent)+len(stopped) != 0 {
		t.Fatalf("before the first limit: sent %v stopped %v", sent, stopped)
	}
	poll(t0.Add(21 * time.Minute))
	if strings.Join(sent, ",") != "auto,busy" || len(stopped) != 0 {
		t.Fatalf("first limit: sent %v stopped %v", sent, stopped)
	}
	// busy makes progress in between; auto and the owner's run stay silent.
	report("busy", true, t0.Add(30*time.Minute))
	poll(t0.Add(41 * time.Minute))
	if strings.Join(stopped, ",") != "A" || strings.Join(resumed, ",") != "A" {
		t.Fatalf("second limit: stopped %v resumed %v", stopped, resumed)
	}
	if len(holds) != 1 || holds[0].Run != "O" {
		t.Fatalf("holds %+v, want the attended run's", holds)
	}
	var lines int
	for _, l := range log {
		if strings.Contains(l, "run A") && strings.Contains(l, "stopped and resumed") {
			lines++
		}
	}
	if lines != 1 {
		t.Errorf("want one stop line for A, log %q", log)
	}
	for _, m := range mem {
		if m.Run == "A" && m.Spent != 1 {
			t.Errorf("A spent %d, want 1", m.Spent)
		}
		if m.Run == "B" && m.Spent != 0 {
			t.Errorf("B spent %d, want 0", m.Spent)
		}
	}
	// The budget is spent: the next stall of A is held, not retried. B, silent
	// again since its progress, is retried once, within its own budget.
	poll(t0.Add(90 * time.Minute))
	if strings.Join(stopped, ",") != "A,B" {
		t.Fatalf("over budget: stopped %v", stopped)
	}
	if !slicesHoldFor(holds, "A") || !slicesHoldFor(holds, "O") {
		t.Fatalf("holds %+v, want A and O", holds)
	}
	for _, id := range []string{"O"} {
		for _, s := range append(append([]string{}, stopped...), resumed...) {
			if s == id {
				t.Errorf("attended run %s was stopped or resumed", id)
			}
		}
	}
	if items := watchdogHoldItems(holds); len(items) != len(holds) || items[0].Kind != InboxStalled {
		t.Errorf("hold items %+v", items)
	}
}

func slicesHoldFor(holds []daemon.WatchdogHold, id string) bool {
	for _, h := range holds {
		if h.Run == id {
			return true
		}
	}
	return false
}

// TestWatchdogTransient drives the transient-failure recovery: each known
// signature leads to one resume --retry, an unknown failure is left alone,
// and so is a run that is live.
func TestWatchdogTransient(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	failed := func(id, out string, auto bool) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: auto}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\n" + out},
			{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
		}}
	}
	js := []run.Journal{
		failed("work", "open /tmp/go-build123/b001/importcfg: no such file or directory", true),
		failed("lock", "fatal: Unable to create '/r/.git/index.lock': File exists.", true),
		failed("gh", "gh: HTTP 502 Bad Gateway", true),
		failed("real", "FAIL: TestThing expected 1 got 2", true),
		failed("live", "gh: HTTP 503 Service Unavailable", true),
		failed("owner", "gh: HTTP 502 Bad Gateway", false),
	}
	var log, resumed []string
	d := watchdogDeps{
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	live := func(id string) bool { return id == "live" }
	mem, err := watchdogTransient(pol, js, live, nil, d)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(resumed, ","); got != "work,lock,gh" {
		t.Fatalf("resumed %s, want work,lock,gh", got)
	}
	if len(log) != 3 {
		t.Errorf("want one line per resume, log %q", log)
	}
	// The budget is one: a second poll over the same failures acts on none.
	if _, err := watchdogTransient(pol, js, live, mem, d); err != nil {
		t.Fatal(err)
	}
	if len(resumed) != 3 {
		t.Errorf("over budget resumed %v", resumed)
	}
}

// TestWatchdogCIRerun drives the CI rerun: red checks on a head are rerun
// once, the run is resumed once they pass, and a head already at attempt 2,
// or already rerun, gets no second rerun, while a new head gets its own.
func TestWatchdogCIRerun(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	const def = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\n"
	failed := func(id string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
			{Type: run.TypeEntered, At: t0, Step: "wait", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "wait", Visit: 1, Outcome: "fail", Detail: "failed: test"},
			{Type: run.TypeFinished, At: t0, Step: "wait", End: "failed", Detail: "step wait failed"},
		}}
	}
	red := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "FAILURE", RunID: "7"}}, Attempts: map[string]int{"7": 1}}
	green := watchdogPR{Head: "aaaa111", Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS", RunID: "7"}}, Attempts: map[string]int{"7": 2}}
	pr := red
	var reran, resumed, log []string
	d := watchdogDeps{
		ci:     func(run.State) (watchdogPR, error) { return pr, nil },
		rerun:  func(_ run.State, id string) error { reran = append(reran, id); return nil },
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 2, CIReruns: 1, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	js := []run.Journal{failed("r")}
	poll := func(mem []daemon.WatchdogRun) []daemon.WatchdogRun {
		t.Helper()
		mem, err := watchdogCIRerun(pol, js, nolive, mem, d)
		if err != nil {
			t.Fatal(err)
		}
		return mem
	}

	mem := poll(nil)
	if strings.Join(reran, ",") != "7" || len(resumed) != 0 {
		t.Fatalf("red head: reran %v resumed %v", reran, resumed)
	}
	// Still red after the rerun (attempt 2): no second rerun, no resume.
	pr = red
	pr.Attempts = map[string]int{"7": 2}
	mem = poll(mem)
	if len(reran) != 1 || len(resumed) != 0 {
		t.Fatalf("red again: reran %v resumed %v", reran, resumed)
	}
	// The rerun passed: the run is resumed once.
	pr = green
	mem = poll(mem)
	mem = poll(mem)
	if strings.Join(resumed, ",") != "r" {
		t.Fatalf("green: resumed %v, want one", resumed)
	}
	if len(log) != 2 {
		t.Errorf("want one line per action, log %q", log)
	}
	// A head whose run is already at attempt 2 (the plan step's rerun) is not
	// rerun again.
	reran = nil
	pr = red
	pr.Attempts = map[string]int{"7": 2}
	poll(nil)
	if len(reran) != 0 {
		t.Errorf("attempt above 1 reran %v", reran)
	}
	// A new head is allowed its own rerun.
	pr = red
	pr.Head = "bbbb222"
	poll(mem)
	if strings.Join(reran, ",") != "7" {
		t.Errorf("new head reran %v, want 7", reran)
	}
}

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

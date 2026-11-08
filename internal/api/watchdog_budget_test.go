package api

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogBudget runs polls over a run that keeps failing on a transient
// error, restarting the daemon between them (the memory goes through its
// saved JSON): each attempt is journalled and counted across restarts, the
// watchdog acts no more once the budget is spent, and the owner gets one
// inbox item naming every attempt.
func TestWatchdogBudget(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	js := []run.Journal{{ID: "gh", Entries: []run.Entry{
		{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true}},
		{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
		{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\ngh: HTTP 502 Bad Gateway"},
		{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: "step build failed"},
	}}}
	now := t0.Add(time.Hour)
	var resumed []string
	d := watchdogDeps{
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		now:    func() time.Time { return now },
		log:    func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 2, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	var mem []daemon.WatchdogRun
	var holds []daemon.WatchdogHold
	poll := func() {
		t.Helper()
		var err error
		if mem, holds, err = watchdogPasses(pol, js, nolive, now, mem, d); err != nil {
			t.Fatal(err)
		}
		// A daemon restart: the memory is read back from its saved state.
		b, err := json.Marshal(daemon.State{Watchdog: mem, WatchdogHolds: holds})
		if err != nil {
			t.Fatal(err)
		}
		var st daemon.State
		if err := json.Unmarshal(b, &st); err != nil {
			t.Fatal(err)
		}
		mem, holds = st.Watchdog, st.WatchdogHolds
		now = now.Add(time.Minute)
	}
	poll()
	if len(holds) != 0 || len(mem) != 1 || len(mem[0].Attempts) != 1 {
		t.Fatalf("after one attempt: holds %+v memory %+v", holds, mem)
	}
	for range 4 {
		poll()
	}
	if len(resumed) != 2 {
		t.Fatalf("resumed %v, want two attempts and no more", resumed)
	}
	m := mem[0]
	if m.Spent != 2 || len(m.Attempts) != 2 || !m.Attempts[0].At.Equal(t0.Add(time.Hour)) || !m.Attempts[1].At.Equal(t0.Add(time.Hour+time.Minute)) {
		t.Fatalf("attempts %+v spent %d", m.Attempts, m.Spent)
	}
	items := watchdogHoldItems(holds)
	if len(items) != 1 {
		t.Fatalf("want one inbox item, got %+v", items)
	}
	for _, want := range []string{"budget of 2 is spent", "502", "2030-01-01T01:00:00Z", "2030-01-01T01:01:00Z"} {
		if !strings.Contains(items[0].Detail, want) {
			t.Errorf("item %q lacks %q", items[0].Detail, want)
		}
	}
	id := items[0].ID
	poll()
	if again := watchdogHoldItems(holds); len(again) != 1 || again[0].ID != id {
		t.Errorf("item changed across polls: %+v", again)
	}
}

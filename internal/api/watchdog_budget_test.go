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

// TestWatchdogHolds: a run with its budget spent keeps its one hold, and
// loses it once abandoned or once its change set is withdrawn or rejected.
func TestWatchdogHolds(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	journal := func(id, change, detail string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Change: change}},
			{Type: run.TypeEntered, At: t0, Step: "build", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "build", Visit: 1, Outcome: "fail", Detail: "exit 1\ngh: HTTP 502 Bad Gateway"},
			{Type: run.TypeFinished, At: t0, Step: "build", End: "failed", Detail: detail},
		}}
	}
	states := map[string]string{"c-ok": "accepted", "c-aband": "accepted", "c-with": "withdrawn", "c-rej": "rejected"}
	js := []run.Journal{
		journal("ok", "c-ok", "step build failed"),
		journal("aband", "c-aband", run.AbandonedPrefix+"me"),
		journal("with", "c-with", "step build failed"),
		journal("rej", "c-rej", "step build failed"),
	}
	spent := func(ids ...string) []daemon.WatchdogRun {
		var mem []daemon.WatchdogRun
		for _, id := range ids {
			mem = append(mem, daemon.WatchdogRun{Run: id, Spent: 2, Attempts: []daemon.WatchdogAttempt{
				{At: t0, Action: "a"}, {At: t0.Add(time.Minute), Action: "b"}}})
		}
		return mem
	}
	d := watchdogDeps{
		changeState: func(c string) (string, bool) { s, ok := states[c]; return s, ok },
		now:         func() time.Time { return t0.Add(time.Hour) },
		log:         func(string, ...any) {},
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 2, AutonomousOnly: true}
	// A hold saved for each run before; the next poll drops all but one.
	prev := spent("ok", "aband", "with", "rej")
	_, holds, err := watchdogPasses(pol, js, func(string) bool { return false }, t0.Add(time.Hour), prev, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(holds) != 1 || holds[0].Run != "ok" {
		t.Fatalf("holds %+v, want only run ok's", holds)
	}
	if items := watchdogHoldItems(holds); len(items) != 1 {
		t.Fatalf("items %+v", items)
	}
}

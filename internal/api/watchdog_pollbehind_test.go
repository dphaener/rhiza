package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// TestWatchdogPollBehind drives the poll-behind check with a fake clock: a
// slow poll is logged once and warns in health, a poll within the limit
// clears it.
func TestWatchdogPollBehind(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	pol := WatchdogPolicy{Enabled: true, PollBehind: 3}
	var state daemon.State
	var log []string
	logf := func(l string) { log = append(log, l) }
	health := func() HealthPart {
		behind := ""
		if state.PollBehind != nil {
			behind = pollBehindText(state.PollBehind)
		}
		return healthPollBehind(HealthPart{Part: HealthDaemon, State: HealthOK, Summary: "running"}, behind)
	}

	watchdogPollBehind(&state, pol, 5*time.Second, 15*time.Second, t0, nil, logf)
	if len(log) != 0 || health().State != HealthOK {
		t.Fatalf("a poll at the limit: log %q, health %+v", log, health())
	}

	watchdogPollBehind(&state, pol, 5*time.Second, 40*time.Second, t0.Add(time.Minute), nil, logf)
	if len(log) != 1 || !strings.Contains(log[0], "40s") || !strings.Contains(log[0], "15s") {
		t.Fatalf("slow poll log %q, want one line naming 40s and 15s", log)
	}
	if h := health(); h.State != HealthWarn || !strings.Contains(h.Summary, "polling behind") {
		t.Fatalf("health after slow poll %+v, want warn", h)
	}

	watchdogPollBehind(&state, pol, 5*time.Second, 2*time.Second, t0.Add(2*time.Minute), nil, logf)
	if len(log) != 1 || health().State != HealthOK || state.PollBehind != nil {
		t.Fatalf("after a poll within the limit: log %q, health %+v", log, health())
	}

	pol.Enabled = false
	watchdogPollBehind(&state, pol, 5*time.Second, time.Hour, t0, nil, logf)
	if len(log) != 1 || state.PollBehind != nil {
		t.Fatalf("a disabled watchdog logged or recorded: %q %+v", log, state.PollBehind)
	}
}

// TestWatchdogPollBehindSteps times named steps with a fake clock and checks
// the logged line and the recorded watchdog action name the slowest ones.
func TestWatchdogPollBehindSteps(t *testing.T) {
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	now := t0
	ps := newPollSteps(func() time.Time { return now })
	step := func(project, name string, d time.Duration) {
		ps.run(project, name, func() { now = now.Add(d) })
	}
	step("", "registry load", 10*time.Millisecond)
	step("rhiza", "demand", 21*time.Second)
	step("rhiza", "watchdog", 5*time.Minute+58*time.Second)
	step("other", "settings", 2*time.Second)
	step("other", "events", 3*time.Second)
	step("rhiza", "notify", 0)

	pol := WatchdogPolicy{Enabled: true, PollBehind: 3}
	var state daemon.State
	var log []string
	took := now.Sub(t0)
	watchdogPollBehind(&state, pol, 5*time.Second, took, now, ps.all(), func(l string) {
		log = append(log, l)
		state.WatchdogActions = watchdogRecord(state.WatchdogActions, now, l)
	})
	want := "watchdog: poll took 6m24.01s, over its limit of 15s (3 times the 5s poll interval); slowest: rhiza watchdog 5m58s, rhiza demand 21s, other events 3s"
	if len(log) != 1 || log[0] != want {
		t.Fatalf("log %q, want %q", log, want)
	}
	if len(state.WatchdogActions) != 1 || !strings.Contains(state.WatchdogActions[0].Line, "slowest: rhiza watchdog 5m58s") {
		t.Fatalf("recorded actions %+v", state.WatchdogActions)
	}

	log = nil
	watchdogPollBehind(&state, pol, 5*time.Second, 15*time.Second, now, ps.all(), func(l string) { log = append(log, l) })
	if len(log) != 0 {
		t.Fatalf("a poll within the limit logged %q", log)
	}
}

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

	watchdogPollBehind(&state, pol, 5*time.Second, 15*time.Second, t0, logf)
	if len(log) != 0 || health().State != HealthOK {
		t.Fatalf("a poll at the limit: log %q, health %+v", log, health())
	}

	watchdogPollBehind(&state, pol, 5*time.Second, 40*time.Second, t0.Add(time.Minute), logf)
	if len(log) != 1 || !strings.Contains(log[0], "40s") || !strings.Contains(log[0], "15s") {
		t.Fatalf("slow poll log %q, want one line naming 40s and 15s", log)
	}
	if h := health(); h.State != HealthWarn || !strings.Contains(h.Summary, "polling behind") {
		t.Fatalf("health after slow poll %+v, want warn", h)
	}

	watchdogPollBehind(&state, pol, 5*time.Second, 2*time.Second, t0.Add(2*time.Minute), logf)
	if len(log) != 1 || health().State != HealthOK || state.PollBehind != nil {
		t.Fatalf("after a poll within the limit: log %q, health %+v", log, health())
	}

	pol.Enabled = false
	watchdogPollBehind(&state, pol, 5*time.Second, time.Hour, t0, logf)
	if len(log) != 1 || state.PollBehind != nil {
		t.Fatalf("a disabled watchdog logged or recorded: %q %+v", log, state.PollBehind)
	}
}

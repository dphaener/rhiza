package api

import (
	"slices"
	"strings"
	"sync"
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

// TestWatchdogPollOverrun runs polls against a shortened limit: one whose
// step blocks past it logs one line naming the step and project before the
// step is released, a poll within the limit and a disabled check log none.
func TestWatchdogPollOverrun(t *testing.T) {
	const interval = 10 * time.Millisecond
	var mu sync.Mutex
	var lines []string
	logf := func(l string) { mu.Lock(); lines = append(lines, l); mu.Unlock() }
	got := func() []string { mu.Lock(); defer mu.Unlock(); return slices.Clone(lines) }
	reset := func() { mu.Lock(); lines = nil; mu.Unlock() }

	pol := WatchdogPolicy{Enabled: true, PollBehind: 2}
	ps := newPollSteps(time.Now)
	stop, _ := watchPollOverrun(ps, pol, interval, logf)
	release := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		ps.run("rhiza", "watchdog", func() { <-release })
		close(finished)
	}()
	deadline := time.Now().Add(10 * time.Second)
	for len(got()) == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	// The line is in before the step is released.
	if l := got(); len(l) != 1 || !strings.Contains(l[0], "in rhiza watchdog") || !strings.Contains(l[0], "after 20ms") {
		t.Fatalf("lines %q, want one naming rhiza watchdog and 20ms", l)
	}
	time.Sleep(60 * time.Millisecond)
	close(release)
	<-finished
	stop()
	if l := got(); len(l) != 1 {
		t.Fatalf("lines %q, want exactly one", l)
	}

	reset()
	stop, _ = watchPollOverrun(newPollSteps(time.Now), WatchdogPolicy{Enabled: true, PollBehind: 100}, interval, logf)
	time.Sleep(30 * time.Millisecond)
	stop()
	time.Sleep(10 * time.Millisecond)
	if l := got(); len(l) != 0 {
		t.Fatalf("a poll within the limit logged %q", l)
	}

	for _, off := range []WatchdogPolicy{{Enabled: false, PollBehind: 2}, {Enabled: true, PollBehind: 0}} {
		stop, _ = watchPollOverrun(newPollSteps(time.Now), off, interval, logf)
		time.Sleep(40 * time.Millisecond)
		stop()
		if l := got(); len(l) != 0 {
			t.Fatalf("policy %+v logged %q", off, l)
		}
	}
}

// TestWatchdogPollOverrunFirstPoll covers the first poll of a daemon image:
// it starts under the shipped defaults, which watch it, and a policy read
// mid-poll re-arms the watch, timed from the poll's start.
func TestWatchdogPollOverrunFirstPoll(t *testing.T) {
	if def := defaultWatchdogPolicy(); !def.Enabled || def.PollBehind != DefaultPollBehind {
		t.Fatalf("default policy %+v, want enabled with poll-behind %d", def, DefaultPollBehind)
	}
	const interval = 10 * time.Millisecond
	var mu sync.Mutex
	var lines []string
	logf := func(l string) { mu.Lock(); lines = append(lines, l); mu.Unlock() }
	got := func() []string { mu.Lock(); defer mu.Unlock(); return slices.Clone(lines) }

	// Armed off, then re-armed on: the line comes at the new limit.
	stop, rearm := watchPollOverrun(newPollSteps(time.Now), WatchdogPolicy{}, interval, logf)
	rearm(WatchdogPolicy{Enabled: true, PollBehind: 2})
	deadline := time.Now().Add(10 * time.Second)
	for len(got()) == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(30 * time.Millisecond)
	stop()
	if l := got(); len(l) != 1 || !strings.Contains(l[0], "after 20ms") {
		t.Fatalf("lines %q, want one after 20ms", l)
	}

	// Armed on, then re-armed off: nothing is logged.
	mu.Lock()
	lines = nil
	mu.Unlock()
	stop, rearm = watchPollOverrun(newPollSteps(time.Now), WatchdogPolicy{Enabled: true, PollBehind: 2}, interval, logf)
	rearm(WatchdogPolicy{Enabled: false, PollBehind: 2})
	time.Sleep(40 * time.Millisecond)
	stop()
	rearm(WatchdogPolicy{Enabled: true, PollBehind: 1})
	time.Sleep(20 * time.Millisecond)
	if l := got(); len(l) != 0 {
		t.Fatalf("a watch re-armed off, or stopped, logged %q", l)
	}
}

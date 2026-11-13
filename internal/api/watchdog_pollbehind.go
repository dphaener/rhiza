package api

// This file is the daemon watchdog's own-poll check (decision
// watchdog-scope-and-defaults): a poll that takes longer than pollBehind
// times the poll interval is logged, one line per slow poll, and kept in the
// daemon state so health warns until a poll is within the limit again.

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// pollStep is how long one named step of a poll took, in one project (empty
// for a step of the poll as a whole).
type pollStep struct {
	Project string
	Name    string
	Took    time.Duration
}

// pollSteps times the steps of one poll and knows which are running. A nil
// *pollSteps times nothing. It is safe for the poll and the overrun watcher
// to use at once.
type pollSteps struct {
	now     func() time.Time
	mu      sync.Mutex
	steps   []pollStep
	running []pollStep // started steps, innermost last; Took unused
}

func newPollSteps(now func() time.Time) *pollSteps { return &pollSteps{now: now} }

// run runs f and records how long it took as the step name of project.
func (ps *pollSteps) run(project, name string, f func()) {
	if ps == nil {
		f()
		return
	}
	start := ps.now()
	ps.mu.Lock()
	ps.running = append(ps.running, pollStep{Project: project, Name: name})
	ps.mu.Unlock()
	defer func() {
		took := ps.now().Sub(start)
		ps.mu.Lock()
		ps.steps = append(ps.steps, pollStep{Project: project, Name: name, Took: took})
		ps.running = ps.running[:len(ps.running)-1]
		ps.mu.Unlock()
	}()
	f()
}

// all is the recorded steps.
func (ps *pollSteps) all() []pollStep {
	if ps == nil {
		return nil
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return slices.Clone(ps.steps)
}

// current is the innermost step running now, and false when none is.
func (ps *pollSteps) current() (pollStep, bool) {
	if ps == nil {
		return pollStep{}, false
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if len(ps.running) == 0 {
		return pollStep{}, false
	}
	return ps.running[len(ps.running)-1], true
}

// watchPollOverrun starts the live check of a poll: when the poll is still
// running after pol.PollBehind times interval, it calls log once with a line
// naming the step and project it is in. The returned stop ends the poll's
// watch; after it returns nothing is logged. The returned rearm swaps in the
// policy the poll learns once its settings are read: the limit stays timed
// from the watch's start, and a line already logged is not logged again. A
// disabled watchdog or poll-behind check watches nothing.
func watchPollOverrun(ps *pollSteps, pol WatchdogPolicy, interval time.Duration, log func(line string)) (stop func(), rearm func(WatchdogPolicy)) {
	if log == nil {
		return func() {}, func(WatchdogPolicy) {}
	}
	start := time.Now()
	var mu sync.Mutex
	var t *time.Timer
	done, logged := false, false
	arm := func(pol WatchdogPolicy) {
		if t != nil {
			t.Stop()
			t = nil
		}
		if done || logged || !pol.Enabled || pol.PollBehind < 1 {
			return
		}
		limit := interval * time.Duration(pol.PollBehind)
		var self *time.Timer
		self = time.AfterFunc(limit-time.Since(start), func() {
			mu.Lock()
			defer mu.Unlock()
			if done || logged || t != self {
				return
			}
			logged = true
			line := fmt.Sprintf("watchdog: poll still running after %s", limit)
			if st, ok := ps.current(); ok {
				where := st.Name
				if st.Project != "" {
					where = st.Project + " " + st.Name
				}
				line += ", in " + where
			}
			log(line)
		})
		t = self
	}
	mu.Lock()
	arm(pol)
	mu.Unlock()
	stop = func() {
		mu.Lock()
		defer mu.Unlock()
		done = true
		if t != nil {
			t.Stop()
		}
	}
	rearm = func(pol WatchdogPolicy) {
		mu.Lock()
		defer mu.Unlock()
		arm(pol)
	}
	return stop, rearm
}

// pollSlowestShown is how many steps the poll-behind line names.
const pollSlowestShown = 3

// slowestText names the slowest steps, slowest first, as "rhiza watchdog
// 5m58s, rhiza demand 21s": steps of the same project and name are summed,
// and the empty string is returned when none took any time.
func slowestText(steps []pollStep) string {
	var sums []pollStep
	for _, st := range steps {
		if i := slices.IndexFunc(sums, func(o pollStep) bool { return o.Project == st.Project && o.Name == st.Name }); i >= 0 {
			sums[i].Took += st.Took
		} else {
			sums = append(sums, st)
		}
	}
	sums = slices.DeleteFunc(sums, func(st pollStep) bool { return st.Took <= 0 })
	slices.SortStableFunc(sums, func(a, b pollStep) int { return int(b.Took - a.Took) })
	if len(sums) > pollSlowestShown {
		sums = sums[:pollSlowestShown]
	}
	parts := make([]string, 0, len(sums))
	for _, st := range sums {
		label := st.Name
		if st.Project != "" {
			label = st.Project + " " + st.Name
		}
		d := st.Took.Round(time.Second)
		if st.Took < time.Second {
			d = st.Took.Round(time.Millisecond)
		}
		parts = append(parts, label+" "+d.String())
	}
	return strings.Join(parts, ", ")
}

// watchdogPollBehind records on state whether a poll that took took ran past
// pol.PollBehind times interval, calling log with one line when it did. A
// poll within the limit clears the record. A disabled watchdog records none.
// The line ends with the slowest of steps, when any took time.
func watchdogPollBehind(state *daemon.State, pol WatchdogPolicy, interval, took time.Duration, now time.Time, steps []pollStep, log func(line string)) {
	limit := interval * time.Duration(pol.PollBehind)
	if !pol.Enabled || pol.PollBehind < 1 || took <= limit {
		state.PollBehind = nil
		return
	}
	state.PollBehind = &daemon.PollBehind{At: now.UTC(), Took: took, Limit: limit}
	if log != nil {
		line := fmt.Sprintf("watchdog: poll took %s, over its limit of %s (%d times the %s poll interval)", took.Round(time.Millisecond), limit, pol.PollBehind, interval)
		if slow := slowestText(steps); slow != "" {
			line += "; slowest: " + slow
		}
		log(line)
	}
}

// pollBehindText is the status line for a slow poll.
func pollBehindText(pb *daemon.PollBehind) string {
	return fmt.Sprintf("last poll took %s, over the limit of %s", pb.Took.Round(time.Millisecond), pb.Limit)
}

// healthPollBehind turns the daemon part to warn while behind (the status
// line of a slow poll) is set.
func healthPollBehind(p HealthPart, behind string) HealthPart {
	if behind == "" {
		return p
	}
	p.State = HealthWarn
	p.Summary += "; polling behind: " + behind
	p.Detail = append(p.Detail, "poll behind: "+behind)
	return p
}

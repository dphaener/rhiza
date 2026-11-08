package api

// This file is the daemon watchdog's own-poll check (decision
// watchdog-scope-and-defaults): a poll that takes longer than pollBehind
// times the poll interval is logged, one line per slow poll, and kept in the
// daemon state so health warns until a poll is within the limit again.

import (
	"fmt"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// watchdogPollBehind records on state whether a poll that took took ran past
// pol.PollBehind times interval, calling log with one line when it did. A
// poll within the limit clears the record. A disabled watchdog records none.
func watchdogPollBehind(state *daemon.State, pol WatchdogPolicy, interval, took time.Duration, now time.Time, log func(line string)) {
	limit := interval * time.Duration(pol.PollBehind)
	if !pol.Enabled || pol.PollBehind < 1 || took <= limit {
		state.PollBehind = nil
		return
	}
	state.PollBehind = &daemon.PollBehind{At: now.UTC(), Took: took, Limit: limit}
	if log != nil {
		log(fmt.Sprintf("watchdog: poll took %s, over its limit of %s (%d times the %s poll interval)", took.Round(time.Millisecond), limit, pol.PollBehind, interval))
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

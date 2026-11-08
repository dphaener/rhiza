package api

// This file is the watchdog's action journal (decision
// watchdog-scope-and-defaults): each line the watchdog logs, saying what it
// saw, what it did and the run or request, is also kept in the daemon state,
// newest first and bounded, for daemon status and rhiza health.

import (
	"fmt"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// watchdogActionsKept bounds the journal.
const watchdogActionsKept = 20

// watchdogRecord puts line, stamped at, first in actions. A line equal to
// one already kept replaces it rather than repeating: a hold is logged every
// poll and is one action.
func watchdogRecord(actions []daemon.WatchdogAction, at time.Time, line string) []daemon.WatchdogAction {
	line = strings.TrimPrefix(line, "watchdog: ")
	out := []daemon.WatchdogAction{{At: at.UTC(), Line: line}}
	for _, a := range actions {
		if a.Line != line && len(out) < watchdogActionsKept {
			out = append(out, a)
		}
	}
	return out
}

// watchdogJournalLog wraps log so each line it logs is also journalled on
// state.
func watchdogJournalLog(log func(format string, args ...any), state *daemon.State, now func() time.Time) func(format string, args ...any) {
	return func(format string, args ...any) {
		log(format, args...)
		var at time.Time
		if now != nil {
			at = now()
		}
		state.WatchdogActions = watchdogRecord(state.WatchdogActions, at, fmt.Sprintf(format, args...))
	}
}

// watchdogActionLines are the journal as one line each: time, then what.
func watchdogActionLines(actions []daemon.WatchdogAction) []string {
	var out []string
	for _, a := range actions {
		out = append(out, a.At.UTC().Format(time.RFC3339)+" "+a.Line)
	}
	return out
}

// healthWatchdogActions lists the watchdog's last actions, newest first, as
// detail of the daemon part.
func healthWatchdogActions(p HealthPart, lines []string) HealthPart {
	for i, l := range lines {
		if i == 0 {
			p.Detail = append(p.Detail, "watchdog's last actions, newest first:")
		}
		p.Detail = append(p.Detail, "  "+l)
	}
	return p
}

package api

// This file is the watchdog's runaway pass (decision
// runaway-fails-not-retries): a live run's stall session whose process tree
// holds more than runawayProcs processes or runawayMemory of resident memory
// is nudged once with the reading; when a poll at least runawayGrace after
// the nudge still reads it over a threshold, the run is stopped (the stall
// sessions' panes closed and the parents waiting on their spawns, the run's
// process terminated) and its current step journalled stopped with the
// reading, never resumed or retried: a retry would rerun the same runaway.
// A reading back under every threshold clears the nudge. An attended run
// under autonomousOnly gets a log line only.

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/proctree"
	"github.com/dphaener/rhiza/internal/run"
)

// runawayReading is one stall session's tree read over a threshold.
type runawayReading struct {
	session string
	tree    proctree.Tree
	// over names each threshold the tree passed, e.g. "runawayProcs 400".
	over []string
}

// detail is the reading as the stopped entry, the bug and the log name it:
// the session with its role, the process count, the resident memory and
// each threshold passed.
func (r runawayReading) detail(st run.State) string {
	return fmt.Sprintf("session %s (%s): %d processes, %s resident memory, over %s",
		r.session, sessionRole(st, r.session), r.tree.Count, daemon.FormatSize(uint64(r.tree.RSS)), strings.Join(r.over, " and "))
}

// runawayOver names each threshold of pol that t passes; none is under
// both. A threshold of 0 is off and never passed.
func runawayOver(pol WatchdogPolicy, t proctree.Tree) []string {
	var over []string
	if pol.RunawayProcs > 0 && t.Count > pol.RunawayProcs {
		over = append(over, fmt.Sprintf("runawayProcs %d", pol.RunawayProcs))
	}
	if pol.RunawayMemory > 0 && uint64(t.RSS) > pol.RunawayMemory {
		over = append(over, "runawayMemory "+daemon.FormatSize(pol.RunawayMemory))
	}
	return over
}

// runawayNudgeText is the nudge a runaway session gets: its reading, the
// thresholds it passed and what happens after the grace.
func runawayNudgeText(r runawayReading, grace time.Duration) string {
	return fmt.Sprintf("The Rhiza watchdog reads your process tree at %d processes and %s resident memory, over its threshold %s. "+
		"Stop the processes you started and no longer need (servers, watchers, test runs left behind) and start no more in parallel. "+
		"If your tree is still over in %s, the watchdog stops this run.",
		r.tree.Count, daemon.FormatSize(uint64(r.tree.RSS)), strings.Join(r.over, " and "), grace)
}

// watchdogRunaway checks the stall sessions of each live run of journals
// against pol's runaway thresholds, reading the process table once through
// d.procs. It returns the memory after (prev with each run's runaway
// nudges updated), the runs it stopped, which the stalled-step pass leaves
// be this poll, and the first error. A failed nudge or stop is not
// recorded, so the next poll tries again. With both thresholds off, no
// table reader or no stall session, it reads nothing.
func watchdogRunaway(pol WatchdogPolicy, journals []run.Journal, now time.Time, prev []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, []string, error) {
	if d.procs == nil || pol.RunawayProcs <= 0 && pol.RunawayMemory == 0 {
		return prev, nil, nil
	}
	if !slices.ContainsFunc(journals, func(j run.Journal) bool { return len(stallSessions(j.State())) > 0 }) {
		return prev, nil, nil
	}
	table, err := d.procs()
	if err != nil {
		return prev, nil, fmt.Errorf("read the process table: %w", err)
	}
	mem := slices.Clone(prev)
	var stopped []string
	var first error
	fail := func(err error) {
		if first == nil {
			first = err
		}
	}
	for _, j := range journals {
		st := j.State()
		sids := stallSessions(st)
		var readings []runawayReading
		for _, sid := range sids {
			t, err := proctree.SessionTree(table, sid)
			if err != nil {
				continue
			}
			if over := runawayOver(pol, t); len(over) > 0 {
				readings = append(readings, runawayReading{session: sid, tree: t, over: over})
			}
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		var nudged []daemon.RunawayNudge
		if i >= 0 {
			nudged = mem[i].Runaway
		}
		// Only the sessions still over keep their nudge.
		var keep []daemon.RunawayNudge
		for _, n := range nudged {
			if slices.ContainsFunc(readings, func(r runawayReading) bool { return r.session == n.Session }) {
				keep = append(keep, n)
			}
		}
		set := func(ns []daemon.RunawayNudge) {
			if i < 0 {
				if len(ns) == 0 {
					return
				}
				mem = append(mem, daemon.WatchdogRun{Run: j.ID})
				i = len(mem) - 1
			}
			mem[i].Runaway = ns
		}
		if len(readings) == 0 {
			set(nil)
			continue
		}
		var parts []string
		for _, r := range readings {
			parts = append(parts, r.detail(st))
		}
		reading := strings.Join(parts, "; ")
		if !st.Autonomous && pol.AutonomousOnly {
			set(nil)
			d.log("watchdog: run %s step %s runaway: %s: attended, left to the owner", j.ID, st.Current, reading)
			continue
		}
		due := slices.ContainsFunc(keep, func(n daemon.RunawayNudge) bool { return now.Sub(n.At) >= pol.RunawayGrace })
		if !due {
			for _, r := range readings {
				if slices.ContainsFunc(keep, func(n daemon.RunawayNudge) bool { return n.Session == r.session }) {
					continue
				}
				if err := d.send(r.session, runawayNudgeText(r, pol.RunawayGrace)); err != nil {
					fail(fmt.Errorf("nudge session %s: %w", r.session, err))
					continue
				}
				keep = append(keep, daemon.RunawayNudge{Session: r.session, At: now})
				d.log("watchdog: run %s step %s runaway: %s: nudged", j.ID, st.Current, r.detail(st))
			}
			set(keep)
			continue
		}
		set(keep)
		closing := slices.Clone(sids)
		for _, sid := range spawnWaiting(st) {
			if !slices.Contains(closing, sid) {
				closing = append(closing, sid)
			}
		}
		if err := d.stop(st, closing); err != nil {
			fail(fmt.Errorf("stop run %s: %w", j.ID, err))
			continue
		}
		set(nil)
		stopped = append(stopped, j.ID)
		detail := daemon.RunawayStopped + reading
		if d.stopped != nil {
			if err := d.stopped(j.ID, st.Current, detail); err != nil {
				fail(fmt.Errorf("journal run %s stopped: %w", j.ID, err))
			}
		}
		d.fileBug(pol, j.ID, bugRunaway, st.Current, reading)
		d.log("watchdog: run %s step %s runaway: %s: stopped, not resumed", j.ID, st.Current, reading)
	}
	return mem, stopped, first
}

// watchdogProcTable is the daemon's process-table reader: one ps call
// through internal/proctree. Under RHIZA_TEST_NOW it is nil, so no test
// daemon reads the host's processes.
func watchdogProcTable(ctx *Ctx) func() ([]proctree.Process, error) {
	if ctxGetenv(ctx)("RHIZA_TEST_NOW") != "" {
		return nil
	}
	return proctree.ReadTable
}

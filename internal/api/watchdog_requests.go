package api

// This file is the watchdog's handling of plan requests (decision
// watchdog-scope-and-defaults): a request whose plan run is gone while the
// request still says planning, or whose plan run failed on a transient
// signature, is recovered like any run, by resume --retry within the run's
// budget. A request that failed for any other reason is left for the owner.
// A stalled live plan run is the stalled-step pass's. A milestone request
// always needs the owner and is never touched.

import (
	"fmt"
	"slices"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// watchdogRequests recovers the plan runs of reqs. A plan run the daemon
// started for a request is the daemon's own, so it is recovered whether or
// not it is autonomous; the failed-transient pass already takes the ones it
// would, and those are skipped here. It returns the memory after, and the
// first error.
func watchdogRequests(pol WatchdogPolicy, reqs []daemon.Request, journals []run.Journal, live func(id string) bool, mem []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, error) {
	var first error
	fail := func(err error) {
		if first == nil {
			first = err
		}
	}
	for _, r := range reqs {
		if r.Run == "" || r.Kind == daemon.KindMilestone {
			continue
		}
		ji := slices.IndexFunc(journals, func(j run.Journal) bool { return j.ID == r.Run })
		if ji < 0 || len(journals[ji].Entries) == 0 {
			continue
		}
		st := journals[ji].State()
		if live(r.Run) || st.Abandoned() {
			continue
		}
		var why string
		switch {
		case !st.Ended:
			why = "its plan run is gone while the request still says planning"
		case st.Result == run.Failed:
			if st.Autonomous || !pol.AutonomousOnly {
				continue // the failed-transient pass's
			}
			step, err := st.FailedStep()
			if err != nil {
				continue
			}
			sig := transientSignature(failedDetail(st, step))
			if sig == "" {
				continue
			}
			why = fmt.Sprintf("its plan run failed at step %s on a transient error (%q)", step, sig)
		default:
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == r.Run })
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: r.Run})
			i = len(mem) - 1
		}
		if mem[i].Spent >= pol.RetryBudget || d.queued != nil && d.queued(r.Run) {
			continue
		}
		if err := d.resume(r.Run); err != nil {
			fail(fmt.Errorf("resume plan run %s: %w", r.Run, err))
			continue
		}
		d.charge(&mem[i], fmt.Sprintf("request %s: %s: resumed with --retry", r.ID, why))
		d.log("watchdog: request %s: %s: resumed run %s with --retry (budget %d of %d)", r.ID, why, r.Run, mem[i].Spent, pol.RetryBudget)
	}
	return mem, first
}

package daemon

import "github.com/dphaener/rhiza/internal/run"

// DefaultMaxRuns is how many runs the daemon keeps running at once.
const DefaultMaxRuns = 2

// Pick returns the change sets the daemon starts now, in the order of
// candidates, which are the accepted, ready change sets that serve a
// milestone and have a leaf, in roadmap order. runs are every run of the
// repository, oldest first. host and alive say which runs' processes are
// running (see run.State.Liveness).
//
// A change set with a live run, running or interrupted, is not picked, and
// neither is one whose latest run failed and was neither abandoned nor
// retried, which stays in the owner's inbox until cleared. Only running
// runs take one of max slots: an interrupted run takes none. pending are
// change sets the daemon started whose runs have no journal yet; each
// counts as running, and is not picked again, until a journal of it shows.
func Pick(candidates []string, runs []run.Journal, pending []string, host string, alive func(pid int) bool, max int) []string {
	live := map[string]bool{}
	latest := map[string]run.State{}
	running := 0
	for _, j := range runs {
		st := j.State()
		latest[st.Change] = st
		switch st.Liveness(host, alive) {
		case run.LiveRunning:
			live[st.Change] = true
			running++
		case run.LiveInterrupted:
			live[st.Change] = true
		}
	}
	held := map[string]bool{}
	for _, cs := range pending {
		if !live[cs] && !held[cs] {
			held[cs] = true
			running++
		}
	}
	var picked []string
	for _, cs := range candidates {
		if running >= max {
			break
		}
		if live[cs] || held[cs] {
			continue
		}
		if st, ok := latest[cs]; ok && st.Ended && st.Result == run.Failed && !st.Abandoned() {
			continue
		}
		picked = append(picked, cs)
		running++
	}
	return picked
}

// Resumable returns the IDs of the runs the daemon resumes now: those
// interrupted (see run.State.Liveness) that it has not resumed before. A
// run interrupted again after its one resume is left alone, and a finished
// run, failed or not, is never resumed.
func Resumable(runs []run.Journal, host string, alive func(pid int) bool) []string {
	var ids []string
	for _, j := range runs {
		st := j.State()
		if st.Liveness(host, alive) == run.LiveInterrupted && !st.DaemonResumed {
			ids = append(ids, j.ID)
		}
	}
	return ids
}

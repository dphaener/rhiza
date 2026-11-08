package daemon

import (
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// DefaultMaxRuns is how many runs the daemon keeps running at once.
const DefaultMaxRuns = 2

// Pick returns the change sets the daemon starts now, in the order of
// candidates, which are the accepted, ready change sets that serve a
// milestone and have a leaf, in roadmap order. runs are every run of the
// repository, oldest first. host and alive say which runs' processes are
// running (see run.State.Liveness).
//
// A change set with a live run, running or interrupted, is not picked, and
// neither is one whose latest run shipped it (see Ship) or failed and was neither abandoned nor
// retried, which stays in the owner's inbox until cleared. Only running
// runs take one of max slots: an interrupted run takes none. pending are
// change sets the daemon started whose runs have no journal yet; each
// counts as running, and is not picked again, until a journal of it shows.
func Pick(candidates []string, runs []run.Journal, pending []string, host string, alive func(pid int) bool, max int, ship Ship) []string {
	picked, _ := decide(candidates, runs, pending, host, alive, max, ship)
	return picked
}

// Skip reasons, fixed: the state file and its readers match on them.
const (
	SkipNoMilestone  = "no-milestone"
	SkipNoLeaves     = "no-leaves"
	SkipNotReady     = "not-ready"
	SkipLive         = "live"
	SkipFailed       = "failed"
	SkipShipped      = "shipped"
	SkipHeld         = "held"
	SkipUnjournaled  = "unjournaled"
	SkipNoSlot       = "no-slot"
	SkipDiskLow      = "disk-low"
	SkipUsageReserve = "usage-reserve"
	SkipShipOwner    = "ship-owner"
	SkipPaused       = "paused"
)

// Ship says which succeeded runs mean a change set has shipped: those of
// Workflow (the ship-change workflow's name, passed in so this package does
// not import api) that ended after the change set was last accepted or
// reopened, which Since reports (the zero time when it does not know).
type Ship struct {
	Workflow string
	Since    func(change string) time.Time
}

// shipped reports whether st is a succeeded run of the ship workflow that
// ended after its change set was last accepted or reopened.
func (s Ship) shipped(st run.State) bool {
	if s.Workflow == "" || s.Since == nil || !st.Ended || st.Result != run.Succeeded || st.Workflow != s.Workflow {
		return false
	}
	return st.EndedAt.After(s.Since(st.Change))
}

// Skip is an accepted change set the daemon does not start, and why.
type Skip struct {
	Change string `json:"change"`
	// Title is the change set's display title, filled in by the poll.
	Title string `json:"title,omitempty"`
	// Reason is one of the Skip constants.
	Reason string `json:"reason"`
	// Detail holds the readiness reasons of a not-ready change set, the free
	// space of a disk-low one, or why a no-slot one waits when a slot is kept
	// for a queued plan request.
	Detail []string `json:"detail,omitempty"`
}

// Skips returns, in the order of candidates, the candidates Pick (called
// with the same arguments) does not pick, each with its reason: live (a
// running or interrupted run), failed (the latest run failed and was not
// abandoned or retried, so it waits in the inbox), shipped (the latest run
// is a succeeded ship run that ended after the change set was last accepted
// or reopened; see Ship), held (started and not
// yet journaled) or no-slot (max runs are running). Candidates are the
// accepted, ready change sets that serve a milestone and have a leaf.
func Skips(candidates []string, runs []run.Journal, pending []string, host string, alive func(pid int) bool, max int, ship Ship) []Skip {
	_, skips := decide(candidates, runs, pending, host, alive, max, ship)
	return skips
}

func decide(candidates []string, runs []run.Journal, pending []string, host string, alive func(pid int) bool, max int, ship Ship) ([]string, []Skip) {
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
	var skips []Skip
	for _, cs := range candidates {
		switch st, ok := latest[cs]; {
		case live[cs]:
			skips = append(skips, Skip{Change: cs, Reason: SkipLive})
		case held[cs]:
			skips = append(skips, Skip{Change: cs, Reason: SkipHeld})
		case ok && ship.shipped(st):
			skips = append(skips, Skip{Change: cs, Reason: SkipShipped})
		case ok && st.Ended && st.Result == run.Failed && !st.Abandoned():
			skips = append(skips, Skip{Change: cs, Reason: SkipFailed})
		case running >= max:
			skips = append(skips, Skip{Change: cs, Reason: SkipNoSlot})
		default:
			picked = append(picked, cs)
			running++
		}
	}
	return picked, skips
}

// Unrecorded is how long after its start a run with no process record is
// taken to be starting rather than interrupted: a new run is journaled
// before the process running it records itself.
const Unrecorded = time.Minute

// Resumable returns the IDs of the runs the daemon resumes now: those
// interrupted (see run.State.Liveness) that it has not resumed before. A
// run interrupted again after its one resume is left alone, and a finished
// run, failed or not, is never resumed. Nor is a run whose lock another
// process holds (locked reports it), which is running whatever its journal
// says, or one with no process record started less than Unrecorded before
// now, whose process may not have recorded itself yet. Nor is a run whose
// change set state reports as neither accepted nor landed, where landed
// counts only for a run of the plan workflow, or of another workflow that
// entered its land step: a run on a change set that has been rejected,
// withdrawn or superseded, is still proposed or draft, or has landed without
// a non-plan run reaching land waits in the inbox for the owner.
func Resumable(runs []run.Journal, host string, alive func(pid int) bool, locked func(id string) bool, state func(change string) string, plan string, now time.Time) []string {
	var ids []string
	for _, j := range runs {
		st := j.State()
		if st.Liveness(host, alive) != run.LiveInterrupted || st.DaemonResumed {
			continue
		}
		if st.Process == nil && now.Sub(st.StartedAt) < Unrecorded {
			continue
		}
		if locked(j.ID) {
			continue
		}
		switch state(st.Change) {
		case "accepted":
		case "landed":
			if st.Workflow != plan && st.Workflow != plan+"-milestone" && st.Visits["land"] == 0 {
				continue
			}
		default:
			continue
		}
		ids = append(ids, j.ID)
	}
	return ids
}

// LiveRuns are the runs that are running or interrupted, in the order of
// runs.
func LiveRuns(runs []run.Journal, host string, alive func(pid int) bool) []LiveRun {
	out := []LiveRun{}
	for _, j := range runs {
		st := j.State()
		switch st.Liveness(host, alive) {
		case run.LiveRunning:
			out = append(out, LiveRun{ID: j.ID, Change: st.Change, State: "running"})
		case run.LiveInterrupted:
			out = append(out, LiveRun{ID: j.ID, Change: st.Change, State: "interrupted"})
		}
	}
	return out
}

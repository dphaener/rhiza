package api

// This file is the watchdog's handling of a green PR its run never merged
// (decision watchdog-scope-and-defaults). The daemon never merges outside a
// workflow: a failed run that can reach its own merge step without the owner
// is resumed with --retry so the workflow merges; a run that ended succeeded
// with the PR unmerged (run resume refuses --retry on it), or whose way to
// the merge step needs the owner, gets an inbox item. A live run is left
// alone.

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// mergeStepReach reports, for st's workflow, whether it has a step that runs
// `rhiza ship merge` (hasMerge) and whether that step can be reached from
// from without passing a signoff step (direct): the retry then merges with
// nobody's answer.
func mergeStepReach(st run.State, from string) (hasMerge, direct bool) {
	def, err := st.ParseDefinition()
	if err != nil {
		return false, false
	}
	isMerge := func(s workflow.Step) bool {
		return strings.Contains(s.Run, "rhiza ship merge")
	}
	for _, s := range def.Steps {
		if isMerge(s) {
			hasMerge = true
		}
	}
	if !hasMerge || from == "" {
		return hasMerge, false
	}
	seen := map[string]bool{}
	queue := []string{from}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if seen[n] {
			continue
		}
		seen[n] = true
		s, ok := def.Steps[n]
		if !ok {
			continue
		}
		if isMerge(s) {
			return true, true
		}
		if s.Kind == workflow.KindSignoff {
			continue
		}
		for _, next := range s.On {
			queue = append(queue, next)
		}
	}
	return true, false
}

// watchdogGreenUnmerged handles each ended run of journals whose change set's
// PR is open with every check green: a failed run that reaches its merge step
// on its own is resumed with --retry within the run's budget; a succeeded run,
// or one that needs the owner before it merges, is held as an inbox item.
// Runs with a live process, abandoned runs, runs of a workflow with no merge
// step, and attended runs under an autonomous-only policy are left alone. A
// run whose PR was found not open is remembered and not read again while its
// journal stays as it was, so the pass's gh reads stay bounded by the runs
// with an open PR plus those that ended or moved since the last poll. It
// returns the memory and holds after, and the first error.
func watchdogGreenUnmerged(pol WatchdogPolicy, journals []run.Journal, live func(id string) bool, mem []daemon.WatchdogRun, d watchdogDeps) ([]daemon.WatchdogRun, []daemon.WatchdogHold, error) {
	var holds []daemon.WatchdogHold
	var first error
	if d.ci == nil {
		return mem, nil, nil
	}
	for _, j := range journals {
		st := j.State()
		if live(j.ID) || !st.Ended || st.Abandoned() || len(j.Entries) == 0 {
			continue
		}
		if st.Result != run.Failed && st.Result != run.Succeeded {
			continue
		}
		if !st.Autonomous && pol.AutonomousOnly {
			continue
		}
		step, from := "", ""
		if st.Result == run.Failed {
			s, err := st.FailedStep()
			if err != nil {
				continue
			}
			step, from = s, s
		}
		hasMerge, direct := mergeStepReach(st, from)
		if !hasMerge {
			continue
		}
		last := j.Entries[len(j.Entries)-1].At
		k := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if k >= 0 && !mem[k].PRGone.IsZero() && mem[k].PRGone.Equal(last) {
			continue
		}
		pr, err := d.ci(st)
		if err != nil {
			if first == nil {
				first = fmt.Errorf("read the PR of run %s: %w", j.ID, err)
			}
			continue
		}
		if pr.Head == "" {
			// No open PR: an ended run's PR is not reopened by itself, so
			// later polls skip it until the journal moves on (a resume).
			if k < 0 {
				mem = append(mem, daemon.WatchdogRun{Run: j.ID})
				k = len(mem) - 1
			}
			mem[k].PRGone = last
			continue
		}
		if _, green, _ := checksVerdict(pr.Checks); !green {
			continue
		}
		hold := func(detail string) {
			holds = append(holds, daemon.WatchdogHold{
				Run: j.ID, Change: st.Change, Step: step, Since: last,
				Detail: detail,
			})
			d.log("watchdog: run %s: PR is green at %s but unmerged: left to the owner", j.ID, shortSHA(pr.Head))
		}
		if st.Result == run.Succeeded {
			hold(fmt.Sprintf("the run ended succeeded but its PR is green at %s and unmerged: merge it by hand, a resume with --retry is refused on a succeeded run", shortSHA(pr.Head)))
			continue
		}
		if !direct {
			hold(fmt.Sprintf("the run failed at step %s with its PR green at %s and unmerged: the workflow needs you before it merges", step, shortSHA(pr.Head)))
			continue
		}
		i := slices.IndexFunc(mem, func(w daemon.WatchdogRun) bool { return w.Run == j.ID })
		if i < 0 {
			mem = append(mem, daemon.WatchdogRun{Run: j.ID})
			i = len(mem) - 1
		}
		if mem[i].Spent >= pol.RetryBudget {
			hold(fmt.Sprintf("the run failed at step %s with its PR green at %s and unmerged, and the watchdog's retry budget of %d is spent", step, shortSHA(pr.Head), pol.RetryBudget))
			continue
		}
		if d.queued != nil && d.queued(j.ID) {
			continue
		}
		if err := d.resume(j.ID); err != nil {
			if first == nil {
				first = fmt.Errorf("resume run %s: %w", j.ID, err)
			}
			continue
		}
		d.charge(&mem[i], fmt.Sprintf("step %s failed with its PR green at %s and unmerged: resumed with --retry", step, shortSHA(pr.Head)))
		d.log("watchdog: run %s step %s failed with its PR green at %s and unmerged: resumed with --retry for the workflow's merge (budget %d of %d)", j.ID, step, shortSHA(pr.Head), mem[i].Spent, pol.RetryBudget)
		d.fileBug(pol, j.ID, bugUnmerged, step, fmt.Sprintf("PR green at %s and unmerged after step %s failed: %s", shortSHA(pr.Head), step, failedDetail(st, step)))
	}
	return mem, holds, first
}

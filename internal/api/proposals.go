package api

import (
	"slices"

	"github.com/dphaener/rhiza/internal/run"
)

// planReviewWorkflow is the built-in workflow that plans a change set
// already proposed: plan without its interview.
const planReviewWorkflow = "plan-review"

// planCoverWorkflow reports whether a run of workflow name is one whose end
// settles the proposals it gated: plan, plan-milestone or plan-review.
func planCoverWorkflow(name string) bool {
	return name == run.PlanWorkflow || name == run.PlanMilestoneWorkflow || name == planReviewWorkflow
}

// uncoveredProposals is the proposed change sets (proposed, with leaves) that
// no run covers, in the order given. It is the one predicate the daemon,
// which starts a plan-review run for each, and the inbox, which lists them,
// share. held are the change sets whose start is pending or not journaled.
// A change set is covered when
//   - a run on it has not ended (live, or interrupted and awaiting its
//     resume), or its start is held;
//   - a plan, plan-milestone or plan-review run (or any plan run), whatever
//     its own change set, has a latest gated entry naming it: a plan-milestone
//     run on the record change set gates every change set it serves;
//   - a run on it ended failed and was not abandoned (the owner retries it),
//     or a run of those workflows on it ended without failing.
//
// An abandoned run covers nothing, but the gated entry it journaled still
// does.
func uncoveredProposals(proposed []string, runs []run.Journal, held map[string]bool) []string {
	covered := map[string]bool{}
	for cs := range held {
		covered[cs] = true
	}
	for _, j := range runs {
		st := j.State()
		if planCoverWorkflow(st.Workflow) || st.IsPlanRun() {
			for _, g := range st.Gated {
				covered[g] = true
			}
		}
		if st.Change == "" {
			continue
		}
		switch {
		case !st.Ended:
			covered[st.Change] = true
		case st.Abandoned():
		case st.Result == run.Failed:
			covered[st.Change] = true
		case planCoverWorkflow(st.Workflow) || st.IsPlanRun():
			covered[st.Change] = true
		}
	}
	var out []string
	for _, cs := range proposed {
		if !covered[cs] && !slices.Contains(out, cs) {
			out = append(out, cs)
		}
	}
	return out
}

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
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

// inboxUncoveredProposals is one accept item per proposed change set with
// leaves that no run covers (uncoveredProposals), whatever
// rhiza.autonomy.accept says: nothing is going to accept it. Only the
// proposed change sets are folded when the fold has an index; a leaf is then
// read from the change set's own ops. A change set that already has an
// accept item in items (a gate left it for the owner) is not listed again.
func inboxUncoveredProposals(items []InboxItem, runs []run.Journal, fold *inboxFold) []InboxItem {
	if fold.load() != nil {
		return nil
	}
	var proposed []string
	for cs, st := range fold.states {
		if st != changeset.StateProposed {
			continue
		}
		c, ok := fold.changeSet(cs.String())
		if !ok {
			continue
		}
		sets := fold.list
		if sets == nil {
			sets = []changeset.ChangeSet{c}
		}
		if len(changeLeaves(newModelSets(sets, cs), c.Ops)) > 0 {
			proposed = append(proposed, cs.String())
		}
	}
	slices.Sort(proposed)
	return uncoveredAcceptItems(items, proposed, runs)
}

// uncoveredAcceptItems is one accept item per change set of proposed that
// uncoveredProposals leaves uncovered and that has no accept item in items.
func uncoveredAcceptItems(items []InboxItem, proposed []string, runs []run.Journal) []InboxItem {
	listed := map[string]bool{}
	for _, it := range items {
		if it.Kind == InboxAccept {
			listed[it.Change] = true
		}
	}
	var out []InboxItem
	for _, cs := range uncoveredProposals(proposed, runs, nil) {
		if listed[cs] {
			continue
		}
		cid, err := id.Parse(cs)
		if err != nil {
			continue
		}
		at := ulid.Time(ulid.ULID(cid).Time()).UTC().Format(time.RFC3339)
		sum := sha256.Sum256([]byte(strings.Join([]string{InboxAccept, cs, "uncovered"}, "\x00")))
		out = append(out, InboxItem{
			ID: hex.EncodeToString(sum[:])[:7], Kind: InboxAccept, At: at, Change: cs,
			Detail: "waiting for acceptance: no plan run covers it",
			Next:   ownerAcceptNext(cs),
		})
	}
	return out
}

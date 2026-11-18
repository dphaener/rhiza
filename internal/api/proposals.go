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
//   - a run on it ended failed and was not abandoned (the owner retries it).
//
// An abandoned run covers nothing, but the gated entry it journaled still
// does.
func uncoveredProposals(proposed []string, runs []run.Journal, held map[string]bool) []string {
	if len(proposed) == 0 {
		return nil
	}
	return uncoveredProposalsOf(proposed, runStates(runs), held)
}

// uncoveredProposalsOf is uncoveredProposals with each run's state already
// folded (runStates).
func uncoveredProposalsOf(proposed []string, states []run.State, held map[string]bool) []string {
	if len(proposed) == 0 {
		return nil
	}
	covered := map[string]bool{}
	for cs := range held {
		covered[cs] = true
	}
	for _, st := range states {
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
// rhiza.autonomy.accept says: nothing is going to accept it. The proposed
// change sets a run covers, or that already have an accept item in items (a
// gate left it for the owner), are dropped before any view is built, so
// only the change sets it reports cost a what-if overlay: a leaf is read
// from each remaining change set's own view over the projection's stored
// as-is rows.
func inboxUncoveredProposals(items []InboxItem, runs []run.Journal, fold *inboxFold) []InboxItem {
	return inboxUncoveredProposalsOf(items, runStates(runs), fold)
}

// inboxUncoveredProposalsOf is inboxUncoveredProposals with each run's
// state already folded (runStates).
func inboxUncoveredProposalsOf(items []InboxItem, states []run.State, fold *inboxFold) []InboxItem {
	if fold.load() != nil {
		return nil
	}
	var proposed []string
	for cs, st := range fold.states {
		if st == changeset.StateProposed {
			proposed = append(proposed, cs.String())
		}
	}
	slices.Sort(proposed)
	listed := acceptListed(items)
	var withLeaves []string
	for _, cs := range uncoveredProposalsOf(proposed, states, nil) {
		if !listed[cs] && fold.proposalHasLeaf(cs) {
			withLeaves = append(withLeaves, cs)
		}
	}
	return uncoveredAcceptItemsOf(items, withLeaves, states)
}

// proposalHasLeaf reports whether proposed change set change has a leaf in
// its own view over the projection's stored as-is rows, so a leaf under a
// landed parent counts. Only whether it has a leaf matters, so the leaves
// are not ranked.
func (f *inboxFold) proposalHasLeaf(change string) bool {
	c, ok := f.changeSet(change)
	if !ok {
		return false
	}
	if f.onProposalView != nil {
		f.onProposalView(c.ID)
	}
	m := f.p.model(f.ws.pending, c.ID)
	leaves := changeLeavesIn(m, c.Ops, func([]Node) map[id.ID]int { return nil })
	return f.p.failure() == nil && len(leaves) > 0
}

// acceptListed is the change sets that have an accept item in items.
func acceptListed(items []InboxItem) map[string]bool {
	listed := map[string]bool{}
	for _, it := range items {
		if it.Kind == InboxAccept {
			listed[it.Change] = true
		}
	}
	return listed
}

// uncoveredAcceptItems is one accept item per change set of proposed that
// uncoveredProposals leaves uncovered and that has no accept item in items.
func uncoveredAcceptItems(items []InboxItem, proposed []string, runs []run.Journal) []InboxItem {
	if len(proposed) == 0 {
		return nil
	}
	return uncoveredAcceptItemsOf(items, proposed, runStates(runs))
}

// uncoveredAcceptItemsOf is uncoveredAcceptItems with each run's state
// already folded (runStates).
func uncoveredAcceptItemsOf(items []InboxItem, proposed []string, states []run.State) []InboxItem {
	listed := acceptListed(items)
	var out []InboxItem
	for _, cs := range uncoveredProposalsOf(proposed, states, nil) {
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

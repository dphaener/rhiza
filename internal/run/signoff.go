package run

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

// SignoffWait is the detail of the stopped entry a signoff step writes
// while it waits for the owner's sign-off.
const SignoffWait = "waiting for the owner's sign-off: rhiza run approve | rhiza run request-changes <file>; rhiza run summary"

// signoff runs a signoff step: the review's last, which waits for the
// owner. A sign-off of this visit decides it: approved passes, and
// changes ends it with outcome changes, whose on.changes goes back to the
// checks step that fixes them. With none, an autonomous run signs off in
// the owner's place, journalling a signoff entry of rule autonomous, and
// passes; any other run stops, waiting (SignoffWait). It never polls:
// rhiza run approve and rhiza run request-changes journal the owner's
// sign-off and continue the run. The error return is for journal
// failures.
func (r *running) signoff(step workflow.Step) (result, error) {
	st := r.j.State()
	if s, ok := st.SignoffFor(step.Name, st.CurrentVisit); ok {
		who := s.Actor
		if who == "" {
			who = "rule " + s.Rule
		}
		if s.Outcome == SignoffChanges {
			return result{outcome: workflow.OutcomeChanges, detail: fmt.Sprintf("changes asked by %s: %d %s",
				who, len(s.Findings), plural(len(s.Findings), "finding", "findings"))}, nil
		}
		return result{outcome: workflow.OutcomePass, detail: "approved by " + who}, nil
	}
	if wf, owner := ownerDeclared(st); owner {
		r.log("step %s: waiting for the owner: workflow %s declares autonomy signoff owner", step.Name, wf)
		return result{stopped: SignoffWait}, nil
	}
	if st.Autonomous {
		if err := r.append(Entry{
			Type: TypeSignoff, Step: step.Name, Visit: st.CurrentVisit, Rule: RuleAutonomous, Outcome: SignoffApproved,
		}); err != nil {
			return result{}, err
		}
		r.log("step %s: signed off under rule autonomous", step.Name)
		return result{outcome: workflow.OutcomePass, detail: "approved by rule autonomous"}, nil
	}
	return result{stopped: SignoffWait}, nil
}

// ownerDeclared reports whether the run's workflow declares the signoff
// stage owner, and the name of the workflow that does. It reads the run's
// recorded stage autonomy (the top-level workflow's, so a signoff step of
// a nested workflow obeys it too); a journal that predates that record
// falls back to the started definition's declaration. A signoff that is
// owner only by the rhiza.autonomy.signoff setting is not a declaration.
func ownerDeclared(st State) (string, bool) {
	if len(st.Autonomy) > 0 {
		for _, a := range st.Autonomy {
			if a.Stage == workflow.StageSignoff && a.Source == SourceWorkflow && a.Mode == workflow.AutonomyOwner {
				return cmp.Or(a.Workflow, st.Workflow), true
			}
		}
		return "", false
	}
	def, err := st.ParseDefinition()
	if err != nil {
		return "", false
	}
	if m, ok := def.Autonomy(workflow.StageSignoff); ok && m == workflow.AutonomyOwner {
		return st.Workflow, true
	}
	return "", false
}

// plural is one when n is 1, else many.
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// ownerCycle opens checks step step's owner cycle when a request for
// changes waits for it (State.PendingChanges, from a signoff step whose
// on.changes names step): a check-cycle of kind owner, numbered after the
// step's latest cycle (from 2, so it has a fix phase), over the checks
// the owner's findings name, in the step's order, carrying the findings.
// Journalling it consumes the request. The error return is for journal
// failures.
func (r *running) ownerCycle(step workflow.Step) error {
	st := r.j.State()
	so, ok := st.PendingChanges()
	if !ok {
		return nil
	}
	if from, ok := r.def.Steps[so.Step]; !ok || from.On[workflow.OutcomeChanges] != step.Name {
		return nil
	}
	var checks []string
	for _, c := range step.Checks {
		if slices.ContainsFunc(so.Findings, func(f Finding) bool { return f.Check == c.Name }) {
			checks = append(checks, c.Name)
		}
	}
	cycle := 2
	if cur, ok := st.LatestCycle(step.Name); ok {
		cycle = max(cycle, cur.Cycle+1)
	}
	if err := r.append(Entry{
		Type: TypeCheckCycle, Step: step.Name, Cycle: cycle, Checks: checks, Kind: CycleOwner, Findings: so.Findings,
	}); err != nil {
		return err
	}
	r.log("step %s: cycle %d, the owner's changes: %s", step.Name, cycle, strings.Join(checks, ", "))
	return nil
}

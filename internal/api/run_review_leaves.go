package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// leafContractFields are the fields every live leaf of a plan must carry,
// in the order a problem names them.
var leafContractFields = []string{"statement", "acceptance", "files", "must_not"}

// leafContractProblem is one live leaf of a plan that lacks part of its
// contract: Lacks names each missing field, then "a failure case" when its
// acceptance holds no failure case and it has no failure_waiver.
type leafContractProblem struct {
	Node   id.ID
	Handle string
	Lacks  []string
}

func (p leafContractProblem) line() string {
	name := p.Handle
	if name == "" {
		name = p.Node.String()
	} else {
		name += " (" + p.Node.String() + ")"
	}
	return "- " + name + " lacks " + strings.Join(p.Lacks, ", ")
}

// leafContractProblems checks each live leaf (a live capability with no
// live child capability in any plan change set's view, see hasPlanChild) that a change set of plan creates or changes
// with an op: it must have a non-empty statement, acceptance, files and
// must_not, and its acceptance must hold the words "failure case" (any
// case) unless it has a failure_waiver. A node several plan change sets
// touch is read in the view of the last of them, in plan order; a landed
// leaf no plan op touches is not checked. Problems come in plan order.
func leafContractProblems(w *workspace, plan []changeset.ChangeSet) []leafContractProblem {
	var order []id.ID
	view := map[id.ID]id.ID{}
	for _, c := range plan {
		for _, o := range c.Ops {
			switch o.Type {
			case op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeMove:
			default:
				continue
			}
			if _, ok := view[o.Node]; !ok {
				order = append(order, o.Node)
			}
			view[o.Node] = c.ID
		}
	}
	var out []leafContractProblem
	for _, n := range order {
		m := w.model(view[n])
		node, ok := m.Node(n)
		if !ok || node.Kind != op.KindCapability || node.Retired || hasPlanChild(w, plan, view, n) {
			continue
		}
		var lacks []string
		for _, f := range leafContractFields {
			if leafField(node, f) == "" {
				lacks = append(lacks, f)
			}
		}
		if !strings.Contains(strings.ToLower(leafField(node, "acceptance")), "failure case") && leafField(node, "failure_waiver") == "" {
			lacks = append(lacks, "a failure case (or a failure_waiver)")
		}
		if len(lacks) > 0 {
			out = append(out, leafContractProblem{Node: n, Handle: node.Handle, Lacks: lacks})
		}
	}
	return out
}

// hasPlanChild reports whether capability n has a live child capability
// in any plan change set's view: a parent whose child another plan change
// set adds is not a leaf. Each child is judged in the view of the last
// plan change set that touches it (view), so a child a later one retires
// or moves away does not count.
func hasPlanChild(w *workspace, plan []changeset.ChangeSet, view map[id.ID]id.ID, n id.ID) bool {
	for _, c := range plan {
		m := w.model(c.ID)
		for _, ch := range m.Children(op.KindCapability, n) {
			cm := m
			if v, ok := view[ch.ID]; ok {
				cm = w.model(v)
			}
			if got, ok := cm.Node(ch.ID); ok && !got.Retired && got.Kind == op.KindCapability && got.Parent == n {
				return true
			}
		}
	}
	return false
}

// leafField is node's string field f, trimmed; "" when it is absent,
// empty or not a string.
func leafField(n Node, f string) string {
	var s string
	if raw, ok := n.Fields[f]; ok {
		_ = json.Unmarshal(raw, &s)
	}
	return strings.TrimSpace(s)
}

// planLeafContract is the plan-review leaf-contract check of review run
// rr: a line per live leaf of the run's plan (planChangeSets) that lacks
// part of its contract, none when every leaf has it.
func planLeafContract(ctx *Ctx, cmd string, rr reviewRun) ([]string, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	defer w.closeProjection()
	if _, ok := w.changeSet(rr.change); !ok {
		return nil, NotFoundf("run %s's change set %s is not in the log", rr.id, rr.change).WithNext("rhiza run status " + rr.id)
	}
	plan, err := planChangeSets(w.changeSets(), w.model(rr.change), rr.change)
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err).WithNext("rhiza run status " + rr.id)
	}
	var lines []string
	for _, p := range leafContractProblems(w, plan) {
		lines = append(lines, p.line())
	}
	return lines, nil
}

// leafContractError is review-check's failure for the leaf-contract
// problems leaves of run rid, at any round; findings, when not empty, is
// the round's own failure, said after it.
func leafContractError(cmd, rid string, change id.ID, leaves []string, findings string) error {
	msg := fmt.Sprintf("%s: run %s's plan has %d leaf(s) without their contract (statement, acceptance, files, must_not, and an acceptance holding a `failure case` or a failure_waiver); revise the plan and review it again:\n%s",
		cmd, rid, len(leaves), strings.Join(leaves, "\n"))
	if findings != "" {
		msg += "\n" + findings
	}
	return Errorf("%s", msg).WithNext("rhiza slice --change " + change.String())
}

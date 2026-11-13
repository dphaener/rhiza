package modellint

import (
	"slices"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// UnappliedDecisions reports each live Decision cs adds or changes that is
// never applied: no live node links to it with decided-by, or every node
// that does is a leaf of cs (a capability cs creates with no live child
// capability) none of whose handle or ID is in committed. committed names
// the leaves with a branch commit carrying their Leaf: trailer. A Decision
// whose only link is answers on a question is not applied by that link and
// is exempt from the first clause. Decisions cs did not write are never
// looked at. Hits are sorted by node.
func UnappliedDecisions(m model.Model, cs changeset.ChangeSet, committed []string) []Hit {
	created := map[id.ID]bool{}
	for _, o := range cs.Ops {
		if o.Type == op.TypeNodeCreate {
			created[o.Node] = true
		}
	}
	hasChild := map[id.ID]bool{}
	for _, n := range m.Nodes() {
		if n.Kind == op.KindCapability && !n.Retired {
			hasChild[n.Parent] = true
		}
	}
	// unbuilt reports a leaf of cs with no Leaf: commit.
	unbuilt := func(n model.Node) bool {
		if n.Kind != op.KindCapability || !created[n.ID] || hasChild[n.ID] {
			return false
		}
		return !slices.Contains(committed, n.Handle) && !slices.Contains(committed, n.ID.String())
	}

	appliers := map[id.ID][]model.Node{}
	answered := map[id.ID]bool{}
	for _, l := range m.Links() {
		if l.Retired {
			continue
		}
		switch l.Type {
		case op.LinkDecidedBy:
			// The decided-by link runs from the node to the Decision.
			if n, ok := m.Node(l.From); ok && !n.Retired {
				appliers[l.To] = append(appliers[l.To], n)
			}
		case op.LinkAnswers:
			answered[l.From], answered[l.To] = true, true
		}
	}

	var out []Hit
	for _, i := range writtenNodes(cs) {
		n, ok := m.Node(i)
		if !ok || n.Retired || n.Kind != op.KindDecision {
			continue
		}
		as := appliers[n.ID]
		if len(as) == 0 {
			if answered[n.ID] {
				continue
			}
		} else if slices.ContainsFunc(as, func(a model.Node) bool { return !unbuilt(a) }) {
			continue
		}
		out = append(out, Hit{Node: n.ID, Handle: n.Handle, Kind: KindUnapplied})
	}
	return out
}

package changeset

import (
	"slices"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// Resting is the ops of one change set that rest on a dropped change set:
// one that was withdrawn, rejected or superseded.
type Resting struct {
	// ChangeSet is the dropped change set and State its state.
	ChangeSet id.ID
	State     State
	// Ops are the resting op IDs, sorted.
	Ops []id.ID
}

// RestsOn finds the current ops of change set cs that rest on a dropped
// change set X in sets. An op rests on X when its target, parent, from, to
// or base is a node or link X's ops create, or the ID of one of X's ops,
// and no landed or in-flight change set (cs included) also creates it: with
// X gone, nothing the op names would exist. An op can rest on more than
// one dropped change set. The result is sorted by X, each Resting's Ops by
// op ID, and is nil when nothing rests. A cs not in sets has no ops, so
// nothing rests.
func RestsOn(sets []ChangeSet, cs id.ID) []Resting {
	return Analyze(sets).RestsOn(cs)
}

// RestsOn is RestsOn(sets, cs) for a's sets.
func (a *Analysis) RestsOn(cs id.ID) []Resting {
	s, _ := a.set(cs)
	mine := s.Ops
	if len(mine) == 0 {
		return nil
	}
	a.rest()

	var out []Resting
	for _, x := range a.dropped {
		if x.ID == cs {
			continue
		}
		var ops []id.ID
		for _, o := range mine {
			if slices.ContainsFunc(Names(o), func(n id.ID) bool { return x.gone[n] }) {
				ops = append(ops, o.ID)
			}
		}
		if len(ops) > 0 {
			slices.SortFunc(ops, id.Compare)
			out = append(out, Resting{ChangeSet: x.ID, State: x.State, Ops: ops})
		}
	}
	slices.SortFunc(out, func(a, b Resting) int {
		return id.Compare(a.ChangeSet, b.ChangeSet)
	})
	return out
}

// creates returns the node or link o creates, if it is a create.
func creates(o op.Op) (id.ID, bool) {
	if o.Type == op.TypeNodeCreate || o.Type == op.TypeLinkCreate {
		return o.Target(), true
	}
	return id.ID{}, false
}

// Names lists the IDs o rests on: its target, parent, from, to and base.
// The zero ID is left out.
func Names(o op.Op) []id.ID {
	var out []id.ID
	for _, n := range []id.ID{o.Target(), o.Parent, o.From, o.To, o.Base} {
		if !n.IsZero() {
			out = append(out, n)
		}
	}
	return out
}

package changeset

import (
	"slices"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
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
	var mine []op.Op
	kept := map[id.ID]bool{}
	for _, s := range sets {
		if s.ID == cs {
			mine = s.Ops
		}
		if s.State == StateLanded || s.State.inFlight() {
			for _, o := range s.Ops {
				if created, ok := creates(o); ok {
					kept[created] = true
				}
			}
		}
	}
	if len(mine) == 0 {
		return nil
	}

	var out []Resting
	for _, x := range sets {
		switch x.State {
		case StateWithdrawn, StateRejected, StateSuperseded:
		default:
			continue
		}
		if x.ID == cs {
			continue
		}
		gone := map[id.ID]bool{}
		for _, o := range x.Ops {
			gone[o.ID] = true
			if created, ok := creates(o); ok && !kept[created] {
				gone[created] = true
			}
		}
		var ops []id.ID
		for _, o := range mine {
			if slices.ContainsFunc(names(o), func(n id.ID) bool { return gone[n] }) {
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

// names lists the IDs o rests on: its target, parent, from, to and base.
// The zero ID is left out.
func names(o op.Op) []id.ID {
	var out []id.ID
	for _, n := range []id.ID{o.Target(), o.Parent, o.From, o.To, o.Base} {
		if !n.IsZero() {
			out = append(out, n)
		}
	}
	return out
}

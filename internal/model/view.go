package model

import (
	"fmt"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
)

// AsIs returns the as-is view of sets: the landed change sets in apply order,
// each one's operations in their own order. The problems are the order's and
// the operations', sorted, nil when there are none.
func AsIs(sets []changeset.ChangeSet) (Model, []Problem) {
	return asIs(sets).finish()
}

// ToBe returns the to-be view of sets: as-is, then the accepted change sets
// in their own apply order on top. The problems are as-is's and the accepted
// change sets', each once.
func ToBe(sets []changeset.ChangeSet) (Model, []Problem) {
	b := asIs(sets)
	layer(b, inState(sets, changeset.StateAccepted))
	return b.finish()
}

// WhatIf returns the what-if view of change set cs: as-is, then cs on top.
// When cs is not among sets the result is as-is plus an ErrMissing problem
// on cs; when it is not draft, proposed or accepted, as-is plus an ErrState
// problem on cs.
func WhatIf(sets []changeset.ChangeSet, cs id.ID) (Model, []Problem) {
	b := asIs(sets)
	i := -1
	for j, s := range sets {
		if s.ID == cs {
			i = j
			break
		}
	}
	switch {
	case i < 0:
		b.problems = append(b.problems, Problem{ChangeSet: cs,
			Err: fmt.Errorf("%w: change set %s is not in the log", ErrMissing, cs)})
	case !inFlight(sets[i].State):
		b.problems = append(b.problems, Problem{ChangeSet: cs,
			Err: fmt.Errorf("%w: change set %s is %q, not draft, proposed or accepted", ErrState, cs, sets[i].State)})
	default:
		applySets(b, sets[i:i+1])
	}
	return b.finish()
}

// asIs returns a builder that knows every change set in sets, whatever its
// state, with the landed change sets applied.
func asIs(sets []changeset.ChangeSet) *builder {
	known := make([]id.ID, 0, len(sets))
	for _, s := range sets {
		known = append(known, s.ID)
	}
	b := newBuilder(known)
	layer(b, inState(sets, changeset.StateLanded))
	return b
}

// layer applies sets to b in their apply order and records the order's
// problems in b.
func layer(b *builder, sets []changeset.ChangeSet) {
	ordered, problems := order(sets)
	b.problems = append(b.problems, problems...)
	applySets(b, ordered)
}

// applySets applies every operation of sets to b, change set by change set in
// the order given, each one's operations in their own order.
func applySets(b *builder, sets []changeset.ChangeSet) {
	for _, s := range sets {
		for _, o := range s.Ops {
			b.apply(s.ID, o)
		}
	}
}

// inState returns the change sets in sets whose state is s, nil when there
// are none.
func inState(sets []changeset.ChangeSet, s changeset.State) []changeset.ChangeSet {
	var out []changeset.ChangeSet
	for _, cs := range sets {
		if cs.State == s {
			out = append(out, cs)
		}
	}
	return out
}

// inFlight reports whether s is draft, proposed or accepted.
func inFlight(s changeset.State) bool {
	return s == changeset.StateDraft || s == changeset.StateProposed || s == changeset.StateAccepted
}

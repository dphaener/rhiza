package changeset

import (
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// An Analysis answers DependsOn, Closure, Dependencies, RestsOn and Revised
// for many change sets of one sets, with the same results as those
// functions: what they would each work out from sets again (the landed
// links, the ops by change set, what is kept and what each dropped change
// set takes with it) it works out once, on first use. It is not safe for
// concurrent use, and sets must not change while it is in use.
type Analysis struct {
	sets []ChangeSet

	// Set by graphed: each change set's state and ops, and the landed
	// depends-on links.
	graphed bool
	states  map[id.ID]State
	ops     map[id.ID][]op.Op
	landed  links

	// byID is each change set by ID, the last of a repeated ID; set by
	// sets.
	byID map[id.ID]ChangeSet

	// Set by rested: the dropped change sets, in sets order, each with the
	// IDs it takes with it.
	rested  bool
	dropped []droppedSet
}

// droppedSet is a dropped change set and the IDs that go with it: its op
// IDs and what it creates that nothing kept also creates.
type droppedSet struct {
	ChangeSet
	gone map[id.ID]bool
}

// Analyze returns an Analysis of sets.
func Analyze(sets []ChangeSet) *Analysis {
	return &Analysis{sets: sets}
}

// graph sets the states, the ops by change set and the landed links.
func (a *Analysis) graph() {
	if a.graphed {
		return
	}
	a.graphed = true
	a.states = make(map[id.ID]State, len(a.sets))
	for _, s := range a.sets {
		a.states[s.ID] = s.State
	}
	a.landed, a.ops = landedLinks(a.sets), opsByID(a.sets)
}

// set returns the change set cs, the last of a repeated ID.
func (a *Analysis) set(cs id.ID) (ChangeSet, bool) {
	if a.byID == nil {
		a.byID = make(map[id.ID]ChangeSet, len(a.sets))
		for _, s := range a.sets {
			a.byID[s.ID] = s
		}
	}
	s, ok := a.byID[cs]
	return s, ok
}

// rest sets the dropped change sets with what each takes with it.
func (a *Analysis) rest() {
	if a.rested {
		return
	}
	a.rested = true
	kept := map[id.ID]bool{}
	for _, s := range a.sets {
		if s.State == StateLanded || s.State.inFlight() {
			for _, o := range s.Ops {
				if created, ok := creates(o); ok {
					kept[created] = true
				}
			}
		}
	}
	for _, x := range a.sets {
		if !Dropped(x.State) {
			continue
		}
		gone := map[id.ID]bool{}
		for _, o := range x.Ops {
			gone[o.ID] = true
			if created, ok := creates(o); ok && !kept[created] {
				gone[created] = true
			}
		}
		a.dropped = append(a.dropped, droppedSet{ChangeSet: x, gone: gone})
	}
}

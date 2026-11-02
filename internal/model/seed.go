package model

import (
	"errors"
	"slices"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// State is what building the as-is view knows beyond the model and its
// problems, and what a what-if needs to go on from the stored as-is: the
// aspects each history entry of a node or link wrote, where each node's
// handle or parent was last set, the conflicts later writes may still
// repair, and the count of operations seen.
type State struct {
	// Seq counts the operations applied or skipped.
	Seq int
	// Aspects holds, by node or link ID, the aspects each entry of its
	// History wrote, parallel to it.
	Aspects map[id.ID][][]string
	// Placed holds, by node ID, when its handle or parent was last set.
	Placed map[id.ID]Mark
	// Held are the conflicts that later writes may still repair.
	Held []Held
}

// Mark is the operation that last set a node's handle or parent, and its
// position in the apply order.
type Mark struct {
	Seq int
	Ref Ref
}

// Held is a conflict of operation Op of ChangeSet on Target that later
// writes may still repair, with the aspects of the operation not yet
// rewritten.
type Held struct {
	Target    id.ID
	ChangeSet id.ID
	Op        id.ID
	Left      []string
}

// AsIsState is AsIs, with the State a what-if goes on from.
func AsIsState(sets []changeset.ChangeSet) (Model, []Problem, State) {
	b := asIs(sets)
	st := State{Seq: b.seq, Aspects: make(map[id.ID][][]string, len(b.m.nodes)+len(b.m.links)),
		Placed: make(map[id.ID]Mark, len(b.placed))}
	for k, r := range b.m.nodes {
		st.Aspects[k] = r.aspects
	}
	for k, r := range b.m.links {
		st.Aspects[k] = r.aspects
	}
	for k, m := range b.placed {
		st.Placed[k] = Mark{Seq: m.seq, Ref: m.ref}
	}
	for t, hs := range b.held {
		for _, h := range hs {
			p := b.problems[h.problem]
			st.Held = append(st.Held, Held{Target: t, ChangeSet: p.ChangeSet, Op: p.Op, Left: slices.Clone(h.left)})
		}
	}
	slices.SortFunc(st.Held, func(a, b Held) int {
		if c := id.Compare(a.Target, b.Target); c != 0 {
			return c
		}
		if c := id.Compare(a.ChangeSet, b.ChangeSet); c != 0 {
			return c
		}
		return id.Compare(a.Op, b.Op)
	})
	m, ps := b.finish()
	return m, ps, st
}

// Seed is the as-is view as stored: the landed change sets already applied.
// A what-if goes on from it instead of applying them again.
type Seed struct {
	// Known are the IDs of every change set in the log.
	Known []id.ID
	// Nodes and Links are the as-is nodes and links, each with the aspects
	// of its History entries.
	Nodes []SeedNode
	Links []SeedLink
	// Problems are the as-is problems without the duplicate-handle ones,
	// which a view reports on its final state.
	Problems []Problem
	// State is the as-is state, but for Aspects, which Nodes and Links hold.
	State State
}

// SeedNode is a node and the aspects of its History entries.
type SeedNode struct {
	Node    Node
	Aspects [][]string
	Placed  Mark
}

// SeedLink is a link and the aspects of its History entries.
type SeedLink struct {
	Link    Link
	Aspects [][]string
}

// WhatIfFrom is WhatIf over a seed: the as-is view is the seed instead of
// the landed change sets of sets. sets need hold only what WhatIf reads of
// them: every change set's ID and state; the whole of cs, of the in-flight
// change sets it depends on and of the landed ones it depends on, through
// others or not; and the depends-on link.create and link.retire ops of the
// landed change sets. With those it gives WhatIf's result over the whole
// log whose as-is view is the seed.
func WhatIfFrom(seed Seed, sets []changeset.ChangeSet, cs id.ID) (Model, []Problem) {
	return whatIf(seed.builder(), sets, cs, false)
}

// builder returns a builder holding the seed.
func (s Seed) builder() *builder {
	b := newBuilder(s.Known)
	b.seq = s.State.Seq
	b.m.nodes = make(map[id.ID]*record[Node], len(s.Nodes))
	b.m.links = make(map[id.ID]*record[Link], len(s.Links))
	b.placed = make(map[id.ID]mark, len(s.Nodes))
	for _, n := range s.Nodes {
		b.m.nodes[n.Node.ID] = &record[Node]{v: n.Node, aspects: n.Aspects}
		if op.HasParent(n.Node.Kind) {
			b.addChild(n.Node.Parent, n.Node.ID)
		}
		b.placed[n.Node.ID] = mark{seq: n.Placed.Seq, ref: n.Placed.Ref}
	}
	for _, l := range s.Links {
		b.m.links[l.Link.ID] = &record[Link]{v: l.Link, aspects: l.Aspects}
	}
	b.problems = slices.Clone(s.Problems)
	type key struct{ cs, op id.ID }
	at := make(map[key]int, len(b.problems))
	for i, p := range b.problems {
		if !p.Op.IsZero() {
			at[key{p.ChangeSet, p.Op}] = i
		}
	}
	for _, h := range s.State.Held {
		i, ok := at[key{h.ChangeSet, h.Op}]
		if !ok {
			continue
		}
		b.held[h.Target] = append(b.held[h.Target], held{problem: i, left: slices.Clone(h.Left)})
	}
	return b
}

// IsDuplicateHandle reports whether err is the duplicate-handle problem of
// a view's last step.
func IsDuplicateHandle(err error) bool {
	return errors.Is(err, ErrHandle)
}

// Source gives a delta computation the as-is rows it needs, one at a time,
// from wherever they are stored.
type Source interface {
	// SeedNode returns the as-is node i with its aspects and placement.
	SeedNode(i id.ID) (SeedNode, bool, error)
	// SeedLink returns the as-is link i with its aspects.
	SeedLink(i id.ID) (SeedLink, bool, error)
}

// Delta is what a what-if changes over the as-is view: every node and link
// its operations read or wrote, as they stand in the what-if. A node or
// link the what-if view has that Delta lacks is as the as-is view has it.
// Nothing is ever removed from a view, only retired, so Delta names every
// row that differs.
type Delta struct {
	Nodes []Node
	Links []Link
}

// WhatIfDelta is WhatIfFrom over a seed with no Nodes and no Links, whose
// as-is rows src gives as the operations reach them: it never holds the
// whole as-is view. It returns the nodes and links that were read, sorted
// by ID, and the problems of the change sets' own operations; the
// duplicate-handle problems are left out, as they need the whole view.
// The error is src's.
func WhatIfDelta(seed Seed, src Source, sets []changeset.ChangeSet, cs id.ID) (Delta, []Problem, error) {
	b := seed.builder()
	b.src = src
	whatIfLayer(b, sets, cs)
	if b.srcErr != nil {
		return Delta{}, nil, b.srcErr
	}
	var d Delta
	for _, r := range b.m.nodes {
		d.Nodes = append(d.Nodes, copyNode(r.v))
	}
	for _, r := range b.m.links {
		d.Links = append(d.Links, copyLink(r.v))
	}
	slices.SortFunc(d.Nodes, func(a, b Node) int { return id.Compare(a.ID, b.ID) })
	slices.SortFunc(d.Links, func(a, b Link) int { return id.Compare(a.ID, b.ID) })
	var ps []Problem
	for i, p := range b.problems {
		if _, ok := b.repaired[i]; !ok {
			ps = append(ps, p)
		}
	}
	return d, sortProblems(ps), nil
}

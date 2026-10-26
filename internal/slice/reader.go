package slice

import (
	"fmt"
	"slices"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
)

// Direction picks which links of a node Links returns.
type Direction int

// The directions.
const (
	// From is the links whose From is the node.
	From Direction = iota
	// To is the links whose To is the node.
	To
	// Both is the links with the node at either end, each once.
	Both
)

// Reader is one view of the model, as Build reads it. Every method can fail,
// since a view may be read from a database.
type Reader interface {
	// Node returns the node with ID i, live or retired.
	Node(i id.ID) (model.Node, bool, error)
	// Ancestors returns the parent chain of node i, the root first, without
	// i; nil for a root or a node outside the tree. An unknown i is an error
	// wrapping ErrNotFound.
	Ancestors(i id.ID) ([]model.Node, error)
	// Children returns the nodes directly under parent, live and retired,
	// sorted by key then ID.
	Children(parent id.ID) ([]model.Node, error)
	// Links returns the links at node i in direction dir, live and retired,
	// sorted by ID. i may be a change set.
	Links(i id.ID, dir Direction) ([]model.Link, error)
	// ShortID returns the short ID of i: see the package doc.
	ShortID(i id.ID) (string, error)
}

// LinkReader is a Reader that can also look a link up by ID. Build uses it,
// when the Reader has it, to find the endpoints of the links a change set
// retires.
type LinkReader interface {
	Link(i id.ID) (model.Link, bool, error)
}

// FromModel returns a Reader over m, with its links indexed by endpoint.
func FromModel(m model.Model) Reader {
	r := modelReader{m: m, nodes: m.Nodes(), from: map[id.ID][]model.Link{}, to: map[id.ID][]model.Link{}}
	// m.Links is sorted by ID, so each endpoint's list is too.
	for _, l := range m.Links() {
		r.from[l.From] = append(r.from[l.From], l)
		r.to[l.To] = append(r.to[l.To], l)
	}
	return r
}

type modelReader struct {
	m     model.Model
	nodes []model.Node // sorted by ID
	from  map[id.ID][]model.Link
	to    map[id.ID][]model.Link
}

func (r modelReader) Node(i id.ID) (model.Node, bool, error) {
	n, ok := r.m.Node(i)
	return n, ok, nil
}

func (r modelReader) Link(i id.ID) (model.Link, bool, error) {
	l, ok := r.m.Link(i)
	return l, ok, nil
}

func (r modelReader) Ancestors(i id.ID) ([]model.Node, error) {
	n, ok := r.m.Node(i)
	if !ok {
		return nil, fmt.Errorf("%w: no node %s", ErrNotFound, i)
	}
	var out []model.Node
	// The bound ends even a cycle the model should never produce.
	for p := n.Parent; !p.IsZero() && len(out) < 1000; {
		pn, ok := r.m.Node(p)
		if !ok {
			break
		}
		out = append(out, pn)
		p = pn.Parent
	}
	slices.Reverse(out)
	return out, nil
}

func (r modelReader) Children(parent id.ID) ([]model.Node, error) {
	return r.m.Children(parent), nil
}

func (r modelReader) Links(i id.ID, dir Direction) ([]model.Link, error) {
	switch dir {
	case From:
		return slices.Clone(r.from[i]), nil
	case To:
		return slices.Clone(r.to[i]), nil
	case Both:
		out := slices.Concat(r.from[i], r.to[i])
		slices.SortFunc(out, func(a, b model.Link) int { return id.Compare(a.ID, b.ID) })
		return slices.CompactFunc(out, func(a, b model.Link) bool { return a.ID == b.ID }), nil
	}
	return nil, fmt.Errorf("slice: unknown direction %d", dir)
}

func (r modelReader) ShortID(i id.ID) (string, error) {
	return shortID(r.nodes, i)
}

// shortID returns the short ID of i among nodes, which are sorted by ID: one
// more character than it shares with the nearest other ID on either side,
// at least six. i need not be among nodes.
func shortID(nodes []model.Node, i id.ID) (string, error) {
	k, found := slices.BinarySearchFunc(nodes, i, func(n model.Node, t id.ID) int { return id.Compare(n.ID, t) })
	var prev, next id.ID
	if k > 0 {
		prev = nodes[k-1].ID
	}
	if found {
		k++
	}
	if k < len(nodes) {
		next = nodes[k].ID
	}
	return Short(i, prev, next), nil
}

// Short returns the short ID of i given the IDs just before and just after
// it in the view (zero when there is none): one more character than its
// longest common prefix with either, at least 6 and at most 26.
func Short(i, prev, next id.ID) string {
	s := i.String()
	n := 6
	for _, o := range []id.ID{prev, next} {
		if o.IsZero() || o == i {
			continue
		}
		if c := commonPrefix(s, o.String()) + 1; c > n {
			n = c
		}
	}
	return s[:min(n, len(s))]
}

func commonPrefix(a, b string) int {
	n := 0
	for n < len(a) && n < len(b) && a[n] == b[n] {
		n++
	}
	return n
}

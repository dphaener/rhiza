// Package model computes the model: the nodes and links that the change sets
// in the log describe, as they stand once their operations are applied.
//
// The vocabulary:
//
//   - The model is a set of nodes (with their tree of parents and keys) and
//     links. A Model is one computed view of it.
//   - A view is the model under a choice of change sets. As-is is the landed
//     change sets applied. To-be is as-is plus every accepted change set.
//     What-if is as-is plus one in-flight change set, for review.
//   - The apply order is the order operations are applied in: change set by
//     change set, each change set's operations in their own order.
//   - Every operation except a create names a base: the version of its
//     target it was written against.
//   - A node's or link's version is the ID of the last operation applied to
//     it in this view. In as-is, that is the version new writes name as
//     their base.
//   - Its history is every operation applied to it, as (change set, op)
//     references in apply order, the create first.
//   - A base is stale when it is in the target's history but is not its
//     version: other operations were applied to the target since.
//   - Each operation writes aspects of its target: a field, the handle, the
//     position (parent and key) or its retirement. Two operations overlap
//     when they write a common aspect; a retirement overlaps everything.
//   - A conflict is an operation with a stale base that overlaps one of the
//     operations applied after its base. It is reported and still applied,
//     so the later operation in apply order wins; a small change set written
//     against the new version repairs it. A stale operation that overlaps
//     nothing merges silently: editing different fields is not a conflict.
//   - A skipped operation breaks a rule of the model (a missing target, a
//     base not in its history, a cycle, ...). It is reported and changes
//     nothing.
//   - Two in-flight change sets overlap when they write overlapping aspects
//     of the same target, whatever their bases.
//
// Computing a view never fails and never panics. Whatever is wrong becomes a
// Problem and the computation goes on.
//
// A view is built in layers on one engine. The as-is layer is the landed
// change sets in their apply order (see order.go). To-be adds the accepted
// change sets on top, in their own apply order; what-if adds the one change
// set asked for, which must be draft, proposed or accepted (ErrState) and
// among the input (ErrMissing). A link endpoint may name any change set in
// the input, whatever its state. Change sets with no state (no valid open)
// are in no view. An operation in an upper layer whose base a landed
// operation has since superseded is stale like any other: it merges or
// conflicts by the rules below. The problems are every layer's, each once.
//
// An operation is applied in these steps.
//
//  1. A create (node.create, link.create) needs an ID that is not already a
//     node or a link (ErrExists). A node's non-zero parent must be a node
//     (ErrMissing), live (ErrRetired), of a kind with children (ErrParent).
//     A link's endpoints must each be a node, live or retired, or a known
//     change set (ErrMissing); a non-zero decision must be a node
//     (ErrMissing). The record starts with the create as its history and
//     version, and its change set as its creator.
//  2. Every other operation needs a target of the right sort: a node for
//     node operations, a link for link.retire (ErrMissing). Its base must be
//     in the target's history (ErrBase). When the base is the version and
//     the target is retired, it is skipped (ErrRetired). When the base is
//     stale, the aspects of every operation applied after it are compared
//     with this operation's.
//  3. The aspects: node.update writes "field:<name>" for each name it sets
//     or unsets, and "handle" when it names a handle. node.move writes
//     "position". node.retire and link.retire write "retired", which
//     overlaps every other aspect. Creates write none.
//  4. The effects. An update sets and unsets fields and replaces the handle
//     when it names one. A move sets parent and key once its node's kind
//     sits in the tree (ErrParent) and its new parent, if any, is a node
//     (ErrMissing), live (ErrRetired), of a kind with children (ErrParent)
//     and neither the node nor one of its descendants (ErrCycle). A retire
//     sets Retired and Reason; link.retire sets Retired. A skipped
//     operation reports only why it was skipped. An applied one is reported
//     as ErrConflict when step 2 found an overlap, is appended to the
//     history and becomes the version.
//  5. After the last operation, live nodes are grouped by (parent, handle)
//     when their kind sits in the tree and by (kind, handle) otherwise. In
//     each group of two or more, ordered by when each node's handle or
//     parent was last set, every node after the first is reported
//     (ErrHandle) on the operation that last set it. Nothing is skipped.
package model

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

var (
	// ErrConflict is wrapped for an operation with a stale base that
	// overlaps an operation applied after its base. The operation applied.
	ErrConflict = errors.New("model: conflict")
	// ErrMissing is wrapped for a target, parent, endpoint or decision that
	// is not in the model, and for an operation of an unknown type.
	ErrMissing = errors.New("model: not found")
	// ErrBase is wrapped for a base that is not in the target's history.
	ErrBase = errors.New("model: base not in history")
	// ErrExists is wrapped for a create of an ID already used by a node or a
	// link.
	ErrExists = errors.New("model: already exists")
	// ErrRetired is wrapped for an operation on a retired target from its
	// current version, and for a create or move under a retired parent.
	ErrRetired = errors.New("model: retired")
	// ErrParent is wrapped for a create or move under a parent whose kind
	// has no children, and for a move of a node whose kind is not in the
	// tree.
	ErrParent = errors.New("model: no children here")
	// ErrCycle is wrapped for a move that would make a node its own
	// ancestor.
	ErrCycle = errors.New("model: cycle")
	// ErrHandle is wrapped for a live node whose handle another live node
	// already holds in the same place.
	ErrHandle = errors.New("model: duplicate handle")
	// ErrState is wrapped for a change set whose state does not allow the
	// view asked for.
	ErrState = errors.New("model: wrong state")
)

// Ref names one operation of one change set.
type Ref struct {
	ChangeSet id.ID
	Op        id.ID
}

// Node is a node as it stands in a view. Empty Fields and History are nil.
type Node struct {
	ID     id.ID
	Kind   op.Kind
	Handle string
	// Parent is zero for a root and for a kind outside the tree; Key is ""
	// for the latter.
	Parent  id.ID
	Key     string
	Fields  map[string]json.RawMessage
	Retired bool
	Reason  string
	// Version is the last operation applied to the node, Created the change
	// set of its create and History every applied operation, the create
	// first.
	Version id.ID
	Created id.ID
	History []Ref
}

// Link is a link as it stands in a view. From and To are nodes or change
// sets.
type Link struct {
	ID       id.ID
	Type     op.LinkType
	From     id.ID
	To       id.ID
	Note     string
	Decision id.ID
	Retired  bool
	Version  id.ID
	Created  id.ID
	History  []Ref
}

// record is a node or a link with the aspects each history entry wrote,
// parallel to its History, so that stale bases need no operations.
type record[T any] struct {
	v       T
	aspects [][]string
}

// Model is one view of the model. The zero Model is empty. Its accessors
// return copies, so nothing a caller does to them changes the Model.
type Model struct {
	nodes map[id.ID]*record[Node]
	links map[id.ID]*record[Link]
	// children maps a parent, or the zero ID for the roots, to the nodes
	// under it whose kind sits in the tree, live or retired.
	children map[id.ID]map[id.ID]struct{}
}

// Node returns the node with ID i.
func (m Model) Node(i id.ID) (Node, bool) {
	r, ok := m.nodes[i]
	if !ok {
		return Node{}, false
	}
	return copyNode(r.v), true
}

// Link returns the link with ID i.
func (m Model) Link(i id.ID) (Link, bool) {
	r, ok := m.links[i]
	if !ok {
		return Link{}, false
	}
	return copyLink(r.v), true
}

// Nodes returns every node, live and retired, sorted by ID; nil when there
// are none.
func (m Model) Nodes() []Node {
	var out []Node
	for _, r := range m.nodes {
		out = append(out, copyNode(r.v))
	}
	slices.SortFunc(out, func(a, b Node) int { return id.Compare(a.ID, b.ID) })
	return out
}

// Links returns every link, live and retired, sorted by ID; nil when there
// are none.
func (m Model) Links() []Link {
	var out []Link
	for _, r := range m.links {
		out = append(out, copyLink(r.v))
	}
	slices.SortFunc(out, func(a, b Link) int { return id.Compare(a.ID, b.ID) })
	return out
}

// Children returns the nodes directly under parent, live and retired, sorted
// by key then ID. The zero parent gives the roots of the tree: nodes of a
// kind that sits in the tree, with no parent. It is nil when there are none.
func (m Model) Children(parent id.ID) []Node {
	var out []Node
	for c := range m.children[parent] {
		out = append(out, copyNode(m.nodes[c].v))
	}
	slices.SortFunc(out, func(a, b Node) int {
		return cmp.Or(strings.Compare(a.Key, b.Key), id.Compare(a.ID, b.ID))
	})
	return out
}

// clone returns a copy of m that shares nothing with it.
func (m Model) clone() Model {
	out := Model{
		nodes:    make(map[id.ID]*record[Node], len(m.nodes)),
		links:    make(map[id.ID]*record[Link], len(m.links)),
		children: make(map[id.ID]map[id.ID]struct{}, len(m.children)),
	}
	for k, r := range m.nodes {
		out.nodes[k] = &record[Node]{v: copyNode(r.v), aspects: slices.Clone(r.aspects)}
	}
	for k, r := range m.links {
		out.links[k] = &record[Link]{v: copyLink(r.v), aspects: slices.Clone(r.aspects)}
	}
	for k, cs := range m.children {
		out.children[k] = make(map[id.ID]struct{}, len(cs))
		for c := range cs {
			out.children[k][c] = struct{}{}
		}
	}
	return out
}

// copyFields deep-copies a field map, returning nil when it is empty.
func copyFields(m map[string]json.RawMessage) map[string]json.RawMessage {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		out[k] = json.RawMessage(bytes.Clone(v))
	}
	return out
}

func copyNode(n Node) Node {
	n.Fields = copyFields(n.Fields)
	n.History = slices.Clone(n.History)
	return n
}

func copyLink(l Link) Link {
	l.History = slices.Clone(l.History)
	return l
}

// Problem is one thing wrong with applying a change set's operations. Op is
// zero for a problem with the change set as a whole. Err wraps exactly one
// of the package's sentinels and is never nil.
type Problem struct {
	ChangeSet id.ID
	Op        id.ID
	Err       error
}

func (p Problem) String() string {
	return fmt.Sprintf("%s/%s: %v", p.ChangeSet, p.Op, p.Err)
}

// sortProblems sorts by change set, op, then error text, and returns nil
// when there are none.
func sortProblems(ps []Problem) []Problem {
	if len(ps) == 0 {
		return nil
	}
	slices.SortFunc(ps, func(a, b Problem) int {
		return cmp.Or(id.Compare(a.ChangeSet, b.ChangeSet), id.Compare(a.Op, b.Op),
			strings.Compare(a.Err.Error(), b.Err.Error()))
	})
	return ps
}

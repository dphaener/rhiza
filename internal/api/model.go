package api

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// Model is what the write operations need to know about the graph: refs,
// base versions, siblings, handles and actors. It sees the write view: the
// as-is graph (every landed change set), then the in-flight change sets the
// current one depends on, directly or through others, then the current
// change set's ops on top.
//
// The implementation behind it, newModel, is interim until task 4's
// projection exists; then swap the constructor and keep this interface.
// Nothing outside internal/api may depend on the implementation.
type Model interface {
	// Resolve finds a node by ref (see the ref rules on newModel). Retired
	// nodes resolve by ID and ID prefix, never by handle. No match is a
	// not_found Error with next "genesis ls"; several are a usage Error
	// listing the candidates.
	Resolve(ref string) (Node, error)
	// Node returns the node with ID n, retired or not.
	Node(n id.ID) (Node, bool)
	// ResolveLink finds a link by its full ID or a unique ID prefix of at
	// least four characters, retired or not.
	ResolveLink(ref string) (Link, error)
	// LinkBetween finds the live link of type t from from to to. No match
	// is a not_found Error; several are a usage Error listing them.
	LinkBetween(from, to id.ID, t op.LinkType) (Link, error)
	// Children returns the live nodes of kind whose parent is parent, sorted
	// by key, then ID. A zero parent gives the roots of that kind.
	Children(kind op.Kind, parent id.ID) []Node
	// HandleTaken reports which live node has handle h, if any.
	HandleTaken(h string) (id.ID, bool)
	// ActorByEmail finds the live actor node whose email field is the JSON
	// string email.
	ActorByEmail(email string) (Node, bool)
}

// Node is a node in the to-be view.
type Node struct {
	ID     id.ID
	Kind   op.Kind
	Handle string
	// Parent is zero for a root and for kinds without a tree; Key is empty
	// for the latter.
	Parent id.ID
	Key    string
	Fields map[string]json.RawMessage
	// Retired is true once a node.retire applies, landed or current.
	Retired bool
	// Version is the head over the landed ops and the dependencies' ops on
	// the node: the op that no other of them names as base. It is zero for a
	// node created only in the current change set, and the current change
	// set never changes it.
	Version id.ID
	// Conflict is true when the node has more than one such head. Version
	// is then the largest head ID.
	Conflict bool
	// Pending lists the IDs of the current change set's ops on the node, in
	// change set order; nil when it has none. A dependency's ops are never
	// in it.
	Pending []id.ID
}

// Link is a link in the to-be view. Version, Conflict and Pending mean what
// they mean on Node.
type Link struct {
	ID       id.ID
	Type     op.LinkType
	From     id.ID
	To       id.ID
	Note     string
	Decision id.ID
	Retired  bool
	Version  id.ID
	Conflict bool
	Pending  []id.ID
}

// model is the interim Model: two maps rebuilt from the whole log each time.
type model struct {
	nodes map[id.ID]*Node
	links map[id.ID]*Link
}

// newModel builds the interim Model from l and the current change set (zero
// for none).
//
// As-is is the ops of every change set whose folded state is landed.
// Proposed and accepted change sets are not in it. The write view adds the
// change sets in changeset.Closure of the current one (its in-flight
// dependencies, direct or not), then the current change set's ops in their
// listed order. A terminal current change set adds neither (a landed one
// is already in as-is). The dependencies apply after the landed ops, in
// dependencyOrder, each one's ops in their listed order.
//
// Landed ops on one target apply along their base chain: the create (no
// base) first, then each op after the op it names as base, siblings in op ID
// order. With two landed heads both branches apply, the later op ID last, so
// its fields win; Version is the largest head and Conflict is set. Ops whose
// base is not a landed op on the same target apply last, in op ID order.
// Version and Conflict count the dependencies' ops too: a node a dependency
// created has that create as its version, and a dependency's op on a landed
// node moves its version on.
//
// A ref is, in this order: the full canonical ID of a node; a unique prefix,
// four characters or more, of the node IDs, matched as given (IDs are upper
// case, so a lower-case ref never matches one); the handle of a live node.
//
// Log and fold problems are ignored here: a node with no create, or an entry
// that does not decode, is simply left out. The check command reports them.
func newModel(l store.Log, current id.ID) Model {
	m := &model{nodes: map[id.ID]*Node{}, links: map[id.ID]*Link{}}
	sets, _ := changeset.FoldLog(l)
	var cur *changeset.ChangeSet
	for i, cs := range sets {
		if cs.ID == current && !current.IsZero() && cs.State != "" && !cs.State.Terminal() {
			cur = &sets[i]
		}
	}

	byTarget := map[id.ID][]op.Op{}
	for _, cs := range sets {
		if cs.State != changeset.StateLanded {
			continue
		}
		for _, o := range cs.Ops {
			byTarget[o.Target()] = append(byTarget[o.Target()], o)
		}
	}
	for _, t := range slices.SortedFunc(maps.Keys(byTarget), id.Compare) {
		for _, o := range chain(byTarget[t]) {
			m.apply(o)
		}
	}
	if cur != nil {
		byID := map[id.ID]changeset.ChangeSet{}
		for _, cs := range sets {
			byID[cs.ID] = cs
		}
		for _, d := range dependencyOrder(sets, current) {
			for _, o := range byID[d].Ops {
				m.apply(o)
				byTarget[o.Target()] = append(byTarget[o.Target()], o)
			}
		}
	}
	for t, ops := range byTarget {
		version, conflict := head(ops)
		if n := m.nodes[t]; n != nil {
			n.Version, n.Conflict = version, conflict
		}
		if k := m.links[t]; k != nil {
			k.Version, k.Conflict = version, conflict
		}
	}

	if cur != nil {
		for _, o := range cur.Ops {
			m.apply(o)
			if n := m.nodes[o.Target()]; n != nil {
				n.Pending = append(n.Pending, o.ID)
			}
			if k := m.links[o.Target()]; k != nil {
				k.Pending = append(k.Pending, o.ID)
			}
		}
	}
	return m
}

// dependencyOrder returns changeset.Closure(sets, current) with every change
// set after the ones it depends on among them, ties by ID. In a loop no
// member is first by that rule, so the smallest ID left goes next.
func dependencyOrder(sets []changeset.ChangeSet, current id.ID) []id.ID {
	left := changeset.Closure(sets, current)
	deps := map[id.ID][]id.ID{}
	for _, cs := range left {
		deps[cs] = changeset.DependsOn(sets, cs)
	}
	var out []id.ID
	for len(left) > 0 {
		i := slices.IndexFunc(left, func(cs id.ID) bool {
			return !slices.ContainsFunc(deps[cs], func(d id.ID) bool {
				return slices.Contains(left, d)
			})
		})
		if i < 0 {
			i = 0
		}
		out = append(out, left[i])
		left = slices.Delete(left, i, i+1)
	}
	return out
}

// chain orders one target's landed ops for application: depth first along
// base links from the ops with no landed base, siblings in op ID order.
func chain(ops []op.Op) []op.Op {
	ids := map[id.ID]bool{}
	for _, o := range ops {
		ids[o.ID] = true
	}
	next := map[id.ID][]op.Op{}
	var roots, orphans []op.Op
	for _, o := range ops {
		switch {
		case o.Base.IsZero():
			roots = append(roots, o)
		case ids[o.Base] && o.Base != o.ID:
			next[o.Base] = append(next[o.Base], o)
		default:
			orphans = append(orphans, o)
		}
	}
	byID := func(a, b op.Op) int { return id.Compare(a.ID, b.ID) }
	var out []op.Op
	seen := map[id.ID]bool{}
	var walk func(level []op.Op)
	walk = func(level []op.Op) {
		slices.SortFunc(level, byID)
		for _, o := range level {
			if seen[o.ID] {
				continue
			}
			seen[o.ID] = true
			out = append(out, o)
			walk(next[o.ID])
		}
	}
	walk(roots)
	slices.SortFunc(orphans, byID)
	for _, o := range orphans {
		if !seen[o.ID] {
			seen[o.ID] = true
			out = append(out, o)
			walk(next[o.ID])
		}
	}
	return out
}

// head returns the largest op ID that no other op names as base, and
// whether there is more than one such op.
func head(ops []op.Op) (id.ID, bool) {
	named := map[id.ID]bool{}
	for _, o := range ops {
		if o.Base != o.ID {
			named[o.Base] = true
		}
	}
	var heads []id.ID
	for _, o := range ops {
		if !named[o.ID] {
			heads = append(heads, o.ID)
		}
	}
	if len(heads) == 0 {
		return id.ID{}, false
	}
	return slices.MaxFunc(heads, id.Compare), len(heads) > 1
}

// apply applies o to the node or link it targets. Ops on a target with no
// create are ignored.
func (m *model) apply(o op.Op) {
	switch o.Type {
	case op.TypeNodeCreate:
		n := m.nodes[o.Node]
		if n == nil {
			n = &Node{ID: o.Node}
			m.nodes[o.Node] = n
		}
		n.Kind, n.Handle, n.Parent, n.Key = o.Kind, o.Handle, o.Parent, o.Key
		n.Fields = maps.Clone(o.Fields)
		if n.Fields == nil {
			n.Fields = map[string]json.RawMessage{}
		}
	case op.TypeNodeUpdate:
		if n := m.nodes[o.Node]; n != nil {
			if o.Handle != "" {
				n.Handle = o.Handle
			}
			maps.Copy(n.Fields, o.Set)
			for _, k := range o.Unset {
				delete(n.Fields, k)
			}
		}
	case op.TypeNodeMove:
		if n := m.nodes[o.Node]; n != nil {
			n.Parent, n.Key = o.Parent, o.Key
		}
	case op.TypeNodeRetire:
		if n := m.nodes[o.Node]; n != nil {
			n.Retired = true
		}
	case op.TypeLinkCreate:
		k := m.links[o.Link]
		if k == nil {
			k = &Link{ID: o.Link}
			m.links[o.Link] = k
		}
		k.Type, k.From, k.To, k.Note, k.Decision = o.LinkType, o.From, o.To, o.Note, o.Decision
	case op.TypeLinkRetire:
		if k := m.links[o.Link]; k != nil {
			k.Retired = true
		}
	}
}

// copyNode returns n with its map and slice copied, so callers cannot change
// the model.
func copyNode(n *Node) Node {
	c := *n
	c.Fields = maps.Clone(n.Fields)
	c.Pending = slices.Clone(n.Pending)
	return c
}

func copyLink(k *Link) Link {
	c := *k
	c.Pending = slices.Clone(k.Pending)
	return c
}

// sortedNodes returns the nodes that keep says to, sorted by ID.
func (m *model) sortedNodes(keep func(*Node) bool) []Node {
	var out []Node
	for _, n := range m.nodes {
		if keep(n) {
			out = append(out, copyNode(n))
		}
	}
	slices.SortFunc(out, func(a, b Node) int { return id.Compare(a.ID, b.ID) })
	return out
}

func (m *model) sortedLinks(keep func(*Link) bool) []Link {
	var out []Link
	for _, k := range m.links {
		if keep(k) {
			out = append(out, copyLink(k))
		}
	}
	slices.SortFunc(out, func(a, b Link) int { return id.Compare(a.ID, b.ID) })
	return out
}

func (m *model) Resolve(ref string) (Node, error) {
	if v, err := id.Parse(ref); err == nil {
		if n, ok := m.nodes[v]; ok {
			return copyNode(n), nil
		}
	}
	if len(ref) >= minPrefix {
		matches := m.sortedNodes(func(n *Node) bool { return strings.HasPrefix(n.ID.String(), ref) })
		if len(matches) == 1 {
			return matches[0], nil
		}
		if len(matches) > 1 {
			return Node{}, ambiguousNodes(ref, matches)
		}
	}
	matches := m.sortedNodes(func(n *Node) bool { return !n.Retired && n.Handle == ref })
	switch len(matches) {
	case 0:
		return Node{}, NotFoundf("no node matches %q", ref).WithNext(lsNext)
	case 1:
		return matches[0], nil
	}
	return Node{}, ambiguousNodes(ref, matches)
}

func ambiguousNodes(ref string, ns []Node) error {
	names := make([]string, len(ns))
	for i, n := range ns {
		names[i] = n.ID.String() + " (" + n.Handle + ")"
	}
	return Usagef("%q matches %d nodes: %s", ref, len(ns), strings.Join(names, ", "))
}

func ambiguousLinks(what string, ls []Link) error {
	names := make([]string, len(ls))
	for i, k := range ls {
		names[i] = k.ID.String()
	}
	return Usagef("%s matches %d links: %s", what, len(ls), strings.Join(names, ", "))
}

func (m *model) Node(n id.ID) (Node, bool) {
	if v, ok := m.nodes[n]; ok {
		return copyNode(v), true
	}
	return Node{}, false
}

func (m *model) ResolveLink(ref string) (Link, error) {
	if v, err := id.Parse(ref); err == nil {
		if k, ok := m.links[v]; ok {
			return copyLink(k), nil
		}
	}
	if len(ref) >= minPrefix {
		matches := m.sortedLinks(func(k *Link) bool { return strings.HasPrefix(k.ID.String(), ref) })
		if len(matches) == 1 {
			return matches[0], nil
		}
		if len(matches) > 1 {
			return Link{}, ambiguousLinks(strconv.Quote(ref), matches)
		}
	}
	return Link{}, NotFoundf("no link matches %q", ref)
}

func (m *model) LinkBetween(from, to id.ID, t op.LinkType) (Link, error) {
	matches := m.sortedLinks(func(k *Link) bool {
		return !k.Retired && k.From == from && k.To == to && k.Type == t
	})
	switch len(matches) {
	case 0:
		return Link{}, NotFoundf("no %s link from %s to %s", t, from, to)
	case 1:
		return matches[0], nil
	}
	return Link{}, ambiguousLinks(fmt.Sprintf("%s from %s to %s", t, from, to), matches)
}

func (m *model) Children(kind op.Kind, parent id.ID) []Node {
	out := m.sortedNodes(func(n *Node) bool { return !n.Retired && n.Kind == kind && n.Parent == parent })
	slices.SortStableFunc(out, func(a, b Node) int { return cmp.Compare(a.Key, b.Key) })
	return out
}

func (m *model) HandleTaken(h string) (id.ID, bool) {
	matches := m.sortedNodes(func(n *Node) bool { return !n.Retired && n.Handle == h })
	if len(matches) == 0 {
		return id.ID{}, false
	}
	return matches[0].ID, true
}

func (m *model) ActorByEmail(email string) (Node, bool) {
	matches := m.sortedNodes(func(n *Node) bool {
		var s string
		return !n.Retired && n.Kind == op.KindActor &&
			json.Unmarshal(n.Fields["email"], &s) == nil && s == email
	})
	if len(matches) == 0 {
		return Node{}, false
	}
	return matches[0], true
}

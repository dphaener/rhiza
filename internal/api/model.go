package api

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
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
	// not_found Error with next "rhiza ls"; several are a usage Error
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
	// LinksTo returns the live links of type t whose To is to, sorted by
	// ID. It is nil when there are none.
	LinksTo(to id.ID, t op.LinkType) []Link
	// LinksFrom returns the live links of type t whose From is from, sorted
	// by ID. It is nil when there are none.
	LinksFrom(from id.ID, t op.LinkType) []Link
	// DependsOn returns what n depends on: the To of each live depends-on
	// link from n, sorted by ID, each once. It is nil when there are none.
	DependsOn(n id.ID) []id.ID
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
	// Version is the last operation in the node's strict history that is
	// not the current change set's own: the one new writes name as base. It
	// is zero for a node created only in the current change set, and the
	// current change set never changes it.
	Version id.ID
	// Conflict is true when the view holds an unrepaired conflict
	// (model.ErrConflict) on the node from an operation not of the current
	// change set: a stale write that overlaps one it did not see.
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

// model is the interim Model: two maps taken from one of internal/model's
// views of the whole log.
type model struct {
	nodes map[id.ID]*Node
	links map[id.ID]*Link
}

// newModel builds the interim Model from l and the current change set (zero
// for none).
//
// The nodes and links are internal/model's (the strict engine's): the
// as-is view (model.AsIs, every landed change set) when there is no
// current change set or it is not in flight, else the current change set's
// what-if view (model.WhatIf: as-is, then the current change set and the
// in-flight change sets it depends on, directly or through others, in
// their apply order). An op the strict engine skips (a stale or unknown
// base, a retired or missing target, parent or endpoint, a taken ID) is
// left out, so a node whose create no longer applies is not in the view.
//
// Version and Conflict follow the strict rule: Version is the last op in
// the target's strict history that is not the current change set's own,
// and Conflict is set when the view holds an unrepaired model.ErrConflict
// problem on the target from an op not of the current change set. Pending
// lists the current change set's ops on each target the view holds,
// whether the strict engine applied them or not.
//
// A ref is, in this order: the full canonical ID of a node; a unique prefix,
// four characters or more, of the node IDs, matched as given (IDs are upper
// case, so a lower-case ref never matches one); the handle of a live node.
//
// Log and fold problems are ignored here: a node with no create, or an entry
// that does not decode, is simply left out. The check command reports them.
func newModel(l store.Log, current id.ID) Model {
	sets, _ := changeset.FoldLog(l)
	return newModelSets(sets, current)
}

// newModelSets is newModel over the log already folded: sets are
// changeset.FoldLog's change sets, which it only reads. A caller that builds
// several models of one log (the daemon's, one per candidate) builds a
// modelBase once instead, so as-is is built once.
func newModelSets(sets []changeset.ChangeSet, current id.ID) Model {
	return newModelBase(sets).model(current)
}

// modelBase is the as-is view of one folded log, built once, that the
// interim models of its change sets go on from (model.AsIsSeed, then
// model.WhatIfFrom per change set).
type modelBase struct {
	sets     []changeset.ChangeSet
	asIs     gmodel.Model
	problems []gmodel.Problem
	seed     gmodel.Seed
	// target is each op's target, by op ID.
	target map[id.ID]id.ID
	// asIsModel is the interim as-is Model, built on first use. Nothing
	// changes a model once built, so every as-is view shares it and every
	// what-if starts from a copy of its maps.
	asIsModel *model
}

// newModelBase builds the as-is view of sets, which it only reads.
func newModelBase(sets []changeset.ChangeSet) *modelBase {
	b := &modelBase{sets: sets, target: map[id.ID]id.ID{}}
	b.asIs, b.problems, b.seed = gmodel.AsIsSeed(sets)
	for _, cs := range sets {
		for _, o := range cs.Ops {
			b.target[o.ID] = o.Target()
		}
	}
	return b
}

// model is newModelSets(b's sets, current).
//
// A what-if differs from as-is only on the targets of the ops it layers:
// those of the current change set and of its in-flight dependencies. Its
// other nodes and links, and their versions and conflicts, are as-is's,
// so only those targets are taken from the what-if view.
func (b *modelBase) model(current id.ID) Model {
	var cur *changeset.ChangeSet
	for i, cs := range b.sets {
		if cs.ID == current && !current.IsZero() && inFlightState(cs.State) {
			cur = &b.sets[i]
		}
	}
	if b.asIsModel == nil {
		b.asIsModel = &model{nodes: map[id.ID]*Node{}, links: map[id.ID]*Link{}}
		v := strictVersions{conflict: b.conflicts(b.problems, id.ID{})}
		for _, n := range b.asIs.Nodes() {
			b.asIsModel.nodes[n.ID] = v.node(n)
		}
		for _, l := range b.asIs.Links() {
			b.asIsModel.links[l.ID] = v.link(l)
		}
	}
	if cur == nil {
		return b.asIsModel
	}
	w, ps := gmodel.WhatIfFrom(b.seed, b.sets, current)
	v := strictVersions{current: current, conflict: b.conflicts(ps, current)}
	touched := map[id.ID]bool{}
	byID := map[id.ID]*changeset.ChangeSet{}
	for i := range b.sets {
		byID[b.sets[i].ID] = &b.sets[i]
	}
	for _, d := range append(changeset.Closure(b.sets, current), current) {
		if s := byID[d]; s != nil {
			for _, o := range s.Ops {
				touched[o.Target()] = true
			}
		}
	}
	m := &model{nodes: maps.Clone(b.asIsModel.nodes), links: maps.Clone(b.asIsModel.links)}
	for t := range touched {
		if n, ok := w.Node(t); ok {
			m.nodes[t] = v.node(n)
		}
		if l, ok := w.Link(t); ok {
			m.links[t] = v.link(l)
		}
	}
	for _, o := range cur.Ops {
		if n := m.nodes[o.Target()]; n != nil {
			n.Pending = append(n.Pending, o.ID)
		}
		if k := m.links[o.Target()]; k != nil {
			k.Pending = append(k.Pending, o.ID)
		}
	}
	return m
}

// conflicts returns the targets ps holds an unrepaired model.ErrConflict
// problem on from an op not of change set current (zero for none).
func (b *modelBase) conflicts(ps []gmodel.Problem, current id.ID) map[id.ID]bool {
	out := map[id.ID]bool{}
	for _, p := range ps {
		if !errors.Is(p.Err, gmodel.ErrConflict) || (!current.IsZero() && p.ChangeSet == current) {
			continue
		}
		var stale *gmodel.StaleBaseError
		if errors.As(p.Err, &stale) {
			out[stale.Node] = true
		} else if t, ok := b.target[p.Op]; ok {
			out[t] = true
		}
	}
	return out
}

// strictVersions turns the strict engine's nodes and links into the
// interim model's, with Version and Conflict by the strict rule for the
// view of change set current (zero for as-is).
type strictVersions struct {
	current  id.ID
	conflict map[id.ID]bool
}

// version is the last op of h that is not the current change set's own.
func (v strictVersions) version(h []gmodel.Ref) id.ID {
	for i := len(h) - 1; i >= 0; i-- {
		if v.current.IsZero() || h[i].ChangeSet != v.current {
			return h[i].Op
		}
	}
	return id.ID{}
}

func (v strictVersions) node(n gmodel.Node) *Node {
	a := apiNode(n)
	if a.Fields == nil {
		a.Fields = map[string]json.RawMessage{}
	}
	a.Version, a.Conflict = v.version(n.History), v.conflict[n.ID]
	return &a
}

func (v strictVersions) link(l gmodel.Link) *Link {
	return &Link{ID: l.ID, Type: l.Type, From: l.From, To: l.To, Note: l.Note, Decision: l.Decision,
		Retired: l.Retired, Version: v.version(l.History), Conflict: v.conflict[l.ID]}
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

func (m *model) LinksTo(to id.ID, t op.LinkType) []Link {
	return m.sortedLinks(func(k *Link) bool { return !k.Retired && k.To == to && k.Type == t })
}

func (m *model) LinksFrom(from id.ID, t op.LinkType) []Link {
	return m.sortedLinks(func(k *Link) bool { return !k.Retired && k.From == from && k.Type == t })
}

func (m *model) DependsOn(n id.ID) []id.ID {
	var out []id.ID
	for _, k := range m.sortedLinks(func(k *Link) bool {
		return !k.Retired && k.From == n && k.Type == op.LinkDependsOn
	}) {
		out = append(out, k.To)
	}
	slices.SortFunc(out, id.Compare)
	return slices.Compact(out)
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

package model

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// The aspects an operation writes, as described at the top of model.go.
const (
	aspectField    = "field:"
	aspectHandle   = "handle"
	aspectPosition = "position"
	aspectRetired  = "retired"
)

// aspects returns the aspects o writes, sorted and unique; nil for a create.
func aspects(o op.Op) []string {
	var out []string
	switch o.Type {
	case op.TypeNodeUpdate:
		for k := range o.Set {
			out = append(out, aspectField+k)
		}
		for _, k := range o.Unset {
			out = append(out, aspectField+k)
		}
		if o.Handle != "" {
			out = append(out, aspectHandle)
		}
	case op.TypeNodeMove:
		out = append(out, aspectPosition)
	case op.TypeNodeRetire, op.TypeLinkRetire:
		out = append(out, aspectRetired)
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// overlaps reports whether two operations' aspects overlap: they share one,
// or either is a retirement and the other writes anything.
func overlaps(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if slices.Contains(a, aspectRetired) || slices.Contains(b, aspectRetired) {
		return true
	}
	return slices.ContainsFunc(a, func(s string) bool { return slices.Contains(b, s) })
}

// mark is the operation that last set a node's handle or parent, and when it
// was applied. Duplicate handles are reported on it.
type mark struct {
	seq int
	ref Ref
}

// held is a conflict that later writes may still repair: the index of its
// problem and the aspects of its operation not yet rewritten.
type held struct {
	problem int
	left    []string
}

// builder applies operations one at a time to a model. A view is a new
// builder that knows the input's change set IDs, apply for every operation of
// every change set in apply order, then finish. finish does not change the
// builder, so more operations can be applied after it, and clone forks it.
type builder struct {
	m Model
	// known are the change set IDs a link endpoint may name.
	known map[id.ID]struct{}
	// seq counts applied and skipped operations, in apply order.
	seq      int
	placed   map[id.ID]mark
	problems []Problem
	// held are the conflicts each target's later writes may still repair,
	// and repaired the indexes in problems of those they did.
	held     map[id.ID][]held
	repaired map[int]struct{}
	// src, when set, holds the as-is rows b.m lacks, which node and link
	// load on first use (see WhatIfDelta). srcErr is the first error it
	// gave.
	src    Source
	srcErr error
}

// newBuilder returns an empty builder whose links may name the change sets
// in known.
func newBuilder(known []id.ID) *builder {
	b := &builder{
		m: Model{
			nodes:    map[id.ID]*record[Node]{},
			links:    map[id.ID]*record[Link]{},
			children: map[id.ID]map[id.ID]struct{}{},
		},
		known:    make(map[id.ID]struct{}, len(known)),
		placed:   map[id.ID]mark{},
		held:     map[id.ID][]held{},
		repaired: map[int]struct{}{},
	}
	for _, k := range known {
		b.known[k] = struct{}{}
	}
	return b
}

// clone returns a copy of b that shares nothing with it.
func (b *builder) clone() *builder {
	out := &builder{
		m:        b.m.clone(),
		known:    maps.Clone(b.known),
		seq:      b.seq,
		placed:   maps.Clone(b.placed),
		problems: slices.Clone(b.problems),
		held:     make(map[id.ID][]held, len(b.held)),
		repaired: maps.Clone(b.repaired),
	}
	for t, hs := range b.held {
		hs = slices.Clone(hs)
		for i := range hs {
			hs[i].left = slices.Clone(hs[i].left)
		}
		out.held[t] = hs
	}
	return out
}

// finish reports duplicate handles (step 5) and returns a copy of the model
// with every problem so far but the repaired conflicts, sorted, nil when
// there are none. The builder itself is unchanged.
func (b *builder) finish() (Model, []Problem) {
	return b.finishModel(true)
}

// finishModel is finish; with clone false the model returned is b's own,
// which b must not be used for again.
func (b *builder) finishModel(clone bool) (Model, []Problem) {
	type place struct {
		parent id.ID
		kind   op.Kind
		handle string
	}
	groups := map[place][]id.ID{}
	for k, r := range b.m.nodes {
		n := r.v
		if n.Retired {
			continue
		}
		p := place{parent: n.Parent, handle: n.Handle}
		if !op.HasParent(n.Kind) {
			p = place{kind: n.Kind, handle: n.Handle}
		}
		groups[p] = append(groups[p], k)
	}
	var ps []Problem
	for i, p := range b.problems {
		if _, ok := b.repaired[i]; !ok {
			ps = append(ps, p)
		}
	}
	for p, ns := range groups {
		if len(ns) < 2 {
			continue
		}
		slices.SortFunc(ns, func(x, y id.ID) int {
			return cmp.Or(cmp.Compare(b.placed[x].seq, b.placed[y].seq), id.Compare(x, y))
		})
		for _, n := range ns[1:] {
			ref := b.placed[n].ref
			ps = append(ps, Problem{ChangeSet: ref.ChangeSet, Op: ref.Op,
				Err: fmt.Errorf("%w: node %s has handle %q, already held by node %s", ErrHandle, n, p.handle, ns[0])})
		}
	}
	if !clone {
		return b.m, sortProblems(ps)
	}
	return b.m.clone(), sortProblems(ps)
}

// apply applies operation o of change set cs, as described at the top of
// model.go, and records any problem.
func (b *builder) apply(cs id.ID, o op.Op) {
	b.seq++
	err := b.applyOp(cs, o)
	if err == nil {
		return
	}
	// A retirement's conflict cannot be repaired: its target takes no more
	// writes.
	if w := aspects(o); errors.Is(err, ErrConflict) && !slices.Contains(w, aspectRetired) {
		t := o.Node
		if o.Type == op.TypeLinkRetire {
			t = o.Link
		}
		b.held[t] = append(b.held[t], held{problem: len(b.problems), left: w})
	}
	b.problems = append(b.problems, Problem{ChangeSet: cs, Op: o.ID, Err: err})
}

// repair takes the aspects an operation from t's current version wrote off
// every conflict held on t, and marks repaired those with none left.
func (b *builder) repair(t id.ID, wrote []string) {
	var kept []held
	for _, h := range b.held[t] {
		h.left = slices.DeleteFunc(h.left, func(a string) bool { return slices.Contains(wrote, a) })
		if len(h.left) == 0 {
			b.repaired[h.problem] = struct{}{}
			continue
		}
		kept = append(kept, h)
	}
	if kept == nil {
		delete(b.held, t)
		return
	}
	b.held[t] = kept
}

// applyOp returns why o was skipped, ErrConflict when it applied with a
// conflict, or nil.
func (b *builder) applyOp(cs id.ID, o op.Op) error {
	ref := Ref{ChangeSet: cs, Op: o.ID}
	switch o.Type {
	case op.TypeNodeCreate:
		return b.createNode(ref, o)
	case op.TypeLinkCreate:
		return b.createLink(ref, o)

	case op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire:
		r, ok := b.node(o.Node)
		if !ok {
			return fmt.Errorf("%w: node %s", ErrMissing, o.Node)
		}
		conflict, err := checkBase(cs, o, r.v.ID, r.v.Retired, r.v.History, r.aspects)
		if err != nil {
			return err
		}
		if o.Type == op.TypeNodeMove {
			if err := b.checkMove(r.v, o); err != nil {
				return err
			}
		}
		current := o.Base == r.v.Version
		b.effect(r, ref, o)
		r.v.History = append(r.v.History, ref)
		r.aspects = append(r.aspects, aspects(o))
		r.v.Version = o.ID
		if current {
			b.repair(o.Node, aspects(o))
		}
		return conflict

	case op.TypeLinkRetire:
		r, ok := b.link(o.Link)
		if !ok {
			return fmt.Errorf("%w: link %s", ErrMissing, o.Link)
		}
		conflict, err := checkBase(cs, o, r.v.ID, r.v.Retired, r.v.History, r.aspects)
		if err != nil {
			return err
		}
		r.v.Retired = true
		r.v.History = append(r.v.History, ref)
		r.aspects = append(r.aspects, aspects(o))
		r.v.Version = o.ID
		return conflict
	}
	return fmt.Errorf("%w: unknown operation type %q", ErrMissing, o.Type)
}

// checkBase is step 2 for target t with the given history: it returns the
// error that skips o, or else ErrConflict when o's base is stale and o
// overlaps an operation of another change set applied after it. Operations
// of cs itself, o's own change set, never make a base stale.
func checkBase(cs id.ID, o op.Op, t id.ID, retired bool, history []Ref, written [][]string) (conflict, skip error) {
	i := len(history) - 1
	for i >= 0 && history[i].Op != o.Base {
		i--
	}
	switch {
	case i < 0:
		return nil, fmt.Errorf("%w: base %s is not an operation on %s", ErrBase, o.Base, t)
	case i == len(history)-1:
		if retired {
			return nil, fmt.Errorf("%w: %s", ErrRetired, t)
		}
		return nil, nil
	}
	mine := aspects(o)
	for j := i + 1; j < len(history); j++ {
		if history[j].ChangeSet == cs {
			// An operation is never stale against its own change set's.
			continue
		}
		if overlaps(mine, written[j]) {
			return fmt.Errorf("%w: base %s of %s is stale: op %s of change set %s since wrote %s",
				ErrConflict, o.Base, t, history[j].Op, history[j].ChangeSet, strings.Join(written[j], ", ")), nil
		}
	}
	if retired {
		// Only the change set's own operations came after the base, so a
		// retirement among them is its own: a retired target takes no writes.
		return nil, fmt.Errorf("%w: %s", ErrRetired, t)
	}
	return nil, nil
}

func (b *builder) exists(i id.ID) bool {
	_, n := b.node(i)
	_, l := b.link(i)
	return n || l
}

// checkParent checks that p can take a child: it is a live node of a kind
// with children.
func (b *builder) checkParent(p id.ID) error {
	r, ok := b.node(p)
	switch {
	case !ok:
		return fmt.Errorf("%w: parent %s", ErrMissing, p)
	case r.v.Retired:
		return fmt.Errorf("%w: parent %s", ErrRetired, p)
	case !op.HasParent(r.v.Kind):
		return fmt.Errorf("%w: parent %s is a %s", ErrParent, p, r.v.Kind)
	}
	return nil
}

func (b *builder) createNode(ref Ref, o op.Op) error {
	if b.exists(o.Node) {
		return fmt.Errorf("%w: %s", ErrExists, o.Node)
	}
	if !o.Parent.IsZero() {
		if err := b.checkParent(o.Parent); err != nil {
			return err
		}
	}
	b.m.nodes[o.Node] = &record[Node]{
		v: Node{ID: o.Node, Kind: o.Kind, Handle: o.Handle, Parent: o.Parent, Key: o.Key, Fields: copyFields(o.Fields),
			Version: o.ID, Created: ref.ChangeSet, History: []Ref{ref}},
		aspects: [][]string{nil},
	}
	if op.HasParent(o.Kind) {
		b.addChild(o.Parent, o.Node)
	}
	b.placed[o.Node] = mark{seq: b.seq, ref: ref}
	return nil
}

func (b *builder) createLink(ref Ref, o op.Op) error {
	if b.exists(o.Link) {
		return fmt.Errorf("%w: %s", ErrExists, o.Link)
	}
	for _, e := range []struct {
		name string
		id   id.ID
	}{{"from", o.From}, {"to", o.To}} {
		_, node := b.node(e.id)
		_, cs := b.known[e.id]
		if !node && !cs {
			return fmt.Errorf("%w: %s %s is neither a node nor a known change set", ErrMissing, e.name, e.id)
		}
	}
	if _, ok := b.node(o.Decision); !o.Decision.IsZero() && !ok {
		return fmt.Errorf("%w: decision %s", ErrMissing, o.Decision)
	}
	b.m.links[o.Link] = &record[Link]{
		v: Link{ID: o.Link, Type: o.LinkType, From: o.From, To: o.To, Note: o.Note, Decision: o.Decision,
			Version: o.ID, Created: ref.ChangeSet, History: []Ref{ref}},
		aspects: [][]string{nil},
	}
	return nil
}

// checkMove checks a move of n: n sits in the tree, and its new parent can
// take a child and is neither n nor below it. Only the new parent's
// ancestors are walked, never the whole model.
func (b *builder) checkMove(n Node, o op.Op) error {
	if !op.HasParent(n.Kind) {
		return fmt.Errorf("%w: %s is a %s, which is not in the tree", ErrParent, n.ID, n.Kind)
	}
	if o.Parent.IsZero() {
		return nil
	}
	if err := b.checkParent(o.Parent); err != nil {
		return err
	}
	// Moves are checked, so parents never form a cycle; the bound is only a
	// guard.
	for p, steps := o.Parent, 0; !p.IsZero() && steps <= len(b.m.nodes); steps++ {
		if p == n.ID {
			return fmt.Errorf("%w: %s is %s or below it", ErrCycle, o.Parent, n.ID)
		}
		r, ok := b.node(p)
		if !ok {
			break
		}
		p = r.v.Parent
	}
	return nil
}

// effect is step 4 for an operation on node record r that passed its checks.
func (b *builder) effect(r *record[Node], ref Ref, o op.Op) {
	n := &r.v
	switch o.Type {
	case op.TypeNodeUpdate:
		if len(o.Set) > 0 && n.Fields == nil {
			n.Fields = make(map[string]json.RawMessage, len(o.Set))
		}
		for k, v := range o.Set {
			n.Fields[k] = json.RawMessage(bytes.Clone(v))
		}
		for _, k := range o.Unset {
			delete(n.Fields, k)
		}
		if len(n.Fields) == 0 {
			n.Fields = nil
		}
		if o.Handle != "" {
			n.Handle = o.Handle
			b.placed[n.ID] = mark{seq: b.seq, ref: ref}
		}
	case op.TypeNodeMove:
		b.removeChild(n.Parent, n.ID)
		n.Parent, n.Key = o.Parent, o.Key
		b.addChild(n.Parent, n.ID)
		b.placed[n.ID] = mark{seq: b.seq, ref: ref}
	case op.TypeNodeRetire:
		n.Retired, n.Reason = true, o.Reason
	}
}

func (b *builder) addChild(parent, child id.ID) {
	cs, ok := b.m.children[parent]
	if !ok {
		cs = map[id.ID]struct{}{}
		b.m.children[parent] = cs
	}
	cs[child] = struct{}{}
}

// removeChild drops child from parent, and drops parent's set once empty so
// that equal models compare equal.
func (b *builder) removeChild(parent, child id.ID) {
	delete(b.m.children[parent], child)
	if len(b.m.children[parent]) == 0 {
		delete(b.m.children, parent)
	}
}

// node returns the record of node i, loading it from the source on first
// use when b has one.
func (b *builder) node(i id.ID) (*record[Node], bool) {
	if r, ok := b.m.nodes[i]; ok {
		return r, true
	}
	if b.src == nil || b.srcErr != nil || i.IsZero() {
		return nil, false
	}
	sn, ok, err := b.src.SeedNode(i)
	if err != nil {
		b.srcErr = err
		return nil, false
	}
	if !ok {
		return nil, false
	}
	r := &record[Node]{v: sn.Node, aspects: sn.Aspects}
	b.m.nodes[i] = r
	if op.HasParent(sn.Node.Kind) {
		b.addChild(sn.Node.Parent, i)
	}
	b.placed[i] = mark{seq: sn.Placed.Seq, ref: sn.Placed.Ref}
	return r, true
}

// link returns the record of link i, loading it from the source on first
// use when b has one.
func (b *builder) link(i id.ID) (*record[Link], bool) {
	if r, ok := b.m.links[i]; ok {
		return r, true
	}
	if b.src == nil || b.srcErr != nil || i.IsZero() {
		return nil, false
	}
	sl, ok, err := b.src.SeedLink(i)
	if err != nil {
		b.srcErr = err
		return nil, false
	}
	if !ok {
		return nil, false
	}
	r := &record[Link]{v: sl.Link, aspects: sl.Aspects}
	b.m.links[i] = r
	return r, true
}

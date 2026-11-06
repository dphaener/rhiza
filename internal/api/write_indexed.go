package api

// This file holds the write path over the index: a write validates against
// the index's folded state plus its own new entries, so it never reads or
// folds the log. The Model a builder sees is the what-if overlay of the
// target change set (the stored as-is rows, corrected by what the change
// set and its in-flight dependencies wrote), with the change set taken from
// the entries the dry run holds in pending when there are any.

import (
	"cmp"
	"encoding/json"
	"errors"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// writeIndex is the index of one write. A write to disk makes it stale; the
// next read opens it again, which catches it up with the entry written.
type writeIndex struct {
	ctx   *Ctx
	x     *index.Index
	stale bool
	err   error
	// models, when set (a driver's reads, see run_model_cache.go), keeps
	// the model of each change set until the caller clears it.
	models map[id.ID]Model
	// relevant, kept with models, is the change sets each model reads: the
	// change set and its in-flight dependencies.
	relevant map[id.ID][]id.ID
}

// openWriteIndex opens the index of ctx's repository for a write.
func openWriteIndex(ctx *Ctx) (*writeIndex, error) {
	x, err := index.Open(ctx.Root, ctx.IndexPath)
	if err != nil {
		return nil, Errorf("%s", err.Error())
	}
	return &writeIndex{ctx: ctx, x: x}, nil
}

func (ix *writeIndex) close() {
	if ix != nil && ix.x != nil {
		ix.x.Close()
	}
}

// failure is the first index error a read met, nil for none. Reads answer
// nothing after one, and the write stops at it.
func (ix *writeIndex) failure() error {
	if ix == nil || ix.err == nil {
		return nil
	}
	return Errorf("%s", ix.err.Error())
}

func (ix *writeIndex) fail(err error) {
	if err != nil && ix.err == nil {
		ix.err = err
	}
}

// index returns the open index, opened again first if a write made it stale.
func (ix *writeIndex) index() *index.Index {
	if ix.stale {
		ix.x.Close()
		x, err := index.Open(ix.ctx.Root, ix.ctx.IndexPath)
		if err != nil {
			ix.fail(err)
			return ix.x
		}
		ix.x, ix.stale = x, false
	}
	return ix.x
}

// entries returns change set cs's stored entries and the dry run's pending
// ones, sorted by ID.
func (ix *writeIndex) entries(w *workspace, cs id.ID) []entry.Entry {
	var es []entry.Entry
	if !cs.IsZero() {
		var err error
		if es, err = ix.index().Entries(cs); err != nil {
			ix.fail(err)
		}
	}
	for _, e := range w.pending {
		if e.ChangeSet == cs {
			es = append(es, e)
		}
	}
	slices.SortFunc(es, func(a, b entry.Entry) int { return id.Compare(a.ID, b.ID) })
	return es
}

// setIDs returns the IDs of the change sets in the index and in pending,
// sorted.
func (ix *writeIndex) setIDs(w *workspace) []id.ID {
	infos, err := ix.index().ChangeSets()
	ix.fail(err)
	ids := infoIDs(infos)
	for _, e := range w.pending {
		if !slices.Contains(ids, e.ChangeSet) {
			ids = append(ids, e.ChangeSet)
		}
	}
	slices.SortFunc(ids, id.Compare)
	return ids
}

// changeSets folds every change set: the stored ones, with those pending
// entries touch folded again. Only a depends-on link between change sets
// asks for it.
func (ix *writeIndex) changeSets(w *workspace) []changeset.ChangeSet {
	ids := ix.setIDs(w)
	got, err := ix.index().FoldedOf(ids)
	ix.fail(err)
	touched := map[id.ID]bool{}
	for _, e := range w.pending {
		touched[e.ChangeSet] = true
	}
	var out []changeset.ChangeSet
	for _, cs := range ids {
		if touched[cs] {
			c, _ := changeset.Fold(cs, ix.entries(w, cs))
			out = append(out, c)
		} else if c, ok := got[cs]; ok {
			out = append(out, c)
		}
	}
	return out
}

// model is the write view of change set current (zero for none), as
// newModel builds it: the stored as-is rows, then the ops of current's
// in-flight dependencies in dependencyOrder, then current's own, each
// applied as newModel's model.apply does. The strict what-if view is not
// it: that one drops a node whose create no longer applies (its parent's
// create gone with a retirement in the change set, or a move under a node
// created later), which the interim model keeps.
func (ix *writeIndex) model(w *workspace, current id.ID) Model {
	if ix.models != nil && len(w.pending) == 0 {
		if m, ok := ix.models[current]; ok {
			return m
		}
		m := ix.buildModel(w, current)
		if ix.err == nil {
			ix.models[current] = m
		}
		return m
	}
	return ix.buildModel(w, current)
}

// buildModel is model without the memo.
func (ix *writeIndex) buildModel(w *workspace, current id.ID) Model {
	if ix.ctx != nil && ix.ctx.builtModel != nil {
		ix.ctx.builtModel(current)
	}
	m := &writeModel{ix: ix, cs: current}
	var over *changeset.ChangeSet
	if !current.IsZero() {
		if es := ix.entries(w, current); len(es) > 0 {
			c, _ := changeset.Fold(current, es)
			if inFlightState(c.State) {
				m.ops = c.Ops
				m.pend = map[id.ID][]id.ID{}
				for _, o := range c.Ops {
					m.pend[o.Target()] = append(m.pend[o.Target()], o.ID)
				}
			}
			over = &c
		}
	}
	if ix.relevant != nil {
		ix.relevant[current] = []id.ID{current}
	}
	l := &writeLayer{ix: ix, nodes: map[id.ID]*gmodel.Node{}, links: map[id.ID]*gmodel.Link{}}
	if over != nil && inFlightState(over.State) {
		sets, err := ix.index().WhatIfSets(current, over)
		ix.fail(err)
		byID := map[id.ID]changeset.ChangeSet{}
		for _, c := range sets {
			byID[c.ID] = c
		}
		order := dependencyOrder(sets, current)
		if ix.relevant != nil {
			ix.relevant[current] = append([]id.ID{current}, order...)
		}
		for _, d := range order {
			for _, o := range byID[d].Ops {
				l.apply(d, o)
			}
		}
		for _, o := range over.Ops {
			l.apply(current, o)
		}
	}
	m.ov = ix.index().OverlayOf(l.list())
	return m
}

// writeLayer is the nodes and links in-flight ops touch, each loaded from
// the stored as-is rows on first use, with the ops applied on top.
type writeLayer struct {
	ix    *writeIndex
	nodes map[id.ID]*gmodel.Node
	links map[id.ID]*gmodel.Link
}

// node returns node i of the layer, loaded from as-is on first use; nil
// when neither holds it.
func (l *writeLayer) node(i id.ID) *gmodel.Node {
	if n, ok := l.nodes[i]; ok {
		return n
	}
	n, ok, err := l.ix.index().Node(index.AsIs, i)
	l.ix.fail(err)
	if !ok || err != nil {
		return nil
	}
	n.Fields, n.History = maps.Clone(n.Fields), slices.Clone(n.History)
	l.nodes[i] = &n
	return &n
}

// link is node for links.
func (l *writeLayer) link(i id.ID) *gmodel.Link {
	if k, ok := l.links[i]; ok {
		return k
	}
	k, ok, err := l.ix.index().Link(index.AsIs, i)
	l.ix.fail(err)
	if !ok || err != nil {
		return nil
	}
	k.History = slices.Clone(k.History)
	l.links[i] = &k
	return &k
}

// apply applies op o of change set cs as model.apply does: a create sets
// its target whatever it was, and an op on a target with no create is
// ignored. Each op applied joins its target's History, which versionOf
// reads.
func (l *writeLayer) apply(cs id.ID, o op.Op) {
	ref := gmodel.Ref{ChangeSet: cs, Op: o.ID}
	switch o.Type {
	case op.TypeNodeCreate:
		n := l.node(o.Node)
		if n == nil {
			n = &gmodel.Node{ID: o.Node, Created: cs}
			l.nodes[o.Node] = n
		}
		n.Kind, n.Handle, n.Parent, n.Key = o.Kind, o.Handle, o.Parent, o.Key
		n.Fields = maps.Clone(o.Fields)
		n.History = append(n.History, ref)
	case op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire:
		n := l.node(o.Node)
		if n == nil {
			return
		}
		switch o.Type {
		case op.TypeNodeUpdate:
			if o.Handle != "" {
				n.Handle = o.Handle
			}
			if n.Fields == nil && len(o.Set) > 0 {
				n.Fields = map[string]json.RawMessage{}
			}
			maps.Copy(n.Fields, o.Set)
			for _, k := range o.Unset {
				delete(n.Fields, k)
			}
		case op.TypeNodeMove:
			n.Parent, n.Key = o.Parent, o.Key
		default:
			n.Retired = true
		}
		n.History = append(n.History, ref)
	case op.TypeLinkCreate:
		k := l.link(o.Link)
		if k == nil {
			k = &gmodel.Link{ID: o.Link, Created: cs}
			l.links[o.Link] = k
		}
		k.Type, k.From, k.To, k.Note, k.Decision = o.LinkType, o.From, o.To, o.Note, o.Decision
		k.History = append(k.History, ref)
	case op.TypeLinkRetire:
		if k := l.link(o.Link); k != nil {
			k.Retired = true
			k.History = append(k.History, ref)
		}
	}
}

// list returns the layer's nodes and links.
func (l *writeLayer) list() ([]gmodel.Node, []gmodel.Link) {
	ns := make([]gmodel.Node, 0, len(l.nodes))
	for _, n := range l.nodes {
		ns = append(ns, *n)
	}
	ls := make([]gmodel.Link, 0, len(l.links))
	for _, k := range l.links {
		ls = append(ls, *k)
	}
	return ns, ls
}

// writeModel is the Model of a write over the index.
type writeModel struct {
	ix   *writeIndex
	ov   *index.Overlay
	cs   id.ID
	ops  []op.Op
	pend map[id.ID][]id.ID
}

func (m *writeModel) ok() bool { return m.ov != nil && m.ix.err == nil }

// versionOf gives the version and conflict of a node or link from its
// history: the head over the operations other than the change set's own,
// as the interim model counts them.
func (m *writeModel) versionOf(hist []gmodel.Ref) (id.ID, bool) {
	var refs []gmodel.Ref
	for _, h := range hist {
		if h.ChangeSet != m.cs || m.cs.IsZero() {
			refs = append(refs, h)
		}
	}
	switch len(refs) {
	case 0:
		return id.ID{}, false
	case 1:
		return refs[0].Op, false
	}
	var css []id.ID
	want := map[id.ID]bool{}
	for _, r := range refs {
		want[r.Op] = true
		if !slices.Contains(css, r.ChangeSet) {
			css = append(css, r.ChangeSet)
		}
	}
	got, err := m.ix.index().FoldedOf(css)
	if err != nil {
		m.ix.fail(err)
		return id.ID{}, false
	}
	var ops []op.Op
	for _, c := range got {
		for _, o := range c.Ops {
			if want[o.ID] {
				ops = append(ops, o)
			}
		}
	}
	return head(ops)
}

// creator is the actor that created node n (gmodel.CreatedBy): the one
// stamped on the entry that first wrote its create, read from the change
// set that holds the create (History[0]) alone.
func (m *writeModel) creator(n id.ID) (id.ID, bool) {
	if !m.ok() {
		return id.ID{}, false
	}
	g, ok, err := m.ov.Node(n)
	m.ix.fail(err)
	if !ok || err != nil || len(g.History) == 0 {
		return id.ID{}, false
	}
	cs := g.History[0].ChangeSet
	got, err := m.ix.index().FoldedOf([]id.ID{cs})
	if err != nil {
		m.ix.fail(err)
		return id.ID{}, false
	}
	c, ok := got[cs]
	if !ok {
		return id.ID{}, false
	}
	return gmodel.CreatedBy(g, []changeset.ChangeSet{c})
}

func (m *writeModel) node(n gmodel.Node, version bool) Node {
	out := apiNode(n)
	out.Pending = slices.Clone(m.pend[n.ID])
	if version {
		out.Version, out.Conflict = m.versionOf(n.History)
	}
	return out
}

func (m *writeModel) link(l gmodel.Link, version bool) Link {
	out := Link{ID: l.ID, Type: l.Type, From: l.From, To: l.To, Note: l.Note, Decision: l.Decision, Retired: l.Retired}
	out.Pending = slices.Clone(m.pend[l.ID])
	if version {
		out.Version, out.Conflict = m.versionOf(l.History)
	}
	return out
}

// queryError maps an index error to an Error as readError does, and keeps
// any other for the write to stop at.
func (m *writeModel) queryError(err error) error {
	if errors.Is(err, index.ErrNotFound) || errors.Is(err, index.ErrAmbiguous) {
		return readError(err)
	}
	m.ix.fail(err)
	return Errorf("%s", err.Error())
}

func (m *writeModel) Resolve(ref string) (Node, error) {
	if !m.ok() {
		return Node{}, NotFoundf("no node matches %q", ref).WithNext(lsNext)
	}
	n, err := m.ov.Resolve(ref)
	if err != nil {
		return Node{}, m.queryError(err)
	}
	return m.node(n, true), nil
}

func (m *writeModel) Node(i id.ID) (Node, bool) {
	if !m.ok() {
		return Node{}, false
	}
	n, ok, err := m.ov.Node(i)
	m.ix.fail(err)
	if !ok || err != nil {
		return Node{}, false
	}
	return m.node(n, false), true
}

func (m *writeModel) ResolveLink(ref string) (Link, error) {
	if !m.ok() {
		return Link{}, NotFoundf("no link matches %q", ref)
	}
	l, err := m.ov.ResolveLink(ref)
	switch {
	case errors.Is(err, index.ErrNotFound):
		return Link{}, NotFoundf("no link matches %q", ref)
	case errors.Is(err, index.ErrAmbiguous):
		return Link{}, Usagef("%s", strings.TrimPrefix(err.Error(), index.ErrAmbiguous.Error()+": "))
	case err != nil:
		return Link{}, m.queryError(err)
	}
	return m.link(l, true), nil
}

// links gives the live links of type t at node i in direction dir, sorted
// by ID.
func (m *writeModel) links(i id.ID, dir index.Direction, t op.LinkType) []gmodel.Link {
	if !m.ok() {
		return nil
	}
	var ls []gmodel.Link
	var err error
	if dir == index.To {
		ls, err = m.ov.LinksTo([]id.ID{i}, t)
	} else {
		ls, err = m.ov.Links(i, dir)
	}
	m.ix.fail(err)
	return slices.DeleteFunc(ls, func(l gmodel.Link) bool { return l.Retired || l.Type != t })
}

func (m *writeModel) LinkBetween(from, to id.ID, t op.LinkType) (Link, error) {
	var matches []Link
	for _, l := range m.links(from, index.From, t) {
		if l.To == to {
			matches = append(matches, m.link(l, true))
		}
	}
	switch len(matches) {
	case 0:
		return Link{}, NotFoundf("no %s link from %s to %s", t, from, to)
	case 1:
		return matches[0], nil
	}
	return Link{}, ambiguousLinks(strings.Join([]string{string(t), "from", from.String(), "to", to.String()}, " "), matches)
}

func (m *writeModel) LinksTo(to id.ID, t op.LinkType) []Link {
	return m.linkList(m.links(to, index.To, t))
}

func (m *writeModel) LinksFrom(from id.ID, t op.LinkType) []Link {
	return m.linkList(m.links(from, index.From, t))
}

func (m *writeModel) linkList(ls []gmodel.Link) []Link {
	var out []Link
	for _, l := range ls {
		out = append(out, m.link(l, false))
	}
	return out
}

func (m *writeModel) DependsOn(n id.ID) []id.ID {
	var out []id.ID
	for _, l := range m.links(n, index.From, op.LinkDependsOn) {
		out = append(out, l.To)
	}
	slices.SortFunc(out, id.Compare)
	return slices.Compact(out)
}

func (m *writeModel) Children(kind op.Kind, parent id.ID) []Node {
	if !m.ok() {
		return nil
	}
	ns, err := m.ov.Children(parent)
	m.ix.fail(err)
	var out []Node
	for _, n := range ns {
		if !n.Retired && n.Kind == kind {
			out = append(out, m.node(n, false))
		}
	}
	slices.SortStableFunc(out, func(a, b Node) int { return cmp.Compare(a.Key, b.Key) })
	return out
}

func (m *writeModel) HandleTaken(h string) (id.ID, bool) {
	if !m.ok() {
		return id.ID{}, false
	}
	n, ok, err := m.ov.HandleTaken(h)
	m.ix.fail(err)
	return n, ok && err == nil
}

func (m *writeModel) ActorByEmail(email string) (Node, bool) {
	if !m.ok() {
		return Node{}, false
	}
	ns, err := m.ov.List(index.ListQuery{Kinds: []op.Kind{op.KindActor}, Limit: 1 << 20})
	m.ix.fail(err)
	for _, n := range ns {
		var s string
		if !n.Retired && json.Unmarshal(n.Fields["email"], &s) == nil && s == email {
			return m.node(n, false), true
		}
	}
	return Node{}, false
}

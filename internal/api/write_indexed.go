package api

// This file holds the write path over the index: a write validates against
// the index's folded state plus its own new entries, so it never reads or
// folds the log. The Model a builder sees is the what-if overlay of the
// target change set (the stored as-is rows, corrected by what the change
// set and its in-flight dependencies wrote, as the strict engine applies
// them), with the change set taken from the entries the dry run holds in
// pending when there are any.

import (
	"cmp"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// writeIndex is the projection (projection.go): the index of one call with
// the dry run's pending entries laid over it. A write to disk makes it
// stale; the next read opens it again, which catches it up with the entry
// written.
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
	if ctx.openedProjection != nil {
		ctx.openedProjection()
	}
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

// wrote marks the index stale after a write to disk.
func (ix *writeIndex) wrote() { ix.stale = true }

// problems returns the problems the index stores.
func (ix *writeIndex) problems() []index.Problem {
	ps, err := ix.index().Problems()
	ix.fail(err)
	return ps
}

// asIs returns the as-is model the index assembles from its stored rows.
func (ix *writeIndex) asIs() (gmodel.Model, []gmodel.Problem) {
	m, ps, err := ix.index().AsIs()
	ix.fail(err)
	return m, ps
}

func (ix *writeIndex) fail(err error) {
	if err != nil && ix.err == nil {
		ix.err = err
	}
}

// store is the open index for the reads the projection's other methods
// do not cover.
func (ix *writeIndex) store() *index.Index { return ix.index() }

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
func (ix *writeIndex) entries(pending []entry.Entry, cs id.ID) []entry.Entry {
	var es []entry.Entry
	if !cs.IsZero() {
		var err error
		if es, err = ix.index().Entries(cs); err != nil {
			ix.fail(err)
		}
	}
	for _, e := range pending {
		if e.ChangeSet == cs {
			es = append(es, e)
		}
	}
	slices.SortFunc(es, func(a, b entry.Entry) int { return id.Compare(a.ID, b.ID) })
	return es
}

// changeSetIDs returns the IDs of the change sets in the index and in pending,
// sorted.
func (ix *writeIndex) changeSetIDs(pending []entry.Entry) []id.ID {
	infos, err := ix.index().ChangeSets()
	ix.fail(err)
	ids := infoIDs(infos)
	for _, e := range pending {
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
func (ix *writeIndex) changeSets(pending []entry.Entry) []changeset.ChangeSet {
	ids := ix.changeSetIDs(pending)
	got, err := ix.index().FoldedOf(ids)
	ix.fail(err)
	touched := map[id.ID]bool{}
	for _, e := range pending {
		touched[e.ChangeSet] = true
	}
	var out []changeset.ChangeSet
	for _, cs := range ids {
		if touched[cs] {
			c, _ := changeset.Fold(cs, ix.entries(pending, cs))
			out = append(out, c)
		} else if c, ok := got[cs]; ok {
			out = append(out, c)
		}
	}
	return out
}

// changeSetInfos returns every change set's stored row, sorted by ID, with
// the change sets pending entries touch folded again from their stored and
// pending entries.
func (ix *writeIndex) changeSetInfos(pending []entry.Entry) []index.ChangeSetInfo {
	infos, err := ix.index().ChangeSets()
	ix.fail(err)
	if len(pending) == 0 {
		return infos
	}
	touched := map[id.ID]bool{}
	for _, e := range pending {
		touched[e.ChangeSet] = true
	}
	out := slices.DeleteFunc(infos, func(c index.ChangeSetInfo) bool { return touched[c.ID] })
	for cs := range touched {
		c, _ := changeset.Fold(cs, ix.entries(pending, cs))
		out = append(out, index.ChangeSetInfo{ID: cs, State: c.State, Kind: c.Open.Kind, Intent: c.Open.Intent,
			OpenEntry: c.OpenEntry, SupersededBy: c.SupersededBy, TerminalEntry: c.TerminalEntry})
	}
	slices.SortFunc(out, func(a, b index.ChangeSetInfo) int { return id.Compare(a.ID, b.ID) })
	return out
}

// model is the write view of change set current (zero for none), as
// newModel builds it: the strict engine's, taken from the index. When
// current is in flight it is the strict what-if of current (the stored
// as-is rows, then current and its in-flight dependencies in their apply
// order, an op the strict engine skips left out), with current folded from
// its stored entries and the dry run's pending ones; else it is as-is.
// Version and Conflict follow the strict rule from the view's histories
// and problems, as newModel's do.
func (ix *writeIndex) model(pending []entry.Entry, current id.ID) Model {
	if ix.models != nil && len(pending) == 0 {
		if m, ok := ix.models[current]; ok {
			return m
		}
		m := ix.buildModel(pending, current)
		if ix.err == nil {
			ix.models[current] = m
		}
		return m
	}
	return ix.buildModel(pending, current)
}

// buildModel is model without the memo.
func (ix *writeIndex) buildModel(pending []entry.Entry, current id.ID) Model {
	if ix.ctx != nil && ix.ctx.builtModel != nil {
		ix.ctx.builtModel(current)
	}
	m := &writeModel{ix: ix}
	var over *changeset.ChangeSet
	if !current.IsZero() {
		if es := ix.entries(pending, current); len(es) > 0 {
			if c, _ := changeset.Fold(current, es); inFlightState(c.State) {
				over = &c
			}
		}
	}
	if ix.relevant != nil {
		ix.relevant[current] = []id.ID{current}
	}
	view := id.ID{}
	if over != nil {
		view, m.cs = current, current
		m.pend = map[id.ID][]id.ID{}
		for _, o := range over.Ops {
			m.pend[o.Target()] = append(m.pend[o.Target()], o.ID)
		}
		if ix.relevant != nil {
			sets, err := ix.index().WhatIfSets(current, over)
			ix.fail(err)
			ix.relevant[current] = append([]id.ID{current}, dependencyOrder(sets, current)...)
		}
	}
	ov, ps, err := ix.index().WhatIfOverlayOver(view, over)
	if err != nil {
		ix.fail(err)
		return m
	}
	m.ov = ov
	m.conflict = ix.conflicts(ps, m.cs)
	return m
}

// conflicts returns the targets ps holds an unrepaired model.ErrConflict
// problem on from an op not of change set current (zero for none), as
// modelBase.conflicts does. A stored conflict that is not a stale base
// names only its op, whose target is read from its change set.
func (ix *writeIndex) conflicts(ps []gmodel.Problem, current id.ID) map[id.ID]bool {
	out := map[id.ID]bool{}
	var rest []gmodel.Problem
	for _, p := range ps {
		if !errors.Is(p.Err, gmodel.ErrConflict) || (!current.IsZero() && p.ChangeSet == current) {
			continue
		}
		var stale *gmodel.StaleBaseError
		if errors.As(p.Err, &stale) {
			out[stale.Node] = true
		} else if !p.Op.IsZero() {
			rest = append(rest, p)
		}
	}
	if len(rest) == 0 {
		return out
	}
	var css []id.ID
	for _, p := range rest {
		if !slices.Contains(css, p.ChangeSet) {
			css = append(css, p.ChangeSet)
		}
	}
	got, err := ix.index().OpsOf(css)
	ix.fail(err)
	for _, p := range rest {
		for _, o := range got[p.ChangeSet].Ops {
			if o.ID == p.Op {
				out[o.Target()] = true
			}
		}
	}
	return out
}

// writeModel is the Model of a write over the index.
type writeModel struct {
	ix *writeIndex
	ov *index.Overlay
	// cs is the current change set: the view's while it is in flight, else
	// none.
	cs   id.ID
	pend map[id.ID][]id.ID
	// conflict is the targets the view holds an unrepaired conflict on
	// from an op not of cs.
	conflict map[id.ID]bool
}

func (m *writeModel) ok() bool { return m.ov != nil && m.ix.err == nil }

// versionOf gives the version and conflict of node or link i from its
// strict history: the last operation that is not the current change set's
// own, and whether the view holds an unrepaired conflict on it.
func (m *writeModel) versionOf(i id.ID, hist []gmodel.Ref) (id.ID, bool) {
	for j := len(hist) - 1; j >= 0; j-- {
		if m.cs.IsZero() || hist[j].ChangeSet != m.cs {
			return hist[j].Op, m.conflict[i]
		}
	}
	return id.ID{}, m.conflict[i]
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
		out.Version, out.Conflict = m.versionOf(n.ID, n.History)
	}
	return out
}

func (m *writeModel) link(l gmodel.Link, version bool) Link {
	out := Link{ID: l.ID, Type: l.Type, From: l.From, To: l.To, Note: l.Note, Decision: l.Decision, Retired: l.Retired}
	out.Pending = slices.Clone(m.pend[l.ID])
	if version {
		out.Version, out.Conflict = m.versionOf(l.ID, l.History)
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

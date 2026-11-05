package api

// This file holds change show over the index: the same answers as the
// log's fold, read from the stored change sets, entries and rows, so the
// log is not read and a change set's own view is the what-if overlay, not
// a model built from every entry.

import (
	"cmp"
	"encoding/json"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// indexModel is a change set's view as the index gives it: the what-if
// overlay for a change set in flight, the as-is rows for any other. It
// answers what the interim Model answers about the same change set. An
// index error stops the answers (nothing found) and is kept in err.
type indexModel struct {
	x   *index.Index
	ov  *index.Overlay
	err error
}

// newIndexModel is the view of change set cs, in state state.
func newIndexModel(x *index.Index, cs id.ID, state changeset.State) (*indexModel, error) {
	m := &indexModel{x: x}
	if inFlightState(state) {
		ov, _, err := x.WhatIfOverlay(cs)
		if err != nil {
			return nil, Errorf("%s", err.Error())
		}
		m.ov = ov
	}
	return m, nil
}

func (m *indexModel) fail(err error) {
	if err != nil && m.err == nil {
		m.err = err
	}
}

func (m *indexModel) gnode(i id.ID) (gmodel.Node, bool) {
	if m.err != nil {
		return gmodel.Node{}, false
	}
	var n gmodel.Node
	var ok bool
	var err error
	if m.ov != nil {
		n, ok, err = m.ov.Node(i)
	} else {
		n, ok, err = m.x.Node(index.AsIs, i)
	}
	m.fail(err)
	return n, ok && err == nil
}

func (m *indexModel) glink(i id.ID) (gmodel.Link, bool) {
	if m.err != nil {
		return gmodel.Link{}, false
	}
	var l gmodel.Link
	var ok bool
	var err error
	if m.ov != nil {
		l, ok, err = m.ov.Link(i)
	} else {
		l, ok, err = m.x.Link(index.AsIs, i)
	}
	m.fail(err)
	return l, ok && err == nil
}

func apiNode(n gmodel.Node) Node {
	return Node{ID: n.ID, Kind: n.Kind, Handle: n.Handle, Parent: n.Parent, Key: n.Key, Fields: n.Fields, Retired: n.Retired}
}

func (m *indexModel) Node(i id.ID) (Node, bool) {
	n, ok := m.gnode(i)
	return apiNode(n), ok
}

func (m *indexModel) Children(kind op.Kind, parent id.ID) []Node {
	if m.err != nil {
		return nil
	}
	var ns []gmodel.Node
	var err error
	if m.ov != nil {
		ns, err = m.ov.Children(parent)
	} else {
		ns, err = m.x.Children(index.AsIs, parent)
	}
	m.fail(err)
	var out []Node
	for _, n := range ns {
		if !n.Retired && n.Kind == kind {
			out = append(out, apiNode(n))
		}
	}
	return out
}

// links gives the live links of type t at node i in direction dir.
func (m *indexModel) links(i id.ID, dir index.Direction, t op.LinkType) []Link {
	if m.err != nil {
		return nil
	}
	var ls []gmodel.Link
	var err error
	switch {
	case m.ov != nil && dir == index.To:
		ls, err = m.ov.LinksTo([]id.ID{i}, t)
	case m.ov != nil:
		ls, err = m.ov.Links(i, dir)
	case dir == index.To:
		ls, err = m.x.LinksTo(index.AsIs, []id.ID{i}, t)
	default:
		ls, err = m.x.Links(index.AsIs, i, dir)
	}
	m.fail(err)
	var out []Link
	for _, l := range ls {
		if !l.Retired && l.Type == t {
			out = append(out, Link{ID: l.ID, Type: l.Type, From: l.From, To: l.To, Note: l.Note, Decision: l.Decision})
		}
	}
	return out
}

func (m *indexModel) LinksFrom(from id.ID, t op.LinkType) []Link {
	return m.links(from, index.From, t)
}

func (m *indexModel) LinksTo(to id.ID, t op.LinkType) []Link {
	return m.links(to, index.To, t)
}

func (m *indexModel) DependsOn(n id.ID) []id.ID {
	var out []id.ID
	for _, l := range m.links(n, index.From, op.LinkDependsOn) {
		out = append(out, l.To)
	}
	slices.SortFunc(out, id.Compare)
	return slices.Compact(out)
}

// ancestors is the parent chain of n, the root first.
func (m *indexModel) ancestors(n id.ID) []gmodel.Node {
	if m.err != nil {
		return nil
	}
	var as []gmodel.Node
	var err error
	if m.ov != nil {
		as, err = m.ov.Ancestors(n)
	} else {
		as, err = m.x.Ancestors(index.AsIs, n)
	}
	m.fail(err)
	return as
}

// rank gives each of leaves that a depth-first walk of the capability tree
// reaches (live capabilities all the way up) its place among them, by
// comparing their paths from the root as the walk orders siblings: by key,
// then ID. A leaf has no children, so no path is the start of another.
func (m *indexModel) rank(leaves []Node) map[id.ID]int {
	type path struct {
		leaf  id.ID
		chain []gmodel.Node
	}
	var ps []path
	for _, l := range leaves {
		n, ok := m.gnode(l.ID)
		if !ok {
			continue
		}
		chain := append(m.ancestors(l.ID), n)
		if !chain[0].Parent.IsZero() || slices.ContainsFunc(chain, func(c gmodel.Node) bool {
			return c.Retired || c.Kind != op.KindCapability
		}) {
			continue
		}
		ps = append(ps, path{l.ID, chain})
	}
	slices.SortFunc(ps, func(a, b path) int {
		for i := 0; i < len(a.chain) && i < len(b.chain); i++ {
			if c := cmp.Or(cmp.Compare(a.chain[i].Key, b.chain[i].Key), id.Compare(a.chain[i].ID, b.chain[i].ID)); c != 0 {
				return c
			}
		}
		return cmp.Compare(len(a.chain), len(b.chain))
	})
	rank := map[id.ID]int{}
	for i, p := range ps {
		rank[p.leaf] = i
	}
	return rank
}

// runChangeShowIndexed is runChangeShow over the index.
func runChangeShowIndexed(ctx *Ctx, w *workspace, r ChangeShowRequest) (any, error) {
	x, err := index.Open(ctx.Root, ctx.IndexPath)
	if err != nil {
		return nil, Errorf("%s", err.Error())
	}
	defer x.Close()
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	infos, err := x.ChangeSets()
	if err != nil {
		return nil, Errorf("%s", err.Error())
	}
	var cs id.ID
	if ref := strings.TrimSpace(r.Ref); ref != "" {
		cs, err = matchChangeSet(infoIDs(infos), ref)
	} else {
		cs, err = indexedEffective(ctx, current, infos)
	}
	if err != nil {
		return nil, err
	}
	if cs.IsZero() {
		return nil, Usagef("no change set: name one or make one current").WithNext(changeLsNext)
	}
	es, err := x.Entries(cs)
	if err != nil {
		return nil, Errorf("%s", err.Error())
	}
	if len(es) == 0 {
		return nil, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	c, problems := changeset.Fold(cs, es)
	// The diff's landed change sets do not depend on the view, so they
	// load while the view is built.
	landed := prefetchLanded(x, c)
	m, err := newIndexModel(x, cs, c.State)
	if err != nil {
		<-landed
		return nil, err
	}
	serves := []id.ID{}
	for _, l := range m.LinksFrom(cs, op.LinkServes) {
		serves = append(serves, l.To)
	}
	pre := <-landed
	if pre.err != nil {
		return nil, Errorf("%s", pre.err.Error())
	}
	diff, err := indexedDiff(x, m, c, pre.sets)
	if err != nil {
		return nil, err
	}
	notes, err := indexedNotes(x, infos, cs, current, m)
	if err != nil {
		return nil, err
	}
	resp := ChangeShowResponse{
		ID:           cs,
		Kind:         c.Open.Kind,
		Intent:       c.Open.Intent,
		Title:        c.Open.DisplayTitle(),
		Sources:      orEmpty(c.Open.Sources),
		Serves:       serves,
		State:        c.State,
		Current:      cs == current,
		Ops:          []json.RawMessage{},
		Leaves:       changeLeavesIn(m, c.Ops, m.rank),
		SupersededBy: c.SupersededBy,
		Heads:        orEmpty(c.Heads),
		Applied:      orEmpty(c.Applied),
		History:      []LifecycleEvent{},
		Problems:     foldProblems(problems),
		Changes:      buildChangeDiff(cs, diff, modelHandleOf(m)),
		Notes:        notes,
	}
	for _, e := range c.Lifecycle {
		resp.History = append(resp.History, LifecycleEvent{Type: e.Type, Entry: e.Entry,
			At: e.At.UTC().Format(lineageAtLayout), Actor: e.Actor})
	}
	for _, o := range c.Ops {
		raw, err := op.Encode(showOp(m, o))
		if err != nil {
			return nil, Errorf("change set %s: op %s: %v", cs, o.ID, err)
		}
		resp.Ops = append(resp.Ops, raw)
	}
	if d := c.Decided; d != nil {
		resp.Decided = &DecisionInfo{Type: d.Type, Entry: d.Entry, Actor: d.Actor, For: d.For, Rule: d.Rule, Reason: d.Reason}
	}
	if m.err != nil {
		return nil, Errorf("%s", m.err.Error())
	}
	return resp, nil
}

func infoIDs(infos []index.ChangeSetInfo) []id.ID {
	out := make([]id.ID, len(infos))
	for i, c := range infos {
		out[i] = c.ID
	}
	return out
}

// indexedEffective is effectiveChangeSet over the index's change sets:
// ctx.ChangeSet (it must be in the log), else ctx.ChangeRef, else current.
func indexedEffective(ctx *Ctx, current id.ID, infos []index.ChangeSetInfo) (id.ID, error) {
	ref := strings.TrimSpace(ctx.ChangeRef)
	if cs := ctx.ChangeSet; !cs.IsZero() {
		if ref != "" {
			return id.ID{}, Usagef("set either ChangeSet or ChangeRef, not both")
		}
		if !slices.ContainsFunc(infos, func(c index.ChangeSetInfo) bool { return c.ID == cs }) {
			return id.ID{}, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
		}
		return cs, nil
	}
	if ref != "" {
		return matchChangeSet(infoIDs(infos), ref)
	}
	return current, nil
}

func stateOf(infos []index.ChangeSetInfo, cs id.ID) changeset.State {
	for _, c := range infos {
		if c.ID == cs {
			return c.State
		}
	}
	return ""
}

// landedSets is what prefetchLanded loads: folded change sets by ID.
type landedSets struct {
	sets map[id.ID]changeset.ChangeSet
	err  error
}

// prefetchLanded loads, in the background, the change sets that the as-is
// history of the nodes and links mine's operations write or link to names:
// most of what indexedDiff reads, found without mine's own view. The
// channel yields once.
func prefetchLanded(x *index.Index, mine changeset.ChangeSet) <-chan landedSets {
	out := make(chan landedSets, 1)
	go func() {
		want := map[id.ID]bool{}
		for _, o := range mine.Ops {
			for _, i := range []id.ID{o.Target(), o.From, o.To} {
				if i.IsZero() {
					continue
				}
				var hist []gmodel.Ref
				if n, ok, err := x.Node(index.AsIs, i); err != nil {
					out <- landedSets{err: err}
					return
				} else if ok {
					hist = n.History
				} else if l, ok, err := x.Link(index.AsIs, i); err != nil {
					out <- landedSets{err: err}
					return
				} else if ok {
					hist = l.History
				}
				for _, h := range hist {
					if h.ChangeSet != mine.ID {
						want[h.ChangeSet] = true
					}
				}
			}
		}
		if len(want) == 0 {
			out <- landedSets{}
			return
		}
		got, err := x.FoldedOf(slices.Collect(maps.Keys(want)))
		out <- landedSets{sets: got, err: err}
	}()
	return out
}

// indexedDiff is changeset.BuildDiff for change set mine, over the change
// sets that can bear on it: those that wrote a node or link it writes or
// links to, found by the history the index keeps of each in mine's view
// (landed change sets and the in-flight ones it depends on). A change set
// the history does not name wrote nothing mine builds on, since a diff
// follows only the ops that lead to the ones mine is based on, so the diff
// is BuildDiff's over all of them; loaded holds change sets already
// loaded, and one of them that wrote none of those is left out. A landed
// operation that never applied has no history, so it is not found.
func indexedDiff(x *index.Index, m *indexModel, mine changeset.ChangeSet, loaded map[id.ID]changeset.ChangeSet) (changeset.Diff, error) {
	want := map[id.ID]bool{}
	have := map[id.ID]changeset.ChangeSet{}
	seen := map[id.ID]bool{}
	// touch adds the change sets that wrote node or link i.
	touch := func(i id.ID) {
		if i.IsZero() || seen[i] {
			return
		}
		seen[i] = true
		var hist []gmodel.Ref
		if n, ok := m.gnode(i); ok {
			hist = n.History
		} else if l, ok := m.glink(i); ok {
			hist = l.History
		}
		for _, h := range hist {
			if h.ChangeSet != mine.ID {
				want[h.ChangeSet] = true
			}
		}
	}
	var fail error
	load := func() {
		var need []id.ID
		for cs := range want {
			if _, ok := have[cs]; ok {
				continue
			}
			if c, ok := loaded[cs]; ok {
				have[cs] = c
				continue
			}
			need = append(need, cs)
		}
		if len(need) == 0 || fail != nil {
			return
		}
		got, err := x.FoldedOf(need)
		if err != nil {
			fail = err
			return
		}
		for cs, c := range got {
			have[cs] = c
		}
	}
	for _, o := range mine.Ops {
		touch(o.Target())
		touch(o.From)
		touch(o.To)
	}
	load()
	// A retired link's ends come from the op that created it.
	for _, o := range mine.Ops {
		if o.Type != op.TypeLinkRetire {
			continue
		}
		for _, c := range have {
			for _, p := range c.Ops {
				if p.Type == op.TypeLinkCreate && p.Link == o.Link {
					touch(p.From)
					touch(p.To)
				}
			}
		}
	}
	load()
	if fail == nil {
		fail = m.err
	}
	if fail != nil {
		return changeset.Diff{}, Errorf("%s", fail.Error())
	}
	sets := []changeset.ChangeSet{mine}
	for _, c := range have {
		sets = append(sets, c)
	}
	return changeset.BuildDiff(sets, mine.ID), nil
}

// indexedNotes is changeNotes over the index: m is the view of cs, and the
// current change set's view is added when it is another.
func indexedNotes(x *index.Index, infos []index.ChangeSetInfo, cs, current id.ID, m *indexModel) ([]ShowNote, error) {
	notes := []ShowNote{}
	views := []*indexModel{m}
	if !current.IsZero() && current != cs {
		// A note about cs written into current shows through current's
		// own view, which costs an overlay: build it only when current
		// has a link to cs (a note in a change set it depends on is
		// not looked for).
		es, err := x.Entries(current)
		if err != nil {
			return nil, Errorf("%s", err.Error())
		}
		cc, _ := changeset.Fold(current, es)
		if slices.ContainsFunc(cc.Ops, func(o op.Op) bool {
			return o.Type == op.TypeLinkCreate && o.LinkType == op.LinkDiscusses && o.To == cs
		}) {
			cm, err := newIndexModel(x, current, cc.State)
			if err != nil {
				return nil, err
			}
			views = append(views, cm)
		}
	}
	seen := map[id.ID]bool{}
	for _, v := range views {
		handle := func(a id.ID) string {
			n, _ := v.Node(a)
			return n.Handle
		}
		for _, l := range v.LinksTo(cs, op.LinkDiscusses) {
			g, ok := v.gnode(l.From)
			if !ok || g.Kind != op.KindNote || g.Retired || seen[g.ID] {
				continue
			}
			seen[g.ID] = true
			note := ShowNote{ID: g.ID, Text: gmodel.ShowText(fieldText(g.Fields["statement"]), modelHandleOf(v)), Op: noteOp(fieldText(g.Fields["op"]))}
			if len(g.History) > 0 {
				ncs := g.History[0].ChangeSet
				es, err := x.Entries(ncs)
				if err != nil {
					return nil, Errorf("%s", err.Error())
				}
				nc, _ := changeset.Fold(ncs, es)
				if rs := nc.History[g.History[0].Op]; len(rs) > 0 {
					note.Actor, note.At = noteWho(rs[0].Actor, rs[0].Entry, es, handle)
				}
			}
			notes = append(notes, note)
		}
		if v.err != nil {
			return nil, Errorf("%s", v.err.Error())
		}
	}
	sortNotes(notes)
	return notes, nil
}

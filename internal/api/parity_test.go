package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/index/indextest"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
	"pgregory.net/rapid"
)

// This file is the projection parity harness: on drawn logs, the
// projection (projection.go) and a fresh strict fold (changeset.FoldLog,
// then internal/model's AsIs, or WhatIf for an in-flight change set) must
// give the same write views. The strict engine is the spec.
//
// A view is the as-is view (no current change set) or the write view of one
// change set of the log. A terminal change set's write view is as-is with no
// current change set, as writeIndex.model documents, so "the current change set"
// below is the view's change set only while it is in flight.
//
// Version and Conflict follow the strict rule (Decision
// parity-version-rule): Version is the last operation in the target's
// strict history that is not the current change set's own, and Conflict is
// true when the view holds an unrepaired model.ErrConflict problem on the
// target from an operation not of the current change set. Pending lists
// the current change set's ops on each target the view holds, in change
// set order, whether the strict engine applied them or not.

// side names one of the two models compared.
type side string

const (
	sideIndex  side = "index"
	sideStrict side = "strict"
)

// parityView is one side's answer for one view: every node and link it
// holds, keyed by ID.
type parityView struct {
	nodes map[id.ID]Node
	links map[id.ID]Link
}

// disagreement is one difference between two sides on one view: on node or
// link Target, Field (presence, version, conflict, pending, or a node's or
// link's field name: kind, handle, parent, key, retired, type, from, to,
// note, decision, or fields.<name>) has value A on side SideA and B on
// SideB.
type disagreement struct {
	View         string
	Link         bool
	Target       id.ID
	Field        string
	SideA, SideB side
	A, B         string
}

func (d disagreement) what() string {
	if d.Link {
		return "link"
	}
	return "node"
}

func (d disagreement) String() string {
	return fmt.Sprintf("%s: %s %s: %s: %s %s, %s %s", d.View, d.what(), d.Target, d.Field, d.SideA, d.A, d.SideB, d.B)
}

// viewName names the view of change set cs (zero for as-is).
func viewName(cs id.ID) string {
	if cs.IsZero() {
		return "as-is"
	}
	return "write view of " + cs.String()
}

// show renders a value for a report.
func show(v any) string {
	switch v := v.(type) {
	case id.ID:
		if v.IsZero() {
			return "none"
		}
		return v.String()
	case []id.ID:
		if len(v) == 0 {
			return "[]"
		}
		s := make([]string, len(v))
		for i, x := range v {
			s[i] = x.String()
		}
		return "[" + strings.Join(s, " ") + "]"
	case string:
		return fmt.Sprintf("%q", v)
	case rawText:
		return string(v)
	}
	return fmt.Sprint(v)
}

// compareViews lists every difference between a (side sa) and b (side sb)
// on the view named view, nodes before links, each by ID. Pending is
// compared only when pending is true.
func compareViews(view string, sa side, a parityView, sb side, b parityView, pending bool) []disagreement {
	var out []disagreement
	add := func(link bool, t id.ID, field string, x, y any) {
		if sx, sy := show(x), show(y); sx != sy {
			out = append(out, disagreement{View: view, Link: link, Target: t, Field: field, SideA: sa, SideB: sb, A: sx, B: sy})
		}
	}
	nodeIDs := slices.Collect(maps.Keys(a.nodes))
	for k := range b.nodes {
		if _, ok := a.nodes[k]; !ok {
			nodeIDs = append(nodeIDs, k)
		}
	}
	slices.SortFunc(nodeIDs, id.Compare)
	for _, k := range nodeIDs {
		x, okx := a.nodes[k]
		y, oky := b.nodes[k]
		if okx != oky {
			add(false, k, "presence", present(okx), present(oky))
			continue
		}
		add(false, k, "kind", string(x.Kind), string(y.Kind))
		add(false, k, "handle", x.Handle, y.Handle)
		add(false, k, "parent", x.Parent, y.Parent)
		add(false, k, "key", x.Key, y.Key)
		add(false, k, "retired", x.Retired, y.Retired)
		for _, f := range fieldNames(x.Fields, y.Fields) {
			add(false, k, "fields."+f, fieldValue(x.Fields, f), fieldValue(y.Fields, f))
		}
		add(false, k, "version", x.Version, y.Version)
		add(false, k, "conflict", x.Conflict, y.Conflict)
		if pending {
			add(false, k, "pending", x.Pending, y.Pending)
		}
	}
	linkIDs := slices.Collect(maps.Keys(a.links))
	for k := range b.links {
		if _, ok := a.links[k]; !ok {
			linkIDs = append(linkIDs, k)
		}
	}
	slices.SortFunc(linkIDs, id.Compare)
	for _, k := range linkIDs {
		x, okx := a.links[k]
		y, oky := b.links[k]
		if okx != oky {
			add(true, k, "presence", present(okx), present(oky))
			continue
		}
		add(true, k, "type", string(x.Type), string(y.Type))
		add(true, k, "from", x.From, y.From)
		add(true, k, "to", x.To, y.To)
		add(true, k, "note", x.Note, y.Note)
		add(true, k, "decision", x.Decision, y.Decision)
		add(true, k, "retired", x.Retired, y.Retired)
		add(true, k, "version", x.Version, y.Version)
		add(true, k, "conflict", x.Conflict, y.Conflict)
		if pending {
			add(true, k, "pending", x.Pending, y.Pending)
		}
	}
	return out
}

func present(ok bool) string {
	if ok {
		return "present"
	}
	return "absent"
}

func fieldNames(a, b map[string]json.RawMessage) []string {
	names := slices.Collect(maps.Keys(a))
	for k := range b {
		if _, ok := a[k]; !ok {
			names = append(names, k)
		}
	}
	slices.Sort(names)
	return names
}

// rawText is shown as it is, unquoted: a field's JSON value.
type rawText string

func fieldValue(fs map[string]json.RawMessage, name string) rawText {
	v, ok := fs[name]
	if !ok {
		return "(unset)"
	}
	return rawText(bytes.TrimSpace(v))
}

// parityLog is one log as the two sides read it: the entries on disk
// plus a dry run's pending ones, folded.
type parityLog struct {
	entries []entry.Entry
	pending []entry.Entry
	sets    []changeset.ChangeSet
	// target is each op's target, for the strict side's conflicts.
	target map[id.ID]id.ID
}

func newParityLog(gdir string, pending []entry.Entry) parityLog {
	l := store.ReadLog(os.DirFS(gdir))
	p := parityLog{entries: l.Entries, pending: pending}
	l.Entries = append(slices.Clone(l.Entries), pending...)
	p.sets, _ = changeset.FoldLog(l)
	p.target = map[id.ID]id.ID{}
	for _, cs := range p.sets {
		for _, o := range cs.Ops {
			p.target[o.ID] = o.Target()
		}
	}
	return p
}

// current is the current change set of cs's write view: cs while it is in
// flight, else none.
func (p parityLog) current(cs id.ID) id.ID {
	for _, s := range p.sets {
		if s.ID == cs && inFlightState(s.State) {
			return cs
		}
	}
	return id.ID{}
}

// String lists the log for a report: each entry by change set, ID, type,
// and for an ops entry each op's ID, type, target and base.
func (p parityLog) String() string {
	var b strings.Builder
	line := func(e entry.Entry, tag string) {
		fmt.Fprintf(&b, "  %s%s %s after %s: %s", tag, e.ChangeSet, e.ID, show(e.After), e.Type)
		if e.Type == entry.TypeOps {
			if body, err := changeset.DecodeBody(e); err == nil {
				if ops, ok := body.(changeset.Ops); ok {
					for _, o := range ops.Ops {
						fmt.Fprintf(&b, "\n      %s %s %s base %s", o.ID, o.Type, o.Target(), show(o.Base))
					}
				}
			}
		}
		b.WriteString("\n")
	}
	for _, e := range p.entries {
		line(e, "")
	}
	for _, e := range p.pending {
		line(e, "pending ")
	}
	for _, s := range p.sets {
		fmt.Fprintf(&b, "  folded %s: %s\n", s.ID, s.State)
	}
	return b.String()
}

// indexView is the projection's view of cs, every node and link the drawn
// logs can hold queried by ID.
func indexView(t fataler, p projection, pending []entry.Entry, cs id.ID) parityView {
	m := p.model(pending, cs).(*writeModel)
	if err := p.failure(); err != nil {
		t.Fatalf("index model of %s: %v", viewName(cs), err)
	}
	v := parityView{nodes: map[id.ID]Node{}, links: map[id.ID]Link{}}
	for _, k := range indextest.Nodes {
		n, ok, err := m.ov.Node(k)
		if err != nil {
			t.Fatalf("index node %s: %v", k, err)
		}
		if ok {
			v.nodes[k] = m.node(n, true)
		}
	}
	for _, k := range indextest.Links {
		l, ok, err := m.ov.Link(k)
		if err != nil {
			t.Fatalf("index link %s: %v", k, err)
		}
		if ok {
			v.links[k] = m.link(l, true)
		}
	}
	return v
}

// fataler is what indexView needs of a *testing.T or a *rapid.T.
type fataler interface {
	Fatalf(format string, args ...any)
}

// strictView is the strict engine's view of cs by the rule at the top of
// this file.
func strictView(p parityLog, cs id.ID) (parityView, gmodel.Model, []gmodel.Problem) {
	cur := p.current(cs)
	var m gmodel.Model
	var ps []gmodel.Problem
	if cur.IsZero() {
		m, ps = gmodel.AsIs(p.sets)
	} else {
		m, ps = gmodel.WhatIf(p.sets, cur)
	}
	conflict := map[id.ID]bool{}
	for _, pr := range ps {
		if !errors.Is(pr.Err, gmodel.ErrConflict) || pr.ChangeSet == cur {
			continue
		}
		var stale *gmodel.StaleBaseError
		if errors.As(pr.Err, &stale) {
			conflict[stale.Node] = true
		} else if t, ok := p.target[pr.Op]; ok {
			conflict[t] = true
		}
	}
	version := func(h []gmodel.Ref) id.ID {
		for i := len(h) - 1; i >= 0; i-- {
			if cur.IsZero() || h[i].ChangeSet != cur {
				return h[i].Op
			}
		}
		return id.ID{}
	}
	v := parityView{nodes: map[id.ID]Node{}, links: map[id.ID]Link{}}
	for _, n := range m.Nodes() {
		a := apiNode(n)
		a.Version, a.Conflict = version(n.History), conflict[n.ID]
		v.nodes[n.ID] = a
	}
	for _, l := range m.Links() {
		v.links[l.ID] = Link{ID: l.ID, Type: l.Type, From: l.From, To: l.To, Note: l.Note, Decision: l.Decision,
			Retired: l.Retired, Version: version(l.History), Conflict: conflict[l.ID]}
	}
	for _, s := range p.sets {
		if cur.IsZero() || s.ID != cur {
			continue
		}
		for _, o := range s.Ops {
			if n, ok := v.nodes[o.Target()]; ok {
				n.Pending = append(n.Pending, o.ID)
				v.nodes[o.Target()] = n
			}
			if l, ok := v.links[o.Target()]; ok {
				l.Pending = append(l.Pending, o.ID)
				v.links[o.Target()] = l
			}
		}
	}
	return v, m, ps
}

// viewFacts is what a known case reads of one view's log shape: the log,
// the view's change set and its current change set, and the strict
// engine's view and problems.
type viewFacts struct {
	log      parityLog
	cs, cur  id.ID
	strict   gmodel.Model
	problems []gmodel.Problem
}

func newViewFacts(p parityLog, cs id.ID, m gmodel.Model, ps []gmodel.Problem) *viewFacts {
	return &viewFacts{log: p, cs: cs, cur: p.current(cs), strict: m, problems: ps}
}

// knownCase is a disagreement the harness has found and tolerates until
// the leaf that fixes it empties this table. The index's write model takes
// its views from the strict engine (parity-write-index-from-strict), so
// the table is empty: any disagreement fails. Name names the log shape that gives
// it; match recognises a disagreement of that shape.
type knownCase struct {
	name  string
	match func(f *viewFacts, d disagreement) bool
}

// knownCases is the table of known disagreements (Decision
// parity-known-cases-table), matched in order. Each one is reported with
// its count. The cases the index's write model had, each fixed and kept as
// a TestParityIndex test, are in the notes on leaf
// parity-write-index-from-strict.
var knownCases = []knownCase{}

// known returns the known case d is one of, or "" for none.
func known(f *viewFacts, d disagreement) string {
	for _, k := range knownCases {
		if k.match(f, d) {
			return k.name
		}
	}
	return ""
}

// parityCounts is what TestPropertyProjectionParity drew, to fail a vacuous
// run.
type parityCounts struct {
	draws, twoHeads, rebuilds, pending, views int
	known                                     map[string]int
}

// twoLandedHeads reports whether some target has two ops among the landed
// change sets' that no other landed op names as base.
func twoLandedHeads(sets []changeset.ChangeSet) bool {
	byTarget := map[id.ID][]op.Op{}
	for _, cs := range sets {
		if cs.State == changeset.StateLanded {
			for _, o := range cs.Ops {
				byTarget[o.Target()] = append(byTarget[o.Target()], o)
			}
		}
	}
	for _, ops := range byTarget {
		named := map[id.ID]bool{}
		for _, o := range ops {
			if o.Base != o.ID {
				named[o.Base] = true
			}
		}
		heads := 0
		for _, o := range ops {
			if !named[o.ID] {
				heads++
			}
		}
		if heads > 1 {
			return true
		}
	}
	return false
}

// checkParity compares the two sides on every view of p: as-is and each
// change set's write view, and, when p has pending entries, only the write
// view of their change set.
func checkParity(t *rapid.T, ix projection, p parityLog, c *parityCounts) {
	views := []id.ID{{}}
	for _, s := range p.sets {
		views = append(views, s.ID)
	}
	if len(p.pending) > 0 {
		views = []id.ID{p.pending[0].ChangeSet}
	}
	var unknown []string
	for _, cs := range views {
		c.views++
		name := viewName(cs)
		strict, sm, ps := strictView(p, cs)
		facts := newViewFacts(p, cs, sm, ps)
		indexed := indexView(t, ix, p.pending, cs)
		for _, d := range compareViews(name, sideIndex, indexed, sideStrict, strict, true) {
			if k := known(facts, d); k != "" {
				c.known[k]++
				continue
			}
			unknown = append(unknown, d.String())
		}
	}
	if len(unknown) > 0 {
		t.Fatalf("the models disagree:\n  %s\nlog:\n%s", strings.Join(unknown, "\n  "), p)
	}
}

// TestPropertyProjectionParity: on drawn logs, caught up step by step
// (rebuilds included), and with dry-run pending entries laid over at the
// end, the index's write model and the strict engine give the same as-is view and the same write view of every change set,
// but for the known cases.
func TestPropertyProjectionParity(t *testing.T) {
	c := &parityCounts{known: map[string]int{}}
	base := t.TempDir()
	rapid.Check(t, func(t *rapid.T) {
		c.draws++
		lg := indextest.Draw(t)
		root, err := os.MkdirTemp(base, "parity-")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(root)
		gdir := filepath.Join(root, store.Dir)
		if err := os.MkdirAll(filepath.Join(gdir, store.LogDir), 0o755); err != nil {
			t.Fatal(err)
		}
		ctx := &Ctx{Root: root, IndexPath: filepath.Join(root, "index.db")}
		x, err := index.Open(ctx.Root, ctx.IndexPath)
		if err != nil {
			t.Fatalf("index.Open: %v", err)
		}
		var ix projection = &writeIndex{ctx: ctx, x: x}
		defer ix.close()

		rebuilt, heads := false, false
		st := indextest.NewSteps(gdir, lg.Files)
		for !st.Done() {
			st.Next(t)
			if err := x.Refresh(); err != nil {
				t.Fatalf("Refresh: %v", err)
			}
			rebuilt = rebuilt || x.LastCatchUp().Rebuilt
			p := newParityLog(gdir, nil)
			heads = heads || twoLandedHeads(p.sets)
			checkParity(t, ix, p, c)
		}
		if rebuilt {
			c.rebuilds++
		}
		if heads {
			c.twoHeads++
		}

		// A dry run's pending entries on one in-flight change set.
		p := newParityLog(gdir, nil)
		var inFlight []id.ID
		for _, s := range p.sets {
			if inFlightState(s.State) {
				inFlight = append(inFlight, s.ID)
			}
		}
		if len(inFlight) == 0 {
			return
		}
		cs := rapid.SampledFrom(inFlight).Draw(t, "pending change set")
		var pending []entry.Entry
		for range rapid.IntRange(0, 2).Draw(t, "pending entries") {
			if e, ok := lg.Pending(t, cs); ok {
				pending = append(pending, e)
			}
		}
		if len(pending) == 0 {
			return
		}
		c.pending++
		checkParity(t, ix, newParityLog(gdir, pending), c)
	})
	t.Logf("draws %d, with two landed heads on one target %d, with a rebuild %d, with pending entries %d, views %d",
		c.draws, c.twoHeads, c.rebuilds, c.pending, c.views)
	for _, k := range knownCases {
		t.Logf("known case %q: %d", k.name, c.known[k.name])
	}
	if c.twoHeads == 0 || c.rebuilds == 0 || c.pending == 0 {
		t.Errorf("vacuous run: two landed heads %d, rebuilds %d, pending %d", c.twoHeads, c.rebuilds, c.pending)
	}
}

// TestParityReportNamesNode: two views that differ in one node's field
// give one disagreement, and its report names the view, the node, the
// field and both values.
func TestParityReportNamesNode(t *testing.T) {
	n := indextest.Node(1)
	mk := func(statement string) parityView {
		return parityView{
			nodes: map[id.ID]Node{n: {ID: n, Kind: op.KindCapability, Handle: "a", Key: "i",
				Fields: map[string]json.RawMessage{"statement": json.RawMessage(statement)}}},
			links: map[id.ID]Link{},
		}
	}
	cs := indextest.ChangeSet(1)
	ds := compareViews(viewName(cs), sideIndex, mk(`"old"`), sideStrict, mk(`"new"`), true)
	if len(ds) != 1 {
		t.Fatalf("disagreements = %v, want one", ds)
	}
	got := ds[0].String()
	for _, want := range []string{"write view of " + cs.String(), "node " + n.String(), "fields.statement",
		`index "old"`, `strict "new"`} {
		if !strings.Contains(got, want) {
			t.Errorf("report %q lacks %q", got, want)
		}
	}
}

// parityOf is the parity log of b's entries, folded without an index.
func parityOf(t *testing.T, b *logBuilder) parityLog {
	t.Helper()
	sets, problems := changeset.FoldLog(store.Log{Entries: b.entries})
	if len(problems) > 0 {
		t.Fatalf("fold problems: %v", problems)
	}
	p := parityLog{entries: b.entries, sets: sets, target: map[id.ID]id.ID{}}
	for _, cs := range sets {
		for _, o := range cs.Ops {
			p.target[o.ID] = o.Target()
		}
	}
	return p
}

// TestParityCreateNoLongerApplies: a draft creates a node under a landed
// parent that a later landed change set retires. The strict what-if skips
// the create (model.ErrRetired), so the draft's write view has no such
// node, and neither has the projection's.
func TestParityCreateNoLongerApplies(t *testing.T) {
	b := newLogBuilder(t)
	parent := b.create(op.KindCapability, "parent", nil)
	b.landed(parent)
	child := b.create(op.KindCapability, "child", nil)
	child.Parent = parent.Node
	draft := b.draft(child)
	b.landed(op.Op{ID: b.id(), Type: op.TypeNodeRetire, Node: parent.Node, Base: parent.ID, Reason: "gone"})

	if n, ok := b.model(draft).Node(child.Node); ok {
		t.Errorf("write view of the draft holds node %s (%s), whose create no longer applies", n.ID, n.Handle)
	}
	indexAgrees(t, b, draft)
}

// TestParityDisjointHeads: two landed updates of one node, both on its
// create, write different fields. The strict engine merges them with no
// conflict and takes the last one applied as the version (Decision
// parity-version-rule), and so does the projection.
func TestParityDisjointHeads(t *testing.T) {
	b := newLogBuilder(t)
	c := b.create(op.KindCapability, "fork", statement("a"))
	b.landed(c)
	u1 := b.update(c.Node, c.ID, statement("b"))
	b.landed(u1)
	u2 := b.update(c.Node, c.ID, map[string]json.RawMessage{"acceptance": json.RawMessage(`"x"`)})
	b.landed(u2)

	n, err := b.model(id.ID{}).Resolve(c.Node.String())
	if err != nil {
		t.Fatalf("as-is lacks node %s: %v", c.Node, err)
	}
	if n.Conflict || n.Version != u2.ID {
		t.Errorf("conflict %v version %s, want false %s", n.Conflict, n.Version, u2.ID)
	}
	if string(n.Fields["statement"]) != `"b"` || string(n.Fields["acceptance"]) != `"x"` {
		t.Errorf("fields %v, want both updates' fields", n.Fields)
	}
	indexAgrees(t, b, id.ID{})
}

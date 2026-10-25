package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/index"
	gmodel "github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

// planID is a fixed ID whose short ID tells n apart.
func planID(n int) id.ID {
	i, err := id.Parse(fmt.Sprintf("01M5%04d000000000000000000", n))
	if err != nil {
		panic(err)
	}
	return i
}

// planFix builds planInputs by hand. Every change set's ops go in in.ops,
// whatever its state; landed holds the landed stand-in's ops and history
// any op's revisions other than the one each op gets by default.
type planFix struct {
	in      planInputs
	n       int
	landed  []op.Op
	history map[id.ID][]changeset.Revision
}

func newPlanFix() *planFix {
	return &planFix{in: planInputs{ops: map[id.ID][]op.Op{}, nodes: map[id.ID]gmodel.Node{},
		ancestors: map[id.ID][]id.ID{}, intentPath: map[id.ID][]id.ID{}, milestones: map[id.ID]bool{}},
		history: map[id.ID][]changeset.Revision{}}
}

// core fills in.sets from the change sets and their ops, as planGather
// does, and runs the core. An op's history is one version written by the
// entry with its own ID, unless f.history says otherwise.
func (f *planFix) core() planResult {
	f.in.sets = []changeset.ChangeSet{{State: changeset.StateLanded, Ops: f.landed}}
	for _, c := range f.in.changeSets {
		s := changeset.ChangeSet{ID: c.ID, State: c.State, Ops: f.in.ops[c.ID], History: map[id.ID][]changeset.Revision{}}
		for _, o := range s.Ops {
			s.History[o.ID] = []changeset.Revision{{Entry: o.ID, Op: o}}
			if h, ok := f.history[o.ID]; ok {
				s.History[o.ID] = h
			}
		}
		f.in.sets = append(f.in.sets, s)
	}
	return planCore(f.in)
}

// dep adds a landed depends-on link from cs to to and returns its ID.
func (f *planFix) dep(cs, to id.ID) id.ID {
	l := f.next()
	f.landed = append(f.landed, op.Op{ID: f.next(), Type: op.TypeLinkCreate, Link: l,
		LinkType: op.LinkDependsOn, From: cs, To: to})
	return l
}

// undep retires link l: landed, or in change set cs when cs is not zero.
func (f *planFix) undep(cs, l id.ID) {
	o := op.Op{ID: f.next(), Type: op.TypeLinkRetire, Link: l}
	if cs.IsZero() {
		f.landed = append(f.landed, o)
		return
	}
	f.in.ops[cs] = append(f.in.ops[cs], o)
}

// add appends o to change set cs's ops with a fresh op ID, and returns it.
func (f *planFix) add(cs id.ID, o op.Op) op.Op {
	o.ID = f.next()
	f.in.ops[cs] = append(f.in.ops[cs], o)
	return o
}

func (f *planFix) next() id.ID {
	f.n++
	return planID(f.n)
}

// cs adds a change set in state s whose intent is name.
func (f *planFix) cs(name string, s changeset.State) id.ID {
	i := f.next()
	f.in.changeSets = append(f.in.changeSets, index.ChangeSetInfo{ID: i, State: s,
		Kind: changeset.KindIntention, Intent: name})
	return i
}

// node adds a live node of kind with handle under parent (zero for none).
func (f *planFix) node(kind op.Kind, handle string, parent id.ID, key string) id.ID {
	i := f.next()
	f.in.nodes[i] = gmodel.Node{ID: i, Kind: kind, Handle: handle, Parent: parent, Key: key}
	return i
}

// retire marks node i retired.
func (f *planFix) retire(i id.ID) {
	n := f.in.nodes[i]
	n.Retired = true
	f.in.nodes[i] = n
}

// link adds a view link and returns its ID.
func (f *planFix) link(t op.LinkType, from, to id.ID, retired bool) id.ID {
	i := f.next()
	f.in.links = append(f.in.links, gmodel.Link{ID: i, Type: t, From: from, To: to, Retired: retired})
	return i
}

// own adds a link.create op to change set cs and returns the link's ID.
func (f *planFix) own(cs id.ID, t op.LinkType, from, to id.ID) id.ID {
	l := f.next()
	f.in.ops[cs] = append(f.in.ops[cs], op.Op{ID: f.next(), Type: op.TypeLinkCreate, Link: l,
		LinkType: t, From: from, To: to})
	return l
}

// touch adds a node.update op on node to change set cs.
func (f *planFix) touch(cs, node id.ID, ancestors ...id.ID) {
	f.in.ops[cs] = append(f.in.ops[cs], op.Op{ID: f.next(), Type: op.TypeNodeUpdate, Node: node})
	f.in.ancestors[node] = ancestors
}

// intent adds a served intent's path.
func (f *planFix) intent(handle string, parent id.ID, key string, path ...id.ID) id.ID {
	i := f.node(op.KindIntent, handle, parent, key)
	f.in.intentPath[i] = append([]id.ID{}, path...)
	return i
}

// milestone marks node i a milestone (its milestone field non-empty).
func (f *planFix) milestone(i id.ID) {
	f.in.milestones[i] = true
}

// serve adds a live view serves link from cs to each of intents.
func (f *planFix) serve(cs id.ID, intents ...id.ID) {
	for _, i := range intents {
		f.link(op.LinkServes, cs, i, false)
	}
}

// verdict runs the core and returns cs's verdict.
func (f *planFix) verdict(t *testing.T, cs id.ID) *planChangeSet {
	t.Helper()
	for _, p := range f.core().order {
		if p.info.ID == cs {
			return p
		}
	}
	t.Fatalf("change set %s not in the ready order", cs)
	return nil
}

func wantReasons(t *testing.T, p *planChangeSet, want ...string) {
	t.Helper()
	if len(want) == 0 {
		want = nil
	}
	if !reflect.DeepEqual(p.reasons, want) || p.ready != (len(want) == 0) {
		t.Errorf("%s: ready %v, reasons %q; want %q", p.info.Intent, p.ready, p.reasons, want)
	}
}

func TestPlanReadyWithLandedDep(t *testing.T) {
	f := newPlanFix()
	landed := f.cs("L", changeset.StateLanded)
	a := f.cs("A", changeset.StateAccepted)
	f.dep(a, landed)
	p := f.verdict(t, a)
	wantReasons(t, p)
	if !reflect.DeepEqual(p.dependsOn, []id.ID{landed}) {
		t.Errorf("dependsOn = %v", p.dependsOn)
	}
}

func TestPlanDependencyNotLanded(t *testing.T) {
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	gone := planID(999)
	f.dep(b, a)
	f.own(b, op.LinkDependsOn, b, gone)
	wantReasons(t, f.verdict(t, b),
		"depends on "+shortID(a)+" (accepted)",
		"depends on "+shortID(gone)+" (not in the log)")
	// A retired depends-on link does not count, retired where it landed or
	// by the change set's own ops.
	c := f.cs("C", changeset.StateAccepted)
	f.undep(id.ID{}, f.dep(c, a))
	f.undep(c, f.dep(c, b))
	f.undep(c, f.own(c, op.LinkDependsOn, c, b))
	wantReasons(t, f.verdict(t, c))
}

func TestPlanDependencyOnlyItsOwn(t *testing.T) {
	// A depends-on link from B written inside a third change set does not
	// make B depend (D1), though the to-be view has it.
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	l := f.own(c, op.LinkDependsOn, b, a)
	f.in.links = append(f.in.links, gmodel.Link{ID: l, Type: op.LinkDependsOn, From: b, To: a, Created: c})
	if p := f.verdict(t, b); len(p.dependsOn) != 0 || len(p.reasons) != 0 {
		t.Errorf("B: dependsOn %v, reasons %q", p.dependsOn, p.reasons)
	}
}

func TestPlanRestsOnDropped(t *testing.T) {
	// B updates a node withdrawn A created: rests-on comes after the
	// depends-on reasons and before the blocking questions, and stays when
	// B no longer depends on A.
	f := newPlanFix()
	a := f.cs("A", changeset.StateWithdrawn)
	b := f.cs("B", changeset.StateProposed)
	n := f.next()
	created := f.add(a, op.Op{Type: op.TypeNodeCreate, Node: n, Kind: op.KindCapability})
	l := f.own(b, op.LinkDependsOn, b, a)
	f.add(b, op.Op{Type: op.TypeNodeUpdate, Node: n, Base: created.ID})
	q := f.node(op.KindQuestion, "q", id.ID{}, "")
	f.link(op.LinkBlocks, q, b, false)
	wantReasons(t, f.verdict(t, b),
		"not accepted: proposed",
		"depends on "+shortID(a)+" (withdrawn)",
		"rests on "+shortID(a)+" (withdrawn): 1 operations",
		"blocked by question q on this change set")
	f.undep(b, l)
	wantReasons(t, f.verdict(t, b),
		"not accepted: proposed",
		"rests on "+shortID(a)+" (withdrawn): 1 operations",
		"blocked by question q on this change set")
}

func TestPlanRevisedUnder(t *testing.T) {
	// A revises in place, after B wrote it, a field B sets: one reason per
	// dependency, counting B's ops, after the depends-on reasons.
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	n := f.next()
	v1 := op.Op{Type: op.TypeNodeCreate, Node: n, Kind: op.KindCapability,
		Fields: map[string]json.RawMessage{"statement": json.RawMessage(`"one"`)}}
	created := f.add(a, v1)
	v1.ID = created.ID
	f.own(b, op.LinkDependsOn, b, a)
	set := map[string]json.RawMessage{"statement": json.RawMessage(`"b"`)}
	u := f.add(b, op.Op{Type: op.TypeNodeUpdate, Node: n, Base: created.ID, Set: set})
	f.add(b, op.Op{Type: op.TypeNodeUpdate, Node: n, Base: u.ID, Set: set})
	v2 := v1
	v2.Fields = map[string]json.RawMessage{"statement": json.RawMessage(`"two"`)}
	f.in.ops[a][0] = v2
	f.history[created.ID] = []changeset.Revision{{Entry: created.ID, Op: v1}, {Entry: f.next(), Op: v2}}
	wantReasons(t, f.verdict(t, b),
		"depends on "+shortID(a)+" (accepted)",
		shortID(a)+" revised under it: 2 operations")
}

func TestPlanNotAccepted(t *testing.T) {
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	p := f.cs("P", changeset.StateProposed)
	d := f.cs("D", changeset.StateDraft)
	f.dep(p, a)
	wantReasons(t, f.verdict(t, p), "not accepted: proposed", "depends on "+shortID(a)+" (accepted)")
	wantReasons(t, f.verdict(t, d), "not accepted: draft")
}

func TestPlanBlockedByQuestion(t *testing.T) {
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	root := f.node(op.KindCapability, "store", id.ID{}, "s")
	mid := f.node(op.KindCapability, "cache", root, "c")
	leaf := f.node(op.KindCapability, "ram", mid, "r")
	f.touch(a, leaf, root, mid)

	onCS := f.node(op.KindQuestion, "q-cs", id.ID{}, "")
	onLeaf := f.node(op.KindQuestion, "q-leaf", id.ID{}, "")
	onRoot := f.node(op.KindQuestion, "q-root", id.ID{}, "")
	f.link(op.LinkBlocks, onCS, a, false)
	f.link(op.LinkBlocks, onLeaf, leaf, false)
	f.link(op.LinkBlocks, onRoot, root, false)
	wantReasons(t, f.verdict(t, a),
		"blocked by question q-cs on this change set",
		"blocked by question q-leaf on ram",
		"blocked by question q-root on store")

	// A live decision answering a question unblocks it; a retired answers
	// link does not, nor a retired blocks link or a retired question.
	dec := f.node(op.KindDecision, "d", id.ID{}, "")
	f.link(op.LinkAnswers, dec, onCS, false)
	f.link(op.LinkAnswers, dec, onLeaf, true)
	f.retire(onRoot)
	wantReasons(t, f.verdict(t, a), "blocked by question q-leaf on ram")
}

func TestPlanBlocksOwnOps(t *testing.T) {
	// A blocks link written inside the change set counts too (D7).
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	q := f.node(op.KindQuestion, "q", id.ID{}, "")
	f.own(a, op.LinkBlocks, q, a)
	wantReasons(t, f.verdict(t, a), "blocked by question q on this change set")
}

func TestPlanServesInsideChangeSet(t *testing.T) {
	// A serves link written inside a proposed change set groups it (D7);
	// a retired view copy of the same link wins over the op.
	f := newPlanFix()
	fast := f.intent("fast", id.ID{}, "f")
	f.milestone(fast)
	p := f.cs("P", changeset.StateProposed)
	q := f.cs("Q", changeset.StateProposed)
	f.own(p, op.LinkServes, p, fast)
	f.own(q, op.LinkServes, q, fast)
	o := f.in.ops[q][0]
	f.in.links = append(f.in.links, gmodel.Link{ID: o.Link, Type: o.LinkType, From: q, To: fast, Retired: true})
	res := f.core()
	if len(res.groups) != 2 || res.groups[0].intent != fast || len(res.groups[0].members) != 1 ||
		res.groups[0].members[0].info.ID != p || res.groups[0].backlog ||
		!res.groups[1].intent.IsZero() || !res.groups[1].backlog || res.groups[1].members[0].info.ID != q {
		t.Errorf("groups = %+v", res.groups)
	}
}

// planGroupIDs is each group's intent and member IDs.
func planGroupIDs(res planResult) [][]id.ID {
	var out [][]id.ID
	for _, g := range res.groups {
		row := []id.ID{g.intent}
		for _, p := range g.members {
			row = append(row, p.info.ID)
		}
		out = append(out, row)
	}
	return out
}

func TestPlanTopologicalOrder(t *testing.T) {
	f := newPlanFix()
	i := f.intent("fast", id.ID{}, "f")
	f.milestone(i)
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	for _, cs := range []id.ID{a, b, c} {
		f.link(op.LinkServes, cs, i, false)
	}
	f.dep(a, c)
	got := planGroupIDs(f.core())
	want := [][]id.ID{{i, b, c, a}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
}

func TestPlanCycle(t *testing.T) {
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	f.dep(a, b)
	f.dep(b, a)
	res := f.core()
	if got, want := planGroupIDs(res), [][]id.ID{{{}, c, a, b}}; !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	wantReasons(t, res.order[1], "depends on "+shortID(b)+" (accepted)", planCycleReason)
	wantReasons(t, res.order[2], "depends on "+shortID(a)+" (accepted)", planCycleReason)
	wantReasons(t, res.order[0])
}

func TestPlanRoadmapGroupsAndReadyOrder(t *testing.T) {
	// Milestone groups in intent-tree pre-order (siblings by key), then one
	// backlog: work on plain, retired and no intents, each once; drafts only
	// in the ready order, after.
	f := newPlanFix()
	late := f.intent("late", id.ID{}, "z")
	early := f.intent("early", id.ID{}, "a")
	child := f.intent("child", early, "m", early)
	plain := f.intent("plain", id.ID{}, "p")
	gone := f.intent("gone", id.ID{}, "b")
	f.milestone(late)
	f.milestone(child)
	f.retire(gone)
	draft := f.cs("D", changeset.StateDraft)
	both := f.cs("Both", changeset.StateAccepted)
	alone := f.cs("Alone", changeset.StateProposed)
	l := f.cs("L", changeset.StateAccepted)
	g := f.cs("G", changeset.StateAccepted)
	e := f.cs("E", changeset.StateProposed)
	f.serve(draft, child)
	f.serve(both, late, child)
	f.serve(l, late)
	f.serve(g, gone)
	f.serve(e, early, plain)
	res := f.core()
	want := [][]id.ID{{child, both}, {late, both, l}, {{}, alone, g, e}}
	if got := planGroupIDs(res); !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	if got := planBacklogs(res); !reflect.DeepEqual(got, []bool{false, false, true}) {
		t.Errorf("backlog = %v", got)
	}
	if got, want := planOrderIDs(res), []id.ID{both, l, alone, g, e, draft}; !reflect.DeepEqual(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// planBacklogs is each group's backlog mark.
func planBacklogs(res planResult) []bool {
	var out []bool
	for _, g := range res.groups {
		out = append(out, g.backlog)
	}
	return out
}

// planOrderIDs is the ready order's IDs.
func planOrderIDs(res planResult) []id.ID {
	var out []id.ID
	for _, p := range res.order {
		out = append(out, p.info.ID)
	}
	return out
}

// planWantGroups checks the core's groups against want, rows of intent then
// member IDs, the backlog's intent zero.
func planWantGroups(t *testing.T, f *planFix, want ...[]id.ID) planResult {
	t.Helper()
	res := f.core()
	if got := planGroupIDs(res); !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	for _, g := range res.groups {
		if g.backlog != g.intent.IsZero() {
			t.Errorf("group %s: backlog %v", g.intent, g.backlog)
		}
	}
	return res
}

func TestPlanMilestonesThenBacklog(t *testing.T) {
	f := newPlanFix()
	m2 := f.intent("m2", id.ID{}, "b")
	m1 := f.intent("m1", id.ID{}, "a")
	plain := f.intent("plain", id.ID{}, "c")
	f.milestone(m1)
	f.milestone(m2)
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	f.serve(a, m2)
	f.serve(b, plain)
	f.serve(c, m1)
	planWantGroups(t, f, []id.ID{m1, c}, []id.ID{m2, a}, []id.ID{{}, b})
}

func TestPlanMilestoneAndPlainIntent(t *testing.T) {
	// Serving a milestone and a plain intent: the milestone group only.
	f := newPlanFix()
	m := f.intent("m", id.ID{}, "a")
	plain := f.intent("plain", id.ID{}, "b")
	f.milestone(m)
	a := f.cs("A", changeset.StateAccepted)
	f.serve(a, m, plain)
	planWantGroups(t, f, []id.ID{m, a})
}

func TestPlanDonePointServes(t *testing.T) {
	// Serving an intent under a milestone puts the change set in the
	// milestone's group, whose intent is the milestone.
	f := newPlanFix()
	m := f.node(op.KindIntent, "m", id.ID{}, "a")
	mid := f.node(op.KindIntent, "mid", m, "b")
	done := f.intent("done", mid, "c", m, mid)
	f.milestone(m)
	a := f.cs("A", changeset.StateAccepted)
	f.serve(a, done)
	planWantGroups(t, f, []id.ID{m, a})
}

func TestPlanNestedMilestones(t *testing.T) {
	// Nested milestones: the nearest at or above a served intent only.
	f := newPlanFix()
	outer := f.intent("outer", id.ID{}, "a")
	inner := f.node(op.KindIntent, "inner", outer, "b")
	done := f.intent("done", inner, "c", outer, inner)
	f.milestone(outer)
	f.milestone(inner)
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	f.serve(a, done)
	f.serve(b, outer)
	planWantGroups(t, f, []id.ID{outer, b}, []id.ID{inner, a})
}

func TestPlanTwoMilestones(t *testing.T) {
	// Serving two milestones: in both groups, once in the ready order.
	f := newPlanFix()
	m1 := f.intent("m1", id.ID{}, "a")
	m2 := f.intent("m2", id.ID{}, "b")
	f.milestone(m1)
	f.milestone(m2)
	a := f.cs("A", changeset.StateAccepted)
	f.serve(a, m2, m1)
	res := planWantGroups(t, f, []id.ID{m1, a}, []id.ID{m2, a})
	if got := planOrderIDs(res); !reflect.DeepEqual(got, []id.ID{a}) {
		t.Errorf("order = %v", got)
	}
}

func TestPlanRetiredMilestone(t *testing.T) {
	// A retired milestone, or one under a retired node, is no milestone:
	// its work goes to the backlog. A live milestone above a retired one
	// does not take it either, the served intent being outside the tree.
	f := newPlanFix()
	gone := f.intent("gone", id.ID{}, "a")
	f.milestone(gone)
	f.retire(gone)
	dead := f.node(op.KindIntent, "dead", id.ID{}, "b")
	f.retire(dead)
	under := f.intent("under", dead, "c", dead)
	f.milestone(under)
	top := f.node(op.KindIntent, "top", id.ID{}, "d")
	f.milestone(top)
	mid := f.node(op.KindIntent, "mid", top, "e")
	f.retire(mid)
	leaf := f.intent("leaf", mid, "f", top, mid)
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	f.serve(a, gone)
	f.serve(b, under)
	f.serve(c, leaf)
	planWantGroups(t, f, []id.ID{{}, a, b, c})
}

func TestPlanServesNothingInBacklog(t *testing.T) {
	f := newPlanFix()
	m := f.intent("m", id.ID{}, "a")
	f.milestone(m)
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateProposed)
	f.serve(a, m)
	planWantGroups(t, f, []id.ID{m, a}, []id.ID{{}, b})
}

func TestPlanBacklogOrder(t *testing.T) {
	// The backlog follows depends-on, ties by ID; cycle members come last
	// with the cycle reason.
	f := newPlanFix()
	plain := f.intent("plain", id.ID{}, "a")
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	x := f.cs("X", changeset.StateAccepted)
	y := f.cs("Y", changeset.StateAccepted)
	f.serve(a, plain)
	f.serve(c, plain)
	f.dep(a, c)
	f.dep(x, y)
	f.dep(y, x)
	res := planWantGroups(t, f, []id.ID{{}, b, c, a, x, y})
	wantReasons(t, res.order[3], "depends on "+shortID(y)+" (accepted)", planCycleReason)
	wantReasons(t, res.order[4], "depends on "+shortID(x)+" (accepted)", planCycleReason)
}

func TestPlanNoBacklogGroup(t *testing.T) {
	// No backlog group when every change set serves a milestone (drafts
	// never show), and no groups when nothing is in flight.
	f := newPlanFix()
	m := f.intent("m", id.ID{}, "a")
	f.milestone(m)
	a := f.cs("A", changeset.StateAccepted)
	d := f.cs("D", changeset.StateDraft)
	f.serve(a, m)
	res := planWantGroups(t, f, []id.ID{m, a})
	if got := planOrderIDs(res); !reflect.DeepEqual(got, []id.ID{a, d}) {
		t.Errorf("order = %v", got)
	}

	f = newPlanFix()
	m = f.intent("m", id.ID{}, "a")
	f.milestone(m)
	f.cs("L", changeset.StateLanded)
	if res := f.core(); len(res.groups) != 0 || len(res.order) != 0 {
		t.Errorf("groups = %v, order = %v", planGroupIDs(res), planOrderIDs(res))
	}
}

func TestPlanReadyFollowsGroups(t *testing.T) {
	// Ready lists the milestone work first, then the backlog, then the
	// drafts by ID, whatever the IDs.
	f := newPlanFix()
	d1 := f.cs("D1", changeset.StateDraft)
	back := f.cs("Back", changeset.StateAccepted)
	m := f.intent("m", id.ID{}, "a")
	f.milestone(m)
	d2 := f.cs("D2", changeset.StateDraft)
	work := f.cs("Work", changeset.StateAccepted)
	f.serve(work, m)
	f.serve(d2, m)
	res := planWantGroups(t, f, []id.ID{m, work}, []id.ID{{}, back})
	if got, want := planOrderIDs(res), []id.ID{work, back, d1, d2}; !reflect.DeepEqual(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// planLink links from to to through the link operation, inside the
// effective change set.
func planLink(t *testing.T, ctx *Ctx, lt op.LinkType, from, to id.ID) {
	t.Helper()
	req := fmt.Sprintf(`{"from":%q,"type":%q,"to":%q}`, from, lt, to)
	if _, err := callOp(t, ctx, "link", req); err != nil {
		t.Fatalf("link %s %s %s: %v", from, lt, to, err)
	}
}

func readyOf(t *testing.T, ctx *Ctx, req string) ReadyResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "ready", req)
	if err != nil {
		t.Fatalf("ready: %v", err)
	}
	return resp.(ReadyResponse)
}

func TestReadyEndToEnd(t *testing.T) {
	ctx := readCtx(t)
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"fast"}`)
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	recordWrite(t, ctx, "add", `{"kind":"capability","handle":"cache","parent":"store"}`)
	fast := nodeAt(t, ctx, "fast").ID

	a := newChange(t, ctx, `{"intent":"First"}`).ChangeSet
	planLink(t, ctx, op.LinkServes, a, fast)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	b := newChange(t, ctx, `{"intent":"Second"}`).ChangeSet
	planLink(t, ctx, op.LinkServes, b, fast)
	planLink(t, ctx, op.LinkDependsOn, b, a)
	mustWrite(t, ctx, "set", `{"ref":"cache","fields":["statement=Fast"]}`)

	// B is a draft (its own ops say what it serves and depends on); the
	// current change set does not move the default view off to-be (D17).
	r := readyOf(t, ctx, `{}`)
	if r.View != "to-be" || r.Next != a || len(r.ChangeSets) != 2 || r.ChangeSets[1].ID != b {
		t.Fatalf("ready = %+v", r)
	}
	if want := []string{"not accepted: draft", "depends on " + shortID(a) + " (accepted)"}; !reflect.DeepEqual(r.ChangeSets[1].Reasons, want) {
		t.Errorf("B reasons = %q, want %q", r.ChangeSets[1].Reasons, want)
	}
	if out := human(t, r); !strings.HasPrefix(out, "next: "+shortID(a)+" First\n") {
		t.Errorf("human = %q", out)
	}

	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	la := *ctx
	la.ChangeSet = a
	move(t, &la, "land", `{}`)
	recordWrite(t, ctx, "add", `{"kind":"question","handle":"which"}`)
	recordWrite(t, ctx, "link", `{"from":"which","type":"blocks","to":"store"}`)
	r = readyOf(t, ctx, `{}`)
	if !r.Next.IsZero() || len(r.ChangeSets) != 1 ||
		!reflect.DeepEqual(r.ChangeSets[0].Reasons, []string{"blocked by question which on store"}) {
		t.Fatalf("ready = %+v", r)
	}
	if out := human(t, r); !strings.HasPrefix(out, "Nothing is ready.\n"+shortID(b)+" Second: blocked by question which on store\n") {
		t.Errorf("human = %q", out)
	}
	recordWrite(t, ctx, "add", `{"kind":"decision","handle":"this"}`)
	recordWrite(t, ctx, "link", `{"from":"this","type":"answers","to":"which"}`)
	if r = readyOf(t, ctx, `{}`); r.Next != b {
		t.Errorf("ready = %+v", r)
	}
	// as-is has none of B's ops, but its own ops still give its links.
	if r = readyOf(t, ctx, `{"view":"as-is"}`); r.View != "as-is" || r.Next != b {
		t.Errorf("as-is ready = %+v", r)
	}
}

func TestReadyRestsOnEndToEnd(t *testing.T) {
	// B depends on A and sets a field on the node A adds; A is withdrawn.
	ctx := readCtx(t)
	a := newChange(t, ctx, `{"intent":"A"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"base","fields":["statement=Base"]}`)
	b := newChange(t, ctx, `{"intent":"B"}`).ChangeSet
	planLink(t, ctx, op.LinkDependsOn, b, a)
	mustWrite(t, ctx, "set", `{"ref":"base","fields":["statement=Changed"]}`)
	ca := *ctx
	ca.ChangeRef = a.String()
	move(t, &ca, "withdraw", `{}`)
	r := readyOf(t, ctx, `{}`)
	if len(r.ChangeSets) != 1 || r.ChangeSets[0].ID != b {
		t.Fatalf("ready = %+v", r)
	}
	want := []string{"not accepted: draft", "depends on " + shortID(a) + " (withdrawn)",
		"rests on " + shortID(a) + " (withdrawn): 1 operations"}
	if !reflect.DeepEqual(r.ChangeSets[0].Reasons, want) {
		t.Errorf("B reasons = %q, want %q", r.ChangeSets[0].Reasons, want)
	}
}

func TestRoadmapEndToEnd(t *testing.T) {
	ctx := readCtx(t)
	resp, err := callOp(t, ctx, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if out := human(t, resp); !strings.Contains(out, "Nothing planned yet.") {
		t.Errorf("empty roadmap = %q", out)
	}
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"org"}`)
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"grow","parent":"org","json_fields":["milestone=null"]}`)
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"fast","parent":"grow"}`)
	fast := nodeAt(t, ctx, "fast").ID
	a := newChange(t, ctx, `{"intent":"First"}`).ChangeSet
	planLink(t, ctx, op.LinkServes, a, fast)
	move(t, ctx, "propose", `{}`)
	newChange(t, ctx, `{"intent":"Draft"}`)
	resp, err = callOp(t, ctx, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	// No milestone is at or above fast (grow's is null), so First is in
	// the backlog.
	r := resp.(RoadmapResponse)
	if len(r.Groups) != 1 || r.Groups[0].Kind != "backlog" || !r.Groups[0].Intent.ID.IsZero() ||
		r.Groups[0].ChangeSets[0].ID != a || !slices.Equal(r.Groups[0].ChangeSets[0].Serves, []id.ID{fast}) {
		t.Fatalf("roadmap = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "\n## Backlog\n\n- [proposed] "+shortID(a)+" First\n") {
		t.Errorf("human = %q", out)
	}
	// grow, a milestone, is only an ancestor of fast; its group still has
	// its path.
	recordWrite(t, ctx, "set", `{"ref":"grow","fields":["milestone=2"]}`)
	resp, err = callOp(t, ctx, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(RoadmapResponse)
	if g := r.Groups; len(g) != 1 || g[0].Kind != "milestone" || g[0].Intent.Handle != "grow" ||
		!slices.Equal(g[0].Intent.Path, []string{"org"}) || g[0].ChangeSets[0].ID != a {
		t.Fatalf("roadmap = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "\n## org / grow\n\n- [proposed] "+shortID(a)+" First\n") {
		t.Errorf("human = %q", out)
	}
	_, err = callOp(t, ctx, "roadmap", `{"view":"what-if:nothing"}`)
	wantExit(t, err, 4)
}

func TestPlanMilestoneField(t *testing.T) {
	for raw, want := range map[string]bool{`"2"`: true, `1.5`: true, `""`: false, `null`: false, ` null `: false, ``: false} {
		if got := planMilestoneField(map[string]json.RawMessage{"milestone": json.RawMessage(raw)}); got != want {
			t.Errorf("milestone %q = %v, want %v", raw, got, want)
		}
	}
	if planMilestoneField(map[string]json.RawMessage{"statement": json.RawMessage(`"x"`)}) || planMilestoneField(nil) {
		t.Error("no milestone field: is a milestone")
	}
}

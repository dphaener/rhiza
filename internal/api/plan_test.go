package api

import (
	"fmt"
	"reflect"
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

// planFix builds planInputs by hand.
type planFix struct {
	in planInputs
	n  int
}

func newPlanFix() *planFix {
	return &planFix{in: planInputs{ops: map[id.ID][]op.Op{}, nodes: map[id.ID]gmodel.Node{},
		ancestors: map[id.ID][]id.ID{}, intentPath: map[id.ID][]id.ID{}}}
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

// own adds a link.create op to change set cs.
func (f *planFix) own(cs id.ID, t op.LinkType, from, to id.ID) {
	f.in.ops[cs] = append(f.in.ops[cs], op.Op{ID: f.next(), Type: op.TypeLinkCreate, Link: f.next(),
		LinkType: t, From: from, To: to})
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

// verdict runs the core and returns cs's verdict.
func (f *planFix) verdict(t *testing.T, cs id.ID) *planChangeSet {
	t.Helper()
	for _, p := range planCore(f.in).order {
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
	f.link(op.LinkDependsOn, a, landed, false)
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
	f.link(op.LinkDependsOn, b, a, false)
	f.own(b, op.LinkDependsOn, b, gone)
	wantReasons(t, f.verdict(t, b),
		"depends on "+shortID(a)+" (accepted)",
		"depends on "+shortID(gone)+" (not in the log)")
	// A retired depends-on link does not count.
	c := f.cs("C", changeset.StateAccepted)
	f.link(op.LinkDependsOn, c, a, true)
	wantReasons(t, f.verdict(t, c))
}

func TestPlanNotAccepted(t *testing.T) {
	f := newPlanFix()
	a := f.cs("A", changeset.StateAccepted)
	p := f.cs("P", changeset.StateProposed)
	d := f.cs("D", changeset.StateDraft)
	f.link(op.LinkDependsOn, p, a, false)
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
	p := f.cs("P", changeset.StateProposed)
	q := f.cs("Q", changeset.StateProposed)
	f.own(p, op.LinkServes, p, fast)
	f.own(q, op.LinkServes, q, fast)
	o := f.in.ops[q][0]
	f.in.links = append(f.in.links, gmodel.Link{ID: o.Link, Type: o.LinkType, From: q, To: fast, Retired: true})
	res := planCore(f.in)
	if len(res.groups) != 2 || res.groups[0].intent != fast || len(res.groups[0].members) != 1 ||
		res.groups[0].members[0].info.ID != p || !res.groups[1].intent.IsZero() || res.groups[1].members[0].info.ID != q {
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
	a := f.cs("A", changeset.StateAccepted)
	b := f.cs("B", changeset.StateAccepted)
	c := f.cs("C", changeset.StateAccepted)
	for _, cs := range []id.ID{a, b, c} {
		f.link(op.LinkServes, cs, i, false)
	}
	f.link(op.LinkDependsOn, a, c, false)
	got := planGroupIDs(planCore(f.in))
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
	f.link(op.LinkDependsOn, a, b, false)
	f.link(op.LinkDependsOn, b, a, false)
	res := planCore(f.in)
	if got, want := planGroupIDs(res), [][]id.ID{{{}, c, a, b}}; !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	wantReasons(t, res.order[1], "depends on "+shortID(b)+" (accepted)", planCycleReason)
	wantReasons(t, res.order[2], "depends on "+shortID(a)+" (accepted)", planCycleReason)
	wantReasons(t, res.order[0])
}

func TestPlanRoadmapGroupsAndReadyOrder(t *testing.T) {
	// Groups in intent-tree pre-order (siblings by key), retired intents
	// after, No intent last; drafts only in the ready order, after.
	f := newPlanFix()
	late := f.intent("late", id.ID{}, "z")
	early := f.intent("early", id.ID{}, "a")
	child := f.intent("child", early, "m", early)
	gone := f.intent("gone", id.ID{}, "b")
	f.retire(gone)
	draft := f.cs("D", changeset.StateDraft)
	both := f.cs("Both", changeset.StateAccepted)
	alone := f.cs("Alone", changeset.StateProposed)
	l := f.cs("L", changeset.StateAccepted)
	g := f.cs("G", changeset.StateAccepted)
	e := f.cs("E", changeset.StateProposed)
	f.link(op.LinkServes, draft, early, false)
	f.link(op.LinkServes, both, late, false)
	f.link(op.LinkServes, both, child, false)
	f.link(op.LinkServes, l, late, false)
	f.link(op.LinkServes, g, gone, false)
	f.link(op.LinkServes, e, early, false)
	res := planCore(f.in)
	want := [][]id.ID{{early, e}, {child, both}, {late, both, l}, {gone, g}, {{}, alone}}
	if got := planGroupIDs(res); !reflect.DeepEqual(got, want) {
		t.Errorf("groups = %v, want %v", got, want)
	}
	var order []id.ID
	for _, p := range res.order {
		order = append(order, p.info.ID)
	}
	if want := []id.ID{e, both, l, g, alone, draft}; !reflect.DeepEqual(order, want) {
		t.Errorf("order = %v, want %v", order, want)
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

func TestRoadmapEndToEnd(t *testing.T) {
	ctx := readCtx(t)
	resp, err := callOp(t, ctx, "roadmap", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if out := human(t, resp); !strings.Contains(out, "Nothing planned yet.") {
		t.Errorf("empty roadmap = %q", out)
	}
	recordWrite(t, ctx, "add", `{"kind":"intent","handle":"grow"}`)
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
	r := resp.(RoadmapResponse)
	if len(r.Groups) != 1 || r.Groups[0].Intent.Handle != "fast" ||
		!reflect.DeepEqual(r.Groups[0].Intent.Path, []string{"grow"}) || r.Groups[0].ChangeSets[0].ID != a {
		t.Fatalf("roadmap = %+v", r)
	}
	want := "# Roadmap\n\n## grow / fast\n\n- [proposed] " + shortID(a) + " First\n"
	if out := human(t, r); out != want {
		t.Errorf("human = %q, want %q", out, want)
	}
	_, err = callOp(t, ctx, "roadmap", `{"view":"what-if:nothing"}`)
	wantExit(t, err, 4)
}

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"pgregory.net/rapid"
)

// opGen draws valid-shape operations over a pool of nodes and links. Bases
// are drawn mostly from the operations already drawn on the same target, so
// current, stale and unknown bases all occur.
type opGen struct {
	// next is the number of the last op ID handed out: op IDs are fid(next).
	next int
	// onTarget lists the ops drawn so far per target, all of them in drawn.
	onTarget map[id.ID][]id.ID
	drawn    []id.ID
	// handle draws a handle for a node.
	handle func(t *rapid.T, node id.ID) string
}

var (
	propKinds  = []op.Kind{op.KindCapability, op.KindIntent, op.KindDecision, op.KindQuestion}
	propFields = []string{"body", "statement", "title"}
	propKeys   = []string{"8", "i", "r"}
	// propTypes weights updates up, so that conflicts are common.
	propTypes = []op.Type{op.TypeNodeCreate, op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeUpdate,
		op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire, op.TypeLinkCreate, op.TypeLinkRetire}
)

func newOpGen(handle func(t *rapid.T, node id.ID) string) *opGen {
	return &opGen{onTarget: map[id.ID][]id.ID{}, handle: handle}
}

// base draws a base for target: usually one of its ops, sometimes any op
// drawn so far or one that never exists.
func (g *opGen) base(t *rapid.T, target id.ID) id.ID {
	switch k := rapid.IntRange(0, 9).Draw(t, "base from"); {
	case k > 1 && len(g.onTarget[target]) > 0:
		return rapid.SampledFrom(g.onTarget[target]).Draw(t, "base")
	case k == 1 && len(g.drawn) > 0:
		return rapid.SampledFrom(g.drawn).Draw(t, "any base")
	}
	return fid(499)
}

// parent draws the zero parent or a node of the pool.
func parentFrom(t *rapid.T, nodes []id.ID) id.ID {
	if rapid.Bool().Draw(t, "root?") {
		return id.ID{}
	}
	return rapid.SampledFrom(nodes).Draw(t, "parent")
}

// draw returns a valid op on nodes and links whose link endpoints come from
// ends, and records it; ok is false when the drawn op is not valid.
func (g *opGen) draw(t *rapid.T, nodes, links, ends []id.ID) (o op.Op, ok bool) {
	g.next++
	o = op.Op{ID: fid(g.next)}
	switch o.Type = rapid.SampledFrom(propTypes).Draw(t, "type"); o.Type {
	case op.TypeNodeCreate:
		o.Node = rapid.SampledFrom(nodes).Draw(t, "node")
		o.Kind = rapid.SampledFrom(propKinds).Draw(t, "kind")
		o.Handle = g.handle(t, o.Node)
		if op.HasParent(o.Kind) {
			o.Parent = parentFrom(t, nodes)
			o.Key = rapid.SampledFrom(propKeys).Draw(t, "key")
		}
		if rapid.Bool().Draw(t, "fields?") {
			o.Fields = fieldsOf("statement", fmt.Sprintf(`"s%d"`, g.next))
		}
	case op.TypeNodeUpdate:
		o.Node = rapid.SampledFrom(nodes).Draw(t, "node")
		o.Base = g.base(t, o.Node)
		switch rapid.IntRange(0, 3).Draw(t, "update") {
		case 0:
			o.Unset = []string{rapid.SampledFrom(propFields).Draw(t, "unset")}
		case 1:
			o.Handle = g.handle(t, o.Node)
		case 2:
			o.Handle = g.handle(t, o.Node)
			fallthrough
		default:
			// statement twice as often, so that two updates often write it.
			name := rapid.SampledFrom(append([]string{"statement"}, propFields...)).Draw(t, "set")
			o.Set = fieldsOf(name, fmt.Sprintf(`"v%d"`, g.next))
		}
	case op.TypeNodeMove:
		o.Node = rapid.SampledFrom(nodes).Draw(t, "node")
		o.Base = g.base(t, o.Node)
		o.Parent = parentFrom(t, nodes)
		o.Key = rapid.SampledFrom(propKeys).Draw(t, "key")
	case op.TypeNodeRetire:
		o.Node = rapid.SampledFrom(nodes).Draw(t, "node")
		o.Base = g.base(t, o.Node)
		o.Reason = "gone"
	case op.TypeLinkCreate:
		o.Link = rapid.SampledFrom(links).Draw(t, "link")
		o.LinkType = op.LinkServes
		o.From = rapid.SampledFrom(ends).Draw(t, "from")
		o.To = rapid.SampledFrom(ends).Draw(t, "to")
	case op.TypeLinkRetire:
		o.Link = rapid.SampledFrom(links).Draw(t, "link")
		o.Base = g.base(t, o.Link)
	}
	if o.Validate() != nil {
		return op.Op{}, false
	}
	g.record(o)
	return o, true
}

// record notes o as drawn, so later bases can name it.
func (g *opGen) record(o op.Op) {
	g.onTarget[o.Target()] = append(g.onTarget[o.Target()], o.ID)
	g.drawn = append(g.drawn, o.ID)
}

func fieldsOf(name, value string) map[string]json.RawMessage {
	return map[string]json.RawMessage{name: raw(value)}
}

// propStates are the states drawn for change sets after the landed ones.
var propStates = []changeset.State{changeset.StateAccepted, changeset.StateAccepted, changeset.StateProposed,
	changeset.StateDraft, changeset.StateRejected}

// genLog draws 1 to 6 landed change sets and 0 to 3 others (mostly in
// flight) over three nodes and two links. Change set IDs are drawn in a random
// order, so the ID tie-break and cycles both come into play.
func genLog(t *rapid.T) []changeset.ChangeSet {
	landed := rapid.IntRange(1, 6).Draw(t, "landed")
	n := landed + rapid.IntRange(0, 3).Draw(t, "others")
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	nums = rapid.Permutation(nums).Draw(t, "change set IDs")

	nodes := []id.ID{nd(1), nd(2), nd(3)}
	links := []id.ID{ln(1), ln(2)}
	ends := append([]id.ID{}, nodes...)
	for _, k := range nums {
		ends = append(ends, csID(k))
	}
	g := newOpGen(func(t *rapid.T, _ id.ID) string {
		return rapid.SampledFrom([]string{"a", "b", "c"}).Draw(t, "handle")
	})

	sets := make([]changeset.ChangeSet, n)
	for i := range sets {
		s := changeset.ChangeSet{ID: csID(nums[i]), State: changeset.StateLanded}
		if i >= landed {
			s.State = rapid.SampledFrom(propStates).Draw(t, "state")
		}
		if i == 0 && rapid.IntRange(0, 3).Draw(t, "seed?") > 0 {
			// Usually the first set drawn creates every node, so that
			// later updates find them and conflict rather than skip.
			for _, node := range nodes {
				g.next++
				o := op.Op{ID: fid(g.next), Type: op.TypeNodeCreate, Node: node, Kind: op.KindCapability,
					Handle: g.handle(t, node), Key: "i"}
				g.record(o)
				s.Ops = append(s.Ops, o)
			}
		}
		for range rapid.IntRange(1, 5).Draw(t, "ops") {
			if o, ok := g.draw(t, nodes, links, ends); ok {
				s.Ops = append(s.Ops, o)
			}
		}
		sets[i] = s
	}
	return sets
}

// countKinds adds to conflicts and skips the problems of each kind in ps.
func countKinds(ps []Problem, conflicts, skips *int) {
	for _, p := range ps {
		switch {
		case errors.Is(p.Err, ErrConflict):
			*conflicts++
		case errors.Is(p.Err, ErrMissing), errors.Is(p.Err, ErrBase), errors.Is(p.Err, ErrExists),
			errors.Is(p.Err, ErrRetired), errors.Is(p.Err, ErrParent):
			*skips++
		}
	}
}

func sameModel(t *rapid.T, view string, m1, m2 Model, p1, p2 []Problem) {
	t.Helper()
	if !reflect.DeepEqual(m1, m2) {
		t.Fatalf("%s: models differ:\n%+v\n%+v\nvs\n%+v\n%+v", view, m1.Nodes(), m1.Links(), m2.Nodes(), m2.Links())
	}
	if !reflect.DeepEqual(p1, p2) {
		t.Fatalf("%s: problems differ:\n%v\nvs\n%v", view, p1, p2)
	}
}

// TestPropertyDeterminism: every view and Overlaps depend only on the set of
// change sets, never on the order they are given in.
func TestPropertyDeterminism(t *testing.T) {
	var conflicts, skips, inFlightViews int
	rapid.Check(t, func(t *rapid.T) {
		sets := genLog(t)
		perm := rapid.Permutation(sets).Draw(t, "permuted")

		m1, p1 := AsIs(sets)
		m2, p2 := AsIs(perm)
		sameModel(t, "AsIs", m1, m2, p1, p2)
		countKinds(p1, &conflicts, &skips)

		m1, p1 = ToBe(sets)
		m2, p2 = ToBe(perm)
		sameModel(t, "ToBe", m1, m2, p1, p2)
		countKinds(p1, &conflicts, &skips)

		// Every in-flight set, plus a landed one (ErrState) and an unknown
		// one (ErrMissing).
		for _, s := range append(sets, changeset.ChangeSet{ID: csID(99)}) {
			if inFlight(s.State) {
				inFlightViews++
			}
			m1, p1 = WhatIf(sets, s.ID)
			m2, p2 = WhatIf(perm, s.ID)
			sameModel(t, "WhatIf "+s.ID.String(), m1, m2, p1, p2)
		}

		if o1, o2 := Overlaps(sets), Overlaps(perm); !reflect.DeepEqual(o1, o2) {
			t.Fatalf("Overlaps differ:\n%v\nvs\n%v", o1, o2)
		}
	})
	t.Logf("conflicts %d, skips %d, in-flight what-ifs %d", conflicts, skips, inFlightViews)
	if conflicts == 0 || skips == 0 || inFlightViews == 0 {
		t.Errorf("vacuous run: conflicts %d, skips %d, in-flight what-ifs %d", conflicts, skips, inFlightViews)
	}
}

// genDisjoint draws 1 to 6 landed change sets, each over its own nodes and
// links, whose links name only its own nodes or itself. Each set keeps only
// the drawn ops that apply cleanly in order on its own, so the sets have no
// problems; stale bases that merge silently survive. Handles embed the node,
// so no two nodes ever share one.
func genDisjoint(t *rapid.T) []changeset.ChangeSet {
	n := rapid.IntRange(1, 6).Draw(t, "sets")
	known := make([]id.ID, n)
	for i := range known {
		known[i] = csID(i + 1)
	}
	g := newOpGen(func(t *rapid.T, node id.ID) string {
		return fmt.Sprintf("n%x-%d", node[12:], rapid.IntRange(0, 2).Draw(t, "handle"))
	})
	sets := make([]changeset.ChangeSet, n)
	for i := range sets {
		cs := known[i]
		nodes := []id.ID{nd(10*i + 1), nd(10*i + 2), nd(10*i + 3)}
		links := []id.ID{ln(10*i + 1), ln(10*i + 2)}
		ends := append([]id.ID{cs}, nodes...)
		s := changeset.ChangeSet{ID: cs, State: changeset.StateLanded}
		b := newBuilder(known)
		for range rapid.IntRange(1, 10).Draw(t, "ops") {
			o, ok := g.draw(t, nodes, links, ends)
			if !ok {
				continue
			}
			c := b.clone()
			c.apply(cs, o)
			if len(c.problems) == len(b.problems) {
				b = c
				s.Ops = append(s.Ops, o)
			}
		}
		if _, ps := b.finish(); ps != nil {
			t.Fatalf("set %d on its own has problems: %v", i, ps)
		}
		sets[i] = s
	}
	return sets
}

// TestPropertyCommutativity: change sets that touch different nodes and
// links give the same model in any apply order, with no problems.
func TestPropertyCommutativity(t *testing.T) {
	var ops, updates int
	rapid.Check(t, func(t *rapid.T) {
		sets := genDisjoint(t)
		want, ps := AsIs(sets)
		if ps != nil {
			t.Fatalf("AsIs problems: %v", ps)
		}
		for _, s := range sets {
			ops += len(s.Ops)
			for _, o := range s.Ops {
				if o.Type == op.TypeNodeUpdate {
					updates++
				}
			}
		}

		perm := rapid.Permutation(sets).Draw(t, "apply order")
		b := newBuilder(setIDs(sets))
		applySets(b, perm)
		got, ps := b.finish()
		if ps != nil {
			t.Fatalf("problems in order %v: %v", setIDs(perm), ps)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("order %v: model differs:\n%+v\nvs AsIs\n%+v", setIDs(perm), got.Nodes(), want.Nodes())
		}
	})
	t.Logf("ops %d, updates %d", ops, updates)
	if updates == 0 {
		t.Errorf("vacuous run: ops %d, updates %d", ops, updates)
	}
}

// genAny draws anything at all: op types, IDs, bases, parents, endpoints,
// kinds and keys from small pools including the zero ID, so op IDs and
// change set IDs repeat and most shapes are invalid.
func genAny(t *rapid.T) []changeset.ChangeSet {
	pool := []id.ID{{}, fid(1), fid(2), fid(3), fid(4), nd(1), nd(2), nd(3), ln(1), ln(2), csID(1), csID(2), csID(3)}
	pick := func(label string) id.ID { return rapid.SampledFrom(pool).Draw(t, label) }
	types := append(op.Types(), "bogus", "")
	kinds := append(op.Kinds(), "", "bogus")
	states := append(changeset.States(), "")
	strs := []string{"", "a", "b", "0", "i"}

	sets := make([]changeset.ChangeSet, rapid.IntRange(0, 5).Draw(t, "sets"))
	for i := range sets {
		s := changeset.ChangeSet{ID: rapid.SampledFrom([]id.ID{{}, csID(1), csID(2), csID(3)}).Draw(t, "cs"),
			State: rapid.SampledFrom(states).Draw(t, "state")}
		for range rapid.IntRange(0, 6).Draw(t, "ops") {
			o := op.Op{ID: pick("id"), Type: rapid.SampledFrom(types).Draw(t, "type"),
				Node: pick("node"), Link: pick("link"), Base: pick("base"),
				Kind: rapid.SampledFrom(kinds).Draw(t, "kind"), Handle: rapid.SampledFrom(strs).Draw(t, "handle"),
				Parent: pick("parent"), Key: rapid.SampledFrom(strs).Draw(t, "key"),
				Reason: rapid.SampledFrom(strs).Draw(t, "reason"), LinkType: op.LinkServes,
				From: pick("from"), To: pick("to"), Decision: pick("decision")}
			if rapid.Bool().Draw(t, "set?") {
				o.Set = fieldsOf(rapid.SampledFrom(propFields).Draw(t, "set"), `"x"`)
			}
			if rapid.Bool().Draw(t, "unset?") {
				o.Unset = rapid.SliceOfN(rapid.SampledFrom(propFields), 0, 3).Draw(t, "unset")
			}
			if rapid.Bool().Draw(t, "fields?") {
				o.Fields = fieldsOf("statement", `"y"`)
			}
			s.Ops = append(s.Ops, o)
		}
		sets[i] = s
	}
	return sets
}

// TestPropertyNoPanic: no input makes a view panic, and every problem wraps
// exactly one of the package's sentinels.
func TestPropertyNoPanic(t *testing.T) {
	var problems int
	rapid.Check(t, func(t *rapid.T) {
		sets := genAny(t)
		check := func(view string, ps []Problem) {
			for _, p := range ps {
				problems++
				if p.Err == nil {
					t.Fatalf("%s: problem %v has a nil Err", view, p)
				}
				if n := sentinels(p.Err); n != 1 {
					t.Fatalf("%s: problem %v wraps %d sentinels, want 1", view, p, n)
				}
			}
		}
		_, ps := AsIs(sets)
		check("AsIs", ps)
		_, ps = ToBe(sets)
		check("ToBe", ps)
		for _, s := range append(sets, changeset.ChangeSet{ID: csID(99)}) {
			_, ps = WhatIf(sets, s.ID)
			check("WhatIf", ps)
		}
		Overlaps(sets)
	})
	if problems == 0 {
		t.Errorf("vacuous run: no problems seen")
	}
}

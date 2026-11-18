package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// genID is ID n of space s, apart from every fixture ID of fid and csID.
func genID(s byte, n int) id.ID {
	var i id.ID
	i[0] = s
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

// bigLog is a landed log of nodes capabilities, a tenth of them roots and
// the rest spread under them, each with four fields and a second history
// entry, and a relates-to link from each node to its root; and a draft change set
// whose one op updates a field of node 1. It returns the log and the
// draft's ID.
func bigLog(nodes int) ([]changeset.ChangeSet, id.ID) {
	node := func(n int) id.ID { return genID(3, n) }
	link := func(n int) id.ID { return genID(4, n) }
	opID := func(n int) id.ID { return genID(5, n) }
	var ops, updates []op.Op
	roots := nodes / 10
	for n := range nodes {
		var parent id.ID
		if n >= roots {
			parent = node(n % roots)
		}
		fields := map[string]json.RawMessage{}
		for f := range 4 {
			fields[fmt.Sprintf("f%d", f)] = json.RawMessage(fmt.Sprintf("%q", fmt.Sprintf("node %d field %d", n, f)))
		}
		ops = append(ops, op.Op{ID: opID(n), Type: op.TypeNodeCreate, Node: node(n), Kind: op.KindCapability,
			Handle: fmt.Sprintf("n%d", n), Parent: parent, Key: fmt.Sprintf("k%d", n), Fields: fields})
		updates = append(updates, op.Op{ID: opID(nodes + n), Type: op.TypeNodeUpdate, Node: node(n), Base: opID(n),
			Set: map[string]json.RawMessage{"f0": json.RawMessage(`"updated"`)}})
		if n >= roots {
			ops = append(ops, op.Op{ID: opID(3*nodes + n), Type: op.TypeLinkCreate, Link: link(n),
				LinkType: op.LinkRelatesTo, From: node(n), To: parent})
		}
	}
	landed := genID(6, 1)
	draft := genID(6, 2)
	return []changeset.ChangeSet{
		vSet(landed, changeset.StateLanded, append(ops, updates...)...),
		vSet(draft, changeset.StateDraft, op.Op{ID: opID(5 * nodes), Type: op.TypeNodeUpdate, Node: node(1),
			Base: opID(nodes + 1), Set: map[string]json.RawMessage{"f1": json.RawMessage(`"draft"`)}}),
	}, draft
}

// TestWhatIfFromAllocs: on a model of a few thousand nodes, WhatIfFrom
// with a one-op change set allocates at least 80% less than building the
// same view from a deep copy of the seed, as it did before it shared the
// seed's records, and gives the same view.
func TestWhatIfFromAllocs(t *testing.T) {
	sets, draft := bigLog(3000)
	_, _, seed := AsIsSeed(sets)
	deep := func() (Model, []Problem) { return whatIf(seed.builder(false), sets, draft, false) }
	shared := func() (Model, []Problem) { return WhatIfFrom(seed, sets, draft) }
	dm, dps := deep()
	sm, sps := shared()
	sameView(t, "WhatIfFrom", sm, sps, dm, dps)
	if n, _ := sm.Node(genID(3, 1)); string(n.Fields["f1"]) != `"draft"` {
		t.Fatalf("the draft's update is missing: %s", n.Fields)
	}
	before := testing.AllocsPerRun(5, func() { deep() })
	after := testing.AllocsPerRun(5, func() { shared() })
	fall := 100 * (before - after) / before
	t.Logf("WhatIfFrom allocations, 3000 nodes, one op: deep copy %.0f, shared %.0f, %.1f%% fewer", before, after, fall)
	if fall < 80 {
		t.Errorf("allocations fell %.1f%% (deep copy %.0f, shared %.0f), want at least 80%%", fall, before, after)
	}
}

// BenchmarkWhatIfFrom is WhatIfFrom with a one-op change set on a model of
// 3000 nodes.
func BenchmarkWhatIfFrom(b *testing.B) {
	sets, draft := bigLog(3000)
	_, _, seed := AsIsSeed(sets)
	b.ReportAllocs()
	for b.Loop() {
		WhatIfFrom(seed, sets, draft)
	}
}

// TestWhatIfFromSharedSeedUnchanged: two what-if views built from one seed,
// each with ops updating the same node's fields (set and unset), handle,
// place and so its history and aspects, and retiring the same link, leave
// the seed as a fresh AsIsSeed's and the first view as it was before the
// second was built.
func TestWhatIfFromSharedSeedUnchanged(t *testing.T) {
	fields := map[string]json.RawMessage{"a": raw(`"1"`), "b": raw(`"2"`)}
	landed := vSet(csA, changeset.StateLanded,
		mkCap(1, 1, "root", 0, "a"),
		mkCap(2, 2, "other", 0, "b"),
		op.Op{ID: fid(3), Type: op.TypeNodeCreate, Node: nd(3), Kind: op.KindCapability, Handle: "leaf",
			Parent: nd(1), Key: "a", Fields: fields},
		op.Op{ID: fid(4), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkDependsOn, From: nd(3), To: nd(1)},
	)
	edits := func(n int, value string) []op.Op {
		return []op.Op{
			{ID: fid(n), Type: op.TypeNodeUpdate, Node: nd(3), Base: fid(3),
				Set: map[string]json.RawMessage{"a": raw(value)}, Unset: []string{"b"}, Handle: "leaf-" + value[1:2]},
			{ID: fid(n + 1), Type: op.TypeNodeMove, Node: nd(3), Base: fid(n), Parent: nd(2), Key: "z"},
			{ID: fid(n + 2), Type: op.TypeLinkRetire, Link: ln(1), Base: fid(4)},
		}
	}
	sets := []changeset.ChangeSet{
		landed,
		vSet(csB, changeset.StateDraft, edits(10, `"x"`)...),
		vSet(csC, changeset.StateProposed, edits(20, `"y"`)...),
	}
	_, _, seed := AsIsSeed(cloneSets(sets))

	m1, ps1 := WhatIfFrom(seed, sets, csB)
	n1, _ := m1.Node(nd(3))
	l1, _ := m1.Link(ln(1))
	snap := m1.clone()

	m2, ps2 := WhatIfFrom(seed, sets, csC)
	if len(ps1) != 0 || len(ps2) != 0 {
		t.Fatalf("problems: %v; %v", ps1, ps2)
	}
	n2, _ := m2.Node(nd(3))
	l2, _ := m2.Link(ln(1))

	sameView(t, "the first view after the second", m1, nil, snap, nil)
	if _, _, fresh := AsIsSeed(cloneSets(sets)); !reflect.DeepEqual(seed, fresh) {
		t.Fatalf("the what-ifs changed the shared seed:\n%+v\nvs a fresh one\n%+v", seed, fresh)
	}
	for _, c := range []struct {
		name   string
		n      Node
		l      Link
		value  string
		handle string
		op     int
	}{{"first", n1, l1, `"x"`, "leaf-x", 10}, {"second", n2, l2, `"y"`, "leaf-y", 20}} {
		want := map[string]json.RawMessage{"a": raw(c.value)}
		if !reflect.DeepEqual(c.n.Fields, want) || c.n.Handle != c.handle || c.n.Parent != nd(2) ||
			len(c.n.History) != 3 || c.n.Version != fid(c.op+1) {
			t.Errorf("%s view's node: %+v", c.name, c.n)
		}
		if !c.l.Retired || len(c.l.History) != 2 || c.l.Version != fid(c.op+2) {
			t.Errorf("%s view's link: %+v", c.name, c.l)
		}
	}
	asis, _ := AsIs(sets)
	if n, _ := asis.Node(nd(3)); !reflect.DeepEqual(n.Fields, fields) || len(n.History) != 1 {
		t.Errorf("as-is node: %+v", n)
	}
}

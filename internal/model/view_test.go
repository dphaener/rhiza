package model

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// vSet is a folded change set cs in state s with the given operations.
func vSet(cs id.ID, s changeset.State, ops ...op.Op) changeset.ChangeSet {
	return changeset.ChangeSet{ID: cs, State: s, Ops: ops}
}

// notInViews are the states whose change sets are in neither as-is nor
// to-be. The zero State is a change set whose fold had no valid open.
var notInViews = []changeset.State{
	changeset.StateDraft, changeset.StateProposed, changeset.StateRejected,
	changeset.StateWithdrawn, changeset.StateSuperseded, "",
}

func stateName(s changeset.State) string {
	if s == "" {
		return "no state"
	}
	return string(s)
}

// cloneSets deep-copies sets, as far as views could reach into them.
func cloneSets(sets []changeset.ChangeSet) []changeset.ChangeSet {
	out := slices.Clone(sets)
	for i := range out {
		out[i].Ops = slices.Clone(out[i].Ops)
		for j, o := range out[i].Ops {
			out[i].Ops[j].Fields = copyFields(o.Fields)
			out[i].Ops[j].Set = copyFields(o.Set)
			out[i].Ops[j].Unset = slices.Clone(o.Unset)
		}
	}
	return out
}

// sameView fails unless the two results are equal.
func sameView(t *testing.T, what string, m1 Model, ps1 []Problem, m2 Model, ps2 []Problem) {
	t.Helper()
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("%s: models differ:\n%+v\n%+v", what, m1.Nodes(), m2.Nodes())
	}
	if !reflect.DeepEqual(ps1, ps2) {
		t.Errorf("%s: problems differ:\n%v\n%v", what, ps1, ps2)
	}
}

func TestAsIs(t *testing.T) {
	landed := changeset.StateLanded

	for _, s := range append(slices.Clone(notInViews), changeset.StateAccepted) {
		t.Run(stateName(s)+" set is not in as-is", func(t *testing.T) {
			m, ps := AsIs([]changeset.ChangeSet{
				vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
				vSet(csB, s, mkCap(2, 2, "b", 0, ""), set(3, 1, 1, "status", `"x"`)),
			})
			checkProblems(t, ps, nil)
			if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, nds(1)) {
				t.Errorf("nodes = %v, want only node 1", g)
			}
			if x := mustNode(t, m, 1); x.Fields != nil || x.Version != fid(1) {
				t.Errorf("node 1 = %+v, changed by a %s set", x, stateName(s))
			}
		})
	}

	t.Run("version is the landed op ID", func(t *testing.T) {
		m, ps := AsIs([]changeset.ChangeSet{
			vSet(csB, landed, set(2, 1, 1, "status", `"x"`), set(3, 1, 2, "owner", `"me"`)),
			vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
			vSet(csC, changeset.StateAccepted, set(4, 1, 3, "status", `"y"`)),
		})
		checkProblems(t, ps, nil)
		x := mustNode(t, m, 1)
		if x.Version != fid(3) {
			t.Errorf("version = %v, want op 3", x.Version)
		}
		if want := history(r(csA, 1), r(csB, 2), r(csB, 3)); !reflect.DeepEqual(x.History, want) {
			t.Errorf("history = %v, want %v", x.History, want)
		}
		if x.Created != csA {
			t.Errorf("created = %v, want csA", x.Created)
		}
	})

	t.Run("link to a change set ID resolves, a rejected one too", func(t *testing.T) {
		m, ps := AsIs([]changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, ""), mkLink(2, 1, nd(1), csA), mkLink(3, 2, nd(1), csB)),
			vSet(csB, changeset.StateRejected, mkCap(4, 2, "b", 0, "")),
		})
		checkProblems(t, ps, nil)
		if l := mustLink(t, m, 1); l.To != csA {
			t.Errorf("link 1 to = %v, want csA", l.To)
		}
		if l := mustLink(t, m, 2); l.To != csB {
			t.Errorf("link 2 to = %v, want csB", l.To)
		}
	})

	t.Run("link to an unknown change set ID is missing", func(t *testing.T) {
		_, ps := AsIs([]changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, ""), mkLink(2, 1, nd(1), csC)),
		})
		checkProblems(t, ps, []want{{csA, 2, ErrMissing}})
	})

	t.Run("cycle problems surface", func(t *testing.T) {
		// csB builds on csA's op 1 and csA on csB's op 3: a cycle, broken
		// at csA, which then updates a node not yet created.
		_, ps := AsIs([]changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, ""), set(2, 2, 4, "status", `"x"`)),
			vSet(csB, landed, mkCap(3, 2, "b", 0, ""), set(4, 2, 3, "status", `"y"`), set(5, 1, 1, "status", `"z"`)),
		})
		checkProblems(t, ps, []want{{csA, 0, ErrCycle}, {csA, 2, ErrMissing}, {csB, 0, ErrCycle}})
	})

	t.Run("landed sets apply in causal order", func(t *testing.T) {
		m, ps := AsIs([]changeset.ChangeSet{
			vSet(csA, landed, mkCap(2, 2, "child", 1, "a")),
			vSet(csB, landed, mkCap(1, 1, "root", 0, "")),
		})
		checkProblems(t, ps, nil)
		if g := childIDs(m.Children(nd(1))); !reflect.DeepEqual(g, nds(2)) {
			t.Errorf("children = %v, want node 2", g)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		m, ps := AsIs(nil)
		checkProblems(t, ps, nil)
		if m.Nodes() != nil || m.Links() != nil {
			t.Errorf("model not empty")
		}
	})

	t.Run("input order does not matter", func(t *testing.T) {
		sets := []changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, ""), set(2, 2, 4, "status", `"x"`)),
			vSet(csB, landed, mkCap(3, 2, "b", 0, ""), set(4, 2, 3, "status", `"y"`), set(5, 1, 1, "status", `"z"`)),
			vSet(csC, landed, mkCap(6, 3, "a", 0, ""), mkLink(7, 1, nd(3), csID(4))),
			vSet(csID(4), changeset.StateWithdrawn, mkCap(8, 4, "d", 0, "")),
		}
		m0, ps0 := AsIs(sets)
		rnd := rand.New(rand.NewPCG(3, 4))
		for range 30 {
			in := slices.Clone(sets)
			rnd.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
			m, ps := AsIs(in)
			sameView(t, "shuffled", m0, ps0, m, ps)
		}
	})
}

func TestToBe(t *testing.T) {
	landed, accepted := changeset.StateLanded, changeset.StateAccepted

	for _, s := range notInViews {
		t.Run(stateName(s)+" set is not in to-be", func(t *testing.T) {
			m, ps := ToBe([]changeset.ChangeSet{
				vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
				vSet(csB, s, mkCap(2, 2, "b", 0, ""), set(3, 1, 1, "status", `"x"`)),
			})
			checkProblems(t, ps, nil)
			if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, nds(1)) {
				t.Errorf("nodes = %v, want only node 1", g)
			}
		})
	}

	t.Run("accepted set is in to-be", func(t *testing.T) {
		m, ps := ToBe([]changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
			vSet(csB, accepted, mkCap(2, 2, "b", 1, "a"), set(3, 1, 1, "status", `"x"`)),
			vSet(csC, changeset.StateProposed, set(4, 1, 1, "owner", `"me"`)),
		})
		checkProblems(t, ps, nil)
		if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, nds(1, 2)) {
			t.Errorf("nodes = %v, want 1 and 2", g)
		}
		x := mustNode(t, m, 1)
		if x.Version != fid(3) || string(x.Fields["status"]) != `"x"` || x.Fields["owner"] != nil {
			t.Errorf("node 1 = %+v, want the accepted update only", x)
		}
	})

	t.Run("accepted sets apply in their own causal order", func(t *testing.T) {
		m, ps := ToBe([]changeset.ChangeSet{
			vSet(csA, accepted, set(3, 2, 2, "status", `"x"`)),
			vSet(csB, accepted, mkCap(2, 2, "b", 1, "a")),
			vSet(csC, landed, mkCap(1, 1, "a", 0, "")),
		})
		checkProblems(t, ps, nil)
		if x := mustNode(t, m, 2); x.Version != fid(3) {
			t.Errorf("node 2 = %+v, want version op 3", x)
		}
	})

	t.Run("stale accepted set conflicts in to-be but not in as-is", func(t *testing.T) {
		sets := []changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
			vSet(csB, landed, set(2, 1, 1, "status", `"x"`)),
			vSet(csC, accepted, set(3, 1, 1, "status", `"y"`), set(4, 1, 1, "owner", `"me"`)),
		}
		_, ps := AsIs(sets)
		checkProblems(t, ps, nil)
		m, ps := ToBe(sets)
		checkProblems(t, ps, []want{{csC, 3, ErrConflict}})
		x := mustNode(t, m, 1)
		if string(x.Fields["status"]) != `"y"` || string(x.Fields["owner"]) != `"me"` {
			t.Errorf("node 1 fields = %v, want the accepted writes applied", x.Fields)
		}
	})

	t.Run("as-is problems are not duplicated", func(t *testing.T) {
		sets := []changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "x", 0, ""), set(2, 2, 4, "status", `"x"`)),
			vSet(csB, landed, mkCap(3, 2, "x", 0, ""), set(4, 2, 3, "status", `"y"`), set(5, 1, 1, "status", `"z"`)),
			vSet(csC, accepted, set(6, 1, 5, "owner", `"me"`), mkCap(7, 3, "x", 0, "")),
		}
		_, asIs := AsIs(sets)
		checkProblems(t, asIs, []want{{csA, 0, ErrCycle}, {csA, 2, ErrMissing}, {csB, 0, ErrCycle}, {csB, 3, ErrHandle}})
		_, ps := ToBe(sets)
		checkProblems(t, ps, []want{
			{csA, 0, ErrCycle}, {csA, 2, ErrMissing}, {csB, 0, ErrCycle}, {csB, 3, ErrHandle}, {csC, 7, ErrHandle},
		})
	})

	t.Run("input order does not matter", func(t *testing.T) {
		sets := []changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
			vSet(csB, landed, set(2, 1, 1, "status", `"x"`)),
			vSet(csC, accepted, set(3, 1, 1, "status", `"y"`)),
			vSet(csID(4), accepted, mkCap(4, 2, "b", 1, "a"), mkLink(5, 1, nd(2), csC)),
			vSet(csID(5), changeset.StateProposed, set(6, 1, 2, "owner", `"me"`)),
		}
		m0, ps0 := ToBe(sets)
		rnd := rand.New(rand.NewPCG(5, 6))
		for range 30 {
			in := slices.Clone(sets)
			rnd.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
			m, ps := ToBe(in)
			sameView(t, "shuffled", m0, ps0, m, ps)
		}
	})

	t.Run("views share no state", func(t *testing.T) {
		sets := []changeset.ChangeSet{
			vSet(csA, landed, mkDec(1, 1, "jwt"), set(2, 1, 1, "status", `"x"`)),
			vSet(csB, accepted, set(3, 1, 2, "status", `"y"`), mkDec(4, 2, "b")),
			vSet(csC, changeset.StateDraft, set(5, 1, 2, "owner", `"me"`)),
		}
		in := cloneSets(sets)
		asIs, asIsPs := AsIs(sets)
		toBe, _ := ToBe(sets)
		whatIf, _ := WhatIf(sets, csC)
		if !reflect.DeepEqual(sets, in) {
			t.Fatalf("views changed their input")
		}

		// Views built from the same input are independent of each other.
		fresh, freshPs := AsIs(sets)
		sameView(t, "as-is after to-be and what-if", asIs, asIsPs, fresh, freshPs)
		if field(t, toBe, 1, "status") != `"y"` || field(t, whatIf, 1, "owner") != `"me"` {
			t.Errorf("to-be or what-if wrong")
		}
		if field(t, asIs, 1, "owner") != "" || field(t, asIs, 1, "status") != `"x"` {
			t.Errorf("as-is changed by another view")
		}

		// Changing the input afterwards changes no view.
		sets[0].Ops[0].Fields["statement"][1] = 'X'
		sets[0].Ops[1].Set["status"][1] = 'Q'
		sets[1].Ops[0].Set["status"][1] = 'Q'
		if got := field(t, asIs, 1, "statement"); got != `"use JWT"` {
			t.Errorf("statement = %s, changed through the input", got)
		}
		if got := field(t, asIs, 1, "status"); got != `"x"` {
			t.Errorf("as-is status = %s, changed through the input", got)
		}
		if got := field(t, toBe, 1, "status"); got != `"y"` {
			t.Errorf("to-be status = %s, changed through the input", got)
		}
	})
}

func TestWhatIf(t *testing.T) {
	landed, accepted := changeset.StateLanded, changeset.StateAccepted
	csD, csE := csID(4), csID(5)
	sets := []changeset.ChangeSet{
		vSet(csA, landed, mkCap(1, 1, "a", 0, "")),
		vSet(csB, changeset.StateDraft, mkCap(2, 2, "b", 1, "a"), set(3, 1, 1, "status", `"draft"`)),
		vSet(csC, accepted, set(4, 1, 1, "owner", `"me"`)),
		vSet(csD, changeset.StateProposed, set(5, 1, 1, "status", `"proposed"`), mkLink(6, 1, nd(1), csE)),
		vSet(csE, changeset.StateRejected, set(7, 1, 1, "status", `"rejected"`)),
		vSet(csID(6), landed, set(8, 1, 1, "landed", `true`)),
		vSet(csID(7), changeset.StateWithdrawn),
		vSet(csID(8), changeset.StateSuperseded),
		vSet(csID(9), ""),
	}
	asIs, asIsPs := AsIs(sets)
	checkProblems(t, asIsPs, nil)

	t.Run("draft shows its ops on as-is", func(t *testing.T) {
		m, ps := WhatIf(sets, csB)
		checkProblems(t, ps, nil)
		if g := childIDs(m.Children(nd(1))); !reflect.DeepEqual(g, nds(2)) {
			t.Errorf("children = %v, want node 2", g)
		}
		x := mustNode(t, m, 1)
		if string(x.Fields["status"]) != `"draft"` || string(x.Fields["landed"]) != "true" || x.Fields["owner"] != nil {
			t.Errorf("node 1 fields = %v, want as-is plus the draft only", x.Fields)
		}
	})

	t.Run("proposed shows its ops and its link to a rejected set", func(t *testing.T) {
		m, ps := WhatIf(sets, csD)
		checkProblems(t, ps, nil)
		if got := field(t, m, 1, "status"); got != `"proposed"` {
			t.Errorf("status = %s", got)
		}
		if l := mustLink(t, m, 1); l.To != csE {
			t.Errorf("link to = %v, want csE", l.To)
		}
	})

	t.Run("accepted shows its ops alone", func(t *testing.T) {
		m, ps := WhatIf(sets, csC)
		checkProblems(t, ps, nil)
		x := mustNode(t, m, 1)
		if string(x.Fields["owner"]) != `"me"` || x.Fields["status"] != nil || x.Version != fid(4) {
			t.Errorf("node 1 = %+v, want as-is plus csC only", x)
		}
	})

	t.Run("stale what-if conflicts", func(t *testing.T) {
		in := append(slices.Clone(sets), vSet(csID(10), changeset.StateProposed, set(9, 1, 1, "landed", `false`)))
		_, ps := WhatIf(in, csID(10))
		checkProblems(t, ps, []want{{csID(10), 9, ErrConflict}})
	})

	for _, tc := range []struct {
		name string
		cs   id.ID
		err  error
	}{
		{"landed set gives as-is and ErrState", csID(6), ErrState},
		{"rejected set gives as-is and ErrState", csE, ErrState},
		{"withdrawn set gives as-is and ErrState", csID(7), ErrState},
		{"superseded set gives as-is and ErrState", csID(8), ErrState},
		{"set with no state gives as-is and ErrState", csID(9), ErrState},
		{"unknown set gives as-is and ErrMissing", csID(99), ErrMissing},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m, ps := WhatIf(sets, tc.cs)
			if !reflect.DeepEqual(m, asIs) {
				t.Errorf("model is not as-is:\n%+v\n%+v", m.Nodes(), asIs.Nodes())
			}
			checkProblems(t, ps, []want{{tc.cs, 0, tc.err}})
		})
	}

	t.Run("what-if keeps as-is problems", func(t *testing.T) {
		in := []changeset.ChangeSet{
			vSet(csA, landed, mkCap(1, 1, "x", 0, ""), mkCap(2, 2, "x", 0, "")),
			vSet(csB, changeset.StateDraft, set(3, 1, 1, "status", `"x"`)),
		}
		_, ps := WhatIf(in, csB)
		checkProblems(t, ps, []want{{csA, 2, ErrHandle}})
		_, ps = WhatIf(in, csC)
		checkProblems(t, ps, []want{{csA, 2, ErrHandle}, {csC, 0, ErrMissing}})
	})

	t.Run("input order does not matter", func(t *testing.T) {
		for _, cs := range []id.ID{csB, csC, csD, csID(6), csID(99)} {
			m0, ps0 := WhatIf(sets, cs)
			rnd := rand.New(rand.NewPCG(7, 8))
			for range 10 {
				in := slices.Clone(sets)
				rnd.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
				m, ps := WhatIf(in, cs)
				sameView(t, "shuffled", m0, ps0, m, ps)
			}
		}
	})
}

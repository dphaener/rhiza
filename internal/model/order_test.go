package model

import (
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// orderSet is a folded change set holding ops, for ordering.
func orderSet(i id.ID, ops ...op.Op) changeset.ChangeSet {
	return changeset.ChangeSet{ID: i, State: changeset.StateLanded, Ops: ops}
}

func setIDs(sets []changeset.ChangeSet) []id.ID {
	var out []id.ID
	for _, s := range sets {
		out = append(out, s.ID)
	}
	return out
}

func TestOrder(t *testing.T) {
	csD := csID(4)
	tests := []struct {
		name  string
		sets  []changeset.ChangeSet
		want  []id.ID
		probs []want
	}{
		{
			name: "independent sets in ID order",
			sets: []changeset.ChangeSet{
				orderSet(csC, mkCap(3, 3, "c", 0, "")),
				orderSet(csA, mkCap(1, 1, "a", 0, "")),
				orderSet(csB, mkCap(2, 2, "b", 0, "")),
			},
			want: []id.ID{csA, csB, csC},
		},
		{
			name: "base edge puts a smaller ID after a larger one",
			sets: []changeset.ChangeSet{
				orderSet(csA, set(2, 1, 1, "status", `"x"`)),
				orderSet(csB, mkCap(1, 1, "a", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "base edge through a chain of three",
			sets: []changeset.ChangeSet{
				orderSet(csA, set(3, 1, 2, "status", `"y"`)),
				orderSet(csB, set(2, 1, 1, "status", `"x"`)),
				orderSet(csC, mkCap(1, 1, "a", 0, "")),
			},
			want: []id.ID{csC, csB, csA},
		},
		{
			name: "existence edge: link to a node created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkLink(2, 1, nd(1), csA)),
				orderSet(csB, mkCap(1, 1, "a", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: link from a node created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkLink(2, 1, csA, nd(1))),
				orderSet(csB, mkCap(1, 1, "a", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: child under a parent created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(2, 2, "child", 1, "a")),
				orderSet(csB, mkCap(1, 1, "parent", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: move under a parent created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(1, 1, "a", 0, ""), move(3, 1, 1, 2, "a")),
				orderSet(csB, mkCap(2, 2, "b", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: update of a node created elsewhere",
			sets: []changeset.ChangeSet{
				// The base is not an op of csB, so only the target orders.
				orderSet(csA, set(2, 1, 99, "status", `"x"`)),
				orderSet(csB, mkCap(1, 1, "a", 0, "")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: retire of a link created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, retireLink(2, 1, 99)),
				orderSet(csB, mkLink(1, 1, csA, csB)),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "existence edge: decision created elsewhere",
			sets: []changeset.ChangeSet{
				orderSet(csA, op.Op{ID: fid(2), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkServes,
					From: csA, To: csC, Decision: nd(1)}),
				orderSet(csB, mkDec(1, 1, "jwt")),
			},
			want: []id.ID{csB, csA},
		},
		{
			name: "references within one set add no edge",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(1, 1, "a", 0, ""), set(2, 1, 1, "status", `"x"`)),
				orderSet(csB, mkCap(3, 2, "b", 0, "")),
			},
			want: []id.ID{csA, csB},
		},
		{
			name: "unknown references add no edge",
			sets: []changeset.ChangeSet{
				orderSet(csB, set(2, 1, 1, "status", `"x"`)),
				orderSet(csA, mkLink(3, 1, nd(9), nd(8))),
			},
			want: []id.ID{csA, csB},
		},
		{
			name: "duplicate create: the smallest change set owns",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(1, 1, "a", 0, "")),
				orderSet(csB, mkCap(2, 1, "a", 0, "")),
				orderSet(csC, set(3, 1, 99, "status", `"x"`)),
			},
			// csB names nd(1), which csA owns, so csB follows csA too.
			want: []id.ID{csA, csB, csC},
		},
		{
			name: "duplicate create owned by the smaller ID orders the larger after it",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(2, 1, "a", 0, ""), set(4, 2, 3, "status", `"x"`)),
				orderSet(csB, mkCap(1, 1, "a", 0, ""), mkCap(3, 2, "b", 0, "")),
			},
			// csB creates nd(1), which csA owns; csA updates from csB's op
			// fid(3): a cycle, broken at csA.
			want:  []id.ID{csA, csB},
			probs: []want{{csA, 0, ErrCycle}, {csB, 0, ErrCycle}},
		},
		{
			name: "two-set cycle reported on both and broken at the smaller ID",
			sets: []changeset.ChangeSet{
				orderSet(csB, mkCap(2, 2, "b", 0, ""), set(4, 1, 1, "status", `"x"`)),
				orderSet(csA, mkCap(1, 1, "a", 0, ""), set(3, 2, 2, "status", `"y"`)),
			},
			want:  []id.ID{csA, csB},
			probs: []want{{csA, 0, ErrCycle}, {csB, 0, ErrCycle}},
		},
		{
			name: "cycle reports sets blocked behind it and orders the rest",
			sets: []changeset.ChangeSet{
				orderSet(csA, mkCap(1, 1, "a", 0, "")),
				orderSet(csB, mkCap(2, 2, "b", 0, ""), set(4, 3, 3, "status", `"x"`)),
				orderSet(csC, mkCap(3, 3, "c", 0, ""), set(5, 2, 2, "status", `"y"`)),
				orderSet(csD, set(6, 2, 4, "status", `"z"`), set(7, 1, 1, "status", `"w"`)),
			},
			want:  []id.ID{csA, csB, csC, csD},
			probs: []want{{csB, 0, ErrCycle}, {csC, 0, ErrCycle}, {csD, 0, ErrCycle}},
		},
		{
			name: "cycle broken at the smallest, then order resumes",
			sets: []changeset.ChangeSet{
				// csA and csC form a cycle; csB depends on csC only.
				orderSet(csA, mkCap(1, 1, "a", 0, ""), set(4, 3, 3, "status", `"x"`)),
				orderSet(csB, set(5, 3, 3, "status", `"y"`)),
				orderSet(csC, mkCap(3, 3, "c", 0, ""), set(6, 1, 1, "status", `"z"`)),
			},
			want:  []id.ID{csA, csC, csB},
			probs: []want{{csA, 0, ErrCycle}, {csB, 0, ErrCycle}, {csC, 0, ErrCycle}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.sets)
			got, probs := order(tt.sets)
			if !reflect.DeepEqual(tt.sets, in) {
				t.Fatalf("order changed its input")
			}
			if g := setIDs(got); !reflect.DeepEqual(g, tt.want) {
				t.Errorf("order = %v, want %v", g, tt.want)
			}
			checkProblems(t, probs, tt.probs)
			for _, s := range got {
				i := slices.IndexFunc(tt.sets, func(x changeset.ChangeSet) bool { return x.ID == s.ID })
				if !reflect.DeepEqual(s, tt.sets[i]) {
					t.Errorf("change set %s changed by ordering", s.ID)
				}
			}
		})
	}

	t.Run("empty input gives nil, nil", func(t *testing.T) {
		got, probs := order(nil)
		if got != nil || probs != nil {
			t.Fatalf("order(nil) = %v, %v; want nil, nil", got, probs)
		}
		got, probs = order([]changeset.ChangeSet{})
		if got != nil || probs != nil {
			t.Fatalf("order(empty) = %v, %v; want nil, nil", got, probs)
		}
	})

	t.Run("shuffled input gives the same order and problems", func(t *testing.T) {
		csD, csE := csID(4), csID(5)
		sets := []changeset.ChangeSet{
			orderSet(csA, mkCap(1, 1, "a", 0, ""), set(4, 3, 3, "status", `"x"`)),
			orderSet(csB, set(5, 3, 3, "status", `"y"`), mkLink(8, 1, nd(1), nd(2))),
			orderSet(csC, mkCap(3, 3, "c", 0, ""), set(6, 1, 1, "status", `"z"`)),
			orderSet(csD, mkCap(2, 2, "d", 1, "a")),
			orderSet(csE, mkCap(7, 4, "e", 0, "")),
		}
		wantSets, wantProbs := order(sets)
		// csE is ready at once. csA and csC wait on each other, csD on csA,
		// and csB on csA, csC and csD (its link names nd(2)). The cycle is
		// broken at csA, which frees csC and csD.
		if g := setIDs(wantSets); !reflect.DeepEqual(g, []id.ID{csE, csA, csC, csD, csB}) {
			t.Fatalf("order = %v, want [E A C D B]", g)
		}
		checkProblems(t, wantProbs, []want{{csA, 0, ErrCycle}, {csB, 0, ErrCycle}, {csC, 0, ErrCycle}, {csD, 0, ErrCycle}})
		r := rand.New(rand.NewPCG(1, 2))
		for range 50 {
			in := slices.Clone(sets)
			r.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
			got, probs := order(in)
			if !reflect.DeepEqual(got, wantSets) || !reflect.DeepEqual(probs, wantProbs) {
				t.Fatalf("order(%v) = %v %v, want %v %v", setIDs(in), setIDs(got), probs, setIDs(wantSets), wantProbs)
			}
		}
	})
}

package model

import (
	"cmp"
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// ovSet is a folded change set cs in state s with the given operations.
func ovSet(cs id.ID, s changeset.State, ops ...op.Op) changeset.ChangeSet {
	return changeset.ChangeSet{ID: cs, State: s, Ops: ops}
}

func TestOverlaps(t *testing.T) {
	csD := csID(4)
	prop, acc := changeset.StateProposed, changeset.StateAccepted
	tests := []struct {
		name string
		sets []changeset.ChangeSet
		want []Overlap
	}{
		{
			name: "same field",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, prop, set(2, 1, 0, "title", `"b"`)),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title"}}},
		},
		{
			name: "same field from different bases",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 7, "title", `"a"`)),
				ovSet(csB, acc, set(2, 1, 8, "title", `"b"`)),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title"}}},
		},
		{
			name: "only the shared fields",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, op.Op{ID: fid(1), Type: op.TypeNodeUpdate, Node: nd(1), Base: fid(0),
					Set: map[string]json.RawMessage{"title": raw(`"a"`), "body": raw(`"x"`)}}),
				ovSet(csB, prop, unset(2, 1, 0, "title", "owner"), rename(3, 1, 0, "h")),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title"}}},
		},
		{
			name: "different fields",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, acc, set(2, 1, 0, "body", `"b"`)),
			},
		},
		{
			name: "different targets",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, prop, set(2, 2, 0, "title", `"b"`)),
			},
		},
		{
			name: "handle",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, rename(1, 1, 0, "x")),
				ovSet(csB, prop, rename(2, 1, 0, "y")),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"handle"}}},
		},
		{
			name: "move vs move",
			sets: []changeset.ChangeSet{
				ovSet(csA, acc, move(1, 1, 0, 2, "a")),
				ovSet(csB, acc, move(2, 1, 0, 3, "b")),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"position"}}},
		},
		{
			name: "move vs update",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, move(1, 1, 0, 2, "a")),
				ovSet(csB, prop, set(2, 1, 0, "title", `"b"`)),
			},
		},
		{
			name: "retire vs update",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`), rename(2, 1, 1, "h")),
				ovSet(csB, acc, retire(3, 1, 0)),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title", "handle", "retired"}}},
		},
		{
			name: "retire vs retire",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, retire(1, 1, 0)),
				ovSet(csB, prop, retire(2, 1, 0)),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"retired"}}},
		},
		{
			name: "link.retire vs link.retire",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, retireLink(1, 1, 0)),
				ovSet(csB, acc, retireLink(2, 1, 0)),
			},
			want: []Overlap{{A: csA, B: csB, Target: ln(1), Aspects: []string{"retired"}}},
		},
		{
			name: "retire vs link create to the node",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, mkLink(1, 1, nd(2), nd(1))),
				ovSet(csB, acc, retire(2, 1, 0)),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"linked", "retired"}}},
		},
		{
			name: "retire vs link create from the node",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, retire(1, 1, 0)),
				ovSet(csB, prop, mkLink(2, 1, nd(1), nd(2))),
			},
			want: []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"linked", "retired"}}},
		},
		{
			name: "retire and link in one set",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, retire(1, 1, 0), mkLink(2, 1, nd(1), nd(2))),
			},
		},
		{
			name: "two link creates to one node",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, mkLink(1, 1, nd(2), nd(1))),
				ovSet(csB, prop, mkLink(2, 2, nd(3), nd(1))),
			},
		},
		{
			name: "link create vs update",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, mkLink(1, 1, nd(2), nd(1))),
				ovSet(csB, prop, set(2, 1, 0, "title", `"b"`)),
			},
		},
		{
			name: "creates never overlap",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, mkCap(1, 1, "x", 0, "a"), mkLink(2, 1, nd(1), nd(1))),
				ovSet(csB, prop, mkCap(3, 1, "x", 0, "a"), mkLink(4, 1, nd(1), nd(1)), set(5, 1, 3, "title", `"b"`)),
			},
		},
		{
			name: "one set alone",
			sets: []changeset.ChangeSet{
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`), set(2, 1, 1, "title", `"b"`), retire(3, 1, 2)),
			},
		},
		{
			name: "draft ignored",
			sets: []changeset.ChangeSet{
				ovSet(csA, changeset.StateDraft, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, prop, set(2, 1, 0, "title", `"b"`)),
				ovSet(csC, changeset.StateDraft, retire(3, 1, 0)),
			},
		},
		{
			name: "terminal sets ignored",
			sets: []changeset.ChangeSet{
				ovSet(csA, changeset.StateLanded, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, changeset.StateRejected, set(2, 1, 0, "title", `"b"`)),
				ovSet(csC, changeset.StateWithdrawn, set(3, 1, 0, "title", `"c"`)),
				ovSet(csD, changeset.StateSuperseded, set(4, 1, 0, "title", `"d"`)),
				ovSet(csID(5), acc, set(5, 1, 0, "title", `"e"`)),
				ovSet(csID(6), "", set(6, 1, 0, "title", `"f"`)),
			},
		},
		{
			name: "three-way overlap gives three pairs",
			sets: []changeset.ChangeSet{
				ovSet(csC, prop, set(3, 1, 0, "title", `"c"`)),
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`)),
				ovSet(csB, acc, set(2, 1, 0, "title", `"b"`)),
			},
			want: []Overlap{
				{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title"}},
				{A: csA, B: csC, Target: nd(1), Aspects: []string{"field:title"}},
				{A: csB, B: csC, Target: nd(1), Aspects: []string{"field:title"}},
			},
		},
		{
			name: "one per pair and target, sorted by target",
			sets: []changeset.ChangeSet{
				ovSet(csB, prop, move(4, 2, 0, 3, "b"), set(5, 1, 0, "title", `"b"`), set(6, 1, 5, "body", `"b"`)),
				ovSet(csA, prop, set(1, 1, 0, "title", `"a"`), move(2, 2, 0, 4, "a"), set(3, 1, 1, "body", `"a"`)),
			},
			want: []Overlap{
				{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:body", "field:title"}},
				{A: csA, B: csB, Target: nd(2), Aspects: []string{"position"}},
			},
		},
		{
			name: "empty",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Overlaps(tc.sets)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Overlaps = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestOverlapsInputOrder(t *testing.T) {
	sets := []changeset.ChangeSet{
		ovSet(csA, changeset.StateProposed, set(1, 1, 0, "title", `"a"`), retireLink(2, 1, 0), move(3, 2, 0, 3, "a")),
		ovSet(csB, changeset.StateAccepted, retire(4, 1, 0), retireLink(5, 1, 0)),
		ovSet(csC, changeset.StateProposed, move(6, 2, 0, 4, "c"), set(7, 1, 0, "body", `"c"`)),
		ovSet(csID(4), changeset.StateDraft, retire(8, 2, 0)),
	}
	want := Overlaps(sets)
	if len(want) != 4 {
		t.Fatalf("Overlaps = %v, want 4 overlaps", want)
	}
	if !slices.IsSortedFunc(want, func(x, y Overlap) int {
		return cmp.Or(id.Compare(x.A, y.A), id.Compare(x.B, y.B), id.Compare(x.Target, y.Target))
	}) {
		t.Fatalf("Overlaps = %v, not sorted", want)
	}
	for _, perm := range ovPermutations(len(sets)) {
		in := make([]changeset.ChangeSet, len(sets))
		for i, j := range perm {
			in[i] = sets[j]
		}
		if got := Overlaps(in); !reflect.DeepEqual(got, want) {
			t.Fatalf("order %v: Overlaps = %v, want %v", perm, got, want)
		}
	}
}

// ovPermutations returns every ordering of 0..n-1.
func ovPermutations(n int) [][]int {
	if n == 0 {
		return [][]int{nil}
	}
	var out [][]int
	for _, p := range ovPermutations(n - 1) {
		for i := 0; i <= len(p); i++ {
			out = append(out, slices.Insert(slices.Clone(p), i, n-1))
		}
	}
	return out
}

func TestOverlapsNamedDraft(t *testing.T) {
	sets := []changeset.ChangeSet{
		ovSet(csA, changeset.StateDraft, set(1, 1, 0, "title", `"a"`)),
		ovSet(csB, changeset.StateProposed, set(2, 1, 0, "title", `"b"`)),
		ovSet(csC, changeset.StateDraft, set(3, 1, 0, "title", `"c"`)),
	}
	if got := Overlaps(sets); got != nil {
		t.Fatalf("no draft named: %v", got)
	}
	want := []Overlap{{A: csA, B: csB, Target: nd(1), Aspects: []string{"field:title"}}}
	if got := Overlaps(sets, csA); !reflect.DeepEqual(got, want) {
		t.Fatalf("draft A named: %v, want %v", got, want)
	}
	// Naming a proposed set changes nothing; naming an absent one neither.
	if got := Overlaps(sets, csB, csID(9)); got != nil {
		t.Fatalf("proposed named: %v", got)
	}
}

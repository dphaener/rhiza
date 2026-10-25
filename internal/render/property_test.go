package render

import (
	"bytes"
	"maps"
	"slices"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/model"
	"pgregory.net/rapid"
)

// TestPropertyPermutation checks that the order of the change sets given to
// Render and AsIs changes nothing: every permutation renders the same bytes,
// for the hand-built model and every corpus case.
func TestPropertyPermutation(t *testing.T) {
	cases := map[string][]changeset.ChangeSet{"hand-built": handBuilt()}
	for name, root := range corpusCases(t) {
		cases[name] = corpusSets(t, root)
	}
	for name, sets := range cases {
		t.Run(name, func(t *testing.T) {
			want := renderSets(t, sets)
			rapid.Check(t, func(rt *rapid.T) {
				perm := rapid.Permutation(sets).Draw(rt, "perm")
				m, _ := model.AsIs(perm)
				got := Render(m, perm)
				if !maps.EqualFunc(got, want, bytes.Equal) {
					keys := slices.Sorted(maps.Keys(got))
					rt.Fatalf("permutation renders differently; paths %q", keys)
				}
			})
		})
	}
}

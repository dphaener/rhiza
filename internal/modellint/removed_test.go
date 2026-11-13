package modellint

import (
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func removed(t *testing.T, other, written []op.Op, fsys fstest.MapFS, deleted ...string) []string {
	t.Helper()
	sets := []changeset.ChangeSet{
		{ID: nid(200), State: changeset.StateLanded, Ops: other},
		{ID: nid(201), State: changeset.StateAccepted, Ops: written},
	}
	m, ps := model.ToBe(sets)
	if len(ps) != 0 {
		t.Fatalf("problems: %v", ps)
	}
	var out []string
	for _, h := range RemovedCode(m, sets[1], fsys, deleted) {
		out = append(out, h.Handle+"/"+h.Kind+"/"+h.Ref)
	}
	return out
}

func TestRemovedCode(t *testing.T) {
	mine := []op.Op{capOp(9, "mine", "Does y.", "")}
	for _, tc := range []struct {
		name    string
		other   []op.Op
		wrote   []op.Op
		deleted []string
		want    []string
	}{
		{"only path deleted", []op.Op{capOp(1, "old", "Lives in internal/gone/x.go.", "")}, mine,
			[]string{"internal/gone/x.go"}, []string{"old/removed-code-capability/internal/gone/x.go"}},
		{"directory deleted", []op.Op{capOp(1, "old", "Lives in internal/gone.", "")}, mine,
			[]string{"internal/gone/x.go"}, []string{"old/removed-code-capability/internal/gone"}},
		{"another path present", []op.Op{capOp(1, "old", "Lives in internal/gone/x.go and internal/api/run.go.", "")}, mine,
			[]string{"internal/gone/x.go"}, nil},
		{"missing but not deleted", []op.Op{capOp(1, "old", "Lives in internal/gone/x.go.", "")}, mine,
			[]string{"internal/other.go"}, nil},
		{"no path references", []op.Op{capOp(1, "old", "Does x.", "")}, mine,
			[]string{"internal/gone/x.go"}, nil},
		{"retired", []op.Op{capOp(1, "old", "Lives in internal/gone/x.go.", ""),
			{ID: nid(150), Type: op.TypeNodeRetire, Node: nid(1), Base: nid(101), Reason: "gone"}}, mine,
			[]string{"internal/gone/x.go"}, nil},
		{"written by the change set", nil,
			[]op.Op{capOp(1, "old", "Lives in internal/gone/x.go.", "")},
			[]string{"internal/gone/x.go"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := removed(t, tc.other, tc.wrote, head(), tc.deleted...); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("hits = %v, want %v", got, tc.want)
			}
		})
	}
}

package modellint

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func decOp(n byte, handle string) op.Op {
	return op.Op{ID: nid(100 + n), Type: op.TypeNodeCreate, Node: nid(n), Kind: op.KindDecision, Handle: handle,
		Fields: map[string]json.RawMessage{"statement": js("Decided.")}}
}

func linkOp(n byte, typ op.LinkType, from, to byte) op.Op {
	return op.Op{ID: nid(100 + n), Type: op.TypeLinkCreate, Link: nid(n), LinkType: typ, From: nid(from), To: nid(to)}
}

func qOp(n byte, handle string) op.Op {
	return op.Op{ID: nid(100 + n), Type: op.TypeNodeCreate, Node: nid(n), Kind: op.KindQuestion, Handle: handle,
		Fields: map[string]json.RawMessage{"statement": js("Why?")}}
}

func unapplied(t *testing.T, other, written []op.Op, committed ...string) []string {
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
	for _, h := range UnappliedDecisions(m, sets[1], committed) {
		out = append(out, h.Handle+"/"+h.Kind)
	}
	return out
}

func TestUnappliedDecisions(t *testing.T) {
	one := []string{"d/unapplied-decision"}
	for _, tc := range []struct {
		name         string
		other, wrote []op.Op
		committed    []string
		want         []string
	}{
		{"no link", nil, []op.Op{decOp(1, "d")}, nil, one},
		{"answers only", nil, []op.Op{decOp(1, "d"), qOp(2, "q"), linkOp(3, op.LinkAnswers, 1, 2)}, nil, nil},
		{"landed capability applies", []op.Op{capOp(2, "old", "x", "")},
			[]op.Op{decOp(1, "d"), linkOp(3, op.LinkDecidedBy, 2, 1)}, nil, nil},
		{"leaf without commit", nil,
			[]op.Op{decOp(1, "d"), capOp(2, "leaf", "x", ""), linkOp(3, op.LinkDecidedBy, 2, 1)}, nil, one},
		{"leaf committed by handle", nil,
			[]op.Op{decOp(1, "d"), capOp(2, "leaf", "x", ""), linkOp(3, op.LinkDecidedBy, 2, 1)}, []string{"leaf"}, nil},
		{"leaf committed by id", nil,
			[]op.Op{decOp(1, "d"), capOp(2, "leaf", "x", ""), linkOp(3, op.LinkDecidedBy, 2, 1)}, []string{nid(2).String()}, nil},
		{"other change set's decision", []op.Op{decOp(1, "d")}, []op.Op{capOp(2, "leaf", "x", "")}, nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := unapplied(t, tc.other, tc.wrote, tc.committed...); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("hits = %v, want %v", got, tc.want)
			}
		})
	}
}

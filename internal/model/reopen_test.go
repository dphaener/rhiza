package model

import (
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	ent "github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// logEntry is log entry entry(n) of csA with body b after entry(n-1), so
// the entries form one chain.
func logEntry(t *testing.T, n int, b changeset.Body) ent.Entry {
	t.Helper()
	body, err := changeset.EncodeBody(b)
	if err != nil {
		t.Fatalf("entry %d: %v", n, err)
	}
	after := []id.ID{}
	if n > 1 {
		after = []id.ID{entry(n - 1)}
	}
	return ent.Entry{V: changeset.EntryVersion(b), ID: entry(n), ChangeSet: csA, Actor: fid(900), For: fid(901),
		At: time.Unix(1_800_000_000+int64(n), 0).UTC(), After: after, Type: b.Type(), Body: body}
}

// TestViewsReopen folds a change set landed, reopened, revised and relanded,
// and checks the views at each step: as-is has its ops only while landed,
// to-be has them at every step, what-if works only while it is reopened,
// and Lifecycle lists every move so far.
func TestViewsReopen(t *testing.T) {
	es := []ent.Entry{
		logEntry(t, 1, changeset.Open{Kind: changeset.KindIntention, Intent: "i"}),
		logEntry(t, 2, changeset.Ops{Ops: []op.Op{mkCap(1, 1, "a", 0, "i")}}),
		logEntry(t, 3, changeset.Propose{}),
		logEntry(t, 4, changeset.Accept{Rule: "owner-approves"}),
		logEntry(t, 5, changeset.Land{}),
		logEntry(t, 6, changeset.Reopen{}),
		logEntry(t, 7, changeset.Ops{Ops: []op.Op{mkCap(2, 2, "b", 0, "j")}}),
		logEntry(t, 8, changeset.Land{}),
	}
	landed := []ent.Type{ent.TypePropose, ent.TypeAccept, ent.TypeLand}
	steps := []struct {
		name      string
		upTo      int
		state     changeset.State
		nodes     []id.ID
		lifecycle []ent.Type
	}{
		{"landed", 5, changeset.StateLanded, nds(1), landed},
		{"reopened", 6, changeset.StateAccepted, nds(1), append(slices.Clone(landed), ent.TypeReopen)},
		{"reopened with an op", 7, changeset.StateAccepted, nds(1, 2), append(slices.Clone(landed), ent.TypeReopen)},
		{"relanded", 8, changeset.StateLanded, nds(1, 2), append(slices.Clone(landed), ent.TypeReopen, ent.TypeLand)},
	}
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			cs, fps := changeset.Fold(csA, es[:s.upTo])
			if fps != nil {
				t.Fatalf("fold problems = %v", fps)
			}
			if cs.State != s.state {
				t.Fatalf("State = %s, want %s", cs.State, s.state)
			}
			var types []ent.Type
			for _, m := range cs.Lifecycle {
				types = append(types, m.Type)
			}
			if !reflect.DeepEqual(types, s.lifecycle) {
				t.Errorf("Lifecycle = %v, want %v", types, s.lifecycle)
			}
			sets := []changeset.ChangeSet{cs}

			m, ps := AsIs(sets)
			checkProblems(t, ps, nil)
			wantAsIs := []id.ID(nil)
			if s.state == changeset.StateLanded {
				wantAsIs = s.nodes
			}
			if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, wantAsIs) {
				t.Errorf("as-is nodes = %v, want %v", g, wantAsIs)
			}

			m, ps = ToBe(sets)
			checkProblems(t, ps, nil)
			if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, s.nodes) {
				t.Errorf("to-be nodes = %v, want %v", g, s.nodes)
			}

			m, ps = WhatIf(sets, csA)
			if s.state == changeset.StateLanded {
				if len(ps) != 1 || !errors.Is(ps[0].Err, ErrState) {
					t.Errorf("what-if problems = %v, want one ErrState", ps)
				}
				return
			}
			checkProblems(t, ps, nil)
			if g := childIDs(m.Nodes()); !reflect.DeepEqual(g, s.nodes) {
				t.Errorf("what-if nodes = %v, want %v", g, s.nodes)
			}
		})
	}
}

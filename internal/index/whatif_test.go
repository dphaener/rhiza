package index

import (
	"errors"
	"os"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

// land writes a landed change set of one ops entry.
func land(f *fixture, cs id.ID, ops ...op.Op) {
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "landed"})
	f.add(cs, changeset.Ops{Ops: ops})
	f.add(cs, changeset.Propose{})
	f.add(cs, changeset.Accept{Rule: "owner"})
	f.add(cs, changeset.Land{})
}

// draft writes a draft change set of one ops entry.
func draft(f *fixture, cs id.ID, ops ...op.Op) {
	f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "draft"})
	f.add(cs, changeset.Ops{Ops: ops})
}

func dependsOn(n int, from, to id.ID) op.Op {
	return op.Op{ID: opID(n), Type: op.TypeLinkCreate, Link: ln(n), LinkType: op.LinkDependsOn, From: from, To: to}
}

// TestWhatIfRepairsLandedConflict: a conflict left in as-is that a draft's
// write from the current version repairs is gone from its what-if, as in
// the full fold.
func TestWhatIfRepairsLandedConflict(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	land(f, csID(2), update(10, 1, opID(1), "a"))
	land(f, csID(3), update(11, 1, opID(1), "b"))
	draft(f, csID(4), update(12, 1, opID(11), "c"))

	x := openT(t, f, f.dbPath(), nil)
	sets, _ := changeset.FoldLog(store.ReadLog(os.DirFS(f.gdir)))
	_, asis := model.AsIs(sets)
	if len(asis) != 1 {
		t.Fatalf("want one as-is conflict, got %v", asis)
	}
	_, ps, err := x.WhatIf(csID(4))
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Fatalf("what-if problems %v, want the conflict repaired", ps)
	}
	checkModel(t, x, f.gdir)
}

// TestWhatIfLandedSkippedCreate: a dependency's dropped create is no
// conflict for a change set whose view still has a landed create of the
// node, even one the model skipped.
func TestWhatIfLandedSkippedCreate(t *testing.T) {
	for _, skipped := range []bool{false, true} {
		f := newFixture(t)
		landedSet(f, 1)
		if skipped {
			land(f, csID(2), mkCap(21, 5, "five", 99, "i")) // parent 99 is missing: skipped
		}
		a, b := csID(3), csID(4)
		f.add(a, changeset.Open{Kind: changeset.KindIntention, Intent: "a"})
		f.add(a, changeset.Ops{Ops: []op.Op{mkCap(20, 5, "five", 0, "i")}})
		f.add(b, changeset.Open{Kind: changeset.KindIntention, Intent: "b"})
		f.add(b, changeset.Ops{Ops: []op.Op{dependsOn(7, b, a), mkCap(30, 6, "six", 5, "i")}})
		f.add(a, changeset.Ops{Drop: []id.ID{opID(20)}})

		x := openT(t, f, f.dbPath(), nil)
		_, ps, err := x.WhatIf(b)
		if err != nil {
			t.Fatal(err)
		}
		dropped := 0
		for _, p := range ps {
			if p.Op == opID(30) && errors.Is(p.Err, model.ErrConflict) {
				dropped++
			}
		}
		if want := map[bool]int{false: 1, true: 0}[skipped]; dropped != want {
			t.Errorf("skipped %v: %d conflicts on op 30, want %d (%v)", skipped, dropped, want, ps)
		}
		checkModel(t, x, f.gdir)
	}
}

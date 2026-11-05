package index

import (
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// A signature over some change sets moves with their entries and with the
// change sets that are no longer in flight, and with nothing else; Refresh
// catches up and reports what it did.
func TestSignature(t *testing.T) {
	f := newFixture(t)
	landedSet(f, 1)
	mine, other := csID(2), csID(3)
	for _, cs := range []id.ID{mine, other} {
		f.add(cs, changeset.Open{Kind: changeset.KindIntention, Intent: "draft"})
	}
	x := openT(t, f, f.dbPath(), nil)
	sig := func() string {
		t.Helper()
		s, err := x.Signature([]id.ID{mine})
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	refresh := func() CatchUp {
		t.Helper()
		if err := x.Refresh(); err != nil {
			t.Fatal(err)
		}
		return x.LastCatchUp()
	}
	stamp := func() string {
		t.Helper()
		s, err := x.Stamp()
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	if c := refresh(); !c.Fresh || c.Read != 0 {
		t.Fatalf("Refresh of a fresh index did %+v", c)
	}
	base, st := sig(), stamp()

	f.add(other, changeset.Ops{Ops: []op.Op{mkCap(10, 10, "elsewhere", 0, "i")}})
	if c := refresh(); c.Fresh || c.Read != 1 || c.Folded != 1 {
		t.Fatalf("Refresh after an entry did %+v, want that change set alone", c)
	}
	if stamp() == st {
		t.Error("the stamp did not move with an entry")
	}
	if sig() != base {
		t.Error("an entry in another change set in flight moved the signature")
	}

	f.add(mine, changeset.Ops{Ops: []op.Op{mkCap(11, 11, "mine", 0, "i")}})
	refresh()
	if sig() == base {
		t.Error("an entry of the change set left the signature as it was")
	}
	base = sig()

	f.add(other, changeset.Propose{})
	f.add(other, changeset.Accept{Rule: "owner"})
	f.add(other, changeset.Land{})
	refresh()
	if sig() == base {
		t.Error("a change set landing left the signature as it was")
	}
}

package changeset

import (
	"reflect"
	"testing"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// Fixture nodes and links for the rests-on tests.
var (
	restNode  = fid(600) // created by A
	restNode2 = fid(601) // created by B
	restLink  = fid(602) // created by B
	restOther = fid(603) // created before any of these change sets
)

func restCreate(n int, node id.ID) op.Op {
	return op.Op{ID: opID(n), Type: op.TypeNodeCreate, Node: node, Kind: op.KindCapability, Handle: "h", Key: "i"}
}

func restUpdate(n int, node, base id.ID) op.Op {
	return op.Op{ID: opID(n), Type: op.TypeNodeUpdate, Node: node, Base: base}
}

func restSet(cs id.ID, s State, ops ...op.Op) ChangeSet {
	return ChangeSet{ID: cs, State: s, Ops: ops}
}

func TestRestsOn(t *testing.T) {
	csC := csID(3)
	aCreates := restSet(csA, StateWithdrawn, restCreate(1, restNode))
	bUpdates := restSet(csB, StateDraft, restUpdate(10, restNode, opID(1)))
	// bTargets names A's node only, not A's op.
	bTargets := restSet(csB, StateDraft, restUpdate(10, restNode, opBase))
	cases := []struct {
		name string
		sets []ChangeSet
		want []Resting
	}{
		{"withdrawn", []ChangeSet{aCreates, bUpdates},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"rejected", []ChangeSet{restSet(csA, StateRejected, restCreate(1, restNode)), bUpdates},
			[]Resting{{ChangeSet: csA, State: StateRejected, Ops: opIDsOf(10)}}},
		{"superseded", []ChangeSet{restSet(csA, StateSuperseded, restCreate(1, restNode)), bUpdates},
			[]Resting{{ChangeSet: csA, State: StateSuperseded, Ops: opIDsOf(10)}}},
		{"landed", []ChangeSet{restSet(csA, StateLanded, restCreate(1, restNode)), bUpdates}, nil},
		{"in flight", []ChangeSet{restSet(csA, StateAccepted, restCreate(1, restNode)), bUpdates}, nil},
		{"parent", []ChangeSet{aCreates, restSet(csB, StateDraft,
			op.Op{ID: opID(10), Type: op.TypeNodeCreate, Node: restNode2, Kind: op.KindCapability,
				Handle: "child", Parent: restNode, Key: "i"})},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"link to", []ChangeSet{aCreates, restSet(csB, StateDraft,
			op.Op{ID: opID(10), Type: op.TypeLinkCreate, Link: restLink, LinkType: op.LinkServes,
				From: restOther, To: restNode})},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"link from", []ChangeSet{aCreates, restSet(csB, StateDraft,
			op.Op{ID: opID(10), Type: op.TypeLinkCreate, Link: restLink, LinkType: op.LinkServes,
				From: restNode, To: restOther})},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"base names an op of the dropped set", []ChangeSet{
			restSet(csA, StateWithdrawn, restUpdate(1, restOther, opBase)),
			restSet(csB, StateDraft, restUpdate(10, restOther, opID(1)))},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"own creates", []ChangeSet{aCreates, restSet(csB, StateDraft,
			restCreate(10, restNode2), restUpdate(11, restNode2, opID(10)))}, nil},
		{"target only", []ChangeSet{aCreates, bTargets},
			[]Resting{{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)}}},
		{"also created by an in-flight set", []ChangeSet{aCreates, bTargets,
			restSet(csC, StateProposed, restCreate(20, restNode))}, nil},
		{"also created by a landed set", []ChangeSet{aCreates, bTargets,
			restSet(csC, StateLanded, restCreate(20, restNode))}, nil},
		{"created again by another dropped set", []ChangeSet{aCreates, bTargets,
			restSet(csC, StateRejected, restCreate(20, restNode))},
			[]Resting{
				{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10)},
				{ChangeSet: csC, State: StateRejected, Ops: opIDsOf(10)},
			}},
		{"sorted by change set then op", []ChangeSet{
			restSet(csC, StateRejected, restCreate(20, restNode2)),
			restSet(csB, StateDraft,
				restUpdate(12, restNode2, opID(20)), restUpdate(11, restNode, opID(1)), restUpdate(10, restNode, opID(1))),
			aCreates},
			[]Resting{
				{ChangeSet: csA, State: StateWithdrawn, Ops: opIDsOf(10, 11)},
				{ChangeSet: csC, State: StateRejected, Ops: opIDsOf(12)},
			}},
		{"unrelated", []ChangeSet{aCreates, restSet(csB, StateDraft, restUpdate(10, restOther, opBase))}, nil},
		{"unknown change set", []ChangeSet{aCreates}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RestsOn(c.sets, csB)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("RestsOn = %+v, want %+v", got, c.want)
			}
		})
	}
}

// TestRestsOnSelf: a dropped change set's own ops never rest on itself.
func TestRestsOnSelf(t *testing.T) {
	a := restSet(csA, StateWithdrawn, restCreate(1, restNode), restUpdate(2, restNode, opID(1)))
	if got := RestsOn([]ChangeSet{a}, csA); got != nil {
		t.Fatalf("RestsOn = %+v, want nil", got)
	}
}

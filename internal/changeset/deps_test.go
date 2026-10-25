package changeset

import (
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// dep is a depends-on link.create with op ID fid(600+n) and link ID
// fid(650+n), from change set from to change set to.
func dep(n int, from, to id.ID) op.Op {
	return op.Op{ID: fid(600 + n), Type: op.TypeLinkCreate, Link: fid(650 + n), LinkType: op.LinkDependsOn, From: from, To: to}
}

// undep retires the link dep(n, ...) created, with op ID fid(620+n).
func undep(n int) op.Op {
	return op.Op{ID: fid(620 + n), Type: op.TypeLinkRetire, Link: fid(650 + n), Base: fid(600 + n)}
}

// folded is change set cs in state s with ops.
func folded(cs id.ID, s State, ops ...op.Op) ChangeSet {
	return ChangeSet{ID: cs, State: s, Ops: ops}
}

func idList(xs ...id.ID) []id.ID { return xs }

var csC, csD, csE = csID(3), csID(4), csID(5)

func TestDepsDependsOn(t *testing.T) {
	relates := dep(9, csA, csB)
	relates.LinkType = op.LinkRelatesTo
	tests := []struct {
		name string
		sets []ChangeSet
		cs   id.ID
		want []id.ID
	}{
		{"none", []ChangeSet{folded(csA, StateDraft)}, csA, nil},
		{"landed link counts", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateDraft), folded(csC, StateLanded, dep(1, csA, csB)),
		}, csA, idList(csB)},
		{"own link counts", []ChangeSet{
			folded(csA, StateProposed, dep(1, csA, csB)), folded(csB, StateDraft),
		}, csA, idList(csB)},
		{"own retire cancels own link", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB), undep(1)), folded(csB, StateDraft),
		}, csA, nil},
		{"landed retire cancels landed link", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateDraft),
			folded(csC, StateLanded, dep(1, csA, csB)), folded(csD, StateLanded, undep(1)),
		}, csA, nil},
		{"own retire cancels landed link", []ChangeSet{
			folded(csA, StateDraft, undep(1)), folded(csB, StateDraft), folded(csC, StateLanded, dep(1, csA, csB)),
		}, csA, nil},
		{"in-flight retire elsewhere does not cancel", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateDraft),
			folded(csC, StateLanded, dep(1, csA, csB)), folded(csD, StateAccepted, undep(1)),
		}, csA, idList(csB)},
		{"link in a third in-flight change set does not count", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateDraft), folded(csC, StateAccepted, dep(1, csA, csB)),
		}, csA, nil},
		{"link in a dropped change set does not count", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateDraft), folded(csC, StateRejected, dep(1, csA, csB)),
		}, csA, nil},
		{"link from another change set does not count", []ChangeSet{
			folded(csA, StateDraft, dep(1, csB, csC)), folded(csB, StateDraft), folded(csC, StateDraft),
		}, csA, nil},
		{"other link types do not count", []ChangeSet{
			folded(csA, StateDraft, relates), folded(csB, StateDraft),
		}, csA, nil},
		{"dangling target kept", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csE)),
		}, csA, idList(csE)},
		{"self excluded", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csA)),
		}, csA, nil},
		{"sorted and once", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csD), dep(2, csA, csB)),
			folded(csC, StateLanded, dep(3, csA, csB), dep(4, csA, csC)),
		}, csA, idList(csB, csC, csD)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DependsOn(tt.sets, tt.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DependsOn = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepsDropped(t *testing.T) {
	want := map[State]bool{StateWithdrawn: true, StateRejected: true, StateSuperseded: true}
	for _, s := range States() {
		if got := Dropped(s); got != want[s] {
			t.Errorf("Dropped(%s) = %v, want %v", s, got, want[s])
		}
	}
}

func TestClosure(t *testing.T) {
	tests := []struct {
		name string
		sets []ChangeSet
		cs   id.ID
		want []id.ID
	}{
		{"no dependencies", []ChangeSet{folded(csA, StateDraft)}, csA, nil},
		{"chain through draft, proposed and accepted", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateProposed, dep(2, csB, csC)),
			folded(csC, StateAccepted, dep(3, csC, csD)),
			folded(csD, StateDraft),
		}, csA, idList(csB, csC, csD)},
		{"stops at dropped", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB), dep(2, csA, csD)),
			folded(csB, StateWithdrawn, dep(3, csB, csC)),
			folded(csC, StateDraft),
			folded(csD, StateDraft),
		}, csA, idList(csD)},
		{"rejected and superseded are dropped too", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB), dep(2, csA, csC)),
			folded(csB, StateRejected), folded(csC, StateSuperseded),
		}, csA, nil},
		{"stops at landed", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateLanded, dep(2, csB, csC)),
			folded(csC, StateDraft),
		}, csA, nil},
		{"not in the log skipped", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csE), dep(2, csA, csB)), folded(csB, StateDraft),
		}, csA, idList(csB)},
		{"tolerates a cycle", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateDraft, dep(2, csB, csC)),
			folded(csC, StateDraft, dep(3, csC, csA)),
		}, csA, idList(csB, csC)},
		{"diamond, each once", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csC), dep(2, csA, csB)),
			folded(csB, StateDraft, dep(3, csB, csD)),
			folded(csC, StateDraft, dep(4, csC, csD)),
			folded(csD, StateProposed),
		}, csA, idList(csB, csC, csD)},
		{"landed dependency link reaches in-flight", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateAccepted),
			folded(csC, StateLanded, dep(1, csA, csB)),
		}, csA, idList(csB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Closure(tt.sets, tt.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Closure = %v, want %v", got, tt.want)
			}
			rev := slices.Clone(tt.sets)
			slices.Reverse(rev)
			if got := Closure(rev, tt.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Closure reversed = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDependencies(t *testing.T) {
	tests := []struct {
		name string
		sets []ChangeSet
		want []id.ID
	}{
		{"no dependencies", []ChangeSet{folded(csA, StateDraft)}, nil},
		{"goes through landed", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateLanded, dep(2, csB, csC)),
			folded(csC, StateLanded),
		}, idList(csB, csC)},
		{"landed and in flight", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB), dep(2, csA, csC)),
			folded(csB, StateLanded), folded(csC, StateAccepted, dep(3, csC, csD)),
			folded(csD, StateLanded),
		}, idList(csB, csC, csD)},
		{"stops at dropped", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB), dep(2, csA, csD)),
			folded(csB, StateWithdrawn, dep(3, csB, csC)),
			folded(csC, StateLanded), folded(csD, StateRejected),
		}, nil},
		{"not in the log skipped", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csE), dep(2, csA, csB)), folded(csB, StateLanded),
		}, idList(csB)},
		{"tolerates a cycle", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateDraft, dep(2, csB, csA)),
		}, idList(csB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dependencies(tt.sets, csA); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dependencies = %v, want %v", got, tt.want)
			}
			rev := slices.Clone(tt.sets)
			slices.Reverse(rev)
			if got := Dependencies(rev, csA); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dependencies reversed = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWouldCycle(t *testing.T) {
	chain := []ChangeSet{
		folded(csA, StateDraft, dep(1, csA, csB)),
		folded(csB, StateLanded, dep(2, csB, csC)),
		folded(csC, StateWithdrawn),
		folded(csD, StateDraft),
	}
	tests := []struct {
		name     string
		from, to id.ID
		want     bool
	}{
		{"self", csD, csD, true},
		{"closes a loop through every state", csC, csA, true},
		{"direct back edge", csB, csA, true},
		{"forward edge", csA, csC, false},
		{"unrelated", csD, csA, false},
		{"into unrelated", csA, csD, false},
		{"not in the log", csA, csE, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WouldCycle(chain, tt.from, tt.to); got != tt.want {
				t.Errorf("WouldCycle(%v, %v) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestCycles(t *testing.T) {
	tests := []struct {
		name string
		sets []ChangeSet
		want [][]id.ID
	}{
		{"none", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)), folded(csB, StateDraft),
		}, nil},
		{"self-loop", []ChangeSet{
			folded(csA, StateDraft), folded(csB, StateLanded, dep(1, csA, csA)),
		}, [][]id.ID{idList(csA)}},
		{"three-cycle", []ChangeSet{
			folded(csC, StateDraft, dep(1, csC, csA)),
			folded(csA, StateProposed, dep(2, csA, csB)),
			folded(csB, StateDraft, dep(3, csB, csC)),
			folded(csD, StateDraft, dep(4, csD, csA)),
		}, [][]id.ID{idList(csA, csB, csC)}},
		{"two cycles across states", []ChangeSet{
			folded(csA, StateDraft, dep(1, csA, csB)),
			folded(csB, StateRejected, dep(2, csB, csA)),
			folded(csC, StateLanded, dep(3, csC, csD), dep(4, csE, csE)),
			folded(csD, StateAccepted, dep(5, csD, csC)),
			folded(csE, StateDraft),
		}, [][]id.ID{idList(csA, csB), idList(csC, csD), idList(csE)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cycles(tt.sets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cycles = %v, want %v", got, tt.want)
			}
			rev := slices.Clone(tt.sets)
			slices.Reverse(rev)
			if got := Cycles(rev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cycles reversed = %v, want %v", got, tt.want)
			}
		})
	}
}

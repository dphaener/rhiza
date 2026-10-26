package changeset

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// revNode is the node A creates and B writes in the revised tests.
var revNode = fid(500)

// fields builds a field map from name, JSON value pairs.
func fields(kv ...string) map[string]json.RawMessage {
	m := map[string]json.RawMessage{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i]] = json.RawMessage(kv[i+1])
	}
	return m
}

// revCreate is A's create of revNode, op opID(1).
func revCreate(handle, key string, fs map[string]json.RawMessage) op.Op {
	return op.Op{ID: opID(1), Type: op.TypeNodeCreate, Node: revNode, Kind: op.KindCapability,
		Handle: handle, Key: key, Fields: fs}
}

// revUpdate is an update of revNode with ID opID(n), based on A's create.
func revUpdate(n int, set map[string]json.RawMessage, unset ...string) op.Op {
	return op.Op{ID: opID(n), Type: op.TypeNodeUpdate, Node: revNode, Base: opID(1), Set: set, Unset: unset}
}

// step is entry n with body b.
type step struct {
	n int
	b Body
}

// chain folds change set cs from steps, each entry after the one before.
func chain(t *testing.T, cs id.ID, steps ...step) ChangeSet {
	t.Helper()
	var es []entry.Entry
	for i, s := range steps {
		if i == 0 {
			es = append(es, mk(t, cs, s.n, s.b))
			continue
		}
		es = append(es, mk(t, cs, s.n, s.b, steps[i-1].n))
	}
	got, ps := Fold(cs, es)
	checkProblems(t, ps, cs, nil)
	return got
}

func ops(os ...op.Op) Ops { return Ops{Ops: os} }

func TestRevised(t *testing.T) {
	statement := fields("statement", `"a"`)
	created := revCreate("h", "i", statement)
	bWrites := revUpdate(10, fields("statement", `"b"`))
	aOpen := step{1, intention()}
	aCreate := step{2, ops(created)}
	bOpen := step{3, intention()}
	bWrite := step{4, ops(dep(1, csB, csA), bWrites)}
	conflict := func(aspects ...string) []RevisedConflict {
		return []RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
			Target: revNode, Aspects: aspects}}
	}
	dropped := []RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
		Target: revNode, Dropped: true}}

	cases := []struct {
		name string
		a, b []step
		want []RevisedConflict
	}{
		{"overlap",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "i", fields("statement", `"a2"`)))}},
			[]step{bOpen, bWrite}, conflict("field:statement")},
		{"disjoint fields",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "i", fields("statement", `"a"`, "title", `"t"`)))}},
			[]step{bOpen, bWrite}, nil},
		{"revision before the write",
			[]step{aOpen, aCreate, {3, ops(revCreate("h", "i", fields("statement", `"a2"`)))}},
			[]step{{4, intention()}, {5, ops(dep(1, csB, csA), bWrites)}}, nil},
		{"rewritten after the revision",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "i", fields("statement", `"a2"`)))}},
			[]step{bOpen, bWrite, {6, ops(bWrites)}}, nil},
		{"dropped after the write",
			[]step{aOpen, aCreate, {5, Ops{Drop: []id.ID{opID(1)}}}},
			[]step{bOpen, bWrite}, dropped},
		{"dropped before the write",
			[]step{aOpen, aCreate, {3, Ops{Drop: []id.ID{opID(1)}}}},
			[]step{{4, intention()}, {5, ops(dep(1, csB, csA), bWrites)}}, nil},
		{"written and dropped after the write",
			[]step{aOpen, {5, ops(created)}, {6, Ops{Drop: []id.ID{opID(1)}}}},
			[]step{bOpen, bWrite}, nil},
		{"written after the write",
			[]step{aOpen, {5, ops(created)}, {6, ops(revCreate("h", "i", fields("statement", `"a2"`)))}},
			[]step{bOpen, bWrite}, nil},
		{"not_a_dependency",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "i", fields("statement", `"a2"`)))}},
			[]step{bOpen, {4, ops(bWrites)}}, nil},
		{"handle",
			[]step{aOpen, aCreate, {5, ops(revCreate("h2", "i", statement))}},
			[]step{bOpen, {4, ops(dep(1, csB, csA), op.Op{ID: opID(10), Type: op.TypeNodeUpdate,
				Node: revNode, Base: opID(1), Handle: "h3"})}}, conflict("handle")},
		{"handle against a field write",
			[]step{aOpen, aCreate, {5, ops(revCreate("h2", "i", statement))}},
			[]step{bOpen, bWrite}, nil},
		{"position",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "j", statement))}},
			[]step{bOpen, {4, ops(dep(1, csB, csA), op.Op{ID: opID(10), Type: op.TypeNodeMove,
				Node: revNode, Base: opID(1), Key: "k"})}}, conflict("position")},
		{"retire overlaps any change",
			[]step{aOpen, aCreate, {5, ops(revCreate("h2", "i", fields("statement", `"a2"`)))}},
			[]step{bOpen, {4, ops(dep(1, csB, csA), op.Op{ID: opID(10), Type: op.TypeNodeRetire,
				Node: revNode, Base: opID(1), Reason: "gone"})}}, conflict("field:statement", "handle")},
		{"update set then unset",
			[]step{aOpen, {2, ops(created, revUpdate(2, statement))},
				{5, ops(revUpdate(2, nil, "statement"))}},
			[]step{bOpen, bWrite}, []RevisedConflict{
				{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(2),
					Target: revNode, Aspects: []string{"field:statement"}}}},
		{"withdrawn",
			[]step{aOpen, aCreate, {5, ops(revCreate("h", "i", fields("statement", `"a2"`)))},
				{6, Withdraw{}}},
			[]step{bOpen, bWrite}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sets := []ChangeSet{chain(t, csA, c.a...), chain(t, csB, c.b...)}
			got := Revised(sets, csB)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Revised = %+v, want %+v", got, c.want)
			}
		})
	}
}

// TestRevisedLanded: a dependency's revision or drop under B's write still
// counts once the dependency has landed, until B writes its op again. The
// fold keeps a landed change set's History and Drops; only the state is
// set here.
func TestRevisedLanded(t *testing.T) {
	created := revCreate("h", "i", fields("statement", `"a"`))
	bWrites := revUpdate(10, fields("statement", `"b"`))
	bSteps := []step{{3, intention()}, {4, ops(dep(1, csB, csA), bWrites)}}
	landed := func(s ChangeSet) ChangeSet { s.State = StateLanded; return s }
	cases := []struct {
		name string
		a    ChangeSet
		b    []step
		want []RevisedConflict
	}{
		{"revised", landed(chain(t, csA, step{1, intention()}, step{2, ops(created)},
			step{5, ops(revCreate("h", "i", fields("statement", `"a2"`)))})), bSteps,
			[]RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
				Target: revNode, Aspects: []string{"field:statement"}}}},
		{"dropped", landed(chain(t, csA, step{1, intention()}, step{2, ops(created)},
			step{5, Ops{Drop: []id.ID{opID(1)}}})), bSteps,
			[]RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
				Target: revNode, Dropped: true}}},
		{"rewritten after the revision", landed(chain(t, csA, step{1, intention()}, step{2, ops(created)},
			step{5, ops(revCreate("h", "i", fields("statement", `"a2"`)))})),
			append(slices.Clone(bSteps), step{6, ops(bWrites)}), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Revised([]ChangeSet{c.a, chain(t, csB, c.b...)}, csB)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Revised = %+v, want %+v", got, c.want)
			}
		})
	}
}

// TestRevisedSorted: conflicts come by B's op, then A, then A's op.
func TestRevisedSorted(t *testing.T) {
	a := chain(t, csA, step{1, intention()},
		step{2, ops(revCreate("h", "i", fields("statement", `"a"`, "title", `"t"`)))},
		step{5, ops(revCreate("h", "i", fields("statement", `"a2"`, "title", `"t2"`)))})
	b := chain(t, csB, step{3, intention()},
		step{4, ops(dep(1, csB, csA), revUpdate(11, fields("title", `"x"`)), revUpdate(10, fields("statement", `"x"`)))})
	got := Revised([]ChangeSet{b, a}, csB)
	want := []RevisedConflict{
		{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1), Target: revNode,
			Aspects: []string{"field:statement"}},
		{ChangeSet: csB, Op: opID(11), Dependency: csA, DependencyOp: opID(1), Target: revNode,
			Aspects: []string{"field:title"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Revised = %+v, want %+v", got, want)
	}
	if got := Revised([]ChangeSet{a, b}, csID(9)); got != nil {
		t.Errorf("unknown change set: Revised = %+v, want nil", got)
	}
}

// TestRevisedRecreated: a drop of A's create is no conflict for B while
// something in B's view still creates the node under the same ID, as after
// a carve moves the create to C and A depends on C. A drop whose op ID B's
// op names, as base, still is.
func TestRevisedRecreated(t *testing.T) {
	created := revCreate("h", "i", fields("statement", `"a"`))
	again := created
	again.ID = opID(20)
	unbased := revUpdate(10, fields("statement", `"b"`))
	unbased.Base = opID(99) // not A's create
	carved := func(depOn ...op.Op) []step {
		return []step{{1, intention()}, {2, ops(created)}, {5, Ops{Ops: depOn, Drop: []id.ID{opID(1)}}}}
	}
	bSteps := func(o op.Op) []step { return []step{{3, intention()}, {4, ops(dep(1, csB, csA), o)}} }
	cSteps := []step{{6, intention()}, {7, ops(again)}}
	landed := func(s ChangeSet) ChangeSet { s.State = StateLanded; return s }
	dropped := []RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
		Target: revNode, Dropped: true}}
	cases := []struct {
		name string
		sets func(t *testing.T) []ChangeSet
		want []RevisedConflict
	}{
		{"re-created in the closure", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved(dep(2, csA, csC))...), chain(t, csB, bSteps(unbased)...),
				chain(t, csC, cSteps...)}
		}, nil},
		{"re-created by a landed change set", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved()...), chain(t, csB, bSteps(unbased)...),
				landed(chain(t, csC, cSteps...))}
		}, nil},
		{"re-created by B", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved()...),
				chain(t, csB, step{3, intention()}, step{4, ops(dep(1, csB, csA), unbased)}, step{8, ops(again)})}
		}, nil},
		{"not re-created", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved()...), chain(t, csB, bSteps(unbased)...)}
		}, dropped},
		{"re-created outside the view", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved()...), chain(t, csB, bSteps(unbased)...),
				chain(t, csC, cSteps...)}
		}, dropped},
		{"base names the dropped op", func(t *testing.T) []ChangeSet {
			return []ChangeSet{chain(t, csA, carved(dep(2, csA, csC))...),
				chain(t, csB, bSteps(revUpdate(10, fields("statement", `"b"`)))...), chain(t, csC, cSteps...)}
		}, []RevisedConflict{{ChangeSet: csB, Op: opID(10), Dependency: csA, DependencyOp: opID(1),
			Target: opID(1), Dropped: true}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Revised(c.sets(t), csB); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Revised = %+v, want %+v", got, c.want)
			}
		})
	}
}

package model

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// fid is fixture ID n. IDs grow with n, so fid(1) sorts before fid(2).
func fid(n int) id.ID {
	var i id.ID
	i[0] = 1
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

// csID is fixture change set n, distinct from every fid.
func csID(n int) id.ID {
	var i id.ID
	i[0] = 2
	binary.BigEndian.PutUint32(i[12:], uint32(n))
	return i
}

var (
	csA = csID(1)
	csB = csID(2)
	csC = csID(3)
)

// Op IDs are fid(n) for n below 500, nodes are nd(n), links ln(n).
func nd(n int) id.ID { return fid(500 + n) }
func ln(n int) id.ID { return fid(700 + n) }

// ref turns a fixture number into an ID, 0 into the zero ID.
func ref(f func(int) id.ID, n int) id.ID {
	if n == 0 {
		return id.ID{}
	}
	return f(n)
}

func raw(s string) json.RawMessage { return json.RawMessage(s) }

// mkCap creates capability nd(node) under nd(parent) (0: a root).
func mkCap(n, node int, handle string, parent int, key string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeCreate, Node: nd(node), Kind: op.KindCapability,
		Handle: handle, Parent: ref(nd, parent), Key: key}
}

// mkDec creates decision nd(node), a kind with no place in the tree.
func mkDec(n, node int, handle string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeCreate, Node: nd(node), Kind: op.KindDecision, Handle: handle,
		Fields: map[string]json.RawMessage{"statement": raw(`"use JWT"`)}}
}

// set updates field name of nd(node) to value from base fid(base).
func set(n, node, base int, name, value string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeUpdate, Node: nd(node), Base: fid(base),
		Set: map[string]json.RawMessage{name: raw(value)}}
}

// setTwo updates fields statement and a of nd(node) from base fid(base).
func setTwo(n, node, base int) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeUpdate, Node: nd(node), Base: fid(base),
		Set: map[string]json.RawMessage{"statement": raw(fmt.Sprintf(`"%d"`, n)), "a": raw(fmt.Sprint(n))}}
}

func unset(n, node, base int, names ...string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeUpdate, Node: nd(node), Base: fid(base), Unset: names}
}

func rename(n, node, base int, handle string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeUpdate, Node: nd(node), Base: fid(base), Handle: handle}
}

func move(n, node, base, parent int, key string) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeMove, Node: nd(node), Base: fid(base), Parent: ref(nd, parent), Key: key}
}

func retire(n, node, base int) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeNodeRetire, Node: nd(node), Base: fid(base), Reason: "gone"}
}

func mkLink(n, link int, from, to id.ID) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeLinkCreate, Link: ln(link), LinkType: op.LinkServes, From: from, To: to}
}

func retireLink(n, link, base int) op.Op {
	return op.Op{ID: fid(n), Type: op.TypeLinkRetire, Link: ln(link), Base: fid(base)}
}

// step is one op applied by change set cs.
type step struct {
	cs id.ID
	o  op.Op
}

func a(o op.Op) step { return step{csA, o} }
func b(o op.Op) step { return step{csB, o} }
func c(o op.Op) step { return step{csC, o} }

// want is one expected problem: op fid(n) of change set cs with an error
// matching err.
type want struct {
	cs  id.ID
	n   int
	err error
}

func checkProblems(t *testing.T, got []Problem, ws []want) {
	t.Helper()
	if len(ws) == 0 {
		if got != nil {
			t.Fatalf("problems = %v, want nil", got)
		}
		return
	}
	if len(got) != len(ws) {
		t.Fatalf("got %d problems %v, want %d", len(got), got, len(ws))
	}
	for i, w := range ws {
		p := got[i]
		if p.ChangeSet != w.cs || p.Op != ref(fid, w.n) || !errors.Is(p.Err, w.err) {
			t.Errorf("problem %d = %v, want %s/op %d %v", i, p, w.cs, w.n, w.err)
		}
		if n := sentinels(p.Err); n != 1 {
			t.Errorf("problem %d wraps %d sentinels, want 1: %v", i, n, p.Err)
		}
	}
}

func sentinels(err error) int {
	n := 0
	for _, s := range []error{ErrConflict, ErrMissing, ErrBase, ErrExists, ErrRetired, ErrParent, ErrCycle, ErrHandle, ErrState} {
		if errors.Is(err, s) {
			n++
		}
	}
	return n
}

// build applies steps to a new builder that knows change sets csA, csB and
// csC, and finishes it.
func build(steps []step) (Model, []Problem) {
	bld := newBuilder([]id.ID{csA, csB, csC})
	for _, s := range steps {
		bld.apply(s.cs, s.o)
	}
	return bld.finish()
}

func mustNode(t *testing.T, m Model, n int) Node {
	t.Helper()
	x, ok := m.Node(nd(n))
	if !ok {
		t.Fatalf("node %d missing", n)
	}
	return x
}

func mustLink(t *testing.T, m Model, n int) Link {
	t.Helper()
	x, ok := m.Link(ln(n))
	if !ok {
		t.Fatalf("link %d missing", n)
	}
	return x
}

func field(t *testing.T, m Model, n int, name string) string {
	t.Helper()
	return string(mustNode(t, m, n).Fields[name])
}

func history(rs ...Ref) []Ref { return rs }

func r(cs id.ID, n int) Ref { return Ref{ChangeSet: cs, Op: fid(n)} }

func childIDs(ns []Node) []id.ID {
	var out []id.ID
	for _, x := range ns {
		out = append(out, x.ID)
	}
	return out
}

func nds(ns ...int) []id.ID {
	var out []id.ID
	for _, n := range ns {
		out = append(out, nd(n))
	}
	return out
}

type applyCase struct {
	name     string
	steps    []step
	problems []want
	check    func(t *testing.T, m Model)
}

var applyCases = []applyCase{
	// Rule 1: creates.
	{
		name:  "create node",
		steps: []step{a(mkCap(1, 1, "login", 0, "m")), b(mkCap(2, 2, "token", 1, "m"))},
		check: func(t *testing.T, m Model) {
			x := mustNode(t, m, 2)
			want := Node{ID: nd(2), Kind: op.KindCapability, Handle: "token", Parent: nd(1), Key: "m",
				Version: fid(2), Created: csB, History: history(r(csB, 2))}
			if !reflect.DeepEqual(x, want) {
				t.Errorf("node = %+v, want %+v", x, want)
			}
		},
	},
	{
		name:     "create of an existing node ID is skipped",
		steps:    []step{a(mkCap(1, 1, "login", 0, "m")), b(mkCap(2, 1, "other", 0, "n"))},
		problems: []want{{csB, 2, ErrExists}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); x.Handle != "login" || x.Created != csA {
				t.Errorf("node = %+v, want the first create", x)
			}
		},
	},
	{
		name: "create of a link ID used by a node is skipped",
		steps: []step{a(mkCap(1, 1, "login", 0, "m")), a(mkDec(2, 2, "jwt")),
			a(op.Op{ID: fid(3), Type: op.TypeLinkCreate, Link: nd(1), LinkType: op.LinkServes, From: nd(2), To: nd(1)})},
		problems: []want{{csA, 3, ErrExists}},
	},
	{
		name:     "create under a missing parent is skipped",
		steps:    []step{a(mkCap(1, 1, "login", 9, "m"))},
		problems: []want{{csA, 1, ErrMissing}},
		check: func(t *testing.T, m Model) {
			if _, ok := m.Node(nd(1)); ok {
				t.Error("skipped create is in the model")
			}
		},
	},
	{
		name:     "create under a retired parent is skipped",
		steps:    []step{a(mkCap(1, 1, "login", 0, "m")), a(retire(2, 1, 1)), b(mkCap(3, 2, "token", 1, "m"))},
		problems: []want{{csB, 3, ErrRetired}},
	},
	{
		name:     "create under a parent whose kind has no children is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(mkCap(2, 2, "token", 1, "m"))},
		problems: []want{{csA, 2, ErrParent}},
	},
	{
		name:  "link between nodes and to a change set",
		steps: []step{a(mkCap(1, 1, "login", 0, "m")), a(mkDec(2, 2, "jwt")), b(mkLink(3, 1, nd(2), nd(1))), b(mkLink(4, 2, nd(2), csC))},
		check: func(t *testing.T, m Model) {
			want := Link{ID: ln(1), Type: op.LinkServes, From: nd(2), To: nd(1), Version: fid(3), Created: csB,
				History: history(r(csB, 3))}
			if x := mustLink(t, m, 1); !reflect.DeepEqual(x, want) {
				t.Errorf("link = %+v, want %+v", x, want)
			}
			if x := mustLink(t, m, 2); x.To != csC {
				t.Errorf("link 2 = %+v", x)
			}
		},
	},
	{
		name:     "link to a missing endpoint is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(mkLink(2, 1, nd(1), nd(9))), a(mkLink(3, 2, csID(9), nd(1)))},
		problems: []want{{csA, 2, ErrMissing}, {csA, 3, ErrMissing}},
		check: func(t *testing.T, m Model) {
			if m.Links() != nil {
				t.Errorf("Links = %v, want nil", m.Links())
			}
		},
	},
	{
		name: "link with a missing decision is skipped",
		steps: []step{a(mkDec(1, 1, "jwt")), a(mkCap(2, 2, "login", 0, "m")),
			a(op.Op{ID: fid(3), Type: op.TypeLinkCreate, Link: ln(1), LinkType: op.LinkExempts, From: nd(1), To: nd(2),
				Note: "n", Decision: nd(9)}),
			a(op.Op{ID: fid(4), Type: op.TypeLinkCreate, Link: ln(2), LinkType: op.LinkExempts, From: nd(1), To: nd(2),
				Note: "n", Decision: nd(1)})},
		problems: []want{{csA, 3, ErrMissing}},
		check: func(t *testing.T, m Model) {
			if x := mustLink(t, m, 2); x.Decision != nd(1) || x.Note != "n" {
				t.Errorf("link = %+v", x)
			}
		},
	},
	{
		name:  "link to a retired endpoint applies",
		steps: []step{a(mkDec(1, 1, "jwt")), a(mkDec(2, 2, "pki")), a(retire(3, 2, 2)), a(mkLink(4, 1, nd(1), nd(2)))},
	},

	// Rule 2: target and base.
	{
		name: "op on a missing target is skipped",
		steps: []step{a(set(1, 9, 1, "a", "1")), a(retireLink(2, 9, 1)),
			a(mkDec(3, 1, "jwt")), a(retireLink(4, 1, 3)), a(mkLink(5, 1, nd(1), csB)),
			a(op.Op{ID: fid(6), Type: op.TypeNodeRetire, Node: ln(1), Base: fid(5), Reason: "a link"})},
		problems: []want{{csA, 1, ErrMissing}, {csA, 2, ErrMissing}, {csA, 4, ErrMissing}, {csA, 6, ErrMissing}},
	},
	{
		name:     "unknown op type is skipped",
		steps:    []step{a(op.Op{ID: fid(1), Type: "node.explode"})},
		problems: []want{{csA, 1, ErrMissing}},
	},
	{
		name:     "base not in the target's history is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(mkDec(2, 2, "pki")), b(set(3, 1, 2, "a", "1"))},
		problems: []want{{csB, 3, ErrBase}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); x.Version != fid(1) {
				t.Errorf("Version = %s, want the create", x.Version)
			}
		},
	},
	{
		name:     "op on a retired node from its current version is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(retire(2, 1, 1)), b(set(3, 1, 2, "a", "1")), b(retire(4, 1, 2))},
		problems: []want{{csB, 3, ErrRetired}, {csB, 4, ErrRetired}},
	},
	{
		name:     "link.retire of a retired link from its current version is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(mkLink(2, 1, nd(1), csA)), a(retireLink(3, 1, 2)), b(retireLink(4, 1, 3))},
		problems: []want{{csB, 4, ErrRetired}},
	},
	{
		name:  "a base from the node's own earlier op in the same change set is current",
		steps: []step{a(mkDec(1, 1, "jwt")), a(set(2, 1, 1, "a", "1")), a(set(3, 1, 2, "a", "2")), a(rename(4, 1, 3, "pki"))},
		check: func(t *testing.T, m Model) {
			x := mustNode(t, m, 1)
			if x.Version != fid(4) || x.Handle != "pki" || field(t, m, 1, "a") != "2" {
				t.Errorf("node = %+v", x)
			}
			if want := history(r(csA, 1), r(csA, 2), r(csA, 3), r(csA, 4)); !reflect.DeepEqual(x.History, want) {
				t.Errorf("History = %v, want %v", x.History, want)
			}
		},
	},

	// Rules 2 and 3: stale bases, aspects, conflicts.
	{
		name:  "stale update of a different field merges",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "a", "1")), c(set(3, 1, 1, "b", "2"))},
		check: func(t *testing.T, m Model) {
			x := mustNode(t, m, 1)
			if field(t, m, 1, "a") != "1" || field(t, m, 1, "b") != "2" || x.Version != fid(3) {
				t.Errorf("node = %+v", x)
			}
		},
	},
	{
		name:     "stale update of the same field conflicts and the later wins",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "statement", `"b"`)), c(set(3, 1, 1, "statement", `"c"`))},
		problems: []want{{csC, 3, ErrConflict}},
		check: func(t *testing.T, m Model) {
			if got := field(t, m, 1, "statement"); got != `"c"` {
				t.Errorf("statement = %s, want the later write", got)
			}
			if x := mustNode(t, m, 1); !reflect.DeepEqual(x.History, history(r(csA, 1), r(csB, 2), r(csC, 3))) {
				t.Errorf("History = %v", x.History)
			}
		},
	},
	{
		name:     "stale unset of a set field conflicts",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "a", "1")), c(unset(3, 1, 1, "a"))},
		problems: []want{{csC, 3, ErrConflict}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); x.Fields["a"] != nil {
				t.Errorf("Fields = %v, want a unset", x.Fields)
			}
		},
	},
	{
		name:     "stale handle change conflicts with an earlier one",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(rename(2, 1, 1, "pki")), c(rename(3, 1, 1, "oidc"))},
		problems: []want{{csC, 3, ErrConflict}},
	},
	{
		name:     "stale update after a retire conflicts and still applies",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(retire(2, 1, 1)), c(set(3, 1, 1, "a", "1"))},
		problems: []want{{csC, 3, ErrConflict}},
		check: func(t *testing.T, m Model) {
			x := mustNode(t, m, 1)
			if !x.Retired || x.Reason != "gone" || field(t, m, 1, "a") != "1" || x.Version != fid(3) {
				t.Errorf("node = %+v", x)
			}
		},
	},
	{
		name:     "stale retire after an update conflicts",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "a", "1")), c(retire(3, 1, 1))},
		problems: []want{{csC, 3, ErrConflict}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); !x.Retired {
				t.Error("not retired")
			}
		},
	},
	{
		name: "stale move conflicts with an earlier move",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 0, "n")),
			b(move(3, 2, 2, 1, "a")), c(move(4, 2, 2, 0, "z"))},
		problems: []want{{csC, 4, ErrConflict}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 2); !x.Parent.IsZero() || x.Key != "z" {
				t.Errorf("node = %+v, want the later move", x)
			}
		},
	},
	{
		name: "stale update merges with a move",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 0, "n")),
			b(move(3, 2, 2, 1, "a")), c(set(4, 2, 2, "a", "1")), c(move(5, 1, 1, 0, "b"))},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 2); x.Parent != nd(1) || field(t, m, 2, "a") != "1" {
				t.Errorf("node = %+v", x)
			}
		},
	},
	{
		name: "link.retire twice from the same base conflicts",
		steps: []step{a(mkDec(1, 1, "jwt")), a(mkLink(2, 1, nd(1), csA)),
			b(retireLink(3, 1, 2)), c(retireLink(4, 1, 2))},
		problems: []want{{csC, 4, ErrConflict}},
		check: func(t *testing.T, m Model) {
			x := mustLink(t, m, 1)
			if !x.Retired || x.Version != fid(4) || !reflect.DeepEqual(x.History, history(r(csA, 2), r(csB, 3), r(csC, 4))) {
				t.Errorf("link = %+v", x)
			}
		},
	},
	{
		name: "a stale move that breaks a rule is skipped, not a conflict",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 0, "n")),
			b(move(3, 2, 2, 1, "a")), c(move(4, 2, 2, 9, "z"))},
		problems: []want{{csC, 4, ErrMissing}},
	},

	// Repair: later writes from the current version clear a conflict.
	{
		name: "rewriting the conflicted field from the current version repairs the conflict",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "statement", `"b"`)), c(set(3, 1, 1, "statement", `"c"`)),
			c(set(4, 1, 3, "statement", `"d"`))},
		check: func(t *testing.T, m Model) {
			if got := field(t, m, 1, "statement"); got != `"d"` {
				t.Errorf("statement = %s, want the repair", got)
			}
		},
	},
	{
		name: "an unrelated edit from the current version does not repair a conflict",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "statement", `"b"`)), c(set(3, 1, 1, "statement", `"c"`)),
			c(set(4, 1, 3, "other", "1"))},
		problems: []want{{csC, 3, ErrConflict}},
	},
	{
		name: "a conflict on two fields stays while only one is rewritten",
		steps: []step{a(mkDec(1, 1, "jwt")), b(setTwo(2, 1, 1)), c(setTwo(3, 1, 1)),
			c(set(4, 1, 3, "statement", `"d"`))},
		problems: []want{{csC, 3, ErrConflict}},
	},
	{
		name: "a conflict on two fields is repaired once both are rewritten, by two ops",
		steps: []step{a(mkDec(1, 1, "jwt")), b(setTwo(2, 1, 1)), c(setTwo(3, 1, 1)),
			c(set(4, 1, 3, "statement", `"d"`)), a(set(5, 1, 4, "a", "2"))},
	},
	{
		name: "a stale rewrite does not repair a conflict: it conflicts itself",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "statement", `"b"`)), c(set(3, 1, 1, "statement", `"c"`)),
			c(set(4, 1, 2, "statement", `"d"`))},
		problems: []want{{csC, 3, ErrConflict}, {csC, 4, ErrConflict}},
	},
	{
		name: "a handle conflict is repaired by a rename from the current version",
		steps: []step{a(mkDec(1, 1, "jwt")), b(rename(2, 1, 1, "pki")), c(rename(3, 1, 1, "oidc")),
			c(rename(4, 1, 3, "sso"))},
	},
	{
		name: "a stale retire's conflict stays: nothing can write the retired node",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "a", "1")), c(retire(3, 1, 1)),
			c(set(4, 1, 3, "a", "2"))},
		problems: []want{{csC, 3, ErrConflict}, {csC, 4, ErrRetired}},
	},
	{
		name: "a conflict whose target is then retired stays",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "statement", `"b"`)), c(set(3, 1, 1, "statement", `"c"`)),
			c(retire(4, 1, 3)), c(set(5, 1, 4, "statement", `"d"`))},
		problems: []want{{csC, 3, ErrConflict}, {csC, 5, ErrRetired}},
	},

	// Rule 4: effects of moves and retires.
	{
		name: "move under a missing, retired or childless parent is skipped",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 0, "n")), a(retire(3, 1, 1)),
			a(mkDec(4, 3, "jwt")),
			b(move(5, 2, 2, 9, "a")), b(move(6, 2, 2, 1, "a")), b(move(7, 2, 2, 3, "a"))},
		problems: []want{{csB, 5, ErrMissing}, {csB, 6, ErrRetired}, {csB, 7, ErrParent}},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 2); x.Version != fid(2) || !x.Parent.IsZero() {
				t.Errorf("node = %+v", x)
			}
		},
	},
	{
		name:     "move of a node whose kind has no place in the tree is skipped",
		steps:    []step{a(mkDec(1, 1, "jwt")), a(move(2, 1, 1, 0, "a"))},
		problems: []want{{csA, 2, ErrParent}},
	},
	{
		name: "move under itself or a descendant is skipped",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "mid", 1, "m")), a(mkCap(3, 3, "leaf", 2, "m")),
			b(move(4, 1, 1, 1, "a")), b(move(5, 1, 1, 3, "a")), b(move(6, 3, 3, 0, "a"))},
		problems: []want{{csB, 4, ErrCycle}, {csB, 5, ErrCycle}},
		check: func(t *testing.T, m Model) {
			if got := childIDs(m.Children(id.ID{})); !reflect.DeepEqual(got, nds(3, 1)) {
				t.Errorf("roots = %v, want %v", got, nds(3, 1))
			}
		},
	},
	{
		name:  "retire of a node with live children applies",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 1, "m")), b(retire(3, 1, 1))},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); !x.Retired || x.Reason != "gone" || x.Version != fid(3) {
				t.Errorf("node = %+v", x)
			}
			if got := childIDs(m.Children(nd(1))); !reflect.DeepEqual(got, nds(2)) {
				t.Errorf("children = %v", got)
			}
		},
	},

	// Rule 5: duplicate handles.
	{
		name: "duplicate handle under one parent is reported on the later node",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), b(mkCap(2, 3, "login", 1, "b")), a(mkCap(3, 2, "login", 1, "a")),
			a(mkCap(4, 4, "login", 0, "c"))},
		problems: []want{{csA, 3, ErrHandle}},
		check: func(t *testing.T, m Model) {
			if _, ok := m.Node(nd(2)); !ok {
				t.Error("duplicate handle was skipped")
			}
		},
	},
	{
		name:     "duplicate handle among nodes outside the tree is per kind",
		steps:    []step{a(mkDec(1, 1, "jwt")), b(mkDec(2, 2, "jwt")), b(op.Op{ID: fid(3), Type: op.TypeNodeCreate, Node: nd(3), Kind: op.KindActor, Handle: "jwt"})},
		problems: []want{{csB, 2, ErrHandle}},
	},
	{
		name: "a rename or move into a taken handle is reported on that op",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "login", 1, "a")), a(mkCap(3, 3, "token", 1, "b")),
			a(mkCap(4, 4, "login", 0, "n")), b(rename(5, 3, 3, "login")), c(move(6, 4, 4, 1, "c"))},
		problems: []want{{csB, 5, ErrHandle}, {csC, 6, ErrHandle}},
	},
	{
		name:  "retired nodes do not hold their handle",
		steps: []step{a(mkDec(1, 1, "jwt")), a(retire(2, 1, 1)), b(mkDec(3, 2, "jwt"))},
	},

	// The model's accessors.
	{
		name: "children are sorted by key then ID, retired ones included",
		steps: []step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 4, "d", 1, "b")), a(mkCap(3, 3, "c", 1, "a")),
			a(mkCap(4, 2, "b", 1, "b")), a(retire(5, 3, 3)), a(op.Op{ID: fid(6), Type: op.TypeNodeCreate, Node: nd(5),
				Kind: op.KindIntent, Handle: "aim", Key: "a"}), a(mkDec(7, 6, "jwt"))},
		check: func(t *testing.T, m Model) {
			if got := childIDs(m.Children(nd(1))); !reflect.DeepEqual(got, nds(3, 2, 4)) {
				t.Errorf("children = %v, want %v", got, nds(3, 2, 4))
			}
			if got := childIDs(m.Children(id.ID{})); !reflect.DeepEqual(got, nds(5, 1)) {
				t.Errorf("roots = %v, want %v", got, nds(5, 1))
			}
			if got := m.Children(nd(2)); got != nil {
				t.Errorf("children of a leaf = %v, want nil", got)
			}
			if got := childIDs(m.Nodes()); !reflect.DeepEqual(got, nds(1, 2, 3, 4, 5, 6)) {
				t.Errorf("Nodes = %v", got)
			}
		},
	},
	{
		name: "a moved node leaves its old parent",
		steps: []step{a(mkCap(1, 1, "a", 0, "m")), a(mkCap(2, 2, "b", 0, "n")), a(mkCap(3, 3, "c", 1, "m")),
			a(move(4, 3, 3, 2, "z"))},
		check: func(t *testing.T, m Model) {
			if got := m.Children(nd(1)); got != nil {
				t.Errorf("old parent's children = %v", got)
			}
			if got := childIDs(m.Children(nd(2))); !reflect.DeepEqual(got, nds(3)) {
				t.Errorf("new parent's children = %v", got)
			}
		},
	},
	{
		name:  "update sets and unsets fields and replaces the handle",
		steps: []step{a(mkDec(1, 1, "jwt")), b(set(2, 1, 1, "a", "1")), b(unset(3, 1, 2, "statement")), b(rename(4, 1, 3, "pki"))},
		check: func(t *testing.T, m Model) {
			x := mustNode(t, m, 1)
			if x.Handle != "pki" || !reflect.DeepEqual(x.Fields, map[string]json.RawMessage{"a": raw("1")}) {
				t.Errorf("node = %+v", x)
			}
		},
	},
	{
		name:  "unsetting every field leaves nil fields",
		steps: []step{a(mkDec(1, 1, "jwt")), a(unset(2, 1, 1, "statement"))},
		check: func(t *testing.T, m Model) {
			if x := mustNode(t, m, 1); x.Fields != nil {
				t.Errorf("Fields = %#v, want nil", x.Fields)
			}
		},
	},
}

func TestApplyCases(t *testing.T) {
	for _, tc := range applyCases {
		t.Run(tc.name, func(t *testing.T) {
			m, ps := build(tc.steps)
			checkProblems(t, ps, tc.problems)
			if tc.check != nil {
				tc.check(t, m)
			}
		})
	}
}

// TestApplySentinels checks that every sentinel apply can report has a row.
// ErrState belongs to the views.
func TestApplySentinels(t *testing.T) {
	for _, s := range []error{ErrConflict, ErrMissing, ErrBase, ErrExists, ErrRetired, ErrParent, ErrCycle, ErrHandle} {
		found := false
		for _, tc := range applyCases {
			for _, w := range tc.problems {
				found = found || w.err == s
			}
		}
		if !found {
			t.Errorf("no case reports %v", s)
		}
	}
}

func TestApplyAccessorCopiesAreIndependent(t *testing.T) {
	m, _ := build([]step{a(mkCap(1, 1, "root", 0, "m")), a(mkCap(2, 2, "leaf", 1, "m")), b(set(3, 2, 2, "a", `"x"`)),
		a(mkDec(4, 3, "jwt")), a(mkLink(5, 1, nd(3), nd(1)))})
	before := []any{mustNode(t, m, 2), m.Nodes(), m.Links(), m.Children(nd(1))}

	x := mustNode(t, m, 2)
	x.Fields["a"][1] = 'y'
	x.Fields["b"] = raw("2")
	x.History[0] = Ref{}
	ns := m.Nodes()
	ns[0].History[0] = Ref{}
	ns[1].Fields["c"] = raw("3")
	ls := m.Links()
	ls[0].History[0] = Ref{}
	cs := m.Children(nd(1))
	cs[0].Fields["a"][1] = 'z'

	after := []any{mustNode(t, m, 2), m.Nodes(), m.Links(), m.Children(nd(1))}
	if !reflect.DeepEqual(before, after) {
		t.Errorf("model changed through its accessors:\nbefore %+v\nafter  %+v", before, after)
	}
}

func TestApplyInputsAreNotShared(t *testing.T) {
	o := mkDec(1, 1, "jwt")
	u := set(2, 1, 1, "a", `"x"`)
	m, _ := build([]step{a(o), a(u)})
	o.Fields["statement"][1] = 'X'
	u.Set["a"][1] = 'Y'
	if got := field(t, m, 1, "statement"); got != `"use JWT"` {
		t.Errorf("statement = %s, changed through the create op", got)
	}
	if got := field(t, m, 1, "a"); got != `"x"` {
		t.Errorf("a = %s, changed through the update op", got)
	}
}

func TestApplyFinishIsASnapshot(t *testing.T) {
	bld := newBuilder(nil)
	bld.apply(csA, mkDec(1, 1, "jwt"))
	m1, ps1 := bld.finish()
	cl := bld.clone()
	bld.apply(csA, set(2, 1, 1, "a", "1"))
	bld.apply(csB, mkDec(3, 2, "jwt"))
	m2, ps2 := bld.finish()

	if x := mustNode(t, m1, 1); x.Version != fid(1) || x.Fields["a"] != nil {
		t.Errorf("first snapshot changed: %+v", x)
	}
	if ps1 != nil {
		t.Errorf("first problems = %v", ps1)
	}
	checkProblems(t, ps2, []want{{csB, 3, ErrHandle}})
	if x := mustNode(t, m2, 1); x.Version != fid(2) {
		t.Errorf("second snapshot = %+v", x)
	}
	m3, ps3 := cl.finish()
	if !reflect.DeepEqual(m3, m1) || ps3 != nil {
		t.Errorf("clone changed with the original")
	}
}

func TestApplyCloneKeepsPendingConflicts(t *testing.T) {
	bld := newBuilder(nil)
	bld.apply(csA, mkDec(1, 1, "jwt"))
	bld.apply(csB, set(2, 1, 1, "statement", `"b"`))
	bld.apply(csC, set(3, 1, 1, "statement", `"c"`))
	cl := bld.clone()
	bld.apply(csC, set(4, 1, 3, "statement", `"d"`))

	_, ps := bld.finish()
	checkProblems(t, ps, nil)
	_, ps = cl.finish()
	checkProblems(t, ps, []want{{csC, 3, ErrConflict}})
	cl.apply(csC, set(5, 1, 3, "statement", `"e"`))
	_, ps = cl.finish()
	checkProblems(t, ps, nil)
}

func TestApplyZeroModel(t *testing.T) {
	var m Model
	if _, ok := m.Node(nd(1)); ok {
		t.Error("Node found in the zero model")
	}
	if _, ok := m.Link(ln(1)); ok {
		t.Error("Link found in the zero model")
	}
	if m.Nodes() != nil || m.Links() != nil || m.Children(id.ID{}) != nil {
		t.Error("zero model is not empty")
	}
	if _, ps := newBuilder(nil).finish(); ps != nil {
		t.Errorf("problems = %v", ps)
	}
}

func TestApplyProblemsAreSorted(t *testing.T) {
	_, ps := build([]step{c(set(9, 1, 1, "a", "1")), a(set(5, 1, 1, "a", "1")), a(set(2, 1, 1, "a", "1")), b(set(1, 1, 1, "a", "1"))})
	checkProblems(t, ps, []want{{csA, 2, ErrMissing}, {csA, 5, ErrMissing}, {csB, 1, ErrMissing}, {csC, 9, ErrMissing}})
	if got := ps[0].String(); got != csA.String()+"/"+fid(2).String()+": "+ps[0].Err.Error() {
		t.Errorf("String = %q", got)
	}
}

func TestApplyAspects(t *testing.T) {
	for _, tc := range []struct {
		o    op.Op
		want []string
	}{
		{mkDec(1, 1, "jwt"), nil},
		{op.Op{Type: op.TypeNodeUpdate, Handle: "x", Set: map[string]json.RawMessage{"b": raw("1"), "a": raw("1")}, Unset: []string{"c"}},
			[]string{"field:a", "field:b", "field:c", "handle"}},
		{move(1, 1, 1, 0, "a"), []string{"position"}},
		{retire(1, 1, 1), []string{"retired"}},
		{retireLink(1, 1, 1), []string{"retired"}},
	} {
		if got := aspects(tc.o); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("aspects(%s) = %v, want %v", tc.o.Type, got, tc.want)
		}
	}
	for _, tc := range []struct {
		a, b []string
		want bool
	}{
		{[]string{"field:a"}, []string{"field:a", "handle"}, true},
		{[]string{"field:a"}, []string{"field:b"}, false},
		{[]string{"position"}, []string{"handle"}, false},
		{[]string{"retired"}, []string{"field:b"}, true},
		{[]string{"position"}, []string{"retired"}, true},
		{nil, []string{"retired"}, false},
		{nil, nil, false},
	} {
		if got := overlaps(tc.a, tc.b); got != tc.want {
			t.Errorf("overlaps(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

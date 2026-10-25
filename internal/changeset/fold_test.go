package changeset

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
	"pgregory.net/rapid"
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
	csA   = csID(1)
	csB   = csID(2)
	actor = fid(900)
	owner = fid(901)
)

// opNode and opBase are fixture nodes and base versions for ops.
func opNode(n int) id.ID { return fid(800 + n) }

var opBase = fid(850)

// retire is a node.retire op with ID fid(700+n) on node opNode(node).
func retire(n, node int, reason string) op.Op {
	return op.Op{ID: fid(700 + n), Type: op.TypeNodeRetire, Node: opNode(node), Base: opBase, Reason: reason}
}

// move is a node.move op with ID fid(700+n) on node opNode(node).
func move(n, node int) op.Op {
	return op.Op{ID: fid(700 + n), Type: op.TypeNodeMove, Node: opNode(node), Base: opBase, Key: "i"}
}

func opID(n int) id.ID { return fid(700 + n) }

// mk builds entry n of change set cs with body b after the given entries.
func mk(t testing.TB, cs id.ID, n int, b Body, after ...int) entry.Entry {
	t.Helper()
	raw, err := EncodeBody(b)
	if err != nil {
		t.Fatalf("entry %d: %v", n, err)
	}
	var as []id.ID
	for _, a := range after {
		as = append(as, fid(a))
	}
	slices.SortFunc(as, id.Compare)
	as = slices.Compact(as)
	if as == nil {
		as = []id.ID{}
	}
	return entry.Entry{
		V: entry.FormatVersion, ID: fid(n), ChangeSet: cs, Actor: actor, For: owner,
		At: time.Unix(1_800_000_000+int64(n), 0).UTC(), After: as, Type: b.Type(), Body: raw,
	}
}

// e builds entry n of csA.
func e(t testing.TB, n int, b Body, after ...int) entry.Entry {
	t.Helper()
	return mk(t, csA, n, b, after...)
}

func intention() Open { return Open{Kind: KindIntention, Intent: "Refresh tokens"} }
func record() Open    { return Open{Kind: KindRecord, Intent: "Refreshed tokens"} }

var (
	acceptBody = Accept{Rule: "owner-approves", Reason: "ok"}
	rejectBody = Reject{Rule: "owner-approves"}
)

// want is one expected problem: entry fid(n) (or the zero ID when n is 0)
// with an error matching err.
type want struct {
	n   int
	err error
}

func checkProblems(t *testing.T, got []Problem, cs id.ID, ws []want) {
	t.Helper()
	if len(ws) == 0 {
		if got != nil {
			t.Fatalf("problems = %v, want nil", got)
		}
		return
	}
	if len(got) != len(ws) {
		t.Fatalf("got %d problems %v, want %d %v", len(got), got, len(ws), ws)
	}
	for i, w := range ws {
		var eid id.ID
		if w.n != 0 {
			eid = fid(w.n)
		}
		p := got[i]
		if p.ChangeSet != cs || p.Entry != eid || !errors.Is(p.Err, w.err) {
			t.Errorf("problem %d = {%s %s %v}, want entry %d %v", i, p.ChangeSet, p.Entry, p.Err, w.n, w.err)
		}
	}
}

func opIDs(ops []op.Op) []id.ID {
	var out []id.ID
	for _, o := range ops {
		out = append(out, o.ID)
	}
	return out
}

func opIDsOf(ns ...int) []id.ID {
	var out []id.ID
	for _, n := range ns {
		out = append(out, opID(n))
	}
	return out
}

func TestStateTerminal(t *testing.T) {
	for s, want := range map[State]bool{
		StateDraft: false, StateProposed: false, StateAccepted: false,
		StateLanded: true, StateRejected: true, StateWithdrawn: true, StateSuperseded: true,
	} {
		if got := s.Terminal(); got != want {
			t.Errorf("%s.Terminal() = %v, want %v", s, got, want)
		}
	}
}

// TestFoldLifecycle runs every entry type from every state, for both kinds,
// against the lifecycle table written out independently here.
func TestFoldLifecycle(t *testing.T) {
	// prefixes reach each state from a fresh open, as a chain of bodies.
	prefixes := map[State][]Body{
		StateDraft:      nil,
		StateProposed:   {Propose{}},
		StateAccepted:   {Propose{}, acceptBody},
		StateLanded:     {Propose{}, acceptBody, Land{}},
		StateRejected:   {Propose{}, rejectBody},
		StateWithdrawn:  {Withdraw{}},
		StateSuperseded: {Supersede{By: csB}},
	}
	inFlight := []State{StateDraft, StateProposed, StateAccepted}
	type move struct {
		name string
		body Body
		kind Kind
		from []State
		to   State // "" keeps the state
	}
	moves := []move{
		{"ops", Ops{Ops: []op.Op{retire(9, 1, "x")}}, KindIntention, inFlight, ""},
		{"propose", Propose{}, KindIntention, []State{StateDraft}, StateProposed},
		{"accept", acceptBody, KindIntention, []State{StateProposed}, StateAccepted},
		{"reject", rejectBody, KindIntention, []State{StateProposed}, StateRejected},
		{"withdraw", Withdraw{}, KindIntention, inFlight, StateWithdrawn},
		{"supersede", Supersede{By: csB}, KindIntention, inFlight, StateSuperseded},
		{"land intention", Land{}, KindIntention, []State{StateAccepted}, StateLanded},
		{"land record", Land{}, KindRecord, inFlight, StateLanded},
		{"note", Note{Text: "hi"}, KindIntention, States(), ""},
	}
	for _, m := range moves {
		for _, from := range States() {
			t.Run(fmt.Sprintf("%s from %s", m.name, from), func(t *testing.T) {
				open := intention()
				if m.kind == KindRecord {
					open = record()
				}
				es := []entry.Entry{e(t, 1, open)}
				for i, b := range prefixes[from] {
					es = append(es, e(t, i+2, b, i+1))
				}
				last := len(es) + 1
				es = append(es, e(t, last, m.body, last-1))
				got, ps := Fold(csA, es)
				allowed := slices.Contains(m.from, from)
				wantState := from
				if allowed && m.to != "" {
					wantState = m.to
				}
				if got.State != wantState {
					t.Errorf("State = %s, want %s", got.State, wantState)
				}
				if allowed {
					checkProblems(t, ps, csA, nil)
				} else {
					checkProblems(t, ps, csA, []want{{last, ErrTransition}})
				}
				if len(got.Applied) != last || got.Applied[last-1] != fid(last) {
					t.Errorf("Applied = %v, want all %d entries", got.Applied, last)
				}
				if allowed && m.name == "ops" && !reflect.DeepEqual(opIDs(got.Ops), opIDsOf(9)) {
					t.Errorf("Ops = %v", got.Ops)
				}
				if !allowed && m.name == "ops" && got.Ops != nil {
					t.Errorf("Ops = %v, want nil after a skipped ops entry", got.Ops)
				}
			})
		}
	}
}

// TestFoldFull checks every field of a folded change set.
func TestFoldFull(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()),
		e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b"), retire(3, 3, "c")}}, 1),
		e(t, 3, Ops{Ops: []op.Op{retire(2, 2, "b2"), retire(4, 4, "d")}, Drop: []id.ID{opID(1)}}, 2),
		e(t, 4, Note{Text: "side note"}, 2),
		e(t, 5, Propose{}, 3, 4),
		e(t, 6, acceptBody, 5),
		e(t, 7, Land{}, 6),
	}
	want := ChangeSet{
		ID:        csA,
		Open:      intention(),
		OpenEntry: fid(1),
		State:     StateLanded,
		Ops:       []op.Op{retire(2, 2, "b2"), retire(3, 3, "c"), retire(4, 4, "d")},
		Decided: &Decision{
			Entry: fid(6), Actor: actor, For: owner, Type: entry.TypeAccept,
			Rule: acceptBody.Rule, Reason: acceptBody.Reason,
		},
		History: map[id.ID][]Revision{
			opID(2): {{Entry: fid(2), Op: retire(2, 2, "b")}, {Entry: fid(3), Op: retire(2, 2, "b2")}},
			opID(3): {{Entry: fid(2), Op: retire(3, 3, "c")}},
			opID(4): {{Entry: fid(3), Op: retire(4, 4, "d")}},
		},
		Drops:         map[id.ID]Revision{opID(1): {Entry: fid(3), Op: retire(1, 1, "a")}},
		TerminalEntry: fid(7),
		Heads:         []id.ID{fid(7)},
		Applied:       []id.ID{fid(1), fid(2), fid(3), fid(4), fid(5), fid(6), fid(7)},
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, nil)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Fold =\n%+v\nwant\n%+v", got, want)
	}
}

// TestFoldHistory: an op written in entry 2 and replaced in entry 4 keeps
// both versions with their entries; an op dropped in entry 5 keeps its
// drop entry and its last version; a skipped revision adds nothing.
func TestFoldHistory(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()),
		e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b")}}, 1),
		e(t, 3, Ops{Ops: []op.Op{retire(2, 2, "b2")}}, 2),
		e(t, 4, Ops{Ops: []op.Op{retire(1, 1, "a2")}}, 3),
		e(t, 5, Ops{Drop: []id.ID{opID(2)}}, 4),
		e(t, 6, Ops{Ops: []op.Op{retire(1, 1, "a3"), retire(2, 2, "b3")}}, 5),
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, []want{{6, ErrRevision}})
	wantH := map[id.ID][]Revision{
		opID(1): {{Entry: fid(2), Op: retire(1, 1, "a")}, {Entry: fid(4), Op: retire(1, 1, "a2")}},
	}
	if !reflect.DeepEqual(got.History, wantH) {
		t.Errorf("History = %+v, want %+v", got.History, wantH)
	}
	wantD := map[id.ID]Revision{opID(2): {Entry: fid(5), Op: retire(2, 2, "b2")}}
	if !reflect.DeepEqual(got.Drops, wantD) {
		t.Errorf("Drops = %+v, want %+v", got.Drops, wantD)
	}

	got, _ = Fold(csA, es[:1])
	if got.History != nil || got.Drops != nil {
		t.Errorf("no ops: History = %v, Drops = %v, want nil", got.History, got.Drops)
	}
}

func TestFoldRejectedAndSuperseded(t *testing.T) {
	got, ps := Fold(csA, []entry.Entry{e(t, 1, intention()), e(t, 2, Propose{}, 1), e(t, 3, rejectBody, 2)})
	checkProblems(t, ps, csA, nil)
	if got.State != StateRejected || got.TerminalEntry != fid(3) || got.Decided == nil ||
		got.Decided.Type != entry.TypeReject || got.Decided.Rule != rejectBody.Rule {
		t.Errorf("rejected fold = %+v", got)
	}
	got, ps = Fold(csA, []entry.Entry{e(t, 1, intention()), e(t, 2, Supersede{By: csB}, 1)})
	checkProblems(t, ps, csA, nil)
	if got.State != StateSuperseded || got.SupersededBy != csB || got.TerminalEntry != fid(2) || got.Decided != nil {
		t.Errorf("superseded fold = %+v", got)
	}
}

type foldCase struct {
	name     string
	entries  func(t testing.TB) []entry.Entry
	state    State
	ops      []id.ID
	applied  []int
	heads    []int
	problems []want
}

func nums(ns []int) []id.ID {
	var out []id.ID
	for _, n := range ns {
		out = append(out, fid(n))
	}
	return out
}

func badBody(t testing.TB, n int, after ...int) entry.Entry {
	x := e(t, n, Note{Text: "x"}, after...)
	x.Body = json.RawMessage(`{}`)
	return x
}

var foldCases = []foldCase{
	{
		name: "second open",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, Propose{}, 1), e(t, 3, record())}
		},
		state: StateProposed, applied: []int{1, 2}, heads: []int{2},
		problems: []want{{3, ErrOrder}},
	},
	{
		name: "second open with a smaller ID wins; entries after the other leave",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 2, intention()), e(t, 1, record()), e(t, 3, Propose{}, 2), e(t, 4, Note{Text: "n"}, 1)}
		},
		state: StateDraft, applied: []int{1, 4}, heads: []int{4},
		problems: []want{{2, ErrOrder}, {3, ErrOrder}},
	},
	{
		name: "open not first",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, record(), 1), e(t, 3, Propose{}, 2)}
		},
		state: StateDraft, applied: []int{1}, heads: []int{1},
		problems: []want{{2, ErrOrder}, {3, ErrOrder}},
	},
	{
		name: "empty after",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, Propose{})}
		},
		state: StateDraft, applied: []int{1}, heads: []int{1},
		problems: []want{{2, ErrOrder}},
	},
	{
		name: "dangling after",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, Propose{}, 1, 99), e(t, 3, Withdraw{}, 2)}
		},
		state: StateDraft, applied: []int{1}, heads: []int{1},
		problems: []want{{2, ErrOrder}, {3, ErrOrder}},
	},
	{
		name: "a cycle and entries after it",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()), e(t, 2, Note{Text: "a"}, 1, 3), e(t, 3, Note{Text: "b"}, 2),
				e(t, 4, Propose{}, 3), e(t, 5, Note{Text: "self"}, 5), e(t, 6, Note{Text: "fine"}, 1),
			}
		},
		state: StateDraft, applied: []int{1, 6}, heads: []int{6},
		problems: []want{{2, ErrOrder}, {3, ErrOrder}, {4, ErrOrder}, {5, ErrOrder}},
	},
	{
		name: "entry after a bad-body entry leaves the graph",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), badBody(t, 2, 1), e(t, 3, Propose{}, 2), e(t, 4, Note{Text: "n"}, 3)}
		},
		state: StateDraft, applied: []int{1}, heads: []int{1},
		problems: []want{{2, ErrBody}, {3, ErrOrder}, {4, ErrOrder}},
	},
	{
		name: "entry after an illegal transition is still applied",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, acceptBody, 1), e(t, 3, Propose{}, 2)}
		},
		state: StateProposed, applied: []int{1, 2, 3}, heads: []int{3},
		problems: []want{{2, ErrTransition}},
	},
	{
		name: "entry from another change set is ignored",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), mk(t, csB, 2, Propose{}, 1)}
		},
		state: StateDraft, applied: []int{1}, heads: []int{1},
		problems: []want{{2, ErrOrder}},
	},
	{
		name: "no open at all",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 2, Propose{}, 1), e(t, 3, Withdraw{}, 2)}
		},
		problems: []want{{0, ErrOrder}, {2, ErrOrder}, {3, ErrOrder}},
	},
	{
		name:     "no entries",
		entries:  func(t testing.TB) []entry.Entry { return nil },
		problems: []want{{0, ErrOrder}},
	},
	{
		name: "only open has a bad body",
		entries: func(t testing.TB) []entry.Entry {
			x := e(t, 1, intention())
			x.Body = json.RawMessage(`{}`)
			return []entry.Entry{x, e(t, 2, Propose{}, 1)}
		},
		problems: []want{{0, ErrOrder}, {1, ErrBody}, {2, ErrOrder}},
	},
	{
		name: "duplicate entry ID keeps one copy",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, Propose{}, 1), e(t, 2, Propose{}, 1)}
		},
		state: StateProposed, applied: []int{1, 2}, heads: []int{2},
		problems: []want{{2, ErrOrder}},
	},
	{
		name: "supersede by self",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{e(t, 1, intention()), e(t, 2, Supersede{By: csA}, 1)}
		},
		state: StateDraft, applied: []int{1, 2}, heads: []int{2},
		problems: []want{{2, ErrTransition}},
	},
	{
		name: "replace in place",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b")}}, 1),
				e(t, 3, Ops{Ops: []op.Op{retire(3, 3, "c"), retire(1, 1, "a2")}}, 2),
			}
		},
		state: StateDraft, ops: opIDsOf(1, 2, 3), applied: []int{1, 2, 3}, heads: []int{3},
	},
	{
		name: "replace with a changed type",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1),
				e(t, 3, Ops{Ops: []op.Op{retire(2, 2, "new"), move(1, 1)}}, 2),
			}
		},
		state: StateDraft, ops: opIDsOf(1), applied: []int{1, 2, 3}, heads: []int{3},
		problems: []want{{3, ErrRevision}},
	},
	{
		name: "replace with a changed target",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1),
				e(t, 3, Ops{Ops: []op.Op{retire(1, 2, "a")}}, 2),
			}
		},
		state: StateDraft, ops: opIDsOf(1), applied: []int{1, 2, 3}, heads: []int{3},
		problems: []want{{3, ErrRevision}},
	},
	{
		name: "drop unknown",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1),
				e(t, 3, Ops{Ops: []op.Op{retire(2, 2, "b")}, Drop: []id.ID{opID(9)}}, 2),
			}
		},
		state: StateDraft, ops: opIDsOf(1), applied: []int{1, 2, 3}, heads: []int{3},
		problems: []want{{3, ErrRevision}},
	},
	{
		name: "drop twice",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b")}}, 1),
				e(t, 3, Ops{Drop: []id.ID{opID(1)}}, 2),
				e(t, 4, Ops{Drop: []id.ID{opID(1)}}, 3),
			}
		},
		state: StateDraft, ops: opIDsOf(2), applied: []int{1, 2, 3, 4}, heads: []int{4},
		problems: []want{{4, ErrRevision}},
	},
	{
		name: "resurrect dropped",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()),
				e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1),
				e(t, 3, Ops{Drop: []id.ID{opID(1)}}, 2),
				e(t, 4, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 3),
			}
		},
		state: StateDraft, applied: []int{1, 2, 3, 4}, heads: []int{4},
		problems: []want{{4, ErrRevision}},
	},
	{
		name: "ops kept through proposed and accepted",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()), e(t, 2, Propose{}, 1),
				e(t, 3, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 2), e(t, 4, acceptBody, 3),
				e(t, 5, Ops{Ops: []op.Op{retire(2, 2, "b")}}, 4),
			}
		},
		state: StateAccepted, ops: opIDsOf(1, 2), applied: []int{1, 2, 3, 4, 5}, heads: []int{5},
	},
	{
		name: "concurrent note and propose need no report",
		entries: func(t testing.TB) []entry.Entry {
			return []entry.Entry{
				e(t, 1, intention()), e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1),
				e(t, 3, Note{Text: "n"}, 1), e(t, 4, Propose{}, 1),
			}
		},
		state: StateProposed, ops: opIDsOf(1), applied: []int{1, 2, 3, 4}, heads: []int{2, 3, 4},
	},
}

func TestFoldCases(t *testing.T) {
	for _, c := range foldCases {
		t.Run(c.name, func(t *testing.T) {
			got, ps := Fold(csA, c.entries(t))
			checkProblems(t, ps, csA, c.problems)
			if got.ID != csA {
				t.Errorf("ID = %s", got.ID)
			}
			if got.State != c.state {
				t.Errorf("State = %q, want %q", got.State, c.state)
			}
			if !reflect.DeepEqual(opIDs(got.Ops), c.ops) {
				t.Errorf("Ops = %v, want IDs %v", opIDs(got.Ops), c.ops)
			}
			if !reflect.DeepEqual(got.Applied, nums(c.applied)) {
				t.Errorf("Applied = %v, want %v", got.Applied, nums(c.applied))
			}
			if !reflect.DeepEqual(got.Heads, nums(c.heads)) {
				t.Errorf("Heads = %v, want %v", got.Heads, nums(c.heads))
			}
		})
	}
}

func TestFoldReplaceKeepsPositionAndValue(t *testing.T) {
	got, ps := Fold(csA, foldCase_(t, "replace in place").entries(t))
	checkProblems(t, ps, csA, nil)
	want := []op.Op{retire(1, 1, "a2"), retire(2, 2, "b"), retire(3, 3, "c")}
	if !reflect.DeepEqual(got.Ops, want) {
		t.Fatalf("Ops = %+v, want %+v", got.Ops, want)
	}
}

func foldCase_(t *testing.T, name string) foldCase {
	for _, c := range foldCases {
		if c.name == name {
			return c
		}
	}
	t.Fatalf("no fold case %q", name)
	return foldCase{}
}

func TestFoldBodyErrorKeepsCause(t *testing.T) {
	x := e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a")}}, 1)
	x.Body = json.RawMessage(strings.Replace(string(x.Body), `"node.retire"`, `"node.nope"`, 1))
	_, ps := Fold(csA, []entry.Entry{e(t, 1, intention()), x})
	if len(ps) != 1 || !errors.Is(ps[0].Err, ErrBody) || !errors.Is(ps[0].Err, ErrInvalid) || !errors.Is(ps[0].Err, op.ErrInvalid) {
		t.Fatalf("problems = %v", ps)
	}
}

// TestFoldConcurrentLandWithdraw is the plan's merge contradiction: a land
// and a withdraw both after the accept. Fold order is by entry ID, so the
// land (4) applies first; the withdraw (5) is reported as concurrent with
// it and, since landed is terminal, as an illegal move. No panic.
func TestFoldConcurrentLandWithdraw(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()), e(t, 2, Propose{}, 1), e(t, 3, acceptBody, 2),
		e(t, 4, Land{}, 3), e(t, 5, Withdraw{}, 3),
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, []want{{5, ErrConcurrent}, {5, ErrTransition}})
	if n := countIs(ps, ErrConcurrent); n != 1 {
		t.Errorf("%d ErrConcurrent problems, want 1", n)
	}
	if got.State != StateLanded || got.TerminalEntry != fid(4) {
		t.Errorf("State = %s by %s, want landed by entry 4", got.State, got.TerminalEntry)
	}
	if !reflect.DeepEqual(got.Heads, nums([]int{4, 5})) {
		t.Errorf("Heads = %v", got.Heads)
	}
	if !strings.Contains(ps[0].Err.Error(), fid(4).String()) {
		t.Errorf("concurrency problem %q does not name entry 4", ps[0].Err)
	}
}

func countIs(ps []Problem, target error) int {
	n := 0
	for _, p := range ps {
		if errors.Is(p.Err, target) {
			n++
		}
	}
	return n
}

// TestFoldConcurrentAcceptWithdraw: the later entry in fold order wins when
// its move is legal from the earlier one's result.
func TestFoldConcurrentAcceptWithdraw(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()), e(t, 2, Propose{}, 1),
		e(t, 3, acceptBody, 2), e(t, 4, Withdraw{}, 2),
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, []want{{4, ErrConcurrent}})
	if got.State != StateWithdrawn || got.Decided == nil || got.TerminalEntry != fid(4) {
		t.Errorf("fold = %+v", got)
	}

	// Swapped IDs: withdraw folds first, so the accept is illegal too.
	es = []entry.Entry{
		e(t, 1, intention()), e(t, 2, Propose{}, 1),
		e(t, 3, Withdraw{}, 2), e(t, 4, acceptBody, 2),
	}
	got, ps = Fold(csA, es)
	checkProblems(t, ps, csA, []want{{4, ErrConcurrent}, {4, ErrTransition}})
	if got.State != StateWithdrawn || got.Decided != nil {
		t.Errorf("fold = %+v", got)
	}
}

func TestFoldConcurrentThreeWay(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()), e(t, 2, Propose{}, 1), e(t, 3, Withdraw{}, 1), e(t, 4, Supersede{By: csB}, 1),
	}
	_, ps := Fold(csA, es)
	if countIs(ps, ErrConcurrent) != 3 {
		t.Fatalf("problems = %v, want 3 concurrent pairs", ps)
	}
}

func TestFoldConcurrentReplacements(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()),
		e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b")}}, 1),
		e(t, 3, Ops{Ops: []op.Op{retire(1, 1, "from 3")}}, 2),
		e(t, 4, Ops{Ops: []op.Op{retire(1, 1, "from 4")}}, 2),
		e(t, 5, Ops{Ops: []op.Op{retire(2, 2, "b2")}}, 2),
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, []want{{4, ErrConcurrent}})
	want := []op.Op{retire(1, 1, "from 4"), retire(2, 2, "b2")}
	if !reflect.DeepEqual(got.Ops, want) {
		t.Errorf("Ops = %+v, want %+v", got.Ops, want)
	}
	if !strings.Contains(ps[0].Err.Error(), opID(1).String()) {
		t.Errorf("problem %q does not name the shared op", ps[0].Err)
	}
}

func TestFoldConcurrentDropAndReplace(t *testing.T) {
	es := []entry.Entry{
		e(t, 1, intention()),
		e(t, 2, Ops{Ops: []op.Op{retire(1, 1, "a"), retire(2, 2, "b")}}, 1),
		e(t, 3, Ops{Drop: []id.ID{opID(1)}}, 2),
		e(t, 4, Ops{Ops: []op.Op{retire(1, 1, "again")}}, 2),
	}
	got, ps := Fold(csA, es)
	checkProblems(t, ps, csA, []want{{4, ErrConcurrent}, {4, ErrRevision}})
	if !reflect.DeepEqual(opIDs(got.Ops), opIDsOf(2)) {
		t.Errorf("Ops = %v", opIDs(got.Ops))
	}
}

func TestFoldLog(t *testing.T) {
	l := store.Log{Entries: []entry.Entry{
		mk(t, csB, 10, intention()),
		mk(t, csA, 1, intention()),
		mk(t, csA, 2, Supersede{By: csB}, 1),
		mk(t, csB, 11, Withdraw{}),
	}}
	got, ps := FoldLog(l)
	if len(got) != 2 || got[0].ID != csA || got[1].ID != csB {
		t.Fatalf("FoldLog = %+v", got)
	}
	if got[0].State != StateSuperseded || got[1].State != StateDraft {
		t.Errorf("states = %s, %s", got[0].State, got[1].State)
	}
	if len(ps) != 1 || ps[0].ChangeSet != csB || ps[0].Entry != fid(11) || !errors.Is(ps[0].Err, ErrOrder) {
		t.Errorf("problems = %v", ps)
	}
}

func TestFoldLogSupersedeUnknown(t *testing.T) {
	l := store.Log{Entries: []entry.Entry{
		mk(t, csA, 1, intention()),
		mk(t, csA, 2, Supersede{By: csID(77)}, 1),
	}}
	got, ps := FoldLog(l)
	if len(got) != 1 || got[0].State != StateSuperseded {
		t.Fatalf("FoldLog = %+v", got)
	}
	if len(ps) != 1 || ps[0].ChangeSet != csA || ps[0].Entry != fid(2) || !errors.Is(ps[0].Err, ErrUnknownChangeSet) {
		t.Fatalf("problems = %v", ps)
	}
}

func TestFoldLogEmpty(t *testing.T) {
	got, ps := FoldLog(store.Log{})
	if got != nil || ps != nil {
		t.Fatalf("FoldLog(empty) = %v, %v", got, ps)
	}
}

// genID draws a random non-zero ID.
func genID(t *rapid.T, label string) id.ID {
	var i id.ID
	copy(i[:], rapid.SliceOfN(rapid.Byte(), 16, 16).Draw(t, label))
	i[0] |= 0x80 // never zero, never equal to a fixture
	return i
}

// genHistory draws a change set history: an open, then entries each after a
// random non-empty set of earlier ones, so branches are concurrent. Most
// entries are well formed; lifecycle moves are random, so illegal moves,
// bad revisions and concurrency reports all occur. A few entries have bad
// bodies, dangling after links or belong to another change set.
func genHistory(t *rapid.T) []entry.Entry {
	n := rapid.IntRange(0, 12).Draw(t, "entries")
	eids := rapid.SliceOfNDistinct(rapid.Custom(func(t *rapid.T) id.ID { return genID(t, "id") }), n+1, n+1,
		func(i id.ID) id.ID { return i }).Draw(t, "ids")
	kind := rapid.SampledFrom(Kinds()).Draw(t, "kind")
	at := time.Unix(1_800_000_000, 0).UTC()
	build := func(i int, cs id.ID, b Body, after []id.ID) entry.Entry {
		raw, err := EncodeBody(b)
		if err != nil {
			t.Fatalf("encode %T: %v", b, err)
		}
		slices.SortFunc(after, id.Compare)
		if after == nil {
			after = []id.ID{}
		}
		return entry.Entry{V: 1, ID: eids[i], ChangeSet: cs, Actor: actor, For: owner, At: at, After: after, Type: b.Type(), Body: raw}
	}
	es := []entry.Entry{build(0, csA, Open{Kind: kind, Intent: "i"}, nil)}
	types := []entry.Type{entry.TypeOps, entry.TypeOps, entry.TypeOps, entry.TypeNote, entry.TypePropose, entry.TypeAccept,
		entry.TypeReject, entry.TypeWithdraw, entry.TypeSupersede, entry.TypeLand, entry.TypeOpen}
	for i := 1; i <= n; i++ {
		var after []id.ID
		for _, j := range rapid.SliceOfNDistinct(rapid.IntRange(0, i-1), 1, 3, func(j int) int { return j }).Draw(t, "after") {
			after = append(after, eids[j])
		}
		var b Body
		switch typ := rapid.SampledFrom(types).Draw(t, "type"); typ {
		case entry.TypeOps:
			pool := rapid.SliceOfNDistinct(rapid.IntRange(1, 5), 1, 3, func(j int) int { return j }).Draw(t, "op ids")
			var o Ops
			for _, k := range pool {
				if rapid.IntRange(0, 3).Draw(t, "drop?") == 0 {
					o.Drop = append(o.Drop, opID(k))
				} else {
					node := rapid.IntRange(1, 4).Draw(t, "node")
					if node == 4 {
						o.Ops = append(o.Ops, move(k, 1))
					} else {
						o.Ops = append(o.Ops, retire(k, node, rapid.SampledFrom([]string{"a", "b"}).Draw(t, "reason")))
					}
				}
			}
			slices.SortFunc(o.Drop, id.Compare)
			b = o
		case entry.TypeNote:
			b = Note{Text: "n"}
		case entry.TypePropose:
			b = Propose{}
		case entry.TypeAccept:
			b = acceptBody
		case entry.TypeReject:
			b = rejectBody
		case entry.TypeWithdraw:
			b = Withdraw{}
		case entry.TypeSupersede:
			b = Supersede{By: rapid.SampledFrom([]id.ID{csA, csB}).Draw(t, "by")}
		case entry.TypeLand:
			b = Land{}
		case entry.TypeOpen:
			b = Open{Kind: kind, Intent: "again"}
			if rapid.Bool().Draw(t, "open after?") {
				after = nil
			}
		}
		cs := csA
		x := build(i, cs, b, after)
		switch rapid.IntRange(0, 19).Draw(t, "damage") {
		case 0:
			x.Body = json.RawMessage(`{"bad": true}`)
		case 1:
			x.After = append(slices.Clone(x.After), genID(t, "dangling"))
			slices.SortFunc(x.After, id.Compare)
		case 2:
			x.ChangeSet = csB
		}
		es = append(es, x)
	}
	return es
}

type textProblem struct {
	ChangeSet, Entry id.ID
	Err              string
}

func textProblems(ps []Problem) []textProblem {
	var out []textProblem
	for _, p := range ps {
		out = append(out, textProblem{p.ChangeSet, p.Entry, p.Err.Error()})
	}
	return out
}

func TestFoldOrderIndependentProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		es := genHistory(t)
		want, wantP := Fold(csA, es)
		got, gotP := Fold(csA, rapid.Permutation(es).Draw(t, "shuffled"))
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ChangeSet differs:\n got %+v\nwant %+v", got, want)
		}
		if g, w := textProblems(gotP), textProblems(wantP); !reflect.DeepEqual(g, w) {
			t.Fatalf("problems differ:\n got %v\nwant %v", g, w)
		}
	})
}

func TestFoldLogOrderIndependentProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		es := genHistory(t)
		want, wantP := FoldLog(store.Log{Entries: es})
		got, gotP := FoldLog(store.Log{Entries: rapid.Permutation(es).Draw(t, "shuffled")})
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("change sets differ:\n got %+v\nwant %+v", got, want)
		}
		if g, w := textProblems(gotP), textProblems(wantP); !reflect.DeepEqual(g, w) {
			t.Fatalf("problems differ:\n got %v\nwant %v", g, w)
		}
	})
}

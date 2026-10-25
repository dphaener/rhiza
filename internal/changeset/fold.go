package changeset

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
	"github.com/dphaener/genesis/internal/store"
)

// Folding turns a change set's entries into a ChangeSet. It is a pure
// function of the set of entries: the order they are passed in never
// matters, and it never fails. Whatever is wrong with an entry becomes a
// Problem and the fold goes on.
//
// Fold works in five steps.
//
//  1. Each body is decoded with DecodeBody. An entry whose body does not
//     decode is reported (ErrBody) and leaves the graph.
//  2. The graph is the entries that are structurally sound. The open entry
//     is the open with an empty after list and the smallest entry ID; every
//     other open is reported (ErrOrder) and leaves the graph. Every other
//     entry needs a non-empty after list naming entries of this change set
//     that are in the graph; an entry with an empty after list, an after
//     naming an entry that is missing or has left the graph, or an after
//     that leads into a cycle is reported (ErrOrder) and leaves the graph.
//     The causal order is Kahn's algorithm over after links with ties broken
//     by the smallest entry ID.
//  3. Each entry in the graph is applied in causal order to the lifecycle
//     state, which starts as draft after the open:
//
//     entry      allowed in                         result
//     ops        draft, proposed, accepted          revision applied, state kept
//     propose    draft                              proposed
//     accept     proposed                           accepted
//     reject     proposed                           rejected
//     withdraw   draft, proposed, accepted          withdrawn
//     supersede  draft, proposed, accepted          superseded
//     land       intention: accepted                landed
//                record: draft, proposed, accepted
//     note       any state                          nothing
//
//     Any other move, and a supersede by the change set itself, is skipped
//     and reported (ErrTransition). A skipped entry stays in the graph, so
//     the entries after it are still applied.
//  4. An ops entry is a revision. A new op ID is appended; a repeated ID
//     replaces the op in place and must keep its type and target; a drop
//     removes an op, and a dropped ID may never be written or dropped again;
//     a drop of an ID never written is an error. Any violation skips the
//     whole entry (ErrRevision): an entry applies entirely or not at all.
//  5. Two entries in the graph are concurrent when neither reaches the other
//     through after links. Each concurrent pair of lifecycle entries
//     (propose, accept, reject, withdraw, supersede, land) is reported, as
//     is each concurrent pair of ops entries that write (add, replace or
//     drop) a common op ID (ErrConcurrent). The report goes on the later
//     entry in fold order. Both entries still went through steps 3 and 4 in
//     fold order, so the later one wins when its move is legal after the
//     earlier one, and is reported as ErrTransition too when it is not (two
//     concurrent terminal moves: the first in fold order stands).

var (
	// ErrBody is wrapped by a Problem's error for an entry whose body does
	// not decode. The error also wraps DecodeBody's error, so errors.Is
	// matches ErrInvalid (and op.ErrInvalid for a bad operation).
	ErrBody = errors.New("changeset: bad body")
	// ErrOrder is wrapped for an entry that leaves the graph: it belongs to
	// another change set, repeats an entry ID, is a second open or an open
	// with after links, has an empty after list, names a missing or skipped
	// entry, or is in or after a cycle. It is also the error of the
	// change-set-level Problem (zero Entry) reported when there is no open.
	ErrOrder = errors.New("changeset: entry out of order")
	// ErrTransition is wrapped for a lifecycle move the state does not
	// allow, including a supersede by the change set itself.
	ErrTransition = errors.New("changeset: illegal lifecycle move")
	// ErrRevision is wrapped for an ops entry that changes an op's type or
	// target, writes or drops a dropped op ID, or drops an ID never written.
	ErrRevision = errors.New("changeset: illegal revision")
	// ErrConcurrent is wrapped for two concurrent lifecycle entries, or two
	// concurrent ops entries writing the same op ID.
	ErrConcurrent = errors.New("changeset: concurrent entries")
	// ErrUnknownChangeSet is wrapped by FoldLog for an applied supersede
	// whose by names a change set with no entries in the log.
	ErrUnknownChangeSet = errors.New("changeset: unknown change set")
)

// State is a change set's lifecycle state.
type State string

// The seven lifecycle states.
const (
	StateDraft      State = "draft"
	StateProposed   State = "proposed"
	StateAccepted   State = "accepted"
	StateLanded     State = "landed"
	StateRejected   State = "rejected"
	StateWithdrawn  State = "withdrawn"
	StateSuperseded State = "superseded"
)

// States returns the seven states in their documented order.
func States() []State {
	return []State{StateDraft, StateProposed, StateAccepted, StateLanded, StateRejected, StateWithdrawn, StateSuperseded}
}

// Terminal reports whether s is final: landed, rejected, withdrawn or
// superseded.
func (s State) Terminal() bool {
	switch s {
	case StateLanded, StateRejected, StateWithdrawn, StateSuperseded:
		return true
	}
	return false
}

// inFlight reports whether s is draft, proposed or accepted.
func (s State) inFlight() bool {
	return s == StateDraft || s == StateProposed || s == StateAccepted
}

// Decision is an applied accept or reject: who decided, for whom, under
// which rule.
type Decision struct {
	// Entry is the accept or reject entry.
	Entry id.ID
	Actor id.ID
	For   id.ID
	// Type is entry.TypeAccept or entry.TypeReject.
	Type   entry.Type
	Rule   string
	Reason string
}

// ChangeSet is where a change set stands after folding. Empty slices are
// nil, so folded values compare with reflect.DeepEqual.
type ChangeSet struct {
	ID id.ID
	// Open is the open entry's body and OpenEntry its ID. Both are zero, and
	// so is State, when the change set has no valid open; then nothing else
	// is set either.
	Open      Open
	OpenEntry id.ID
	State     State
	// Ops are the change set's current operations: in the order their IDs
	// were first written, replacements in place, dropped ones removed.
	Ops []op.Op
	// History holds, for each current op ID, every version of the op that
	// was written, in fold order, each with the entry that wrote it. The
	// last version is the one in Ops. Drops holds, for each dropped op ID,
	// the entry that dropped it, the entry that first wrote it and the last
	// version before the drop. Only ops entries that applied contribute.
	// Both live in memory only: they are worked out from the log, never
	// written to it, and are nil when empty.
	History map[id.ID][]Revision
	Drops   map[id.ID]Revision
	// Decided is the applied accept or reject, nil when there is none.
	Decided *Decision
	// SupersededBy is the applied supersede's by, zero unless the state is
	// superseded.
	SupersededBy id.ID
	// TerminalEntry is the entry that moved the change set to its terminal
	// state, zero while it is not terminal.
	TerminalEntry id.ID
	// Heads are the entries in the graph that no entry in the graph lists in
	// its after, sorted by ID. A new entry would name them in its after.
	Heads []id.ID
	// Applied is every entry in the graph in fold (causal) order, the open
	// first. It includes entries skipped in step 3 or 4 (ErrTransition,
	// ErrRevision), which were folded but changed nothing; it excludes
	// entries that left the graph.
	Applied []id.ID
}

// Revision is one version of an op, as written by an ops entry.
type Revision struct {
	// Entry is the ops entry that wrote the version, or that dropped the op.
	Entry id.ID
	// First is set in a Drops record only: the entry that first wrote the
	// dropped op, so that a reader can tell whether it existed at a given
	// entry. It is zero in History, whose first element says the same.
	First id.ID
	Op    op.Op
}

// Problem is one thing wrong with a change set's entries. Entry is zero for
// a problem with the change set as a whole. Err wraps one of the package's
// sentinels and is never nil.
type Problem struct {
	ChangeSet id.ID
	Entry     id.ID
	Err       error
}

func (p Problem) String() string {
	return fmt.Sprintf("%s/%s: %v", p.ChangeSet, p.Entry, p.Err)
}

func sortProblems(ps []Problem) []Problem {
	if len(ps) == 0 {
		return nil
	}
	slices.SortFunc(ps, func(a, b Problem) int {
		return cmp.Or(id.Compare(a.ChangeSet, b.ChangeSet), id.Compare(a.Entry, b.Entry),
			strings.Compare(a.Err.Error(), b.Err.Error()))
	})
	return ps
}

// node is one decoded entry.
type node struct {
	e    entry.Entry
	body Body
	// after is e.After sorted and unique.
	after []id.ID
}

type folder struct {
	cs       id.ID
	problems []Problem
}

func (f *folder) report(eid id.ID, err error) {
	f.problems = append(f.problems, Problem{ChangeSet: f.cs, Entry: eid, Err: err})
}

// compareEntries orders two entries with the same ID, so that which copy of
// a duplicate is kept does not depend on input order.
func compareEntries(a, b entry.Entry) int {
	return cmp.Or(
		strings.Compare(string(a.Type), string(b.Type)),
		bytes.Compare(a.Body, b.Body),
		slices.CompareFunc(a.After, b.After, id.Compare),
		id.Compare(a.Actor, b.Actor),
		id.Compare(a.For, b.For),
		a.At.Compare(b.At),
		cmp.Compare(a.V, b.V),
	)
}

// Fold folds the entries of change set cs into a ChangeSet, as described at
// the top of this file. Entries of another change set are reported
// (ErrOrder) and ignored. Problems are sorted by change set, entry ID, then
// error text, and are nil when there are none. Fold never panics.
func Fold(cs id.ID, entries []entry.Entry) (ChangeSet, []Problem) {
	f := &folder{cs: cs}
	out := ChangeSet{ID: cs}

	// Step 1: keep this change set's entries, one per ID, and decode them.
	byID := make(map[id.ID]entry.Entry, len(entries))
	for _, e := range entries {
		if e.ChangeSet != cs {
			f.report(e.ID, fmt.Errorf("%w: entry belongs to change set %s", ErrOrder, e.ChangeSet))
			continue
		}
		if prev, ok := byID[e.ID]; ok {
			f.report(e.ID, fmt.Errorf("%w: entry ID appears more than once", ErrOrder))
			if compareEntries(e, prev) < 0 {
				byID[e.ID] = e
			}
			continue
		}
		byID[e.ID] = e
	}
	eids := make([]id.ID, 0, len(byID))
	for x := range byID {
		eids = append(eids, x)
	}
	slices.SortFunc(eids, id.Compare)

	nodes := make(map[id.ID]*node, len(eids))
	var opens []id.ID
	for _, x := range eids {
		e := byID[x]
		b, err := DecodeBody(e)
		if err != nil {
			f.report(x, fmt.Errorf("%w: %w", ErrBody, err))
			continue
		}
		after := slices.Clone(e.After)
		slices.SortFunc(after, id.Compare)
		after = slices.Compact(after)
		nodes[x] = &node{e: e, body: b, after: after}
		if _, ok := b.(Open); ok {
			opens = append(opens, x)
		}
	}

	// Step 2: the structural graph. opens is sorted by ID.
	var root id.ID
	for _, x := range opens {
		if len(nodes[x].after) == 0 {
			root = x
			break
		}
	}
	for _, x := range opens {
		switch {
		case x == root:
		case len(nodes[x].after) > 0:
			f.report(x, fmt.Errorf("%w: open must come first, with an empty after", ErrOrder))
			delete(nodes, x)
		default:
			f.report(x, fmt.Errorf("%w: a second open; the change set is opened by %s", ErrOrder, root))
			delete(nodes, x)
		}
	}
	if root.IsZero() {
		f.report(id.ID{}, fmt.Errorf("%w: change set has no open entry with an empty after", ErrOrder))
		for _, x := range eids {
			if _, ok := nodes[x]; ok {
				f.report(x, fmt.Errorf("%w: change set has no open entry", ErrOrder))
			}
		}
		return out, sortProblems(f.problems)
	}

	// Entries with an empty after, then everything that names an entry not
	// in the graph, transitively. The result is a set, so it does not depend
	// on iteration order; the reasons are worked out from the final set.
	out1 := map[id.ID]bool{}
	for _, x := range eids {
		if n, ok := nodes[x]; ok && x != root && len(n.after) == 0 {
			f.report(x, fmt.Errorf("%w: after is empty; only the open entry has no after", ErrOrder))
			out1[x] = true
			delete(nodes, x)
		}
	}
	inGraph := func(x id.ID) bool { _, ok := nodes[x]; return ok }
	for changed := true; changed; {
		changed = false
		for _, x := range eids {
			n, ok := nodes[x]
			if !ok {
				continue
			}
			if slices.ContainsFunc(n.after, func(a id.ID) bool { return !inGraph(a) }) {
				out1[x] = true
				changed = true
			}
		}
		for x := range out1 {
			delete(nodes, x)
		}
	}
	for _, x := range eids {
		if !out1[x] {
			continue
		}
		e := byID[x]
		after := slices.Clone(e.After)
		slices.SortFunc(after, id.Compare)
		for _, a := range after {
			if _, ok := byID[a]; !ok {
				f.report(x, fmt.Errorf("%w: after names %s, which is not an entry of this change set", ErrOrder, a))
				break
			}
			if !inGraph(a) {
				f.report(x, fmt.Errorf("%w: after names %s, which was left out of the fold", ErrOrder, a))
				break
			}
		}
		// An entry with an empty after was reported above and has no
		// offending link.
	}

	// Kahn's algorithm, smallest ID first.
	indeg := make(map[id.ID]int, len(nodes))
	children := make(map[id.ID][]id.ID, len(nodes))
	for _, x := range eids {
		n, ok := nodes[x]
		if !ok {
			continue
		}
		indeg[x] = len(n.after)
		for _, a := range n.after {
			children[a] = append(children[a], x)
		}
	}
	var order []*node
	ready := []id.ID{root}
	for len(ready) > 0 {
		i := 0
		for j := range ready {
			if id.Compare(ready[j], ready[i]) < 0 {
				i = j
			}
		}
		x := ready[i]
		ready = slices.Delete(ready, i, i+1)
		order = append(order, nodes[x])
		for _, c := range children[x] {
			indeg[c]--
			if indeg[c] == 0 {
				ready = append(ready, c)
			}
		}
	}
	for _, x := range eids {
		if _, ok := nodes[x]; ok && indeg[x] > 0 {
			f.report(x, fmt.Errorf("%w: in or after a cycle of after links", ErrOrder))
			delete(nodes, x)
		}
	}

	// Ancestors, as bit sets over positions in order.
	pos := make(map[id.ID]int, len(order))
	for i, n := range order {
		pos[n.e.ID] = i
	}
	words := (len(order) + 63) / 64
	anc := make([][]uint64, len(order))
	for i, n := range order {
		anc[i] = make([]uint64, words)
		for _, a := range n.after {
			p := pos[a]
			anc[i][p/64] |= 1 << (p % 64)
			for w := range anc[i] {
				anc[i][w] |= anc[p][w]
			}
		}
	}

	// Step 5: concurrency reports.
	for j := range order {
		for i := range j {
			if anc[j][i/64]&(1<<(i%64)) != 0 {
				continue
			}
			a, b := order[i], order[j]
			switch {
			case lifecycle(a.e.Type) && lifecycle(b.e.Type):
				f.report(b.e.ID, fmt.Errorf("%w: %s %s and %s %s", ErrConcurrent, a.e.Type, a.e.ID, b.e.Type, b.e.ID))
			case a.e.Type == entry.TypeOps && b.e.Type == entry.TypeOps:
				if shared, ok := sharedOp(a.body.(Ops), b.body.(Ops)); ok {
					f.report(b.e.ID, fmt.Errorf("%w: ops %s and ops %s both write op %s", ErrConcurrent, a.e.ID, b.e.ID, shared))
				}
			}
		}
	}

	// Steps 3 and 4: apply in order.
	rootNode := nodes[root]
	out.Open = rootNode.body.(Open)
	out.OpenEntry = root
	out.State = StateDraft
	rev := revision{pos: map[id.ID]int{}, dropped: map[id.ID]bool{},
		history: map[id.ID][]Revision{}, drops: map[id.ID]Revision{}}
	for _, n := range order[1:] {
		if err := f.apply(&out, &rev, n); err != nil {
			f.report(n.e.ID, err)
		}
	}
	if len(rev.ops) > 0 {
		out.Ops = rev.ops
	}
	if len(rev.history) > 0 {
		out.History = rev.history
	}
	if len(rev.drops) > 0 {
		out.Drops = rev.drops
	}

	listed := map[id.ID]bool{}
	for _, n := range order {
		out.Applied = append(out.Applied, n.e.ID)
		for _, a := range n.after {
			listed[a] = true
		}
	}
	for _, n := range order {
		if !listed[n.e.ID] {
			out.Heads = append(out.Heads, n.e.ID)
		}
	}
	slices.SortFunc(out.Heads, id.Compare)
	return out, sortProblems(f.problems)
}

func lifecycle(t entry.Type) bool {
	switch t {
	case entry.TypePropose, entry.TypeAccept, entry.TypeReject, entry.TypeWithdraw, entry.TypeSupersede, entry.TypeLand:
		return true
	}
	return false
}

// written lists the op IDs an ops body adds, replaces or drops.
func written(b Ops) []id.ID {
	out := slices.Clone(b.Drop)
	for _, o := range b.Ops {
		out = append(out, o.ID)
	}
	slices.SortFunc(out, id.Compare)
	return out
}

// sharedOp returns the smallest op ID both bodies write.
func sharedOp(a, b Ops) (id.ID, bool) {
	bs := written(b)
	for _, x := range written(a) {
		if _, ok := slices.BinarySearchFunc(bs, x, id.Compare); ok {
			return x, true
		}
	}
	return id.ID{}, false
}

func transition(t entry.Type, s State) error {
	return fmt.Errorf("%w: %s in state %s", ErrTransition, t, s)
}

// apply applies one entry in the graph to out, or returns why it is skipped.
func (f *folder) apply(out *ChangeSet, rev *revision, n *node) error {
	from := out.State
	terminal := func(s State) {
		out.State = s
		out.TerminalEntry = n.e.ID
	}
	decided := func(t entry.Type, rule, reason string) {
		out.Decided = &Decision{Entry: n.e.ID, Actor: n.e.Actor, For: n.e.For, Type: t, Rule: rule, Reason: reason}
	}
	switch b := n.body.(type) {
	case Ops:
		if !from.inFlight() {
			return transition(n.e.Type, from)
		}
		if err := rev.check(b); err != nil {
			return err
		}
		rev.apply(n.e.ID, b)
	case Propose:
		if from != StateDraft {
			return transition(n.e.Type, from)
		}
		out.State = StateProposed
	case Accept:
		if from != StateProposed {
			return transition(n.e.Type, from)
		}
		out.State = StateAccepted
		decided(entry.TypeAccept, b.Rule, b.Reason)
	case Reject:
		if from != StateProposed {
			return transition(n.e.Type, from)
		}
		terminal(StateRejected)
		decided(entry.TypeReject, b.Rule, b.Reason)
	case Withdraw:
		if !from.inFlight() {
			return transition(n.e.Type, from)
		}
		terminal(StateWithdrawn)
	case Supersede:
		if b.By == out.ID {
			return fmt.Errorf("%w: a change set cannot supersede itself", ErrTransition)
		}
		if !from.inFlight() {
			return transition(n.e.Type, from)
		}
		terminal(StateSuperseded)
		out.SupersededBy = b.By
	case Land:
		if out.Open.Kind == KindRecord && !from.inFlight() ||
			out.Open.Kind != KindRecord && from != StateAccepted {
			return fmt.Errorf("%w: land of %s change set in state %s", ErrTransition, out.Open.Kind, from)
		}
		terminal(StateLanded)
	case Note:
	}
	return nil
}

// revision is the op list as revisions build it.
type revision struct {
	ops []op.Op
	// pos maps each current op ID to its index in ops.
	pos map[id.ID]int
	// dropped holds every dropped op ID.
	dropped map[id.ID]bool
	// history and drops are ChangeSet.History and ChangeSet.Drops.
	history map[id.ID][]Revision
	drops   map[id.ID]Revision
}

// check reports the first rule b breaks against the current ops. Body
// validation already rules out repeated IDs inside b.
func (r *revision) check(b Ops) error {
	for _, o := range b.Ops {
		if r.dropped[o.ID] {
			return fmt.Errorf("%w: op %s was dropped and cannot be written again", ErrRevision, o.ID)
		}
		i, ok := r.pos[o.ID]
		if !ok {
			continue
		}
		old := r.ops[i]
		if old.Type != o.Type {
			return fmt.Errorf("%w: op %s changes type from %s to %s", ErrRevision, o.ID, old.Type, o.Type)
		}
		if old.Target() != o.Target() {
			return fmt.Errorf("%w: op %s changes target from %s to %s", ErrRevision, o.ID, old.Target(), o.Target())
		}
	}
	for _, x := range b.Drop {
		if r.dropped[x] {
			return fmt.Errorf("%w: op %s is already dropped", ErrRevision, x)
		}
		if _, ok := r.pos[x]; !ok {
			return fmt.Errorf("%w: drop names op %s, which this change set never wrote", ErrRevision, x)
		}
	}
	return nil
}

// apply applies b, written by entry eid, which check accepted.
func (r *revision) apply(eid id.ID, b Ops) {
	for _, o := range b.Ops {
		r.history[o.ID] = append(r.history[o.ID], Revision{Entry: eid, Op: o})
		if i, ok := r.pos[o.ID]; ok {
			r.ops[i] = o
			continue
		}
		r.pos[o.ID] = len(r.ops)
		r.ops = append(r.ops, o)
	}
	if len(b.Drop) == 0 {
		return
	}
	for _, x := range b.Drop {
		r.dropped[x] = true
		r.drops[x] = Revision{Entry: eid, First: r.history[x][0].Entry, Op: r.ops[r.pos[x]]}
		delete(r.history, x)
	}
	r.ops = slices.DeleteFunc(r.ops, func(o op.Op) bool { return r.dropped[o.ID] })
	clear(r.pos)
	for i, o := range r.ops {
		r.pos[o.ID] = i
	}
}

// FoldLog folds every change set in l, sorted by change set ID, and returns
// their problems sorted by change set, entry ID, then error text. On top of
// Fold's problems it reports an applied supersede whose by names a change
// set with no entries in l (ErrUnknownChangeSet), on the supersede entry.
// The log's own Problems are not copied: the caller already has them. Both
// results are nil when there is nothing to report.
func FoldLog(l store.Log) ([]ChangeSet, []Problem) {
	groups := map[id.ID][]entry.Entry{}
	for _, e := range l.Entries {
		groups[e.ChangeSet] = append(groups[e.ChangeSet], e)
	}
	keys := make([]id.ID, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, id.Compare)

	var sets []ChangeSet
	var problems []Problem
	for _, k := range keys {
		c, ps := Fold(k, groups[k])
		sets = append(sets, c)
		problems = append(problems, ps...)
		if c.State == StateSuperseded {
			if _, ok := groups[c.SupersededBy]; !ok {
				problems = append(problems, Problem{ChangeSet: k, Entry: c.TerminalEntry,
					Err: fmt.Errorf("%w: superseded by %s, which has no entries in the log", ErrUnknownChangeSet, c.SupersededBy)})
			}
		}
	}
	return sets, sortProblems(problems)
}

package changeset

// Revisions in place under a dependent. A change set B sees the latest
// content of every change set it depends on, revisions in place included:
// an op revised under its own ID keeps that ID, so a base naming it still
// matches. That would let B overwrite a value it never saw, so a revision
// is checked against B's writes.
//
// For each change set A that B depends on, directly or through others,
// landed or in flight (Dependencies), and each current op o of B, let eB
// be the entry that last wrote o. An op X of A on the same node or link as
// o has a version as of eB: the last one written by an entry of A that
// sorts before eB. When X had no version yet, B did not write against it
// and nothing is checked. Otherwise the aspects X's current version changes
// from that one (each field whose value differs, the handle, the position
// as parent and key) are compared with the aspects o writes, under the
// same rule as the model's conflicts: they overlap when they share one, or
// when o retires and X changed anything. A retire's revision changes only
// its reason, which is no aspect. Overlap is a conflict; no overlap, no
// conflict. When A drops, after eB, an op X it had written before eB, and o
// names (as target, parent, from, to or base) X's ID or the node or link X
// creates, that is a conflict too, unless X created what o names and a
// current op in B's view still creates it under the same ID: an op of B,
// of a change set in B's in-flight closure (Closure) or of a landed one.
// That is a carve: the create moved to another change set A depends on,
// and B still sees the node or link through it. o naming X's ID itself
// still conflicts.
//
// B writing o again clears its conflicts: eB moves past A's revision. Which
// entry came first is decided by entry ID order only.

import (
	"bytes"
	"cmp"
	"encoding/json"
	"maps"
	"slices"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// The aspects an operation writes, spelled as internal/model spells them.
const (
	aspectField    = "field:"
	aspectHandle   = "handle"
	aspectPosition = "position"
	aspectRetired  = "retired"
)

// RevisedConflict is one op of a change set that a dependency's revision
// in place, or drop, has gone under.
type RevisedConflict struct {
	// ChangeSet is the dependent change set and Op its op.
	ChangeSet id.ID
	Op        id.ID
	// Dependency is the change set ChangeSet depends on (landed or in
	// flight) that revised or dropped DependencyOp.
	Dependency   id.ID
	DependencyOp id.ID
	// Target is the node or link both ops write, or for a drop the name Op
	// rests on that went with DependencyOp.
	Target id.ID
	// Aspects are the overlapping aspects the revision changed, sorted; nil
	// for a drop.
	Aspects []string
	// Dropped is set when DependencyOp was dropped.
	Dropped bool
}

// Revised returns the conflicts between the current ops of change set cs
// and the revisions and drops of its landed and in-flight dependencies, as
// described at the top of this file. They are sorted by cs's op, then the
// dependency, then its op, and are nil when there are none. A cs not in
// sets has no ops and so no conflicts.
func Revised(sets []ChangeSet, cs id.ID) []RevisedConflict {
	return Analyze(sets).Revised(cs)
}

// Revised is Revised(sets, cs) for a's sets.
func (a *Analysis) Revised(cs id.ID) []RevisedConflict {
	b, ok := a.set(cs)
	if !ok || len(b.Ops) == 0 {
		return nil
	}
	var out []RevisedConflict
	var view map[id.ID]bool
	created := func(n id.ID) bool {
		if view == nil {
			view = a.viewCreates(cs)
		}
		return view[n]
	}
	for _, depID := range a.Dependencies(cs) {
		dep, _ := a.set(depID)
		for _, o := range b.Ops {
			h := b.History[o.ID]
			if len(h) == 0 {
				continue
			}
			eB := h[len(h)-1].Entry
			out = append(out, revisedUnder(cs, o, eB, dep)...)
			out = append(out, droppedUnder(cs, o, eB, dep, created)...)
		}
	}
	if len(out) == 0 {
		return nil
	}
	slices.SortFunc(out, func(x, y RevisedConflict) int {
		return cmp.Or(id.Compare(x.Op, y.Op), id.Compare(x.Dependency, y.Dependency),
			id.Compare(x.DependencyOp, y.DependencyOp))
	})
	return out
}

// revisedUnder finds a's ops on o's target revised after eB in a way that
// overlaps what o writes.
func revisedUnder(cs id.ID, o op.Op, eB id.ID, a ChangeSet) []RevisedConflict {
	writes := aspects(o)
	var out []RevisedConflict
	for _, x := range a.Ops {
		if x.Target() != o.Target() {
			continue
		}
		h := a.History[x.ID]
		i := slices.IndexFunc(h, func(r Revision) bool { return id.Compare(r.Entry, eB) >= 0 })
		if i == -1 {
			i = len(h)
		}
		if i == 0 {
			continue // X did not exist when o was written.
		}
		changed := changes(h[i-1].Op, x)
		if !overlaps(changed, writes) {
			continue
		}
		var shown []string
		if slices.Contains(writes, aspectRetired) {
			shown = changed
		} else {
			for _, s := range changed {
				if slices.Contains(writes, s) {
					shown = append(shown, s)
				}
			}
		}
		out = append(out, RevisedConflict{ChangeSet: cs, Op: o.ID, Dependency: a.ID, DependencyOp: x.ID,
			Target: o.Target(), Aspects: shown})
	}
	return out
}

// droppedUnder finds a's ops dropped after eB, written before it, whose ID
// or created node or link o names; a created node or link for which
// created reports true does not count.
func droppedUnder(cs id.ID, o op.Op, eB id.ID, a ChangeSet, created func(id.ID) bool) []RevisedConflict {
	var out []RevisedConflict
	for xID, d := range a.Drops {
		if id.Compare(d.Entry, eB) <= 0 || id.Compare(d.First, eB) >= 0 {
			continue
		}
		made, makes := creates(d.Op)
		for _, n := range Names(o) {
			if n == xID || makes && n == made && !created(n) {
				out = append(out, RevisedConflict{ChangeSet: cs, Op: o.ID, Dependency: a.ID, DependencyOp: xID,
					Target: n, Dropped: true})
				break
			}
		}
	}
	return out
}

// viewCreates is every node and link a current op in cs's view creates:
// cs's own ops, those of its in-flight closure and those of landed change
// sets.
func (a *Analysis) viewCreates(cs id.ID) map[id.ID]bool {
	a.graph()
	out := map[id.ID]bool{}
	add := func(os []op.Op) {
		for _, o := range os {
			if n, ok := creates(o); ok {
				out[n] = true
			}
		}
	}
	add(a.ops[cs])
	for _, c := range a.Closure(cs) {
		add(a.ops[c])
	}
	for _, s := range a.sets {
		if s.State == StateLanded {
			add(s.Ops)
		}
	}
	return out
}

// aspects returns the aspects o writes, sorted and unique; nil for a
// create. It matches internal/model.
func aspects(o op.Op) []string {
	var out []string
	switch o.Type {
	case op.TypeNodeUpdate:
		for k := range o.Set {
			out = append(out, aspectField+k)
		}
		for _, k := range o.Unset {
			out = append(out, aspectField+k)
		}
		if o.Handle != "" {
			out = append(out, aspectHandle)
		}
	case op.TypeNodeMove:
		out = append(out, aspectPosition)
	case op.TypeNodeRetire, op.TypeLinkRetire:
		out = append(out, aspectRetired)
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// overlaps reports whether two sets of aspects overlap: they share one, or
// either is a retirement and the other writes anything. It matches
// internal/model.
func overlaps(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if slices.Contains(a, aspectRetired) || slices.Contains(b, aspectRetired) {
		return true
	}
	return slices.ContainsFunc(a, func(s string) bool { return slices.Contains(b, s) })
}

// changes returns the aspects that differ between two versions of one op,
// sorted: each field whose written value differs, the handle, and the
// position (parent and key). A node.create writes its fields; a
// node.update writes a field it sets or unsets, and an unset differs from
// any set. A retire's or a link's revision changes no aspect.
func changes(old, cur op.Op) []string {
	var out []string
	switch cur.Type {
	case op.TypeNodeCreate:
		for _, k := range keys(old.Fields, cur.Fields, nil, nil) {
			if !sameValue(old.Fields, cur.Fields, k) {
				out = append(out, aspectField+k)
			}
		}
	case op.TypeNodeUpdate:
		for _, k := range keys(old.Set, cur.Set, old.Unset, cur.Unset) {
			if slices.Contains(old.Unset, k) != slices.Contains(cur.Unset, k) || !sameValue(old.Set, cur.Set, k) {
				out = append(out, aspectField+k)
			}
		}
	}
	if cur.Type == op.TypeNodeCreate || cur.Type == op.TypeNodeUpdate {
		if old.Handle != cur.Handle {
			out = append(out, aspectHandle)
		}
	}
	if cur.Type == op.TypeNodeCreate || cur.Type == op.TypeNodeMove {
		if old.Parent != cur.Parent || old.Key != cur.Key {
			out = append(out, aspectPosition)
		}
	}
	slices.Sort(out)
	return out
}

// keys is every name in the two maps and the two lists, sorted and unique.
func keys(a, b map[string]json.RawMessage, c, d []string) []string {
	out := slices.Concat(slices.Collect(maps.Keys(a)), slices.Collect(maps.Keys(b)), c, d)
	slices.Sort(out)
	return slices.Compact(out)
}

// sameValue reports whether a and b hold the same JSON value under k, or
// both lack it. Insignificant whitespace does not count.
func sameValue(a, b map[string]json.RawMessage, k string) bool {
	va, oka := a[k]
	vb, okb := b[k]
	if oka != okb {
		return false
	}
	if !oka {
		return true
	}
	var ca, cb bytes.Buffer
	if json.Compact(&ca, va) != nil || json.Compact(&cb, vb) != nil {
		return bytes.Equal(va, vb)
	}
	return bytes.Equal(ca.Bytes(), cb.Bytes())
}

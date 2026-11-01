package changeset

// The field-level diff of a change set. BuildDiff folds the change set's ops
// into one entry per node it writes, merging every write to that node, and
// reads old values off the node as it stood at the first op's base.
//
// The node at a base is rebuilt by replaying its history: the ops of the
// other change sets on that node, along the chain of bases that leads to
// the base, the create first. Ops that landed after the base, or that no
// chain through the base reaches, are not in it, so another change set
// landing a different field first does not show up as this change set's
// change. New values are what the change set's own ops leave once applied
// in order on top of the old node. Only what differs is reported.

import (
	"bytes"
	"encoding/json"
	"slices"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// Action is what a change set does to a node.
type Action string

// The four actions of a node entry.
const (
	ActionCreated Action = "created"
	ActionUpdated Action = "updated"
	ActionMoved   Action = "moved"
	ActionRetired Action = "retired"
)

// FieldChange is one aspect of a node that differs. Name is a field name,
// or "handle", "parent" or "key" prefixed with "@" for the aspects that are
// not fields. Old and New are canonical JSON; nil means absent (a created
// node's fields have no Old, an unset field no New).
type FieldChange struct {
	Name string
	Old  json.RawMessage
	New  json.RawMessage
}

// NodeDiff is the change set's writes to one node.
type NodeDiff struct {
	Node   id.ID
	Action Action
	Kind   op.Kind
	// Handle is the node's handle after the change set.
	Handle string
	// Reason is the retire's reason, for a retired node.
	Reason  string
	Changes []FieldChange
}

// LinkDiff is a link the change set adds or retires.
type LinkDiff struct {
	Link         id.ID
	Type         op.LinkType
	From         id.ID
	To           id.ID
	FromHandle   string
	ToHandle     string
	Note         string
	Retired      bool
	RetireReason string
}

// Diff is the field-level diff of a change set.
type Diff struct {
	Nodes []NodeDiff
	Links []LinkDiff
}

// state is a node as replay rebuilds it.
type state struct {
	exists  bool
	kind    op.Kind
	handle  string
	parent  id.ID
	key     string
	fields  map[string]json.RawMessage
	retired bool
}

func (s *state) apply(o op.Op) {
	switch o.Type {
	case op.TypeNodeCreate:
		s.exists = true
		s.kind, s.handle, s.parent, s.key = o.Kind, o.Handle, o.Parent, o.Key
		s.fields = map[string]json.RawMessage{}
		for k, v := range o.Fields {
			s.fields[k] = v
		}
	case op.TypeNodeUpdate:
		if o.Handle != "" {
			s.handle = o.Handle
		}
		if s.fields == nil {
			s.fields = map[string]json.RawMessage{}
		}
		for k, v := range o.Set {
			s.fields[k] = v
		}
		for _, k := range o.Unset {
			delete(s.fields, k)
		}
	case op.TypeNodeMove:
		s.parent, s.key = o.Parent, o.Key
	case op.TypeNodeRetire:
		s.retired = true
	}
}

func (s state) clone() state {
	c := s
	c.fields = make(map[string]json.RawMessage, len(s.fields))
	for k, v := range s.fields {
		c.fields[k] = v
	}
	return c
}

// BuildDiff returns the diff of change set cs among sets. A cs not in sets has
// an empty diff.
func BuildDiff(sets []ChangeSet, cs id.ID) Diff {
	var mine ChangeSet
	found := false
	others := map[id.ID][]op.Op{}
	names := map[id.ID]string{}
	for _, s := range sets {
		if s.ID == cs {
			mine, found = s, true
			continue
		}
		if s.State == "" || s.State == StateRejected || s.State == StateWithdrawn || s.State == StateSuperseded {
			continue
		}
		for _, o := range s.Ops {
			others[o.Target()] = append(others[o.Target()], o)
		}
	}
	if !found {
		return Diff{}
	}
	for _, os := range others {
		for _, o := range os {
			noteName(names, o)
		}
	}
	for _, o := range mine.Ops {
		noteName(names, o)
	}

	var res Diff
	byNode := map[id.ID]*NodeDiff{}
	var order []id.ID
	opsOf := map[id.ID][]op.Op{}
	for _, o := range mine.Ops {
		switch o.Type {
		case op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire:
			if _, ok := opsOf[o.Node]; !ok {
				order = append(order, o.Node)
			}
			opsOf[o.Node] = append(opsOf[o.Node], o)
		}
	}
	for _, n := range order {
		nd := nodeDiff(n, opsOf[n], others[n])
		byNode[n] = &nd
		res.Nodes = append(res.Nodes, nd)
	}

	linkOps := map[id.ID]op.Op{}
	for _, os := range others {
		for _, o := range os {
			if o.Type == op.TypeLinkCreate {
				linkOps[o.Link] = o
			}
		}
	}
	for _, o := range mine.Ops {
		if o.Type == op.TypeLinkCreate {
			linkOps[o.Link] = o
		}
	}
	for _, o := range mine.Ops {
		switch o.Type {
		case op.TypeLinkCreate:
			res.Links = append(res.Links, linkDiff(o, names, false))
		case op.TypeLinkRetire:
			c, ok := linkOps[o.Link]
			if !ok {
				c = op.Op{Link: o.Link}
			}
			res.Links = append(res.Links, linkDiff(c, names, true))
		}
	}
	// A link created and retired within the change set nets out.
	res.Links = slices.DeleteFunc(res.Links, func(l LinkDiff) bool {
		if l.Retired {
			return false
		}
		return slices.ContainsFunc(res.Links, func(r LinkDiff) bool { return r.Retired && r.Link == l.Link })
	})
	res.Links = slices.DeleteFunc(res.Links, func(l LinkDiff) bool {
		if !l.Retired {
			return false
		}
		return slices.ContainsFunc(mine.Ops, func(o op.Op) bool { return o.Type == op.TypeLinkCreate && o.Link == l.Link })
	})
	return res
}

func noteName(names map[id.ID]string, o op.Op) {
	switch o.Type {
	case op.TypeNodeCreate:
		names[o.Node] = o.Handle
	case op.TypeNodeUpdate:
		if o.Handle != "" {
			names[o.Node] = o.Handle
		}
	}
}

func linkDiff(c op.Op, names map[id.ID]string, retired bool) LinkDiff {
	return LinkDiff{Link: c.Link, Type: c.LinkType, From: c.From, To: c.To,
		FromHandle: names[c.From], ToHandle: names[c.To], Note: c.Note, Retired: retired}
}

// oldState rebuilds the node at base from the other change sets' ops on
// it. With no first base (a create) it is the zero state. When base is not
// found among them, every op of theirs on the node applies, in base-chain
// order.
func oldState(base id.ID, others []op.Op) state {
	if base.IsZero() {
		return state{}
	}
	byID := map[id.ID]op.Op{}
	for _, o := range others {
		byID[o.ID] = o
	}
	var chain []op.Op
	seen := map[id.ID]bool{}
	for cur := base; !cur.IsZero() && !seen[cur]; {
		seen[cur] = true
		o, ok := byID[cur]
		if !ok {
			break
		}
		chain = append(chain, o)
		cur = o.Base
	}
	if len(chain) == 0 {
		// Base unknown: fall back to everything known on the node.
		all := slices.Clone(others)
		slices.SortStableFunc(all, func(a, b op.Op) int {
			if (a.Type == op.TypeNodeCreate) != (b.Type == op.TypeNodeCreate) {
				if a.Type == op.TypeNodeCreate {
					return -1
				}
				return 1
			}
			return id.Compare(a.ID, b.ID)
		})
		var s state
		for _, o := range all {
			s.apply(o)
		}
		return s
	}
	slices.Reverse(chain)
	var s state
	for _, o := range chain {
		s.apply(o)
	}
	return s
}

func nodeDiff(n id.ID, mine, others []op.Op) NodeDiff {
	created := slices.ContainsFunc(mine, func(o op.Op) bool { return o.Type == op.TypeNodeCreate })
	var old state
	if !created {
		for _, o := range mine {
			old = oldState(o.Base, others)
			break
		}
	}
	cur := old.clone()
	for _, o := range mine {
		cur.apply(o)
	}
	nd := NodeDiff{Node: n, Kind: cur.kind, Handle: cur.handle}
	if nd.Kind == "" {
		nd.Kind = old.kind
	}

	retired := false
	for _, o := range mine {
		if o.Type == op.TypeNodeRetire {
			retired = true
			nd.Reason = o.Reason
		}
	}
	switch {
	case created:
		nd.Action = ActionCreated
	case retired:
		nd.Action = ActionRetired
	}

	var changes []FieldChange
	names := map[string]bool{}
	for k := range old.fields {
		names[k] = true
	}
	for k := range cur.fields {
		names[k] = true
	}
	keys := make([]string, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		ov, oh := old.fields[k]
		nv, nh := cur.fields[k]
		if oh && nh && bytes.Equal(ov, nv) {
			continue
		}
		c := FieldChange{Name: k}
		if oh {
			c.Old = ov
		}
		if nh {
			c.New = nv
		}
		changes = append(changes, c)
	}
	if !created && !retired {
		moved := false
		if old.handle != cur.handle {
			changes = append(changes, FieldChange{Name: "@handle", Old: jsonString(old.handle), New: jsonString(cur.handle)})
		}
		if old.parent != cur.parent {
			moved = true
			c := FieldChange{Name: "@parent", New: parentJSON(cur.parent)}
			c.Old = parentJSON(old.parent)
			changes = append(changes, c)
		}
		if old.key != cur.key {
			moved = true
			changes = append(changes, FieldChange{Name: "@key", Old: jsonString(old.key), New: jsonString(cur.key)})
		}
		if nd.Action == "" {
			if moved && !hasNonPosition(changes) {
				nd.Action = ActionMoved
			} else {
				nd.Action = ActionUpdated
			}
		}
	}
	if nd.Action == "" {
		nd.Action = ActionUpdated
	}
	if created {
		// A create shows every aspect, with no old value.
		nd.Changes = append(nd.Changes, FieldChange{Name: "@handle", New: jsonString(cur.handle)})
		if cur.parent != (id.ID{}) {
			nd.Changes = append(nd.Changes, FieldChange{Name: "@parent", New: parentJSON(cur.parent)})
		}
		if cur.key != "" {
			nd.Changes = append(nd.Changes, FieldChange{Name: "@key", New: jsonString(cur.key)})
		}
	}
	nd.Changes = append(nd.Changes, changes...)
	return nd
}

func hasNonPosition(cs []FieldChange) bool {
	return slices.ContainsFunc(cs, func(c FieldChange) bool { return c.Name != "@parent" && c.Name != "@key" })
}

func jsonString(s string) json.RawMessage {
	b, _ := json.Marshal(s)
	return b
}

func parentJSON(p id.ID) json.RawMessage {
	if p.IsZero() {
		return nil
	}
	return jsonString(p.String())
}

package api

// This file holds the node writes: add, set, mv and retire. Each has a
// builder (build*) that turns its request into an edit against a writeBase,
// following the revision rules (D9): when the change set already has an op
// on the node, the write revises that op instead of adding another.

import (
	"encoding/json"
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/key"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "add",
		Summary:  "Add a node",
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[AddRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(AddRequest)
			return runWrite(ctx, "add", commandLine("add", r.Kind, r.Handle),
				func(b writeBase) (edit, error) { return buildAdd(b, r) })
		},
	})
	register(Operation{
		Name:     "set",
		Summary:  "Change a node's handle or fields",
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[SetRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(SetRequest)
			return runWrite(ctx, "set", commandLine("set", r.Ref),
				func(b writeBase) (edit, error) { return buildSet(b, r) })
		},
	})
	register(Operation{
		Name:     "mv",
		Summary:  "Move a capability or intent to another parent or position",
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[MvRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(MvRequest)
			return runWrite(ctx, "mv", commandLine("mv", r.Ref),
				func(b writeBase) (edit, error) { return buildMv(b, r) })
		},
	})
	register(Operation{
		Name:     "retire",
		Summary:  "Retire a node",
		Write:    true,
		Record:   true,
		Request:  reflect.TypeFor[RetireRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(RetireRequest)
			return runWrite(ctx, "retire", commandLine("retire", r.Ref),
				func(b writeBase) (edit, error) { return buildRetire(b, r) })
		},
	})
}

// AddRequest adds a node.
type AddRequest struct {
	Kind       string   `json:"kind" arg:"0" enum:"intent,capability,constraint,decision,question,evidence,source,actor,anchor" desc:"The node's kind"`
	Handle     string   `json:"handle" arg:"1" desc:"The node's handle: lower-case words joined by single dashes, unique among live nodes"`
	Parent     string   `json:"parent,omitempty" desc:"The parent node (capabilities and intents only, same kind); default a root"`
	Before     string   `json:"before,omitempty" desc:"Place it just before this sibling"`
	After      string   `json:"after,omitempty" desc:"Place it just after this sibling; default after the last one"`
	Fields     []string `json:"fields,omitempty" flag:"field" desc:"A text field, name=text"`
	JSONFields []string `json:"json_fields,omitempty" flag:"json-field" desc:"A JSON field, name=<json>"`
}

// SetRequest changes a node's handle or fields.
type SetRequest struct {
	Ref        string   `json:"ref" arg:"0" desc:"The node: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Handle     string   `json:"handle,omitempty" desc:"A new handle"`
	Fields     []string `json:"fields,omitempty" flag:"field" desc:"Set a text field, name=text"`
	JSONFields []string `json:"json_fields,omitempty" flag:"json-field" desc:"Set a JSON field, name=<json>"`
	Unset      []string `json:"unset,omitempty" desc:"Remove a field by name"`
}

// MvRequest moves a node in its tree.
type MvRequest struct {
	Ref    string `json:"ref" arg:"0" desc:"The node: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Parent string `json:"parent,omitempty" desc:"The new parent, of the same kind"`
	Root   bool   `json:"root,omitempty" desc:"Make it a root instead"`
	Before string `json:"before,omitempty" desc:"Place it just before this sibling"`
	After  string `json:"after,omitempty" desc:"Place it just after this sibling; default after the last one"`
}

// RetireRequest retires a node.
type RetireRequest struct {
	Ref    string `json:"ref" arg:"0" desc:"The node: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Reason string `json:"reason" desc:"Why it is retired"`
}

// buildAdd creates a node. A taken handle is a conflict; a parent or a
// position on a kind outside the tree is a usage error. Tree nodes go after
// their last sibling unless --before or --after names one.
func buildAdd(b writeBase, r AddRequest) (edit, error) {
	kind := op.Kind(r.Kind)
	if !slices.Contains(op.Kinds(), kind) {
		return edit{}, Usagef("unknown kind %q: want one of %s", r.Kind, joinKinds()).WithNext("rhiza describe add")
	}
	fields, err := parseFields(r.Fields, r.JSONFields)
	if err != nil {
		return edit{}, err
	}
	if err := knownFieldsOnly("add", kind, fields); err != nil {
		return edit{}, err
	}
	var parent id.ID
	var k string
	if op.HasParent(kind) {
		if r.Parent != "" {
			p, err := treeParent(b, kind, r.Parent)
			if err != nil {
				return edit{}, err
			}
			parent = p.ID
		}
		if k, err = place(b, kind, parent, id.ID{}, r.Before, r.After); err != nil {
			return edit{}, err
		}
	} else if r.Parent != "" || r.Before != "" || r.After != "" {
		return edit{}, Usagef("%s %s is not in a tree: --parent, --before and --after are for capabilities and intents", article(string(kind)), kind)
	}
	node, err := b.newID()
	if err != nil {
		return edit{}, err
	}
	opID, err := b.newID()
	if err != nil {
		return edit{}, err
	}
	o := op.Op{ID: opID, Type: op.TypeNodeCreate, Node: node, Kind: kind, Handle: r.Handle, Parent: parent, Key: k, Fields: fields}
	if err := o.Validate(); err != nil {
		return edit{}, Usagef("%s", err.Error()).WithNext("rhiza describe add")
	}
	if err := handleFree(b, r.Handle, id.ID{}); err != nil {
		return edit{}, err
	}
	e := edit{ops: []op.Op{o}}
	e.touchNode(node, r.Handle)
	return e, nil
}

// buildSet revises the node's pending create or update, else adds a
// node.update against its version.
func buildSet(b writeBase, r SetRequest) (edit, error) {
	if r.Handle == "" && len(r.Fields)+len(r.JSONFields)+len(r.Unset) == 0 {
		return edit{}, Usagef("nothing to change: give --handle, --field, --json-field or --unset").
			WithNext("rhiza describe set")
	}
	n, err := writableNode(b, r.Ref)
	if err != nil {
		return edit{}, err
	}
	set, err := parseFields(r.Fields, r.JSONFields)
	if err != nil {
		return edit{}, err
	}
	if err := knownFieldsOnly("set", n.Kind, set); err != nil {
		return edit{}, err
	}
	unset := slices.Compact(slices.Sorted(slices.Values(r.Unset)))
	for _, name := range unset {
		// Check the name by op's own rules: a probe update unsetting it.
		probe := op.Op{ID: n.ID, Type: op.TypeNodeUpdate, Node: n.ID, Base: n.ID, Unset: []string{name}}
		if err := probe.Validate(); err != nil {
			return edit{}, Usagef("--unset: %v", err).WithNext("rhiza describe set")
		}
		if _, ok := set[name]; ok {
			return edit{}, Usagef("field %q is both set and unset", name)
		}
	}
	if r.Handle != "" {
		if err := handleFree(b, r.Handle, n.ID); err != nil {
			return edit{}, err
		}
	}
	handle := n.Handle
	if r.Handle != "" {
		handle = r.Handle
	}

	var o op.Op
	if c, ok := b.pendingOp(n.Pending, op.TypeNodeCreate); ok {
		o = c
		if r.Handle != "" {
			o.Handle = r.Handle
		}
		o.Fields = mergeFields(c.Fields, set, unset)
	} else if u, ok := b.pendingOp(n.Pending, op.TypeNodeUpdate); ok {
		o = u
		if r.Handle != "" {
			o.Handle = r.Handle
		}
		o.Set = mergeFields(u.Set, set, unset)
		var keep []string
		for _, name := range u.Unset {
			if _, ok := set[name]; !ok {
				keep = append(keep, name)
			}
		}
		o.Unset = slices.Compact(slices.Sorted(slices.Values(append(keep, unset...))))
		if len(o.Unset) == 0 {
			o.Unset = nil
		}
	} else {
		if n.Version.IsZero() {
			return edit{}, Errorf("node %s has no version to change", n.ID)
		}
		opID, err := b.newID()
		if err != nil {
			return edit{}, err
		}
		o = op.Op{ID: opID, Type: op.TypeNodeUpdate, Node: n.ID, Base: n.Version, Handle: r.Handle, Set: set}
		if len(unset) > 0 {
			o.Unset = unset
		}
	}
	e := edit{ops: []op.Op{o}}
	e.touchNode(n.ID, handle)
	return e, nil
}

// buildMv revises the node's pending create or move, else adds a node.move
// against its version. The parent stays unless --parent or --root changes
// it; the node goes after the last sibling unless --before or --after names
// one.
func buildMv(b writeBase, r MvRequest) (edit, error) {
	n, err := writableNode(b, r.Ref)
	if err != nil {
		return edit{}, err
	}
	if !op.HasParent(n.Kind) {
		return edit{}, Usagef("%s %s is not in a tree: only capabilities and intents move", article(string(n.Kind)), n.Kind)
	}
	if r.Root && r.Parent != "" {
		return edit{}, Usagef("give --parent or --root, not both")
	}
	if !r.Root && r.Parent == "" && r.Before == "" && r.After == "" {
		return edit{}, Usagef("nothing to move: give --parent, --root, --before or --after").WithNext("rhiza describe mv")
	}
	parent := n.Parent
	switch {
	case r.Root:
		parent = id.ID{}
	case r.Parent != "":
		p, err := treeParent(b, n.Kind, r.Parent)
		if err != nil {
			return edit{}, err
		}
		if under(b.m, p.ID, n.ID) {
			return edit{}, Usagef("cannot move %s under itself or one of its descendants", r.Ref)
		}
		parent = p.ID
	}
	k, err := place(b, n.Kind, parent, n.ID, r.Before, r.After)
	if err != nil {
		return edit{}, err
	}

	var o op.Op
	if c, ok := b.pendingOp(n.Pending, op.TypeNodeCreate); ok {
		o = c
		o.Parent, o.Key = parent, k
	} else if m, ok := b.pendingOp(n.Pending, op.TypeNodeMove); ok {
		o = m
		o.Parent, o.Key = parent, k
	} else {
		if n.Version.IsZero() {
			return edit{}, Errorf("node %s has no version to move", n.ID)
		}
		opID, err := b.newID()
		if err != nil {
			return edit{}, err
		}
		o = op.Op{ID: opID, Type: op.TypeNodeMove, Node: n.ID, Base: n.Version, Parent: parent, Key: k}
	}
	e := edit{ops: []op.Op{o}}
	e.touchNode(n.ID, n.Handle)
	return e, nil
}

// buildRetire adds a node.retire against the node's version. A node created
// in this change set is dropped instead: its create, every other op on it,
// and the links created here that end at it.
func buildRetire(b writeBase, r RetireRequest) (edit, error) {
	if strings.TrimSpace(r.Reason) == "" {
		return edit{}, Usagef("retire needs a --reason").WithNext("rhiza describe retire")
	}
	n, err := writableNode(b, r.Ref)
	if err != nil {
		return edit{}, err
	}
	var e edit
	e.touchNode(n.ID, n.Handle)
	if _, ok := b.pendingOp(n.Pending, op.TypeNodeCreate); ok {
		e.drop = append(e.drop, n.Pending...)
		for _, o := range b.ops {
			if o.Type == op.TypeLinkCreate && (o.From == n.ID || o.To == n.ID) {
				e.touchLink(o.Link)
			}
		}
		for _, o := range b.ops {
			if slices.Contains(e.links, o.Link) {
				e.drop = append(e.drop, o.ID)
			}
		}
		return e, nil
	}
	if n.Version.IsZero() {
		return edit{}, Errorf("node %s has no version to retire", n.ID)
	}
	opID, err := b.newID()
	if err != nil {
		return edit{}, err
	}
	e.ops = []op.Op{{ID: opID, Type: op.TypeNodeRetire, Node: n.ID, Base: n.Version, Reason: r.Reason}}
	return e, nil
}

// writableNode resolves ref to a node a write may change: not retired and
// without conflicting heads.
func writableNode(b writeBase, ref string) (Node, error) {
	n, err := b.m.Resolve(ref)
	if err != nil {
		return Node{}, err
	}
	if n.Retired {
		return Node{}, Conflictf("node %s (%s) is retired", n.ID, n.Handle)
	}
	if n.Conflict {
		return Node{}, Conflictf("node %s (%s) has conflicting changes", n.ID, n.Handle).WithNext("rhiza check")
	}
	return n, nil
}

// handleFree reports a conflict if a live node other than self holds h.
func handleFree(b writeBase, h string, self id.ID) error {
	if holder, ok := b.m.HandleTaken(h); ok && holder != self {
		return Conflictf("handle %q is taken by node %s", h, holder)
	}
	return nil
}

// treeParent resolves ref to a live node of kind, a parent for kind.
func treeParent(b writeBase, kind op.Kind, ref string) (Node, error) {
	p, err := b.m.Resolve(ref)
	if err != nil {
		return Node{}, err
	}
	if p.Retired {
		return Node{}, Usagef("parent %s (%s) is retired", p.ID, p.Handle)
	}
	if p.Kind != kind {
		return Node{}, Usagef("parent %s is %s %s; %s %s's parent must be %s %s", ref, article(string(p.Kind)), p.Kind, article(string(kind)), kind, article(string(kind)), kind)
	}
	return p, nil
}

// under reports whether node x is n or one of n's descendants.
func under(m Model, x, n id.ID) bool {
	seen := map[id.ID]bool{}
	for !x.IsZero() && !seen[x] {
		if x == n {
			return true
		}
		seen[x] = true
		p, ok := m.Node(x)
		if !ok {
			return false
		}
		x = p.Parent
	}
	return false
}

// place returns a sibling key for a node of kind under parent: just before
// the sibling before names, just after the one after names, else after the
// last sibling. self, the node being moved (zero for a new one), is not its
// own sibling.
func place(b writeBase, kind op.Kind, parent, self id.ID, before, after string) (string, error) {
	if before != "" && after != "" {
		return "", Usagef("give --before or --after, not both")
	}
	sibs := slices.DeleteFunc(b.m.Children(kind, parent), func(n Node) bool { return n.ID == self })
	var lo, hi string
	switch {
	case before != "" || after != "":
		ref, flag := before, "before"
		if after != "" {
			ref, flag = after, "after"
		}
		x, err := b.m.Resolve(ref)
		if err != nil {
			return "", err
		}
		i := slices.IndexFunc(sibs, func(n Node) bool { return n.ID == x.ID })
		if i < 0 {
			return "", Usagef("--%s %s: not a sibling of the node at its new place", flag, ref)
		}
		if before != "" {
			hi = sibs[i].Key
			if i > 0 {
				lo = sibs[i-1].Key
			}
		} else {
			lo = sibs[i].Key
			if i+1 < len(sibs) {
				hi = sibs[i+1].Key
			}
		}
	case len(sibs) > 0:
		lo = sibs[len(sibs)-1].Key
	}
	k, err := key.Between(lo, hi)
	if err != nil {
		return "", Errorf("no key between siblings %q and %q: %v", lo, hi, err)
	}
	return k, nil
}

func joinKinds() string {
	s := make([]string, 0, len(op.Kinds()))
	for _, k := range op.Kinds() {
		s = append(s, string(k))
	}
	return strings.Join(s, ", ")
}

// knownFieldsOnly refuses a field the kind does not have, as a usage error
// naming the known ones. Only new writes are checked: nodes already in the
// log with other fields still read and fold.
func knownFieldsOnly(cmd string, kind op.Kind, fields map[string]json.RawMessage) error {
	for _, name := range slices.Sorted(maps.Keys(fields)) {
		if op.FieldKnown(kind, name) {
			continue
		}
		known := "it has no fields"
		if ks := op.KnownFields(kind); len(ks) > 0 {
			known = "known fields: " + strings.Join(ks, ", ")
		}
		return Usagef("%s %s has no field %q (%s)", article(string(kind)), kind, name, known).WithNext("rhiza describe " + cmd)
	}
	return nil
}

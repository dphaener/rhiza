package api

// This file resolves the references in a node write's text fields: @handle,
// an ID or ID prefix, or @{ID} in the input become @{<full ID>} in the log,
// and @@ becomes a literal @ (see package textref and decision
// text-references-shape).

import (
	"encoding/json"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/textref"
)

// resolveFieldTexts resolves the references in each name=text of fields, in
// the view of b. current holds the node's stored fields (nil for a new node):
// a reference to a retired node is kept only when the field's current value
// already holds that ID. Strict, an unresolvable reference refuses the whole
// write; lenient, it is stored as a literal @ and the write goes on. A source
// node's fields and the verbatim fields are returned as given.
func resolveFieldTexts(b writeBase, kind op.Kind, fields []string, current map[string]json.RawMessage, lenient bool) ([]string, error) {
	if kind == op.KindSource || len(fields) == 0 {
		return fields, nil
	}
	out := make([]string, len(fields))
	var bad []string
	for i, kv := range fields {
		name, text, ok := strings.Cut(kv, "=")
		if !ok || gmodel.TextVerbatim(kind, name) {
			out[i] = kv
			continue
		}
		held := storedIDs(fieldText(current[name]))
		resolved, unresolved := textref.Resolve(text, func(r textref.Ref) (string, bool) {
			n, err := b.m.Resolve(r.Token)
			if err != nil {
				return "", false
			}
			if n.Retired && !held[n.ID.String()] {
				return "", false
			}
			return n.ID.String(), true
		})
		for _, u := range unresolved {
			bad = append(bad, "@"+u.Token+" in field "+name)
		}
		out[i] = name + "=" + resolved
	}
	if len(bad) > 0 && !lenient {
		return nil, NotFoundf("no node to refer to for: %s (a write may refer to landed nodes and nodes of its own change set, not retired ones)", strings.Join(bad, ", ")).
			WithNext("fix the handle, or write @@ for a literal @")
	}
	return out, nil
}

// storedIDs returns the full IDs stored as @{<ID>} in text.
func storedIDs(text string) map[string]bool {
	var out map[string]bool
	for _, r := range textref.Find(text) {
		if r.Braced && len(r.Token) == textref.IDLen {
			if out == nil {
				out = map[string]bool{}
			}
			out[r.Token] = true
		}
	}
	return out
}

// modelHandleOf gives the current handle of a live node of m, for showing the
// references in texts.
func modelHandleOf(m Model) gmodel.HandleOf {
	return func(i id.ID) (string, bool) {
		n, ok := m.Node(i)
		return n.Handle, ok && !n.Retired
	}
}

// showOp is o with the texts of its fields shown by current handle in m.
func showOp(m Model, o op.Op) op.Op {
	kind := o.Kind
	if n, ok := m.Node(o.Node); ok {
		kind = n.Kind
	}
	h := modelHandleOf(m)
	o.Fields = gmodel.ShowFields(kind, o.Fields, h)
	o.Set = gmodel.ShowFields(kind, o.Set, h)
	return o
}

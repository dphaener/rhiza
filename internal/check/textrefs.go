package check

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/textref"
)

// textRefs fails, under the text-ref rule, on each stored reference
// (@{<full 26-character ID>}) in a text field's current value that names no
// node in the view of the change set that last wrote that field, and warns
// on each that names a retired node. A landed change set is read in the
// as-is view, a draft, proposed or accepted one in its what-if view; other
// change sets are not read. Only the folded nodes' current fields are
// scanned, so a value a later write superseded is never read. Verbatim
// fields (an owner's words, a source) are not scanned, and any other @{...}
// is literal text.
func (c *checker) textRefs(asIs model.Model, sets []changeset.ChangeSet) {
	ops := map[id.ID]op.Op{}
	for _, s := range sets {
		for _, o := range s.Ops {
			ops[o.ID] = o
		}
	}
	for _, s := range sets {
		touched := map[id.ID]bool{}
		for _, o := range s.Ops {
			if o.Type == op.TypeNodeCreate || o.Type == op.TypeNodeUpdate {
				touched[o.Node] = true
			}
		}
		if len(touched) == 0 {
			continue
		}
		var view model.Model
		switch {
		case s.State == changeset.StateLanded:
			view = asIs
		case s.State == changeset.StateDraft || s.State == changeset.StateProposed || s.State == changeset.StateAccepted:
			view, _ = model.WhatIf(sets, s.ID)
		default:
			continue
		}
		var uses []textUse
		for node := range touched {
			n, ok := view.Node(node)
			if !ok {
				continue
			}
			for name, raw := range n.Fields {
				if model.TextVerbatim(n.Kind, name) || lastWriter(n, name, ops) != s.ID {
					continue
				}
				uses = appendTextUses(uses, node, n.Kind, name, raw)
			}
		}
		slices.SortFunc(uses, func(a, b textUse) int {
			return cmp.Or(id.Compare(a.node, b.node), strings.Compare(a.field, b.field), strings.Compare(a.ref, b.ref))
		})
		for _, u := range uses {
			name := u.node.String()
			if n, ok := view.Node(u.node); ok && n.Handle != "" {
				name = n.Handle
			}
			ref, err := id.Parse(u.ref)
			target, ok := model.Node{}, false
			if err == nil {
				target, ok = view.Node(ref)
			}
			f := Finding{Rule: RuleTextRef, ChangeSet: s.ID.String()}
			switch {
			case !ok:
				f.Severity = SeverityFail
				f.Message = fmt.Sprintf("%s field %s holds @{%s}, which names no node", name, u.field, u.ref)
			case target.Retired:
				f.Severity = SeverityWarn
				f.Message = fmt.Sprintf("%s field %s holds @{%s}, which names a retired node", name, u.field, u.ref)
			default:
				continue
			}
			c.add(f)
		}
	}
}

// lastWriter is the change set of the last operation in n's history that
// wrote field name, zero when none of ops did.
func lastWriter(n model.Node, name string, ops map[id.ID]op.Op) id.ID {
	for i := len(n.History) - 1; i >= 0; i-- {
		h := n.History[i]
		o, ok := ops[h.Op]
		if !ok {
			continue
		}
		switch o.Type {
		case op.TypeNodeCreate:
			if _, ok := o.Fields[name]; ok {
				return h.ChangeSet
			}
		case op.TypeNodeUpdate:
			if _, ok := o.Set[name]; ok {
				return h.ChangeSet
			}
		}
	}
	return id.ID{}
}

// textUse is one stored reference in a field of a node.
type textUse struct {
	node  id.ID
	kind  op.Kind
	field string
	ref   string
}

func appendTextUses(uses []textUse, node id.ID, kind op.Kind, name string, raw json.RawMessage) []textUse {
	if len(raw) == 0 || raw[0] != '"' {
		return uses
	}
	var s string
	if json.Unmarshal(raw, &s) != nil || !strings.Contains(s, "@{") {
		return uses
	}
	var seen []string
	for _, r := range textref.Find(s) {
		if r.Braced && len(r.Token) == textref.IDLen && !slices.Contains(seen, r.Token) {
			seen = append(seen, r.Token)
			uses = append(uses, textUse{node, kind, name, r.Token})
		}
	}
	return uses
}

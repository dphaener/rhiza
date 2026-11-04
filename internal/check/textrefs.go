package check

import (
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
// (@{<full 26-character ID>}) in a text field written by a change set that
// names no node in that change set's view, and warns on each that names a
// retired node. A landed change set is read in the as-is view, a draft,
// proposed or accepted one in its what-if view; other change sets are not
// read. Verbatim fields (an owner's words, a source) are not scanned, and
// any other @{...} is literal text.
func (c *checker) textRefs(asIs model.Model, sets []changeset.ChangeSet) {
	for _, s := range sets {
		var uses []textUse
		for _, o := range s.Ops {
			switch o.Type {
			case op.TypeNodeCreate:
				for name, raw := range o.Fields {
					uses = appendTextUses(uses, o.Node, o.Kind, name, raw)
				}
			case op.TypeNodeUpdate:
				for name, raw := range o.Set {
					uses = appendTextUses(uses, o.Node, "", name, raw)
				}
			}
		}
		if len(uses) == 0 {
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
		for _, u := range uses {
			kind := u.kind
			if n, ok := view.Node(u.node); ok && kind == "" {
				kind = n.Kind
			}
			if kind != "" && model.TextVerbatim(kind, u.field) {
				continue
			}
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

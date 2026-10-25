package api

// This file holds the link writes: link creates a link between two live
// nodes, unlink retires one, or drops it if this change set created it.

import (
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/op"
)

func init() {
	register(Operation{
		Name:     "link",
		Summary:  "Link two nodes",
		Write:    true,
		Request:  reflect.TypeFor[LinkRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(LinkRequest)
			return runWrite(ctx, "link", commandLine("link", r.From, r.Type, r.To),
				func(b writeBase) (edit, error) { return buildLink(b, r) })
		},
	})
	register(Operation{
		Name:     "unlink",
		Summary:  "Retire a link",
		Write:    true,
		Request:  reflect.TypeFor[UnlinkRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			r := req.(UnlinkRequest)
			return runWrite(ctx, "unlink", commandLine("unlink", r.Ref, r.Type, r.To),
				func(b writeBase) (edit, error) { return buildUnlink(b, r) })
		},
	})
}

// LinkRequest creates a link.
type LinkRequest struct {
	From     string `json:"from" arg:"0" desc:"The node the link starts at: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Type     string `json:"type" arg:"1" enum:"serves,constrains,exempts,depends-on,relates-to,decided-by,answers,blocks,supersedes,derived-from,discovered-from,implements,verifies,observes" desc:"What the link says"`
	To       string `json:"to" arg:"2" desc:"The node the link ends at"`
	Note     string `json:"note,omitempty" desc:"Why; required on relates-to and exempts"`
	Decision string `json:"decision,omitempty" desc:"The decision node behind an exempts link (exempts only, required there)"`
}

// UnlinkRequest names a link to retire: a link ref, or from, type and to.
type UnlinkRequest struct {
	Ref  string `json:"ref" arg:"0" desc:"The link (a full ID or a unique prefix of 4+ characters), or the node it starts at when type and to follow"`
	Type string `json:"type,omitempty" arg:"1" desc:"The link type, with to"`
	To   string `json:"to,omitempty" arg:"2" desc:"The node the link ends at, with type"`
}

// buildLink creates a link between two live nodes. An identical live link is
// a conflict.
func buildLink(b writeBase, r LinkRequest) (edit, error) {
	t, err := linkType(r.Type, "link")
	if err != nil {
		return edit{}, err
	}
	from, err := writableNode(b, r.From)
	if err != nil {
		return edit{}, err
	}
	to, err := writableNode(b, r.To)
	if err != nil {
		return edit{}, err
	}
	if r.Note == "" && (t == op.LinkRelatesTo || t == op.LinkExempts) {
		return edit{}, Usagef("%s %s link needs --note", article(string(t)), t).WithNext("genesis describe link")
	}
	o := op.Op{Type: op.TypeLinkCreate, LinkType: t, From: from.ID, To: to.ID, Note: r.Note}
	switch {
	case t == op.LinkExempts && r.Decision == "":
		return edit{}, Usagef("an exempts link needs --decision <decision>").WithNext("genesis describe link")
	case t != op.LinkExempts && r.Decision != "":
		return edit{}, Usagef("only an exempts link takes --decision, not %s", t)
	case r.Decision != "":
		d, err := writableNode(b, r.Decision)
		if err != nil {
			return edit{}, err
		}
		if d.Kind != op.KindDecision {
			return edit{}, Usagef("--decision %s is %s %s, not a decision", r.Decision, article(string(d.Kind)), d.Kind)
		}
		o.Decision = d.ID
	}
	if from.ID == to.ID {
		return edit{}, Usagef("a link cannot start and end at the same node %s", from.ID)
	}
	switch k, err := b.m.LinkBetween(from.ID, to.ID, t); {
	case err == nil:
		return edit{}, Conflictf("%s %s link from %s to %s already exists: %s", article(string(t)), t, r.From, r.To, k.ID)
	case AsError(err).Code != CodeNotFound:
		return edit{}, err
	}
	if o.Link, err = b.newID(); err != nil {
		return edit{}, err
	}
	if o.ID, err = b.newID(); err != nil {
		return edit{}, err
	}
	e := edit{ops: []op.Op{o}}
	e.touchLink(o.Link)
	return e, nil
}

// buildUnlink retires a link against its version, or drops it, with every
// op on it, if this change set created it.
func buildUnlink(b writeBase, r UnlinkRequest) (edit, error) {
	var k Link
	switch {
	case r.Type == "" && r.To == "":
		var err error
		if k, err = b.m.ResolveLink(r.Ref); err != nil {
			return edit{}, err
		}
	case r.Type != "" && r.To != "":
		t, err := linkType(r.Type, "unlink")
		if err != nil {
			return edit{}, err
		}
		from, err := b.m.Resolve(r.Ref)
		if err != nil {
			return edit{}, err
		}
		to, err := b.m.Resolve(r.To)
		if err != nil {
			return edit{}, err
		}
		if k, err = b.m.LinkBetween(from.ID, to.ID, t); err != nil {
			return edit{}, err
		}
	default:
		return edit{}, Usagef("unlink takes a link, or <from> <type> <to>").WithNext("genesis describe unlink")
	}
	if k.Retired {
		return edit{}, Conflictf("link %s is retired", k.ID)
	}
	if k.Conflict {
		return edit{}, Conflictf("link %s has conflicting landed changes; a change set resolving them must land first", k.ID)
	}
	var e edit
	e.touchLink(k.ID)
	if _, ok := b.pendingOp(k.Pending, op.TypeLinkCreate); ok {
		e.drop = append(e.drop, k.Pending...)
		return e, nil
	}
	if k.Version.IsZero() {
		return edit{}, Errorf("link %s has no landed version to retire", k.ID)
	}
	opID, err := b.newID()
	if err != nil {
		return edit{}, err
	}
	e.ops = []op.Op{{ID: opID, Type: op.TypeLinkRetire, Link: k.ID, Base: k.Version}}
	return e, nil
}

// linkType checks a link type name.
func linkType(s, opName string) (op.LinkType, error) {
	t := op.LinkType(s)
	if !slices.Contains(op.LinkTypes(), t) {
		names := make([]string, 0, len(op.LinkTypes()))
		for _, x := range op.LinkTypes() {
			names = append(names, string(x))
		}
		return "", Usagef("unknown link type %q: want one of %s", s, strings.Join(names, ", ")).
			WithNext("genesis describe " + opName)
	}
	return t, nil
}

package api

// This file holds the link writes: link creates a link between two live
// nodes, or a node and a change set, or two change sets, refusing a
// depends-on between change sets, or between nodes, that would close a
// loop; unlink retires one, or drops it if this change set created it.

import (
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "link",
		Summary:  "Link two nodes, two change sets, or one of each",
		Write:    true,
		Record:   true,
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
		Record:   true,
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
	From     string `json:"from" arg:"0" desc:"The node or change set the link starts at: an ID, a unique ID prefix of 4+ characters, or a live node's handle"`
	Type     string `json:"type" arg:"1" enum:"serves,constrains,exempts,depends-on,relates-to,decided-by,answers,blocks,supersedes,derived-from,discovered-from,implements,verifies,observes" desc:"What the link says"`
	To       string `json:"to" arg:"2" desc:"The node or change set the link ends at"`
	Note     string `json:"note,omitempty" desc:"Why; required on relates-to and exempts"`
	Decision string `json:"decision,omitempty" desc:"The decision node behind an exempts link (exempts only, required there)"`
}

// UnlinkRequest names a link to retire: a link ref, or from, type and to.
type UnlinkRequest struct {
	Ref  string `json:"ref" arg:"0" desc:"The link (a full ID or a unique prefix of 4+ characters), or the node it starts at when type and to follow"`
	Type string `json:"type,omitempty" arg:"1" desc:"The link type, with to"`
	To   string `json:"to,omitempty" arg:"2" desc:"The node the link ends at, with type"`
}

// buildLink creates a link between two live nodes, or change sets, or one
// of each. An identical live link is a conflict, and so is a depends-on
// between change sets, or between nodes, that would close a loop.
func buildLink(b writeBase, r LinkRequest) (edit, error) {
	t, err := linkType(r.Type, "link")
	if err != nil {
		return edit{}, err
	}
	from, fromSet, err := linkEnd(b, r.From)
	if err != nil {
		return edit{}, err
	}
	to, toSet, err := linkEnd(b, r.To)
	if err != nil {
		return edit{}, err
	}
	if r.Note == "" && (t == op.LinkRelatesTo || t == op.LinkExempts) {
		return edit{}, Usagef("%s %s link needs --note", article(string(t)), t).WithNext("rhiza describe link")
	}
	o := op.Op{Type: op.TypeLinkCreate, LinkType: t, From: from, To: to, Note: r.Note}
	switch {
	case t == op.LinkExempts && r.Decision == "":
		return edit{}, Usagef("an exempts link needs --decision <decision>").WithNext("rhiza describe link")
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
	if from == to {
		return edit{}, Usagef("a link cannot start and end at the same node %s", from)
	}
	switch k, err := b.m.LinkBetween(from, to, t); {
	case err == nil:
		return edit{}, Conflictf("%s %s link from %s to %s already exists: %s", article(string(t)), t, r.From, r.To, k.ID)
	case AsError(err).Code != CodeNotFound:
		return edit{}, err
	}
	if t == op.LinkDependsOn && fromSet && toSet && b.changeSets != nil && wouldCycle(b, from, to) {
		return edit{}, Conflictf("depends-on from %s to %s would close a loop: %s already depends on %s", from, to, to, from).
			WithNext("rhiza change show " + shortID(to))
	}
	if t == op.LinkDependsOn && !fromSet && !toSet {
		if path := dependsPath(b.m, to, from); path != nil {
			return edit{}, Conflictf("depends-on from %s to %s would close a loop: %s",
				nodeName(b.m, from), nodeName(b.m, to), loopText(b.m, path)).
				WithNext("rhiza show " + nodeName(b.m, to))
		}
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

// linkEnd resolves one end of a new link: a writable node, or else a
// change set in the log, such as either end of a discovered-from link,
// reporting which. A ref that is neither keeps the node's error.
func linkEnd(b writeBase, ref string) (id.ID, bool, error) {
	n, err := writableNode(b, ref)
	if err == nil {
		return n.ID, false, nil
	}
	if AsError(err).Code != CodeNotFound || b.changeSet == nil {
		return id.ID{}, false, err
	}
	if cs, csErr := b.changeSet(ref); csErr == nil {
		return cs, true, nil
	}
	return id.ID{}, false, err
}

// wouldCycle reports whether change set from depending on change set to
// would close a loop. The graph is every change set's dependencies in the
// log, plus each depends-on link the change set being written holds, from
// whichever change set, as it will count once that change set lands: so an
// apply whose links together close a loop is refused at the one that does.
func wouldCycle(b writeBase, from, to id.ID) bool {
	sets := b.changeSets()
	retired := map[id.ID]bool{}
	for _, o := range b.ops {
		if o.Type == op.TypeLinkRetire {
			retired[o.Link] = true
		}
	}
	// The written change set's links go in as one landed change set with
	// no ID of its own, which nothing can depend on.
	var held []op.Op
	for _, o := range b.ops {
		if o.Type == op.TypeLinkCreate && o.LinkType == op.LinkDependsOn && !retired[o.Link] {
			held = append(held, o)
		}
	}
	if len(held) > 0 {
		sets = append(sets, changeset.ChangeSet{State: changeset.StateLanded, Ops: held})
	}
	return changeset.WouldCycle(sets, from, to)
}

// dependsPath returns the shortest path from node from to node to over the
// live depends-on links of m, both ends included, choosing among equally
// short ones by ID; nil when to cannot be reached. m holds the change set's
// own links, and in an apply every earlier item's, so a loop that an apply
// closes is refused at the item that closes it.
func dependsPath(m Model, from, to id.ID) []id.ID {
	prev := map[id.ID]id.ID{from: {}}
	next := []id.ID{from}
	for len(next) > 0 {
		at := next[0]
		next = next[1:]
		for _, n := range m.DependsOn(at) {
			if _, ok := prev[n]; ok {
				continue
			}
			prev[n] = at
			if n == to {
				path := []id.ID{to}
				for p := at; p != from; p = prev[p] {
					path = append(path, p)
				}
				path = append(path, from)
				slices.Reverse(path)
				return path
			}
			next = append(next, n)
		}
	}
	return nil
}

// loopText says what path says, one step at a time: "a depends on b, b
// depends on c".
func loopText(m Model, path []id.ID) string {
	steps := make([]string, 0, len(path)-1)
	for i := range len(path) - 1 {
		steps = append(steps, nodeName(m, path[i])+" depends on "+nodeName(m, path[i+1]))
	}
	return strings.Join(steps, ", ")
}

// nodeName is node n's handle, or its ID when it has none.
func nodeName(m Model, n id.ID) string {
	if v, ok := m.Node(n); ok && v.Handle != "" {
		return v.Handle
	}
	return n.String()
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
		return edit{}, Usagef("unlink takes a link, or <from> <type> <to>").WithNext("rhiza describe unlink")
	}
	if k.Retired {
		return edit{}, Conflictf("link %s is retired", k.ID)
	}
	if k.Conflict {
		return edit{}, Conflictf("link %s has conflicting changes", k.ID).WithNext("rhiza check")
	}
	var e edit
	e.touchLink(k.ID)
	if _, ok := b.pendingOp(k.Pending, op.TypeLinkCreate); ok {
		e.drop = append(e.drop, k.Pending...)
		return e, nil
	}
	if k.Version.IsZero() {
		return edit{}, Errorf("link %s has no version to retire", k.ID)
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
			WithNext("rhiza describe " + opName)
	}
	return t, nil
}

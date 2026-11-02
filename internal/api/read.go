package api

// This file holds the read layer every read operation shares: views, the
// default view, and a reader that answers queries in one view from the
// index.
//
// The vocabulary:
//
//   - A view is the model under a choice of change sets (see package
//     model): "as-is" (the landed change sets), "to-be" (as-is plus every
//     accepted change set) or "what-if:<change set>" (as-is plus that one
//     change set, which must be draft, proposed or accepted). A read's view
//     request field takes one of these, the change set as any change set
//     ref; empty means the operation's default.
//   - The effective change set is the one a call acts on (see
//     effectiveChangeSet): Ctx.ChangeSet, else Ctx.ChangeRef (the CLI's
//     --change or RHIZA_CHANGE), else the current change set, else none.
//
// The default view (defaultView) follows the effective change set's state:
// draft or proposed gives what-if of it, so a read shows the work in hand;
// accepted gives to-be, where it already is; none, landed or any other
// terminal state gives as-is. The rule is slice.DefaultView's, which
// defaultView maps onto a view, so every read and slice agree (D4). An
// operation whose question is about the plan rather than the work in hand
// passes its own rule instead: roadmap and ready pass toBeView, to-be
// always. An explicit view always wins, and then the effective change set
// is not looked at.
//
// The as-is and to-be views are stored in the index; what-if is computed on
// demand from its stored entries (index.WhatIf), so a what-if read costs a
// fold of the whole log. The reader answers every query identically from
// either backend: the model's are written to the index queries' rules.

import (
	"cmp"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/slice"
)

// The view names.
const (
	viewAsIs       = "as-is"
	viewToBe       = "to-be"
	viewWhatIfPref = "what-if:"
)

// lsNext is the next command for a ref that names nothing.
const lsNext = "rhiza ls"

// view is a parsed view: stored is index.AsIs or index.ToBe, or "" for the
// what-if of whatIf. The zero view stands for the operation's default.
type view struct {
	stored index.View
	whatIf id.ID
}

// String is the view's name: "as-is", "to-be" or "what-if:<full ID>"; ""
// for the zero view.
func (v view) String() string {
	switch {
	case v.stored == index.AsIs:
		return viewAsIs
	case v.stored == index.ToBe:
		return viewToBe
	case !v.whatIf.IsZero():
		return viewWhatIfPref + v.whatIf.String()
	}
	return ""
}

// viewRule picks an operation's default view from the effective change set
// cs and its folded state; both are zero when there is none.
type viewRule func(cs id.ID, state changeset.State) view

// defaultView is the default view of every read but roadmap and ready:
// slice.DefaultView's rule (what-if of a draft or proposed change set,
// to-be for an accepted one, as-is otherwise), as a view.
func defaultView(cs id.ID, state changeset.State) view {
	var c *changeset.ChangeSet
	if !cs.IsZero() {
		c = &changeset.ChangeSet{ID: cs, State: state}
	}
	switch slice.DefaultView(c) {
	case slice.ViewWhatIf:
		return view{whatIf: cs}
	case slice.ViewToBe:
		return view{stored: index.ToBe}
	}
	return view{stored: index.AsIs}
}

// toBeView is the default view of roadmap and ready, whose question is the
// plan: to-be, whatever the effective change set.
func toBeView(id.ID, changeset.State) view {
	return view{stored: index.ToBe}
}

// parseView parses the view field of operation opName: "" (the zero view),
// "as-is", "to-be" or "what-if:<ref>", the ref resolved over w's change
// sets (not_found or usage as resolveChangeSet says). Anything else is a
// usage error.
func parseView(w *workspace, opName, s string) (view, error) {
	return parseViewWith(w.resolveChangeSet, opName, s)
}

// parseViewWith is parseView with resolve for a what-if's change set ref.
func parseViewWith(resolve func(string) (id.ID, error), opName, s string) (view, error) {
	switch s = strings.TrimSpace(s); {
	case s == "":
		return view{}, nil
	case s == viewAsIs:
		return view{stored: index.AsIs}, nil
	case s == viewToBe:
		return view{stored: index.ToBe}, nil
	case strings.HasPrefix(s, viewWhatIfPref):
		ref := strings.TrimPrefix(s, viewWhatIfPref)
		if strings.TrimSpace(ref) == "" {
			return view{}, Usagef("view %q names no change set: give what-if:<change set>", s).
				WithNext("rhiza describe " + opName)
		}
		cs, err := resolve(ref)
		if err != nil {
			return view{}, err
		}
		return view{whatIf: cs}, nil
	}
	return view{}, Usagef("unknown view %q: give as-is, to-be or what-if:<change set>", s).
		WithNext("rhiza describe " + opName)
}

// reader answers the queries of every read in one view: from the index for
// as-is and to-be, from the computed model for what-if. Errors are already
// Errors: an unknown ref not_found with next "rhiza ls", an ambiguous one
// usage.
type reader struct {
	x    *index.Index
	view view
	// ov is the what-if overlay on the as-is rows, nil for a stored view.
	ov *index.Overlay
}

// openReader opens a reader for operation opName in the view named by s,
// or by rule (defaultView, or the operation's own) when s is empty. A
// what-if of a change set not in the log is not_found, of one not in
// flight a conflict, both with next "rhiza change ls". An empty
// ctx.IndexPath is an error.
func openReader(ctx *Ctx, opName, s string, rule viewRule) (*reader, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	return openReaderIn(ctx, w, opName, s, rule)
}

// openReaderIn is openReader over the open workspace w, whose pending
// entries a no-view read sees.
func openReaderIn(ctx *Ctx, w *workspace, opName, s string, rule viewRule) (*reader, error) {
	var err error
	// The index answers a what-if ref without reading the log; it is opened
	// once, and only when the view needs it or the read goes on.
	var x *index.Index
	open := func() (*index.Index, error) {
		if x != nil {
			return x, nil
		}
		if ctx.IndexPath == "" {
			return nil, Errorf("no index path")
		}
		xi, err := index.Open(ctx.Root, ctx.IndexPath)
		if err != nil {
			return nil, Errorf("%s", err.Error())
		}
		x = xi
		return x, nil
	}
	closeX := func() {
		if x != nil {
			x.Close()
			x = nil
		}
	}
	resolve := w.resolveChangeSet
	if len(w.pending) == 0 && ctx.IndexPath != "" {
		resolve = func(ref string) (id.ID, error) {
			xi, err := open()
			if err != nil {
				return id.ID{}, err
			}
			infos, err := xi.ChangeSets()
			if err != nil {
				return id.ID{}, Errorf("%s", err.Error())
			}
			sets := make([]id.ID, len(infos))
			for i, c := range infos {
				sets[i] = c.ID
			}
			return matchChangeSet(sets, ref)
		}
	}
	v, err := parseViewWith(resolve, opName, s)
	if err != nil {
		closeX()
		return nil, err
	}
	if v == (view{}) {
		var cs id.ID
		var state changeset.State
		if len(w.pending) == 0 && ctx.IndexPath != "" {
			cs, state, err = indexedChangeSet(ctx, w, open)
		} else {
			cs, state, err = foldedChangeSet(ctx, w)
		}
		if err != nil {
			closeX()
			return nil, err
		}
		v = rule(cs, state)
	}
	x, err = open()
	if err != nil {
		return nil, err
	}
	r := &reader{x: x, view: v}
	if v.whatIf.IsZero() {
		return r, nil
	}
	ov, ps, err := x.WhatIfOverlay(v.whatIf)
	if err == nil {
		err = whatIfViewError(v.whatIf, ps)
	}
	if err != nil {
		x.Close()
		return nil, AsError(err)
	}
	r.ov = ov
	return r, nil
}

// foldedChangeSet is the effective change set and its state, folded from
// the log: the path a workspace with pending writes takes.
func foldedChangeSet(ctx *Ctx, w *workspace) (id.ID, changeset.State, error) {
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return id.ID{}, "", err
	}
	if cs.IsZero() {
		return cs, "", nil
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return id.ID{}, "", nil
	}
	return cs, c.State, nil
}

// indexedChangeSet is foldedChangeSet over the index's change set table: it
// folds nothing, and answers the same change set and state when nothing is
// pending.
func indexedChangeSet(ctx *Ctx, w *workspace, open func() (*index.Index, error)) (id.ID, changeset.State, error) {
	ref := strings.TrimSpace(ctx.ChangeRef)
	if !ctx.ChangeSet.IsZero() && ref != "" {
		return id.ID{}, "", Usagef("set either ChangeSet or ChangeRef, not both")
	}
	cs := ctx.ChangeSet
	if cs.IsZero() && ref == "" {
		cur, err := w.current()
		if err != nil || cur.IsZero() {
			return id.ID{}, "", err
		}
		cs = cur
	}
	xi, err := open()
	if err != nil {
		return id.ID{}, "", err
	}
	infos, err := xi.ChangeSets()
	if err != nil {
		return id.ID{}, "", Errorf("%s", err.Error())
	}
	switch {
	case !ctx.ChangeSet.IsZero():
		if !slices.ContainsFunc(infos, func(c index.ChangeSetInfo) bool { return c.ID == cs }) {
			return id.ID{}, "", NotFoundf("no change set %s", cs).WithNext(changeLsNext)
		}
	case ref != "":
		sets := make([]id.ID, len(infos))
		for i, c := range infos {
			sets[i] = c.ID
		}
		if cs, err = matchChangeSet(sets, ref); err != nil {
			return id.ID{}, "", err
		}
	}
	for _, c := range infos {
		if c.ID == cs {
			return cs, c.State, nil
		}
	}
	return id.ID{}, "", nil
}

// whatIfViewError is the Error for the what-if problems ps of change set cs:
// ErrMissing on cs itself is not_found, ErrState a conflict. Any other
// problem is the change set's own (a conflict, a skipped op) and no error:
// the view shows what it would make.
func whatIfViewError(cs id.ID, ps []gmodel.Problem) error {
	for _, p := range ps {
		if p.ChangeSet != cs || !p.Op.IsZero() {
			continue
		}
		switch {
		case errors.Is(p.Err, gmodel.ErrMissing):
			return NotFoundf("no change set %s", cs).WithNext(changeLsNext)
		case errors.Is(p.Err, gmodel.ErrState):
			return Conflictf("%s, so it has no what-if view",
				strings.TrimPrefix(p.Err.Error(), gmodel.ErrState.Error()+": ")).WithNext(changeLsNext)
		}
	}
	return nil
}

// Close closes the index.
func (r *reader) Close() error {
	return r.x.Close()
}

// View is the resolved view's name, such as "to-be" or "what-if:<ID>".
func (r *reader) View() string {
	return r.view.String()
}

// readError maps an index or model query error to an Error: ErrNotFound is
// not_found with next "rhiza ls", ErrAmbiguous usage, anything else
// error; the message drops the sentinel's prefix.
func readError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case errors.Is(err, index.ErrNotFound):
		return NotFoundf("%s", strings.TrimPrefix(msg, index.ErrNotFound.Error()+": ")).WithNext(lsNext)
	case errors.Is(err, index.ErrAmbiguous):
		return Usagef("%s", strings.TrimPrefix(msg, index.ErrAmbiguous.Error()+": "))
	}
	return Errorf("%s", msg)
}

// Node returns the node with ID i, live or retired.
func (r *reader) Node(i id.ID) (gmodel.Node, bool, error) {
	if r.ov != nil {
		n, ok, err := r.ov.Node(i)
		return n, ok, readError(err)
	}
	n, ok, err := r.x.Node(r.view.stored, i)
	return n, ok, readError(err)
}

// Link returns the link with ID i, live or retired.
func (r *reader) Link(i id.ID) (gmodel.Link, bool, error) {
	if r.ov != nil {
		l, ok, err := r.ov.Link(i)
		return l, ok, readError(err)
	}
	l, ok, err := r.x.Link(r.view.stored, i)
	return l, ok, readError(err)
}

// Resolve finds a node by ref as index.Index.Resolve does: a full ID, a
// unique ID prefix of 4+ characters as given, or a live node's handle.
func (r *reader) Resolve(ref string) (gmodel.Node, error) {
	if r.ov != nil {
		n, err := r.ov.Resolve(ref)
		return n, readError(err)
	}
	n, err := r.x.Resolve(r.view.stored, ref)
	return n, readError(err)
}

// ResolveLink finds a link by ref as index.Index.ResolveLink does: a full
// ID or a unique ID prefix of 4+ characters, live or retired.
func (r *reader) ResolveLink(ref string) (gmodel.Link, error) {
	if r.ov != nil {
		l, err := r.ov.ResolveLink(ref)
		return l, readError(err)
	}
	l, err := r.x.ResolveLink(r.view.stored, ref)
	return l, readError(err)
}

// Children returns the nodes directly under parent, live and retired,
// sorted by key then ID; the zero parent gives the roots of the tree.
func (r *reader) Children(parent id.ID) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.Children(parent)
		return ns, readError(err)
	}
	ns, err := r.x.Children(r.view.stored, parent)
	return ns, readError(err)
}

// maxDepth bounds every walk of the tree, as the index bounds its own.
const maxDepth = 1000

// Ancestors returns node i's parent chain, the root first, without i; nil
// for a root or a node outside the tree. An unknown i is not_found.
func (r *reader) Ancestors(i id.ID) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.Ancestors(i)
		return ns, readError(err)
	}
	ns, err := r.x.Ancestors(r.view.stored, i)
	return ns, readError(err)
}

// Links returns the links at node in direction dir, live and retired,
// sorted by ID; node may be a change set.
func (r *reader) Links(node id.ID, dir index.Direction) ([]gmodel.Link, error) {
	if r.ov != nil {
		ls, err := r.ov.Links(node, dir)
		return ls, readError(err)
	}
	ls, err := r.x.Links(r.view.stored, node, dir)
	return ls, readError(err)
}

// Subtree returns the nodes under root, pre-order, as index.Index.Subtree
// does.
func (r *reader) Subtree(root id.ID, depth int, retired bool) ([]index.TreeNode, error) {
	if r.ov != nil {
		ts, err := r.ov.Subtree(root, depth, retired)
		return ts, readError(err)
	}
	ts, err := r.x.Subtree(r.view.stored, root, depth, retired)
	return ts, readError(err)
}

// SubtreeLight is Subtree with light nodes, as index.Index.SubtreeLight
// gives them: Fields and History nil.
func (r *reader) SubtreeLight(root id.ID, depth int, retired bool) ([]index.TreeNode, error) {
	if r.ov != nil {
		ts, err := r.ov.SubtreeLight(root, depth, retired)
		return ts, readError(err)
	}
	ts, err := r.x.SubtreeLight(r.view.stored, root, depth, retired)
	return ts, readError(err)
}

// NodesLight returns the nodes with the given IDs, live or retired, sorted
// by ID, each once, as index.Index.NodesLight does.
func (r *reader) NodesLight(ids []id.ID) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.NodesLight(ids)
		return ns, readError(err)
	}
	ns, err := r.x.NodesLight(r.view.stored, ids)
	return ns, readError(err)
}

// NodesFields is NodesLight with each node's Fields, as
// index.Index.NodesFields gives them: History nil.
func (r *reader) NodesFields(ids []id.ID) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.NodesFields(ids)
		return ns, readError(err)
	}
	ns, err := r.x.NodesFields(r.view.stored, ids)
	return ns, readError(err)
}

// LinksAt returns the links at any of ids in direction dir, live and
// retired, sorted by ID, each once, as index.Index.LinksAt does.
func (r *reader) LinksAt(ids []id.ID, dir index.Direction) ([]gmodel.Link, error) {
	if r.ov != nil {
		ls, err := r.ov.LinksAt(ids, dir)
		return ls, readError(err)
	}
	ls, err := r.x.LinksAt(r.view.stored, ids, dir)
	return ls, readError(err)
}

// LinksTo returns the links of type t whose To end is one of ids, live and
// retired, sorted by ID, as index.Index.LinksTo does.
func (r *reader) LinksTo(ids []id.ID, t op.LinkType) ([]gmodel.Link, error) {
	if r.ov != nil {
		ls, err := r.ov.LinksTo(ids, t)
		return ls, readError(err)
	}
	ls, err := r.x.LinksTo(r.view.stored, ids, t)
	return ls, readError(err)
}

// List returns the nodes q selects, sorted by ID, at most q.Limit, as
// index.Index.List does; a Limit of zero or less is an error.
func (r *reader) List(q index.ListQuery) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.List(q)
		return ns, readError(err)
	}
	ns, err := r.x.List(r.view.stored, q)
	return ns, readError(err)
}

// lightNode is n without its Fields and History.
func lightNode(n gmodel.Node) gmodel.Node {
	n.Fields, n.History = nil, nil
	return n
}

// ChangeSets returns every folded change set, sorted by ID, from the index
// whatever the view.
func (r *reader) ChangeSets() ([]index.ChangeSetInfo, error) {
	cs, err := r.x.ChangeSets()
	return cs, readError(err)
}

// Entries returns the good entries of change set cs, sorted by ID, from
// the index whatever the view.
func (r *reader) Entries(cs id.ID) ([]entry.Entry, error) {
	es, err := r.x.Entries(cs)
	return es, readError(err)
}

// EntriesOf returns the good entries of each of the change sets css,
// sorted by ID, keyed by change set (none for one without entries), from
// the index whatever the view.
func (r *reader) EntriesOf(css []id.ID) (map[id.ID][]entry.Entry, error) {
	es, err := r.x.EntriesOf(css)
	return es, readError(err)
}

// shortIDLen is how many leading characters of an ID displayName shows:
// the 10 of the ULID's millisecond time plus 4 of its random part (20
// bits), enough to tell apart nodes made in the same millisecond in all but
// rare cases, and long enough to be a ref prefix. It is for people to read;
// output that a program reads carries the full ID.
const shortIDLen = 14

// shortID is the first shortIDLen characters of i.
func shortID(i id.ID) string {
	s := i.String()
	return s[:min(len(s), shortIDLen)]
}

// displayName names a node for people: "<handle> (<short ID>)".
func displayName(n gmodel.Node) string {
	return n.Handle + " (" + shortID(n.ID) + ")"
}

// linkBackNames reads each link type from its To end: A serves B is B
// served by A.
var linkBackNames = map[op.LinkType]string{
	op.LinkServes:         "served by",
	op.LinkConstrains:     "constrained by",
	op.LinkExempts:        "exempted by",
	op.LinkDependsOn:      "depended on by",
	op.LinkRelatesTo:      "related to",
	op.LinkDecidedBy:      "decides",
	op.LinkAnswers:        "answered by",
	op.LinkBlocks:         "blocked by",
	op.LinkSupersedes:     "superseded by",
	op.LinkDerivedFrom:    "source of",
	op.LinkDiscoveredFrom: "led to",
	op.LinkImplements:     "implemented by",
	op.LinkVerifies:       "verified by",
	op.LinkObserves:       "observed by",
	op.LinkDiscusses:      "discussed by",
}

// backName is link type t read backwards, for a link shown at its To end;
// an unknown type is "<t> (from)".
func backName(t op.LinkType) string {
	return cmp.Or(linkBackNames[t], string(t)+" (from)")
}

// fieldText is a field's value as text: a JSON string is its text, any
// other JSON its canonical form without the final newline (raw as given if
// it does not parse).
func fieldText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	body, err := canon.Canonicalize(raw)
	if err != nil {
		return string(raw)
	}
	return strings.TrimRight(string(body), "\n")
}

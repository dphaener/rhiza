package api

// This file holds the slice operation: the composed context for one node
// or for a change set, read through the SQLite index and fitted to a token
// budget (see internal/slice).
//
// A fresh read stays fast by never reading the log through the store: the
// workspace is opened only to check the repository and read the current
// change set file, and everything else comes from the index, whose open
// catches up with the log from the folder listing. Change set refs resolve
// over index.ChangeSets() by the workspace's rules (matchChangeSet), and
// the change sets are folded from index.Entries. Only what-if reads the
// whole log, refolded by index.WhatIf.
//
// The response is the slice's JSON document field for field, so --json
// prints exactly the bytes slice.Fit measured; Human prints the markdown
// Fit measured for its own budget, since each format is fitted on its own.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/slice"
)

// defaultSliceBudget is the token budget when the request gives none.
const defaultSliceBudget = 4000

// SliceBudget is a token budget. Zero means none given, so the default; a
// value given must be at least 1, else the request does not decode (a
// usage error), so `--budget 0` is refused rather than taken as the
// default.
type SliceBudget int

// UnmarshalJSON decodes an integer of at least 1.
func (b *SliceBudget) UnmarshalJSON(data []byte) error {
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("budget must be at least 1, not %d", n)
	}
	*b = SliceBudget(n)
	return nil
}

// SliceRequest asks for a node slice (Ref set) or a change slice (Ref
// empty: the effective change set's).
type SliceRequest struct {
	Ref    string      `json:"ref,omitempty" arg:"0" desc:"The node: an ID, a unique ID prefix of 4+ characters, or a live node's handle; none slices the change set (--change, else the current one)"`
	View   string      `json:"view,omitempty" enum:"as-is,to-be,what-if" desc:"The view to read; default what-if for a draft or proposed change set, to-be for an accepted one, else as-is"`
	Budget SliceBudget `json:"budget,omitempty" desc:"The token budget (a token is 4 bytes of output), at least 1; default 4000"`
}

// SliceResponse is a slice as printed: the fields of slice.Doc, so its
// canonical JSON is the document slice.Fit fitted. markdown is the
// markdown Fit fitted, for Human.
type SliceResponse struct {
	View        string            `json:"view" desc:"The view sliced: as-is, to-be or what-if"`
	Change      *slice.DocChange  `json:"change,omitempty" desc:"The change set sliced or viewed through; absent for none"`
	Ref         string            `json:"ref" desc:"The focus handle of a node slice, or the change set's short ID of a change slice"`
	Budget      int               `json:"budget" desc:"The token budget"`
	Tokens      int               `json:"tokens" desc:"The estimated tokens of this whole document: its bytes divided by 4, rounded up"`
	OverBudget  bool              `json:"over_budget" desc:"True when the slice is over budget even with every section that can be cut cut to a count"`
	Focus       []slice.DocFocus  `json:"focus" desc:"The focus nodes, never cut"`
	Context     []slice.DocItem   `json:"context" desc:"The ancestors, the root first, then the intents served"`
	Constraints []slice.DocItem   `json:"constraints" desc:"The constraints that apply, never cut"`
	Decisions   []slice.DocItem   `json:"decisions" desc:"The current decisions that apply"`
	Questions   []slice.DocItem   `json:"questions" desc:"The open questions, blocking ones first, never cut"`
	Anchors     []slice.DocItem   `json:"anchors" desc:"The anchors that implement the focus"`
	Evidence    []slice.DocItem   `json:"evidence" desc:"The evidence that verifies the focus"`
	Children    []slice.DocItem   `json:"children" desc:"The live children of the focus"`
	Neighbours  []slice.DocItem   `json:"neighbours" desc:"Nodes the focus needs, is needed by or relates to"`
	InFlight    []slice.DocChange `json:"in_flight" desc:"Other in-flight change sets that target a focus node"`
	Omitted     []slice.Omission  `json:"omitted" desc:"Each section cut to fit the budget, in section order"`

	markdown []byte
}

func init() {
	register(Operation{
		Name:     "slice",
		Summary:  "Print the composed context for a node or a change set",
		Request:  reflect.TypeFor[SliceRequest](),
		Response: reflect.TypeFor[SliceResponse](),
		Run:      runSlice,
	})
}

// sliceNext is the next command for a change slice with no change set.
const sliceNext = "rhiza slice <ref>"

// openIndex opens the index of the worktree at root, caught up with its
// log. No repository is not_found with next "rhiza init". The caller
// closes it.
func openIndex(root string) (*index.Index, error) {
	if _, err := openStore(root); err != nil {
		return nil, err
	}
	p, err := index.DefaultPath(root)
	if err != nil {
		return nil, Errorf("%v", err)
	}
	x, err := index.Open(root, p)
	if err != nil {
		return nil, Errorf("%v", err)
	}
	return x, nil
}

// runSlice picks the change set and the view, resolves the ref, builds
// the slice and fits it to the budget in both formats.
func runSlice(ctx *Ctx, req any) (any, error) {
	r := req.(SliceRequest)
	budget := int(r.Budget)
	switch {
	case budget == 0:
		budget = defaultSliceBudget
	case budget < 0:
		return nil, Usagef("budget must be at least 1, not %d", budget).WithNext("rhiza describe slice")
	}
	if r.View != "" && !slices.Contains([]string{slice.ViewAsIs, slice.ViewToBe, slice.ViewWhatIf}, r.View) {
		return nil, Usagef("unknown view %q: want as-is, to-be or what-if", r.View).WithNext("rhiza describe slice")
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	x, err := openIndex(ctx.Root)
	if err != nil {
		return nil, err
	}
	defer x.Close()

	infos, err := x.ChangeSets()
	if err != nil {
		return nil, Errorf("read change sets: %v", err)
	}
	all := make([]id.ID, len(infos))
	for i, c := range infos {
		all[i] = c.ID
	}
	csID, err := sliceChangeSet(ctx, w, all)
	if err != nil {
		return nil, err
	}
	var change *changeset.ChangeSet
	if !csID.IsZero() {
		cs, err := foldIndexed(x, csID)
		if err != nil {
			return nil, err
		}
		change = &cs
	}
	if r.Ref == "" && change == nil {
		return nil, Usagef("no change set to slice: name a node, or give --change or make one current").WithNext(sliceNext)
	}
	var inFlight []slice.InFlight
	for _, c := range infos {
		switch c.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			cs, err := foldIndexed(x, c.ID)
			if err != nil {
				return nil, err
			}
			inFlight = append(inFlight, slice.InFlightOf(cs))
		}
	}

	def := slice.DefaultView(change)
	view := r.View
	if view == "" {
		view = def
	}
	asIs := slice.FromIndex(x, index.AsIs)
	var rd slice.Reader
	var resolve func(ref string) (gmodel.Node, error)
	switch view {
	case slice.ViewAsIs:
		rd = asIs
		resolve = func(ref string) (gmodel.Node, error) { return x.Resolve(index.AsIs, ref) }
	case slice.ViewToBe:
		rd = slice.FromIndex(x, index.ToBe)
		resolve = func(ref string) (gmodel.Node, error) { return x.Resolve(index.ToBe, ref) }
	case slice.ViewWhatIf:
		if change == nil || !inFlightState(change.State) {
			return nil, whatIfError(change)
		}
		m, _, err := x.WhatIf(change.ID)
		if err != nil {
			return nil, Errorf("what-if of %s: %v", change.ID, err)
		}
		rd = slice.FromModel(m)
		resolve = func(ref string) (gmodel.Node, error) { return resolveInModel(m, ref) }
	}

	sreq := slice.Request{Change: change, View: view, InFlight: inFlight, ChangeSets: all}
	if r.Ref != "" {
		n, err := resolve(r.Ref)
		if err != nil {
			return nil, refError(err, r.Ref, view)
		}
		sreq.Focus = n.ID
	}
	s, err := slice.Build(rd, asIs, sreq)
	if err != nil {
		return nil, refError(err, r.Ref, view)
	}
	opts := slice.Options{Budget: budget, DefaultView: def}
	md, err := slice.Fit(s, opts)
	if err != nil {
		return nil, Errorf("slice: %v", err)
	}
	opts.JSON = true
	js, err := slice.Fit(s, opts)
	if err != nil {
		return nil, Errorf("slice: %v", err)
	}
	return sliceResponse(*js.Doc, md.Bytes), nil
}

// sliceChangeSet is the effective change set, as
// workspace.effectiveChangeSet finds it but over the change sets in the
// index (sets, sorted): ctx.ChangeSet, else ctx.ChangeRef, else the
// current change set, else zero.
func sliceChangeSet(ctx *Ctx, w *workspace, sets []id.ID) (id.ID, error) {
	ref := strings.TrimSpace(ctx.ChangeRef)
	if cs := ctx.ChangeSet; !cs.IsZero() {
		if ref != "" {
			return id.ID{}, Usagef("set either ChangeSet or ChangeRef, not both")
		}
		if _, ok := slices.BinarySearchFunc(sets, cs, id.Compare); !ok {
			return id.ID{}, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
		}
		return cs, nil
	}
	if ref != "" {
		return matchChangeSet(sets, ref)
	}
	return w.current()
}

// foldIndexed folds change set cs from its entries in the index.
func foldIndexed(x *index.Index, cs id.ID) (changeset.ChangeSet, error) {
	es, err := x.Entries(cs)
	if err != nil {
		return changeset.ChangeSet{}, Errorf("read change set %s: %v", cs, err)
	}
	c, _ := changeset.Fold(cs, es)
	return c, nil
}

// inFlightState reports whether a change set in state s has a what-if
// view: draft, proposed or accepted.
func inFlightState(s changeset.State) bool {
	return s == changeset.StateDraft || s == changeset.StateProposed || s == changeset.StateAccepted
}

// whatIfError is the usage error for --view what-if without an in-flight
// change set.
func whatIfError(change *changeset.ChangeSet) error {
	if change == nil {
		return Usagef("--view what-if needs a change set: give --change or make one current").WithNext(changeLsNext)
	}
	return Usagef("--view what-if needs a draft, proposed or accepted change set; %s is %s", change.ID, change.State).
		WithNext(changeLsNext)
}

// refError maps a resolve or build error: no such node, or a retired one,
// is not_found (the view has no live node by that ref to slice); several
// matches are a usage error; anything else is an error.
func refError(err error, ref, view string) error {
	switch {
	case errors.Is(err, index.ErrAmbiguous):
		return Usagef("%s", strings.TrimPrefix(err.Error(), index.ErrAmbiguous.Error()+": ")).
			WithNext("rhiza slice <id>")
	case errors.Is(err, index.ErrNotFound), errors.Is(err, slice.ErrNotFound):
		return NotFoundf("no node matches %q in the %s view", ref, view)
	case errors.Is(err, slice.ErrRetired):
		return NotFoundf("%s", strings.TrimPrefix(err.Error(), slice.ErrRetired.Error()+": "))
	}
	return Errorf("slice: %v", err)
}

// resolveInModel finds a node of m by ref, by index.Resolve's rules: the
// full ID of a node, live or retired; then a unique ID prefix of 4+
// characters, matched as given, live or retired; then a live node's
// handle. It fails with index.ErrNotFound or index.ErrAmbiguous as
// Resolve does. The what-if view is not stored, so its nodes are resolved
// in its model, where a node the change set creates is found too.
func resolveInModel(m gmodel.Model, ref string) (gmodel.Node, error) {
	if i, err := id.Parse(ref); err == nil {
		if n, ok := m.Node(i); ok {
			return n, nil
		}
	}
	nodes := m.Nodes()
	var matches []gmodel.Node
	if len(ref) >= minPrefix {
		for _, n := range nodes {
			if strings.HasPrefix(n.ID.String(), ref) {
				matches = append(matches, n)
			}
		}
	}
	if len(matches) == 0 {
		for _, n := range nodes {
			if !n.Retired && n.Handle == ref {
				matches = append(matches, n)
			}
		}
	}
	switch len(matches) {
	case 0:
		return gmodel.Node{}, fmt.Errorf("%w: no node matches %q", index.ErrNotFound, ref)
	case 1:
		return matches[0], nil
	}
	names := make([]string, len(matches))
	for i, n := range matches {
		names[i] = n.ID.String() + " (" + n.Handle + ")"
	}
	return gmodel.Node{}, fmt.Errorf("%w: %q matches %d nodes: %s", index.ErrAmbiguous, ref, len(matches), strings.Join(names, ", "))
}

// sliceResponse is the response for document d and markdown md.
func sliceResponse(d slice.Doc, md []byte) SliceResponse {
	return SliceResponse{
		View: d.View, Change: d.Change, Ref: d.Ref, Budget: d.Budget, Tokens: d.Tokens, OverBudget: d.OverBudget,
		Focus: d.Focus, Context: d.Context, Constraints: d.Constraints, Decisions: d.Decisions,
		Questions: d.Questions, Anchors: d.Anchors, Evidence: d.Evidence, Children: d.Children,
		Neighbours: d.Neighbours, InFlight: d.InFlight, Omitted: d.Omitted,
		markdown: md,
	}
}

// Human writes the slice's markdown.
func (r SliceResponse) Human(w io.Writer) {
	w.Write(r.markdown)
}

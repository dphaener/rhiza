package api

// This file holds the change set operations: change.new opens a change set
// and makes it current, change.use sets or clears the current one, and
// change.show and change.ls report change sets as their entries fold. None
// of them writes nodes or links.

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "change.new",
		Summary:  "Open a change set and make it current",
		Write:    true,
		Request:  reflect.TypeFor[ChangeNewRequest](),
		Response: reflect.TypeFor[ChangeNewResponse](),
		Run:      runChangeNew,
	})
	register(Operation{
		Name:     "change.use",
		Summary:  "Set or clear the current change set",
		Request:  reflect.TypeFor[ChangeUseRequest](),
		Response: reflect.TypeFor[ChangeUseResponse](),
		Run:      runChangeUse,
	})
	register(Operation{
		Name:     "change.show",
		Summary:  "Show a change set: its intent, state, operations and decision",
		Request:  reflect.TypeFor[ChangeShowRequest](),
		Response: reflect.TypeFor[ChangeShowResponse](),
		Run:      runChangeShow,
	})
	register(Operation{
		Name:     "change.ls",
		Summary:  "List every change set with its kind, state and intent",
		Request:  reflect.TypeFor[ChangeLsRequest](),
		Response: reflect.TypeFor[ChangeLsResponse](),
		Run:      runChangeLs,
	})
}

// changeLsNext is the next command when a change set must be named.
const changeLsNext = "rhiza change ls"

// ChangeNewRequest opens a change set.
type ChangeNewRequest struct {
	Intent  string   `json:"intent" arg:"0" desc:"What the change set is for, in a sentence"`
	Kind    string   `json:"kind,omitempty" enum:"intention,record" desc:"intention (reviewed before it lands) or record (lands as is); default intention"`
	Sources []string `json:"sources,omitempty" flag:"source" desc:"A source node the change set comes from: an ID, ID prefix or handle"`
	Title   string   `json:"title,omitempty" desc:"A short name for the change set, one line of at most 30 characters; default the intent's first clause"`
}

// ChangeNewResponse names the new change set and its open entry.
type ChangeNewResponse struct {
	ChangeSet id.ID          `json:"change_set" desc:"The new change set's ID"`
	Entry     id.ID          `json:"entry,omitzero" desc:"The open entry's ID; absent in a dry run"`
	Kind      changeset.Kind `json:"kind" desc:"intention or record"`
	Intent    string         `json:"intent" desc:"The change set's intent"`
	Title     string         `json:"title" desc:"The change set's title, else its intent's first clause"`
	Sources   []id.ID        `json:"sources,omitempty" desc:"The source node IDs, sorted"`
	Actor     id.ID          `json:"actor,omitzero" desc:"Who opened it; absent in a dry run before any actor is recorded"`
	DryRun    bool           `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// runChangeNew validates the request, resolves the sources and the actor
// (recording one on first use), writes the open and makes the change set
// current. ctx.ChangeSet and ctx.ChangeRef are ignored: a new change set is
// not written into another one.
func runChangeNew(ctx *Ctx, req any) (any, error) {
	r := req.(ChangeNewRequest)
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	kind := changeset.Kind(r.Kind)
	if kind == "" {
		kind = changeset.KindIntention
	}
	if !slices.Contains(changeset.Kinds(), kind) {
		return nil, Usagef("unknown change set kind %q: want intention or record", r.Kind).
			WithNext("rhiza describe change.new")
	}
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	m := w.model(current)
	var sources []id.ID
	for _, ref := range r.Sources {
		n, err := m.Resolve(ref)
		if err != nil {
			return nil, err
		}
		if n.Kind != op.KindSource {
			return nil, Usagef("%s is %s %s, not a source", ref, article(string(n.Kind)), n.Kind)
		}
		sources = append(sources, n.ID)
	}
	slices.SortFunc(sources, id.Compare)
	sources = slices.Compact(sources)
	open := changeset.Open{Kind: kind, Intent: r.Intent, Sources: sources, Title: r.Title}
	if err := open.Validate(); err != nil {
		return nil, Usagef("%s", err.Error()).WithNext("rhiza describe change.new")
	}
	actor, forID, err := w.resolveActor(current, ctx)
	if err != nil {
		return nil, err
	}
	cs, e, err := w.openChangeSet(open, actor, forID)
	if err != nil {
		return nil, err
	}
	if err := w.setCurrent(cs); err != nil {
		return nil, err
	}
	resp := ChangeNewResponse{ChangeSet: cs, Entry: e.ID, Kind: kind, Intent: r.Intent,
		Title: open.DisplayTitle(), Sources: sources, Actor: actor, DryRun: ctx.DryRun}
	if ctx.DryRun {
		resp.Entry = id.ID{}
	}
	return resp, nil
}

// Human reports the new change set.
func (r ChangeNewResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would open change set %s (%s): %s\n\n", r.ChangeSet, r.Kind, r.Intent)
		fmt.Fprintf(w, "This was a dry run: nothing was written and the current change set is unchanged.\n")
		return
	}
	fmt.Fprintf(w, "Opened change set %s (%s): %s\n\n", r.ChangeSet, r.Kind, r.Intent)
	fmt.Fprintf(w, "It is now the current change set; writes go into it.\n")
}

// ChangeUseRequest names the change set to make current, or asks to clear
// it.
type ChangeUseRequest struct {
	Ref  string `json:"ref,omitempty" arg:"0" desc:"The change set: a full ID or a unique prefix of 4+ characters"`
	None bool   `json:"none,omitempty" desc:"Clear the current change set instead"`
}

// ChangeUseResponse is the current change set after the call.
type ChangeUseResponse struct {
	ChangeSet id.ID  `json:"change_set,omitzero" desc:"The current change set; absent when cleared"`
	Intent    string `json:"intent,omitempty" desc:"Its intent"`
	DryRun    bool   `json:"dry_run,omitempty" desc:"True if the current change set was left unchanged"`
}

// runChangeUse sets or clears the current change set. A terminal change
// set cannot be current: nothing can be written to it.
func runChangeUse(ctx *Ctx, req any) (any, error) {
	r := req.(ChangeUseRequest)
	ref := strings.TrimSpace(r.Ref)
	if (ref == "") == !r.None {
		return nil, Usagef("change use takes a change set or --none, not both or neither").
			WithNext("rhiza change ls")
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	if r.None {
		if err := w.clearCurrent(); err != nil {
			return nil, err
		}
		return ChangeUseResponse{DryRun: ctx.DryRun}, nil
	}
	cs, err := w.resolveChangeSet(ref)
	if err != nil {
		return nil, err
	}
	c, _ := w.changeSet(cs)
	if c.State.Terminal() {
		return nil, Conflictf("change set %s is %s; nothing can be written to it", cs, c.State).
			WithNext(newEntryNext)
	}
	if err := w.setCurrent(cs); err != nil {
		return nil, err
	}
	return ChangeUseResponse{ChangeSet: cs, Intent: c.Open.Intent, DryRun: ctx.DryRun}, nil
}

// Human reports the current change set.
func (r ChangeUseResponse) Human(w io.Writer) {
	switch {
	case r.DryRun && r.ChangeSet.IsZero():
		fmt.Fprintln(w, "Would clear the current change set. This was a dry run: nothing changed.")
	case r.DryRun:
		fmt.Fprintf(w, "Would make %s current: %s. This was a dry run: nothing changed.\n", r.ChangeSet, r.Intent)
	case r.ChangeSet.IsZero():
		fmt.Fprintln(w, "No current change set.")
	default:
		fmt.Fprintf(w, "Current change set: %s: %s\n", r.ChangeSet, r.Intent)
	}
}

// ChangeShowRequest names the change set to show.
type ChangeShowRequest struct {
	Ref string `json:"ref,omitempty" arg:"0" desc:"The change set: a full ID or a unique prefix of 4+ characters; default --change, else the current one"`
}

// ChangeShowResponse is one change set as its entries fold.
type ChangeShowResponse struct {
	ID           id.ID              `json:"id" desc:"The change set's ID"`
	Kind         changeset.Kind     `json:"kind" desc:"intention or record; empty without a valid open"`
	Intent       string             `json:"intent" desc:"The change set's intent"`
	Title        string             `json:"title" desc:"The change set's title, else its intent's first clause"`
	Sources      []id.ID            `json:"sources" desc:"The source node IDs"`
	Serves       []id.ID            `json:"serves" desc:"The intents its live serves links point to, in its write view, by link ID; empty when it serves none yet"`
	State        changeset.State    `json:"state" desc:"The lifecycle state; empty without a valid open"`
	Current      bool               `json:"current" desc:"Whether it is the current change set"`
	Ops          []json.RawMessage  `json:"ops" desc:"The current operations, as canonical op JSON, in change set order"`
	Leaves       []ChangeLeaf       `json:"leaves" desc:"The live capabilities it adds that have no live child capability, in build order: each after every leaf it depends on, ties in tree order"`
	Decided      *DecisionInfo      `json:"decided,omitempty" desc:"The applied accept or reject, if any"`
	SupersededBy id.ID              `json:"superseded_by,omitzero" desc:"The change set that superseded it, if any"`
	Heads        []id.ID            `json:"heads" desc:"The entries a new entry would follow"`
	Applied      []id.ID            `json:"applied" desc:"Every folded entry, in causal order"`
	History      []LifecycleEvent   `json:"history" desc:"The lifecycle entries that applied, in causal order"`
	Problems     []string           `json:"problems" desc:"Everything wrong with its entries"`
	Changes      ChangeDiffResponse `json:"changes" desc:"The change set's field-level diff, as change.diff reports it"`
	Notes        []ShowNote         `json:"notes" desc:"The live notes that discuss it, oldest first"`
}

// ChangeLeaf is one leaf of a change set.
type ChangeLeaf struct {
	ID     id.ID    `json:"id" desc:"The capability's ID"`
	Handle string   `json:"handle" desc:"The capability's handle"`
	Needs  []string `json:"needs" desc:"The handles of every live capability it depends on, in or out of the change set, sorted"`
}

// LifecycleEvent is one applied lifecycle entry: propose, accept, reject,
// withdraw, supersede, land or reopen.
type LifecycleEvent struct {
	Type  entry.Type `json:"type" desc:"propose, accept, reject, withdraw, supersede, land or reopen"`
	Entry id.ID      `json:"entry" desc:"The lifecycle entry"`
	At    string     `json:"at" desc:"When it was written (RFC 3339, UTC)"`
	Actor id.ID      `json:"actor" desc:"Who wrote it"`
}

// DecisionInfo is an applied accept or reject.
type DecisionInfo struct {
	Type   entry.Type `json:"type" desc:"accept or reject"`
	Entry  id.ID      `json:"entry" desc:"The accept or reject entry"`
	Actor  id.ID      `json:"actor" desc:"Who decided"`
	For    id.ID      `json:"for" desc:"Whom they decided for"`
	Rule   string     `json:"rule" desc:"The rule the decision was made under"`
	Reason string     `json:"reason,omitempty" desc:"Why"`
}

// runChangeShow folds the change set named by the request, else the one
// effectiveChangeSet chooses from ctx.
func runChangeShow(ctx *Ctx, req any) (any, error) {
	r := req.(ChangeShowRequest)
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	var cs id.ID
	if ref := strings.TrimSpace(r.Ref); ref != "" {
		cs, err = w.resolveChangeSet(ref)
	} else {
		cs, err = w.effectiveChangeSet(ctx)
	}
	if err != nil {
		return nil, err
	}
	if cs.IsZero() {
		return nil, Usagef("no change set: name one or make one current").WithNext(changeLsNext)
	}
	es := entriesOf(w.log(), cs)
	if len(es) == 0 {
		return nil, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	c, problems := changeset.Fold(cs, es)
	m := w.model(cs)
	serves := []id.ID{}
	for _, l := range m.LinksFrom(cs, op.LinkServes) {
		serves = append(serves, l.To)
	}
	resp := ChangeShowResponse{
		ID:           cs,
		Kind:         c.Open.Kind,
		Intent:       c.Open.Intent,
		Title:        c.Open.DisplayTitle(),
		Sources:      orEmpty(c.Open.Sources),
		Serves:       serves,
		State:        c.State,
		Current:      cs == current,
		Ops:          []json.RawMessage{},
		Leaves:       changeLeaves(m, c.Ops),
		SupersededBy: c.SupersededBy,
		Heads:        orEmpty(c.Heads),
		Applied:      orEmpty(c.Applied),
		History:      []LifecycleEvent{},
		Problems:     foldProblems(problems),
		Changes:      buildChangeDiff(cs, changeset.BuildDiff(w.changeSets(), cs)),
		Notes:        changeNotes(w, cs, current),
	}
	for _, m := range c.Lifecycle {
		resp.History = append(resp.History, LifecycleEvent{Type: m.Type, Entry: m.Entry,
			At: m.At.UTC().Format(lineageAtLayout), Actor: m.Actor})
	}
	for _, o := range c.Ops {
		raw, err := op.Encode(o)
		if err != nil {
			return nil, Errorf("change set %s: op %s: %v", cs, o.ID, err)
		}
		resp.Ops = append(resp.Ops, raw)
	}
	if d := c.Decided; d != nil {
		resp.Decided = &DecisionInfo{Type: d.Type, Entry: d.Entry, Actor: d.Actor, For: d.For, Rule: d.Rule, Reason: d.Reason}
	}
	return resp, nil
}

// Human shows the change set as a heading, a field list, its operations,
// its lifecycle history and its problems.
func (r ChangeShowResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "# Change set %s\n\n", r.ID)
	fmt.Fprintf(w, "- Title: %s\n", r.Title)
	fmt.Fprintf(w, "- Intent: %s\n", r.Intent)
	fmt.Fprintf(w, "- Kind: %s\n", r.Kind)
	state := string(r.State)
	if r.Current {
		state += " (current)"
	}
	fmt.Fprintf(w, "- State: %s\n", state)
	if len(r.Sources) > 0 {
		fmt.Fprintf(w, "- Sources: %s\n", joinIDs(r.Sources))
	}
	if d := r.Decided; d != nil {
		line := fmt.Sprintf("- Decided: %s by %s under rule %s", d.Type, d.Actor, d.Rule)
		if d.Reason != "" {
			line += ": " + d.Reason
		}
		fmt.Fprintln(w, line)
	}
	if !r.SupersededBy.IsZero() {
		fmt.Fprintf(w, "- Superseded by: %s\n", r.SupersededBy)
	}
	fmt.Fprintf(w, "- Heads: %s\n", joinIDs(r.Heads))
	fmt.Fprintf(w, "- Entries: %d\n", len(r.Applied))
	fmt.Fprintf(w, "\n## Operations\n\n")
	if len(r.Ops) == 0 {
		fmt.Fprintln(w, "No operations yet.")
	}
	for _, raw := range r.Ops {
		o, err := op.Decode(raw)
		if err != nil {
			fmt.Fprintf(w, "- %s\n", raw)
			continue
		}
		line := fmt.Sprintf("- %s %s", o.Type, o.Target())
		if o.Handle != "" {
			line += " (" + o.Handle + ")"
		}
		fmt.Fprintln(w, line)
	}
	if len(r.Leaves) > 0 {
		fmt.Fprintf(w, "\n## Leaves\n\n")
	}
	for _, l := range r.Leaves {
		line := "- " + l.Handle
		if len(l.Needs) > 0 {
			line += " (needs " + strings.Join(l.Needs, ", ") + ")"
		}
		fmt.Fprintln(w, line)
	}
	fmt.Fprintf(w, "\n## History\n\n")
	if len(r.History) == 0 {
		fmt.Fprintln(w, "No lifecycle moves yet.")
	}
	for _, e := range r.History {
		fmt.Fprintf(w, "- %s %s by %s\n", e.Type, e.At, e.Actor)
	}
	writeProblems(w, r.Problems)
	fmt.Fprintf(w, "\n## Changes\n\n")
	r.Changes.Human(w)
	writeNotes(w, r.Notes)
}

// changeLeaves finds the leaves among the capabilities that ops create, as
// m shows them: live, with no live child capability. A change set's own
// view is w.model(cs): its ops on top for an open one, as-is for a landed
// one. A rejected, withdrawn or superseded change set's creates are in
// neither, so it has no leaves.
//
// The order is Kahn's over depends-on among the leaves alone, the ready
// leaf first in tree order (depth first from the roots, siblings by key,
// then ID). If a loop leaves none ready, the first left in tree order goes next. A
// leaf the walk does not reach, under a retired parent, sorts after the
// rest by ID.
func changeLeaves(m Model, ops []op.Op) []ChangeLeaf {
	rank := map[id.ID]int{}
	var walk func(parent id.ID)
	walk = func(parent id.ID) {
		for _, n := range m.Children(op.KindCapability, parent) {
			rank[n.ID] = len(rank)
			walk(n.ID)
		}
	}
	walk(id.ID{})

	var leaves []Node
	for _, o := range ops {
		if o.Type != op.TypeNodeCreate || o.Kind != op.KindCapability {
			continue
		}
		n, ok := m.Node(o.Node)
		if !ok || n.Retired || n.Kind != op.KindCapability || len(m.Children(op.KindCapability, n.ID)) > 0 {
			continue
		}
		leaves = append(leaves, n)
	}
	slices.SortFunc(leaves, func(a, b Node) int {
		ra, oka := rank[a.ID]
		rb, okb := rank[b.ID]
		switch {
		case oka && okb:
			return cmp.Compare(ra, rb)
		case oka != okb:
			if oka {
				return -1
			}
			return 1
		}
		return id.Compare(a.ID, b.ID)
	})

	in := map[id.ID]bool{}
	for _, n := range leaves {
		in[n.ID] = true
	}
	placed := map[id.ID]bool{}
	ready := func(n Node) bool {
		return !slices.ContainsFunc(m.DependsOn(n.ID), func(d id.ID) bool {
			return in[d] && !placed[d]
		})
	}
	out := []ChangeLeaf{}
	for len(leaves) > 0 {
		i := slices.IndexFunc(leaves, ready)
		if i < 0 {
			i = 0
		}
		n := leaves[i]
		placed[n.ID] = true
		leaves = slices.Delete(leaves, i, i+1)
		out = append(out, ChangeLeaf{ID: n.ID, Handle: n.Handle, Needs: leafNeeds(m, n.ID)})
	}
	return out
}

// leafNeeds is the sorted handles of the live capabilities n depends on.
func leafNeeds(m Model, n id.ID) []string {
	needs := []string{}
	for _, d := range m.DependsOn(n) {
		if t, ok := m.Node(d); ok && !t.Retired && t.Kind == op.KindCapability {
			needs = append(needs, t.Handle)
		}
	}
	slices.Sort(needs)
	return needs
}

// ChangeLsRequest takes nothing.
type ChangeLsRequest struct{}

// ChangeLsResponse lists every change set in the log.
type ChangeLsResponse struct {
	ChangeSets []ChangeSetSummary `json:"change_sets" desc:"Every change set, sorted by ID"`
	Problems   []string           `json:"problems" desc:"Everything wrong with the log and its change sets"`
}

// ChangeSetSummary is one line of change.ls.
type ChangeSetSummary struct {
	ID      id.ID           `json:"id" desc:"The change set's ID"`
	Kind    changeset.Kind  `json:"kind" desc:"intention or record; empty without a valid open"`
	State   changeset.State `json:"state" desc:"The lifecycle state; empty without a valid open"`
	Intent  string          `json:"intent" desc:"The change set's intent"`
	Title   string          `json:"title" desc:"The change set's title, else its intent's first clause"`
	Current bool            `json:"current" desc:"Whether it is the current change set"`
}

// runChangeLs folds the whole log. Log and fold problems are listed beside
// the change sets, never hidden.
func runChangeLs(ctx *Ctx, _ any) (any, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	l := w.log()
	sets, problems := changeset.FoldLog(l)
	resp := ChangeLsResponse{ChangeSets: []ChangeSetSummary{}, Problems: logProblems(l.Problems)}
	resp.Problems = append(resp.Problems, foldProblems(problems)...)
	for _, c := range sets {
		resp.ChangeSets = append(resp.ChangeSets, ChangeSetSummary{
			ID: c.ID, Kind: c.Open.Kind, State: c.State, Intent: c.Open.Intent, Title: c.Open.DisplayTitle(),
			Current: c.ID == current,
		})
	}
	return resp, nil
}

// Human lists the change sets, one per line, then any problems.
func (r ChangeLsResponse) Human(w io.Writer) {
	if len(r.ChangeSets) == 0 {
		fmt.Fprintf(w, "No change sets yet.\n\nnext: %s\n", newEntryNext)
	}
	for _, c := range r.ChangeSets {
		line := fmt.Sprintf("- %s %s %s: %s", c.ID, c.State, c.Kind, c.Intent)
		// A title the intent does not start with is news; the fallback,
		// the intent's first clause, never is.
		if !strings.HasPrefix(strings.TrimSpace(c.Intent), c.Title) {
			line = fmt.Sprintf("- %s %s %s: %s (%s)", c.ID, c.State, c.Kind, c.Title, c.Intent)
		}
		if c.Current {
			line += " (current)"
		}
		fmt.Fprintln(w, line)
	}
	writeProblems(w, r.Problems)
}

// writeProblems writes a Problems section, or nothing if there are none.
func writeProblems(w io.Writer, problems []string) {
	if len(problems) == 0 {
		return
	}
	fmt.Fprintf(w, "\n## Problems\n\n")
	for _, p := range problems {
		fmt.Fprintf(w, "- %s\n", p)
	}
}

// logProblems renders store problems as "log: <path>: <error>".
func logProblems(ps []store.Problem) []string {
	out := []string{}
	for _, p := range ps {
		out = append(out, fmt.Sprintf("log: %s: %v", p.Path, p.Err))
	}
	return out
}

// foldProblems renders fold problems as "fold: <change set>/<entry>:
// <error>", or "fold: <change set>: <error>" for the change set as a whole.
func foldProblems(ps []changeset.Problem) []string {
	out := []string{}
	for _, p := range ps {
		where := p.ChangeSet.String()
		if !p.Entry.IsZero() {
			where += "/" + p.Entry.String()
		}
		out = append(out, fmt.Sprintf("fold: %s: %v", where, p.Err))
	}
	return out
}

// orEmpty returns ids, or an empty slice for nil, so required list fields
// encode as [] rather than null.
func orEmpty(ids []id.ID) []id.ID {
	if ids == nil {
		return []id.ID{}
	}
	return ids
}

func joinIDs(ids []id.ID) string {
	if len(ids) == 0 {
		return "none"
	}
	s := make([]string, len(ids))
	for i, x := range ids {
		s[i] = x.String()
	}
	return strings.Join(s, ", ")
}

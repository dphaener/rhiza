package api

// This file holds carve: it moves nodes from one unlanded change set to
// another under the same node and link IDs. What moves is each named
// node's create, every op of the source on that node, and every link of the
// source that starts or ends at it, with every op on that link; with
// --with-children, the nodes created in the source below a named one too,
// transitively.
//
// A carve is two entries that name each other: one in the source dropping
// the moved ops, and one in the target re-creating them under new op IDs,
// with each base among them rewritten to the new ID. When an op left in the
// source, or a current op of an in-flight change set that depends on the
// source directly or through others, still rests on a moved node or link,
// the source's entry also links it depends-on the target, unless it
// already does: a dependent sees the source's dependencies, so it keeps
// seeing what moved. A link that would close a loop is refused. A source
// left with no ops is not withdrawn; the response says how to.
//
// Before anything is written, and before the actor is resolved, both
// entries are built and put in the workspace's pending entries for a
// moment: the what-if views of the target and of the source must then show
// no problem that they did not show before. That one check covers a parent,
// link end or decision missing from the target's view, a handle taken
// there, and a node or link that already exists there.

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	gmodel "github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

func init() {
	register(Operation{
		Name:     "carve",
		Summary:  "Move nodes from one change set to another, keeping their IDs",
		Write:    true,
		Request:  reflect.TypeFor[CarveRequest](),
		Response: reflect.TypeFor[CarveResponse](),
		Run: func(ctx *Ctx, req any) (any, error) {
			return runCarve(ctx, req.(CarveRequest))
		},
	})
}

// CarveRequest names the nodes to carve and the two change sets.
type CarveRequest struct {
	Nodes        []string `json:"nodes" arg:"0" desc:"The nodes to move, at least one: IDs, unique ID prefixes of 4+ characters, or handles, as the source sees them; each must be created in the source"`
	From         string   `json:"from" desc:"The change set to carve from: an ID or a unique prefix of 4+ characters"`
	To           string   `json:"to" desc:"The change set to carve into: an ID or a unique prefix of 4+ characters"`
	WithChildren bool     `json:"with_children,omitempty" desc:"Also move every node created in the source below a named one, transitively"`
}

// CarveResponse reports a carve.
type CarveResponse struct {
	Source      id.ID             `json:"source" desc:"The change set carved from"`
	Target      id.ID             `json:"target" desc:"The change set carved into"`
	SourceEntry id.ID             `json:"source_entry,omitzero" desc:"The entry written to the source, dropping the moved ops; absent in a dry run"`
	TargetEntry id.ID             `json:"target_entry,omitzero" desc:"The entry written to the target, re-creating them; absent in a dry run"`
	Nodes       []id.ID           `json:"nodes" desc:"The nodes moved, sorted"`
	Links       []id.ID           `json:"links,omitempty" desc:"The links moved, sorted"`
	Dropped     []id.ID           `json:"dropped" desc:"The op IDs dropped from the source, sorted"`
	Ops         []id.ID           `json:"ops" desc:"The new op IDs in the target, in order"`
	DependsOn   id.ID             `json:"depends_on,omitzero" desc:"The depends-on link added from the source to the target, if one was"`
	DryRun      bool              `json:"dry_run,omitempty" desc:"True if nothing was written"`
	SourceOps   []json.RawMessage `json:"source_ops,omitempty" desc:"Dry run only: the ops the source's entry would carry, as canonical op JSON"`
	TargetOps   []json.RawMessage `json:"target_ops,omitempty" desc:"Dry run only: the ops the target's entry would carry, as canonical op JSON"`
	Next        string            `json:"next,omitempty" desc:"A command that would help: withdrawing the source, when the carve left it with no ops"`

	handles map[id.ID]string
	restsIn id.ID
}

// Human reports the carve: the nodes and links moved, the dependency
// added, and the withdraw hint.
func (r CarveResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Dry run: would carve %d node(s) from %s into %s; nothing was written.\n", len(r.Nodes), r.Source, r.Target)
	} else {
		fmt.Fprintf(w, "Carved %d node(s) from %s into %s (entries %s and %s).\n",
			len(r.Nodes), r.Source, r.Target, r.SourceEntry, r.TargetEntry)
	}
	fmt.Fprintln(w)
	for _, n := range r.Nodes {
		if h := r.handles[n]; h != "" {
			fmt.Fprintf(w, "- node %s (%s)\n", n, h)
		} else {
			fmt.Fprintf(w, "- node %s\n", n)
		}
	}
	for _, l := range r.Links {
		fmt.Fprintf(w, "- link %s\n", l)
	}
	fmt.Fprintf(w, "\nMoved %d op(s) under new op IDs.\n", len(r.Ops))
	if !r.DependsOn.IsZero() {
		why := "ops left there"
		if r.restsIn != r.Source {
			why = "ops in " + shortID(r.restsIn)
		}
		fmt.Fprintf(w, "%s now depends on %s (link %s): %s rest on what moved.\n",
			shortID(r.Source), shortID(r.Target), r.DependsOn, why)
	}
	if r.Next != "" {
		fmt.Fprintf(w, "\n%s has no ops left.\nnext: %s\n", shortID(r.Source), r.Next)
	}
}

// carvePlan is a carve built but not written: what moves, both bodies and
// both entry IDs.
type carvePlan struct {
	from, to           id.ID
	nodes, links       []id.ID
	handles            map[id.ID]string
	dropped, created   []id.ID
	source, target     changeset.Ops
	srcEntry, tgtEntry id.ID
	dependsOn          id.ID
	restsIn            id.ID
}

// runCarve builds and checks the carve, resolves the actor once (building
// again if that recorded one), then writes the source's entry and the
// target's.
func runCarve(ctx *Ctx, r CarveRequest) (any, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	build := func() (carvePlan, error) {
		p, err := planCarve(w, r)
		if err != nil {
			return carvePlan{}, err
		}
		return p, checkCarve(w, p)
	}
	p, err := build()
	if err != nil {
		return nil, err
	}
	before := len(w.log().Entries)
	actor, forID, err := w.resolveActor(p.from, ctx)
	if err != nil {
		return nil, err
	}
	if len(w.log().Entries) != before {
		if p, err = build(); err != nil {
			return nil, err
		}
	}
	if _, err := w.appendEntryID(p.from, p.source, actor, forID, p.srcEntry); err != nil {
		return nil, err
	}
	if _, err := w.appendEntryID(p.to, p.target, actor, forID, p.tgtEntry); err != nil {
		return nil, err
	}

	resp := CarveResponse{Source: p.from, Target: p.to, SourceEntry: p.srcEntry, TargetEntry: p.tgtEntry,
		Nodes: slices.SortedFunc(slices.Values(p.nodes), id.Compare), Dropped: p.dropped, Ops: p.created,
		DependsOn: p.dependsOn, DryRun: ctx.DryRun, handles: p.handles, restsIn: p.restsIn}
	if len(p.links) > 0 {
		resp.Links = slices.SortedFunc(slices.Values(p.links), id.Compare)
	}
	if c, _ := w.changeSet(p.from); len(c.Ops) == 0 {
		resp.Next = "genesis withdraw --change " + shortID(p.from)
	}
	if ctx.DryRun {
		resp.SourceEntry, resp.TargetEntry = id.ID{}, id.ID{}
		for _, side := range []struct {
			ops []op.Op
			out *[]json.RawMessage
		}{{p.source.Ops, &resp.SourceOps}, {p.target.Ops, &resp.TargetOps}} {
			for _, o := range side.ops {
				raw, err := op.Encode(o)
				if err != nil {
					return nil, Usagef("%s", err.Error())
				}
				*side.out = append(*side.out, raw)
			}
		}
	}
	return resp, nil
}

// carveSide resolves ref as the source or target (role) of a carve: it
// must be in the log and not terminal.
func carveSide(w *workspace, ref, role string) (changeset.ChangeSet, error) {
	cs, err := w.resolveChangeSet(ref)
	if err != nil {
		return changeset.ChangeSet{}, err
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return changeset.ChangeSet{}, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	if c.State.Terminal() {
		return changeset.ChangeSet{}, Conflictf("cannot carve %s change set %s: it is %s", role, cs, c.State).
			WithNext(changeLsNext)
	}
	return c, nil
}

// planCarve selects what moves and builds both bodies, drawing every ID in
// a fixed order (the target's op IDs, the depends-on link's IDs, then the
// source's entry ID and the target's), so a dry run predicts the real one.
func planCarve(w *workspace, r CarveRequest) (carvePlan, error) {
	if len(r.Nodes) == 0 {
		return carvePlan{}, Usagef("carve needs at least one node").WithNext("genesis describe carve")
	}
	src, err := carveSide(w, r.From, "from")
	if err != nil {
		return carvePlan{}, err
	}
	tgt, err := carveSide(w, r.To, "into")
	if err != nil {
		return carvePlan{}, err
	}
	if src.ID == tgt.ID {
		return carvePlan{}, Usagef("--from and --to are both change set %s", src.ID).WithNext("genesis describe carve")
	}
	p := carvePlan{from: src.ID, to: tgt.ID, handles: map[id.ID]string{}}

	created := map[id.ID]op.Op{}
	for _, o := range src.Ops {
		if o.Type == op.TypeNodeCreate {
			created[o.Node] = o
		}
	}
	m := w.model(src.ID)
	moved := map[id.ID]bool{}
	add := func(n Node) {
		if !moved[n.ID] {
			moved[n.ID] = true
			p.nodes = append(p.nodes, n.ID)
			p.handles[n.ID] = n.Handle
		}
	}
	for _, ref := range r.Nodes {
		n, err := m.Resolve(ref)
		if err != nil {
			return carvePlan{}, err
		}
		if _, ok := created[n.ID]; !ok {
			return carvePlan{}, Usagef("%s (%s) is not created in change set %s; carve moves only nodes the source creates",
				ref, n.ID, src.ID).WithNext("genesis change show " + shortID(src.ID))
		}
		add(n)
	}
	for grew := r.WithChildren; grew; {
		grew = false
		for _, o := range src.Ops {
			if o.Type == op.TypeNodeCreate && moved[o.Parent] && !moved[o.Node] {
				n, _ := m.Node(o.Node)
				add(n)
				grew = true
			}
		}
	}
	for _, o := range src.Ops {
		if o.Type == op.TypeLinkCreate && (moved[o.From] || moved[o.To]) && !moved[o.Link] {
			moved[o.Link] = true
			p.links = append(p.links, o.Link)
		}
	}

	var kept, carried []op.Op
	for _, o := range src.Ops {
		if moved[o.Target()] {
			carried = append(carried, o)
		} else {
			kept = append(kept, o)
		}
	}
	rests := func(o op.Op) bool {
		return moved[o.Decision] || slices.ContainsFunc(changeset.Names(o), func(n id.ID) bool { return moved[n] })
	}
	// restsIn is the source if an op left there rests on what moves, else
	// the first in-flight dependent of the source (by ID) whose op does.
	sets := w.changeSets()
	if slices.ContainsFunc(kept, rests) {
		p.restsIn = src.ID
	} else {
		a := changeset.Analyze(sets)
		for _, x := range sets {
			if x.ID != src.ID && !x.State.Terminal() && slices.Contains(a.Closure(x.ID), src.ID) &&
				slices.ContainsFunc(x.Ops, rests) {
				p.restsIn = x.ID
				break
			}
		}
	}
	needDep := !p.restsIn.IsZero() && !slices.Contains(changeset.DependsOn(sets, src.ID), tgt.ID)
	if needDep && changeset.WouldCycle(sets, src.ID, tgt.ID) {
		why := "ops left in " + src.ID.String()
		if p.restsIn != src.ID {
			why = "ops in " + p.restsIn.String()
		}
		return carvePlan{}, Conflictf("%s rest on what moves, so %s would depend on %s, and that would close a loop: %s already depends on %s",
			why, src.ID, tgt.ID, tgt.ID, src.ID).WithNext("genesis change show " + shortID(tgt.ID))
	}

	newID := func() (id.ID, error) {
		v, err := w.ids.New()
		if err != nil {
			return id.ID{}, Errorf("%s", err.Error())
		}
		return v, nil
	}
	renamed := map[id.ID]id.ID{}
	for _, o := range carried {
		v, err := newID()
		if err != nil {
			return carvePlan{}, err
		}
		renamed[o.ID] = v
		p.dropped = append(p.dropped, o.ID)
		p.created = append(p.created, v)
		o.ID = v
		if b, ok := renamed[o.Base]; ok {
			o.Base = b
		}
		p.target.Ops = append(p.target.Ops, o)
	}
	slices.SortFunc(p.dropped, id.Compare)
	p.source.Ops = []op.Op{}
	if needDep {
		l := op.Op{Type: op.TypeLinkCreate, LinkType: op.LinkDependsOn, From: src.ID, To: tgt.ID}
		if l.Link, err = newID(); err != nil {
			return carvePlan{}, err
		}
		if l.ID, err = newID(); err != nil {
			return carvePlan{}, err
		}
		p.source.Ops = append(p.source.Ops, l)
		p.dependsOn = l.Link
	}
	if p.srcEntry, err = newID(); err != nil {
		return carvePlan{}, err
	}
	if p.tgtEntry, err = newID(); err != nil {
		return carvePlan{}, err
	}
	p.source.Drop = p.dropped
	p.source.Carve = &changeset.Carve{To: tgt.ID, Entry: p.tgtEntry}
	p.target.Drop = []id.ID{}
	p.target.Carve = &changeset.Carve{From: src.ID, Entry: p.srcEntry}
	for _, b := range []changeset.Ops{p.source, p.target} {
		if err := b.Validate(); err != nil {
			return carvePlan{}, Usagef("%s", err.Error()).WithNext("genesis describe carve")
		}
	}
	return p, nil
}

// checkCarve refuses a carve whose ops would not apply: with both entries
// pending (actor left zero, which no view reads), the what-if view of the
// target, then of the source, must show no problem it did not show before.
// The first new problem is a conflict naming its op.
func checkCarve(w *workspace, p carvePlan) error {
	problems := func(cs id.ID) []gmodel.Problem {
		_, ps := gmodel.WhatIf(w.changeSets(), cs)
		return ps
	}
	before := map[id.ID]map[string]bool{}
	for _, cs := range []id.ID{p.to, p.from} {
		before[cs] = map[string]bool{}
		for _, x := range problems(cs) {
			before[cs][x.String()] = true
		}
	}
	se, err := w.buildEntry(p.from, p.source, id.ID{}, id.ID{}, p.srcEntry)
	if err != nil {
		return err
	}
	te, err := w.buildEntry(p.to, p.target, id.ID{}, id.ID{}, p.tgtEntry)
	if err != nil {
		return err
	}
	saved := w.pending
	w.pending = append(slices.Clone(saved), []entry.Entry{se, te}...)
	defer func() { w.pending = saved }()
	for _, cs := range []id.ID{p.to, p.from} {
		for _, x := range problems(cs) {
			if before[cs][x.String()] || x.ChangeSet != cs {
				continue
			}
			what := "change set " + cs.String()
			if !x.Op.IsZero() {
				what = "op " + x.Op.String()
				if c, ok := w.changeSet(cs); ok {
					if i := slices.IndexFunc(c.Ops, func(o op.Op) bool { return o.ID == x.Op }); i >= 0 {
						o := c.Ops[i]
						what = fmt.Sprintf("%s %s on %s", o.Type, o.ID, o.Target())
					}
				}
			}
			return Conflictf("the carve would not apply in change set %s: %s: %v", cs, what, x.Err).
				WithNext("genesis slice --change " + shortID(cs) + " --view what-if")
		}
	}
	return nil
}

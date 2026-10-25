package api

// This file holds the plan reads: roadmap shows the proposed and accepted
// change sets grouped by the intent they serve, and ready says which change
// sets can be worked on now and why the others can't. Neither writes
// anything.
//
// The vocabulary:
//
//   - A change set is in flight when it is draft, proposed or accepted.
//   - A change set's links (D7) are the view's links that start or end at
//     it, plus the link.create ops inside the change set itself that start
//     or end at it, so that a proposal already says what it serves before
//     it is accepted. Duplicates (the same link ID) count once, the view's
//     copy first; retired links do not count.
//   - A change set serves the intents its serves links point to. It
//     depends on the change sets changeset.DependsOn gives (D1): its live
//     depends-on links that have landed or are in its own ops, net of its
//     own retires, whatever the view. A link another change set writes
//     does not count.
//   - A change set rests on a dropped one as changeset.RestsOn says (D7),
//     and a dependency revised under it as changeset.Revised says (D13).
//   - The nodes a change set touches are the targets of its node ops.
//   - An open question is a live question node with no live answers link
//     to it.
//   - A change set is blocked by an open question that has a live blocks
//     link to the change set, to a node it touches, or to a tree ancestor
//     of such a node (in the view, or under the new parent its own ops
//     give it).
//   - A change set is ready (D8) when it is accepted, every change set it
//     depends on has landed, none of its ops rests on a dropped change set
//     or has a dependency revised under it, and no open question blocks
//     it. One that isn't
//     says why, in a fixed order: "not accepted: <state>", then "depends
//     on <change set> (<state>)" for each dependency not landed ("(not in
//     the log)" for one that isn't), then "rests on <change set> (<state>):
//     N operations" for each dropped change set it rests on, then
//     "<change set> revised under it: N operations" for each dependency
//     whose revisions or drops went under N of its ops, then "blocked by
//     question <handle> on <this change set|node handle>", then
//     "depends-on cycle".
//
// The roadmap (D9) puts milestone work first. A milestone is an intent
// in the live tree whose milestone field is non-empty. A served intent
// maps to the nearest milestone at or above it: itself, else its nearest
// ancestor that is one (so serving a done point under a milestone counts
// for it, and with nested milestones only the inner one counts); an
// intent outside the live tree, such as a retired one, maps to none. The
// proposed and accepted change sets appear in each milestone group their
// served intents map to; groups follow the intent tree in pre-order,
// siblings by key then ID. Every other one (serving no milestone, or
// nothing) appears once in a single flat backlog after them. Inside a
// group, change sets follow depends-on (one lands before those that
// depend on it), ties by ID; members left in or after a cycle come last
// by ID with the reason "depends-on cycle". Ready lists every in-flight
// change set once and follows that order: the roadmap's ones at the first
// group each appears in, then the drafts by ID; its next is the first
// ready one.
//
// Both default to the to-be view whatever the effective change set (D17):
// their question is the plan, not the work in hand. The view field still
// overrides.
//
// The intent order is computed from each served intent's ancestor chain
// rather than a walk of the whole tree: comparing chains of (key, ID) is
// the same pre-order, and the queries stay bounded by the change sets in
// flight, not the size of the model.

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/index"
	gmodel "github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/op"
)

func init() {
	register(Operation{
		Name:     "roadmap",
		Summary:  "Show the proposed and accepted change sets: milestone groups first, then the backlog",
		Request:  reflect.TypeFor[RoadmapRequest](),
		Response: reflect.TypeFor[RoadmapResponse](),
		Run:      runRoadmap,
	})
	register(Operation{
		Name:     "ready",
		Summary:  "Say which change sets can be worked on now, and why the others can't",
		Request:  reflect.TypeFor[ReadyRequest](),
		Response: reflect.TypeFor[ReadyResponse](),
		Run:      runReady,
	})
}

// planNewNext is the next command when nothing is in flight.
const planNewNext = `genesis change new "<intent>"`

// RoadmapRequest asks for the roadmap.
type RoadmapRequest struct {
	View string `json:"view,omitempty" desc:"The view: as-is, to-be or what-if:<change set>; default to-be, whatever the current change set"`
}

// RoadmapResponse is the roadmap: its groups in order, and the view.
type RoadmapResponse struct {
	Groups []RoadmapGroup `json:"groups" desc:"One group per milestone served, in intent-tree order, then the backlog"`
	View   string         `json:"view" desc:"The view read: as-is, to-be or what-if:<ID>"`
}

// RoadmapGroup is the change sets of one milestone, or of the backlog, in
// depends-on order.
type RoadmapGroup struct {
	Kind       string             `json:"kind" desc:"milestone, or backlog for the change sets that serve no milestone"`
	Intent     RoadmapIntent      `json:"intent,omitzero" desc:"The milestone; absent for the backlog"`
	ChangeSets []RoadmapChangeSet `json:"change_sets" desc:"The group's change sets, in depends-on order, ties by ID"`
}

// RoadmapIntent names a group's intent.
type RoadmapIntent struct {
	ID     id.ID    `json:"id" desc:"The intent's ID"`
	Handle string   `json:"handle" desc:"Its handle; empty if the view does not have it"`
	Path   []string `json:"path" desc:"The handles of its ancestors in the intent tree, the root first"`
}

// RoadmapChangeSet is one change set in a group.
type RoadmapChangeSet struct {
	ID        id.ID           `json:"id" desc:"The change set's ID"`
	State     changeset.State `json:"state" desc:"proposed or accepted"`
	Kind      changeset.Kind  `json:"kind" desc:"intention or record"`
	Intent    string          `json:"intent" desc:"The change set's intent"`
	Ready     bool            `json:"ready" desc:"Whether it can be worked on now (see ready)"`
	DependsOn []id.ID         `json:"depends_on" desc:"The change sets it depends on, sorted"`
	Serves    []id.ID         `json:"serves" desc:"The intents it serves, sorted"`
}

// ReadyRequest asks which change sets are ready.
type ReadyRequest struct {
	View string `json:"view,omitempty" desc:"The view: as-is, to-be or what-if:<change set>; default to-be, whatever the current change set"`
}

// ReadyResponse is every in-flight change set with ready or its reasons.
type ReadyResponse struct {
	Next       id.ID            `json:"next,omitzero" desc:"The first ready change set; absent when none is"`
	ChangeSets []ReadyChangeSet `json:"change_sets" desc:"Every draft, proposed and accepted change set: roadmap order, then drafts by ID"`
	View       string           `json:"view" desc:"The view read: as-is, to-be or what-if:<ID>"`
}

// ReadyChangeSet is one change set, ready or not, with why not.
type ReadyChangeSet struct {
	ID      id.ID           `json:"id" desc:"The change set's ID"`
	State   changeset.State `json:"state" desc:"draft, proposed or accepted"`
	Intent  string          `json:"intent" desc:"The change set's intent"`
	Ready   bool            `json:"ready" desc:"Whether it can be worked on now"`
	Reasons []string        `json:"reasons" desc:"Why it is not ready, in a fixed order; empty when it is"`
}

// planInputs is everything the plan core reads, gathered from a view by
// planGather or written by hand in tests.
type planInputs struct {
	// changeSets is every folded change set.
	changeSets []index.ChangeSetInfo
	// ops is each in-flight change set's folded ops.
	ops map[id.ID][]op.Op
	// sets is what the dependency analyses read: every in-flight and
	// dropped change set folded, and planLanded's stand-in for the landed
	// ones.
	sets []changeset.ChangeSet
	// links holds the view's links the core needs, live or retired, in any
	// order and possibly repeated: those at in-flight change sets, blocks
	// links to touched nodes and their ancestors, and answers links to the
	// questions those come from.
	links []gmodel.Link
	// nodes holds the nodes the core names: served intents and their
	// ancestors, questions, touched nodes and their ancestors.
	nodes map[id.ID]gmodel.Node
	// ancestors maps each node the view has that a change set touches, or
	// moves or creates a node under, to its tree ancestors, in any order.
	ancestors map[id.ID][]id.ID
	// intentPath maps a served intent to its ancestor chain, the root
	// first; absent for an intent the view does not have.
	intentPath map[id.ID][]id.ID
	// milestones is the nodes whose milestone field is non-empty; the core
	// reads no fields, so planGather says which they are.
	milestones map[id.ID]bool
}

// planChangeSet is the core's verdict on one in-flight change set.
type planChangeSet struct {
	info      index.ChangeSetInfo
	serves    []id.ID
	dependsOn []id.ID
	ready     bool
	reasons   []string
}

// planGroup is one roadmap group: a milestone's, or the backlog's, whose
// intent is zero.
type planGroup struct {
	intent  id.ID
	backlog bool
	members []*planChangeSet
}

// planResult is the core's output: the roadmap groups and the ready order.
type planResult struct {
	groups []planGroup
	order  []*planChangeSet
}

// planInFlight reports whether s is draft, proposed or accepted.
func planInFlight(s changeset.State) bool {
	return s == changeset.StateDraft || s == changeset.StateProposed || s == changeset.StateAccepted
}

// planOwnLinks is the links the link.create ops of change set cs make that
// start or end at cs (D7).
func planOwnLinks(cs id.ID, ops []op.Op) []gmodel.Link {
	var out []gmodel.Link
	for _, o := range ops {
		if o.Type == op.TypeLinkCreate && (o.From == cs || o.To == cs) {
			out = append(out, gmodel.Link{ID: o.Link, Type: o.LinkType, From: o.From, To: o.To,
				Note: o.Note, Decision: o.Decision, Created: cs, Version: o.ID})
		}
	}
	return out
}

// planTouched is the nodes the node ops of ops target, in op order, each
// once, with the parents create and move ops give them.
func planTouched(ops []op.Op) (nodes []id.ID, parents map[id.ID][]id.ID) {
	seen := map[id.ID]bool{}
	parents = map[id.ID][]id.ID{}
	for _, o := range ops {
		switch o.Type {
		case op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeMove, op.TypeNodeRetire:
		default:
			continue
		}
		if !seen[o.Node] {
			seen[o.Node] = true
			nodes = append(nodes, o.Node)
		}
		if (o.Type == op.TypeNodeCreate || o.Type == op.TypeNodeMove) && !o.Parent.IsZero() {
			parents[o.Node] = append(parents[o.Node], o.Parent)
		}
	}
	return nodes, parents
}

// planBlocker is one open question blocking a change set on a target: the
// change set itself or a node.
type planBlocker struct {
	question id.ID
	target   id.ID
}

// planCore computes every in-flight change set's verdict, the roadmap
// groups and the ready order from in, as the file comment says.
func planCore(in planInputs) planResult {
	// The view's live links, each once, by end.
	seen := map[id.ID]bool{}
	byEnd := map[id.ID][]gmodel.Link{}
	var live []gmodel.Link
	for _, l := range in.links {
		if l.Retired || seen[l.ID] {
			seen[l.ID] = true
			continue
		}
		seen[l.ID] = true
		live = append(live, l)
	}
	slices.SortFunc(live, func(a, b gmodel.Link) int { return id.Compare(a.ID, b.ID) })
	byTo := map[id.ID][]gmodel.Link{}
	for _, l := range live {
		byEnd[l.From] = append(byEnd[l.From], l)
		if l.To != l.From {
			byEnd[l.To] = append(byEnd[l.To], l)
		}
		byTo[l.To] = append(byTo[l.To], l)
	}
	open := func(q id.ID) bool {
		n, ok := in.nodes[q]
		if !ok || n.Retired || n.Kind != op.KindQuestion {
			return false
		}
		for _, l := range byTo[q] {
			if l.Type == op.LinkAnswers {
				return false
			}
		}
		return true
	}
	states := map[id.ID]changeset.State{}
	for _, c := range in.changeSets {
		states[c.ID] = c.State
	}

	var all []*planChangeSet
	for _, c := range in.changeSets {
		if !planInFlight(c.State) {
			continue
		}
		p := &planChangeSet{info: c}
		links := planLinksOf(c.ID, byEnd[c.ID], seen, in.ops[c.ID])
		var blockers []planBlocker
		for _, l := range links {
			switch {
			case l.From == c.ID && l.Type == op.LinkServes:
				p.serves = append(p.serves, l.To)
			case l.To == c.ID && l.Type == op.LinkBlocks && open(l.From):
				blockers = append(blockers, planBlocker{question: l.From, target: c.ID})
			}
		}
		touched, parents := planTouched(in.ops[c.ID])
		targets := map[id.ID]bool{}
		for _, t := range touched {
			targets[t] = true
			for _, a := range in.ancestors[t] {
				targets[a] = true
			}
			for _, p := range parents[t] {
				if as, ok := in.ancestors[p]; ok {
					targets[p] = true
					for _, a := range as {
						targets[a] = true
					}
				}
			}
		}
		for t := range targets {
			for _, l := range byTo[t] {
				if l.Type == op.LinkBlocks && open(l.From) {
					blockers = append(blockers, planBlocker{question: l.From, target: t})
				}
			}
		}
		p.serves = planSortIDs(p.serves)
		p.dependsOn = changeset.DependsOn(in.sets, c.ID)

		if c.State != changeset.StateAccepted {
			p.reasons = append(p.reasons, "not accepted: "+string(c.State))
		}
		for _, d := range p.dependsOn {
			switch s, ok := states[d]; {
			case !ok:
				p.reasons = append(p.reasons, fmt.Sprintf("depends on %s (not in the log)", shortID(d)))
			case s != changeset.StateLanded:
				p.reasons = append(p.reasons, fmt.Sprintf("depends on %s (%s)", shortID(d), s))
			}
		}
		for _, r := range changeset.RestsOn(in.sets, c.ID) {
			p.reasons = append(p.reasons, fmt.Sprintf("rests on %s (%s): %d operations", shortID(r.ChangeSet), r.State, len(r.Ops)))
		}
		p.reasons = append(p.reasons, planRevisedReasons(changeset.Revised(in.sets, c.ID))...)
		for _, b := range planSortBlockers(c.ID, blockers, in.nodes) {
			on := "this change set"
			if b.target != c.ID {
				on = planHandle(b.target, in.nodes)
			}
			p.reasons = append(p.reasons, fmt.Sprintf("blocked by question %s on %s", planHandle(b.question, in.nodes), on))
		}
		p.ready = len(p.reasons) == 0
		all = append(all, p)
	}

	// The roadmap: proposed and accepted change sets by the milestones
	// their served intents map to, then the rest as the backlog.
	in.intentPath = planAncestorPaths(in.intentPath)
	members := map[id.ID][]*planChangeSet{}
	var milestones []id.ID
	var backlog []*planChangeSet
	for _, p := range all {
		if p.info.State == changeset.StateDraft {
			continue
		}
		var ms []id.ID
		for _, i := range p.serves {
			if m, ok := planMilestoneOf(i, in); ok {
				ms = append(ms, m)
			}
		}
		ms = planSortIDs(ms)
		if len(ms) == 0 {
			backlog = append(backlog, p)
		}
		for _, m := range ms {
			if _, ok := members[m]; !ok {
				milestones = append(milestones, m)
			}
			members[m] = append(members[m], p)
		}
	}
	slices.SortFunc(milestones, func(a, b id.ID) int { return planCompareIntents(a, b, in) })
	var res planResult
	for _, m := range milestones {
		res.groups = append(res.groups, planGroup{intent: m, members: planTopo(members[m])})
	}
	if len(backlog) > 0 {
		res.groups = append(res.groups, planGroup{backlog: true, members: planTopo(backlog)})
	}

	listed := map[id.ID]bool{}
	for _, g := range res.groups {
		for _, p := range g.members {
			if !listed[p.info.ID] {
				listed[p.info.ID] = true
				res.order = append(res.order, p)
			}
		}
	}
	for _, p := range all {
		if !listed[p.info.ID] {
			res.order = append(res.order, p)
		}
	}
	return res
}

// planAncestorPaths is paths with each ancestor on a chain given its own
// chain, the prefix before it, where paths lacks one; paths is not changed.
func planAncestorPaths(paths map[id.ID][]id.ID) map[id.ID][]id.ID {
	out := maps.Clone(paths)
	for _, path := range paths {
		for k, a := range path {
			if _, ok := out[a]; !ok {
				out[a] = path[:k:k]
			}
		}
	}
	return out
}

// planMilestoneOf is the nearest milestone at or above served intent i:
// i itself, then its ancestors from the nearest up. ok is false when i is
// outside the live tree (so none of its ancestors' milestones count) or no
// milestone is at or above it.
func planMilestoneOf(i id.ID, in planInputs) (m id.ID, ok bool) {
	if _, live := planTreeKey(i, in); !live {
		return id.ID{}, false
	}
	if in.milestones[i] {
		return i, true
	}
	path := in.intentPath[i]
	for k := len(path) - 1; k >= 0; k-- {
		if a := path[k]; in.milestones[a] {
			if _, live := planTreeKey(a, in); live {
				return a, true
			}
		}
	}
	return id.ID{}, false
}

// planLinksOf is change set cs's live links (D7): view, the live view
// links at cs, sorted by ID, then those its own ops create whose ID is
// not in seen (every link the view has, live or retired), each once.
func planLinksOf(cs id.ID, view []gmodel.Link, inView map[id.ID]bool, ops []op.Op) []gmodel.Link {
	out := slices.Clone(view)
	seen := maps.Clone(inView)
	for _, l := range planOwnLinks(cs, ops) {
		if !seen[l.ID] {
			seen[l.ID] = true
			out = append(out, l)
		}
	}
	return out
}

// planRevisedReasons is one reason per dependency in rs, by ID, counting
// the dependent's ops it went under.
func planRevisedReasons(rs []changeset.RevisedConflict) []string {
	under := map[id.ID][]id.ID{}
	for _, r := range rs {
		under[r.Dependency] = append(under[r.Dependency], r.Op)
	}
	var out []string
	for _, d := range slices.SortedFunc(maps.Keys(under), id.Compare) {
		out = append(out, fmt.Sprintf("%s revised under it: %d operations", shortID(d), len(planSortIDs(under[d]))))
	}
	return out
}

// planSortIDs sorts ids and drops repeats.
func planSortIDs(ids []id.ID) []id.ID {
	slices.SortFunc(ids, id.Compare)
	return slices.Compact(ids)
}

// planHandle is node i's handle, or its short ID when nodes lacks it.
func planHandle(i id.ID, nodes map[id.ID]gmodel.Node) string {
	if n, ok := nodes[i]; ok && n.Handle != "" {
		return n.Handle
	}
	return shortID(i)
}

// planSortBlockers sorts bs, each once: those on change set cs first, then
// by target handle and ID, then question handle and ID.
func planSortBlockers(cs id.ID, bs []planBlocker, nodes map[id.ID]gmodel.Node) []planBlocker {
	slices.SortFunc(bs, func(a, b planBlocker) int {
		ac, bc := a.target == cs, b.target == cs
		switch {
		case ac && !bc:
			return -1
		case bc && !ac:
			return 1
		}
		return cmp.Or(
			strings.Compare(planHandle(a.target, nodes), planHandle(b.target, nodes)),
			id.Compare(a.target, b.target),
			strings.Compare(planHandle(a.question, nodes), planHandle(b.question, nodes)),
			id.Compare(a.question, b.question))
	})
	return slices.Compact(bs)
}

// planTreeKey is an intent's place in the live tree: the (key, ID) of each
// node from its root down to itself; ok is false for an intent outside the
// live tree (retired, under a retired node, or unknown to the view).
func planTreeKey(i id.ID, in planInputs) (key []gmodel.Node, ok bool) {
	n, found := in.nodes[i]
	path, inPath := in.intentPath[i]
	if !found || !inPath || n.Retired {
		return nil, false
	}
	for _, a := range path {
		an, found := in.nodes[a]
		if !found || an.Retired {
			return nil, false
		}
		key = append(key, an)
	}
	key = append(key, n)
	if !key[0].Parent.IsZero() {
		return nil, false
	}
	return key, true
}

// planCompareIntents orders intents a and b: the live tree in pre-order,
// siblings by key then ID, then the rest by ID.
func planCompareIntents(a, b id.ID, in planInputs) int {
	ka, oka := planTreeKey(a, in)
	kb, okb := planTreeKey(b, in)
	switch {
	case oka && !okb:
		return -1
	case okb && !oka:
		return 1
	case !oka:
		return id.Compare(a, b)
	}
	return slices.CompareFunc(ka, kb, func(x, y gmodel.Node) int {
		return cmp.Or(strings.Compare(x.Key, y.Key), id.Compare(x.ID, y.ID))
	})
}

// planTopo orders a group's change sets by depends-on among them: Kahn's
// algorithm, a dependency first, ties by ID. Those left in or after a
// cycle come last by ID, with the reason "depends-on cycle".
func planTopo(ps []*planChangeSet) []*planChangeSet {
	ps = slices.Clone(ps)
	slices.SortFunc(ps, func(a, b *planChangeSet) int { return id.Compare(a.info.ID, b.info.ID) })
	in := map[id.ID]bool{}
	for _, p := range ps {
		in[p.info.ID] = true
	}
	waiting := map[id.ID]int{}
	dependents := map[id.ID][]*planChangeSet{}
	for _, p := range ps {
		for _, d := range p.dependsOn {
			if in[d] && d != p.info.ID {
				waiting[p.info.ID]++
				dependents[d] = append(dependents[d], p)
			}
		}
	}
	var out, queue []*planChangeSet
	for _, p := range ps {
		if waiting[p.info.ID] == 0 {
			queue = append(queue, p)
		}
	}
	done := map[id.ID]bool{}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		out = append(out, p)
		done[p.info.ID] = true
		for _, q := range dependents[p.info.ID] {
			if waiting[q.info.ID]--; waiting[q.info.ID] == 0 {
				i, _ := slices.BinarySearchFunc(queue, q, func(a, b *planChangeSet) int { return id.Compare(a.info.ID, b.info.ID) })
				queue = slices.Insert(queue, i, q)
			}
		}
	}
	for _, p := range ps {
		if done[p.info.ID] {
			continue
		}
		if !slices.Contains(p.reasons, planCycleReason) {
			p.reasons = append(p.reasons, planCycleReason)
			p.ready = false
		}
		out = append(out, p)
	}
	return out
}

// planCycleReason is the reason of a change set in or after a depends-on
// cycle.
const planCycleReason = "depends-on cycle"

// planGather reads planInputs from r with queries bounded by the change
// sets in flight: each one's entries and links, the nodes its ops touch,
// their ancestors and the blocks links to them, the questions those come
// from and the answers to them, and the served intents with their
// ancestors. The entries and the links at the change sets are read with
// one query each, the nodes light (the core reads no fields) in a batch
// per level of the tree, and the blocks and answers links with one query
// each: every query has a cost of its own (a read transaction's locks), so
// the count stays fixed, not one per change set or node.
func planGather(r *reader) (planInputs, error) {
	in := planInputs{ops: map[id.ID][]op.Op{}, nodes: map[id.ID]gmodel.Node{},
		ancestors: map[id.ID][]id.ID{}, intentPath: map[id.ID][]id.ID{}, milestones: map[id.ID]bool{}}
	var err error
	if in.changeSets, err = r.ChangeSets(); err != nil {
		return planInputs{}, err
	}
	// missing holds the IDs asked for that the view lacks, so that none is
	// asked for twice.
	missing := map[id.ID]bool{}
	// fetch reads the nodes of ids not read yet in one batch, and returns
	// the ones it read.
	fetch := func(ids []id.ID) ([]gmodel.Node, error) {
		var want []id.ID
		for _, i := range ids {
			if _, ok := in.nodes[i]; !ok && !missing[i] && !i.IsZero() {
				want = append(want, i)
			}
		}
		if len(want) == 0 {
			return nil, nil
		}
		ns, err := r.NodesLight(want)
		if err != nil {
			return nil, err
		}
		for _, n := range ns {
			in.nodes[n.ID] = n
		}
		for _, i := range want {
			if _, ok := in.nodes[i]; !ok {
				missing[i] = true
			}
		}
		return ns, nil
	}
	// prefetch reads starts and their ancestors, a level of the tree per
	// batch.
	prefetch := func(starts []id.ID) error {
		level := starts
		for range maxDepth {
			ns, err := fetch(level)
			if err != nil || len(ns) == 0 {
				return err
			}
			level = level[:0:0]
			for _, n := range ns {
				level = append(level, n.Parent)
			}
		}
		return nil
	}
	node := func(i id.ID) (gmodel.Node, bool, error) {
		if _, err := fetch([]id.ID{i}); err != nil {
			return gmodel.Node{}, false, err
		}
		n, ok := in.nodes[i]
		return n, ok, nil
	}
	var flying, dropped []id.ID
	for _, c := range in.changeSets {
		switch {
		case planInFlight(c.State):
			flying = append(flying, c.ID)
		case changeset.Dropped(c.State):
			dropped = append(dropped, c.ID)
		}
	}
	folding := slices.Concat(flying, dropped)
	entries, err := r.EntriesOf(folding)
	if err != nil {
		return planInputs{}, err
	}
	csLinks, err := r.LinksAt(flying, index.Both)
	if err != nil {
		return planInputs{}, err
	}
	in.links = append(in.links, csLinks...)
	csLinks = slices.Clone(csLinks)
	// Folding is CPU work over entries already read: one change set per
	// goroutine.
	folded := make([]changeset.ChangeSet, len(folding))
	var wg sync.WaitGroup
	for k, c := range folding {
		wg.Go(func() {
			folded[k], _ = changeset.Fold(c, entries[c])
		})
	}
	wg.Wait()
	for k, c := range flying {
		in.ops[c] = folded[k].Ops
		csLinks = append(csLinks, planOwnLinks(c, folded[k].Ops)...)
	}
	landed, err := planLanded(r, flying, folded[len(flying):])
	if err != nil {
		return planInputs{}, err
	}
	in.sets = append(folded, landed)
	if in.sets, err = planLandedDeps(r, in.changeSets, flying, in.sets); err != nil {
		return planInputs{}, err
	}

	// Every node a chain below starts from: the touched nodes, the parents
	// ops give them, and the served intents; read with their ancestors.
	var starts []id.ID
	for _, c := range in.changeSets {
		touched, parents := planTouched(in.ops[c.ID])
		for _, t := range touched {
			starts = append(starts, t)
			starts = append(starts, parents[t]...)
		}
	}
	for _, l := range csLinks {
		if l.Type == op.LinkServes && !l.Retired {
			starts = append(starts, l.To)
		}
	}
	if err := prefetch(starts); err != nil {
		return planInputs{}, err
	}

	// Touched nodes, their ancestors, and the blocks links to them.
	blocked := map[id.ID]bool{}
	var questions []id.ID
	addBlocks := func(ls []gmodel.Link) {
		for _, l := range ls {
			if l.Type == op.LinkBlocks && !blocked[l.From] {
				blocked[l.From] = true
				questions = append(questions, l.From)
			}
		}
	}
	addBlocks(csLinks)
	// chain walks up from start through the view's parents, start first,
	// since change sets touch nodes in the same parts of the tree; nil when
	// the view has no start.
	chain := func(start id.ID) ([]id.ID, error) {
		var out []id.ID
		for p := start; !p.IsZero() && !slices.Contains(out, p) && len(out) < maxDepth; {
			n, ok, err := node(p)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
			out = append(out, p)
			p = n.Parent
		}
		return out, nil
	}
	targets := map[id.ID]bool{}
	for _, c := range in.changeSets {
		touched, parents := planTouched(in.ops[c.ID])
		for _, t := range touched {
			for _, s := range append([]id.ID{t}, parents[t]...) {
				if _, done := in.ancestors[s]; done {
					continue
				}
				ch, err := chain(s)
				if err != nil {
					return planInputs{}, err
				}
				if len(ch) == 0 {
					continue
				}
				in.ancestors[s] = planSortIDs(slices.Clone(ch[1:]))
				for _, a := range ch {
					targets[a] = true
				}
			}
		}
	}
	ls, err := r.LinksTo(slices.Collect(maps.Keys(targets)), op.LinkBlocks)
	if err != nil {
		return planInputs{}, err
	}
	in.links = append(in.links, ls...)
	addBlocks(ls)

	// The questions, and the answers links to them.
	if _, err := fetch(questions); err != nil {
		return planInputs{}, err
	}
	var asked []id.ID
	for _, q := range questions {
		if _, ok := in.nodes[q]; ok {
			asked = append(asked, q)
		}
	}
	if ls, err = r.LinksTo(asked, op.LinkAnswers); err != nil {
		return planInputs{}, err
	}
	in.links = append(in.links, ls...)

	// The served intents and their ancestors; one the view lacks but an
	// in-flight change set creates is named from its create.
	for _, l := range csLinks {
		if l.Type != op.LinkServes || l.Retired {
			continue
		}
		if _, done := in.intentPath[l.To]; done {
			continue
		}
		ch, err := chain(l.To)
		if err != nil {
			return planInputs{}, err
		}
		if len(ch) == 0 {
			planCreatedNode(l.To, in)
			continue
		}
		path := slices.Clone(ch[1:])
		slices.Reverse(path)
		in.intentPath[l.To] = path
	}

	// The milestones among the served intents and their ancestors: their
	// fields, read in one batch.
	var intents []id.ID
	for i, path := range in.intentPath {
		for _, n := range append([]id.ID{i}, path...) {
			if in.nodes[n].Kind == op.KindIntent {
				intents = append(intents, n)
			}
		}
	}
	ns, err := r.NodesFields(intents)
	if err != nil {
		return planInputs{}, err
	}
	for _, n := range ns {
		if n.Kind == op.KindIntent && planMilestoneField(n.Fields) {
			in.milestones[n.ID] = true
		}
	}
	return in, nil
}

// planMilestoneField reports whether fields has a non-empty milestone:
// present, and neither JSON "" nor null.
func planMilestoneField(fields map[string]json.RawMessage) bool {
	v, ok := fields["milestone"]
	if !ok {
		return false
	}
	v = bytes.TrimSpace(v)
	return len(v) > 0 && string(v) != `""` && string(v) != "null"
}

// planLanded is a stand-in for the landed change sets, so that the
// analyses over planInputs.sets need not read their entries (too many to
// fold within a read's budget). It holds what they ask of landed ops: a
// link.create for each depends-on link from a flying change set that is
// live in as-is, whatever r's view (a landed retire leaves it out), and a
// create of each node or link a dropped change set creates that as-is has
// (a landed change set creates it too). Its ID is zero.
func planLanded(r *reader, flying []id.ID, dropped []changeset.ChangeSet) (changeset.ChangeSet, error) {
	out := changeset.ChangeSet{State: changeset.StateLanded}
	var nodes []id.ID
	links := map[id.ID]bool{}
	from := slices.Clone(flying)
	for _, s := range dropped {
		for _, o := range s.Ops {
			switch o.Type {
			case op.TypeNodeCreate:
				nodes = append(nodes, o.Node)
			case op.TypeLinkCreate:
				links[o.Link] = true
				from = append(from, o.From)
			}
		}
	}
	ns, err := r.x.NodesLight(index.AsIs, nodes)
	if err != nil {
		return changeset.ChangeSet{}, readError(err)
	}
	for _, n := range ns {
		out.Ops = append(out.Ops, op.Op{Type: op.TypeNodeCreate, Node: n.ID})
	}
	ls, err := r.x.LinksAt(index.AsIs, planSortIDs(from), index.From)
	if err != nil {
		return changeset.ChangeSet{}, readError(err)
	}
	for _, l := range ls {
		if links[l.ID] || l.Type == op.LinkDependsOn && !l.Retired {
			out.Ops = append(out.Ops, op.Op{Type: op.TypeLinkCreate, Link: l.ID, LinkType: l.Type, From: l.From, To: l.To})
		}
	}
	return out, nil
}

// planLandedDeps adds to sets, folded with their real IDs and History, the
// landed change sets that a flying change set depends on, directly or
// through landed ones, so that changeset.Revised sees a landed
// dependency's revisions and drops (D13). Only those are read; the walk
// follows the depends-on links in sets and, for each landed change set
// added, its live as-is depends-on links, which join the stand-in (the
// zero-ID landed change set in sets).
func planLandedDeps(r *reader, all []index.ChangeSetInfo, flying []id.ID, sets []changeset.ChangeSet) ([]changeset.ChangeSet, error) {
	states := make(map[id.ID]changeset.State, len(all))
	for _, c := range all {
		states[c.ID] = c.State
	}
	stand := slices.IndexFunc(sets, func(s changeset.ChangeSet) bool { return s.ID.IsZero() })
	seen := map[id.ID]bool{}
	from := flying
	for len(from) > 0 {
		var next []id.ID
		for _, c := range from {
			for _, d := range changeset.DependsOn(sets, c) {
				if !seen[d] && states[d] == changeset.StateLanded {
					seen[d] = true
					next = append(next, d)
				}
			}
		}
		if len(next) == 0 {
			break
		}
		slices.SortFunc(next, id.Compare)
		entries, err := r.EntriesOf(next)
		if err != nil {
			return nil, err
		}
		for _, c := range next {
			s, _ := changeset.Fold(c, entries[c])
			sets = append(sets, s)
		}
		ls, err := r.x.LinksAt(index.AsIs, next, index.From)
		if err != nil {
			return nil, readError(err)
		}
		for _, l := range ls {
			if l.Type == op.LinkDependsOn && !l.Retired && stand >= 0 {
				sets[stand].Ops = append(sets[stand].Ops,
					op.Op{Type: op.TypeLinkCreate, Link: l.ID, LinkType: l.Type, From: l.From, To: l.To})
			}
		}
		from = next
	}
	return sets, nil
}

// planCreatedNode adds node i to in.nodes from the node.create op of an
// in-flight change set, if one creates it, so it can be named; an intent
// it creates with a non-empty milestone field is a milestone.
func planCreatedNode(i id.ID, in planInputs) {
	for _, c := range in.changeSets {
		for _, o := range in.ops[c.ID] {
			if o.Type == op.TypeNodeCreate && o.Node == i {
				in.nodes[i] = gmodel.Node{ID: i, Kind: o.Kind, Handle: o.Handle, Parent: o.Parent, Key: o.Key}
				if o.Kind == op.KindIntent && planMilestoneField(o.Fields) {
					in.milestones[i] = true
				}
				return
			}
		}
	}
}

// planRead opens a reader for opName in view (to-be by default, D17),
// gathers the inputs and runs the core.
func planRead(ctx *Ctx, opName, view string) (planResult, planInputs, string, error) {
	r, err := openReader(ctx, opName, view, toBeView)
	if err != nil {
		return planResult{}, planInputs{}, "", err
	}
	defer r.Close()
	in, err := planGather(r)
	if err != nil {
		return planResult{}, planInputs{}, "", err
	}
	return planCore(in), in, r.View(), nil
}

// runRoadmap answers roadmap.
func runRoadmap(ctx *Ctx, req any) (any, error) {
	res, in, v, err := planRead(ctx, "roadmap", req.(RoadmapRequest).View)
	if err != nil {
		return nil, err
	}
	resp := RoadmapResponse{Groups: []RoadmapGroup{}, View: v}
	// A milestone may be only an ancestor of the intents served.
	paths := planAncestorPaths(in.intentPath)
	for _, g := range res.groups {
		rg := RoadmapGroup{Kind: "backlog", ChangeSets: []RoadmapChangeSet{}}
		if !g.backlog {
			rg.Kind = "milestone"
			rg.Intent = RoadmapIntent{ID: g.intent, Handle: in.nodes[g.intent].Handle, Path: []string{}}
			for _, a := range paths[g.intent] {
				rg.Intent.Path = append(rg.Intent.Path, planHandle(a, in.nodes))
			}
		}
		for _, p := range g.members {
			rg.ChangeSets = append(rg.ChangeSets, RoadmapChangeSet{ID: p.info.ID, State: p.info.State,
				Kind: p.info.Kind, Intent: p.info.Intent, Ready: p.ready, DependsOn: orEmpty(p.dependsOn),
				Serves: orEmpty(slices.Clone(p.serves))})
		}
		resp.Groups = append(resp.Groups, rg)
	}
	return resp, nil
}

// Human shows the roadmap: a heading per group, a line per change set with
// its ready mark or state.
func (r RoadmapResponse) Human(w io.Writer) {
	fmt.Fprintln(w, "# Roadmap")
	if len(r.Groups) == 0 {
		fmt.Fprintf(w, "\nNothing planned yet.\n\nnext: %s\n", planNewNext)
		return
	}
	for _, g := range r.Groups {
		name := "Backlog"
		if g.Kind != "backlog" {
			name = strings.Join(append(slices.Clone(g.Intent.Path), cmp.Or(g.Intent.Handle, shortID(g.Intent.ID))), " / ")
		}
		fmt.Fprintf(w, "\n## %s\n\n", name)
		for _, c := range g.ChangeSets {
			mark := string(c.State)
			if c.Ready {
				mark = "ready"
			}
			fmt.Fprintf(w, "- [%s] %s %s\n", mark, shortID(c.ID), c.Intent)
		}
	}
}

// runReady answers ready.
func runReady(ctx *Ctx, req any) (any, error) {
	res, _, v, err := planRead(ctx, "ready", req.(ReadyRequest).View)
	if err != nil {
		return nil, err
	}
	resp := ReadyResponse{ChangeSets: []ReadyChangeSet{}, View: v}
	for _, p := range res.order {
		if p.ready && resp.Next.IsZero() {
			resp.Next = p.info.ID
		}
		resp.ChangeSets = append(resp.ChangeSets, ReadyChangeSet{ID: p.info.ID, State: p.info.State,
			Intent: p.info.Intent, Ready: p.ready, Reasons: append([]string{}, p.reasons...)})
	}
	return resp, nil
}

// Human names the next change set (or why nothing is ready), then lists
// the ready ones and the others with their reasons.
func (r ReadyResponse) Human(w io.Writer) {
	var ready, notReady []ReadyChangeSet
	for _, c := range r.ChangeSets {
		if c.Ready {
			ready = append(ready, c)
		} else {
			notReady = append(notReady, c)
		}
	}
	switch {
	case len(r.ChangeSets) == 0:
		fmt.Fprintf(w, "Nothing is ready: no change set is in flight.\n\nnext: %s\n", planNewNext)
		return
	case !r.Next.IsZero():
		for _, c := range ready {
			if c.ID == r.Next {
				fmt.Fprintf(w, "next: %s %s\n", shortID(c.ID), c.Intent)
			}
		}
	default:
		c := notReady[0]
		fmt.Fprintf(w, "Nothing is ready.\n%s %s: %s\n", shortID(c.ID), c.Intent, c.Reasons[0])
	}
	if len(ready) > 0 {
		fmt.Fprintf(w, "\n## Ready\n\n")
		for _, c := range ready {
			fmt.Fprintf(w, "- %s %s\n", shortID(c.ID), c.Intent)
		}
	}
	if len(notReady) > 0 {
		fmt.Fprintf(w, "\n## Not ready\n\n")
		for _, c := range notReady {
			fmt.Fprintf(w, "- [%s] %s %s\n", c.State, shortID(c.ID), c.Intent)
			for _, why := range c.Reasons {
				fmt.Fprintf(w, "  - %s\n", why)
			}
		}
	}
}

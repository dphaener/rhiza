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
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
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
	// sets is what the dependency analyses read: every in-flight change
	// set folded, the dropped ones too when an in-flight op may rest on
	// one (planMayRest), and planLanded's stand-in for the landed ones.
	sets []changeset.ChangeSet
	// links holds the view's links the core needs, live or retired, in any
	// order and possibly repeated: those at in-flight change sets, blocks
	// links to touched nodes and their ancestors, and answers links to the
	// questions those come from.
	links []gmodel.Link
	// nodes holds the nodes the core names: served intents and their
	// ancestors, questions, and the nodes blocks links point to.
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

	// The analyses share what they work out from in.sets.
	an := changeset.Analyze(in.sets)
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
		p.dependsOn = an.DependsOn(c.ID)

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
		for _, r := range an.RestsOn(c.ID) {
			p.reasons = append(p.reasons, fmt.Sprintf("rests on %s (%s): %d operations", shortID(r.ChangeSet), r.State, len(r.Ops)))
		}
		p.reasons = append(p.reasons, planRevisedReasons(an.Revised(c.ID))...)
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
	var own map[id.ID]bool
	for _, l := range planOwnLinks(cs, ops) {
		if !inView[l.ID] && !own[l.ID] {
			if own == nil {
				own = map[id.ID]bool{}
			}
			own[l.ID] = true
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
// ancestors. Every query has a cost of its own (a read transaction's
// locks), so the count stays bounded by the depth of the tree, not one per
// change set or node, and the two connections the index keeps are both
// busy: the in-flight entries are read, decoded and folded on one
// (index.EachInFlight) while the change sets and the links at them are
// read on the other, and then, while the fold still runs, a planWalk on
// that connection climbs from the nodes each folded change set touches,
// a level of the tree per query, probing the blocks links to each level.
// Once the fold ends the walk probes on the connection the fold held while
// it climbs on its own; every query stays within index.BatchOne, so no
// third connection is opened.
func planGather(r *reader) (planInputs, error) {
	in := planInputs{ops: map[id.ID][]op.Op{}, nodes: map[id.ID]gmodel.Node{},
		ancestors: map[id.ID][]id.ID{}, intentPath: map[id.ID][]id.ID{}, milestones: map[id.ID]bool{}}
	w := newPlanWalk(r)
	// flown is the in-flight change sets folded, keyed by ID: their entries
	// are read and folded on another connection while the change sets are
	// listed, since the index knows which are in flight. Each change set
	// folded hands the walk the nodes it touches, the parents its ops give
	// them and the intents its own links serve.
	var flown struct {
		sets map[id.ID]changeset.ChangeSet
		err  error
	}
	var flownWG sync.WaitGroup
	defer flownWG.Wait()
	flownWG.Go(func() {
		var mu sync.Mutex
		sets := map[id.ID]changeset.ChangeSet{}
		err := r.x.EachInFlight(func(c id.ID, es []entry.Entry) {
			// Fold problems are dropped, as plan reads always dropped Fold's: these
			// entries come from the index's store, whose bodies entry.Decode checked
			// to be canonical at catch-up, and fold problems are reported by check's
			// fold rule and change show/ls (and stored by the index at catch-up),
			// not by plan reads.
			s, _ := changeset.FoldCanonical(c, es)
			w.add(planStarts(c, s.Ops, nil))
			mu.Lock()
			sets[c] = s
			mu.Unlock()
		})
		if err != nil {
			flown.err = readError(err)
			return
		}
		flown.sets = sets
	})
	var err error
	if in.changeSets, err = r.ChangeSets(); err != nil {
		return planInputs{}, err
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
	csLinks, err := r.LinksAt(flying, index.Both)
	if err != nil {
		return planInputs{}, err
	}
	in.links = append(in.links, csLinks...)
	csLinks = slices.Clone(csLinks)
	w.add(planServed(csLinks))
	// The connection the change sets were read on is free until the fold
	// ends: the walk takes it.
	w.start()
	defer w.wait()
	flownWG.Wait()
	w.folded()
	if flown.err != nil {
		return planInputs{}, flown.err
	}
	// A change set in flight has entries; one without folds to no open.
	folded := make([]changeset.ChangeSet, len(flying))
	for k, c := range flying {
		var ok bool
		if folded[k], ok = flown.sets[c]; !ok {
			// Fold problems are dropped, as plan reads always dropped Fold's: these
			// entries come from the index's store, whose bodies entry.Decode checked
			// to be canonical at catch-up, and fold problems are reported by check's
			// fold rule and change show/ls (and stored by the index at catch-up),
			// not by plan reads.
			folded[k], _ = changeset.FoldCanonical(c, nil)
		}
		in.ops[c] = folded[k].Ops
		csLinks = append(csLinks, planOwnLinks(c, folded[k].Ops)...)
	}
	// The dropped change sets are read only when an in-flight op may rest
	// on one: they pile up with the log, and most reads need none.
	var gone []changeset.ChangeSet
	if len(dropped) > 0 {
		may, err := planMayRest(r, in.changeSets, folded)
		if err != nil {
			return planInputs{}, err
		}
		if may {
			entries, err := r.EntriesOf(dropped)
			if err != nil {
				return planInputs{}, err
			}
			gone = planFold(dropped, entries)
		}
	}
	landed, err := planLanded(r, flying, gone)
	if err != nil {
		return planInputs{}, err
	}
	in.sets = slices.Concat(folded, gone, []changeset.ChangeSet{landed})
	if in.sets, err = planLandedDeps(r, in.changeSets, flying, in.sets); err != nil {
		return planInputs{}, err
	}

	// Every node a chain below starts from: the touched nodes and the
	// parents ops give them (chainStarts), and the served intents.
	var chainStarts []id.ID
	for _, c := range in.changeSets {
		chainStarts = planStarts(c.ID, in.ops[c.ID], chainStarts)
	}
	served := planServed(csLinks)
	if err := w.wait(); err != nil {
		return planInputs{}, err
	}
	// chainBuf holds the last chain, which each caller copies what it keeps
	// from before the next.
	var chainBuf []id.ID
	// chain walks up from start through the view's parents, start first;
	// empty when the view has no start. The walk has read the parents of
	// every start it was handed; one it was not is read now.
	chain := func(start id.ID) ([]id.ID, error) {
		out := chainBuf[:0]
		defer func() { chainBuf = out[:0] }()
		for p := start; !p.IsZero() && !slices.Contains(out, p) && len(out) < maxDepth; {
			parent, ok, err := w.parent(p)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
			out = append(out, p)
			p = parent
		}
		return out, nil
	}
	// Touched nodes, their ancestors, and the blocks links to them.
	targets := map[id.ID]bool{}
	for _, s := range chainStarts {
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
	blocked := map[id.ID]bool{}
	var questions, named []id.ID
	addBlocks := func(ls []gmodel.Link) {
		for _, l := range ls {
			if l.Type == op.LinkBlocks && !blocked[l.From] {
				blocked[l.From] = true
				questions = append(questions, l.From)
			}
		}
	}
	addBlocks(csLinks)
	blocks, err := w.blocks(targets)
	if err != nil {
		return planInputs{}, err
	}
	for _, l := range blocks {
		named = append(named, l.To)
	}
	in.links = append(in.links, blocks...)
	addBlocks(blocks)

	// The served intents and their ancestors, with their fields, since
	// the milestones are among them; one the view lacks but an in-flight
	// change set creates is named from its create.
	var onPaths []id.ID
	for _, i := range served {
		if _, done := in.intentPath[i]; done {
			continue
		}
		ch, err := chain(i)
		if err != nil {
			return planInputs{}, err
		}
		if len(ch) == 0 {
			planCreatedNode(i, in)
			continue
		}
		onPaths = append(onPaths, ch...)
		path := slices.Clone(ch[1:])
		slices.Reverse(path)
		in.intentPath[i] = path
	}
	ns, err := r.NodesFields(onPaths)
	if err != nil {
		return planInputs{}, err
	}
	for _, n := range ns {
		in.nodes[n.ID] = lightNode(n)
		if n.Kind == op.KindIntent && planMilestoneField(n.Fields) {
			in.milestones[n.ID] = true
		}
	}

	// The questions and the nodes the blocks links name, and the answers
	// links to the questions the view has.
	var want []id.ID
	for _, i := range slices.Concat(questions, named) {
		if _, ok := in.nodes[i]; !ok && !i.IsZero() {
			want = append(want, i)
		}
	}
	if ns, err = r.NodesLight(want); err != nil {
		return planInputs{}, err
	}
	for _, n := range ns {
		in.nodes[n.ID] = n
	}
	var asked []id.ID
	for _, q := range questions {
		if _, ok := in.nodes[q]; ok {
			asked = append(asked, q)
		}
	}
	ls, err := r.LinksTo(asked, op.LinkAnswers)
	if err != nil {
		return planInputs{}, err
	}
	in.links = append(in.links, ls...)
	return in, nil
}

// planStarts appends to starts the nodes the node ops of change set cs's
// ops touch, each followed by the parents its create and move ops give
// it: where a chain starts from for cs.
func planStarts(cs id.ID, ops []op.Op, starts []id.ID) []id.ID {
	touched, parents := planTouched(ops)
	for _, t := range touched {
		starts = append(starts, t)
		starts = append(starts, parents[t]...)
	}
	return starts
}

// planServed is the intents the live serves links among ls point to.
func planServed(ls []gmodel.Link) []id.ID {
	var out []id.ID
	for _, l := range ls {
		if l.Type == op.LinkServes && !l.Retired {
			out = append(out, l.To)
		}
	}
	return out
}

// planWalk climbs the view's tree from the nodes it is handed while the
// in-flight entries are read, and reads the blocks links to the nodes it
// finds. It reads the parent of each node pending, a level per query
// (planParents), queues the parents it has not seen, and probes the nodes
// found for blocks links, until nothing is pending and it has been told no
// more starts will come. While the entries hold the other connection it
// does both on one goroutine, each query bounded by index.BatchOne so that
// the index opens no third connection; once told, a second goroutine
// probes while the first climbs, each query still bounded, so the two
// connections stay busy without a third. What it read is kept for the
// gather: each node's parent, the nodes the view lacks, the nodes probed
// and the blocks links to them.
type planWalk struct {
	r    *reader
	done chan struct{}

	mu sync.Mutex
	// more is signalled when pending grows or ended is set, probes when
	// toProbe grows or climbed is set.
	more, probes sync.Cond
	// pending is the nodes to read the parents of, and seen those handed
	// to the walk so far; toProbe the nodes found and not probed yet.
	pending, toProbe []id.ID
	seen             map[id.ID]bool
	// ended is set once no more starts will come, climbed once nothing is
	// pending after that.
	ended, climbed bool
	// parents is each node found and its parent, missing the nodes the
	// view lacks, probed the nodes whose blocks links are in links. The
	// first two are the climbing goroutine's alone until done is closed;
	// the gather reads them all after.
	parents map[id.ID]id.ID
	missing map[id.ID]bool
	probed  map[id.ID]bool
	links   []gmodel.Link
	err     error
}

func newPlanWalk(r *reader) *planWalk {
	w := &planWalk{r: r, done: make(chan struct{}), seen: map[id.ID]bool{},
		parents: map[id.ID]id.ID{}, missing: map[id.ID]bool{}, probed: map[id.ID]bool{}}
	w.more.L, w.probes.L = &w.mu, &w.mu
	return w
}

// add hands the walk nodes to climb from; those it has seen are dropped.
func (w *planWalk) add(ids []id.ID) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, i := range ids {
		if !i.IsZero() && !w.seen[i] {
			w.seen[i] = true
			w.pending = append(w.pending, i)
		}
	}
	w.more.Signal()
}

// folded tells the walk that no more starts will come.
func (w *planWalk) folded() {
	w.mu.Lock()
	w.ended = true
	w.mu.Unlock()
	w.more.Signal()
}

// start runs the walk on a goroutine of its own, until done.
func (w *planWalk) start() {
	go w.run()
}

// run climbs until nothing is pending and the walk has ended, probing
// between levels until then, and on the probing goroutine from then on,
// which it waits for.
func (w *planWalk) run() {
	defer close(w.done)
	var probing sync.WaitGroup
	defer probing.Wait()
	alone := true
	for {
		w.mu.Lock()
		for len(w.pending) == 0 && !w.ended {
			w.more.Wait()
		}
		if w.ended && alone {
			alone = false
			probing.Go(w.probe)
		}
		if len(w.pending) == 0 || w.err != nil {
			w.climbed = true
			w.mu.Unlock()
			w.probes.Signal()
			return
		}
		level := w.take(&w.pending)
		w.mu.Unlock()
		if _, err := w.climb(level); err != nil {
			w.fail(err)
			continue
		}
		if alone {
			w.probeOnce()
		}
	}
}

// take removes and returns the first index.BatchOne of *ids at most;
// the caller holds mu.
func (w *planWalk) take(ids *[]id.ID) []id.ID {
	n := min(len(*ids), index.BatchOne)
	out := slices.Clone((*ids)[:n])
	*ids = slices.Delete(*ids, 0, n)
	return out
}

// probe reads the blocks links to the nodes found, a batch at a time,
// until nothing is left to probe and the climbing is done.
func (w *planWalk) probe() {
	for {
		w.mu.Lock()
		for len(w.toProbe) == 0 && !w.climbed {
			w.probes.Wait()
		}
		if len(w.toProbe) == 0 || w.err != nil {
			w.mu.Unlock()
			return
		}
		w.mu.Unlock()
		w.probeOnce()
	}
}

// probeOnce probes a batch of the nodes found, if any.
func (w *planWalk) probeOnce() {
	w.mu.Lock()
	batch := w.take(&w.toProbe)
	w.mu.Unlock()
	if len(batch) == 0 {
		return
	}
	ls, err := w.r.LinksTo(batch, op.LinkBlocks)
	if err != nil {
		w.fail(err)
		return
	}
	w.mu.Lock()
	w.links = append(w.links, ls...)
	for _, i := range batch {
		w.probed[i] = true
	}
	w.mu.Unlock()
}

// fail records the walk's first error, which ends it.
func (w *planWalk) fail(err error) {
	w.mu.Lock()
	if w.err == nil {
		w.err = err
	}
	w.mu.Unlock()
	w.more.Signal()
	w.probes.Signal()
}

// climb reads the parents of level, remembers them and the nodes the view
// lacks, queues the parents not seen yet and the nodes found for probing,
// and returns the nodes found.
func (w *planWalk) climb(level []id.ID) ([]id.ID, error) {
	ps, err := planParents(w.r, level)
	if err != nil {
		return nil, err
	}
	found := make([]id.ID, 0, len(ps))
	next := make([]id.ID, 0, len(ps))
	for _, p := range ps {
		w.parents[p.ID] = p.Parent
		found = append(found, p.ID)
		if !p.Parent.IsZero() {
			next = append(next, p.Parent)
		}
	}
	for _, i := range level {
		if _, ok := w.parents[i]; !ok {
			w.missing[i] = true
		}
	}
	w.add(next)
	w.mu.Lock()
	w.toProbe = append(w.toProbe, found...)
	w.mu.Unlock()
	w.probes.Signal()
	return found, nil
}

// wait returns once the walk is done, with its error if it failed. The
// walk ends only once told that no more starts will come (folded), which
// wait makes sure of, for a gather that returns early.
func (w *planWalk) wait() error {
	w.folded()
	<-w.done
	return w.err
}

// parent is node i's parent in the view, zero for a root; ok is false
// when the view lacks i. A node the walk did not climb from is read now,
// so the gather's answer does not depend on what the walk was handed.
func (w *planWalk) parent(i id.ID) (parent id.ID, ok bool, err error) {
	if p, ok := w.parents[i]; ok {
		return p, true, nil
	}
	if w.missing[i] {
		return id.ID{}, false, nil
	}
	if _, err := w.climb([]id.ID{i}); err != nil {
		return id.ID{}, false, err
	}
	p, ok := w.parents[i]
	return p, ok, nil
}

// blocks is the blocks links to the nodes in targets, after the walk is
// done: those probed already, and a probe now of any it did not reach.
func (w *planWalk) blocks(targets map[id.ID]bool) ([]gmodel.Link, error) {
	var left []id.ID
	for t := range targets {
		if !w.probed[t] {
			left = append(left, t)
		}
	}
	ls, err := w.r.LinksTo(left, op.LinkBlocks)
	if err != nil {
		return nil, err
	}
	var out []gmodel.Link
	for _, l := range slices.Concat(w.links, ls) {
		if targets[l.To] {
			out = append(out, l)
		}
	}
	return out, nil
}

// planParents is the ID and parent of each of ids that r's view has,
// sorted by ID, each once, as index.Index.Parents gives them; from the
// what-if model when r reads one.
func planParents(r *reader, ids []id.ID) ([]index.Parent, error) {
	if r.m == nil {
		ps, err := r.x.Parents(r.view.stored, ids)
		return ps, readError(err)
	}
	var out []index.Parent
	for _, i := range planSortIDs(slices.Clone(ids)) {
		if n, ok := r.m.Node(i); ok {
			out = append(out, index.Parent{ID: i, Parent: n.Parent})
		}
	}
	return out, nil
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

// planFold folds each of css from its entries, which the index decoded
// from canonical bytes (changeset.FoldCanonical). Folding is CPU work over
// entries already read: one worker per usable CPU takes the change sets in
// turn, which costs fewer goroutines than one per change set on few CPUs.
func planFold(css []id.ID, entries map[id.ID][]entry.Entry) []changeset.ChangeSet {
	out := make([]changeset.ChangeSet, len(css))
	var next atomic.Int64
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), len(css)) {
		wg.Go(func() {
			for k := int(next.Add(1)) - 1; k < len(css); k = int(next.Add(1)) - 1 {
				// Fold problems are dropped, as plan reads always dropped Fold's: these
				// entries come from the index's store, whose bodies entry.Decode checked
				// to be canonical at catch-up, and fold problems are reported by check's
				// fold rule and change show/ls (and stored by the index at catch-up),
				// not by plan reads.
				out[k], _ = changeset.FoldCanonical(css[k], entries[css[k]])
			}
		})
	}
	wg.Wait()
	return out
}

// planMayRest reports whether an op of the flying change sets may rest on
// a dropped one (changeset.RestsOn): whether one of the IDs it names
// (changeset.Names) is left over once the IDs no dropped change set can
// take with it are set aside. Those are the flying change sets' op IDs
// and what they create, the change sets in all, and, in as-is whatever
// r's view, the nodes and links (a landed create keeps them) and their
// current versions (landed ops). Node, link, op and change set IDs are
// unique in the log, so none of those is a dropped change set's op ID,
// and what it creates is kept. As-is is asked only about the IDs left
// over: their nodes in one query, then one query per link a flying link
// op retires that is still left over.
func planMayRest(r *reader, all []index.ChangeSetInfo, flying []changeset.ChangeSet) (bool, error) {
	known := map[id.ID]bool{}
	for _, c := range all {
		known[c.ID] = true
	}
	for _, s := range flying {
		for _, o := range s.Ops {
			known[o.ID] = true
			if o.Type == op.TypeNodeCreate || o.Type == op.TypeLinkCreate {
				known[o.Target()] = true
			}
		}
	}
	// left is the named IDs not known yet.
	left := func() []id.ID {
		var out []id.ID
		for _, s := range flying {
			for _, o := range s.Ops {
				for _, n := range changeset.Names(o) {
					if !known[n] {
						out = append(out, n)
					}
				}
			}
		}
		return planSortIDs(out)
	}
	named := left()
	if len(named) == 0 {
		return false, nil
	}
	ns, err := r.x.NodesLight(index.AsIs, named)
	if err != nil {
		return false, readError(err)
	}
	for _, n := range ns {
		known[n.ID], known[n.Version] = true, true
	}
	for _, s := range flying {
		for _, o := range s.Ops {
			if o.Type != op.TypeLinkRetire || known[o.Link] {
				continue
			}
			l, ok, err := r.x.Link(index.AsIs, o.Link)
			if err != nil {
				return false, readError(err)
			}
			if !ok {
				return true, nil
			}
			known[l.ID], known[l.Version] = true, true
		}
	}
	return len(left()) > 0, nil
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
		an := changeset.Analyze(sets)
		for _, c := range from {
			for _, d := range an.DependsOn(c) {
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
			// Fold problems are dropped, as plan reads always dropped Fold's: these
			// entries come from the index's store, whose bodies entry.Decode checked
			// to be canonical at catch-up, and fold problems are reported by check's
			// fold rule and change show/ls (and stored by the index at catch-up),
			// not by plan reads.
			s, _ := changeset.FoldCanonical(c, entries[c])
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

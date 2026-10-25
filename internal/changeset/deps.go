package changeset

// Dependencies between change sets. Change set S depends on change set T
// when a depends-on link runs from S to T. Such a link counts in two places
// only: in the landed ops, where it is live unless a landed link.retire or
// S's own ops retire it; and in S's own current ops, unless S's own ops
// retire it. A link written inside any other change set that has not landed
// does not count, so no change set can add to another's dependencies. A
// dropped op is gone from S's ops and counts for nothing.
//
// A write in S sees as-is, plus every change set S depends on, directly or
// through others, that is still in flight (draft, proposed or accepted),
// plus S itself. The walk does not enter a dropped change set (withdrawn,
// rejected or superseded) or a landed one, whose ops are already in as-is.
//
// Depends-on is acyclic: a link that would close a loop is refused, and a
// loop already in the log is reported. Loops are found over the
// dependencies of every change set, whatever its state. Every answer here
// is sorted by ID and depends only on the set of change sets, never on the
// order they are given in.

import (
	"slices"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// Dropped reports whether s is withdrawn, rejected or superseded.
func Dropped(s State) bool {
	return s == StateWithdrawn || s == StateRejected || s == StateSuperseded
}

// DependsOn returns the change sets cs depends on directly, as described
// above, sorted by ID. cs itself is never in it. An ID may name a change
// set that is not in sets.
func DependsOn(sets []ChangeSet, cs id.ID) []id.ID {
	var out []id.ID
	for _, to := range dependsOn(landedLinks(sets), cs, opsByID(sets)[cs]) {
		if to != cs {
			out = append(out, to)
		}
	}
	return out
}

// Closure returns the change sets cs depends on, directly or through
// others, that are in flight, sorted by ID. It excludes cs, skips IDs not
// in sets, and does not enter dropped or landed change sets.
func Closure(sets []ChangeSet, cs id.ID) []id.ID {
	states := make(map[id.ID]State, len(sets))
	for _, s := range sets {
		states[s.ID] = s.State
	}
	landed, ops := landedLinks(sets), opsByID(sets)
	seen := map[id.ID]bool{cs: true}
	var out []id.ID
	next := []id.ID{cs}
	for len(next) > 0 {
		at := next[len(next)-1]
		next = next[:len(next)-1]
		for _, to := range dependsOn(landed, at, ops[at]) {
			if seen[to] {
				continue
			}
			seen[to] = true
			if s, ok := states[to]; ok && s.inFlight() {
				out = append(out, to)
				next = append(next, to)
			}
		}
	}
	slices.SortFunc(out, id.Compare)
	return out
}

// WouldCycle reports whether a new link "from depends-on to" would close a
// loop: from and to are the same, or to already reaches from over the
// dependencies of every change set in sets.
func WouldCycle(sets []ChangeSet, from, to id.ID) bool {
	if from == to {
		return true
	}
	g := graph(sets)
	seen := map[id.ID]bool{to: true}
	next := []id.ID{to}
	for len(next) > 0 {
		at := next[len(next)-1]
		next = next[:len(next)-1]
		for _, n := range g[at] {
			if n == from {
				return true
			}
			if !seen[n] {
				seen[n] = true
				next = append(next, n)
			}
		}
	}
	return false
}

// Cycles returns every depends-on loop among the change sets in sets: each
// group of two or more that reach one another, and each change set that
// depends on itself. Each loop is sorted by ID and the loops by their first
// ID. It is nil when there are none.
func Cycles(sets []ChangeSet) [][]id.ID {
	g := graph(sets)
	nodes := make([]id.ID, 0, len(g))
	for n := range g {
		nodes = append(nodes, n)
	}
	slices.SortFunc(nodes, id.Compare)

	// Tarjan's strongly connected components, visiting nodes and edges in
	// ID order.
	index := map[id.ID]int{}
	low := map[id.ID]int{}
	onStack := map[id.ID]bool{}
	var stack []id.ID
	var out [][]id.ID
	var visit func(v id.ID)
	visit = func(v id.ID) {
		index[v] = len(index)
		low[v] = index[v]
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range g[v] {
			if _, ok := index[w]; !ok {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] != index[v] {
			return
		}
		var scc []id.ID
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		if len(scc) > 1 || slices.Contains(g[v], v) {
			slices.SortFunc(scc, id.Compare)
			out = append(out, scc)
		}
	}
	for _, n := range nodes {
		if _, ok := index[n]; !ok {
			visit(n)
		}
	}
	slices.SortFunc(out, func(a, b []id.ID) int { return id.Compare(a[0], b[0]) })
	return out
}

// graph maps each change set in sets to its dependencies that are in sets,
// self-loops included, each sorted by ID.
func graph(sets []ChangeSet) map[id.ID][]id.ID {
	landed, ops := landedLinks(sets), opsByID(sets)
	g := make(map[id.ID][]id.ID, len(ops))
	for c := range ops {
		g[c] = nil
		for _, to := range dependsOn(landed, c, ops[c]) {
			if _, ok := ops[to]; ok {
				g[c] = append(g[c], to)
			}
		}
	}
	return g
}

// links are depends-on link.creates by link ID, and the link IDs retired.
type links struct {
	created map[id.ID]op.Op
	retired map[id.ID]bool
}

func (ls *links) add(o op.Op) {
	switch {
	case o.Type == op.TypeLinkCreate && o.LinkType == op.LinkDependsOn:
		ls.created[o.Link] = o
	case o.Type == op.TypeLinkRetire:
		ls.retired[o.Link] = true
	}
}

func newLinks() links {
	return links{created: map[id.ID]op.Op{}, retired: map[id.ID]bool{}}
}

// landedLinks gathers the depends-on links and retires of every landed
// change set in sets.
func landedLinks(sets []ChangeSet) links {
	ls := newLinks()
	for _, s := range sets {
		if s.State == StateLanded {
			for _, o := range s.Ops {
				ls.add(o)
			}
		}
	}
	return ls
}

// opsByID maps each change set in sets to its ops.
func opsByID(sets []ChangeSet) map[id.ID][]op.Op {
	m := make(map[id.ID][]op.Op, len(sets))
	for _, s := range sets {
		m[s.ID] = append(m[s.ID], s.Ops...)
	}
	return m
}

// dependsOn is the dependencies of cs, whose own ops are ops, over landed,
// the landed links: those live in landed and not retired by ops, and those
// in ops not retired there. Sorted by ID, each once, cs included if it
// depends on itself.
func dependsOn(landed links, cs id.ID, ops []op.Op) []id.ID {
	own := newLinks()
	for _, o := range ops {
		own.add(o)
	}
	var out []id.ID
	for _, from := range []links{landed, own} {
		for l, o := range from.created {
			if o.From == cs && !landed.retired[l] && !own.retired[l] {
				out = append(out, o.To)
			}
		}
	}
	slices.SortFunc(out, id.Compare)
	return slices.Compact(out)
}

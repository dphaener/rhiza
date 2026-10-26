package model

// The order between change sets. A view applies change sets one after
// another, and a change set must come after the change sets it builds on.
// Change set S comes after T when an operation of S names as its base an
// operation of T, or names as its target, parent, endpoint (from, to) or
// decision a node or link that T creates. These are the only edges: no
// timestamps, no land entries, no states. When an op ID or a created ID
// belongs to several change sets, the smallest change set ID owns it.
// Among change sets that are free to go, the smallest ID goes first. When
// none is free but some remain, the remaining ones are in or behind a cycle:
// each is reported (ErrCycle) once and the one with the smallest ID goes
// next, as if its edges were not there. The order depends only on the set of
// change sets, never on the order they are given in.

import (
	"container/heap"
	"fmt"
	"slices"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// order returns sets in apply order, as described above, with an ErrCycle
// problem for every change set that had to be released from a cycle or from
// behind one. It does not change sets. Problems are sorted by change set and
// are nil when there are none.
func order(sets []changeset.ChangeSet) ([]changeset.ChangeSet, []Problem) {
	if len(sets) == 0 {
		return nil, nil
	}
	// Index i is the i-th change set by ID, so smaller index is smaller ID.
	sorted := slices.Clone(sets)
	slices.SortStableFunc(sorted, func(a, b changeset.ChangeSet) int { return id.Compare(a.ID, b.ID) })

	// Ownership: the first (smallest) change set to hold an op or create an
	// ID owns it.
	ops := make(map[id.ID]int)
	created := make(map[id.ID]int)
	for i, s := range sorted {
		for _, o := range s.Ops {
			if _, ok := ops[o.ID]; !ok {
				ops[o.ID] = i
			}
			if o.Type == op.TypeNodeCreate || o.Type == op.TypeLinkCreate {
				if _, ok := created[o.Target()]; !ok {
					created[o.Target()] = i
				}
			}
		}
	}

	// Edges from owner to dependent, each once.
	after := make([][]int, len(sorted))
	indeg := make([]int, len(sorted))
	for i, s := range sorted {
		seen := map[int]bool{}
		edge := func(owners map[id.ID]int, x id.ID) {
			if x.IsZero() {
				return
			}
			t, ok := owners[x]
			if !ok || t == i || seen[t] {
				return
			}
			seen[t] = true
			after[t] = append(after[t], i)
			indeg[i]++
		}
		for _, o := range s.Ops {
			edge(ops, o.Base)
			edge(created, o.Target())
			edge(created, o.Parent)
			edge(created, o.From)
			edge(created, o.To)
			edge(created, o.Decision)
		}
	}

	// Kahn's algorithm, smallest index first.
	ready := &intHeap{}
	for i, d := range indeg {
		if d == 0 {
			heap.Push(ready, i)
		}
	}
	out := make([]changeset.ChangeSet, 0, len(sorted))
	done := make([]bool, len(sorted))
	reported := make([]bool, len(sorted))
	var problems []Problem
	next := 0 // every index below next is done
	for len(out) < len(sorted) {
		if ready.Len() == 0 {
			for next < len(sorted) && done[next] {
				next++
			}
			for i := next; i < len(sorted); i++ {
				if !done[i] && !reported[i] {
					reported[i] = true
					problems = append(problems, Problem{ChangeSet: sorted[i].ID,
						Err: fmt.Errorf("%w: change set is in or behind a cycle of change sets; ordered by ID", ErrCycle)})
				}
			}
			// Released: its remaining edges no longer hold it back, and
			// they can never bring it to zero again.
			indeg[next] = -1
			heap.Push(ready, next)
		}
		i := heap.Pop(ready).(int)
		done[i] = true
		out = append(out, sorted[i])
		for _, j := range after[i] {
			indeg[j]--
			if indeg[j] == 0 {
				heap.Push(ready, j)
			}
		}
	}
	return out, sortProblems(problems)
}

// intHeap is a min-heap of ints for container/heap.
type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

package model

import (
	"cmp"
	"maps"
	"slices"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// aspectLinked marks a node that an operation links from or to. It overlaps
// only with a retirement of the node, never with another link.
const aspectLinked = "linked"

// Overlap is two in-flight change sets writing overlapping aspects of the
// same target, whatever their bases.
type Overlap struct {
	// A and B are the change sets, A sorting before B.
	A, B id.ID
	// Target is the node or link both write.
	Target id.ID
	// Aspects are the aspects of A's and B's operations on Target that take
	// part in the overlap, sorted and unique: those both write, and every
	// aspect of two operations when either is a retirement (a link created
	// to or from the node counts as the aspect "linked").
	Aspects []string
}

// Overlaps returns every overlap between the proposed and accepted change
// sets in sets: one per pair of change sets and target, sorted by (A, B,
// Target), nil when there are none. Drafts and terminal change sets are
// ignored, except the drafts named in drafts, which count as in flight. No model is computed and bases are not compared: two proposals
// writing the same field overlap whatever their bases. The result does not
// depend on the order of sets.
func Overlaps(sets []changeset.ChangeSet, drafts ...id.ID) []Overlap {
	// For each target, the aspects each change set's operations write, one
	// slice per operation, creates left out.
	byTarget := map[id.ID]map[id.ID][][]string{}
	for _, cs := range sets {
		named := cs.State == changeset.StateDraft && slices.Contains(drafts, cs.ID)
		if cs.State != changeset.StateProposed && cs.State != changeset.StateAccepted && !named {
			continue
		}
		for _, o := range cs.Ops {
			if o.Type == op.TypeLinkCreate {
				// A new link touches both endpoints: only a retirement
				// of one conflicts with it.
				for _, n := range []id.ID{o.From, o.To} {
					if byTarget[n] == nil {
						byTarget[n] = map[id.ID][][]string{}
					}
					byTarget[n][cs.ID] = append(byTarget[n][cs.ID], []string{aspectLinked})
				}
				continue
			}
			as := aspects(o)
			if as == nil {
				continue
			}
			t := o.Target()
			if byTarget[t] == nil {
				byTarget[t] = map[id.ID][][]string{}
			}
			byTarget[t][cs.ID] = append(byTarget[t][cs.ID], as)
		}
	}
	var out []Overlap
	for t, writers := range byTarget {
		if len(writers) < 2 {
			continue
		}
		ids := slices.SortedFunc(maps.Keys(writers), id.Compare)
		for i, a := range ids {
			for _, b := range ids[i+1:] {
				if as := overlapAspects(writers[a], writers[b]); as != nil {
					out = append(out, Overlap{A: a, B: b, Target: t, Aspects: as})
				}
			}
		}
	}
	slices.SortFunc(out, func(x, y Overlap) int {
		return cmp.Or(id.Compare(x.A, y.A), id.Compare(x.B, y.B), id.Compare(x.Target, y.Target))
	})
	return out
}

// overlapAspects returns the aspects taking part in the overlaps between
// operations writing as and operations writing bs, sorted and unique, nil
// when none overlap.
func overlapAspects(as, bs [][]string) []string {
	var out []string
	for _, a := range as {
		for _, b := range bs {
			switch {
			case !overlaps(a, b):
			case slices.Contains(a, aspectRetired) || slices.Contains(b, aspectRetired):
				out = append(append(out, a...), b...)
			default:
				for _, s := range a {
					if s != aspectLinked && slices.Contains(b, s) {
						out = append(out, s)
					}
				}
			}
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

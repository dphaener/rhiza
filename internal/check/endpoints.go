package check

import "github.com/dphaener/rhiza/internal/op"

// End is what one end of a link is: a node of some kind, or a change set.
// A node's End is its kind's name (NodeEnd); a change set's is ChangeSet.
type End string

// ChangeSet is the End of a link endpoint that names a change set.
const ChangeSet End = "change set"

// NodeEnd is the End of a node of kind k.
func NodeEnd(k op.Kind) End { return End(k) }

// Ends returns every End: the ten node kinds in op.Kinds() order, then
// ChangeSet.
func Ends() []End {
	var out []End
	for _, k := range op.Kinds() {
		out = append(out, NodeEnd(k))
	}
	return append(out, ChangeSet)
}

// endSet is a set of ends one side of a row allows.
type endSet []End

func nodes(ks ...op.Kind) endSet {
	var out endSet
	for _, k := range ks {
		out = append(out, NodeEnd(k))
	}
	return out
}

var (
	anyNode  = nodes(op.Kinds()...)
	anything = append(nodes(op.Kinds()...), ChangeSet)
	change   = endSet{ChangeSet}
)

// row allows every pair of an end in from with an end in to.
type row struct {
	from, to endSet
}

// rows is the endpoint table: for each link type, the rows whose pairs it
// allows. A type with several rows allows the union of their pairs, never
// the cross product of their sides.
var rows = map[op.LinkType][]row{
	op.LinkServes:         {{append(nodes(op.KindCapability), ChangeSet), nodes(op.KindIntent)}},
	op.LinkConstrains:     {{nodes(op.KindConstraint), nodes(op.KindCapability, op.KindIntent)}},
	op.LinkExempts:        {{nodes(op.KindCapability), nodes(op.KindConstraint)}},
	op.LinkDependsOn:      {{nodes(op.KindCapability), nodes(op.KindCapability)}, {change, change}},
	op.LinkRelatesTo:      {{anything, anything}},
	op.LinkDecidedBy:      {{anything, nodes(op.KindDecision)}},
	op.LinkAnswers:        {{nodes(op.KindDecision), nodes(op.KindQuestion)}},
	op.LinkBlocks:         {{nodes(op.KindQuestion), append(nodes(op.KindCapability), ChangeSet)}},
	op.LinkSupersedes:     sameKind(anyNode),
	op.LinkDerivedFrom:    {{anything, nodes(op.KindSource)}},
	op.LinkDiscoveredFrom: {{change, change}},
	op.LinkImplements:     {{nodes(op.KindAnchor), nodes(op.KindCapability)}},
	op.LinkVerifies:       {{nodes(op.KindEvidence), nodes(op.KindCapability)}},
	op.LinkObserves:       {{nodes(op.KindEvidence), nodes(op.KindAnchor)}},
	op.LinkDiscusses:      {{nodes(op.KindNote), anything}},
}

// sameKind is one row per end in s, from it to itself.
func sameKind(s endSet) []row {
	var out []row
	for _, e := range s {
		out = append(out, row{endSet{e}, endSet{e}})
	}
	return out
}

type pair struct {
	t        op.LinkType
	from, to End
}

// allowed is rows flattened into the set of allowed (type, from, to).
var allowed = func() map[pair]bool {
	m := map[pair]bool{}
	for t, rs := range rows {
		for _, r := range rs {
			for _, f := range r.from {
				for _, to := range r.to {
					m[pair{t, f, to}] = true
				}
			}
		}
	}
	return m
}()

// Allowed reports whether a link of type t may go from an end of kind from
// to an end of kind to. An unknown type or End is never allowed.
func Allowed(t op.LinkType, from, to End) bool {
	return allowed[pair{t, from, to}]
}

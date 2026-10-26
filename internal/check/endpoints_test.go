package check

import (
	"testing"

	"github.com/dphaener/rhiza/internal/op"
)

// expect is the endpoint table written out a second way, as a switch, so
// TestAllowedExhaustive compares two independent spellings of D10.
func expect(t op.LinkType, from, to End) bool {
	isNode := func(e End) bool { return e != ChangeSet }
	is := func(e End, ks ...op.Kind) bool {
		for _, k := range ks {
			if e == End(k) {
				return true
			}
		}
		return false
	}
	switch t {
	case op.LinkServes:
		return (is(from, op.KindCapability) || from == ChangeSet) && is(to, op.KindIntent)
	case op.LinkConstrains:
		return is(from, op.KindConstraint) && is(to, op.KindCapability, op.KindIntent)
	case op.LinkExempts:
		return is(from, op.KindCapability) && is(to, op.KindConstraint)
	case op.LinkDependsOn:
		return (is(from, op.KindCapability) && is(to, op.KindCapability)) || (from == ChangeSet && to == ChangeSet)
	case op.LinkRelatesTo:
		return true
	case op.LinkDecidedBy:
		return is(to, op.KindDecision)
	case op.LinkAnswers:
		return is(from, op.KindDecision) && is(to, op.KindQuestion)
	case op.LinkBlocks:
		return is(from, op.KindQuestion) && (is(to, op.KindCapability) || to == ChangeSet)
	case op.LinkSupersedes:
		return isNode(from) && from == to
	case op.LinkDerivedFrom:
		return is(to, op.KindSource)
	case op.LinkDiscoveredFrom:
		return from == ChangeSet && to == ChangeSet
	case op.LinkImplements:
		return is(from, op.KindAnchor) && is(to, op.KindCapability)
	case op.LinkVerifies:
		return is(from, op.KindEvidence) && is(to, op.KindCapability)
	case op.LinkObserves:
		return is(from, op.KindEvidence) && is(to, op.KindAnchor)
	}
	return false
}

func TestEnds(t *testing.T) {
	es := Ends()
	if len(es) != len(op.Kinds())+1 || es[len(es)-1] != ChangeSet {
		t.Fatalf("Ends() = %v", es)
	}
}

func TestAllowedExhaustive(t *testing.T) {
	n := 0
	for _, lt := range op.LinkTypes() {
		for _, f := range Ends() {
			for _, to := range Ends() {
				got, want := Allowed(lt, f, to), expect(lt, f, to)
				if got != want {
					t.Errorf("Allowed(%s, %s, %s) = %v, want %v", lt, f, to, got, want)
				}
				if got {
					n++
				}
			}
		}
	}
	if n == 0 {
		t.Fatal("nothing allowed")
	}
	if len(rows) != len(op.LinkTypes()) {
		t.Errorf("table has %d link types, want %d", len(rows), len(op.LinkTypes()))
	}
}

// TestAllowedRows checks each row of D10 and its nearest disallowed
// neighbours by name.
func TestAllowedRows(t *testing.T) {
	cap, intent, cons := NodeEnd(op.KindCapability), NodeEnd(op.KindIntent), NodeEnd(op.KindConstraint)
	dec, q, ev := NodeEnd(op.KindDecision), NodeEnd(op.KindQuestion), NodeEnd(op.KindEvidence)
	src, actor, anchor := NodeEnd(op.KindSource), NodeEnd(op.KindActor), NodeEnd(op.KindAnchor)
	cs := ChangeSet
	cases := []struct {
		name     string
		t        op.LinkType
		from, to End
		want     bool
	}{
		{"serves capability to intent", op.LinkServes, cap, intent, true},
		{"serves change set to intent", op.LinkServes, cs, intent, true},
		{"serves intent to intent", op.LinkServes, intent, intent, false},
		{"serves capability to capability", op.LinkServes, cap, cap, false},
		{"serves capability to change set", op.LinkServes, cap, cs, false},
		{"constrains constraint to capability", op.LinkConstrains, cons, cap, true},
		{"constrains constraint to intent", op.LinkConstrains, cons, intent, true},
		{"constrains constraint to decision", op.LinkConstrains, cons, dec, false},
		{"constrains capability to capability", op.LinkConstrains, cap, cap, false},
		{"exempts capability to constraint", op.LinkExempts, cap, cons, true},
		{"exempts intent to constraint", op.LinkExempts, intent, cons, false},
		{"exempts capability to capability", op.LinkExempts, cap, cap, false},
		{"depends-on capability to capability", op.LinkDependsOn, cap, cap, true},
		{"depends-on change set to change set", op.LinkDependsOn, cs, cs, true},
		{"depends-on capability to change set", op.LinkDependsOn, cap, cs, false},
		{"depends-on change set to capability", op.LinkDependsOn, cs, cap, false},
		{"depends-on intent to intent", op.LinkDependsOn, intent, intent, false},
		{"relates-to actor to change set", op.LinkRelatesTo, actor, cs, true},
		{"relates-to change set to anchor", op.LinkRelatesTo, cs, anchor, true},
		{"decided-by change set to decision", op.LinkDecidedBy, cs, dec, true},
		{"decided-by capability to decision", op.LinkDecidedBy, cap, dec, true},
		{"decided-by capability to question", op.LinkDecidedBy, cap, q, false},
		{"answers decision to question", op.LinkAnswers, dec, q, true},
		{"answers question to decision", op.LinkAnswers, q, dec, false},
		{"answers evidence to question", op.LinkAnswers, ev, q, false},
		{"blocks question to capability", op.LinkBlocks, q, cap, true},
		{"blocks question to change set", op.LinkBlocks, q, cs, true},
		{"blocks question to intent", op.LinkBlocks, q, intent, false},
		{"blocks decision to capability", op.LinkBlocks, dec, cap, false},
		{"supersedes decision to decision", op.LinkSupersedes, dec, dec, true},
		{"supersedes source to source", op.LinkSupersedes, src, src, true},
		{"supersedes across kinds", op.LinkSupersedes, dec, cap, false},
		{"supersedes change set to change set", op.LinkSupersedes, cs, cs, false},
		{"derived-from change set to source", op.LinkDerivedFrom, cs, src, true},
		{"derived-from capability to source", op.LinkDerivedFrom, cap, src, true},
		{"derived-from capability to evidence", op.LinkDerivedFrom, cap, ev, false},
		{"discovered-from change set to change set", op.LinkDiscoveredFrom, cs, cs, true},
		{"discovered-from capability to change set", op.LinkDiscoveredFrom, cap, cs, false},
		{"implements anchor to capability", op.LinkImplements, anchor, cap, true},
		{"implements anchor to intent", op.LinkImplements, anchor, intent, false},
		{"implements evidence to capability", op.LinkImplements, ev, cap, false},
		{"verifies evidence to capability", op.LinkVerifies, ev, cap, true},
		{"verifies from an intent", op.LinkVerifies, intent, cap, false},
		{"verifies evidence to intent", op.LinkVerifies, ev, intent, false},
		{"observes evidence to anchor", op.LinkObserves, ev, anchor, true},
		{"observes anchor to anchor", op.LinkObserves, anchor, anchor, false},
		{"observes evidence to capability", op.LinkObserves, ev, cap, false},
		{"unknown link type", op.LinkType("likes"), cap, cap, false},
		{"unknown end", op.LinkRelatesTo, End("widget"), cap, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Allowed(c.t, c.from, c.to); got != c.want {
				t.Errorf("Allowed(%s, %s, %s) = %v, want %v", c.t, c.from, c.to, got, c.want)
			}
		})
	}
}

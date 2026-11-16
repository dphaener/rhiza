package check

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// bogusRef is a full ID that names no node in the fixtures.
const bogusRef = "01ZZZZZZZZZZZZZZZZZZZZZZZZ"

// whatIfFindings is the findings of the rules that read what-if views, in
// report order, as severity, rule, change set, op, link and message.
func whatIfFindings(rep Report) []string {
	out := []string{}
	for _, f := range rep.Findings {
		if f.Rule != RuleInFlight && f.Rule != RuleTextRef {
			continue
		}
		out = append(out, strings.Join([]string{f.Severity, f.Rule, f.ChangeSet, f.Op, f.Link, f.Message}, " | "))
	}
	return out
}

// TestWhatIfOnceMatches checks the in-flight findings of the passes that
// share the what-if views check builds once from the as-is seed: the
// in-flight problems, the in-flight links and the text refs.
func TestWhatIfOnceMatches(t *testing.T) {
	t.Run("Listed", func(t *testing.T) {
		// A drafts a node and B depends on A and updates it, cleanly; C
		// updates the capability on a base a later landed record moved; D
		// reuses the capability's handle; E links an intent as serving a
		// capability; F holds a reference to no node.
		r := newRepo(t)
		intent, capab, capOp := cleanOps(r)
		a := r.open(changeset.KindIntention, "A")
		n, a1 := r.next(), r.next()
		a.then(changeset.Ops{Ops: []op.Op{mkNode(a1, n, op.KindDecision, "pick")}})
		b := r.open(changeset.KindIntention, "B")
		b.then(changeset.Ops{Ops: []op.Op{dependsOn(r, b.id, a.id), setField(r.next(), n, a1, "two")}})
		c := r.open(changeset.KindIntention, "C")
		c1 := r.next()
		c.then(changeset.Ops{Ops: []op.Op{setField(c1, capab, capOp, "stale")}})
		m1 := r.next()
		move := r.record("Move", setField(m1, capab, capOp, "moved"))
		d := r.open(changeset.KindIntention, "D")
		d1, dup := r.next(), r.next()
		d.then(changeset.Ops{Ops: []op.Op{mkNode(d1, dup, op.KindCapability, "thing")}})
		e := r.open(changeset.KindIntention, "E")
		e1, el := r.next(), r.next()
		e.then(changeset.Ops{Ops: []op.Op{mkLink(e1, el, op.LinkServes, intent, capab)}})
		f := r.open(changeset.KindIntention, "F")
		f.then(changeset.Ops{Ops: []op.Op{textNode(r, r.next(), "holder", "see @{"+bogusRef+"}")}})

		rep := Run(r.input())
		want := []string{
			fmt.Sprintf("warn | in-flight | %s | %s |  | model: conflict: base %s of %s is stale: op %s of change set %s since wrote field:statement",
				c.id, c1, capOp, capab, m1, move.id),
			fmt.Sprintf("warn | in-flight | %s | %s |  | model: duplicate handle: node %s has handle \"thing\", already held by node %s",
				d.id, d1, dup, capab),
			fmt.Sprintf("warn | in-flight | %s | %s | %s | link type \"serves\" does not go from intent to capability", e.id, e1, el),
			fmt.Sprintf("fail | text-ref | %s |  |  | holder field statement holds @{%s}, which names no node", f.id, bogusRef),
		}
		got := whatIfFindings(rep)
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Errorf("in-flight findings:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
		if len(rep.Findings) != len(want) {
			t.Errorf("findings %v, want only %d what-if ones", short(rep.Findings), len(want))
		}
	})
	t.Run("Once", func(t *testing.T) {
		// A has a stale op, a bad link and a broken text ref; B and C depend
		// on A, so their what-if views hold all three too. Each is reported
		// once, on A.
		r := newRepo(t)
		intent, capab, capOp := cleanOps(r)
		a := r.open(changeset.KindIntention, "A")
		a1, a2, al := r.next(), r.next(), r.next()
		a.then(changeset.Ops{Ops: []op.Op{
			setField(a1, capab, capOp, "stale"),
			mkLink(a2, al, op.LinkServes, intent, capab),
			textNode(r, r.next(), "holder", "see @{"+bogusRef+"}"),
		}})
		r.record("Move", setField(r.next(), capab, capOp, "moved"))
		var deps []id.ID
		for _, name := range []string{"B", "C"} {
			s := r.open(changeset.KindIntention, name)
			s.then(changeset.Ops{Ops: []op.Op{dependsOn(r, s.id, a.id)}})
			deps = append(deps, s.id)
		}
		got := whatIfFindings(Run(r.input()))
		if len(got) != 3 {
			t.Fatalf("in-flight findings:\n%s\nwant 3, one each", strings.Join(got, "\n"))
		}
		for i, part := range []string{"model: conflict: base", "link type \"serves\"", "which names no node"} {
			if !strings.Contains(got[i], part) || !strings.Contains(got[i], " | "+a.id.String()+" | ") {
				t.Errorf("finding %d %q: want %q on %s", i, got[i], part, a.id)
			}
			for _, d := range deps {
				if strings.Contains(got[i], d.String()) {
					t.Errorf("finding %d %q names dependent %s", i, got[i], d)
				}
			}
		}
	})
}

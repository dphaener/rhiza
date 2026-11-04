package check

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// textNode is a decision whose statement is text.
func textNode(r *repo, node id.ID, handle, text string) op.Op {
	b, _ := json.Marshal(text)
	return op.Op{ID: r.next(), Type: op.TypeNodeCreate, Node: node, Kind: op.KindDecision, Handle: handle,
		Fields: map[string]json.RawMessage{"statement": b}}
}

func textRefFindings(rep Report) []Finding {
	var out []Finding
	for _, f := range rep.Findings {
		if f.Rule == RuleTextRef {
			out = append(out, f)
		}
	}
	return out
}

func TestTextRefs(t *testing.T) {
	// a 26-character ID that names nothing in the fixtures.
	const bogus = "01ZZZZZZZZZZZZZZZZZZZZZZZZ"
	cases := []struct {
		name string
		text func(live, retired id.ID) string
		want []string // severities
		ok   bool
	}{
		{"dangling", func(_, _ id.ID) string { return "see @{" + bogus + "}" }, []string{"fail"}, false},
		{"retired", func(_, r id.ID) string { return "see @{" + r.String() + "}" }, []string{"warn"}, true},
		{"live", func(l, _ id.ID) string { return "see @{" + l.String() + "}" }, nil, true},
		{"plain", func(_, _ id.ID) string { return "no refs, a@b.c" }, nil, true},
		{"look-alikes", func(l, _ id.ID) string {
			return "'@{ID}' '@{BOGUS}' '@{" + l.String() + "}' @@{" + bogus + "}"
		}, nil, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newRepo(t)
			clean(r)
			live, dead, deadOp := r.next(), r.next(), r.next()
			r.record("Targets", mkNode(r.next(), live, op.KindDecision, "live-one"),
				mkNode(deadOp, dead, op.KindDecision, "dead-one"))
			r.record("Retire", op.Op{ID: r.next(), Type: op.TypeNodeRetire, Node: dead, Base: deadOp, Reason: "x"})
			// "look-alikes" reads literally: the '@{<full ID>}' there is not preceded by a stored form.
			n := r.next()
			r.record("Holder", textNode(r, n, "holder", c.text(live, dead)))
			rep := Run(r.input())
			got := textRefFindings(rep)
			if len(got) != len(c.want) {
				t.Fatalf("findings %+v, want %v", got, c.want)
			}
			for i, f := range got {
				if f.Severity != c.want[i] {
					t.Errorf("severity %q, want %q", f.Severity, c.want[i])
				}
				if !strings.Contains(f.Message, "holder") || !strings.Contains(f.Message, "statement") {
					t.Errorf("message %q names no node or field", f.Message)
				}
			}
			if len(got) == 1 && !strings.Contains(got[0].Message, bogus) && !strings.Contains(got[0].Message, dead.String()) {
				t.Errorf("message %q names no ID", got[0].Message)
			}
			if c.ok != !hasFail(rep) {
				t.Errorf("OK %v, want %v: %+v", rep.OK, c.ok, rep.Findings)
			}
		})
	}
}

func hasFail(r Report) bool { return !r.OK }

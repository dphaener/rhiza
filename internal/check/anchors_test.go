package check

import (
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

const anchorSrc = "package p\n\nfunc Keep() {}\n"

// anchorOp creates anchor handle on path and symbol.
func anchorOp(r *repo, node id.ID, handle, path, symbol string) op.Op {
	x := mkNode(r.next(), node, op.KindAnchor, handle)
	x.Fields = map[string]json.RawMessage{}
	x.Fields["path"], _ = json.Marshal(path)
	if symbol != "" {
		x.Fields["symbol"], _ = json.Marshal(symbol)
	}
	return x
}

func TestAnchorRule(t *testing.T) {
	tree := fstest.MapFS{"p.go": {Data: []byte(anchorSrc)}, "notes.md": {Data: []byte("x")}}
	cases := []struct {
		name         string
		path, symbol string
		want         string // "" for no finding
	}{
		{"resolves", "p.go", "Keep", ""},
		{"file only", "notes.md", "", ""},
		{"missing path", "gone.go", "", "path gone.go is gone"},
		{"missing symbol", "p.go", "Gone", "symbol Gone of p.go is gone"},
		{"symbol on md", "notes.md", "T", "symbol T of notes.md is gone"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRepo(t)
			_, capab := clean(r)
			a := r.next()
			r.record("Anchor", anchorOp(r, a, "pin", tc.path, tc.symbol),
				mkLink(r.next(), r.next(), op.LinkImplements, a, capab))
			in := r.input()
			in.Tree = tree
			rep := Run(in)
			fs := ofRule(rep, RuleAnchor)
			if tc.want == "" {
				if len(fs) != 0 || !rep.OK {
					t.Fatalf("OK %v, findings %+v", rep.OK, rep.Findings)
				}
				return
			}
			if len(fs) != 1 || rep.OK || fs[0].Severity != SeverityFail {
				t.Fatalf("OK %v, findings %+v", rep.OK, rep.Findings)
			}
			m := fs[0].Message
			if !strings.Contains(m, tc.want) || !strings.Contains(m, "anchor pin") ||
				!strings.Contains(m, "implements thing") {
				t.Errorf("message %q", m)
			}
		})
	}
}

func TestAnchorRuleRetiredAndSkipped(t *testing.T) {
	r := newRepo(t)
	a := r.next()
	create := anchorOp(r, a, "pin", "gone.go", "")
	r.record("Anchor", create)
	in := r.input()
	if rep := Run(Input{Log: in.Log, Rendered: in.Rendered}); len(ofRule(rep, RuleAnchor)) != 0 {
		t.Errorf("rule ran without a tree: %+v", rep.Findings)
	}
	r.record("Retire", op.Op{ID: r.next(), Type: op.TypeNodeRetire, Node: a, Base: create.ID, Reason: "gone"})
	in = r.input()
	in.Tree = fstest.MapFS{}
	if rep := Run(in); len(ofRule(rep, RuleAnchor)) != 0 {
		t.Errorf("retired anchor found: %+v", rep.Findings)
	}
}

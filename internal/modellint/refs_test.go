package modellint

import (
	"encoding/json"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func nid(n byte) id.ID { return id.ID{15: n} }

func js(s string) json.RawMessage {
	b, _ := json.Marshal(s)
	return b
}

func capOp(n byte, handle, statement, acceptance string) op.Op {
	f := map[string]json.RawMessage{"statement": js(statement)}
	if acceptance != "" {
		f["acceptance"] = js(acceptance)
	}
	return op.Op{ID: nid(100 + n), Type: op.TypeNodeCreate, Node: nid(n), Kind: op.KindCapability, Handle: handle, Fields: f}
}

func head() fstest.MapFS {
	return fstest.MapFS{
		"internal/api/run.go":                  {Data: []byte("package api\n")},
		"internal/api/run_test.go":             {Data: []byte("package api\nfunc TestRun(t *testing.T) {}\n")},
		"cmd/rhiza/testdata/script/real.txtar": {Data: []byte("")},
		"cmd/rhiza/main_test.go":               {Data: []byte("func TestScripts(t *testing.T) {}\n")},
		"Makefile":                             {Data: []byte("")},
	}
}

// lint folds a landed set holding the other ops and an accepted set holding
// written, and lints the accepted one.
func lint(t *testing.T, other, written []op.Op, fsys fstest.MapFS) []Hit {
	t.Helper()
	sets := []changeset.ChangeSet{
		{ID: nid(200), State: changeset.StateLanded, Ops: other},
		{ID: nid(201), State: changeset.StateAccepted, Ops: written},
	}
	m, ps := model.ToBe(sets)
	if len(ps) != 0 {
		t.Fatalf("problems: %v", ps)
	}
	return StaleRefs(m, sets[1], fsys)
}

func refs(hs []Hit) []string {
	var out []string
	for _, h := range hs {
		out = append(out, h.Handle+"/"+h.Field+"/"+h.Kind+"/"+h.Ref)
	}
	return out
}

func TestStaleRefsPresent(t *testing.T) {
	hs := lint(t, nil, []op.Op{capOp(1, "ok",
		"Edits internal/api/run.go and ./internal/api/... and Makefile, and/or https://example.com/x.",
		"go test -trimpath ./internal/api/... -run TestRun passes; go test ./cmd/rhiza -run TestScripts/real$ passes (see internal/api/run.go:12).")}, head())
	if len(hs) != 0 {
		t.Fatalf("hits = %v, want none", refs(hs))
	}
}

func TestStaleRefsMissing(t *testing.T) {
	for _, tc := range []struct {
		name, field, text string
		want              []string
	}{
		{"path", "statement", "Edits internal/api/gone.go.", []string{"n/statement/path/internal/api/gone.go"}},
		{"path in acceptance", "acceptance", "see cmd/rhiza/missing.go:3", []string{"n/acceptance/path/cmd/rhiza/missing.go"}},
		{"test", "acceptance", "TestGone passes", []string{"n/acceptance/test/TestGone"}},
		{"script", "acceptance", "-run TestScripts/nope passes", []string{"n/acceptance/script/TestScripts/nope"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var st, acc string
			st = "writes things"
			if tc.field == "statement" {
				st = tc.text
			} else {
				acc = tc.text
			}
			// A second, present path keeps the capability from being unretired.
			st += " also internal/api/run.go"
			hs := lint(t, nil, []op.Op{capOp(1, "n", st, acc)}, head())
			if got := refs(hs); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("hits = %v, want %v", got, tc.want)
			}
			if hs[0].Node != nid(1) {
				t.Errorf("node = %v", hs[0].Node)
			}
		})
	}
}

func TestStaleRefsUnretired(t *testing.T) {
	hs := lint(t, nil, []op.Op{capOp(1, "n", "Lives in internal/old/a.go and internal/old/b.go.", "")}, head())
	want := []string{
		"n/statement/path/internal/old/a.go",
		"n/statement/path/internal/old/b.go",
		"n/statement/unretired-capability/",
	}
	if got := refs(hs); !reflect.DeepEqual(got, want) {
		t.Fatalf("hits = %v, want %v", got, want)
	}
}

func TestStaleRefsRetiredCapabilityNotReported(t *testing.T) {
	hs := lint(t, nil, []op.Op{
		capOp(1, "n", "Lives in internal/old/a.go.", ""),
		{ID: nid(150), Type: op.TypeNodeRetire, Node: nid(1), Base: nid(101), Reason: "gone"},
	}, head())
	if len(hs) != 0 {
		t.Fatalf("hits = %v, want none", refs(hs))
	}
}

func TestStaleRefsUpdateCounts(t *testing.T) {
	hs := lint(t,
		[]op.Op{capOp(1, "n", "Fine internal/api/run.go.", "")},
		[]op.Op{{ID: nid(150), Type: op.TypeNodeUpdate, Node: nid(1), Base: nid(101),
			Set: map[string]json.RawMessage{"acceptance": js("TestGone passes")}}}, head())
	if got, want := refs(hs), []string{"n/acceptance/test/TestGone"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("hits = %v, want %v", got, want)
	}
}

func TestStaleRefsOtherNodesNotRead(t *testing.T) {
	hs := lint(t,
		[]op.Op{capOp(1, "old", "Names internal/api/gone.go and TestGone.", "")},
		[]op.Op{capOp(2, "new", "Names internal/api/run.go.", "")}, head())
	if len(hs) != 0 {
		t.Fatalf("hits = %v, want none", refs(hs))
	}
}

func TestStaleRefsNotReferences(t *testing.T) {
	hs := lint(t, nil, []op.Op{capOp(1, "n",
		"Reads and/or writes; see https://example.com/internal/x, docs-site/page, @{x}/y and internal/api/run.go.", "")}, head())
	if len(hs) != 0 {
		t.Fatalf("hits = %v, want none", refs(hs))
	}
}

package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/op"
)

// treeLog lands two capability roots (alpha, beta) with nested children
// whose keys are out of creation order, an intent tree and a decision, then
// lands a second change set retiring beta, which keeps a live child.
func treeLog(t *testing.T) *Ctx {
	t.Helper()
	ctx := readCtx(t)
	newChange(t, ctx, `{"intent":"Base"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"beta"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"alpha","before":"beta"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"a-two","parent":"alpha"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"a-one","parent":"alpha","before":"a-two"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"a-one-x","parent":"a-one"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"b-one","parent":"beta"}`)
	mustWrite(t, ctx, "add", `{"kind":"intent","handle":"goal"}`)
	mustWrite(t, ctx, "add", `{"kind":"intent","handle":"sub-goal","parent":"goal"}`)
	mustWrite(t, ctx, "add", `{"kind":"decision","handle":"pick"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)
	newChange(t, ctx, `{"intent":"Retire"}`)
	mustWrite(t, ctx, "retire", `{"ref":"beta","reason":"gone"}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{"reason":"ok"}`)
	move(t, ctx, "land", `{}`)
	return ctx
}

// treeCall runs tree with reqJSON.
func treeCall(t *testing.T, ctx *Ctx, reqJSON string) (TreeResponse, error) {
	t.Helper()
	resp, err := callOp(t, ctx, "tree", reqJSON)
	if err != nil {
		return TreeResponse{}, err
	}
	return resp.(TreeResponse), nil
}

// treeShape is each node as "<depth> <handle> <children>[ retired]".
func treeShape(r TreeResponse) []string {
	var out []string
	for _, n := range r.Nodes {
		s := strings.Repeat(" ", n.Depth) + n.Handle + " " + string(rune('0'+n.Children))
		if n.Retired {
			s += " retired"
		}
		out = append(out, s)
	}
	return out
}

func TestTree(t *testing.T) {
	ctx := treeLog(t)
	cases := []struct {
		req  string
		kind op.Kind
		want []string
	}{
		{`{}`, op.KindCapability, []string{"alpha 2", " a-one 1", "  a-one-x 0", " a-two 0"}},
		{`{"depth":1}`, op.KindCapability, []string{"alpha 2", " a-one 0", " a-two 0"}},
		{`{"retired":true}`, op.KindCapability, []string{"alpha 2", " a-one 1", "  a-one-x 0", " a-two 0",
			"beta 1 retired", " b-one 0"}},
		{`{"ref":"a-one"}`, op.KindCapability, []string{"a-one 1", " a-one-x 0"}},
		{`{"ref":"alpha","kind":"capability","depth":1}`, op.KindCapability, []string{"alpha 2", " a-one 0", " a-two 0"}},
		{`{"kind":"intent"}`, op.KindIntent, []string{"goal 1", " sub-goal 0"}},
		{`{"ref":"sub-goal"}`, op.KindIntent, []string{"sub-goal 0"}},
	}
	for _, c := range cases {
		r, err := treeCall(t, ctx, c.req)
		if err != nil {
			t.Fatalf("tree %s: %v", c.req, err)
		}
		if got := treeShape(r); !reflect.DeepEqual(got, c.want) {
			t.Errorf("tree %s = %q, want %q", c.req, got, c.want)
		}
		if r.Kind != c.kind || r.View != "as-is" {
			t.Errorf("tree %s: kind %q view %q", c.req, r.Kind, r.View)
		}
	}

	r, _ := treeCall(t, ctx, `{}`)
	if !r.Root.IsZero() {
		t.Errorf("whole tree has root %s", r.Root)
	}
	want := "- alpha (" + shortID(r.Nodes[0].ID) + ")\n" +
		"  - a-one (" + shortID(r.Nodes[1].ID) + ")\n" +
		"    - a-one-x (" + shortID(r.Nodes[2].ID) + ")\n" +
		"  - a-two (" + shortID(r.Nodes[3].ID) + ")\n"
	if got := human(t, r); got != want {
		t.Errorf("Human =\n%s\nwant\n%s", got, want)
	}
	sub, _ := treeCall(t, ctx, `{"ref":"a-one"}`)
	if sub.Root != r.Nodes[1].ID {
		t.Errorf("subtree root %s, want %s", sub.Root, r.Nodes[1].ID)
	}

	// The retired beta by ID: hidden without --retired, shown with it.
	all, _ := treeCall(t, ctx, `{"retired":true}`)
	beta := all.Nodes[4].ID.String()
	hidden, err := treeCall(t, ctx, `{"ref":"`+beta+`"}`)
	if err != nil || len(hidden.Nodes) != 0 {
		t.Errorf("retired root without --retired: %+v, %v", hidden, err)
	}
	if got := human(t, hidden); !strings.Contains(got, "is retired") || !strings.Contains(got, "--retired") {
		t.Errorf("Human of a hidden root = %q", got)
	}
	shown, _ := treeCall(t, ctx, `{"ref":"`+beta+`","retired":true}`)
	if got := treeShape(shown); !reflect.DeepEqual(got, []string{"beta 1 retired", " b-one 0"}) {
		t.Errorf("retired subtree = %q", got)
	}
	if got := human(t, shown); !strings.HasPrefix(got, "- beta (") || !strings.Contains(got, ") [retired]\n  - b-one (") {
		t.Errorf("Human of a retired subtree = %q", got)
	}
}

func TestTreeEmpty(t *testing.T) {
	ctx := readCtx(t)
	for _, c := range []struct{ req, want string }{
		{`{}`, "No capabilities yet.\n\nnext: genesis add capability <handle>\n"},
		{`{"kind":"intent"}`, "No intents yet.\n\nnext: genesis add intent <handle>\n"},
	} {
		r, err := treeCall(t, ctx, c.req)
		if err != nil {
			t.Fatalf("tree %s: %v", c.req, err)
		}
		if r.Nodes == nil || len(r.Nodes) != 0 {
			t.Errorf("tree %s nodes = %#v, want empty", c.req, r.Nodes)
		}
		if got := human(t, r); got != c.want {
			t.Errorf("Human = %q, want %q", got, c.want)
		}
	}
}

func TestTreeErrors(t *testing.T) {
	ctx := treeLog(t)
	for _, c := range []struct {
		req  string
		code int
		msg  string
		next string
	}{
		{`{"ref":"pick"}`, 2, "not in a tree", "genesis show pick"},
		{`{"ref":"goal","kind":"capability"}`, 2, "not a capability", "genesis tree goal"},
		{`{"ref":"alpha","kind":"intent"}`, 2, "not an intent", "genesis tree alpha"},
		{`{"kind":"decision"}`, 2, "unknown tree kind", "genesis describe tree"},
		{`{"depth":-1}`, 2, "negative", "genesis describe tree"},
		{`{"ref":"nope"}`, 4, `no node matches "nope"`, "genesis ls"},
		{`{"view":"sideways"}`, 2, "unknown view", "genesis describe tree"},
	} {
		_, err := treeCall(t, ctx, c.req)
		if ExitCode(err) != c.code {
			t.Errorf("tree %s: %v (exit %d), want exit %d", c.req, err, ExitCode(err), c.code)
			continue
		}
		e := AsError(err)
		if !strings.Contains(e.Message, c.msg) || e.Next != c.next {
			t.Errorf("tree %s: %q next %q; want %q next %q", c.req, e.Message, e.Next, c.msg, c.next)
		}
	}
}

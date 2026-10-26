package api

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// carveCtx opens change set src ("Bundle") and then tgt ("Milestone"), the
// current one, and returns ctx writing into src.
func carveCtx(t *testing.T) (ctx, inSrc *Ctx, src, tgt id.ID) {
	t.Helper()
	ctx = initCtx(t)
	src = newChange(t, ctx, `{"intent":"Bundle"}`).ChangeSet
	tgt = newChange(t, ctx, `{"intent":"Milestone"}`).ChangeSet
	c := *ctx
	c.ChangeRef = src.String()
	return ctx, &c, src, tgt
}

// carve runs the carve operation from src to tgt.
func carve(t *testing.T, ctx *Ctx, src, tgt id.ID, children bool, nodes ...string) (CarveResponse, error) {
	t.Helper()
	ns, _ := json.Marshal(nodes)
	req := fmt.Sprintf(`{"nodes":%s,"from":%q,"to":%q,"with_children":%t}`, ns, src, tgt, children)
	resp, err := callOp(t, ctx, "carve", req)
	if err != nil {
		return CarveResponse{}, err
	}
	return resp.(CarveResponse), nil
}

func mustCarve(t *testing.T, ctx *Ctx, src, tgt id.ID, children bool, nodes ...string) CarveResponse {
	t.Helper()
	r, err := carve(t, ctx, src, tgt, children, nodes...)
	if err != nil {
		t.Fatalf("carve %v: %v", nodes, err)
	}
	return r
}

// carveBody decodes the ops body of entry e.
func carveBody(t *testing.T, ctx *Ctx, e id.ID) changeset.Ops {
	t.Helper()
	en := entryOf(t, ctx, e)
	b, err := changeset.DecodeBody(en)
	if err != nil {
		t.Fatalf("decode %s: %v", e, err)
	}
	ops, ok := b.(changeset.Ops)
	if !ok {
		t.Fatalf("entry %s is %T", e, b)
	}
	return ops
}

// createOf returns the node.create of n among ops.
func createOf(t *testing.T, ops []op.Op, n id.ID) op.Op {
	t.Helper()
	for _, o := range ops {
		if o.Type == op.TypeNodeCreate && o.Node == n {
			return o
		}
	}
	t.Fatalf("no create of %s in %+v", n, ops)
	return op.Op{}
}

// dependsOnOps returns the live depends-on link creates among ops.
func dependsOnOps(ops []op.Op) []op.Op {
	var out []op.Op
	for _, o := range ops {
		if o.Type == op.TypeLinkCreate && o.LinkType == op.LinkDependsOn {
			out = append(out, o)
		}
	}
	return out
}

func TestCarve(t *testing.T) {
	t.Run("KeepsIDs", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		p := mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`).Nodes[0]
		mustWrite(t, in, "add", `{"kind":"capability","handle":"q"}`)
		old := createOf(t, opsOf(t, ctx, src), p)

		r := mustCarve(t, ctx, src, tgt, false, "p")
		if r.Source != src || r.Target != tgt || !slices.Equal(r.Nodes, []id.ID{p}) ||
			!slices.Equal(r.Dropped, []id.ID{old.ID}) || r.DependsOn != (id.ID{}) || r.Next != "" || r.DryRun {
			t.Fatalf("resp = %+v", r)
		}
		moved := onlyOp(t, opsOf(t, ctx, tgt))
		if moved.Node != p || moved.ID == old.ID || moved.Handle != "p" || !slices.Equal(r.Ops, []id.ID{moved.ID}) {
			t.Errorf("moved = %+v, old %+v", moved, old)
		}
		left := onlyOp(t, opsOf(t, ctx, src))
		if left.Handle != "q" {
			t.Errorf("left in source = %+v", left)
		}

		sb, tb := carveBody(t, ctx, r.SourceEntry), carveBody(t, ctx, r.TargetEntry)
		if sb.Carve == nil || sb.Carve.To != tgt || sb.Carve.Entry != r.TargetEntry || len(sb.Ops) != 0 ||
			!slices.Equal(sb.Drop, []id.ID{old.ID}) {
			t.Errorf("source body = %+v (carve %+v)", sb, sb.Carve)
		}
		if tb.Carve == nil || tb.Carve.From != src || tb.Carve.Entry != r.SourceEntry || len(tb.Drop) != 0 {
			t.Errorf("target body = %+v (carve %+v)", tb, tb.Carve)
		}
		se, te := entryOf(t, ctx, r.SourceEntry), entryOf(t, ctx, r.TargetEntry)
		if se.V != 2 || te.V != 2 || se.ChangeSet != src || te.ChangeSet != tgt {
			t.Errorf("entries v%d %s, v%d %s", se.V, se.ChangeSet, te.V, te.ChangeSet)
		}
		actor := entriesOf(ctxWorkspace(t, ctx).log(), src)[0].Actor
		if se.Actor != actor || te.Actor != actor || se.For != actor || te.For != actor {
			t.Errorf("actors %s/%s, %s/%s; want %s", se.Actor, se.For, te.Actor, te.For, actor)
		}
		n, err := ctxWorkspace(t, ctx).model(tgt).Resolve("p")
		if err != nil || n.ID != p {
			t.Errorf("p in target = %+v, %v", n, err)
		}
		if _, err := ctxWorkspace(t, ctx).model(src).Resolve("p"); ExitCode(err) != 4 {
			t.Errorf("p still in source: %v", err)
		}
		if out := human(t, r); !strings.Contains(out, "Carved 1 node(s) from "+src.String()+" into "+tgt.String()) ||
			!strings.Contains(out, "- node "+p.String()+" (p)") {
			t.Errorf("human:\n%s", out)
		}
	})

	t.Run("BasesRewritten", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		d := landNode(t, ctx, op.KindDecision, "d")
		p := mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`).Nodes[0]
		l := mustWrite(t, in, "link", `{"from":"p","type":"decided-by","to":"d"}`).Links[0]
		create := createOf(t, opsOf(t, ctx, src), p)
		// An update resting on the create, as a dependent's write would.
		w := ctxWorkspace(t, ctx)
		u := op.Op{ID: w.mustID(t), Type: op.TypeNodeUpdate, Node: p, Base: create.ID,
			Set: map[string]json.RawMessage{"statement": raw(`"P"`)}}
		actor := entriesOf(w.log(), src)[0].Actor
		if _, err := w.appendEntry(src, changeset.Ops{Ops: []op.Op{u}}, actor, actor); err != nil {
			t.Fatal(err)
		}

		r := mustCarve(t, ctx, src, tgt, false, "p")
		got := opsOf(t, ctx, tgt)
		if len(got) != 3 || got[0].Type != op.TypeNodeCreate || got[1].Type != op.TypeLinkCreate ||
			got[2].Type != op.TypeNodeUpdate {
			t.Fatalf("target ops = %+v", got)
		}
		if got[2].Base != got[0].ID || got[2].ID == u.ID || got[1].Link != l || got[1].To != d {
			t.Errorf("target ops = %+v", got)
		}
		if !slices.Equal(r.Links, []id.ID{l}) || len(r.Dropped) != 3 || len(opsOf(t, ctx, src)) != 0 {
			t.Errorf("resp = %+v", r)
		}
		n, _ := ctxWorkspace(t, ctx).model(tgt).Resolve("p")
		if string(n.Fields["statement"]) != `"P"` {
			t.Errorf("p in target = %+v", n)
		}
	})

	t.Run("WithChildren", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		var ids []id.ID
		for _, req := range []string{
			`{"kind":"capability","handle":"p"}`,
			`{"kind":"capability","handle":"c","parent":"p"}`,
			`{"kind":"capability","handle":"g","parent":"c"}`,
		} {
			ids = append(ids, mustWrite(t, in, "add", req).Nodes[0])
		}
		r := mustCarve(t, ctx, src, tgt, true, "p")
		if want := slices.SortedFunc(slices.Values(ids), id.Compare); !slices.Equal(r.Nodes, want) {
			t.Errorf("nodes = %v, want %v", r.Nodes, want)
		}
		if len(opsOf(t, ctx, tgt)) != 3 || len(opsOf(t, ctx, src)) != 0 {
			t.Errorf("target %d ops, source %d", len(opsOf(t, ctx, tgt)), len(opsOf(t, ctx, src)))
		}
		// Emptied: nothing rests on the moved nodes, so no dependency.
		if r.Next != "genesis withdraw --change "+shortID(src) || r.DependsOn != (id.ID{}) {
			t.Errorf("resp = %+v", r)
		}
		if s := foldedState(t, ctx, src).State; s != changeset.StateDraft {
			t.Errorf("source is %s", s)
		}
		if out := human(t, r); !strings.Contains(out, "next: genesis withdraw --change "+shortID(src)) {
			t.Errorf("human:\n%s", out)
		}
	})

	t.Run("ChildLeftBehind", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		p := mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`).Nodes[0]
		c := mustWrite(t, in, "add", `{"kind":"capability","handle":"c","parent":"p"}`).Nodes[0]
		r := mustCarve(t, ctx, src, tgt, false, "p")
		dep := dependsOnOps(opsOf(t, ctx, src))
		if len(dep) != 1 || dep[0].From != src || dep[0].To != tgt || dep[0].Link != r.DependsOn {
			t.Fatalf("depends-on = %+v, resp %+v", dep, r)
		}
		if got := changeset.DependsOn(ctxWorkspace(t, ctx).changeSets(), src); !slices.Equal(got, []id.ID{tgt}) {
			t.Errorf("source depends on %v", got)
		}
		n, err := ctxWorkspace(t, ctx).model(src).Resolve("c")
		if err != nil || n.ID != c || n.Parent != p {
			t.Errorf("c in source = %+v, %v", n, err)
		}
		if r.Next != "" {
			t.Errorf("next %q", r.Next)
		}
		if out := human(t, r); !strings.Contains(out, r.DependsOn.String()) {
			t.Errorf("human:\n%s", out)
		}
	})

	t.Run("AlreadyDependsOnTarget", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"c","parent":"p"}`)
		planLink(t, in, op.LinkDependsOn, src, tgt)
		r := mustCarve(t, ctx, src, tgt, false, "p")
		if r.DependsOn != (id.ID{}) || len(dependsOnOps(opsOf(t, ctx, src))) != 1 {
			t.Errorf("resp %+v, depends-on %+v", r, dependsOnOps(opsOf(t, ctx, src)))
		}
		if sb := carveBody(t, ctx, r.SourceEntry); len(sb.Ops) != 0 {
			t.Errorf("source entry ops = %+v", sb.Ops)
		}
	})
}

func TestCarveRefused(t *testing.T) {
	t.Run("Cycle", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"c","parent":"p"}`)
		planLink(t, ctx, op.LinkDependsOn, tgt, src)
		before := logSize(t, ctx)
		_, err := carve(t, ctx, src, tgt, false, "p")
		wantExit(t, err, 3)
		if e := AsError(err); !strings.Contains(e.Message, "loop") || e.Next == "" {
			t.Errorf("err = %q, next %q", e.Message, e.Next)
		}
		if n := logSize(t, ctx); n != before {
			t.Errorf("wrote %d entries", n-before)
		}
	})

	t.Run("SameChangeSet", func(t *testing.T) {
		ctx, in, src, _ := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		_, err := carve(t, ctx, src, src, false, "p")
		wantExit(t, err, 2)
	})

	t.Run("NoNodes", func(t *testing.T) {
		ctx, _, src, tgt := carveCtx(t)
		_, err := carve(t, ctx, src, tgt, false)
		wantExit(t, err, 2)
	})

	t.Run("TerminalTarget", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		move(t, ctx, "withdraw", `{}`)
		_, err := carve(t, ctx, src, tgt, false, "p")
		wantExit(t, err, 3)
		if AsError(err).Next == "" {
			t.Error("no next")
		}
	})

	t.Run("LandedSource", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		for _, s := range []string{"propose", "accept", "land"} {
			move(t, in, s, `{}`)
		}
		_, err := carve(t, ctx, src, tgt, false, "p")
		wantExit(t, err, 3)
		if AsError(err).Next == "" {
			t.Error("no next")
		}
	})

	t.Run("NotCreatedInSource", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		landNode(t, ctx, op.KindDecision, "old")
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		before := logSize(t, ctx)
		_, err := carve(t, ctx, src, tgt, false, "old")
		wantExit(t, err, 2)
		if e := AsError(err); !strings.Contains(e.Message, "not created in") || e.Next == "" {
			t.Errorf("err = %q, next %q", e.Message, e.Next)
		}
		if n := logSize(t, ctx); n != before {
			t.Errorf("wrote %d entries", n-before)
		}
	})

	t.Run("LinkToNodeLeftBehind", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"q"}`)
		l := mustWrite(t, in, "link", `{"from":"p","type":"depends-on","to":"q"}`).Links[0]
		before := logSize(t, ctx)
		_, err := carve(t, ctx, src, tgt, false, "p")
		wantExit(t, err, 3)
		if e := AsError(err); !strings.Contains(e.Message, "would not apply") || !strings.Contains(e.Message, "link.create") ||
			!strings.Contains(e.Message, l.String()) {
			t.Errorf("err = %q", e.Message)
		}
		if n := logSize(t, ctx); n != before {
			t.Errorf("wrote %d entries", n-before)
		}
	})

	t.Run("HandleTakenInTarget", func(t *testing.T) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		mustWrite(t, ctx, "add", `{"kind":"capability","handle":"p"}`)
		_, err := carve(t, ctx, src, tgt, false, "p")
		wantExit(t, err, 3)
		if e := AsError(err); !strings.Contains(e.Message, "handle") {
			t.Errorf("err = %q", e.Message)
		}
	})
}

func TestCarveDryRun(t *testing.T) {
	// setup builds the same log in a fresh repository each time: the seed
	// and clock are fixed, so both runs draw the same IDs.
	setup := func(t *testing.T) (*Ctx, id.ID, id.ID) {
		ctx, in, src, tgt := carveCtx(t)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"p"}`)
		mustWrite(t, in, "add", `{"kind":"capability","handle":"c","parent":"p"}`)
		return ctx, src, tgt
	}
	dctx, src, tgt := setup(t)
	dctx.DryRun = true
	before := logSize(t, dctx)
	files := len(logFiles(t, dctx.Root))
	dry := mustCarve(t, dctx, src, tgt, false, "p")
	if logSize(t, dctx) != before || len(logFiles(t, dctx.Root)) != files {
		t.Fatal("a dry run wrote")
	}
	if !dry.DryRun || !dry.SourceEntry.IsZero() || !dry.TargetEntry.IsZero() || len(dry.SourceOps) != 1 ||
		len(dry.TargetOps) != 1 {
		t.Fatalf("dry = %+v", dry)
	}

	rctx, rsrc, rtgt := setup(t)
	real := mustCarve(t, rctx, rsrc, rtgt, false, "p")
	if rsrc != src || rtgt != tgt {
		t.Fatalf("setups differ: %s %s vs %s %s", src, tgt, rsrc, rtgt)
	}
	enc := func(ops []op.Op) []string {
		var out []string
		for _, o := range ops {
			b, err := op.Encode(o)
			if err != nil {
				t.Fatal(err)
			}
			out = append(out, string(b))
		}
		return out
	}
	str := func(raws []json.RawMessage) []string {
		var out []string
		for _, r := range raws {
			out = append(out, string(r))
		}
		return out
	}
	if got, want := str(dry.SourceOps), enc(carveBody(t, rctx, real.SourceEntry).Ops); !slices.Equal(got, want) {
		t.Errorf("source ops\n%v\nwant\n%v", got, want)
	}
	if got, want := str(dry.TargetOps), enc(carveBody(t, rctx, real.TargetEntry).Ops); !slices.Equal(got, want) {
		t.Errorf("target ops\n%v\nwant\n%v", got, want)
	}
	real.SourceEntry, real.TargetEntry, real.SourceOps, real.TargetOps, real.DryRun = id.ID{}, id.ID{}, dry.SourceOps, dry.TargetOps, true
	a, _ := json.Marshal(dry)
	b, _ := json.Marshal(real)
	if string(a) != string(b) {
		t.Errorf("dry run\n%s\nreal run\n%s", a, b)
	}
	if out := human(t, dry); !strings.Contains(out, "nothing was written") {
		t.Errorf("human:\n%s", out)
	}
}

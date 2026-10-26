package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
)

// entryOf returns the log entry with ID e.
func entryOf(t *testing.T, ctx *Ctx, e id.ID) entry.Entry {
	t.Helper()
	for _, en := range ctxWorkspace(t, ctx).log().Entries {
		if en.ID == e {
			return en
		}
	}
	t.Fatalf("no entry %s", e)
	return entry.Entry{}
}

// recordAda records an actor node for ada@example.com and returns it.
func recordAda(t *testing.T, ctx *Ctx, email string) id.ID {
	t.Helper()
	a, _, err := resolveRefs(ctxWorkspace(t, ctx), "", "", gitConfig(map[string]string{"user.email": email}))
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// twoChanges opens two change sets, the second current, and returns both.
func twoChanges(t *testing.T) (*Ctx, id.ID, id.ID) {
	t.Helper()
	ctx := initCtx(t)
	a := newChange(t, ctx, `{"intent":"A"}`).ChangeSet
	b := newChange(t, ctx, `{"intent":"B"}`).ChangeSet
	return ctx, a, b
}

func TestCtxChangeSet(t *testing.T) {
	t.Run("NoCurrent", func(t *testing.T) {
		ctx, a, _ := twoChanges(t)
		if _, err := callOp(t, ctx, "change.use", `{"none":true}`); err != nil {
			t.Fatal(err)
		}
		ctx.ChangeSet = a
		r := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if r.ChangeSet != a || len(opsOf(t, ctx, a)) != 1 {
			t.Fatalf("wrote into %s, want %s", r.ChangeSet, a)
		}
		if !currentOf(t, ctx).IsZero() {
			t.Error("the write made a change set current")
		}
	})

	t.Run("BeatsCurrent", func(t *testing.T) {
		ctx, a, b := twoChanges(t)
		ctx.ChangeSet = a
		r := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if r.ChangeSet != a || len(opsOf(t, ctx, b)) != 0 {
			t.Fatalf("wrote into %s, want %s", r.ChangeSet, a)
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		ctx, _, _ := twoChanges(t)
		ctx.ChangeSet = ctxWorkspace(t, ctx).mustID(t)
		_, err := callOp(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		wantExit(t, err, 4)
		wantNext(t, err, changeLsNext)
	})

	t.Run("AndRef", func(t *testing.T) {
		ctx, a, b := twoChanges(t)
		ctx.ChangeSet, ctx.ChangeRef = a, b.String()
		for _, name := range []string{"add", "propose", "change.show"} {
			req := `{}`
			if name == "add" {
				req = `{"kind":"capability","handle":"store"}`
			}
			_, err := callOp(t, ctx, name, req)
			wantExit(t, err, 2)
		}
	})

	t.Run("Lifecycle", func(t *testing.T) {
		ctx, a, b := twoChanges(t)
		ctx.ChangeSet = a
		resp, err := callOp(t, ctx, "propose", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(LifecycleResponse); r.ChangeSet != a {
			t.Fatalf("proposed %s, want %s", r.ChangeSet, a)
		}
		c, _ := ctxWorkspace(t, ctx).changeSet(b)
		if c.State != "draft" {
			t.Errorf("current is %s, want draft", c.State)
		}
	})

	t.Run("Apply", func(t *testing.T) {
		ctx, a, b := twoChanges(t)
		ctx.ChangeSet = a
		r, err := applyFile(t, ctx, treeWrites)
		if err != nil {
			t.Fatal(err)
		}
		if r.ChangeSet != a || len(opsOf(t, ctx, a)) == 0 || len(opsOf(t, ctx, b)) != 0 {
			t.Fatalf("applied into %s, want %s", r.ChangeSet, a)
		}
	})

	t.Run("Show", func(t *testing.T) {
		ctx, a, _ := twoChanges(t)
		ctx.ChangeSet = a
		resp, err := callOp(t, ctx, "change.show", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(ChangeShowResponse); r.ID != a || r.Current {
			t.Fatalf("showed %s (current %v), want %s", r.ID, r.Current, a)
		}
	})
}

func TestCtxActor(t *testing.T) {
	t.Run("Honoured", func(t *testing.T) {
		ctx := initCtx(t)
		ada := recordAda(t, ctx, "ada@example.com")
		ctx.Actor = ada
		r := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if e := entryOf(t, ctx, r.Entry); e.Actor != ada || e.For != ada {
			t.Fatalf("entry by %s for %s, want %s", e.Actor, e.For, ada)
		}
		if _, ok := ctxWorkspace(t, ctx).model(id.ID{}).ActorByEmail("test@example.com"); ok {
			t.Error("the git actor was recorded")
		}
	})

	t.Run("For", func(t *testing.T) {
		ctx := initCtx(t)
		ada := recordAda(t, ctx, "ada@example.com")
		bob := recordAda(t, ctx, "bob@example.com")
		ctx.Actor, ctx.For = ada, bob
		r := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		if e := entryOf(t, ctx, r.Entry); e.Actor != ada || e.For != bob {
			t.Fatalf("entry by %s for %s, want %s for %s", e.Actor, e.For, ada, bob)
		}
	})

	t.Run("ForOverGit", func(t *testing.T) {
		ctx, cs := writeCtx(t)
		bob := recordAda(t, ctx, "bob@example.com")
		ctx.For = bob
		r := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		e := entryOf(t, ctx, r.Entry)
		if e.ChangeSet != cs || e.For != bob || e.Actor == bob {
			t.Fatalf("entry %+v, want for %s", e, bob)
		}
	})

	t.Run("NotAnActor", func(t *testing.T) {
		ctx := initCtx(t)
		node := recordWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`).Nodes[0]
		ctx.Actor = node
		_, err := callOp(t, ctx, "add", `{"kind":"capability","handle":"log"}`)
		wantExit(t, err, 2)
		c := *ctx
		c.Actor, c.For = id.ID{}, node
		_, err = callOp(t, &c, "add", `{"kind":"capability","handle":"log"}`)
		wantExit(t, err, 2)
	})

	t.Run("Absent", func(t *testing.T) {
		ctx := initCtx(t)
		ctx.Record = true
		ctx.Actor = ctxWorkspace(t, ctx).mustID(t)
		_, err := callOp(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		wantExit(t, err, 4)
	})

	t.Run("AndRef", func(t *testing.T) {
		ctx := initCtx(t)
		ada := recordAda(t, ctx, "ada@example.com")
		ctx.Record = true
		ctx.Actor, ctx.ActorRef = ada, ada.String()
		_, err := callOp(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		wantExit(t, err, 2)
		ctx.Actor, ctx.ActorRef, ctx.For, ctx.ForRef = id.ID{}, "", ada, ada.String()
		_, err = callOp(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
		wantExit(t, err, 2)
	})
}

func TestRecordRefusedOnNonWrites(t *testing.T) {
	ctx, _ := writeCtx(t)
	ctx.Record = true
	for _, c := range []struct{ op, req string }{
		{"change.ls", `{}`}, {"change.show", `{}`}, {"change.use", `{"none":true}`},
		{"change.new", `{"intent":"x"}`}, {"propose", `{}`}, {"init", `{}`}, {"describe", `{}`},
	} {
		_, err := callOp(t, ctx, c.op, c.req)
		wantExit(t, err, 2)
		wantNext(t, err, "rhiza describe "+c.op)
	}
	for _, o := range Operations() {
		want := map[string]bool{"add": true, "set": true, "mv": true, "retire": true, "link": true, "unlink": true, "apply": true, "source.add": true}[o.Name]
		if o.Record != want {
			t.Errorf("%s: Record = %v, want %v", o.Name, o.Record, want)
		}
	}
}

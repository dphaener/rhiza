package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/op"
)

// move runs a lifecycle operation and fails the test on error.
func move(t *testing.T, ctx *Ctx, name, req string) LifecycleResponse {
	t.Helper()
	resp, err := callOp(t, ctx, name, req)
	if err != nil {
		t.Fatalf("%s %s: %v", name, req, err)
	}
	return resp.(LifecycleResponse)
}

// foldedState folds cs from the log on disk.
func foldedState(t *testing.T, ctx *Ctx, cs id.ID) changeset.ChangeSet {
	t.Helper()
	c, ok := ctxWorkspace(t, ctx).changeSet(cs)
	if !ok {
		t.Fatalf("change set %s has no entries", cs)
	}
	return c
}

// logSize counts the entries in the log on disk.
func logSize(t *testing.T, ctx *Ctx) int {
	t.Helper()
	return len(ctxWorkspace(t, ctx).log().Entries)
}

// wantNext checks err's next command.
func wantNext(t *testing.T, err error, next string) {
	t.Helper()
	e := AsError(err)
	if e == nil {
		t.Errorf("err = nil, want next %q", next)
		return
	}
	if e.Next != next {
		t.Errorf("err = %v, next %q; want %q", err, e.Next, next)
	}
}

func TestLifecycleIntentionPath(t *testing.T) {
	ctx := initCtx(t)
	cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet

	r := move(t, ctx, "propose", `{}`)
	if r.ChangeSet != cs || r.From != changeset.StateDraft || r.To != changeset.StateProposed || r.Entry.IsZero() || r.ClearedCurrent {
		t.Errorf("propose = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "Proposed change set "+cs.String()+" (draft → proposed).") {
		t.Errorf("propose human:\n%s", out)
	}
	if got := foldedState(t, ctx, cs).State; got != changeset.StateProposed {
		t.Errorf("state = %s", got)
	}

	r = move(t, ctx, "accept", `{"reason":"ok"}`)
	if r.From != changeset.StateProposed || r.To != changeset.StateAccepted || r.Rule != "owner" || r.Reason != "ok" {
		t.Errorf("accept = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "next: genesis land") {
		t.Errorf("accept human:\n%s", out)
	}
	c := foldedState(t, ctx, cs)
	if c.State != changeset.StateAccepted || c.Decided == nil || c.Decided.Type != entry.TypeAccept ||
		c.Decided.Rule != "owner" || c.Decided.Reason != "ok" || c.Decided.Entry != r.Entry {
		t.Errorf("folded after accept = %+v / %+v", c.State, c.Decided)
	}
	if currentOf(t, ctx) != cs {
		t.Error("accept cleared the current change set")
	}

	r = move(t, ctx, "land", `{}`)
	if r.From != changeset.StateAccepted || r.To != changeset.StateLanded || !r.ClearedCurrent {
		t.Errorf("land = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "No current change set now.") {
		t.Errorf("land human:\n%s", out)
	}
	if got := foldedState(t, ctx, cs).State; got != changeset.StateLanded {
		t.Errorf("state = %s", got)
	}
	if got := currentOf(t, ctx); !got.IsZero() {
		t.Errorf("current after land = %s", got)
	}
}

func TestLifecycleCustomRule(t *testing.T) {
	ctx := initCtx(t)
	cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
	move(t, ctx, "propose", `{}`)
	r := move(t, ctx, "reject", `{"rule":"review","reason":"too big"}`)
	if r.Rule != "review" || r.Reason != "too big" || r.To != changeset.StateRejected || !r.ClearedCurrent {
		t.Errorf("reject = %+v", r)
	}
	c := foldedState(t, ctx, cs)
	if c.State != changeset.StateRejected || c.Decided == nil || c.Decided.Type != entry.TypeReject ||
		c.Decided.Rule != "review" || c.Decided.Reason != "too big" {
		t.Errorf("folded = %s / %+v", c.State, c.Decided)
	}
	if got := currentOf(t, ctx); !got.IsZero() {
		t.Errorf("current after reject = %s", got)
	}
}

func TestLifecycleWithdraw(t *testing.T) {
	for _, c := range []struct {
		state changeset.State
		steps []string
	}{
		{changeset.StateDraft, nil},
		{changeset.StateProposed, []string{"propose"}},
		{changeset.StateAccepted, []string{"propose", "accept"}},
	} {
		t.Run(string(c.state), func(t *testing.T) {
			ctx := initCtx(t)
			cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
			for _, s := range c.steps {
				move(t, ctx, s, `{}`)
			}
			r := move(t, ctx, "withdraw", `{}`)
			if r.From != c.state || r.To != changeset.StateWithdrawn || !r.ClearedCurrent {
				t.Errorf("withdraw = %+v", r)
			}
			if got := foldedState(t, ctx, cs).State; got != changeset.StateWithdrawn {
				t.Errorf("state = %s", got)
			}
			if got := currentOf(t, ctx); !got.IsZero() {
				t.Errorf("current after withdraw = %s", got)
			}
		})
	}
}

func TestLifecycleIllegal(t *testing.T) {
	for _, c := range []struct {
		name  string
		steps []string
		op    string
		next  string
	}{
		{"land draft intention", nil, "land", "genesis propose"},
		{"land proposed intention", []string{"propose"}, "land", "genesis accept"},
		{"accept draft", nil, "accept", "genesis propose"},
		{"reject draft", nil, "reject", "genesis propose"},
		{"accept accepted", []string{"propose", "accept"}, "accept", "genesis land"},
		{"propose proposed", []string{"propose"}, "propose", "genesis accept"},
		{"withdraw landed", []string{"propose", "accept", "land"}, "withdraw", newEntryNext},
		{"accept rejected", []string{"propose", "reject"}, "accept", newEntryNext},
		{"land withdrawn", []string{"withdraw"}, "land", newEntryNext},
	} {
		t.Run(c.name, func(t *testing.T) {
			ctx := initCtx(t)
			cs := newChange(t, ctx, `{"intent":"X"}`).ChangeSet
			for _, s := range c.steps {
				move(t, ctx, s, `{}`)
			}
			before := foldedState(t, ctx, cs).State
			size := logSize(t, ctx)
			// A terminal move cleared the current change set; name it.
			ctx.ChangeRef = cs.String()
			_, err := callOp(t, ctx, c.op, `{}`)
			wantExit(t, err, 3)
			wantNext(t, err, c.next)
			if !strings.Contains(err.Error(), "cannot "+c.op+" change set "+cs.String()+": it is "+string(before)) {
				t.Errorf("message = %q", err.Error())
			}
			if got := logSize(t, ctx); got != size {
				t.Errorf("log grew from %d to %d entries", size, got)
			}
		})
	}
}

func TestLifecycleRecordLandsFromDraft(t *testing.T) {
	ctx := initCtx(t)
	cs := newChange(t, ctx, `{"intent":"R","kind":"record"}`).ChangeSet
	r := move(t, ctx, "land", `{}`)
	if r.From != changeset.StateDraft || r.To != changeset.StateLanded || !r.ClearedCurrent {
		t.Errorf("land = %+v", r)
	}
	if got := foldedState(t, ctx, cs).State; got != changeset.StateLanded {
		t.Errorf("state = %s", got)
	}
}

func TestLifecycleOtherChangeSet(t *testing.T) {
	ctx := initCtx(t)
	other := newChange(t, ctx, `{"intent":"Other"}`).ChangeSet
	cur := newChange(t, ctx, `{"intent":"Current"}`).ChangeSet
	ctx.ChangeRef = other.String()
	r := move(t, ctx, "withdraw", `{}`)
	if r.ChangeSet != other || r.ClearedCurrent {
		t.Errorf("withdraw = %+v", r)
	}
	if out := human(t, r); strings.Contains(out, "No current change set") {
		t.Errorf("human:\n%s", out)
	}
	if got := currentOf(t, ctx); got != cur {
		t.Errorf("current = %s, want %s", got, cur)
	}
	if got := foldedState(t, ctx, other).State; got != changeset.StateWithdrawn {
		t.Errorf("other state = %s", got)
	}
	if got := foldedState(t, ctx, cur).State; got != changeset.StateDraft {
		t.Errorf("current state = %s", got)
	}
}

func TestLifecycleDryRun(t *testing.T) {
	ctx := initCtx(t)
	cs := newChange(t, ctx, `{"intent":"R","kind":"record"}`).ChangeSet
	size := logSize(t, ctx)
	ctx.DryRun = true
	r := move(t, ctx, "land", `{}`)
	if !r.DryRun || r.ClearedCurrent || r.To != changeset.StateLanded || !r.Entry.IsZero() {
		t.Errorf("dry land = %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "dry run") {
		t.Errorf("human:\n%s", out)
	}
	if got := logSize(t, ctx); got != size {
		t.Errorf("log grew from %d to %d entries", size, got)
	}
	if got := currentOf(t, ctx); got != cs {
		t.Errorf("current = %s, want %s", got, cs)
	}
	if got := foldedState(t, ctx, cs).State; got != changeset.StateDraft {
		t.Errorf("state = %s", got)
	}
}

func TestLifecycleNoChangeSet(t *testing.T) {
	ctx := initCtx(t)
	for _, name := range []string{"propose", "accept", "reject", "withdraw", "land"} {
		_, err := callOp(t, ctx, name, `{}`)
		wantExit(t, err, 2)
		wantNext(t, err, changeLsNext)
		if !strings.HasPrefix(err.Error(), "no change set") {
			t.Errorf("%s: %v", name, err)
		}
	}
	ctx.ChangeRef = "ZZZZZZ"
	_, err := callOp(t, ctx, "propose", `{}`)
	wantExit(t, err, 4)
	if _, err := os.Stat(ctx.Root + "/.genesis/log"); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleRecordFlag(t *testing.T) {
	ctx := initCtx(t)
	newChange(t, ctx, `{"intent":"X"}`)
	ctx.Record = true
	_, err := callOp(t, ctx, "propose", `{}`)
	wantExit(t, err, 2)
}

// stackedPair makes change sets A and B, accepted, with B depending on A
// through a depends-on link written in B. It returns them with neither
// current.
func stackedPair(t *testing.T) (ctx *Ctx, a, b id.ID) {
	t.Helper()
	ctx = initCtx(t)
	a = newChange(t, ctx, `{"intent":"A"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"store"}`)
	b = newChange(t, ctx, `{"intent":"B"}`).ChangeSet
	mustWrite(t, ctx, "link", fmt.Sprintf(`{"from":%q,"type":"depends-on","to":%q}`, b, a))
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"log"}`)
	for _, cs := range []id.ID{a, b} {
		ctx.ChangeRef = cs.String()
		move(t, ctx, "propose", `{}`)
		move(t, ctx, "accept", `{}`)
	}
	return ctx, a, b
}

// wantLandRefused lands cs and checks the refusal: exit 3, the message
// starting with want, the next command, and no entry written.
func wantLandRefused(t *testing.T, ctx *Ctx, cs id.ID, want, next string) {
	t.Helper()
	size := logSize(t, ctx)
	c := *ctx
	c.ChangeRef = cs.String()
	_, err := callOp(t, &c, "land", `{}`)
	wantExit(t, err, 3)
	wantNext(t, err, next)
	if e := AsError(err); !strings.HasPrefix(e.Message, want) {
		t.Errorf("message = %q, want prefix %q", e.Message, want)
	}
	if got := logSize(t, ctx); got != size {
		t.Errorf("log grew from %d to %d entries", size, got)
	}
}

func TestLandDependencies(t *testing.T) {
	t.Run("Accepted", func(t *testing.T) {
		ctx, a, b := stackedPair(t)
		want := "cannot land change set " + b.String() + ": it depends on " + a.String() + " (accepted)"
		wantLandRefused(t, ctx, b, want, "genesis land --change "+shortID(a))
		ctx.DryRun = true
		wantLandRefused(t, ctx, b, want, "genesis land --change "+shortID(a))
	})

	t.Run("InOrder", func(t *testing.T) {
		ctx, a, b := stackedPair(t)
		ctx.ChangeRef = a.String()
		move(t, ctx, "land", `{}`)
		ctx.ChangeRef = b.String()
		if r := move(t, ctx, "land", `{}`); r.To != changeset.StateLanded {
			t.Errorf("land B = %+v", r)
		}
	})

	t.Run("Withdrawn", func(t *testing.T) {
		ctx, a, b := stackedPair(t)
		ctx.ChangeRef = a.String()
		move(t, ctx, "withdraw", `{}`)
		wantLandRefused(t, ctx, b,
			"cannot land change set "+b.String()+": it depends on "+a.String()+" (withdrawn)",
			"genesis change show "+shortID(a))
	})

	t.Run("Proposed", func(t *testing.T) {
		ctx := initCtx(t)
		a := newChange(t, ctx, `{"intent":"A"}`).ChangeSet
		move(t, ctx, "propose", `{}`)
		b := newChange(t, ctx, `{"intent":"R","kind":"record"}`).ChangeSet
		mustWrite(t, ctx, "link", fmt.Sprintf(`{"from":%q,"type":"depends-on","to":%q}`, b, a))
		wantLandRefused(t, ctx, b,
			"cannot land change set "+b.String()+": it depends on "+a.String()+" (proposed)",
			"genesis change show "+shortID(a))
	})

	t.Run("NotInLog", func(t *testing.T) {
		ctx := initCtx(t)
		b := newChange(t, ctx, `{"intent":"R","kind":"record"}`).ChangeSet
		w := ctxWorkspace(t, ctx)
		actor, forID, err := w.resolveActor(b, ctx)
		if err != nil {
			t.Fatal(err)
		}
		newID := func() id.ID {
			v, err := w.ids.New()
			if err != nil {
				t.Fatal(err)
			}
			return v
		}
		ghost := newID()
		o := op.Op{ID: newID(), Type: op.TypeLinkCreate, Link: newID(),
			LinkType: op.LinkDependsOn, From: b, To: ghost}
		if _, err := w.appendEntry(b, changeset.Ops{Ops: []op.Op{o}}, actor, forID); err != nil {
			t.Fatal(err)
		}
		wantLandRefused(t, ctx, b,
			"cannot land change set "+b.String()+": it depends on "+ghost.String()+" (not in the log)",
			"genesis change show "+shortID(ghost))
	})
}

func TestLandConflict(t *testing.T) {
	ctx := initCtx(t)
	recordWrite(t, ctx, "add", `{"kind":"decision","handle":"base","fields":["statement=Old"]}`)
	b := newChange(t, ctx, `{"intent":"B"}`).ChangeSet
	mustWrite(t, ctx, "set", `{"ref":"base","fields":["statement=Mine"]}`)
	move(t, ctx, "propose", `{}`)
	move(t, ctx, "accept", `{}`)
	recordWrite(t, ctx, "set", `{"ref":"base","fields":["statement=Theirs"]}`)

	want := "cannot land change set " + b.String() + ": model: conflict"
	wantLandRefused(t, ctx, b, want, "genesis check")
	ctx.DryRun = true
	wantLandRefused(t, ctx, b, want, "genesis check")
}

// TestLandRevisedLanded: A revises in place a field B wrote; once A has
// landed, land still refuses B (D13) until B writes the field again.
func TestLandRevisedLanded(t *testing.T) {
	ctx := initCtx(t)
	a := newChange(t, ctx, `{"intent":"A"}`).ChangeSet
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"base","fields":["statement=Base"]}`)
	b := newChange(t, ctx, `{"intent":"B"}`).ChangeSet
	mustWrite(t, ctx, "link", fmt.Sprintf(`{"from":%q,"type":"depends-on","to":%q}`, b, a))
	mustWrite(t, ctx, "set", `{"ref":"base","fields":["statement=FromB"]}`)
	ctx.ChangeRef = a.String()
	mustWrite(t, ctx, "set", `{"ref":"base","fields":["statement=Revised"]}`)
	for _, cs := range []id.ID{a, b} {
		ctx.ChangeRef = cs.String()
		move(t, ctx, "propose", `{}`)
		move(t, ctx, "accept", `{}`)
	}
	ctx.ChangeRef = a.String()
	move(t, ctx, "land", `{}`)

	want := "cannot land change set " + b.String() + ": model: conflict: change set " + a.String() + " revised op "
	wantLandRefused(t, ctx, b, want, "genesis check")

	ctx.ChangeRef = b.String()
	mustWrite(t, ctx, "set", `{"ref":"base","fields":["statement=FromBAgain"]}`)
	if r := move(t, ctx, "land", `{}`); r.To != changeset.StateLanded {
		t.Errorf("land B = %+v", r)
	}
}

// TestLandProblem: an op on a node that is not there is not a conflict, but
// land refuses it all the same, pointing at the what-if view.
func TestLandProblem(t *testing.T) {
	ctx := initCtx(t)
	b := newChange(t, ctx, `{"intent":"R","kind":"record"}`).ChangeSet
	w := ctxWorkspace(t, ctx)
	actor, forID, err := w.resolveActor(b, ctx)
	if err != nil {
		t.Fatal(err)
	}
	newID := func() id.ID {
		v, err := w.ids.New()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	ghost, base := newID(), newID()
	o := op.Op{ID: newID(), Type: op.TypeNodeUpdate, Node: ghost, Base: base,
		Set: map[string]json.RawMessage{"statement": json.RawMessage(`"x"`)}}
	if _, err := w.appendEntry(b, changeset.Ops{Ops: []op.Op{o}}, actor, forID); err != nil {
		t.Fatal(err)
	}
	want := "cannot land change set " + b.String() + ": model: not found"
	wantLandRefused(t, ctx, b, want, "genesis slice --change "+shortID(b)+" --view what-if")
	ctx.DryRun = true
	wantLandRefused(t, ctx, b, want, "genesis slice --change "+shortID(b)+" --view what-if")
}

package api

import (
	"os"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/entry"
	"github.com/dphaener/genesis/internal/id"
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
	if !r.DryRun || r.ClearedCurrent || r.To != changeset.StateLanded || r.Entry.IsZero() {
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

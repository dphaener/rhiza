package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/quiettest"
	"github.com/dphaener/rhiza/internal/run"
)

// inboxGateBudget is the most the gate-question read of a plan run may take
// over a large log: the owner's bar is well under a second.
const inboxGateBudget = time.Second

// TestInboxGateQuestionsBudget stops a plan run at its gate with open
// questions spread over several planned change sets, in a repository with a
// large log, and requires that inboxGateQuestions answers every question
// with its statement without decoding the log, within the budget.
func TestInboxGateQuestionsBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("generates a large log; skipped under -short")
	}
	q := quiettest.Wait(t)
	_ = q
	budgetLock(t)
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", cache)
	root := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gen.Write(root, gen.Options{Seed: 7, Nodes: 3000, Links: 1200, Landed: 80, InFlight: 8}); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	p, err := index.DefaultPath(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx.IndexPath = p
	seededIDs(ctx, 5000)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	// The plan's own change set creates the intent and lands; three planned
	// change sets serve it, each asking open questions.
	own := call("change.new", `{"intent":"Plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"gate-plan-intent"}`)
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	call("land", `{}`)
	var handles []string
	var stopped string
	for c := range 3 {
		cs := call("change.new", fmt.Sprintf(`{"intent":"Planned %d"}`, c)).(ChangeNewResponse).ChangeSet
		call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"gate-plan-intent"}`, cs.String()))
		for k := range 2 {
			h := fmt.Sprintf("gate-q-%d-%d", c, k)
			call("add", fmt.Sprintf(`{"kind":"question","handle":%q,"fields":["statement=Statement of %s"]}`, h, h))
			handles = append(handles, h)
		}
		call("propose", `{}`)
	}
	for i, h := range handles {
		if i > 0 {
			stopped += ", "
		}
		stopped += h
	}
	st := run.State{Change: own.String(), Worktree: root, Stopped: run.OpenStop + stopped}

	// The worktree's index is current, as the run's own steps leave it.
	x, err := index.Open(root, p)
	if err != nil {
		t.Fatal(err)
	}
	x.Close()

	reads := countingCtx(ctx)
	for attempt := 1; ; attempt++ {
		start := time.Now()
		got := inboxGateQuestions(ctx, st)
		d := time.Since(start)
		if len(got) != len(handles) {
			t.Fatalf("got %d questions, want %d", len(got), len(handles))
		}
		for i, g := range got {
			if want := "Statement of " + handles[i]; g.Statement != want {
				t.Errorf("question %s: statement %q, want %q", handles[i], g.Statement, want)
			}
		}
		if *reads != 0 {
			t.Fatalf("the gate read decoded the full log %d times", *reads)
		}
		if d < inboxGateBudget {
			break
		}
		if attempt > readBudgetRetries {
			t.Fatalf("gate read took %v, over %v on every attempt", d, inboxGateBudget)
		}
	}

	x, err = index.Open(root, p)
	if err != nil {
		t.Fatal(err)
	}
	defer x.Close()
	if c := x.LastCatchUp(); c.Decoded != 0 || c.Rebuilt {
		t.Errorf("index catch-up after the read: %+v, want none decoded", c)
	}
}

// TestInboxGateQuestionsKept stops a plan run at its gate on two questions
// whose worktree does not read, and on one whose handle resolves nowhere:
// the inbox still lists every handle, statements empty.
func TestInboxGateQuestionsKept(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	root := t.TempDir() // no .rhiza: the worktree does not read
	ctx := testCtx(root)
	st := run.State{Change: "01M6RXK2SB0M12FXVEGXER022R", Worktree: root, Stopped: run.OpenStop + "kept-a, kept-b"}
	got := inboxGateQuestions(ctx, st)
	if len(got) != 2 || got[0].Handle != "kept-a" || got[1].Handle != "kept-b" {
		t.Fatalf("got %+v, want kept-a and kept-b", got)
	}
	for _, g := range got {
		if g.Statement != "" {
			t.Errorf("%s: statement %q, want empty", g.Handle, g.Statement)
		}
	}
}

package api

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/op"
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
	var planIDs []id.ID
	var stopped string
	for c := range 3 {
		cs := call("change.new", fmt.Sprintf(`{"intent":"Planned %d"}`, c)).(ChangeNewResponse).ChangeSet
		planIDs = append(planIDs, cs)
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
	var builds map[id.ID]int
	ctx.builtModel = func(cs id.ID) { builds[cs]++ }
	planned := map[id.ID]bool{own: true}
	for _, c := range planIDs {
		planned[c] = true
	}
	for attempt := 1; ; attempt++ {
		builds = map[id.ID]int{}
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
		// Each change set's model is folded at most once per call, and
		// every planned change set (the plan's own and the three it
		// plans) is read.
		for cs, n := range builds {
			if n > 1 {
				t.Errorf("the gate read folded change set %s's model %d times, want at most once", cs, n)
			}
		}
		for cs := range planned {
			if builds[cs] != 1 {
				t.Errorf("the gate read folded planned change set %s's model %d times, want once", cs, builds[cs])
			}
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
// whose worktree does not read: the inbox still lists both handles,
// statements empty.
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

// TestInboxGateQuestionsUnresolved stops a plan run at its gate on two
// questions whose worktree reads, one of whose handles resolves in no
// planned change set: the inbox lists both handles, the unresolved one with
// an empty statement and the resolved one with its full statement.
func TestInboxGateQuestionsUnresolved(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	ctx := initCtx(t)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	own := call("change.new", `{"intent":"Plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"kept-plan-intent"}`)
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	call("land", `{}`)
	cs := call("change.new", `{"intent":"Planned"}`).(ChangeNewResponse).ChangeSet
	call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"kept-plan-intent"}`, cs.String()))
	call("add", `{"kind":"question","handle":"kept-resolved","fields":["statement=Statement of kept-resolved"]}`)
	call("propose", `{}`)
	st := run.State{Change: own.String(), Worktree: ctx.Root, Stopped: run.OpenStop + "kept-resolved, kept-nowhere"}
	got := inboxGateQuestions(ctx, st)
	if len(got) != 2 || got[0].Handle != "kept-resolved" || got[1].Handle != "kept-nowhere" {
		t.Fatalf("got %+v, want kept-resolved and kept-nowhere", got)
	}
	if want := "Statement of kept-resolved"; got[0].Statement != want {
		t.Errorf("kept-resolved: statement %q, want %q", got[0].Statement, want)
	}
	if got[1].Statement != "" {
		t.Errorf("kept-nowhere: statement %q, want empty", got[1].Statement)
	}
}

// TestGateAnswerUsesWorktreeIndex answers a question at a plan run's gate
// (question.answer --run, which the Inbox app's gate submit calls once per
// question): the worktree's own index answers the gate's reads and the
// answer's write, so the full disk log is never decoded.
func TestGateAnswerUsesWorktreeIndex(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", cache)
	ctx := initCtx(t)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	own := call("change.new", `{"intent":"Plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"answer-plan-intent"}`)
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	call("land", `{}`)
	cs := call("change.new", `{"intent":"Planned"}`).(ChangeNewResponse).ChangeSet
	call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"answer-plan-intent"}`, cs.String()))
	call("add", `{"kind":"question","handle":"answer-q","fields":["statement=Statement of answer-q"]}`)
	call("propose", `{}`)

	common := t.TempDir()
	runs := run.NewStore(run.Dir(common))
	at := ctx.Now()
	if err := runs.Create("RGATE", run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan", Change: own.String(), Worktree: ctx.Root}}); err != nil {
		t.Fatal(err)
	}
	if err := runs.Append("RGATE", run.Entry{Type: run.TypeStopped, At: at, Step: "gate", Detail: run.OpenStop + "answer-q; answer each with a Decision, then rhiza run resume RGATE"}); err != nil {
		t.Fatal(err)
	}
	ctx.Git = func(args ...string) (string, error) { return common + "\n", nil }
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		return []byte(common + "\n"), nil
	}
	// The worktree's index is current, as the run's own steps leave it.
	p, err := index.DefaultPath(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	x, err := index.Open(ctx.Root, p)
	if err != nil {
		t.Fatal(err)
	}
	x.Close()

	reads := countingCtx(ctx)
	call("question.answer", `{"question":"answer-q","text":"Yes, do it.","run":"RGATE"}`)
	if *reads != 0 {
		t.Fatalf("question.answer --run decoded the full log %d times", *reads)
	}
	if d, err := ctxWorkspace(t, ctx).model(cs).Resolve("answer-answer-q"); err != nil || d.Kind != op.KindDecision {
		t.Errorf("no answering Decision in %s: %+v, %v", cs, d, err)
	}
}

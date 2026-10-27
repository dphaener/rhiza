package api

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// arbiterFixture is an autonomous run's step runner on a fresh change set
// whose worktree is open, with a ctx that writes into that worktree (its
// own ID seed, so its IDs never meet the launcher's).
func arbiterFixture(t *testing.T) (*sessionFixture, *stepRunner, *Ctx) {
	t.Helper()
	f := newSessionFixture(t, "arbiter-runs")
	cs := newChange(t, f.ctx, `{"intent":"Arbiter runs","title":"Arbiter runs"}`).ChangeSet
	r := &stepRunner{
		ctx: f.ctx, change: cs.String(), run: "01RUNARBITERAAAAAAAAAAAAAA", worktree: f.worktree,
		mode: "manual", autonomous: true, sessions: session.NewStore(session.Dir(f.common)),
	}
	t0, err := resolveChangeWorktree(f.ctx, r.change, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := t0.open(f.ctx); err != nil {
		t.Fatal(err)
	}
	wctx := testCtx(f.worktree)
	wctx.IDs = id.NewGenerator(wctx.Now, rand.New(rand.NewSource(77)))
	return f, r, wctx
}

// TestStartArbiterPrompt launches the arbiter with a prompt that is the
// fixed text over the question's ID and names no other node, change set
// or run.
func TestStartArbiterPrompt(t *testing.T) {
	f, r, wctx := arbiterFixture(t)
	q := mustWrite(t, wctx, "add", `{"kind":"question","handle":"which-store","fields":["statement=Which store?"]}`).Nodes[0]
	step := workflow.Step{Name: workflow.ArbiterName, Kind: workflow.KindArbiter, Interval: workflow.DefaultInterval}
	sid, err := r.StartArbiter(step, q.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := r.sessions.Read(sid)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf(arbiterPrompt, q)
	if rec.Prompt != want || want != "Decide question "+q.String()+": follow rhiza guide arbiter." {
		t.Errorf("prompt %q, want %q", rec.Prompt, want)
	}
	for _, other := range []string{"which-store", r.change, r.run, "Arbiter runs", "arbiter-runs", f.worktree} {
		if strings.Contains(strings.Replace(rec.Prompt, "rhiza guide arbiter", "", 1), other) {
			t.Errorf("prompt %q names %q", rec.Prompt, other)
		}
	}
	if rec.Step != workflow.ArbiterName || rec.Run != r.run || rec.Change != r.change || !rec.Autonomous || rec.PermissionMode != "manual" {
		t.Errorf("session record %+v", rec)
	}
}

// TestArbiterAnswer accepts only a live Decision with a live answers link
// to the question, the assumed wording and a reasoning.
func TestArbiterAnswer(t *testing.T) {
	_, r, wctx := arbiterFixture(t)
	n := 0
	question := func() id.ID {
		n++
		return mustWrite(t, wctx, "add", fmt.Sprintf(`{"kind":"question","handle":"q%d","fields":["statement=Q%d?"]}`, n, n)).Nodes[0]
	}
	decide := func(h, wording, reasoning string, answers id.ID) {
		fields := []string{`"statement=Use SQLite"`}
		if wording != "" {
			fields = append(fields, fmt.Sprintf("%q", "wording="+wording))
		}
		if reasoning != "" {
			fields = append(fields, fmt.Sprintf("%q", "reasoning="+reasoning))
		}
		mustWrite(t, wctx, "add", `{"kind":"decision","handle":"`+h+`","fields":[`+strings.Join(fields, ",")+`]}`)
		if !answers.IsZero() {
			mustWrite(t, wctx, "link", fmt.Sprintf(`{"from":%q,"type":"answers","to":%q}`, h, answers))
		}
	}
	answer := func(q id.ID) string {
		t.Helper()
		a, err := r.Answer(q.String(), "")
		if err != nil {
			t.Fatal(err)
		}
		if a.OK {
			return "ok " + a.Decision + ": " + a.Statement
		}
		if a.Decision != "" || a.Statement != "" || a.Why == "" {
			t.Errorf("not OK answer %+v", a)
		}
		return a.Why
	}

	ok := question()
	decide("d-ok", arbiterWording, "The code already uses it", ok)
	if got := answer(ok); got != "ok d-ok: Use SQLite" {
		t.Errorf("valid: %q", got)
	}

	noReason := question()
	decide("d-no-reason", arbiterWording, "", noReason)
	if got := answer(noReason); !strings.Contains(got, "decision d-no-reason answers question") || !strings.HasSuffix(got, "but has no reasoning") {
		t.Errorf("no reasoning: %q", got)
	}
	blank := question()
	decide("d-blank", arbiterWording, "   ", blank)
	if got := answer(blank); !strings.HasSuffix(got, "but has no reasoning") {
		t.Errorf("blank reasoning: %q", got)
	}

	wording := question()
	decide("d-wording", "the owner's words", "Because", wording)
	if got := answer(wording); !strings.Contains(got, `its wording is "the owner's words", not "(assumed, would have escalated)"`) {
		t.Errorf("wrong wording: %q", got)
	}

	noLink := question()
	decide("d-no-link", arbiterWording, "Because", id.ID{})
	if got, want := answer(noLink), "no live Decision answers question "+noLink.String(); got != want {
		t.Errorf("no link: %q, want %q", got, want)
	}

	retired := question()
	decide("d-retired", arbiterWording, "Because", retired)
	mustWrite(t, wctx, "retire", `{"ref":"d-retired","reason":"Wrong"}`)
	if got, want := answer(retired), "no live Decision answers question "+retired.String(); got != want {
		t.Errorf("retired: %q, want %q", got, want)
	}

	// A valid Decision among invalid ones answers.
	mixed := question()
	decide("d-mixed-bad", "", "Because", mixed)
	decide("d-mixed-good", arbiterWording, "Because", mixed)
	if got := answer(mixed); got != "ok d-mixed-good: Use SQLite" {
		t.Errorf("mixed: %q", got)
	}

	unknown, err := id.Parse("01M5ZZZZZZZZZZZZZZZZZZZZZZ")
	if err != nil {
		t.Fatal(err)
	}
	if got := answer(unknown); !strings.HasSuffix(got, "is not a live node") {
		t.Errorf("unknown question: %q", got)
	}
}

// TestRunStartChecksArbiter checks an autonomous run's arbiter, built-in
// or declared, before anything is written. A run that is not autonomous
// never launches it, nor does one with no agent step to ask, so their
// arbiter is not checked.
func TestRunStartChecksArbiter(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "cmd", "[step.a]\nkind = \"command\"\nrun = \"true\"\n")
	writeWorkflow(t, f.ctx.Root, "moded", agentWorkflow+"permission_mode = \"manual\"\n")
	writeWorkflow(t, f.ctx.Root, "bad", agentWorkflow+"[step.decide]\nkind = \"arbiter\"\nagent = \"nope\"\n")
	for _, tt := range []struct {
		name, req, want string
	}{
		{"builtin needs a mode", `{"workflow":"moded","autonomous":true}`,
			`step "arbiter" is the arbiter step with no permission mode`},
		{"declared arbiter's adapter", `{"workflow":"bad","permission_mode":"manual","autonomous":true}`,
			`step "decide": unknown agent "nope"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := callOp(t, f.ctx, "run.start", tt.req)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
	f.ctx.DryRun = true
	for _, req := range []string{
		`{"workflow":"cmd","autonomous":true}`, `{"workflow":"moded"}`, `{"workflow":"bad","permission_mode":"manual"}`,
	} {
		if _, err := callOp(t, f.ctx, "run.start", req); err != nil {
			t.Errorf("%s: %v", req, err)
		}
	}
}

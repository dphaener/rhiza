package api

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
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

// holdPane gives r's run a session with a pane, block:a1 of window w9, and
// forgets the rex calls so far.
func holdPane(t *testing.T, f *sessionFixture, r *stepRunner) {
	t.Helper()
	if err := r.sessions.Write(session.Record{ID: "S-first", Run: r.run, Change: r.change,
		Tab: session.Tab{Window: "window:w9", Block: "block:a1"}}); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
}

// An arbiter's session opens its own tab and splits nothing.
func TestStartArbiterOwnTab(t *testing.T) {
	f, r, wctx := arbiterFixture(t)
	q := mustWrite(t, wctx, "add", `{"kind":"question","handle":"which-store","fields":["statement=Which store?"]}`).Nodes[0]
	holdPane(t, f, r)
	step := workflow.Step{Name: workflow.ArbiterName, Kind: workflow.KindArbiter, Interval: workflow.DefaultInterval}
	if _, err := r.StartArbiter(step, q.String(), nil); err != nil {
		t.Fatal(err)
	}
	wantOwnTab(t, f)
}

// An advisor's session opens its own tab and splits nothing.
func TestStartAdvisorOwnTab(t *testing.T) {
	f, r, _ := arbiterFixture(t)
	holdPane(t, f, r)
	step := workflow.Step{Name: workflow.ArbiterName, Kind: workflow.KindArbiter, Interval: workflow.DefaultInterval}
	if _, err := r.StartAdvisor(step, "S-first", "Which store?", nil); err != nil {
		t.Fatal(err)
	}
	wantOwnTab(t, f)
}

// run recommend closes the advisor's pane once the recommendation is
// journalled; a dry run journals nothing and closes nothing.
func TestRunRecommendClosesAdvisorPane(t *testing.T) {
	f, rs := signoffRaceFixture(t)
	if err := session.NewStore(session.Dir(f.common)).Write(session.Record{ID: "S-advisor", Run: signoffRaceRun, Step: "review", Advises: "S-asker",
		Tab: session.Tab{Window: "window:w9", Block: "block:a2"}}); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
	// Each rex call sees whether the recommendation is journalled yet.
	var journalledAtClose []bool
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		j, err := rs.Read(signoffRaceRun)
		if err != nil {
			t.Fatal(err)
		}
		journalledAtClose = append(journalledAtClose, slices.ContainsFunc(j.Entries, func(e run.Entry) bool { return e.Type == run.TypeRecommended }))
		return f.rex.run(name, args, dir, env)
	}
	body := `{"session":"S-advisor","recommendation":"Use blue","reasons":"Brand"}`
	f.ctx.DryRun = true
	if _, err := callOp(t, f.ctx, "run.recommend", body); err != nil {
		t.Fatal(err)
	}
	f.ctx.DryRun = false
	if len(f.rex.calls) != 0 {
		t.Fatalf("dry run ran rex: %q", f.rex.calls)
	}
	if _, err := callOp(t, f.ctx, "run.recommend", body); err != nil {
		t.Fatal(err)
	}
	j, err := rs.Read(signoffRaceRun)
	if err != nil {
		t.Fatal(err)
	}
	var journalled bool
	for _, e := range j.Entries {
		journalled = journalled || e.Type == run.TypeRecommended
	}
	if !journalled {
		t.Error("recommendation not journalled")
	}
	if len(f.rex.calls) != 1 || !slices.Equal(f.rex.calls[0][1:3], []string{"block", "close"}) || f.rex.calls[0][len(f.rex.calls[0])-1] != "block:a2" {
		t.Errorf("rex calls %q, want one block close of block:a2", f.rex.calls)
	}
	if !slices.Equal(journalledAtClose, []bool{true}) {
		t.Errorf("recommendation journalled at each rex call: %v, want it journalled before the close", journalledAtClose)
	}
}

// arbiterRexRunner is an autonomous run's step runner whose agent and
// arbiter sessions are scripted, while ClosePane is the step runner's own:
// a `rex block close` of the session record's block, through the fake rex.
// Session S-asker blocks on a question once, then exits 0; arbiter
// session S-arb polls through arb, then says done, and Answer is answer.
type arbiterRexRunner struct {
	*stepRunner
	arb    []run.AgentStatus
	answer run.Answer
	polls  map[string]int
	// closedWhileRunning counts the block closes seen while the arbiter
	// was still running (polled, not yet said done).
	closedWhileRunning int
	rex                *fakeExec
	// lost loses the arbiter's session (AgentStatus fails) once arb runs
	// out, while it is still running.
	lost bool
}

func (r *arbiterRexRunner) StartAgent(workflow.Step, string, []string) (string, error) {
	return "S-asker", nil
}

func (r *arbiterRexRunner) StartArbiter(_ workflow.Step, _ string, _ []string) (string, error) {
	return "S-arb", r.sessions.Write(session.Record{ID: "S-arb", Run: r.run, Change: r.change, Step: workflow.ArbiterName,
		Tab: session.Tab{Window: "window:w9", Block: "block:arb"}})
}

func (r *arbiterRexRunner) AgentStatus(sid string) (run.AgentStatus, error) {
	n := r.polls[sid]
	r.polls[sid]++
	if sid == "S-asker" {
		if n == 0 {
			return run.AgentStatus{Blocked: &run.Block{Kind: "question", Question: "Which store?"}}, nil
		}
		return run.AgentStatus{Done: true, Outcome: workflow.OutcomePass, Detail: "exit 0"}, nil
	}
	r.closedWhileRunning += len(blockCloses(r.rex))
	if n < len(r.arb) {
		return r.arb[n], nil
	}
	if r.lost {
		return run.AgentStatus{}, fmt.Errorf("session %s is gone", sid)
	}
	return run.AgentStatus{SaidDone: true}, nil
}

func (r *arbiterRexRunner) AddQuestion(string, string) (run.Question, error) {
	return run.Question{ID: "Q1", Handle: "which-store"}, nil
}

func (r *arbiterRexRunner) Answer(string, string) (run.Answer, error)    { return r.answer, nil }
func (r *arbiterRexRunner) Deliver(string, string, string, string) error { return nil }
func (r *arbiterRexRunner) AgentEnded(workflow.Step, string) error       { return nil }
func (r *arbiterRexRunner) Holds(workflow.Step, []string) (bool, error)  { return false, nil }
func (r *arbiterRexRunner) StartAdvisor(workflow.Step, string, string, []string) (string, error) {
	return "", fmt.Errorf("no advisor in an autonomous run")
}

// blockCloses is each `rex block close` call f ran.
func blockCloses(f *fakeExec) [][]string {
	var out [][]string
	for _, c := range f.calls {
		if len(c) > 2 && c[1] == "block" && c[2] == "close" {
			out = append(out, c)
		}
	}
	return out
}

// An autonomous run's arbiter, driven through the engine to its verdict,
// has its pane closed with `rex block close` once the verdict is
// journalled, pass or fail, never while it is still running; one whose
// session is lost while still running keeps its pane.
func TestArbiterPaneClosesAfterVerdict(t *testing.T) {
	const wf = "[step.build]\nkind = \"agent\"\nprompt = \"build\"\non.pass = \"done\"\n"
	running := run.AgentStatus{}
	for _, tt := range []struct {
		name    string
		answer  run.Answer
		lost    bool // the arbiter's session is lost while still running
		closed  bool
		outcome string
	}{
		{"verdict", run.Answer{OK: true, Decision: "use-sqlite", Statement: "Use SQLite"}, false, true, "succeeded"},
		{"failed", run.Answer{Why: "no live Decision answers question Q1"}, false, true, "failed"},
		{"still running", run.Answer{OK: true, Decision: "use-sqlite"}, true, false, "failed"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, sr, _ := arbiterFixture(t)
			def, err := workflow.Parse("wf", []byte(wf))
			if err != nil {
				t.Fatal(err)
			}
			rs := run.NewStore(run.Dir(f.common))
			if err := run.Start(rs, sr.run, run.Started{Workflow: "wf", Definition: wf, Change: sr.change,
				Worktree: f.worktree, PermissionMode: "manual", Autonomous: true}, time.Unix(100, 0), def); err != nil {
				t.Fatal(err)
			}
			f.rex.calls = nil
			// Each block close sees whether the verdict is journalled yet.
			var journalledAtClose []bool
			f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
				if len(args) > 1 && args[0] == "block" && args[1] == "close" {
					j, err := rs.Read(sr.run)
					if err != nil {
						t.Fatal(err)
					}
					a, ok := j.State().Arbitration("Q1")
					journalledAtClose = append(journalledAtClose, ok && !a.Open())
				}
				return f.rex.run(name, args, dir, env)
			}
			r := &arbiterRexRunner{stepRunner: sr, arb: []run.AgentStatus{running, {Idle: true}},
				answer: tt.answer, polls: map[string]int{}, rex: f.rex, lost: tt.lost}
			e := &run.Engine{Store: rs, Runner: r, Clock: noSleep{}, Now: func() time.Time { return time.Unix(200, 0) }}
			st, err := e.Run(sr.run)
			if err != nil {
				t.Fatal(err)
			}
			if st.Result != tt.outcome {
				t.Errorf("run ended %q (%s), want %q", st.Result, st.Detail, tt.outcome)
			}
			if r.closedWhileRunning != 0 {
				t.Errorf("the arbiter's pane closed while it was still running: %q", f.rex.calls)
			}
			closes := blockCloses(f.rex)
			if !tt.closed {
				if len(closes) != 0 {
					t.Errorf("rex calls %q, want no block close", closes)
				}
				return
			}
			if len(closes) != 1 || closes[0][len(closes[0])-1] != "block:arb" {
				t.Fatalf("rex calls %q, want one block close of block:arb", f.rex.calls)
			}
			if !slices.Equal(journalledAtClose, []bool{true}) {
				t.Errorf("verdict journalled at each close: %v, want it journalled before the one close", journalledAtClose)
			}
		})
	}
}

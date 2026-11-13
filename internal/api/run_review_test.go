package api

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A check or critic session's RHIZA_CYCLE names its round's cycle when it
// is one of the step's cycles that runs the check; anything else is 0, so
// the cycle is inferred.
func TestEnvCheckCycle(t *testing.T) {
	st := run.State{Cycles: []run.CheckCycle{
		{Step: "review", Cycle: 1, Checks: []string{"alpha", "beta"}},
		{Step: "review", Cycle: 2, Checks: []string{"beta"}},
	}}
	for _, tt := range []struct {
		env, check string
		want       int
	}{
		{"1", "alpha", 1},
		{" 2 ", "beta", 2},
		{"2", "alpha", 0},
		{"3", "beta", 0},
		{"0", "alpha", 0},
		{"x", "alpha", 0},
		{"", "alpha", 0},
	} {
		ctx := &Ctx{Getenv: func(k string) string {
			if k == envCycle {
				return tt.env
			}
			return ""
		}}
		if got := envCheckCycle(ctx, st, "review", tt.check); got != tt.want {
			t.Errorf("RHIZA_CYCLE=%q for %s: %d, want %d", tt.env, tt.check, got, tt.want)
		}
	}
}

// Two run IDs that share their first 8 characters (the same start time)
// get different question handles carrying the full IDs.
func TestReviewQuestionHandleFullID(t *testing.T) {
	a := "01M6VY7E0AAAAAAAAAAAAAAAAA"
	b := "01M6VY7E0BBBBBBBBBBBBBBBBB"
	ha, hb := reviewQuestionHandle(a, 3, "f1"), reviewQuestionHandle(b, 3, "f1")
	if ha == hb {
		t.Fatalf("handles collide: %s", ha)
	}
	if want := "review-01m6vy7e0aaaaaaaaaaaaaaaaa-r3-f1"; ha != want {
		t.Errorf("handle = %s, want %s", ha, want)
	}
}

// The asked entry for a finding is found by its journalled handle, the base
// or base-N, and the autonomous gate takes it for the question left.
func TestReviewAskedFor(t *testing.T) {
	base := reviewQuestionHandle("01M6VY7E0AAAAAAAAAAAAAAAAA", 3, "f1")
	other := reviewQuestionHandle("01M6VY7E0BBBBBBBBBBBBBBBBB", 3, "f1")
	asked := []run.Asked{{Handle: other, Question: "x"}, {Handle: base + "-2", Question: "q2"}}
	a, ok := reviewAskedFor(asked, base)
	if !ok || a.Question != "q2" {
		t.Fatalf("reviewAskedFor = %v, %v; want the -2 entry", a, ok)
	}
	if _, ok := reviewAskedFor(asked[:1], base); ok {
		t.Error("another run's entry matched")
	}
	if _, ok := reviewAskedFor([]run.Asked{{Handle: base + "-x"}}, base); ok {
		t.Error("a non-numeric suffix matched")
	}
	st := run.State{ID: "01M6VY7E0AAAAAAAAAAAAAAAAA", Asked: []run.Asked{{Handle: base + "-2", Question: "01M6W7QGXS2VNCJ3KMFRMW78XE"}},
		Rounds: []run.Round{{Number: 3, Findings: []run.Finding{{ID: "f1", Node: "n"}}, Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed}}}}}
	m := newModel(store.Log{}, id.ID{})
	open, err := gateReview(nil, m, reviewRun{id: st.ID, state: st})
	if err != nil {
		t.Fatalf("gateReview with the -2 question: %v", err)
	}
	if !open {
		t.Error("gateReview: the unanswered -2 question is not open")
	}
	st.Asked = []run.Asked{{Handle: other, Question: "01M6W7QGXS2VNCJ3KMFRMW78XE"}}
	if _, err := gateReview(nil, m, reviewRun{id: st.ID, state: st}); err == nil {
		t.Error("gateReview took another run's question for the finding")
	}
}

// review-check, and its --fresh test, refuse a round journaled before the
// latest revise visit, adding nothing, even when the round is 3.
func TestReviewCheckRefusesStaleRound(t *testing.T) {
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R8", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Plan it"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R8", e); err != nil {
			t.Fatal(err)
		}
	}
	round := func() {
		app(run.Entry{Type: run.TypeEntered, Step: "review"})
		app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{{ID: "f1", Node: "n"}}})
		app(run.Entry{Type: run.TypeVerdict, Step: "critic", Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed}}})
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R8"
		}
		return ""
	}
	round()
	round()
	app(run.Entry{Type: run.TypeEntered, Step: "revise"})
	for n := 2; n <= 3; n++ {
		if n == 3 {
			// A third stale round: the revise follows round 3.
			app(run.Entry{Type: run.TypeEntered, Step: "review"})
			app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{{ID: "f1", Node: "n"}}})
			app(run.Entry{Type: run.TypeVerdict, Step: "critic", Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed}}})
			app(run.Entry{Type: run.TypeEntered, Step: "revise"})
		}
		for _, req := range []string{`{}`, `{"fresh":true}`} {
			_, err := callOp(t, f.ctx, "run.review-check", req)
			e := AsError(err)
			if e == nil || !strings.Contains(e.Message, "round "+strconv.Itoa(n)) || !strings.Contains(e.Message, "predates the latest revise") {
				t.Fatalf("round %d, %s: %v", n, req, err)
			}
		}
		st, _ := store.Read("R8")
		if len(st.State().Asked) != 0 {
			t.Fatalf("round %d: a stale check asked a question", n)
		}
	}
	// A round after the revise visit is fresh.
	app(run.Entry{Type: run.TypeEntered, Step: "review"})
	app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}})
	if _, err := callOp(t, f.ctx, "run.review-check", `{"fresh":true}`); err != nil {
		t.Fatalf("--fresh after a fresh round: %v", err)
	}
	if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err != nil {
		t.Fatalf("review-check of a fresh round: %v", err)
	}
}

// A critic may downgrade a blocker or major finding to a lower severity,
// journalled with the finding's original; a downgrade to the same or a
// higher severity, to none, or of a finding that is not serious is usage.
func TestRunVerdictDowngrade(t *testing.T) {
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	def, ok := workflow.Builtin("review")
	if !ok {
		t.Fatal("no built-in review workflow")
	}
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R9", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "review", Definition: string(def), Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Review it"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R9", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "review"})
	app(run.Entry{Type: run.TypeFindings, Step: "review", Check: "quality", Findings: []run.Finding{
		{ID: "f1", Check: "quality", Severity: run.SeverityMajor, Reason: "a"},
		{ID: "f2", Check: "quality", Severity: run.SeverityBlocker, Reason: "b"},
		{ID: "f3", Check: "quality", Severity: run.SeverityMinor, Reason: "c"},
	}})
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R9"
		}
		return ""
	}
	verdict := func(body string) error {
		f.ctx.Stdin = strings.NewReader(body)
		_, err := callOp(t, f.ctx, "run.verdict", `{"file":"-","check":"quality"}`)
		return err
	}
	for name, body := range map[string]string{
		"to the same severity":   `{"verdicts":[{"finding":"f1","verdict":"downgraded","severity":"major","reason":"r"}]}`,
		"to a higher severity":   `{"verdicts":[{"finding":"f1","verdict":"downgraded","severity":"blocker","reason":"r"}]}`,
		"of a blocker to major+": `{"verdicts":[{"finding":"f2","verdict":"downgraded","severity":"blocker","reason":"r"}]}`,
		"with no severity":       `{"verdicts":[{"finding":"f1","verdict":"downgraded","reason":"r"}]}`,
		"of a minor":             `{"verdicts":[{"finding":"f3","verdict":"downgraded","severity":"nit","reason":"r"}]}`,
		"severity on confirmed":  `{"verdicts":[{"finding":"f1","verdict":"confirmed","severity":"minor","reason":"r"}]}`,
	} {
		err := verdict(body)
		if e := AsError(err); e == nil || e.Code != CodeUsage {
			t.Errorf("downgrade %s: err = %v, want usage", name, err)
		}
	}
	if err := verdict(`{"verdicts":[{"finding":"f1","verdict":"downgraded","severity":"minor","reason":"small"},{"finding":"f2","verdict":"downgraded","severity":"major","reason":"less"}]}`); err != nil {
		t.Fatal(err)
	}
	st, err := store.Read("R9")
	if err != nil {
		t.Fatal(err)
	}
	cs, _ := st.State().Check("review", "quality")
	rd, _ := cs.LatestRound()
	want := []run.Verdict{
		{Finding: "f1", Verdict: run.Downgraded, Reason: "small", Severity: "minor", Original: "major"},
		{Finding: "f2", Verdict: run.Downgraded, Reason: "less", Severity: "major", Original: "blocker"},
	}
	if len(rd.Verdicts) != 2 || rd.Verdicts[0] != want[0] || rd.Verdicts[1] != want[1] {
		t.Fatalf("verdicts = %+v, want %+v", rd.Verdicts, want)
	}
	if s := rd.Standing(); len(s) != 1 || s[0].ID != "f2" {
		t.Fatalf("standing = %+v, want only f2 (a major still stands)", s)
	}
}

// A plan-review finding carries its kind's severity: a plan critic may
// downgrade a major kind (missing-acceptance) but not a minor one
// (missing-provenance), and review-check and the gate count only confirmed
// major findings, so a confirmed minor one passes.
func TestPlanReviewDowngradeByKind(t *testing.T) {
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	rs := run.NewStore(run.Dir(f.common))
	if err := rs.Create("R12", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Definition: "name = \"plan\"\n[step.review]\nkind = \"command\"\nrun = \"true\"\n", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Plan it"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := rs.Append("R12", e); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R12"
		}
		return ""
	}
	app(run.Entry{Type: run.TypeEntered, Step: "review"})
	app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{
		{ID: "f1", Node: "n", Kind: run.KindMissingAcceptance, Reason: "a"},
		{ID: "f2", Node: "n", Kind: run.KindMissingProvenance, Reason: "b"},
	}})
	verdict := func(body string) error {
		f.ctx.Stdin = strings.NewReader(body)
		_, err := callOp(t, f.ctx, "run.verdict", `{"file":"-"}`)
		return err
	}
	if e := AsError(verdict(`{"verdicts":[{"finding":"f2","verdict":"downgraded","severity":"nit","reason":"r"}]}`)); e == nil || e.Code != CodeUsage {
		t.Fatalf("downgrade of a minor kind: %v, want usage", e)
	}
	if err := verdict(`{"verdicts":[{"finding":"f1","verdict":"downgraded","severity":"minor","reason":"small"},{"finding":"f2","verdict":"confirmed","reason":"real"}]}`); err != nil {
		t.Fatalf("downgrade of a major kind: %v", err)
	}
	st, _ := rs.Read("R12")
	rd, _ := st.State().LatestRound()
	if v, _ := rd.Verdict("f1"); v.Verdict != run.Downgraded || v.Original != run.SeverityMajor || v.Severity != run.SeverityMinor {
		t.Fatalf("f1 verdict = %+v, want downgraded major to minor", v)
	}
	if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err != nil {
		t.Fatalf("review-check with a downgrade and a confirmed minor: %v", err)
	}
	m := newModel(store.Log{}, id.ID{})
	if _, err := gateReview(nil, m, reviewRun{id: "R12", state: st.State()}); err != nil {
		t.Fatalf("gate with no confirmed major: %v", err)
	}
	// A confirmed major kind still fails the check.
	app(run.Entry{Type: run.TypeVerdict, Step: "critic", Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed, Reason: "r"}}})
	if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err == nil || !strings.Contains(err.Error(), "f1") {
		t.Fatalf("review-check with a confirmed missing-acceptance: %v", err)
	}
}

// A run in an agent step declaring checks plan, code and qa journals
// findings and verdicts under --check plan as that check's rounds, and
// refuses --check scope with a message naming the declared checks.
func TestRunFindingsCheckNames(t *testing.T) {
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	def := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
		"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"code\"\n[[step.conduct.check]]\nname = \"qa\"\n"
	if err := store.Create("R12", run.Entry{
		Type: run.TypeStarted, Step: "conduct", At: t0,
		Started: &run.Started{Workflow: "conduct", Definition: def, Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Conduct"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("R12", run.Entry{Type: run.TypeEntered, Step: "conduct", Visit: 1, At: t0}); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R12"
		}
		return ""
	}
	f.ctx.Stdin = strings.NewReader(`{"findings":[{"severity":"major","place":"a.go:1","reason":"r"}]}`)
	resp, err := callOp(t, f.ctx, "run.findings", `{"check":"plan","file":"-"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RunFindingsResponse); r.Check != "plan" || r.Round != 1 || len(r.Findings) != 1 {
		t.Errorf("findings %+v", r)
	}
	f.ctx.Stdin = strings.NewReader(`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"real"}]}`)
	resp, err = callOp(t, f.ctx, "run.verdict", `{"check":"plan","file":"-"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RunVerdictResponse); r.Check != "plan" || r.Round != 1 {
		t.Errorf("verdict %+v", r)
	}
	for op, body := range map[string]string{"run.findings": `{"findings":[]}`, "run.verdict": `{"verdicts":[]}`} {
		f.ctx.Stdin = strings.NewReader(body)
		_, err = callOp(t, f.ctx, op, `{"check":"scope","file":"-"}`)
		e := AsError(err)
		if e == nil || !strings.Contains(e.Message, "plan, code, qa") {
			t.Errorf("%s --check scope: %v", op, err)
		}
	}
}

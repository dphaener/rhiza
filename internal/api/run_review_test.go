package api

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A check or critic session's RHIZA_CYCLE names its round's cycle when it
// is one of the step's cycles that runs the check; anything else is 0, so
// the cycle is inferred.
func TestEnvCheckCycle(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	m := projectionModel(t, nil, id.ID{})
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
	t.Parallel()
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R8", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Change: currentOf(t, f.ctx).String(), Title: "Plan it"},
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
	t.Parallel()
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
	t.Parallel()
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	rs := run.NewStore(run.Dir(f.common))
	if err := rs.Create("R12", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Definition: "name = \"plan\"\n[step.review]\nkind = \"command\"\nrun = \"true\"\n", Change: currentOf(t, f.ctx).String(), Title: "Plan it"},
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
	m := projectionModel(t, nil, id.ID{})
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
// refuses --check scope with a message naming the declared checks. An
// agent step that declares no checks refuses --check.
func TestRunFindingsAgentStepChecks(t *testing.T) {
	t.Parallel()
	f, sess, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	def := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
		"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"code\"\n[[step.conduct.check]]\nname = \"qa\"\n"
	if err := store.Create("R12", run.Entry{
		Type: run.TypeStarted, Step: "conduct", At: t0,
		Started: &run.Started{Workflow: "conduct", Definition: def, Change: sess.Change, Title: "Conduct"},
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
	// An agent step that declares no checks refuses --check.
	if err := store.Create("R12b", run.Entry{
		Type: run.TypeStarted, Step: "plain", At: t0,
		Started: &run.Started{Workflow: "plain", Definition: "[step.plain]\nkind = \"agent\"\nprompt = \"go\"\n", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Plain"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("R12b", run.Entry{Type: run.TypeEntered, Step: "plain", Visit: 1, At: t0}); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R12b"
		}
		return ""
	}
	for op, body := range map[string]string{"run.findings": `{"findings":[]}`, "run.verdict": `{"verdicts":[]}`} {
		f.ctx.Stdin = strings.NewReader(body)
		_, err = callOp(t, f.ctx, op, `{"check":"plan","file":"-"}`)
		if e := AsError(err); e == nil || !strings.Contains(e.Message, "not in a checks step") {
			t.Errorf("%s --check at an agent step with no checks: %v", op, err)
		}
	}
}

// In a conductor run (conduct declares plan, code and qa) review-check,
// run from the later gate step, counts each check's latest round under
// conduct: a declared check with no round fails, minor rounds pass, a blocker with no verdict or a confirmed
// major fails, a downgraded blocker or dropped major passes.
func TestReviewCheckConductor(t *testing.T) {
	t.Parallel()
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	def := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
		"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"code\"\n[[step.conduct.check]]\nname = \"qa\"\n" +
		"[step.gate]\nkind = \"command\"\nrun = \"true\"\n"
	if err := store.Create("R13", run.Entry{
		Type: run.TypeStarted, Step: "conduct", At: t0,
		Started: &run.Started{Workflow: "conduct", Definition: def, Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Conduct"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R13", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "conduct", Visit: 1})
	app(run.Entry{Type: run.TypeEntered, Step: "gate", Visit: 1})
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R13"
		}
		return ""
	}
	findings := func(check string, fs ...run.Finding) {
		app(run.Entry{Type: run.TypeFindings, Step: "conduct", Check: check, Cycle: 1, Findings: fs})
	}
	verdict := func(check string, v run.Verdict) {
		app(run.Entry{Type: run.TypeVerdict, Step: "conduct", Check: check, Cycle: 1, Verdicts: []run.Verdict{v}})
	}
	fd := func(id, sev string) run.Finding {
		return run.Finding{ID: id, Severity: sev, Place: "a.go:1", Reason: "r " + id}
	}
	pass := func(what string) {
		t.Helper()
		if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err != nil {
			t.Fatalf("%s: %v", what, err)
		}
	}
	fail := func(what string, names ...string) {
		t.Helper()
		_, err := callOp(t, f.ctx, "run.review-check", `{}`)
		e := AsError(err)
		if e == nil || strings.Contains(e.Message, "has no review round yet") {
			t.Fatalf("%s: %v", what, err)
		}
		for _, n := range names {
			if !strings.Contains(e.Message, n) {
				t.Errorf("%s: %q missing from %v", what, n, e.Message)
			}
		}
	}
	// A declared check with no round fails: the conductor cannot skip one.
	fail("no rounds", "plan", "code", "qa")
	findings("plan", fd("f1", "minor"), fd("f2", "nit"))
	findings("code")
	fail("qa has no round", "qa")
	findings("qa", fd("f1", "minor"))
	pass("minor rounds")
	findings("code", fd("f1", "blocker"))
	fail("blocker with no verdict", "code", "r f1")
	findings("code", fd("f1", "major"))
	verdict("code", run.Verdict{Finding: "f1", Verdict: run.Confirmed, Reason: "real"})
	fail("confirmed major", "code")
	findings("code", fd("f1", "blocker"))
	verdict("code", run.Verdict{Finding: "f1", Verdict: run.Downgraded, Severity: "minor", Original: "blocker", Reason: "small"})
	pass("downgraded blocker")
	findings("qa", fd("f1", "major"))
	verdict("qa", run.Verdict{Finding: "f1", Verdict: run.Dropped, Reason: "no"})
	pass("dropped major")
	// Only the latest round counts: an old major is superseded.
	findings("plan", fd("f1", "major"))
	fail("new major in plan", "plan")
	findings("plan")
	pass("latest plan round clean")
}

// review-check on a run past an agent step declaring plan, code and qa
// passes with minor findings only, fails naming a confirmed major of the
// code check or a declared check with no round, and passes once the critic
// downgraded that major to minor or dropped it.
func TestReviewCheckDeclaredChecks(t *testing.T) {
	t.Parallel()
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	def := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
		"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"code\"\n[[step.conduct.check]]\nname = \"qa\"\n" +
		"[step.gate]\nkind = \"command\"\nrun = \"true\"\n"
	if err := store.Create("R14", run.Entry{
		Type: run.TypeStarted, Step: "conduct", At: t0,
		Started: &run.Started{Workflow: "conduct", Definition: def, Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Conduct"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R14", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "conduct", Visit: 1})
	app(run.Entry{Type: run.TypeEntered, Step: "gate", Visit: 1})
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R14"
		}
		return ""
	}
	round := func(check string, fs ...run.Finding) {
		app(run.Entry{Type: run.TypeFindings, Step: "conduct", Check: check, Cycle: 1, Findings: fs})
	}
	verdict := func(v run.Verdict) {
		app(run.Entry{Type: run.TypeVerdict, Step: "conduct", Check: "code", Cycle: 1, Verdicts: []run.Verdict{v}})
	}
	review := func() (any, error) { return callOp(t, f.ctx, "run.review-check", `{}`) }

	if _, err := review(); err == nil {
		t.Error("declared checks with no round passed")
	}
	minor := run.Finding{ID: "f1", Severity: "minor", Place: "a.go:1", Reason: "small thing"}
	round("plan", minor)
	round("code", minor)
	round("qa")
	if _, err := review(); err != nil {
		t.Fatalf("minor findings only: %v", err)
	}
	major := run.Finding{ID: "f1", Severity: "major", Place: "b.go:2", Reason: "real defect"}
	round("code", major)
	verdict(run.Verdict{Finding: "f1", Verdict: run.Confirmed, Reason: "real"})
	_, err := review()
	e := AsError(err)
	if e == nil || ExitCode(err) != 1 || !strings.Contains(e.Message, "real defect") {
		t.Fatalf("confirmed major: %v", err)
	}
	round("code", major)
	verdict(run.Verdict{Finding: "f1", Verdict: run.Downgraded, Severity: "minor", Original: "major", Reason: "small"})
	if _, err := review(); err != nil {
		t.Errorf("downgraded major: %v", err)
	}
	round("code", major)
	verdict(run.Verdict{Finding: "f1", Verdict: run.Dropped, Reason: "no"})
	if _, err := review(); err != nil {
		t.Errorf("dropped major: %v", err)
	}
}

// On the plan-review path review-check fails, at every round, while a live
// leaf the plan creates or changes lacks statement, acceptance, files or
// must_not, or a failure case in its acceptance with no failure_waiver,
// naming the leaf and what it lacks; a clean round over leaves with all of
// them passes. At round 3 the failure still fails and adds no question. A
// parent capability and a conductor run's review-check are not checked.
func TestReviewCheckLeafContract(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "leafcontract")
	a := newChange(t, f.ctx, `{"intent":"Plan A"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"parent","fields":["statement=Holds the leaves"]}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-a","parent":"parent","fields":["statement=Does A","acceptance=go test passes; (Failure Case) a bad input is refused","files=internal/a.go","must_not=touch b"]}`)
	t0 := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	runID := "R20"
	start := func(rid, def, step string) {
		t.Helper()
		if err := store.Create(rid, run.Entry{
			Type: run.TypeStarted, Step: step, At: t0,
			Started: &run.Started{Workflow: "plan", Change: a.ChangeSet.String(), Title: "Plan it", Definition: def},
		}); err != nil {
			t.Fatal(err)
		}
		runID = rid
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append(runID, e); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return runID
		}
		return ""
	}
	cleanRound := func(visit int) {
		app(run.Entry{Type: run.TypeEntered, Step: "review", Visit: visit})
		app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}})
	}
	set := func(fields string) {
		t.Helper()
		mustWrite(t, f.ctx, "set", `{"ref":"leaf-a",`+fields+`}`)
	}
	pass := func(what string) {
		t.Helper()
		if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err != nil {
			t.Fatalf("%s: %v", what, err)
		}
	}
	fail := func(what string, names ...string) {
		t.Helper()
		_, err := callOp(t, f.ctx, "run.review-check", `{}`)
		e := AsError(err)
		if e == nil {
			t.Fatalf("%s: passed, want a failure", what)
		}
		for _, n := range names {
			if !strings.Contains(e.Message, n) {
				t.Errorf("%s: %q missing from %v", what, n, e.Message)
			}
		}
		if strings.Contains(e.Message, "parent (") {
			t.Errorf("%s: the parent capability was checked: %v", what, e.Message)
		}
	}

	plan := "name = \"plan\"\n[step.review]\nkind = \"command\"\nrun = \"true\"\n"
	start("R20", plan, "review")
	cleanRound(1)
	pass("a clean round over a leaf with its whole contract")

	set(`"unset":["files"]`)
	fail("a leaf missing files", "leaf-a", "lacks files")
	set(`"fields":["files=internal/a.go","acceptance=go test passes"]`)
	fail("an acceptance with no failure case", "leaf-a", "a failure case")
	set(`"fields":["failure_waiver=A pure rename has no failure path"]`)
	pass("the same leaf with a failure_waiver")

	// Round 3: the leaf-contract failure still fails, leaving no question.
	set(`"unset":["failure_waiver","must_not"]`)
	cleanRound(2)
	app(run.Entry{Type: run.TypeEntered, Step: "review", Visit: 3})
	app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{{ID: "f1", Node: "leaf-a", Kind: run.KindContradiction, Reason: "wrong"}}})
	app(run.Entry{Type: run.TypeVerdict, Step: "critic", Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed, Reason: "right"}}})
	fail("round 3 with a leaf lacking its contract", "leaf-a", "must_not", "a failure case")
	st, _ := store.Read(runID)
	if n := len(st.State().Asked); n != 0 {
		t.Fatalf("round 3 leaf-contract failure asked %d question(s)", n)
	}
	if rd, _ := st.State().LatestRound(); rd.Number != 3 {
		t.Fatalf("latest round %d, want 3", rd.Number)
	}

	// A conductor run on the same plan reads only its checks' rounds.
	conduct := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n[[step.conduct.check]]\nname = \"plan\"\n" +
		"[step.gate]\nkind = \"command\"\nrun = \"true\"\n"
	start("R21", conduct, "conduct")
	app(run.Entry{Type: run.TypeEntered, Step: "conduct", Visit: 1})
	app(run.Entry{Type: run.TypeFindings, Step: "conduct", Check: "plan", Cycle: 1, Findings: []run.Finding{}})
	pass("a conductor run with a leaf lacking its contract")
}

// A parent capability created in one plan change set whose child
// capability is added only in another (built on the first) is not a leaf:
// review-check reads children across all the plan's change sets and
// passes over the child's whole contract.
func TestReviewCheckLeafContractChildInOtherChangeSet(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "leafcontractsplit")
	a := newChange(t, f.ctx, `{"intent":"Plan A"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"intent","handle":"goal","fields":["statement=Ship the leaves"]}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"parent","fields":["statement=Holds the leaves"]}`)
	b := newChange(t, f.ctx, `{"intent":"Plan B"}`)
	mustWrite(t, f.ctx, "link", `{"from":"`+b.ChangeSet.String()+`","type":"depends-on","to":"`+a.ChangeSet.String()+`"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-b","parent":"parent","fields":["statement=Does B","acceptance=go test passes; failure case: a bad input is refused","files=internal/b.go","must_not=touch a"]}`)
	mustWrite(t, f.ctx, "link", `{"from":"`+b.ChangeSet.String()+`","type":"serves","to":"goal"}`)

	t0 := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R30", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Change: a.ChangeSet.String(), Title: "Plan it",
			Definition: "name = \"plan\"\n[step.review]\nkind = \"command\"\nrun = \"true\"\n"},
	}); err != nil {
		t.Fatal(err)
	}
	for _, e := range []run.Entry{
		{Type: run.TypeEntered, Step: "review", Visit: 1, At: t0},
		{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}, At: t0},
	} {
		if err := store.Append("R30", e); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R30"
		}
		return ""
	}
	if _, err := callOp(t, f.ctx, "run.review-check", `{}`); err != nil {
		t.Fatalf("a parent whose child is in another plan change set: %v", err)
	}
}

// run.review-check --rulings passes while every passed arbiter ruling on a
// review question has a live node linking decided-by to its Decision, and
// fails naming the question and the Decision otherwise. An open question
// and a session's question are not looked at.
func TestRunReviewCheckRulings(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "rulings")
	a := newChange(t, f.ctx, `{"intent":"Plan A"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-a","fields":["statement=Does A"]}`)
	q1 := mustWrite(t, f.ctx, "add", `{"kind":"question","handle":"review-r40-r1-f1","fields":["statement=How?"]}`).Nodes[0]
	mustWrite(t, f.ctx, "add", `{"kind":"decision","handle":"use-z","fields":["statement=Use Z"]}`)
	mustWrite(t, f.ctx, "link", `{"from":"use-z","type":"answers","to":"review-r40-r1-f1"}`)
	q2 := mustWrite(t, f.ctx, "add", `{"kind":"question","handle":"review-r40-r1-f2","fields":["statement=Open?"]}`).Nodes[0]
	q3 := mustWrite(t, f.ctx, "add", `{"kind":"question","handle":"asked-by-session","fields":["statement=Mine?"]}`).Nodes[0]

	t0 := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R40", run.Entry{
		Type: run.TypeStarted, Step: "rulings", At: t0,
		Started: &run.Started{Workflow: "plan", Change: a.ChangeSet.String(), Title: "Plan it",
			Definition: "name = \"plan\"\n[step.rulings]\nkind = \"command\"\nrun = \"true\"\n"},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R40", e); err != nil {
			t.Fatal(err)
		}
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R40"
		}
		return ""
	}
	pass := func(what string, applied ...string) {
		t.Helper()
		resp, err := callOp(t, f.ctx, "run.review-check", `{"rulings":true}`)
		if err != nil {
			t.Fatalf("%s: %v", what, err)
		}
		r := resp.(RunReviewCheckResponse)
		got := []string{}
		for _, rl := range r.Rulings {
			got = append(got, rl.Question)
			if len(rl.By) == 0 {
				t.Errorf("%s: ruling %s applied by no node", what, rl.Question)
			}
		}
		if !r.Passed || strings.Join(got, ",") != strings.Join(applied, ",") {
			t.Fatalf("%s: passed %v, applied %v, want %v", what, r.Passed, got, applied)
		}
	}
	fail := func(what string, names ...string) {
		t.Helper()
		_, err := callOp(t, f.ctx, "run.review-check", `{"rulings":true}`)
		e := AsError(err)
		if e == nil {
			t.Fatalf("%s: passed, want a failure", what)
		}
		if e.Next != "rhiza run status R40" {
			t.Errorf("%s: next %q", what, e.Next)
		}
		for _, n := range names {
			if !strings.Contains(e.Message, n) {
				t.Errorf("%s: %q missing from %v", what, n, e.Message)
			}
		}
	}

	pass("no asked entry")

	app(run.Entry{Type: run.TypeEntered, Step: "review", Visit: 3})
	app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{
		{ID: "f1", Node: "leaf-a", Handle: "leaf-a", Kind: run.KindContradiction, Reason: "wrong"},
		{ID: "f2", Node: "leaf-a", Handle: "leaf-a", Kind: run.KindOutsideIntent, Reason: "out"},
	}})
	app(run.Entry{Type: run.TypeAsked, Step: "check", Question: q1.String(), Handle: "review-r40-r1-f1"})
	app(run.Entry{Type: run.TypeAsked, Step: "check", Question: q2.String(), Handle: "review-r40-r1-f2"})
	pass("open (unarbitrated) review questions")

	// A session's question with a passed, unapplied ruling is not counted.
	app(run.Entry{Type: run.TypeAsked, Step: "plan", Session: "S1", Question: q3.String(), Handle: "asked-by-session"})
	app(run.Entry{Type: run.TypeArbitrated, Step: "plan", Question: q3.String(), Outcome: workflow.OutcomePass, Decision: "use-z"})
	pass("a session's question")

	// The arbiter passes the review question's ruling; with only the answers link
	// it is not applied, so the check fails.
	app(run.Entry{Type: run.TypeArbitrated, Step: "check", Question: q1.String(), Outcome: workflow.OutcomePass, Decision: "use-z"})
	fail("a ruling with only its answers link", "review-r40-r1-f1", "use-z", "leaf-a")

	mustWrite(t, f.ctx, "link", `{"from":"leaf-a","type":"decided-by","to":"use-z"}`)
	pass("a capability linked decided-by the Decision", "review-r40-r1-f1")

	mustWrite(t, f.ctx, "unlink", `{"ref":"leaf-a","type":"decided-by","to":"use-z"}`)
	fail("the only decided-by link retired", "review-r40-r1-f1", "use-z")
	mustWrite(t, f.ctx, "link", `{"from":"leaf-a","type":"decided-by","to":"use-z"}`)

	// A Decision node in a landed record, linked by a record write.
	recordWrite(t, f.ctx, "add", `{"kind":"decision","handle":"keep-flat","fields":["statement=Keep it flat"]}`)
	recordWrite(t, f.ctx, "add", `{"kind":"decision","handle":"nest-one","fields":["statement=Nest one level"]}`)
	app(run.Entry{Type: run.TypeArbitrated, Step: "check", Question: q2.String(), Outcome: workflow.OutcomePass, Decision: "nest-one"})
	fail("the second ruling unapplied", "review-r40-r1-f2", "nest-one")
	recordWrite(t, f.ctx, "link", `{"from":"keep-flat","type":"decided-by","to":"nest-one"}`)
	pass("a landed Decision linked decided-by the ruling", "review-r40-r1-f1", "review-r40-r1-f2")
}

// A round-3 question goes into the change set that holds the finding's
// node (decision review-ruling-in-node-change-set): with the run's change
// set open, a serving draft's leaf gets its question (and blocks link) in
// that draft and the asked entry names the draft; a node the run's change
// set created keeps its question there, the asked entry naming none. With
// the run's change set landed, a finding on a node no in-flight plan
// change set created is still a record write.
func TestReviewCheckQuestionInNodeChangeSet(t *testing.T) {
	t.Parallel()
	contract := `"acceptance=go test passes; failure case: a bad input is refused","files=internal/a.go","must_not=touch b"`
	t0 := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	def := "name = \"plan\"\n[step.check]\nkind = \"command\"\nrun = \"true\"\n"
	// review runs review-check on run rid of change set cs after two clean
	// rounds and a round 3 whose findings on nodes are confirmed.
	// reviewFindings is review with findings given whole (IDs f1, f2, ...
	// and kind, reason set here), such as one on a link.
	var reviewFindings func(t *testing.T, f *sessionFixture, rid string, cs id.ID, findings ...run.Finding) (RunReviewCheckResponse, run.State)
	review := func(t *testing.T, f *sessionFixture, rid string, cs id.ID, nodes ...id.ID) (RunReviewCheckResponse, run.State) {
		t.Helper()
		var findings []run.Finding
		for _, n := range nodes {
			findings = append(findings, run.Finding{Node: n.String()})
		}
		return reviewFindings(t, f, rid, cs, findings...)
	}
	reviewFindings = func(t *testing.T, f *sessionFixture, rid string, cs id.ID, findings ...run.Finding) (RunReviewCheckResponse, run.State) {
		t.Helper()
		store := run.NewStore(run.Dir(f.common))
		if err := store.Create(rid, run.Entry{Type: run.TypeStarted, Step: "check", At: t0,
			Started: &run.Started{Workflow: "plan", Change: cs.String(), Title: "Plan it", Definition: def, Autonomous: true}}); err != nil {
			t.Fatal(err)
		}
		var verdicts []run.Verdict
		for i := range findings {
			fid := fmt.Sprintf("f%d", i+1)
			findings[i].ID, findings[i].Kind, findings[i].Reason = fid, run.KindContradiction, "wrong"
			verdicts = append(verdicts, run.Verdict{Finding: fid, Verdict: run.Confirmed, Reason: "right"})
		}
		for _, e := range []run.Entry{
			{Type: run.TypeEntered, Step: "review", Visit: 1},
			{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}},
			{Type: run.TypeEntered, Step: "review", Visit: 2},
			{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}},
			{Type: run.TypeEntered, Step: "review", Visit: 3},
			{Type: run.TypeFindings, Step: "review", Findings: findings},
			{Type: run.TypeVerdict, Step: "critic", Verdicts: verdicts},
		} {
			e.At = t0
			if err := store.Append(rid, e); err != nil {
				t.Fatal(err)
			}
		}
		f.ctx.Getenv = func(k string) string {
			switch k {
			case envRun:
				return rid
			case "RHIZA_STEP":
				return "check"
			}
			return ""
		}
		resp, err := callOp(t, f.ctx, "run.review-check", `{}`)
		if err != nil {
			t.Fatalf("review-check: %v", err)
		}
		j, err := store.Read(rid)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(RunReviewCheckResponse), j.State()
	}
	// holder is the change set whose ops create question node q.
	holder := func(t *testing.T, f *sessionFixture, q string) changeset.ChangeSet {
		t.Helper()
		for _, c := range ctxWorkspace(t, f.ctx).changeSets() {
			for _, o := range c.Ops {
				if o.Type == op.TypeNodeCreate && o.Node.String() == q {
					return c
				}
			}
		}
		t.Fatalf("no change set creates question %s", q)
		return changeset.ChangeSet{}
	}
	blocks := func(c changeset.ChangeSet, q, node string) bool {
		return slices.ContainsFunc(c.Ops, func(o op.Op) bool {
			return o.Type == op.TypeLinkCreate && o.LinkType == op.LinkBlocks && o.From.String() == q && o.To.String() == node
		})
	}

	t.Run("open", func(t *testing.T) {
		t.Parallel()
		f := newSessionFixture(t, "qinnodeopen")
		a := newChange(t, f.ctx, `{"intent":"Plan A"}`).ChangeSet
		mustWrite(t, f.ctx, "add", `{"kind":"intent","handle":"goal","fields":["statement=Ship it"]}`)
		own := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"own-leaf","fields":["statement=Does A",`+contract+`]}`).Nodes[0]
		b := newChange(t, f.ctx, `{"intent":"Serve A"}`).ChangeSet
		mustWrite(t, f.ctx, "link", `{"from":"`+b.String()+`","type":"depends-on","to":"`+a.String()+`"}`)
		mustWrite(t, f.ctx, "link", `{"from":"`+b.String()+`","type":"serves","to":"goal"}`)
		leaf := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-b","fields":["statement=Does B",`+contract+`]}`).Nodes[0]

		resp, st := review(t, f, "R50", a, leaf, own)
		if !resp.Passed || len(resp.Questions) != 2 {
			t.Fatalf("passed %v questions %+v", resp.Passed, resp.Questions)
		}
		qb, qa := resp.Questions[0], resp.Questions[1]
		if c := holder(t, f, qb.ID); c.ID != b || !blocks(c, qb.ID, leaf.String()) {
			t.Errorf("the draft leaf's question is in %s (blocks %v), want %s blocking %s", c.ID, blocks(c, qb.ID, leaf.String()), b, leaf)
		}
		if c := holder(t, f, qa.ID); c.ID != a || !blocks(c, qa.ID, own.String()) {
			t.Errorf("the run's own node's question is in %s, want %s", c.ID, a)
		}
		changes := map[string]string{}
		for _, as := range st.Asked {
			changes[as.Question] = as.Change
		}
		if len(st.Asked) != 2 || changes[qb.ID] != b.String() || changes[qa.ID] != "" {
			t.Errorf("asked changes %v, want %s=%s and %s=\"\"", changes, qb.ID, b, qa.ID)
		}
	})

	t.Run("link", func(t *testing.T) {
		t.Parallel()
		f := newSessionFixture(t, "qinnodelink")
		a := newChange(t, f.ctx, `{"intent":"Plan A"}`).ChangeSet
		mustWrite(t, f.ctx, "add", `{"kind":"intent","handle":"goal","fields":["statement=Ship it"]}`)
		own := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"own-leaf","fields":["statement=Does A",`+contract+`]}`).Nodes[0]
		b := newChange(t, f.ctx, `{"intent":"Serve A"}`).ChangeSet
		mustWrite(t, f.ctx, "link", `{"from":"`+b.String()+`","type":"depends-on","to":"`+a.String()+`"}`)
		mustWrite(t, f.ctx, "link", `{"from":"`+b.String()+`","type":"serves","to":"goal"}`)
		leaf := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-b","fields":["statement=Does B",`+contract+`]}`).Nodes[0]
		lk := mustWrite(t, f.ctx, "link", `{"from":"leaf-b","type":"depends-on","to":"own-leaf"}`).Links[0]

		resp, st := reviewFindings(t, f, "R52", a, run.Finding{Link: lk.String(), LinkType: string(op.LinkDependsOn),
			LinkFrom: leaf.String(), LinkTo: own.String()})
		if !resp.Passed || len(resp.Questions) != 1 {
			t.Fatalf("passed %v questions %+v", resp.Passed, resp.Questions)
		}
		q := resp.Questions[0]
		if c := holder(t, f, q.ID); c.ID != b {
			t.Errorf("the draft link's question is in %s, want the draft %s that created the link", c.ID, b)
		}
		if len(st.Asked) != 1 || st.Asked[0].Change != b.String() {
			t.Errorf("asked %+v, want one naming %s", st.Asked, b)
		}
	})

	t.Run("landed", func(t *testing.T) {
		t.Parallel()
		f := newSessionFixture(t, "qinnodelanded")
		a := newChange(t, f.ctx, `{"intent":"Plan A"}`).ChangeSet
		lc := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"landed-cap","fields":["statement=Does A",`+contract+`]}`).Nodes[0]
		for _, s := range []string{"propose", "accept", "land"} {
			move(t, f.ctx, s, `{}`)
		}
		resp, st := review(t, f, "R51", a, lc)
		if !resp.Passed || len(resp.Questions) != 1 {
			t.Fatalf("passed %v questions %+v", resp.Passed, resp.Questions)
		}
		q := resp.Questions[0]
		if c := holder(t, f, q.ID); c.ID == a || c.Open.Kind != changeset.KindRecord || c.State != changeset.StateLanded {
			t.Errorf("the landed node's question is in %s (%s, %s), want a landed record", c.ID, c.Open.Kind, c.State)
		}
		if len(st.Asked) != 1 || st.Asked[0].Change != "" {
			t.Errorf("asked %+v, want one naming no change set", st.Asked)
		}
	})
}

// In a conductor step run findings --check refuses, journaling nothing, a
// round that leaves a leaf of the run's change set uncovered, naming each
// uncovered leaf, and accepts one where each leaf has a finding naming it
// or an evidence line, journaling the lines on the round. An evidence line
// for a node that is not a leaf, or with a blank line, is usage. A change
// set with no leaves, and a classic checks step, accept {"findings":[]}.
func TestRunCheckFindingsLeafCoverage(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "leafcoverage")
	a := newChange(t, f.ctx, `{"intent":"Plan A"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"parent","fields":["statement=Holds the leaves"]}`)
	leafA := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-a","parent":"parent","fields":["statement=Does A"]}`).Nodes[0].String()
	leafB := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"leaf-b","parent":"parent","fields":["statement=Does B"]}`).Nodes[0].String()
	empty := newChange(t, f.ctx, `{"intent":"No leaves"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"decision","handle":"a-ruling","fields":["statement=Rule it"]}`)

	t0 := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	conduct := "[step.conduct]\nkind = \"agent\"\nprompt = \"go\"\n" +
		"[[step.conduct.check]]\nname = \"plan\"\n[[step.conduct.check]]\nname = \"qa\"\n"
	checks, ok := workflow.Builtin("review")
	if !ok {
		t.Fatal("no built-in review workflow")
	}
	runID := ""
	start := func(rid, change, def, step string) {
		t.Helper()
		if err := store.Create(rid, run.Entry{
			Type: run.TypeStarted, Step: step, At: t0,
			Started: &run.Started{Workflow: "w", Definition: def, Change: change, Title: "Cover it"},
		}); err != nil {
			t.Fatal(err)
		}
		if err := store.Append(rid, run.Entry{Type: run.TypeEntered, Step: step, Visit: 1, At: t0}); err != nil {
			t.Fatal(err)
		}
		runID = rid
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return runID
		}
		return ""
	}
	entries := func() int {
		t.Helper()
		j, err := store.Read(runID)
		if err != nil {
			t.Fatal(err)
		}
		return len(j.Entries)
	}
	record := func(check, body string) error {
		t.Helper()
		f.ctx.Stdin = strings.NewReader(body)
		_, err := callOp(t, f.ctx, "run.findings", `{"check":"`+check+`","file":"-"}`)
		return err
	}
	refused := func(what, check, body string, want ...string) *Error {
		t.Helper()
		before := entries()
		err := record(check, body)
		e := AsError(err)
		if e == nil || e.Code != CodeUsage {
			t.Fatalf("%s: %v, want a usage error", what, err)
		}
		for _, w := range want {
			if !strings.Contains(e.Message, w) {
				t.Errorf("%s: %q missing from %q", what, w, e.Message)
			}
		}
		if n := entries(); n != before {
			t.Errorf("%s: journal grew from %d to %d entries", what, before, n)
		}
		return e
	}

	start("RC1", a.ChangeSet.String(), conduct, "conduct")
	e := refused("a clean round with no evidence", "plan", `{"findings":[]}`, "leaf-a", "leaf-b", "uncovered")
	if e.Next != runFindingsNext {
		t.Errorf("next %q, want %q", e.Next, runFindingsNext)
	}
	refused("evidence for one leaf only", "plan", `{"findings":[],"evidence":[{"leaf":"leaf-a","line":"go test: ok"}]}`, "leaf-b")
	refused("evidence on a node that is not a leaf", "plan",
		`{"findings":[],"evidence":[{"leaf":"leaf-a","line":"ok"},{"leaf":"parent","line":"read it"},{"leaf":"leaf-b","line":"ok"}]}`,
		"evidence 2", "parent", "not one of")
	refused("evidence naming nothing", "plan", `{"findings":[],"evidence":[{"leaf":"no-such-leaf","line":"ok"}]}`, "evidence 1", "no-such-leaf")
	refused("a blank evidence line", "plan", `{"findings":[],"evidence":[{"leaf":"leaf-a","line":"ok"},{"leaf":"leaf-b","line":"  "}]}`, "evidence 2", "no line")
	refused("an evidence line with no leaf", "plan", `{"findings":[],"evidence":[{"leaf":"","line":"ok"}]}`, "evidence 1", "no leaf")
	// A finding on a node outside the leaves covers nothing.
	refused("a finding naming no leaf", "plan", `{"findings":[{"severity":"nit","place":"a.go:1","reason":"r"}],"evidence":[{"leaf":"leaf-a","line":"ok"}]}`, "leaf-b")
	if st, _ := store.Read(runID); len(st.State().Checks) != 0 {
		t.Fatalf("a refused round left check state %+v", st.State().Checks)
	}

	// One evidence line per leaf: journaled, with the lines on the round.
	if err := record("plan", `{"findings":[],"evidence":[{"leaf":"leaf-a","line":"go test ./a: ok"},{"leaf":"leaf-b","line":"read b.go:3: holds"},{"leaf":"leaf-b","line":"go test ./b: ok"}]}`); err != nil {
		t.Fatalf("a round with evidence per leaf: %v", err)
	}
	j, _ := store.Read(runID)
	cs, ok := j.State().Check("conduct", "plan")
	if !ok || len(cs.Rounds) != 1 {
		t.Fatalf("plan check %+v, %v", cs, ok)
	}
	want := []run.Evidence{
		{Leaf: leafA, LeafHandle: "leaf-a", Line: "go test ./a: ok"},
		{Leaf: leafB, LeafHandle: "leaf-b", Line: "read b.go:3: holds"},
		{Leaf: leafB, LeafHandle: "leaf-b", Line: "go test ./b: ok"},
	}
	if !reflect.DeepEqual(cs.Rounds[0].Evidence, want) {
		t.Errorf("round evidence %+v, want %+v", cs.Rounds[0].Evidence, want)
	}
	if last := j.Entries[len(j.Entries)-1]; !reflect.DeepEqual(last.Evidence, want) {
		t.Errorf("findings entry evidence %+v", last.Evidence)
	}

	// A finding naming leaf A and evidence for leaf B covers both.
	if err := record("qa", `{"findings":[{"severity":"minor","leaf":"leaf-a","place":"a.go:1","reason":"r"}],"evidence":[{"leaf":"leaf-b","line":"ran it: ok"}]}`); err != nil {
		t.Fatalf("a finding on A and evidence on B: %v", err)
	}

	// A change set with no leaves needs no evidence.
	start("RC2", empty.ChangeSet.String(), conduct, "conduct")
	if err := record("plan", `{"findings":[]}`); err != nil {
		t.Fatalf("a change set with no leaves: %v", err)
	}

	// A classic checks step still accepts a clean round with no evidence.
	start("RC3", a.ChangeSet.String(), string(checks), "review")
	if err := record("compliance", `{"findings":[]}`); err != nil {
		t.Fatalf("a classic checks step: %v", err)
	}
	if j, _ := store.Read(runID); len(j.Entries[len(j.Entries)-1].Evidence) != 0 {
		t.Errorf("classic round journaled evidence %+v", j.Entries[len(j.Entries)-1].Evidence)
	}
}

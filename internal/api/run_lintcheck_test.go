package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// lint-check reads the latest code-check entry: two lints exit 1 with
// both hits in --json and the human output; an entry with none, and a run
// with no code-check entry, exit 0.
func TestLintCheck(t *testing.T) {
	f, _, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	err := store.Create("R7", run.Entry{
		Type: run.TypeStarted, Step: "code", At: t0,
		Started: &run.Started{Workflow: "review", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Lint it"},
	})
	if err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R7"
		}
		return ""
	}
	exit := func() (RunLintCheckResponse, string) {
		t.Helper()
		out, err := callOp(t, f.ctx, "run.lint-check", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		r := out.(RunLintCheckResponse)
		var sb strings.Builder
		r.Human(&sb)
		return r, sb.String()
	}

	r, _ := exit()
	if !r.Passed || r.ExitCode() != 0 || r.Lints == nil || len(r.Lints) != 0 {
		t.Errorf("no code-check entry: %+v; want exit 0 and empty lints", r)
	}

	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R7", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeCodeCheck, Step: "code", Outcome: run.CodeCheckPass})
	if r, _ = exit(); !r.Passed || r.ExitCode() != 0 {
		t.Errorf("entry with no lints: %+v; want exit 0", r)
	}

	app(run.Entry{Type: run.TypeCodeCheck, Step: "code", Outcome: run.CodeCheckPass, Lints: []run.LintHit{
		{Node: "N1", Handle: "alpha", Field: "acceptance", Kind: "test", Ref: "TestGone"},
		{Node: "N2", Handle: "beta", Kind: "unapplied-decision", Ref: "use-x"},
	}})
	r, human := exit()
	if r.Passed || r.ExitCode() != 1 || len(r.Lints) != 2 || r.Lints[0].Ref != "TestGone" || r.Lints[1].Kind != "unapplied-decision" {
		t.Errorf("two lints: %+v; want exit 1 with both hits", r)
	}
	for _, want := range []string{"test alpha acceptance: TestGone", "unapplied-decision beta: use-x"} {
		if !strings.Contains(human, want) {
			t.Errorf("human output %q lacks %q", human, want)
		}
	}

	// A later clean entry supersedes the hits.
	app(run.Entry{Type: run.TypeCodeCheck, Step: "code", Outcome: run.CodeCheckPass})
	if r, _ = exit(); !r.Passed {
		t.Errorf("latest entry clean: %+v; want exit 0", r)
	}
}

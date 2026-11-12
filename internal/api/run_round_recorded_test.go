package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// A plan review visit passes round-recorded, and its reviewer's session
// done, only once the visit journaled a findings round, empty for a clean
// review; a round of an earlier visit does not count.
func TestReviewVisitNeedsRound(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	err := store.Create("R9", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Plan it"},
	})
	if err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := store.Append("R9", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "review", Visit: 1})
	ss := session.NewStore(session.Dir(f.common))
	if _, err := ss.Update(s.ID, func(r *session.Record) { r.Run, r.Step = "R9", "review" }); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		switch k {
		case envRun:
			return "R9"
		case envSession:
			return s.ID
		}
		return ""
	}

	_, err = callOp(t, f.ctx, "run.round-recorded", `{}`)
	if e := AsError(err); e == nil || !strings.Contains(e.Message, "rhiza run findings") {
		t.Fatalf("no round: %v", err)
	}
	_, err = callOp(t, f.ctx, "session.done", `{}`)
	if e := AsError(err); e == nil || !strings.Contains(e.Message, "rhiza run findings") {
		t.Fatalf("done with no round: %v", err)
	}
	if rec, _ := ss.Read(s.ID); rec.Status.SaidDone() {
		t.Fatal("the session said done with no round")
	}

	app(run.Entry{Type: run.TypeFindings, Step: "review", Findings: []run.Finding{}})
	if _, err := callOp(t, f.ctx, "run.round-recorded", `{}`); err != nil {
		t.Fatalf("after an empty round: %v", err)
	}
	if _, err := callOp(t, f.ctx, "session.done", `{}`); err != nil {
		t.Fatalf("done after an empty round: %v", err)
	}

	// The next visit starts with that round behind it: it does not count.
	app(run.Entry{Type: run.TypeEntered, Step: "review", Visit: 2})
	if _, err := callOp(t, f.ctx, "run.round-recorded", `{}`); err == nil {
		t.Fatal("an earlier visit's round counted")
	}
}

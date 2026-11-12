package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// A review check's session done is refused until its check and cycle have a
// findings round (empty for a clean check); a critic's until every blocker
// and major finding of the check's latest round has a verdict; a fixer is
// never refused for this.
func TestSessionDoneReviewRecords(t *testing.T) {
	r, _, ctx, _, ada := evidenceCommitRunner(t)
	common, err := gitCommonDir(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ss := session.NewStore(session.Dir(common))
	rs := run.NewStore(run.Dir(common))
	t0 := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	if err := rs.Create("R2", run.Entry{
		Type: run.TypeStarted, Step: "checks", At: t0,
		Started: &run.Started{Workflow: "review", Change: r.change, Title: "Review", Worktree: r.worktree},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := rs.Append("R2", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "checks", Visit: 1})
	mk := func(sid, role string) {
		t.Helper()
		rec := session.Record{ID: sid, Actor: ada.String(), Change: r.change, Run: "R2", Step: "checks", Check: "bugs", Cycle: 1, Role: role}
		if err := ss.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	mk("S-c", "check")
	mk("S-k", "critic")
	mk("S-f", "fixer")
	done := func(sid string) error {
		_, err := callOp(t, ctx, "session.done", `{"id":"`+sid+`"}`)
		return err
	}
	said := func(sid string) bool { rec, _ := ss.Read(sid); return rec.Status.SaidDone() }

	err = done("S-c")
	if e := AsError(err); e == nil || !strings.Contains(e.Message, "rhiza run findings") || said("S-c") {
		t.Fatalf("check with no round: %v", err)
	}
	app(run.Entry{Type: run.TypeFindings, Step: "checks", Check: "bugs", Cycle: 1, Findings: []run.Finding{
		{ID: "f1", Check: "bugs", Severity: run.SeverityBlocker, Place: "a.go:1", Reason: "x"},
		{ID: "f2", Check: "bugs", Severity: run.SeverityMajor, Place: "a.go:2", Reason: "y"},
		{ID: "f3", Check: "bugs", Severity: run.SeverityNit, Place: "a.go:3", Reason: "z"},
	}})
	if err := done("S-c"); err != nil || !said("S-c") {
		t.Fatalf("check after round: %v", err)
	}

	err = done("S-k")
	if e := AsError(err); e == nil || !strings.Contains(e.Message, "rhiza run verdict") || said("S-k") {
		t.Fatalf("critic with none judged: %v", err)
	}
	app(run.Entry{Type: run.TypeVerdict, Step: "checks", Check: "bugs", Verdicts: []run.Verdict{{Finding: "f1", Verdict: run.Confirmed, Reason: "r"}}})
	if err := done("S-k"); err == nil || said("S-k") {
		t.Fatalf("critic with f2 unjudged: %v", err)
	}
	app(run.Entry{Type: run.TypeVerdict, Step: "checks", Check: "bugs", Verdicts: []run.Verdict{{Finding: "f2", Verdict: run.Dropped, Reason: "r"}}})
	if err := done("S-k"); err != nil || !said("S-k") {
		t.Fatalf("critic after verdicts: %v", err)
	}

	if err := done("S-f"); err != nil || !said("S-f") {
		t.Fatalf("fixer: %v", err)
	}
}

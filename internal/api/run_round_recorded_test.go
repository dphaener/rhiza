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

// A plan review finding may name a link, by ID: a model link or a depends-on
// between change sets. It records the round naming the link and its type,
// status shows it, a verdict confirms it, and the round-3 review-check adds
// a question naming it that blocks nothing; an ID that is neither a node
// nor a link is usage and writes nothing.
func TestReviewFindingsNameLinks(t *testing.T) {
	f := newSessionFixture(t, "links")
	a := newChange(t, f.ctx, `{"intent":"Plan A"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"x"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"y"}`)
	modelLink := mustWrite(t, f.ctx, "link", `{"from":"x","type":"depends-on","to":"y"}`).Links[0]
	b := newChange(t, f.ctx, `{"intent":"Plan B"}`)
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"z"}`)
	if _, err := callOp(t, f.ctx, "change.use", `{"ref":"`+a.ChangeSet.String()+`"}`); err != nil {
		t.Fatal(err)
	}
	csLink := mustWrite(t, f.ctx, "link", `{"from":"`+a.ChangeSet.String()+`","type":"depends-on","to":"`+b.ChangeSet.String()+`"}`).Links[0]

	t0 := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	store := run.NewStore(run.Dir(f.common))
	if err := store.Create("R10", run.Entry{
		Type: run.TypeStarted, Step: "review", At: t0,
		Started: &run.Started{Workflow: "plan", Change: a.ChangeSet.String(), Title: "Plan it",
			Definition: "name = \"plan\"\n[step.review]\nkind = \"command\"\nrun = \"true\"\n"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("R10", run.Entry{Type: run.TypeEntered, Step: "review", At: t0, Visit: 1}); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R10"
		}
		return ""
	}
	findings := func(node string) (any, error) {
		f.ctx.Stdin = strings.NewReader(`{"findings":[{"node":"` + node + `","kind":"unreal-depends-on","reason":"Not real"}]}`)
		return callOp(t, f.ctx, "run.findings", `{"file":"-"}`)
	}

	// Neither a node nor a link: usage, nothing journalled.
	_, err := findings("01M76F8S1FKFPKAAYVM3GD01CN")
	wantExit(t, err, 2)
	if st, _ := store.Read("R10"); len(st.State().Rounds) != 0 {
		t.Fatal("a refused finding journalled a round")
	}

	for i, c := range []struct{ link, typ string }{{csLink.String(), "depends-on"}, {modelLink.String(), "depends-on"}} {
		resp, err := findings(c.link)
		if err != nil {
			t.Fatalf("link %s: %v", c.link, err)
		}
		fd := resp.(RunFindingsResponse).Findings[0]
		if fd.Link != c.link || fd.LinkType != c.typ || fd.Node != "" {
			t.Errorf("finding %+v", fd)
		}
		if st, _ := store.Read("R10"); len(st.State().Rounds) != i+1 {
			t.Fatalf("rounds %d, want %d", len(st.State().Rounds), i+1)
		}
		f.ctx.Stdin = strings.NewReader(`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"Right"}]}`)
		if _, err := callOp(t, f.ctx, "run.verdict", `{"file":"-"}`); err != nil {
			t.Fatalf("verdict: %v", err)
		}
	}
	st, _ := store.Read("R10")
	rd, _ := st.State().LatestRound()
	if f1, _ := rd.Finding("f1"); f1.Link != modelLink.String() || f1.LinkType != "depends-on" {
		t.Fatalf("journalled finding %+v", f1)
	}

	// Round 3: the check leaves a question naming the link, blocking nothing.
	if err := store.Append("R10", run.Entry{Type: run.TypeEntered, Step: "review", At: t0, Visit: 3}); err != nil {
		t.Fatal(err)
	}
	if _, err := findings(csLink.String()); err != nil {
		t.Fatal(err)
	}
	f.ctx.Stdin = strings.NewReader(`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"Right"}]}`)
	if _, err := callOp(t, f.ctx, "run.verdict", `{"file":"-"}`); err != nil {
		t.Fatal(err)
	}
	resp, err := callOp(t, f.ctx, "run.review-check", `{}`)
	if err != nil {
		t.Fatalf("round 3 review-check: %v", err)
	}
	r := resp.(RunReviewCheckResponse)
	if !r.Passed || len(r.Questions) != 1 || r.Questions[0].Blocks != "" || r.Questions[0].ID == "" {
		t.Fatalf("response %+v", r)
	}
	if got := r.Confirmed[0].Link; got != csLink.String() {
		t.Errorf("confirmed link %q", got)
	}
}

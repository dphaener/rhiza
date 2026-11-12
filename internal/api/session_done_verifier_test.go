package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// A leaf verifier's session done is refused, naming rhiza run findings,
// until its attempt has journaled findings or Evidence its own actor wrote;
// Evidence by another actor alone does not count, and a builder is never
// refused.
func TestSessionDoneVerifierNeedsRecord(t *testing.T) {
	r, _, ctx, leaf, ada := evidenceCommitRunner(t)
	common, err := gitCommonDir(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ss := session.NewStore(session.Dir(common))
	bob := recordAda(t, ctx, "bob@example.com")
	t0 := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	rs := run.NewStore(run.Dir(common))
	if err := rs.Create("R1", run.Entry{
		Type: run.TypeStarted, Step: "build", At: t0,
		Started: &run.Started{Workflow: "build", Change: r.change, Title: "Build", Worktree: r.worktree},
	}); err != nil {
		t.Fatal(err)
	}
	app := func(e run.Entry) {
		t.Helper()
		e.At = t0
		if err := rs.Append("R1", e); err != nil {
			t.Fatal(err)
		}
	}
	app(run.Entry{Type: run.TypeEntered, Step: "build", Visit: 1})
	mk := func(sid, role string, actor id.ID) {
		t.Helper()
		rec := session.Record{ID: sid, Actor: actor.String(), Change: r.change, Run: "R1", Step: "build", Leaf: "leaf-a", Attempt: 1, Role: role}
		if err := ss.Write(rec); err != nil {
			t.Fatal(err)
		}
		app(run.Entry{Type: run.TypeLeafSession, Step: "build", Leaf: leaf.ID, Handle: "leaf-a", Role: role, Attempt: 1, Session: sid})
	}
	mk("S-b", "builder", ada)
	mk("S-v", "verifier", ada)
	done := func(sid string) error {
		_, err := callOp(t, ctx, "session.done", `{"id":"`+sid+`"}`)
		return err
	}
	said := func(sid string) bool { rec, _ := ss.Read(sid); return rec.Status.SaidDone() }

	if err := done("S-b"); err != nil {
		t.Fatalf("builder: %v", err)
	}
	err = done("S-v")
	if e := AsError(err); e == nil || !strings.Contains(e.Message, "rhiza run findings") {
		t.Fatalf("nothing recorded: %v", err)
	}
	if said("S-v") {
		t.Fatal("done mark written with nothing recorded")
	}

	// Another actor's Evidence alone is still refused.
	v := *ctx
	v.Actor = bob
	mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-bob","fields":["statement=ok","result=PASS"]}`)
	mustWrite(t, &v, "link", `{"from":"leaf-a-bob","type":"verifies","to":"leaf-a"}`)
	if err := done("S-v"); err == nil || said("S-v") {
		t.Fatalf("another actor's Evidence counted: %v", err)
	}

	// Findings on the attempt let it say done.
	app(run.Entry{Type: run.TypeFindings, Step: "build", Leaf: leaf.ID, Handle: "leaf-a", Attempt: 1, Findings: []run.Finding{{ID: "f1", Reason: "bad"}}})
	if err := done("S-v"); err != nil || !said("S-v") {
		t.Fatalf("after findings: %v", err)
	}

	// Its own Evidence does too.
	v.Actor = ada
	mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-ada","fields":["statement=ok","result=PASS"]}`)
	mustWrite(t, &v, "link", `{"from":"leaf-a-ada","type":"verifies","to":"leaf-a"}`)
	// Attempt 1's findings count only for that attempt: use attempt 2.
	if err := ss.Write(session.Record{ID: "S-v3", Actor: ada.String(), Change: r.change, Run: "R1", Step: "build", Leaf: "leaf-a", Attempt: 2, Role: "verifier"}); err != nil {
		t.Fatal(err)
	}
	if err := done("S-v3"); err != nil {
		t.Fatalf("own Evidence: %v", err)
	}
}

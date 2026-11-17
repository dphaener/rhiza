package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// TestStepRunnerChild: a restarted driver asks the runner where a spawned
// child stands: running, ended with its outcome, or gone with its record.
func TestStepRunnerChild(t *testing.T) {
	t.Parallel()
	f := newSpawnFixture(t)
	sessions := session.NewStore(session.Dir(f.f.common))
	zero, three := 0, 3
	said := &session.Report{State: session.StateDone}
	for id, st := range map[string]session.Status{
		"live": {},
		"ok":   {Exit: &session.Report{State: session.StateExited, Code: &zero}},
		"red":  {Exit: &session.Report{State: session.StateExited, Code: &three}},
		// A done mark ends the child with pass, whatever exit follows it.
		"done":      {Done: said},
		"done-red":  {Done: said, Exit: &session.Report{State: session.StateExited, Code: &three}},
		"done-lost": {Done: said, Hook: &session.Report{State: session.StateDone, At: f.f.ctx.Now().Add(-time.Hour)}},
	} {
		if err := sessions.Write(session.Record{ID: id, Run: "R1", Status: st}); err != nil {
			t.Fatal(err)
		}
	}
	r := &stepRunner{ctx: f.f.ctx, sessions: sessions, run: "R1"}
	for _, c := range []struct {
		id      string
		ended   bool
		gone    bool
		outcome string
		detail  string
	}{
		{"live", false, false, "", ""},
		{"ok", true, false, workflow.OutcomePass, "exit 0"},
		{"red", true, false, workflow.OutcomeFail, "exit 3"},
		{"done", true, false, workflow.OutcomePass, "said done"},
		{"done-red", true, false, workflow.OutcomePass, "said done"},
		{"done-lost", true, false, workflow.OutcomePass, "said done"},
		{"nobody", false, true, "", "record is gone"},
	} {
		got, err := r.Child(c.id)
		if err != nil {
			t.Fatalf("%s: %v", c.id, err)
		}
		if got.Ended != c.ended || got.Gone != c.gone || got.Outcome != c.outcome || !strings.Contains(got.Detail, c.detail) {
			t.Errorf("%s: %+v", c.id, got)
		}
	}
}

// markLostIfEnded marks a done child whose hook ended lost, and the done
// mark still wins: TestStepRunnerChild's done-lost reads pass.
func TestStepRunnerChildDoneLost(t *testing.T) {
	t.Parallel()
	f := newSpawnFixture(t)
	sessions := session.NewStore(session.Dir(f.f.common))
	if err := sessions.Write(session.Record{ID: "c", Run: "R1", Status: session.Status{
		Done: &session.Report{State: session.StateDone},
		Hook: &session.Report{State: session.StateDone, At: f.f.ctx.Now().Add(-time.Hour)},
	}}); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: f.f.ctx, sessions: sessions, run: "R1"}
	got, err := r.Child("c")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := sessions.Read("c")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Exit == nil || rec.Status.Exit.State != session.StateLost {
		t.Fatalf("not marked lost: %+v", rec.Status)
	}
	if !got.Ended || got.Outcome != workflow.OutcomePass || got.Detail != "said done" {
		t.Errorf("child %+v", got)
	}
}

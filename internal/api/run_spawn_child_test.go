package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// TestStepRunnerChild: a restarted driver asks the runner where a spawned
// child stands: running, ended with its outcome, or gone with its record.
func TestStepRunnerChild(t *testing.T) {
	f := newSpawnFixture(t)
	sessions := session.NewStore(session.Dir(f.f.common))
	zero, three := 0, 3
	for id, exit := range map[string]*session.Report{
		"live": nil,
		"ok":   {State: session.StateExited, Code: &zero},
		"red":  {State: session.StateExited, Code: &three},
	} {
		if err := sessions.Write(session.Record{ID: id, Run: "R1", Status: session.Status{Exit: exit}}); err != nil {
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

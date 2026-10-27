package api

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// sendMessagePayload is a SendMessage PostToolUse payload for sid, with
// made-up content in the real shape.
func sendMessagePayload(sid string) string {
	return `{"session_id":"` + sid + `","transcript_path":"/tmp/x.jsonl","cwd":"/tmp/w",` +
		`"hook_event_name":"PostToolUse","tool_name":"SendMessage",` +
		`"tool_input":{"to":"reviewer","summary":"Ready","message":"All done"},` +
		`"tool_response":{"success":true},"tool_use_id":"toolu_01X"}`
}

// A session in a run that has exited still gets the messages it sent
// journaled, while the hook's status is not recorded: a late hook does not
// revive the session.
func TestSessionHookJournalsAfterExit(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual","run":"run-1","step":"talk"}`)
	rs := run.NewStore(run.Dir(f.common))
	if err := rs.Create("run-1", run.Entry{Type: run.TypeStarted, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	st := session.NewStore(session.Dir(f.common))
	code := 0
	if _, err := st.Update(s.ID, func(r *session.Record) {
		r.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now().UTC()}
	}); err != nil {
		t.Fatal(err)
	}

	r, errOut := hook(t, f.ctx, "PostToolUse", sendMessagePayload(s.ID))
	if errOut != "" {
		t.Errorf("stderr %q", errOut)
	}
	if r.Recorded || !r.Journaled || r.Session != s.ID {
		t.Errorf("response %+v: want journaled, not recorded", r)
	}
	rec, err := st.Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Hook != nil {
		t.Errorf("hook report after the exit %+v", rec.Status.Hook)
	}

	j, err := rs.Read("run-1")
	if err != nil {
		t.Fatal(err)
	}
	var msgs []run.Entry
	for _, e := range j.Entries {
		if e.Type == run.TypeMessage {
			msgs = append(msgs, e)
		}
	}
	if len(msgs) != 1 {
		t.Fatalf("message entries %+v", msgs)
	}
	m := msgs[0]
	if m.Step != "talk" || m.Session != s.ID || m.Actor != rec.Actor || m.Actor == "" ||
		m.To != "reviewer" || m.Text != "All done" || m.Summary != "Ready" || m.Structured || m.At.IsZero() {
		t.Errorf("message entry %+v", m)
	}
}

// A session in no run journals nothing and reports nothing.
func TestSessionHookNoRunJournalsNothing(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	r, errOut := hook(t, f.ctx, "PostToolUse", sendMessagePayload(s.ID))
	if errOut != "" {
		t.Errorf("stderr %q", errOut)
	}
	if r.Journaled || !r.Recorded || r.State != session.StateWorking {
		t.Errorf("response %+v", r)
	}
}

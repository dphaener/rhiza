package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
)

// startClock makes a run's session out of the message fixture's one and
// puts its context on a clock the test moves; it returns the advance.
func startClock(t *testing.T, f *messageFixture) func(time.Duration) {
	t.Helper()
	if _, err := f.st.Update(f.sid, func(x *session.Record) { x.Run = "run-1" }); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 13, 9, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	return func(d time.Duration) { now = now.Add(d) }
}

// feed runs the agent's hook for event on the fixture's session, as the
// claude-code hook does, and wants it recorded.
func feed(t *testing.T, f *messageFixture, event string) {
	t.Helper()
	r, errOut := hook(t, f.ctx, event, `{"session_id":"`+f.sid+`","hook_event_name":"`+event+`"}`)
	if !r.Recorded || errOut != "" {
		t.Fatalf("%s: hook %+v, stderr %q", event, r, errOut)
	}
}

// SessionStart then idle past the grace with no prompt submitted fails the
// session with that reason, as an environment failure, and again on every
// later read.
func TestStartNotSubmittedFailsSession(t *testing.T) {
	f, r := lostFixture(t)
	advance := startClock(t, f)
	feed(t, f, "SessionStart")
	rec, err := f.st.Read(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Started == nil || rec.Status.Prompted != nil {
		t.Fatalf("after SessionStart: started %+v, prompted %+v", rec.Status.Started, rec.Status.Prompted)
	}
	advance(session.StartGrace + time.Second)

	for i := 0; i < 2; i++ {
		s, err := r.AgentStatus(f.sid)
		if err != nil {
			t.Fatal(err)
		}
		if !s.Done || !s.Environment || s.Outcome != "fail" || s.Detail != session.NeverSubmittedDetail {
			t.Fatalf("status %+v, want an environment failure: %s", s, session.NeverSubmittedDetail)
		}
	}
	if rec, err = f.st.Read(f.sid); err != nil {
		t.Fatal(err)
	}
	if x := rec.Status.Exit; x == nil || x.State != session.StateError || !strings.Contains(x.Detail, "never submitted") {
		t.Errorf("exit %+v", x)
	}
}

// SessionStart, UserPromptSubmit then idle (Stop) past the grace does not
// fail.
func TestStartNotSubmittedPromptedIsFine(t *testing.T) {
	f, r := lostFixture(t)
	advance := startClock(t, f)
	feed(t, f, "SessionStart")
	advance(time.Second)
	feed(t, f, "UserPromptSubmit")
	rec, err := f.st.Read(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Prompted == nil {
		t.Fatal("UserPromptSubmit did not mark the session prompted")
	}
	advance(time.Second)
	feed(t, f, "Stop")
	advance(session.StartGrace + time.Hour)

	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if s.Done || s.Environment || !s.Idle {
		t.Errorf("status %+v, want idle and live", s)
	}
}

// SessionStart then idle within the grace does not fail.
func TestStartNotSubmittedWithinGrace(t *testing.T) {
	f, r := lostFixture(t)
	advance := startClock(t, f)
	feed(t, f, "SessionStart")
	advance(session.StartGrace - time.Second)

	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if s.Done || s.Environment {
		t.Errorf("status %+v, want live", s)
	}
}

package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
)

// startedAt records a SessionStart hook report age before the run's clock,
// and, when prompted, a working one after it.
func startedAt(t *testing.T, f *messageFixture, r *stepRunner, age time.Duration, prompted bool) {
	t.Helper()
	at := r.ctx.Now().UTC().Add(-age)
	if _, err := f.st.Update(f.sid, func(x *session.Record) {
		x.Run = "run-1"
		x.Status.Hook = &session.Report{State: session.StateIdle, At: at}
		x.Status.Started = &session.Report{State: session.StateIdle, At: at}
		if prompted {
			x.Status.Hook = &session.Report{State: session.StateWorking, At: at}
			x.Status.Prompted = x.Status.Hook
		}
	}); err != nil {
		t.Fatal(err)
	}
}

// SessionStart then idle past the grace with no prompt submitted fails the
// session with that reason, as an environment failure, and again on every
// later read.
func TestStartNotSubmittedFailsSession(t *testing.T) {
	f, r := lostFixture(t)
	startedAt(t, f, r, session.StartGrace+time.Minute, false)

	for i := 0; i < 2; i++ {
		s, err := r.AgentStatus(f.sid)
		if err != nil {
			t.Fatal(err)
		}
		if !s.Done || !s.Environment || s.Outcome != "fail" || s.Detail != session.NeverSubmittedDetail {
			t.Fatalf("status %+v, want an environment failure: %s", s, session.NeverSubmittedDetail)
		}
	}
	rec, err := f.st.Read(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if x := rec.Status.Exit; x == nil || x.State != session.StateError || !strings.Contains(x.Detail, "never submitted") {
		t.Errorf("exit %+v", x)
	}
}

// SessionStart, UserPromptSubmit then idle does not fail.
func TestStartNotSubmittedPromptedIsFine(t *testing.T) {
	f, r := lostFixture(t)
	startedAt(t, f, r, session.StartGrace+time.Hour, true)
	if _, err := f.st.Update(f.sid, func(x *session.Record) {
		x.Status.Hook = &session.Report{State: session.StateIdle, At: r.ctx.Now().UTC()}
	}); err != nil {
		t.Fatal(err)
	}
	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if s.Done || s.Environment || !s.Idle {
		t.Errorf("status %+v, want idle and live", s)
	}
}

// Idle within the grace does not fail.
func TestStartNotSubmittedWithinGrace(t *testing.T) {
	f, r := lostFixture(t)
	startedAt(t, f, r, session.StartGrace-time.Second, false)
	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if s.Done || s.Environment {
		t.Errorf("status %+v, want live", s)
	}
}

// A hook records the start and the first non-idle report.
func TestStartNotSubmittedHookMarks(t *testing.T) {
	f, _ := lostFixture(t)
	if _, err := f.st.Update(f.sid, func(x *session.Record) { x.Run = "run-1" }); err != nil {
		t.Fatal(err)
	}
	f.setHook(t, session.StateIdle, "")
	rec, _ := f.st.Read(f.sid)
	if rec.Status.Prompted != nil {
		t.Errorf("an idle report marked the session prompted")
	}
}

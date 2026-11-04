package api

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
)

func lostFixture(t *testing.T) (*messageFixture, *stepRunner) {
	t.Helper()
	f := newMessageFixture(t)
	return f, &stepRunner{ctx: f.ctx, sessions: f.st}
}

func hookEnd(t *testing.T, f *messageFixture) {
	t.Helper()
	f.setHook(t, session.StateDone, "")
}

// hookEndAt records the SessionEnd mark at age before the run's clock.
func hookEndAt(t *testing.T, f *messageFixture, r *stepRunner, age time.Duration) {
	t.Helper()
	if _, err := f.st.Update(f.sid, func(x *session.Record) {
		x.Status.Hook = &session.Report{State: session.StateDone, At: r.ctx.Now().UTC().Add(-age)}
	}); err != nil {
		t.Fatal(err)
	}
}

func wantLost(t *testing.T, f *messageFixture, detail string) {
	t.Helper()
	rec, err := f.st.Read(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	x := rec.Status.Exit
	if x == nil || x.State != session.StateLost || x.Code != nil || !strings.Contains(x.Detail, detail) {
		t.Fatalf("exit %+v, want lost naming %q", x, detail)
	}
}

// A rejoined session with the SessionEnd mark and no exit is lost and
// judged failed; Rex is asked nothing for it.
func TestRejoinHookEndedIsLost(t *testing.T) {
	f, r := lostFixture(t)
	hookEnd(t, f)
	if err := r.Rejoin(f.sid); err != nil {
		t.Fatal(err)
	}
	wantLost(t, f, lostHookEnded)
	f.noRex(t, "a record that already says it ended")
	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Done || !s.Lost || s.Outcome != "fail" || !strings.Contains(s.Detail, "lost") {
		t.Errorf("status %+v", s)
	}
}

// A rejoined session whose block Rex does not know is lost.
func TestRejoinBlockNotFoundIsLost(t *testing.T) {
	f, r := lostFixture(t)
	f.setHook(t, session.StateWorking, "")
	f.rex.err = errors.New("block not found: block:b1")
	if err := r.Rejoin(f.sid); err != nil {
		t.Fatal(err)
	}
	wantLost(t, f, lostNoBlock)
	if s, _ := r.AgentStatus(f.sid); !s.Done || !s.Lost || s.Outcome != "fail" {
		t.Errorf("status %+v", s)
	}
}

// A block Rex still has leaves the session polled as before; later polls
// make no Rex call.
func TestRejoinBlockExists(t *testing.T) {
	f, r := lostFixture(t)
	f.setHook(t, session.StateWorking, "")
	if err := r.Rejoin(f.sid); err != nil {
		t.Fatal(err)
	}
	if len(f.rex.calls) != 1 || f.rex.calls[0][1] != "block" || f.rex.calls[0][2] != "inspect" {
		t.Fatalf("rex calls %q", f.rex.calls)
	}
	f.rex.calls = nil
	if s, err := r.AgentStatus(f.sid); err != nil || s.Done {
		t.Fatalf("status %+v %v", s, err)
	}
	f.noRex(t, "a poll after the block was found")
}

// Rex failing otherwise judges nothing and is retried on later polls until
// it answers; then not found makes the session lost.
func TestRejoinRexUnreachableRetried(t *testing.T) {
	f, r := lostFixture(t)
	f.setHook(t, session.StateWorking, "")
	f.rex.err = errors.New("connection refused")
	if err := r.Rejoin(f.sid); err == nil {
		t.Fatal("no error for an unreachable Rex")
	}
	if s, err := r.AgentStatus(f.sid); err != nil || s.Done {
		t.Fatalf("status %+v %v", s, err)
	}
	if rec, _ := f.st.Read(f.sid); rec.Status.Exit != nil {
		t.Fatalf("exit written: %+v", rec.Status.Exit)
	}
	f.rex.err = errors.New("block not found")
	if s, _ := r.AgentStatus(f.sid); !s.Done || !s.Lost {
		t.Fatalf("status %+v", s)
	}
	wantLost(t, f, lostNoBlock)
}

// Mid-run, a record whose SessionEnd mark has stood past the grace with
// no exit is lost on the next poll, with no Rex call.
func TestPollHookEndedIsLostWithoutRex(t *testing.T) {
	f, r := lostFixture(t)
	f.setHook(t, session.StateWorking, "")
	if s, _ := r.AgentStatus(f.sid); s.Done {
		t.Fatalf("done too early: %+v", s)
	}
	hookEndAt(t, f, r, hookEndGrace)
	s, err := r.AgentStatus(f.sid)
	if err != nil || !s.Done || !s.Lost || s.Outcome != "fail" {
		t.Fatalf("status %+v %v", s, err)
	}
	f.noRex(t, "per-poll check")
	wantLost(t, f, lostHookEnded)
}

// A fresh SessionEnd mark is not judged on a poll: a session ending
// normally has it before its wrapper records the exit. The exit that
// follows is what the poll then reads.
func TestPollFreshHookEndNotJudged(t *testing.T) {
	f, r := lostFixture(t)
	hookEndAt(t, f, r, hookEndGrace-time.Second)
	s, err := r.AgentStatus(f.sid)
	if err != nil || s.Done {
		t.Fatalf("status %+v %v", s, err)
	}
	if rec, _ := f.st.Read(f.sid); rec.Status.Exit != nil {
		t.Fatalf("exit written: %+v", rec.Status.Exit)
	}
	code := 0
	if _, err := f.st.Update(f.sid, func(x *session.Record) {
		x.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now()}
	}); err != nil {
		t.Fatal(err)
	}
	if s, _ := r.AgentStatus(f.sid); !s.Done || s.Lost || s.Outcome != "pass" {
		t.Errorf("status %+v", s)
	}
	f.noRex(t, "per-poll check")
}

// An exit already on the record is final: the lost check never replaces it.
func TestLostNeverReplacesExit(t *testing.T) {
	f, r := lostFixture(t)
	code := 0
	f.st.Update(f.sid, func(x *session.Record) {
		x.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now()}
		x.Status.Hook = &session.Report{State: session.StateDone}
	})
	if err := r.Rejoin(f.sid); err != nil {
		t.Fatal(err)
	}
	if s, _ := r.AgentStatus(f.sid); !s.Done || s.Lost || s.Outcome != "pass" {
		t.Errorf("status %+v", s)
	}
}

// A wrapper whose session a run judged lost while it ran keeps that lost
// exit: it does not overwrite it with its own, so the record and the run
// agree.
func TestSessionRunKeepsLostExit(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"manual"}`)
	fs.agent.code = 0
	st := session.NewStore(session.Dir(f.common))
	f.ctx.Start = func(p Program) (Process, error) {
		if p.Terminal {
			if _, err := st.Update(s.ID, func(x *session.Record) {
				x.Status.Exit = &session.Report{State: session.StateLost, Detail: lostHookEnded, At: time.Now()}
			}); err != nil {
				t.Fatal(err)
			}
		}
		return fs.start(p)
	}
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	rec, err := st.Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if x := rec.Status.Exit; x == nil || x.State != session.StateLost || x.Code != nil {
		t.Errorf("exit %+v, want the lost exit kept", x)
	}
}

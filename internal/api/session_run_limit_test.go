package api

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// limitAdapter is a fake agent whose transcript is $HOME/limit/<session>.txt
// and which is usage-limited while that file says "limited". Every command
// it builds goes to limitLaunches.
type limitAdapter struct{ interviewAdapter }

const limitAgent = "limit-fake"

var (
	limitLaunches []adapter.Launch
	limitRegister sync.Once
)

func (limitAdapter) Name() string { return limitAgent }
func (limitAdapter) Command(l adapter.Launch) ([]string, error) {
	limitLaunches = append(limitLaunches, l)
	return []string{limitAgent}, nil
}
func (limitAdapter) TranscriptGlobs(sid string, getenv func(string) string) []string {
	return []string{filepath.Join(getenv("HOME"), "limit", sid+".txt")}
}
func (limitAdapter) UsageLimit(t []byte) (adapter.UsageLimit, bool) {
	return adapter.UsageLimit{}, strings.Contains(string(t), "limited")
}

// HookStatus reports idle for a Stop of the session the payload's
// session_id names.
func (limitAdapter) HookStatus(event string, raw []byte) (adapter.HookReport, bool, error) {
	var p struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return adapter.HookReport{}, false, err
	}
	if event != "Stop" {
		return adapter.HookReport{}, false, nil
	}
	return adapter.HookReport{SessionID: p.SessionID, State: session.StateIdle}, true, nil
}

// removeClock removes a file when the run sleeps.
type removeClock struct{ path string }

func (c removeClock) Now() time.Time      { return time.Unix(0, 0) }
func (c removeClock) Sleep(time.Duration) { os.Remove(c.path) }

func TestSessionRunFreshSessionWhenTranscriptGone(t *testing.T) {
	f, s, path := limitFixture(t)
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	st := session.NewStore(session.Dir(f.common))
	if _, err := st.Update(s.ID, func(r *session.Record) { r.Run, r.Step = "run-1", "build" }); err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	if err := rs.Create("run-1", run.Entry{Type: run.TypeStarted, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	useClock(t, removeClock{path})

	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionRunResponse)
	launches := limitLaunches
	if len(launches) != 2 || launches[1].Resume || launches[1].Prompt != "do the step" || launches[1].SessionID != s.ID {
		t.Errorf("launches %+v: want a second, fresh launch of the same session with its first prompt", launches)
	}
	if !strings.Contains(strings.Join(r.Warnings, "\n"), "started a fresh session") {
		t.Errorf("warnings %v do not record the fresh session", r.Warnings)
	}

	// The restart outlives the session's terminal: its record and its
	// run's journal keep it.
	rec, err := st.Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Fresh) != 1 {
		t.Errorf("record's fresh %v: want one restart", rec.Fresh)
	}
	j, err := rs.Read("run-1")
	if err != nil {
		t.Fatal(err)
	}
	var fresh []run.Entry
	for _, e := range j.Entries {
		if e.Type == run.TypeFreshSession {
			fresh = append(fresh, e)
		}
	}
	if len(fresh) != 1 || fresh[0].Session != s.ID || fresh[0].Step != "build" || !strings.Contains(fresh[0].Detail, "started a fresh session") {
		t.Errorf("journal's fresh-session entries %+v: want one for session %s", fresh, s.ID)
	}
}

// blockingProcess is an agent that stays open until killed, as an
// interactive agent does at its usage limit.
type blockingProcess struct {
	killed chan struct{}
	once   sync.Once
}

func (p *blockingProcess) Output() io.Reader { return nil }
func (p *blockingProcess) Wait() (int, error) {
	<-p.killed
	return 137, nil
}
func (p *blockingProcess) Kill() { p.once.Do(func() { close(p.killed) }) }

// writeClock writes text to a file when the run sleeps.
type writeClock struct{ path, text string }

func (c writeClock) Now() time.Time { return time.Unix(0, 0) }
func (c writeClock) Sleep(time.Duration) {
	os.WriteFile(c.path, []byte(c.text), 0o644)
}

// limitFixture is a session of the limit-fake agent, with no launches yet,
// and the path of its transcript.
func limitFixture(t *testing.T) (*sessionFixture, SessionInfo, string) {
	t.Helper()
	f, s, _ := interviewFixture(t)
	limitRegister.Do(func() { adapter.Register(limitAdapter{}) })
	limitLaunches = nil
	if _, err := session.NewStore(session.Dir(f.common)).Update(s.ID, func(r *session.Record) {
		r.Agent, r.Prompt = limitAgent, "do the step"
	}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.ctx.Home, "limit", s.ID+".txt")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	return f, s, path
}

// An agent that stays open at its usage limit is stopped, and the session
// resumed once the limit resets.
func TestSessionRunStopsAgentOpenAtLimit(t *testing.T) {
	f, s, path := limitFixture(t)
	old := limitPoll
	limitPoll = time.Millisecond
	t.Cleanup(func() { limitPoll = old })
	useClock(t, writeClock{path, "done"})
	open := &blockingProcess{killed: make(chan struct{})}
	starts := 0
	f.ctx.Start = func(p Program) (Process, error) {
		starts++
		if starts > 1 {
			return &fakeProcess{}, nil
		}
		if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
			t.Fatal(err)
		}
		return open, nil
	}

	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	if len(limitLaunches) != 2 || !limitLaunches[1].Resume {
		t.Errorf("launches %+v: want the open agent stopped and the session resumed", limitLaunches)
	}
}

// A session still at its usage limit after the last pause fails: the step
// did not finish.
func TestSessionRunFailsWhenPausesRunOut(t *testing.T) {
	f, s, path := limitFixture(t)
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	useClock(t, writeClock{path, "limited"})
	f.ctx.Start = func(Program) (Process, error) { return &fakeProcess{}, nil }

	_, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err == nil || !strings.Contains(err.Error(), "still stopped at its usage limit") {
		t.Fatalf("error %v: want the session failed at its usage limit", err)
	}
	if len(limitLaunches) != maxUsageLimitPauses+1 {
		t.Errorf("%d launches, want %d", len(limitLaunches), maxUsageLimitPauses+1)
	}
	rec, rerr := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if rec.Status.Exit == nil || rec.Status.Exit.State != session.StateError {
		t.Errorf("exit %+v, want an error", rec.Status.Exit)
	}
}

// An agent that stops at its usage limit is recorded as working, paused,
// not idle: a run must not take it for ready at its prompt. Once it goes
// on, its Stop is idle again.
func TestSessionHookStopAtLimitIsPaused(t *testing.T) {
	f, s, path := limitFixture(t)
	st := session.NewStore(session.Dir(f.common))
	stop := func() session.Report {
		t.Helper()
		f.ctx.Stdin = strings.NewReader(`{"session_id":"` + s.ID + `"}`)
		if _, err := callOp(t, f.ctx, "session.hook", `{"event":"Stop"}`); err != nil {
			t.Fatal(err)
		}
		rec, err := st.Read(s.ID)
		if err != nil {
			t.Fatal(err)
		}
		rep, _ := session.Resolve(rec)
		return rep
	}

	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if rep := stop(); rep.State != session.StateWorking || rep.Detail != limitPausedDetail {
		t.Errorf("report %+v at the limit: want working, %q", rep, limitPausedDetail)
	}
	if err := os.WriteFile(path, []byte("went on"), 0o644); err != nil {
		t.Fatal(err)
	}
	if rep := stop(); rep.State != session.StateIdle {
		t.Errorf("report %+v past the limit: want idle", rep)
	}
}

// pauseClock records the session's report when the run sleeps, then
// ends the limit.
type pauseClock struct {
	path string
	st   *session.Store
	sid  string
	seen *session.Report
}

func (c pauseClock) Now() time.Time { return time.Unix(0, 0) }
func (c pauseClock) Sleep(time.Duration) {
	if rec, err := c.st.Read(c.sid); err == nil {
		*c.seen, _ = session.Resolve(rec)
	}
	os.WriteFile(c.path, []byte("done"), 0o644)
}

// While a session is paused at its usage limit its record says so.
func TestSessionRunRecordsPause(t *testing.T) {
	f, s, path := limitFixture(t)
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	st := session.NewStore(session.Dir(f.common))
	if _, err := st.Update(s.ID, func(r *session.Record) {
		r.Status.Hook = &session.Report{State: session.StateIdle}
	}); err != nil {
		t.Fatal(err)
	}
	var seen session.Report
	useClock(t, pauseClock{path, st, s.ID, &seen})
	f.ctx.Start = func(Program) (Process, error) { return &fakeProcess{}, nil }

	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	if seen.State != session.StateWorking || !strings.HasPrefix(seen.Detail, limitPausedDetail+" until ") {
		t.Errorf("report during the pause %+v: want working, %s until ...", seen, limitPausedDetail)
	}
}

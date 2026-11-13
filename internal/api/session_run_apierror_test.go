package api

import (
	"bytes"
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

// apiErrAdapter is the limit-fake agent whose transcript, lines of text,
// ends in an API error when its last line starts "ERR ".
type apiErrAdapter struct{ limitAdapter }

const apiErrAgent = "apierr-fake"

var apiErrRegister sync.Once

func (apiErrAdapter) Name() string { return apiErrAgent }
func (apiErrAdapter) APIError(t []byte) (string, bool) {
	lines := strings.Split(strings.TrimSpace(string(t)), "\n")
	if last := lines[len(lines)-1]; strings.HasPrefix(last, "ERR ") {
		return strings.TrimPrefix(last, "ERR "), true
	}
	return "", false
}

// apiErrRun is one session.run of a session whose agent stays open.
type apiErrRun struct {
	t     *testing.T
	f     *sessionFixture
	sid   string
	path  string
	st    *session.Store
	proc  *blockingProcess
	done  chan error
	mu    sync.Mutex
	msgs  []string
	seen  []session.Report // the record's report at each send
	wrote chan struct{}
}

// startAPIErr sets up the session (autonomous or not, last reported working
// by its hook, with transcript before) and starts running it.
func startAPIErr(t *testing.T, autonomous bool, before string) *apiErrRun {
	t.Helper()
	return startAPIErrWith(t, autonomous, before, nil)
}

// startAPIErrWith is startAPIErr with its own Start, built from the run.
func startAPIErrWith(t *testing.T, autonomous bool, before string, start func(*apiErrRun) func(Program) (Process, error)) *apiErrRun {
	t.Helper()
	f, s, path := limitFixture(t)
	apiErrRegister.Do(func() { adapter.Register(apiErrAdapter{}) })
	st := session.NewStore(session.Dir(f.common))
	if _, err := st.Update(s.ID, func(r *session.Record) {
		r.Agent, r.Autonomous = apiErrAgent, autonomous
		r.Status.Hook = &session.Report{State: session.StateWorking, At: time.Now()}
	}); err != nil {
		t.Fatal(err)
	}
	if before != "" {
		if err := os.WriteFile(path, []byte(before), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	oldPoll, oldBackoff, oldSend := limitPoll, apiErrorBackoff, sendAPIErrorMessage
	limitPoll = time.Millisecond
	apiErrorBackoff = []time.Duration{40 * time.Millisecond, 40 * time.Millisecond, 40 * time.Millisecond}
	t.Cleanup(func() { limitPoll, apiErrorBackoff, sendAPIErrorMessage = oldPoll, oldBackoff, oldSend })
	r := &apiErrRun{t: t, f: f, sid: s.ID, path: path, st: st, proc: &blockingProcess{killed: make(chan struct{})}, done: make(chan error, 1)}
	sendAPIErrorMessage = func(_ *Ctx, sid, text string) error {
		rec, err := st.Read(sid)
		if err != nil {
			t.Error(err)
		}
		rep, _ := session.Resolve(rec)
		r.mu.Lock()
		r.msgs, r.seen = append(r.msgs, text), append(r.seen, rep)
		r.mu.Unlock()
		return nil
	}
	started, once := make(chan struct{}), sync.Once{}
	startFn := func(Program) (Process, error) { return r.proc, nil }
	if start != nil {
		startFn = start(r)
	}
	f.ctx.Start = func(p Program) (Process, error) {
		proc, err := startFn(p)
		once.Do(func() { close(started) })
		return proc, err
	}
	go func() {
		_, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
		r.done <- err
	}()
	<-started
	return r
}

func (r *apiErrRun) write(text string) {
	r.t.Helper()
	if err := os.WriteFile(r.path, []byte(text), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r *apiErrRun) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.msgs)
}

// waitMsgs waits until n messages were sent.
func (r *apiErrRun) waitMsgs(n int) {
	r.t.Helper()
	for end := time.Now().Add(5 * time.Second); r.count() < n; {
		if time.Now().After(end) {
			r.t.Fatalf("%d messages, want %d", r.count(), n)
		}
		time.Sleep(time.Millisecond)
	}
}

// finish ends the agent and the run, and returns the session record.
func (r *apiErrRun) finish() session.Record {
	r.t.Helper()
	r.proc.Kill()
	if err := <-r.done; err != nil {
		r.t.Fatal(err)
	}
	rec, err := r.st.Read(r.sid)
	if err != nil {
		r.t.Fatal(err)
	}
	return rec
}

// apiErrored is the inbox's view of the session as a step of a live run.
func (r *apiErrRun) apiErrored() []RunAPIError {
	state := run.State{Current: "build", Session: map[string]string{"build": r.sid}}
	return runAPIErrorAll(state, true, r.st)
}

func TestSessionRunAPIErrorRepromptsAfterBackoff(t *testing.T) {
	r := startAPIErr(t, true, "")
	r.write("did work\nERR 403 Cloudflare")
	time.Sleep(15 * time.Millisecond)
	if n := r.count(); n != 0 {
		t.Fatalf("%d messages before the backoff, want none", n)
	}
	r.waitMsgs(1)
	r.mu.Lock()
	msg, rep := r.msgs[0], r.seen[0]
	r.mu.Unlock()
	if !strings.Contains(msg, "403 Cloudflare") {
		t.Errorf("message %q does not name the error", msg)
	}
	if rep.State != session.StateIdle || strings.HasPrefix(rep.Detail, session.APIErrorDetail) {
		t.Errorf("record at the send %+v: want plain idle", rep)
	}
	if got := r.apiErrored(); len(got) != 0 {
		t.Errorf("inbox lists %+v while retries remain", got)
	}
	r.write("did work\nERR 403 Cloudflare\nrecovered")
	r.finish()
}

func TestSessionRunAPIErrorCancelledByAssistantLine(t *testing.T) {
	r := startAPIErr(t, true, "")
	r.write("ERR boom")
	time.Sleep(10 * time.Millisecond)
	r.write("ERR boom\nfine")
	time.Sleep(120 * time.Millisecond)
	if n := r.count(); n != 0 {
		t.Errorf("%d messages, want the send cancelled", n)
	}
	r.finish()
}

func TestSessionRunAPIErrorGivesUpAfterThree(t *testing.T) {
	r := startAPIErr(t, true, "")
	tr := "ERR 403 Cloudflare"
	r.write(tr)
	for i := 1; i <= 3; i++ {
		r.waitMsgs(i)
		tr += "\nmsg\nERR 403 Cloudflare"
		r.write(tr)
	}
	for end := time.Now().Add(5 * time.Second); ; {
		rep, _ := session.Resolve(func() session.Record { rec, _ := r.st.Read(r.sid); return rec }())
		if strings.HasPrefix(rep.Detail, session.APIErrorDetail) {
			break
		}
		if time.Now().After(end) {
			t.Fatal("record never left on the API error")
		}
		time.Sleep(time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	if n := r.count(); n != 3 {
		t.Errorf("%d messages, want 3", n)
	}
	rec := r.finish()
	if rec.Status.Hook == nil || rec.Status.Hook.Detail != session.APIErrorDetail+"403 Cloudflare" {
		t.Errorf("hook report %+v: want %q", rec.Status.Hook, session.APIErrorDetail+"403 Cloudflare")
	}
}

func TestSessionRunAPIErrorAgainAfterRecovery(t *testing.T) {
	r := startAPIErr(t, true, "")
	r.write("ERR one")
	r.waitMsgs(1)
	r.write("ERR one\nok")
	time.Sleep(20 * time.Millisecond)
	r.write("ERR one\nok\nERR two")
	r.waitMsgs(2)
	r.mu.Lock()
	second := r.msgs[1]
	r.mu.Unlock()
	if !strings.Contains(second, "two") {
		t.Errorf("second message %q does not name the new error", second)
	}
	r.write("ERR one\nok\nERR two\nok")
	r.finish()
}

func TestSessionRunAPIErrorAttendedRecordsAtOnce(t *testing.T) {
	r := startAPIErr(t, false, "")
	r.write("ERR 403 Cloudflare")
	for end := time.Now().Add(5 * time.Second); ; {
		if got := r.apiErrored(); len(got) == 1 && strings.Contains(got[0].Detail, "403 Cloudflare") {
			break
		}
		if time.Now().After(end) {
			t.Fatal("the attended session's API error was never recorded")
		}
		time.Sleep(time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	if n := r.count(); n != 0 {
		t.Errorf("%d messages to an owner-attended session, want none", n)
	}
	r.finish()
}

func TestSessionRunAPIErrorBeforeBaseIsNothing(t *testing.T) {
	r := startAPIErr(t, true, "ERR old\n")
	r.write("ERR old\nstill working")
	time.Sleep(100 * time.Millisecond)
	if n := r.count(); n != 0 {
		t.Errorf("%d messages for an error before the start, want none", n)
	}
	if got := r.apiErrored(); len(got) != 0 {
		t.Errorf("inbox lists %+v", got)
	}
	r.finish()
}

// lateWriteClock writes text to the file path points at when the run
// sleeps; the path is known once the run is set up.
type lateWriteClock struct {
	path *string
	text string
}

func (c lateWriteClock) Now() time.Time { return time.Unix(0, 0) }
func (c lateWriteClock) Sleep(time.Duration) {
	os.WriteFile(*c.path, []byte(c.text), 0o644)
}

func TestSessionRunAPIErrorUsageLimitKeepsPause(t *testing.T) {
	var path string
	useClock(t, lateWriteClock{&path, "done"})
	starts := 0
	r := startAPIErrWith(t, true, "", func(r *apiErrRun) func(Program) (Process, error) {
		path = r.path
		return func(Program) (Process, error) {
			starts++
			if starts == 1 {
				return r.proc, nil
			}
			return &fakeProcess{}, nil
		}
	})
	// A usage-limit line: the pause happens, no message is sent.
	r.write("limited")
	if err := <-r.done; err != nil {
		t.Fatal(err)
	}
	if len(limitLaunches) != 2 || !limitLaunches[1].Resume {
		t.Errorf("launches %+v: want the usage-limit pause and a resume", limitLaunches)
	}
	if n := r.count(); n != 0 {
		t.Errorf("%d messages at the usage limit, want none", n)
	}
}

// watchFixture is an apiErrorWatch on a fresh session store in dir, with
// the messages it sends counted and its diagnostics kept.
func watchFixture(t *testing.T, autonomous bool) (w *apiErrorWatch, dir string, stderr *bytes.Buffer, sent *int) {
	t.Helper()
	dir = filepath.Join(t.TempDir(), "sessions")
	st := session.NewStore(dir)
	if err := st.Write(session.Record{ID: "s1"}); err != nil {
		t.Fatal(err)
	}
	stderr, sent = &bytes.Buffer{}, new(int)
	old, oldBackoff := sendAPIErrorMessage, apiErrorBackoff
	t.Cleanup(func() { sendAPIErrorMessage, apiErrorBackoff = old, oldBackoff })
	sendAPIErrorMessage = func(*Ctx, string, string) error { *sent++; return nil }
	ctx := &Ctx{Now: time.Now, Stderr: stderr}
	w = &apiErrorWatch{ctx: ctx, a: apiErrAdapter{}, st: st, sid: "s1", autonomous: autonomous}
	return w, dir, stderr, sent
}

// A record that cannot be written is reported, and written on a later poll.
func TestAPIErrorWatchReportsFailedRecord(t *testing.T) {
	w, dir, stderr, _ := watchFixture(t, false)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	w.poll([]byte("ERR boom"))
	if !strings.Contains(stderr.String(), "session s1") {
		t.Errorf("stderr %q: want the failed record reported", stderr.String())
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	w.poll([]byte("ERR boom"))
	rec, err := w.st.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Hook == nil || rec.Status.Hook.Detail != session.APIErrorDetail+"boom" {
		t.Errorf("hook report %+v: want the API error recorded on a later poll", rec.Status.Hook)
	}
}

// A fourth consecutive API error goes to the inbox at once, whatever the
// third wait.
func TestAPIErrorWatchFourthErrorAtOnce(t *testing.T) {
	w, _, _, sent := watchFixture(t, true)
	apiErrorBackoff = []time.Duration{0, 0, 0}
	tr := "ERR a"
	for i := 0; i < 3; i++ {
		w.poll([]byte(tr))
		if i == 2 {
			apiErrorBackoff = []time.Duration{0, 0, time.Hour}
		}
		w.poll([]byte(tr))
		if *sent != i+1 {
			t.Fatalf("%d messages, want %d", *sent, i+1)
		}
		tr += "\nmsg\nERR again"
	}
	w.poll([]byte(tr))
	rec, err := w.st.Read("s1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status.Hook == nil || rec.Status.Hook.Detail != session.APIErrorDetail+"again" {
		t.Errorf("hook report %+v: want the fourth error recorded at once", rec.Status.Hook)
	}
	if *sent != 3 {
		t.Errorf("%d messages, want 3", *sent)
	}
}

// Each later wait runs from when the new error is seen, not from the send.
func TestAPIErrorWatchWaitsFromErrorSeen(t *testing.T) {
	w, _, _, sent := watchFixture(t, true)
	apiErrorBackoff = []time.Duration{0, time.Second, time.Second}
	w.poll([]byte("ERR a"))
	w.poll([]byte("ERR a"))
	if *sent != 1 {
		t.Fatalf("%d messages, want 1", *sent)
	}
	time.Sleep(700 * time.Millisecond)
	tr := []byte("ERR a\nmsg\nERR b")
	w.poll(tr)
	time.Sleep(400 * time.Millisecond)
	w.poll(tr)
	if *sent != 1 {
		t.Errorf("%d messages 400ms after the second error, want the wait to run from it", *sent)
	}
}

package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
)

// messageFixture is a session started on a change set, its tab open on
// block:b1, with the launch's rex call forgotten.
type messageFixture struct {
	*sessionFixture
	st  *session.Store
	sid string
}

func newMessageFixture(t *testing.T) *messageFixture {
	t.Helper()
	f := newSessionFixture(t, "palette")
	newChange(t, f.ctx, `{"intent":"Pick a palette","title":"Palette"}`)
	s := startSession(t, f.ctx, `{"permission_mode":"manual"}`).Session
	f.rex.calls, f.rex.dirs = nil, nil
	return &messageFixture{sessionFixture: f, st: session.NewStore(session.Dir(f.common)), sid: s.ID}
}

// setHook records a hook report on the session, as session.hook would.
func (f *messageFixture) setHook(t *testing.T, state, kind string) {
	t.Helper()
	if _, err := f.st.Update(f.sid, func(r *session.Record) {
		r.Status.Hook = &session.Report{State: state, Kind: kind, At: time.Now().UTC()}
	}); err != nil {
		t.Fatal(err)
	}
}

// messageFile is the absolute path of message name in the session's folder.
func (f *messageFixture) messageFile(name string) string {
	return filepath.Join(f.worktree, ".rhiza", "local", "messages", f.sid, name)
}

func (f *messageFixture) noRex(t *testing.T, when string) {
	t.Helper()
	if len(f.rex.calls) != 0 {
		t.Fatalf("%s: rex ran %q", when, f.rex.calls)
	}
}

func TestSessionMessageWaitsForThePrompt(t *testing.T) {
	f := newMessageFixture(t)
	f.setHook(t, session.StateBlocked, session.KindQuestion)
	c := &scriptClock{t: t, now: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)}
	c.script = []func(){
		func() { f.noRex(t, "poll 1, a question open") },
		func() {
			f.noRex(t, "poll 2, a question open")
			f.setHook(t, session.StateWorking, "")
		},
		func() {
			f.noRex(t, "poll 3, working")
			f.setHook(t, session.StateIdle, "")
		},
	}
	useClock(t, c)

	resp, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid[:8]+`","text":"Use the blue palette\nand keep it calm"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionMessageResponse)
	file := f.messageFile("0001.md")
	rel := ".rhiza/local/messages/" + f.sid + "/0001.md"
	if r.Session != f.sid || r.File != file || r.Waited != "750ms" || r.DryRun {
		t.Errorf("response %+v", r)
	}
	if r.From != "test" {
		t.Errorf("from %q, want the sender's handle test", r.From)
	}
	pointer := "rhiza: message from test: read " + rel + " and act on it."
	if r.Pointer != pointer {
		t.Errorf("pointer %q, want %q", r.Pointer, pointer)
	}
	want := [][]string{
		{"rex", "send", "-b", "block:b1", pointer},
		{"rex", "send-key", "-b", "block:b1", "enter"},
	}
	if !slices.EqualFunc(f.rex.calls, want, slices.Equal) {
		t.Fatalf("rex calls %q, want %q", f.rex.calls, want)
	}
	if f.rex.dirs[0] != f.worktree {
		t.Errorf("rex ran in %s, want the worktree", f.rex.dirs[0])
	}

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	for _, s := range []string{
		"# Message from test\n\n",
		"\n- From: test (01",
		"\n- Session: " + f.sid + "\n",
		"\n- Step: build\n",
		"\n- At: 2026-10-24T12:00:00Z\n\nUse the blue palette\nand keep it calm\n",
	} {
		if !strings.Contains(got, s) {
			t.Errorf("message file lacks %q:\n%s", s, got)
		}
	}
	if !strings.HasSuffix(got, "calm\n") {
		t.Errorf("message file does not end in one newline: %q", got)
	}
	h := human(t, r)
	if !strings.Contains(h, "Sent session "+f.sid+" a message from test\n") || !strings.Contains(h, "- Pointer: "+pointer+"\n") {
		t.Errorf("human:\n%s", h)
	}

	// A second message takes the next number, and an idle session gets it
	// without waiting.
	f.rex.calls = nil
	resp, err = callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"And rounder corners"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionMessageResponse); r.File != f.messageFile("0002.md") || r.Waited != "0s" || len(f.rex.calls) != 2 {
		t.Errorf("second message %+v, rex %q", r, f.rex.calls)
	}
}

func TestSessionMessageExitedSession(t *testing.T) {
	f := newMessageFixture(t)
	code := 0
	if _, err := f.st.Update(f.sid, func(r *session.Record) {
		r.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now().UTC()}
	}); err != nil {
		t.Fatal(err)
	}
	_, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"hello"}`)
	wantExit(t, err, 3)
	if e := AsError(err); !strings.Contains(e.Message, "has exited; it can take no message") || e.Next != "rhiza session show "+f.sid {
		t.Errorf("error %+v", e)
	}
	f.noRex(t, "exited")
	if _, err := os.Stat(filepath.Dir(f.messageFile("0001.md"))); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("an exited session got a message folder: %v", err)
	}
}

func TestSessionMessageExitWhileWaiting(t *testing.T) {
	f := newMessageFixture(t)
	f.setHook(t, session.StateBlocked, session.KindPermission)
	code := 1
	useClock(t, &scriptClock{t: t, script: []func(){func() {
		if _, err := f.st.Update(f.sid, func(r *session.Record) {
			r.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now().UTC()}
		}); err != nil {
			t.Fatal(err)
		}
	}}})
	_, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"hello"}`)
	wantExit(t, err, 3)
	if msg := AsError(err).Message; !strings.Contains(msg, "exited before it reached its prompt; the message is in "+f.messageFile("0001.md")) {
		t.Errorf("message %q", msg)
	}
	f.noRex(t, "exited while waiting")
}

func TestSessionMessageTimeout(t *testing.T) {
	f := newMessageFixture(t)
	f.setHook(t, session.StateBlocked, session.KindQuestion)
	c := &scriptClock{t: t}
	for range 4 {
		c.script = append(c.script, func() { f.noRex(t, "a question open") })
	}
	useClock(t, c)
	_, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"hello","timeout":"1s"}`)
	wantExit(t, err, 3)
	e := AsError(err)
	file := f.messageFile("0001.md")
	if !strings.Contains(e.Message, "did not reach its prompt within 1s; the message is in "+file) || e.Next != "rhiza session show "+f.sid {
		t.Errorf("error %+v", e)
	}
	if c.sleeps != 4 {
		t.Errorf("polled %d times, want 4 in 1s at 250ms", c.sleeps)
	}
	f.noRex(t, "timed out")
	if _, err := os.Stat(file); err != nil {
		t.Errorf("the message file was not kept: %v", err)
	}
}

func TestSessionMessageDryRun(t *testing.T) {
	f := newMessageFixture(t)
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionMessageResponse)
	if !r.DryRun || r.File != f.messageFile("0001.md") || r.From != "test" || !strings.Contains(r.Pointer, "0001.md") {
		t.Errorf("response %+v", r)
	}
	f.noRex(t, "dry run")
	if _, err := os.Stat(filepath.Dir(r.File)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a dry run made the message folder: %v", err)
	}
	if h := human(t, r); !strings.Contains(h, "Would send session ") || !strings.Contains(h, "This was a dry run: nothing was written or sent.") {
		t.Errorf("human:\n%s", h)
	}
}

func TestSessionMessageRefusals(t *testing.T) {
	cases := []struct {
		name, req string
		exit      int
	}{
		{"BlankText", `{"session":"SID","text":"  \n "}`, 2},
		{"BadTimeout", `{"session":"SID","text":"hi","timeout":"soon"}`, 2},
		{"ZeroTimeout", `{"session":"SID","text":"hi","timeout":"0s"}`, 2},
		{"UnknownSession", `{"session":"deadbeef","text":"hi"}`, 4},
		{"NoTab", `{"session":"SID","text":"hi"}`, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newMessageFixture(t)
			if c.name == "NoTab" {
				if _, err := f.st.Update(f.sid, func(r *session.Record) { r.Tab.Block = "" }); err != nil {
					t.Fatal(err)
				}
			}
			_, err := callOp(t, f.ctx, "session.message", strings.ReplaceAll(c.req, "SID", f.sid))
			wantExit(t, err, c.exit)
			f.noRex(t, c.name)
			if _, err := os.Stat(filepath.Dir(f.messageFile("0001.md"))); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("a refused message made the message folder: %v", err)
			}
		})
	}
}

func TestTypePointerNoExec(t *testing.T) {
	f := newMessageFixture(t)
	f.ctx.Exec = nil
	_, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"hello"}`)
	if err == nil || !strings.Contains(err.Error(), "this caller runs no programs") {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(filepath.Dir(f.messageFile("0001.md"))); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a caller that runs no programs got a message folder: %v", err)
	}
	if _, err := typePointer(f.ctx, f.st, f.sid, "x", time.Second); err == nil {
		t.Error("typePointer without Exec did not fail")
	}
}

func TestTypePointerOneLine(t *testing.T) {
	if got := messagePointer("a", "x\ny\tz.md"); got != "rhiza: message from a: read x y z.md and act on it." {
		t.Errorf("pointer %q", got)
	}
}

// A session that is not working keeps nothing for its next tool-call hook:
// one blocked on a question or menu, idle at its prompt or still starting
// is left to typePointer, which waits for the prompt.
func TestQueuePointerLeavesNonWorkingSessions(t *testing.T) {
	for _, tc := range []struct{ state, kind string }{
		{session.StateBlocked, session.KindQuestion},
		{session.StateBlocked, session.KindPermission},
		{session.StateIdle, ""},
	} {
		t.Run(tc.state+"/"+tc.kind, func(t *testing.T) {
			f := newMessageFixture(t)
			f.setHook(t, tc.state, tc.kind)
			queued, err := queuePointer(f.st, f.sid, "rhiza: message from test: read x and act on it.")
			if err != nil || queued {
				t.Fatalf("queued %v, err %v: want neither", queued, err)
			}
			rec, err := f.st.Read(f.sid)
			if err != nil {
				t.Fatal(err)
			}
			if len(rec.Pending) != 0 {
				t.Errorf("pending %v, want none", rec.Pending)
			}
		})
	}
}

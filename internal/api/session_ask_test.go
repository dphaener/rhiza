package api

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// askAdapter is a fake agent whose hook payload is
// {"session_id":..., "ask":...}: a non-empty ask is a question to the
// owner, which its AskReply denies as {"deny":<reason>}.
type askAdapter struct{}

const askAgent = "ask-fake"

func init() { adapter.Register(askAdapter{}) }

type askPayload struct {
	SessionID string `json:"session_id"`
	Ask       string `json:"ask"`
}

func (askAdapter) Name() string                             { return askAgent }
func (askAdapter) PermissionModes() []string                { return []string{"manual"} }
func (askAdapter) Bypass(string) bool                       { return false }
func (askAdapter) Command(adapter.Launch) ([]string, error) { return []string{askAgent}, nil }
func (askAdapter) HookStatus(event string, raw []byte) (adapter.HookReport, bool, error) {
	var p askPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return adapter.HookReport{}, false, err
	}
	switch {
	case event == "Stop":
		return adapter.HookReport{SessionID: p.SessionID, State: "idle"}, true, nil
	case p.Ask != "":
		return adapter.HookReport{SessionID: p.SessionID, State: "blocked", Kind: "question", Detail: p.Ask}, true, nil
	}
	return adapter.HookReport{SessionID: p.SessionID, State: "working"}, true, nil
}
func (askAdapter) AskReply(event string, raw []byte, reason string) ([]byte, bool, error) {
	var p askPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, false, err
	}
	if event == "Stop" || p.Ask == "" {
		return nil, false, nil
	}
	b, err := json.Marshal(map[string]string{"deny": reason})
	return b, true, err
}
func (askAdapter) HookMessage(string, []byte) (adapter.Message, bool, error) {
	return adapter.Message{}, false, nil
}
func (askAdapter) TranscriptGlobs(string, func(string) string) []string { return nil }
func (askAdapter) Answers([]byte) ([]adapter.Answer, error)             { return nil, nil }
func (askAdapter) TokenUse([]byte) (adapter.TokenUse, error)            { return adapter.TokenUse{}, nil }

// askFixture is a message fixture whose session is the ask-fake agent's,
// autonomous or not.
func askFixture(t *testing.T, autonomous bool) *messageFixture {
	t.Helper()
	f := newMessageFixture(t)
	if _, err := f.st.Update(f.sid, func(r *session.Record) { r.Agent, r.Autonomous = askAgent, autonomous }); err != nil {
		t.Fatal(err)
	}
	return f
}

// askHook runs session.hook for event with payload on stdin and returns
// the response and what it prints for the agent.
func askHook(t *testing.T, ctx *Ctx, event, payload string) (SessionHookResponse, string) {
	t.Helper()
	var stderr strings.Builder
	ctx.Stdin, ctx.Stderr = strings.NewReader(payload), &stderr
	resp, err := callOp(t, ctx, "session.hook", `{"event":"`+event+`"}`)
	if err != nil || stderr.Len() != 0 {
		t.Fatalf("session.hook: err %v, stderr %q", err, stderr.String())
	}
	r := resp.(SessionHookResponse)
	return r, human(t, r)
}

// TestSessionHookAutonomousAskDenied: in an autonomous session a question
// to the owner is denied on stdout with the fixed reason and kept as the
// asked mark, which later hooks leave alone; the JSON response is as
// before.
func TestSessionHookAutonomousAskDenied(t *testing.T) {
	f := askFixture(t, true)
	r, out := askHook(t, f.ctx, "PreToolUse", `{"session_id":"`+f.sid+`","ask":"Which store?"}`)
	want := `{"deny":"` + autonomousAskReason + `"}` + "\n"
	if out != want {
		t.Errorf("hook printed %q\nwant %q", out, want)
	}
	if !r.Recorded || r.State != "blocked" || r.Kind != "question" {
		t.Errorf("response %+v", r)
	}
	if b, _ := json.Marshal(r); strings.Contains(string(b), "deny") {
		t.Errorf("JSON response carries the reply: %s", b)
	}
	rec, err := f.st.Read(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	a := rec.Status.Asked
	if a == nil || a.State != session.StateBlocked || a.Kind != session.KindQuestion || a.Detail != "Which store?" || a.At.IsZero() {
		t.Fatalf("asked %+v", a)
	}

	// The session goes idle: the hook moves on, the asked mark stays, and
	// nothing is printed.
	if _, out := askHook(t, f.ctx, "Stop", `{"session_id":"`+f.sid+`"}`); out != "" {
		t.Errorf("Stop printed %q", out)
	}
	rec, _ = f.st.Read(f.sid)
	if rec.Status.Hook.State != session.StateIdle || rec.Status.Asked == nil || *rec.Status.Asked != *a {
		t.Errorf("status %+v asked %+v", rec.Status.Hook, rec.Status.Asked)
	}
	if rep, _ := session.Resolve(rec); rep.State != session.StateIdle {
		t.Errorf("Resolve weighs the asked mark: %+v", rep)
	}
}

// TestSessionHookNormalAskUnchanged: outside an autonomous session the
// question is the owner's: nothing printed, no asked mark.
func TestSessionHookNormalAskUnchanged(t *testing.T) {
	f := askFixture(t, false)
	r, out := askHook(t, f.ctx, "PreToolUse", `{"session_id":"`+f.sid+`","ask":"Which store?"}`)
	if out != "" || !r.Recorded || r.State != "blocked" {
		t.Errorf("hook printed %q, response %+v", out, r)
	}
	rec, _ := f.st.Read(f.sid)
	if rec.Status.Asked != nil || rec.Status.Hook.Detail != "Which store?" {
		t.Errorf("status %+v", rec.Status)
	}
}

// TestAgentStatusAsked: a live session with an asked mark is blocked on
// its question whatever its hook says; an exited one is done.
func TestAgentStatusAsked(t *testing.T) {
	f := askFixture(t, true)
	r := &stepRunner{ctx: f.ctx, sessions: f.st}
	askHook(t, f.ctx, "PreToolUse", `{"session_id":"`+f.sid+`","ask":"Which store?"}`)
	askHook(t, f.ctx, "Stop", `{"session_id":"`+f.sid+`"}`)
	s, err := r.AgentStatus(f.sid)
	if err != nil {
		t.Fatal(err)
	}
	if b := s.Blocked; b == nil || b.Kind != "question" || b.Question != "Which store?" || b.TabBlock != "block:b1" || s.Done {
		t.Fatalf("status %+v block %+v", s, s.Blocked)
	}
	rec, _ := f.st.Read(f.sid)
	st := run.State{Current: "build", Session: map[string]string{"build": f.sid}}
	if rb := runBlockedNow(st, f.st); rb == nil || rb.Kind != "question" || rb.Question != "Which store?" {
		t.Errorf("run blocked %+v", rb)
	}
	code := 0
	if _, err := f.st.Update(rec.ID, func(x *session.Record) {
		x.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
	}); err != nil {
		t.Fatal(err)
	}
	if s, _ := r.AgentStatus(f.sid); !s.Done || s.Blocked != nil {
		t.Errorf("exited status %+v", s)
	}
}

// TestDeliver sends the asking session the arbiter's answer as a message
// once it is idle, and clears its asked mark.
func TestDeliver(t *testing.T) {
	f := askFixture(t, true)
	r := &stepRunner{ctx: f.ctx, sessions: f.st}
	askHook(t, f.ctx, "PreToolUse", `{"session_id":"`+f.sid+`","ask":"Which store?"}`)
	askHook(t, f.ctx, "Stop", `{"session_id":"`+f.sid+`"}`)
	f.rex.calls = nil
	if err := r.Deliver(f.sid, "01QUESTION", "use-sqlite", "Use SQLite.\n"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.messageFile("0001.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "\n- At: 2026-10-24T12:00:00Z\n\nThe arbiter answered your question 01QUESTION with Decision use-sqlite:\n\nUse SQLite.\n\nCarry on with your work.\n"
	if !strings.HasSuffix(string(data), want) {
		t.Errorf("message:\n%s", data)
	}
	pointer := "rhiza: message from test: read .rhiza/local/messages/" + f.sid + "/0001.md and act on it."
	wantRex := [][]string{
		{"rex", "send", "-b", "block:b1", pointer},
		{"rex", "send-key", "-b", "block:b1", "enter"},
	}
	if !slices.EqualFunc(f.rex.calls, wantRex, slices.Equal) {
		t.Errorf("rex calls %q", f.rex.calls)
	}
	if rec, _ := f.st.Read(f.sid); rec.Status.Asked != nil {
		t.Errorf("asked mark kept: %+v", rec.Status.Asked)
	}

	// An exited session takes no answer, and keeps its mark.
	askHook(t, f.ctx, "PreToolUse", `{"session_id":"`+f.sid+`","ask":"Again?"}`)
	code := 0
	f.st.Update(f.sid, func(x *session.Record) { x.Status.Exit = &session.Report{State: session.StateExited, Code: &code} })
	if err := r.Deliver(f.sid, "01QUESTION", "use-sqlite", "Use SQLite."); err == nil {
		t.Error("delivered to an exited session")
	}
}

func (askAdapter) ContextReply(string, string) ([]byte, bool, error) { return nil, false, nil }
func (askAdapter) UsageLimit([]byte) (adapter.UsageLimit, bool)      { return adapter.UsageLimit{}, false }

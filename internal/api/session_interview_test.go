package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/session"
)

// interviewAdapter is a fake agent whose transcript is
// $HOME/transcripts/<session>.json, a JSON array of adapter.Answer, so the
// capture is tested without any real agent's format.
type interviewAdapter struct{}

const interviewAgent = "interview-fake"

func init() { adapter.Register(interviewAdapter{}) }

func (interviewAdapter) Name() string              { return interviewAgent }
func (interviewAdapter) PermissionModes() []string { return []string{"manual"} }
func (interviewAdapter) Bypass(string) bool        { return false }
func (interviewAdapter) Command(adapter.Launch) ([]string, error) {
	return []string{"interview-fake"}, nil
}
func (interviewAdapter) HookStatus(string, []byte) (adapter.HookReport, bool, error) {
	return adapter.HookReport{}, false, nil
}
func (interviewAdapter) AskReply(string, []byte, string) ([]byte, bool, error) {
	return nil, false, nil
}
func (interviewAdapter) HookMessage(string, []byte) (adapter.Message, bool, error) {
	return adapter.Message{}, false, nil
}
func (interviewAdapter) TranscriptGlobs(sid string, getenv func(string) string) []string {
	if getenv("HOME") == "" {
		return nil
	}
	return []string{filepath.Join(getenv("HOME"), "transcripts", sid+".json")}
}
func (interviewAdapter) Answers(transcript []byte) ([]adapter.Answer, error) {
	var out []adapter.Answer
	err := json.Unmarshal(transcript, &out)
	return out, err
}
func (interviewAdapter) TokenUse([]byte) (adapter.TokenUse, error) {
	return adapter.TokenUse{}, nil
}

var (
	pickedAnswer = adapter.Answer{ToolUseID: "t1", Question: "Which store?", Header: "Store",
		Options: []string{"SQLite (Recommended)", "Postgres"}, Picked: []string{"SQLite (Recommended)"}}
	typedAnswer = adapter.Answer{ToolUseID: "t1", Question: "Name the tool?",
		Options: []string{"rhiza", "root"}, Typed: "Call it rhizome, for now"}
	mixedAnswer = adapter.Answer{ToolUseID: "t2", Question: "Which checks?", MultiSelect: true,
		Options: []string{"vet", "fmt, strict"}, Picked: []string{"vet", "fmt, strict"}, Typed: "staticcheck",
		Notes: "Only on CI, please."}
)

// fields reads an add request's text and JSON fields into one map.
func fields(t *testing.T, r AddRequest) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, kv := range append(slices.Clone(r.Fields), r.JSONFields...) {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			t.Fatalf("field %q has no =", kv)
		}
		if _, dup := out[k]; dup {
			t.Fatalf("field %s twice", k)
		}
		out[k] = v
	}
	return out
}

func TestAnswerWrites(t *testing.T) {
	sid := "ABCDEF12-3456-4789-8abc-def012345678"
	writes, err := answerWrites(sid, []adapter.Answer{pickedAnswer, typedAnswer, mixedAnswer})
	if err != nil {
		t.Fatal(err)
	}
	var handles []string
	for _, w := range writes {
		handles = append(handles, w.add.Handle)
		if w.add.Kind != "decision" {
			t.Errorf("%s kind %s", w.add.Handle, w.add.Kind)
		}
	}
	if want := []string{"answer-abcdef12-1", "answer-abcdef12-2", "answer-abcdef12-3"}; !slices.Equal(handles, want) {
		t.Fatalf("handles %v, want %v", handles, want)
	}
	where := "session " + sid

	t.Run("PickedOnly", func(t *testing.T) {
		f := fields(t, writes[0].add)
		want := map[string]string{
			"statement": "Which store? — SQLite (Recommended)",
			"question":  "Which store?",
			"session":   sid,
			"picked":    "SQLite (Recommended)",
			"wording":   "the agent's option, picked by the owner in " + where,
			"options":   `["SQLite (Recommended)","Postgres"]`,
		}
		if !mapsEqual(f, want) {
			t.Errorf("fields %v\nwant %v", f, want)
		}
	})
	t.Run("Typed", func(t *testing.T) {
		f := fields(t, writes[1].add)
		want := map[string]string{
			"statement":   "Name the tool? — Call it rhizome, for now",
			"question":    "Name the tool?",
			"session":     sid,
			"quote":       "Call it rhizome, for now",
			"quote_where": where + `, answering "Name the tool?"`,
			"options":     `["rhiza","root"]`,
		}
		if !mapsEqual(f, want) {
			t.Errorf("fields %v\nwant %v", f, want)
		}
	})
	t.Run("MixedWithNotes", func(t *testing.T) {
		f := fields(t, writes[2].add)
		want := map[string]string{
			"statement":   "Which checks? — vet, fmt, strict, staticcheck",
			"question":    "Which checks?",
			"session":     sid,
			"picked":      "vet, fmt, strict",
			"wording":     "the agent's option, picked by the owner in " + where,
			"quote":       "staticcheck",
			"quote_where": where + `, answering "Which checks?"`,
			"notes_quote": "Only on CI, please.",
			"options":     `["vet","fmt, strict"]`,
		}
		if !mapsEqual(f, want) {
			t.Errorf("fields %v\nwant %v", f, want)
		}
	})
	t.Run("NoOptionsNoNotes", func(t *testing.T) {
		w, err := answerWrites(sid, []adapter.Answer{{Question: "Why?", Typed: "Because"}})
		if err != nil {
			t.Fatal(err)
		}
		f := fields(t, w[0].add)
		for _, k := range []string{"options", "notes_quote", "picked", "wording"} {
			if _, ok := f[k]; ok {
				t.Errorf("field %s set: %v", k, f)
			}
		}
	})
	t.Run("Empty", func(t *testing.T) {
		w, err := answerWrites(sid, nil)
		if err != nil || len(w) != 0 {
			t.Errorf("writes %v, %v", w, err)
		}
	})
	t.Run("ShortSession", func(t *testing.T) {
		if _, err := answerWrites("abc", []adapter.Answer{pickedAnswer}); err == nil {
			t.Error("a session ID under 8 characters was accepted")
		}
	})
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || w != v {
			return false
		}
	}
	return true
}

// interviewFixture is a launched session switched to the fake agent, with
// a home folder for its transcript.
func interviewFixture(t *testing.T) (*sessionFixture, SessionInfo, *fakeStart) {
	t.Helper()
	f, s, fs := runFixture(t, `{"permission_mode":"manual"}`)
	if _, err := session.NewStore(session.Dir(f.common)).Update(s.ID, func(r *session.Record) {
		r.Agent = interviewAgent
	}); err != nil {
		t.Fatal(err)
	}
	f.ctx.Home = t.TempDir()
	return f, s, fs
}

// writeTranscript writes the fake agent's transcript of session sid.
func writeTranscript(t *testing.T, home, sid string, answers ...adapter.Answer) {
	t.Helper()
	data, err := json.Marshal(answers)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(home, "transcripts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, sid+".json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func interview(t *testing.T, ctx *Ctx, sid string) SessionInterviewResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "session.interview", `{"id":"`+sid+`"}`)
	if err != nil {
		t.Fatalf("session.interview: %v", err)
	}
	return resp.(SessionInterviewResponse)
}

func proposedHandles(r SessionInterviewResponse) []string {
	var out []string
	for _, p := range r.Proposed {
		out = append(out, p.Handle)
	}
	return out
}

func TestSessionInterview(t *testing.T) {
	f, s, _ := interviewFixture(t)
	writeTranscript(t, f.ctx.Home, s.ID, pickedAnswer, typedAnswer, mixedAnswer)
	cs, err := id.Parse(s.Change)
	if err != nil {
		t.Fatal(err)
	}
	h := func(n string) string { return "answer-" + s.ID[:8] + "-" + n }
	all := []string{h("1"), h("2"), h("3")}
	logSize := func() int { return len(ctxWorkspace(t, testCtx(f.worktree)).log().Entries) }
	before := logSize()

	// A dry run lists everything and writes nothing.
	f.ctx.DryRun = true
	dry := interview(t, f.ctx, s.ID[:8])
	f.ctx.DryRun = false
	if !dry.DryRun || !slices.Equal(proposedHandles(dry), all) || len(dry.Skipped) != 0 {
		t.Errorf("dry run %+v", dry)
	}
	if logSize() != before {
		t.Fatal("a dry run wrote to the worktree's log")
	}
	if out := human(t, dry); !strings.Contains(out, "Would propose 3") || !strings.Contains(out, "nothing was written") {
		t.Errorf("dry run human %q", out)
	}

	r := interview(t, f.ctx, s.ID)
	if r.DryRun || r.Session != s.ID || r.Change != s.Change || !slices.Equal(proposedHandles(r), all) || len(r.Skipped) != 0 {
		t.Fatalf("capture %+v", r)
	}
	if p := r.Proposed[2]; p.Question != "Which checks?" || !slices.Equal(p.Picked, []string{"vet", "fmt, strict"}) || p.Typed != "staticcheck" {
		t.Errorf("proposed %+v", p)
	}
	if logSize() != before+1 {
		t.Errorf("log grew by %d entries, want one", logSize()-before)
	}

	// The Decisions are in the session's change set in the worktree,
	// written by the agent's actor for the launcher.
	ww := ctxWorkspace(t, testCtx(f.worktree))
	m := ww.model(cs)
	for _, handle := range all {
		n, err := m.Resolve(handle)
		if err != nil {
			t.Fatalf("%s: %v", handle, err)
		}
		if n.Kind != "decision" {
			t.Errorf("%s kind %s", handle, n.Kind)
		}
	}
	picked, _ := m.Resolve(h("1"))
	if _, ok := picked.Fields["quote"]; ok {
		t.Errorf("a picked option was quoted as the owner's: %s", picked.Fields["quote"])
	}
	typed, _ := m.Resolve(h("2"))
	if string(typed.Fields["quote"]) != `"Call it rhizome, for now"` {
		t.Errorf("quote %s", typed.Fields["quote"])
	}
	if _, ok := typed.Fields["wording"]; ok {
		t.Error("typed text was recorded as the agent's wording")
	}
	launcher := ctxWorkspace(t, f.ctx)
	if _, err := launcher.model(cs).Resolve(h("1")); err == nil {
		t.Error("the Decision went into the launcher's log, not the worktree's")
	}
	agent, err := m.Resolve(s.Actor)
	if err != nil {
		t.Fatal(err)
	}
	var last *idEntry
	for _, e := range ww.log().Entries {
		if e.ChangeSet == cs {
			last = &idEntry{actor: e.Actor, forID: e.For}
		}
	}
	if last == nil || last.actor != agent.ID || last.forID.String() != s.For {
		t.Errorf("entry by %+v, want actor %s for %s", last, agent.ID, s.For)
	}
	if out := human(t, r); !strings.Contains(out, "Proposed 3 decision(s)") || !strings.Contains(out, h("3")) {
		t.Errorf("human %q", out)
	}

	// A second capture of the same transcript proposes nothing.
	again := interview(t, f.ctx, s.ID)
	if len(again.Proposed) != 0 || !slices.Equal(again.Skipped, all) {
		t.Errorf("second capture %+v", again)
	}
	if logSize() != before+1 {
		t.Error("a second capture wrote to the log")
	}
	if out := human(t, again); !strings.Contains(out, "Skipped "+h("1")) {
		t.Errorf("human %q", out)
	}

	// A later answer is proposed alone.
	writeTranscript(t, f.ctx.Home, s.ID, pickedAnswer, typedAnswer, mixedAnswer,
		adapter.Answer{Question: "Ship it?", Options: []string{"Yes", "No"}, Picked: []string{"Yes"}})
	more := interview(t, f.ctx, s.ID)
	if !slices.Equal(proposedHandles(more), []string{h("4")}) || !slices.Equal(more.Skipped, all) {
		t.Errorf("third capture %+v", more)
	}
}

type idEntry struct{ actor, forID id.ID }

func TestSessionInterviewNoAnswers(t *testing.T) {
	f, s, _ := interviewFixture(t)
	writeTranscript(t, f.ctx.Home, s.ID)
	r := interview(t, f.ctx, s.ID)
	if r.Proposed == nil || len(r.Proposed) != 0 || len(r.Skipped) != 0 {
		t.Errorf("no answers: %+v", r)
	}
	if out := human(t, r); !strings.Contains(out, "nothing to propose") {
		t.Errorf("human %q", out)
	}
}

func TestSessionInterviewNoTranscript(t *testing.T) {
	f, s, _ := interviewFixture(t)
	_, err := callOp(t, f.ctx, "session.interview", `{"id":"`+s.ID+`"}`)
	e := AsError(err)
	if e == nil || e.Code != CodeError || ExitCode(err) != 1 {
		t.Fatalf("no transcript: %v", err)
	}
	if !strings.Contains(e.Message, "no transcript") || e.Next != "rhiza session show "+s.ID {
		t.Errorf("error %q next %q", e.Message, e.Next)
	}

	// No home folder and no environment: nowhere to look.
	f.ctx.Home = ""
	_, err = callOp(t, f.ctx, "session.interview", `{"id":"`+s.ID+`"}`)
	if e := AsError(err); e == nil || e.Code != CodeError || !strings.Contains(e.Message, "nowhere") {
		t.Errorf("no home: %v", err)
	}
}

func TestSessionInterviewRefusals(t *testing.T) {
	f, _, _ := interviewFixture(t)
	_, err := callOp(t, f.ctx, "session.interview", `{"id":""}`)
	wantExit(t, err, 2)
	_, err = callOp(t, f.ctx, "session.interview", `{"id":"ffffffff"}`)
	if e := AsError(err); e == nil || e.Code != CodeNotFound {
		t.Errorf("unknown session: %v", err)
	}
}

func TestSessionRunCapturesAnswers(t *testing.T) {
	f, s, _ := interviewFixture(t)
	writeTranscript(t, f.ctx.Home, s.ID, pickedAnswer, typedAnswer)
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionRunResponse)
	want := []string{"answer-" + s.ID[:8] + "-1", "answer-" + s.ID[:8] + "-2"}
	if r.Code != 3 || !slices.Equal(r.Decisions, want) || len(r.Warnings) != 0 {
		t.Errorf("code %d decisions %v warnings %v", r.Code, r.Decisions, r.Warnings)
	}
	out := human(t, r)
	for _, h := range want {
		if !strings.Contains(out, "- Proposed decision "+h+"\n") {
			t.Errorf("human %q lacks %s", out, h)
		}
	}
	if r.Session.Source != "exit" {
		t.Errorf("status from %s", r.Session.Source)
	}
}

func TestSessionRunCaptureFailureIsAWarning(t *testing.T) {
	f, s, _ := interviewFixture(t)
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatalf("a failed capture failed the run: %v", err)
	}
	r := resp.(SessionRunResponse)
	if r.Code != 3 || r.ExitCode() != 3 || len(r.Decisions) != 0 {
		t.Errorf("code %d exit %d decisions %v", r.Code, r.ExitCode(), r.Decisions)
	}
	// No transcript: the token use is not recorded either, a warning of
	// its own.
	if len(r.Warnings) != 2 || !strings.HasPrefix(r.Warnings[0], "record token use: no transcript") ||
		!strings.Contains(r.Warnings[1], "no transcript") {
		t.Errorf("warnings %v", r.Warnings)
	}
	if !strings.Contains(human(t, r), "- Warning: propose decisions") {
		t.Errorf("human %q", human(t, r))
	}
	rec, err := session.NewStore(session.Dir(f.common)).Read(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if x := rec.Status.Exit; x == nil || x.Code == nil || *x.Code != 3 {
		t.Errorf("exit report %+v", x)
	}
}

func TestSessionRunDryRunCapturesNothing(t *testing.T) {
	f, s, fs := interviewFixture(t)
	writeTranscript(t, f.ctx.Home, s.ID, pickedAnswer)
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionRunResponse); len(r.Decisions) != 0 || len(fs.programs) != 0 {
		t.Errorf("dry run %+v started %v", r, fs.programs)
	}
	f.ctx.DryRun = false
	if r := interview(t, f.ctx, s.ID); len(r.Proposed) != 1 {
		t.Errorf("after a dry run: %+v", r)
	}
}

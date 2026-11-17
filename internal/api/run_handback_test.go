package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// resumeAdapter is the ask-fake agent whose transcript can be found at
// resumeTranscript, when that is set.
type resumeAdapter struct{ askAdapter }

const resumeAgent = "resume-fake"

var resumeTranscript string

func init() { adapter.Register(resumeAdapter{}) }

func (resumeAdapter) Name() string { return resumeAgent }
func (resumeAdapter) TranscriptGlobs(string, func(string) string) []string {
	if resumeTranscript == "" {
		return nil
	}
	return []string{resumeTranscript}
}

func handbackStep() workflow.Step {
	return workflow.Step{Name: "build", Kind: workflow.KindAgent, Handback: "Fix it, see {handback}"}
}

func handbackFixture(t *testing.T, agent string) (*messageFixture, *stepRunner, string) {
	t.Helper()
	f := askFixture(t, false)
	dir := t.TempDir()
	if _, err := f.st.Update(f.sid, func(r *session.Record) { r.Agent, r.Brief = agent, dir }); err != nil {
		t.Fatal(err)
	}
	return f, &stepRunner{ctx: f.ctx, sessions: f.st, change: "C1", run: "R1", worktree: f.worktree}, dir
}

// Serial: it sets the package's resumeTranscript, which the resume fake
// agent registered for the whole binary reads.
func TestHandbackDeliver(t *testing.T) {
	gate := workflow.Step{Name: "gate", Kind: workflow.KindGate, Commands: []string{"make gate"}}
	cmd := workflow.Step{Name: "check", Kind: workflow.KindCommand, Run: "make check"}
	ship := workflow.Step{Name: "ship-check", Kind: workflow.KindCommand, Run: "rhiza ship check --json"}

	t.Run("handback.md names the failing gate command and its output file", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		detail := "gate failed: 1 of 1 commands\ncommand 1 `make gate`: exit 2; full output: /run/gate-gate-1-1.log"
		f.setHook(t, session.StateIdle, "")
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(gate, "fail", detail)); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(dir, "handback.md"))
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Failing step: gate", "`make gate`: output in /run/gate-gate-1-1.log"} {
			if !strings.Contains(string(data), want) {
				t.Errorf("handback.md lacks %q:\n%s", want, data)
			}
		}
	})

	t.Run("a command step's handback.md names the command and its output tail", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(cmd, "fail", "exit 1\nFAIL TestX\nboom")); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(filepath.Join(dir, "handback.md"))
		for _, want := range []string{"`make check`", "FAIL TestX\nboom"} {
			if !strings.Contains(string(data), want) {
				t.Errorf("handback.md lacks %q:\n%s", want, data)
			}
		}
	})

	t.Run("a ship check step's handback.md carries the ship check cause", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(ship, "fail", `exit 1
{"findings":["ci red"]}`)); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(filepath.Join(dir, "handback.md"))
		if !strings.Contains(string(data), `## Cause`) || !strings.Contains(string(data), `"findings":["ci red"]`) {
			t.Errorf("handback.md:\n%s", data)
		}
	})

	t.Run("a live session is messaged with the handback text", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		f.rex.calls = nil
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom")); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(f.messageFile("0001.md"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "Fix it, see " + filepath.Join(dir, "handback.md"); !strings.HasSuffix(strings.TrimSpace(string(data)), want) {
			t.Errorf("message:\n%s\nwant it to end %q", data, want)
		}
		if len(f.rex.calls) != 2 || f.rex.calls[0][1] != "send" {
			t.Errorf("rex calls %q", f.rex.calls)
		}
	})

	t.Run("a live session that said done on an earlier visit has its done mark cleared", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		if _, err := f.st.Update(f.sid, func(x *session.Record) {
			x.Status.Done = &session.Report{State: session.StateDone, At: f.ctx.Now().UTC()}
		}); err != nil {
			t.Fatal(err)
		}
		f.rex.calls = nil
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom")); err != nil {
			t.Fatal(err)
		}
		rec, err := f.st.Read(f.sid)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Status.SaidDone() {
			t.Errorf("done mark %+v kept after the hand-back", rec.Status.Done)
		}
		if rec.Status.Hook == nil || rec.Status.Hook.State != session.StateIdle {
			t.Errorf("hook report %+v: want it left idle", rec.Status.Hook)
		}
		data, err := os.ReadFile(f.messageFile("0001.md"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "Fix it, see " + filepath.Join(dir, "handback.md"); !strings.HasSuffix(strings.TrimSpace(string(data)), want) {
			t.Errorf("message:\n%s\nwant it to end %q", data, want)
		}
		if len(f.rex.calls) != 2 || f.rex.calls[0][1] != "send" {
			t.Errorf("rex calls %q", f.rex.calls)
		}
		if st, err := r.AgentStatus(f.sid); err != nil || st.SaidDone {
			t.Fatalf("status after the hand-back %+v, %v: want not said done", st, err)
		}
		// A fresh session done marks the visit done again.
		resp, err := callOp(t, f.ctx, "session.done", `{"id":"`+f.sid+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		if d := resp.(SessionDoneResponse); d.Already {
			t.Errorf("session done %+v: want a fresh mark, not already done", d)
		}
		if st, err := r.AgentStatus(f.sid); err != nil || !st.SaidDone {
			t.Errorf("status after session done %+v, %v: want said done", st, err)
		}
	})

	t.Run("a failed hand-back leaves a live session's done mark in place", func(t *testing.T) {
		f, r, _ := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		mark := &session.Report{State: session.StateDone, At: f.ctx.Now().UTC()}
		if _, err := f.st.Update(f.sid, func(x *session.Record) { x.Status.Done = mark }); err != nil {
			t.Fatal(err)
		}
		f.rex.calls = nil
		step := workflow.Step{Name: "build", Kind: workflow.KindAgent}
		err := r.Handback(step, f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom"))
		if err == nil || !strings.Contains(err.Error(), "step build") {
			t.Fatalf("error %v: want it to name step build", err)
		}
		rec, err := f.st.Read(f.sid)
		if err != nil {
			t.Fatal(err)
		}
		if !rec.Status.SaidDone() {
			t.Error("the failed hand-back cleared the done mark")
		}
		if _, err := os.Stat(f.messageFile("0001.md")); !os.IsNotExist(err) {
			t.Errorf("message written: %v", err)
		}
		if len(f.rex.calls) != 0 {
			t.Errorf("rex calls %q", f.rex.calls)
		}
	})

	t.Run("an ended session is resumed with the same session ID", func(t *testing.T) {
		f, r, dir := handbackFixture(t, resumeAgent)
		code := 0
		if _, err := f.st.Update(f.sid, func(x *session.Record) {
			x.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
		}); err != nil {
			t.Fatal(err)
		}
		resumeTranscript = filepath.Join(t.TempDir(), f.sid+".jsonl")
		t.Cleanup(func() { resumeTranscript = "" })
		if err := os.WriteFile(resumeTranscript, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		f.rex.calls = nil
		if err := r.Handback(handbackStep(), f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom")); err != nil {
			t.Fatal(err)
		}
		rec, err := f.st.Read(f.sid)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Status.Exited() || rec.Resume != "Fix it, see "+filepath.Join(dir, "handback.md") {
			t.Errorf("record exit %v resume %q", rec.Status.Exit, rec.Resume)
		}
		if len(f.rex.calls) == 0 || !strings.Contains(strings.Join(f.rex.calls[0], " "), f.sid) {
			t.Errorf("rex calls %q: want a tab running session %s", f.rex.calls, f.sid)
		}
	})

	t.Run("an agent step with no handback template fails naming the step", func(t *testing.T) {
		f, r, dir := handbackFixture(t, askAgent)
		f.setHook(t, session.StateIdle, "")
		f.rex.calls = nil
		step := workflow.Step{Name: "build", Kind: workflow.KindAgent}
		err := r.Handback(step, f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom"))
		if err == nil || !strings.Contains(err.Error(), "step build") || !strings.Contains(err.Error(), "no handback template") {
			t.Fatalf("error %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, "handback.md")); !os.IsNotExist(err) {
			t.Errorf("handback.md written: %v", err)
		}
		if _, err := os.Stat(f.messageFile("0001.md")); !os.IsNotExist(err) {
			t.Errorf("message written: %v", err)
		}
		if len(f.rex.calls) != 0 {
			t.Errorf("rex calls %q", f.rex.calls)
		}
	})

	t.Run("a session that cannot be resumed fails naming it", func(t *testing.T) {
		f, r, _ := handbackFixture(t, resumeAgent)
		code := 0
		f.st.Update(f.sid, func(x *session.Record) {
			x.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
		})
		resumeTranscript = ""
		err := r.Handback(handbackStep(), f.sid, run.CauseOf(cmd, "fail", "exit 1\nboom"))
		if err == nil || !strings.Contains(err.Error(), f.sid) || !strings.Contains(err.Error(), "cannot be resumed") {
			t.Fatalf("error %v", err)
		}
		if rec, _ := f.st.Read(f.sid); !rec.Status.Exited() || rec.Resume != "" {
			t.Errorf("record changed: %+v", rec.Status)
		}
	})
}

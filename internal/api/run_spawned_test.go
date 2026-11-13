package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// spawnedRun journals run R1 on the fixture's change set, at step build,
// with two children of one step: the fixture's session (running) and a
// second that failed.
func spawnedRun(t *testing.T, f *messageFixture, change string) {
	t.Helper()
	rs := run.NewStore(run.Dir(f.common))
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	if err := rs.Create("R1", run.Entry{
		Type: run.TypeStarted, Step: "build", At: t0,
		Started: &run.Started{Workflow: "w", Change: change, PermissionMode: "manual", Agent: "claude-code"},
	}); err != nil {
		t.Fatal(err)
	}
	for _, e := range []run.Entry{
		{Type: run.TypeSpawn, Step: "build", Parent: "parent-1", Session: f.sid, Role: "builder", Leaf: "spawned-leaf", Model: "sonnet"},
		{Type: run.TypeSpawn, Step: "build", Parent: "parent-1", Session: "child-2", Role: "verifier", Leaf: "spawned-leaf"},
		{Type: run.TypeSpawnEnded, Step: "build", Parent: "parent-1", Session: "child-2", Outcome: workflow.OutcomeFail, Detail: "exit 3"},
	} {
		e.At = t0.Add(time.Minute)
		if err := rs.Append("R1", e); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := f.st.Update(f.sid, func(r *session.Record) { r.Run, r.Step, r.Parent, r.Role = "R1", "build", "parent-1", "builder" }); err != nil {
		t.Fatal(err)
	}
}

func TestRunStatusSpawned(t *testing.T) {
	f := newMessageFixture(t)
	cs := currentOf(t, f.ctx).String()
	if _, err := callOp(t, f.ctx, "add", `{"kind":"capability","handle":"spawned-leaf","fields":["statement=Show spawned sessions"]}`); err != nil {
		t.Fatal(err)
	}
	spawnedRun(t, f, cs)
	resp, err := callOp(t, f.ctx, "run.status", `{"run":"R1"}`)
	if err != nil {
		t.Fatal(err)
	}
	info := resp.(RunInfo)
	want := []RunSpawned{
		{Step: "build", Parent: "parent-1", Session: f.sid, Role: "builder", Leaf: "spawned-leaf", Statement: "Show spawned sessions", Model: "sonnet", State: "running"},
		{Step: "build", Parent: "parent-1", Session: "child-2", Role: "verifier", Leaf: "spawned-leaf", Statement: "Show spawned sessions", State: "ended", Outcome: "fail", Detail: "exit 3"},
	}
	if len(info.Spawned) != 2 || info.Spawned[0] != want[0] || info.Spawned[1] != want[1] {
		t.Fatalf("spawned %+v\nwant %+v", info.Spawned, want)
	}
	text := human(t, info)
	for _, w := range []string{"## Spawned sessions", "- build, builder on spawned-leaf, model sonnet, session " + f.sid + ": running", "verifier on spawned-leaf, session child-2: ended, fail (exit 3)"} {
		if !strings.Contains(text, w) {
			t.Errorf("human text lacks %q:\n%s", w, text)
		}
	}

	// The same through --json's response, and rhiza session show takes the child's ID.
	show, err := callOp(t, f.ctx, "session.show", `{"id":"`+f.sid+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if rec := show.(SessionInfo); rec.ID != f.sid || rec.Run != "R1" {
		t.Errorf("session show %+v", rec)
	}
}

// A spawned child is opened, messaged and resumed as any session is.
func TestSpawnedSessionOpenResume(t *testing.T) {
	f := askFixture(t, false)
	spawnedRun(t, f, currentOf(t, f.ctx).String())
	f.setHook(t, session.StateIdle, "")

	rexOut := f.rex.out
	f.rex.out = "child working\n> "
	resp, err := callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionPaneResponse); r.Session != f.sid || r.Text != "child working\n> " {
		t.Errorf("pane %+v", r)
	}

	f.rex.calls = nil
	if _, err := callOp(t, f.ctx, "session.message", `{"session":"`+f.sid+`","text":"Also check the web page"}`); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(f.messageFile("0001.md"))
	if err != nil || !strings.Contains(string(data), "Also check the web page") || len(f.rex.calls) == 0 {
		t.Errorf("message %q err %v, rex calls %q", data, err, f.rex.calls)
	}

	// Once the child ended, the Runner's hand-back resume path continues the same conversation.
	dir := t.TempDir()
	code := 0
	if _, err := f.st.Update(f.sid, func(r *session.Record) {
		r.Agent, r.Brief = resumeAgent, dir
		r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
	}); err != nil {
		t.Fatal(err)
	}
	resumeTranscript = filepath.Join(t.TempDir(), f.sid+".jsonl")
	t.Cleanup(func() { resumeTranscript = "" })
	if err := os.WriteFile(resumeTranscript, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
	f.rex.out = rexOut
	r := &stepRunner{ctx: f.ctx, sessions: f.st, change: "C1", run: "R1", worktree: f.worktree}
	cause := run.CauseOf(workflow.Step{Name: "check", Kind: workflow.KindCommand, Run: "make check"}, "fail", "exit 1\nboom")
	if err := r.Handback(handbackStep(), f.sid, cause); err != nil {
		t.Fatal(err)
	}
	rec, err := f.st.Read(f.sid)
	if err != nil || rec.ID != f.sid || rec.Status.Exited() || rec.Parent != "parent-1" || rec.Resume == "" {
		t.Errorf("record %+v err %v", rec, err)
	}
	if len(f.rex.calls) == 0 || !strings.Contains(strings.Join(f.rex.calls[0], " "), f.sid) {
		t.Errorf("rex calls %q: want a tab running session %s", f.rex.calls, f.sid)
	}
}

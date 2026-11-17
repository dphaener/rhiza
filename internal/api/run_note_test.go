package api

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// noteFixture is a session fixture with run R1 at step build, and env the
// environment its calls see.
func noteFixture(t *testing.T) (*sessionFixture, *run.Store, map[string]string) {
	t.Helper()
	f := newSessionFixture(t, "notes")
	store := run.NewStore(run.Dir(f.common))
	t0 := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	if err := store.Create("R1", run.Entry{
		Type: run.TypeStarted, Step: "build", At: t0,
		Started: &run.Started{Workflow: "conduct", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Notes"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("R1", run.Entry{Type: run.TypeEntered, Step: "build", Visit: 1, At: t0}); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{}
	f.ctx.Getenv = func(k string) string { return env[k] }
	return f, store, env
}

// Two notes journaled by run note come back in run status's notes, oldest
// first, with the session, role and leaf of the noting session; the human
// status lists them under Notes.
func TestRunNoteInStatus(t *testing.T) {
	t.Parallel()
	f, _, env := noteFixture(t)
	env[envRun], env[envSession], env[run.EnvRole], env[run.EnvLeaf] = "R1", "s1", "builder", "leaf-a"
	resp, err := callOp(t, f.ctx, "run.note", `{"text":"go test ./internal/api takes 6 minutes: scope it with -run"}`)
	if err != nil {
		t.Fatal(err)
	}
	nr := resp.(RunNoteResponse)
	if nr.Run != "R1" || nr.Note.Role != "builder" || nr.Note.Leaf != "leaf-a" || nr.Note.Session != "s1" || nr.Note.Step != "build" {
		t.Fatalf("note %+v", nr)
	}
	if h := human(t, nr); !strings.HasPrefix(h, "Noted on run R1 (builder on leaf-a, session s1): go test") {
		t.Errorf("human %q", h)
	}
	// Outside the run's env, --run names it; role and leaf follow the env.
	delete(env, envRun)
	env[envSession], env[run.EnvRole], env[run.EnvLeaf] = "s2", "verifier", "leaf-b"
	if _, err := callOp(t, f.ctx, "run.note", `{"text":"the fixture needs a git repo","run":"R1"}`); err != nil {
		t.Fatal(err)
	}

	st, err := callOp(t, f.ctx, "run.status", `{"run":"R1"}`)
	if err != nil {
		t.Fatal(err)
	}
	notes := st.(RunInfo).Notes
	if len(notes) != 2 {
		t.Fatalf("notes %+v", notes)
	}
	want := []RunNote{
		{Step: "build", Session: "s1", Role: "builder", Leaf: "leaf-a", Text: "go test ./internal/api takes 6 minutes: scope it with -run"},
		{Step: "build", Session: "s2", Role: "verifier", Leaf: "leaf-b", Text: "the fixture needs a git repo"},
	}
	for i, n := range notes {
		if n.At == "" {
			t.Errorf("note %d has no time", i)
		}
		n.At = ""
		if n != want[i] {
			t.Errorf("note %d = %+v, want %+v", i, n, want[i])
		}
	}
	h := human(t, st)
	if !strings.Contains(h, "\n## Notes\n\n- "+notes[0].At+" (builder on leaf-a, session s1): go test") ||
		!strings.Contains(h, "(verifier on leaf-b, session s2): the fixture needs a git repo\n") {
		t.Errorf("human status lacks the notes:\n%s", h)
	}
}

// An empty, a multi-line and a 201-character note are usage and journal
// nothing; a note with no --run outside a run names RHIZA_RUN.
func TestRunNoteRefused(t *testing.T) {
	t.Parallel()
	f, store, env := noteFixture(t)
	env[envRun] = "R1"
	for _, text := range []string{"", "  ", "one\ntwo", strings.Repeat("x", 201)} {
		_, err := callOp(t, f.ctx, "run.note", `{"text":`+strconv.Quote(text)+`}`)
		wantExit(t, err, 2)
	}
	delete(env, envRun)
	_, err := callOp(t, f.ctx, "run.note", `{"text":"fine"}`)
	wantExit(t, err, 2)
	if e := AsError(err); !strings.Contains(e.Message, "RHIZA_RUN") {
		t.Errorf("no run: %q does not name RHIZA_RUN", e.Message)
	}
	j, err := store.Read("R1")
	if err != nil {
		t.Fatal(err)
	}
	if n := j.State().Notes; len(n) != 0 {
		t.Fatalf("a refused note was journaled: %+v", n)
	}
	st, err := callOp(t, f.ctx, "run.status", `{"run":"R1"}`)
	if err != nil {
		t.Fatal(err)
	}
	if st.(RunInfo).Notes != nil || strings.Contains(human(t, st), "## Notes") {
		t.Error("status shows notes with none journaled")
	}
}

// A dry run journals nothing.
func TestRunNoteDryRun(t *testing.T) {
	t.Parallel()
	f, store, env := noteFixture(t)
	env[envRun] = "R1"
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "run.note", `{"text":"fine"}`)
	if err != nil || !resp.(RunNoteResponse).DryRun {
		t.Fatalf("dry run: %+v, %v", resp, err)
	}
	j, _ := store.Read("R1")
	if len(j.State().Notes) != 0 {
		t.Fatal("a dry run journaled a note")
	}
}

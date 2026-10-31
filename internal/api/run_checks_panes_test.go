package api

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// paneRunner is a stepRunner of run R1 whose session store holds one
// session, S-first, with a pane in window w9: the anchor a check or fixer
// session splits.
func paneRunner(t *testing.T) (*sessionFixture, *stepRunner) {
	t.Helper()
	f := newSessionFixture(t, "readable-tabs")
	cs := newChange(t, f.ctx, `{"intent":"Readable tabs"}`).ChangeSet
	st := session.NewStore(session.Dir(f.common))
	if err := st.Write(session.Record{ID: "S-first", Run: "R1", Change: cs.String(),
		Tab: session.Tab{Window: "window:w9", Block: "block:a1"}}); err != nil {
		t.Fatal(err)
	}
	return f, &stepRunner{ctx: f.ctx, change: cs.String(), run: "R1", worktree: f.worktree, mode: "manual", sessions: st}
}

// wantSplitOf fails unless the fake rex was asked to split anchor's block
// and never to open a window.
func wantSplitOf(t *testing.T, f *sessionFixture, block string) {
	t.Helper()
	var split bool
	for _, c := range f.rex.calls {
		if slices.Contains(c, "window") && slices.Contains(c, "new") {
			t.Errorf("rex opened a window: %q", c)
		}
		if len(c) > 3 && c[1] == "split" && c[2] == "-b" && c[3] == block {
			split = true
		}
	}
	if !split {
		t.Errorf("no rex split of %s in %q", block, f.rex.calls)
	}
}

// A check session whose spec has an Anchor splits that session's pane and
// opens no window.
func TestStartCheck(t *testing.T) {
	f, r := paneRunner(t)
	f.rex.calls = nil
	step := workflow.Step{Name: "review", Split: workflow.SplitRight}
	sid, err := r.StartCheck(run.CheckSpec{
		Step: step, Role: workflow.RoleCheck, Cycle: 1, Anchor: "S-first",
		Check: workflow.Check{Name: "alpha", Role: workflow.Role{Prompt: "check {check}"}},
		Dir:   filepath.Join(t.TempDir(), "unused"),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantSplitOf(t, f, "block:a1")
	if rec, err := r.sessions.Read(sid); err != nil || strings.TrimSpace(rec.Tab.Block) == "" {
		t.Errorf("record %+v, %v", rec, err)
	}
}

// A change set fixer's session splits its Anchor's pane the same way.
func TestStartFix(t *testing.T) {
	f, r := paneRunner(t)
	f.rex.calls = nil
	step := workflow.Step{Name: "review", Split: workflow.SplitBelow, Fixer: &workflow.Role{Prompt: "fix"}}
	if _, err := r.StartFix(run.FixSpec{Step: step, Cycle: 2, Anchor: "S-first"}); err != nil {
		t.Fatal(err)
	}
	wantSplitOf(t, f, "block:a1")
}

// A session with no Anchor splits the pane of the run's latest session,
// not its first: an earlier step's window may be gone.
func TestStartCheckLatestAnchor(t *testing.T) {
	f, r := paneRunner(t)
	if err := r.sessions.Write(session.Record{ID: "S-later", Run: "R1", Change: r.change, StartedAt: time.Now(),
		Tab: session.Tab{Window: "window:w10", Block: "block:a2"}}); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
	step := workflow.Step{Name: "review", Split: workflow.SplitRight}
	if _, err := r.StartCheck(run.CheckSpec{
		Step: step, Role: workflow.RoleCheck, Cycle: 1,
		Check: workflow.Check{Name: "alpha", Role: workflow.Role{Prompt: "check {check}"}},
		Dir:   filepath.Join(t.TempDir(), "unused"),
	}); err != nil {
		t.Fatal(err)
	}
	wantSplitOf(t, f, "block:a2")
}

// An agent step's session, such as a plan run's interview, splits the pane
// of the run's latest session and opens no window.
func TestStartAgent(t *testing.T) {
	f, r := paneRunner(t)
	f.rex.calls = nil
	if _, err := r.StartAgent(workflow.Step{Name: "interview", Prompt: "go", Split: workflow.SplitRight}, "", nil); err != nil {
		t.Fatal(err)
	}
	wantSplitOf(t, f, "block:a1")
}

// With no session of the run holding a pane, an agent step opens exactly
// one window.
func TestStartAgentFirstSession(t *testing.T) {
	f, r := paneRunner(t)
	r.run = "R-other"
	f.rex.calls = nil
	if _, err := r.StartAgent(workflow.Step{Name: "interview", Prompt: "go"}, "", nil); err != nil {
		t.Fatal(err)
	}
	var n int
	for _, c := range f.rex.calls {
		if slices.Contains(c, "window") && slices.Contains(c, "new") {
			n++
		}
		if len(c) > 1 && c[1] == "split" {
			t.Errorf("rex split: %q", c)
		}
	}
	if n != 1 {
		t.Errorf("window new calls = %d, want 1: %q", n, f.rex.calls)
	}
}

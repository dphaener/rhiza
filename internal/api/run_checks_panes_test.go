package api

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

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

// wantOwnTab fails unless the fake rex opened exactly one window for the
// session and never split a pane.
func wantOwnTab(t *testing.T, f *sessionFixture) {
	t.Helper()
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

// A check session opens its own tab even when an earlier session of the
// run has a pane, and splits nothing.
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
	wantOwnTab(t, f)
	rec, err := r.sessions.Read(sid)
	if err != nil || strings.TrimSpace(rec.Tab.Block) == "" || rec.Tab.Block == "block:a1" {
		t.Errorf("record %+v, %v", rec, err)
	}
}

// A change set fixer's session opens its own tab the same way.
func TestStartFix(t *testing.T) {
	f, r := paneRunner(t)
	f.rex.calls = nil
	step := workflow.Step{Name: "review", Split: workflow.SplitBelow, Fixer: &workflow.Role{Prompt: "fix"}}
	if _, err := r.StartFix(run.FixSpec{Step: step, Cycle: 2, Anchor: "S-first"}); err != nil {
		t.Fatal(err)
	}
	wantOwnTab(t, f)
}

// An agent step's session, such as a plan run's interview, opens its own
// tab and splits nothing.
func TestStartAgent(t *testing.T) {
	f, r := paneRunner(t)
	f.rex.calls = nil
	if _, err := r.StartAgent(workflow.Step{Name: "interview", Prompt: "go", Split: workflow.SplitRight}, "", nil); err != nil {
		t.Fatal(err)
	}
	wantOwnTab(t, f)
}

package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

func startLiveRun(t *testing.T, f *sessionFixture, cs string) {
	t.Helper()
	src := "[step.s]\nkind = \"leaves\"\n[step.s.builder]\nprompt = \"b\"\n[step.s.verifier]\nprompt = \"v\"\n"
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
	if err := run.Start(run.NewStore(run.Dir(f.common)), "R1", run.Started{Workflow: "wf", Definition: src, Change: cs, Worktree: f.worktree}, now, def); err != nil {
		t.Fatal(err)
	}
}

func wtLogFiles(t *testing.T, root, cs string) []string {
	t.Helper()
	es, err := os.ReadDir(filepath.Join(root, ".rhiza", "log", cs))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	var out []string
	for _, e := range es {
		out = append(out, e.Name())
	}
	return out
}

// An accept made in the main checkout while the change set's run is live
// also lands in the run's worktree, byte for byte.
func TestLiveRunWriteReachesWorktree(t *testing.T) {
	f := newSessionFixture(t, "live")
	cs := newChange(t, f.ctx, `{"intent":"Live"}`).ChangeSet.String()
	if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
		t.Fatal(err)
	}
	startLiveRun(t, f, cs)
	resp, err := callOp(t, f.ctx, "accept", `{"reason":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	name := resp.(LifecycleResponse).Entry.String() + ".json"
	a, err := os.ReadFile(filepath.Join(f.ctx.Root, ".rhiza", "log", cs, name))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(f.worktree, ".rhiza", "log", cs, name))
	if err != nil || string(a) != string(b) {
		t.Fatalf("worktree copy of the accept: %v, equal=%v", err, string(a) == string(b))
	}
	if got := wtLogFiles(t, f.worktree, cs); len(got) != 1 {
		t.Errorf("worktree holds %v, want only the accept", got)
	}
}

// A change set with no live run, and a dry run, copy nothing.
func TestLiveRunWriteNoRunNoCopy(t *testing.T) {
	f := newSessionFixture(t, "none")
	cs := newChange(t, f.ctx, `{"intent":"Quiet"}`).ChangeSet.String()
	if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := wtLogFiles(t, f.worktree, cs); len(got) != 0 {
		t.Errorf("no live run, worktree holds %v", got)
	}
	startLiveRun(t, f, cs)
	f.ctx.DryRun = true
	if _, err := callOp(t, f.ctx, "accept", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := wtLogFiles(t, f.worktree, cs); len(got) != 0 {
		t.Errorf("dry run, worktree holds %v", got)
	}
}

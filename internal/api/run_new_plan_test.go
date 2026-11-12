package api

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
)

// run start --new opens the change set in the run's worktree and copies the
// same bytes into the main checkout; a dry run predicts the same change set
// ID and worktree and writes nothing.
func TestNewPlanChange(t *testing.T) {
	f := newSessionFixture(t, "plan-reports")
	plan, err := os.ReadFile(filepath.Join("..", "workflow", "builtin", "plan.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(f.ctx.Root, ".rhiza", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.ctx.Root, ".rhiza", "workflows", "plan.toml"), plan, 0o644); err != nil {
		t.Fatal(err)
	}
	newChange(t, f.ctx, `{"intent":"Seed","title":"Seed"}`) // records the actor, so both runs draw alike
	req := `{"workflow":"plan","new":"Plan reports","permission_mode":"manual"}`

	// A real run is another process: it draws from the same seed.
	seed := func() *id.Generator { return id.NewGenerator(f.ctx.Now, rand.New(rand.NewSource(11))) }
	f.ctx.IDs = seed()
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "run.start", req)
	if err != nil {
		t.Fatal(err)
	}
	dry := resp.(RunResponse)
	if !dry.DryRun || dry.Run.Change == "" || dry.Run.Worktree != f.worktree {
		t.Fatalf("dry run = %+v, want a change set in worktree %s", dry.Run, f.worktree)
	}
	if _, err := os.Stat(filepath.Join(f.ctx.Root, ".rhiza", "log", dry.Run.Change)); !os.IsNotExist(err) {
		t.Fatalf("a dry run wrote the change set to main: %v", err)
	}

	f.ctx.IDs = seed()
	f.ctx.DryRun = false
	f.ctx.Start = (&shStart{}).start
	resp, err = callOp(t, f.ctx, "run.start", req)
	if err != nil {
		t.Fatal(err)
	}
	real := resp.(RunResponse)
	if real.Run.Change != dry.Run.Change || real.Run.Worktree != dry.Run.Worktree {
		t.Errorf("real run = %s in %s, dry run predicted %s in %s", real.Run.Change, real.Run.Worktree, dry.Run.Change, dry.Run.Worktree)
	}
	inWT := filepath.Join(f.worktree, ".rhiza", "log", real.Run.Change)
	inMain := filepath.Join(f.ctx.Root, ".rhiza", "log", real.Run.Change)
	wt, err := os.ReadDir(inWT)
	if err != nil || len(wt) != 1 {
		t.Fatalf("worktree log %s = %v, %v; want the open entry", inWT, wt, err)
	}
	a, _ := os.ReadFile(filepath.Join(inWT, wt[0].Name()))
	b, err := os.ReadFile(filepath.Join(inMain, wt[0].Name()))
	if err != nil || string(a) != string(b) {
		t.Errorf("main copy differs from the worktree's entry: %v", err)
	}
}

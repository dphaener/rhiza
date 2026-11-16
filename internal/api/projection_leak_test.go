package api

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// throwawayProjections makes ctx open throwaway indexes (no index path,
// no user cache folder for another root's) under a fresh TMPDIR, and
// returns a check that none is left there.
func throwawayProjections(t *testing.T, ctx *Ctx) func() {
	t.Helper()
	tmp := t.TempDir()
	ctx.IndexPath = ""
	inner := ctx.Getenv
	ctx.Getenv = func(k string) string {
		switch k {
		case "TMPDIR":
			return tmp
		case "XDG_CACHE_HOME", "HOME", "LocalAppData":
			return ""
		}
		if inner == nil {
			return ""
		}
		return inner(k)
	}
	return func() {
		t.Helper()
		left, err := filepath.Glob(filepath.Join(tmp, "rhiza-projection-*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(left) > 0 {
			t.Errorf("projections left open: %v", left)
		}
	}
}

func TestShipCheckModelConflictClosesProjection(t *testing.T) {
	t.Parallel()
	f, _, _ := modelConflictFixture(t)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err == nil {
		t.Fatal("land succeeded")
	}
	check := throwawayProjections(t, f.ctx)
	if r := shipCheck(t, f, `{"model_conflict":true}`); r.ModelConflict == nil || r.ModelConflict.Intent != "Ship it" {
		t.Fatalf("check = %+v", r)
	}
	check()
}

func TestRunStatusGatedClosesProjection(t *testing.T) {
	t.Parallel()
	f := newMessageFixture(t)
	cs := currentOf(t, f.ctx).String()
	rs := run.NewStore(run.Dir(f.common))
	t0 := time.Date(2026, 11, 13, 12, 0, 0, 0, time.UTC)
	if err := rs.Create("R1", run.Entry{Type: run.TypeStarted, At: t0, Started: &run.Started{Workflow: "w", Change: cs}}); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append("R1", run.Entry{Type: run.TypeGated, At: t0.Add(time.Minute), Changes: []string{cs}}); err != nil {
		t.Fatal(err)
	}
	check := throwawayProjections(t, f.ctx)
	if _, err := callOp(t, f.ctx, "run.status", `{"run":"R1"}`); err != nil {
		t.Fatal(err)
	}
	check()
}

func TestWorktreeWorkspaceClosesReplacedProjection(t *testing.T) {
	t.Parallel()
	main := newMessageFixture(t)
	other := newSessionFixture(t, "other")
	cs := currentOf(t, main.ctx)
	check := throwawayProjections(t, main.ctx)
	w, err := openCtxWorkspace(main.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer w.closeProjection()
	// The other root lacks cs, so its first workspace is replaced after
	// the copy; that first one's projection must be closed.
	ww, err := worktreeWorkspace(main.ctx, w, other.ctx.Root, cs)
	if err != nil {
		t.Fatal(err)
	}
	ww.closeProjection()
	w.closeProjection()
	check()

	// A change set the main log lacks fails the copy; the open workspace
	// is closed on that error too.
	if _, err := worktreeWorkspace(main.ctx, w, other.ctx.Root, newID(t, main.ctx)); err == nil {
		t.Fatal("unknown change set opened")
	}
	w.closeProjection()
	check()
}

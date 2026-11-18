package api

import (
	"slices"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// The daemon's own resume and the watchdog's retry journal no
// intervention: the watchdog queues its retry through queueContinue, which
// journals nothing, and every run resume --foreground the daemon detaches
// (its own resume, a queued continuation) carries envDaemonResume, so the
// resumed process journals none either.
func TestInterventionDaemonResumesJournalNothing(t *testing.T) {
	t.Parallel()
	ctx := testCtx(t, t.TempDir())
	common := t.TempDir()
	ctx.Git = func(args ...string) (string, error) { return common, nil }
	ctx.Self = "rhiza"
	ctx.repoSession = true
	var detached []Program
	ctx.Detach = func(p Program, _ string) error { detached = append(detached, p); return nil }

	rs := runStore(ctx, common)
	at := time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)
	for _, id := range []string{"R1", "R2"} {
		if err := rs.Create(id, run.Entry{
			Type: run.TypeStarted, At: at,
			Started: &run.Started{Workflow: "hold", Change: "01M76F4G0JQ649RR29M9QPFX68", Title: "Hold"},
		}); err != nil {
			t.Fatal(err)
		}
		if err := rs.Append(id, run.Entry{Type: run.TypeEntered, At: at, Step: "hold", Visit: 1}); err != nil {
			t.Fatal(err)
		}
	}

	// The daemon's own resume of R1.
	if err := daemonResume(ctx, "R1"); err != nil {
		t.Fatal(err)
	}
	// The watchdog's retry of R2, queued, then driven by the daemon.
	if _, err := queueContinue(ctx, common, "R2", true, false); err != nil {
		t.Fatal(err)
	}
	runs, err := rs.List()
	if err != nil {
		t.Fatal(err)
	}
	continued, err := daemonQueuedContinues(ctx, rs, runs, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(continued, []string{"R2"}) {
		t.Fatalf("continued %v, want [R2]", continued)
	}

	if len(detached) != 2 {
		t.Fatalf("detached %d programs, want 2", len(detached))
	}
	for _, p := range detached {
		if !slices.Contains(p.Env, envDaemonResume+"=1") {
			t.Errorf("%v: env %v lacks %s=1", p.Args, p.Env, envDaemonResume)
		}
	}
	for _, id := range []string{"R1", "R2"} {
		j, err := rs.Read(id)
		if err != nil {
			t.Fatal(err)
		}
		if n := len(j.State().Interventions); n != 0 {
			t.Errorf("run %s: %d interventions, want 0", id, n)
		}
	}
}

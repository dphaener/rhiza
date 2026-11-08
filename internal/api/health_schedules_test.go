package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

func schedRuns(t *testing.T, ctx *Ctx, rid string, ended bool, result string, warn string) {
	t.Helper()
	common, err := gitCommonDir(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(common))
	at := time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC)
	if err := rs.Create(rid, run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "schedule:tick", Trigger: run.TriggerSchedule, TriggerName: "tick"}}); err != nil {
		t.Fatal(err)
	}
	if ended {
		if err := rs.Append(rid, run.Entry{Type: run.TypeFinished, At: at.Add(time.Minute), End: result, Detail: "boom"}); err != nil {
			t.Fatal(err)
		}
	}
	if warn != "" {
		if err := rs.Append(rid, run.Entry{Type: run.TypeWarning, At: at.Add(2 * time.Minute), Detail: warn}); err != nil {
			t.Fatal(err)
		}
	}
}

// schedHealthCtx is a project whose rhiza check passes (its rendered
// folder written), so no other part of health is bad and the exit code is
// health's own over the schedules part. The disk part reads the pinned
// 1 TiB, never the host's free space.
func schedHealthCtx(t *testing.T, body string) *Ctx {
	t.Helper()
	f := newSessionFixture(t, "s")
	pinFreeBytes(f.ctx)
	f.ctx.Now = func() time.Time { return time.Date(2026, 11, 8, 10, 30, 0, 0, time.UTC) }
	if body != "" {
		writeDispatch(t, f.ctx.Root, body)
	}
	if _, err := callOp(t, f.ctx, "render", `{}`); err != nil {
		t.Fatal(err)
	}
	return f.ctx
}

// schedHealth runs the health op, requires it to exit 0 with no bad part,
// and returns its schedules part, false when it lists none.
func schedHealth(t *testing.T, ctx *Ctx) (HealthPart, bool) {
	t.Helper()
	resp, err := callOp(t, ctx, "health", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(HealthResponse)
	if r.ExitCode() != 0 {
		t.Errorf("health exits %d (state %s): %+v", r.ExitCode(), r.State, r.Parts)
	}
	for _, p := range r.Parts {
		if p.Part == HealthSchedules {
			if p.State == HealthBad {
				t.Errorf("schedules part bad: %+v", p)
			}
			return p, true
		}
	}
	return HealthPart{}, false
}

func TestHealthSchedulesPart(t *testing.T) {
	t.Run("no schedules, no part", func(t *testing.T) {
		ctx := schedHealthCtx(t, "")
		if _, ok := healthSchedules(ctx, true); ok {
			t.Fatal("part present")
		}
		if _, ok := schedHealth(t, ctx); ok {
			t.Fatal("health lists a schedules part")
		}
	})
	t.Run("not yet run is ok with its next run", func(t *testing.T) {
		ctx := schedHealthCtx(t, hourlyCmd)
		p, ok := healthSchedules(ctx, true)
		if !ok || p.State != HealthOK || !strings.Contains(p.Summary, "1 schedule") {
			t.Fatalf("got %+v", p)
		}
		if len(p.Detail) != 1 || !strings.Contains(p.Detail[0], "none yet") || !strings.Contains(p.Detail[0], "next run 2026-11-08T11:00:00Z") {
			t.Errorf("detail %q", p.Detail)
		}
		// Through the health op no daemon runs, so the part warns, but
		// health still exits 0 and the line is the same.
		hp, ok := schedHealth(t, ctx)
		if !ok || len(hp.Detail) != 1 || hp.Detail[0] != p.Detail[0] {
			t.Errorf("health part %+v", hp)
		}
	})
	t.Run("failed last run warns with its log", func(t *testing.T) {
		ctx := schedHealthCtx(t, hourlyCmd)
		schedRuns(t, ctx, "01M6WSQMY5Q33CWG319V9ARAAA", true, run.Failed, "")
		p, _ := healthSchedules(ctx, true)
		if p.State != HealthWarn || p.Next != "rhiza run log 01M6WSQMY5Q33CWG319V9ARAAA" || !strings.Contains(p.Detail[0], "failed") {
			t.Fatalf("got %+v", p)
		}
		hp, ok := schedHealth(t, ctx)
		if !ok || hp.State != HealthWarn || hp.Next != p.Next {
			t.Errorf("health part %+v", hp)
		}
	})
	t.Run("a run that changed the checkout warns", func(t *testing.T) {
		ctx := schedHealthCtx(t, hourlyCmd)
		schedRuns(t, ctx, "01M6WSQMY5Q33CWG319V9ARBBB", true, run.Succeeded, "main checkout changed: HEAD moved")
		p, _ := healthSchedules(ctx, true)
		if p.State != HealthWarn || !strings.Contains(p.Summary, "main checkout") || !strings.Contains(p.Detail[0], "HEAD moved") {
			t.Fatalf("got %+v", p)
		}
		hp, ok := schedHealth(t, ctx)
		if !ok || hp.State != HealthWarn || !strings.Contains(hp.Summary, "main checkout") {
			t.Errorf("health part %+v", hp)
		}
	})
	t.Run("no daemon warns, and health exits 0", func(t *testing.T) {
		ctx := schedHealthCtx(t, hourlyCmd)
		p, _ := healthSchedules(ctx, false)
		if p.State != HealthWarn || !strings.Contains(p.Summary, "no daemon") {
			t.Fatalf("got %+v", p)
		}
		hp, ok := schedHealth(t, ctx)
		if !ok || hp.State != HealthWarn || !strings.Contains(hp.Summary, "no daemon") {
			t.Errorf("health part %+v", hp)
		}
	})
}

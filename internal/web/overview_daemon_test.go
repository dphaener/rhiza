package web_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/web"
)

// The Overview at / reads the daemon from the repository: a live daemon
// (its lock held) whose recorded state says paused, with two of its three
// runs running against a limit of four, and every health part of the
// repository's health report as a pill with its state.
func TestOverviewReadsDaemonFromRepo(t *testing.T) {
	ctx := fixtureCtx(t)
	common := filepath.Join(ctx.Root, ".git")
	lock, err := daemon.TryLock(common, ctx.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(lock.Unlock)
	if err := daemon.WriteState(common, daemon.State{
		Process:  daemon.ReadHolder(common),
		LastPoll: ctx.Now(),
		MaxRuns:  4,
		Settings: []daemon.SettingValue{{Key: "rhiza.daemon.paused", Value: "true", Source: "git config"}},
		Runs: []daemon.LiveRun{
			{ID: "R1", Change: "C1", State: "running"},
			{ID: "R2", Change: "C2", State: "interrupted"},
			{ID: "R3", Change: "C3", State: "running"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	o, _ := api.Lookup("health")
	out, err := o.Run(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	parts := out.(api.HealthResponse).Parts
	if len(parts) == 0 {
		t.Fatal("health has no parts")
	}

	body := serve(web.HandlerWith(api.WebDaemonSource(ctx), nil), "GET", "/").Body.String()
	wants := []string{
		"Paused, 2 of 4 slots in use",
		`<div class="rz-stat-val">Paused</div>`,
		`<div class="rz-stat-val">2<small>/ 4</small></div>`,
		`<div class="rz-slots" aria-hidden="true"><i class="on"></i><i class="on"></i><i></i><i></i></div>`,
		"daemon · paused · 2/4",
		`rz-pill--health rz-pill--ok" title="daemon: ok`,
	}
	for _, p := range parts {
		wants = append(wants, `rz-pill--health rz-pill--`+p.State+`" title="`+p.Part+`: `+p.State)
	}
	for _, w := range wants {
		if !strings.Contains(body, w) {
			t.Errorf("/ lacks %q:\n%s", w, body)
		}
	}
	if n := strings.Count(body, "rz-pill--health "); n != len(parts) {
		t.Errorf("%d health pills, want %d", n, len(parts))
	}
}

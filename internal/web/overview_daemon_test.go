package web_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/web"
)

// The Overview at / reads the daemon from the repository: a live daemon
// (its lock held) whose recorded state says paused, with two of its three
// runs running, one start not journaled against a limit of four, and every health part of the
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
		Starting: 1,
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

	h := web.HandlerWith(api.WebDaemonSource(ctx, api.NewWebHealth(ctx)), nil)
	body := serve(h, "GET", "/").Body.String()
	for end := time.Now().Add(30 * time.Second); strings.Contains(body, "rz-pill--pending") && time.Now().Before(end); {
		time.Sleep(50 * time.Millisecond)
		body = serve(h, "GET", "/").Body.String()
	}
	wants := []string{
		"Paused, 3 of 4 slots in use",
		`<div class="rz-stat-val">Paused</div>`,
		`<div class="rz-stat-val">3<small>/ 4</small></div>`,
		`<div class="rz-slots" aria-hidden="true"><i class="on"></i><i class="on"></i><i class="on"></i><i></i></div>`,
		"runs: 2/4 live, 1 starting",
		"daemon · paused · 3/4",
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

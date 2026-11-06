package web

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func overviewBody(t *testing.T, src DaemonSource, path string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	HandlerWith(src, nil).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
	if rec.Code != 200 {
		t.Fatalf("%s: %d", path, rec.Code)
	}
	return rec.Body.String()
}

func TestOverviewShowsDaemon(t *testing.T) {
	src := func() DaemonView {
		return DaemonView{Running: true, Paused: true, Slots: 2, Limit: 4, Health: []HealthPill{
			{Part: "daemon", State: "ok", Summary: "running"},
			{Part: "disk", State: "warn", Summary: "81% used"},
			{Part: "check", State: "bad", Summary: "2 problems"},
		}}
	}
	for _, path := range []string{"/", FragmentURL("daemon")} {
		body := overviewBody(t, src, path)
		for _, want := range []string{
			"Paused, 2 of 4 slots in use",
			`<div class="rz-stat-val">Paused</div>`,
			`<div class="rz-stat-val">2<small>/ 4</small></div>`,
			`<div class="rz-slots" aria-hidden="true"><i class="on"></i><i class="on"></i><i></i><i></i></div>`,
			`rz-pill--health rz-pill--ok" title="daemon: ok · running"`,
			`rz-pill--health rz-pill--warn" title="disk: warn · 81% used"`,
			`rz-pill--health rz-pill--bad" title="check: bad · 2 problems"`,
		} {
			if !strings.Contains(body, want) {
				t.Errorf("%s lacks %q:\n%s", path, want, body)
			}
		}
	}
	page := overviewBody(t, src, "/")
	if !strings.Contains(page, `data-rz-fragment="/fragments/daemon"`) || !strings.Contains(page, "daemon · paused · 2/4") {
		t.Errorf("page lacks the live daemon section or header pill")
	}
}

func TestOverviewDaemonRunningAndStopped(t *testing.T) {
	run := overviewBody(t, func() DaemonView { return DaemonView{Running: true, Slots: 1, Limit: 2} }, FragmentURL("daemon"))
	if !strings.Contains(run, "Running, 1 of 2 slots in use") || !strings.Contains(run, "rz-live-dot") {
		t.Errorf("running: %s", run)
	}
	off := overviewBody(t, nil, FragmentURL("daemon"))
	if !strings.Contains(off, "Not running") || strings.Contains(off, "rz-live-dot") {
		t.Errorf("stopped: %s", off)
	}
}

func TestOverviewListsLiveRuns(t *testing.T) {
	src := func() []RunView {
		return []RunView{
			{ID: "01RUN1", Title: "Retry webhook delivery", Workflow: "build", Step: "leaves", State: "running", Started: "2026-11-06T10:00:00Z"},
			{ID: "01RUN2", Title: "Idempotency keys", Workflow: "plan", Step: "gate", State: "waiting"},
		}
	}
	for _, path := range []string{"/", FragmentURL("runs")} {
		rec := httptest.NewRecorder()
		HandlerWithRuns(nil, src, nil).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		body := rec.Body.String()
		for _, want := range []string{
			`href="/runs/01RUN1"`, "Retry webhook delivery", "<small>build</small>", ">leaves</span>",
			`href="/runs/01RUN2"`, "Idempotency keys", "<small>plan</small>", ">gate</span>", "waiting on you",
		} {
			if !strings.Contains(body, want) {
				t.Errorf("%s lacks %q:\n%s", path, want, body)
			}
		}
	}
	rec := httptest.NewRecorder()
	HandlerWithRuns(nil, func() []RunView { return nil }, nil).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if b := rec.Body.String(); !strings.Contains(b, "No live runs") || !strings.Contains(b, `data-rz-fragment="/fragments/runs"`) {
		t.Errorf("empty: %s", b)
	}
}

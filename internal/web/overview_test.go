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

func TestOverviewListsInbox(t *testing.T) {
	src := func() []InboxItemView {
		return []InboxItemView{
			{Kind: "failed", Title: "Idempotency keys", Detail: "build failed at leaves", Run: "01RUN3", Step: "leaves", At: "2026-11-06T09:00:00Z"},
			{Kind: "question", Title: "Which queue should retries use?", Detail: "Retry webhook delivery", Run: "01RUN1", At: "2026-11-06T07:00:00Z"},
			{Kind: "signoff", Title: "Rate limits", Run: "01RUN2", Step: "signoff", At: "2026-11-06T08:00:00Z"},
		}
	}
	for _, path := range []string{"/", FragmentURL("inbox")} {
		rec := httptest.NewRecorder()
		HandlerWithInbox(nil, nil, src, nil).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		body := rec.Body.String()
		q, s, f := strings.Index(body, "Which queue should retries use?"), strings.Index(body, "Rate limits"), strings.Index(body, "Idempotency keys")
		if q < 0 || s < 0 || f < 0 || !(q < s && s < f) {
			t.Errorf("%s: items missing or not oldest first (%d %d %d):\n%s", path, q, s, f, body)
		}
		for _, want := range []string{
			"3 waiting on you", "rz-panel--waiting", `href="/runs/01RUN2"`,
			"rz-pill--waiting", ">question</span>", "sign-off</span>", "failed</span>",
		} {
			if !strings.Contains(body, want) {
				t.Errorf("%s lacks %q:\n%s", path, want, body)
			}
		}
		if strings.Contains(body, "01RUN2</") {
			t.Errorf("%s shows an ID", path)
		}
	}
	rec := httptest.NewRecorder()
	HandlerWithInbox(nil, nil, func() []InboxItemView { return nil }, nil).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if b := rec.Body.String(); !strings.Contains(b, "Nothing waits on you") || strings.Contains(b, "rz-panel--waiting") || !strings.Contains(b, `data-rz-fragment="/fragments/inbox"`) {
		t.Errorf("empty: %s", b)
	}
}

func TestOverviewShowsMilestoneProgress(t *testing.T) {
	src := func() []MilestoneView {
		return []MilestoneView{
			{Name: "Durable delivery", Total: 5, Counts: map[string]int{"accepted": 3, "proposed": 2}},
			{Name: "Single change set", Total: 1, Counts: map[string]int{"proposed": 1}},
		}
	}
	for _, path := range []string{"/", FragmentURL("milestones")} {
		rec := httptest.NewRecorder()
		HandlerWithMilestones(nil, nil, nil, src, nil).ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		body := rec.Body.String()
		for _, want := range []string{
			"2 open", "Durable delivery", "<small>5 change sets</small>", "<small>1 change set</small>",
			"rz-pill--accepted", "rz-pill--proposed", `<span class="rz-num">3</span>`, `<span class="rz-num">2</span>`,
			`href="/roadmap"`,
		} {
			if !strings.Contains(body, want) {
				t.Errorf("%s lacks %q:\n%s", path, want, body)
			}
		}
	}
	rec := httptest.NewRecorder()
	HandlerWithMilestones(nil, nil, nil, func() []MilestoneView { return nil }, nil).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if b := rec.Body.String(); !strings.Contains(b, "No open milestones") || !strings.Contains(b, `data-rz-fragment="/fragments/milestones"`) {
		t.Errorf("empty: %s", b)
	}
}

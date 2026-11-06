package web

import (
	"net/http"
	"strings"
	"testing"
)

const fixturePlanRun = `{"title":"Plan the thing","workflow":"plan","state":"running","current":"gate","started_at":"2026-10-24T10:00:00Z",
"steps":[
{"step":"interview","outcome":"pass","at":"2026-10-24T10:05:00Z"},
{"step":"draft","outcome":"fail","detail":"exit 2","at":"2026-10-24T10:06:30Z"}],
"sessions":[{"step":"interview","session":"01SESSIONA","tokens":{"session":{"input":10,"output":5}}},{"step":"gate","session":"01SESSIONB","no_tokens":"still running"}]}`

const fixtureReviewRun = `{"title":"Review the thing","workflow":"review","state":"succeeded","started_at":"2026-10-24T10:00:00Z",
"steps":[{"step":"checks","outcome":"pass","at":"2026-10-24T10:20:00Z"}],
"checks":[{"check":"correctness","outcome":"failed","rounds":[{"findings":[
{"severity":"blocker","place":"a.go:3","reason":"Nil dereference","verdict":"confirmed"},
{"severity":"nit","reason":"Spelling"}]}]}]}`

func runAPI() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/run.status" {
			http.NotFound(w, r)
			return
		}
		switch r.URL.Query().Get("run") {
		case "PLAN":
			w.Write([]byte(fixturePlanRun))
		case "REVIEW":
			w.Write([]byte(fixtureReviewRun))
		default:
			http.NotFound(w, r)
		}
	})
}

func inOrder(t *testing.T, body string, parts ...string) {
	t.Helper()
	last := -1
	for _, p := range parts {
		i := strings.Index(body, p)
		if i < 0 {
			t.Errorf("missing %q", p)
			continue
		}
		if i < last {
			t.Errorf("%q is out of order", p)
		}
		last = i
	}
}

func TestRunPagePlanRun(t *testing.T) {
	rec := get(Handler(runAPI()), "/runs/PLAN")
	if rec.Code != 200 {
		t.Fatalf("%d %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	inOrder(t, body, "Plan the thing", ">interview<", ">draft<", ">gate<")
	for _, want := range []string{
		"rz-pill--done", "rz-pill--failed", "5m 0s", "1m 30s", "exit 2", "Oct 24 10:05",
		"rz-pill--running", "Session 1", "finished · 15 tokens", "running",
		`data-rz-fragment="/fragments/run-detail?run=PLAN"`, `data-rz-events="run"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("lacks %q", want)
		}
	}
	if strings.Contains(body, "01SESSION") || strings.Contains(body, "Review</h2>") {
		t.Error("shows a session ID or a review section")
	}
	frag := get(Handler(runAPI()), "/fragments/run-detail?run=PLAN")
	if frag.Code != 200 || !strings.Contains(frag.Body.String(), ">interview<") {
		t.Errorf("fragment: %d %s", frag.Code, frag.Body)
	}
}

func TestRunPageReviewRun(t *testing.T) {
	body := get(Handler(runAPI()), "/runs/REVIEW").Body.String()
	inOrder(t, body, ">checks<", "Review</h2>", "correctness", ">blocker<", "Nil dereference", "a.go:3", ">nit<", "Spelling")
	if !strings.Contains(body, "rz-pill--failed") {
		t.Error("check outcome pill missing")
	}
}

func TestRunPageUnknownIs404(t *testing.T) {
	if rec := get(Handler(runAPI()), "/runs/NOPE"); rec.Code != 404 {
		t.Errorf("unknown run: %d", rec.Code)
	}
}

package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fixtureRuns serves run.ls over fixture journals: one running, one waiting
// on the owner, one done, one failed, newest started first.
const fixtureRuns = `{"runs":[
{"id":"01DONE000","title":"Tidy the done run","state":"succeeded","started_at":"2026-10-24T11:00:00Z","live":false,"workflow":"ship"},
{"id":"01LIVE000","title":"Build the live run","state":"running","current":"build","started_at":"2026-10-24T10:00:00Z","live":true,"workflow":"build"},
{"id":"01FAIL000","title":"Break the failed run","state":"failed","current":"review","started_at":"2026-10-24T09:00:00Z","live":false,"workflow":"review"},
{"id":"01WAIT000","title":"Gate the waiting run","state":"blocked","current":"gate","started_at":"2026-10-24T08:00:00Z","live":true,"workflow":"plan"}
]}`

func fakeAPI(body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/run.ls" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(body))
	})
}

func get(h http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", target, nil))
	return rec
}

func TestRunsPageOrderFieldsAndLinks(t *testing.T) {
	rec := get(Handler(fakeAPI(fixtureRuns)), "/runs")
	if rec.Code != 200 {
		t.Fatalf("/runs: %d %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	// live first (in run.ls order), then recent newest first
	order := []string{"Build the live run", "Gate the waiting run", "Tidy the done run", "Break the failed run"}
	last := -1
	for _, title := range order {
		i := strings.Index(body, title)
		if i < 0 {
			t.Fatalf("missing %q", title)
		}
		if i < last {
			t.Errorf("%q is out of order", title)
		}
		last = i
	}
	if strings.Index(body, ">Live<") > strings.Index(body, order[0]) || strings.Index(body, ">Recent<") < strings.Index(body, order[1]) ||
		strings.Index(body, ">Recent<") > strings.Index(body, order[2]) {
		t.Error("Live and Recent group labels misplaced")
	}
	for _, want := range []string{
		`href="/runs/01LIVE000"`, `href="/runs/01WAIT000"`, `href="/runs/01DONE000"`, `href="/runs/01FAIL000"`,
		"rz-pill--running", "rz-pill--waiting", "waiting on you", "rz-pill--done", "rz-pill--failed",
		"<small>build</small>", "<small>review</small>", ">gate<", ">build<",
		`<time datetime="2026-10-24T10:00:00Z" data-rel>Oct 24 10:00</time>`,
		`data-rz-fragment="/fragments/runs"`, `data-rz-events="run"`, `aria-current="page"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("/runs lacks %q", want)
		}
	}
	// constraint owner-surfaces-no-ids: IDs are in links only, never text
	for _, id := range []string{"01LIVE000", "01DONE000", "01FAIL000", "01WAIT000"} {
		if strings.Count(body, id) != 1 {
			t.Errorf("run ID %s shown outside its link", id)
		}
	}
}

func TestRunsFragmentMatchesSection(t *testing.T) {
	h := Handler(fakeAPI(fixtureRuns))
	frag := get(h, "/fragments/runs")
	if frag.Code != 200 || !strings.Contains(frag.Body.String(), "Build the live run") {
		t.Fatalf("fragment: %d %s", frag.Code, frag.Body)
	}
	if !strings.Contains(get(h, "/runs").Body.String(), frag.Body.String()) {
		t.Error("the page does not embed the fragment's content")
	}
	if f := fragments["runs"]; len(f.Events) != 1 || f.Events[0] != EventRun {
		t.Errorf("runs section events %v", f.Events)
	}
}

func TestRunsPageEmptyAndUnavailable(t *testing.T) {
	if b := get(Handler(fakeAPI(`{"runs":[]}`)), "/runs").Body.String(); !strings.Contains(b, "No runs yet") {
		t.Errorf("empty: %s", b)
	}
	if rec := get(Handler(nil), "/runs"); rec.Code != 200 || !strings.Contains(rec.Body.String(), "Runs are unavailable") {
		t.Errorf("no api: %d %s", rec.Code, rec.Body)
	}
}

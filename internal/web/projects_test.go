package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/web"
)

func TestProjectsHomeListsEveryProject(t *testing.T) {
	src := func() []web.ProjectView {
		return []web.ProjectView{
			{Name: "alpha", Initialised: true, Running: 2, Waiting: 1, Inbox: 3, Health: []web.HealthPill{{Part: "git", State: "ok"}}},
			{Name: "beta", Initialised: true, Inbox: 0, Health: []web.HealthPill{{Part: "git", State: "ok"}, {Part: "disk", State: "bad"}}},
			{Name: "gamma"},
		}
	}
	h := web.RootWithProjects(web.Handler(nil), func(string) http.Handler { return nil }, src)
	rec := serve(h, "GET", "/projects")
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`href="/p/alpha/"`, "2 running, 1 waiting on you · 3 in inbox", "healthy",
		`href="/p/beta/"`, "none live · 0 in inbox", "unhealthy", "rz-pill--bad",
		`href="/p/gamma/"`, "not initialised",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "gamma</a>") && strings.Contains(body, "none live · 0 in inbox</small></a><span class=\"rz-pill-cell\"><span class=\"rz-pill rz-pill--draft\"") {
		t.Errorf("bare project shown with counts")
	}
}

func TestProjectsHomeEmpty(t *testing.T) {
	rec := serve(web.RootWithProjects(web.Handler(nil), nil, nil), "GET", "/projects")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "No projects") {
		t.Errorf("%d %s", rec.Code, rec.Body)
	}
}

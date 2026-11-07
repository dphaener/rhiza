package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/web"
)

func twoProjects(t *testing.T) http.Handler {
	t.Helper()
	a, b := fixtureCtx(t), fixtureCtx(t)
	plan(t, a, "Alpha only work", "", true)
	plan(t, b, "Beta only work", "", true)
	hs := map[string]http.Handler{"a": web.Handler(api.ReadOnlyHandler(a), web.Events(nil, web.PollEvery)), "b": web.Handler(api.ReadOnlyHandler(b))}
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "a", Initialised: true}, {Name: "b", Initialised: true}, {Name: "c"}}
	}
	return web.RootWithProjects(web.Handler(nil), func(n string) http.Handler { return hs[n] }, list)
}

func TestProjectPagesAreScopedWithSwitcher(t *testing.T) {
	h := twoProjects(t)
	for _, c := range []struct{ project, other, page, otherPage string }{
		{"a", "b", "", "/"},
		{"a", "b", "runs", "/runs"},
		{"b", "a", "roadmap", "/roadmap"},
		{"b", "a", "runs", "/runs"},
	} {
		rec := serve(h, "GET", "/p/"+c.project+"/"+c.page)
		if rec.Code != 200 {
			t.Fatalf("%s/%s: %d", c.project, c.page, rec.Code)
		}
		body := rec.Body.String()
		if want := `href="/p/` + c.other + c.otherPage + `"`; !strings.Contains(body, want) {
			t.Errorf("/p/%s/%s lacks switcher link %s", c.project, c.page, want)
		}
		if strings.Contains(body, `href="/p/c`) {
			t.Errorf("uninitialised project in switcher")
		}
		for _, bad := range []string{`href="/runs"`, `href="/roadmap"`, `src="/static/`, `data-rz-fragment="/fragments`} {
			if strings.Contains(body, bad) {
				t.Errorf("/p/%s/%s has unscoped %s", c.project, c.page, bad)
			}
		}
		if !strings.Contains(body, `data-rz-base="/p/`+c.project+`"`) {
			t.Errorf("no base on /p/%s/%s", c.project, c.page)
		}
	}
	// each project's roadmap shows only its own change set
	if body := serve(h, "GET", "/p/a/roadmap").Body.String(); !strings.Contains(body, "Alpha only work") || strings.Contains(body, "Beta only work") {
		t.Errorf("project a roadmap:\n%s", body)
	}
	if body := serve(h, "GET", "/p/b/roadmap").Body.String(); !strings.Contains(body, "Beta only work") || strings.Contains(body, "Alpha only work") {
		t.Errorf("project b roadmap:\n%s", body)
	}
	if rec := serve(h, "GET", "/p/c/runs"); rec.Code != 404 {
		t.Errorf("uninitialised: %d", rec.Code)
	}
	if rec := serve(h, "GET", "/p/nope/runs"); rec.Code != 404 {
		t.Errorf("unknown: %d", rec.Code)
	}
	// fragments are scoped too
	if body := serve(h, "GET", "/p/a/fragments/inbox").Body.String(); strings.Contains(body, `href="/runs`) || strings.Contains(body, `href="/changes`) {
		t.Errorf("fragment unscoped: %s", body)
	}
}

func TestSwitcherKeepsListKindOnDetailPages(t *testing.T) {
	h := twoProjects(t)
	rec := serve(h, "GET", "/p/a/runs/nope")
	if strings.Contains(rec.Body.String(), `href="/p/b/runs/nope"`) {
		t.Errorf("detail id carried to the other project")
	}
}

func TestEventStreamScopedToProject(t *testing.T) {
	h := twoProjects(t)
	if rec := serve(h, "GET", "/p/b/events"); rec.Code == 200 {
		t.Errorf("project b has no stream, got %d", rec.Code)
	}
}

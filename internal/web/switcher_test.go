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
		if strings.Contains(body, `href="/p/`+c.other) {
			t.Errorf("/p/%s/%s lists project %s in the page itself", c.project, c.page, c.other)
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

func TestEventStreamScopedToProject(t *testing.T) {
	h := twoProjects(t)
	if rec := serve(h, "GET", "/p/b/events"); rec.Code == 200 {
		t.Errorf("project b has no stream, got %d", rec.Code)
	}
}

// The picker's list is not in the page: assets, fragments and page renders
// under /p/<project>/ never read the projects source; one GET /picker reads
// it exactly once.
func TestSwitcherNamesReadOnlyForPages(t *testing.T) {
	a := fixtureCtx(t)
	hs := map[string]http.Handler{"a": web.Handler(api.ReadOnlyHandler(a))}
	projects := 0
	list := func() []web.ProjectView { projects++; return []web.ProjectView{{Name: "a", Initialised: true}} }
	h := web.RootWithInbox(web.Handler(nil), "", func(n string) http.Handler { return hs[n] }, list, nil)
	for _, target := range []string{"/p/a/static/css/rhiza.css", "/p/a/fragments/runs", "/p/a/runs", "/p/a/"} {
		serve(h, "GET", target)
	}
	if projects != 0 {
		t.Errorf("assets, fragments and pages read projects %d times, want none", projects)
	}
	if rec := serve(h, "GET", "/picker?path=/p/a/runs"); rec.Code != 200 || projects != 1 {
		t.Errorf("GET /picker: %d, read projects %d times, want 200 and once", rec.Code, projects)
	}
}

func TestProjectPicker(t *testing.T) {
	list := func(n int) web.ProjectsSource {
		return func() []web.ProjectView {
			vs := []web.ProjectView{
				{Name: "a", Initialised: true},
				{Name: "b", Initialised: true, Running: 2, NeedsYou: 3},
				{Name: "c"},
			}
			for i := 3; i < n; i++ {
				vs = append(vs, web.ProjectView{Name: string(rune('a' + i)), Initialised: true})
			}
			return vs
		}
	}
	h := web.RootWithProjects(web.Handler(nil), func(string) http.Handler { return nil }, list(2))
	rec := serve(h, "GET", "/picker?path=/p/a/roadmap")
	if rec.Code != 200 {
		t.Fatalf("picker: %d", rec.Code)
	}
	body := rec.Body.String()
	for _, w := range []string{`href="/p/a/roadmap"`, `href="/p/b/roadmap"`, `href="/inbox"`, "2 running", "idle", "3 need you", "All projects"} {
		if !strings.Contains(body, w) {
			t.Errorf("picker lacks %q:\n%s", w, body)
		}
	}
	if strings.Contains(body, `href="/p/c`) {
		t.Errorf("uninitialised project listed")
	}
	if strings.Contains(body, "type=\"search\"") {
		t.Errorf("search with two projects")
	}
	if strings.Index(body, `href="/p/b/`) > strings.Index(body, `href="/p/a/`) {
		t.Errorf("project with something waiting not first")
	}
	// a run's own page keeps its list
	if body := serve(h, "GET", "/picker?path=/p/a/runs/xyz").Body.String(); !strings.Contains(body, `href="/p/b/runs"`) {
		t.Errorf("detail page not kept as its list:\n%s", body)
	}
	// five projects: no search; six: search
	five := serve(web.RootWithProjects(web.Handler(nil), func(string) http.Handler { return nil }, list(6)), "GET", "/picker").Body.String()
	six := serve(web.RootWithProjects(web.Handler(nil), func(string) http.Handler { return nil }, list(7)), "GET", "/picker").Body.String()
	if strings.Contains(five, `type="search"`) || !strings.Contains(six, `type="search"`) {
		t.Errorf("search field: with five %v, with six %v", strings.Contains(five, `type="search"`), strings.Contains(six, `type="search"`))
	}
}

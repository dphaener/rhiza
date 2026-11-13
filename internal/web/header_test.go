package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/web"
)

// headerRoot is a two-project root whose project a has n things waiting.
func headerRoot(t *testing.T, n int) http.Handler {
	t.Helper()
	a, b := fixtureCtx(t), fixtureCtx(t)
	hs := map[string]http.Handler{
		"a": web.HandlerWithNeedsYou(nil, nil, nil, nil, func() int { return n }, api.ReadOnlyHandler(a)),
		"b": web.Handler(api.ReadOnlyHandler(b)),
	}
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "a", Initialised: true}, {Name: "b", Initialised: true}}
	}
	return web.RootWithProjects(web.Handler(nil), func(n string) http.Handler { return hs[n] }, list)
}

func TestHeaderNoProjectList(t *testing.T) {
	body := serve(headerRoot(t, 3), "GET", "/p/a/runs").Body.String()
	if strings.Contains(body, "/p/b") || strings.Contains(body, ">b<") {
		t.Errorf("page names project b:\n%s", body)
	}
	for _, want := range []string{
		`data-rz-picker="/picker"`, `<span class="rz-count is-due" aria-label="3 need you">3</span>`,
		`<a href="/p/a/">Home</a>`, `<a href="/p/a/runs" aria-current="page">Runs</a>`, `<a href="/p/a/roadmap">Roadmap</a>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("header lacks %s", want)
		}
	}
}

func TestHeaderNoBadgeWhenNothingWaits(t *testing.T) {
	body := serve(headerRoot(t, 0), "GET", "/p/a/runs").Body.String()
	if strings.Contains(body, `rz-count`) {
		t.Errorf("a project with nothing waiting shows a count badge")
	}
}

func TestHeaderPickerOnRootPages(t *testing.T) {
	h := headerRoot(t, 3)
	for _, path := range []string{"/inbox", "/projects"} {
		body := serve(h, "GET", path).Body.String()
		if !strings.Contains(body, `data-rz-picker="/picker"`) || !strings.Contains(body, `<span class="rz-picker-cur">All projects</span>`) {
			t.Errorf("%s has no picker button for all projects:\n%s", path, body)
		}
	}
}

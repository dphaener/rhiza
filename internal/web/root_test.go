package web_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/web"
)

func TestRootRoutesProjects(t *testing.T) {
	say := func(s string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(s + " " + r.URL.Path)) })
	}
	h := web.Root(say("home"), func(n string) http.Handler {
		if n == "a" || n == "b" {
			return say(n)
		}
		return nil
	})
	for target, want := range map[string]string{
		"/":            "home /",
		"/runs":        "home /runs",
		"/p/a/":        "a /",
		"/p/b/runs":    "b /runs",
		"/p/b/api/run": "b /api/run",
	} {
		rec := serve(h, "GET", target)
		if rec.Code != 200 || rec.Body.String() != want {
			t.Errorf("%s: %d %q, want %q", target, rec.Code, rec.Body, want)
		}
	}
	if rec := serve(h, "GET", "/p/c/"); rec.Code != 404 {
		t.Errorf("unknown project: %d", rec.Code)
	}
	if rec := serve(h, "GET", "/p/a"); rec.Code != 301 || rec.Header().Get("Location") != "/p/a/" {
		t.Errorf("bare project: %d %q", rec.Code, rec.Header().Get("Location"))
	}
}

// With a home project named, / is the all-projects home, an unscoped page
// redirects to that page of the home project, and the shared assets and
// streams stay served by home.
func TestRootRedirectsUnscopedToHome(t *testing.T) {
	say := func(s string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(s + " " + r.URL.Path)) })
	}
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "alpha", Initialised: true}, {Name: "beta", Initialised: true}}
	}
	h := web.RootWithInbox(say("home"), "alpha", func(n string) http.Handler { return say(n) }, list, nil)
	rec := serve(h, "GET", "/")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `href="/p/beta/"`) {
		t.Errorf("/: %d, want the projects home:\n%s", rec.Code, rec.Body)
	}
	for target, want := range map[string]string{
		"/runs":          "/p/alpha/runs",
		"/roadmap":       "/p/alpha/roadmap",
		"/runs/r1?x=1":   "/p/alpha/runs/r1?x=1",
		"/changes/c%2F1": "/p/alpha/changes/c%2F1",
	} {
		rec := serve(h, "GET", target)
		if rec.Code != http.StatusFound || rec.Header().Get("Location") != want {
			t.Errorf("%s: %d %q, want a redirect to %q", target, rec.Code, rec.Header().Get("Location"), want)
		}
	}
	for _, target := range []string{"/static/css/rhiza.css", "/events", "/api/status", "/site.webmanifest"} {
		if rec := serve(h, "GET", target); rec.Code != 200 || rec.Body.String() != "home "+target {
			t.Errorf("%s: %d %q, want served by home", target, rec.Code, rec.Body)
		}
	}
	if rec := serve(h, "POST", "/api/run"); rec.Code != 200 || rec.Body.String() != "home /api/run" {
		t.Errorf("POST /api/run: %d %q", rec.Code, rec.Body)
	}
}

package web_test

import (
	"net/http"
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

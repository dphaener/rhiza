package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A root of two projects serves both through one handler: each project's
// pages and API answer from its own model, and an unknown project is not
// found.
func TestWebRootHandlerServesEveryProject(t *testing.T) {
	root := t.TempDir()
	intents := map[string]string{"alpha": "Alpha intent", "beta": "Beta intent"}
	var home *Ctx
	for name, intent := range intents {
		c := testCtx(filepath.Join(root, name))
		pinFreeBytes(c)
		if err := os.MkdirAll(c.Root, 0o755); err != nil {
			t.Fatal(err)
		}
		c.IndexPath = filepath.Join(t.TempDir(), "index.db")
		if _, err := callOp(t, c, "init", `{}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, c, "change.new", `{"intent":"`+intent+`"}`); err != nil {
			t.Fatal(err)
		}
		if name == "alpha" {
			home = c
		}
	}
	// Registered after the temp folders, so it runs before they are removed.
	t.Cleanup(healthRefreshing.Wait)
	if err := os.WriteFile(filepath.Join(root, "rhiza.toml"), []byte("[projects.alpha]\npath = \"alpha\"\n[projects.beta]\npath = \"beta\"\n[projects.bare]\npath = \"bare\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := WebRootHandler(home, root, "alpha", WebProjectHandler(home))
	get := func(target string) (int, string) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		return rec.Code, rec.Body.String()
	}
	for name := range intents {
		code, body := get("/p/" + name + "/api/roadmap")
		if code != 200 {
			t.Fatalf("/p/%s/api/roadmap: %d %s", name, code, body)
		}
		code, _ = get("/p/" + name + "/runs")
		if code != 200 {
			t.Errorf("/p/%s/runs: %d", name, code)
		}
	}
	// The web UI opens on the all-projects home; an old unscoped page
	// redirects into the home project.
	if code, body := get("/"); code != 200 || !strings.Contains(body, "/p/alpha/") || !strings.Contains(body, "/p/beta/") {
		t.Errorf("/: %d, want the projects home listing both:\n%s", code, body)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/runs", nil))
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/p/alpha/runs" {
		t.Errorf("/runs: %d %q, want a redirect to /p/alpha/runs", rec.Code, rec.Header().Get("Location"))
	}
	for _, n := range []string{"bare", "nope"} {
		if code, _ := get("/p/" + n + "/runs"); code != 404 {
			t.Errorf("/p/%s/runs: %d, want 404", n, code)
		}
	}
	// The two projects answer from their own models.
	_, a := get("/p/alpha/api/status")
	_, b := get("/p/beta/api/status")
	if !strings.Contains(a, "Alpha intent") || strings.Contains(a, "Beta intent") ||
		!strings.Contains(b, "Beta intent") || strings.Contains(b, "Alpha intent") {
		t.Errorf("projects mixed up:\nalpha %s\nbeta %s", a, b)
	}
}

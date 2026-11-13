package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Under a root, a project's web roadmap and its API (even with all=true) list
// only that project's milestone and change sets, never another project's.
func TestWebRootRoadmapIsProjectScoped(t *testing.T) {
	root := t.TempDir()
	names := []string{"alpha", "beta"}
	var home *Ctx
	for _, name := range names {
		c := testCtx(filepath.Join(root, name))
		pinFreeBytes(c)
		if err := os.MkdirAll(c.Root, 0o755); err != nil {
			t.Fatal(err)
		}
		c.IndexPath = filepath.Join(t.TempDir(), "index.db")
		if _, err := callOp(t, c, "init", `{}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, c, "change.new", `{"intent":"`+name+` setup"}`); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, c, "add", `{"kind":"intent","handle":"`+name+`-ms","fields":["milestone=1","statement=`+name+` milestone"]}`); err != nil {
			t.Fatal(err)
		}
		for _, s := range []string{"propose", "accept", "land"} {
			move(t, c, s, `{}`)
		}
		for _, n := range []string{"one", "two"} {
			resp, err := callOp(t, c, "change.new", `{"intent":"`+name+` change `+n+`"}`)
			if err != nil {
				t.Fatal(err)
			}
			cs := resp.(ChangeNewResponse).ChangeSet.String()
			if _, err := callOp(t, c, "link", `{"from":"`+cs+`","type":"serves","to":"`+name+`-ms"}`); err != nil {
				t.Fatal(err)
			}
			if _, err := callOp(t, c, "propose", `{}`); err != nil {
				t.Fatal(err)
			}
		}
		if name == "alpha" {
			home = c
		}
	}
	t.Cleanup(healthRefreshing.Wait)
	if err := os.WriteFile(filepath.Join(root, "rhiza.toml"), []byte("[projects.alpha]\npath = \"alpha\"\n[projects.beta]\npath = \"beta\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := WebRootHandler(home, root, "alpha", WebProjectHandler(home))
	get := func(target string) string {
		t.Helper()
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		if rec.Code != 200 {
			t.Fatalf("%s: %d %s", target, rec.Code, rec.Body.String())
		}
		return rec.Body.String()
	}
	for i, name := range names {
		other := names[1-i]
		page := get("/p/" + name + "/roadmap")
		for _, want := range []string{name + " milestone", name + " change one", name + " change two"} {
			if !strings.Contains(page, want) {
				t.Errorf("/p/%s/roadmap lacks %q:\n%s", name, want, page)
			}
		}
		if strings.Contains(page, other+" ") {
			t.Errorf("/p/%s/roadmap shows %s's rows:\n%s", name, other, page)
		}
		var resp RoadmapResponse
		if err := json.Unmarshal([]byte(get("/p/"+name+"/api/roadmap?all=true")), &resp); err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, g := range resp.Groups {
			for _, cs := range g.ChangeSets {
				n++
				if !strings.HasPrefix(cs.Intent, name+" ") {
					t.Errorf("/p/%s/api/roadmap?all=true returned %q", name, cs.Intent)
				}
			}
		}
		if n != 2 {
			t.Errorf("/p/%s/api/roadmap?all=true: %d change sets, want 2", name, n)
		}
	}
}

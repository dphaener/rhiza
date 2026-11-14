package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var hrefRe = regexp.MustCompile(`href="(/[^"]*)"`)

// Under a root of two projects, every row on each project's roadmap and every
// depends-on link on those change pages opens a page in that same project.
func TestWebRootRoadmapRowsResolve(t *testing.T) {
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
		var made []string
		for i, n := range []string{"milestone", "backlog"} {
			resp, err := callOp(t, c, "change.new", `{"intent":"`+name+` `+n+` change"}`)
			if err != nil {
				t.Fatal(err)
			}
			cs := resp.(ChangeNewResponse).ChangeSet.String()
			made = append(made, cs)
			if i == 0 {
				if _, err := callOp(t, c, "link", `{"from":"`+cs+`","type":"serves","to":"`+name+`-ms"}`); err != nil {
					t.Fatal(err)
				}
			} else if _, err := callOp(t, c, "link", `{"from":"`+cs+`","type":"depends-on","to":"`+made[0]+`"}`); err != nil {
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
	get := func(target string) (int, string) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		return rec.Code, rec.Body.String()
	}
	for i, name := range names {
		other := names[1-i]
		prefix := "/p/" + name + "/changes/"
		code, page := get("/p/" + name + "/roadmap")
		if code != 200 {
			t.Fatalf("/p/%s/roadmap: %d", name, code)
		}
		if strings.Contains(page, other+" ") {
			t.Errorf("/p/%s/roadmap names %s's change sets:\n%s", name, other, page)
		}
		var rows []string
		for _, m := range hrefRe.FindAllStringSubmatch(page, -1) {
			if strings.Contains(m[1], "/changes/") {
				rows = append(rows, m[1])
			}
		}
		if len(rows) < 2 {
			t.Fatalf("/p/%s/roadmap has %d change rows, want at least 2:\n%s", name, len(rows), page)
		}
		for _, href := range rows {
			if !strings.HasPrefix(href, prefix) {
				t.Errorf("/p/%s/roadmap row %s is outside %s", name, href, prefix)
			}
			code, body := get(href)
			if code != 200 {
				t.Errorf("GET %s: %d", href, code)
				continue
			}
			for _, m := range hrefRe.FindAllStringSubmatch(body, -1) {
				if !strings.Contains(m[1], "/changes/") {
					continue
				}
				if !strings.HasPrefix(m[1], prefix) {
					t.Errorf("%s links %s, outside %s", href, m[1], prefix)
				}
				if code, _ := get(m[1]); code != 200 {
					t.Errorf("GET %s (linked from %s): %d", m[1], href, code)
				}
			}
		}
	}
}

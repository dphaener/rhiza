package web_test

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

var hrefRe = regexp.MustCompile(`(?:href|src|data-rz-fragment)="([^"]*)"`)

// TestLinksResolve crawls a project's Home, roadmap, a change set page, a run
// page and a node page, in each project of a two-project root, and follows
// every internal link: each must answer 200 inside the same project.
func TestLinksResolve(t *testing.T) {
	hs := map[string]http.Handler{}
	type proj struct{ cs, run string }
	ps := map[string]proj{}
	for i, name := range []string{"a", "b"} {
		ctx := fixtureCtx(t)
		ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(int64(30+i))))
		cs := newChange(t, ctx, name+" links work", "intention")
		callJSON(t, ctx, "add", `{"kind":"intent","handle":"links-intent","fields":["statement=Links work in `+name+`"]}`)
		callJSON(t, ctx, "add", `{"kind":"capability","handle":"links-cap","fields":["statement=A link opens in `+name+`","acceptance=Crawl it"]}`)
		callJSON(t, ctx, "link", `{"from":"links-cap","type":"serves","to":"links-intent"}`)
		callJSON(t, ctx, "add", `{"kind":"decision","handle":"links-dec","fields":["statement=Links stay in `+name+`","question=Where do links go?","wording=In the project"]}`)
		callJSON(t, ctx, "link", `{"from":"links-cap","type":"decided-by","to":"links-dec"}`)
		callJSON(t, ctx, "propose", `{}`)
		callJSON(t, ctx, "accept", `{}`)
		callJSON(t, ctx, "land", `{}`)
		runID := "01LINKRUN" + strings.ToUpper(name)
		rs := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
		if err := rs.Create(runID, run.Entry{Type: run.TypeStarted, Started: &run.Started{Workflow: "build", Change: cs, Title: name + " run"}}); err != nil {
			t.Fatal(err)
		}
		if err := rs.Append(runID, run.Entry{Type: run.TypeEntered, Step: "leaves", Visit: 1}); err != nil {
			t.Fatal(err)
		}
		hs[name] = web.HandlerWithRuns(nil, api.WebRunsSource(ctx), api.ReadOnlyHandler(ctx))
		ps[name] = proj{cs, runID}
	}
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "a", Initialised: true}, {Name: "b", Initialised: true}}
	}
	h := web.RootWithProjects(web.Handler(nil), func(n string) http.Handler { return hs[n] }, list)

	for name, p := range ps {
		base := "/p/" + name
		seen := map[string]bool{}
		var queue []string
		queue = append(queue, base+"/", base+"/roadmap", base+"/changes/"+p.cs, base+"/runs/"+p.run)
		queue = append(queue, base+"/nodes/links-cap")
		for len(queue) > 0 {
			page := queue[0]
			queue = queue[1:]
			if seen[page] {
				continue
			}
			seen[page] = true
			for _, l := range links(t, h, page) {
				if seen[l] {
					continue
				}
				seen[l] = true
				if !inProject(l, base) {
					if !shared(l) {
						t.Errorf("%s links to %s, outside project %s", page, l, name)
					}
					continue
				}
				if rec := serve(h, "GET", l); rec.Code != 200 {
					t.Errorf("%s links to %s: %d", page, l, rec.Code)
					continue
				}
				queue = append(queue, l)
			}
		}
		if len(seen) < 8 {
			t.Errorf("%s: crawled only %d pages", name, len(seen))
		}
	}
}

// links is every internal link of a page.
func links(t *testing.T, h http.Handler, target string) []string {
	t.Helper()
	rec := serve(h, "GET", target)
	if rec.Code != 200 {
		t.Fatalf("%s: %d", target, rec.Code)
	}
	var out []string
	for _, m := range hrefRe.FindAllStringSubmatch(rec.Body.String(), -1) {
		l := strings.ReplaceAll(m[1], "&amp;", "&")
		if strings.HasPrefix(l, "/") && !strings.HasPrefix(l, "//") {
			out = append(out, l)
		}
	}
	return out
}

func inProject(l, base string) bool { return strings.HasPrefix(l, base+"/") || l == base }

// shared is a root path that serves the same in every project.
func shared(l string) bool {
	for _, p := range []string{"/static/", "/site.webmanifest", "/projects", "/picker", "/inbox"} {
		if strings.HasPrefix(l, p) {
			return true
		}
	}
	return false
}

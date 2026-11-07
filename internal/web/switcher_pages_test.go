package web_test

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// scopedProject is one project of a two-project root: a change set and a
// run of its own, each with a title no other project has.
type scopedProject struct {
	name, cs, run, csTitle, runTitle string
}

func twoProjectsWithRuns(t *testing.T) (http.Handler, [2]scopedProject) {
	t.Helper()
	ps := [2]scopedProject{
		{name: "a", run: "01ALPHARUN", csTitle: "Alpha only work", runTitle: "Alpha run title"},
		{name: "b", run: "01BETARUN", csTitle: "Beta only work", runTitle: "Beta run title"},
	}
	hs := map[string]http.Handler{}
	for i := range ps {
		p := &ps[i]
		ctx := fixtureCtx(t)
		// a seed of its own, so the two projects' change sets differ in ID too
		ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(int64(10+i))))
		p.cs = newChange(t, ctx, p.csTitle, "intention")
		callJSON(t, ctx, "propose", `{}`)
		rs := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
		if err := rs.Create(p.run, run.Entry{Type: run.TypeStarted, Started: &run.Started{Workflow: "build", Change: p.cs, Title: p.runTitle}}); err != nil {
			t.Fatal(err)
		}
		if err := rs.Append(p.run, run.Entry{Type: run.TypeEntered, Step: "leaves", Visit: 1}); err != nil {
			t.Fatal(err)
		}
		hs[p.name] = web.HandlerWithRuns(nil, api.WebRunsSource(ctx), api.ReadOnlyHandler(ctx))
	}
	list := func() []web.ProjectView {
		return []web.ProjectView{{Name: "a", Initialised: true}, {Name: "b", Initialised: true}, {Name: "c"}}
	}
	return web.RootWithProjects(web.Handler(nil), func(n string) http.Handler { return hs[n] }, list), ps
}

// Every page kind under /p/<a>/ and /p/<b>/ shows only that project's data,
// links the other project's page of the same kind (a detail page its list),
// and an unknown project's same page is a 404.
func TestEveryPageKindIsScopedToItsProject(t *testing.T) {
	h, ps := twoProjectsWithRuns(t)
	for i, p := range ps {
		o := ps[1-i]
		for _, c := range []struct {
			page, otherKind string
			own             []string
		}{
			{"", "/", []string{p.runTitle, `href="/p/` + p.name + `/runs/` + p.run + `"`}},
			{"roadmap", "/roadmap", []string{p.csTitle, `href="/p/` + p.name + `/changes/` + p.cs + `"`}},
			{"runs", "/runs", []string{p.runTitle, `href="/p/` + p.name + `/runs/` + p.run + `"`}},
			{"runs/" + p.run, "/runs", []string{p.runTitle, "leaves"}},
			{"changes/" + p.cs, "/roadmap", []string{p.csTitle}},
		} {
			target := "/p/" + p.name + "/" + c.page
			rec := serve(h, "GET", target)
			if rec.Code != 200 {
				t.Errorf("%s: %d %s", target, rec.Code, rec.Body)
				continue
			}
			body := rec.Body.String()
			for _, w := range c.own {
				if !strings.Contains(body, w) {
					t.Errorf("%s lacks its own %q", target, w)
				}
			}
			for _, bad := range []string{o.csTitle, o.runTitle, o.cs, o.run} {
				if strings.Contains(body, bad) {
					t.Errorf("%s shows project %s's %q", target, o.name, bad)
				}
			}
			if w := `href="/p/` + o.name + c.otherKind + `"`; !strings.Contains(body, w) {
				t.Errorf("%s lacks switcher link %s", target, w)
			}
			if w := `href="/p/` + p.name + c.otherKind + `"`; !strings.Contains(body, w) {
				t.Errorf("%s lacks its own switcher entry %s", target, w)
			}
			if rec := serve(h, "GET", "/p/nope/"+c.page); rec.Code != 404 {
				t.Errorf("/p/nope/%s: %d, want 404", c.page, rec.Code)
			}
		}
		// another project's run or change set is not found under this one
		for _, target := range []string{"/p/" + p.name + "/runs/" + o.run, "/p/" + p.name + "/changes/" + o.cs} {
			if rec := serve(h, "GET", target); rec.Code != 404 {
				t.Errorf("%s: %d, want 404", target, rec.Code)
			}
		}
	}
}

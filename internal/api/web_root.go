package api

import (
	"net/http"
	"sync"

	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/web"
)

// WebProjectHandler is the web UI of the project base works in: its pages,
// events and read-only API, all over that project's own data. base is
// copied, as the poll goes on changing the ctx it came from.
func WebProjectHandler(base *Ctx) http.Handler {
	snap := *base
	health := NewWebHealth(&snap)
	return web.HandlerWithMilestones(WebDaemonSource(&snap, health), WebRunsSource(&snap), WebInboxSource(&snap), WebMilestonesSource(&snap), ReadOnlyHandler(&snap), web.Events(WebEventSources(&snap, health), web.PollEvery))
}

// WebRootHandler is the one web UI of the root daemon: home, the home
// project's UI, at /, and every other registered, initialised project's at
// /p/<project>/ (the home project also answers there). Projects are read
// from the root's rhiza.toml on each request, so a newly registered one
// is served at once; each project's handler is built once and kept.
func WebRootHandler(ctx *Ctx, root, homeName string, home http.Handler) http.Handler {
	var mu sync.Mutex
	type built struct {
		dir string
		h   http.Handler
	}
	cache := map[string]built{}
	return web.RootWithInbox(home, func(name string) http.Handler {
		reg, err := rootcfg.Load(root)
		if err != nil {
			return nil
		}
		if _, ok := reg.Projects[name]; !ok {
			return nil
		}
		if name == homeName {
			return home
		}
		dir := reg.Dir(root, name)
		if !projectInitialised(dir) {
			return nil
		}
		mu.Lock()
		defer mu.Unlock()
		if b, ok := cache[name]; ok && b.dir == dir {
			return b.h
		}
		h := WebProjectHandler(projectCtx(ctx, dir))
		cache[name] = built{dir, h}
		return h
	}, WebProjectsSource(ctx, root, homeName), WebMergedInboxSource(ctx, root))
}

// WebProjectsSource reads every project registered in the root for the web
// UI's all-projects home: live runs, inbox count and health of each
// initialised one, a project with no .rhiza/ yet marked as such. Each
// project's health report is kept, as WebProjectHandler's is.
func WebProjectsSource(ctx *Ctx, root, homeName string) web.ProjectsSource {
	var mu sync.Mutex
	healths := map[string]*WebHealth{}
	return func() []web.ProjectView {
		reg, err := rootcfg.Load(root)
		if err != nil {
			return nil
		}
		var views []web.ProjectView
		for _, name := range reg.Names() {
			v := web.ProjectView{Name: name}
			dir := reg.Dir(root, name)
			if projectInitialised(dir) {
				v.Initialised = true
				c := projectCtx(ctx, dir)
				for _, r := range WebRunsSource(c)() {
					if r.State == "running" {
						v.Running++
					} else {
						v.Waiting++
					}
				}
				v.Inbox = len(WebInboxSource(c)())
				mu.Lock()
				h := healths[name]
				if h == nil {
					h = NewWebHealth(c)
					healths[name] = h
				}
				mu.Unlock()
				v.Health = WebDaemonSource(c, h)().Health
			}
			views = append(views, v)
		}
		return views
	}
}

// WebMergedInboxSource reads the inbox of every initialised project
// registered in the root, for the web UI's merged inbox.
func WebMergedInboxSource(ctx *Ctx, root string) web.MergedInboxSource {
	return func() []web.ProjectInbox {
		reg, err := rootcfg.Load(root)
		if err != nil {
			return nil
		}
		var out []web.ProjectInbox
		for _, name := range reg.Names() {
			dir := reg.Dir(root, name)
			if !projectInitialised(dir) {
				continue
			}
			out = append(out, web.ProjectInbox{Name: name, Items: WebInboxSource(projectCtx(ctx, dir))()})
		}
		return out
	}
}

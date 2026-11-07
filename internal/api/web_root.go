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
	return web.Root(home, func(name string) http.Handler {
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
	})
}

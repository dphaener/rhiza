package web

import (
	"net/http"
	"net/url"
	"strings"
)

// ProjectPrefix is where a project's pages live on the root's listener:
// /p/<project>/ serves that project's UI, so /p/<project>/runs is its Runs.
const ProjectPrefix = "/p/"

// ProjectResolver gives the handler serving the project name, or nil when no
// such project is registered and initialised.
type ProjectResolver func(name string) http.Handler

// Root is the handler of the root daemon's one listener. A request under
// /p/<project>/ is served by that project's handler with the prefix taken
// off; an unknown project is not found; everything else is served by home,
// the handler of the daemon's home project.
func Root(home http.Handler, project ProjectResolver) http.Handler {
	return RootWithProjects(home, project, nil)
}

// ProjectsPath is the all-projects home on the root's listener.
const ProjectsPath = "/projects"

// InboxPath is the merged inbox on the root's listener.
const InboxPath = "/inbox"

// RootWithProjects is Root that also serves the all-projects home at
// /projects, listing what projects reads.
func RootWithProjects(home http.Handler, project ProjectResolver, projects ProjectsSource) http.Handler {
	return RootWithInbox(home, "", project, projects, nil)
}

// shared reports whether path is one of the assets or streams every page
// loads from the root, whichever project it shows: these stay served by
// home rather than redirected into a project.
func shared(path string) bool {
	return strings.HasPrefix(path, "/static/") || strings.HasPrefix(path, "/api/") ||
		path == "/events" || path == "/site.webmanifest"
}

// ProjectNames lists the initialised projects a project page's header
// switcher offers, by name. It is read only when a page is rendered, never
// for an asset, fragment or stream, so it should stay cheap: names, not the
// projects' runs or inboxes.
type ProjectNames func() []string

// RootWithInbox is RootWithProjects that also serves the merged inbox at
// /inbox, listing what inbox reads. The switcher's projects are the
// initialised ones projects reads; RootWithSwitcher takes a cheaper list.
func RootWithInbox(home http.Handler, homeName string, project ProjectResolver, projects ProjectsSource, inbox MergedInboxSource) http.Handler {
	var names ProjectNames
	if projects != nil {
		names = func() []string {
			var out []string
			for _, p := range projects() {
				if p.Initialised {
					out = append(out, p.Name)
				}
			}
			return out
		}
	}
	return RootWithSwitcher(home, homeName, project, projects, inbox, names)
}

// RootWithSwitcher is RootWithInbox with the switcher's projects listed by
// names, read only when a project page is rendered. With homeName set, the web UI opens on
// the all-projects home: / serves it, and a page request on any other
// unscoped path is redirected to that page of the home project, under
// /p/<homeName>/, so every project page carries its switcher. With no
// homeName, unscoped paths are served by home as they are.
func RootWithSwitcher(home http.Handler, homeName string, project ProjectResolver, projects ProjectsSource, inbox MergedInboxSource, names ProjectNames) http.Handler {
	list := ProjectsHome(projects)
	merged := MergedInbox(inbox)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		get := r.Method == "GET" || r.Method == "HEAD"
		if get && (r.URL.Path == ProjectsPath || (homeName != "" && r.URL.Path == "/")) {
			list.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == InboxPath && get {
			merged.ServeHTTP(w, r)
			return
		}
		rest, ok := strings.CutPrefix(r.URL.Path, ProjectPrefix)
		if !ok {
			if homeName != "" && get && !shared(r.URL.Path) {
				to := ProjectPrefix + url.PathEscape(homeName) + r.URL.EscapedPath()
				if r.URL.RawQuery != "" {
					to += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, to, http.StatusFound)
				return
			}
			home.ServeHTTP(w, r)
			return
		}
		name, tail, _ := strings.Cut(rest, "/")
		var h http.Handler
		if name != "" && project != nil {
			h = project(name)
		}
		if h == nil {
			http.NotFound(w, r)
			return
		}
		if !strings.Contains(rest, "/") {
			http.Redirect(w, r, ProjectPrefix+name+"/", http.StatusMovedPermanently)
			return
		}
		sc := &scope{Name: name, Prefix: ProjectPrefix + url.PathEscape(name), Tail: tail, names: names}
		r2 := withScope(r.Clone(r.Context()), sc)
		r2.URL.Path = "/" + tail
		r2.URL.RawPath = ""
		h.ServeHTTP(w, r2)
	})
}

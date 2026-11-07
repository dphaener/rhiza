package web

import (
	"net/http"
	"net/url"
	"slices"
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
	return RootWithInbox(home, project, projects, nil)
}

// RootWithInbox is RootWithProjects that also serves the merged inbox at
// /inbox, listing what inbox reads.
func RootWithInbox(home http.Handler, project ProjectResolver, projects ProjectsSource, inbox MergedInboxSource) http.Handler {
	list := ProjectsHome(projects)
	merged := MergedInbox(inbox)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == ProjectsPath && (r.Method == "GET" || r.Method == "HEAD") {
			list.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == InboxPath && (r.Method == "GET" || r.Method == "HEAD") {
			merged.ServeHTTP(w, r)
			return
		}
		rest, ok := strings.CutPrefix(r.URL.Path, ProjectPrefix)
		if !ok {
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
		sc := &scope{Name: name, Prefix: ProjectPrefix + url.PathEscape(name), Tail: tail, Names: []string{name}}
		if projects != nil {
			sc.Names = sc.Names[:0]
			for _, p := range projects() {
				if p.Initialised {
					sc.Names = append(sc.Names, p.Name)
				}
			}
			if !slices.Contains(sc.Names, name) {
				sc.Names = append(sc.Names, name)
			}
		}
		r2 := withScope(r.Clone(r.Context()), sc)
		r2.URL.Path = "/" + tail
		r2.URL.RawPath = ""
		h.ServeHTTP(w, r2)
	})
}

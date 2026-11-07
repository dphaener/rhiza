package web

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// A scope is the project a request is served for under /p/<project>/: the
// prefix its links carry, and the projects its header switcher offers.
type scope struct {
	Name   string
	Prefix string // /p/<project>, no trailing slash
	Tail   string // the page's path inside the project, without a leading slash
	Names  []string
}

type scopeKey struct{}

func scopeOf(r *http.Request) *scope {
	s, _ := r.Context().Value(scopeKey{}).(*scope)
	return s
}

func withScope(r *http.Request, s *scope) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), scopeKey{}, s))
}

// SwitchItem is one project in the switcher: where the same kind of page is
// in that project.
type SwitchItem struct {
	Name, Href string
	Current    bool
}

// Switcher is the header's project switcher.
type Switcher struct{ Items []SwitchItem }

// pageKind is the list page of the kind of page tail is: a run's or change
// set's own page does not exist in another project, so it keeps its list.
func pageKind(tail string) string {
	first, _, _ := strings.Cut(tail, "/")
	switch first {
	case "runs", "roadmap":
		return "/" + first
	case "changes":
		return "/roadmap"
	}
	return "/"
}

// Switcher lists the projects with the page kind kept, the current one marked.
func (s *scope) Switcher() *Switcher {
	kind := pageKind(s.Tail)
	sw := &Switcher{}
	for _, n := range s.Names {
		sw.Items = append(sw.Items, SwitchItem{Name: n, Href: ProjectPrefix + url.PathEscape(n) + kind, Current: n == s.Name})
	}
	return sw
}

var rootLink = regexp.MustCompile(`(href|src|data-rz-fragment)="/([^"]*)"`)

// rewrite points the root-relative links of a page served under /p/<project>/
// at that project, leaving links already into /p/, /projects or /inbox as
// they are. A nil scope leaves the page untouched.
func (s *scope) rewrite(b []byte) []byte {
	if s == nil {
		return b
	}
	return rootLink.ReplaceAllFunc(b, func(m []byte) []byte {
		sub := rootLink.FindSubmatch(m)
		path := string(sub[2])
		if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "p/") || path == "projects" || path == "inbox" {
			return m
		}
		return []byte(string(sub[1]) + `="` + s.Prefix + "/" + path + `"`)
	})
}

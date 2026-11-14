package web

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// A scope is the project a request is served for under /p/<project>/: the
// prefix its links carry, and the projects its header switcher offers.
type scope struct {
	Name   string
	Prefix string // /p/<project>, no trailing slash
	Tail   string // the page's path inside the project, without a leading slash
}

type scopeKey struct{}

func scopeOf(r *http.Request) *scope {
	s, _ := r.Context().Value(scopeKey{}).(*scope)
	return s
}

func withScope(r *http.Request, s *scope) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), scopeKey{}, s))
}

// pageKind is the list page of the kind of page tail is: a run's or change
// set's own page does not exist in another project, so it keeps its list.
func pageKind(tail string) string {
	first, _, _ := strings.Cut(tail, "/")
	switch first {
	case "runs", "roadmap":
		return "/" + first
	case "changes", "nodes":
		return "/roadmap"
	}
	return "/"
}

var rootLink = regexp.MustCompile(`(href|src|action|data-rz-fragment)="/([^"]*)"`)

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
		if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "p/") || path == "projects" || path == "picker" || path == "inbox" {
			return m
		}
		return []byte(string(sub[1]) + `="` + s.Prefix + "/" + path + `"`)
	})
}

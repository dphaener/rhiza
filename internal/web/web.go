// Package web is the daemon's read-only web UI server: the port setting, the
// loopback address, the handler and the listener manager that follows the
// setting from poll to poll.
package web

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// The setting that names the port, its default and its off value.
const (
	PortKey     = "rhiza.web.port"
	DefaultPort = "7787"
	Off         = "off"
	// Host is the only address the UI is ever served on.
	Host = "127.0.0.1"
)

// ParsePort checks a rhiza.web.port value: a port number from 1 to 65535, or
// off. It returns the port, 0 for off.
func ParsePort(v string) (int, error) {
	v = strings.TrimSpace(v)
	if v == Off {
		return 0, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > 65535 || strconv.Itoa(n) != v {
		return 0, fmt.Errorf("must be a port number from 1 to 65535, or %s", Off)
	}
	return n, nil
}

// Addr is the address to listen on for a setting value: always on 127.0.0.1,
// and empty when the setting is off.
func Addr(v string) (string, error) {
	n, err := ParsePort(v)
	if err != nil || n == 0 {
		return "", err
	}
	return net.JoinHostPort(Host, strconv.Itoa(n)), nil
}

//go:embed site.webmanifest static
var assets embed.FS

//go:embed templates
var templateFS embed.FS

// pages holds each page's template set: the shared layout plus the page's
// own "content". They are parsed once, from the binary.
var pages = map[string]*template.Template{
	"overview": template.Must(template.New("").Funcs(template.FuncMap{"section": section}).ParseFS(templateFS, "templates/layout.html", "templates/partials.html", "templates/overview.html", "templates/daemon.html", "templates/runs.html", "templates/inbox.html", "templates/milestones.html")),
}

// page is what a template is rendered with.
type page struct {
	Title, Nav string
	// Daemon is the header's daemon pill; empty Text leaves it out.
	Daemon Pill
	// DaemonView is the Overview's daemon panel.
	DaemonView DaemonView
	// Runs is the Overview's live runs panel.
	Runs RunsPanel
	// Inbox is the Overview's inbox panel.
	Inbox InboxPanel
	// Milestones is the Overview's milestone progress panel.
	Milestones MilestonesPanel
}

func render(w http.ResponseWriter, name string, p page) {
	var buf bytes.Buffer
	if err := pages[name].ExecuteTemplate(&buf, "layout", p); err != nil {
		http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

// Handler serves the UI's embedded assets, read-only: GET and HEAD only.
// api, if not nil, serves /api/ (see api.ReadOnlyHandler); without it
// /api/ is not found. events, if given, serves /events (see Events).
func Handler(api http.Handler, events ...http.Handler) http.Handler {
	return HandlerWith(nil, api, events...)
}

// HandlerWith is Handler with the source of the daemon the Overview shows;
// nil shows a daemon that is not running.
func HandlerWith(daemon DaemonSource, api http.Handler, events ...http.Handler) http.Handler {
	return HandlerWithRuns(daemon, nil, api, events...)
}

// HandlerWithRuns is HandlerWith with the source of the live runs the
// Overview lists; nil lists none.
func HandlerWithRuns(daemon DaemonSource, runs RunsSource, api http.Handler, events ...http.Handler) http.Handler {
	return HandlerWithInbox(daemon, runs, nil, api, events...)
}

// HandlerWithInbox is HandlerWithRuns with the source of the inbox the
// Overview lists; nil lists nothing.
func HandlerWithInbox(daemon DaemonSource, runs RunsSource, inbox InboxSource, api http.Handler, events ...http.Handler) http.Handler {
	return HandlerWithMilestones(daemon, runs, inbox, nil, api, events...)
}

// HandlerWithMilestones is HandlerWithInbox with the source of the open
// milestones the Overview shows the progress of; nil shows none.
func HandlerWithMilestones(daemon DaemonSource, runs RunsSource, inbox InboxSource, milestones MilestonesSource, api http.Handler, events ...http.Handler) http.Handler {
	mux := http.NewServeMux()
	for _, e := range events {
		if e != nil {
			mux.Handle("/events", e)
		}
	}
	if api != nil {
		mux.Handle("/api/", api)
	}
	mux.HandleFunc("GET "+FragmentPrefix+"{name}", serveFragment)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		d := daemonOf(r)
		render(w, "overview", page{Title: "Overview", Nav: "overview", DaemonView: d, Daemon: d.Pill(), Runs: runsOf(r), Inbox: inboxOf(r), Milestones: milestonesOf(r)})
	})
	static, _ := fs.Sub(assets, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		b, _ := assets.ReadFile("site.webmanifest")
		w.Header().Set("Content-Type", "application/manifest+json")
		w.Write(b)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if daemon != nil {
			r = r.WithContext(context.WithValue(r.Context(), daemonKey{}, daemon))
		}
		if runs != nil {
			r = r.WithContext(context.WithValue(r.Context(), runsKey{}, runs))
		}
		if inbox != nil {
			r = r.WithContext(context.WithValue(r.Context(), inboxKey{}, inbox))
		}
		if milestones != nil {
			r = r.WithContext(context.WithValue(r.Context(), milestonesKey{}, milestones))
		}
		mux.ServeHTTP(w, r)
	})
}

// Listener is a running listener, which can be stopped.
type Listener interface{ Close() error }

// Factory starts serving h on addr.
type Factory func(addr string, h http.Handler) (Listener, error)

// NetFactory is the real Factory: an HTTP server on a TCP listener.
func NetFactory(addr string, h http.Handler) (Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &http.Server{Handler: h}
	go srv.Serve(ln)
	return srv, nil
}

// Manager keeps one listener in step with the port setting. It is used from
// one goroutine, the daemon's poll.
type Manager struct {
	Factory Factory
	Handler http.Handler

	addr string
	l    Listener
}

// Addr is the address now served, empty when none.
func (m *Manager) Addr() string { return m.addr }

// Apply makes the listener match the setting value: it starts one, moves it
// to a new port, stops it on off, and leaves a matching one alone. A start
// that fails leaves nothing listening and is tried again at the next call.
func (m *Manager) Apply(v string) error {
	want, err := Addr(v)
	if err != nil {
		return fmt.Errorf("%s: %w", PortKey, err)
	}
	if want == m.addr && (m.l != nil || want == "") {
		return nil
	}
	m.Close()
	if want == "" {
		return nil
	}
	h := m.Handler
	if h == nil {
		h = Handler(nil)
	}
	l, err := m.Factory(want, h)
	if err != nil {
		return fmt.Errorf("web UI on %s: %w", want, err)
	}
	m.addr, m.l = want, l
	return nil
}

// Close stops the listener, if any.
func (m *Manager) Close() {
	if m.l != nil {
		m.l.Close()
	}
	m.l, m.addr = nil, ""
}

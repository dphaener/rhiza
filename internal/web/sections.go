package web

import (
	"bytes"
	"html"
	"html/template"
	"net/http"
	"strings"
)

// A Fragment is a live section of a page: the event kinds it listens to, and
// the template that renders just its content. A page places it with
// {{section "name" .}}; GET /fragments/<name> serves the same content alone,
// so the script can swap it when a matching event arrives.
type Fragment struct {
	// Events are the kinds of event (EventRun, ...) that make the section
	// stale. The script re-fetches the section when one arrives.
	Events []string
	// Template renders the section's content.
	Template *template.Template
	// Data builds what Template is rendered with, from the fragment request.
	// Nil means no data.
	Data func(r *http.Request) any
}

// FragmentPrefix is where a fragment's content is served.
const FragmentPrefix = "/fragments/"

// fragments is the registry of live sections, by name.
var fragments = map[string]*Fragment{}

// RegisterFragment adds a live section under name. Call it at init.
func RegisterFragment(name string, f *Fragment) { fragments[name] = f }

// FragmentURL is the URL that serves a section's content.
func FragmentURL(name string) string { return FragmentPrefix + name }

// renderFragment renders a section's content alone.
func renderFragment(f *Fragment, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := f.Template.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// section renders a page's live section: a wrapper that declares its events
// and fragment URL as data attributes, around its content.
func section(name string, data any) (template.HTML, error) {
	// A section may carry a query (name?run=...) that its fragment reads.
	base, _, _ := strings.Cut(name, "?")
	f, ok := fragments[base]
	if !ok {
		return "", &unknownSection{name}
	}
	body, err := renderFragment(f, data)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString(`<div class="rz-section" data-rz-fragment="` + html.EscapeString(FragmentURL(name)) +
		`" data-rz-events="` + html.EscapeString(strings.Join(f.Events, " ")) + `">`)
	b.Write(body)
	b.WriteString(`</div>`)
	return template.HTML(b.String()), nil
}

type unknownSection struct{ name string }

func (e *unknownSection) Error() string { return "unknown section " + e.name }

func serveFragment(w http.ResponseWriter, r *http.Request) {
	f, ok := fragments[r.PathValue("name")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	var data any
	if f.Data != nil {
		data = f.Data(r)
	}
	body, err := renderFragment(f, data)
	if err != nil {
		http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(body)
}

package web

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// searchFrom is how many projects the picker needs before it offers a
// search field: six or more.
const searchFrom = 6

var pickerTemplate = template.Must(template.ParseFS(templateFS, "templates/picker.html"))

// PickerItem is one project in the picker: where the same kind of page is in
// that project, whether it is live and how many things need the owner there.
type PickerItem struct {
	Name, Href, State string
	Live, Current     bool
	NeedsYou          int
}

// PickerView is what the picker renders: the initialised projects with
// something waiting first, the all-projects entry and its total.
type PickerView struct {
	Items    []PickerItem
	Search   bool
	AllHref  string
	AllCount int
}

// Picker serves the project picker for the page at the path query: every
// initialised project the source reads, linking the same kind of page, and an
// All projects entry for the merged inbox. It is the only reader of src on
// the pages' behalf, and reads it once per request.
func Picker(src ProjectsSource) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var projects []ProjectView
		if src != nil {
			projects = src()
		}
		current, tail := "", strings.TrimPrefix(r.URL.Query().Get("path"), "/")
		if rest, ok := strings.CutPrefix(r.URL.Query().Get("path"), ProjectPrefix); ok {
			current, tail, _ = strings.Cut(rest, "/")
			if n, err := url.PathUnescape(current); err == nil {
				current = n
			}
		}
		kind := pageKind(tail)
		v := PickerView{AllHref: InboxPath}
		for _, p := range projects {
			if !p.Initialised {
				continue
			}
			it := PickerItem{Name: p.Name, Href: ProjectPrefix + url.PathEscape(p.Name) + kind, Current: p.Name == current, NeedsYou: p.NeedsYou, State: p.Runs(), Live: p.Running+p.Waiting > 0}
			if !it.Live {
				it.State = "idle"
			}
			v.AllCount += p.NeedsYou
			v.Items = append(v.Items, it)
		}
		sort.SliceStable(v.Items, func(i, j int) bool { return (v.Items[i].NeedsYou > 0) && !(v.Items[j].NeedsYou > 0) })
		v.Search = len(v.Items) >= searchFrom
		var buf bytes.Buffer
		if err := pickerTemplate.ExecuteTemplate(&buf, "picker", v); err != nil {
			http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	})
}

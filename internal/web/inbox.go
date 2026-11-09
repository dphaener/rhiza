package web

import (
	"html/template"
	"net/http"
	"sort"
)

// InboxItemView is one thing waiting on the owner. Run and Change only build
// the link; neither is shown (constraint owner-surfaces-no-ids).
type InboxItemView struct {
	// Kind is question, signoff, failed, gave-up or stalled.
	Kind string
	// Title is what to read: a question's full text, else the change set's title.
	Title string
	// Detail is the line under the title: why it waits.
	Detail   string
	Run      string
	Change   string
	Step     string
	At       string // RFC 3339, when it began to wait
	Relative string
}

// InboxSource reads what waits on the owner as it is now.
type InboxSource func() []InboxItemView

var inboxLabels = map[string]string{
	"question": "question", "signoff": "sign-off", "failed": "failed",
	"gave-up": "gave up", "stalled": "stalled", "permission": "permission", "open": "open question",
}

// inboxRows is the items as dense rows, oldest first, each marked in the
// warning accent (the waiting pill).
func inboxRows(items []InboxItemView) []Row {
	items = append([]InboxItemView(nil), items...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].At < items[j].At })
	rows := make([]Row, 0, len(items))
	for _, it := range items {
		label := inboxLabels[it.Kind]
		if label == "" {
			label = it.Kind
		}
		href := "/changes/" + it.Change
		if it.Run != "" {
			href = "/runs/" + it.Run
		}
		rows = append(rows, Row{
			Href: href, Title: it.Title, Sub: it.Detail,
			Pill: Pill{State: "waiting", Text: label}, Step: it.Step,
			Time: it.Relative, Datetime: it.At,
		})
	}
	return rows
}

// InboxPanel is what the inbox section renders.
type InboxPanel struct {
	Rows []Row
}

type inboxKey struct{}

func inboxOf(r *http.Request) InboxPanel {
	if src, _ := r.Context().Value(inboxKey{}).(InboxSource); src != nil {
		return InboxPanel{Rows: inboxRows(src())}
	}
	return InboxPanel{}
}

func init() {
	RegisterFragment("inbox", &Fragment{
		Events:   []string{EventInbox, EventRun, EventChange},
		Template: template.Must(template.New("inbox").ParseFS(templateFS, "templates/partials.html", "templates/inbox.html")).Lookup("inbox"),
		Data:     func(r *http.Request) any { return inboxOf(r) },
	})
}

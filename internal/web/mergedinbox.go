package web

import (
	"net/http"
	"net/url"
	"sort"
)

// ProjectInbox is one project's inbox: what waits on the owner there. Name
// is the project's registered name.
type ProjectInbox struct {
	Name  string
	Items []InboxItemView
}

// MergedInboxSource reads every initialised project's inbox as it is now.
type MergedInboxSource func() []ProjectInbox

// mergedInboxRows is every project's items as one list, oldest first (ties
// by project name, then as the project listed them), each row tagged with
// its project and linking into that project's pages under /p/<project>/.
func mergedInboxRows(inboxes []ProjectInbox) []Row {
	type item struct {
		project string
		row     Row
	}
	var all []item
	for _, in := range inboxes {
		prefix := ProjectPrefix + url.PathEscape(in.Name)
		for _, row := range inboxRows(in.Items) {
			row.Project = in.Name
			row.Href = prefix + row.Href
			all = append(all, item{in.Name, row})
		}
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].row.Datetime != all[j].row.Datetime {
			return all[i].row.Datetime < all[j].row.Datetime
		}
		return all[i].project < all[j].project
	})
	rows := make([]Row, len(all))
	for i, it := range all {
		rows[i] = it.row
	}
	return rows
}

// MergedInbox serves the merged inbox page at every path it is given.
func MergedInbox(src MergedInboxSource) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rows []Row
		if src != nil {
			rows = mergedInboxRows(src())
		}
		render(w, "mergedinbox", page{Title: "Inbox", Nav: "inbox", Inbox: InboxPanel{Rows: rows}})
	})
}

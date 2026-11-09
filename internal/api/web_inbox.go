package api

import (
	"github.com/dphaener/rhiza/internal/web"
)

// webInboxKinds are the inbox items the Overview lists: what waits on the
// owner in a run. The rest (blockers, overlaps, disk and usage holds) stay
// in the CLI inbox.
var webInboxKinds = map[string]bool{
	InboxQuestion: true, InboxOpen: true, InboxSignoff: true,
	InboxFailed: true, InboxGaveUp: true, InboxStalled: true, InboxPermission: true,
}

// WebInboxSource reads the inbox for the web UI's Overview from the
// repository base works in, as rhiza inbox does, oldest first.
func WebInboxSource(base *Ctx) web.InboxSource {
	return func() []web.InboxItemView {
		c := *base
		c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
		c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
		out, err := runInbox(&c, InboxRequest{})
		if err != nil {
			return nil
		}
		return webInboxViews(out.(InboxResponse).Items)
	}
}

// webInboxViews are the inbox items in the owner's words: a question by its
// full text, the rest by their change set's title with why each waits. The
// CLI's details name handles, session and run IDs and commands; none of
// them reaches the page (constraint owner-surfaces-no-ids). An open item is
// one row per question it stopped on, each by its statement.
func webInboxViews(items []InboxItem) []web.InboxItemView {
	var views []web.InboxItemView
	for _, it := range items {
		if !webInboxKinds[it.Kind] {
			continue
		}
		view := func(title, detail string) web.InboxItemView {
			if title == "" {
				title = it.ChangeTitle
			}
			if title == "" {
				title = "a change set"
			}
			return web.InboxItemView{
				Kind: it.Kind, Title: oneLine(title), Detail: oneLine(detail),
				Run: it.Run, Change: it.Change, Step: it.Step, At: it.At,
			}
		}
		switch it.Kind {
		case InboxQuestion:
			views = append(views, view(it.Question, it.ChangeTitle))
		case InboxOpen:
			if len(it.Questions) == 0 {
				views = append(views, view("", "open questions wait on your answers"))
				continue
			}
			for _, q := range it.Questions {
				detail := it.ChangeTitle
				if q.Blocks != "" {
					detail = "blocks " + q.Blocks
				}
				views = append(views, view(q.Statement, detail))
			}
		case InboxSignoff:
			views = append(views, view("", "waiting for your sign-off"))
		case InboxFailed:
			// The failure's output names Rex windows, commits and paths:
			// the page says where it failed, and the run's page has the rest.
			detail := "the run failed"
			if it.Step != "" {
				detail = "the run failed at its " + it.Step + " step"
			}
			views = append(views, view("", detail))
		case InboxGaveUp:
			views = append(views, view("", "interrupted again after the daemon resumed it once"))
		case InboxPermission:
			// The prompt's message can hold commands and paths: the page only
			// says the run waits on a permission.
			views = append(views, view("", "a session waits for your permission"))
		case InboxStalled:
			views = append(views, view("", "its session has stopped reporting"))
		}
	}
	return views
}

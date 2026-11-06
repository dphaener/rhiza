package api

import (
	"github.com/dphaener/rhiza/internal/web"
)

// webInboxKinds are the inbox items the Overview lists: what waits on the
// owner in a run. The rest (blockers, overlaps, disk and usage holds) stay
// in the CLI inbox.
var webInboxKinds = map[string]bool{
	InboxQuestion: true, InboxOpen: true, InboxSignoff: true,
	InboxFailed: true, InboxGaveUp: true, InboxStalled: true,
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
		var items []web.InboxItemView
		for _, it := range out.(InboxResponse).Items {
			if !webInboxKinds[it.Kind] {
				continue
			}
			title, detail := it.ChangeTitle, it.Detail
			if it.Kind == InboxQuestion && it.Question != "" {
				title, detail = it.Question, it.ChangeTitle
			}
			if title == "" {
				title = it.Kind
			}
			items = append(items, web.InboxItemView{
				Kind: it.Kind, Title: title, Detail: oneLine(detail),
				Run: it.Run, Change: it.Change, Step: it.Step, At: it.At,
			})
		}
		return items
	}
}

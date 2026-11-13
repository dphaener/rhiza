package api

import (
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// needsYouStuckAfter is how long a live run's agent step may go without a
// journal entry before the run counts as stuck.
const needsYouStuckAfter = 20 * time.Minute

// needsYouKinds are the inbox kinds a project's needs-you count holds: open
// questions, change sets waiting for the owner's sign-off and failed runs.
// Stuck runs are added from the journals. The daemon's machine-wide holds
// (disk-low, usage-reserve, load-high, memory-low) belong to no project and
// are in no count (decision needs-you-count-excludes-machine-holds).
var needsYouKinds = map[string]bool{
	InboxQuestion: true, InboxOpen: true, InboxSignoff: true, InboxFailed: true,
}

// StuckKind is the kind of a needs-you item for a stuck run.
const StuckKind = "stuck"

// projectNeedsYou is the one function that counts and lists what waits on
// the owner in a project, for the header, the picker and the project's other
// surfaces: the inbox items of the kinds in needsYouKinds, each in the
// owner's words, then each live run whose current agent step (one with a
// session) has had no journal entry for needsYouStuckAfter. A run waiting on
// a block or stopped is not stuck: it already shows as a question or a
// sign-off. The count is len of the result.
func projectNeedsYou(items []InboxItem, runs []run.Journal, live func(id string) bool, now time.Time) []web.InboxItemView {
	var out []web.InboxItemView
	for _, v := range webInboxViews(items) {
		if needsYouKinds[v.Kind] {
			out = append(out, v)
		}
	}
	for _, j := range runs {
		st := j.State()
		if st.Ended || st.Current == "" || st.Blocked != nil || st.Stopped != "" || len(j.Entries) == 0 {
			continue
		}
		if st.Session[st.Current] == "" || (live != nil && !live(j.ID)) {
			continue
		}
		last := j.Entries[len(j.Entries)-1].At
		if now.Sub(last) <= needsYouStuckAfter {
			continue
		}
		title := st.Title
		if title == "" {
			title = st.Workflow
		}
		if title == "" {
			title = "a run"
		}
		out = append(out, web.InboxItemView{
			Kind: StuckKind, Title: oneLine(title), Detail: "its " + st.Current + " step has stopped making progress",
			Run: j.ID, Change: st.Change, Step: st.Current, At: last.UTC().Format(time.RFC3339),
		})
	}
	return out
}

// WebNeedsYou lists what waits on the owner in the project base works in.
func WebNeedsYou(base *Ctx, now time.Time) []web.InboxItemView {
	c := *base
	c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
	c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
	var items []InboxItem
	if out, err := runInbox(&c, InboxRequest{}); err == nil {
		items = out.(InboxResponse).Items
	}
	common, err := gitCommonDir(&c)
	if err != nil {
		return projectNeedsYou(items, nil, nil, now)
	}
	rs := run.NewStore(run.Dir(common))
	runs, _ := rs.List()
	return projectNeedsYou(items, runs, rs.Live, now)
}

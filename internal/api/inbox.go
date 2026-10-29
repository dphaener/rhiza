package api

// This file holds inbox: everything across the repository's runs that
// waits on the owner, oldest first. It reads the run journals only and
// never writes.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "inbox",
		Summary:  "List what waits on the owner across every run: questions, sign-offs, failed runs and runs the daemon gave up on, oldest first",
		Request:  reflect.TypeFor[InboxRequest](),
		Response: reflect.TypeFor[InboxResponse](),
		Run:      runInbox,
	})
}

// The kinds of an InboxItem.
const (
	InboxQuestion = "question"
	InboxSignoff  = "signoff"
	InboxFailed   = "failed"
	InboxGaveUp   = "gave-up"
)

// InboxRequest takes nothing.
type InboxRequest struct{}

// InboxItem is one thing waiting on the owner.
type InboxItem struct {
	ID       string `json:"id" desc:"The item's short ID: stable while the item waits"`
	Kind     string `json:"kind" desc:"question, signoff, failed or gave-up"`
	At       string `json:"at" desc:"When it began to wait, RFC 3339 UTC"`
	Run      string `json:"run" desc:"The run it belongs to"`
	Change   string `json:"change" desc:"The change set of the run"`
	Step     string `json:"step,omitempty" desc:"The step: the asking session's, the sign-off's, the failed step, or where the run was interrupted"`
	Question string `json:"question,omitempty" desc:"A question: what the agent asked"`
	Tab      string `json:"tab,omitempty" desc:"A question: the title of the session's tab"`
	TabBlock string `json:"tab_block,omitempty" desc:"A question: the Rex block of the session's tab"`
	Session  string `json:"session,omitempty" desc:"A question: the asking session"`
	Detail   string `json:"detail,omitempty" desc:"A failed run: why it failed; a sign-off: what it waits for"`
	Next     string `json:"next,omitempty" desc:"The command that clears it, when one does"`
}

// InboxResponse is the inbox.
type InboxResponse struct {
	Items []InboxItem `json:"items" desc:"The items waiting on the owner, oldest first; empty when none"`
}

// Human writes the inbox as text.
func (r InboxResponse) Human(w io.Writer) {
	if len(r.Items) == 0 {
		fmt.Fprintln(w, "Nothing waits on you")
		return
	}
	for _, it := range r.Items {
		fmt.Fprintf(w, "%s %s run %s, change set %s", it.ID, it.Kind, it.Run, it.Change)
		if it.Step != "" {
			fmt.Fprintf(w, ", step %s", it.Step)
		}
		fmt.Fprintln(w)
		switch it.Kind {
		case InboxQuestion:
			if it.Tab != "" {
				fmt.Fprintf(w, "  tab: %s\n", it.Tab)
			}
			fmt.Fprintf(w, "  question: %s\n", it.Question)
		default:
			if it.Detail != "" {
				fmt.Fprintf(w, "  %s\n", it.Detail)
			}
		}
		if it.Next != "" {
			fmt.Fprintf(w, "  next: %s\n", it.Next)
		}
	}
}

// inboxItems are the items of the runs, oldest first. runs are the
// repository's runs, oldest first; host and alive say which processes run
// and locked reports a run another process holds.
func inboxItems(runs []run.Journal, host string, alive func(pid int) bool, locked func(id string) bool) []InboxItem {
	items := []InboxItem{}
	newer := map[string]bool{}
	for i := len(runs) - 1; i >= 0; i-- {
		j := runs[i]
		st := j.State()
		retried := newer[st.Change]
		newer[st.Change] = true
		add := func(kind, disc string, at time.Time, it InboxItem) {
			it.Kind, it.Run, it.Change = kind, j.ID, st.Change
			it.At = at.UTC().Format(time.RFC3339)
			sum := sha256.Sum256([]byte(strings.Join([]string{j.ID, kind, disc, it.At}, "\x00")))
			it.ID = hex.EncodeToString(sum[:])[:7]
			items = append(items, it)
		}
		switch {
		case st.Ended:
			if st.Result != run.Failed || st.Abandoned() || retried {
				continue
			}
			step, _ := st.FailedStep()
			add(InboxFailed, "", st.EndedAt, InboxItem{
				Step: step, Detail: st.Detail, Next: "rhiza run resume --retry " + j.ID,
			})
		case locked(j.ID):
			continue
		case st.Stopped == run.SignoffWait && st.CurrentVisit > 0 && !hasSignoff(st):
			add(InboxSignoff, st.Current, st.StoppedAt, InboxItem{
				Step: st.Current, Detail: st.Stopped, Next: "rhiza run approve " + j.ID,
			})
		default:
			if st.Autonomous {
				continue
			}
			for _, b := range st.Blocks {
				if b.Block.Kind != "question" {
					continue
				}
				add(InboxQuestion, b.Session, b.At, InboxItem{
					Step: b.Step, Question: b.Block.Question, Tab: b.Block.TabTitle,
					TabBlock: b.Block.TabBlock, Session: b.Session,
				})
			}
			if st.DaemonResumed && st.ProcessAt.After(st.DaemonResumedAt) && st.Liveness(host, alive) == run.LiveInterrupted && len(st.Blocks) == 0 {
				add(InboxGaveUp, "", st.DaemonResumedAt, InboxItem{
					Step: st.Current, Detail: "interrupted again after the daemon's one resume",
					Next: "rhiza run resume " + j.ID,
				})
			}
		}
	}
	slices.SortStableFunc(items, func(a, b InboxItem) int {
		if c := strings.Compare(a.At, b.At); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return items
}

// hasSignoff reports whether the current visit of the current step is
// signed off already.
func hasSignoff(st run.State) bool {
	_, ok := st.SignoffFor(st.Current, st.CurrentVisit)
	return ok
}

// runInbox reads every run's journal.
func runInbox(ctx *Ctx, req any) (any, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	runs, err := rs.List()
	if err != nil {
		return nil, Errorf("read runs: %v", err)
	}
	host, _ := os.Hostname()
	return InboxResponse{Items: inboxItems(runs, host, nil, rs.Live)}, nil
}

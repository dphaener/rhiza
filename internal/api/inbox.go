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

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "inbox",
		Summary:  "List what waits on the owner across every run: questions, sign-offs, failed runs, runs the daemon gave up on and interrupted runs it skipped because their change set is no longer accepted, oldest first",
		Request:  reflect.TypeFor[InboxRequest](),
		Response: reflect.TypeFor[InboxResponse](),
		Run:      runInbox,
	})
	register(Operation{
		Name:     "inbox.approve",
		Summary:  "Approve the sign-off an inbox item waits at, as run approve does on that item's run; exit 2 on any other kind of item",
		Write:    true,
		Request:  reflect.TypeFor[InboxApproveRequest](),
		Response: reflect.TypeFor[RunSignoffResponse](),
		Run:      runInboxApprove,
	})
	register(Operation{
		Name:     "inbox.request-changes",
		Summary:  "Ask for changes to the sign-off an inbox item waits at, as run request-changes does on that item's run; exit 2 on any other kind of item",
		Write:    true,
		Request:  reflect.TypeFor[InboxRequestChangesRequest](),
		Response: reflect.TypeFor[RunSignoffResponse](),
		Input: &Input{
			Summary: `{"findings":[{"check":...,"severity":...,"leaf":...,"place":...,"reason":...},...]}, at least one; see rhiza describe run.request-changes`,
			Schema:  changesFileSchema,
		},
		Run: runInboxRequestChanges,
	})
	register(Operation{
		Name:     "inbox.open",
		Summary:  "Focus the Rex tab an inbox item waits in: the asking session's for a question, the run's journalled tab for the rest",
		Request:  reflect.TypeFor[InboxOpenRequest](),
		Response: reflect.TypeFor[InboxOpenResponse](),
		Run:      runInboxOpen,
	})
}

// The kinds of an InboxItem.
const (
	InboxQuestion = "question"
	InboxSignoff  = "signoff"
	InboxFailed   = "failed"
	InboxGaveUp   = "gave-up"
	InboxSkipped  = "skipped"
)

// InboxRequest takes nothing.
type InboxRequest struct{}

// InboxItem is one thing waiting on the owner.
type InboxItem struct {
	ID             string `json:"id" desc:"The item's short ID: stable while the item waits"`
	Kind           string `json:"kind" desc:"question, signoff, failed, gave-up or skipped"`
	At             string `json:"at" desc:"When it began to wait, RFC 3339 UTC"`
	Run            string `json:"run" desc:"The run it belongs to"`
	Change         string `json:"change" desc:"The change set of the run"`
	ChangeTitle    string `json:"change_title,omitempty" desc:"The change set's title; absent when the log has no such change set"`
	Step           string `json:"step,omitempty" desc:"The step: the asking session's, the sign-off's, the failed step, or where the run was interrupted"`
	Question       string `json:"question,omitempty" desc:"A question: what the agent asked"`
	Tab            string `json:"tab,omitempty" desc:"A question: the title of the session's tab"`
	TabBlock       string `json:"tab_block,omitempty" desc:"A question: the Rex block of the session's tab"`
	TabSession     string `json:"tab_session,omitempty" desc:"A question: the Rex session of the session's tab"`
	Session        string `json:"session,omitempty" desc:"A question: the asking session"`
	Recommendation string `json:"recommendation,omitempty" desc:"A question: what the advisor recommends answering; absent while it has not, or when it failed"`
	Reasons        string `json:"reasons,omitempty" desc:"A question: why the advisor recommends it"`
	Detail         string `json:"detail,omitempty" desc:"A failed run: why it failed; a sign-off: what it waits for"`
	Next           string `json:"next,omitempty" desc:"The command that clears it, when one does"`
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
			if it.Recommendation != "" {
				fmt.Fprintf(w, "  recommendation: %s\n", it.Recommendation)
				fmt.Fprintf(w, "  reasons: %s\n", it.Reasons)
			}
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
// locked reports a run another process holds. changeState gives the state
// of a change set and whether the log has it at all.
func inboxItems(runs []run.Journal, host string, alive func(pid int) bool, locked func(id string) bool, changeState func(change string) (string, bool)) []InboxItem {
	// gone reports a change set that is landed, withdrawn, rejected or
	// superseded: a failed or given-up run on one is left out, as its next
	// command could only fail.
	gone := func(change string) bool {
		state, ok := changeState(change)
		return ok && changeset.State(state).Terminal()
	}
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
		// A run another process holds is still going: its questions wait on
		// the owner, but it is neither at sign-off nor interrupted.
		live := locked(j.ID)
		switch {
		case st.Ended:
			if st.Result != run.Failed || st.Abandoned() || retried || gone(st.Change) {
				continue
			}
			step, _ := st.FailedStep()
			add(InboxFailed, "", st.EndedAt, InboxItem{
				Step: step, Detail: st.Detail, Next: "rhiza run resume --retry " + j.ID,
			})
		case !live && st.Stopped == run.SignoffWait && st.CurrentVisit > 0 && !hasSignoff(st):
			add(InboxSignoff, st.Current, st.StoppedAt, InboxItem{
				Step: st.Current, Detail: st.Stopped, Next: "rhiza run approve " + j.ID,
			})
		default:
			// The daemon resumes only a run whose change set is accepted: an
			// interrupted run on one no longer in flight waits for the owner
			// to abandon it, autonomous or not. A draft or proposed change
			// set's run (a plan run's) is the owner's to resume by hand.
			if !live && !st.DaemonResumed && st.Process != nil && st.Liveness(host, alive) == run.LiveInterrupted && len(st.Blocks) == 0 {
				if state, ok := changeState(st.Change); ok && changeset.State(state).Terminal() {
					add(InboxSkipped, state, st.ProcessAt, InboxItem{
						Step:   st.Current,
						Detail: "interrupted, and not resumed: its change set is " + state,
						Next:   "rhiza run abandon " + j.ID,
					})
				}
			}
			// Only questions are the orchestrator's in an autonomous run.
			if st.Autonomous {
				continue
			}
			for _, b := range st.Blocks {
				if b.Block.Kind != "question" {
					continue
				}
				it := InboxItem{
					Step: b.Step, Question: b.Block.Question, Tab: b.Block.TabTitle,
					TabBlock: b.Block.TabBlock, TabSession: b.Block.TabSession, Session: b.Session,
				}
				if a, ok := st.AdviceFor(b); ok {
					it.Recommendation, it.Reasons = a.Recommendation, a.Reasons
				}
				add(InboxQuestion, b.Session, b.At, it)
			}
			if !live && !gone(st.Change) && st.DaemonResumed && st.ProcessAt.After(st.DaemonResumedAt) && st.Liveness(host, alive) == run.LiveInterrupted && len(st.Blocks) == 0 {
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

// changeStateOf gives the state of a change set by ID, and whether the log
// has it: a made-up or unparseable ID, or a workspace that does not open,
// reports none.
func changeStateOf(ctx *Ctx) func(change string) (string, bool) {
	ws, err := openCtxWorkspace(ctx)
	return func(change string) (string, bool) {
		if err != nil {
			return "", false
		}
		cs, perr := id.Parse(change)
		if perr != nil {
			return "", false
		}
		c, ok := ws.changeSet(cs)
		return string(c.State), ok
	}
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
	items := inboxItems(runs, host, nil, rs.Live, changeStateOf(ctx))
	title := changeTitleOf(ctx)
	for i := range items {
		items[i].ChangeTitle = title(items[i].Change)
	}
	return InboxResponse{Items: items}, nil
}

// changeTitleOf gives the display title of a change set by ID, empty when
// the log has none or the ID does not parse.
func changeTitleOf(ctx *Ctx) func(change string) string {
	ws, err := openCtxWorkspace(ctx)
	return func(change string) string {
		if err != nil {
			return ""
		}
		cs, perr := id.Parse(change)
		if perr != nil {
			return ""
		}
		if c, ok := ws.changeSet(cs); ok {
			return c.Open.DisplayTitle()
		}
		return ""
	}
}

// InboxOpenRequest names the item to open.
type InboxOpenRequest struct {
	Item string `json:"item" arg:"0" desc:"The inbox item: its short ID from rhiza inbox, or a unique prefix of it"`
}

// InboxOpenResponse says which tab was focused, or that it is gone.
type InboxOpenResponse struct {
	Item    string `json:"item" desc:"The item's short ID"`
	Kind    string `json:"kind" desc:"question, signoff, failed, gave-up or skipped"`
	Run     string `json:"run" desc:"The run it belongs to"`
	Opened  bool   `json:"opened" desc:"True when the tab was focused"`
	Tab     string `json:"tab,omitempty" desc:"The Rex block of the tab, when one is recorded"`
	TabGone bool   `json:"tab_gone,omitempty" desc:"True when the tab is no longer in Rex, or none was recorded"`
}

// Human writes the response as text.
func (r InboxOpenResponse) Human(w io.Writer) {
	switch {
	case r.Opened:
		fmt.Fprintf(w, "Focused %s, the tab of %s run %s\n", r.Tab, r.Kind, r.Run)
	case r.Tab == "":
		fmt.Fprintf(w, "%s run %s: tab gone (no tab is recorded for it)\n", r.Kind, r.Run)
	default:
		fmt.Fprintf(w, "%s run %s: tab gone (Rex no longer has %s)\n", r.Kind, r.Run, r.Tab)
	}
}

// itemTab is where item's tab is: a question's session tab, else the run's
// journalled tab, scoped by the Rex session it records: Rex finds a block
// only within its session, so the owner's shell may be in another. The
// Rex session is empty when the journal has none (one written before it
// was recorded, or by a process with no REX_SESSION).
func itemTab(item InboxItem, st run.State) (rexSession, block string) {
	if item.Kind == InboxQuestion {
		return item.TabSession, item.TabBlock
	}
	if st.Process != nil {
		return st.Process.RexSession, st.Process.RexBlock
	}
	return "", ""
}

// findInboxItem is the inbox item ref names, with the state of its run: a
// short ID or a unique prefix of one. No item, or more than one, is usage.
func findInboxItem(ctx *Ctx, ref string) (InboxItem, run.State, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return InboxItem{}, run.State{}, err
	}
	rs := run.NewStore(run.Dir(common))
	runs, err := rs.List()
	if err != nil {
		return InboxItem{}, run.State{}, Errorf("read runs: %v", err)
	}
	host, _ := os.Hostname()
	var found []InboxItem
	for _, it := range inboxItems(runs, host, nil, rs.Live, changeStateOf(ctx)) {
		if strings.HasPrefix(it.ID, ref) {
			found = append(found, it)
		}
	}
	switch len(found) {
	case 0:
		return InboxItem{}, run.State{}, Usagef("no inbox item %q", ref).WithNext("rhiza inbox")
	case 1:
	default:
		return InboxItem{}, run.State{}, Usagef("%q matches %d inbox items: give more of the ID", ref, len(found)).WithNext("rhiza inbox")
	}
	var st run.State
	for _, j := range runs {
		if j.ID == found[0].Run {
			st = j.State()
		}
	}
	return found[0], st, nil
}

// runInboxOpen focuses the tab of one inbox item.
func runInboxOpen(ctx *Ctx, req any) (any, error) {
	r := req.(InboxOpenRequest)
	if r.Item == "" {
		return nil, Usagef("inbox open needs an item: see rhiza inbox").WithNext("rhiza inbox")
	}
	item, st, err := findInboxItem(ctx, r.Item)
	if err != nil {
		return nil, err
	}
	rexSession, block := itemTab(item, st)
	resp := InboxOpenResponse{Item: item.ID, Kind: item.Kind, Run: item.Run, Tab: block}
	if block == "" || ctx.Exec == nil {
		resp.TabGone = true
		return resp, nil
	}
	scope := []string{}
	if rexSession != "" {
		scope = []string{"-s", rexSession}
	}
	inspect := append(append([]string{"block", "inspect"}, scope...), "--json", block)
	if _, err := ctx.Exec("rex", inspect, "", nil); err != nil {
		resp.TabGone = true
		return resp, nil
	}
	if _, err := ctx.Exec("rex", append(append([]string{"focus"}, scope...), block), "", nil); err != nil {
		return nil, Errorf("rex could not focus %s: %v", block, err)
	}
	resp.Opened = true
	return resp, nil
}

// InboxApproveRequest names the sign-off item to approve.
type InboxApproveRequest struct {
	Item string `json:"item" arg:"0" desc:"The inbox item: its short ID from rhiza inbox, or a unique prefix of it; it must be a sign-off"`
}

// InboxRequestChangesRequest names the sign-off item and the file of the
// owner's findings.
type InboxRequestChangesRequest struct {
	Item string `json:"item" arg:"0" desc:"The inbox item: its short ID from rhiza inbox, or a unique prefix of it; it must be a sign-off"`
	File string `json:"file" arg:"1" desc:"A JSON file of the owner's findings, or - for stdin (see the input schema)"`
}

// signoffItemRun is the run of the sign-off item ref names. Any other
// kind of item is usage, naming what clears it.
func signoffItemRun(ctx *Ctx, cmd, ref string) (string, error) {
	if strings.TrimSpace(ref) == "" {
		return "", Usagef("%s needs an item: see rhiza inbox", cmd).WithNext("rhiza inbox")
	}
	item, _, err := findInboxItem(ctx, ref)
	if err != nil {
		return "", err
	}
	if item.Kind != InboxSignoff {
		e := Usagef("%s works on a sign-off item, and item %s is a %s", cmd, item.ID, item.Kind)
		if item.Next != "" {
			e = e.WithNext(item.Next)
		} else {
			e = e.WithNext("rhiza inbox open " + item.ID)
		}
		return "", e
	}
	return item.Run, nil
}

// runInboxApprove approves the run an inbox sign-off item belongs to.
func runInboxApprove(ctx *Ctx, req any) (any, error) {
	r := req.(InboxApproveRequest)
	id, err := signoffItemRun(ctx, "inbox approve", r.Item)
	if err != nil {
		return nil, err
	}
	return runRunApprove(ctx, RunApproveRequest{Run: id})
}

// runInboxRequestChanges asks for changes on the run an inbox sign-off
// item belongs to.
func runInboxRequestChanges(ctx *Ctx, req any) (any, error) {
	r := req.(InboxRequestChangesRequest)
	id, err := signoffItemRun(ctx, "inbox request-changes", r.Item)
	if err != nil {
		return nil, err
	}
	return runRunRequestChanges(ctx, RunRequestChangesRequest{File: r.File, Run: id})
}

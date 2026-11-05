package api

// This file holds inbox: everything across the repository's runs that
// waits on the owner, oldest first. It reads the run journals only and
// never writes.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/oklog/ulid/v2"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "inbox",
		ReadOnly: true,
		Summary:  "List what waits on the owner across every run: questions, sign-offs, failed runs, runs the daemon gave up on, interrupted runs it skipped because their change set is no longer accepted and stalled sessions of live runs, and the daemon's low disk space, oldest first",
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
	InboxQuestion     = "question"
	InboxOpen         = "open"
	InboxSignoff      = "signoff"
	InboxFailed       = "failed"
	InboxGaveUp       = "gave-up"
	InboxSkipped      = "skipped"
	InboxBlocker      = "blocker"
	InboxOverlap      = "overlap"
	InboxStalled      = "stalled"
	InboxDiskLow      = "disk-low"
	InboxUsageReserve = "usage-reserve"
	InboxMainStale    = "main-stale"
	InboxAccept       = "accept"
	InboxInterview    = "interview"
)

// InboxRequest takes nothing.
type InboxRequest struct{}

// InboxItem is one thing waiting on the owner.
type InboxItem struct {
	ID             string              `json:"id" desc:"The item's short ID: stable while the item waits"`
	Kind           string              `json:"kind" desc:"question, open, signoff, failed, gave-up, skipped, stalled, blocker, overlap, disk-low, usage-reserve, main-stale, accept or interview"`
	At             string              `json:"at" desc:"When it began to wait, RFC 3339 UTC"`
	Run            string              `json:"run,omitempty" desc:"The run it belongs to; absent for a blocker and an overlap"`
	Change         string              `json:"change" desc:"The change set of the run; for a blocker, the change set the question blocks; for an overlap, the first of the two (A)"`
	ChangeTitle    string              `json:"change_title,omitempty" desc:"The change set's title; absent when the log has no such change set"`
	Step           string              `json:"step,omitempty" desc:"The step: the asking session's, the sign-off's, the failed step, or where the run was interrupted"`
	Question       string              `json:"question,omitempty" desc:"A question: what the agent asked; a blocker: the question's statement"`
	Tab            string              `json:"tab,omitempty" desc:"A question or an interview: the title of the session's tab"`
	TabBlock       string              `json:"tab_block,omitempty" desc:"A question or an interview: the Rex block of the session's tab"`
	TabSession     string              `json:"tab_session,omitempty" desc:"A question or an interview: the Rex session of the session's tab"`
	Session        string              `json:"session,omitempty" desc:"A question: the asking session; an interview: the interview's session; a stalled item: the stalled session"`
	Recommendation string              `json:"recommendation,omitempty" desc:"A question or a blocker: what the advisor recommends answering; absent while it has not, or when it failed"`
	Reasons        string              `json:"reasons,omitempty" desc:"A question or a blocker: why the advisor recommends it"`
	Restatement    string              `json:"restatement,omitempty" desc:"A blocker: what is asked in plain words, as the advisor restated it; absent until it has"`
	Handle         string              `json:"handle,omitempty" desc:"A blocker: the question's handle"`
	Blocks         string              `json:"blocks,omitempty" desc:"A blocker: what the question blocks, this change set or the handle of the node"`
	BlocksTitle    string              `json:"blocks_title,omitempty" desc:"A blocker: what the question blocks by name, the node's statement or the change set's title; absent when it has none"`
	Questions      []InboxGateQuestion `json:"questions,omitempty" desc:"An open item at a plan run's gate: each open question it stopped on, in the order it names them; absent when the run has no worktree to read them from"`
	Other          string              `json:"other,omitempty" desc:"An overlap: the second change set (B)"`
	OtherTitle     string              `json:"other_title,omitempty" desc:"An overlap: the second change set's title; absent when the log has no such change set"`
	Target         string              `json:"target,omitempty" desc:"An overlap: the node or link both change sets write, its handle when it has one"`
	TargetTitle    string              `json:"target_title,omitempty" desc:"An overlap: the target node's statement; absent when it has none"`
	Detail         string              `json:"detail,omitempty" desc:"A failed run: why it failed; a sign-off: what it waits for"`
	Next           string              `json:"next,omitempty" desc:"The command that clears it, when one does"`
}

// InboxGateQuestion is one open question a plan run's gate stopped on.
type InboxGateQuestion struct {
	Handle         string `json:"handle" desc:"The question's handle, for the commands that take it"`
	Statement      string `json:"statement" desc:"The question's full statement"`
	Blocks         string `json:"blocks,omitempty" desc:"What it blocks by name: the node's statement or the change set's title; absent when it blocks nothing"`
	Recommendation string `json:"recommendation,omitempty" desc:"What the advisor recommends answering; absent until it has written one"`
	Reasons        string `json:"reasons,omitempty" desc:"Why the advisor recommends it"`
	Restatement    string `json:"restatement,omitempty" desc:"What is asked in plain words, as the advisor restated it; absent until it has"`
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
		if it.Kind == InboxBlocker {
			fmt.Fprintf(w, "%s blocker change set %s\n", it.ID, it.Change)
			if it.ChangeTitle != "" {
				fmt.Fprintf(w, "  change set: %s\n", it.ChangeTitle)
			}
			fmt.Fprintf(w, "  question: %s (%s)\n", it.Question, it.Handle)
			if it.Restatement != "" {
				fmt.Fprintf(w, "  in plain words: %s\n", it.Restatement)
			}
			if it.BlocksTitle != "" {
				fmt.Fprintf(w, "  blocks: %s (%s)\n", it.BlocksTitle, it.Blocks)
			} else {
				fmt.Fprintf(w, "  blocks: %s\n", it.Blocks)
			}
			if it.Recommendation != "" {
				fmt.Fprintf(w, "  recommendation: %s\n", it.Recommendation)
				fmt.Fprintf(w, "  reasons: %s\n", it.Reasons)
			}
			if it.Next != "" {
				fmt.Fprintf(w, "  next: %s\n", it.Next)
			}
			continue
		}
		if it.Kind == InboxDiskLow {
			fmt.Fprintf(w, "%s disk-low\n", it.ID)
			fmt.Fprintf(w, "  %s\n", it.Detail)
			continue
		}
		if it.Kind == InboxUsageReserve {
			fmt.Fprintf(w, "%s usage-reserve\n", it.ID)
			fmt.Fprintf(w, "  %s\n", it.Detail)
			continue
		}
		if it.Kind == InboxMainStale {
			fmt.Fprintf(w, "%s main-stale change set %s\n", it.ID, it.Change)
			fmt.Fprintf(w, "  %s\n", it.Detail)
			continue
		}
		if it.Kind == InboxOverlap {
			fmt.Fprintf(w, "%s overlap change sets %s and %s\n", it.ID, it.Change, it.Other)
			fmt.Fprintf(w, "  target: %s\n", it.Target)
			fmt.Fprintf(w, "  %s\n", it.Detail)
			if it.Next != "" {
				fmt.Fprintf(w, "  next: %s\n", it.Next)
			}
			continue
		}
		fmt.Fprintf(w, "%s %s run %s, change set %s", it.ID, it.Kind, it.Run, it.Change)
		if it.Step != "" {
			fmt.Fprintf(w, ", step %s", it.Step)
		}
		fmt.Fprintln(w)
		switch it.Kind {
		case InboxInterview:
			if it.Tab != "" {
				fmt.Fprintf(w, "  tab: %s\n", it.Tab)
			}
			fmt.Fprintf(w, "  %s\n", it.Detail)
		case InboxOpen:
			first, _, _ := strings.Cut(it.Detail, "\n")
			fmt.Fprintf(w, "  %s\n", first)
			for _, q := range it.Questions {
				fmt.Fprintf(w, "  question: %s (%s)\n", q.Statement, q.Handle)
				if q.Restatement != "" {
					fmt.Fprintf(w, "    in plain words: %s\n", q.Restatement)
				}
				if q.Blocks != "" {
					fmt.Fprintf(w, "    blocks: %s\n", q.Blocks)
				}
				if q.Recommendation != "" {
					fmt.Fprintf(w, "    recommendation: %s\n", q.Recommendation)
					fmt.Fprintf(w, "    reasons: %s\n", q.Reasons)
				}
			}
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
			first, rest, _ := strings.Cut(it.Detail, "\n")
			if first != "" {
				fmt.Fprintf(w, "  %s\n", first)
			}
			for _, l := range strings.Split(rest, "\n") {
				if l != "" {
					fmt.Fprintf(w, "    %s\n", l)
				}
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
//
// tabOf gives the Rex tab of a session, nil for none. A live milestone
// request's plan run, not autonomous, is an interview item from the
// interview session's launch until the run ends.
//
// stalled gives the stalled session of a run, nil when it has none; a nil
// stalled lists none. Only a run another process holds can have one.
func inboxItems(runs []run.Journal, host string, alive func(pid int) bool, locked func(id string) bool, changeState func(change string) (string, bool), stalled func(st run.State) []RunStalled, tabOf func(sid string) session.Tab) []InboxItem {
	// gone reports a change set that is landed, withdrawn, rejected or
	// superseded: a failed or given-up run on one is left out, as its next
	// command could only fail.
	// A landed change set is not gone for a plan run, nor for a non-plan run
	// that entered its land step: the daemon resumes such a run, so it is
	// listed like any other run.
	gone := func(st run.State) bool {
		state, ok := changeState(st.Change)
		if ok && state == string(changeset.StateLanded) && (st.Workflow == planWorkflow || st.Visits["land"] > 0) {
			return false
		}
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
		// A plan run's gate, owner-run or autonomous, left these proposed
		// (the accept setting is owner, or the review did not pass): each
		// waits for acceptance.
		if !retried {
			for _, g := range st.Gated {
				if state, ok := changeState(g); ok && state == string(changeset.StateProposed) {
					add(InboxAccept, g, st.GatedAt, InboxItem{Detail: "waiting for acceptance", Next: ownerAcceptNext(g)})
					items[len(items)-1].Change = g
				}
			}
		}
		switch {
		case st.Ended:
			if st.Result != run.Failed || st.Abandoned() || retried || gone(st) {
				continue
			}
			step, _ := st.FailedStep()
			add(InboxFailed, "", st.EndedAt, InboxItem{
				Step: step, Detail: failedDetail(st, step), Next: "rhiza run resume --retry " + j.ID,
			})
		case !live && strings.HasPrefix(st.Stopped, run.OpenStop):
			add(InboxOpen, st.Stopped, st.StoppedAt, InboxItem{
				Step: st.Current, Detail: st.Stopped, Next: "rhiza run resume " + j.ID,
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
				if state, ok := changeState(st.Change); ok && gone(st) {
					add(InboxSkipped, state, st.ProcessAt, InboxItem{
						Step:   st.Current,
						Detail: "interrupted, and not resumed: its change set is " + state,
						Next:   "rhiza run abandon " + j.ID,
					})
				}
			}
			if live && stalled != nil {
				for _, sl := range stalled(st) {
					at, _ := time.Parse(time.RFC3339, sl.Since)
					add(InboxStalled, sl.Session, at, InboxItem{
						Step: sl.Step, Session: sl.Session,
						Detail: "session " + sl.Session + " has not reported since " + sl.Since,
						Next:   "rhiza session message " + sl.Session + " <text>",
					})
				}
			}
			if sid := st.Session[milestoneInterviewStep]; live && st.MilestonePlan && !st.Autonomous && sid != "" {
				it := InboxItem{
					Step: milestoneInterviewStep, Session: sid,
					Detail: "the milestone interview waits on you",
				}
				if tabOf != nil {
					t := tabOf(sid)
					it.Tab, it.TabBlock, it.TabSession = t.Title, t.Block, t.Session
				}
				add(InboxInterview, sid, st.StartedAt, it)
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
			if !live && !gone(st) && st.DaemonResumed && st.ProcessAt.After(st.DaemonResumedAt) && st.Liveness(host, alive) == run.LiveInterrupted && len(st.Blocks) == 0 {
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

// inboxStalledOf reads each run's session records at the clock's now, for
// every stalled session of the run's current step. The inbox asks only of
// live runs. Without a clock no run is stalled.
func inboxStalledOf(ctx *Ctx, common string) func(st run.State) []RunStalled {
	if ctx == nil || ctx.Now == nil {
		return nil
	}
	sessions := session.NewStore(session.Dir(common))
	now := ctx.Now()
	return func(st run.State) []RunStalled { return runStalledAll(st, true, sessions, now) }
}

// inboxTabOf reads a session's Rex tab from its record; the zero tab when
// there is none.
func inboxTabOf(ctx *Ctx, common string) func(sid string) session.Tab {
	sessions := session.NewStore(session.Dir(common))
	return func(sid string) session.Tab {
		rec, err := sessions.Read(sid)
		if err != nil {
			return session.Tab{}
		}
		return rec.Tab
	}
}

// failedDetail is why a failed run failed, followed by the detail of the
// last fail the step ended with when it has more than one line: a command
// step's exit code and the tail of its output.
func failedDetail(st run.State, step string) string {
	for i := len(st.Steps) - 1; i >= 0; i-- {
		s := st.Steps[i]
		if s.Step != step || s.Outcome != workflow.OutcomeFail {
			continue
		}
		if strings.Contains(s.Detail, "\n") {
			return strings.TrimSuffix(st.Detail, "\n") + "\n" + s.Detail
		}
		break
	}
	return st.Detail
}

// changeStateOf gives the state of a change set by ID, and whether the log
// has it: a made-up or unparseable ID, or a workspace that does not open,
// reports none. It is for one lookup: the index it opens is closed after.
func changeStateOf(ctx *Ctx) func(change string) (string, bool) {
	f := newInboxFold(ctx)
	return func(change string) (string, bool) {
		defer f.close()
		return f.stateOf(change)
	}
}

// inboxFold is one read of the change sets, shared by every lookup of one
// inbox call: a change set's state and title, and the overlaps. Reading the
// log decodes every entry, so a log read per lookup made rhiza inbox cost
// seconds. With an index (ctx.IndexPath) the states come from its change
// set rows and a change set is folded from its own entries only when a
// lookup needs it; without one the log is read and folded once. Either
// way nothing is read until the first lookup. close releases the index.
type inboxFold struct {
	ctx  *Ctx
	done bool
	err  error
	// x is the index, nil when the log was read instead.
	x      *index.Index
	states map[id.ID]changeset.State
	sets   map[id.ID]changeset.ChangeSet
	// list is every change set folded, from the log; nil with an index.
	list []changeset.ChangeSet
}

// newInboxFold returns the fold of ctx's change sets; nothing is read
// until a lookup asks.
func newInboxFold(ctx *Ctx) *inboxFold { return &inboxFold{ctx: ctx} }

// close closes the index, if one was opened.
func (f *inboxFold) close() {
	if f.x != nil {
		f.x.Close()
		f.x = nil
	}
}

// load reads the change sets the first time it is called.
func (f *inboxFold) load() error {
	if f.done {
		return f.err
	}
	f.done = true
	f.sets = map[id.ID]changeset.ChangeSet{}
	f.states = map[id.ID]changeset.State{}
	ws, err := openCtxWorkspace(f.ctx)
	if err != nil {
		f.err = err
		return err
	}
	if f.ctx.IndexPath != "" && len(ws.pending) == 0 {
		x, err := index.Open(f.ctx.Root, f.ctx.IndexPath)
		if err != nil {
			f.err = Errorf("%s", err.Error())
			return f.err
		}
		f.x = x
		infos, err := x.ChangeSets()
		if err != nil {
			f.err = Errorf("%s", err.Error())
			return f.err
		}
		for _, c := range infos {
			f.states[c.ID] = c.State
		}
		return nil
	}
	f.list, _ = changeset.FoldLog(ws.log())
	for _, c := range f.list {
		f.sets[c.ID] = c
		f.states[c.ID] = c.State
	}
	return nil
}

// changeSet is the folded change set of change, false when the ID does not
// parse, the log has none or the change sets cannot be read.
func (f *inboxFold) changeSet(change string) (changeset.ChangeSet, bool) {
	cs, err := id.Parse(change)
	if err != nil || f.load() != nil {
		return changeset.ChangeSet{}, false
	}
	if c, ok := f.sets[cs]; ok {
		return c, true
	}
	if f.x == nil {
		return changeset.ChangeSet{}, false
	}
	if _, ok := f.states[cs]; !ok {
		return changeset.ChangeSet{}, false
	}
	es, err := f.x.Entries(cs)
	if err != nil || len(es) == 0 {
		return changeset.ChangeSet{}, false
	}
	c, _ := changeset.Fold(cs, es)
	f.sets[cs] = c
	return c, true
}

// stateOf gives the state of a change set by ID, and whether the log has it.
func (f *inboxFold) stateOf(change string) (string, bool) {
	cs, err := id.Parse(change)
	if err != nil || f.load() != nil {
		return "", false
	}
	st, ok := f.states[cs]
	return string(st), ok
}

// titleOf gives the display title of a change set by ID, empty when the log
// has none or the ID does not parse.
func (f *inboxFold) titleOf(change string) string {
	if c, ok := f.changeSet(change); ok {
		return c.Open.DisplayTitle()
	}
	return ""
}

// overlaps is one overlap item per overlap between in-flight change sets
// (model.Overlaps), with no run. It began to wait when the later of the
// two change sets was opened, which its ID dates. With an index only the
// in-flight change sets are folded, and targets are named from its as-is
// view, as the landed model names them.
func (f *inboxFold) overlaps() ([]InboxItem, error) {
	if err := f.load(); err != nil {
		return nil, err
	}
	if f.x == nil {
		ovs := gmodel.Overlaps(f.list)
		if len(ovs) == 0 {
			return []InboxItem{}, nil
		}
		return overlapItems(ovs, newModelSets(f.list, id.ID{})), nil
	}
	// EachInFlight calls back from several goroutines at once.
	var mu sync.Mutex
	var flight []changeset.ChangeSet
	if err := f.x.EachInFlight(func(cs id.ID, es []entry.Entry) {
		c, _ := changeset.Fold(cs, es)
		mu.Lock()
		defer mu.Unlock()
		f.sets[cs] = c
		flight = append(flight, c)
	}); err != nil {
		return nil, Errorf("%s", err.Error())
	}
	// In ID order, as the log's fold lists them.
	slices.SortFunc(flight, func(a, b changeset.ChangeSet) int { return id.Compare(a.ID, b.ID) })
	ovs := gmodel.Overlaps(flight)
	if len(ovs) == 0 {
		return []InboxItem{}, nil
	}
	return overlapItemsNamed(ovs, func(t id.ID) (string, string) {
		if n, ok, err := f.x.Node(index.AsIs, t); err == nil && ok {
			return n.Handle, fieldStatement(n.Fields)
		}
		return "", ""
	}), nil
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
	fold := newInboxFold(ctx)
	defer fold.close()
	items := inboxItems(runs, host, nil, rs.Live, fold.stateOf, inboxStalledOf(ctx, common), inboxTabOf(ctx, common))
	blockers, err := inboxBlockers(ctx, questionAdvice(runs))
	if err != nil {
		return nil, err
	}
	overlaps, err := fold.overlaps()
	if err != nil {
		return nil, err
	}
	items = sortInboxItems(append(append(append(append(append(items, blockers...), overlaps...), inboxDiskLow(ctx)...), inboxUsageHold(ctx)...), inboxMainStale(ctx)...))
	byRun := map[string]run.State{}
	for _, j := range runs {
		byRun[j.ID] = j.State()
	}
	for i := range items {
		items[i].ChangeTitle = fold.titleOf(items[i].Change)
		if items[i].Other != "" {
			items[i].OtherTitle = fold.titleOf(items[i].Other)
		}
		switch {
		case items[i].Kind == InboxBlocker && items[i].Blocks == "this change set":
			items[i].BlocksTitle = items[i].ChangeTitle
		case items[i].Kind == InboxOpen:
			items[i].Questions = inboxGateQuestions(ctx, byRun[items[i].Run])
		}
	}
	return InboxResponse{Items: items}, nil
}

// sortInboxItems orders items oldest first, ties by ID.
func sortInboxItems(items []InboxItem) []InboxItem {
	slices.SortStableFunc(items, func(a, b InboxItem) int {
		if c := strings.Compare(a.At, b.At); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return items
}

// inboxBlockers is one blocker item per open question and proposed or
// accepted change set it blocks, read from the to-be view as ready reads
// it: a question on a draft or landed change set is never listed. A
// blocker began to wait when its question was raised, which its ID dates.
// advice gives the advisor's restatement and recommendation on a question
// by its node ID (questionAdvice), for a question a run raised from a
// template: a failed leaf, a review cycle cap.
func inboxBlockers(ctx *Ctx, advice map[string]run.Advice) ([]InboxItem, error) {
	r, err := openReader(ctx, "inbox", "", toBeView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	in, err := planGather(r)
	if err != nil {
		return nil, err
	}
	res := planCore(in)
	var qs []id.ID
	for _, p := range res.order {
		if p.info.State == changeset.StateDraft {
			continue
		}
		for _, b := range p.blockers {
			qs = append(qs, b.question)
			if _, ok := in.nodes[b.target]; ok {
				qs = append(qs, b.target)
			}
		}
	}
	statements := map[id.ID]string{}
	if len(qs) > 0 {
		ns, err := r.NodesFields(planSortIDs(qs))
		if err != nil {
			return nil, err
		}
		for _, n := range ns {
			var st string
			_ = json.Unmarshal(n.Fields["statement"], &st)
			statements[n.ID] = st
		}
	}
	items := []InboxItem{}
	for _, p := range res.order {
		if p.info.State == changeset.StateDraft {
			continue
		}
		seen := map[id.ID]bool{}
		for _, b := range p.blockers {
			if seen[b.question] {
				continue
			}
			seen[b.question] = true
			handle := planHandle(b.question, in.nodes)
			blocks, blocksTitle := "this change set", ""
			if b.target != p.info.ID {
				blocks = planHandle(b.target, in.nodes)
				blocksTitle = statements[b.target]
			}
			text := statements[b.question]
			if text == "" {
				text = handle
			}
			at := ulid.Time(ulid.ULID(b.question).Time()).UTC().Format(time.RFC3339)
			sum := sha256.Sum256([]byte(strings.Join([]string{InboxBlocker, b.question.String(), p.info.ID.String()}, "\x00")))
			it := InboxItem{
				ID: hex.EncodeToString(sum[:])[:7], Kind: InboxBlocker, At: at, Change: p.info.ID.String(),
				Question: text, Handle: handle, Blocks: blocks, BlocksTitle: blocksTitle, Next: "rhiza show " + handle,
			}
			if a, ok := advice[b.question.String()]; ok {
				it.Recommendation, it.Reasons, it.Restatement = a.Recommendation, a.Reasons, a.Restatement
			}
			items = append(items, it)
		}
	}
	return items, nil
}

// questionAdvice is the latest advice with a recommendation on each
// question node, by the node's ID (the Asker of an advisor the engine
// started on a question it raised), across every run.
func questionAdvice(runs []run.Journal) map[string]run.Advice {
	out := map[string]run.Advice{}
	for _, j := range runs {
		for _, a := range j.State().Advice {
			if a.Recommendation == "" {
				continue
			}
			if old, ok := out[a.Asker]; !ok || !a.At.Before(old.At) {
				out[a.Asker] = a
			}
		}
	}
	return out
}

// overlapItems are the overlap items of ovs, targets named by m.
func overlapItems(ovs []gmodel.Overlap, m Model) []InboxItem {
	return overlapItemsNamed(ovs, func(t id.ID) (string, string) {
		if n, ok := m.Node(t); ok {
			return n.Handle, fieldStatement(n.Fields)
		}
		return "", ""
	})
}

// fieldStatement is the statement field of a node's fields, "" for none.
func fieldStatement(fs map[string]json.RawMessage) string {
	var st string
	_ = json.Unmarshal(fs["statement"], &st)
	return st
}

// overlapItemsNamed are the overlap items of ovs; name gives a target's
// handle and statement, "" for none: the target's ID stands in for no
// handle.
func overlapItemsNamed(ovs []gmodel.Overlap, name func(id.ID) (string, string)) []InboxItem {
	items := make([]InboxItem, 0, len(ovs))
	for _, ov := range ovs {
		target := ov.Target.String()
		h, statement := name(ov.Target)
		if h != "" {
			target = h
		}
		later := max(ulid.ULID(ov.A).Time(), ulid.ULID(ov.B).Time())
		at := ulid.Time(later).UTC().Format(time.RFC3339)
		sum := sha256.Sum256([]byte(strings.Join([]string{InboxOverlap, ov.A.String(), ov.B.String(), ov.Target.String()}, "\x00")))
		items = append(items, InboxItem{
			ID: hex.EncodeToString(sum[:])[:7], Kind: InboxOverlap, At: at,
			Change: ov.A.String(), Other: ov.B.String(), Target: target, TargetTitle: statement,
			Detail: "both write " + strings.Join(ov.Aspects, ", "),
			Next:   "rhiza slice --change " + ov.A.String() + " --view what-if",
		})
	}
	return items
}

// inboxDiskLow is one disk-low item while the daemon's last poll found free
// space under its minimum, and none once a later poll found enough. It began to wait at the first poll that found it low.
func inboxDiskLow(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	st, ok := daemon.ReadState(common)
	if !ok {
		return nil
	}
	return diskLowItems(st.DiskLow)
}

// diskLowItems is one disk-low item for low, none when low is nil.
func diskLowItems(low *daemon.DiskLow) []InboxItem {
	if low == nil {
		return nil
	}
	at := low.Since.UTC().Format(time.RFC3339)
	sum := sha256.Sum256([]byte(strings.Join([]string{InboxDiskLow, at}, "\x00")))
	return []InboxItem{{
		ID: hex.EncodeToString(sum[:])[:7], Kind: InboxDiskLow, At: at,
		Detail: fmt.Sprintf("%s free, under the minimum of %s: the daemon holds new starts and resumes", daemon.FormatSize(low.Free), daemon.FormatSize(low.Min)),
	}}
}

// inboxUsageHold is one usage-reserve item while the daemon's last poll held
// its launches for usage, and none once a later poll did not.
func inboxUsageHold(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	st, ok := daemon.ReadState(common)
	if !ok {
		return nil
	}
	return usageHoldItems(st.UsageHold)
}

// usageHoldItems is one usage-reserve item for h, none when h is nil.
func usageHoldItems(h *daemon.UsageHold) []InboxItem {
	if h == nil {
		return nil
	}
	at := h.Since.UTC().Format(time.RFC3339)
	sum := sha256.Sum256([]byte(strings.Join([]string{InboxUsageReserve, at}, "\x00")))
	return []InboxItem{{
		ID: hex.EncodeToString(sum[:])[:7], Kind: InboxUsageReserve, At: at,
		Detail: usageHoldDetail(*h),
	}}
}

// inboxMainStale is one main-stale item while the last ship close-out could
// not fast-forward the main checkout, and none once a later one did.
func inboxMainStale(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	return mainStaleItems(common)
}

// mainStaleItems is the main-stale item the record under common makes.
func mainStaleItems(common string) []InboxItem {
	ms, ok := daemon.ReadMainStale(common)
	if !ok {
		return nil
	}
	at := ms.Since.UTC().Format(time.RFC3339)
	sum := sha256.Sum256([]byte(strings.Join([]string{InboxMainStale, at}, "\x00")))
	return []InboxItem{{
		ID: hex.EncodeToString(sum[:])[:7], Kind: InboxMainStale, At: at, Change: ms.Change,
		Detail: ms.Detail,
	}}
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
	if item.Kind == InboxQuestion || item.Kind == InboxInterview {
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
	fold := newInboxFold(ctx)
	defer fold.close()
	all := inboxItems(runs, host, nil, rs.Live, fold.stateOf, inboxStalledOf(ctx, common), inboxTabOf(ctx, common))
	blockers, err := inboxBlockers(ctx, questionAdvice(runs))
	if err != nil {
		return InboxItem{}, run.State{}, err
	}
	overlaps, err := fold.overlaps()
	if err != nil {
		return InboxItem{}, run.State{}, err
	}
	for _, it := range append(append(append(append(append(all, blockers...), overlaps...), inboxDiskLow(ctx)...), inboxUsageHold(ctx)...), inboxMainStale(ctx)...) {
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
	if item.Kind == InboxDiskLow {
		return nil, Usagef("item %s is the daemon's low disk space: it has no tab to open", item.ID).WithNext("rhiza daemon status")
	}
	if item.Kind == InboxUsageReserve {
		return nil, Usagef("item %s is the daemon's usage hold: it has no tab to open", item.ID).WithNext("rhiza daemon status")
	}
	if item.Kind == InboxMainStale {
		return nil, Usagef("item %s is the main checkout not being up to date: it has no tab to open", item.ID).WithNext("rhiza health")
	}
	if item.Kind == InboxOverlap {
		return nil, Usagef("item %s is an overlap between change sets %s and %s: it has no tab to open", item.ID, item.Change, item.Other).WithNext(item.Next)
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
	Item       string `json:"item" arg:"0" desc:"The inbox item: its short ID from rhiza inbox, or a unique prefix of it; it must be a sign-off"`
	Foreground bool   `json:"foreground,omitempty" desc:"Drive the run in this process, rather than hand it to the daemon. The daemon passes it to the detached process it starts"`
}

// InboxRequestChangesRequest names the sign-off item and the file of the
// owner's findings.
type InboxRequestChangesRequest struct {
	Item       string `json:"item" arg:"0" desc:"The inbox item: its short ID from rhiza inbox, or a unique prefix of it; it must be a sign-off"`
	File       string `json:"file" arg:"1" desc:"A JSON file of the owner's findings, or - for stdin (see the input schema)"`
	Foreground bool   `json:"foreground,omitempty" desc:"Drive the run in this process, rather than hand it to the daemon. The daemon passes it to the detached process it starts"`
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
	return runRunApprove(ctx, RunApproveRequest{Run: id, Foreground: r.Foreground})
}

// runInboxRequestChanges asks for changes on the run an inbox sign-off
// item belongs to.
func runInboxRequestChanges(ctx *Ctx, req any) (any, error) {
	r := req.(InboxRequestChangesRequest)
	id, err := signoffItemRun(ctx, "inbox request-changes", r.Item)
	if err != nil {
		return nil, err
	}
	return runRunRequestChanges(ctx, RunRequestChangesRequest{File: r.File, Run: id, Foreground: r.Foreground})
}

// inboxGateQuestions are the open questions a plan run stopped on at its
// gate (run.OpenStop), read from the run's worktree as answering one
// reads them (gateAnswerCtx): each question's full statement and what it
// blocks by name, with the advisor's recommendation once it has
// recommended on the question by handle or ID. It is nil when the run is
// not stopped at its gate, has no worktree, or the worktree does not read.
func inboxGateQuestions(ctx *Ctx, st run.State) []InboxGateQuestion {
	handles := gateStopQuestions(st.Stopped)
	if len(handles) == 0 || st.Worktree == "" {
		return nil
	}
	c := *ctx
	c.Root, c.IndexPath = st.Worktree, ""
	w, err := openCtxWorkspace(&c)
	if err != nil {
		return nil
	}
	own, err := id.Parse(st.Change)
	if err != nil {
		return nil
	}
	sets := w.changeSets()
	plan, err := planChangeSets(sets, w.model(own), own)
	if err != nil {
		return nil
	}
	text := func(m Model, n id.ID) string {
		if x, ok := m.Node(n); ok {
			var s string
			_ = json.Unmarshal(x.Fields["statement"], &s)
			return s
		}
		for _, cs := range sets {
			if cs.ID == n {
				return cs.Open.DisplayTitle()
			}
		}
		return ""
	}
	out := make([]InboxGateQuestion, 0, len(handles))
	for _, h := range handles {
		q := InboxGateQuestion{Handle: h}
		asker := h
		for _, cs := range plan {
			m := w.model(cs.ID)
			n, err := m.Resolve(h)
			if err != nil || n.Kind != op.KindQuestion {
				continue
			}
			_ = json.Unmarshal(n.Fields["statement"], &q.Statement)
			asker = n.ID.String()
			if ls := m.LinksFrom(n.ID, op.LinkBlocks); len(ls) > 0 {
				q.Blocks = text(m, ls[0].To)
			}
			break
		}
		for i := len(st.Advice) - 1; i >= 0; i-- {
			if a := st.Advice[i]; (a.Asker == h || a.Asker == asker) && a.Recommendation != "" {
				q.Recommendation, q.Reasons, q.Restatement = a.Recommendation, a.Reasons, a.Restatement
				break
			}
		}
		out = append(out, q)
	}
	return out
}

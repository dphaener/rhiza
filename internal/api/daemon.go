package api

// This file holds the daemon operation: rhiza daemon runs in the
// foreground, one per repository. It holds a lock under the git common dir,
// so a second daemon started from any worktree of the repository is
// refused, naming the first one's process. It binds no socket.

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "daemon",
		Summary:  "Run the daemon in the foreground: one per repository, held by a lock in the git common dir",
		Request:  reflect.TypeFor[DaemonRequest](),
		Response: reflect.TypeFor[DaemonResponse](),
		Run:      runDaemon,
	})
}

// daemonPoll is how often the daemon looks at the roadmap and the run
// journals.
const daemonPoll = 5 * time.Second

// DaemonRequest starts the daemon.
type DaemonRequest struct {
	MaxRuns        int    `json:"max_runs,omitempty" desc:"Start runs until this many are running; default 2. An interrupted run takes no slot"`
	Once           bool   `json:"once,omitempty" desc:"Poll once, start what is picked and stop"`
	Autonomous     bool   `json:"autonomous,omitempty" desc:"Start each run with --autonomous: it goes on without the owner"`
	MinFree        string `json:"min_free,omitempty" desc:"Start and resume nothing while the repo's filesystem has less free space than this: bytes, or a number with K, M, G or T (powers of 1024); default 10G; 0 turns the guard off"`
	PermissionMode string `json:"permission_mode,omitempty" desc:"Start each run with this --permission-mode, for its agent steps that set none; checked against ship-change before the daemon starts"`
	EventCommand   string `json:"event_command,omitempty" desc:"Run this command through sh -c in the repo root once per change set lifecycle move (proposed, accepted, rejected, withdrawn, superseded, landed, reopened) it has not handled, in log order, with the move as JSON on stdin and RHIZA_EVENT, RHIZA_EVENT_CHANGE, RHIZA_EVENT_ENTRY and RHIZA_EVENT_TITLE set; default: git config rhiza.eventCommand. The first run only records the moves already in the log. A failing or timed-out (10 s) command is logged and not retried"`
	NotifyCommand  string `json:"notify_command,omitempty" desc:"Run this command through sh -c in the repo root once per new inbox item, with the item's JSON on stdin and RHIZA_INBOX_KIND, RHIZA_INBOX_ID, RHIZA_INBOX_CHANGE and RHIZA_INBOX_TITLE set; default: git config rhiza.notifyCommand. A failing or timed-out (10 s) command is logged and not retried"`
}

// DaemonResponse is the daemon after it stopped.
type DaemonResponse struct {
	Process string   `json:"process" desc:"The daemon's process, as host:pid"`
	Started []string `json:"started" desc:"The change sets whose runs the daemon started, in order"`
	Resumed []string `json:"resumed" desc:"The interrupted runs the daemon resumed, in order"`
}

// Human writes the response as markdown.
func (r DaemonResponse) Human(w io.Writer) {
	for _, cs := range r.Started {
		fmt.Fprintf(w, "Started a run of %s\n", cs)
	}
	for _, id := range r.Resumed {
		fmt.Fprintf(w, "Resumed run %s\n", id)
	}
	fmt.Fprintf(w, "Daemon %s stopped\n", r.Process)
}

// runDaemon takes the repository's daemon lock and keeps it until the
// caller interrupts (Ctx.Interrupt). A second daemon exits with a usage
// error naming the holder.
func runDaemon(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("daemon runs in the foreground and has no dry run")
	}
	r := req.(DaemonRequest)
	if r.MaxRuns < 0 {
		return nil, Usagef("--max-runs must be at least 1")
	}
	if r.MaxRuns == 0 {
		r.MaxRuns = daemon.DefaultMaxRuns
	}
	minFree := daemon.DefaultMinFree
	if r.MinFree != "" {
		v, err := daemon.ParseSize(r.MinFree)
		if err != nil {
			return nil, Usagef("--min-free: %v", err)
		}
		minFree = v
	}
	// Every run the daemon starts is a ship-change run with this mode: one
	// that run start would refuse is refused here, before any tab opens.
	_, def, _, err := resolveWorkflow(ctx.Root, "")
	if err != nil {
		return nil, err
	}
	if err := checkAgentSteps(def, r.PermissionMode, "", r.Autonomous, "rhiza describe daemon"); err != nil {
		return nil, err
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	lock, err := daemon.TryLock(common, ctx.Now())
	if err != nil {
		if h, ok := err.(*daemon.Held); ok {
			return nil, Usagef("%s", h.Error()).WithNext("rhiza daemon status")
		}
		return nil, err
	}
	defer func() { lock.Unlock() }()
	self := daemon.ReadHolder(common)
	if ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon running, process %s\n", self)
	}
	notifyCmd := strings.TrimSpace(r.NotifyCommand)
	if notifyCmd == "" && ctx.Git != nil {
		if out, err := ctx.Git("config", "--get", "rhiza.notifyCommand"); err == nil {
			notifyCmd = strings.TrimSpace(out)
		}
	}
	eventCmd := strings.TrimSpace(r.EventCommand)
	if eventCmd == "" && ctx.Git != nil {
		if out, err := ctx.Git("config", "--get", "rhiza.eventCommand"); err == nil {
			eventCmd = strings.TrimSpace(out)
		}
	}
	rs := run.NewStore(run.Dir(common))
	resp := DaemonResponse{Process: self.String(), Started: []string{}, Resumed: []string{}}
	pending := map[string]daemonPending{}
	unjournaled := map[string]daemonPending{}
	state := daemon.State{Process: self, Version: ctx.Version, MaxRuns: r.MaxRuns, Autonomous: r.Autonomous, Runs: []daemon.LiveRun{}}
	if prev, ok := daemon.ReadState(common); ok {
		state.LastStarted, state.LastResumed = prev.LastStarted, prev.LastResumed
		state.Notified = prev.Notified
		state.Hooked = prev.Hooked
		state.Events, state.EventsSeeded = prev.Events, prev.EventsSeeded
		state.DiskLow = prev.DiskLow
		// A re-exec, or a restart, keeps the holds of the daemon before.
		unjournaled = daemonRestoreUnjournaled(prev.Unjournaled)
	}
	saveState := func() {
		state.LastPoll = ctx.Now().UTC()
		if runs, err := rs.List(); err == nil {
			host, _ := os.Hostname()
			state.Runs = daemon.LiveRuns(runs, host, nil)
		}
		daemon.WriteState(common, state)
	}
	defer func() {
		state.Stopped = ctx.Now().UTC()
		daemon.WriteState(common, state)
	}()
	poll := func() error {
		var skips []daemon.Skip
		started, resumed, err := daemonPoll1(ctx, rs, pending, unjournaled, r.MaxRuns, minFree, r.Autonomous, r.PermissionMode, &skips)
		if err == nil {
			state.Skips = skips
			if low, free := daemonDiskLow(ctx, minFree); !low {
				state.DiskLow = nil
			} else if state.DiskLow == nil {
				state.DiskLow = &daemon.DiskLow{Since: ctx.Now().UTC(), Free: free, Min: minFree}
			} else {
				state.DiskLow.Free, state.DiskLow.Min = free, minFree
			}
		}
		resp.Started = append(resp.Started, started...)
		resp.Resumed = append(resp.Resumed, resumed...)
		if len(started) > 0 {
			state.LastStarted = &daemon.Event{Change: started[len(started)-1], At: ctx.Now().UTC()}
		}
		if len(resumed) > 0 {
			state.LastResumed = &daemon.Event{Run: resumed[len(resumed)-1], At: ctx.Now().UTC()}
		}
		state.Unjournaled = []daemon.Event{}
		for cs, p := range unjournaled {
			state.Unjournaled = append(state.Unjournaled, daemon.Event{Change: cs, At: p.at.UTC(), Runs: p.runs})
		}
		slices.SortFunc(state.Unjournaled, func(a, b daemon.Event) int { return strings.Compare(a.Change, b.Change) })
		if nerr := daemonNotify(ctx, rs, &state, notifyCmd); err == nil {
			err = nerr
		}
		if eerr := daemonEvents(ctx, &state, eventCmd); err == nil {
			err = eerr
		}
		saveState()
		return err
	}
	if r.Once {
		return resp, poll()
	}
	reexec := daemonReexec()
	var interrupted <-chan struct{}
	if ctx.Interrupts != nil {
		var stop func()
		interrupted, stop = ctx.Interrupts()
		defer stop()
	}
	ticks, stopTicks := daemonTicker()
	defer stopTicks()
	l, err := daemonLoop(ctx, common, lock, poll, func() int { return len(pending) }, reexec, ticks, interrupted)
	if l != nil {
		lock = l
	}
	return resp, err
}

// daemonLoop polls now and at each tick until interrupted. At the end of a
// poll with no start pending (pending), a rebuilt binary takes over
// (reexec): live runs are separate processes and carry on. An exec that
// returns is not tried again. It returns the lock it holds at the end, nil
// with the error when it could not take the lock back after an exec.
func daemonLoop(ctx *Ctx, common string, lock *daemon.Lock, poll func() error, pending func() int, reexec *daemon.Reexec, ticks <-chan time.Time, interrupted <-chan struct{}) (*daemon.Lock, error) {
	for {
		if err := poll(); err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "rhiza: daemon: %v\n", err)
		}
		if reexec.Due(pending()) {
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: binary changed, re-executing\n")
			}
			l, err := reexec.Do(common, lock, ctx.Now())
			if l == nil {
				return nil, err
			}
			lock = l
			if err != nil && ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: re-exec: %v\n", err)
			}
			reexec = nil // do not retry a failed exec every poll
		}
		select {
		case <-interrupted:
			return lock, nil
		case <-ticks:
		}
	}
}

// daemonTicker ticks once per poll interval; tests replace it.
var daemonTicker = func() (<-chan time.Time, func()) {
	t := time.NewTicker(daemonPoll)
	return t.C, t.Stop
}

// daemonReexec is how the daemon re-executes itself; tests replace it.
var daemonReexec = daemon.NewReexec

// daemonNotify sends one Rex notification for each inbox item that waits
// and that state does not list as notified, and records it there. An item
// that no longer waits drops out of state.Notified. A run keeps its own
// blocked-run notice; this is the daemon's, written as a blocked report to
// the terminal Rex shows the daemon in (ctx.Rex). A working report goes
// before each one so that Rex sees each as a new blocked state.
func daemonNotify(ctx *Ctx, rs *run.Store, state *daemon.State, notifyCmd string) error {
	runs, err := rs.List()
	if err != nil {
		return Errorf("read runs: %v", err)
	}
	host, _ := os.Hostname()
	items := inboxItems(runs, host, nil, rs.Live, changeStateOf(ctx), nil)
	if overlaps, err := inboxOverlaps(ctx); err == nil {
		items = append(items, overlaps...)
	}
	// This poll's disk-low state, not the one saved after the last poll.
	items = append(items, diskLowItems(state.DiskLow)...)
	titles := map[string]string{}
	for _, j := range runs {
		if st := j.State(); st.Title != "" {
			titles[st.Change] = st.Title
		}
	}
	state.Hooked = daemonHook(ctx, notifyCmd, items, titles, state.Hooked)
	var kept []string
	for _, it := range items {
		if slices.Contains(state.Notified, it.ID) {
			kept = append(kept, it.ID)
			continue
		}
		if ctx.Rex != nil {
			kept = append(kept, it.ID)
			what := titles[it.Change]
			if what == "" {
				what = it.Change
			}
			ctx.Rex.Write(rexReport("working", "", "", ""))
			ctx.Rex.Write(rexReport("blocked", "", "", inboxNotice(what, it)))
		}
	}
	state.Notified = kept
	return nil
}

// inboxNotice is the notification text for item it of the change set called
// what.
func inboxNotice(what string, it InboxItem) string {
	switch it.Kind {
	case InboxQuestion:
		msg := what + ": a question waits on you"
		if q := strings.TrimSpace(it.Question); q != "" {
			msg += ": " + q
		}
		return msg
	case InboxSignoff:
		return what + ": a sign-off waits on you"
	case InboxOpen:
		return what + ": a plan run stopped for open questions"
	case InboxFailed:
		return what + ": a run failed"
	case InboxSkipped:
		return what + ": an interrupted run was not resumed"
	case InboxDiskLow:
		return "disk space is low: the daemon holds new starts and resumes"
	case InboxOverlap:
		return what + ": an overlap with change set " + it.Other + " waits on you"
	default:
		return what + ": the daemon gave up on a run"
	}
}

// daemonStartGrace is how long a change set the daemon started is held as
// running while its run has no journal yet.
const daemonStartGrace = time.Minute

// daemonUnjournaledHold is how long after its start a change set whose run
// wrote no journal stays out of the daemon's picks before it is picked
// again.
const daemonUnjournaledHold = 15 * time.Minute

// daemonPending is a start of the daemon's: when, and how many runs the
// change set had then, so a new journal shows as one more.
type daemonPending struct {
	at   time.Time
	runs int
}

// daemonRestoreUnjournaled is the unjournaled holds a daemon's state
// recorded, for the daemon that takes over from it.
func daemonRestoreUnjournaled(events []daemon.Event) map[string]daemonPending {
	unjournaled := map[string]daemonPending{}
	for _, e := range events {
		if e.Change != "" {
			unjournaled[e.Change] = daemonPending{at: e.At, runs: e.Runs}
		}
	}
	return unjournaled
}

// daemonHeld settles pending against runs at now and returns the change
// sets still pending, which count as running. A start whose change set has
// a new journal is dropped from pending. One with none once the grace
// passed failed before it journaled (a refused permission mode, a dirty
// worktree, a usage error): it moves to unjournaled, which are not picked
// nor take a slot until a journal of the change set shows, daemonUnjournaledHold
// passes since the start, or the change set is no longer accepted (landed or
// withdrawn: live says whether it still is), whichever comes first.
func daemonHeld(pending map[string]daemonPending, unjournaled map[string]daemonPending, runs []run.Journal, live func(string) bool, now time.Time) []string {
	var held []string
	for cs, p := range pending {
		switch {
		case daemonRunsOf(runs, cs) > p.runs:
			delete(pending, cs)
		case now.Sub(p.at) >= daemonStartGrace:
			delete(pending, cs)
			unjournaled[cs] = p
		default:
			held = append(held, cs)
		}
	}
	for cs, p := range unjournaled {
		switch {
		case daemonRunsOf(runs, cs) > p.runs:
			// A run of it journaled since (the owner started one): its
			// runs say what to do with it now.
			delete(unjournaled, cs)
		case now.Sub(p.at) >= daemonUnjournaledHold, !live(cs):
			delete(unjournaled, cs)
		}
	}
	slices.Sort(held)
	return held
}

// daemonRunsOf counts the runs of change set cs.
func daemonRunsOf(runs []run.Journal, cs string) int {
	n := 0
	for _, j := range runs {
		if j.State().Change == cs {
			n++
		}
	}
	return n
}

// daemonPoll1 picks and starts runs once. pending holds the change sets
// started; unjournaled those whose start never journaled (see daemonHeld).
func daemonPoll1(ctx *Ctx, rs *run.Store, pending map[string]daemonPending, unjournaled map[string]daemonPending, max int, minFree uint64, autonomous bool, mode string, skips *[]daemon.Skip) ([]string, []string, error) {
	res, _, _, err := planRead(ctx, "daemon", "")
	if err != nil {
		return nil, nil, err
	}
	ws, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, nil, err
	}
	cands, planSkips := daemonCandidates(res, func(cs id.ID) bool {
		c, ok := ws.changeSet(cs)
		return ok && len(changeLeaves(ws.model(cs), c.Ops)) > 0
	})
	runs, err := rs.List()
	if err != nil {
		return nil, nil, Errorf("read runs: %v", err)
	}
	host, _ := os.Hostname()
	var resumed []string
	low, free := daemonDiskLow(ctx, minFree)
	accepted := func(change string) bool {
		cs, err := id.Parse(change)
		if err != nil {
			return false
		}
		c, ok := ws.changeSet(cs)
		return ok && c.State == changeset.StateAccepted
	}
	var resumable []string
	if !low {
		resumable = daemon.Resumable(runs, host, nil, rs.Live, accepted, ctx.Now())
	}
	for _, id := range resumable {
		// Journal the resume before the tab opens: a resume is tried once.
		if err := rs.Append(id, run.Entry{Type: run.TypeDaemonResumed, At: ctx.Now(), Detail: "resumed by the daemon"}); err != nil {
			return nil, resumed, Errorf("journal the resume of run %s: %v", id, err)
		}
		if err := daemonResume(ctx, id); err != nil {
			return nil, resumed, err
		}
		resumed = append(resumed, id)
	}
	held := daemonHeld(pending, unjournaled, runs, accepted, ctx.Now())
	var unj []daemon.Skip
	cands = slices.DeleteFunc(cands, func(cs string) bool {
		_, ok := unjournaled[cs]
		if ok {
			unj = append(unj, daemon.Skip{Change: cs, Reason: daemon.SkipUnjournaled})
		}
		return ok
	})
	ship := daemon.Ship{Workflow: shipChangeWorkflow, Since: func(change string) time.Time {
		cs, err := id.Parse(change)
		if err != nil {
			return time.Time{}
		}
		c, ok := ws.changeSet(cs)
		if !ok {
			return time.Time{}
		}
		var at time.Time
		for _, m := range c.Lifecycle {
			if m.Type == entry.TypeAccept || m.Type == entry.TypeReopen {
				at = m.At
			}
		}
		return at
	}}
	var picked []string
	if low {
		detail := fmt.Sprintf("%s free, under the minimum of %s", daemon.FormatSize(free), daemon.FormatSize(minFree))
		var low []daemon.Skip
		for _, cs := range cands {
			low = append(low, daemon.Skip{Change: cs, Reason: daemon.SkipDiskLow, Detail: []string{detail}})
		}
		*skips = append(append(planSkips, unj...), low...)
	} else {
		picked = daemon.Pick(cands, runs, held, host, nil, max, ship)
		*skips = append(append(planSkips, unj...), daemon.Skips(cands, runs, held, host, nil, max, ship)...)
	}
	for i, sk := range *skips {
		if cs, err := id.Parse(sk.Change); err == nil {
			if c, ok := ws.changeSet(cs); ok {
				(*skips)[i].Title = c.Open.DisplayTitle()
			}
		}
	}
	var started []string
	for _, cs := range picked {
		if err := daemonStart(ctx, cs, autonomous, mode); err != nil {
			return started, resumed, err
		}
		pending[cs] = daemonPending{at: ctx.Now(), runs: daemonRunsOf(runs, cs)}
		started = append(started, cs)
	}
	return started, resumed, nil
}

// daemonDiskLow reports whether free space on the repo's filesystem is under
// minFree, and the free space. A minFree of 0 turns the guard off, and so
// does a filesystem whose free space cannot be read. RHIZA_TEST_FREE_BYTES
// stands in for the real figure in tests; like RHIZA_TEST_RENDER_STALE it
// applies only alongside RHIZA_TEST_NOW, so a real daemon never sees it.
func daemonDiskLow(ctx *Ctx, minFree uint64) (bool, uint64) {
	if minFree == 0 {
		return false, 0
	}
	free, err := daemon.FreeBytes(ctx.Root)
	getenv := ctxGetenv(ctx)
	if v := getenv("RHIZA_TEST_FREE_BYTES"); v != "" && getenv("RHIZA_TEST_NOW") != "" {
		n, perr := strconv.ParseUint(v, 10, 64)
		if perr != nil {
			return false, 0
		}
		free, err = n, nil
	}
	if err != nil {
		return false, 0
	}
	return free < minFree, free
}

// daemonCandidates are the accepted, ready change sets of the milestone
// groups that hasLeaves, in roadmap order, each once. skips are the other
// accepted change sets, each once: those serving no milestone, those not
// ready (with their readiness reasons) and those with no leaves.
func daemonCandidates(res planResult, hasLeaves func(id.ID) bool) ([]string, []daemon.Skip) {
	var out []string
	var skips []daemon.Skip
	seen := map[id.ID]bool{}
	for _, g := range res.groups {
		if g.backlog {
			continue
		}
		for _, p := range g.members {
			if seen[p.info.ID] || p.info.State != changeset.StateAccepted {
				continue
			}
			seen[p.info.ID] = true
			switch {
			case !p.ready:
				skips = append(skips, daemon.Skip{Change: p.info.ID.String(), Reason: daemon.SkipNotReady, Detail: slices.Clone(p.reasons)})
			case !hasLeaves(p.info.ID):
				skips = append(skips, daemon.Skip{Change: p.info.ID.String(), Reason: daemon.SkipNoLeaves})
			default:
				out = append(out, p.info.ID.String())
			}
		}
	}
	for _, g := range res.groups {
		if !g.backlog {
			continue
		}
		for _, p := range g.members {
			if seen[p.info.ID] || p.info.State != changeset.StateAccepted {
				continue
			}
			seen[p.info.ID] = true
			skips = append(skips, daemon.Skip{Change: p.info.ID.String(), Reason: daemon.SkipNoMilestone})
		}
	}
	return out, skips
}

// daemonStart starts the run of change set cs in a new Rex tab, with
// permission mode mode.
func daemonStart(ctx *Ctx, cs string, autonomous bool, mode string) error {
	if ctx.Exec == nil || ctx.Self == "" {
		return Errorf("daemon cannot start a run of %s: no rex or rhiza path", cs)
	}
	args := []string{"window", "new"}
	if ctx.RexSession != "" {
		args = append(args, "-s", ctx.RexSession)
	}
	args = append(args, cs[:min(len(cs), 9)], "-c", ctx.Root, "--json", "--", ctx.Self, "run", "start", cs, "--permission-mode", mode)
	if autonomous {
		args = append(args, "--autonomous")
	}
	if _, err := ctx.Exec("rex", args, ctx.Root, nil); err != nil {
		return Errorf("rex window new for %s: %v", cs, err)
	}
	return nil
}

// daemonResume opens a new Rex tab running rhiza run resume on run id.
func daemonResume(ctx *Ctx, id string) error {
	if ctx.Exec == nil || ctx.Self == "" {
		return Errorf("daemon cannot resume run %s: no rex or rhiza path", id)
	}
	args := []string{"window", "new"}
	if ctx.RexSession != "" {
		args = append(args, "-s", ctx.RexSession)
	}
	args = append(args, id[:min(len(id), 9)], "-c", ctx.Root, "--json", "--", ctx.Self, "run", "resume", id)
	if _, err := ctx.Exec("rex", args, ctx.Root, nil); err != nil {
		return Errorf("rex window new for run %s: %v", id, err)
	}
	return nil
}

func init() {
	register(Operation{
		Name:     "daemon.status",
		Summary:  "Show the daemon of this repository: whether it runs, its settings, its live runs and what it last started or resumed",
		Request:  reflect.TypeFor[DaemonStatusRequest](),
		Response: reflect.TypeFor[DaemonStatusResponse](),
		Run:      runDaemonStatus,
	})
}

// DaemonStatusRequest takes nothing.
type DaemonStatusRequest struct{}

// DaemonRun is a run the daemon sees as live.
type DaemonRun struct {
	ID     string `json:"id" desc:"The run ID"`
	Change string `json:"change" desc:"The change set it runs on"`
	State  string `json:"state" desc:"running or interrupted"`
}

// DaemonEvent is the last run the daemon started or resumed.
type DaemonEvent struct {
	Change string `json:"change,omitempty" desc:"The change set whose run was started; absent for a resume"`
	Run    string `json:"run,omitempty" desc:"The run that was resumed; absent for a start"`
	At     string `json:"at" desc:"When, RFC 3339 UTC"`
}

// DaemonSkip is an accepted change set the daemon did not start at its last
// poll.
type DaemonSkip struct {
	Change string   `json:"change" desc:"The change set's ID"`
	Title  string   `json:"title" desc:"The change set's display title"`
	Reason string   `json:"reason" desc:"Why it was not started: no-milestone, no-leaves, not-ready, live, failed, shipped, held, unjournaled, no-slot or disk-low"`
	Detail []string `json:"detail" desc:"The readiness reasons, when the reason is not-ready; else empty"`
}

// skipText is the words for a skip reason in the text status.
var skipText = map[string]string{
	daemon.SkipNoMilestone: "serves no milestone",
	daemon.SkipNoLeaves:    "has no leaves",
	daemon.SkipNotReady:    "not ready",
	daemon.SkipLive:        "has a live run",
	daemon.SkipFailed:      "its latest run failed and waits in the inbox",
	daemon.SkipShipped:     "its latest run shipped it",
	daemon.SkipHeld:        "started, its run not journaled yet",
	daemon.SkipUnjournaled: "its start wrote no journal",
	daemon.SkipNoSlot:      "waiting for a slot",
	daemon.SkipDiskLow:     "disk space is low",
}

// DaemonStatusResponse is the daemon's state. While no daemon holds the
// lock, Running is false and the rest, if any, is what the last daemon left.
type DaemonStatusResponse struct {
	Running     bool          `json:"running" desc:"True while a live daemon holds the repository's lock"`
	Process     string        `json:"process,omitempty" desc:"The daemon's process, as host:pid; absent when none ever ran"`
	Version     string        `json:"version,omitempty" desc:"The version of the binary the daemon runs, as recorded when it started"`
	Started     string        `json:"started,omitempty" desc:"When the daemon started, RFC 3339 UTC"`
	LastPoll    string        `json:"last_poll,omitempty" desc:"When it last polled, RFC 3339 UTC"`
	Stopped     string        `json:"stopped,omitempty" desc:"When the last daemon stopped, RFC 3339 UTC; absent while it runs"`
	MaxRuns     int           `json:"max_runs,omitempty" desc:"How many runs it keeps running at once"`
	Autonomous  bool          `json:"autonomous" desc:"True when it starts each run with --autonomous"`
	Runs        []DaemonRun   `json:"runs" desc:"The live runs (running or interrupted) as of its last poll"`
	LastStarted *DaemonEvent  `json:"last_started,omitempty" desc:"The run it last started; absent if none"`
	LastResumed *DaemonEvent  `json:"last_resumed,omitempty" desc:"The run it last resumed; absent if none"`
	Skipped     []DaemonSkip  `json:"skipped" desc:"The accepted change sets it did not start at its last poll, each with its reason"`
	Unjournaled []DaemonEvent `json:"unjournaled" desc:"The runs it started that never wrote a journal, by change set and when; it starts these again 15 minutes after the start, or sooner when a run of the change set journals"`
}

// Human writes the status as text.
func (r DaemonStatusResponse) Human(w io.Writer) {
	if r.Running {
		fmt.Fprintf(w, "Daemon running, process %s\n", r.Process)
	} else {
		fmt.Fprintln(w, "Daemon not running")
		if r.Process == "" {
			return
		}
		fmt.Fprintf(w, "Last daemon: process %s, stopped %s\n", r.Process, r.Stopped)
	}
	if r.Version != "" {
		fmt.Fprintf(w, "Version: %s\n", r.Version)
	}
	fmt.Fprintf(w, "Started: %s\nLast poll: %s\nMax runs: %d\nAutonomous: %t\n", r.Started, r.LastPoll, r.MaxRuns, r.Autonomous)
	if len(r.Runs) == 0 {
		fmt.Fprintln(w, "Live runs: none")
	} else {
		fmt.Fprintln(w, "Live runs:")
		for _, d := range r.Runs {
			fmt.Fprintf(w, "- %s %s (change set %s)\n", d.ID, d.State, d.Change)
		}
	}
	if e := r.LastStarted; e != nil {
		fmt.Fprintf(w, "Last started: a run of %s at %s\n", e.Change, e.At)
	}
	if e := r.LastResumed; e != nil {
		fmt.Fprintf(w, "Last resumed: run %s at %s\n", e.Run, e.At)
	}
	for _, k := range r.Skipped {
		fmt.Fprintf(w, "Skipped: %s %s: %s", k.Change, k.Title, skipText[k.Reason])
		if len(k.Detail) > 0 {
			fmt.Fprintf(w, " (%s)", strings.Join(k.Detail, "; "))
		}
		fmt.Fprintln(w)
	}
	for _, e := range r.Unjournaled {
		fmt.Fprintf(w, "Not journaled: the run of %s started at %s wrote no journal; retried at %s, or sooner when a run of it journals\n", e.Change, e.At, daemonRetryAt(e.At))
	}
}

// daemonRetryAt is when a start at the RFC 3339 time at is picked again.
func daemonRetryAt(at string) string {
	t, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return at
	}
	return t.Add(daemonUnjournaledHold).UTC().Format(time.RFC3339)
}

// runDaemonStatus reads the lock and the state file. Without a live
// daemon holding the lock it says not running, and shows the state file's
// last daemon, if one left it.
func runDaemonStatus(ctx *Ctx, req any) (any, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	resp := DaemonStatusResponse{Runs: []DaemonRun{}, Skipped: []DaemonSkip{}, Unjournaled: []DaemonEvent{}}
	_, resp.Running = daemon.Running(common)
	st, ok := daemon.ReadState(common)
	if !ok || (resp.Running && st.Process != daemon.ReadHolder(common)) {
		// A live daemon that has not written its state yet, or a lock
		// naming no state: say what the lock says.
		if resp.Running {
			resp.Process = daemon.ReadHolder(common).String()
		}
		return resp, nil
	}
	resp.Process = st.Process.String()
	resp.Version = st.Version
	resp.Started = st.Process.Started.UTC().Format(time.RFC3339)
	resp.LastPoll = st.LastPoll.UTC().Format(time.RFC3339)
	if !st.Stopped.IsZero() && !resp.Running {
		resp.Stopped = st.Stopped.UTC().Format(time.RFC3339)
	}
	resp.MaxRuns, resp.Autonomous = st.MaxRuns, st.Autonomous
	for _, d := range st.Runs {
		resp.Runs = append(resp.Runs, DaemonRun{ID: d.ID, Change: d.Change, State: d.State})
	}
	resp.LastStarted = daemonEvent(st.LastStarted)
	resp.LastResumed = daemonEvent(st.LastResumed)
	for _, k := range st.Skips {
		resp.Skipped = append(resp.Skipped, DaemonSkip{Change: k.Change, Title: k.Title, Reason: k.Reason, Detail: append([]string{}, k.Detail...)})
	}
	for _, e := range st.Unjournaled {
		resp.Unjournaled = append(resp.Unjournaled, *daemonEvent(&e))
	}
	return resp, nil
}

func daemonEvent(e *daemon.Event) *DaemonEvent {
	if e == nil {
		return nil
	}
	return &DaemonEvent{Change: e.Change, Run: e.Run, At: e.At.UTC().Format(time.RFC3339)}
}

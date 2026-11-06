package api

// This file holds the daemon operation: rhiza daemon runs in the
// foreground, one per repository. It holds a lock under the git common dir,
// so a second daemon started from any worktree of the repository is
// refused, naming the first one's process. It binds no socket.

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/gosweep"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
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
	Planned []string `json:"planned" desc:"The plan requests whose plan runs the daemon started, in order"`
}

// Human writes the response as markdown.
func (r DaemonResponse) Human(w io.Writer) {
	for _, cs := range r.Started {
		fmt.Fprintf(w, "Started a run of %s\n", cs)
	}
	for _, id := range r.Resumed {
		fmt.Fprintf(w, "Resumed run %s\n", id)
	}
	for _, q := range r.Planned {
		fmt.Fprintf(w, "Started the plan run of request %s\n", q)
	}
	fmt.Fprintf(w, "Daemon %s stopped\n", r.Process)
}

// runDaemon takes the root's daemon lock and keeps it until the caller
// interrupts (Ctx.Interrupt). A second daemon exits with a usage error
// naming the holder. Each poll reads every project registered in the
// root's rhiza.toml and works the initialised ones as the per-repo daemon
// did, in the project's own checkout and journal; the others are skipped
// and status says so. The project the daemon was started in (else the first
// initialised one) is its home: its git common dir holds the per-repo lock,
// the daemon's state, the re-exec and the web UI.
func runDaemon(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("daemon runs in the foreground and has no dry run")
	}
	r := req.(DaemonRequest)
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return nil, Usagef("%s is not inside a Rhiza root: the daemon runs on the root", ctx.Wd).WithNext("rhiza root init")
	}
	if r.MaxRuns < 0 {
		return nil, Usagef("--max-runs must be at least 1")
	}
	reg, err := rootcfg.Load(root)
	if err != nil {
		return nil, err
	}
	homeName, homeCtx, ok := daemonHome(ctx, root, reg)
	if !ok {
		return nil, Usagef("%s has no registered, initialised project for the daemon to work", root).WithNext("rhiza project add")
	}
	conf, err := daemonConfOf(homeCtx, r)
	if err != nil {
		return nil, err
	}
	// Every run the daemon starts is a ship-change run with this mode: one
	// that run start would refuse is refused here, before any tab opens.
	_, def, _, err := resolveWorkflow(homeCtx.Root, "")
	if err != nil {
		return nil, err
	}
	if err := checkAgentSteps(homeCtx, def, conf.mode, "", r.Autonomous, "rhiza describe daemon"); err != nil {
		return nil, err
	}
	common, err := gitCommonDir(homeCtx)
	if err != nil {
		return nil, err
	}
	// Listen for an interrupt before the lock appears: whoever sees the
	// lock may signal at once, and that must stop the daemon cleanly.
	var interrupted <-chan struct{}
	if !r.Once && ctx.Interrupts != nil {
		var stop func()
		interrupted, stop = ctx.Interrupts()
		defer stop()
	}
	rootLock, err := daemon.TryLockRoot(root, ctx.Now())
	if err != nil {
		if h, ok := err.(*daemon.Held); ok {
			return nil, Usagef("%s", h.Error()).WithNext("rhiza daemon status")
		}
		return nil, err
	}
	defer rootLock.Unlock()
	// Status from anywhere under the root reads the state in the home's
	// git common dir: say which project that is.
	if err := daemon.WriteHome(root, homeName); err != nil {
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
	resp := DaemonResponse{Process: self.String(), Started: []string{}, Resumed: []string{}, Planned: []string{}}
	projs := map[string]*daemonProj{}
	// project is the daemon's working state for a registered project,
	// made when it first sees it initialised and kept from poll to poll.
	project := func(name, dir string) (*daemonProj, error) {
		if p, ok := projs[name]; ok && p.dir == dir {
			return p, nil
		}
		pctx, startConf := homeCtx, &conf
		if name != homeName {
			pctx, startConf = projectCtx(ctx, dir), nil
		}
		p, err := newDaemonProj(pctx, name, dir, self, name == homeName, r, startConf)
		if err != nil {
			return nil, err
		}
		projs[name] = p
		return p, nil
	}
	_, err = project(homeName, homeCtx.Root)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, p := range projs {
			p.stop()
		}
	}()
	selfUpdateOn := false
	// The web UI listener follows rhiza.web.host and rhiza.web.port from poll to poll. A single
	// poll (--once) and a test run (RHIZA_TEST_NOW) bind nothing.
	var webUI *web.Manager
	if !r.Once && ctxGetenv(ctx)("RHIZA_TEST_NOW") == "" {
		// The handler gets its own copy of ctx: the poll goes on changing
		// ctx while requests run.
		snap := *homeCtx
		health := NewWebHealth(&snap)
		webUI = &web.Manager{Factory: webFactory, Handler: web.HandlerWithMilestones(WebDaemonSource(&snap, health), WebRunsSource(&snap), WebInboxSource(&snap), WebMilestonesSource(&snap), ReadOnlyHandler(&snap), web.Events(WebEventSources(&snap, health), web.PollEvery))}
		defer webUI.Close()
	}
	var sweep daemon.SweepSchedule
	// lastTurn is the project that took the last run slot: the next poll's
	// round-robin starts after it, so turns rotate across polls instead of
	// restarting at the home project.
	lastTurn := ""
	// atHome is told the home project's settings at each poll: the root-wide
	// parts of the daemon follow them.
	atHome := func(conf daemonConf) {
		if webUI != nil {
			if werr := webUI.Apply(conf.webHost, conf.webPort); werr != nil && ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: %v\n", werr)
			}
		}
		selfUpdateOn = conf.selfUpdate
		if conf.trimGoCache {
			daemonSweepGoCache(homeCtx, &sweep, conf.minFree)
		}
	}
	poll := func() error {
		reg, err := rootcfg.Load(root)
		if err != nil {
			return err
		}
		var firstErr error
		var skipped []daemon.ProjectSkip
		var polled []string
		// The home project goes first: it holds the daemon's state.
		names := append([]string{homeName}, slices.DeleteFunc(reg.Names(), func(n string) bool { return n == homeName })...)
		for _, name := range names {
			if _, ok := reg.Projects[name]; !ok {
				continue
			}
			dir := reg.Dir(root, name)
			if name == homeName {
				dir = homeCtx.Root
			}
			if !projectInitialised(dir) {
				skipped = append(skipped, daemon.ProjectSkip{Project: name, Reason: daemon.ProjectNotInitialised})
				continue
			}
			_, err := project(name, dir)
			if err != nil {
				var held *daemon.Held
				if errors.As(err, &held) {
					skipped = append(skipped, daemon.ProjectSkip{Project: name, Reason: daemon.ProjectHeld, Detail: err.Error()})
				} else {
					skipped = append(skipped, daemon.ProjectSkip{Project: name, Reason: daemon.ProjectNoGit, Detail: err.Error()})
				}
				continue
			}
			polled = append(polled, name)
		}
		// Every project's demand is read before any starts: the run slots
		// are machine-wide, filled round-robin across the projects (see
		// daemon.RoundRobin), each under its own rhiza.project.maxRuns.
		// The home project's settings give the machine-wide maxRuns.
		maxRuns := daemon.DefaultMaxRuns
		if r.MaxRuns > 0 {
			maxRuns = r.MaxRuns
		}
		turns := make([]daemonTurn, len(polled))
		var picks []daemon.ProjectPick
		running := 0
		for i, name := range polled {
			var onConf func(daemonConf)
			if name == homeName {
				onConf = atHome
			}
			turns[i] = projs[name].prepare(r, onConf)
			if turns[i].err != nil {
				continue
			}
			if name == homeName {
				maxRuns = turns[i].conf.maxRuns
			}
			picks = append(picks, turns[i].pick)
			running += turns[i].pick.Running
		}
		share := map[string]int{}
		picked := daemon.RoundRobin(picks, running, maxRuns, lastTurn)
		lastTurn = daemon.LastTurn(picked, lastTurn)
		for _, pk := range picked {
			share[pk.Project]++
		}
		for i, name := range polled {
			if err := projs[name].poll(r, &resp, turns[i], maxRuns, share[name]); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		// A project dropped from the registry stops being polled.
		for name, p := range projs {
			if !slices.Contains(polled, name) {
				p.stop()
				delete(projs, name)
			}
		}
		// Every project's state says which of the others were skipped.
		for _, p := range projs {
			p.state.ProjectSkips = skipped
			p.saveState()
		}
		return firstErr
	}
	if r.Once {
		return resp, poll()
	}
	reexec := daemonReexec()
	update := &daemon.SelfUpdate{Now: ctx.Now, Every: daemon.SelfUpdateEvery}
	selfUpdate := func() {
		if !selfUpdateOn || ctx.Exec == nil || ctx.Git == nil {
			return
		}
		reg, err := rootcfg.Load(root)
		if err != nil {
			return
		}
		main, ok := daemon.SelfUpdateCheckout(root, reg)
		if !ok {
			return
		}
		update.Git = func(args ...string) (string, error) { return gitAt(ctx, main, args...) }
		update.Fetch = func() error { return fetchMainWithin(ctx, main, selfUpdateFetchTimeout) }
		update.Make = func(args ...string) error {
			_, err := ctx.Exec("make", args, main, nil)
			return err
		}
		update.Log = func(format string, a ...any) {
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: daemon: "+format+"\n", a...)
			}
		}
		update.Run()
	}
	ticks, stopTicks := daemonTicker()
	defer stopTicks()
	pendingAll := func() int {
		n := 0
		for _, p := range projs {
			n += len(p.pending)
		}
		return n
	}
	l, err := daemonLoop(ctx, common, lock, poll, pendingAll, reexec, selfUpdate, ticks, interrupted)
	if l != nil {
		lock = l
	}
	return resp, err
}

// daemonStartState is the state a daemon image writes before its first poll:
// its own process, start time and settings, and what the daemon before it
// (an earlier image of a re-exec, or a restart) left: its last poll, live
// runs, skips, hooks and unjournaled holds. The holds come back as the
// daemon's unjournaled map.
func daemonStartState(common string, self daemon.Holder, version, commit string, maxRuns int, autonomous bool) (daemon.State, map[string]daemonPending) {
	state := daemon.State{Process: self, Version: version, Commit: commit, MaxRuns: maxRuns, Autonomous: autonomous, Runs: []daemon.LiveRun{}}
	unjournaled := map[string]daemonPending{}
	if prev, ok := daemon.ReadState(common); ok {
		state.LastPoll = prev.LastPoll
		if prev.Runs != nil {
			state.Runs = prev.Runs
		}
		state.Skips = prev.Skips
		state.Unjournaled = prev.Unjournaled
		state.LastStarted, state.LastResumed = prev.LastStarted, prev.LastResumed
		state.Notified = prev.Notified
		state.Nudged = prev.Nudged
		state.Hooked = prev.Hooked
		state.Events, state.EventsSeeded = prev.Events, prev.EventsSeeded
		state.Hooks = prev.Hooks
		state.DiskLow = prev.DiskLow
		state.UsageHold = prev.UsageHold
		state.UsageWindows = prev.UsageWindows
		state.Learned = prev.Learned
		// A re-exec, or a restart, keeps the holds of the daemon before.
		unjournaled = daemonRestoreUnjournaled(prev.Unjournaled)
	}
	return state, unjournaled
}

// daemonLoop polls now and at each tick until interrupted. At the end of a
// poll with no start pending (pending), a rebuilt binary takes over
// (reexec): live runs are separate processes and carry on. An exec that
// returns is not tried again. It returns the lock it holds at the end, nil
// with the error when it could not take the lock back after an exec.
func daemonLoop(ctx *Ctx, common string, lock *daemon.Lock, poll func() error, pending func() int, reexec *daemon.Reexec, selfUpdate func(), ticks <-chan time.Time, interrupted <-chan struct{}) (*daemon.Lock, error) {
	for {
		if ctx.Reap != nil {
			ctx.Reap(common)
		}
		if err := poll(); err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "rhiza: daemon: %v\n", err)
		}
		if selfUpdate != nil && pending() == 0 {
			selfUpdate()
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

// webFactory starts the web UI listener; tests replace it.
var webFactory web.Factory = web.NetFactory

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
	stateOf := func(string) (string, bool) { return "", false }
	var overlaps []InboxItem
	fold, ferr := daemonFoldOf(ctx)
	if ferr == nil {
		stateOf = fold.stateOf
		for _, ov := range fold.overlaps {
			if o, err := id.Parse(ov.Other); err == nil {
				if c, ok := fold.changeSet(o); ok {
					ov.OtherTitle = c.Open.DisplayTitle()
				}
			}
			overlaps = append(overlaps, ov)
		}
	}
	items := inboxItems(runs, host, nil, rs.Live, stateOf, nil, nil)
	items = append(items, overlaps...)
	// This poll's disk-low state, not the one saved after the last poll.
	items = append(items, diskLowItems(state.DiskLow)...)
	items = append(items, usageHoldItems(state.UsageHold)...)
	if common, err := gitCommonDir(ctx); err == nil {
		items = append(items, mainStaleItems(common)...)
	}
	titles := map[string]string{}
	for _, j := range runs {
		if st := j.State(); st.Title != "" {
			titles[st.Change] = st.Title
		}
	}
	// A change set with no run, such as either side of an overlap, takes
	// its title from the fold: a notice never names a change set by ID.
	if ferr == nil {
		for _, it := range items {
			if _, ok := titles[it.Change]; ok || it.Change == "" {
				continue
			}
			if cs, err := id.Parse(it.Change); err == nil {
				if c, ok := fold.changeSet(cs); ok {
					titles[it.Change] = c.Open.DisplayTitle()
				}
			}
		}
	}
	for i := range items {
		if items[i].Kind == InboxOpen && items[i].Questions == nil {
			if j := slices.IndexFunc(runs, func(j run.Journal) bool { return j.ID == items[i].Run }); j >= 0 {
				items[i].Questions = inboxGateQuestions(ctx, runs[j].State())
			}
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
			what := noticeSubject(titles, it.Change)
			ctx.Rex.Write(rexReport("working", "", "", ""))
			ctx.Rex.Write(rexReport("blocked", "", "", inboxNotice(what, it)))
		}
	}
	state.Notified = kept
	return nil
}

// noticeSubject is what a notice calls change set change: its title, or
// "a change set" when no title is known, never its ID.
func noticeSubject(titles map[string]string, change string) string {
	if t := titles[change]; t != "" {
		return t
	}
	return "a change set"
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
	case InboxInterview:
		return what + ": a milestone interview waits on you"
	case InboxSignoff:
		return what + ": a sign-off waits on you"
	case InboxOpen:
		if len(it.Questions) > 0 {
			if q := noticeQuestionText(it.Questions[0]); q != "" {
				return what + ": a plan run stopped for open questions, the first: " + q
			}
		}
		return what + ": a plan run stopped for open questions"
	case InboxBlocker:
		msg := what + ": a question blocks it"
		if q := strings.TrimSpace(it.Question); q != "" {
			msg += ": " + q
		}
		return msg
	case InboxAccept:
		return what + ": a plan waits for your acceptance"
	case InboxStalled:
		return what + ": a session stalled"
	case InboxFailed:
		return what + ": a run failed"
	case InboxSkipped:
		return what + ": an interrupted run was not resumed"
	case InboxDiskLow:
		return "disk space is low: the daemon holds new starts and resumes"
	case InboxUsageReserve:
		return "usage reserve: the daemon holds new launches until the window resets"
	case InboxMainStale:
		return what + ": main was not updated after the merge: " + it.Detail
	case InboxOverlap:
		other := it.OtherTitle
		if other == "" {
			other = "another change set"
		}
		return what + ": an overlap with " + other + " waits on you"
	default:
		return what + ": the daemon gave up on a run"
	}
}

// noticeQuestionText is the question in plain words: the advisor's
// restatement when it has one, else the full statement.
func noticeQuestionText(q InboxGateQuestion) string {
	if r := strings.TrimSpace(q.Restatement); r != "" {
		return r
	}
	return strings.TrimSpace(q.Statement)
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

// selfUpdateFetchTimeout bounds the self-update's git fetch, so a stalled
// network cannot hold up the daemon's polls.
const selfUpdateFetchTimeout = 30 * time.Second

// fetchMainWithin runs git fetch origin main in dir, killing it after
// timeout. Git never prompts for credentials: it reads nothing and has
// GIT_TERMINAL_PROMPT=0. Without ctx.Start it runs untimed through ctx.Exec.
func fetchMainWithin(ctx *Ctx, dir string, timeout time.Duration) error {
	args := []string{"fetch", "origin", "main"}
	if ctx.Start == nil {
		_, err := gitAt(ctx, dir, args...)
		return err
	}
	p, err := ctx.Start(Program{Name: "git", Args: args, Dir: dir, Env: []string{"GIT_TERMINAL_PROMPT=0"}})
	if err != nil {
		return err
	}
	timedOut := make(chan struct{})
	timer := time.AfterFunc(timeout, func() {
		close(timedOut)
		p.Kill()
	})
	_, _ = io.Copy(io.Discard, p.Output())
	code, err := p.Wait()
	if !timer.Stop() {
		<-timedOut
		return fmt.Errorf("git fetch origin main: timed out after %s", timeout)
	}
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("git fetch origin main: exit %d", code)
	}
	return nil
}

// daemonSlotKept details a no-slot skip whose slot is kept for a queued
// continuation or plan request.
const daemonSlotKept = "waiting for a slot: kept for queued continuations or a plan request"

// daemonContinueCanStart reports whether a queued continuation could start
// now, as daemonQueuedContinues picks them: its run exists, is not live, and
// has not ended unless the continuation is a retry.
func daemonContinueCanStart(rs *run.Store, runs []run.Journal, c daemon.Continue) bool {
	i := slices.IndexFunc(runs, func(j run.Journal) bool { return j.ID == c.Run })
	if i < 0 || rs.Live(c.Run) {
		return false
	}
	return !runs[i].State().Ended || c.Retry
}

// daemonContinuesWaiting counts the queued continuations that could start
// now. An unreadable queue counts none.
func daemonContinuesWaiting(ctx *Ctx, rs *run.Store, runs []run.Journal) int {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return 0
	}
	q, err := daemon.ReadContinues(common)
	if err != nil {
		return 0
	}
	n := 0
	for _, c := range q {
		if daemonContinueCanStart(rs, runs, c) {
			n++
		}
	}
	return n
}

// daemonRequestWaits reports whether a queued feature or bug plan request
// could start now: one with no run whose retry time has passed, as
// daemonPlanRequests picks them. A milestone request starts without a slot,
// so it never counts. An unreadable queue reports false.
func daemonRequestWaits(ctx *Ctx) bool {
	return daemonRequestsWaiting(ctx) > 0
}

// daemonRequestsWaiting counts the queued feature and bug plan requests
// that could start now (see daemonRequestWaits).
func daemonRequestsWaiting(ctx *Ctx) int {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return 0
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return 0
	}
	now := ctx.Now()
	n := 0
	for _, q := range reqs {
		if q.Kind != daemon.KindMilestone && q.Run == "" && !now.Before(q.RetryAt) {
			n++
		}
	}
	return n
}

// daemonShipOf is what says which of fold's change sets have shipped (see
// daemon.Ship).
func daemonShipOf(fold *daemonFolded) daemon.Ship {
	return daemon.Ship{Workflow: shipChangeWorkflow, Since: func(change string) time.Time {
		cs, err := id.Parse(change)
		if err != nil {
			return time.Time{}
		}
		c, ok := fold.changeSet(cs)
		if !ok {
			return time.Time{}
		}
		return acceptedSince(c)
	}}
}

// daemonPoll1 picks and starts runs once: resumes, then ship-change runs,
// then plan runs of queued plan requests in the slots left (see
// daemonPlanRequests). While a queued request could start, ship picks take
// at most max-1 slots (when max > 1), so one is left for it. pending holds the change sets started; unjournaled
// those whose start never journaled (see daemonHeld). It returns the
// change sets started, the runs resumed and the requests started.
func daemonPoll1(ctx *Ctx, rs *run.Store, pending map[string]daemonPending, unjournaled map[string]daemonPending, max int, minFree uint64, autonomous bool, mode string, paused bool, hold *usageHold, skips *[]daemon.Skip) (started, resumed, planned []string, err error) {
	fold, err := daemonFoldOf(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	// The skips are edited below, and appended to: work on copies.
	cands, planSkips := slices.Clone(fold.cands), slices.Clone(fold.planSkips)
	runs, err := rs.List()
	if err != nil {
		return nil, nil, nil, Errorf("read runs: %v", err)
	}
	host, _ := os.Hostname()
	low, free := daemonDiskLow(ctx, minFree)
	// A usage hold stops everything the daemon would launch but a milestone
	// request's plan run; low disk stops all.
	stopped := low || hold != nil
	accepted := func(change string) bool {
		st, ok := fold.stateOf(change)
		return ok && st == string(changeset.StateAccepted)
	}
	var resumable []string
	if !stopped && !paused {
		resumable = daemon.Resumable(runs, host, nil, rs.Live, func(change string) string {
			st, _ := fold.stateOf(change)
			return st
		}, planWorkflow, ctx.Now())
	}
	// The owner's queued continuations are started below, within max; the
	// daemon's own resume of an interrupted run leaves those runs to them.
	var queuedRuns []string
	if !stopped && !paused {
		if common, err := gitCommonDir(ctx); err == nil {
			if q, err := daemon.ReadContinues(common); err == nil {
				for _, c := range q {
					queuedRuns = append(queuedRuns, c.Run)
				}
			}
		}
	}
	resumable = slices.DeleteFunc(resumable, func(id string) bool { return slices.Contains(queuedRuns, id) })
	for _, id := range resumable {
		// Journal the resume before the tab opens: a resume is tried once.
		if err := rs.Append(id, run.Entry{Type: run.TypeDaemonResumed, At: ctx.Now(), Detail: "resumed by the daemon"}); err != nil {
			return nil, resumed, nil, Errorf("journal the resume of run %s: %v", id, err)
		}
		if err := daemonResume(ctx, id); err != nil {
			return nil, resumed, nil, err
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
	ship := daemonShipOf(fold)
	shipMode, err := autonomyStage(ctx, "ship")
	if err != nil {
		return nil, resumed, nil, err
	}
	signoffMode, err := autonomyStage(ctx, "signoff")
	if err != nil {
		return nil, resumed, nil, err
	}
	if signoffMode == AutonomyOwner {
		autonomous = false
	}
	// A milestone request's plan run is the owner's own interactive start:
	// it goes first, whether or not a run slot is free, and keeps no ship
	// slot; the later picks (continuations, queued starts, plan requests)
	// count it as a running run.
	var milestoned []string
	if !paused && !low {
		milestoned, err = daemonPlanRequests(ctx, rs, runs, 0, max, mode, daemon.KindMilestone)
		if err != nil {
			return nil, resumed, milestoned, err
		}
	}
	var picked []string
	if paused {
		var ps []daemon.Skip
		for _, cs := range cands {
			ps = append(ps, daemon.Skip{Change: cs, Reason: daemon.SkipPaused})
		}
		*skips = append(append(planSkips, unj...), ps...)
	} else if shipMode == AutonomyOwner {
		var owner []daemon.Skip
		for _, cs := range cands {
			owner = append(owner, daemon.Skip{Change: cs, Reason: daemon.SkipShipOwner})
		}
		*skips = append(append(planSkips, unj...), owner...)
	} else if low {
		detail := fmt.Sprintf("%s free, under the minimum of %s", daemon.FormatSize(free), daemon.FormatSize(minFree))
		var low []daemon.Skip
		for _, cs := range cands {
			low = append(low, daemon.Skip{Change: cs, Reason: daemon.SkipDiskLow, Detail: []string{detail}})
		}
		*skips = append(append(planSkips, unj...), low...)
	} else if hold != nil {
		var held []daemon.Skip
		for _, cs := range cands {
			held = append(held, daemon.Skip{Change: cs, Reason: daemon.SkipUsageReserve, Detail: []string{hold.Detail()}})
		}
		*skips = append(append(planSkips, unj...), held...)
	} else {
		// Ship picks leave a slot for each queued continuation that could
		// start now and one for a waiting plan request: continuations
		// (open-pr, a resume after a wait, a retry) and plan runs are short,
		// and must not wait behind every ready change set.
		kept := daemonContinuesWaiting(ctx, rs, runs)
		if max > 1 && daemonRequestWaits(ctx) {
			kept++
		}
		shipMax := max - kept
		if shipMax < 0 {
			shipMax = 0
		}
		picked = daemon.Pick(cands, runs, held, host, nil, shipMax, ship)
		shipSkips := daemon.Skips(cands, runs, held, host, nil, shipMax, ship)
		if shipMax < max {
			full := daemon.Pick(cands, runs, held, host, nil, max, ship)
			for i, sk := range shipSkips {
				if sk.Reason == daemon.SkipNoSlot && slices.Contains(full, sk.Change) {
					shipSkips[i].Detail = []string{daemonSlotKept}
				}
			}
		}
		*skips = append(append(planSkips, unj...), shipSkips...)
	}
	for i, sk := range *skips {
		if cs, err := id.Parse(sk.Change); err == nil {
			if c, ok := fold.changeSet(cs); ok {
				(*skips)[i].Title = c.Open.DisplayTitle()
			}
		}
	}
	for _, cs := range picked {
		if err := daemonStart(ctx, cs, autonomous, mode); err != nil {
			return started, resumed, milestoned, err
		}
		pending[cs] = daemonPending{at: ctx.Now(), runs: daemonRunsOf(runs, cs)}
		started = append(started, cs)
	}
	if paused || low {
		return started, resumed, nil, nil
	}
	running := len(held) + len(started) + len(milestoned)
	for _, j := range runs {
		if j.State().Liveness(host, nil) == run.LiveRunning {
			running++
		}
	}
	if hold != nil {
		// A usage hold stops every launch but a milestone request's plan
		// run, started above.
		return started, resumed, milestoned, nil
	}
	continued, err := daemonQueuedContinues(ctx, rs, runs, running, max)
	resumed = append(resumed, continued...)
	if err != nil {
		return started, resumed, nil, err
	}
	running += len(continued)
	// While a plan request could start, queued starts leave it one slot
	// (when max > 1), as ship picks do.
	startMax := max
	if max > 1 && daemonRequestWaits(ctx) {
		startMax--
	}
	queued, err := daemonQueuedStarts(ctx, rs, pending, runs, running, startMax, mode)
	started = append(started, queued...)
	if err != nil {
		return started, resumed, nil, err
	}
	planned, err = daemonPlanRequests(ctx, rs, runs, running+len(queued), max, mode, "")
	return started, resumed, append(milestoned, planned...), err
}

// daemonQueuedStarts starts the runs the owner queued with rhiza run start,
// oldest first, while fewer than max runs run: each as a detached rhiza
// run start --foreground process on its change set, with the workflow,
// agent and flags it was queued with. They are the owner's starts, not
// picks: no milestone, state or leaf filter applies. An entry whose change
// set has a live run, or whose run the daemon just started, is dropped. It
// returns the change sets started.
func daemonQueuedStarts(ctx *Ctx, rs *run.Store, pending map[string]daemonPending, runs []run.Journal, running, max int, mode string) ([]string, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	var started []string
	err = daemon.UpdateStarts(common, func(starts []daemon.Start) ([]daemon.Start, error) {
		var keep []daemon.Start
		for _, s := range starts {
			if _, held := pending[s.Change]; held {
				continue
			}
			if _, live, err := rs.LiveOf(s.Change, ""); err != nil {
				return nil, Errorf("read runs: %v", err)
			} else if live {
				continue
			}
			if running >= max {
				keep = append(keep, s)
				continue
			}
			m := s.PermissionMode
			if m == "" {
				m = mode
			}
			if err := daemonStartQueued(ctx, s, m); err != nil {
				return nil, err
			}
			pending[s.Change] = daemonPending{at: ctx.Now(), runs: daemonRunsOf(runs, s.Change)}
			started = append(started, s.Change)
			running++
		}
		return keep, nil
	})
	if err != nil {
		var e *Error
		if errors.As(err, &e) {
			return started, err
		}
		return started, Errorf("read queued starts: %v", err)
	}
	return started, nil
}

// daemonStartQueued starts the queued run s as a detached rhiza run start
// --foreground process, with permission mode mode.
func daemonStartQueued(ctx *Ctx, s daemon.Start, mode string) error {
	args := []string{"run", "start", "--foreground", s.Change}
	if s.Workflow != "" {
		args = append(args, "--workflow", s.Workflow)
	}
	if mode != "" {
		args = append(args, "--permission-mode", mode)
	}
	if s.Agent != "" {
		args = append(args, "--agent", s.Agent)
	}
	if s.NoWait {
		args = append(args, "--no-wait")
	}
	if s.Autonomous {
		args = append(args, "--autonomous")
	}
	return daemonDetachStart(ctx, "a run of "+s.Change, s.Change, args)
}

// daemonQueuedContinues continues the runs the owner queued with rhiza run
// resume, approve and request-changes, oldest first, while fewer than max
// runs run: each as a detached rhiza run resume --foreground process. They
// are the owner's actions, not picks. An entry whose run is gone, finished
// (a retry's excepted) or live is dropped. It returns the runs continued.
func daemonQueuedContinues(ctx *Ctx, rs *run.Store, runs []run.Journal, running, max int) ([]string, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	var continued []string
	err = daemon.UpdateContinues(common, func(queue []daemon.Continue) ([]daemon.Continue, error) {
		var keep []daemon.Continue
		for _, c := range queue {
			if !daemonContinueCanStart(rs, runs, c) {
				continue
			}
			if running >= max {
				keep = append(keep, c)
				continue
			}
			args := []string{"run", "resume", "--foreground", c.Run}
			if c.Retry {
				args = append(args, "--retry")
			}
			if c.NoWait {
				args = append(args, "--no-wait")
			}
			if err := daemonDetach(ctx, "run "+c.Run, c.Run, nil, args); err != nil {
				return nil, err
			}
			continued = append(continued, c.Run)
			running++
		}
		return keep, nil
	})
	if err != nil {
		var e *Error
		if errors.As(err, &e) {
			return continued, err
		}
		return continued, Errorf("read queued continuations: %v", err)
	}
	return continued, nil
}

// daemonSweepGoCache trims the Go build cache and stale go-build* temp
// folders when the schedule says so, writing one stderr line with what was
// freed, only when something was.
func daemonSweepGoCache(ctx *Ctx, sched *daemon.SweepSchedule, minFree uint64) {
	now := ctx.Now()
	free, ferr := repoFreeBytes(ctx)
	if !sched.Due(now, free, minFree, ferr == nil) {
		return
	}
	sched.Done(now)
	getenv := ctxGetenv(ctx)
	tmp := getenv("TMPDIR")
	if tmp == "" {
		tmp = os.TempDir()
	}
	f := gosweep.SweepDirs(gosweep.CacheDirFrom(getenv), tmp, now)
	if f.Any() && ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon: go build sweep freed %s: %d cache entries, %d temp folders\n",
			daemon.FormatSize(uint64(f.Bytes())), f.CacheEntries, f.TempFolders)
	}
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
	free, err := repoFreeBytes(ctx)
	if err != nil {
		return false, 0
	}
	return free < minFree, free
}

// repoFreeBytes is the free space on the repo's filesystem.
// RHIZA_TEST_FREE_BYTES stands in for it alongside RHIZA_TEST_NOW.
func repoFreeBytes(ctx *Ctx) (uint64, error) {
	getenv := ctxGetenv(ctx)
	if v := getenv("RHIZA_TEST_FREE_BYTES"); v != "" && getenv("RHIZA_TEST_NOW") != "" {
		return strconv.ParseUint(v, 10, 64)
	}
	return daemon.FreeBytes(ctx.Root)
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

// daemonDetach starts rhiza with args as a detached child in the repo, outside
// Rex, with its output appended to the log file name under <git common
// dir>/rhiza/logs and env added to its environment. The child's agent steps
// open their tabs in the repo's Rex session, which the daemon makes sure of
// first and hands the child as its REX_SESSION; the child has no REX_BLOCK,
// since it sits in no tab of its own. A run driven with --foreground by hand
// keeps the caller's session.
func daemonDetach(ctx *Ctx, what, name string, env, args []string) error {
	if ctx.Detach == nil || ctx.Self == "" {
		return Errorf("daemon cannot start %s: no rhiza path", what)
	}
	if err := ensureRepoSession(ctx); err != nil {
		return err
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return err
	}
	// Whatever the child does to the log, the next poll folds it again.
	if ctx.daemonFold != nil {
		ctx.daemonFold.invalidate()
	}
	p := Program{Name: ctx.Self, Args: args, Dir: ctx.Root, Env: append([]string{"REX_SESSION=" + ctx.RexSession}, env...), Unset: []string{"REX_BLOCK"}}
	if err := ctx.Detach(p, filepath.Join(runLogDir(common), name+".log")); err != nil {
		return Errorf("start %s: %v", what, err)
	}
	return nil
}

// envRunLog names, to a run start --foreground the daemon detached, the
// start log under <git common dir>/rhiza/logs its output goes to. The run
// has no ID until the child draws one, so each start gets a log of its
// own, <change set>.start-<n>.log, which the child renames to <run>.log
// before it journals: a later start of the same change set never appends
// to an earlier run's log.
const envRunLog = "RHIZA_RUN_LOG"

// runLogDir is where run drivers' logs live.
func runLogDir(common string) string {
	return filepath.Join(common, "rhiza", "logs")
}

// daemonDetachStart detaches the run start --foreground args of change set
// cs into a start log of its own, named to the child through envRunLog.
func daemonDetachStart(ctx *Ctx, what, cs string, args []string) error {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return err
	}
	name := ""
	for n := 1; ; n++ {
		name = fmt.Sprintf("%s.start-%d", cs, n)
		if _, err := os.Lstat(filepath.Join(runLogDir(common), name+".log")); errors.Is(err, fs.ErrNotExist) {
			break
		} else if err != nil {
			return Errorf("start %s: %v", what, err)
		}
	}
	return daemonDetach(ctx, what, name, []string{envRunLog + "=" + name + ".log"}, args)
}

// claimStartLog renames the start log the daemon named in envRunLog, if
// any, to run rid's log, so the run's output so far and from here on is
// its own. It is best effort: a run with no log still runs.
func claimStartLog(ctx *Ctx, common, rid string) {
	if ctx.Getenv == nil {
		return
	}
	name := filepath.Base(strings.TrimSpace(ctx.Getenv(envRunLog)))
	if name == "." || name == "/" || !strings.HasSuffix(name, ".log") {
		return
	}
	dir := runLogDir(common)
	_ = os.Rename(filepath.Join(dir, name), filepath.Join(dir, rid+".log"))
}

// daemonStart starts the run of change set cs as a detached rhiza run start
// --foreground process, with permission mode mode.
func daemonStart(ctx *Ctx, cs string, autonomous bool, mode string) error {
	args := []string{"run", "start", "--foreground", cs, "--permission-mode", mode}
	if autonomous {
		args = append(args, "--autonomous")
	}
	return daemonDetachStart(ctx, "a run of "+cs, cs, args)
}

// daemonResume continues run id as a detached rhiza run resume --foreground
// process, with the permission mode the run started with.
func daemonResume(ctx *Ctx, id string) error {
	return daemonDetach(ctx, "run "+id, id, nil, []string{"run", "resume", "--foreground", id})
}

func init() {
	register(Operation{
		Name:     "daemon.status",
		ReadOnly: true,
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
	Reason string   `json:"reason" desc:"Why it was not started: no-milestone, no-leaves, not-ready, live, failed, shipped, held, unjournaled, no-slot, disk-low, usage-reserve or ship-owner"`
	Detail []string `json:"detail" desc:"The readiness reasons, when the reason is not-ready; the free space, when disk-low; the kept slot, when no-slot keeps one for a queued plan request; the reserve hold, when usage-reserve; else empty"`
}

// skipText is the words for a skip reason in the text status.
var skipText = map[string]string{
	daemon.SkipNoMilestone:  "serves no milestone",
	daemon.SkipNoLeaves:     "has no leaves",
	daemon.SkipNotReady:     "not ready",
	daemon.SkipLive:         "has a live run",
	daemon.SkipFailed:       "its latest run failed and waits in the inbox",
	daemon.SkipShipped:      "its latest run shipped it",
	daemon.SkipHeld:         "started, its run not journaled yet",
	daemon.SkipUnjournaled:  "its start wrote no journal",
	daemon.SkipNoSlot:       "waiting for a slot",
	daemon.SkipDiskLow:      "disk space is low",
	daemon.SkipUsageReserve: "usage reserve: held until the window resets",
	daemon.SkipPaused:       "paused",
	daemon.SkipShipOwner:    "ship is set to owner: only the owner starts its run",
}

// DaemonProjectSkip is a registered project the root daemon did not work at
// its last poll.
type DaemonProjectSkip struct {
	Project string `json:"project" desc:"The project's name in the root"`
	Reason  string `json:"reason" desc:"Why it was not worked: not initialised, not a git checkout or held by another daemon"`
	Detail  string `json:"detail,omitempty" desc:"The error behind the reason, when it has one"`
}

// DaemonStatusResponse is the daemon's state. While no daemon holds the
// lock, Running is false and the rest, if any, is what the last daemon left.
type DaemonStatusResponse struct {
	Mode            string              `json:"mode" enum:"service,tab" desc:"How this installation runs the daemon: service when daemon install wrote a systemd unit, else tab (the Rex tab 'rhiza daemon')"`
	Running         bool                `json:"running" desc:"True while a live daemon holds the repository's lock"`
	Process         string              `json:"process,omitempty" desc:"The daemon's process, as host:pid; absent when none ever ran"`
	Version         string              `json:"version,omitempty" desc:"The version of the binary the daemon runs, as recorded when it started"`
	Commit          string              `json:"commit,omitempty" desc:"The git commit the daemon's binary was built from, as recorded when it started; absent when the build carried none"`
	Started         string              `json:"started,omitempty" desc:"When the daemon started, RFC 3339 UTC"`
	LastPoll        string              `json:"last_poll,omitempty" desc:"When it last polled, RFC 3339 UTC"`
	Stopped         string              `json:"stopped,omitempty" desc:"When the last daemon stopped, RFC 3339 UTC; absent while it runs"`
	MaxRuns         int                 `json:"max_runs,omitempty" desc:"How many runs it keeps running at once"`
	Autonomous      bool                `json:"autonomous" desc:"True when it starts each run with --autonomous"`
	Web             string              `json:"web,omitempty" desc:"The web UI as the daemon's last poll recorded it from rhiza.web.port: its URL (http://127.0.0.1:<port>), or off; it is the recorded setting, the listener is not probed; absent when none was recorded"`
	Settings        []Setting           `json:"settings" desc:"The settings as its last poll read them, each with its source: default, git config or flag"`
	Runs            []DaemonRun         `json:"runs" desc:"The live runs (running or interrupted) as of its last poll"`
	LastStarted     *DaemonEvent        `json:"last_started,omitempty" desc:"The run it last started; absent if none"`
	LastResumed     *DaemonEvent        `json:"last_resumed,omitempty" desc:"The run it last resumed; absent if none"`
	Skipped         []DaemonSkip        `json:"skipped" desc:"The accepted change sets it did not start at its last poll, each with its reason"`
	SkippedProjects []DaemonProjectSkip `json:"skipped_projects" desc:"The registered projects the root daemon did not work at its last poll, each with its reason"`
	Hooks           []HookRow           `json:"hooks" desc:"The notify and event commands the daemon took at start, each with its source: daemon flag, git config or unset"`
	Unjournaled     []DaemonEvent       `json:"unjournaled" desc:"The runs it started that never wrote a journal, by change set and when; it starts these again 15 minutes after the start, or sooner when a run of the change set journals"`
}

// Human writes the status as text.
func (r DaemonStatusResponse) Human(w io.Writer) {
	if r.Running {
		fmt.Fprintf(w, "Daemon running, process %s\n", r.Process)
		fmt.Fprintf(w, "Mode: %s\n", r.Mode)
	} else {
		fmt.Fprintln(w, "Daemon not running")
		fmt.Fprintf(w, "Mode: %s\n", r.Mode)
		if r.Process == "" {
			return
		}
		fmt.Fprintf(w, "Last daemon: process %s, stopped %s\n", r.Process, r.Stopped)
	}
	if r.Version != "" {
		fmt.Fprintf(w, "Version: %s\n", r.Version)
	}
	if r.Commit != "" {
		fmt.Fprintf(w, "Commit: %s\n", r.Commit)
	}
	fmt.Fprintf(w, "Started: %s\nLast poll: %s\nMax runs: %d\nAutonomous: %t\n", r.Started, r.LastPoll, r.MaxRuns, r.Autonomous)
	if r.Web != "" {
		fmt.Fprintf(w, "Web UI: %s\n", r.Web)
	}
	for _, s := range r.Settings {
		if s.Key == "rhiza.daemon.paused" {
			fmt.Fprintf(w, "Paused: %s\n", map[bool]string{true: "yes", false: "no"}[s.Value == "true"])
		}
	}
	if len(r.Settings) > 0 {
		fmt.Fprintln(w, "Settings:")
		for _, s := range r.Settings {
			fmt.Fprintf(w, "- %s = %s (%s)\n", s.Key, s.Value, s.Source)
		}
	}
	for _, h := range r.Hooks {
		fmt.Fprintf(w, "Hook %s: %s\n", h.Name, h.text())
	}
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
	for _, k := range r.SkippedProjects {
		fmt.Fprintf(w, "Skipped project: %s: %s", k.Project, k.Reason)
		if k.Detail != "" {
			fmt.Fprintf(w, " (%s)", k.Detail)
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
	resp := DaemonStatusResponse{Settings: []Setting{}, Runs: []DaemonRun{}, Skipped: []DaemonSkip{}, SkippedProjects: []DaemonProjectSkip{}, Hooks: []HookRow{}, Unjournaled: []DaemonEvent{}}
	resp.Mode, _ = daemonMode(ctx)
	common, rooted, err := daemonStatusDir(ctx)
	if err != nil {
		return nil, err
	}
	if rooted != "" {
		// The root's lock says whether its daemon runs, wherever under the
		// root status is asked; the home's own lock says so too.
		var h daemon.Holder
		h, resp.Running = daemon.RunningRoot(rooted)
		if common == "" {
			if resp.Running {
				resp.Process = h.String()
			}
			return resp, nil
		}
		if !resp.Running {
			_, resp.Running = daemon.Running(common)
		}
	} else {
		_, resp.Running = daemon.Running(common)
	}
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
	resp.Commit = st.Commit
	resp.Started = st.Process.Started.UTC().Format(time.RFC3339)
	resp.LastPoll = st.LastPoll.UTC().Format(time.RFC3339)
	if !st.Stopped.IsZero() && !resp.Running {
		resp.Stopped = st.Stopped.UTC().Format(time.RFC3339)
	}
	resp.MaxRuns, resp.Autonomous = st.MaxRuns, st.Autonomous
	for _, v := range st.Settings {
		resp.Settings = append(resp.Settings, Setting{Key: v.Key, Value: v.Value, Source: v.Source, Default: settingDefault(v.Key), Scope: settingScope(v.Key)})
	}
	resp.Web = webStatus(resp.Settings)
	for _, h := range st.Hooks {
		resp.Hooks = append(resp.Hooks, HookRow{Name: h.Name, Command: h.Command, Source: h.Source})
	}
	for _, d := range st.Runs {
		resp.Runs = append(resp.Runs, DaemonRun{ID: d.ID, Change: d.Change, State: d.State})
	}
	resp.LastStarted = daemonEvent(st.LastStarted)
	resp.LastResumed = daemonEvent(st.LastResumed)
	for _, k := range st.Skips {
		resp.Skipped = append(resp.Skipped, DaemonSkip{Change: k.Change, Title: k.Title, Reason: k.Reason, Detail: append([]string{}, k.Detail...)})
	}
	for _, k := range st.ProjectSkips {
		resp.SkippedProjects = append(resp.SkippedProjects, DaemonProjectSkip{Project: k.Project, Reason: k.Reason, Detail: k.Detail})
	}
	for _, e := range st.Unjournaled {
		resp.Unjournaled = append(resp.Unjournaled, *daemonEvent(&e))
	}
	return resp, nil
}

// daemonStatusDir is the git common dir that holds the state status reads.
// Under a Rhiza root it is the root daemon's home project's (the one the
// daemon recorded, else the one daemon would pick from here), and root is
// the root; common is "" when the root has no initialised project. Outside
// a root it is this repo's, and root is "".
func daemonStatusDir(ctx *Ctx) (common, root string, err error) {
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		common, err = gitCommonDir(ctx)
		return common, "", err
	}
	reg, err := rootcfg.Load(root)
	if err != nil {
		return "", "", err
	}
	var home *Ctx
	if name := daemon.ReadHome(root); name != "" && slices.Contains(reg.Names(), name) && projectInitialised(reg.Dir(root, name)) {
		home = projectCtx(ctx, reg.Dir(root, name))
	} else if _, c, ok := daemonHome(ctx, root, reg); ok {
		home = c
	} else {
		return "", root, nil
	}
	common, err = gitCommonDir(home)
	return common, root, err
}

// webStatus is the web UI line for the recorded settings: its URL, off, or
// empty when the state recorded no valid web port.
func webStatus(rows []Setting) string {
	addr, err := web.Addr(settingValue(rows, web.HostKey), settingValue(rows, web.PortKey))
	switch {
	case err != nil || settingValue(rows, web.PortKey) == "":
		return ""
	case addr == "":
		return web.Off
	}
	return "http://" + addr
}

func daemonEvent(e *daemon.Event) *DaemonEvent {
	if e == nil {
		return nil
	}
	return &DaemonEvent{Change: e.Change, Run: e.Run, At: e.At.UTC().Format(time.RFC3339)}
}

// repoSessionName is rhiza-<project>, <project> being the name the project
// is registered under in the root's rhiza.toml (its folder, or project add
// --name): the Rex session every tab of this project's runs opens in. Outside
// a root it is the main checkout's folder name.
func repoSessionName(ctx *Ctx) (string, error) {
	if ctx.ProjectName != "" {
		return "rhiza-" + ctx.ProjectName, nil
	}
	if root, ok := rootcfg.Find(ctx.Root); ok {
		if reg, err := rootcfg.Load(root); err == nil {
			if name, ok := reg.ProjectOf(root, ctx.Root); ok {
				return "rhiza-" + name, nil
			}
		}
	}
	if ctx.Git == nil {
		return "", Errorf("daemon needs git: its Rex session is named after the main checkout")
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return "", Errorf("daemon needs a git repository: %v", err)
	}
	main, ok := firstWorktree(out)
	if !ok {
		return "", Errorf("git worktree list names no main checkout")
	}
	return "rhiza-" + filepath.Base(main), nil
}

// ensureRepoSession makes the Rex session rhiza-<repo> the one the daemon's
// tabs open in, whatever REX_SESSION says, creating it on the Rex server
// when `rex session ls` lacks it. It acts once per process.
func ensureRepoSession(ctx *Ctx) error {
	if ctx.repoSession {
		return nil
	}
	name, err := repoSessionName(ctx)
	if err != nil {
		return err
	}
	out, err := ctx.Exec("rex", []string{"session", "ls", "--json"}, ctx.Root, nil)
	if err != nil {
		return Errorf("rex session ls: %v", err)
	}
	var in struct {
		Sessions []struct {
			Label string `json:"label"`
		} `json:"sessions"`
	}
	if err := json.Unmarshal(out, &in); err != nil {
		return Errorf("rex session ls: unreadable JSON: %v", err)
	}
	if !slices.ContainsFunc(in.Sessions, func(s struct {
		Label string `json:"label"`
	}) bool {
		return s.Label == name
	}) {
		if _, err := ctx.Exec("rex", []string{"session", "new", name, "-c", ctx.Root, "--json"}, ctx.Root, nil); err != nil {
			return Errorf("rex session new %s: %v", name, err)
		}
	}
	ctx.RexSession, ctx.repoSession = name, true
	return nil
}

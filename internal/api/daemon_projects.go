package api

// This file holds what the root daemon keeps for each registered project:
// a context bound to the project's own checkout, and the state of the
// project's part of a poll (its journal, holds and state file).

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/dphaener/rhiza/internal/usage"
)

// projectInitialised says whether dir holds an initialised Rhiza project:
// a .rhiza folder with its repo.json.
func projectInitialised(dir string) bool {
	fi, err := os.Stat(filepath.Join(dir, store.Dir, store.RepoFile))
	return err == nil && !fi.IsDir()
}

// projectCtx is ctx made to work in the checkout dir: its own root, git,
// git config, model index and fold cache, and none of the caller's change
// set, actor or session.
func projectCtx(ctx *Ctx, dir string) *Ctx {
	c := *ctx
	c.Root, c.Wd = dir, dir
	c.ChangeSet, c.Actor, c.For = id.ID{}, id.ID{}, id.ID{}
	c.ChangeRef, c.ActorRef, c.ForRef = "", "", ""
	c.daemonFold, c.readLog, c.builtModel, c.repoSession = nil, nil, nil, false
	c.RexBlock = ""
	c.ProjectName = ""
	c.IndexPath = ""
	if p, err := index.DefaultPath(dir); err == nil {
		c.IndexPath = p
	}
	if ctx.Git != nil && ctx.Exec != nil {
		c.Git = func(args ...string) (string, error) {
			out, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
			return string(out), err
		}
		git := c.Git
		c.GitConfig = func(key string) (string, error) {
			out, err := git("config", "--get", key)
			v := strings.TrimSpace(out)
			if err != nil || v == "" {
				return "", Errorf("git config %s is not set", key)
			}
			return v, nil
		}
	}
	return &c
}

// daemonHome picks the project that holds the daemon's own state: the
// registered project the working directory belongs to, else the first
// initialised one. It returns the project's name and a context for it.
func daemonHome(ctx *Ctx, root string, reg rootcfg.Config) (string, *Ctx, bool) {
	name, ok := reg.ProjectOf(root, ctx.Wd)
	if !ok || !projectInitialised(reg.Dir(root, name)) {
		ok = false
		for _, n := range reg.Names() {
			if projectInitialised(reg.Dir(root, n)) {
				name, ok = n, true
				break
			}
		}
	}
	if !ok {
		return "", nil, false
	}
	dir := reg.Dir(root, name)
	if abs, err := filepath.Abs(ctx.Root); err == nil && abs == filepath.Clean(dir) {
		return name, ctx, true
	}
	return name, projectCtx(ctx, dir), true
}

// daemonProj is the daemon's working state for one project.
type daemonProj struct {
	name, dir, common string
	ctx               *Ctx
	rs                *run.Store
	pending           map[string]daemonPending
	unjournaled       map[string]daemonPending
	state             daemon.State
	scheds            daemonSchedules
	notifyCmd         string
	eventCmd          string
	// lock is the project's per-repo daemon lock, nil for the home project,
	// whose lock runDaemon holds.
	lock *daemon.Lock
}

// newDaemonProj starts the daemon's work in the project at dir, whose
// context is pctx. A project other than the home one gets its own per-repo
// lock, so status in it names this daemon; one held by another process
// is an error.
//
// conf, when not nil, is the settings the caller already read for it.
func newDaemonProj(pctx *Ctx, name, dir string, self daemon.Holder, home bool, r DaemonRequest, conf *daemonConf) (*daemonProj, error) {
	common, err := gitCommonDir(pctx)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		c, err := daemonConfOf(pctx, r)
		if err != nil {
			return nil, err
		}
		conf = &c
	}
	pctx.ProjectName = name
	p := &daemonProj{name: name, dir: dir, common: common, ctx: pctx}
	if !home {
		if p.lock, err = daemon.TryLock(common, pctx.Now()); err != nil {
			return nil, err
		}
		self = daemon.ReadHolder(common)
	}
	hooks := resolveHooks(pctx, r.NotifyCommand, r.EventCommand)
	p.notifyCmd, p.eventCmd = hookCommand(hooks, HookNotify), hookCommand(hooks, HookEvent)
	p.rs = run.NewStore(run.Dir(common))
	p.state, p.unjournaled = daemonStartState(common, self, pctx.Version, pctx.Commit, conf.maxRuns, r.Autonomous)
	p.pending = map[string]daemonPending{}
	p.state.Settings = conf.stateRows()
	p.state.Hooks = hooks
	// Write it before the first poll, so status matches the lock holder
	// at once, through a re-exec.
	daemon.WriteState(common, p.state)
	return p, nil
}

func (p *daemonProj) saveState() {
	p.state.LastPoll = p.ctx.Now().UTC()
	if runs, err := p.rs.List(); err == nil {
		host, _ := os.Hostname()
		p.state.Runs = daemon.LiveRuns(runs, host, nil)
	}
	daemon.WriteState(p.common, p.state)
}

// stop records that the daemon stopped working the project and releases
// the project's lock.
func (p *daemonProj) stop() {
	p.state.Stopped = p.ctx.Now().UTC()
	daemon.WriteState(p.common, p.state)
	p.lock.Unlock()
}

// daemonTurn is a project's part of a poll as read before any project
// starts anything: its settings, its usage hold and what it could start,
// for the machine-wide round-robin. err is a failure to read its settings.
type daemonTurn struct {
	conf daemonConf
	hold *usageHold
	pick daemon.ProjectPick
	err  error
}

// prepare reads the project's part of a poll (see daemonTurn). onConf,
// when set, is told the project's settings first.
func (p *daemonProj) prepare(r DaemonRequest, onConf func(daemonConf)) daemonTurn {
	ctx, state := p.ctx, &p.state
	if ctx.Reap != nil {
		ctx.Reap(p.common)
	}
	conf, err := daemonConfOf(ctx, r)
	if err != nil {
		p.saveState()
		return daemonTurn{err: err}
	}
	if onConf != nil {
		onConf(conf)
	}
	state.Settings = conf.stateRows()
	learnLimits(ctx, p.common, state)
	var hold *usageHold
	if runs, lerr := p.rs.List(); lerr == nil {
		host, _ := os.Hostname()
		var windows map[usage.Kind]usage.Window
		var live map[usage.Kind]int64
		windows, state.UsageWindows, live = usageWindowsNow(ctx, p.common, state.UsageWindows)
		ratchetLimits(ctx, p.common, state, windows)
		hold = usageHoldOf(conf.rows, state.Learned, live, runs, host, windows)
	}
	if hold == nil {
		hold = governorHold(ctx, p.common)
	}
	ready, running := daemonDemand(ctx, p.rs, p.pending, p.unjournaled, conf, hold)
	// A due schedule is demand for a slot after the rest; one started and
	// not journaled yet holds a slot.
	if runs, lerr := p.rs.List(); lerr == nil {
		p.scheds.logEnds(ctx, p.name, runs)
		running += p.scheds.startingCount(ctx.Now(), runs)
		if !p.schedulesStopped(conf, hold) {
			for range p.scheds.due(ctx, p.name, runs) {
				ready = append(ready, "schedule")
			}
		}
	}
	return daemonTurn{conf: conf, hold: hold, pick: daemon.ProjectPick{Project: p.name, Ready: ready, Running: running, Max: conf.projectMax}}
}

// schedulesStopped says a paused or held project, or one on a disk low on
// space, starts no scheduled run, as it starts nothing else.
func (p *daemonProj) schedulesStopped(conf daemonConf, hold *usageHold) bool {
	if conf.paused || hold != nil {
		return true
	}
	low, _ := daemonDiskLow(p.ctx, conf.minFree)
	return low
}

// startSchedules starts the due schedules' runs in the slots left of share.
func (p *daemonProj) startSchedules(conf daemonConf, hold *usageHold, free int) {
	if free <= 0 || p.schedulesStopped(conf, hold) {
		return
	}
	runs, err := p.rs.List()
	if err != nil {
		daemonLogf(p.ctx, "project %s: schedules not started: read runs: %v", p.name, err)
		return
	}
	for _, ds := range p.scheds.due(p.ctx, p.name, runs) {
		if free <= 0 {
			break
		}
		if err := p.scheds.start(p.ctx, p.name, ds, conf.mode); err != nil {
			daemonLogf(p.ctx, "project %s: schedule %s not started: %v", p.name, ds.sc.Name, err)
			continue
		}
		free--
	}
}

// daemonDemand is what the project of ctx could start now, in the order
// its own poll starts them, and how many run slots its runs hold. The
// owner's queued work comes first, a placeholder each (a continuation that
// could start, a queued start, a plan request that could start), then its ready
// change sets in roadmap order. A paused or held project, or one on a
// disk low on space, could start nothing. A slot is held by a running run,
// by a start not yet journaled, as daemon.Pick counts them, and by a plan
// request's start not yet journaled, as daemonPlanRequests counts them.
func daemonDemand(ctx *Ctx, rs *run.Store, pending, unjournaled map[string]daemonPending, conf daemonConf, hold *usageHold) ([]string, int) {
	runs, err := rs.List()
	if err != nil {
		return nil, 0
	}
	host, _ := os.Hostname()
	fold, ferr := daemonFoldOf(ctx)
	accepted := func(change string) bool {
		if ferr != nil {
			return false
		}
		st, ok := fold.stateOf(change)
		return ok && st == string(changeset.StateAccepted)
	}
	held := daemonHeld(pending, unjournaled, runs, accepted, ctx.Now())
	running := 0
	live := map[string]bool{}
	for _, j := range runs {
		st := j.State()
		switch st.Liveness(host, nil) {
		case run.LiveRunning:
			live[st.Change] = true
			running++
		case run.LiveInterrupted:
			live[st.Change] = true
		}
	}
	for _, cs := range held {
		if !live[cs] {
			running++
		}
	}
	running += daemonRequestsStarting(ctx, rs, runs)
	running += daemonResumesStarting(ctx, runs, host)
	if conf.paused || ferr != nil {
		return nil, running
	}
	low, _ := daemonDiskLow(ctx, conf.minFree)
	if low {
		return nil, running
	}
	// A queued milestone request's plan run and the daemon's resume of an
	// interrupted run each take a slot, first (a usage hold stops all but
	// the milestone request).
	var ready []string
	if hold == nil {
		host, _ := os.Hostname()
		for range daemonResumable(ctx, rs, runs, fold, host, false) {
			ready = append(ready, "resume")
		}
	}
	if daemonMilestoneWaits(ctx) {
		ready = append(ready, "milestone request")
	}
	if hold != nil {
		return ready, running
	}
	for range daemonContinuesWaiting(ctx, rs, runs) {
		ready = append(ready, "continue")
	}
	if common, err := gitCommonDir(ctx); err == nil {
		if starts, err := daemon.ReadStarts(common); err == nil {
			for _, s := range starts {
				if _, ok := pending[s.Change]; !ok && !live[s.Change] {
					ready = append(ready, "start")
				}
			}
		}
	}
	for range daemonRequestsWaiting(ctx) {
		ready = append(ready, "plan request")
	}
	// As in daemonPoll1, each candidate's ship stage is the one its
	// dispatched workflow declares: one whose ship stage is the owner's is
	// raised, not picked, so it claims no slot.
	ship := daemonShipOf(ctx.Root, fold)
	cands := slices.DeleteFunc(slices.Clone(fold.cands), func(cs string) bool {
		if _, ok := unjournaled[cs]; ok {
			return true
		}
		mode, err := workflowStageMode(ctx, ship.For(cs), "ship")
		return err != nil || mode == AutonomyOwner
	})
	ready = append(ready, daemon.Pick(cands, runs, held, host, nil, running+len(cands), ship)...)
	return ready, running
}

// poll is the rest of the project's poll, after prepare read t: what the
// per-repo daemon did at each tick, with share the run slots the
// machine-wide round-robin gave it. What the poll started is appended to
// resp.
func (p *daemonProj) poll(r DaemonRequest, resp *DaemonResponse, t daemonTurn, maxRuns, share int) error {
	ctx, state := p.ctx, &p.state
	if t.err != nil {
		return t.err
	}
	conf, hold := t.conf, t.hold
	minFree := conf.minFree
	state.MaxRuns = maxRuns
	var skips []daemon.Skip
	startRuns, _ := p.rs.List()
	started, resumed, planned, err := daemonPoll1(ctx, p.rs, p.pending, p.unjournaled, p.scheds.startingCount(ctx.Now(), startRuns), t.pick.Running+share, minFree, r.Autonomous, conf.mode, conf.paused, hold, &skips)
	resp.Planned = append(resp.Planned, planned...)
	if err == nil {
		// A schedule that does not start is in the daemon log; it never
		// fails the poll.
		p.startSchedules(conf, hold, share-len(started)-len(resumed)-len(planned))
	}
	if err == nil {
		state.Skips = skips
		state.UsageHold = usageHoldState(hold, state.UsageHold, ctx.Now())
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
	for cs, pd := range p.unjournaled {
		state.Unjournaled = append(state.Unjournaled, daemon.Event{Change: cs, At: pd.at.UTC(), Runs: pd.runs})
	}
	slices.SortFunc(state.Unjournaled, func(a, b daemon.Event) int { return strings.Compare(a.Change, b.Change) })
	if nerr := daemonNotify(ctx, p.rs, state, p.notifyCmd); err == nil {
		err = nerr
	}
	if pol := watchdogPolicyOf(t.conf.rows); pol.Enabled {
		if werr := daemonWatchdog(ctx, p.rs, state, pol); err == nil {
			err = werr
		}
	} else if nerr := daemonNudge(ctx, p.rs, state); err == nil {
		err = nerr
	}
	if eerr := daemonEvents(ctx, state, p.eventCmd); err == nil {
		err = eerr
	}
	p.saveState()
	return err
}

// daemonMilestoneWaits reports whether a queued milestone request could
// start now: one with no run whose retry time has passed, and no other
// milestone request's run in flight.
func daemonMilestoneWaits(ctx *Ctx) bool {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return false
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return false
	}
	now := ctx.Now()
	wait := false
	for _, q := range reqs {
		if q.Kind != daemon.KindMilestone {
			continue
		}
		if q.Run != "" {
			return false
		}
		wait = wait || !now.Before(q.RetryAt)
	}
	return wait
}

// daemonRequestsStarting counts the plan requests whose start has not
// journaled yet and is still within daemonStartGrace: each holds a slot
// (see daemonPlanRequests). An unreadable queue counts none.
func daemonRequestsStarting(ctx *Ctx, rs *run.Store, runs []run.Journal) int {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return 0
	}
	reqs, err := daemon.ReadRequests(common)
	if err != nil {
		return 0
	}
	n := 0
	for _, q := range reqs {
		if q.Run == "" || ctx.Now().Sub(q.StartedAt) >= daemonStartGrace {
			continue
		}
		if slices.ContainsFunc(runs, func(j run.Journal) bool { return j.ID == q.Run }) {
			continue
		}
		if _, err := rs.Read(q.Run); err == nil {
			continue
		}
		n++
	}
	return n
}

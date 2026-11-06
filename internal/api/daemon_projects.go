package api

// This file holds what the root daemon keeps for each registered project:
// a context bound to the project's own checkout, and the state of the
// project's part of a poll (its journal, holds and state file).

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

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

// poll is one poll of the project: what the per-repo daemon did at each
// tick. onConf, when set, is told the project's settings first. What the
// poll started is appended to resp.
func (p *daemonProj) poll(r DaemonRequest, resp *DaemonResponse, onConf func(daemonConf)) error {
	ctx, state := p.ctx, &p.state
	if ctx.Reap != nil {
		ctx.Reap(p.common)
	}
	conf, err := daemonConfOf(ctx, r)
	if err != nil {
		p.saveState()
		return err
	}
	if onConf != nil {
		onConf(conf)
	}
	minFree := conf.minFree
	state.MaxRuns, state.Settings = conf.maxRuns, conf.stateRows()
	learnLimits(ctx, p.common, state)
	var hold *usageHold
	if runs, lerr := p.rs.List(); lerr == nil {
		host, _ := os.Hostname()
		var windows map[usage.Kind]usage.Window
		windows, state.UsageWindows = usageWindowsNow(ctx, p.common, state.UsageWindows)
		hold = usageHoldOf(conf.rows, state.Learned, runs, host, windows)
	}
	if hold == nil {
		hold = governorHold(ctx, p.common)
	}
	var skips []daemon.Skip
	started, resumed, planned, err := daemonPoll1(ctx, p.rs, p.pending, p.unjournaled, conf.maxRuns, minFree, r.Autonomous, conf.mode, conf.paused, hold, &skips)
	resp.Planned = append(resp.Planned, planned...)
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
	if nerr := daemonNudge(ctx, p.rs, state); err == nil {
		err = nerr
	}
	if eerr := daemonEvents(ctx, state, p.eventCmd); err == nil {
		err = eerr
	}
	p.saveState()
	return err
}

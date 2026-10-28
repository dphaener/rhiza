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
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
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
	MaxRuns int  `json:"max_runs,omitempty" desc:"Start runs until this many are running; default 2. An interrupted run takes no slot"`
	Once    bool `json:"once,omitempty" desc:"Poll once, start what is picked and stop"`
}

// DaemonResponse is the daemon after it stopped.
type DaemonResponse struct {
	Process string   `json:"process" desc:"The daemon's process, as host:pid"`
	Started []string `json:"started" desc:"The change sets whose runs the daemon started, in order"`
}

// Human writes the response as markdown.
func (r DaemonResponse) Human(w io.Writer) {
	for _, cs := range r.Started {
		fmt.Fprintf(w, "Started a run of %s\n", cs)
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
	defer lock.Unlock()
	self := daemon.ReadHolder(common)
	if ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon running, process %s\n", self)
	}
	rs := run.NewStore(run.Dir(common))
	resp := DaemonResponse{Process: self.String(), Started: []string{}}
	pending := map[string]time.Time{}
	poll := func() error {
		started, err := daemonPoll1(ctx, rs, pending, r.MaxRuns)
		resp.Started = append(resp.Started, started...)
		return err
	}
	if r.Once {
		return resp, poll()
	}
	var interrupted <-chan struct{}
	if ctx.Interrupts != nil {
		var stop func()
		interrupted, stop = ctx.Interrupts()
		defer stop()
	}
	tick := time.NewTicker(daemonPoll)
	defer tick.Stop()
	for {
		if err := poll(); err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "rhiza: daemon: %v\n", err)
		}
		select {
		case <-interrupted:
			return resp, nil
		case <-tick.C:
		}
	}
}

// daemonStartGrace is how long a change set the daemon started is held as
// running while its run has no journal yet.
const daemonStartGrace = time.Minute

// daemonPoll1 picks and starts runs once. pending holds the change sets
// started and when; entries are dropped once their run has a journal or
// the grace passed.
func daemonPoll1(ctx *Ctx, rs *run.Store, pending map[string]time.Time, max int) ([]string, error) {
	res, _, _, err := planRead(ctx, "daemon", "")
	if err != nil {
		return nil, err
	}
	ws, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	cands := daemonCandidates(res, func(cs id.ID) bool {
		c, ok := ws.changeSet(cs)
		return ok && len(changeLeaves(ws.model(cs), c.Ops)) > 0
	})
	runs, err := rs.List()
	if err != nil {
		return nil, Errorf("read runs: %v", err)
	}
	var held []string
	for cs, at := range pending {
		if ctx.Now().Sub(at) >= daemonStartGrace {
			delete(pending, cs)
			continue
		}
		held = append(held, cs)
	}
	host, _ := os.Hostname()
	picked := daemon.Pick(cands, runs, held, host, nil, max)
	var started []string
	for _, cs := range picked {
		if err := daemonStart(ctx, cs); err != nil {
			return started, err
		}
		pending[cs] = ctx.Now()
		started = append(started, cs)
	}
	return started, nil
}

// daemonCandidates are the accepted, ready change sets of the milestone
// groups that hasLeaves, in roadmap order, each once.
func daemonCandidates(res planResult, hasLeaves func(id.ID) bool) []string {
	var out []string
	seen := map[id.ID]bool{}
	for _, g := range res.groups {
		if g.backlog {
			continue
		}
		for _, p := range g.members {
			if seen[p.info.ID] || !p.ready || p.info.State != changeset.StateAccepted {
				continue
			}
			seen[p.info.ID] = true
			if hasLeaves(p.info.ID) {
				out = append(out, p.info.ID.String())
			}
		}
	}
	return out
}

// daemonStart starts the run of change set cs in a new Rex tab.
func daemonStart(ctx *Ctx, cs string) error {
	if ctx.Exec == nil || ctx.Self == "" {
		return Errorf("daemon cannot start a run of %s: no rex or rhiza path", cs)
	}
	args := []string{"window", "new"}
	if ctx.RexSession != "" {
		args = append(args, "-s", ctx.RexSession)
	}
	args = append(args, cs[:min(len(cs), 9)], "-c", ctx.Root, "--json", "--", ctx.Self, "run", "start", cs)
	if _, err := ctx.Exec("rex", args, ctx.Root, nil); err != nil {
		return Errorf("rex window new for %s: %v", cs, err)
	}
	return nil
}

package api

// This file holds the daemon operation: rhiza daemon runs in the
// foreground, one per repository. It holds a lock under the git common dir,
// so a second daemon started from any worktree of the repository is
// refused, naming the first one's process. It binds no socket.

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
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
type DaemonRequest struct{}

// DaemonResponse is the daemon after it stopped.
type DaemonResponse struct {
	Process string `json:"process" desc:"The daemon's process, as host:pid"`
}

// Human writes the response as markdown.
func (r DaemonResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Daemon %s stopped\n", r.Process)
}

// runDaemon takes the repository's daemon lock and keeps it until the
// caller interrupts (Ctx.Interrupt). A second daemon exits with a usage
// error naming the holder.
func runDaemon(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("daemon runs in the foreground and has no dry run")
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
	var interrupted <-chan struct{}
	if ctx.Interrupts != nil {
		var stop func()
		interrupted, stop = ctx.Interrupts()
		defer stop()
	}
	tick := time.NewTicker(daemonPoll)
	defer tick.Stop()
	for {
		select {
		case <-interrupted:
			return DaemonResponse{Process: self.String()}, nil
		case <-tick.C:
		}
	}
}

package api

// This file holds the daemon.start and daemon.stop operations: they start
// or stop this repository's daemon the way its installation runs it. With
// the systemd unit that daemon.install wrote they call systemctl --user;
// without one the daemon runs in the Rex tab "rhiza daemon" of the repo's
// session, and stop signals its process. daemon.status names the mode.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/dphaener/rhiza/internal/daemon"
)

// DaemonTabName is the name of the Rex tab the daemon runs in.
const DaemonTabName = "rhiza daemon"

// The modes a daemon runs in.
const (
	DaemonModeService = "service"
	DaemonModeTab     = "tab"
)

// DaemonStartRequest asks for this repo's daemon started.
type DaemonStartRequest struct{}

// DaemonStartResponse says how the daemon was started.
type DaemonStartResponse struct {
	Mode    string `json:"mode" enum:"service,tab" desc:"How this installation runs the daemon: service (a systemd user unit) or tab (the Rex tab 'rhiza daemon')"`
	Changed bool   `json:"changed" desc:"False when the daemon already ran and nothing was done"`
}

// DaemonStopRequest asks for this repo's daemon stopped.
type DaemonStopRequest struct{}

// DaemonStopResponse says how the daemon was stopped.
type DaemonStopResponse struct {
	Mode    string `json:"mode" enum:"service,tab" desc:"How this installation runs the daemon: service (a systemd user unit) or tab (the Rex tab 'rhiza daemon')"`
	Changed bool   `json:"changed" desc:"False when no daemon ran and nothing was done"`
}

func init() {
	register(Operation{
		Name:     "daemon.start",
		Summary:  "Start this repo's daemon: through its systemd unit when daemon install wrote one, else in the Rex tab 'rhiza daemon'",
		Request:  reflect.TypeFor[DaemonStartRequest](),
		Response: reflect.TypeFor[DaemonStartResponse](),
		Run:      runDaemonStart,
	})
	register(Operation{
		Name:     "daemon.stop",
		Summary:  "Stop this repo's daemon: through its systemd unit when daemon install wrote one, else by ending its process",
		Request:  reflect.TypeFor[DaemonStopRequest](),
		Response: reflect.TypeFor[DaemonStopResponse](),
		Run:      runDaemonStop,
	})
}

// daemonMode is service when the main checkout's unit exists and daemon
// install wrote it, else tab. It also returns the unit's name.
func daemonMode(ctx *Ctx) (mode, unit string) {
	if ctx.Home == "" || ctx.Git == nil {
		return DaemonModeTab, ""
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return DaemonModeTab, ""
	}
	main, ok := firstWorktree(out)
	if !ok {
		return DaemonModeTab, ""
	}
	unit = serviceUnitName(main)
	have, exists, err := readIfExists(filepath.Join(ctx.Home, ".config", "systemd", "user", unit))
	if err != nil || !exists || !bytes.HasPrefix(have, []byte(serviceMarker)) {
		return DaemonModeTab, ""
	}
	return DaemonModeService, unit
}

func daemonRunning(ctx *Ctx) (daemon.Holder, bool, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return daemon.Holder{}, false, err
	}
	h, running := daemon.Running(common)
	if running {
		h = daemon.ReadHolder(common)
	}
	return h, running, nil
}

func runDaemonStart(ctx *Ctx, _ any) (any, error) {
	const cmd = "rhiza daemon start"
	if ctx.DryRun {
		return nil, Usagef("%s does not take --dry-run: it writes no log entry", cmd).WithNext(cmd)
	}
	_, running, err := daemonRunning(ctx)
	if err != nil {
		return nil, err
	}
	mode, unit := daemonMode(ctx)
	resp := DaemonStartResponse{Mode: mode, Changed: !running}
	if ctx.Exec == nil {
		return nil, Errorf("%s: cannot run programs", cmd)
	}
	if mode == DaemonModeService {
		if err := serviceExec(ctx, cmd, "systemctl", "--user", "start", unit); err != nil {
			return nil, err
		}
		return resp, nil
	}
	if running {
		return resp, nil
	}
	if ctx.Self == "" {
		return nil, Errorf("%s: cannot tell where rhiza is installed", cmd)
	}
	if err := ensureRepoSession(ctx); err != nil {
		return nil, err
	}
	args := []string{"window", "new", "-s", ctx.RexSession, DaemonTabName, "-c", ctx.Root, "--keep-open", "--json", "--", ctx.Self, "daemon", "--autonomous"}
	if _, err := ctx.Exec("rex", args, ctx.Root, nil); err != nil {
		return nil, Errorf("rex window new for the daemon: %v", err)
	}
	return resp, nil
}

func runDaemonStop(ctx *Ctx, _ any) (any, error) {
	const cmd = "rhiza daemon stop"
	if ctx.DryRun {
		return nil, Usagef("%s does not take --dry-run: it writes no log entry", cmd).WithNext(cmd)
	}
	h, running, err := daemonRunning(ctx)
	if err != nil {
		return nil, err
	}
	mode, unit := daemonMode(ctx)
	resp := DaemonStopResponse{Mode: mode, Changed: running}
	if ctx.Exec == nil {
		return nil, Errorf("%s: cannot run programs", cmd)
	}
	if mode == DaemonModeService {
		if err := serviceExec(ctx, cmd, "systemctl", "--user", "stop", unit); err != nil {
			return nil, err
		}
		return resp, nil
	}
	if !running {
		return resp, nil
	}
	if host, _ := os.Hostname(); host != h.Host {
		return nil, Errorf("%s: the daemon runs on %s, not this machine", cmd, h.Host)
	}
	if h.Pid <= 0 {
		return nil, Errorf("%s: the daemon's process is unknown", cmd)
	}
	if _, err := ctx.Exec("kill", []string{"-TERM", strconv.Itoa(h.Pid)}, ctx.Root, nil); err != nil {
		return nil, Errorf("%s: kill: %v", cmd, err)
	}
	return resp, nil
}

// Human reports what start did.
func (r DaemonStartResponse) Human(w io.Writer) {
	switch {
	case r.Mode == DaemonModeService:
		fmt.Fprintln(w, "Started the daemon through its systemd unit.")
	case !r.Changed:
		fmt.Fprintln(w, "The daemon already runs.")
	default:
		fmt.Fprintf(w, "Started the daemon in the Rex tab %s.\n", DaemonTabName)
	}
}

// Human reports what stop did.
func (r DaemonStopResponse) Human(w io.Writer) {
	switch {
	case r.Mode == DaemonModeService:
		fmt.Fprintln(w, "Stopped the daemon through its systemd unit.")
	case !r.Changed:
		fmt.Fprintln(w, "No daemon runs.")
	default:
		fmt.Fprintln(w, "Asked the daemon to stop.")
	}
}

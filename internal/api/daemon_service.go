package api

// This file holds the daemon.install and daemon.uninstall operations: they
// run the repository's daemon as a systemd user service, so a server
// installation needs no Rex tab for it. The unit lives in
// $HOME/.config/systemd/user, named for the repository, and runs
// `rhiza daemon --autonomous` in the main checkout, as the Rex tab does.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// serviceMarker is the first line of a unit daemon.install wrote; a file
// without it is the owner's and is never overwritten or removed.
const serviceMarker = "# Written by rhiza daemon install."

// DaemonInstallRequest asks for the daemon installed as a systemd user service.
type DaemonInstallRequest struct{}

// DaemonInstallResponse says where the unit is and what changed.
type DaemonInstallResponse struct {
	Unit    string `json:"unit" desc:"The unit's name, such as rhiza-myrepo.service"`
	Path    string `json:"path" desc:"The unit file, absolute"`
	Written bool   `json:"written" desc:"Whether the unit file changed; false if it was already current"`
}

// DaemonUninstallRequest asks for the systemd user service removed.
type DaemonUninstallRequest struct{}

// DaemonUninstallResponse says what uninstall removed.
type DaemonUninstallResponse struct {
	Unit    string `json:"unit" desc:"The unit's name, such as rhiza-myrepo.service"`
	Path    string `json:"path" desc:"The unit file, absolute"`
	Removed bool   `json:"removed" desc:"Whether a unit file was removed; false if there was none"`
}

func init() {
	register(Operation{
		Name:     "daemon.install",
		Summary:  "Run this repo's daemon as a systemd user service: write and enable the unit and turn on lingering",
		Request:  reflect.TypeFor[DaemonInstallRequest](),
		Response: reflect.TypeFor[DaemonInstallResponse](),
		Run:      runDaemonInstall,
	})
	register(Operation{
		Name:     "daemon.uninstall",
		Summary:  "Stop, disable and remove this repo's systemd user service",
		Request:  reflect.TypeFor[DaemonUninstallRequest](),
		Response: reflect.TypeFor[DaemonUninstallResponse](),
		Run:      runDaemonUninstall,
	})
}

// serviceUnitName is rhiza-<repo>.service, <repo> being the main checkout's
// folder name with every rune a unit name does not allow turned into '-'.
func serviceUnitName(main string) string {
	repo := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '.', r == '-':
			return r
		}
		return '-'
	}, filepath.Base(main))
	return "rhiza-" + repo + ".service"
}

// serviceQuote quotes an ExecStart word the way systemd reads one.
func serviceQuote(s string) string {
	if !strings.ContainsAny(s, " \t\"'\\$%;") {
		return s
	}
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, `%`, `%%`, `$`, `$$`)
	return `"` + r.Replace(s) + `"`
}

// serviceEnvNames are the variables of the installing shell the unit carries,
// so the service starts agents with the Rex tab's tools and temp dir.
var serviceEnvNames = []string{"PATH", "TMPDIR"}

// serviceEnvQuote quotes an Environment= assignment the way systemd reads it:
// the whole NAME=value in double quotes, backslash and double quote escaped,
// % doubled and $ left as is (unlike an ExecStart word).
func serviceEnvQuote(name, value string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, `%`, `%%`)
	return `"` + r.Replace(name+"="+value) + `"`
}

// serviceEnvLines returns the Environment= lines for the variables that are set.
func serviceEnvLines(getenv func(string) string) string {
	var b strings.Builder
	for _, n := range serviceEnvNames {
		if v := getenv(n); v != "" {
			b.WriteString("Environment=" + serviceEnvQuote(n, v) + "\n")
		}
	}
	return b.String()
}

// serviceUnit returns the unit file for the daemon of the repository whose
// main checkout is main, run by the rhiza at bin, with env lines env.
func serviceUnit(main, bin, env string) []byte {
	return []byte(serviceMarker + "\n" +
		"[Unit]\n" +
		"Description=Rhiza daemon for " + filepath.Base(main) + "\n" +
		"After=rex.service\n\n" +
		"[Service]\n" +
		"WorkingDirectory=" + main + "\n" +
		env +
		"ExecStart=" + serviceQuote(bin) + " " + strings.Join(daemonArgs, " ") + "\n" +
		"Restart=always\n\n" +
		"[Install]\n" +
		"WantedBy=default.target\n")
}

// serviceSetup finds what both operations need: systemd (a usage error
// naming it when systemctl is missing), the unit's name and file, and the
// main checkout.
func serviceSetup(ctx *Ctx, cmd string) (unit, path, main string, err error) {
	if ctx.DryRun {
		return "", "", "", Usagef("%s does not take --dry-run: it writes no log entry", cmd).WithNext(cmd)
	}
	if ctx.LookPath == nil {
		return "", "", "", Usagef("%s needs systemd: systemctl is not on PATH", cmd)
	}
	for _, p := range []string{"systemctl", "loginctl"} {
		if _, lerr := ctx.LookPath(p); lerr != nil {
			return "", "", "", Usagef("%s needs systemd: %s is not on PATH", cmd, p)
		}
	}
	if ctx.Exec == nil {
		return "", "", "", Errorf("%s: cannot run systemctl", cmd)
	}
	if ctx.Home == "" {
		return "", "", "", Usagef("%s: no home directory: set HOME", cmd)
	}
	if ctx.Git == nil {
		return "", "", "", Errorf("%s needs git: the daemon runs in the main checkout", cmd)
	}
	out, gerr := ctx.Git("worktree", "list", "--porcelain")
	if gerr != nil {
		return "", "", "", Errorf("%s needs a git repository: %v", cmd, gerr)
	}
	main, ok := firstWorktree(out)
	if !ok {
		return "", "", "", Errorf("%s: git worktree list names no main checkout", cmd)
	}
	unit = serviceUnitName(main)
	path = filepath.Join(ctx.Home, ".config", "systemd", "user", unit)
	return unit, path, main, nil
}

func serviceExec(ctx *Ctx, cmd, name string, args ...string) error {
	if _, err := ctx.Exec(name, args, "", nil); err != nil {
		return Errorf("%s: %v", cmd, err)
	}
	return nil
}

// runDaemonInstall writes the unit (unless it is already current), reloads
// systemd, enables and starts the service and turns on lingering, so the
// service runs without a login. A unit not written by install is a
// conflict. A second install writes and reloads nothing; the enabling
// calls are idempotent and run again, so an install that stopped halfway
// is finished by the next.
func runDaemonInstall(ctx *Ctx, _ any) (any, error) {
	const cmd = "rhiza daemon install"
	unit, path, main, err := serviceSetup(ctx, cmd)
	if err != nil {
		return nil, err
	}
	bin := ctx.Self
	if bin == "" {
		return nil, Errorf("%s: cannot tell where rhiza is installed", cmd)
	}
	want := serviceUnit(main, bin, serviceEnvLines(ctxGetenv(ctx)))
	have, exists, err := readIfExists(path)
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err)
	}
	if exists && !bytes.HasPrefix(have, []byte(serviceMarker)) {
		return nil, Conflictf("%s exists and was not written by rhiza daemon install", path)
	}
	resp := DaemonInstallResponse{Unit: unit, Path: path}
	if !exists || !bytes.Equal(have, want) {
		if err := writeAtomic(path, want); err != nil {
			return nil, Errorf("%s: %v", cmd, err)
		}
		resp.Written = true
		if err := serviceExec(ctx, cmd, "systemctl", "--user", "daemon-reload"); err != nil {
			return nil, err
		}
	}
	if err := serviceExec(ctx, cmd, "systemctl", "--user", "enable", "--now", unit); err != nil {
		return nil, err
	}
	if err := serviceExec(ctx, cmd, "loginctl", "enable-linger"); err != nil {
		return nil, err
	}
	return resp, nil
}

// runDaemonUninstall stops and disables the service, removes its unit and
// reloads systemd. Without a unit file it does nothing. Lingering stays on:
// other services of the owner may need it.
func runDaemonUninstall(ctx *Ctx, _ any) (any, error) {
	const cmd = "rhiza daemon uninstall"
	unit, path, _, err := serviceSetup(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := DaemonUninstallResponse{Unit: unit, Path: path}
	have, exists, err := readIfExists(path)
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err)
	}
	if !exists {
		return resp, nil
	}
	if !bytes.HasPrefix(have, []byte(serviceMarker)) {
		return nil, Conflictf("%s was not written by rhiza daemon install: remove it yourself", path)
	}
	if err := serviceExec(ctx, cmd, "systemctl", "--user", "disable", "--now", unit); err != nil {
		return nil, err
	}
	if err := os.Remove(path); err != nil {
		return nil, Errorf("%s: %v", cmd, err)
	}
	resp.Removed = true
	if err := serviceExec(ctx, cmd, "systemctl", "--user", "daemon-reload"); err != nil {
		return nil, err
	}
	return resp, nil
}

// Human reports what install did.
func (r DaemonInstallResponse) Human(w io.Writer) {
	if r.Written {
		fmt.Fprintf(w, "Installed %s (%s); the daemon runs as a systemd user service.\n", r.Unit, r.Path)
		return
	}
	fmt.Fprintf(w, "%s is already installed (%s); nothing changed.\n", r.Unit, r.Path)
}

// Human reports what uninstall did.
func (r DaemonUninstallResponse) Human(w io.Writer) {
	if !r.Removed {
		fmt.Fprintf(w, "No %s to remove.\n", r.Unit)
		return
	}
	fmt.Fprintf(w, "Stopped, disabled and removed %s (%s).\n", r.Unit, r.Path)
}

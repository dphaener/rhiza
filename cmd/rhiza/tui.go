package main

// `rhiza tui launch` is the program the quick-launch layer runs. It is the
// one command here that is a terminal program rather than an operation, so
// it is not in the api registry; the other rhiza apps are reached through
// it.

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/client"
	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/readapps"
	"github.com/dphaener/rhiza/internal/tui/writeapps"
)

func tuiCommand(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	if len(args) >= 3 && args[1] == "app" {
		return tuiApp(g, args[2], args[3:], stdout, stderr, getenv)
	}
	if len(args) != 2 || args[1] != "launch" {
		fmt.Fprintln(stderr, "usage: rhiza tui launch | rhiza tui app <name> [args] [--once]")
		return 2
	}
	session := getenv("REX_SESSION")
	if session == "" {
		return writeError(stdout, stderr, api.Usagef("tui launch runs inside Rex: REX_SESSION is not set"), false)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rex, err := client.Dial(ctx, client.Options{Getenv: getenv})
	if err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", err), false)
	}
	defer rex.Close()
	// surface.New needs the connection attached to the session.
	if _, err := rex.Attach(ctx, session); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", err), false)
	}
	cancel()
	cwd, _ := os.Getwd()
	l := launcher.Launcher{Rex: rex, Session: session, Bin: selfPath(), Cwd: cwd}
	m := launcher.NewModel(context.Background(), l, launcher.Apps(), colorprofile.Detect(os.Stdout, os.Environ()))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", err), false)
	}
	if m.Err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", m.Err), false)
	}
	return 0
}

// tuiApp runs an app: `rhiza tui app <name> [args]`. It calls the api
// operation of the app in-process, as the matching command does. With --once
// it prints the app's result and exits; otherwise it shows that result in a
// terminal pane. The global flags in g reach the operation as they do the
// command's.
//
// A write app shows a what-if preview first, the operation run with
// --dry-run. With --once it prints that preview and stops, so nothing is
// written; --once --confirm confirms it: the operation runs for real and
// prints what the command prints with --json, byte for byte.
func tuiApp(g globals, name string, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	read, isRead := readapps.Lookup(name)
	write, isWrite := writeapps.Lookup(name)
	if !isRead && !isWrite {
		return writeError(stdout, stderr, api.NotFoundf("unknown app %q", name).WithNext("rhiza tui launch"), false)
	}
	once, confirm := false, false
	var rest []string
	for _, a := range args {
		switch a {
		case "--once":
			once = true
		case "--confirm":
			confirm = true
		default:
			rest = append(rest, a)
		}
	}
	call := func(op string, record, dry bool, args []string) ([]byte, error) {
		o, ok := api.Lookup(op)
		if !ok {
			return nil, api.NotFoundf("unknown operation %q", op)
		}
		req, err := requestFromArgs(o, args)
		if err != nil {
			return nil, err
		}
		gg := g
		gg.record = record
		gg.dryRun = g.dryRun || dry
		ctx, err := newCtx(gg, os.Stdin, getenv)
		if err != nil {
			return nil, err
		}
		resp, err := o.Run(ctx, req)
		if err != nil {
			return nil, err
		}
		return canon.Marshal(resp)
	}
	var m tea.Model
	var result func() ([]byte, error)
	if isWrite {
		wm := writeapps.NewModel(write, func(a writeapps.App, args []string, dry bool) ([]byte, error) {
			return call(a.Op, a.Record, dry, args)
		}, rest)
		m = wm
		result = func() ([]byte, error) {
			out, err := wm.Preview()
			if err == nil && confirm && !write.ReadOnly {
				out, err = wm.Confirm()
			}
			return out, err
		}
	} else {
		rm := readapps.NewModel(read, func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, rest)
		m, result = rm, rm.Result
	}
	if once {
		out, err := result()
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		stdout.Write(out)
		return 0
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	return 0
}

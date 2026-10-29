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
)

func tuiCommand(args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	if len(args) >= 3 && args[1] == "app" {
		return tuiApp(args[2], args[3:], stdout, stderr, getenv)
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
	rex, err := client.Dial(ctx, client.Options{Getenv: getenv})
	cancel()
	if err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", err), false)
	}
	defer rex.Close()
	cwd, _ := os.Getwd()
	l := launcher.Launcher{Rex: rex, Session: getenv("REX_SESSION"), Bin: selfPath(), Cwd: cwd}
	m := launcher.NewModel(context.Background(), l, launcher.Apps(), colorprofile.Detect(os.Stdout, os.Environ()))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", err), false)
	}
	if m.Err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui launch: %w", m.Err), false)
	}
	return 0
}

// tuiApp runs a read app: `rhiza tui app <name> [args]`. It calls the
// api operation of the same name in-process, as `rhiza <name> --json`
// does. With --once it prints that JSON and exits; otherwise it shows the
// result in a terminal pane.
func tuiApp(name string, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	app, ok := readapps.Lookup(name)
	if !ok {
		return writeError(stdout, stderr, api.NotFoundf("unknown app %q", name).WithNext("rhiza tui launch"), false)
	}
	once := false
	var rest []string
	for _, a := range args {
		if a == "--once" {
			once = true
		} else {
			rest = append(rest, a)
		}
	}
	run := func(op string, args []string) ([]byte, error) {
		o, ok := api.Lookup(op)
		if !ok {
			return nil, api.NotFoundf("unknown operation %q", op)
		}
		req, err := requestFromArgs(o, args)
		if err != nil {
			return nil, err
		}
		ctx, err := newCtx(globals{}, os.Stdin, getenv)
		if err != nil {
			return nil, err
		}
		resp, err := o.Run(ctx, req)
		if err != nil {
			return nil, err
		}
		return canon.Marshal(resp)
	}
	if once {
		out, err := run(app.Name, rest)
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		stdout.Write(out)
		return 0
	}
	if _, err := tea.NewProgram(readapps.NewModel(app, run, rest)).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	return 0
}

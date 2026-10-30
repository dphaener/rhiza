package main

// `rhiza tui app <name>` is the program each app's Rex action runs. It is
// the one command here that is a terminal program rather than an operation,
// so it is not in the api registry.

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/inboxapp"
	"github.com/dphaener/rhiza/internal/tui/modelapp"
	"github.com/dphaener/rhiza/internal/tui/readapps"
	"github.com/dphaener/rhiza/internal/tui/writeapps"
)

func tuiCommand(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	if len(args) >= 3 && args[1] == "app" {
		return tuiApp(g, args[2], args[3:], stdout, stderr, getenv)
	}
	fmt.Fprintln(stderr, "usage: rhiza tui app <name> [args] [--once]")
	return 2
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
	if name == modelapp.Name {
		return tuiModel(g, args, stdout, stderr, getenv)
	}
	read, isRead := readapps.Lookup(name)
	isInbox := name == inboxapp.Name
	write, isWrite := writeapps.Lookup(name)
	if !isRead && !isWrite && !isInbox {
		return writeError(stdout, stderr, api.NotFoundf("unknown app %q", name), false)
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
	switch {
	case isInbox:
		im := inboxapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		})
		if ctx, err := newCtx(g, os.Stdin, getenv); err == nil {
			im.Exec, im.Self, im.Root, im.RexSession = ctx.Exec, ctx.Self, ctx.Root, ctx.RexSession
		}
		m, result = im, im.Result
	case isWrite:
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
	default:
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

// tuiModel runs the model app: the as-is tree and the roadmap, each the
// text its command prints; a ref limits the tree to that node's subtree.
// With --once it prints the tree, then the roadmap
// after a blank line, and exits.
func tuiModel(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	once, ref := false, ""
	for _, a := range args {
		switch {
		case a == "--once":
			once = true
		case ref == "" && !strings.HasPrefix(a, "-"):
			ref = a
		default:
			return writeError(stdout, stderr, api.Usagef("tui app model takes [<ref>] [--once]"), false)
		}
	}
	read := func(op string, args []string) (string, error) {
		o, ok := api.Lookup(op)
		if !ok {
			return "", api.NotFoundf("unknown operation %q", op)
		}
		req, err := requestFromArgs(o, args)
		if err != nil {
			return "", err
		}
		ctx, err := newCtx(g, os.Stdin, getenv)
		if err != nil {
			return "", err
		}
		resp, err := o.Run(ctx, req)
		if err != nil {
			return "", err
		}
		var b strings.Builder
		if err := writeResponse(&b, resp, false); err != nil {
			return "", err
		}
		return b.String(), nil
	}
	m := modelapp.NewModel(read, ref)
	if err := m.Err(); err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if once {
		io.WriteString(stdout, m.Text("tree")+"\n"+m.Text("roadmap"))
		return 0
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	return 0
}

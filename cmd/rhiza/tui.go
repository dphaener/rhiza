package main

// `rhiza tui app <name>` is the program each app's Rex action runs. It is
// the one command here that is a terminal program rather than an operation,
// so it is not in the api registry.

import (
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/autonomyapp"
	"github.com/dphaener/rhiza/internal/tui/daemonapp"
	"github.com/dphaener/rhiza/internal/tui/inboxapp"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
	"github.com/dphaener/rhiza/internal/tui/requestapp"
	"github.com/dphaener/rhiza/internal/tui/runsapp"
)

func tuiCommand(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	if len(args) >= 3 && args[1] == "open" {
		return tuiOpen(args[2], args[3:], stdout, stderr, getenv)
	}
	if len(args) >= 3 && args[1] == "app" {
		// The opener names a file for the app to report its content's
		// size to, so the modal follows it (modalshell.SizeFileFlag).
		rest := make([]string, 0, len(args)-3)
		for i := 3; i < len(args); i++ {
			if args[i] == modalshell.SizeFileFlag && i+1 < len(args) {
				modalshell.ReportSizeTo(args[i+1])
				i++
				continue
			}
			rest = append(rest, args[i])
		}
		return tuiApp(g, args[2], rest, stdout, stderr, getenv)
	}
	fmt.Fprintln(stderr, "usage: rhiza tui app <name> [args] [--once]\n       "+tuiOpenUsage)
	return 2
}

// tuiApp runs an app: `rhiza tui app <name> [args]`. It calls the api
// operation of the app in-process, as the matching command does. With --once
// it prints the app's result and exits; otherwise it shows that result in a
// terminal pane. The global flags in g reach the operation as they do the
// command's.
func tuiApp(g globals, name string, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	isInbox := name == inboxapp.Name
	isRuns := name == runsapp.Name
	isRequests := name == requestapp.Name
	isAutonomy := name == autonomyapp.Name
	isDaemon := name == daemonapp.Name
	if !isInbox && !isRuns && !isRequests && !isAutonomy && !isDaemon {
		return writeError(stdout, stderr, api.NotFoundf("unknown app %q", name), false)
	}
	once := false
	for _, a := range args {
		if a == "--once" {
			once = true
		}
	}
	callWith := func(g globals, op string, record, dry bool, args []string) ([]byte, error) {
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
	call := func(op string, record, dry bool, args []string) ([]byte, error) {
		return callWith(g, op, record, dry, args)
	}
	inboxModel := func() *inboxapp.Model {
		im := inboxapp.NewModel(func(op string, args []string) ([]byte, error) {
			if op == "question.answer" && len(args) == 3 {
				// The app passes question, change set, words; the change set
				// is the write target, not a request field.
				gg := g
				gg.change, gg.changeSet = args[1], true
				return callWith(gg, op, false, false, []string{args[0], args[2]})
			}
			if op == "accept" && len(args) == 1 {
				// The app passes the change set to accept: the write target,
				// accepted under the default rule, owner.
				gg := g
				gg.change, gg.changeSet = args[0], true
				return callWith(gg, op, false, false, nil)
			}
			return call(op, false, false, args)
		})
		return im
	}
	runsModel := func() *runsapp.Model {
		return runsapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, nil)
	}
	autonomyModel := func() *autonomyapp.Model {
		return autonomyapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		})
	}
	var m tea.Model
	var result func() ([]byte, error)
	switch {
	case isAutonomy:
		am := autonomyModel()
		result = am.Result
		if once {
			break
		}
		sh := am.Shell()
		if _, err := sh.Run(); err != nil {
			return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
		}
		if _, err := sh.Result(); err != nil {
			return writeError(stdout, stderr, err, false)
		}
		return 0
	case isDaemon:
		if once {
			return writeError(stdout, stderr, api.Usagef("tui app daemon has no --once: it writes the daemon settings from its form"), false)
		}
		da, err := daemonapp.New(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, theme.Dark(), colorprofile.Detect(os.Stdout, os.Environ()))
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		out, err := modalshell.New(da).Run()
		if err != nil {
			return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
		}
		stdout.Write(out)
		return 0
	case isInbox:
		im := inboxModel()
		m, result = im, im.Result
	case isRuns:
		rm := runsModel()
		m, result = modalshell.New(rm), rm.Result
	case isRequests:
		if once {
			return writeError(stdout, stderr, api.Usagef("tui app request has no --once: it files a request from its form"), false)
		}
		out, err := modalshell.New(requestapp.New(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, colorprofile.Detect(os.Stdout, os.Environ()))).Run()
		if err != nil {
			return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
		}
		stdout.Write(out)
		return 0
	}
	if once {
		out, err := result()
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		stdout.Write(out)
		return 0
	}
	if err := modalshell.Run(m); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	return 0
}

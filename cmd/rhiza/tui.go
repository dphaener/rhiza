package main

// `rhiza tui app <name>` is the program each app's Rex action runs. It is
// the one command here that is a terminal program rather than an operation,
// so it is not in the api registry.

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/autonomyapp"
	"github.com/dphaener/rhiza/internal/tui/daemonapp"
	"github.com/dphaener/rhiza/internal/tui/inboxapp"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
	"github.com/dphaener/rhiza/internal/tui/projectctx"
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
	apps := newTUIModels(g, getenv)
	inboxModel, runsModel, autonomyModel := apps.inbox, apps.runs, apps.autonomy
	if !once {
		return runSwitching(name, apps, stdout, stderr)
	}
	var result func() ([]byte, error)
	switch {
	case isAutonomy:
		result = autonomyModel().Result
	case isDaemon:
		return writeError(stdout, stderr, api.Usagef("tui app daemon has no --once: it writes the daemon settings from its form"), false)
	case isInbox:
		result = inboxModel().Result
	case isRuns:
		result = runsModel().Result
	case isRequests:
		return writeError(stdout, stderr, api.Usagef("tui app request has no --once: it files a request from its form"), false)
	}
	if once {
		out, err := result()
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		stdout.Write(out)
		return 0
	}
	return 0
}

// tuiModels builds the apps' models over the api operations, called
// in-process with the global flags in g.
type tuiModels struct {
	call  func(op string, record, dry bool, args []string) ([]byte, error)
	inbox func() *inboxapp.Model
	// mergedInbox is the inbox over every registered project.
	mergedInbox func(projects []string, folder func(string) string) *inboxapp.Model
	runs        func() *runsapp.Model
	// mergedRuns is the runs app over every registered project.
	mergedRuns func(projects []string, folder func(string) string) *runsapp.Model
	autonomy   func() *autonomyapp.Model
}

func newTUIModels(g globals, getenv func(string) string) tuiModels {
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
		return inboxapp.NewModel(func(op string, args []string) ([]byte, error) {
			return inboxRun(g, call, callWith, op, args)
		})
	}
	var mergedInbox func(projects []string, folder func(string) string) *inboxapp.Model
	mergedInbox = func(projects []string, folder func(string) string) *inboxapp.Model {
		// Each call runs in its project's folder, and the app's own
		// directory comes back after it.
		in := func(project, op string, args []string) ([]byte, error) {
			if wd, err := os.Getwd(); err == nil {
				defer os.Chdir(wd)
			}
			if err := os.Chdir(folder(project)); err != nil {
				return nil, err
			}
			return inboxRun(g, call, callWith, op, args)
		}
		return inboxapp.NewMergedModel(in, projects)
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
	mergedRuns := func(projects []string, folder func(string) string) *runsapp.Model {
		in := func(project, op string, args []string) ([]byte, error) {
			if wd, err := os.Getwd(); err == nil {
				defer os.Chdir(wd)
			}
			if err := os.Chdir(folder(project)); err != nil {
				return nil, err
			}
			return call(op, false, false, args)
		}
		return runsapp.NewMergedModel(in, projects, nil)
	}
	return tuiModels{call: call, inbox: inboxModel, mergedInbox: mergedInbox, mergedRuns: mergedRuns, runs: runsModel, autonomy: autonomyModel}
}

// inboxRun runs one of the inbox app's operations. The app passes
// question.answer a question, change set and words, and accept a change
// set: the change set is the write target, not a request field.
func inboxRun(g globals, call func(op string, record, dry bool, args []string) ([]byte, error),
	callWith func(g globals, op string, record, dry bool, args []string) ([]byte, error),
	op string, args []string) ([]byte, error) {
	if op == "question.answer" && len(args) == 3 {
		gg := g
		gg.change, gg.changeSet = args[1], true
		return callWith(gg, op, false, false, []string{args[0], args[2]})
	}
	if op == "accept" && len(args) == 1 {
		gg := g
		gg.change, gg.changeSet = args[0], true
		return callWith(gg, op, false, false, nil)
	}
	return call(op, false, false, args)
}

// runSwitching runs an app in a terminal pane, hosted by appSwitcher.
func runSwitching(name string, apps tuiModels, stdout, stderr io.Writer) int {
	s, err := appSwitcher(name, apps)
	if err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if err := modalshell.Run(s); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	out, err := s.Result()
	if err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if name == requestapp.Name || name == daemonapp.Name {
		stdout.Write(out)
	}
	return 0
}

// appSwitcher is the switcher the named app runs in. Every app names its
// current project, the one the working directory is in, in a header, and
// one key (modalshell.SwitchKey) opens a picker of the registered
// projects: picking one moves the app there and builds it anew, so its
// view is the picked project's. It fails when the daemon app's form cannot
// be built, before any pane opens.
func appSwitcher(name string, apps tuiModels) (*modalshell.Switcher, error) {
	call := apps.call
	wd, _ := os.Getwd()
	names, folder := projectctx.Folders(wd)
	sw := modalshell.Switching{
		Projects: names,
		Folder:   folder,
		Header: func() string {
			if wd, err := os.Getwd(); err == nil {
				return projectctx.Header(wd)
			}
			return projectctx.NoProject
		},
	}
	var daemonErr error
	sw.Build = func() modalshell.Built {
		switch name {
		case autonomyapp.Name:
			sh := apps.autonomy().Shell()
			return modalshell.Built{Model: sh, Result: sh.Result}
		case daemonapp.Name:
			cur, _ := projectctx.Current(wdNow())
			da, err := daemonapp.NewIn(func(op string, args []string) ([]byte, error) {
				return call(op, false, false, args)
			}, theme.Dark(), colorprofile.Detect(os.Stdout, os.Environ()), cur)
			if err != nil {
				daemonErr = err
				return modalshell.Built{Model: modalshell.New(modalshell.Dialog(ui.NewDialog(theme.Dark(), colorprofile.Detect(os.Stdout, os.Environ()), "Daemon", []string{err.Error()})))}
			}
			sh := modalshell.New(da)
			return modalshell.Built{Model: sh, Result: sh.Result}
		case inboxapp.Name:
			im := apps.inbox()
			if names, folder := projectctx.Folders(wdNow()); len(names) > 0 && apps.mergedInbox != nil {
				im = apps.mergedInbox(names, folder)
			}
			return modalshell.Built{Model: im, Result: im.Result}
		case runsapp.Name:
			rm := apps.runs()
			if names, folder := projectctx.Folders(wdNow()); len(names) > 0 && apps.mergedRuns != nil {
				rm = apps.mergedRuns(names, folder)
			}
			sh := modalshell.New(rm)
			return modalshell.Built{Model: sh, Result: rm.Result}
		default:
			wd, _ := os.Getwd()
			current, _ := projectctx.Current(wd)
			names, folder := projectctx.Folders(wd)
			sh := modalshell.New(requestapp.NewIn(func(project, op string, args []string) ([]byte, error) {
				if project != "" && folder != nil {
					if err := os.Chdir(folder(project)); err != nil {
						return nil, err
					}
				}
				return call(op, false, false, args)
			}, colorprofile.Detect(os.Stdout, os.Environ()), current, names))
			return modalshell.Built{Model: sh, Result: sh.Result}
		}
	}
	if name == daemonapp.Name {
		// Fail before opening a pane when the daemon form cannot be built.
		sw.Build()
		if daemonErr != nil {
			return nil, daemonErr
		}
	}
	return modalshell.NewSwitcher(sw), nil
}

func wdNow() string { wd, _ := os.Getwd(); return wd }

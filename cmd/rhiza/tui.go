package main

// `rhiza tui app <name>` is the program each app's Rex action runs. It is
// the one command here that is a terminal program rather than an operation,
// so it is not in the api registry.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/changesapp"
	"github.com/dphaener/rhiza/internal/tui/healthapp"
	"github.com/dphaener/rhiza/internal/tui/inboxapp"
	"github.com/dphaener/rhiza/internal/tui/liveapp"
	"github.com/dphaener/rhiza/internal/tui/modelapp"
	"github.com/dphaener/rhiza/internal/tui/readapps"
	"github.com/dphaener/rhiza/internal/tui/requestsapp"
	"github.com/dphaener/rhiza/internal/tui/roadmapapp"
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
	if name == changesapp.Name {
		return tuiChanges(g, args, stdout, stderr, getenv)
	}
	read, isRead := readapps.Lookup(name)
	isInbox := name == inboxapp.Name
	isRoadmap := name == roadmapapp.Name
	isLive := name == liveapp.Name
	isHealth := name == healthapp.Name
	isRequests := name == requestsapp.Name
	write, isWrite := writeapps.Lookup(name)
	if !isRead && !isWrite && !isInbox && !isRoadmap && !isLive && !isHealth && !isRequests {
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
			return call(op, false, false, args)
		})
		if ctx, err := newCtx(g, os.Stdin, getenv); err == nil {
			im.Exec, im.Self, im.Root, im.RexSession = ctx.Exec, ctx.Self, ctx.Root, ctx.RexSession
		}
		return im
	}
	liveModel := func() *liveapp.Model {
		return liveapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, nil)
	}
	var m tea.Model
	var result func() ([]byte, error)
	switch {
	case isInbox:
		im := inboxModel()
		m, result = im, im.Result
	case isLive:
		lm := liveModel()
		m, result = lm, lm.Result
	case isHealth:
		hm := healthapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		}, func(app string) (tea.Model, error) {
			switch app {
			case inboxapp.Name:
				return inboxModel(), nil
			case liveapp.Name:
				return liveModel(), nil
			}
			return nil, api.NotFoundf("unknown app %q", app)
		})
		m, result = hm, hm.Result
	case isRequests:
		qm := requestsapp.NewModel(func(op string, args []string) ([]byte, error) {
			return call(op, false, false, args)
		})
		m, result = qm, qm.Result
	case isRoadmap:
		am := roadmapModel(g, callWith)
		if ctx, err := newCtx(g, os.Stdin, getenv); err == nil {
			am.SetStarter(roadmapapp.Starter{
				Exec: func(name string, args []string, dir string) ([]byte, error) {
					return ctx.Exec(name, args, dir, nil)
				},
				Self: ctx.Self, Root: ctx.Root, RexSession: ctx.RexSession,
			})
		}
		m, result = am, am.Result
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

// roadmapModel is the roadmap app over call, which runs an operation as
// the matching command does. Its what-ifs are the same operations with
// --dry-run.
func roadmapModel(g globals, call func(g globals, op string, record, dry bool, args []string) ([]byte, error)) *roadmapapp.Model {
	runner := func(dry bool) roadmapapp.Runner {
		return func(op, change string, args []string) ([]byte, error) {
			gg := g
			if change != "" {
				gg.change, gg.changeSet = change, true
			}
			return call(gg, op, false, dry, args)
		}
	}
	am := roadmapapp.NewModel(runner(false))
	am.SetPreview(runner(true))
	return am
}

// tuiChanges runs the changes app: the landed change sets not seen yet,
// newest first. Showing them marks them seen in .git/rhiza/feed.json under
// the git common dir. With --once it prints them as JSON and exits.
func tuiChanges(g globals, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	once := false
	for _, a := range args {
		if a != "--once" {
			return writeError(stdout, stderr, api.Usagef("tui app changes takes [--once]"), false)
		}
		once = true
	}
	ctx, err := newCtx(g, os.Stdin, getenv)
	if err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if ctx.Git == nil {
		return writeError(stdout, stderr, api.Errorf("changes needs git: it remembers what you saw in the git common dir"), false)
	}
	out, err := ctx.Git("rev-parse", "--git-common-dir")
	if err != nil {
		return writeError(stdout, stderr, api.Errorf("changes needs a git repository: %v", err), false)
	}
	common := strings.TrimSpace(out)
	if !filepath.IsAbs(common) {
		common = filepath.Join(ctx.Root, common)
	}
	o, _ := api.Lookup("feed")
	read := func(args []string) ([]byte, error) {
		req, err := requestFromArgs(o, args)
		if err != nil {
			return nil, err
		}
		resp, err := o.Run(ctx, req)
		if err != nil {
			return nil, err
		}
		return canon.Marshal(resp)
	}
	m := changesapp.NewModel(read, filepath.Join(common, "rhiza", changesapp.SeenFile))
	if err := m.Err(); err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if once {
		res, err := m.Result()
		if err == nil {
			err = m.MarkSeen()
		}
		if err != nil {
			return writeError(stdout, stderr, err, false)
		}
		stdout.Write(res)
		return 0
	}
	if err := m.MarkSeen(); err != nil {
		return writeError(stdout, stderr, err, false)
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui app: %w", err), false)
	}
	return 0
}

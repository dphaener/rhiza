package main

// `rhiza tui open <app>` is what each app's Rex action runs: it opens the
// app as a rexuikit Modal over Rex's control connection (internal/tui/modalopen),
// running `rhiza tui app <app>` in it, and resizes the modal as the app's
// content changes until the app exits.

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dphaener/rexuikit/client"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalopen"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

const tuiOpenUsage = "rhiza tui open <app> [--session <id>] [--block <id>]"

// tuiOpen opens the named app's modal. The session and block default to
// REX_SESSION and REX_BLOCK; the modal opens over the block's window and
// starts in this process's directory.
func tuiOpen(name string, args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	session, block := getenv("REX_SESSION"), getenv("REX_BLOCK")
	for i := 0; i < len(args); i++ {
		if i+1 >= len(args) || (args[i] != "--session" && args[i] != "--block") {
			fmt.Fprintln(stderr, "usage: "+tuiOpenUsage)
			return 2
		}
		if args[i] == "--session" {
			session = args[i+1]
		} else {
			block = args[i+1]
		}
		i++
	}
	var app launcher.App
	found := false
	for _, a := range launcher.Apps() {
		if a.Name == name {
			app, found = a, true
		}
	}
	if !found {
		return writeError(stdout, stderr, api.NotFoundf("unknown app %q", name), false)
	}
	if session == "" {
		return writeError(stdout, stderr, api.Usagef("tui open: no Rex session: pass --session or run it in Rex"), false)
	}
	if err := openModal(app, session, block, getenv); err != nil {
		return writeError(stdout, stderr, fmt.Errorf("tui open %s: %w", name, err), false)
	}
	return 0
}

func openModal(app launcher.App, session, block string, getenv func(string) string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	sf, err := os.CreateTemp("", "rhiza-modal-*.size")
	if err != nil {
		return err
	}
	sf.Close()
	defer os.Remove(sf.Name())

	addr, err := modalopen.ServerAddr(getenv)
	if err != nil {
		return err
	}
	dial, cancelDial := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDial()
	c, err := client.Dial(dial, client.Options{Addr: addr, Getenv: getenv, Name: "rhiza"})
	if err != nil {
		return err
	}
	defer c.Close()
	if _, err := c.Attach(dial, session); err != nil {
		return err
	}
	// The modal outlives the dial: it ends when the app exits, or with the
	// connection.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-c.Done()
		cancel()
	}()
	cols, rows := app.Size()
	return modalopen.Open(ctx, c, modalopen.Request{
		SessionID: session,
		BlockID:   block,
		Label:     app.Name,
		Argv:      []string{exe, "tui", "app", app.Name, modalshell.SizeFileFlag, sf.Name()},
		Cwd:       cwd,
		Cols:      cols,
		Rows:      rows,
		Sizes:     modalopen.WatchSize(ctx, sf.Name(), 100*time.Millisecond),
	})
}

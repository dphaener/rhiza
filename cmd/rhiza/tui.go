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
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

func tuiCommand(args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	if len(args) != 2 || args[1] != "launch" {
		fmt.Fprintln(stderr, "usage: rhiza tui launch")
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

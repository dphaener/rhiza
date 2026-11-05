// Package daemonapp is the Daemon app on the modal shell: a rexuikit Form
// over the daemon's settings. It shows whether the daemon is paused as a
// Confirm, prefilled from daemon.settings; submitting a change writes
// rhiza.daemon.paused through daemon.set and closes the app. Submitting
// unchanged, or Esc, writes nothing. A Text field below it holds the run
// limit, prefilled from rhiza.daemon.maxRuns; a valid new limit is written
// through daemon.set the same way, and a non-number or 0 is refused by the
// field's validator, so nothing is called.
package daemonapp

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Name is the app's name.
const Name = "daemon"

// PausedKey is the setting the pause Confirm writes.
const PausedKey = "rhiza.daemon.paused"

// MaxRunsKey is the setting the run-limit field writes.
const MaxRunsKey = "rhiza.daemon.maxRuns"

const (
	pausedField  = "paused"
	maxRunsField = "maxRuns"
)

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// App is the Daemon app's component, hosted by the modal shell.
type App struct {
	modalshell.Component
	run     Runner
	paused  bool
	maxRuns string
}

// validateMaxRuns accepts a whole number of at least 1.
func validateMaxRuns(v string) error {
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fmt.Errorf("must be a whole number of at least 1")
	}
	return nil
}

// New reads daemon.settings through run and returns the app with the
// Confirm prefilled from rhiza.daemon.paused.
func New(run Runner, t theme.Theme, profile colorprofile.Profile) (*App, error) {
	out, err := run("daemon.settings", nil)
	if err != nil {
		return nil, err
	}
	var s api.DaemonSettingsResponse
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, err
	}
	paused := false
	maxRuns := ""
	for _, row := range s.Settings {
		switch row.Key {
		case PausedKey:
			paused = row.Value == "true"
		case MaxRunsKey:
			maxRuns = row.Value
		}
	}
	confirm := ui.NewConfirm(pausedField, "Pause the daemon?").Default(paused)
	limit := ui.NewText(maxRunsField, "Run limit").Default(maxRuns).Validator(validateMaxRuns)
	form := ui.NewForm(t, profile, confirm, limit).
		WithSubtitle("Daemon").
		WithThemeFunc(func(color.Color) theme.Theme { return t })
	return &App{Component: modalshell.Form(form), run: run, paused: paused, maxRuns: maxRuns}, nil
}

// Paused is the value the app opened with.
func (a *App) Paused() bool { return a.paused }

// MaxRuns is the run limit the app opened with.
func (a *App) MaxRuns() string { return a.maxRuns }

// Update writes a submitted change and reports the app done; every other
// message goes to the form.
func (a *App) Update(msg tea.Msg) tea.Cmd {
	sub, ok := msg.(ui.FormSubmittedMsg)
	if !ok {
		return a.Component.Update(msg)
	}
	var out []byte
	if v := sub.Values.Bool(pausedField); v != a.paused {
		var err error
		if out, err = a.run("daemon.set", []string{"--", PausedKey, fmt.Sprint(v)}); err != nil {
			return modalshell.Done(out, err)
		}
	}
	if v := sub.Values.String(maxRunsField); v != a.maxRuns {
		var err error
		if out, err = a.run("daemon.set", []string{"--", MaxRunsKey, v}); err != nil {
			return modalshell.Done(out, err)
		}
	}
	return modalshell.Done(out, nil)
}

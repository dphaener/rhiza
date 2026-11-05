// Package daemonapp is the Daemon app on the modal shell: a rexuikit Form
// over the daemon's settings. It shows whether the daemon is paused as a
// Confirm, prefilled from daemon.settings; submitting a change writes
// rhiza.daemon.paused through daemon.set and closes the app. Submitting
// unchanged, or Esc, writes nothing.
package daemonapp

import (
	"encoding/json"
	"fmt"
	"image/color"

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

const pausedField = "paused"

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// App is the Daemon app's component, hosted by the modal shell.
type App struct {
	modalshell.Component
	run    Runner
	paused bool
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
	for _, row := range s.Settings {
		if row.Key == PausedKey {
			paused = row.Value == "true"
		}
	}
	confirm := ui.NewConfirm(pausedField, "Pause the daemon?").Default(paused)
	form := ui.NewForm(t, profile, confirm).
		WithSubtitle("Daemon").
		WithThemeFunc(func(color.Color) theme.Theme { return t })
	return &App{Component: modalshell.Form(form), run: run, paused: paused}, nil
}

// Paused is the value the app opened with.
func (a *App) Paused() bool { return a.paused }

// Update writes a submitted change and reports the app done; every other
// message goes to the form.
func (a *App) Update(msg tea.Msg) tea.Cmd {
	sub, ok := msg.(ui.FormSubmittedMsg)
	if !ok {
		return a.Component.Update(msg)
	}
	v := sub.Values.Bool(pausedField)
	if v == a.paused {
		return modalshell.Done(nil, nil)
	}
	out, err := a.run("daemon.set", []string{"--", PausedKey, fmt.Sprint(v)})
	return modalshell.Done(out, err)
}

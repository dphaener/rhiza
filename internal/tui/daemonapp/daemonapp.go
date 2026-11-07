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
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Name is the app's name.
const Name = "daemon"

func init() {
	launcher.Register(launcher.App{
		Name:    Name,
		Key:     "d",
		Summary: "Pause the daemon and set its run limit",
		Component: func() modalshell.Component {
			a, err := New(func(string, []string) ([]byte, error) { return []byte("{}"), nil },
				theme.Dark(), colorprofile.TrueColor)
			if err != nil {
				return nil
			}
			return a
		},
	})
}

// PausedKey is the setting the pause Confirm writes.
const PausedKey = "rhiza.daemon.paused"

// MaxRunsKey is the setting the run-limit field writes.
const MaxRunsKey = "rhiza.daemon.maxRuns"

// ProjectPausedKey and ProjectMaxRunsKey are the current project's own
// pause and run cap, written to its git config.
const (
	ProjectPausedKey  = "rhiza.project.paused"
	ProjectMaxRunsKey = "rhiza.project.maxRuns"
)

const (
	pausedField  = "paused"
	maxRunsField = "maxRuns"
	projPaused   = "projectPaused"
	projMaxRuns  = "projectMaxRuns"
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
	// project is the current project's name, empty when there is none; its
	// own pause and run cap are then fields too.
	project  string
	pPaused  bool
	pMaxRuns string
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
	return NewIn(run, t, profile, "")
}

// NewIn is New for the app opened in project: the machine settings show
// once, and the project's own pause and run cap follow, written to its
// git config. An empty project shows the machine settings only.
func NewIn(run Runner, t theme.Theme, profile colorprofile.Profile, project string) (*App, error) {
	out, err := run("daemon.settings", nil)
	if err != nil {
		return nil, err
	}
	var s api.DaemonSettingsResponse
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, err
	}
	paused := false
	maxRuns, pMax, pPause := "", "", false
	for _, row := range s.Settings {
		switch row.Key {
		case PausedKey:
			paused = row.Value == "true"
		case MaxRunsKey:
			maxRuns = row.Value
		case ProjectPausedKey:
			pPause = row.Value == "true"
		case ProjectMaxRunsKey:
			pMax = row.Value
		}
	}
	confirm := ui.NewConfirm(pausedField, "Pause the daemon?").Default(paused)
	limit := ui.NewText(maxRunsField, "Run limit").Default(maxRuns).Validator(validateMaxRuns)
	fields := []ui.Field{confirm, limit}
	if project != "" {
		fields = append(fields,
			ui.NewConfirm(projPaused, "Pause project "+project+"?").Default(pPause),
			ui.NewText(projMaxRuns, "Project run limit").Default(pMax).Validator(func(v string) error {
				if v == "" {
					return nil
				}
				return validateMaxRuns(v)
			}))
	}
	form := ui.NewForm(t, profile, fields...).
		WithSubtitle("Daemon").
		WithThemeFunc(func(color.Color) theme.Theme { return t })
	return &App{Component: modalshell.Form(form), run: run, paused: paused, maxRuns: maxRuns,
		project: project, pPaused: pPause, pMaxRuns: pMax}, nil
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
	if a.project != "" {
		if v := sub.Values.Bool(projPaused); v != a.pPaused {
			var err error
			if out, err = a.run("daemon.set", []string{"--", ProjectPausedKey, fmt.Sprint(v)}); err != nil {
				return modalshell.Done(out, err)
			}
		}
		if v := sub.Values.String(projMaxRuns); v != a.pMaxRuns && v != "" {
			var err error
			if out, err = a.run("daemon.set", []string{"--", ProjectMaxRunsKey, v}); err != nil {
				return modalshell.Done(out, err)
			}
		}
	}
	return modalshell.Done(out, nil)
}

// Package requestapp is the request app: one rexuikit Form, a Select of the
// request's kind (feature, bug or milestone) and a TextArea for its text,
// hosted by the modal shell. Submitting queues the request through the
// request operation, in-process through the registry, and closes the app;
// Esc closes it with nothing queued.
package requestapp

import (
	"errors"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Name is the app's name.
const Name = "request"

// MaxText bounds a request's typed text.
const MaxText = 500

// Field keys of the form.
const (
	kindKey = "kind"
	textKey = "text"
)

func init() {
	// Key is fixed so that adding the app moves no other app's default key.
	launcher.Register(launcher.App{
		Name:    Name,
		Key:     "q",
		Summary: "File a feature, bug or milestone request",
		Component: func() modalshell.Component {
			return New(nil, colorprofile.Detect(os.Stdout, os.Environ()))
		},
	})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// RequestArgs is the argv of a request op for text of kind. The text goes
// after "--" so text starting with a dash is not read as a flag.
func RequestArgs(kind, text string) []string {
	return []string{"--kind", kind, "--", text}
}

// App is the request app as a modal shell component.
type App struct {
	modalshell.Component
	run Runner
}

// New returns the request app over run, drawn for profile.
func New(run Runner, profile colorprofile.Profile) *App {
	kind := ui.NewSelect(kindKey, "Kind").Options(
		ui.Option{Label: "Feature", Value: "feature"},
		ui.Option{Label: "Bug", Value: "bug"},
		ui.Option{Label: "Milestone", Value: "milestone"},
	).Default("feature")
	text := ui.NewTextArea(textKey, "Request").
		Placeholder("What to plan").
		CharLimit(MaxText).
		Validator(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return errors.New("A request needs its text")
			}
			return nil
		})
	f := ui.NewForm(theme.Dark(), profile, kind, text)
	return &App{Component: modalshell.Form(f), run: run}
}

// Update queues the request when the form submits and reports the app done;
// every other message goes to the form.
func (a *App) Update(msg tea.Msg) tea.Cmd {
	if s, ok := msg.(ui.FormSubmittedMsg); ok {
		text := strings.Join(strings.Fields(s.Values.String(textKey)), " ")
		if text == "" || a.run == nil {
			return modalshell.Done(nil, errors.New("request: nothing to queue"))
		}
		return modalshell.Done(a.run("request", RequestArgs(s.Values.String(kindKey), text)))
	}
	return a.Component.Update(msg)
}

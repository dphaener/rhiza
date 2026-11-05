// Package autonomyapp is the autonomy app: one rexuikit Form of the plan,
// accept, ship and sign-off toggles (rhiza.autonomy.*), each auto or owner,
// prefilled from daemon.settings and hosted by the modal shell. Submitting
// writes each changed toggle through daemon.set and closes the app; Esc
// writes nothing. --once prints what `rhiza daemon settings --json` prints.
package autonomyapp

import (
	"encoding/json"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Name is the app's name.
const Name = "autonomy"

// Key is the app's key in the rhiza key mode.
const Key = "y"

// Toggle is one autonomy toggle: its setting key and its label.
type Toggle struct {
	Key, Label string
}

// Toggles are the form's fields, in order.
var Toggles = []Toggle{
	{"rhiza.autonomy.plan", "Plan"},
	{"rhiza.autonomy.accept", "Accept"},
	{"rhiza.autonomy.ship", "Ship"},
	{"rhiza.autonomy.signoff", "Sign-off"},
}

func init() {
	launcher.Register(launcher.App{
		Name:    Name,
		Key:     Key,
		Summary: "Set the plan, accept, ship and sign-off toggles to auto or owner",
		Component: func() modalshell.Component {
			return newComponent(nil, defaultValues(), nil)
		},
	})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// Model is the app: a modalshell.Component over a ui.Form.
type Model struct {
	run     Runner
	form    *ui.Form
	initial map[string]string
	out     []byte
	err     error
}

func defaultValues() map[string]string {
	v := map[string]string{}
	for _, t := range Toggles {
		v[t.Key] = api.AutonomyOwner
	}
	return v
}

func newComponent(run Runner, vals map[string]string, out []byte) *Model {
	fields := make([]ui.Field, 0, len(Toggles))
	for _, t := range Toggles {
		fields = append(fields, ui.NewSelect(t.Key, t.Label).
			Options(
				ui.Option{Label: api.AutonomyAuto, Value: api.AutonomyAuto},
				ui.Option{Label: api.AutonomyOwner, Value: api.AutonomyOwner},
			).Rows(2).Default(vals[t.Key]))
	}
	form := ui.NewForm(theme.Dark(), colorprofile.TrueColor, fields...).
		WithSubtitle("Who decides each step: auto or the owner").
		WithHint("↑/↓ choose · enter next · enter on the last saves · esc cancels")
	return &Model{run: run, form: form, initial: vals, out: out}
}

// NewModel returns the app over run, loaded once from daemon.settings. A
// failure to load is kept for Result and the form shows the defaults.
func NewModel(run Runner) *Model {
	out, err := run("daemon.settings", nil)
	vals := defaultValues()
	if err == nil {
		var s api.DaemonSettingsResponse
		if err = json.Unmarshal(out, &s); err == nil {
			for _, st := range s.Settings {
				if _, ok := vals[st.Key]; ok {
					vals[st.Key] = st.Value
				}
			}
		}
	}
	m := newComponent(run, vals, out)
	m.err = err
	return m
}

// Result is the daemon.settings output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Values are the toggles as the form holds them now.
func (m *Model) Values() map[string]string {
	v := map[string]string{}
	for _, t := range Toggles {
		v[t.Key] = m.form.Values().String(t.Key)
	}
	return v
}

// Shell hosts the app in the modal shell.
func (m *Model) Shell() *modalshell.Shell { return modalshell.New(m) }

func (m *Model) Init() tea.Cmd { return m.form.Init() }

// Update handles the form's submission: each changed toggle is written
// through daemon.set, once, and the app reports done.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if sub, ok := msg.(ui.FormSubmittedMsg); ok {
		return modalshell.Done(m.submit(sub.Values))
	}
	_, cmd := m.form.Update(msg)
	return cmd
}

func (m *Model) submit(v ui.Values) ([]byte, error) {
	for _, t := range Toggles {
		want := v.String(t.Key)
		if want == "" || want == m.initial[t.Key] {
			continue
		}
		if _, err := m.run("daemon.set", []string{"--", t.Key, want}); err != nil {
			return nil, fmt.Errorf("set %s: %w", t.Key, err)
		}
	}
	return m.out, nil
}

func (m *Model) Render() string                { return m.form.Render() }
func (m *Model) NaturalSize() (cols, rows int) { return m.form.NaturalSize() }
func (m *Model) Theme() theme.Theme            { return m.form.Theme() }

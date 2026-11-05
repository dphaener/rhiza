// Package healthapp is the health app: rhiza health's report, part by part,
// each marked ok, warn or bad. It runs the api operation health in-process
// through the registry, so --once prints what `rhiza health --json` prints.
// r reloads; enter opens the app that shows more of the part under the
// cursor (the inbox app for inbox, the runs app for daemon, the autonomy
// app for autonomy) in place of the report, until that app quits.
package healthapp

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name.
const Name = "health"

// Key is the app's key in the rhiza key mode.
const Key = "t"

func init() {
	launcher.Register(launcher.App{Name: Name, Key: Key, Summary: "The whole system's state, part by part: ok, warn or bad"})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// Opener builds the named app, ready to run in place of this one.
type Opener func(app string) (tea.Model, error)

// apps maps a part to the app that shows more of it.
var apps = map[string]string{
	api.HealthInbox:    "inbox",
	api.HealthDaemon:   "runs",
	api.HealthAutonomy: "autonomy",
}

// AppFor is the app enter opens on a part, and false when none.
func AppFor(part string) (string, bool) {
	a, ok := apps[part]
	return a, ok
}

// Model is the app's bubbletea program. q or esc quits.
type Model struct {
	run    Runner
	open   Opener
	out    []byte
	report api.HealthResponse
	err    error
	cursor int
	notice string
	// child is the app enter opened, shown in place of the report.
	child tea.Model
}

// NewModel returns the model over run, loaded once. open builds the app
// enter opens; nil opens none.
func NewModel(run Runner, open Opener) *Model {
	m := &Model{run: run, open: open}
	m.Reload()
	return m
}

// Reload runs health again. The cursor stays on the same part.
func (m *Model) Reload() {
	out, err := m.run("health", nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var r api.HealthResponse
	if err := json.Unmarshal(out, &r); err != nil {
		m.err = err
		return
	}
	m.report = r
	m.cursor = max(0, min(m.cursor, len(r.Parts)-1))
}

// Result is the health output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Selected is the part under the cursor.
func (m *Model) Selected() (api.HealthPart, bool) {
	if m.cursor < 0 || m.cursor >= len(m.report.Parts) {
		return api.HealthPart{}, false
	}
	return m.report.Parts[m.cursor], true
}

// Child is the app enter opened, nil while the report shows.
func (m *Model) Child() tea.Model { return m.child }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.child != nil {
		c, cmd := m.child.Update(msg)
		m.child = c
		return m, cmd
	}
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.cursor = max(0, min(m.cursor+1, len(m.report.Parts)-1))
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "r":
			m.notice = ""
			m.Reload()
		case "enter":
			return m, m.enter()
		}
	}
	return m, nil
}

// enter opens the app of the part under the cursor.
func (m *Model) enter() tea.Cmd {
	p, ok := m.Selected()
	if !ok {
		return nil
	}
	app, ok := AppFor(p.Part)
	if !ok {
		m.notice = "No app shows more of " + p.Part
		return nil
	}
	if m.open == nil {
		m.notice = "Cannot open the " + app + " app here"
		return nil
	}
	c, err := m.open(app)
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return nil
	}
	m.child, m.notice = c, ""
	return c.Init()
}

// marks are how each state shows.
var marks = map[string]string{
	api.HealthOK:   "[ok]  ",
	api.HealthWarn: "[warn]",
	api.HealthBad:  "[bad] ",
}

// Mark is how a state shows in the view.
func Mark(state string) string {
	if s, ok := marks[state]; ok {
		return s
	}
	return "[" + state + "]"
}

func (m *Model) View() tea.View {
	if m.child != nil {
		return m.child.View()
	}
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	} else {
		fmt.Fprintf(&b, "Health: %s\n", m.report.State)
	}
	for i, p := range m.report.Parts {
		cur := "  "
		if i == m.cursor {
			cur = "> "
		}
		fmt.Fprintf(&b, "%s%s %-8s %s\n", cur, Mark(p.State), p.Part, p.Summary)
		if i != m.cursor {
			continue
		}
		for _, d := range p.Detail {
			fmt.Fprintf(&b, "           %s\n", d)
		}
		if p.Next != "" {
			fmt.Fprintf(&b, "           next: %s\n", p.Next)
		}
	}
	if m.notice != "" {
		b.WriteString(m.notice + "\n")
	}
	b.WriteString(Name + " · j/k move · enter open app · r reload · q quit")
	return tea.NewView(b.String())
}

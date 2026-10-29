// Package readapps is the read-only rhiza apps: roadmap, ready, show, why
// and slice. Each wraps the api operation of the same name, run in-process
// through the registry, so an app shows what `rhiza <op> --json` prints. An
// app starts no daemon and binds no socket.
package readapps

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// App is one read app. Name is both the app's name and the api operation it
// runs.
type App struct {
	Name    string
	Summary string
	// Prompt, when set, is what the app asks for when started without
	// arguments: its operation needs a ref.
	Prompt string
}

// Apps are the read apps, in the order they are registered.
var Apps = []App{
	{Name: "roadmap", Summary: "Milestones and the change sets that serve them"},
	{Name: "ready", Summary: "Which change sets can be worked on now"},
	{Name: "show", Summary: "One node or link", Prompt: "show (ID or handle)"},
	{Name: "why", Summary: "Why a node exists", Prompt: "why (ID or handle)"},
	{Name: "slice", Summary: "A budgeted slice of the model"},
}

func init() {
	for _, a := range Apps {
		launcher.Register(launcher.App{Name: a.Name, Summary: a.Summary})
	}
}

// Lookup returns the read app named name.
func Lookup(name string) (App, bool) {
	for _, a := range Apps {
		if a.Name == name {
			return a, true
		}
	}
	return App{}, false
}

// Runner runs an app's operation on args and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// Model is the app's bubbletea program: a prompt for a ref when the app
// needs one, then the result in a scrolling pane. q or esc quits.
type Model struct {
	app    App
	run    Runner
	input  string
	asking bool
	lines  []string
	err    error
	top    int
	rows   int
}

// NewModel returns the model of a over run, started with args. An app that
// has a prompt and no args asks for its ref first.
func NewModel(a App, run Runner, args []string) *Model {
	m := &Model{app: a, run: run, rows: 24}
	if a.Prompt != "" && len(args) == 0 {
		m.asking = true
		return m
	}
	m.load(args)
	return m
}

func (m *Model) load(args []string) {
	m.asking = false
	out, err := m.run(m.app.Name, args)
	m.err = err
	m.top = 0
	m.lines = strings.Split(strings.TrimRight(string(out), "\n"), "\n")
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.rows = max(msg.Height-2, 1)
	case tea.KeyPressMsg:
		if m.asking {
			return m.askKey(msg)
		}
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.scroll(1)
		case "up", "k":
			m.scroll(-1)
		case "pgdown", "space":
			m.scroll(m.rows)
		case "pgup":
			m.scroll(-m.rows)
		case "home", "g":
			m.top = 0
		case "end", "G":
			m.scroll(len(m.lines))
		}
	}
	return m, nil
}

func (m *Model) askKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		return m, tea.Quit
	case "enter":
		if m.input != "" {
			m.load([]string{m.input})
		}
	case "backspace":
		if r := []rune(m.input); len(r) > 0 {
			m.input = string(r[:len(r)-1])
		}
	default:
		if msg.Text != "" {
			m.input += msg.Text
		}
	}
	return m, nil
}

func (m *Model) scroll(d int) {
	m.top = max(0, min(m.top+d, len(m.lines)-m.rows))
}

func (m *Model) View() tea.View {
	if m.asking {
		return tea.NewView(fmt.Sprintf("%s: %s", m.app.Prompt, m.input))
	}
	end := min(m.top+m.rows, len(m.lines))
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	b.WriteString(strings.Join(m.lines[m.top:end], "\n"))
	b.WriteString("\n\n" + m.app.Name + " · j/k scroll · q quit")
	return tea.NewView(b.String())
}

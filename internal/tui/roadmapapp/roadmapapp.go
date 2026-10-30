// Package roadmapapp is the roadmap app: the roadmap as a cursor list, a
// heading per milestone group and the backlog and a row per change set with
// its short ID, state, title and the reasons it is not ready. It reads the
// roadmap and ready operations in-process through the registry. An app
// starts no daemon and binds no socket.
package roadmapapp

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name.
const Name = "roadmap"

func init() {
	launcher.Register(launcher.App{Name: Name, Summary: "Milestones and the change sets that serve them"})
}

// Runner runs an api operation on args and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

const shortIDLen = 14

// row is one line of the list: a heading, or a change set when ID is set.
type row struct {
	heading string
	id      id.ID
	state   string
	title   string
	reasons []string
}

// Model is the app's bubbletea program. j/k or the arrows move the cursor
// over the change set rows, r reloads, q or esc quits.
type Model struct {
	run    Runner
	out    []byte
	rows   []row
	err    error
	cursor int
	top    int
	height int
}

// NewModel returns the model over run, loaded.
func NewModel(run Runner) *Model {
	m := &Model{run: run, height: 24}
	m.load()
	return m
}

// load reads the roadmap and the ready reasons and rebuilds the rows,
// keeping the cursor on the same change set where it is still there.
func (m *Model) load() {
	var keep id.ID
	if m.cursor < len(m.rows) {
		keep = m.rows[m.cursor].id
	}
	out, err := m.run("roadmap", nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var rm api.RoadmapResponse
	if err := json.Unmarshal(out, &rm); err != nil {
		m.err = fmt.Errorf("roadmap: %w", err)
		return
	}
	reasons := map[id.ID][]string{}
	rout, err := m.run("ready", nil)
	if err != nil {
		m.err = err
		return
	}
	var rd api.ReadyResponse
	if err := json.Unmarshal(rout, &rd); err != nil {
		m.err = fmt.Errorf("ready: %w", err)
		return
	}
	for _, c := range rd.ChangeSets {
		reasons[c.ID] = c.Reasons
	}
	m.rows = m.rows[:0]
	for _, g := range rm.Groups {
		name := "Backlog"
		if g.Kind != "backlog" {
			name = strings.Join(append(slices.Clone(g.Intent.Path), cmp.Or(g.Intent.Handle, short(g.Intent.ID))), " / ")
		}
		m.rows = append(m.rows, row{heading: name})
		for _, c := range g.ChangeSets {
			m.rows = append(m.rows, row{id: c.ID, state: string(c.State), title: c.Intent, reasons: reasons[c.ID]})
		}
	}
	m.cursor = m.first()
	for i, r := range m.rows {
		if !keep.IsZero() && r.id == keep {
			m.cursor = i
		}
	}
	m.top = 0
	m.follow()
}

func short(i id.ID) string {
	s := i.String()
	return s[:min(len(s), shortIDLen)]
}

// first is the index of the first change set row, or 0 if there is none.
func (m *Model) first() int {
	for i, r := range m.rows {
		if !r.id.IsZero() {
			return i
		}
	}
	return 0
}

// Result is the roadmap operation's output, byte for byte, or its error.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Cursor is the change set under the cursor; zero when there is none.
func (m *Model) Cursor() id.ID {
	if m.cursor < len(m.rows) {
		return m.rows[m.cursor].id
	}
	return id.ID{}
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = max(msg.Height-2, 1)
		m.follow()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.move(1)
		case "up", "k":
			m.move(-1)
		case "r":
			m.load()
		}
	}
	return m, nil
}

// move puts the cursor on the next change set row in direction d, skipping
// headings; at either end it stays.
func (m *Model) move(d int) {
	for i := m.cursor + d; i >= 0 && i < len(m.rows); i += d {
		if !m.rows[i].id.IsZero() {
			m.cursor = i
			break
		}
	}
	m.follow()
}

// follow scrolls so the cursor row is on screen, with its heading when it
// is the first row of a group.
func (m *Model) follow() {
	if m.cursor < m.top {
		m.top = m.cursor
		if m.top > 0 && m.rows[m.top-1].id.IsZero() {
			m.top--
		}
	}
	if m.cursor >= m.top+m.height {
		m.top = m.cursor - m.height + 1
	}
}

func (m *Model) View() tea.View {
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	if len(m.rows) == 0 && m.err == nil {
		b.WriteString("Nothing planned yet.\n")
	}
	end := min(m.top+m.height, len(m.rows))
	for i := m.top; i < end; i++ {
		r := m.rows[i]
		if r.id.IsZero() {
			fmt.Fprintf(&b, "\n%s\n", r.heading)
			continue
		}
		mark := "  "
		if i == m.cursor {
			mark = "> "
		}
		fmt.Fprintf(&b, "%s%s [%s] %s", mark, short(r.id), r.state, r.title)
		if len(r.reasons) > 0 {
			fmt.Fprintf(&b, " (not ready: %s)", strings.Join(r.reasons, "; "))
		}
		b.WriteString("\n")
	}
	b.WriteString("\nroadmap · j/k move · r reload · q quit")
	return tea.NewView(b.String())
}

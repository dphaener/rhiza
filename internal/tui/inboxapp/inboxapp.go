// Package inboxapp is the inbox app: the items of `rhiza inbox` as a cursor
// list grouped by kind. It runs the api operation in-process through the
// registry, so --once prints what `rhiza inbox --json` prints. The list
// reloads every few seconds and on r.
package inboxapp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name, and the api operation it runs.
const Name = "inbox"

// Interval is how often the list reloads.
const Interval = 3 * time.Second

func init() {
	launcher.Register(launcher.App{Name: Name, Summary: "What waits on you: questions, sign-offs and failed runs"})
}

// Runner runs the inbox operation and returns its result as
// `rhiza inbox --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// TickMsg is the reload timer firing.
type TickMsg struct{}

// groups are the kinds in the order the list shows them, with their headings.
var groups = []struct{ kind, heading string }{
	{api.InboxQuestion, "Questions"},
	{api.InboxSignoff, "Sign-offs"},
	{api.InboxFailed, "Failed runs"},
	{api.InboxGaveUp, "Gave up"},
	{api.InboxSkipped, "Skipped"},
}

// Model is the app's bubbletea program. q or esc quits.
type Model struct {
	run    Runner
	out    []byte
	items  []api.InboxItem // in display order, grouped by kind
	err    error
	cursor int
}

// NewModel returns the model over run, loaded once.
func NewModel(run Runner) *Model {
	m := &Model{run: run}
	m.Reload()
	return m
}

// Reload runs the operation again. The cursor stays on the same item when it
// still waits, else on the nearest row.
func (m *Model) Reload() {
	var kept string
	if it, ok := m.Selected(); ok {
		kept = it.ID
	}
	out, err := m.run(Name, nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var resp api.InboxResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		m.err = err
		return
	}
	m.items = m.items[:0]
	for _, g := range groups {
		for _, it := range resp.Items {
			if it.Kind == g.kind {
				m.items = append(m.items, it)
			}
		}
	}
	for i, it := range m.items {
		if it.ID == kept {
			m.cursor = i
		}
	}
	m.cursor = max(0, min(m.cursor, len(m.items)-1))
}

// Result is what the app shows: the operation's output, byte for byte, or
// the error it ran into.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Selected is the item under the cursor.
func (m *Model) Selected() (api.InboxItem, bool) {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return api.InboxItem{}, false
	}
	return m.items[m.cursor], true
}

func tick() tea.Cmd {
	return tea.Tick(Interval, func(time.Time) tea.Msg { return TickMsg{} })
}

func (m *Model) Init() tea.Cmd { return tick() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.Reload()
		return m, tick()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.cursor = max(0, min(m.cursor+1, len(m.items)-1))
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "r":
			m.Reload()
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	if len(m.items) == 0 && m.err == nil {
		b.WriteString("Nothing waits on you\n")
	}
	i := 0
	for _, g := range groups {
		first := true
		for _, it := range m.items {
			if it.Kind != g.kind {
				continue
			}
			if first {
				fmt.Fprintf(&b, "\n%s\n", g.heading)
				first = false
			}
			mark := "  "
			if i == m.cursor {
				mark = "> "
			}
			b.WriteString(mark + row(it))
			i++
		}
	}
	b.WriteString("\n" + Name + " · j/k move · r reload · q quit")
	return tea.NewView(b.String())
}

// row is one item: its change set's title and step, then what it asks or
// why it failed, and a question's recommendation.
func row(it api.InboxItem) string {
	title := it.ChangeTitle
	if title == "" {
		title = it.Change
	}
	var b strings.Builder
	b.WriteString(title)
	if it.Step != "" {
		b.WriteString(" · " + it.Step)
	}
	b.WriteString("\n")
	text := it.Detail
	if it.Kind == api.InboxQuestion {
		text = it.Question
	}
	if text != "" {
		b.WriteString("    " + text + "\n")
	}
	if it.Recommendation != "" {
		b.WriteString("    recommends: " + it.Recommendation + "\n")
	}
	return b.String()
}

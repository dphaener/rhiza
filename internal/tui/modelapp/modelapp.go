// Package modelapp is the model app: the as-is capability tree and the
// roadmap's groups in one scrolling pane, switched with tab. Both are the
// text `rhiza tree --view as-is` and `rhiza roadmap` print, read in-process
// from the index through the api operations of those names. The app starts
// no daemon and binds no socket.
package modelapp

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name, as `rhiza tui app` takes it.
const Name = "model"

func init() {
	launcher.Register(launcher.App{Name: Name, Summary: "The capability tree and the roadmap"})
}

// Reader runs the api operation op on args and returns its markdown, as
// `rhiza <op> <args>` prints it.
type Reader func(op string, args []string) (string, error)

// pane is one of the app's two views.
type pane struct {
	title string
	op    string
	args  []string
	lines []string
	err   error
	top   int
}

// Model is the app's bubbletea program. tab or left and right switch
// between the tree and the roadmap; q or esc quits.
type Model struct {
	panes []*pane
	cur   int
	rows  int
}

// NewModel reads both panes with read and returns the model, showing the
// tree. A non-empty ref limits the tree to the subtree under that node, as
// `rhiza tree <ref>` does; the whole capability tree of a large model costs
// more than a subtree to read.
func NewModel(read Reader, ref string) *Model {
	treeArgs := []string{"--view", "as-is"}
	if ref != "" {
		treeArgs = append([]string{ref}, treeArgs...)
	}
	m := &Model{rows: 24, panes: []*pane{
		{title: "tree", op: "tree", args: treeArgs},
		{title: "roadmap", op: "roadmap"},
	}}
	for _, p := range m.panes {
		out, err := read(p.op, p.args)
		p.err = err
		p.lines = strings.Split(strings.TrimRight(out, "\n"), "\n")
	}
	return m
}

// Text is the markdown of the pane called title (tree or roadmap), as the
// command prints it; empty if the pane failed to read.
func (m *Model) Text(title string) string {
	for _, p := range m.panes {
		if p.title == title && p.err == nil {
			return strings.Join(p.lines, "\n") + "\n"
		}
	}
	return ""
}

// Err is the first error a pane's read ran into.
func (m *Model) Err() error {
	for _, p := range m.panes {
		if p.err != nil {
			return fmt.Errorf("%s: %w", p.title, p.err)
		}
	}
	return nil
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.rows = max(msg.Height-3, 1)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.cur = (m.cur + 1) % len(m.panes)
		case "shift+tab", "left", "h":
			m.cur = (m.cur + len(m.panes) - 1) % len(m.panes)
		case "down", "j":
			m.scroll(1)
		case "up", "k":
			m.scroll(-1)
		case "pgdown", "space":
			m.scroll(m.rows)
		case "pgup":
			m.scroll(-m.rows)
		case "home", "g":
			m.scroll(-len(m.panes[m.cur].lines))
		case "end", "G":
			m.scroll(len(m.panes[m.cur].lines))
		}
	}
	return m, nil
}

func (m *Model) scroll(d int) {
	p := m.panes[m.cur]
	p.top = max(0, min(p.top+d, len(p.lines)-m.rows))
}

func (m *Model) View() tea.View {
	var b strings.Builder
	for i, p := range m.panes {
		if i == m.cur {
			b.WriteString("[" + p.title + "] ")
		} else {
			b.WriteString(" " + p.title + "  ")
		}
	}
	b.WriteString("\n\n")
	p := m.panes[m.cur]
	if p.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", p.err)
	} else {
		end := min(p.top+m.rows, len(p.lines))
		b.WriteString(strings.Join(p.lines[p.top:end], "\n"))
	}
	b.WriteString("\n\ntab switch · j/k scroll · q quit")
	return tea.NewView(b.String())
}

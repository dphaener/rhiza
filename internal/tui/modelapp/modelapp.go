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

	read Reader
	// filter is the search filter: nil when the panes show as usual.
	filter *filter
}

// filter is the / filter. Its pane holds what the search operation found
// for text, in place of the current pane's lines.
type filter struct {
	text    string
	editing bool
	pane    pane
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
	m := &Model{rows: 24, read: read, panes: []*pane{
		{title: "tree", op: "tree", args: treeArgs},
		{title: "roadmap", op: "roadmap", args: []string{"--all"}},
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
		m.rows = max(msg.Height-4, 1)
	case tea.KeyPressMsg:
		if m.filter != nil {
			return m.filterKey(msg)
		}
		switch msg.String() {
		case "/":
			m.filter = &filter{editing: true, pane: pane{title: "search", op: "search"}}
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

// filterKey handles a key while the filter is open. While editing, text
// types into the filter and enter keeps the hits for scrolling; esc always
// closes the filter and restores the pane.
func (m *Model) filterKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	f := m.filter
	switch k := msg.String(); {
	case k == "ctrl+c":
		return m, tea.Quit
	case k == "esc":
		m.filter = nil
	case f.editing && k == "enter":
		f.editing = false
	case f.editing && k == "backspace":
		if r := []rune(f.text); len(r) > 0 {
			f.text = string(r[:len(r)-1])
			m.search()
		}
	case f.editing && msg.Text != "":
		f.text += msg.Text
		m.search()
	case !f.editing && (k == "/"):
		f.editing = true
	case !f.editing:
		return m.scrollKey(k)
	}
	return m, nil
}

// search reads the hits of the filter's text through the search operation.
func (m *Model) search() {
	p := &m.filter.pane
	p.top = 0
	if strings.TrimSpace(m.filter.text) == "" {
		p.err, p.lines = nil, nil
		return
	}
	out, err := m.read("search", []string{m.filter.text})
	p.err = err
	p.lines = strings.Split(strings.TrimRight(out, "\n"), "\n")
}

// scrollKey scrolls the shown pane, the hits while a filter is open.
func (m *Model) scrollKey(k string) (tea.Model, tea.Cmd) {
	switch k {
	case "down", "j":
		m.scroll(1)
	case "up", "k":
		m.scroll(-1)
	case "pgdown", "space":
		m.scroll(m.rows)
	case "pgup":
		m.scroll(-m.rows)
	case "home", "g":
		m.scroll(-len(m.shown().lines))
	case "end", "G":
		m.scroll(len(m.shown().lines))
	}
	return m, nil
}

// shown is the pane on screen: the filter's hits while one is open.
func (m *Model) shown() *pane {
	if m.filter != nil {
		return &m.filter.pane
	}
	return m.panes[m.cur]
}

// shownRows is how many lines of the shown pane fit: one fewer while the
// filter's own line is on screen.
func (m *Model) shownRows() int {
	if m.filter != nil {
		return max(m.rows-1, 0)
	}
	return m.rows
}

func (m *Model) scroll(d int) {
	p := m.shown()
	p.top = max(0, min(p.top+d, len(p.lines)-m.shownRows()))
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
	p := m.shown()
	if f := m.filter; f != nil {
		cursor := ""
		if f.editing {
			cursor = "_"
		}
		b.WriteString("/" + f.text + cursor + "\n")
	}
	if p.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", p.err)
	} else {
		end := min(p.top+m.shownRows(), len(p.lines))
		b.WriteString(strings.Join(p.lines[min(p.top, end):end], "\n"))
	}
	if m.filter != nil {
		b.WriteString("\n\nesc clear · enter keep · j/k scroll")
	} else {
		b.WriteString("\n\ntab switch · / search · j/k scroll · q quit")
	}
	return tea.NewView(b.String())
}

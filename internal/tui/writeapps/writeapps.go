// Package writeapps is the rhiza apps that write: change show, propose,
// accept and land on a change set, and record add, set, link and retire.
// Each runs the api operation of the same name in-process, twice at most:
// first as a dry run, shown as a what-if preview, then for real only when
// the owner confirms once. Cancelling writes nothing. An app starts no
// daemon and binds no socket.
package writeapps

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// App is one write app. Name is the app's name; Op is the api operation it
// runs.
type App struct {
	Name    string
	Op      string
	Summary string
	// Record runs the operation as `rhiza <op> --record` does: in a record
	// change set landed in one step.
	Record bool
	// ReadOnly apps write nothing, so they show their result and ask for no
	// confirmation.
	ReadOnly bool
	// Prompt, when set, is what the app asks for when started without
	// arguments: the arguments of the matching command, on one line.
	Prompt string
	// Required means an empty answer to Prompt is not accepted.
	Required bool
}

// Apps are the write apps, in the order they are registered.
var Apps = []App{
	{Name: "change-show", Op: "change.show", Summary: "A change set: its intent, state and operations", ReadOnly: true},
	{Name: "propose", Op: "propose", Summary: "Propose the current change set for review"},
	{Name: "accept", Op: "accept", Summary: "Accept a proposed change set", Prompt: "accept (--reason TEXT)"},
	{Name: "land", Op: "land", Summary: "Land the accepted change set"},
	{Name: "record-add", Op: "add", Record: true, Summary: "Record a new node", Prompt: "add --record (KIND HANDLE --field name=text ...)", Required: true},
	{Name: "record-set", Op: "set", Record: true, Summary: "Record a change to a node", Prompt: "set --record (REF --field name=text ...)", Required: true},
	{Name: "record-link", Op: "link", Record: true, Summary: "Record a link between two nodes", Prompt: "link --record (FROM TYPE TO ...)", Required: true},
	{Name: "record-retire", Op: "retire", Record: true, Summary: "Record a node's retirement", Prompt: "retire --record (REF --reason TEXT)", Required: true},
}

func init() {
	for _, a := range Apps {
		launcher.Register(launcher.App{Name: a.Name, Summary: a.Summary})
	}
}

// Lookup returns the write app named name.
func Lookup(name string) (App, bool) {
	for _, a := range Apps {
		if a.Name == name {
			return a, true
		}
	}
	return App{}, false
}

// Runner runs the app's operation on args and returns its result as
// `rhiza <cmd> --json` prints it. With dry it runs the operation as a dry
// run: everything built, nothing written.
type Runner func(a App, args []string, dry bool) ([]byte, error)

type stage int

const (
	asking stage = iota
	previewing
	done
)

// Model is the app's bubbletea program: a prompt for the arguments when the
// app needs them, the what-if preview with a single confirmation, then the
// result of the write. n, esc or q cancels at the preview, and nothing has
// been written.
type Model struct {
	app       App
	run       Runner
	args      []string
	input     string
	stage     stage
	preview   []byte
	out       []byte
	err       error
	cancelled bool
	lines     []string
	top, rows int
}

// NewModel returns the model of a over run, started with args. An app that
// has a prompt and no args asks for them first; any other goes straight to
// its preview.
func NewModel(a App, run Runner, args []string) *Model {
	m := &Model{app: a, run: run, rows: 24}
	if a.Prompt != "" && len(args) == 0 {
		m.stage = asking
		return m
	}
	m.start(args)
	return m
}

func (m *Model) start(args []string) {
	m.args = args
	if m.app.ReadOnly {
		m.out, m.err = m.run(m.app, args, false)
		m.preview = m.out
		m.show(m.out, m.err, done)
		return
	}
	m.preview, m.err = m.run(m.app, args, true)
	if m.err != nil {
		m.show(nil, m.err, done)
		return
	}
	m.show(m.preview, nil, previewing)
}

func (m *Model) show(b []byte, err error, s stage) {
	m.stage, m.err, m.top = s, err, 0
	m.lines = strings.Split(strings.TrimRight(string(b), "\n"), "\n")
}

// Preview is the what-if the app shows: the operation's dry-run result.
// An app still asking for its arguments has none.
func (m *Model) Preview() ([]byte, error) {
	if m.stage == asking {
		return nil, fmt.Errorf("%s needs arguments", m.app.Name)
	}
	return m.preview, m.err
}

// Confirm writes: it runs the operation for real and returns its result,
// byte for byte what `rhiza <cmd> --json` prints. It only follows a
// preview that succeeded.
func (m *Model) Confirm() ([]byte, error) {
	if m.stage != previewing {
		return nil, fmt.Errorf("%s has nothing to confirm", m.app.Name)
	}
	m.out, m.err = m.run(m.app, m.args, false)
	m.show(m.out, m.err, done)
	return m.out, m.err
}

// Cancelled reports whether the owner declined the preview.
func (m *Model) Cancelled() bool { return m.cancelled }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.rows = max(msg.Height-3, 1)
	case tea.KeyPressMsg:
		switch m.stage {
		case asking:
			return m.askKey(msg)
		case previewing:
			switch msg.String() {
			case "y", "enter":
				m.Confirm()
				return m, nil
			case "n", "q", "esc", "ctrl+c":
				m.cancelled = true
				return m, tea.Quit
			}
		case done:
			if k := msg.String(); k == "q" || k == "esc" || k == "ctrl+c" {
				return m, tea.Quit
			}
		}
		m.scrollKey(msg.String())
	}
	return m, nil
}

func (m *Model) scrollKey(k string) {
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
		m.top = 0
	case "end", "G":
		m.scroll(len(m.lines))
	}
}

func (m *Model) askKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.cancelled = true
		return m, tea.Quit
	case "enter":
		args, err := Split(m.input)
		if err != nil {
			m.err = err
		} else if len(args) > 0 || !m.app.Required {
			m.err = nil
			m.start(args)
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
	var b strings.Builder
	if m.stage == asking {
		if m.err != nil {
			fmt.Fprintf(&b, "rhiza: %v\n", m.err)
		}
		fmt.Fprintf(&b, "%s: %s", m.app.Prompt, m.input)
		return tea.NewView(b.String())
	}
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	switch m.stage {
	case previewing:
		b.WriteString("What if: nothing is written until you confirm\n")
	case done:
		if m.app.ReadOnly || m.err != nil {
			break
		}
		b.WriteString("Written\n")
	}
	end := min(m.top+m.rows, len(m.lines))
	b.WriteString(strings.Join(m.lines[m.top:end], "\n"))
	b.WriteString("\n\n" + m.app.Name + " · ")
	if m.stage == previewing {
		b.WriteString("y write · n cancel · j/k scroll")
	} else {
		b.WriteString("j/k scroll · q quit")
	}
	return tea.NewView(b.String())
}

// Split cuts a line into arguments as a shell would for plain words and
// single or double quotes; there is no escaping beyond the quotes.
func Split(line string) ([]string, error) {
	var args []string
	var cur strings.Builder
	var quote rune
	has := false
	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, has = r, true
		case r == ' ' || r == '\t':
			if has || cur.Len() > 0 {
				args = append(args, cur.String())
				cur.Reset()
				has = false
			}
		default:
			cur.WriteRune(r)
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated %c quote", quote)
	}
	if has || cur.Len() > 0 {
		args = append(args, cur.String())
	}
	return args, nil
}

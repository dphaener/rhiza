// Package requestsapp is the requests app: it files a feature, bug or
// milestone plan request from typed text and lists every request with its
// state, run and planned change set. It runs the api operations request,
// request.ls and run.open in-process through the registry, so --once prints
// what `rhiza request ls --json` prints. The view reloads every few seconds
// and on r. Enter focuses the Rex tab of the selected request's plan run.
package requestsapp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name.
const Name = "requests"

// Interval is how often the view reloads.
const Interval = 3 * time.Second

// MaxText bounds a request's typed text.
const MaxText = 500

// kinds are the request kinds in the order tab cycles them.
var kinds = []string{"feature", "bug", "milestone"}

func init() {
	// Key is fixed so that adding the app moves no other app's default key.
	launcher.Register(launcher.App{Name: Name, Key: "q", Summary: "File feature, bug and milestone requests and follow their planning"})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// TickMsg is the reload timer firing.
type TickMsg struct{}

// Model is the app's bubbletea program. q or esc quits; n, f, b and m open
// the text entry.
type Model struct {
	run      Runner
	out      []byte
	requests []api.RequestInfo
	err      error
	cursor   int
	top      int
	height   int
	notice   string
	// typing is true while the text entry is open; kind is its kind and
	// text what was typed.
	typing bool
	kind   string
	text   []rune
	// milestone is the entry's milestone field: a milestone intent's handle
	// or ID, empty for auto (the planner picks the open milestone it fits).
	// onMilestone is true while typing edits it instead of the text.
	milestone   []rune
	onMilestone bool
}

// NewModel returns the model over run, loaded once.
func NewModel(run Runner) *Model {
	m := &Model{run: run, height: 24, kind: kinds[0]}
	m.Reload()
	return m
}

// Reload runs request.ls again. The cursor stays on the same request when it
// is still listed.
func (m *Model) Reload() {
	var kept string
	if q, ok := m.Selected(); ok {
		kept = q.ID
	}
	out, err := m.run("request.ls", nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var ls api.RequestLsResponse
	if err := json.Unmarshal(out, &ls); err != nil {
		m.err = err
		return
	}
	m.requests = ls.Requests
	for i, q := range m.requests {
		if q.ID == kept {
			m.cursor = i
		}
	}
	m.cursor = max(0, min(m.cursor, len(m.requests)-1))
	m.follow()
}

// Result is the request.ls output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Selected is the request under the cursor.
func (m *Model) Selected() (api.RequestInfo, bool) {
	if m.cursor < 0 || m.cursor >= len(m.requests) {
		return api.RequestInfo{}, false
	}
	return m.requests[m.cursor], true
}

// Typing says whether the text entry is open.
func (m *Model) Typing() bool { return m.typing }

// Start opens the text entry for a request of kind (feature, bug or
// milestone).
func (m *Model) Start(kind string) {
	m.typing, m.kind, m.text, m.notice = true, kind, nil, ""
	m.milestone, m.onMilestone = nil, false
}

// Type adds s to the typed text, up to MaxText characters.
func (m *Model) Type(s string) {
	for _, r := range s {
		switch {
		case r < ' ':
		case m.onMilestone:
			if r != ' ' && len(m.milestone) < MaxText {
				m.milestone = append(m.milestone, r)
			}
		case len(m.text) < MaxText:
			m.text = append(m.text, r)
		}
	}
}

// RequestArgs is the argv of a request op for text of kind. A milestone
// (handle or ID) adds --milestone; empty is auto, with no flag. The text
// goes after "--" so text starting with a dash is not read as a flag.
func RequestArgs(kind, text, milestone string) []string {
	args := []string{"--kind", kind}
	if milestone != "" {
		args = append(args, "--milestone", milestone)
	}
	return append(args, "--", text)
}

// Submit files the typed text as a request of the chosen kind and closes
// the entry. Empty text is refused with a notice and the entry stays open.
func (m *Model) Submit() {
	text := strings.Join(strings.Fields(string(m.text)), " ")
	if text == "" {
		m.notice = "A request needs its text"
		return
	}
	out, err := m.run("request", RequestArgs(m.kind, text, string(m.milestone)))
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	var r api.RequestResponse
	if err := json.Unmarshal(out, &r); err != nil {
		m.notice = fmt.Sprintf("request: %v", err)
		return
	}
	m.typing, m.text, m.milestone, m.onMilestone = false, nil, nil, false
	m.notice = fmt.Sprintf("Queued %s request %s", r.Request.Kind, r.Request.ID)
	m.Reload()
	for i, q := range m.requests {
		if q.ID == r.Request.ID {
			m.cursor = i
		}
	}
	m.follow()
}

// rowLines is how many lines a request takes: a blank line, its text, its
// state line.
const rowLines = 3

func (m *Model) visible() int {
	fixed := 2 // the key help and the blank line before it
	if m.err != nil {
		fixed++
	}
	if m.notice != "" {
		fixed++
	}
	if m.typing {
		fixed += 3
	}
	return max((m.height-fixed)/rowLines, 1)
}

func (m *Model) follow() {
	n := m.visible()
	if m.cursor < m.top {
		m.top = m.cursor
	}
	if m.cursor >= m.top+n {
		m.top = m.cursor - n + 1
	}
	m.top = max(0, min(m.top, len(m.requests)-n))
}

func tick() tea.Cmd {
	return tea.Tick(Interval, func(time.Time) tea.Msg { return TickMsg{} })
}

func (m *Model) Init() tea.Cmd { return tick() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if !m.typing {
			m.Reload()
		}
		return m, tick()
	case tea.WindowSizeMsg:
		m.height = max(msg.Height, 1)
		m.follow()
	case tea.KeyPressMsg:
		if m.typing {
			m.typingKey(msg)
			return m, nil
		}
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.cursor = max(0, min(m.cursor+1, len(m.requests)-1))
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "r":
			m.Reload()
		case "n", "f":
			m.Start("feature")
		case "b":
			m.Start("bug")
		case "m":
			m.Start("milestone")
		case "enter":
			m.focus()
		}
		m.follow()
	}
	return m, nil
}

func (m *Model) typingKey(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.typing, m.text, m.notice = false, nil, ""
		m.milestone, m.onMilestone = nil, false
	case "ctrl+o":
		m.onMilestone = !m.onMilestone
	case "enter":
		m.Submit()
	case "tab":
		for i, k := range kinds {
			if k == m.kind {
				m.kind = kinds[(i+1)%len(kinds)]
				break
			}
		}
	case "backspace":
		if m.onMilestone {
			if len(m.milestone) > 0 {
				m.milestone = m.milestone[:len(m.milestone)-1]
			}
		} else if len(m.text) > 0 {
			m.text = m.text[:len(m.text)-1]
		}
	default:
		m.Type(msg.Text)
	}
}

// focus focuses the Rex tab of the selected request's plan run when it is
// live. A queued request has no run; a finished run calls nothing.
func (m *Model) focus() {
	q, ok := m.Selected()
	if !ok {
		return
	}
	if q.Run == "" {
		m.notice = "Request " + q.ID + " is still queued: it has no run yet"
		return
	}
	out, err := m.run("run.open", []string{"--run", q.Run})
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	var o api.RunOpenResponse
	if err := json.Unmarshal(out, &o); err != nil {
		m.notice = fmt.Sprintf("run.open: %v", err)
		return
	}
	switch {
	case o.Opened:
		m.notice = fmt.Sprintf("Focused Rex tab %s of run %s", o.Tab, q.Run)
	case !o.Live:
		m.notice = "Run " + q.Run + " is not live"
	default:
		m.notice = fmt.Sprintf("The tab of run %s is gone", q.Run)
	}
}

func (m *Model) View() tea.View {
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	if len(m.requests) == 0 && m.err == nil {
		b.WriteString("No requests yet: press n for a feature, b for a bug or m for a milestone.\nThe daemon plans the oldest queued request when a slot is free.\n")
	}
	end := min(m.top+m.visible(), len(m.requests))
	for i := m.top; i < end; i++ {
		mark := "  "
		if i == m.cursor {
			mark = "> "
		}
		b.WriteString("\n" + mark + format(m.requests[i]))
	}
	if m.typing {
		cur := [2]string{"_", ""}
		if m.onMilestone {
			cur = [2]string{"", "_"}
		}
		ms := string(m.milestone)
		if ms == "" && !m.onMilestone {
			ms = "auto"
		}
		fmt.Fprintf(&b, "\nNew %s request (tab changes kind):\n> %s%s\nmilestone: %s%s", m.kind, string(m.text), cur[0], ms, cur[1])
	}
	if m.notice != "" {
		b.WriteString("\n" + m.notice)
	}
	if m.typing {
		b.WriteString("\n" + Name + " · enter file · tab kind · ctrl+o milestone (empty = auto) · esc cancel")
	} else {
		b.WriteString("\n" + Name + " · j/k move · n/f feature · b bug · m milestone · enter focus run · r reload · q quit")
	}
	return tea.NewView(b.String())
}

// format is one request: kind and text, then state, run and change set.
func format(q api.RequestInfo) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s · %s\n    %s", q.Kind, q.Text, q.State)
	if q.Run != "" {
		fmt.Fprintf(&b, " · run %s", q.Run)
	}
	if q.Change != "" {
		fmt.Fprintf(&b, " · change set %s", q.Change)
	}
	b.WriteString("\n")
	return b.String()
}

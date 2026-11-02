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

// Exec runs a program, as api.Ctx.Exec does: the app opens run tabs
// through it with rex, so the runs that retry, approve and request changes
// continue run there and the app stays responsive.
type Exec func(name string, args []string, dir string, env []string) ([]byte, error)

// The finding an owner's typed changes become: the review step's quality
// check, at major severity, placed on the owner's own words.
const (
	changesCheck    = "quality"
	changesSeverity = "major"
	changesPlace    = "owner"
)

// The modes the app takes typed text or a y/n in.
const (
	modeNone    = ""
	modeAnswer  = "answer"
	modeChanges = "changes"
	modeAbandon = "abandon"
)

// TickMsg is the reload timer firing.
type TickMsg struct{}

// groups are the kinds in the order the list shows them, with their headings.
var groups = []struct{ kind, heading string }{
	{api.InboxQuestion, "Questions"},
	{api.InboxSignoff, "Sign-offs"},
	{api.InboxOpen, "Open questions"},
	{api.InboxBlocker, "Blocking questions"},
	{api.InboxOverlap, "Overlaps"},
	{api.InboxFailed, "Failed runs"},
	{api.InboxGaveUp, "Gave up"},
	{api.InboxSkipped, "Skipped"},
	{api.InboxStalled, "Stalled"},
	{api.InboxDiskLow, "Disk space"},
	{api.InboxMainStale, "Main checkout"},
}

// Model is the app's bubbletea program. q or esc quits.
type Model struct {
	run    Runner
	out    []byte
	items  []api.InboxItem // in display order, grouped by kind
	err    error
	cursor int

	// Exec, Self, Root and RexSession open a run's tab: rex window new
	// running Self run resume, approve or request-changes. Without Exec
	// or Self, R, y and c say they cannot.
	Exec       Exec
	Self       string
	Root       string
	RexSession string

	mode  string
	input string
	note  string
	// target is the item the mode was entered on, which the mode acts on
	// whatever the cursor is on by then.
	target api.InboxItem
}

// NewModel returns the model over run, loaded once.
func NewModel(run Runner) *Model {
	m := &Model{run: run}
	m.Reload()
	return m
}

// Reload runs the operation again. The cursor stays on the same item when it
// still waits, else on the nearest row. A mode whose item no longer waits
// is cancelled.
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
	if m.mode != modeNone && !m.waits(m.target.ID) {
		m.mode, m.input = modeNone, ""
		m.note = "that item no longer waits"
	}
}

// waits says whether the item with ID id is in the list.
func (m *Model) waits(id string) bool {
	for _, it := range m.items {
		if it.ID == id {
			return true
		}
	}
	return false
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
		if m.mode == modeNone {
			return m, m.normalKey(msg)
		}
		m.modeKey(msg)
	}
	return m, nil
}

// normalKey handles a key outside a mode.
func (m *Model) normalKey(msg tea.KeyPressMsg) tea.Cmd {
	it, has := m.Selected()
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return tea.Quit
	case "down", "j":
		m.cursor = max(0, min(m.cursor+1, len(m.items)-1))
	case "up", "k":
		m.cursor = max(0, m.cursor-1)
	case "r":
		m.note = ""
		m.Reload()
	case "enter", "o":
		if has {
			m.act("open", Name+".open", it.ID)
		}
	case "a":
		if has && (it.Kind == api.InboxQuestion || it.Kind == api.InboxBlocker) {
			m.mode, m.input, m.target = modeAnswer, "", it
		}
	case "t":
		if has && it.Kind == api.InboxQuestion {
			if it.Recommendation == "" {
				m.note = "the advisor has no recommendation yet"
			} else {
				m.act("answer", "session.message", it.Session, it.Recommendation)
			}
		}
	case "y":
		if has && it.Kind == api.InboxSignoff {
			m.tab(it, "approve", m.Self, "run", "approve", it.Run)
		}
	case "c":
		if has && it.Kind == api.InboxSignoff {
			m.mode, m.input, m.target = modeChanges, "", it
		}
	case "n":
		if has && it.Kind == api.InboxStalled {
			m.nudge(it)
		}
	case "R":
		if has && retryable(it) {
			m.resume(it)
		}
	case "x":
		if has && retryable(it) {
			m.mode, m.target = modeAbandon, it
		}
	}
	return nil
}

// carryOn is the message a nudge sends a stalled session.
const carryOn = "You have been silent for a while. Please carry on with your task, or say what is blocking you."

// nudge asks a stalled session to carry on. A session parked on a usage-limit
// prompt (its pane text, as its adapter reads it) gets Escape first.
func (m *Model) nudge(it api.InboxItem) {
	out, err := m.run("session.pane", []string{it.Session})
	if err != nil {
		m.note = fmt.Sprintf("nudge failed: %v", err)
		return
	}
	var pane api.SessionPaneResponse
	if err := json.Unmarshal(out, &pane); err != nil {
		m.note = fmt.Sprintf("nudge failed: %v", err)
		return
	}
	if pane.UsageLimit {
		if _, err := m.run("session.key", []string{it.Session, "escape"}); err != nil {
			m.note = fmt.Sprintf("nudge failed: %v", err)
			return
		}
	}
	m.act("nudge", "session.message", it.Session, carryOn)
}

// retryable says whether R and x work on the item: a failed, gave-up or
// skipped run.
func retryable(it api.InboxItem) bool {
	return it.Kind == api.InboxFailed || it.Kind == api.InboxGaveUp || it.Kind == api.InboxSkipped
}

// modeKey handles a key while typing an answer or changes, or at the
// y/n of an abandon. It acts on the item the mode was entered on.
func (m *Model) modeKey(msg tea.KeyPressMsg) {
	it, has := m.target, m.waits(m.target.ID)
	if m.mode == modeAbandon {
		m.mode = modeNone
		if msg.String() == "y" && has {
			m.act("abandon", "run.abandon", it.Run)
		} else {
			m.note = "not abandoned"
		}
		return
	}
	switch msg.String() {
	case "esc", "ctrl+c":
		m.mode, m.input = modeNone, ""
	case "backspace":
		if r := []rune(m.input); len(r) > 0 {
			m.input = string(r[:len(r)-1])
		}
	case "enter":
		mode, text := m.mode, strings.TrimSpace(m.input)
		m.mode, m.input = modeNone, ""
		if text == "" || !has {
			m.note = "nothing typed, nothing sent"
			return
		}
		if mode == modeAnswer && it.Kind == api.InboxBlocker {
			// A blocker has no asking session: the words become a Decision
			// in the blocked change set, which question.answer writes.
			m.act("answer", "question.answer", it.Handle, it.Change, text)
		} else if mode == modeAnswer {
			m.act("answer", "session.message", it.Session, text)
		} else {
			m.requestChanges(it, text)
		}
	default:
		m.input += msg.Text
	}
}

// act runs an operation for an action, reports its outcome and reloads
// the list.
func (m *Model) act(what, op string, args ...string) {
	if _, err := m.run(op, args); err != nil {
		m.note = fmt.Sprintf("%s failed: %v", what, err)
	} else {
		m.note = what + " done"
	}
	m.Reload()
}

// requestChanges asks for changes with the typed text as the one finding,
// in a new tab: the owner's findings go to rhiza run request-changes on
// stdin.
func (m *Model) requestChanges(it api.InboxItem, text string) {
	body, err := json.Marshal(map[string]any{"findings": []map[string]string{{
		"check": changesCheck, "severity": changesSeverity, "place": changesPlace, "reason": text,
	}}})
	if err != nil {
		m.note = fmt.Sprintf("request changes failed: %v", err)
		return
	}
	m.tab(it, "request changes", "sh", "-c", `printf %s "$2" | "$1" run request-changes - "$3"`, "sh", m.Self, string(body), it.Run)
}

// resume opens a new Rex tab running rhiza run resume on the item's run. A
// failed run has ended, so its resume restarts it at the step that failed.
func (m *Model) resume(it api.InboxItem) {
	cmd := []string{m.Self, "run", "resume", it.Run}
	if it.Kind == api.InboxFailed {
		cmd = append(cmd, "--retry")
	}
	m.tab(it, "retry", cmd...)
}

// tab opens a new Rex tab, named after the item's run, running cmd in the
// repo, so the run cmd continues runs there and not inside the app, and
// reloads the list.
func (m *Model) tab(it api.InboxItem, what string, cmd ...string) {
	if m.Exec == nil || m.Self == "" {
		m.note = "cannot open a tab: no rex or rhiza path"
		return
	}
	args := []string{"window", "new"}
	if m.RexSession != "" {
		args = append(args, "-s", m.RexSession)
	}
	args = append(args, it.Run[:min(len(it.Run), 9)], "-c", m.Root, "--json", "--")
	args = append(args, cmd...)
	if _, err := m.Exec("rex", args, m.Root, nil); err != nil {
		m.note = fmt.Sprintf("%s failed: %v", what, err)
	} else {
		m.note = what + " opened in a new tab"
	}
	m.Reload()
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
	switch m.mode {
	case modeAnswer:
		b.WriteString("\nanswer: " + m.input + "▏ · enter send · esc cancel")
	case modeChanges:
		b.WriteString("\nchanges: " + m.input + "▏ · enter send · esc cancel")
	case modeAbandon:
		b.WriteString("\nAbandon this run? y/n")
	default:
		if m.note != "" {
			b.WriteString("\n" + m.note)
		}
		b.WriteString("\n" + Name + " · j/k move · o open · a answer · t take advice · n nudge · y approve · c changes · R retry · x abandon · r reload · q quit")
	}
	return tea.NewView(b.String())
}

// row is one item: its change set's title and step, then what it asks or
// why it failed, and a question's recommendation.
func row(it api.InboxItem) string {
	title := it.ChangeTitle
	if title == "" {
		title = it.Change
	}
	if it.Kind == api.InboxDiskLow {
		title = "Daemon holds starts"
	}
	if it.Kind == api.InboxMainStale {
		title = "Main not updated"
	}
	var b strings.Builder
	b.WriteString(title)
	if it.Step != "" {
		b.WriteString(" · " + it.Step)
	}
	if it.Kind == api.InboxOverlap {
		b.WriteString(" · with " + it.Other + " · " + it.Target)
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

// Package autonomyapp is the autonomy app: every daemon and autonomy
// setting with its value and source, and the daemon's own state. Enter
// toggles a true/false or auto/owner setting and opens the entry for a
// number or word; each write goes through rhiza daemon set. d starts the
// daemon in a Rex tab named "rhiza daemon" when none runs, or stops the one
// that does. --once prints what `rhiza daemon settings --json` prints.
package autonomyapp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name.
const Name = "autonomy"

// Key is the app's key in the rhiza key mode.
const Key = "y"

// Interval is how often the view reloads.
const Interval = 3 * time.Second

// TabName is the name of the Rex tab the daemon runs in.
const TabName = "rhiza daemon"

// StartChecks is how many reloads after a start the app waits for the
// daemon to report running before it says the start failed.
const StartChecks = 4

func init() {
	launcher.Register(launcher.App{Name: Name, Key: Key, Summary: "Edit the daemon and autonomy settings, start or stop the daemon"})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// Starter is how the app reaches Rex and the daemon's process: Exec runs a
// program and returns its output, Self is the rhiza binary the tab runs,
// Root the project root, RexSession the Rex session (maybe empty) and Host
// this machine's name, which a daemon must share to be stopped from here.
type Starter struct {
	Exec       func(name string, args []string, dir string) ([]byte, error)
	Self       string
	Root       string
	RexSession string
	Host       string
}

// TickMsg is the reload timer firing.
type TickMsg struct{}

// Model is the app's bubbletea program. q or esc quits.
type Model struct {
	run      Runner
	starter  Starter
	out      []byte
	settings []api.Setting
	status   api.DaemonStatusResponse
	err      error
	cursor   int
	notice   string
	// editing is true while a value is typed for the setting under the
	// cursor; text is what was typed.
	editing bool
	text    []rune
	// confirm is "start" or "stop" while the owner is asked y/n.
	confirm string
	// starting counts the reloads left to see a just-started daemon
	// running; 0 when no start is pending.
	starting int
}

// NewModel returns the model over run, loaded once.
func NewModel(run Runner) *Model {
	m := &Model{run: run}
	m.Reload()
	return m
}

// SetStarter gives the model the means to start and stop the daemon.
func (m *Model) SetStarter(s Starter) { m.starter = s }

// DefaultStarter fills Host from the machine.
func DefaultStarter(s Starter) Starter {
	if s.Host == "" {
		s.Host, _ = os.Hostname()
	}
	return s
}

// Reload reads the settings and the daemon's state again. The cursor stays
// on the same setting.
func (m *Model) Reload() {
	out, err := m.run("daemon.settings", nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var s api.DaemonSettingsResponse
	if err := json.Unmarshal(out, &s); err != nil {
		m.err = err
		return
	}
	m.settings = s.Settings
	m.cursor = max(0, min(m.cursor, len(m.settings)-1))
	st, err := m.run("daemon.status", nil)
	if err != nil {
		m.err = err
		return
	}
	var d api.DaemonStatusResponse
	if err := json.Unmarshal(st, &d); err != nil {
		m.err = err
		return
	}
	m.status = d
	m.checkStart()
}

// checkStart settles a pending start: the daemon running, or, once
// StartChecks reloads pass without it, a failure pointing at its tab, which
// stays open with the daemon's error.
func (m *Model) checkStart() {
	if m.starting == 0 {
		return
	}
	if m.status.Running {
		m.starting = 0
		m.notice = "The daemon is running in the Rex tab " + TabName
		return
	}
	m.starting--
	if m.starting == 0 {
		m.notice = "The daemon did not start: its error is in the Rex tab " + TabName
	}
}

// Result is the daemon.settings output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Selected is the setting under the cursor.
func (m *Model) Selected() (api.Setting, bool) {
	if m.cursor < 0 || m.cursor >= len(m.settings) {
		return api.Setting{}, false
	}
	return m.settings[m.cursor], true
}

// Editing says whether a value is being typed.
func (m *Model) Editing() bool { return m.editing }

// Running says whether a daemon runs.
func (m *Model) Running() bool { return m.status.Running }

// Notice is the last message shown.
func (m *Model) Notice() string { return m.notice }

// toggled is the other value of a two-valued setting, and false for a
// setting that takes a number or a word.
func toggled(s api.Setting) (string, bool) {
	switch s.Value {
	case "true":
		return "false", true
	case "false":
		return "true", true
	case api.AutonomyAuto:
		return api.AutonomyOwner, true
	case api.AutonomyOwner:
		return api.AutonomyAuto, true
	}
	return "", false
}

// Set writes value to key through rhiza daemon set and reloads.
func (m *Model) Set(key, value string) {
	if _, err := m.run("daemon.set", []string{"--", key, value}); err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.notice = fmt.Sprintf("Set %s to %s", key, value)
	m.Reload()
}

// Activate is enter: toggle a two-valued setting, else open the entry.
func (m *Model) Activate() {
	s, ok := m.Selected()
	if !ok {
		return
	}
	if v, ok := toggled(s); ok {
		m.Set(s.Key, v)
		return
	}
	m.editing, m.text, m.notice = true, []rune(s.Value), ""
}

// Type adds s to the typed value.
func (m *Model) Type(s string) {
	for _, r := range s {
		if r >= ' ' && len(m.text) < 100 {
			m.text = append(m.text, r)
		}
	}
}

// Submit writes the typed value. A refusal keeps the entry open.
func (m *Model) Submit() {
	s, ok := m.Selected()
	if !ok {
		m.editing = false
		return
	}
	v := strings.TrimSpace(string(m.text))
	if v == "" {
		m.notice = "A setting needs a value"
		return
	}
	if _, err := m.run("daemon.set", []string{"--", s.Key, v}); err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.editing, m.text = false, nil
	m.notice = fmt.Sprintf("Set %s to %s", s.Key, v)
	m.Reload()
}

// Ask puts the y/n question to start the daemon when none runs, or to stop
// the one that does.
func (m *Model) Ask() {
	m.notice = ""
	if m.status.Running {
		m.confirm = "stop"
	} else {
		m.confirm = "start"
	}
}

// Confirming is the pending question's kind, "" when none.
func (m *Model) Confirming() string { return m.confirm }

// Answer settles the pending question: yes runs it.
func (m *Model) Answer(yes bool) {
	kind := m.confirm
	m.confirm = ""
	if !yes {
		return
	}
	switch kind {
	case "start":
		m.startDaemon()
	case "stop":
		m.stopDaemon()
	}
	m.Reload()
}

// StartArgs is the rex call that opens the daemon's tab. The tab stays
// open after the daemon exits, so an error it stops on can be read there.
func (s Starter) StartArgs() []string {
	args := []string{"window", "new"}
	if s.RexSession != "" {
		args = append(args, "-s", s.RexSession)
	}
	return append(args, TabName, "-c", s.Root, "--keep-open", "--json", "--", s.Self, "daemon", "--autonomous")
}

func (m *Model) startDaemon() {
	if m.starter.Exec == nil || m.starter.Self == "" {
		m.notice = "Cannot start the daemon: no rex or rhiza path"
		return
	}
	if _, err := m.starter.Exec("rex", m.starter.StartArgs(), m.starter.Root); err != nil {
		m.notice = fmt.Sprintf("rex window new failed: %v", err)
		return
	}
	m.notice = "Starting the daemon in the Rex tab " + TabName
	m.starting = StartChecks
}

func (m *Model) stopDaemon() {
	host, pid, ok := strings.Cut(m.status.Process, ":")
	if !ok || pid == "" {
		m.notice = "The daemon's process is unknown"
		return
	}
	if m.starter.Exec == nil {
		m.notice = "Cannot stop the daemon: no way to signal it"
		return
	}
	if m.starter.Host != "" && host != m.starter.Host {
		m.notice = fmt.Sprintf("The daemon runs on %s, not this machine", host)
		return
	}
	if _, err := m.starter.Exec("kill", []string{"-TERM", pid}, m.starter.Root); err != nil {
		m.notice = fmt.Sprintf("kill failed: %v", err)
		return
	}
	m.notice = "Asked the daemon " + m.status.Process + " to stop"
}

func tick() tea.Cmd {
	return tea.Tick(Interval, func(time.Time) tea.Msg { return TickMsg{} })
}

func (m *Model) Init() tea.Cmd { return tick() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if !m.editing {
			m.Reload()
		}
		return m, tick()
	case tea.KeyPressMsg:
		switch {
		case m.editing:
			m.editKey(msg)
		case m.confirm != "":
			switch msg.String() {
			case "y":
				m.Answer(true)
			case "n", "esc", "ctrl+c":
				m.Answer(false)
			}
		default:
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			case "down", "j":
				m.cursor = max(0, min(m.cursor+1, len(m.settings)-1))
			case "up", "k":
				m.cursor = max(0, m.cursor-1)
			case "enter", "e":
				m.Activate()
			case "d":
				m.Ask()
			case "r":
				m.Reload()
			}
		}
	}
	return m, nil
}

func (m *Model) editKey(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.editing, m.text, m.notice = false, nil, ""
	case "enter":
		m.Submit()
	case "backspace":
		if len(m.text) > 0 {
			m.text = m.text[:len(m.text)-1]
		}
	default:
		m.Type(msg.Text)
	}
}

func (m *Model) View() tea.View {
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	b.WriteString(m.daemonLine() + "\n\n")
	for i, s := range m.settings {
		mark := "  "
		if i == m.cursor {
			mark = "> "
		}
		v := s.Value
		if v == "" {
			v = "(none)"
		}
		fmt.Fprintf(&b, "%s%-30s %-10s %s\n", mark, s.Key, v, s.Source)
	}
	if m.editing {
		s, _ := m.Selected()
		fmt.Fprintf(&b, "\nNew value for %s:\n> %s_\n", s.Key, string(m.text))
	}
	switch m.confirm {
	case "start":
		b.WriteString("\nStart the daemon in a new Rex tab? y/n\n")
	case "stop":
		fmt.Fprintf(&b, "\nStop the daemon %s? y/n\n", m.status.Process)
	}
	if m.notice != "" {
		b.WriteString("\n" + m.notice + "\n")
	}
	switch {
	case m.editing:
		b.WriteString("\n" + Name + " · enter save · esc cancel")
	case m.confirm != "":
		b.WriteString("\n" + Name + " · y yes · n no")
	default:
		b.WriteString("\n" + Name + " · j/k move · enter toggle or edit · d start/stop daemon · r reload · q quit")
	}
	return tea.NewView(b.String())
}

func (m *Model) daemonLine() string {
	if m.status.Running {
		return fmt.Sprintf("Daemon running, process %s", m.status.Process)
	}
	return "Daemon not running"
}

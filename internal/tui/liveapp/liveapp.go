// Package liveapp is the live app: the daemon's state and each live or
// recent run with its change set title, current step, leaves passed of the
// change set's leaves, checks and elapsed time. It runs the api operations
// daemon.status, run.ls, run.status and change.show in-process through the
// registry, so --once prints what `rhiza run ls --json` prints. The view
// reloads every few seconds and on r, and scrolls to keep the cursor's run
// on screen below the daemon line. change.show folds the model, so each
// change set's leaf total is fetched once, off the update loop, and kept;
// r fetches the totals again.
package liveapp

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
const Name = "live"

// Interval is how often the view reloads.
const Interval = 3 * time.Second

func init() {
	launcher.Register(launcher.App{Name: Name, Summary: "The daemon and each live or recent run's progress"})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// Clock tells the time; tests inject one.
type Clock func() time.Time

// TickMsg is the reload timer firing.
type TickMsg struct{}

// LeavesMsg carries leaf totals fetched off the update loop, by change set.
// A change set whose change.show failed is missing from it.
type LeavesMsg struct{ Totals map[string]int }

// row is one run as the view shows it.
type row struct {
	api.RunLsItem
	info    api.RunInfo
	hasInfo bool
	err     error
}

// rowLines is how many lines a run takes: a blank line, its title, its
// progress.
const rowLines = 3

// Model is the app's bubbletea program. q or esc quits.
type Model struct {
	run    Runner
	now    Clock
	out    []byte
	daemon api.DaemonStatusResponse
	hasD   bool
	rows   []row
	err    error
	cursor int
	top    int
	height int
	notice string
	// leaves holds each change set's leaf total from change.show; asked
	// holds the change sets fetched or being fetched.
	leaves map[string]int
	asked  map[string]bool
}

// NewModel returns the model over run, loaded once. A nil now means the
// wall clock.
func NewModel(run Runner, now Clock) *Model {
	if now == nil {
		now = time.Now
	}
	m := &Model{run: run, now: now, height: 24, leaves: map[string]int{}, asked: map[string]bool{}}
	m.Reload()
	return m
}

// Reload runs the operations again. The cursor stays on the same run when
// it is still listed, else on the nearest row.
func (m *Model) Reload() {
	var kept string
	if r, ok := m.Selected(); ok {
		kept = r.ID
	}
	out, err := m.run("run.ls", nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var ls api.RunLsResponse
	if err := json.Unmarshal(out, &ls); err != nil {
		m.err = err
		return
	}
	m.hasD = false
	if d, err := m.run("daemon.status", nil); err == nil && json.Unmarshal(d, &m.daemon) == nil {
		m.hasD = true
	}
	m.rows = m.rows[:0]
	for _, it := range ls.Runs {
		r := row{RunLsItem: it}
		if b, err := m.run("run.status", []string{it.ID}); err != nil {
			r.err = err
		} else if err := json.Unmarshal(b, &r.info); err != nil {
			r.err = err
		} else {
			r.hasInfo = true
		}
		m.rows = append(m.rows, r)
	}
	for i, r := range m.rows {
		if r.ID == kept {
			m.cursor = i
		}
	}
	m.cursor = max(0, min(m.cursor, len(m.rows)-1))
	m.follow()
}

// fetchLeaves fetches the leaf totals of the listed change sets not yet
// asked for, off the update loop. It is nil when there are none.
func (m *Model) fetchLeaves() tea.Cmd {
	var changes []string
	for _, r := range m.rows {
		if !m.asked[r.Change] {
			m.asked[r.Change] = true
			changes = append(changes, r.Change)
		}
	}
	if len(changes) == 0 {
		return nil
	}
	run := m.run
	return func() tea.Msg {
		totals := map[string]int{}
		for _, ch := range changes {
			b, err := run("change.show", []string{ch})
			var c api.ChangeShowResponse
			if err == nil && json.Unmarshal(b, &c) == nil {
				totals[ch] = len(c.Leaves)
			}
		}
		return LeavesMsg{Totals: totals}
	}
}

// visible is how many runs fit below the daemon line, above the notice
// and the key help.
func (m *Model) visible() int {
	fixed := 3 // the daemon line, the blank line and the key help
	if m.err != nil {
		fixed++
	}
	if m.notice != "" {
		fixed++
	}
	return max((m.height-fixed)/rowLines, 1)
}

// follow scrolls so the cursor's run is on screen.
func (m *Model) follow() {
	n := m.visible()
	if m.cursor < m.top {
		m.top = m.cursor
	}
	if m.cursor >= m.top+n {
		m.top = m.cursor - n + 1
	}
	m.top = max(0, min(m.top, len(m.rows)-n))
}

// Result is the run.ls output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Selected is the run under the cursor.
func (m *Model) Selected() (api.RunLsItem, bool) {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return api.RunLsItem{}, false
	}
	return m.rows[m.cursor].RunLsItem, true
}

func tick() tea.Cmd {
	return tea.Tick(Interval, func(time.Time) tea.Msg { return TickMsg{} })
}

func (m *Model) Init() tea.Cmd { return tea.Batch(tick(), m.fetchLeaves()) }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.Reload()
		return m, tea.Batch(tick(), m.fetchLeaves())
	case LeavesMsg:
		for ch, n := range msg.Totals {
			m.leaves[ch] = n
		}
	case tea.WindowSizeMsg:
		m.height = max(msg.Height, 1)
		m.follow()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.cursor = max(0, min(m.cursor+1, len(m.rows)-1))
		case "up", "k":
			m.cursor = max(0, m.cursor-1)
		case "r":
			clear(m.asked)
			m.Reload()
			m.follow()
			return m, m.fetchLeaves()
		case "enter":
			m.focus()
		case "l":
			m.openLog()
		}
		m.follow()
	}
	return m, nil
}

// focus focuses the Rex tab of the selected run when it is live. A finished
// run calls nothing.
func (m *Model) focus() {
	r, ok := m.Selected()
	if !ok {
		return
	}
	if !r.Live {
		m.notice = "Run " + r.ID + " is not live"
		return
	}
	out, err := m.run("run.open", []string{r.Change})
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
		m.notice = fmt.Sprintf("Focused Rex tab %s of run %s", o.Tab, o.Run)
	case !o.Live:
		m.notice = "No live run for " + r.Change
	default:
		m.notice = fmt.Sprintf("The tab of run %s is gone", o.Run)
	}
}

// openLog opens the selected run's driver log in a Rex layer, running
// rhiza run log for it.
func (m *Model) openLog() {
	r, ok := m.Selected()
	if !ok {
		return
	}
	out, err := m.run("run.log", []string{r.ID, "--layer"})
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	var l api.RunLogResponse
	if err := json.Unmarshal(out, &l); err != nil {
		m.notice = fmt.Sprintf("run.log: %v", err)
		return
	}
	m.notice = "Opened the log of run " + r.ID
}

func (m *Model) View() tea.View {
	var b strings.Builder
	if m.err != nil {
		fmt.Fprintf(&b, "rhiza: %v\n", m.err)
	}
	b.WriteString(m.daemonLine() + "\n")
	if len(m.rows) == 0 && m.err == nil {
		b.WriteString("\nNo runs\n")
	}
	end := min(m.top+m.visible(), len(m.rows))
	for i := m.top; i < end; i++ {
		mark := "  "
		if i == m.cursor {
			mark = "> "
		}
		b.WriteString("\n" + mark + m.format(m.rows[i]))
	}
	if m.notice != "" {
		b.WriteString("\n" + m.notice)
	}
	b.WriteString("\n" + Name + " · j/k move · enter focus tab · l log · r reload · q quit")
	return tea.NewView(b.String())
}

func (m *Model) daemonLine() string {
	switch {
	case !m.hasD:
		return "Daemon: unknown"
	case m.daemon.Running:
		return fmt.Sprintf("Daemon: running · %d live runs", len(m.daemon.Runs))
	default:
		return "Daemon: not running"
	}
}

// format is one run: title and state, then step, leaves, checks, elapsed.
func (m *Model) format(r row) string {
	title := r.Title
	if title == "" {
		title = r.Change
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s · %s\n", title, r.State)
	step := r.Current
	if step == "" {
		step = "-"
	}
	fmt.Fprintf(&b, "    step %s", step)
	if r.err != nil {
		fmt.Fprintf(&b, " · status: %v\n", r.err)
		return b.String()
	}
	if r.hasInfo {
		total, ok := m.leaves[r.Change]
		passed, total := leavesPassed(r, total, ok)
		fmt.Fprintf(&b, " · leaves %d/%d", passed, total)
		fmt.Fprintf(&b, " · checks %s", checksSummary(r.info.Checks))
	}
	fmt.Fprintf(&b, " · %s\n", m.elapsed(r.StartedAt))
	return b.String()
}

// leavesPassed is how many of the change set's leaves passed, and how many
// it has. A leaf counts once however many entries it has: run.status lists
// a leaf again for each checks step's fixer that fixed it, with no outcome.
// total is the change set's leaf count when hasTotal; without it the total
// falls back to the leaves the run lists.
func leavesPassed(r row, total int, hasTotal bool) (passed, _ int) {
	seen := map[string]bool{}
	done := map[string]bool{}
	for _, l := range r.info.Leaves {
		k := l.Leaf
		if k == "" {
			k = l.Handle
		}
		seen[k] = true
		if l.Outcome == "passed" && !done[k] {
			done[k] = true
			passed++
		}
	}
	if !hasTotal {
		total = len(seen)
	}
	return passed, total
}

// checksSummary counts the checks by outcome, e.g. "2 passed, 1 failed, 1
// running", or "none".
func checksSummary(cs []api.RunCheck) string {
	if len(cs) == 0 {
		return "none"
	}
	var passed, failed, open int
	for _, c := range cs {
		switch c.Outcome {
		case "passed":
			passed++
		case "failed":
			failed++
		default:
			open++
		}
	}
	var parts []string
	if passed > 0 {
		parts = append(parts, fmt.Sprintf("%d passed", passed))
	}
	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", failed))
	}
	if open > 0 {
		parts = append(parts, fmt.Sprintf("%d running", open))
	}
	return strings.Join(parts, ", ")
}

// elapsed is the time since the run started, by the injected clock.
func (m *Model) elapsed(started string) string {
	t, err := time.Parse(time.RFC3339, started)
	if err != nil {
		return "elapsed ?"
	}
	d := m.now().Sub(t).Round(time.Second)
	if d < 0 {
		d = 0
	}
	h, mi, s := int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60
	if h > 0 {
		return fmt.Sprintf("elapsed %dh%02dm%02ds", h, mi, s)
	}
	return fmt.Sprintf("elapsed %dm%02ds", mi, s)
}

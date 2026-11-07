// Package runsapp is the Runs app: a read-only glance at the live runs,
// each with its change set title, current step and elapsed time. It is a
// rexuikit InfoPanel hosted by the modal shell, over the api operation
// run.ls only, so --once prints what `rhiza run ls --json` prints. The
// list reloads on its own every few seconds. The app writes nothing.
package runsapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Name is the app's name.
const Name = "runs"

// Key is the app's key in the rhiza key mode.
const Key = "v"

// Interval is how often the list reloads.
const Interval = 3 * time.Second

// readOp is the one operation the app ever runs.
const readOp = "run.ls"

func init() {
	launcher.Register(launcher.App{Name: Name, Key: Key, Summary: "A glance at the live runs and their current steps"})
}

// Runner runs an api operation and returns its result as
// `rhiza <op> --json` prints it.
type Runner func(op string, args []string) ([]byte, error)

// ProjectRunner runs an api operation in the named project.
type ProjectRunner func(project, op string, args []string) ([]byte, error)

// Clock tells the time; tests inject one.
type Clock func() time.Time

// TickMsg is the reload timer firing.
type TickMsg struct{}

// Model is the app: a modalshell.Component over an InfoPanel.
type Model struct {
	run   Runner
	runIn ProjectRunner // set for a merged model: every project's runs
	projs []string
	tags  []string // the project of each row, parallel to rows
	now   Clock
	panel *ui.InfoPanel
	out   []byte
	err   error
	rows  []api.RunLsItem
}

var _ modalshell.Component = (*Model)(nil)

// NewModel returns the model over run, loaded once. A nil now means the
// wall clock.
func NewModel(run Runner, now Clock) *Model {
	if now == nil {
		now = time.Now
	}
	m := &Model{run: run, now: now}
	m.load()
	m.panel = ui.NewInfoPanel(theme.Dark(), colorprofile.TrueColor, "Runs", m.body).
		WithHint("↑↓ scroll · q close")
	return m
}

// NewMergedModel returns the model over every one of projects' live runs,
// each row tagged with its project, newest started first. A project whose
// runs cannot be read is skipped; the error shows only when none can.
func NewMergedModel(run ProjectRunner, projects []string, now Clock) *Model {
	if now == nil {
		now = time.Now
	}
	m := &Model{runIn: run, projs: projects, now: now}
	m.load()
	m.panel = ui.NewInfoPanel(theme.Dark(), colorprofile.TrueColor, "Runs", m.body).
		WithHint("↑↓ scroll · q close")
	return m
}

func (m *Model) loadMerged() {
	var errs []error
	var all []api.RunLsItem
	var tags []string
	for _, p := range m.projs {
		out, err := m.runIn(p, readOp, nil)
		var ls api.RunLsResponse
		if err == nil {
			err = json.Unmarshal(out, &ls)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", p, err))
			continue
		}
		for _, it := range ls.Runs {
			if it.Live {
				all = append(all, it)
				tags = append(tags, p)
			}
		}
	}
	idx := make([]int, len(all))
	for i := range idx {
		idx[i] = i
	}
	slices.SortStableFunc(idx, func(a, b int) int { return strings.Compare(all[b].StartedAt, all[a].StartedAt) })
	m.rows, m.tags = m.rows[:0], m.tags[:0]
	for _, i := range idx {
		m.rows, m.tags = append(m.rows, all[i]), append(m.tags, tags[i])
	}
	m.err = nil
	if len(errs) > 0 && len(errs) == len(m.projs) {
		m.err = errors.Join(errs...)
	}
	m.out, _ = json.Marshal(api.RunLsResponse{Runs: m.rows})
}

// load runs run.ls and keeps the live runs.
func (m *Model) load() {
	if m.runIn != nil {
		m.loadMerged()
		return
	}
	out, err := m.run(readOp, nil)
	m.out, m.err = out, err
	if err != nil {
		return
	}
	var ls api.RunLsResponse
	if err := json.Unmarshal(out, &ls); err != nil {
		m.err = err
		return
	}
	m.rows = m.rows[:0]
	for _, it := range ls.Runs {
		if it.Live {
			m.rows = append(m.rows, it)
		}
	}
}

// body is the panel's lines: one per live run, or the empty-state line.
func (m *Model) body(int) []string {
	var lines []string
	if m.err != nil {
		lines = append(lines, fmt.Sprintf("rhiza: %v", m.err))
	}
	if len(m.rows) == 0 {
		if m.err == nil {
			lines = append(lines, "No live runs")
		}
		return lines
	}
	for i, r := range m.rows {
		title := r.Title
		if title == "" {
			title = r.Change
		}
		step := r.Current
		if step == "" {
			step = "-"
		}
		if m.runIn != nil {
			title = m.tags[i] + " · " + title
		}
		lines = append(lines, fmt.Sprintf("%s · %s · %s", title, step, m.elapsed(r.StartedAt)))
	}
	return lines
}

// Result is the run.ls output, byte for byte, or the error met.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Component is the model as the shell hosts it.
func (m *Model) Component() modalshell.Component { return m }

func tick() tea.Cmd {
	return tea.Tick(Interval, func(time.Time) tea.Msg { return TickMsg{} })
}

func (m *Model) Init() tea.Cmd { return tea.Batch(m.panel.Init(), tick()) }

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if _, ok := msg.(TickMsg); ok {
		m.load()
		m.panel.Rebuild()
		return tick()
	}
	_, cmd := m.panel.Update(msg)
	return cmd
}

func (m *Model) Render() string          { return m.panel.View().Content }
func (m *Model) NaturalSize() (int, int) { return m.panel.NaturalSize() }
func (m *Model) Theme() theme.Theme      { return m.panel.Theme() }

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

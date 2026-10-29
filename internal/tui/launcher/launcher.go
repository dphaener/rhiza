// Package launcher is the Rex quick-launch layer: a picker over the rhiza
// apps that opens the chosen one in a modal. The installed rhiza.lua binds
// the launch key to a layer running `rhiza tui launch`, which runs Model.
package launcher

import (
	"context"
	"sort"
	"sync"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/client"
	"github.com/dphaener/rexuikit/surface"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
)

// App is one rhiza app the layer lists.
type App struct {
	// Name is the app's name, listed in the layer and passed to
	// `rhiza tui app`.
	Name string
	// Summary is the dimmed text after the name.
	Summary string
	// Cols and Rows are the app modal's size in cells; zero means the
	// default.
	Cols, Rows int
}

var (
	mu       sync.Mutex
	registry = map[string]App{}
)

// Register adds a to the apps the layer lists, replacing one of the same name.
func Register(a App) {
	mu.Lock()
	defer mu.Unlock()
	registry[a.Name] = a
}

// Apps returns the registered apps, sorted by name.
func Apps() []App {
	mu.Lock()
	defer mu.Unlock()
	out := make([]App, 0, len(registry))
	for _, a := range registry {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

const (
	defaultCols = 100
	defaultRows = 30
)

// Size returns the app modal's size in cells, the defaults where the app
// sets none.
func (a App) Size() (cols, rows int) {
	cols, rows = a.Cols, a.Rows
	if cols == 0 {
		cols = defaultCols
	}
	if rows == 0 {
		rows = defaultRows
	}
	return cols, rows
}

// Launcher opens the apps' modals in one Rex session.
type Launcher struct {
	Rex     surface.Rex
	Session string
	// Window is the window to open in; empty is the active one.
	Window string
	// Bin is the absolute path of rhiza, run as the modal's program.
	Bin string
	Cwd string
}

// OpenApp opens a's modal: a layer whose block runs `rhiza tui app <name>`.
func (l Launcher) OpenApp(ctx context.Context, a App) (*surface.Handle, error) {
	cols, rows := a.Size()
	return surface.New(l.Rex, l.Session).Open(ctx, surface.Spec{
		Kind:     surface.Modal,
		Argv:     []string{l.Bin, "tui", "app", a.Name},
		Cwd:      l.Cwd,
		Label:    a.Name,
		WindowID: l.Window,
		Cells:    client.Size{Columns: cols, Rows: rows},
		Focus:    true,
	})
}

// Rows are the picker's rows, one per app, in the order of apps.
func Rows(apps []App) []ui.Row {
	rows := make([]ui.Row, len(apps))
	for i, a := range apps {
		rows[i] = ui.Row{Label: a.Name, Detail: a.Summary}
	}
	return rows
}

// Model is the bubbletea program of the layer: a ui.List over the apps.
// Selecting one opens its modal and quits, which closes the layer; esc
// quits too. Err holds the error of a failed open.
type Model struct {
	list     *ui.List
	apps     []App
	launcher Launcher
	ctx      context.Context
	Err      error
}

// NewModel returns the layer's model over apps.
func NewModel(ctx context.Context, l Launcher, apps []App, p colorprofile.Profile) *Model {
	return &Model{
		list:     ui.NewList(theme.Dark(), p, Rows(apps)).WithTitle("Rhiza apps"),
		apps:     apps,
		launcher: l,
		ctx:      ctx,
	}
}

// NaturalSize is the layer's size in cells.
func (m *Model) NaturalSize() (cols, rows int) { return m.list.NaturalSize() }

func (m *Model) Init() tea.Cmd { return m.list.Init() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.ListSelectedMsg:
		if _, err := m.launcher.OpenApp(m.ctx, m.apps[msg.Index]); err != nil {
			m.Err = err
		}
		return m, tea.Quit
	case ui.ListCancelledMsg:
		return m, tea.Quit
	}
	next, cmd := m.list.Update(msg)
	m.list, _ = next.(*ui.List)
	return m, cmd
}

func (m *Model) View() tea.View { return m.list.View() }

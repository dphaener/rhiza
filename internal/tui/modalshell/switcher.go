package modalshell

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
)

// SwitchKey is the key that opens the project picker in every app.
const SwitchKey = "ctrl+o"

// Built is an app as the switcher hosts it: the model, and the result it
// reports when it closes.
type Built struct {
	Model  tea.Model
	Result func() ([]byte, error)
}

// Switching describes how an app switches project.
type Switching struct {
	// Projects are the registered projects the picker lists.
	Projects []string
	// Folder is the folder of a project.
	Folder func(project string) string
	// Chdir moves the process to a folder; os.Chdir when nil.
	Chdir func(dir string) error
	// Header is the header line for the directory the app is in.
	Header func() string
	// Build makes the app for the current directory. It is called once at
	// the start and again after each switch, so the app reloads its view
	// for the picked project.
	Build func() Built
}

// Switcher hosts an app and lets one key switch its current project. It is
// a tea.Model with the shell's NaturalSize and Theme.
type Switcher struct {
	sw     Switching
	cur    Built
	header string
	picker *ui.List
	err    error
}

// NewSwitcher returns a switcher over sw, with the app built for the
// current directory.
func NewSwitcher(sw Switching) *Switcher {
	if sw.Chdir == nil {
		sw.Chdir = os.Chdir
	}
	s := &Switcher{sw: sw}
	s.reload()
	return s
}

func (s *Switcher) reload() {
	s.cur = s.sw.Build()
	s.header = ""
	if s.sw.Header != nil {
		s.header = s.sw.Header()
	}
}

// Header is the line shown above the app.
func (s *Switcher) Header() string { return s.header }

// Picking reports whether the project picker is open.
func (s *Switcher) Picking() bool { return s.picker != nil }

// Model is the app as currently built.
func (s *Switcher) Model() tea.Model { return s.cur.Model }

// Result is the current app's result.
func (s *Switcher) Result() ([]byte, error) {
	if s.cur.Result == nil {
		return nil, nil
	}
	return s.cur.Result()
}

// Err is the error of the last switch that failed, if any.
func (s *Switcher) Err() error { return s.err }

func (s *Switcher) app() model { m, _ := s.cur.Model.(model); return m }

func (s *Switcher) Init() tea.Cmd { return s.cur.Model.Init() }

func (s *Switcher) Theme() theme.Theme {
	if m := s.app(); m != nil {
		return m.Theme()
	}
	return theme.Dark()
}

func (s *Switcher) NaturalSize() (int, int) {
	if s.picker != nil {
		return headedSize(s.header, s.picker.NaturalSize)
	}
	if m := s.app(); m != nil {
		return headedSize(s.header, m.NaturalSize)
	}
	return headedSize(s.header, func() (int, int) { return 0, 0 })
}

func (s *Switcher) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	msg = underHeader(s.header, msg)
	if s.picker != nil {
		return s.updatePicker(msg)
	}
	if k, ok := msg.(tea.KeyPressMsg); ok && k.String() == SwitchKey && len(s.sw.Projects) > 0 {
		return s, s.openPicker()
	}
	m, cmd := s.cur.Model.Update(msg)
	s.cur.Model = m
	return s, cmd
}

func (s *Switcher) openPicker() tea.Cmd {
	rows := make([]ui.Row, 0, len(s.sw.Projects))
	for _, p := range s.sw.Projects {
		rows = append(rows, ui.Row{Label: p})
	}
	s.picker = ui.NewList(s.Theme(), colorprofile.Detect(os.Stdout, os.Environ()), rows)
	return s.picker.Init()
}

func (s *Switcher) updatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case ui.ListCancelledMsg:
		s.picker = nil
		return s, nil
	case tea.KeyPressMsg:
		if m.String() == "esc" {
			s.picker = nil
			return s, nil
		}
	case ui.ListSelectedMsg:
		s.picker = nil
		if m.Index < 0 || m.Index >= len(s.sw.Projects) {
			return s, nil
		}
		if err := s.sw.Chdir(s.sw.Folder(s.sw.Projects[m.Index])); err != nil {
			s.err = err
			return s, nil
		}
		s.err = nil
		s.reload()
		return s, s.cur.Model.Init()
	}
	_, cmd := s.picker.Update(msg)
	return s, cmd
}

func (s *Switcher) View() tea.View {
	if s.picker != nil {
		return tea.NewView(withHeader(s.header, s.picker.View().Content))
	}
	v := s.cur.Model.View()
	v.Content = withHeader(s.header, v.Content)
	return v
}

// Package modalshell is the one shell every Rhiza app runs in: it hosts a
// rexuikit ui component (Form, List, Dialog or InfoPanel), themed from the
// terminal's background by rexuikit's ui/theme, and wraps it in rexuikit's
// surface/program for focus. It closes the app on Esc, on the component's
// cancel and when the app reports its job done, and hands back the result
// the app set.
package modalshell

import (
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rexuikit/surface"
	"github.com/dphaener/rexuikit/surface/program"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
)

// DoneMsg is what an app sends, as a command, when its job is done. The
// shell quits on it and keeps Result as the app's result.
type DoneMsg struct {
	Result []byte
	Err    error
}

// Done is the command that sends a DoneMsg.
func Done(result []byte, err error) tea.Cmd {
	return func() tea.Msg { return DoneMsg{Result: result, Err: err} }
}

// Component is a rexuikit ui component as the shell drives it.
type Component interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	// Render is the component's own render.
	Render() string
	// NaturalSize is the component's size in cells.
	NaturalSize() (cols, rows int)
	Theme() theme.Theme
}

// Shell hosts one Component. It is a tea.Model.
type Shell struct {
	c      Component
	header string
	result []byte
	err    error
	done   bool
}

// New returns a shell hosting c.
func New(c Component) *Shell { return &Shell{c: c} }

// WithHeader shows line above the hosted component, one row of its own.
func (s *Shell) WithHeader(line string) *Shell {
	s.header = line
	return s
}

// Header is the line shown above the component, empty for none.
func (s *Shell) Header() string { return s.header }

// Component is the hosted component.
func (s *Shell) Component() Component { return s.c }

// Theme is the hosted component's current theme.
func (s *Shell) Theme() theme.Theme { return s.c.Theme() }

// NaturalSize is the hosted component's size in cells.
func (s *Shell) NaturalSize() (cols, rows int) {
	return headedSize(s.header, s.c.NaturalSize)
}

// Result is what the app set when it reported done, or the error it did.
func (s *Shell) Result() ([]byte, error) { return s.result, s.err }

// Closed reports whether the shell has asked to quit.
func (s *Shell) Closed() bool { return s.done }

func (s *Shell) Init() tea.Cmd { return s.c.Init() }

// Update closes on Esc, ui.FormCancelledMsg, ui.ListCancelledMsg,
// ui.InfoClosedMsg and DoneMsg; every other message goes to the component.
func (s *Shell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case DoneMsg:
		s.result, s.err = m.Result, m.Err
		return s.quit()
	case ui.FormCancelledMsg, ui.ListCancelledMsg, ui.InfoClosedMsg:
		return s.quit()
	case tea.KeyPressMsg:
		if m.String() == "esc" {
			return s.quit()
		}
	}
	return s, s.c.Update(msg)
}

func (s *Shell) quit() (tea.Model, tea.Cmd) {
	s.done = true
	return s, tea.Quit
}

func (s *Shell) View() tea.View { return tea.NewView(withHeader(s.header, s.c.Render())) }

// withHeader puts line, when there is one, above content.
func withHeader(line, content string) string {
	if line == "" {
		return content
	}
	return line + "\n" + content
}

// headedSize is size with room for the header line.
func headedSize(line string, size func() (int, int)) (cols, rows int) {
	cols, rows = size()
	if line == "" {
		return cols, rows
	}
	return max(cols, utf8.RuneCountInString(line)), rows + 1
}

// headed is a model with a header line above its view.
type headed struct {
	model
	line string
}

func (h headed) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := h.model.Update(msg)
	h.model = m.(model)
	return h, cmd
}

func (h headed) View() tea.View {
	v := h.model.View()
	v.Content = withHeader(h.line, v.Content)
	return v
}

func (h headed) NaturalSize() (int, int) { return headedSize(h.line, h.model.NaturalSize) }

// Program is the shell wrapped for focus: what tea.NewProgram runs. A
// Modal stays open on blur.
func (s *Shell) Program() program.Model {
	return program.Wrap(s, program.Options{Blur: surface.StayOpen})
}

// Run runs s in a bubbletea program and returns its result.
func (s *Shell) Run(opts ...tea.ProgramOption) ([]byte, error) {
	p := program.Wrap(sized(s), program.Options{Blur: surface.StayOpen})
	if _, err := tea.NewProgram(p, opts...).Run(); err != nil {
		return nil, err
	}
	return s.Result()
}

type model interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
	NaturalSize() (int, int)
	Theme() theme.Theme
}

type modelComponent struct{ m model }

func (c modelComponent) Init() tea.Cmd { return c.m.Init() }
func (c modelComponent) Update(msg tea.Msg) tea.Cmd {
	_, cmd := c.m.Update(msg)
	return cmd
}
func (c modelComponent) Render() string          { return c.m.View().Content }
func (c modelComponent) NaturalSize() (int, int) { return c.m.NaturalSize() }
func (c modelComponent) Theme() theme.Theme      { return c.m.Theme() }

// Form hosts a ui.Form.
func Form(f *ui.Form) Component { return formComponent{f} }

type formComponent struct{ *ui.Form }

func (c formComponent) Update(msg tea.Msg) tea.Cmd {
	_, cmd := c.Form.Update(msg)
	return cmd
}

// List hosts a ui.List.
func List(l *ui.List) Component { return modelComponent{l} }

// Info hosts a ui.InfoPanel.
func Info(p *ui.InfoPanel) Component { return modelComponent{p} }

// Dialog hosts a ui.Dialog.
func Dialog(d *ui.Dialog) Component { return dialogComponent{d} }

type dialogComponent struct{ *ui.Dialog }

func (c dialogComponent) Render() string { return c.Dialog.View() }

// Run runs a model that is not yet a hosted component (an app awaiting its
// move to one) through the same surface/program wrapper, so every app has
// the same focus handling. Like Shell.Run, it reports the model's
// NaturalSize, when it has one, to the file ReportSizeTo named.
func Run(m tea.Model, opts ...tea.ProgramOption) error {
	return RunHeaded(m, "", opts...)
}

// RunHeaded is Run with line shown above the model's view. The model must
// have NaturalSize and Theme, as the shell's models do, when line is set.
func RunHeaded(m tea.Model, line string, opts ...tea.ProgramOption) error {
	if hm, ok := m.(model); ok && line != "" {
		m = headed{model: hm, line: line}
	}
	_, err := tea.NewProgram(program.Wrap(sized(m), program.Options{Blur: surface.StayOpen}), opts...).Run()
	return err
}

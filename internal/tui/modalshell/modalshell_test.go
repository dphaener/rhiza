package modalshell

import (
	"errors"
	"image/color"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/surface/program"
	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"
)

var prof = colorprofile.TrueColor

func hosts() map[string]func() Component {
	t := theme.Dark()
	return map[string]func() Component{
		"form": func() Component {
			return Form(ui.NewForm(t, prof, ui.NewText("name", "Name")))
		},
		"list": func() Component {
			return List(ui.NewList(t, prof, []ui.Row{{Label: "one"}, {Label: "two"}}))
		},
		"dialog": func() Component { return Dialog(ui.NewDialog(t, prof, "Title", []string{"body line"})) },
		"info": func() Component {
			return Info(ui.NewInfoPanel(t, prof, "Info", func(int) []string { return []string{"info line"} }))
		},
	}
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestCloses(t *testing.T) {
	esc := tea.KeyPressMsg{Code: tea.KeyEscape}
	for _, tt := range []struct {
		name string
		msg  tea.Msg
		quit bool
	}{
		{"esc", esc, true},
		{"form cancelled", ui.FormCancelledMsg{}, true},
		{"list cancelled", ui.ListCancelledMsg{}, true},
		{"done", DoneMsg{}, true},
		{"other key", tea.KeyPressMsg{Code: 'a', Text: "a"}, false},
		{"window size", tea.WindowSizeMsg{Width: 40, Height: 10}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := New(hosts()["dialog"]())
			_, cmd := s.Update(tt.msg)
			if got := isQuit(cmd); got != tt.quit {
				t.Errorf("quit = %v, want %v", got, tt.quit)
			}
			if s.Closed() != tt.quit {
				t.Errorf("Closed = %v, want %v", s.Closed(), tt.quit)
			}
		})
	}
}

func TestResult(t *testing.T) {
	s := New(hosts()["dialog"]())
	s.Update(DoneMsg{Result: []byte("out")})
	if got, err := s.Result(); string(got) != "out" || err != nil {
		t.Errorf("Result = %q, %v", got, err)
	}
	boom := errors.New("boom")
	s = New(hosts()["dialog"]())
	s.Update(DoneMsg{Err: boom})
	if _, err := s.Result(); err != boom {
		t.Errorf("err = %v", err)
	}
	// Cancel keeps no result.
	s = New(hosts()["form"]())
	s.Update(ui.FormCancelledMsg{})
	if got, err := s.Result(); got != nil || err != nil {
		t.Errorf("cancel result = %q, %v", got, err)
	}
}

func TestRendersComponent(t *testing.T) {
	want := map[string]string{"form": "Name", "list": "one", "dialog": "body line", "info": "info line"}
	for name, mk := range hosts() {
		t.Run(name, func(t *testing.T) {
			s := New(mk())
			if got := s.View().Content; !strings.Contains(got, want[name]) {
				t.Errorf("view lacks %q:\n%s", want[name], got)
			}
			if cols, rows := s.NaturalSize(); cols <= 0 || rows <= 0 {
				t.Errorf("natural size %dx%d", cols, rows)
			}
		})
	}
}

func TestThemeFromBackground(t *testing.T) {
	black, white := color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}
	for name, mk := range hosts() {
		t.Run(name, func(t *testing.T) {
			s := New(mk())
			s.Update(tea.BackgroundColorMsg{Color: white})
			if s.Theme() != theme.Light() {
				t.Errorf("light background: theme is not theme.Light")
			}
			s.Update(tea.BackgroundColorMsg{Color: black})
			if s.Theme() != theme.Dark() {
				t.Errorf("dark background: theme is not theme.Dark")
			}
		})
	}
}

func TestProgramIsSurfaceModel(t *testing.T) {
	s := New(hosts()["list"]())
	var m tea.Model = s.Program()
	if _, ok := m.(program.Model); !ok {
		t.Fatalf("program model is %T", m)
	}
	if program.Unwrap(m) != tea.Model(s) {
		t.Errorf("Unwrap is not the shell")
	}
	if !m.View().ReportFocus {
		t.Errorf("focus reporting is off")
	}
}

// Every hosted component closes the same way, whatever it is.
func TestClosesEveryHost(t *testing.T) {
	for name, mk := range hosts() {
		t.Run(name, func(t *testing.T) {
			s := New(mk())
			if _, cmd := s.Update(tea.KeyPressMsg{Code: 'a', Text: "a"}); isQuit(cmd) || s.Closed() {
				t.Fatal("a plain key closed the shell")
			}
			if _, cmd := s.Update(tea.KeyPressMsg{Code: tea.KeyEscape}); !isQuit(cmd) || !s.Closed() {
				t.Fatal("esc did not close the shell")
			}
		})
	}
}

func TestShellHeader(t *testing.T) {
	for name, mk := range hosts() {
		plain := New(mk())
		s := New(mk()).WithHeader("Project: alpha")
		if got := s.View().Content; !strings.HasPrefix(got, "Project: alpha\n") {
			t.Errorf("%s: view %q lacks the header", name, got)
		}
		if strings.Contains(plain.View().Content, "Project:") {
			t.Errorf("%s: a shell without a header shows one", name)
		}
		c0, r0 := plain.NaturalSize()
		if c, r := s.NaturalSize(); r != r0+1 || c < c0 {
			t.Errorf("%s: size %dx%d, want one more row than %dx%d", name, c, r, c0, r0)
		}
	}
}

func TestRunHeadedWrapsModel(t *testing.T) {
	m := headed{model: New(hosts()["list"]()), line: "No project"}
	if got := m.View().Content; !strings.HasPrefix(got, "No project\n") {
		t.Errorf("view %q lacks the header", got)
	}
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if _, ok := next.(headed); !ok {
		t.Errorf("update lost the header: %T", next)
	}
}

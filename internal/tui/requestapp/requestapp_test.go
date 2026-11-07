package requestapp

import (
	"encoding/json"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

type fake struct{ calls [][]string }

func (f *fake) run(op string, args []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{op}, args...))
	return json.Marshal(api.RequestResponse{Request: api.RequestInfo{ID: "R1", Kind: "bug", State: "queued"}})
}

func newShell(f *fake) *modalshell.Shell {
	return modalshell.New(New(f.run, colorprofile.TrueColor))
}

func press(s *modalshell.Shell, keys ...tea.KeyPressMsg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	for _, k := range keys {
		_, cmd = s.Update(k)
		// Run the command and feed its message back, as a program does.
		for cmd != nil {
			msg := cmd()
			if _, ok := msg.(tea.QuitMsg); ok {
				return nil, true
			}
			if b, ok := msg.(tea.BatchMsg); ok {
				for _, c := range b {
					if c != nil {
						if m := c(); m != nil {
							if _, q := m.(tea.QuitMsg); q {
								return nil, true
							}
							_, next := s.Update(m)
							if next != nil {
								if _, q := next().(tea.QuitMsg); q {
									return nil, true
								}
							}
						}
					}
				}
				break
			}
			_, cmd = s.Update(msg)
		}
	}
	return nil, s.Closed()
}

func typed(text string) []tea.KeyPressMsg {
	var ks []tea.KeyPressMsg
	for _, r := range text {
		ks = append(ks, tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return ks
}

var (
	enter = tea.KeyPressMsg{Code: tea.KeyEnter}
	down  = tea.KeyPressMsg{Code: tea.KeyDown}
	esc   = tea.KeyPressMsg{Code: tea.KeyEscape}
)

func TestSubmitQueuesAndCloses(t *testing.T) {
	f := &fake{}
	s := newShell(f)
	s.Init()
	// Kind: down once picks bug; enter moves to the text.
	keys := []tea.KeyPressMsg{down, enter}
	keys = append(keys, typed("login  fails")...)
	press(s, keys...)
	if _, closed := press(s, enter); !closed {
		t.Fatal("app did not quit after submit")
	}
	if len(f.calls) != 1 {
		t.Fatalf("calls = %q", f.calls)
	}
	want := "request --kind bug -- login fails"
	if got := strings.Join(f.calls[0], " "); got != want {
		t.Fatalf("call = %q, want %q", got, want)
	}
}

func TestEmptyTextRefused(t *testing.T) {
	f := &fake{}
	s := newShell(f)
	s.Init()
	press(s, enter)
	_, closed := press(s, enter)
	if closed || len(f.calls) != 0 {
		t.Fatalf("closed=%v calls=%q", closed, f.calls)
	}
	_, closed = press(s, typed("   ")...)
	if _, closed = press(s, enter); closed || len(f.calls) != 0 {
		t.Fatalf("blank text submitted: closed=%v calls=%q", closed, f.calls)
	}
	if !strings.Contains(s.Component().Render(), "A request needs its text") {
		t.Fatal("no validation message shown")
	}
}

func TestEscCallsNothing(t *testing.T) {
	f := &fake{}
	s := newShell(f)
	s.Init()
	press(s, typed("x")...)
	if _, closed := press(s, esc); !closed {
		t.Fatal("esc did not close")
	}
	if len(f.calls) != 0 {
		t.Fatalf("calls = %q", f.calls)
	}
}

func TestRegistered(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == "request" {
			if a.Component == nil {
				t.Fatal("no component")
			}
			if cols, rows := a.Size(); cols <= 0 || rows <= 0 {
				t.Fatalf("size %dx%d", cols, rows)
			}
			return
		}
	}
	t.Fatal("request not registered")
}

type pfake struct{ calls [][]string }

func (f *pfake) run(project, op string, args []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{project, op}, args...))
	return json.Marshal(api.RequestResponse{Request: api.RequestInfo{ID: "R1", Kind: "bug", State: "queued"}})
}

func TestRequestQueuedInCurrentProject(t *testing.T) {
	f := &pfake{}
	s := modalshell.New(NewIn(f.run, colorprofile.TrueColor, "B", []string{"A", "B"}))
	s.Init()
	press(s, enter)
	press(s, typed("fix it")...)
	if _, closed := press(s, enter); !closed {
		t.Fatal("app did not quit after submit")
	}
	want := "B request --kind feature -- fix it"
	if len(f.calls) != 1 || strings.Join(f.calls[0], " ") != want {
		t.Fatalf("calls = %q, want %q", f.calls, want)
	}
}

func TestPickerAtRootThenQueuesInPicked(t *testing.T) {
	f := &pfake{}
	s := modalshell.New(NewIn(f.run, colorprofile.TrueColor, "", []string{"A", "B"}))
	s.Init()
	if !strings.Contains(s.Component().Render(), "Project") {
		t.Fatal("no project picker shown")
	}
	// Enter takes the first project, A; the request form follows.
	if _, closed := press(s, enter); closed {
		t.Fatal("picker closed the app")
	}
	if len(f.calls) != 0 {
		t.Fatalf("picking queued something: %q", f.calls)
	}
	press(s, enter)
	press(s, typed("hello")...)
	if _, closed := press(s, enter); !closed {
		t.Fatal("app did not quit after submit")
	}
	want := "A request --kind feature -- hello"
	if len(f.calls) != 1 || strings.Join(f.calls[0], " ") != want {
		t.Fatalf("calls = %q, want %q", f.calls, want)
	}
}

func TestPickerEscStartsNothing(t *testing.T) {
	f := &pfake{}
	s := modalshell.New(NewIn(f.run, colorprofile.TrueColor, "", []string{"A", "B"}))
	s.Init()
	if _, closed := press(s, esc); !closed {
		t.Fatal("esc did not close")
	}
	if len(f.calls) != 0 {
		t.Fatalf("calls = %q", f.calls)
	}
}

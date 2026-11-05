package daemonapp

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

type fake struct {
	paused string
	calls  [][]string
}

func (f *fake) run(op string, args []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{op}, args...))
	if op == "daemon.settings" {
		return []byte(fmt.Sprintf(`{"settings":[{"key":"rhiza.daemon.paused","value":%q,"default":"false","source":"default"},{"key":"rhiza.daemon.maxRuns","value":"3","default":"4","source":"git config"}]}`, f.paused)), nil
	}
	return []byte(`{"ok":true}`), nil
}

func (f *fake) sets() [][]string {
	var out [][]string
	for _, c := range f.calls {
		if c[0] == "daemon.set" {
			out = append(out, c)
		}
	}
	return out
}

// drive sends keys through the shell, running returned commands until the
// shell quits.
func drive(t *testing.T, s *modalshell.Shell, keys ...string) {
	t.Helper()
	pump := func(msg tea.Msg) {
		queue := []tea.Msg{msg}
		for len(queue) > 0 && !s.Closed() {
			m := queue[0]
			queue = queue[1:]
			_, cmd := s.Update(m)
			if cmd == nil {
				continue
			}
			if r := cmd(); r != nil {
				if b, ok := r.(tea.BatchMsg); ok {
					for _, c := range b {
						if c != nil {
							if rm := c(); rm != nil {
								queue = append(queue, rm)
							}
						}
					}
				} else if _, quit := r.(tea.QuitMsg); !quit {
					queue = append(queue, r)
				}
			}
		}
	}
	for _, k := range keys {
		pump(key(k))
	}
}

func key(k string) tea.KeyPressMsg {
	switch k {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "backspace":
		return tea.KeyPressMsg{Code: tea.KeyBackspace}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	}
	return tea.KeyPressMsg{Code: rune(k[0]), Text: k}
}

func open(t *testing.T, paused string) (*fake, *modalshell.Shell) {
	t.Helper()
	f := &fake{paused: paused}
	a, err := New(f.run, theme.Dark(), colorprofile.TrueColor)
	if err != nil {
		t.Fatal(err)
	}
	return f, modalshell.New(a)
}

func TestPrefilledFromSettings(t *testing.T) {
	for paused, want := range map[string]string{"false": "No", "true": "Yes"} {
		_, s := open(t, paused)
		if got := s.Component().(*App).Paused(); got != (paused == "true") {
			t.Errorf("paused %s: Paused() = %v", paused, got)
		}
		if !strings.Contains(s.Component().Render(), want) {
			t.Errorf("paused %s: render lacks %q", paused, want)
		}
	}
}

func TestFlipAndSubmitSetsAndQuits(t *testing.T) {
	for paused, want := range map[string]string{"false": "true", "true": "false"} {
		f, s := open(t, paused)
		flip := "y"
		if paused == "true" {
			flip = "n"
		}
		drive(t, s, flip, "enter", "enter")
		sets := f.sets()
		if len(sets) != 1 || strings.Join(sets[0], " ") != "daemon.set -- rhiza.daemon.paused "+want {
			t.Errorf("paused %s: sets = %v", paused, sets)
		}
		if !s.Closed() {
			t.Errorf("paused %s: app did not quit", paused)
		}
		if out, err := s.Result(); err != nil || string(out) != `{"ok":true}` {
			t.Errorf("paused %s: result %q, %v", paused, out, err)
		}
	}
}

func TestUnchangedSubmitSetsNothing(t *testing.T) {
	f, s := open(t, "false")
	drive(t, s, "enter", "enter")
	if len(f.sets()) != 0 {
		t.Errorf("sets = %v", f.sets())
	}
	if !s.Closed() {
		t.Error("app did not quit")
	}
}

func TestEscSetsNothing(t *testing.T) {
	f, s := open(t, "false")
	drive(t, s, "y", "esc")
	if len(f.sets()) != 0 {
		t.Errorf("sets = %v", f.sets())
	}
	if !s.Closed() {
		t.Error("app did not quit")
	}
}

func TestRunLimitStartsFromSettings(t *testing.T) {
	_, s := open(t, "false")
	if got := s.Component().(*App).MaxRuns(); got != "3" {
		t.Errorf("MaxRuns() = %q", got)
	}
}

// typeLimit moves to the run-limit field, clears it and types v.
func typeLimit(t *testing.T, s *modalshell.Shell, v string) {
	t.Helper()
	keys := []string{"enter", "backspace"}
	for _, r := range v {
		keys = append(keys, string(r))
	}
	drive(t, s, keys...)
}

func TestValidLimitSetsAndQuits(t *testing.T) {
	f, s := open(t, "false")
	typeLimit(t, s, "8")
	drive(t, s, "enter")
	sets := f.sets()
	if len(sets) != 1 || strings.Join(sets[0], " ") != "daemon.set -- rhiza.daemon.maxRuns 8" {
		t.Errorf("sets = %v", sets)
	}
	if !s.Closed() {
		t.Error("app did not quit")
	}
}

func TestInvalidLimitCallsNothing(t *testing.T) {
	for _, v := range []string{"0", "x", ""} {
		f, s := open(t, "false")
		typeLimit(t, s, v)
		drive(t, s, "enter")
		if len(f.sets()) != 0 {
			t.Errorf("limit %q: sets = %v", v, f.sets())
		}
		if s.Closed() {
			t.Errorf("limit %q: app quit", v)
		}
	}
}

package autonomyapp

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// fake is a runner over in-memory settings.
type fake struct {
	vals  map[string]string
	order []string
	calls [][]string
	err   error
}

func newFake() *fake {
	f := &fake{vals: map[string]string{
		"rhiza.daemon.maxRuns":   "2",
		"rhiza.autonomy.plan":    "owner",
		"rhiza.autonomy.accept":  "auto",
		"rhiza.autonomy.ship":    "auto",
		"rhiza.autonomy.signoff": "owner",
		"rhiza.daemon.paused":    "false",
	}}
	f.order = []string{"rhiza.daemon.paused", "rhiza.daemon.maxRuns", "rhiza.autonomy.plan",
		"rhiza.autonomy.accept", "rhiza.autonomy.ship", "rhiza.autonomy.signoff"}
	return f
}

func (f *fake) run(op string, args []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{op}, args...))
	switch op {
	case "daemon.settings":
		var rows []api.Setting
		for _, k := range f.order {
			rows = append(rows, api.Setting{Key: k, Value: f.vals[k], Source: "default"})
		}
		return json.Marshal(api.DaemonSettingsResponse{Settings: rows})
	case "daemon.set":
		if f.err != nil {
			return nil, f.err
		}
		if len(args) != 3 || args[0] != "--" {
			return nil, fmt.Errorf("unexpected set args %q", args)
		}
		f.vals[args[1]] = args[2]
		return []byte("{}"), nil
	}
	return nil, nil
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

// send delivers msg to the shell and then every message its commands make.
func send(sh *modalshell.Shell, msg tea.Msg) {
	_, cmd := sh.Update(msg)
	drain(sh, cmd)
}

func drain(sh *modalshell.Shell, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	switch m := cmd().(type) {
	case nil:
	case tea.BatchMsg:
		for _, c := range m {
			drain(sh, c)
		}
	default:
		if _, ok := m.(tea.QuitMsg); ok {
			return
		}
		send(sh, m)
	}
}

func key(sh *modalshell.Shell, s string) {
	switch s {
	case "enter":
		send(sh, tea.KeyPressMsg{Code: tea.KeyEnter})
	case "esc":
		send(sh, tea.KeyPressMsg{Code: tea.KeyEscape})
	case "down":
		send(sh, tea.KeyPressMsg{Code: tea.KeyDown})
	case "up":
		send(sh, tea.KeyPressMsg{Code: tea.KeyUp})
	}
}

func start(f *fake) (*Model, *modalshell.Shell) {
	m := NewModel(f.run)
	sh := m.Shell()
	sh.Init()
	return m, sh
}

func TestRegistered(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == Name {
			if a.Key != "y" {
				t.Errorf("key = %q", a.Key)
			}
			if c, r := a.Size(); c <= 0 || r <= 0 {
				t.Errorf("size = %dx%d", c, r)
			}
			return
		}
	}
	t.Fatal("autonomy not registered")
}

func TestStartsAtDaemonSettings(t *testing.T) {
	f := newFake()
	m, _ := start(f)
	want := map[string]string{
		"rhiza.autonomy.plan": "owner", "rhiza.autonomy.accept": "auto",
		"rhiza.autonomy.ship": "auto", "rhiza.autonomy.signoff": "owner",
	}
	for k, v := range want {
		if got := m.Values()[k]; got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	if len(m.Values()) != 4 {
		t.Errorf("fields = %v", m.Values())
	}
	for _, label := range []string{"Plan", "Accept", "Ship", "Sign-off"} {
		if !strings.Contains(m.Render(), label) {
			t.Errorf("render lacks %q", label)
		}
	}
}

func TestSubmitWritesChangedOnly(t *testing.T) {
	f := newFake()
	_, sh := start(f)
	// plan: owner -> auto (up), accept unchanged, ship: auto -> owner (down),
	// sign-off unchanged; enter on the last field submits.
	key(sh, "up")
	key(sh, "enter")
	key(sh, "enter")
	key(sh, "down")
	key(sh, "enter")
	key(sh, "enter")
	sets := f.sets()
	if len(sets) != 2 {
		t.Fatalf("sets = %q", sets)
	}
	if f.vals["rhiza.autonomy.plan"] != "auto" || f.vals["rhiza.autonomy.ship"] != "owner" {
		t.Errorf("vals = %v", f.vals)
	}
	if !sh.Closed() {
		t.Error("shell did not close after submit")
	}
	if _, err := sh.Result(); err != nil {
		t.Errorf("err = %v", err)
	}
}

func TestSubmitUnchangedWritesNothing(t *testing.T) {
	f := newFake()
	_, sh := start(f)
	for range 4 {
		key(sh, "enter")
	}
	if n := len(f.sets()); n != 0 {
		t.Errorf("sets = %d", n)
	}
	if !sh.Closed() {
		t.Error("shell did not close")
	}
}

func TestEscCallsNothing(t *testing.T) {
	f := newFake()
	_, sh := start(f)
	key(sh, "up")
	key(sh, "esc")
	if n := len(f.sets()); n != 0 {
		t.Errorf("sets = %d", n)
	}
	if !sh.Closed() {
		t.Error("shell did not close")
	}
}

func TestSetErrorIsResult(t *testing.T) {
	f := newFake()
	_, sh := start(f)
	key(sh, "up")
	f.err = fmt.Errorf("refused")
	for range 4 {
		key(sh, "enter")
	}
	if _, err := sh.Result(); err == nil || !strings.Contains(err.Error(), "refused") {
		t.Errorf("err = %v", err)
	}
}

func TestResultIsDaemonSettings(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	want, _ := f.run("daemon.settings", nil)
	got, err := m.Result()
	if err != nil || string(got) != string(want) {
		t.Errorf("Result = %q, %v", got, err)
	}
}

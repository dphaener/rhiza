package autonomyapp

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// fake is a runner over in-memory settings and a daemon state.
type fake struct {
	vals    map[string]string
	order   []string
	running bool
	process string
	calls   [][]string
	stopErr error
}

func newFake() *fake {
	f := &fake{vals: map[string]string{}, process: "host1:42"}
	for k, v := range map[string]string{
		"rhiza.daemon.paused":   "false",
		"rhiza.daemon.maxRuns":  "2",
		"rhiza.autonomy.plan":   "owner",
		"rhiza.autonomy.accept": "auto",
	} {
		f.vals[k] = v
	}
	f.order = []string{"rhiza.daemon.paused", "rhiza.daemon.maxRuns", "rhiza.autonomy.plan", "rhiza.autonomy.accept"}
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
	case "daemon.status":
		return json.Marshal(api.DaemonStatusResponse{Running: f.running, Process: f.process})
	case "daemon.start":
		return []byte("{}"), nil
	case "daemon.stop":
		return []byte("{}"), f.stopErr
	case "daemon.set":
		if len(args) != 3 || args[0] != "--" {
			return nil, fmt.Errorf("unexpected set args %q", args)
		}
		if args[1] == "rhiza.daemon.maxRuns" && args[2] == "0" {
			return nil, fmt.Errorf("must be a whole number of at least 1")
		}
		f.vals[args[1]] = args[2]
		return []byte("{}"), nil
	}
	return nil, nil
}

func key(m *Model, s string) {
	var k tea.KeyPressMsg
	switch s {
	case "enter":
		k = tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		k = tea.KeyPressMsg{Code: tea.KeyEscape}
	case "backspace":
		k = tea.KeyPressMsg{Code: tea.KeyBackspace}
	default:
		k = tea.KeyPressMsg{Code: rune(s[0]), Text: s}
	}
	m.Update(k)
}

func moveTo(m *Model, k string) {
	for i := 0; i < 10; i++ {
		if s, _ := m.Selected(); s.Key == k {
			return
		}
		key(m, "j")
	}
}

func TestRegistered(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == Name {
			if a.Key != "y" {
				t.Errorf("key = %q", a.Key)
			}
			return
		}
	}
	t.Fatal("autonomy not registered")
}

func TestToggleAutonomyPlan(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	moveTo(m, "rhiza.autonomy.plan")
	key(m, "enter")
	if f.vals["rhiza.autonomy.plan"] != "auto" {
		t.Fatalf("plan = %q", f.vals["rhiza.autonomy.plan"])
	}
	if !slices.ContainsFunc(f.calls, func(c []string) bool {
		return slices.Equal(c, []string{"daemon.set", "--", "rhiza.autonomy.plan", "auto"})
	}) {
		t.Errorf("calls = %q", f.calls)
	}
	key(m, "enter")
	if f.vals["rhiza.autonomy.plan"] != "owner" {
		t.Errorf("plan = %q", f.vals["rhiza.autonomy.plan"])
	}
}

func TestTogglePaused(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	key(m, "enter")
	if f.vals["rhiza.daemon.paused"] != "true" {
		t.Fatalf("paused = %q", f.vals["rhiza.daemon.paused"])
	}
}

func TestEditMaxRuns(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	moveTo(m, "rhiza.daemon.maxRuns")
	key(m, "enter")
	if !m.Editing() {
		t.Fatal("entry not open")
	}
	key(m, "backspace")
	key(m, "0")
	key(m, "enter")
	if !m.Editing() || !strings.Contains(m.Notice(), "at least 1") {
		t.Fatalf("a refused value must keep the entry open, notice %q", m.Notice())
	}
	key(m, "backspace")
	key(m, "5")
	key(m, "enter")
	if m.Editing() || f.vals["rhiza.daemon.maxRuns"] != "5" {
		t.Errorf("maxRuns = %q, editing %v", f.vals["rhiza.daemon.maxRuns"], m.Editing())
	}
}

// ops are the operations the fake runner was asked for, in order.
func (f *fake) ops() []string {
	var out []string
	for _, c := range f.calls {
		if c[0] == "daemon.start" || c[0] == "daemon.stop" {
			out = append(out, c[0])
		}
	}
	return out
}

func TestStartDaemonWhenNotRunning(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	key(m, "d")
	if m.Confirming() != "start" {
		t.Fatalf("confirming = %q", m.Confirming())
	}
	key(m, "y")
	if got := f.ops(); !slices.Equal(got, []string{"daemon.start"}) {
		t.Errorf("ops = %q", got)
	}
	if !strings.HasPrefix(m.notice, "Starting the daemon") {
		t.Errorf("notice = %q", m.notice)
	}
	f.running = true
	m.Update(TickMsg{})
	if !strings.HasPrefix(m.notice, "The daemon is running") {
		t.Errorf("notice = %q", m.notice)
	}
}

func TestStartReportsDaemonThatNeverRuns(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	key(m, "d")
	key(m, "y") // the first reload after the start
	for i := 1; i < StartChecks; i++ {
		if strings.Contains(m.notice, "did not start") {
			t.Fatalf("failure reported after %d reloads", i)
		}
		m.Update(TickMsg{})
	}
	if !strings.HasPrefix(m.notice, "The daemon did not start") {
		t.Errorf("notice = %q", m.notice)
	}
}

func TestDeclinedStartRunsNothing(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	key(m, "d")
	key(m, "n")
	if len(f.ops()) != 0 || m.Confirming() != "" {
		t.Errorf("ops = %q", f.ops())
	}
}

func TestStopRunningDaemon(t *testing.T) {
	f := newFake()
	f.running = true
	m := NewModel(f.run)
	key(m, "d")
	if m.Confirming() != "stop" {
		t.Fatalf("confirming = %q", m.Confirming())
	}
	key(m, "y")
	if got := f.ops(); !slices.Equal(got, []string{"daemon.stop"}) {
		t.Errorf("ops = %q", got)
	}
}

func TestStopRefusalIsShown(t *testing.T) {
	f := newFake()
	f.running = true
	f.stopErr = fmt.Errorf("the daemon runs on host1, not this machine")
	m := NewModel(f.run)
	key(m, "d")
	key(m, "y")
	if !strings.Contains(m.Notice(), "host1") {
		t.Errorf("notice %q", m.Notice())
	}
}

func TestResultIsSettingsOutput(t *testing.T) {
	f := newFake()
	m := NewModel(f.run)
	want, _ := f.run("daemon.settings", nil)
	got, err := m.Result()
	if err != nil || string(got) != string(want) {
		t.Errorf("got %s, %v", got, err)
	}
}

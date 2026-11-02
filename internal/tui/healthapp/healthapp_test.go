package healthapp

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// report is a health report with a part in each state.
var report = api.HealthResponse{State: api.HealthBad, Parts: []api.HealthPart{
	{Part: api.HealthDaemon, State: api.HealthWarn, Summary: "not running", Detail: []string{}, Next: "rhiza daemon --autonomous"},
	{Part: api.HealthDisk, State: api.HealthBad, Summary: "100 B free", Detail: []string{}},
	{Part: api.HealthInbox, State: api.HealthWarn, Summary: "1 waiting: 1 failed", Detail: []string{}},
	{Part: api.HealthCheck, State: api.HealthOK, Summary: "ok", Detail: []string{"all good"}},
}}

// fakeRunner serves health from r, counting calls.
func fakeRunner(t *testing.T, r *api.HealthResponse, calls *int) Runner {
	return func(op string, args []string) ([]byte, error) {
		if op != "health" || len(args) != 0 {
			t.Fatalf("ran %s %v", op, args)
		}
		*calls++
		return canon.Marshal(*r)
	}
}

// child is a stand-in app.
type child struct{ name string }

func (c child) Init() tea.Cmd                       { return nil }
func (c child) Update(tea.Msg) (tea.Model, tea.Cmd) { return c, nil }
func (c child) View() tea.View                      { return tea.NewView("the " + c.name + " app") }

func key(m *Model, s string) {
	switch s {
	case "enter":
		m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	default:
		m.Update(tea.KeyPressMsg{Code: rune(s[0]), Text: s})
	}
}

func TestRendersEachPartsState(t *testing.T) {
	var calls int
	r := report
	m := NewModel(fakeRunner(t, &r, &calls), nil)
	v := m.View().Content
	for _, want := range []string{
		"Health: bad",
		"> [warn] daemon   not running",
		"next: rhiza daemon --autonomous",
		"  [bad]  disk     100 B free",
		"  [warn] inbox    1 waiting: 1 failed",
		"  [ok]   check    ok",
	} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
	// Only the cursor's part shows its detail.
	if strings.Contains(v, "all good") {
		t.Errorf("view shows detail off the cursor:\n%s", v)
	}
	key(m, "j")
	key(m, "j")
	key(m, "j")
	if v := m.View().Content; !strings.Contains(v, "> [ok]   check    ok") || !strings.Contains(v, "all good") {
		t.Errorf("cursor on check:\n%s", v)
	}
	out, err := m.Result()
	want, _ := canon.Marshal(report)
	if err != nil || string(out) != string(want) {
		t.Fatalf("result %s, %v", out, err)
	}
}

func TestReloadOnR(t *testing.T) {
	var calls int
	r := report
	m := NewModel(fakeRunner(t, &r, &calls), nil)
	r = api.HealthResponse{State: api.HealthOK, Parts: []api.HealthPart{{Part: api.HealthDaemon, State: api.HealthOK, Summary: "running", Detail: []string{}}}}
	key(m, "r")
	if calls != 2 {
		t.Fatalf("health ran %d times", calls)
	}
	if v := m.View().Content; !strings.Contains(v, "> [ok]   daemon   running") {
		t.Errorf("after r:\n%s", v)
	}
}

func TestEnterOnInboxOpensInboxApp(t *testing.T) {
	var calls int
	r := report
	var opened []string
	m := NewModel(fakeRunner(t, &r, &calls), func(app string) (tea.Model, error) {
		opened = append(opened, app)
		return child{app}, nil
	})
	key(m, "j")
	key(m, "j") // inbox
	key(m, "enter")
	if len(opened) != 1 || opened[0] != "inbox" {
		t.Fatalf("opened %v", opened)
	}
	if v := m.View().Content; v != "the inbox app" {
		t.Fatalf("view %q", v)
	}
	// Keys now go to the inbox app: q does not quit the health app's way.
	if _, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"}); cmd != nil {
		t.Fatalf("q reached the health app")
	}
}

func TestEnterOnDaemonOpensLiveApp(t *testing.T) {
	var calls int
	r := report
	var opened string
	m := NewModel(fakeRunner(t, &r, &calls), func(app string) (tea.Model, error) {
		opened = app
		return child{app}, nil
	})
	key(m, "enter")
	if opened != "live" || m.Child() == nil {
		t.Fatalf("opened %q", opened)
	}
}

func TestEnterWithoutAppSaysSo(t *testing.T) {
	var calls int
	r := report
	m := NewModel(fakeRunner(t, &r, &calls), func(app string) (tea.Model, error) {
		return nil, errors.New("unknown app " + app)
	})
	key(m, "j") // disk: no app
	key(m, "enter")
	if v := m.View().Content; !strings.Contains(v, "No app shows more of disk") {
		t.Errorf("view:\n%s", v)
	}
	key(m, "j") // inbox: the opener fails
	key(m, "enter")
	if m.Child() != nil || !strings.Contains(m.View().Content, "rhiza: unknown app inbox") {
		t.Errorf("view:\n%s", m.View().Content)
	}
}

func TestRegisteredOnT(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == Name {
			if a.Key != "t" {
				t.Fatalf("key %q", a.Key)
			}
			return
		}
	}
	t.Fatal("health app not registered")
}

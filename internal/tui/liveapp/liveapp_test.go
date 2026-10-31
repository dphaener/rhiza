package liveapp

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

type fakeData struct {
	ls     string
	daemon string
	status map[string]string
	calls  int
	opens  []string
	open   string
}

func (f *fakeData) run(op string, args []string) ([]byte, error) {
	f.calls++
	switch op {
	case "run.ls":
		return []byte(f.ls), nil
	case "daemon.status":
		return []byte(f.daemon), nil
	case "run.open":
		f.opens = append(f.opens, args[0])
		return []byte(f.open), nil
	case "run.status":
		return []byte(f.status[args[0]]), nil
	}
	return nil, nil
}

const ls1 = `{"runs":[
{"id":"R2","change":"C2","title":"Second","state":"running","current":"review","started_at":"2026-01-01T00:00:00Z","live":true,"workflow":"ship-change"},
{"id":"R1","change":"C1","title":"First","state":"succeeded","started_at":"2025-12-31T23:00:00Z","live":false,"workflow":"ship-change"}]}`

func newFake() *fakeData {
	return &fakeData{
		ls:     ls1,
		daemon: `{"running":true,"autonomous":false,"runs":[{"id":"R2","change":"C2","state":"running"}],"unjournaled":[]}`,
		status: map[string]string{
			"R2": `{"id":"R2","leaves":[{"leaf":"a","outcome":"passed"},{"leaf":"b","outcome":"passed"},{"leaf":"c"}],"checks":[{"check":"q","outcome":"passed"},{"check":"s","outcome":"failed"},{"check":"t"}]}`,
			"R1": `{"id":"R1","leaves":[{"leaf":"a","outcome":"passed"}],"checks":[]}`,
		},
	}
}

func clock(s string) Clock {
	t, _ := time.Parse(time.RFC3339, s)
	return func() time.Time { return t }
}

func TestViewShowsRunProgress(t *testing.T) {
	m := NewModel(newFake().run, clock("2026-01-01T00:05:07Z"))
	v := m.View().Content
	for _, want := range []string{
		"Daemon: running · 1 live runs",
		"> Second · running",
		"step review · leaves 2/3 · checks 1 passed, 1 failed, 1 running · elapsed 5m07s",
		"First · succeeded",
		"leaves 1/1 · checks none · elapsed 1h05m07s",
	} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
}

func TestTickUpdatesViewKeepingCursorOnRun(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, clock("2026-01-01T00:05:07Z"))
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if r, _ := m.Selected(); r.ID != "R1" {
		t.Fatalf("selected %q", r.ID)
	}
	// A newer run appears on top and R2's results change.
	f.ls = `{"runs":[{"id":"R3","change":"C3","title":"Third","state":"running","current":"build","started_at":"2026-01-01T00:05:00Z","live":true,"workflow":"ship-change"},` + ls1[len(`{"runs":[`):]
	f.status["R3"] = `{"id":"R3","leaves":[],"checks":[]}`
	f.status["R2"] = `{"id":"R2","leaves":[{"outcome":"passed"},{"outcome":"passed"},{"outcome":"passed"}],"checks":[]}`
	f.daemon = `{"running":false,"runs":[],"unjournaled":[]}`
	m.Update(TickMsg{})
	v := m.View().Content
	if r, _ := m.Selected(); r.ID != "R1" {
		t.Fatalf("cursor moved to %q", r.ID)
	}
	for _, want := range []string{"Daemon: not running", "Third · running", "leaves 3/3 · checks none"} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
}

func TestEmptyAndRegistered(t *testing.T) {
	f := newFake()
	f.ls = `{"runs":[]}`
	m := NewModel(f.run, nil)
	if !strings.Contains(m.View().Content, "No runs") {
		t.Fatalf("view: %q", m.View().Content)
	}
	found := false
	for _, a := range launcher.Apps() {
		found = found || a.Name == Name
	}
	if !found {
		t.Fatal("live app not in the launcher registry")
	}
}

func TestResultIsRunLsBytes(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	if out, err := m.Result(); err != nil || string(out) != ls1 {
		t.Fatalf("result %q, %v", out, err)
	}
}

func enter(m *Model) { m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) }

func TestEnterFocusesLiveRunTab(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	f.open = `{"change":"C2","live":true,"run":"R2","opened":true,"tab":"block:b1"}`
	enter(m)
	if len(f.opens) != 1 || f.opens[0] != "C2" {
		t.Fatalf("run.open calls: %v", f.opens)
	}
	if v := m.View().Content; !strings.Contains(v, "Focused Rex tab block:b1 of run R2") {
		t.Errorf("view lacks focused notice:\n%s", v)
	}
	f.open = `{"change":"C2","live":true,"run":"R2","opened":false,"tab_gone":true}`
	enter(m)
	if v := m.View().Content; !strings.Contains(v, "The tab of run R2 is gone") {
		t.Errorf("view lacks tab gone notice:\n%s", v)
	}
}

func TestEnterOnFinishedRunCallsNothing(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	enter(m)
	if len(f.opens) != 0 {
		t.Fatalf("run.open called: %v", f.opens)
	}
	if v := m.View().Content; !strings.Contains(v, "is not live") {
		t.Errorf("view lacks not-live notice:\n%s", v)
	}
}

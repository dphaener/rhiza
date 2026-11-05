package runsapp

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

const ls1 = `{"runs":[
{"id":"R2","change":"C2","title":"Second","state":"running","current":"review","started_at":"2026-01-01T00:00:00Z","live":true,"workflow":"ship-change"},
{"id":"R1","change":"C1","title":"First","state":"succeeded","started_at":"2025-12-31T23:00:00Z","live":false,"workflow":"ship-change"}]}`

type fake struct {
	ls  string
	ops []string
}

func (f *fake) run(op string, args []string) ([]byte, error) {
	f.ops = append(f.ops, op)
	return []byte(f.ls), nil
}

func newModel(f *fake) *Model {
	now := time.Date(2026, 1, 1, 0, 1, 5, 0, time.UTC)
	return NewModel(f.run, func() time.Time { return now })
}

func TestRowsAreLiveRunsOnly(t *testing.T) {
	m := newModel(&fake{ls: ls1})
	v := m.Render()
	if !strings.Contains(v, "Second") || !strings.Contains(v, "review") || !strings.Contains(v, "elapsed 1m05s") {
		t.Errorf("live run missing:\n%s", v)
	}
	if strings.Contains(v, "First") {
		t.Errorf("finished run listed:\n%s", v)
	}
}

func TestTickRefreshes(t *testing.T) {
	f := &fake{ls: ls1}
	m := newModel(f)
	f.ls = `{"runs":[{"id":"R3","change":"C3","title":"Third","state":"running","current":"build","started_at":"2026-01-01T00:00:00Z","live":true,"workflow":"ship-change"}]}`
	if cmd := m.Update(TickMsg{}); cmd == nil {
		t.Fatal("tick did not schedule the next")
	}
	v := m.Render()
	if !strings.Contains(v, "Third") || strings.Contains(v, "Second") {
		t.Errorf("rows not refreshed:\n%s", v)
	}
}

func TestEmptyState(t *testing.T) {
	m := newModel(&fake{ls: `{"runs":[]}`})
	if v := m.Render(); strings.Count(v, "No live runs") != 1 {
		t.Errorf("empty state:\n%s", v)
	}
}

func TestOnlyReadOperations(t *testing.T) {
	f := &fake{ls: ls1}
	m := newModel(f)
	m.Update(TickMsg{})
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if len(f.ops) < 2 {
		t.Fatalf("ops = %v", f.ops)
	}
	for _, op := range f.ops {
		if op != "run.ls" {
			t.Errorf("ran %q", op)
		}
	}
}

func TestResultIsRunLs(t *testing.T) {
	m := newModel(&fake{ls: ls1})
	if out, err := m.Result(); string(out) != ls1 || err != nil {
		t.Errorf("result = %q, %v", out, err)
	}
}

func TestClosesInShell(t *testing.T) {
	s := modalshell.New(newModel(&fake{ls: ls1}))
	_, cmd := s.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	s.Update(cmd())
	if !s.Closed() {
		t.Error("q did not close")
	}
}

func TestRegistered(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == Name && a.Key == Key {
			return
		}
	}
	t.Fatal("runs app not registered")
}

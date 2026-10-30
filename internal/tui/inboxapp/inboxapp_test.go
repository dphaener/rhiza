package inboxapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/tui/launcher"
)

const full = `{"items":[
{"id":"a1","kind":"skipped","change":"C5","change_title":"Five","step":"build","detail":"interrupted, and not resumed"},
{"id":"b2","kind":"gave-up","change":"C4","change_title":"Four","step":"review","detail":"interrupted again"},
{"id":"c3","kind":"failed","change":"C3","change_title":"Three","step":"verify","detail":"boom"},
{"id":"d4","kind":"signoff","change":"C2","change_title":"Two","step":"signoff","detail":"signoff-wait"},
{"id":"e5","kind":"question","change":"C1","change_title":"One","step":"build","question":"Which db?","recommendation":"sqlite","reasons":"simple"}
]}`

func fake(out *string, calls *int) Runner {
	return func(op string, args []string) ([]byte, error) {
		*calls++
		return []byte(*out), nil
	}
}

func TestViewGroupsByKind(t *testing.T) {
	out, calls := full, 0
	m := NewModel(fake(&out, &calls))
	v := m.View().Content
	last := -1
	for _, h := range []string{"Questions", "Sign-offs", "Failed runs", "Gave up", "Skipped"} {
		i := strings.Index(v, h)
		if i < 0 || i < last {
			t.Fatalf("heading %q missing or out of order in:\n%s", h, v)
		}
		last = i
	}
	for _, want := range []string{"One · build", "Which db?", "recommends: sqlite", "Two · signoff", "Three · verify", "boom", "Four · review", "Five · build"} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
}

func TestEmptyInbox(t *testing.T) {
	out, calls := `{"items":[]}`, 0
	m := NewModel(fake(&out, &calls))
	if !strings.Contains(m.View().Content, "Nothing waits on you") {
		t.Fatalf("view: %q", m.View().Content)
	}
}

func TestTickAndRReload(t *testing.T) {
	out, calls := `{"items":[]}`, 0
	m := NewModel(fake(&out, &calls))
	out = full
	m.Update(TickMsg{})
	if calls != 2 || !strings.Contains(m.View().Content, "Which db?") {
		t.Fatalf("tick did not reload: calls=%d", calls)
	}
	m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	if calls != 3 {
		t.Fatalf("r did not reload: calls=%d", calls)
	}
}

func TestCursorMovesAndSurvivesReload(t *testing.T) {
	out, calls := full, 0
	m := NewModel(fake(&out, &calls))
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if it, _ := m.Selected(); it.ID != "d4" {
		t.Fatalf("selected %q", it.ID)
	}
	m.Reload()
	if it, _ := m.Selected(); it.ID != "d4" {
		t.Fatalf("after reload selected %q", it.ID)
	}
}

func TestResultIsOperationOutput(t *testing.T) {
	out, calls := full, 0
	got, err := NewModel(fake(&out, &calls)).Result()
	if err != nil || string(got) != full {
		t.Fatalf("Result() = %q, %v", got, err)
	}
}

func TestRegisteredWithLauncher(t *testing.T) {
	for _, a := range launcher.Apps() {
		if a.Name == Name {
			return
		}
	}
	t.Fatal("inbox is not in the launcher registry")
}

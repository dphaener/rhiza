package readapps

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestAppsRegisterWithLauncher(t *testing.T) {
	for _, name := range []string{"roadmap", "ready", "show", "why", "slice"} {
		if _, ok := Lookup(name); !ok {
			t.Errorf("no app %q", name)
		}
	}
}

func TestModelAsksForRefThenShowsResult(t *testing.T) {
	var got []string
	run := func(op string, args []string) ([]byte, error) {
		got = args
		return []byte("{\n  \"ok\": true\n}\n"), nil
	}
	app, _ := Lookup("show")
	m := NewModel(app, run, nil)
	if !strings.Contains(m.View().Content, "show (ID or handle):") {
		t.Fatalf("want prompt, got %q", m.View().Content)
	}
	for _, r := range "store" {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if len(got) != 1 || got[0] != "store" {
		t.Fatalf("run args = %v", got)
	}
	if !strings.Contains(m.View().Content, `"ok": true`) {
		t.Fatalf("result not shown: %q", m.View().Content)
	}
}

func TestResultIsWhatRunReturned(t *testing.T) {
	want := []byte("{\n  \"ok\": true\n}\n")
	run := func(op string, args []string) ([]byte, error) { return want, nil }
	app, _ := Lookup("roadmap")
	got, err := NewModel(app, run, nil).Result()
	if err != nil || string(got) != string(want) {
		t.Fatalf("Result() = %q, %v", got, err)
	}
	show, _ := Lookup("show")
	if _, err := NewModel(show, run, nil).Result(); err == nil {
		t.Fatal("show without a ref has a result")
	}
}

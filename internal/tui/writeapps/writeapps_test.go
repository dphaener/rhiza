package writeapps

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestAppsRegisterWithLauncher(t *testing.T) {
	for _, name := range []string{"change-show", "propose", "accept", "land", "record-add", "record-set", "record-link", "record-retire"} {
		if _, ok := Lookup(name); !ok {
			t.Errorf("no app %q", name)
		}
	}
}

func TestSplit(t *testing.T) {
	got, err := Split(`capability store --field "statement=Kept safe" ''`)
	want := []string{"capability", "store", "--field", "statement=Kept safe", ""}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("Split = %q, %v", got, err)
	}
	if _, err := Split(`a "b`); err == nil {
		t.Fatal("unterminated quote accepted")
	}
}

type call struct {
	dry  bool
	args []string
}

func recorder(calls *[]call) Runner {
	return func(a App, args []string, dry bool) ([]byte, error) {
		*calls = append(*calls, call{dry, args})
		if dry {
			return []byte("{\n  \"dry_run\": true\n}\n"), nil
		}
		return []byte("{\n  \"entry\": \"E\"\n}\n"), nil
	}
}

func key(m *Model, s string) {
	for _, r := range s {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

func TestCancelWritesNothing(t *testing.T) {
	var calls []call
	app, _ := Lookup("propose")
	m := NewModel(app, recorder(&calls), nil)
	if !strings.Contains(m.View().Content, `"dry_run": true`) {
		t.Fatalf("no preview: %q", m.View().Content)
	}
	key(m, "n")
	if !m.Cancelled() {
		t.Fatal("not cancelled")
	}
	if len(calls) != 1 || !calls[0].dry {
		t.Fatalf("calls = %+v", calls)
	}
}

func TestConfirmRunsForRealOnce(t *testing.T) {
	var calls []call
	app, _ := Lookup("record-add")
	m := NewModel(app, recorder(&calls), nil)
	if !strings.Contains(m.View().Content, "add --record") {
		t.Fatalf("want prompt: %q", m.View().Content)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // empty is refused
	if len(calls) != 0 {
		t.Fatal("ran without arguments")
	}
	key(m, "capability store")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if len(calls) != 1 || !calls[0].dry {
		t.Fatalf("preview calls = %+v", calls)
	}
	key(m, "y")
	if len(calls) != 2 || calls[1].dry || !reflect.DeepEqual(calls[1].args, []string{"capability", "store"}) {
		t.Fatalf("calls = %+v", calls)
	}
	if !strings.Contains(m.View().Content, `"entry": "E"`) {
		t.Fatalf("result not shown: %q", m.View().Content)
	}
}

func TestFailedPreviewOffersNoConfirm(t *testing.T) {
	run := func(App, []string, bool) ([]byte, error) { return nil, errTest }
	app, _ := Lookup("land")
	m := NewModel(app, run, nil)
	key(m, "y")
	if _, err := m.Confirm(); err == nil {
		t.Fatal("confirmed a failed preview")
	}
}

type testErr struct{}

func (testErr) Error() string { return "boom" }

var errTest = testErr{}

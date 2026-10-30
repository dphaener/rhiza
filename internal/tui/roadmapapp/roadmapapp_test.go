package roadmapapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

const (
	a = "01M6730000AAAAAAAAAAAAAAAA"
	b = "01M6730000BBBBBBBBBBBBBBBB"
	c = "01M6730000CCCCCCCCCCCCCCCC"
)

const fakeRoadmap = `{"groups":[
 {"kind":"milestone","intent":{"id":"01M67300000000000000000000","handle":"milestone-5","path":["rhiza"]},
  "change_sets":[
   {"id":"` + a + `","state":"accepted","kind":"intention","intent":"First change","ready":true,"depends_on":[],"serves":[]},
   {"id":"` + b + `","state":"proposed","kind":"intention","intent":"Second change","ready":false,"depends_on":[],"serves":[]}]},
 {"kind":"backlog","change_sets":[
   {"id":"` + c + `","state":"accepted","kind":"intention","intent":"Loose change","ready":false,"depends_on":[],"serves":[]}]}],
 "view":"to-be"}`

const fakeReady = `{"change_sets":[
 {"id":"` + a + `","state":"accepted","intent":"First change","ready":true,"reasons":[]},
 {"id":"` + b + `","state":"proposed","intent":"Second change","ready":false,"reasons":["not accepted"]},
 {"id":"` + c + `","state":"accepted","intent":"Loose change","ready":false,"reasons":["waits on a question"]}],
 "view":"to-be"}`

func fake(op string, args []string) ([]byte, error) {
	if op == "ready" {
		return []byte(fakeReady), nil
	}
	return []byte(fakeRoadmap), nil
}

func key(m *Model, s string) {
	m.Update(tea.KeyPressMsg{Code: rune(s[0]), Text: s})
}

func TestViewIsAListNotJSON(t *testing.T) {
	v := NewModel(fake).View().Content
	for _, want := range []string{"rhiza / milestone-5", "Backlog", a[:14] + " [accepted] First change",
		b[:14] + " [proposed] Second change", "not ready: not accepted", "waits on a question"} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
	if strings.ContainsAny(v, "{}") {
		t.Errorf("view has JSON braces:\n%s", v)
	}
	if strings.Contains(v, "First change (not ready") {
		t.Errorf("a ready change set shows reasons:\n%s", v)
	}
}

func TestCursorSkipsHeadings(t *testing.T) {
	m := NewModel(fake)
	if got := m.Cursor().String(); got != a {
		t.Fatalf("cursor starts on %s", got)
	}
	key(m, "j")
	if got := m.Cursor().String(); got != b {
		t.Fatalf("j moved to %s, want %s", got, b)
	}
	key(m, "j") // crosses the Backlog heading
	if got := m.Cursor().String(); got != c {
		t.Fatalf("j moved to %s, want %s", got, c)
	}
	key(m, "j")
	if got := m.Cursor().String(); got != c {
		t.Fatalf("j past the end moved to %s", got)
	}
	key(m, "k")
	if got := m.Cursor().String(); got != b {
		t.Fatalf("k moved to %s, want %s", got, b)
	}
}

func TestReloadKeepsCursorAndResultIsRoadmapBytes(t *testing.T) {
	m := NewModel(fake)
	key(m, "j")
	key(m, "r")
	if got := m.Cursor().String(); got != b {
		t.Fatalf("reload moved cursor to %s", got)
	}
	out, err := m.Result()
	if err != nil || string(out) != fakeRoadmap {
		t.Fatalf("Result() = %q, %v", out, err)
	}
}

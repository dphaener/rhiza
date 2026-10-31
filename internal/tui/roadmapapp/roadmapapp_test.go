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

func fake(op, _ string, args []string) ([]byte, error) {
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

const fakeShow = `{"id":"` + b + `","kind":"intention","intent":"Second change","title":"Second","state":"proposed",
 "leaves":[{"id":"01M6730000EEEEEEEEEEEEEEEE","handle":"leaf-one","needs":[]}],
 "ops":[{"op":"node.create","node":"01M6730000EEEEEEEEEEEEEEEE","handle":"leaf-one",
  "fields":{"statement":"Do the first thing","acceptance":"A test proves the first thing"}}]}`

const fakeQuestion = `{"kind":"node","node":{"id":"01M6730000FFFFFFFFFFFFFFFF","kind":"question","handle":"which-way",
 "fields":{"statement":"Which way do we go?"}},"view":"to-be"}`

const fakeRun = `{"id":"01M6730000GGGGGGGGGGGGGGGG","workflow":"ship","change":"` + b + `","state":"stopped","current":"build","steps":[]}`

func detailFake(op, change string, args []string) ([]byte, error) {
	switch op {
	case "change.show":
		return []byte(fakeShow), nil
	case "show":
		return []byte(fakeQuestion), nil
	case "run.status":
		return []byte(fakeRun), nil
	case "ready":
		return []byte(`{"change_sets":[
 {"id":"` + a + `","state":"accepted","intent":"First change","ready":true,"reasons":[]},
 {"id":"` + b + `","state":"proposed","intent":"Second change","ready":false,"reasons":["blocked by question which-way on leaf-one"]},
 {"id":"` + c + `","state":"accepted","intent":"Loose change","ready":false,"reasons":[]}],"view":"to-be"}`), nil
	}
	return []byte(fakeRoadmap), nil
}

func TestEnterOpensDetailAndEscReturns(t *testing.T) {
	m := NewModel(detailFake)
	key(m, "j")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	v := m.View().Content
	for _, want := range []string{"Second change", "intention", "proposed", "leaf-one", "Do the first thing",
		"A test proves the first thing", "Which way do we go?", "ship", "stopped", "build"} {
		if !strings.Contains(v, want) {
			t.Errorf("detail lacks %q:\n%s", want, v)
		}
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if got := m.Cursor().String(); got != b {
		t.Fatalf("esc left cursor on %s, want %s", got, b)
	}
	if v := m.View().Content; !strings.Contains(v, "Backlog") {
		t.Errorf("esc did not return to the list:\n%s", v)
	}
}

type call struct {
	name string
	args []string
}

// startModel is a model with a fake exec recording every call and a runner
// answering run.open with runOpen.
func startModel(runOpen string, calls *[]call) *Model {
	m := NewModel(func(op, ch string, args []string) ([]byte, error) {
		if op == "run.open" {
			return []byte(runOpen), nil
		}
		return fake(op, ch, args)
	})
	m.SetStarter(Starter{
		Exec: func(name string, args []string, _ string) ([]byte, error) {
			*calls = append(*calls, call{name, args})
			return []byte("{}"), nil
		},
		Self: "/bin/rhiza", Root: "/repo", RexSession: "session:1",
	})
	return m
}

func TestPlanRunStartsInteractiveInNewTab(t *testing.T) {
	var calls []call
	m := startModel("{}", &calls)
	key(m, "p")
	if len(calls) != 0 {
		t.Fatalf("p called rex before y: %v", calls)
	}
	if v := m.View().Content; !strings.Contains(v, "plan run") || !strings.Contains(v, "y/n") {
		t.Errorf("no y/n question:\n%s", v)
	}
	key(m, "y")
	if len(calls) != 1 || calls[0].name != "rex" {
		t.Fatalf("calls = %v, want one rex call", calls)
	}
	got := strings.Join(calls[0].args, " ")
	want := "window new -s session:1 " + a[:9] + " -c /repo --json -- /bin/rhiza run start " + a + " --workflow plan"
	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
	if strings.Contains(got, "--autonomous") {
		t.Errorf("plan run is autonomous: %s", got)
	}
	if v := m.View().Content; !strings.Contains(v, "Opened Rex tab "+a[:9]) {
		t.Errorf("view does not say which tab opened:\n%s", v)
	}
	if m.Cursor().String() != a {
		t.Errorf("app did not stay on the list")
	}
}

func TestShipRunIsAutonomousWithAutoMode(t *testing.T) {
	var calls []call
	m := startModel("{}", &calls)
	key(m, "s")
	key(m, "y")
	if len(calls) != 1 {
		t.Fatalf("calls = %v", calls)
	}
	got := strings.Join(calls[0].args, " ")
	if want := "run start " + a + " --autonomous --permission-mode auto"; !strings.HasSuffix(got, want) {
		t.Errorf("args = %q, want suffix %q", got, want)
	}
}

func TestNDeclinesAndCallsNothing(t *testing.T) {
	var calls []call
	m := startModel("{}", &calls)
	key(m, "s")
	key(m, "n")
	key(m, "p")
	key(m, "n")
	if len(calls) != 0 {
		t.Errorf("declined runs called rex: %v", calls)
	}
	if strings.Contains(m.View().Content, "y/n") {
		t.Errorf("question still showing")
	}
}

func TestOFocusesLiveRunOrSaysNone(t *testing.T) {
	var calls []call
	m := startModel(`{"change":"`+a+`","live":true,"run":"R1","opened":true,"tab":"block:7"}`, &calls)
	key(m, "o")
	if v := m.View().Content; !strings.Contains(v, "Focused Rex tab block:7 of run R1") {
		t.Errorf("view:\n%s", v)
	}
	m = startModel(`{"change":"`+a+`","live":false,"opened":false}`, &calls)
	key(m, "o")
	if v := m.View().Content; !strings.Contains(v, "No live run") {
		t.Errorf("view:\n%s", v)
	}
	if len(calls) != 0 {
		t.Errorf("o started a run: %v", calls)
	}
}

func TestTickReloadsOnlyTheRunSectionOfAnOpenDetail(t *testing.T) {
	state := "running"
	counts := map[string]int{}
	fk := func(op, change string, args []string) ([]byte, error) {
		counts[op]++
		if op == "run.status" {
			return []byte(`{"id":"01M6730000GGGGGGGGGGGGGGGG","workflow":"ship","change":"` + b + `","state":"` + state + `","steps":[]}`), nil
		}
		return detailFake(op, change, args)
	}
	m := NewModel(fk)
	key(m, "j")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if v := m.View().Content; !strings.Contains(v, "state running") {
		t.Fatalf("detail lacks running:\n%s", v)
	}
	shows := counts["change.show"]
	state = "succeeded"
	_, cmd := m.Update(TickMsg{})
	if cmd == nil {
		t.Fatal("tick did not re-arm")
	}
	if v := m.View().Content; !strings.Contains(v, "state succeeded") {
		t.Fatalf("run section not updated:\n%s", v)
	}
	if counts["change.show"] != shows {
		t.Fatalf("tick called change.show again")
	}
}

func TestTickInListCallsNoRunStatus(t *testing.T) {
	n := 0
	fk := func(op, change string, args []string) ([]byte, error) {
		if op == "run.status" {
			n++
		}
		return detailFake(op, change, args)
	}
	m := NewModel(fk)
	if _, cmd := m.Update(TickMsg{}); cmd == nil {
		t.Fatal("tick did not re-arm")
	}
	if n != 0 {
		t.Fatalf("list tick called run.status %d times", n)
	}
}

func TestAnswerCallsQuestionAnswer(t *testing.T) {
	type rec struct {
		op, change string
		args       []string
	}
	var got []rec
	m := NewModel(func(op, change string, args []string) ([]byte, error) {
		if op == "question.answer" {
			got = append(got, rec{op, change, args})
			return []byte(`{}`), nil
		}
		return detailFake(op, change, args)
	})
	key(m, "j")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // opens the detail
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // starts the answer
	for _, r := range "go west" {
		key(m, string(r))
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if len(got) != 1 || got[0].change != b || strings.Join(got[0].args, "|") != "which-way|go west" {
		t.Fatalf("question.answer calls: %+v", got)
	}
}

const umbRoadmap = `{"groups":[
 {"kind":"backlog","change_sets":[
   {"id":"` + c + `","state":"accepted","kind":"intention","intent":"Loose change","ready":false,"depends_on":[],"serves":[]}]},
 {"kind":"umbrella","intent":{"id":"01M67300000000000000000001","handle":"areas","path":["rhiza"]},
  "change_sets":[
   {"id":"` + a + `","state":"accepted","kind":"intention","intent":"Area one","ready":false,"depends_on":[],"serves":[]},
   {"id":"` + b + `","state":"proposed","kind":"intention","intent":"Area two","ready":false,"depends_on":[],"serves":[]}]}],
 "view":"to-be"}`

func umbFake(op, _ string, args []string) ([]byte, error) {
	if op == "ready" {
		return []byte(`{"change_sets":[],"view":"to-be"}`), nil
	}
	return []byte(umbRoadmap), nil
}

func TestUmbrellaGroupIsOneRowUntilEnter(t *testing.T) {
	m := NewModel(umbFake)
	v := m.View().Content
	if !strings.Contains(v, "rhiza / areas (2 change sets)") {
		t.Errorf("collapsed view lacks the umbrella row:\n%s", v)
	}
	if strings.Contains(v, "Area one") || strings.Contains(v, "Area two") {
		t.Errorf("collapsed view shows members:\n%s", v)
	}
	key(m, "j") // onto the umbrella row
	if !m.Cursor().IsZero() {
		t.Fatalf("cursor on a change set, want the umbrella row")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	v = m.View().Content
	for _, want := range []string{"rhiza / areas (2 change sets)", "Area one", "Area two"} {
		if !strings.Contains(v, want) {
			t.Errorf("expanded view lacks %q:\n%s", want, v)
		}
	}
	key(m, "j")
	if got := m.Cursor().String(); got != a {
		t.Fatalf("j moved to %s, want %s", got, a)
	}
	key(m, "k")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if v = m.View().Content; strings.Contains(v, "Area one") {
		t.Errorf("collapsed again but members show:\n%s", v)
	}
	if out, _ := m.Result(); string(out) != umbRoadmap {
		t.Errorf("Result changed")
	}
}

func typeText(m *Model, s string) {
	for _, r := range s {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

func enter(m *Model) { m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) }

func TestNewWorkStartsPlanRunWithNoChangeSet(t *testing.T) {
	for _, tc := range []struct {
		key, want string
	}{
		{"m", "window new -s session:1 new-milestone -c /repo --json -- /bin/rhiza run start --workflow plan --new Ship the thing --milestone"},
		{"n", "window new -s session:1 new-change-set -c /repo --json -- /bin/rhiza run start --workflow plan --new Ship the thing"},
	} {
		var calls []call
		m := startModel("{}", &calls)
		key(m, tc.key)
		typeText(m, "Ship the thing")
		enter(m)
		if len(calls) != 0 {
			t.Fatalf("%s called rex before y: %v", tc.key, calls)
		}
		if v := m.View().Content; !strings.Contains(v, "y/n") {
			t.Errorf("no y/n question:\n%s", v)
		}
		key(m, "y")
		if len(calls) != 1 || calls[0].name != "rex" {
			t.Fatalf("%s calls = %v", tc.key, calls)
		}
		if got := strings.Join(calls[0].args, " "); got != tc.want {
			t.Errorf("%s args = %q, want %q", tc.key, got, tc.want)
		}
		if v := m.View().Content; !strings.Contains(v, "Opened Rex tab") {
			t.Errorf("view does not say a tab opened:\n%s", v)
		}
	}
}

func TestNewWorkDeclinedOrEmptyCallsNothing(t *testing.T) {
	var calls []call
	m := startModel("{}", &calls)
	key(m, "m")
	enter(m) // empty text stays in the prompt
	if m.intent == nil {
		t.Fatalf("empty enter left the prompt")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	key(m, "n")
	typeText(m, "x")
	enter(m)
	key(m, "n")
	key(m, "m")
	typeText(m, "y")
	enter(m)
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if len(calls) != 0 {
		t.Errorf("declined or empty new work called rex: %v", calls)
	}
}

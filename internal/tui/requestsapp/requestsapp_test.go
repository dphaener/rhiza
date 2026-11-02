package requestsapp

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// fake is a runner over an in-memory queue.
type fake struct {
	reqs  []api.RequestInfo
	calls [][]string
	open  string
}

func (f *fake) run(op string, args []string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{op}, args...))
	switch op {
	case "request.ls":
		return json.Marshal(api.RequestLsResponse{Requests: append([]api.RequestInfo{}, f.reqs...)})
	case "request":
		if len(args) != 4 || args[0] != "--kind" || args[2] != "--" {
			return nil, fmt.Errorf("unexpected request args %q", args)
		}
		q := api.RequestInfo{ID: "R" + string(rune('1'+len(f.reqs))), Kind: args[1], Text: args[3], State: "queued"}
		f.reqs = append(f.reqs, q)
		return json.Marshal(api.RequestResponse{Request: q})
	case "run.open":
		return []byte(f.open), nil
	}
	return nil, nil
}

func key(m *Model, s string) {
	var k tea.KeyPressMsg
	switch s {
	case "enter":
		k = tea.KeyPressMsg{Code: tea.KeyEnter}
	case "tab":
		k = tea.KeyPressMsg{Code: tea.KeyTab}
	case "esc":
		k = tea.KeyPressMsg{Code: tea.KeyEscape}
	default:
		k = tea.KeyPressMsg{Code: rune(s[0]), Text: s}
	}
	m.Update(k)
}

func typeText(m *Model, s string) {
	for _, r := range s {
		key(m, string(r))
	}
}

func TestFilesBugAndFeatureFromTypedText(t *testing.T) {
	f := &fake{}
	m := NewModel(f.run)
	if !strings.Contains(m.View().Content, "No requests yet") {
		t.Fatalf("empty state missing:\n%s", m.View().Content)
	}

	key(m, "b")
	if !m.Typing() {
		t.Fatal("b did not open the entry")
	}
	typeText(m, "login fails")
	key(m, "enter")
	if m.Typing() || len(f.reqs) != 1 || f.reqs[0].Kind != "bug" || f.reqs[0].Text != "login fails" || f.reqs[0].State != "queued" {
		t.Fatalf("bug not filed: typing=%v reqs=%+v", m.Typing(), f.reqs)
	}

	key(m, "n")
	typeText(m, "dark mode")
	key(m, "enter")
	if len(f.reqs) != 2 || f.reqs[1].Kind != "feature" {
		t.Fatalf("feature not filed: %+v", f.reqs)
	}
	v := m.View().Content
	for _, want := range []string{"bug · login fails", "feature · dark mode", "queued"} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
}

func TestDashLeadingTextIsFiledAsText(t *testing.T) {
	f := &fake{}
	m := NewModel(f.run)
	key(m, "b")
	typeText(m, "--json output is wrong")
	key(m, "enter")
	if m.Typing() || len(f.reqs) != 1 || f.reqs[0].Text != "--json output is wrong" || f.reqs[0].Kind != "bug" {
		t.Fatalf("dash-leading bug not filed: typing=%v reqs=%+v calls=%q", m.Typing(), f.reqs, f.calls)
	}
}

func TestTabCyclesKindAndEmptyTextIsRefused(t *testing.T) {
	f := &fake{}
	m := NewModel(f.run)
	key(m, "n")
	key(m, "enter")
	if !m.Typing() || len(f.reqs) != 0 || !strings.Contains(m.View().Content, "needs its text") {
		t.Fatalf("empty text was not refused:\n%s", m.View().Content)
	}
	key(m, "tab")
	key(m, "tab")
	typeText(m, "ship it")
	key(m, "enter")
	if len(f.reqs) != 1 || f.reqs[0].Kind != "milestone" {
		t.Fatalf("kind = %+v", f.reqs)
	}
	key(m, "n")
	typeText(m, "x")
	key(m, "esc")
	if m.Typing() || len(f.reqs) != 1 {
		t.Fatal("esc did not cancel")
	}
}

func TestEnterFocusesPlanRunTab(t *testing.T) {
	f := &fake{
		reqs: []api.RequestInfo{
			{ID: "A", Kind: "feature", Text: "one", State: "queued"},
			{ID: "B", Kind: "bug", Text: "two", State: "planning", Run: "RUN2"},
		},
		open: `{"live":true,"run":"RUN2","opened":true,"tab":"block:7"}`,
	}
	m := NewModel(f.run)
	key(m, "enter")
	if !strings.Contains(m.View().Content, "still queued") {
		t.Fatalf("queued request:\n%s", m.View().Content)
	}
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	key(m, "enter")
	last := f.calls[len(f.calls)-1]
	if strings.Join(last, " ") != "run.open --run RUN2" {
		t.Fatalf("last call = %v", last)
	}
	if !strings.Contains(m.View().Content, "Focused Rex tab block:7") {
		t.Fatalf("notice:\n%s", m.View().Content)
	}
}

func TestResultIsRequestLsAndAppIsRegistered(t *testing.T) {
	f := &fake{reqs: []api.RequestInfo{{ID: "A", Kind: "feature", Text: "one", State: "queued"}}}
	m := NewModel(f.run)
	out, err := m.Result()
	want, _ := f.run("request.ls", nil)
	if err != nil || string(out) != string(want) {
		t.Fatalf("Result = %s, %v", out, err)
	}
	found := false
	for _, a := range launcher.Apps() {
		found = found || a.Name == Name
	}
	if !found {
		t.Fatal("requests not registered")
	}
}

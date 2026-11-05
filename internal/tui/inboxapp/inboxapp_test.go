package inboxapp

import (
	"errors"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rexuikit/ui"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

const full = `{"items":[
{"id":"e5","kind":"question","change":"C1","change_title":"One","step":"build","run":"R1","session":"S1","question":"Which db?"},
{"id":"d4","kind":"signoff","change":"C2","change_title":"Two","step":"signoff","run":"R2","detail":"signoff-wait"},
{"id":"c3","kind":"failed","change":"C3","change_title":"Three","step":"verify","run":"R3","detail":"boom"},
{"id":"h8","kind":"stalled","change":"C9","change_title":"Nine","step":"build","run":"R9","session":"S1","detail":"session S1 is silent"},
{"id":"i9","kind":"disk-low","detail":"1.0 GiB free"},
{"id":"j1","kind":"accept","change":"C7","change_title":"Seven","detail":"ready to accept"}
]}`

type call struct {
	op   string
	args []string
}

func fake(out string, calls *[]call) Runner {
	return func(op string, args []string) ([]byte, error) {
		*calls = append(*calls, call{op, args})
		if op == Name {
			return []byte(out), nil
		}
		return []byte("{}"), nil
	}
}

func press(m *Model, key string) tea.Cmd {
	var msg tea.KeyPressMsg
	switch key {
	case "enter":
		msg = tea.KeyPressMsg{Code: tea.KeyEnter}
	case "down":
		msg = tea.KeyPressMsg{Code: tea.KeyDown}
	case "esc":
		msg = tea.KeyPressMsg{Code: tea.KeyEscape}
	}
	_, cmd := m.Update(msg)
	return cmd
}

// enter presses Enter and feeds the resulting selection back in.
func enter(m *Model) tea.Cmd {
	cmd := press(m, "enter")
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if _, ok := msg.(ui.ListSelectedMsg); !ok {
		return cmd
	}
	_, cmd = m.Update(msg)
	return cmd
}

func TestRowPerItemInOrder(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	rows := m.Rows()
	if len(rows) != 6 {
		t.Fatalf("rows = %d: %q", len(rows), rows)
	}
	for i, want := range []string{"question · One · build  Which db?", "signoff · Two · signoff  signoff-wait", "failed · Three · verify  boom", "stalled · Nine · build", "disk-low · Daemon holds starts  1.0 GiB free", "accept · Seven  ready to accept"} {
		if !strings.HasPrefix(rows[i], want) {
			t.Errorf("row %d = %q, want prefix %q", i, rows[i], want)
		}
	}
	v := m.View().Content
	for _, w := range []string{"question", "One", "Which db?"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
}

func TestEnterOpensActionPerKind(t *testing.T) {
	for i, want := range []Action{ActionAnswer, ActionSignoff, ActionResolve} {
		if i == 0 {
			continue // a question opens the answer form, not the Open hook
		}
		var calls []call
		m := NewModel(fake(full, &calls))
		var opened []Pick
		m.Open = func(p Pick) tea.Cmd { opened = append(opened, p); return nil }
		for range i {
			press(m, "down")
		}
		enter(m)
		if len(opened) != 1 || opened[0].Action != want {
			t.Fatalf("row %d: opened %+v, want %s", i, opened, want)
		}
		if len(calls) != 1 {
			t.Errorf("row %d: runner called for an action: %+v", i, calls)
		}
	}
	var calls []call
	m := NewModel(fake(full, &calls))
	var opened []Pick
	m.Open = func(p Pick) tea.Cmd { opened = append(opened, p); return nil }
	for range 5 {
		press(m, "down")
	}
	enter(m)
	if len(opened) != 1 || opened[0].Action != ActionAccept || opened[0].Item.ID != "j1" {
		t.Fatalf("accept: %+v", opened)
	}
}

func TestActionFor(t *testing.T) {
	for kind, want := range map[string]Action{
		api.InboxQuestion: ActionAnswer, api.InboxOpen: ActionAnswer, api.InboxBlocker: ActionAnswer,
		api.InboxSignoff: ActionSignoff,
		api.InboxFailed:  ActionResolve, api.InboxGaveUp: ActionResolve, api.InboxSkipped: ActionResolve,
		api.InboxAccept:  ActionAccept,
		api.InboxStalled: ActionShow, api.InboxOverlap: ActionShow, api.InboxDiskLow: ActionShow,
		api.InboxUsageReserve: ActionShow, api.InboxMainStale: ActionShow,
	} {
		if got := ActionFor(kind); got != want {
			t.Errorf("ActionFor(%s) = %s, want %s", kind, got, want)
		}
	}
}

func TestOtherKindShowsDetailAndOpensTab(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	for range 3 {
		press(m, "down")
	}
	enter(m)
	if !slices.Equal(calls[len(calls)-1].args, []string{"h8"}) || calls[len(calls)-1].op != "inbox.open" {
		t.Fatalf("calls: %+v", calls)
	}
	v := m.View().Content
	for _, w := range []string{"stalled", "Nine", "session S1 is silent"} {
		if !strings.Contains(v, w) {
			t.Errorf("detail lacks %q:\n%s", w, v)
		}
	}
	if p, _ := m.Picked(); p.Action != ActionShow {
		t.Errorf("pick = %+v", p)
	}
}

func TestOpenTabFailureShown(t *testing.T) {
	run := func(op string, args []string) ([]byte, error) {
		if op == Name {
			return []byte(full), nil
		}
		return nil, errors.New("no tab")
	}
	m := NewModel(run)
	for range 3 {
		press(m, "down")
	}
	enter(m)
	if !strings.Contains(m.View().Content, "no tab") {
		t.Fatalf("view: %s", m.View().Content)
	}
}

func TestEmptyInbox(t *testing.T) {
	var calls []call
	m := NewModel(fake(`{"items":[]}`, &calls))
	if got := strings.TrimSpace(m.View().Content); got != emptyLine {
		t.Fatalf("view: %q", got)
	}
}

func TestEscQuits(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	_, cmd := m.Update(ui.ListCancelledMsg{})
	if cmd == nil {
		t.Fatal("no quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("not a quit")
	}
}

func TestResultIsOperationOutput(t *testing.T) {
	var calls []call
	got, err := NewModel(fake(full, &calls)).Result()
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

func typeAndSubmit(m *Model, text string) tea.Cmd {
	for _, r := range text {
		m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if _, ok := msg.(ui.FormSubmittedMsg); ok {
		_, c := m.Update(msg)
		return c
	}
	return nil
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestAnswerQuestionSendsSessionMessage(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	enter(m)
	if !strings.Contains(m.View().Content, "Your answer") {
		t.Fatalf("no form:\n%s", m.View().Content)
	}
	if !isQuit(typeAndSubmit(m, "Use sqlite")) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "session.message" || !slices.Equal(last.args, []string{"S1", "Use sqlite"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestAnswerBlockerCallsQuestionAnswer(t *testing.T) {
	var calls []call
	out := `{"items":[{"id":"b1","kind":"blocker","change":"C5","change_title":"Five","handle":"which-db","question":"Which db?"}]}`
	m := NewModel(fake(out, &calls))
	enter(m)
	if !isQuit(typeAndSubmit(m, "sqlite")) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "question.answer" || !slices.Equal(last.args, []string{"which-db", "C5", "sqlite"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestAnswerEscAndEmptyCallNothing(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	enter(m)
	if isQuit(typeAndSubmit(m, "")) {
		t.Fatal("empty answer quit")
	}
	_, cmd := m.Update(ui.FormCancelledMsg{})
	if !isQuit(cmd) {
		t.Fatal("esc did not quit")
	}
	if len(calls) != 1 {
		t.Fatalf("calls: %+v", calls)
	}
}

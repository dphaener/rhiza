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
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","detail":"open questions: which-db, how-many; answer each with a Decision, then rhiza run resume R4"}]}`
	for _, c := range []struct {
		name, out string
		downs     int
		want      []string
	}{
		{"question", full, 0, []string{"Your answer"}},
		{"signoff", full, 1, []string{"Sign off", "a approve"}},
		{"failed", full, 2, []string{"Retry", "Abandon"}},
		{"accept", full, 5, []string{"Accept this change set?", "Seven"}},
		{"open", gate, 0, []string{"which-db", "how-many", "Four"}},
	} {
		var calls []call
		m := NewModel(fake(c.out, &calls))
		for range c.downs {
			press(m, "down")
		}
		enter(m)
		v := m.View().Content
		for _, w := range c.want {
			if !strings.Contains(v, w) {
				t.Errorf("%s: action view lacks %q:\n%s", c.name, w, v)
			}
		}
		if len(calls) != 1 {
			t.Errorf("%s: runner called before the action: %+v", c.name, calls)
		}
	}
}

func TestAcceptConfirmedCallsAccept(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	for range 5 {
		press(m, "down")
	}
	enter(m)
	if p, _ := m.Picked(); p.Action != ActionAccept || p.Item.ID != "j1" {
		t.Fatalf("picked %+v", p)
	}
	m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if !isQuit(submit(m)) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "accept" || !slices.Equal(last.args, []string{"C7"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestAcceptDeclinedCallsNothing(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	for range 5 {
		press(m, "down")
	}
	enter(m)
	if !isQuit(submit(m)) || len(calls) != 1 {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestGateAnswersEachThenResumes(t *testing.T) {
	var calls []call
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","detail":"open questions: which-db, how-many; answer each with a Decision, then rhiza run resume R4"}]}`
	m := NewModel(fake(gate, &calls))
	enter(m)
	if isQuit(typeAndSubmit(m, "sqlite")) {
		t.Fatal("quit after the first answer")
	}
	if !isQuit(typeAndSubmit(m, "three")) {
		t.Fatal("app did not quit")
	}
	want := []call{
		{Name, nil},
		{"question.answer", []string{"which-db", "sqlite", "--run", "R4"}},
		{"question.answer", []string{"how-many", "three", "--run", "R4"}},
		{"run.resume", []string{"R4"}},
	}
	if len(calls) != len(want) {
		t.Fatalf("calls: %+v", calls)
	}
	for i := range want {
		if calls[i].op != want[i].op || !slices.Equal(calls[i].args, want[i].args) {
			t.Fatalf("call %d = %+v, want %+v", i, calls[i], want[i])
		}
	}
}

func TestGateResumeFailureShowsNote(t *testing.T) {
	var calls []call
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","detail":"open questions: which-db; answer each with a Decision, then rhiza run resume R4"}]}`
	run := fake(gate, &calls)
	m := NewModel(func(op string, args []string) ([]byte, error) {
		if op == "run.resume" {
			calls = append(calls, call{op, args})
			return nil, errors.New("run is busy")
		}
		return run(op, args)
	})
	enter(m)
	if isQuit(typeAndSubmit(m, "sqlite")) {
		t.Fatal("quit though the resume failed")
	}
	if want := "answered; resume failed: run is busy"; m.Note() != want || !strings.Contains(m.View().Content, want) {
		t.Fatalf("note %q, view:\n%s", m.Note(), m.View().Content)
	}
	if last := calls[len(calls)-1]; last.op != "run.resume" {
		t.Fatalf("called %s after the failed resume: %+v", last.op, calls)
	}
}

func TestGateBlankAnswerSendsNothing(t *testing.T) {
	var calls []call
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","detail":"open questions: which-db; answer each with a Decision, then rhiza run resume R4"}]}`
	m := NewModel(fake(gate, &calls))
	enter(m)
	if isQuit(typeAndSubmit(m, "")) || len(calls) != 1 {
		t.Fatalf("blank answer went through: %+v", calls)
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

func TestAnswerEmptyThenAnswerSends(t *testing.T) {
	var calls []call
	m := NewModel(fake(full, &calls))
	enter(m)
	if isQuit(typeAndSubmit(m, "")) || len(calls) != 1 {
		t.Fatalf("empty answer went through: %+v", calls)
	}
	if !strings.Contains(m.View().Content, "an answer needs words") {
		t.Fatalf("no refusal shown:\n%s", m.View().Content)
	}
	if !isQuit(typeAndSubmit(m, "Use sqlite")) {
		t.Fatal("answer after an empty Enter did not go through")
	}
	last := calls[len(calls)-1]
	if last.op != "session.message" || !slices.Equal(last.args, []string{"S1", "Use sqlite"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestAnswerSendFailureReopensForm(t *testing.T) {
	var calls []call
	fail := true
	m := NewModel(func(op string, args []string) ([]byte, error) {
		calls = append(calls, call{op, args})
		if op == Name {
			return []byte(full), nil
		}
		if fail {
			fail = false
			return nil, errors.New("no session")
		}
		return []byte("{}"), nil
	})
	enter(m)
	if isQuit(typeAndSubmit(m, "Use sqlite")) {
		t.Fatal("failed send quit")
	}
	if !strings.Contains(m.View().Content, "send failed: no session") {
		t.Fatalf("no failure shown:\n%s", m.View().Content)
	}
	if !isQuit(typeAndSubmit(m, "Use pg")) {
		t.Fatal("retry after a failed send did not go through")
	}
	last := calls[len(calls)-1]
	if len(calls) != 3 || last.op != "session.message" || !slices.Equal(last.args, []string{"S1", "Use pg"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func signoffModel(calls *[]call) *Model {
	m := NewModel(fake(full, calls))
	press(m, "down")
	enter(m)
	return m
}

func key(m *Model, s string) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: rune(s[0]), Text: s})
	return cmd
}

func TestSignoffDialogApprove(t *testing.T) {
	var calls []call
	m := signoffModel(&calls)
	v := m.View().Content
	for _, w := range []string{"Two", "R2", "signoff-wait"} {
		if !strings.Contains(v, w) {
			t.Errorf("dialog lacks %q:\n%s", w, v)
		}
	}
	if !isQuit(key(m, "a")) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "run.approve" || !slices.Equal(last.args, []string{"R2"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestSignoffSendBackNote(t *testing.T) {
	var calls []call
	m := signoffModel(&calls)
	key(m, "s")
	if isQuit(typeAndSubmit(m, "")) || len(calls) != 1 {
		t.Fatalf("empty note went through: %+v", calls)
	}
	if !isQuit(typeAndSubmit(m, "fix the tests")) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "run.send-back" || !slices.Equal(last.args, []string{"R2", "--note", "fix the tests"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestSignoffEscCallsNothing(t *testing.T) {
	var calls []call
	m := signoffModel(&calls)
	if !isQuit(press(m, "esc")) {
		t.Fatal("esc did not quit")
	}
	if len(calls) != 1 {
		t.Fatalf("calls: %+v", calls)
	}
}

func resolveModel(calls *[]call) *Model {
	m := NewModel(fake(full, calls))
	press(m, "down")
	press(m, "down")
	enter(m)
	return m
}

// submit presses Enter in the open form and feeds its message back.
func submit(m *Model) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if _, ok := msg.(ui.FormSubmittedMsg); !ok {
		return nil
	}
	_, c := m.Update(msg)
	return c
}

func TestResolveRetryCallsRunResume(t *testing.T) {
	var calls []call
	m := resolveModel(&calls)
	if v := m.View().Content; !strings.Contains(v, "Retry") || !strings.Contains(v, "Abandon") {
		t.Fatalf("no choice:\n%s", v)
	}
	if !isQuit(submit(m)) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "run.resume" || !slices.Equal(last.args, []string{"--retry", "R3"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestResolveAbandonNeedsConfirm(t *testing.T) {
	var calls []call
	m := resolveModel(&calls)
	press(m, "down")
	if isQuit(submit(m)) {
		t.Fatal("quit before the confirm")
	}
	if isQuit(submit(m)) || len(calls) != 1 {
		t.Fatalf("abandon without confirm went through: %+v", calls)
	}
}

func TestResolveAbandonConfirmed(t *testing.T) {
	var calls []call
	m := resolveModel(&calls)
	press(m, "down")
	submit(m)
	m.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
	if !isQuit(submit(m)) {
		t.Fatal("app did not quit")
	}
	last := calls[len(calls)-1]
	if last.op != "run.abandon" || !slices.Equal(last.args, []string{"R3"}) {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestResolveEscCallsNothing(t *testing.T) {
	var calls []call
	m := resolveModel(&calls)
	_, cmd := m.Update(ui.FormCancelledMsg{})
	if !isQuit(cmd) || len(calls) != 1 {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestResolveOnlyForFailedKinds(t *testing.T) {
	for _, k := range []string{api.InboxQuestion, api.InboxSignoff, api.InboxAccept, api.InboxStalled} {
		if ActionFor(k) == ActionResolve {
			t.Errorf("%s offers retry/abandon", k)
		}
	}
	for _, k := range []string{api.InboxFailed, api.InboxGaveUp, api.InboxSkipped} {
		if ActionFor(k) != ActionResolve {
			t.Errorf("%s does not offer retry/abandon", k)
		}
	}
}

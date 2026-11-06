package inboxapp

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/dphaener/rexuikit/ui"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

const full = `{"items":[
{"id":"e5","kind":"question","change":"C1","change_title":"One","step":"build","run":"R1","session":"S1","question":"Which db?"},
{"id":"d4","kind":"signoff","change":"C2","change_title":"Two","step":"signoff","run":"R2","detail":"signoff-wait"},
{"id":"c3","kind":"failed","change":"C3","change_title":"Three","step":"verify","run":"R3","detail":"boom"},
{"id":"h8","kind":"stalled","change":"C9","change_title":"Nine","step":"build","run":"R9","session":"S1","detail":"the session is silent"},
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
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[{"handle":"which-db","statement":"Which db?"},{"handle":"how-many","statement":"How many?"}]}]}`
	for _, c := range []struct {
		name, out string
		downs     int
		want      []string
	}{
		{"question", full, 0, []string{"Your answer"}},
		{"signoff", full, 1, []string{"Sign off", "a approve"}},
		{"failed", full, 2, []string{"Retry", "Abandon"}},
		{"accept", full, 5, []string{"Accept this change set?", "Seven"}},
		{"open", gate, 0, []string{"Which db?", "How many?", "Four"}},
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
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[{"handle":"which-db","statement":"Which db?"},{"handle":"how-many","statement":"How many?"}]}]}`
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
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[{"handle":"which-db","statement":"Which db?"}]}]}`
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
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[{"handle":"which-db","statement":"Which db?"}]}]}`
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
	for _, w := range []string{"stalled", "Nine", "the session is silent"} {
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
	for _, w := range []string{"Two", "signoff-wait"} {
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

const longText = "Should the web UI keep the Operation.Write method on the shared registry, or move writes behind a separate interface so read-only callers can never reach them, given that the review found three callers that already depend on it and the daemon will need a stricter surface later when the owner turns autonomy on for the first time? Answer yes to keep it, no to split it, or say what you would rather do instead."

// words are the words of s.
func words(s string) []string { return strings.Fields(s) }

func TestBlockerShowsEveryWordWrapped(t *testing.T) {
	if len(longText) < 300 {
		t.Fatalf("fixture is %d chars", len(longText))
	}
	out := `{"items":[{"id":"b1","kind":"blocker","change":"C1","change_title":"Web UI","question":` + quote(longText) + `,"handle":"web-ui-write","blocks":"this change set"}]}`
	m := NewModel(fake(out, new([]call)))
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	v := plain(m.View().Content)
	for _, w := range words(longText) {
		if !strings.Contains(v, w) {
			t.Fatalf("view lacks %q:\n%s", w, v)
		}
	}
	for _, l := range strings.Split(v, "\n") {
		if lipgloss.Width(l) > 70 {
			t.Errorf("line wider than the modal: %q", l)
		}
	}
	for _, w := range []string{"Blocks: Web UI", "Your answer"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
}

// A question's answer form shows its full statement, then its
// recommendation, above the TextArea.
func TestQuestionFormShowsStatementAndRecommendationAboveAnswer(t *testing.T) {
	for _, kind := range []string{"question", "blocker"} {
		out := `{"items":[{"id":"q1","kind":"` + kind + `","change":"C1","change_title":"Web UI","run":"R1","session":"S1","question":` + quote(longText) + `,"handle":"web-ui-write","recommendation":"Split it behind a read-only interface"}]}`
		m := NewModel(fake(out, new([]call)))
		enter(m)
		m.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
		v := strings.Join(strings.Fields(plain(m.View().Content)), " ")
		s, r, a := strings.Index(v, "Should the web UI"), strings.Index(v, "Recommended: Split it behind a read-only interface"), strings.Index(v, "Your answer")
		if s < 0 || r < 0 || a < 0 || !(s < r && r < a) {
			t.Errorf("%s: statement at %d, recommendation at %d, answer at %d:\n%s", kind, s, r, a, v)
		}
		for _, w := range words(longText) {
			if !strings.Contains(v, w) {
				t.Fatalf("%s: view lacks %q:\n%s", kind, w, v)
			}
		}
	}
}

// A form opened after the window's size wraps to it, and its natural
// width does not follow the question.
func TestAnswerFormOpenedAfterResizeWraps(t *testing.T) {
	out := `{"items":[{"id":"b1","kind":"blocker","change":"C1","change_title":"Web UI","question":` + quote(longText) + `,"handle":"web-ui-write"}]}`
	m := NewModel(fake(out, new([]call)))
	m.Update(tea.WindowSizeMsg{Width: 60, Height: 40})
	enter(m)
	v := plain(m.View().Content)
	for _, l := range strings.Split(v, "\n") {
		if lipgloss.Width(l) > 60 {
			t.Errorf("line wider than the pane: %q", l)
		}
	}
	for _, w := range words(longText) {
		if !strings.Contains(v, w) {
			t.Fatalf("view lacks %q:\n%s", w, v)
		}
	}
	if w, _ := m.NaturalSize(); w > qvMinWidth {
		t.Errorf("natural width %d follows the question", w)
	}
}

func quote(s string) string { b, _ := json.Marshal(s); return string(b) }

func plain(s string) string { return ansi.Strip(s) }

func TestGateRowShowsQuestionsAndWhatTheyBlock(t *testing.T) {
	out := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","detail":"open questions: which-db; answer each",` +
		`"questions":[{"handle":"which-db","statement":"Which db, formally?","restatement":"Which database should we use?","blocks":"Pick storage","recommendation":"SQLite","reasons":"Single file"},` +
		`{"handle":"how-many","statement":"How many replicas?","blocks":"Scale out"}]}]}`
	m := NewModel(fake(out, new([]call)))
	if r := m.Rows()[0]; !strings.Contains(r, "Which database should we use?") || strings.Contains(r, "formally") {
		t.Errorf("row = %q", r)
	}
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	v := plain(m.View().Content)
	for _, w := range []string{"Which database should we use?", "Blocks: Pick storage", "Recommended: SQLite", "Why: Single file", "How many replicas?", "Blocks: Scale out", "Your answer to question 2 of 2"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
	if strings.Contains(v, "formally") {
		t.Errorf("view shows the statement beside its restatement:\n%s", v)
	}
}

func TestBlockerShowsAdvisorRestatement(t *testing.T) {
	out := `{"items":[{"id":"b2","kind":"blocker","change":"C5","change_title":"Five","question":"Leaf leaf-a failed after 3 attempts: f1",` +
		`"restatement":"Should exports keep the old path?","recommendation":"Keep it","reasons":"Callers use it","handle":"leaf-q","blocks":"this change set"}]}`
	m := NewModel(fake(out, new([]call)))
	if r := m.Rows()[0]; !strings.Contains(r, "Should exports keep the old path?") || strings.Contains(r, "attempts") {
		t.Errorf("row = %q", r)
	}
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	v := plain(m.View().Content)
	for _, w := range []string{"Should exports keep the old path?", "Blocks: Five", "Recommended: Keep it", "Why: Callers use it"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
	if strings.Contains(v, "attempts") {
		t.Errorf("view shows the template text beside its restatement:\n%s", v)
	}
}

func TestScrubKeepsWordsAndDropsOnlyWholeIDs(t *testing.T) {
	out := `{"items":[{"id":"b3","kind":"blocker","change":"C6","change_title":"Six","handle":"disk-q",` +
		`"question":"Which disk should hold the store when the machine sleeps (see disk-q, not disk-qs)?","blocks":"store","blocks_title":"Kept store reachable"}]}`
	m := NewModel(fake(out, new([]call)))
	if r := m.Rows()[0]; !strings.Contains(r, "hold the store when") || strings.Contains(r, "disk-q,") {
		t.Errorf("row = %q", r)
	}
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	v := plain(m.View().Content)
	for _, w := range []string{"hold the store when", "(see …, not disk-qs)?", "Blocks: Kept store reachable"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
}

func TestOverlapRowNamesBothChangeSetsAndTarget(t *testing.T) {
	out := `{"items":[{"id":"o1","kind":"overlap","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"Web UI","other":"01M6PN1SW","other_title":"Exports",` +
		`"target":"export-path","target_title":"Nightly exports go to one folder"}]}`
	m := NewModel(fake(out, new([]call)))
	want := "overlap · Web UI  writes the same part of the plan as Exports: Nightly exports go to one folder"
	if r := m.Rows()[0]; r != want {
		t.Errorf("row = %q, want %q", r, want)
	}
}

// idsIn are the handles and IDs the items carry.
var idsIn = []string{"01M6PF30ABCDEFGHJKMNPQRSTV", "review-01m6pf30-r3-f1", "RUN01M6PF30", "SES01M6PF30", "i7x9k2f", "gate-q-handle", "01M6PN1SW"}

const everyKind = `{"items":[
{"id":"i7x9k2f","kind":"question","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"One","step":"build","run":"RUN01M6PF30","session":"SES01M6PF30","question":"Which db?","tab_block":"01M6PN1SW","recommendation":"SQLite","reasons":"Simple"},
{"id":"i7x9k2f","kind":"blocker","change":"01M6PF30ABCDEFGHJKMNPQRSTV","question":"Keep it?","handle":"review-01m6pf30-r3-f1","blocks":"review-01m6pf30-r3-f1"},
{"id":"i7x9k2f","kind":"open","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"Gate","run":"RUN01M6PF30","detail":"open questions: gate-q-handle; answer each","questions":[{"handle":"gate-q-handle","statement":"Keep it?","blocks":"The plan"}]},
{"id":"i7x9k2f","kind":"open","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"Gate no worktree","run":"RUN01M6PF30","detail":"open questions: gate-q-handle; answer each"},
{"id":"i7x9k2f","kind":"signoff","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"Two","run":"RUN01M6PF30","step":"signoff","detail":"waits on RUN01M6PF30"},
{"id":"i7x9k2f","kind":"failed","change":"01M6PF30ABCDEFGHJKMNPQRSTV","run":"RUN01M6PF30","step":"verify","detail":"boom in 01M6PF30ABCDEFGHJKMNPQRSTV"},
{"id":"i7x9k2f","kind":"stalled","change":"01M6PF30ABCDEFGHJKMNPQRSTV","run":"RUN01M6PF30","session":"SES01M6PF30","detail":"session SES01M6PF30 is silent"},
{"id":"i7x9k2f","kind":"overlap","change":"01M6PF30ABCDEFGHJKMNPQRSTV","other":"01M6PN1SW","target":"gate-q-handle"},
{"id":"i7x9k2f","kind":"accept","change":"01M6PF30ABCDEFGHJKMNPQRSTV","change_title":"Seven","detail":"ready"}
]}`

func TestNoSurfaceShowsAnIDOrHandle(t *testing.T) {
	check := func(what, v string) {
		v = plain(v)
		for _, id := range idsIn {
			if strings.Contains(v, id) {
				t.Errorf("%s shows %q:\n%s", what, id, v)
			}
		}
	}
	n := len(NewModel(fake(everyKind, new([]call))).Items())
	for i := range n {
		m := NewModel(fake(everyKind, new([]call)))
		check("rows", strings.Join(m.Rows(), "\n"))
		for range i {
			press(m, "down")
		}
		enter(m)
		m.Update(tea.WindowSizeMsg{Width: 80, Height: 30})
		check("view of item "+m.items[i].Kind, m.View().Content)
		if m.dlg != nil {
			key(m, "s")
			check("send-back form", m.View().Content)
		}
	}
	// the retry/abandon confirm
	m := NewModel(fake(everyKind, new([]call)))
	for range 5 {
		press(m, "down")
	}
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 30})
	check("resolve form", m.View().Content)
}

func TestGateFormShowsEachQuestionNoHandleAndAnswersRightOne(t *testing.T) {
	var calls []call
	gate := `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[` +
		`{"handle":"review-01m6pf30-r3-f1","statement":"Keep Operation.Write?","restatement":"Should we keep the write operation as it is?","recommendation":"Keep it"},` +
		`{"handle":"claude-design-output","statement":"Is the design output final?"}]}]}`
	m := NewModel(fake(gate, &calls))
	enter(m)
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	v := plain(m.View().Content)
	for _, w := range []string{"1 of 2: Should we keep the write operation as it is?", "Recommended: Keep it", "Your answer to question 1 of 2"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
	if strings.Index(v, "Should we keep") > strings.Index(v, "Your answer to question 1") {
		t.Errorf("question is not above its text box:\n%s", v)
	}
	for _, h := range []string{"review-01m6pf30-r3-f1", "claude-design-output"} {
		if strings.Contains(v, h) {
			t.Errorf("view shows handle %q:\n%s", h, v)
		}
	}
	typeAndSubmit(m, "yes")
	v = plain(m.View().Content)
	if !strings.Contains(v, "2 of 2: Is the design output final?") || strings.Contains(v, "claude-design-output") {
		t.Errorf("second question view:\n%s", v)
	}
	typeAndSubmit(m, "no")
	want := [][]string{{"review-01m6pf30-r3-f1", "yes", "--run", "R4"}, {"claude-design-output", "no", "--run", "R4"}}
	var got [][]string
	for _, c := range calls {
		if c.op == "question.answer" {
			got = append(got, c.args)
		}
	}
	if len(got) != 2 || !slices.Equal(got[0], want[0]) || !slices.Equal(got[1], want[1]) {
		t.Errorf("answers = %v", got)
	}
}

func longGate(n int, statement func(i int) string) string {
	var qs []string
	for i := 1; i <= n; i++ {
		qs = append(qs, `{"handle":"q-handle-`+string(rune('a'+i))+`","statement":`+quote(statement(i))+`}`)
	}
	return `{"items":[{"id":"g1","kind":"open","change":"C4","change_title":"Four","run":"R4","questions":[` + strings.Join(qs, ",") + `]}]}`
}

func scrollText(i int) string {
	return strings.Repeat("word"+string(rune('0'+i))+" ", 80) + "END" + string(rune('0'+i))
}

func TestGateFormScrollsLongQuestions(t *testing.T) {
	var calls []call
	m := NewModel(fake(longGate(4, scrollText), &calls))
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 30})
	enter(m)
	v := plain(m.View().Content)
	if strings.Contains(v, "END4") || !strings.Contains(v, "Your answer to question 1 of 4") {
		t.Fatalf("first view:\n%s", v)
	}
	if got := lipgloss.Height(m.View().Content); got > 30 {
		t.Fatalf("view is %d rows in a 30-row window", got)
	}
	for range 40 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	}
	v = plain(m.View().Content)
	if !strings.Contains(v, "END4") || !strings.Contains(v, "Your answer to question 1 of 4") {
		t.Fatalf("after page-down:\n%s", v)
	}
	for range 40 {
		m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
	}
	if v = plain(m.View().Content); strings.Contains(v, "END4") || !strings.Contains(v, "1-") {
		t.Fatalf("after wheel up:\n%s", v)
	}
	for range 40 {
		m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
	}
	if v = plain(m.View().Content); !strings.Contains(v, "END4") {
		t.Fatalf("after wheel down:\n%s", v)
	}
	if m.View().MouseMode == tea.MouseModeNone {
		t.Error("mouse is off, so the wheel never arrives")
	}
	// the TextAreas still take the keys: four answers submit
	for _, a := range []string{"a", "b", "c"} {
		if isQuit(typeAndSubmit(m, a)) {
			t.Fatal("quit early")
		}
	}
	if !isQuit(typeAndSubmit(m, "d")) {
		t.Fatalf("no quit; calls: %+v", calls)
	}
	if n := len(calls); n != 1+4+1 || calls[n-1].op != "run.resume" {
		t.Fatalf("calls: %+v", calls)
	}
}

func TestGateUnreadableQuestionShowsPositionNotHandle(t *testing.T) {
	var calls []call
	gate := longGate(4, func(i int) string {
		if i == 2 {
			return ""
		}
		return "Question text " + string(rune('0'+i)) + "?"
	})
	m := NewModel(fake(gate, &calls))
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	enter(m)
	v := plain(m.View().Content)
	for _, w := range []string{"question 2 of 4", "This question's text could not be read.", "Your answer to question 2 of 4"} {
		if !strings.Contains(v, w) {
			t.Errorf("view lacks %q:\n%s", w, v)
		}
	}
	if strings.Contains(v, "q-handle") {
		t.Errorf("view shows a handle:\n%s", v)
	}
}

func TestGateAnswerFailureReopensOnRemaining(t *testing.T) {
	var calls []call
	run := fake(longGate(3, func(i int) string { return "Q" + string(rune('0'+i)) + "?" }), &calls)
	m := NewModel(func(op string, args []string) ([]byte, error) {
		if op == "question.answer" && args[1] == "two" {
			calls = append(calls, call{op, args})
			return nil, errors.New("disk full")
		}
		return run(op, args)
	})
	enter(m)
	typeAndSubmit(m, "one")
	if isQuit(typeAndSubmit(m, "two")) {
		t.Fatal("quit")
	}
	typeAndSubmit(m, "three")
	v := plain(m.View().Content)
	if !strings.Contains(v, "answer failed: disk full") || strings.Contains(v, "Q1?") {
		t.Fatalf("view:\n%s", v)
	}
	for _, c := range calls {
		if c.op == "run.resume" {
			t.Fatalf("resumed: %+v", calls)
		}
	}
}

func TestDetailPanelScrolls(t *testing.T) {
	var calls []call
	out := `{"items":[{"id":"h8","kind":"stalled","change":"C9","change_title":"Nine","run":"R9","detail":` + quote(strings.Repeat("filler ", 400)+"LASTWORD") + `}]}`
	m := NewModel(fake(out, &calls))
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 20})
	enter(m)
	if v := plain(m.View().Content); strings.Contains(v, "LASTWORD") {
		t.Fatalf("detail not clipped to the window:\n%s", v)
	}
	for range 30 {
		m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	}
	if v := plain(m.View().Content); !strings.Contains(v, "LASTWORD") {
		t.Fatalf("page-down did not reach the end:\n%s", v)
	}
	for range 30 {
		m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
	}
	if v := plain(m.View().Content); strings.Contains(v, "LASTWORD") {
		t.Fatalf("wheel up did not scroll back:\n%s", v)
	}
}

func TestGateResumeFailureDetailScrolls(t *testing.T) {
	var calls []call
	gate := longGate(1, func(int) string { return "Q?" })
	run := fake(gate, &calls)
	m := NewModel(func(op string, args []string) ([]byte, error) {
		if op == "run.resume" {
			return nil, errors.New("busy")
		}
		return run(op, args)
	})
	m.Update(tea.WindowSizeMsg{Width: 70, Height: 12})
	enter(m)
	typeAndSubmit(m, "x")
	if v := plain(m.View().Content); !strings.Contains(v, "resume failed: busy") {
		t.Fatalf("view:\n%s", v)
	}
	if _, ok := m.cur.(*ui.InfoPanel); !ok {
		t.Fatalf("not a scrolling panel: %T", m.cur)
	}
}

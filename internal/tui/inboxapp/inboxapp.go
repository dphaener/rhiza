// Package inboxapp is the inbox app: the items of `rhiza inbox` as one
// rexuikit List in the modal shell, a row per item in the order the
// operation returns them. It runs the api operation in-process through the
// registry, so --once prints what `rhiza inbox --json` prints.
//
// Enter picks an item. A question, gate, sign-off, failed run or accept
// item opens that item's action in the app; any other kind shows its
// detail and asks the runner to open its run's tab.
package inboxapp

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"github.com/dphaener/rexuikit/ui"
	"github.com/dphaener/rexuikit/ui/theme"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

// Name is the app's name, and the api operation it runs.
const Name = "inbox"

func init() {
	launcher.Register(launcher.App{Name: Name, Summary: "What waits on you: questions, sign-offs and failed runs"})
}

// Runner runs an operation and returns its result as `rhiza <op> --json`
// prints it.
type Runner func(op string, args []string) ([]byte, error)

// Action is what picking an item offers to do.
type Action string

const (
	// ActionAnswer answers a question, a gate or a blocker.
	ActionAnswer Action = "answer"
	// ActionSignoff signs a run off or sends it back.
	ActionSignoff Action = "signoff"
	// ActionResolve retries or abandons a failed, gave-up or skipped run.
	ActionResolve Action = "resolve"
	// ActionAccept accepts a change set.
	ActionAccept Action = "accept"
	// ActionShow shows the item's detail and opens its run's tab.
	ActionShow Action = "show"
)

// ActionFor is the action picking an item of kind offers.
func ActionFor(kind string) Action {
	switch kind {
	case api.InboxQuestion, api.InboxOpen, api.InboxBlocker:
		return ActionAnswer
	case api.InboxSignoff:
		return ActionSignoff
	case api.InboxFailed, api.InboxGaveUp, api.InboxSkipped:
		return ActionResolve
	case api.InboxAccept:
		return ActionAccept
	}
	return ActionShow
}

// Pick is an item picked with Enter and the action it opens.
type Pick struct {
	Item   api.InboxItem
	Action Action
}

// emptyLine is the one line an empty inbox shows.
const emptyLine = "Nothing waits on you"

// component is what the model currently shows.
type component interface {
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
	NaturalSize() (int, int)
	Theme() theme.Theme
}

// Model is the app's bubbletea program. Esc quits.
type Model struct {
	run   Runner
	out   []byte
	err   error
	items []api.InboxItem // in the operation's order
	list  *ui.List
	cur   component
	pick  *Pick
	note  string
	form  *ui.Form
	asked *api.InboxItem
	dlg   *ui.Dialog
	back  bool // the form asks for a send-back note
	// stage is "choose" or "confirm" while the retry/abandon form is up,
	// "gate" while a gate's answer form is up, "accept" while the accept
	// confirm is up.
	stage  string
	gateQs []string // the gate's open questions, by handle, in form order
}

// NewModel returns the model over run, loaded once.
func NewModel(run Runner) *Model {
	m := &Model{run: run}
	m.load()
	return m
}

func (m *Model) load() {
	out, err := m.run(Name, nil)
	m.out, m.err = out, err
	var rows []ui.Row
	if err == nil {
		var resp api.InboxResponse
		if err := json.Unmarshal(out, &resp); err != nil {
			m.err = err
		} else {
			m.items = resp.Items
		}
	}
	for _, it := range m.items {
		rows = append(rows, row(it))
	}
	m.list = ui.NewList(theme.Dark(), colorprofile.TrueColor, rows)
	m.cur = m.list
}

// Result is what --once prints: the operation's output, byte for byte, or
// the error it ran into.
func (m *Model) Result() ([]byte, error) { return m.out, m.err }

// Items are the inbox items in row order.
func (m *Model) Items() []api.InboxItem { return m.items }

// Picked is the item Enter picked and its action, once one was.
func (m *Model) Picked() (Pick, bool) {
	if m.pick == nil {
		return Pick{}, false
	}
	return *m.pick, true
}

// Note is the outcome of the last pick that was not an action: an error
// from opening a tab, or empty.
func (m *Model) Note() string { return m.note }

// Rows is the text of the rows as the list shows them.
func (m *Model) Rows() []string {
	var out []string
	for _, it := range m.items {
		r := row(it)
		out = append(out, strings.TrimSpace(r.Label+"  "+r.Detail))
	}
	return out
}

// Theme is the shown component's theme.
func (m *Model) Theme() theme.Theme { return m.cur.Theme() }

// NaturalSize is the shown component's size in cells.
func (m *Model) NaturalSize() (int, int) { return m.cur.NaturalSize() }

func (m *Model) Init() tea.Cmd { return tea.RequestBackgroundColor }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.ListCancelledMsg, ui.InfoClosedMsg:
		return m, tea.Quit
	case ui.ListSelectedMsg:
		return m, m.choose(msg.Index)
	case ui.FormCancelledMsg:
		return m, tea.Quit
	case ui.FormSubmittedMsg:
		if m.stage == "choose" {
			return m, m.chosen(msg.Values.String("choice"))
		}
		if m.stage == "confirm" {
			return m, m.abandon(msg.Values.Bool("confirm"))
		}
		if m.stage == "gate" {
			return m, m.answerGate(msg.Values)
		}
		if m.stage == "accept" {
			return m, m.accept(msg.Values.Bool("accept"))
		}
		if m.back {
			return m, m.sendBack(msg.Values.String("note"))
		}
		return m, m.send(msg.Values.String("answer"))
	case tea.KeyPressMsg:
		if m.dlg != nil && m.cur == component(dialogView{m}) {
			switch msg.String() {
			case "a":
				return m, m.approve()
			case "s":
				return m, m.openNote("send back")
			}
		}
		if k := msg.String(); k == "ctrl+c" || (k == "esc" && m.cur != component(m.list)) {
			return m, tea.Quit
		}
	}
	_, cmd := m.cur.Update(msg)
	return m, cmd
}

// choose handles Enter on item i.
func (m *Model) choose(i int) tea.Cmd {
	if i < 0 || i >= len(m.items) {
		return nil
	}
	it := m.items[i]
	p := Pick{Item: it, Action: ActionFor(it.Kind)}
	m.pick = &p
	if it.Kind == api.InboxQuestion || it.Kind == api.InboxBlocker {
		return m.openAnswer(it)
	}
	if p.Action == ActionSignoff {
		return m.openSignoff(it)
	}
	if p.Action == ActionResolve {
		return m.openResolve(it)
	}
	if it.Kind == api.InboxOpen {
		if qs := gateQuestions(it.Detail); len(qs) > 0 {
			return m.openGate(it, qs, "")
		}
	}
	if p.Action == ActionAccept {
		return m.openAccept(it)
	}
	return m.show(it)
}

// show shows item it's detail and asks the runner to open its run's tab.
func (m *Model) show(it api.InboxItem) tea.Cmd {
	m.note = ""
	if it.Run != "" {
		if _, err := m.run(Name+".open", []string{it.ID}); err != nil {
			m.note = fmt.Sprintf("open failed: %v", err)
		}
	}
	m.cur = ui.NewInfoPanel(theme.Dark(), colorprofile.TrueColor, title(it), func(int) []string {
		return detail(it, m.note)
	})
	return m.cur.(*ui.InfoPanel).Init()
}

// gateQuestions are the handles of the open questions a gate stopped on,
// from the stop detail "open questions: a, b; answer each ...".
func gateQuestions(detail string) []string {
	rest, ok := strings.CutPrefix(detail, run.OpenStop)
	if !ok {
		return nil
	}
	rest, _, _ = strings.Cut(rest, ";")
	var qs []string
	for _, q := range strings.Split(rest, ",") {
		if q = strings.TrimSpace(q); q != "" {
			qs = append(qs, q)
		}
	}
	return qs
}

// openGate shows a form with one TextArea per open question of the gate.
func (m *Model) openGate(it api.InboxItem, qs []string, msg string) tea.Cmd {
	fields := make([]ui.Field, len(qs))
	for i, q := range qs {
		fields[i] = ui.NewTextArea(fmt.Sprintf("q%d", i), q)
	}
	sub := title(it) + " · answer each, then the run resumes"
	if msg != "" {
		sub = title(it) + " · " + msg
	}
	f := ui.NewForm(theme.Dark(), colorprofile.TrueColor, fields...).WithSubtitle(sub)
	m.form, m.asked, m.stage, m.gateQs, m.cur = f, &it, "gate", qs, f
	return f.Init()
}

// answerGate answers each open question through question.answer into the
// run's change set, then queues the run's continuation with run.resume and
// quits. A blank answer sends nothing and asks again.
func (m *Model) answerGate(v ui.Values) tea.Cmd {
	it, qs := *m.asked, m.gateQs
	texts := make([]string, len(qs))
	for i := range qs {
		texts[i] = strings.TrimSpace(v.String(fmt.Sprintf("q%d", i)))
		if texts[i] == "" {
			return m.openGate(it, qs, "answer each question")
		}
	}
	for i, q := range qs {
		if _, err := m.run("question.answer", []string{q, it.Change, texts[i]}); err != nil {
			return m.openGate(it, qs[i:], "answer failed: "+err.Error())
		}
	}
	if _, err := m.run("run.resume", []string{it.Run}); err != nil {
		m.note = "answered; resume failed: " + err.Error()
		return m.show(it)
	}
	return tea.Quit
}

// openAccept asks to confirm accepting the item's change set.
func (m *Model) openAccept(it api.InboxItem) tea.Cmd {
	f := ui.NewForm(theme.Dark(), colorprofile.TrueColor,
		ui.NewConfirm("accept", "Accept this change set?")).WithSubtitle(title(it) + " · " + waits(it))
	m.form, m.asked, m.stage, m.cur = f, &it, "accept", f
	return f.Init()
}

// accept accepts the change set through accept (rule owner) once
// confirmed; without the Confirm nothing is called and the app quits.
func (m *Model) accept(ok bool) tea.Cmd {
	if !ok {
		return tea.Quit
	}
	it := *m.asked
	if _, err := m.run("accept", []string{it.Change}); err != nil {
		f := ui.NewForm(theme.Dark(), colorprofile.TrueColor,
			ui.NewConfirm("accept", "Accept this change set?")).WithSubtitle(title(it) + " · accept failed: " + err.Error())
		m.form, m.cur = f, f
		return f.Init()
	}
	return tea.Quit
}

// openAnswer shows a form with one TextArea for the owner's words.
func (m *Model) openAnswer(it api.InboxItem) tea.Cmd {
	return m.answerForm(it, waits(it))
}

// answerForm opens a fresh answer form whose subtitle carries msg. A
// submitted Form takes no more keys, so a refused or failed send reopens
// one rather than reusing it.
func (m *Model) answerForm(it api.InboxItem, msg string) tea.Cmd {
	f := ui.NewForm(theme.Dark(), colorprofile.TrueColor,
		ui.NewTextArea("answer", "Your answer")).WithSubtitle(title(it) + " · " + msg)
	m.form, m.asked = f, &it
	m.cur = f
	return f.Init()
}

// dialogView adapts the Dialog to the component the model shows.
type dialogView struct{ m *Model }

func (d dialogView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d.m, d.m.dlg.Update(msg)
}
func (d dialogView) View() tea.View          { return tea.NewView(d.m.dlg.View()) }
func (d dialogView) NaturalSize() (int, int) { return d.m.dlg.NaturalSize() }
func (d dialogView) Theme() theme.Theme      { return d.m.dlg.Theme() }

// openSignoff shows the dialog naming the run and what it waits for.
func (m *Model) openSignoff(it api.InboxItem) tea.Cmd {
	body := []string{"run: " + it.Run, "waits for: " + waits(it)}
	m.asked = &it
	m.dlg = ui.NewDialog(theme.Dark(), colorprofile.TrueColor, "Sign off · "+title(it), body).
		WithHint("a approve · s send back · esc close")
	m.cur = dialogView{m}
	return m.dlg.Init()
}

// approve signs the run off through run.approve and quits.
func (m *Model) approve() tea.Cmd {
	if _, err := m.run("run.approve", []string{m.asked.Run}); err != nil {
		m.dlg.SetBody([]string{"run: " + m.asked.Run, "approve failed: " + err.Error()})
		return nil
	}
	return tea.Quit
}

// openNote asks for the note a send-back carries.
func (m *Model) openNote(msg string) tea.Cmd {
	f := ui.NewForm(theme.Dark(), colorprofile.TrueColor,
		ui.NewTextArea("note", "What to change")).WithSubtitle(title(*m.asked) + " · " + msg)
	m.form, m.back, m.cur = f, true, f
	return f.Init()
}

// sendBack sends the run back through run.send-back with the note. An empty
// note is refused and leaves the form open.
func (m *Model) sendBack(note string) tea.Cmd {
	note = strings.TrimSpace(note)
	if note == "" {
		return m.openNote("a send-back needs a note")
	}
	if _, err := m.run("run.send-back", []string{m.asked.Run, "--note", note}); err != nil {
		return m.openNote("send back failed: " + err.Error())
	}
	return tea.Quit
}

// openResolve offers retry or abandon for a failed, gave-up or skipped run.
func (m *Model) openResolve(it api.InboxItem) tea.Cmd {
	return m.resolveForm(it, "")
}

func (m *Model) resolveForm(it api.InboxItem, msg string) tea.Cmd {
	sub := title(it) + " · " + waits(it)
	if msg != "" {
		sub = title(it) + " · " + msg
	}
	sel := ui.NewSelect("choice", "What to do with the run").Options(
		ui.Option{Label: "Retry", Value: "retry"},
		ui.Option{Label: "Abandon", Value: "abandon"},
	)
	f := ui.NewForm(theme.Dark(), colorprofile.TrueColor, sel).WithSubtitle(sub)
	m.form, m.asked, m.stage, m.cur = f, &it, "choose", f
	return f.Init()
}

// chosen acts on the Select: retry queues run.resume and quits; abandon
// asks for the Confirm first.
func (m *Model) chosen(choice string) tea.Cmd {
	it := *m.asked
	if choice == "abandon" {
		f := ui.NewForm(theme.Dark(), colorprofile.TrueColor,
			ui.NewConfirm("confirm", "Abandon this run?")).WithSubtitle(title(it) + " · " + it.Run)
		m.form, m.stage, m.cur = f, "confirm", f
		return f.Init()
	}
	if _, err := m.run("run.resume", []string{"--retry", it.Run}); err != nil {
		return m.resolveForm(it, "retry failed: "+err.Error())
	}
	return tea.Quit
}

// abandon runs run.abandon once confirmed; without the Confirm nothing is
// called and the choice comes back.
func (m *Model) abandon(ok bool) tea.Cmd {
	it := *m.asked
	if !ok {
		return m.resolveForm(it, "abandon needs the confirm")
	}
	if _, err := m.run("run.abandon", []string{it.Run}); err != nil {
		return m.resolveForm(it, "abandon failed: "+err.Error())
	}
	return tea.Quit
}

// send sends the typed words: session.message to the asking session for a
// run's question, question.answer for a blocker. An empty answer sends
// nothing and leaves the form open; success closes the app.
func (m *Model) send(text string) tea.Cmd {
	if m.asked == nil {
		return nil
	}
	it := *m.asked
	text = strings.TrimSpace(text)
	if text == "" {
		return m.answerForm(it, "an answer needs words")
	}
	var err error
	if it.Kind == api.InboxBlocker {
		_, err = m.run("question.answer", []string{it.Handle, it.Change, text})
	} else if it.Session == "" {
		err = fmt.Errorf("the question has no session to answer")
	} else {
		_, err = m.run("session.message", []string{it.Session, text})
	}
	if err != nil {
		m.note = fmt.Sprintf("send failed: %v", err)
		return m.answerForm(it, m.note)
	}
	return tea.Quit
}

func (m *Model) View() tea.View {
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("rhiza: %v\n", m.err))
	}
	if len(m.items) == 0 {
		return tea.NewView(emptyLine + "\n")
	}
	return m.cur.View()
}

// title is the change set's title, or what the item is about when it has
// no change set.
func title(it api.InboxItem) string {
	switch it.Kind {
	case api.InboxDiskLow:
		return "Daemon holds starts"
	case api.InboxMainStale:
		return "Main not updated"
	case api.InboxUsageReserve:
		return "Usage reserve"
	}
	if it.ChangeTitle != "" {
		return it.ChangeTitle
	}
	return it.Change
}

// waits is what the item waits for.
func waits(it api.InboxItem) string {
	switch it.Kind {
	case api.InboxQuestion, api.InboxBlocker:
		if it.Question != "" {
			return it.Question
		}
	case api.InboxOverlap:
		return strings.TrimSpace("with " + it.Other + " · " + it.Target + " " + it.Detail)
	}
	return it.Detail
}

// row is one item: its kind and change set title, then what it waits for.
func row(it api.InboxItem) ui.Row {
	label := it.Kind + " · " + title(it)
	if it.Step != "" {
		label += " · " + it.Step
	}
	return ui.Row{Label: label, Detail: waits(it)}
}

// detail is the lines of an item's detail panel.
func detail(it api.InboxItem, note string) []string {
	lines := []string{it.Kind + " · " + title(it)}
	add := func(k, v string) {
		if v != "" {
			lines = append(lines, k+": "+v)
		}
	}
	add("step", it.Step)
	add("since", it.At)
	add("run", it.Run)
	add("waits for", waits(it))
	add("next", it.Next)
	add("", note)
	return lines
}

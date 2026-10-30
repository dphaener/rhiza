package roadmapapp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
)

// answer is the owner typing the answer to a blocking question.
type answer struct {
	cs id.ID
	q  question
	// text is what the owner has typed so far.
	text string
}

func (a *answer) prompt() string {
	return fmt.Sprintf("Answer %q: %s_  (enter records, esc cancels)", a.q.text, a.text)
}

// askAccept shows the what-if of accepting the target and waits for y/n.
func (m *Model) askAccept() {
	cs := m.target()
	if cs.IsZero() {
		return
	}
	m.notice = ""
	if m.preview == nil {
		m.notice = "Cannot accept: no what-if available"
		return
	}
	out, err := m.preview("accept", cs.String(), nil)
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	var r api.LifecycleResponse
	if err := json.Unmarshal(out, &r); err != nil {
		m.notice = fmt.Sprintf("accept: %v", err)
		return
	}
	m.pending = &pending{kind: "accept", id: cs,
		what: fmt.Sprintf("What-if: accept %s, %s -> %s. Accept it? y/n", short(cs), r.From, r.To)}
}

// accept accepts cs and reloads the list.
func (m *Model) accept(cs id.ID) {
	if _, err := m.run("accept", cs.String(), nil); err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.reload()
	m.notice = fmt.Sprintf("Accepted %s", short(cs))
}

// pick moves the chosen question of the open detail by d.
func (m *Model) pick(d int) {
	n := len(m.detail.questions)
	if n == 0 {
		return
	}
	m.detail.picked = min(max(m.detail.picked+d, 0), n-1)
}

// startAnswer begins typing the answer to the chosen question.
func (m *Model) startAnswer() {
	d := m.detail
	if d.err != nil || d.picked >= len(d.questions) {
		return
	}
	m.notice = ""
	m.answer = &answer{cs: d.id, q: d.questions[d.picked]}
}

// typeAnswer handles a key while an answer is being typed.
func (m *Model) typeAnswer(msg tea.KeyPressMsg) tea.Cmd {
	a := m.answer
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		m.answer = nil
	case "enter":
		if strings.TrimSpace(a.text) == "" {
			return nil
		}
		m.answer = nil
		m.record(a)
	case "backspace":
		if r := []rune(a.text); len(r) > 0 {
			a.text = string(r[:len(r)-1])
		}
	default:
		a.text += msg.Text
	}
	return nil
}

// record writes a Decision quoting the typed words, linked answers to the
// question, in the question's change set, then reloads.
func (m *Model) record(a *answer) {
	text := strings.TrimSpace(a.text)
	writes := map[string]any{"writes": []any{
		map[string]any{"op": "add", "kind": "decision", "handle": "answer-" + a.q.handle,
			"fields": []string{
				"statement=" + text,
				fmt.Sprintf("quote=%q (owner)", text),
			}},
		map[string]any{"op": "link", "from": "answer-" + a.q.handle, "type": "answers", "to": a.q.handle},
	}}
	data, err := json.Marshal(writes)
	if err == nil {
		err = m.apply(a.cs, data)
	}
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.reload()
	m.notice = fmt.Sprintf("Recorded your answer to %s", a.q.handle)
}

// apply runs the apply operation on a file of data, in change set cs.
func (m *Model) apply(cs id.ID, data []byte) error {
	f, err := os.CreateTemp("", "rhiza-answer-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return err
	}
	_, err = m.run("apply", cs.String(), []string{f.Name()})
	return err
}

// reload reads the list again and, if a detail is open, that too.
func (m *Model) reload() {
	m.load()
	if m.detail == nil {
		return
	}
	for _, r := range m.rows {
		if r.id == m.detail.id {
			m.detail = m.loadDetail(r)
			return
		}
	}
	m.detail = nil
}

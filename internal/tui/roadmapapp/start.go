package roadmapapp

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
)

// SetPreview gives the model the runner of what-ifs: the same operations
// run as --dry-run does, so nothing is written.
func (m *Model) SetPreview(r Runner) { m.preview = r }

// pending is a run waiting for the owner's y or n.
type pending struct {
	kind string // "plan", "ship" or "accept"
	id   id.ID
	what string // accept only: the what-if the question shows
}

// target is the change set p, s and o act on: the open detail's, else the
// one under the cursor; zero when there is none.
func (m *Model) target() id.ID {
	if m.detail != nil {
		return m.detail.id
	}
	return m.Cursor()
}

// ask puts the y/n question for a run of kind on the target.
func (m *Model) ask(kind string) {
	cs := m.target()
	if cs.IsZero() {
		return
	}
	m.pending = &pending{kind: kind, id: cs}
	m.notice = ""
}

// prompt is the question shown while a run waits for confirmation.
func (p *pending) prompt() string {
	if p.kind == "accept" {
		return p.what
	}
	if p.kind == "plan" {
		return fmt.Sprintf("Queue a plan run of %s for the daemon? y/n", short(p.id))
	}
	return fmt.Sprintf("Queue an autonomous ship run of %s (permission mode auto) for the daemon? y/n", short(p.id))
}

// args are the arguments of the run.start that queues p.
func (p *pending) args() []string {
	args := []string{p.id.String()}
	if p.kind == "plan" {
		return append(args, "--workflow", "plan")
	}
	return append(args, "--autonomous", "--permission-mode", "auto")
}

// intentInput is the owner typing the one-line intent of new work.
type intentInput struct {
	milestone bool
	text      string
}

func (i *intentInput) prompt() string {
	noun := "change set"
	if i.milestone {
		noun = "milestone"
	}
	return fmt.Sprintf("Intent of the new %s: %s_  (enter continues, esc cancels)", noun, i.text)
}

// fileRequest queues a plan request for the typed intent in the clone; the
// daemon turns it into a plan run. No Rex tab is opened here.
func (m *Model) fileRequest(text string, milestone bool) {
	kind := "feature"
	if milestone {
		kind = "milestone"
	}
	if _, err := m.run("request", "", []string{text, "--kind", kind}); err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.notice = fmt.Sprintf("Queued a %s request %q for the daemon to plan", kind, text)
}

// typeIntent handles a key while an intent is being typed. Enter with
// empty text does nothing; enter with text files the request.
func (m *Model) typeIntent(msg tea.KeyPressMsg) tea.Cmd {
	in := m.intent
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		m.intent = nil
	case "enter":
		text := strings.TrimSpace(in.text)
		if text == "" {
			return nil
		}
		m.intent = nil
		m.fileRequest(text, in.milestone)
	case "backspace":
		if r := []rune(in.text); len(r) > 0 {
			in.text = string(r[:len(r)-1])
		}
	default:
		in.text += msg.Text
	}
	return nil
}

// confirm queues the pending run for the daemon.
func (m *Model) confirm() {
	p := m.pending
	m.pending = nil
	if p.kind == "accept" {
		m.accept(p.id)
		return
	}
	if _, err := m.run("run.start", "", p.args()); err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	m.notice = fmt.Sprintf("Queued the %s run of %s for the daemon", p.kind, short(p.id))
}

// openRun focuses the Rex tab of the target's live run, or says there is
// none.
func (m *Model) openRun() {
	cs := m.target()
	if cs.IsZero() {
		return
	}
	out, err := m.run("run.open", cs.String(), []string{cs.String()})
	if err != nil {
		m.notice = fmt.Sprintf("rhiza: %v", err)
		return
	}
	var r api.RunOpenResponse
	if err := json.Unmarshal(out, &r); err != nil {
		m.notice = fmt.Sprintf("run.open: %v", err)
		return
	}
	switch {
	case r.Opened:
		m.notice = fmt.Sprintf("Focused Rex tab %s of run %s", r.Tab, r.Run)
	case !r.Live:
		m.notice = fmt.Sprintf("No live run for %s", short(cs))
	default:
		m.notice = fmt.Sprintf("The tab of run %s is gone", r.Run)
	}
}

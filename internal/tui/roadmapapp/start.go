package roadmapapp

import (
	"encoding/json"
	"fmt"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
)

// Starter is how the app reaches Rex: Exec runs a program and returns its
// output, Self is the rhiza binary a new tab runs, Root the project root
// the tab starts in and RexSession the Rex session to open it in, which may
// be empty.
type Starter struct {
	Exec       func(name string, args []string, dir string) ([]byte, error)
	Self       string
	Root       string
	RexSession string
}

// SetPreview gives the model the runner of what-ifs: the same operations
// run as --dry-run does, so nothing is written.
func (m *Model) SetPreview(r Runner) { m.preview = r }

// SetStarter gives the model the means to start runs in Rex tabs.
func (m *Model) SetStarter(s Starter) { m.starter = s }

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
		return fmt.Sprintf("Start an interactive plan run of %s in a new Rex tab? y/n", short(p.id))
	}
	return fmt.Sprintf("Start an autonomous ship run of %s (permission mode auto) in a new Rex tab? y/n", short(p.id))
}

// args are the arguments of the `rex window new` call that starts p.
func (p *pending) args(s Starter) []string {
	cs := p.id.String()
	args := []string{"window", "new"}
	if s.RexSession != "" {
		args = append(args, "-s", s.RexSession)
	}
	args = append(args, cs[:min(len(cs), 9)], "-c", s.Root, "--json", "--", s.Self, "run", "start", cs)
	if p.kind == "plan" {
		return append(args, "--workflow", "plan")
	}
	return append(args, "--autonomous", "--permission-mode", "auto")
}

// confirm starts the pending run and says which tab it opened.
func (m *Model) confirm() {
	p := m.pending
	m.pending = nil
	if p.kind == "accept" {
		m.accept(p.id)
		return
	}
	if m.starter.Exec == nil || m.starter.Self == "" {
		m.notice = "Cannot start a run: no rex or rhiza path"
		return
	}
	if _, err := m.starter.Exec("rex", p.args(m.starter), m.starter.Root); err != nil {
		m.notice = fmt.Sprintf("rex window new failed: %v", err)
		return
	}
	cs := p.id.String()
	m.notice = fmt.Sprintf("Opened Rex tab %s running the %s run of %s", cs[:min(len(cs), 9)], p.kind, short(p.id))
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

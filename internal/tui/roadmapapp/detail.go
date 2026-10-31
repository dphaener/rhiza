package roadmapapp

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
)

// detail is one change set opened from the list: what change.show, the
// ready reasons and the latest run say about it.
type detail struct {
	id        id.ID
	kind      string
	state     string
	intent    string
	leaves    []leafView
	questions []question
	picked    int
	run       *api.RunInfo
	err       error
}

// question is an open question that blocks the change set.
type question struct {
	handle string
	text   string
}

type leafView struct {
	handle     string
	statement  string
	acceptance string
}

// blockedBy matches the ready reason of an open question that blocks a
// change set: "blocked by question <handle> on <target>".
var blockedBy = regexp.MustCompile(`^blocked by question (\S+) on `)

// open loads the detail of the change set under the cursor. Nothing
// happens on a heading.
func (m *Model) open() {
	if m.cursor >= len(m.rows) || m.rows[m.cursor].id.IsZero() {
		return
	}
	m.detail = m.loadDetail(m.rows[m.cursor])
}

func (m *Model) loadDetail(r row) *detail {
	d := &detail{id: r.id, state: r.state, intent: r.title}
	cs := r.id.String()
	out, err := m.run("change.show", cs, []string{cs})
	if err != nil {
		d.err = err
		return d
	}
	var cr api.ChangeShowResponse
	if err := json.Unmarshal(out, &cr); err != nil {
		d.err = fmt.Errorf("change.show: %w", err)
		return d
	}
	d.kind, d.state, d.intent = string(cr.Kind), string(cr.State), cr.Intent
	d.leaves = leafViews(cr)
	for _, reason := range r.reasons {
		if sub := blockedBy.FindStringSubmatch(reason); sub != nil {
			d.questions = append(d.questions, question{handle: sub[1], text: m.questionText(cs, sub[1])})
		}
	}
	d.run = m.latestRun(cs)
	return d
}

// latestRun is the change set's latest run, or nil. With no run for the
// change set, run.status answers with the latest run of another one, or an
// error: either way there is none to show.
func (m *Model) latestRun(cs string) *api.RunInfo {
	rout, err := m.run("run.status", cs, nil)
	if err != nil {
		return nil
	}
	var ri api.RunInfo
	if json.Unmarshal(rout, &ri) != nil || ri.Change != cs {
		return nil
	}
	return &ri
}

// reloadRun refreshes the open detail's latest-run section alone.
func (m *Model) reloadRun() {
	if m.detail == nil || m.detail.err != nil {
		return
	}
	m.detail.run = m.latestRun(m.detail.id.String())
}

// leafViews pairs each leaf with the statement and acceptance of the
// node.create operation that adds it.
func leafViews(cr api.ChangeShowResponse) []leafView {
	byNode := map[string]map[string]json.RawMessage{}
	for _, raw := range cr.Ops {
		var o struct {
			Op     string                     `json:"op"`
			Node   string                     `json:"node"`
			Fields map[string]json.RawMessage `json:"fields"`
		}
		if json.Unmarshal(raw, &o) == nil && o.Op == "node.create" {
			byNode[o.Node] = o.Fields
		}
	}
	var out []leafView
	for _, l := range cr.Leaves {
		f := byNode[l.ID.String()]
		out = append(out, leafView{handle: l.Handle, statement: field(f, "statement"), acceptance: field(f, "acceptance")})
	}
	return out
}

func field(f map[string]json.RawMessage, name string) string {
	var s string
	_ = json.Unmarshal(f[name], &s)
	return s
}

// questionText is the statement of the question with handle h, or h when
// it cannot be read.
func (m *Model) questionText(cs, h string) string {
	out, err := m.run("show", cs, []string{h})
	if err != nil {
		return h
	}
	var sr api.ShowResponse
	if json.Unmarshal(out, &sr) != nil || sr.Node == nil {
		return h
	}
	if s := field(sr.Node.Fields, "statement"); s != "" {
		return s
	}
	return h
}

func (d *detail) write(b *strings.Builder) {
	fmt.Fprintf(b, "%s [%s]\n", short(d.id), d.state)
	if d.err != nil {
		fmt.Fprintf(b, "rhiza: %v\n", d.err)
		return
	}
	fmt.Fprintf(b, "%s\nintent: %s\n", d.kind, d.intent)
	b.WriteString("\nLeaves\n")
	if len(d.leaves) == 0 {
		b.WriteString("  none\n")
	}
	for _, l := range d.leaves {
		fmt.Fprintf(b, "  %s: %s\n    acceptance: %s\n", l.handle, l.statement, l.acceptance)
	}
	b.WriteString("\nBlocking questions\n")
	if len(d.questions) == 0 {
		b.WriteString("  none\n")
	}
	for i, q := range d.questions {
		mark := "  "
		if i == d.picked && len(d.questions) > 1 {
			mark = "> "
		}
		fmt.Fprintf(b, "%s%s\n", mark, q.text)
	}
	b.WriteString("\nLatest run\n")
	if d.run == nil {
		b.WriteString("  none\n")
		return
	}
	fmt.Fprintf(b, "  %s · state %s", d.run.Workflow, d.run.State)
	if d.run.Current != "" {
		fmt.Fprintf(b, " · step %s", d.run.Current)
	}
	b.WriteString("\n")
}

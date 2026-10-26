// Package render writes the as-is model as markdown: an index and one page
// per live node, the files `rhiza render` keeps under .rhiza/rendered/.
// The output is a pure function of its input: no file I/O, no clock, and the
// same model always gives the same bytes.
//
// The vocabulary:
//
//   - A page is the markdown file of one live node. Retired nodes get no
//     page; the index lists them.
//   - A node's path is where its page goes, relative to .rhiza/rendered/
//     and written with slashes. A capability or an intent sits at
//     <kind>/<ancestor handles...>/<handle>.md, its ancestors' handles from
//     the root down, whether or not they are live and whatever their kind.
//     Every other node sits at <kind>/<handle>.md.
//   - Two live nodes collide when they would get the same path (a duplicate
//     handle, which the model reports). The one with the smallest ID keeps
//     the path; each other gets <handle>.<ID>.md in the same folder.
//   - A link's endpoint is a node, live or retired, or a change set. Live
//     nodes are written as relative links to their page, retired ones as
//     "<handle> (retired)", change sets as "change set `<ID>` (<intent>)"
//     and anything else as "`<ID>` (unknown)".
//   - A one-line text (an intent, a reason, a note, a field name) has every
//     run of white space, newlines included, written as one space.
//
// The index, index.md, is "# Model", then these sections, each only when it
// has an entry:
//
//   - "## Capabilities" and "## Intents": the tree as nested bullets, two
//     spaces per level, siblings in Children order (key, then ID). A live
//     node of the section's kind is a link to its page. Any other node is
//     shown only when a live node of the section's kind sits under it: a
//     retired one as "<handle> (retired)", a live one of the other kind as
//     its link followed by " (<kind>)".
//   - One section per other kind, in op.Kinds() order, headed by its plural
//     (Constraints, Decisions, Questions, Evidence, Sources, Actors,
//     Anchors): a bullet with a link per live node, sorted by handle then ID.
//   - "## Retired": "- <kind> <handle> `<ID>`", then " — <reason>" when there
//     is one, sorted by op.Kinds() order then ID.
//
// With no node at all it is exactly "# Model\n\nNothing has landed yet.\n".
//
// A page is these blocks, each only when it has content:
//
//   - "# <handle>", then "<kind> · `<ID>`", plus " · parent [<handle>](<path>)"
//     when it has a parent ("<handle> (retired)" for a retired one).
//   - One "## <name>" per field, sorted by name. A JSON string is decoded
//     and written as markdown, line endings made LF and leading and trailing
//     blank lines trimmed; an empty one leaves the heading alone. Any other
//     JSON is its canonical form in a ```json fence, the fence longer than
//     any run of backticks inside.
//   - "## Links": "### Out" (links from the node) then "### In" (links to
//     it), each only when it has a link. Live links only, grouped by type in
//     op.LinkTypes() order, then by link ID; each "- <type> <endpoint>",
//     plus " — <note>" when the link has one. A link from a node to itself
//     is in both.
//   - "## History": every change set in the node's history once, in apply
//     order (first appearance), as "- `<ID>` <intent>".
//
// Blocks are separated by one blank line; every file is LF only and ends
// with exactly one newline. Nothing carries a timestamp or an absolute path.
package render

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// IndexPath is the index's path.
const IndexPath = "index.md"

// empty is the whole index of a model with no nodes.
const empty = "# Model\n\nNothing has landed yet.\n"

// plurals heads the index section of each kind.
var plurals = map[op.Kind]string{
	op.KindIntent:     "Intents",
	op.KindCapability: "Capabilities",
	op.KindConstraint: "Constraints",
	op.KindDecision:   "Decisions",
	op.KindQuestion:   "Questions",
	op.KindEvidence:   "Evidence",
	op.KindSource:     "Sources",
	op.KindActor:      "Actors",
	op.KindAnchor:     "Anchors",
}

// Render returns the rendered files of m, the as-is view, keyed by their
// slash path relative to .rhiza/rendered/. sets gives the intents of the
// change sets that link endpoints and histories name; the change sets of the
// view are expected among them. The result always holds index.md.
func Render(m model.Model, sets []changeset.ChangeSet) map[string][]byte {
	r := newRenderer(m, sets)
	out := map[string][]byte{IndexPath: r.index()}
	for i, p := range r.paths {
		out[p] = r.page(r.nodes[i], p)
	}
	return out
}

// renderer holds what every file needs.
type renderer struct {
	m     model.Model
	nodes map[id.ID]model.Node
	// paths maps each live node to its page's path.
	paths map[id.ID]string
	// out and in map each node or change set to the live links from it and
	// to it, each slice in link ID order.
	out, in map[id.ID][]model.Link
	// intents maps each change set in the input to its intent.
	intents map[id.ID]string
}

func newRenderer(m model.Model, sets []changeset.ChangeSet) *renderer {
	r := &renderer{m: m, nodes: map[id.ID]model.Node{}, paths: map[id.ID]string{}, intents: map[id.ID]string{},
		out: map[id.ID][]model.Link{}, in: map[id.ID][]model.Link{}}
	for _, s := range sets {
		if _, ok := r.intents[s.ID]; !ok {
			r.intents[s.ID] = s.Open.Intent
		}
	}
	all := m.Nodes()
	for _, n := range all {
		r.nodes[n.ID] = n
	}
	// all is sorted by ID, so the first to claim a path has the smallest ID.
	taken := map[string]bool{}
	for _, n := range all {
		if n.Retired {
			continue
		}
		dir := r.dir(n)
		p := dir + n.Handle + ".md"
		if taken[p] {
			p = dir + n.Handle + "." + n.ID.String() + ".md"
		}
		taken[p] = true
		r.paths[n.ID] = p
	}
	for _, l := range m.Links() {
		if !l.Retired {
			r.out[l.From] = append(r.out[l.From], l)
			r.in[l.To] = append(r.in[l.To], l)
		}
	}
	return r
}

// dir returns the folder of n's page, with its final slash: its kind, then
// for a node in the tree its ancestors' handles.
func (r *renderer) dir(n model.Node) string {
	var handles []string
	if op.HasParent(n.Kind) {
		seen := map[id.ID]bool{n.ID: true}
		for p := n.Parent; !p.IsZero() && !seen[p]; {
			seen[p] = true
			a, ok := r.nodes[p]
			if !ok {
				break
			}
			handles = append(handles, a.Handle)
			p = a.Parent
		}
		slices.Reverse(handles)
	}
	return strings.Join(append([]string{string(n.Kind)}, handles...), "/") + "/"
}

// index returns index.md.
func (r *renderer) index() []byte {
	if len(r.nodes) == 0 {
		return []byte(empty)
	}
	blocks := []string{"# Model"}
	for _, k := range []op.Kind{op.KindCapability, op.KindIntent} {
		var lines []string
		r.tree(id.ID{}, k, 0, &lines)
		blocks = appendSection(blocks, "## "+plurals[k], lines)
	}
	for _, k := range op.Kinds() {
		if op.HasParent(k) {
			continue
		}
		var ns []model.Node
		for _, n := range r.nodes {
			if n.Kind == k && !n.Retired {
				ns = append(ns, n)
			}
		}
		slices.SortFunc(ns, func(a, b model.Node) int {
			return cmp.Or(strings.Compare(a.Handle, b.Handle), id.Compare(a.ID, b.ID))
		})
		var lines []string
		for _, n := range ns {
			lines = append(lines, "- "+r.nodeLink(n, IndexPath))
		}
		blocks = appendSection(blocks, "## "+plurals[k], lines)
	}
	blocks = appendSection(blocks, "## Retired", r.retired())
	return finish(blocks)
}

// tree appends the bullets for the nodes under parent to lines, for the
// section of kind k, and reports whether it appended any.
func (r *renderer) tree(parent id.ID, k op.Kind, depth int, lines *[]string) bool {
	found := false
	for _, c := range r.m.Children(parent) {
		var sub []string
		hasSub := depth < len(r.nodes) && r.tree(c.ID, k, depth+1, &sub)
		mine := c.Kind == k && !c.Retired
		if !mine && !hasSub {
			continue
		}
		var text string
		switch {
		case mine:
			text = r.nodeLink(c, IndexPath)
		case c.Retired:
			text = c.Handle + " (retired)"
		default:
			text = r.nodeLink(c, IndexPath) + " (" + string(c.Kind) + ")"
		}
		*lines = append(*lines, strings.Repeat("  ", depth)+"- "+text)
		*lines = append(*lines, sub...)
		found = true
	}
	return found
}

// retired returns the Retired section's lines.
func (r *renderer) retired() []string {
	var ns []model.Node
	for _, n := range r.nodes {
		if n.Retired {
			ns = append(ns, n)
		}
	}
	slices.SortFunc(ns, func(a, b model.Node) int {
		return cmp.Or(cmp.Compare(kindRank(a.Kind), kindRank(b.Kind)), id.Compare(a.ID, b.ID))
	})
	var lines []string
	for _, n := range ns {
		l := fmt.Sprintf("- %s %s `%s`", n.Kind, n.Handle, n.ID)
		if reason := oneLine(n.Reason); reason != "" {
			l += " — " + reason
		}
		lines = append(lines, l)
	}
	return lines
}

// page returns the page of n, whose path is self.
func (r *renderer) page(n model.Node, self string) []byte {
	head := fmt.Sprintf("%s · `%s`", n.Kind, n.ID)
	if !n.Parent.IsZero() {
		head += " · parent " + r.endpoint(n.Parent, self)
	}
	blocks := []string{"# " + n.Handle, head}
	names := make([]string, 0, len(n.Fields))
	for name := range n.Fields {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		blocks = append(blocks, "## "+oneLine(name))
		if v := fieldValue(n.Fields[name]); v != "" {
			blocks = append(blocks, v)
		}
	}
	out := r.linkLines(self, r.out[n.ID], func(l model.Link) id.ID { return l.To })
	in := r.linkLines(self, r.in[n.ID], func(l model.Link) id.ID { return l.From })
	if len(out) > 0 || len(in) > 0 {
		blocks = append(blocks, "## Links")
		blocks = appendSection(blocks, "### Out", out)
		blocks = appendSection(blocks, "### In", in)
	}
	blocks = appendSection(blocks, "## History", r.history(n))
	return finish(blocks)
}

// linkLines returns a bullet per link of links, which are in link ID order,
// with the endpoint end names, grouped by type in op.LinkTypes() order. A
// link of a type not in that list is left out.
func (r *renderer) linkLines(self string, links []model.Link, end func(model.Link) id.ID) []string {
	links = slices.DeleteFunc(slices.Clone(links), func(l model.Link) bool { return linkRank(l.Type) < 0 })
	slices.SortStableFunc(links, func(a, b model.Link) int { return cmp.Compare(linkRank(a.Type), linkRank(b.Type)) })
	lines := make([]string, 0, len(links))
	for _, l := range links {
		line := "- " + string(l.Type) + " " + r.endpoint(end(l), self)
		if note := oneLine(l.Note); note != "" {
			line += " — " + note
		}
		lines = append(lines, line)
	}
	return lines
}

// linkRanks maps each link type to its place in op.LinkTypes().
var linkRanks = func() map[op.LinkType]int {
	m := map[op.LinkType]int{}
	for i, t := range op.LinkTypes() {
		m[t] = i
	}
	return m
}()

// linkRank returns t's place in op.LinkTypes(), -1 when unknown.
func linkRank(t op.LinkType) int {
	if i, ok := linkRanks[t]; ok {
		return i
	}
	return -1
}

// history returns the History section's lines for n.
func (r *renderer) history(n model.Node) []string {
	var lines []string
	seen := map[id.ID]bool{}
	for _, h := range n.History {
		if seen[h.ChangeSet] {
			continue
		}
		seen[h.ChangeSet] = true
		l := fmt.Sprintf("- `%s`", h.ChangeSet)
		if intent := oneLine(r.intents[h.ChangeSet]); intent != "" {
			l += " " + intent
		}
		lines = append(lines, l)
	}
	return lines
}

// endpoint writes i, a node or a change set, as seen from the page at self.
func (r *renderer) endpoint(i id.ID, self string) string {
	if n, ok := r.nodes[i]; ok {
		if n.Retired {
			return n.Handle + " (retired)"
		}
		return r.nodeLink(n, self)
	}
	if intent, ok := r.intents[i]; ok {
		s := fmt.Sprintf("change set `%s`", i)
		if intent = oneLine(intent); intent != "" {
			s += " (" + intent + ")"
		}
		return s
	}
	return fmt.Sprintf("`%s` (unknown)", i)
}

// nodeLink returns a markdown link to live node n's page from the page at
// self.
func (r *renderer) nodeLink(n model.Node, self string) string {
	return "[" + n.Handle + "](" + rel(self, r.paths[n.ID]) + ")"
}

// fieldValue writes a field's JSON value as markdown: a string as itself,
// anything else fenced in its canonical form.
func fieldValue(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		s = strings.ReplaceAll(s, "\r\n", "\n")
		s = strings.ReplaceAll(s, "\r", "\n")
		return strings.TrimRight(strings.TrimLeft(s, "\n"), " \t\n")
	}
	body, err := canon.Canonicalize(raw)
	if err != nil {
		body = bytes.ReplaceAll(raw, []byte("\r"), nil)
	}
	text := strings.TrimRight(string(body), "\n")
	fence := strings.Repeat("`", max(3, longestRun(text, '`')+1))
	return fence + "json\n" + text + "\n" + fence
}

// longestRun returns the length of the longest run of c in s.
func longestRun(s string, c byte) int {
	best, run := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			run++
			best = max(best, run)
		} else {
			run = 0
		}
	}
	return best
}

// oneLine writes s on one line: each run of white space becomes one space.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// kindRank returns k's place in op.Kinds(), past the end when unknown.
func kindRank(k op.Kind) int {
	if i := slices.Index(op.Kinds(), k); i >= 0 {
		return i
	}
	return len(op.Kinds())
}

// appendSection appends heading and the lines as one block after it, or
// nothing when there are no lines.
func appendSection(blocks []string, heading string, lines []string) []string {
	if len(lines) == 0 {
		return blocks
	}
	return append(blocks, heading, strings.Join(lines, "\n"))
}

// finish joins blocks with blank lines and ends the file with one newline.
func finish(blocks []string) []byte {
	return []byte(strings.Join(blocks, "\n\n") + "\n")
}

// rel returns the relative slash path from the file at from to the file at
// to, both relative to the same root.
func rel(from, to string) string {
	f := strings.Split(from, "/")
	f = f[:len(f)-1]
	t := strings.Split(to, "/")
	i := 0
	for i < len(f) && i < len(t)-1 && f[i] == t[i] {
		i++
	}
	parts := slices.Repeat([]string{".."}, len(f)-i)
	return strings.Join(append(parts, t[i:]...), "/")
}

package slice

// This file writes a Slice as the bytes printed: markdown, or canonical
// JSON. Both are pure functions of the Slice, the tiers Fit chose and the
// Options: no clock, no map order and no path reaches them.
//
// Markdown is these blocks, separated by one blank line, LF only, ending
// with exactly one newline and with no trailing white space on any line:
//
//   - The heading, "# Slice · <title> · <view>". A node slice viewed
//     through a change set adds ` (change <short> "<intent>")`; a change
//     slice's title already names its change set.
//   - Each section with something to show, as "## <section>" and then its
//     items, in this order: Focus, Context, Constraints, Decisions, Open
//     questions, Anchors, Evidence, Children, Neighbours, In flight. An
//     empty section, or one cut to a count, is left out.
//   - Each Focus item is a block of its own: "**<handle>** `<short>` ·
//     <kind> · evidence: <status>", plus " · retired" (" · retired:
//     <reason>" with a reason) for a retired one; then its statement; then
//     every other field, sorted by name, as "<name>: <value>"; then one
//     "Proposed: <line>" per line of its proposed change. The statement
//     field is left out of the fields when it is a JSON string, since the
//     statement line shows it. A JSON string is its text, line endings made
//     LF, trailing white space trimmed from every line, leading and trailing
//     blank lines dropped and runs of blank lines made one; any other JSON
//     is its canonical form. A value of one line follows "<name>: "
//     ("<name>:" when empty); a longer one goes on the lines after "<name>:"
//     in a fence, ```json for JSON that is not a string, the fence longer
//     than any run of backticks inside.
//   - Every other section is one block of lines. An item at full is "-
//     <handle> `<short>` — <statement>", plus " · from <handle>", " ·
//     <status>", " · <relation>" and " · <note>", each when it is set; the
//     relation "ancestor" is not written, since Context says it. An In
//     flight item is "- `<short>` <state> — <intent>" (no " — " with no
//     intent). At brief an item is "- <handle> `<short>`", an In flight one
//     "- `<short>`".
//   - The footer, only when a section was cut or the slice is over budget:
//     "---", then "Omitted: <section> <count> — `<command>`" for each cut
//     section in section order, the section in lower case, with " (brief)"
//     after the count for one cut to brief; then "Still over budget by <N>
//     tokens." when it is over, N the tokens past the budget, the footer
//     included. The lines follow one another with no blank line between.
//
// Every free text (a handle, statement, note, intent, reason, field name or
// proposed line) is one line: each run of white space written as one space.
//
// JSON is one canonical object (internal/canon), Doc:
//
//	{view, change, ref, budget, tokens, over_budget, focus, context,
//	 constraints, decisions, questions, anchors, evidence, children,
//	 neighbours, in_flight, omitted}
//
// "change" is the change set sliced or viewed through, absent with none.
// Every section is an array, [] when empty or cut to a count. An item has
// id (the full ID), short, handle, kind and statement, plus note, from,
// status and relation, each absent when empty; at brief it has only id,
// short, handle and kind. A Focus item adds fields (every field, as stored,
// {} for none), retired, reason (absent when empty) and proposed ([] for
// none); its status is its evidence status. An In flight item has id,
// short, state and intent, and at brief only id and short. "omitted" is
// [{section, count, tier, next}] for each cut section in section order,
// section being the section's key ("in_flight" for In flight). "tokens" is
// Tokens of the whole printed document, "tokens" field included: Fit
// renders until the number it writes is the number it measures.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// Doc is a slice as JSON prints it; canon.Marshal of a Doc is the printed
// bytes.
type Doc struct {
	View        string      `json:"view" desc:"The view sliced: as-is, to-be or what-if"`
	Change      *DocChange  `json:"change,omitempty" desc:"The change set sliced or viewed through; absent for none"`
	Ref         string      `json:"ref" desc:"The focus handle of a node slice, or the change set's short ID of a change slice"`
	Budget      int         `json:"budget" desc:"The token budget"`
	Tokens      int         `json:"tokens" desc:"The estimated tokens of this whole document: its bytes divided by 4, rounded up"`
	OverBudget  bool        `json:"over_budget" desc:"True when the slice is over budget even with every section that can be cut cut to a count"`
	Focus       []DocFocus  `json:"focus" desc:"The focus nodes, never cut"`
	Context     []DocItem   `json:"context" desc:"The ancestors, the root first, then the intents served"`
	Constraints []DocItem   `json:"constraints" desc:"The constraints that apply, never cut"`
	Decisions   []DocItem   `json:"decisions" desc:"The current decisions that apply"`
	Questions   []DocItem   `json:"questions" desc:"The open questions, blocking ones first, never cut"`
	Anchors     []DocItem   `json:"anchors" desc:"The anchors that implement the focus"`
	Evidence    []DocItem   `json:"evidence" desc:"The evidence that verifies the focus"`
	Children    []DocItem   `json:"children" desc:"The live children of the focus"`
	Neighbours  []DocItem   `json:"neighbours" desc:"Nodes the focus needs, is needed by or relates to"`
	InFlight    []DocChange `json:"in_flight" desc:"Other in-flight change sets that target a focus node"`
	Omitted     []Omission  `json:"omitted" desc:"Each section cut to fit the budget, in section order"`
}

// DocItem is an item as JSON prints it.
type DocItem struct {
	ID        id.ID   `json:"id" desc:"The node's full ID"`
	Short     string  `json:"short" desc:"The node's short ID"`
	Handle    string  `json:"handle" desc:"The node's handle"`
	Kind      op.Kind `json:"kind" desc:"The node's kind"`
	Statement string  `json:"statement,omitempty" desc:"The node's statement, or its handle; absent at brief"`
	Note      string  `json:"note,omitempty" desc:"The link's note"`
	From      string  `json:"from,omitempty" desc:"The ancestor it is inherited from"`
	Status    string  `json:"status,omitempty" desc:"An evidence status: verified, failing, unverified or none"`
	Relation  string  `json:"relation,omitempty" desc:"How it relates to the focus: ancestor, serves, blocks, relates-to, needs or needed by"`
}

// DocFocus is a Focus item as JSON prints it.
type DocFocus struct {
	ID        id.ID                      `json:"id" desc:"The node's full ID"`
	Short     string                     `json:"short" desc:"The node's short ID"`
	Handle    string                     `json:"handle" desc:"The node's handle"`
	Kind      op.Kind                    `json:"kind" desc:"The node's kind"`
	Statement string                     `json:"statement" desc:"The node's statement, or its handle"`
	Status    string                     `json:"status" desc:"The node's evidence status: verified, failing, unverified or none"`
	Fields    map[string]json.RawMessage `json:"fields" desc:"Every field of the node, as stored"`
	Retired   bool                       `json:"retired" desc:"True when the node is retired in this view"`
	Reason    string                     `json:"reason,omitempty" desc:"Why it was retired"`
	Proposed  []string                   `json:"proposed" desc:"The proposed change against as-is, one line each"`
}

// DocChange is a change set as JSON prints it.
type DocChange struct {
	ID     id.ID           `json:"id" desc:"The change set's full ID"`
	Short  string          `json:"short" desc:"The change set's short ID"`
	State  changeset.State `json:"state,omitempty" desc:"Its state; absent at brief"`
	Intent string          `json:"intent,omitempty" desc:"Its intent; absent at brief or when empty"`
}

// part is a section, by its place in the order.
type part int

const (
	secFocus part = iota
	secContext
	secConstraints
	secDecisions
	secQuestions
	secAnchors
	secEvidence
	secChildren
	secNeighbours
	secInFlight
	numSections
)

// sectionKeys names each section in JSON and in Omission.
var sectionKeys = [numSections]string{"focus", "context", "constraints", "decisions", "questions",
	"anchors", "evidence", "children", "neighbours", "in_flight"}

// headings heads each section in markdown.
var headings = [numSections]string{"Focus", "Context", "Constraints", "Decisions", "Open questions",
	"Anchors", "Evidence", "Children", "Neighbours", "In flight"}

// items returns the items of an item section: not Focus or In flight.
func (s Slice) items(sec part) []Item {
	switch sec {
	case secContext:
		return s.Context
	case secConstraints:
		return s.Constraints
	case secDecisions:
		return s.Decisions
	case secQuestions:
		return s.Questions
	case secAnchors:
		return s.Anchors
	case secEvidence:
		return s.Evidence
	case secChildren:
		return s.Children
	case secNeighbours:
		return s.Neighbours
	}
	return nil
}

// count returns how many items section sec holds.
func (s Slice) count(sec part) int {
	switch sec {
	case secFocus:
		return len(s.Focus)
	case secInFlight:
		return len(s.InFlight)
	}
	return len(s.items(sec))
}

// layout is what Fit chose: each section's tier, the omissions they make,
// and the numbers the output prints about itself.
type layout struct {
	tiers   [numSections]Tier
	omitted []Omission
	budget  int
	tokens  int  // the JSON tokens field
	over    bool // over budget
	overBy  int  // the tokens past the budget, for the markdown footer
}

// markdown writes s as markdown under lay.
func markdown(s Slice, lay layout) []byte {
	blocks := []string{heading(s)}
	if len(s.Focus) > 0 {
		blocks = append(blocks, "## "+headings[secFocus])
		for _, f := range s.Focus {
			blocks = append(blocks, focusBlock(f))
		}
	}
	for sec := secContext; sec < secInFlight; sec++ {
		items := s.items(sec)
		if len(items) == 0 || lay.tiers[sec] == TierCount {
			continue
		}
		lines := make([]string, len(items))
		for i, it := range items {
			lines[i] = itemLine(it, lay.tiers[sec])
		}
		blocks = append(blocks, "## "+headings[sec], strings.Join(lines, "\n"))
	}
	if len(s.InFlight) > 0 && lay.tiers[secInFlight] != TierCount {
		lines := make([]string, len(s.InFlight))
		for i, c := range s.InFlight {
			lines[i] = changeLine(c, lay.tiers[secInFlight])
		}
		blocks = append(blocks, "## "+headings[secInFlight], strings.Join(lines, "\n"))
	}
	if len(lay.omitted) > 0 || lay.over {
		lines := []string{"---"}
		for _, o := range lay.omitted {
			brief := ""
			if o.Tier == TierBrief.String() {
				brief = " (brief)"
			}
			lines = append(lines, fmt.Sprintf("Omitted: %s %d%s — `%s`", strings.ReplaceAll(o.Section, "_", " "), o.Count, brief, o.Next))
		}
		if lay.over {
			lines = append(lines, fmt.Sprintf("Still over budget by %d tokens.", lay.overBy))
		}
		blocks = append(blocks, strings.Join(lines, "\n"))
	}
	return []byte(strings.Join(blocks, "\n\n") + "\n")
}

// heading is the slice's first line.
func heading(s Slice) string {
	h := "# Slice · " + oneLine(s.Title)
	if v := oneLine(s.View); v != "" {
		h += " · " + v
	}
	if !s.ChangeSlice && s.Change != nil {
		h += fmt.Sprintf(" (change %s %q)", s.Change.Short, oneLine(s.Change.Intent))
	}
	return h
}

// focusBlock writes Focus item f.
func focusBlock(f FocusItem) string {
	head := fmt.Sprintf("**%s** `%s` · %s", oneLine(f.Handle), f.Short, f.Kind)
	if f.Status != "" {
		head += " · evidence: " + f.Status
	}
	if f.Retired {
		head += " · retired"
		if r := oneLine(f.Reason); r != "" {
			head += ": " + r
		}
	}
	lines := []string{head}
	if st := oneLine(f.Statement); st != "" {
		lines = append(lines, st)
	}
	for _, k := range slices.Sorted(maps.Keys(f.Fields)) {
		v := f.Fields[k]
		var str string
		if k == "statement" && json.Unmarshal(v, &str) == nil {
			continue
		}
		lines = append(lines, fieldLines(oneLine(k), v))
	}
	for _, p := range f.Proposed {
		lines = append(lines, "Proposed: "+oneLine(p))
	}
	return strings.Join(lines, "\n")
}

// fieldLines writes field name with value raw: see the file comment.
func fieldLines(name string, raw json.RawMessage) string {
	var text, lang string
	var str string
	if json.Unmarshal(raw, &str) == nil {
		text = textBlock(str)
	} else {
		body, err := canon.Canonicalize(raw)
		if err != nil {
			body = bytes.ReplaceAll(raw, []byte("\r"), nil)
		}
		text, lang = textBlock(string(body)), "json"
	}
	switch {
	case text == "":
		return name + ":"
	case !strings.Contains(text, "\n"):
		return name + ": " + text
	}
	fence := strings.Repeat("`", max(3, longestRun(text, '`')+1))
	return name + ":\n" + fence + lang + "\n" + text + "\n" + fence
}

// textBlock makes s LF only, trims trailing white space from every line,
// drops leading and trailing blank lines and makes each run of blank lines
// one.
func textBlock(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	var out []string
	blank := false
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimRight(line, " \t\v\f")
		if line == "" {
			blank = len(out) > 0
			continue
		}
		if blank {
			out = append(out, "")
			blank = false
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
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

// itemLine writes item it at tier t.
func itemLine(it Item, t Tier) string {
	line := fmt.Sprintf("- %s `%s`", oneLine(it.Handle), it.Short)
	if t == TierBrief {
		return line
	}
	if st := oneLine(it.Statement); st != "" {
		line += " — " + st
	}
	if it.From != "" {
		line += " · from " + oneLine(it.From)
	}
	if it.Status != "" {
		line += " · " + it.Status
	}
	if it.Relation != "" && it.Relation != RelAncestor {
		line += " · " + it.Relation
	}
	if n := oneLine(it.Note); n != "" {
		line += " · " + n
	}
	return line
}

// changeLine writes In flight item c at tier t.
func changeLine(c Change, t Tier) string {
	line := fmt.Sprintf("- `%s`", c.Short)
	if t == TierBrief {
		return line
	}
	if c.State != "" {
		line += " " + string(c.State)
	}
	if in := oneLine(c.Intent); in != "" {
		line += " — " + in
	}
	return line
}

// doc returns s as a Doc under lay, ref being its ref.
func doc(s Slice, lay layout, view, ref string) *Doc {
	d := &Doc{View: view, Ref: ref, Budget: lay.budget, Tokens: lay.tokens, OverBudget: lay.over,
		Focus: []DocFocus{}, InFlight: []DocChange{}, Omitted: slices.Clone(lay.omitted)}
	if d.Omitted == nil {
		d.Omitted = []Omission{}
	}
	if s.Change != nil {
		d.Change = &DocChange{ID: s.Change.ID, Short: s.Change.Short, State: s.Change.State, Intent: oneLine(s.Change.Intent)}
	}
	for _, f := range s.Focus {
		fields := maps.Clone(f.Fields)
		if fields == nil {
			fields = map[string]json.RawMessage{}
		}
		proposed := []string{}
		for _, p := range f.Proposed {
			proposed = append(proposed, oneLine(p))
		}
		d.Focus = append(d.Focus, DocFocus{ID: f.ID, Short: f.Short, Handle: f.Handle, Kind: f.Kind, Statement: oneLine(f.Statement),
			Status: f.Status, Fields: fields, Retired: f.Retired, Reason: oneLine(f.Reason), Proposed: proposed})
	}
	for sec, dst := range map[part]*[]DocItem{secContext: &d.Context, secConstraints: &d.Constraints,
		secDecisions: &d.Decisions, secQuestions: &d.Questions, secAnchors: &d.Anchors, secEvidence: &d.Evidence,
		secChildren: &d.Children, secNeighbours: &d.Neighbours} {
		// Each section fills only its own array, so the map's order does
		// not reach the output.
		*dst = []DocItem{}
		if lay.tiers[sec] == TierCount {
			continue
		}
		for _, it := range s.items(sec) {
			di := DocItem{ID: it.ID, Short: it.Short, Handle: it.Handle, Kind: it.Kind}
			if lay.tiers[sec] == TierFull {
				di.Statement, di.Note, di.From, di.Status, di.Relation = oneLine(it.Statement), oneLine(it.Note), oneLine(it.From), it.Status, it.Relation
			}
			*dst = append(*dst, di)
		}
	}
	if lay.tiers[secInFlight] != TierCount {
		for _, c := range s.InFlight {
			dc := DocChange{ID: c.ID, Short: c.Short}
			if lay.tiers[secInFlight] == TierFull {
				dc.State, dc.Intent = c.State, oneLine(c.Intent)
			}
			d.InFlight = append(d.InFlight, dc)
		}
	}
	return d
}

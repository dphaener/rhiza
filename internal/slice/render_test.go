package slice

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

// ulid parses a hand-made ULID.
func ulid(t testing.TB, s string) id.ID {
	t.Helper()
	i, err := id.Parse(s)
	if err != nil {
		t.Fatalf("id.Parse(%q): %v", s, err)
	}
	return i
}

// small is a node slice of hard-pause, viewed to-be through an accepted
// change set, with something in most sections.
func small(t testing.TB) Slice {
	return Slice{
		Title:  "hard-pause",
		View:   "to-be",
		Change: &Change{ID: ulid(t, "01JB7A0000000000000000000C"), Short: "01JB7A", State: changeset.StateAccepted, Intent: "Pause at 90%"},
		Focus: []FocusItem{{
			Item: Item{ID: ulid(t, "01JB2Q0000000000000000000F"), Short: "01JB2Q", Handle: "hard-pause", Kind: op.KindCapability,
				Statement: "Stop new work at 90%.", Status: StatusUnverified},
			Fields: fields(
				"statement", q("Stop new work at 90%."),
				"threshold", "90",
				"notes", q("\nline one  \r\nline two\n\n"),
				"tags", `["a","b"]`,
				"empty", q(""),
				"owner", q("ops team  "),
			),
			Proposed: []string{"threshold: 95 → 90"},
		}},
		Context: []Item{
			{ID: ulid(t, "01J8AA0000000000000000000R"), Short: "01J8AA", Handle: "runtime", Kind: op.KindCapability, Statement: "Runs agents", Relation: RelAncestor},
			{ID: ulid(t, "01J8AB0000000000000000000N"), Short: "01J8AB", Handle: "reserve", Kind: op.KindIntent, Statement: "Keep a reserve", Relation: RelServes},
		},
		Constraints: []Item{
			{ID: ulid(t, "01J9G00000000000000000000K"), Short: "01J9G0", Handle: "owner-reserve", Kind: op.KindConstraint, Statement: "keep 10%", From: "runtime"},
		},
		Questions: []Item{
			{ID: ulid(t, "01J9Q00000000000000000000W"), Short: "01J9Q0", Handle: "when-resume", Kind: op.KindQuestion, Statement: "When to resume?", Relation: RelBlocks},
		},
		Evidence: []Item{
			{ID: ulid(t, "01J9E00000000000000000000T"), Short: "01J9E0", Handle: "pause-test", Kind: op.KindEvidence, Statement: "pause-test", Status: StatusFailing},
		},
		Neighbours: []Item{
			{ID: ulid(t, "01J9R00000000000000000000S"), Short: "01J9R0", Handle: "resume", Kind: op.KindCapability, Statement: "Resume after reset", Relation: RelNeeds, Note: "after reset"},
		},
		InFlight: []Change{
			{ID: ulid(t, "01JB9C0000000000000000000D"), Short: "01JB9C", State: changeset.StateDraft, Intent: "Tune"},
		},
	}
}

const smallMarkdown = "# Slice · hard-pause · to-be (change 01JB7A \"Pause at 90%\")\n" +
	"\n" +
	"## Focus\n" +
	"\n" +
	"**hard-pause** `01JB2Q` · capability · evidence: unverified\n" +
	"Stop new work at 90%.\n" +
	"empty:\n" +
	"notes:\n" +
	"```\n" +
	"line one\n" +
	"line two\n" +
	"```\n" +
	"owner: ops team\n" +
	"tags:\n" +
	"```json\n" +
	"[\n" +
	"  \"a\",\n" +
	"  \"b\"\n" +
	"]\n" +
	"```\n" +
	"threshold: 90\n" +
	"Proposed: threshold: 95 → 90\n" +
	"\n" +
	"## Context\n" +
	"\n" +
	"- runtime `01J8AA` — Runs agents\n" +
	"- reserve `01J8AB` — Keep a reserve · serves\n" +
	"\n" +
	"## Constraints\n" +
	"\n" +
	"- owner-reserve `01J9G0` — keep 10% · from runtime\n" +
	"\n" +
	"## Open questions\n" +
	"\n" +
	"- when-resume `01J9Q0` — When to resume? · blocks\n" +
	"\n" +
	"## Evidence\n" +
	"\n" +
	"- pause-test `01J9E0` — pause-test · failing\n" +
	"\n" +
	"## Neighbours\n" +
	"\n" +
	"- resume `01J9R0` — Resume after reset · needs · after reset\n" +
	"\n" +
	"## In flight\n" +
	"\n" +
	"- `01JB9C` draft — Tune\n"

// tiny is a node slice in as-is with no change set: one focus field, one
// constraint, one child.
func tiny(t testing.TB) Slice {
	return Slice{
		Title: "pause",
		View:  "as-is",
		Focus: []FocusItem{{
			Item:   Item{ID: ulid(t, "01JB2Q0000000000000000000F"), Short: "01JB2Q", Handle: "pause", Kind: op.KindCapability, Statement: "pause", Status: StatusNone},
			Fields: fields("size", "3"),
		}},
		Constraints: []Item{
			{ID: ulid(t, "01J9G00000000000000000000K"), Short: "01J9G0", Handle: "no-spend", Kind: op.KindConstraint, Statement: "no spend", From: "runtime"},
		},
		Children: []Item{
			{ID: ulid(t, "01J9C00000000000000000000V"), Short: "01J9C0", Handle: "step", Kind: op.KindCapability, Statement: "one step", Status: StatusVerified},
		},
	}
}

const tinyMarkdown = "# Slice · pause · as-is\n" +
	"\n" +
	"## Focus\n" +
	"\n" +
	"**pause** `01JB2Q` · capability · evidence: none\n" +
	"pause\n" +
	"size: 3\n" +
	"\n" +
	"## Constraints\n" +
	"\n" +
	"- no-spend `01J9G0` — no spend · from runtime\n" +
	"\n" +
	"## Children\n" +
	"\n" +
	"- step `01J9C0` — one step · verified\n"

const tinyJSON = `{
  "anchors": [],
  "budget": 4000,
  "children": [
    {
      "handle": "step",
      "id": "01J9C00000000000000000000V",
      "kind": "capability",
      "short": "01J9C0",
      "statement": "one step",
      "status": "verified"
    }
  ],
  "constraints": [
    {
      "from": "runtime",
      "handle": "no-spend",
      "id": "01J9G00000000000000000000K",
      "kind": "constraint",
      "short": "01J9G0",
      "statement": "no spend"
    }
  ],
  "context": [],
  "decisions": [],
  "evidence": [],
  "focus": [
    {
      "fields": {
        "size": 3
      },
      "handle": "pause",
      "id": "01JB2Q0000000000000000000F",
      "kind": "capability",
      "proposed": [],
      "retired": false,
      "short": "01JB2Q",
      "statement": "pause",
      "status": "none"
    }
  ],
  "in_flight": [],
  "neighbours": [],
  "omitted": [],
  "over_budget": false,
  "questions": [],
  "ref": "pause",
  "tokens": 241,
  "view": "as-is"
}
`

func TestRender(t *testing.T) {
	cases := []struct {
		name string
		s    Slice
		json bool
		want string
	}{
		{"markdown small", small(t), false, smallMarkdown},
		{"markdown tiny", tiny(t), false, tinyMarkdown},
		{"json tiny", tiny(t), true, tinyJSON},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := Fit(c.s, Options{Budget: 4000, JSON: c.json, DefaultView: c.s.View})
			if err != nil {
				t.Fatalf("Fit: %v", err)
			}
			if got := string(res.Bytes); got != c.want {
				t.Errorf("bytes =\n%s\nwant\n%s", got, c.want)
			}
			if res.Over || len(res.Omitted) != 0 {
				t.Errorf("Over = %v, Omitted = %v; want neither", res.Over, res.Omitted)
			}
			if res.Tokens != Tokens(res.Bytes) {
				t.Errorf("Tokens = %d, want %d", res.Tokens, Tokens(res.Bytes))
			}
			if c.json {
				checkJSON(t, res)
			} else {
				checkMarkdown(t, res.Bytes)
			}
		})
	}
}

func TestRenderFocus(t *testing.T) {
	t.Run("retired with reason", func(t *testing.T) {
		s := tiny(t)
		s.Focus[0].Retired, s.Focus[0].Reason = true, "folded  into\nresume"
		s.Focus[0].Proposed = []string{"retired: folded into resume"}
		res := fit(t, s, Options{Budget: 4000})
		want := "**pause** `01JB2Q` · capability · evidence: none · retired: folded into resume\npause\nsize: 3\nProposed: retired: folded into resume\n"
		if !strings.Contains(string(res.Bytes), want) {
			t.Errorf("focus block missing %q in\n%s", want, res.Bytes)
		}
		res = fit(t, s, Options{Budget: 4000, JSON: true})
		var doc map[string]any
		if err := json.Unmarshal(res.Bytes, &doc); err != nil {
			t.Fatal(err)
		}
		f := doc["focus"].([]any)[0].(map[string]any)
		if f["retired"] != true || f["reason"] != "folded into resume" {
			t.Errorf("focus retired=%v reason=%q", f["retired"], f["reason"])
		}
	})
	t.Run("retired without reason", func(t *testing.T) {
		s := tiny(t)
		s.Focus[0].Retired = true
		res := fit(t, s, Options{Budget: 4000})
		if !strings.Contains(string(res.Bytes), "· evidence: none · retired\n") {
			t.Errorf("no retired mark in\n%s", res.Bytes)
		}
	})
	t.Run("non-string statement field is shown", func(t *testing.T) {
		s := tiny(t)
		s.Focus[0].Fields = fields("statement", "7")
		res := fit(t, s, Options{Budget: 4000})
		if !strings.Contains(string(res.Bytes), "\npause\nstatement: 7\n") {
			t.Errorf("statement field missing in\n%s", res.Bytes)
		}
	})
	t.Run("backticks lengthen the fence", func(t *testing.T) {
		s := tiny(t)
		s.Focus[0].Fields = fields("code", q("a\n````\nb"))
		res := fit(t, s, Options{Budget: 4000})
		if !strings.Contains(string(res.Bytes), "code:\n`````\na\n````\nb\n`````\n") {
			t.Errorf("fence wrong in\n%s", res.Bytes)
		}
	})
	t.Run("several focus items are separate blocks", func(t *testing.T) {
		s := tiny(t)
		other := s.Focus[0]
		other.Handle, other.Short, other.Statement, other.Fields = "resume", "01JB2R", "resume", nil
		s.Focus = append(s.Focus, other)
		res := fit(t, s, Options{Budget: 4000})
		want := "size: 3\n\n**resume** `01JB2R` · capability · evidence: none\nresume\n\n## Constraints"
		if !strings.Contains(string(res.Bytes), want) {
			t.Errorf("missing %q in\n%s", want, res.Bytes)
		}
	})
}

func TestRenderChangeSlice(t *testing.T) {
	f := tree()
	cs := f.begin(changeset.StateDraft, "Tune\n pause")
	f.set(3, "statement", q("pause at 90%"))
	r := f.whatIf(t, cs)
	s := build(t, r, f.asIs(t), Request{Change: f.cs(cs), View: "what-if"})
	if !s.ChangeSlice {
		t.Fatal("ChangeSlice = false for a change slice")
	}
	res := fit(t, s, Options{Budget: 4000, DefaultView: "what-if"})
	head := "# Slice · change " + s.Change.Short + " \"Tune pause\" · what-if\n\n## Focus\n\n**pause** `"
	if !strings.HasPrefix(string(res.Bytes), head) {
		t.Errorf("head =\n%s\nwant prefix\n%s", res.Bytes, head)
	}
	if !strings.Contains(string(res.Bytes), "Proposed: statement: pause does things → pause at 90%\n") {
		t.Errorf("no proposed line in\n%s", res.Bytes)
	}

	ns := build(t, r, f.asIs(t), Request{Focus: nd(3), Change: f.cs(cs), View: "what-if"})
	if ns.ChangeSlice {
		t.Fatal("ChangeSlice = true for a node slice")
	}
	res = fit(t, ns, Options{Budget: 4000, DefaultView: "what-if"})
	head = "# Slice · pause · what-if (change " + s.Change.Short + " \"Tune pause\")\n"
	if !strings.HasPrefix(string(res.Bytes), head) {
		t.Errorf("head =\n%s\nwant prefix\n%s", res.Bytes, head)
	}

	res = fit(t, s, Options{Budget: 4000, DefaultView: "what-if", JSON: true})
	var doc struct {
		Ref    string
		Change map[string]any
	}
	if err := json.Unmarshal(res.Bytes, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Ref != s.Change.Short || doc.Change["id"] != s.Change.ID.String() || doc.Change["state"] != "draft" || doc.Change["intent"] != "Tune pause" {
		t.Errorf("ref = %q, change = %v", doc.Ref, doc.Change)
	}
}

func TestRenderDeterministic(t *testing.T) {
	s := small(t)
	for k := range 30 {
		s.Focus[0].Fields["f"+string(rune('a'+k%26))+strings.Repeat("x", k)] = json.RawMessage(q("v"))
	}
	for _, js := range []bool{false, true} {
		for _, budget := range []int{4000, 300, 1} {
			first := fit(t, s, Options{Budget: budget, JSON: js, DefaultView: "as-is"})
			for range 20 {
				again := fit(t, s, Options{Budget: budget, JSON: js, DefaultView: "as-is"})
				if !bytes.Equal(first.Bytes, again.Bytes) {
					t.Fatalf("json=%v budget=%d: two renders differ:\n%s\n---\n%s", js, budget, first.Bytes, again.Bytes)
				}
			}
		}
	}
}

func TestDefaultView(t *testing.T) {
	cases := []struct {
		cs   *changeset.ChangeSet
		want string
	}{
		{nil, "as-is"},
		{&changeset.ChangeSet{State: changeset.StateDraft}, "what-if"},
		{&changeset.ChangeSet{State: changeset.StateProposed}, "what-if"},
		{&changeset.ChangeSet{State: changeset.StateAccepted}, "to-be"},
		{&changeset.ChangeSet{State: changeset.StateLanded}, "as-is"},
		{&changeset.ChangeSet{State: changeset.StateRejected}, "as-is"},
		{&changeset.ChangeSet{State: changeset.StateWithdrawn}, "as-is"},
		{&changeset.ChangeSet{State: changeset.StateSuperseded}, "as-is"},
	}
	for _, c := range cases {
		if got := DefaultView(c.cs); got != c.want {
			t.Errorf("DefaultView(%v) = %q, want %q", c.cs, got, c.want)
		}
	}
}

func fit(t testing.TB, s Slice, o Options) Result {
	t.Helper()
	if o.DefaultView == "" {
		o.DefaultView = s.View
	}
	res, err := Fit(s, o)
	if err != nil {
		t.Fatalf("Fit: %v", err)
	}
	if o.JSON {
		checkJSON(t, res)
	} else {
		checkMarkdown(t, res.Bytes)
	}
	return res
}

// checkMarkdown checks the byte rules every markdown slice keeps.
func checkMarkdown(t testing.TB, b []byte) {
	t.Helper()
	s := string(b)
	switch {
	case strings.Contains(s, "\r"):
		t.Errorf("markdown has a CR:\n%s", s)
	case !strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\n\n"):
		t.Errorf("markdown does not end with exactly one newline:\n%q", s)
	case strings.Contains(s, "\n\n\n"):
		t.Errorf("markdown has two blank lines in a row:\n%s", s)
	}
	for i, line := range strings.Split(s, "\n") {
		if strings.TrimRight(line, " \t") != line {
			t.Errorf("line %d has trailing white space: %q", i+1, line)
		}
	}
}

// checkJSON checks that JSON output is canonical, that Doc marshals to the
// same bytes, and that its tokens field is the estimate of those bytes.
func checkJSON(t testing.TB, res Result) {
	t.Helper()
	if err := canon.Check(res.Bytes); err != nil {
		t.Errorf("not canonical: %v\n%s", err, res.Bytes)
	}
	if res.Doc == nil {
		t.Fatal("Doc is nil for JSON")
	}
	b, err := canon.Marshal(res.Doc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, res.Bytes) {
		t.Errorf("canon.Marshal(Doc) differs from Bytes")
	}
	var doc struct {
		Tokens     int  `json:"tokens"`
		OverBudget bool `json:"over_budget"`
	}
	if err := json.Unmarshal(res.Bytes, &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Tokens != Tokens(res.Bytes) || doc.Tokens != res.Tokens {
		t.Errorf("tokens field = %d, Result.Tokens = %d, Tokens(bytes) = %d", doc.Tokens, res.Tokens, Tokens(res.Bytes))
	}
	if doc.OverBudget != res.Over {
		t.Errorf("over_budget = %v, Over = %v", doc.OverBudget, res.Over)
	}
}

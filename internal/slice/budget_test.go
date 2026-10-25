package slice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/op"
)

func TestTokens(t *testing.T) {
	for _, c := range []struct{ n, want int }{{0, 0}, {1, 1}, {3, 1}, {4, 1}, {5, 2}, {8, 2}, {9, 3}} {
		if got := Tokens(make([]byte, c.n)); got != c.want {
			t.Errorf("Tokens(%d bytes) = %d, want %d", c.n, got, c.want)
		}
	}
}

// pad pads a ULID prefix with zeros to 26 characters.
func pad(s string) string { return s + strings.Repeat("0", 26-len(s)) }

// big is a node slice with four long items in every section.
func big(t testing.TB) Slice {
	s := small(t)
	s.Context, s.Constraints, s.Decisions, s.Questions, s.Anchors, s.Evidence, s.Children, s.Neighbours = nil, nil, nil, nil, nil, nil, nil, nil
	s.InFlight = nil
	items := func(prefix string, k op.Kind) []Item {
		var out []Item
		for n := range 4 {
			out = append(out, Item{
				ID: ulid(t, pad(fmt.Sprintf("01J%c%02d", prefix[0]-32, n))), Short: fmt.Sprintf("01J%c%02d", prefix[0]-32, n),
				Handle: fmt.Sprintf("%s-%d", prefix, n), Kind: k,
				Statement: fmt.Sprintf("%s number %d, with a statement long enough that briefing it saves more than a footer line costs", prefix, n),
			})
		}
		return out
	}
	s.Context = items("context", op.KindCapability)
	s.Constraints = items("constraint", op.KindConstraint)
	s.Decisions = items("decision", op.KindDecision)
	s.Questions = items("question", op.KindQuestion)
	s.Anchors = items("anchor", op.KindAnchor)
	s.Evidence = items("evidence", op.KindEvidence)
	s.Children = items("child", op.KindCapability)
	s.Neighbours = items("neighbour", op.KindCapability)
	for n := range 4 {
		s.InFlight = append(s.InFlight, Change{ID: ulid(t, pad(fmt.Sprintf("01JB%02d", n))), Short: fmt.Sprintf("01JB%02d", n),
			State: changeset.StateDraft, Intent: fmt.Sprintf("in-flight change number %d, with an intent long enough to matter for the cut order", n)})
	}
	return s
}

// states is the cut sequence Fit walks: all full, then each of cutOrder to
// brief in turn, then each to count, as section key → tier for those not
// full.
func states() []map[string]string {
	order := []string{"in_flight", "neighbours", "children", "evidence", "anchors", "decisions", "context"}
	out := []map[string]string{{}}
	for _, tier := range []string{"brief", "count"} {
		for _, sec := range order {
			next := map[string]string{}
			for k, v := range out[len(out)-1] {
				next[k] = v
			}
			next[sec] = tier
			out = append(out, next)
		}
	}
	return out
}

func stateOf(res Result) map[string]string {
	m := map[string]string{}
	for _, o := range res.Omitted {
		m[o.Section] = o.Tier
	}
	return m
}

func stateIndex(t *testing.T, res Result) int {
	t.Helper()
	got := stateOf(res)
	for i, s := range states() {
		if reflect.DeepEqual(s, got) {
			return i
		}
	}
	t.Fatalf("Omitted %v is not a state of the cut sequence", res.Omitted)
	return -1
}

// kept returns the never-cut part of an output: the Focus, Constraints and
// Open questions blocks in markdown, or those arrays in JSON.
func kept(t *testing.T, b []byte, js bool) string {
	t.Helper()
	if js {
		var doc map[string]json.RawMessage
		if err := json.Unmarshal(b, &doc); err != nil {
			t.Fatal(err)
		}
		return string(doc["focus"]) + string(doc["constraints"]) + string(doc["questions"])
	}
	var out []string
	for _, h := range []string{"## Focus\n", "## Constraints\n", "## Open questions\n"} {
		_, rest, ok := strings.Cut(string(b), h)
		if !ok {
			t.Fatalf("no %q in\n%s", h, b)
		}
		block, _, _ := strings.Cut(rest, "\n## ")
		block, _, _ = strings.Cut(block, "\n---\n")
		out = append(out, h+block)
	}
	return strings.Join(out, "")
}

func TestFit(t *testing.T) {
	for _, js := range []bool{false, true} {
		format := map[bool]string{false: "markdown", true: "json"}[js]
		s := big(t)
		uncut := fit(t, s, Options{Budget: 1 << 30, JSON: js})
		if len(uncut.Omitted) != 0 || uncut.Over {
			t.Fatalf("%s: uncut has cuts: %v", format, uncut.Omitted)
		}
		want := kept(t, uncut.Bytes, js)

		t.Run("cuts in order "+format, func(t *testing.T) {
			// Sweep the budget down; the state may only move forward,
			// and every state is reached.
			last, seen := 0, map[int]bool{}
			for b := uncut.Tokens; b >= 1; b-- {
				res := fit(t, s, Options{Budget: b, JSON: js})
				i := stateIndex(t, res)
				if res.Over {
					i = len(states())
				} else if res.Tokens > b {
					t.Fatalf("budget %d: %d tokens and not over", b, res.Tokens)
				}
				if i < last {
					t.Fatalf("budget %d: state %d after state %d", b, i, last)
				}
				if !seen[i] {
					t.Logf("budget %d: %d tokens, over %v, cut %v", b, res.Tokens, res.Over, res.Omitted)
				}
				last, seen[i] = i, true
				if got := kept(t, res.Bytes, js); got != want {
					t.Fatalf("budget %d: never-cut sections changed:\n%s\nwant\n%s", b, got, want)
				}
			}
			for i := range len(states()) + 1 {
				if !seen[i] {
					t.Errorf("state %d (%v) never reached", i, func() any {
						if i < len(states()) {
							return states()[i]
						}
						return "over"
					}())
				}
			}
		})
	}

	t.Run("footer", func(t *testing.T) {
		s := big(t)
		uncut := fit(t, s, Options{Budget: 1 << 30})
		next := (uncut.Tokens/1000 + 1) * 1000
		// Find the first budget that briefs In flight and Neighbours.
		for b := uncut.Tokens; b >= 1; b-- {
			res := fit(t, s, Options{Budget: b})
			if len(res.Omitted) < 2 {
				continue
			}
			want := fmt.Sprintf("\n\n---\n"+
				"Omitted: neighbours 4 (brief) — `genesis slice hard-pause --budget %d`\n"+
				"Omitted: in flight 4 (brief) — `genesis slice hard-pause --budget %d`\n", next, next)
			if !strings.HasSuffix(string(res.Bytes), want) {
				t.Errorf("footer:\n%s\nwant suffix\n%s", res.Bytes, want)
			}
			if !strings.Contains(string(res.Bytes), "## Neighbours\n\n- neighbour-0 `01JN00`\n- neighbour-1 `01JN01`\n") {
				t.Errorf("neighbours not brief:\n%s", res.Bytes)
			}
			if !strings.Contains(string(res.Bytes), "## In flight\n\n- `01JB00`\n") {
				t.Errorf("in flight not brief:\n%s", res.Bytes)
			}
			return
		}
		t.Fatal("no budget briefs two sections")
	})

	t.Run("commands", func(t *testing.T) {
		node := big(t)
		change := big(t)
		change.ChangeSlice, change.Title = true, `change 01JB7A "Pause at 90%"`
		odd := big(t)
		odd.Focus[0].Handle, odd.Title = "it's odd", "it's odd"
		next := func(s Slice, js bool) int {
			return (fit(t, s, Options{Budget: 1 << 30, JSON: js}).Tokens/1000 + 1) * 1000
		}
		cases := []struct {
			name     string
			s        Slice
			opts     Options
			children string
			other    string
		}{
			{"node default view", node, Options{DefaultView: "to-be"},
				"genesis tree hard-pause", "genesis slice hard-pause --budget %d"},
			{"node other view", node, Options{DefaultView: "as-is"},
				"genesis tree hard-pause", "genesis slice hard-pause --view to-be --budget %d"},
			{"node json", node, Options{DefaultView: "to-be", JSON: true},
				"genesis tree hard-pause", "genesis slice hard-pause --budget %d"},
			{"change default view", change, Options{DefaultView: "to-be"},
				"genesis slice --change 01JB7A --budget %d", "genesis slice --change 01JB7A --budget %d"},
			{"change other view", change, Options{DefaultView: "what-if"},
				"genesis slice --change 01JB7A --view to-be --budget %d", "genesis slice --change 01JB7A --view to-be --budget %d"},
			{"ref given", node, Options{DefaultView: "to-be", Ref: "cap/hard-pause"},
				"genesis tree hard-pause", "genesis slice cap/hard-pause --budget %d"},
			{"quoted", odd, Options{DefaultView: "to-be"},
				`genesis tree 'it'\''s odd'`, `genesis slice 'it'\''s odd' --budget %d`},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				c.opts.Budget = 1
				res := fit(t, c.s, c.opts)
				n := next(c.s, c.opts.JSON)
				if len(res.Omitted) != 7 {
					t.Fatalf("Omitted = %v", res.Omitted)
				}
				for _, o := range res.Omitted {
					want := fmt.Sprintf(c.other, n)
					if o.Section == "children" {
						want = strings.ReplaceAll(c.children, "%d", fmt.Sprint(n))
					}
					if o.Next != want {
						t.Errorf("%s next = %q, want %q", o.Section, o.Next, want)
					}
					if o.Count != 4 || o.Tier != "count" {
						t.Errorf("%s count = %d tier = %q", o.Section, o.Count, o.Tier)
					}
				}
			})
		}
	})

	t.Run("empty sections are never cut", func(t *testing.T) {
		s := tiny(t)
		s.Children = append(s.Children, s.Children[0], s.Children[0], s.Children[0])
		for b := fit(t, s, Options{Budget: 1 << 30}).Tokens; b >= 1; b-- {
			res := fit(t, s, Options{Budget: b})
			for _, o := range res.Omitted {
				if o.Section != "children" {
					t.Fatalf("budget %d: empty section %s cut", b, o.Section)
				}
			}
		}
	})

	t.Run("over budget markdown", func(t *testing.T) {
		s := big(t)
		res := fit(t, s, Options{Budget: 1})
		if !res.Over {
			t.Fatal("Over = false")
		}
		out := string(res.Bytes)
		for _, c := range s.Constraints {
			line := fmt.Sprintf("- %s `%s` — %s\n", c.Handle, c.Short, c.Statement)
			if !strings.Contains(out, line) {
				t.Errorf("constraint %s missing", c.Handle)
			}
		}
		for _, q := range s.Questions {
			if !strings.Contains(out, fmt.Sprintf("- %s `%s` — %s\n", q.Handle, q.Short, q.Statement)) {
				t.Errorf("question %s missing", q.Handle)
			}
		}
		for _, h := range []string{"## Context", "## Decisions", "## Anchors", "## Evidence", "## Children", "## Neighbours", "## In flight"} {
			if strings.Contains(out, h) {
				t.Errorf("%s printed over budget", h)
			}
		}
		want := fmt.Sprintf("Still over budget by %d tokens.\n", res.Tokens-1)
		if !strings.HasSuffix(out, want) {
			t.Errorf("footer:\n%s\nwant suffix %q", out, want)
		}
		if !strings.Contains(out, "\n---\nOmitted: context 4 — `genesis slice hard-pause --budget ") {
			t.Errorf("no context line in footer:\n%s", out)
		}
		if res.Tokens != Tokens(res.Bytes) {
			t.Errorf("Tokens = %d, want %d", res.Tokens, Tokens(res.Bytes))
		}
	})

	t.Run("over budget json", func(t *testing.T) {
		s := big(t)
		res := fit(t, s, Options{Budget: 10, JSON: true})
		var doc struct {
			Budget      int               `json:"budget"`
			Tokens      int               `json:"tokens"`
			OverBudget  bool              `json:"over_budget"`
			Constraints []json.RawMessage `json:"constraints"`
			Questions   []json.RawMessage `json:"questions"`
			Context     []json.RawMessage `json:"context"`
			InFlight    []json.RawMessage `json:"in_flight"`
			Omitted     []Omission        `json:"omitted"`
		}
		if err := json.Unmarshal(res.Bytes, &doc); err != nil {
			t.Fatal(err)
		}
		if !doc.OverBudget || !res.Over || doc.Budget != 10 {
			t.Errorf("over_budget = %v, Over = %v, budget = %d", doc.OverBudget, res.Over, doc.Budget)
		}
		if len(doc.Constraints) != 4 || len(doc.Questions) != 4 || len(doc.Context) != 0 || len(doc.InFlight) != 0 {
			t.Errorf("constraints %d, questions %d, context %d, in flight %d", len(doc.Constraints), len(doc.Questions), len(doc.Context), len(doc.InFlight))
		}
		if len(doc.Omitted) != 7 || doc.Omitted[0].Section != "context" || doc.Omitted[6].Section != "in_flight" {
			t.Errorf("omitted = %v", doc.Omitted)
		}
	})

	t.Run("over budget with nothing to cut", func(t *testing.T) {
		s := tiny(t)
		s.Children = nil
		res := fit(t, s, Options{Budget: 5})
		if !res.Over || len(res.Omitted) != 0 {
			t.Fatalf("Over = %v, Omitted = %v", res.Over, res.Omitted)
		}
		if want := fmt.Sprintf("\n\n---\nStill over budget by %d tokens.\n", res.Tokens-5); !strings.HasSuffix(string(res.Bytes), want) {
			t.Errorf("footer:\n%s\nwant suffix %q", res.Bytes, want)
		}
	})
}

// TestFitWithin checks, for the small test slices, every format and every
// budget, that the output fits the budget unless it is over. TestFit's
// sweeps check the same of big.
func TestFitWithin(t *testing.T) {
	slices := map[string]Slice{"small": small(t), "tiny": tiny(t)}
	for _, name := range []string{"small", "tiny"} {
		s := slices[name]
		for _, js := range []bool{false, true} {
			top := fit(t, s, Options{Budget: 1 << 30, JSON: js}).Tokens + 2
			for b := 1; b <= top; b++ {
				res := fit(t, s, Options{Budget: b, JSON: js})
				if res.Tokens != Tokens(res.Bytes) {
					t.Fatalf("%s json=%v budget %d: Tokens = %d, want %d", name, js, b, res.Tokens, Tokens(res.Bytes))
				}
				if !res.Over && res.Tokens > b {
					t.Fatalf("%s json=%v budget %d: %d tokens and not over", name, js, b, res.Tokens)
				}
				// JSON's "over_budget": true is a byte shorter than false,
				// so an over-budget JSON slice can land on the budget.
				if res.Over && (res.Tokens < b || !js && res.Tokens == b) {
					t.Fatalf("%s json=%v budget %d: over with %d tokens", name, js, b, res.Tokens)
				}
			}
		}
	}
}

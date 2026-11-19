package guide

import (
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// handbackTarget is the on.* target that hands a step back to conduct.
const handbackTarget = "handback:conduct"

// handbackSteps returns the steps of def with an on.* target of
// handback:conduct, sorted.
func handbackSteps(def *workflow.Definition) []string {
	var out []string
	for name, s := range def.Steps {
		for _, target := range s.On {
			if target == handbackTarget {
				out = append(out, name)
				break
			}
		}
	}
	slices.Sort(out)
	return out
}

// routes returns the Routes section of body as its top-level entries,
// "- `<step>`: ...", each with its indented continuation lines, keyed by
// step.
func routes(body string) map[string]string {
	out := map[string]string{}
	step := ""
	for _, l := range strings.Split(sectionOf(body, "## Routes"), "\n") {
		if rest, ok := strings.CutPrefix(l, "- `"); ok {
			name, _, ok := strings.Cut(rest, "`")
			if !ok {
				step = ""
				continue
			}
			step = name
			out[step] = l
			continue
		}
		if step != "" && strings.HasPrefix(l, "  ") {
			out[step] += "\n" + l
			continue
		}
		step = ""
	}
	return out
}

// handbackDrift returns each problem of body against def: a hand-back step
// with no route, a gate command the topic does not name, a route for a step
// that never hands back, a wait-ci route that spawns ci-fix, and a route that
// names rhiza ship rebase (the run's own rebase step runs it; a conductor
// never does).
func handbackDrift(body string, def *workflow.Definition) []string {
	var out []string
	rs := routes(body)
	hb := handbackSteps(def)
	for _, s := range hb {
		if _, ok := rs[s]; !ok {
			out = append(out, "no route for hand-back step "+s)
		}
	}
	var names []string
	for s := range def.Steps {
		names = append(names, s)
	}
	slices.Sort(names)
	for _, s := range names {
		step := def.Steps[s]
		if step.Kind == workflow.KindGate && slices.Contains(hb, s) {
			for _, c := range step.Commands {
				if !strings.Contains(rs[s], "`"+c+"`") {
					out = append(out, "gate step "+s+" runs "+c+", not named in its route")
				}
			}
		}
		if _, ok := rs[s]; ok && !slices.Contains(hb, s) {
			out = append(out, "route for step "+s+", which never hands back")
		}
	}
	if strings.Contains(rs["wait-ci"], "ci-fix") {
		out = append(out, "the wait-ci route names the ci-fix role")
	}
	var routed []string
	for s := range rs {
		routed = append(routed, s)
	}
	slices.Sort(routed)
	for _, s := range routed {
		if strings.Contains(rs[s], "rhiza ship rebase") {
			out = append(out, "the "+s+" route tells the conductor to run rhiza ship rebase")
		}
	}
	return out
}

func TestConductHandback(t *testing.T) {
	c, ok := Lookup("conduct-handback")
	if !ok || c.Skill {
		t.Fatalf("conduct-handback: found %v, skill %v; want a topic", ok, c.Skill)
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct-handback costs %d tokens, cap %d", n, TopicCap)
	}
	def := shipChange(t)
	if hb := handbackSteps(def); len(hb) == 0 {
		t.Fatalf("ship-change hands back from no step")
	}
	for _, p := range handbackDrift(c.Body, def) {
		t.Errorf("conduct-handback: %s", p)
	}
	if def.HandbackCap != workflow.DefaultHandbackCap || workflow.DefaultHandbackCap != 3 {
		t.Errorf("ship-change handback_cap %d, default %d; the topic says 3", def.HandbackCap, workflow.DefaultHandbackCap)
	}
	checkPhrases(t, "conduct-handback", map[string][]string{
		// The section titles RenderHandback writes in internal/run/brief.go.
		"handback.md": {
			"`Failing step`", "`Outcome`", "`## Failing commands`", "`## Output tail`",
			"`## Cause`", "`## Failed CI runs`", "`## Conflicted files`", "`output in <path>`",
		},
		"commands": {
			"`rhiza run build-check`", "`rhiza run review-check`", "`rhiza check`",
			"`rhiza ship gate`", "`rhiza ship open-pr`", "`rhiza ship check`",
			"Never run `rhiza ship rebase` yourself",
		},
		"roles": {"a verifier", "builder retry", "a critic", "a fixer", "the gate-fix role", "the ci-fix role", "the conflict-fix role"},
		"ship rebase": {
			"- `pr-green`:", "- `rebase`:", "An unresolved merge of main",
			"A code conflict left in progress", "failing on the merged head",
			"a gh, fetch or git error", "The merge step never hands back",
		},
		"limits": {"`handback_cap` (3 unless", "parks the run stopped with an inbox item", "never hand back"},
		"finish": {"lands, renders and commits again on the branch with `rhiza land --commit`", "`rhiza session done`"},
		"spawn":  {"never fix anything yourself", "you spawn no role"},
	})
	for _, old := range oldLandWording {
		if strings.Contains(c.Body, old) {
			t.Errorf("conduct-handback: still says %q; land with rhiza land --commit", old)
		}
	}
	if strings.Contains(c.Body, "--no-push") {
		t.Errorf("conduct-handback still names rhiza ship rebase --no-push")
	}
	rs := routes(c.Body)
	for _, s := range []string{"pr-green", "rebase"} {
		for _, role := range []string{"the conflict-fix role", "the gate-fix role"} {
			if !strings.Contains(rs[s], role) {
				t.Errorf("conduct-handback: the %s route names no %s", s, role)
			}
		}
	}
	if !strings.Contains(rs["pr-green"], "the ci-fix role") {
		t.Errorf("conduct-handback: the pr-green route names no ci-fix role")
	}
	if !strings.Contains(rs["rebase"], "spawn no role") {
		t.Errorf("conduct-handback: the rebase route has no no-role case")
	}
	if strings.Contains(c.Body+c.Description, "Claude") {
		t.Errorf("conduct-handback names an agent product")
	}
}

// The conduct topic points to conduct-handback in Read first and stays
// under its cap.
func TestConductReadsHandback(t *testing.T) {
	c, _ := Lookup("conduct")
	if !strings.Contains(sectionOf(c.Body, "## Read first"), "`rhiza guide conduct-handback`") {
		t.Errorf("conduct: Read first does not name rhiza guide conduct-handback")
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct costs %d tokens, cap %d", n, TopicCap)
	}
}

// A hand-back step with no route, a gate command the topic does not name,
// a route for a step that never hands back (merge included), a dropped
// rebase route, ci-fix on wait-ci and a route running rhiza ship rebase are
// each reported by name.
func TestConductHandbackCatchesDrift(t *testing.T) {
	c, _ := Lookup("conduct-handback")

	def := shipChange(t)
	mainCI := def.Steps["main-ci"]
	mainCI.On = map[string]string{"pass": "close-out", "timeout": handbackTarget}
	def.Steps["main-ci"] = mainCI
	if got := handbackDrift(c.Body, def); !slices.Equal(got, []string{"no route for hand-back step main-ci"}) {
		t.Errorf("drift with main-ci handing back = %q", got)
	}

	def = shipChange(t)
	gate := def.Steps["gate"]
	gate.Commands = append(slices.Clone(gate.Commands), "make lint")
	def.Steps["gate"] = gate
	if got := handbackDrift(c.Body, def); !slices.Equal(got, []string{"gate step gate runs make lint, not named in its route"}) {
		t.Errorf("drift with a new gate command = %q", got)
	}

	for _, s := range []string{"merge", "main-ci", "close-out", "ship", "deploy-check"} {
		broken := strings.Replace(c.Body, "\n- `rebase`:", "\n- `"+s+"`: stop.\n- `rebase`:", 1)
		want := []string{"route for step " + s + ", which never hands back"}
		if got := handbackDrift(broken, shipChange(t)); !slices.Equal(got, want) {
			t.Errorf("drift with a %s route = %q, want %q", s, got, want)
		}
	}

	rs := routes(c.Body)
	broken := strings.Replace(c.Body, rs["rebase"]+"\n", "", 1)
	if broken == c.Body {
		t.Fatal("the topic has no rebase route to drop")
	}
	if got := handbackDrift(broken, shipChange(t)); !slices.Equal(got, []string{"no route for hand-back step rebase"}) {
		t.Errorf("drift with the rebase route dropped = %q", got)
	}

	broken = strings.Replace(c.Body, "Change nothing and finish.", "Run `rhiza ship rebase` and finish.", 1)
	if broken == c.Body {
		t.Fatal("the rebase route lost its change-nothing sentence")
	}
	if got := handbackDrift(broken, shipChange(t)); !slices.Equal(got, []string{"the rebase route tells the conductor to run rhiza ship rebase"}) {
		t.Errorf("drift with a route running rhiza ship rebase = %q", got)
	}

	broken = strings.Replace(c.Body, "you spawn no role.", "you spawn the ci-fix role.", 1)
	if broken == c.Body {
		t.Fatal("the wait-ci route lost its no-role sentence")
	}
	if got := handbackDrift(broken, shipChange(t)); !slices.Equal(got, []string{"the wait-ci route names the ci-fix role"}) {
		t.Errorf("drift with ci-fix on wait-ci = %q", got)
	}
}

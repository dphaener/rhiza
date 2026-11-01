package guide

import (
	"strings"
	"testing"
)

// checkPhrases reports each phrase name's text lacks, under the leaf it
// belongs to.
func checkPhrases(t *testing.T, name string, leaves map[string][]string) {
	t.Helper()
	tx, ok := Lookup(name)
	if !ok {
		t.Fatalf("Lookup(%q) found nothing", name)
	}
	for leaf, phrases := range leaves {
		for _, p := range phrases {
			if !strings.Contains(tx.Body, p) {
				t.Errorf("%s (%s): no %q", name, leaf, p)
			}
		}
	}
}

func TestPlanKinds(t *testing.T) {
	plan, ok := Lookup("plan")
	if !ok || !plan.Skill {
		t.Fatalf("plan: found %v, skill %v; want a skill", ok, plan.Skill)
	}
	if !strings.HasPrefix(plan.Body, "# rhiza-plan\n") || !strings.Contains(plan.Body, "\n## Steps\n") {
		t.Errorf("plan: want the # rhiza-plan heading and a ## Steps section")
	}
	// The description says when to load it: planning a milestone or a
	// change set's leaves with the owner, in a plan run.
	for _, w := range []string{"plan run", "plan a milestone", "change set's leaves", "with the owner"} {
		if !strings.Contains(plan.Description, w) {
			t.Errorf("plan description: no %q", w)
		}
	}
	review, ok := Lookup("plan-review")
	if !ok || review.Skill {
		t.Fatalf("plan-review: found %v, skill %v; want a topic", ok, review.Skill)
	}
	// Constraint agents-named-only-in-adapters: no agent product by name.
	for _, tx := range []Text{plan, review} {
		if strings.Contains(tx.Body+tx.Description, "Claude") {
			t.Errorf("%s names an agent product", tx.Name)
		}
	}
}

func TestPlanText(t *testing.T) {
	checkPhrases(t, "plan", map[string][]string{
		"plan-guide": {
			"rhiza apply", `"acceptance=`, `"statement=`, `"type":"decided-by"`,
			"no plan file", "no checkboxes", "rhiza session done",
		},
		"plan-no-intent-yet": {
			"serves no intent yet", "pick a landed intent", `"kind":"intent"`,
			`rhiza link "$RHIZA_CHANGE" serves <intent>`, "before any leaf",
			"`\"serves\": []` in `rhiza change show --json`", "own serves links",
		},
		"plan-interview": {
			"rhiza slice", "rhiza why", "Never ask a fact", "One topic per message",
			"recommendation", "AskUserQuestion", "Only typed words are quoted",
			"`wording`", "`quote`", `"kind":"question"`, `"type":"blocks"`,
		},
		"plan-scope": {
			"`intention`", "`record`", "milestone", "serves", "child intents",
			"rhiza change new", "intent only", "rhiza land", "before review",
			"depends-on", "drafts",
		},
		"plan-autonomous-interview": {
			"RHIZA_AUTONOMOUS=1", "arbiter", "(assumed, would have escalated)",
			"ask the owner nothing", "Quote no owner words",
		},
		"plan-source": {
			"rhiza guide scrub", `--raw-session "$RHIZA_SESSION"`, "rhiza source add",
			"Before you exit",
		},
		"plan-decisions-from-source": {
			"derived-from", "quotes the owner's typed words", "rhiza session interview",
		},
	})
}

func TestPlanReviewText(t *testing.T) {
	checkPhrases(t, "plan-review", map[string][]string{
		"plan-review-guide": {
			"missing-acceptance", "missing-provenance", "unreal-depends-on",
			"outside-intent", "contradiction",
			"rhiza overlaps --change <cs> --json", "naming the other change set",
			"rhiza run findings", `{"findings":[{"node":"`, `"kind":`, `"reason":`, `{"findings":[]}`,
			"rhiza run verdict", `{"verdicts":[{"finding":"f1","verdict":`, "`confirmed` or `dropped`",
			"serving an intent it adds", "rhiza change show", "rhiza slice", "rhiza why",
			"Never add a finding", "confirmed finding", "--record", "rhiza session done",
		},
	})
}

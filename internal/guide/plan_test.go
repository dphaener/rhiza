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

func TestPlanNoMandatoryMilestone(t *testing.T) {
	tx, _ := Lookup("plan")
	if strings.Contains(tx.Body, "builds only a change set that serves a milestone") {
		t.Errorf("plan still says every change set must serve a milestone")
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
		"plan-serves-milestone": {
			"RHIZA_MILESTONE", "only when `RHIZA_MILESTONE` names one",
			"or an open milestone fits", ".rhiza/policy.toml", "`backlog` rule",
			`"skip"`, "umbrella or area intent under no", "milestone",
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

func TestPlanWhoRecords(t *testing.T) {
	checkPhrases(t, "plan", map[string][]string{
		"plan-guide-who-records": {
			"The run's capture records the owner's AskUserQuestion picks",
			"writes no Decision for a plain owner pick",
			"a proxy such as an orchestrator or arbiter answered",
			"the proxy's wording",
			"copy the asked question\nverbatim into its `question` field",
		},
	})
	tx, _ := Lookup("plan")
	for _, line := range strings.Split(tx.Body, "\n") {
		l := strings.ToLower(line)
		if strings.Contains(l, "plain") && strings.Contains(l, "write a decision") {
			t.Errorf("plan tells the planner to write a Decision for a plain pick: %q", line)
		}
	}
}

func TestPlanReviewText(t *testing.T) {
	checkPhrases(t, "plan-review", map[string][]string{
		"plan-review-guide": {
			"missing-acceptance", "missing-provenance", "unreal-depends-on",
			"outside-intent", "contradiction",
			"rhiza overlaps --change <cs> --json", "naming the other change set",
			"rhiza run findings - --json <<'EOF'", "\nEOF\n", `{"findings":[{"node":"`, `"kind":`, `"reason":`, `{"findings":[]}`,
			"rhiza run verdict - --json <<'EOF'", `{"verdicts":[{"finding":"f1","verdict":`, "`confirmed` or `dropped`",
			"serving an intent it adds", "rhiza change show", "rhiza slice", "rhiza why",
			"Never add a finding", "confirmed finding", "--record", "rhiza session done",
		},
	})
}

func TestPlanAcceptanceScoped(t *testing.T) {
	checkPhrases(t, "plan", map[string][]string{
		"plan-acceptance-scoped": {
			"the tests covering that leaf, never the repo's full suite",
		},
	})
	checkPhrases(t, "plan-review", map[string][]string{
		"plan-review-acceptance-scoped": {
			"runs the repo's full suite", "tests covering that leaf",
		},
	})
	for _, name := range []string{"plan", "plan-review"} {
		tx, _ := Lookup(name)
		for _, bad := range []string{"go test", "make check", "make test", "testscript"} {
			if strings.Contains(tx.Body, bad) {
				t.Errorf("%s names %q", name, bad)
			}
		}
	}
}

func TestPlanReviewDoneUnchained(t *testing.T) {
	checkPhrases(t, "plan-review", map[string][]string{
		"plan-review-guide": {"as its own\ncommand only once the recording command", "has succeeded", "Never chain it after one with `;`"},
	})
	tx, _ := Lookup("plan-review")
	if n := Tokens(tx.Body); n > TopicCap {
		t.Errorf("plan-review costs %d tokens, cap %d", n, TopicCap)
	}
}

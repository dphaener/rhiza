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
	if !strings.HasPrefix(plan.Body, "# Plan\n") {
		t.Errorf("plan: want the # Plan heading")
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
	for _, name := range []string{"plan", "planner", "plan-recipes", "plan-interview"} {
		tx, _ := Lookup(name)
		if strings.Contains(tx.Body, "builds only a change set that serves a milestone") {
			t.Errorf("%s still says every change set must serve a milestone", name)
		}
	}
}

// The phrases the plan skill once held, checked where they now live:
// autonomous runs in planner; scope, the apply recipe, no plan file and
// who records a pick in plan-recipes; the interview, Source and milestone
// in plan-interview.
func TestPlanText(t *testing.T) {
	checkPhrases(t, "plan-recipes", map[string][]string{
		"plan-guide": {
			"rhiza apply", `"acceptance=`, `"statement=`, `"type":"decided-by"`,
			"no plan file", "no checkboxes",
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
		"plan-scope": {"`intention`", "`record`", "serves"},
		"plan-decisions-from-source": {
			"derived-from", "quotes the owner's typed words", "rhiza session interview",
		},
	})
	checkPhrases(t, "plan-interview", map[string][]string{
		"plan-interview": {
			"rhiza slice", "rhiza why", "Never ask a fact", "One topic per message",
			"recommendation", "AskUserQuestion", "Only typed words are quoted",
			"`wording`", "`quote`", `"kind":"question"`, `"type":"blocks"`,
		},
		"plan-scope": {
			"milestone", "child intents", "rhiza change new", "intent only", "rhiza land",
			"before review", "depends-on", "drafts",
		},
		"plan-source": {
			"rhiza guide scrub", `--raw-session "$RHIZA_SESSION"`, "rhiza source add",
			"Before you exit",
		},
	})
	checkPhrases(t, "planner", map[string][]string{
		"plan-guide": {"rhiza apply", "rhiza session done"},
		"plan-autonomous-interview": {
			"RHIZA_AUTONOMOUS=1", "arbiter", "(assumed, would have escalated)",
			"ask the owner nothing", "Quote no owner words",
		},
	})
}

func TestPlanWhoRecords(t *testing.T) {
	checkPhrases(t, "plan-recipes", map[string][]string{
		"plan-guide-who-records": {
			"The run's capture records the owner's AskUserQuestion picks",
			"writes no Decision for a plain owner pick",
			"a proxy such as an orchestrator or arbiter answered",
			"the proxy's wording",
			"copy the asked question\nverbatim into its `question` field",
		},
	})
	for _, name := range []string{"planner", "plan-recipes", "plan-interview"} {
		tx, _ := Lookup(name)
		for _, line := range strings.Split(tx.Body, "\n") {
			l := strings.ToLower(line)
			if strings.Contains(l, "plain") && strings.Contains(l, "write a decision") {
				t.Errorf("%s tells the planner to write a Decision for a plain pick: %q", name, line)
			}
		}
	}
}

// The plan and plan-review indexes name the four role topics with both
// routes, and stay short.
func TestPlanIndexes(t *testing.T) {
	roles := []string{"planner", "plan-reviewer", "plan-critic", "plan-reviser"}
	for _, tc := range []struct {
		index string
		roles []string
	}{{"plan", roles}, {"plan-review", roles[1:]}} {
		tx, ok := Lookup(tc.index)
		if !ok {
			t.Fatalf("no %s", tc.index)
		}
		for _, r := range tc.roles {
			if _, ok := Lookup(r); !ok {
				t.Errorf("%s names the %s topic, which does not exist", tc.index, r)
			}
			for _, w := range []string{"`/rhiza-" + r + "`", "`rhiza guide " + r + "`"} {
				if !strings.Contains(tx.Body, w) {
					t.Errorf("%s: no %s", tc.index, w)
				}
			}
		}
		if n := Tokens(tx.Body); n > TopicCap/3 {
			t.Errorf("%s costs %d tokens, over a third of the cap %d", tc.index, n, TopicCap)
		}
		if !strings.Contains(tx.Body, "rhiza session done") {
			t.Errorf("%s: no rhiza session done", tc.index)
		}
	}
	tx, _ := Lookup("plan")
	for _, w := range []string{"/rhiza-planner", "rhiza guide plan-review"} {
		if !strings.Contains(tx.Body, w) {
			t.Errorf("plan: no %q", w)
		}
	}
}

func TestPlanAcceptanceScoped(t *testing.T) {
	checkPhrases(t, "plan-recipes", map[string][]string{
		"plan-acceptance-scoped": {
			"the tests covering that leaf, never the repo's full suite",
		},
	})
	for _, name := range []string{"plan", "plan-review", "planner", "plan-recipes", "plan-interview"} {
		tx, _ := Lookup(name)
		for _, bad := range []string{"go test", "make check", "make test", "testscript"} {
			if strings.Contains(tx.Body, bad) {
				t.Errorf("%s names %q", name, bad)
			}
		}
	}
}

func TestPlanMilestoneNumber(t *testing.T) {
	checkPhrases(t, "plan-interview", map[string][]string{
		"plan-milestone-number": {"RHIZA_MILESTONE_NUMBER", "milestone-N", "milestone=N", "`mN-`", "Never count milestones"},
	})
}

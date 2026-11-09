package guide

import (
	"strings"
	"testing"
)

func TestReviewText(t *testing.T) {
	review, ok := Lookup("review")
	if !ok || review.Skill {
		t.Fatalf("review: found %v, skill %v; want a topic", ok, review.Skill)
	}
	if !strings.Contains(review.Description, "check, critic or fixer of a review run") {
		t.Errorf("review description %q does not say when to load it", review.Description)
	}
	checkPhrases(t, "review", map[string][]string{
		"start": {
			`rhiza slice --change "$RHIZA_CHANGE" --json`, "RHIZA_CHECK", "RHIZA_CYCLE", "RHIZA_ROLE",
			"throwaway checkout",
		},
		"merge-base": {"merge-base", "default branch"},
		"five-checks": {
			"`compliance`", "each leaf's `acceptance`", "`quality`", "`qa`", "built binary",
			"fakes alone do not count", "`scope`", "inside the intent and the leaves", "`model`",
		},
		"fakes-only": {
			"Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`",
		},
		"severities": {"`blocker`", "`major`", "`minor`", "`nit`"},
		"findings": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":`,
			`{"findings":[]}`, "rhiza describe run.findings",
		},
		"verdict": {
			`rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":`,
			"`confirmed` or `dropped`",
		},
		"critic":   {"Never add a finding"},
		"dropped":  {"stay in the run journal", "no Decision"},
		"trailers": {"`Leaf: <leaf>`", "`Review: <run>`", "$RHIZA_RUN", "Never rewrite"},
		"outside-work": {
			"outside the intent", "question", "discovered-from draft change set",
			"rhiza change new", "discovered-from <this>", "rhiza change use <this>",
		},
		"evidence": {
			"Evidence written as yourself on each leaf it exercised",
			`"handle":"<leaf>-qa-<cycle>"`, "`<cycle>` being $RHIZA_CYCLE",
			`"type":"verifies","to":"<leaf>"`,
			"missing constraints",
		},
		"done": {"rhiza session done"},
		"scoped-tests": {
			"Checks, critics and fixers run only the tests covering what the change touches",
			"scoped commands", "never the repo's full\nsuite",
			"the PR's CI runs it before merge",
		},
		"unstable-ci": {
			"While CI is unstable", "not explained by your change",
			"do not push or rerun CI to see", "Reproduce the\nfailing test locally",
			"push only after it and the\ntouched packages pass locally",
		},
		"delay-repro": {
			"injecting a delay at the suspected race", "never by busy loops or parallel stress on the host",
		},
	})
	for _, banned := range []string{"go test", "make check", "make test", "testscript"} {
		if strings.Contains(review.Body, banned) {
			t.Errorf("review names %q", banned)
		}
	}
	// Constraint agents-named-only-in-adapters: no agent product by name.
	if strings.Contains(review.Body+review.Description, "Claude") {
		t.Errorf("review names an agent product")
	}
}

func TestReviewFixerTrailersInOneFinalParagraph(t *testing.T) {
	review, _ := Lookup("review")
	for _, p := range []string{"Co-Authored-By:", "one final paragraph", "no blank line"} {
		if !strings.Contains(strings.ReplaceAll(review.Body, "\n", " "), p) {
			t.Errorf("review fixer: no %q", p)
		}
	}
}

package guide

import (
	"slices"
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
		"merge-base": {"`git diff $(git merge-base origin/main HEAD)`", "default branch"},
		"five-checks": {
			"`compliance`", "each leaf's `acceptance`", "`quality`", "`qa`", "built binary",
			"fakes alone do not count", "`scope`", "inside the intent and the leaves", "`model`",
		},
		"fakes-only": {
			"Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`",
		},
		"severities": {"`blocker`: behaviour broken", "`major`: an acceptance point unmet", "`minor`: style, naming", "`nit`", "`minor` and `nit` are never fixed"},
		"findings": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":`,
			`{"findings":[]}`, "(a check with none fails)", "rhiza describe run.findings",
		},
		"verdict": {
			`rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":`,
			"`confirmed` keeps the severity", "`downgraded` lowers it", "`dropped` removes the finding",
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
		"done": {"rhiza session done", "as its own command", "after its recording\ncommand succeeded", "`;` or `&&`"},
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

func TestReviewFixerPreHandoffChecklist(t *testing.T) {
	review, _ := Lookup("review")
	i := strings.Index(review.Body, "## Fixer")
	j := strings.Index(review.Body, "## Outside work")
	if i < 0 || j < i {
		t.Fatalf("review: no fixer section")
	}
	fixer := review.Body[i:j]
	for _, want := range []string{"reproduce your findings", "`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code"} {
		if !strings.Contains(fixer, want) {
			t.Errorf("review fixer section lacks %q", want)
		}
	}
}

// severityTopics are the topics that grade or fix findings: plan-review,
// review, and every check-*, critic and fixer topic.
func severityTopics(t *testing.T) []Text {
	t.Helper()
	var out []Text
	for _, tx := range List() {
		n := tx.Name
		if n == "plan-review" || n == "review" || strings.HasPrefix(n, "check-") ||
			n == "critic" || n == "fixer" {
			out = append(out, tx)
		}
	}
	for _, want := range []string{"plan-review", "review"} {
		if !slices.ContainsFunc(out, func(tx Text) bool { return tx.Name == want }) {
			t.Fatalf("no %s topic", want)
		}
	}
	return out
}

func TestSeverityBlockOnceInEachTopic(t *testing.T) {
	block := strings.TrimRight(sharedFiles(t)["severity.md"], "\n")
	if block == "" {
		t.Fatal("no severity shared block")
	}
	for _, tx := range severityTopics(t) {
		if n := strings.Count(tx.Body, block); n != 1 {
			t.Errorf("%s holds the severity block %d times, want 1", tx.Name, n)
		}
		if n := strings.Count(tx.Body, "## Severities"); n != 1 {
			t.Errorf("%s has %d severity headings, want 1", tx.Name, n)
		}
	}
	for _, f := range []string{"plan-review.md", "review.md"} {
		b, err := files.ReadFile("text/" + f)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "<!-- shared: severity -->") {
			t.Errorf("%s does not include the severity block", f)
		}
	}
}

func TestSeverityBlockContent(t *testing.T) {
	block := sharedFiles(t)["severity.md"]
	for _, sev := range []string{"`blocker`:", "`major`:", "`minor`:", "`nit`:"} {
		if strings.Count(block, sev) != 1 {
			t.Errorf("severity block defines %q %d times, want 1", sev, strings.Count(block, sev))
		}
	}
	for _, c := range []string{"`compliance`", "`quality`", "`qa`", "`model`"} {
		if !strings.Contains(block, c) {
			t.Errorf("severity mapping lacks %s", c)
		}
	}
	for _, k := range []string{"`missing-acceptance`", "`contradiction`", "`outside-intent`"} {
		if !strings.Contains(block, k) {
			t.Errorf("severity mapping lacks plan review kind %s", k)
		}
	}
	if strings.Contains(block, "`scope`") {
		t.Errorf("severity mapping names a `scope` check")
	}
	for _, out := range []string{"confirmed", "downgraded", "dropped"} {
		if strings.Contains(block, out) {
			t.Errorf("severity block names critic outcome %q", out)
		}
	}
	for _, banned := range []string{"go ", "Go ", "make ", "rhiza check", "rhiza render", "rhiza ship", "testscript"} {
		if strings.Contains(block, banned) {
			t.Errorf("severity block names %q", banned)
		}
	}
	for _, tx := range severityTopics(t) {
		if tx.Name != "plan-review" && tx.Name != "review" && tx.Name != "critic" {
			continue
		}
		i := strings.Index(tx.Body, "## Critic")
		if tx.Name == "critic" {
			i = 0
		}
		if i < 0 {
			t.Fatalf("%s: no critic section", tx.Name)
		}
		critic := tx.Body[i:]
		if j := strings.Index(critic[3:], "\n## "); j >= 0 {
			critic = critic[:j+3]
		}
		for _, want := range []string{"confirmed", "downgraded", "dropped", "Never add a finding"} {
			if !strings.Contains(critic, want) {
				t.Errorf("%s critic section lacks %q", tx.Name, want)
			}
		}
	}
}

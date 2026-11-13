package guide

import (
	"slices"
	"strings"
	"testing"
)

// reviewRoles are the nine review role topics the review index names.
var reviewRoles = []string{"check-compliance", "check-quality", "check-qa", "check-model", "critic", "fixer", "ci-fix", "gate-fix", "conflict-fix"}

// reviewRoleText is the bodies of the nine review role topics, joined.
func reviewRoleText(t *testing.T) Text {
	t.Helper()
	var b strings.Builder
	for _, n := range reviewRoles {
		g, ok := Lookup(n)
		if !ok || g.Skill {
			t.Fatalf("%s: found %v, skill %v; want a topic", n, ok, g.Skill)
		}
		b.WriteString(g.Body)
		b.WriteString("\n")
	}
	return Text{Name: "review roles", Body: b.String()}
}

// TestReviewIsAnIndex: the review topic is a short index naming the nine
// role topics, each by its skill and its rhiza guide route, and routing a
// check named scope to check-compliance.
func TestReviewIsAnIndex(t *testing.T) {
	review, ok := Lookup("review")
	if !ok || review.Skill {
		t.Fatalf("review: found %v, skill %v; want a topic", ok, review.Skill)
	}
	for _, n := range reviewRoles {
		want := "`/rhiza-" + n + "` (`rhiza guide " + n + "`)"
		if !strings.Contains(review.Body, want) {
			t.Errorf("review index: no %q", want)
		}
	}
	checkPhrases(t, "review", map[string][]string{
		"scope": {"A check\n  named `scope` loads this one too"},
		"done":  {"rhiza session done"},
	})
	if n := Tokens(review.Body); n > TopicCap/3 {
		t.Errorf("review index is %d tokens, want under a third of the cap (%d)", n, TopicCap/3)
	}
	for _, banned := range []string{"go test", "make check", "Record agent actor", "rhiza run findings", "Claude"} {
		if strings.Contains(review.Body, banned) {
			t.Errorf("review index holds %q: that belongs in a role topic", banned)
		}
	}
}

// checkRolePhrases fails for each phrase none of the role topics holds.
func checkRolePhrases(t *testing.T, roles Text, leaves map[string][]string) {
	t.Helper()
	for leaf, phrases := range leaves {
		for _, p := range phrases {
			if !strings.Contains(roles.Body, p) {
				t.Errorf("review role topics (%s): no %q", leaf, p)
			}
		}
	}
}

// TestReviewRoleText finds in the role topics what the single review topic
// used to say.
func TestReviewRoleText(t *testing.T) {
	roles := reviewRoleText(t)
	checkRolePhrases(t, roles, map[string][]string{
		"start": {
			`rhiza slice --change "$RHIZA_CHANGE" --json`, "RHIZA_CHECK", "RHIZA_CYCLE", "RHIZA_ROLE",
			"throwaway checkout",
		},
		"merge-base": {"git diff $(git merge-base origin/main HEAD)", "merge-base"},
		"five-checks": {
			"`compliance`", "each leaf's `acceptance`", "`quality`", "`qa`", "fakes alone do not count", "`model`",
		},
		"severities": {"`blocker`: behaviour broken", "`major`: an acceptance point unmet", "`minor`: style, naming", "`nit`", "`minor` and `nit` are never fixed"},
		"findings": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":`,
			`{"findings":[]}`, "(a check with no round fails)", "rhiza describe run.findings",
		},
		"verdict": {
			`rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'`, "\nEOF\n",
			`{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":`,
			"`confirmed`: the defect is real", "`downgraded`: the defect is real", "`dropped`: the defect is not real",
		},
		"critic":   {"Never add a finding"},
		"dropped":  {"stays in the run journal", "Decision"},
		"trailers": {"`Leaf: <leaf>`", "`Review: <run>`", "$RHIZA_RUN", "rewrite"},
		"outside-work": {
			"outside the intent", "question", "discovered-from draft change set",
			"rhiza change new", "discovered-from <this>", "rhiza change use <this>",
		},
		"evidence": {
			"Evidence written as yourself on",
			`"handle":"<leaf>-qa-<cycle>"`,
			`"type":"verifies","to":"<leaf>"`,
			"Missing constraint",
		},
		"done": {"rhiza session done", "as its own command", "after its recording\ncommand succeeded", "`;` or `&&`"},
		"scoped-tests": {
			"only the tests covering what the change touches",
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
}

func TestReviewFixerTrailersInOneFinalParagraph(t *testing.T) {
	for _, name := range []string{"fixer", "ci-fix", "gate-fix", "conflict-fix"} {
		g, _ := Lookup(name)
		for _, p := range []string{"Co-Authored-By:", "one final paragraph", "no blank line"} {
			if !strings.Contains(strings.ReplaceAll(g.Body, "\n", " "), p) {
				t.Errorf("%s: no %q", name, p)
			}
		}
	}
}

func TestReviewFixerPreHandoffChecklist(t *testing.T) {
	fixer, _ := Lookup("fixer")
	for _, want := range []string{"reproduce", "`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code"} {
		if !strings.Contains(fixer.Body, want) {
			t.Errorf("fixer topic lacks %q", want)
		}
	}
}

// severityTopics are the topics that grade or fix findings: every check-*,
// critic and fixer topic, and the plan reviewer and plan critic.
func severityTopics(t *testing.T) []Text {
	t.Helper()
	var out []Text
	for _, tx := range List() {
		n := tx.Name
		if strings.HasPrefix(n, "check-") || n == "critic" || n == "fixer" ||
			n == "plan-reviewer" || n == "plan-critic" {
			out = append(out, tx)
		}
	}
	for _, want := range []string{"plan-reviewer", "plan-critic", "check-compliance", "critic", "fixer"} {
		if !slices.ContainsFunc(out, func(tx Text) bool { return tx.Name == want }) {
			t.Fatalf("no %s topic", want)
		}
	}
	return out
}

func TestSeverityBlockOnceInEachTopic(t *testing.T) {
	block := strings.TrimRight(sharedFiles(t)["severities.md"], "\n")
	if block == "" {
		t.Fatal("no severities shared block")
	}
	for _, tx := range severityTopics(t) {
		if n := strings.Count(tx.Body, block); n != 1 { // matches the nested (demoted) form too
			t.Errorf("%s holds the severities block %d times, want 1", tx.Name, n)
		}
		if n := strings.Count(tx.Body, "## Severities"); n != 1 {
			t.Errorf("%s has %d severity headings, want 1", tx.Name, n)
		}
	}
	for _, f := range []string{"plan-reviewer.md", "plan-critic.md", "check-compliance.md", "check-quality.md"} {
		b, err := files.ReadFile("text/" + f)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "<!-- shared: severities") {
			t.Errorf("%s does not include the severities block", f)
		}
	}
}

func TestSeveritiesBlockShared(t *testing.T) {
	ents, err := files.ReadDir("text")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		raw, err := files.ReadFile("text/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), "## Severities") {
			t.Errorf("%s holds its own ## Severities section", e.Name())
		}
	}
	block := sharedFiles(t)["severities.md"]
	if !strings.HasPrefix(block, "## Severities") {
		t.Fatalf("shared/severities.md does not open with ## Severities")
	}
	for _, tx := range List() {
		raw, _ := files.ReadFile("text/" + tx.Name + ".md")
		if strings.Contains(string(raw), "<!-- shared: severities") && !strings.Contains(tx.Body, "## Severities") {
			t.Errorf("%s names the severities block but does not show it", tx.Name)
		}
	}
	for _, sev := range []string{"`blocker`:", "`major`:", "`minor`:", "`nit`:"} {
		if strings.Count(block, sev) != 1 {
			t.Errorf("severities block defines %q %d times, want 1", sev, strings.Count(block, sev))
		}
	}
	for _, c := range []string{"`compliance`", "`quality`", "`qa`", "`model`"} {
		if !strings.Contains(block, c) {
			t.Errorf("severities mapping lacks %s", c)
		}
	}
	for _, k := range []string{"`missing-acceptance`", "`contradiction`", "`outside-intent`", "`uncovered-intent`", "`oversized-leaf`"} {
		if !strings.Contains(block, k) {
			t.Errorf("severities mapping lacks plan review kind %s", k)
		}
	}
	if !strings.Contains(block, "`uncovered-intent` and `oversized-leaf` major") {
		t.Errorf("severities block does not make uncovered-intent and oversized-leaf major")
	}
	if strings.Contains(block, "`scope`") {
		t.Errorf("severities mapping names a `scope` check")
	}
	for _, out := range []string{"confirmed", "downgraded", "dropped"} {
		if strings.Contains(block, out) {
			t.Errorf("severities block names critic outcome %q", out)
		}
	}
	for _, banned := range []string{"go ", "Go ", "make ", "rhiza check", "rhiza render", "rhiza ship", "testscript"} {
		if strings.Contains(block, banned) {
			t.Errorf("severities block names %q", banned)
		}
	}
	_, err = Expand("x\n<!-- shared: severity -->\n", loadShared)
	if err == nil || !strings.Contains(err.Error(), "severity") {
		t.Errorf("misspelt block: Expand = %v, want an error naming severity", err)
	}
}

func TestReviewCriticSections(t *testing.T) {
	for _, name := range []string{"critic"} {
		tx, _ := Lookup(name)
		i := strings.Index(tx.Body, "## Critic")
		if tx.Name == "critic" {
			i = 0
		}
		if i < 0 {
			t.Fatalf("%s: no critic section", tx.Name)
		}
		critic := tx.Body[i:]
		if j := strings.Index(critic[3:], "\n## "); j >= 0 && tx.Name != "critic" {
			critic = critic[:j+3]
		}
		for _, want := range []string{"confirmed", "downgraded", "dropped", "Never add a finding"} {
			if !strings.Contains(critic, want) {
				t.Errorf("%s critic section lacks %q", tx.Name, want)
			}
		}
	}
}

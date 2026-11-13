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

// severityBlock is a topic's expanded "## Severities" section.
func severityBlock(t *testing.T, name string) string {
	t.Helper()
	tx, ok := Lookup(name)
	if !ok {
		t.Fatalf("%s: no topic", name)
	}
	i := strings.Index(tx.Body, "## Severities")
	if i < 0 {
		t.Fatalf("%s: no severity block", name)
	}
	rest := tx.Body[i+len("## Severities"):]
	if j := strings.Index(rest, "\n## "); j >= 0 {
		rest = rest[:j]
	}
	return rest
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
	shared := sharedFiles(t)["severities.md"]
	if !strings.HasPrefix(shared, "## Severities") {
		t.Fatalf("shared/severities.md does not open with ## Severities")
	}
	want := strings.TrimSuffix(strings.TrimPrefix(shared, "## Severities"), "\n")
	for _, name := range []string{"review"} {
		if got := severityBlock(t, name); strings.TrimSuffix(got, "\n") != want {
			t.Errorf("%s: severity block differs from the shared block:\n%s", name, got)
		}
	}
	for _, tx := range List() {
		raw, _ := files.ReadFile("text/" + tx.Name + ".md")
		if strings.Contains(string(raw), "<!-- shared: severities") && !strings.Contains(tx.Body, "## Severities") {
			t.Errorf("%s names the severities block but does not show it", tx.Name)
		}
	}
	for _, kind := range []string{"`uncovered-intent`", "`oversized-leaf`"} {
		if !strings.Contains(shared, kind) {
			t.Errorf("severities block lacks %s", kind)
		}
	}
	if !strings.Contains(shared, "`uncovered-intent` and `oversized-leaf` major") {
		t.Errorf("severities block does not make uncovered-intent and oversized-leaf major")
	}
	for _, banned := range []string{"go ", "Go ", "make ", "rhiza check", "rhiza render", "rhiza ship", "testscript"} {
		if strings.Contains(shared, banned) {
			t.Errorf("severity block names %q", banned)
		}
	}
	_, err = Expand("x\n<!-- shared: severity -->\n", loadShared)
	if err == nil || !strings.Contains(err.Error(), "severity") {
		t.Errorf("misspelt block: Expand = %v, want an error naming severity", err)
	}
}

func TestReviewCriticSections(t *testing.T) {
	for _, name := range []string{"review"} {
		tx, _ := Lookup(name)
		i := strings.Index(tx.Body, "## Critic")
		if i < 0 {
			t.Fatalf("%s: no critic section", name)
		}
		critic := tx.Body[i:]
		if j := strings.Index(critic[3:], "\n## "); j >= 0 {
			critic = critic[:j+3]
		}
		for _, want := range []string{"confirmed", "downgraded", "dropped", "Never add a finding"} {
			if !strings.Contains(critic, want) {
				t.Errorf("%s critic section lacks %q", name, want)
			}
		}
	}
}

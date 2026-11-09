package guide

import (
	"regexp"
	"strings"
	"testing"
)

func TestBuildText(t *testing.T) {
	build, ok := Lookup("build")
	if !ok || build.Skill {
		t.Fatalf("build: found %v, skill %v; want a topic", ok, build.Skill)
	}
	checkPhrases(t, "build", map[string][]string{
		"start": {
			`rhiza slice --change "$RHIZA_CHANGE" --json`, `rhiza show "$RHIZA_LEAF" --json`,
			"RHIZA_ATTEMPT", "acceptance",
		},
		"trailer": {
			"exactly one new commit", "\nLeaf: export-csv\n", "\nChange: ", "\nActor: ",
			"rhiza run status --json", "fix commit with the same `Leaf:` trailer",
		},
		"outside-work": {
			"raise a question", "linked `blocks` to the leaf",
			"discovered-from draft change set, never built in this run",
			"rhiza change new", "discovered-from <this>", "rhiza change use <this>",
		},
		"evidence-or-findings": {
			"detached checkout", "Never edit code", "Evidence written as yourself",
			`"type":"verifies","to":"<leaf>"`, "rhiza run findings --leaf <leaf> - <<'EOF'", "\n  EOF\n",
			`{"findings":[{"reason":`, "neither Evidence nor findings fails",
			"Evidence\nwritten by anyone else does not count",
		},
		"fakes-only": {"with fakes only, never real Rex,\nagents or `gh`"},
		"done":       {"rhiza session done"},
		"scoped-tests": {
			"Builders and verifiers run only the tests covering what the change touches",
			"scoped commands", "never the repo's full\nsuite",
			"the PR's CI runs it before merge",
		},
		"unstable-ci": {
			"While CI is unstable", "not explained by your change",
			"do not push or rerun CI to see", "Reproduce the\nfailing test locally",
			"push only after it and the\ntouched packages pass locally",
		},
	})
	for _, banned := range []string{"go test", "make check", "make test", "testscript"} {
		if strings.Contains(build.Body, banned) {
			t.Errorf("build names %q", banned)
		}
	}
	// Constraint agents-named-only-in-adapters: no agent product by name.
	if strings.Contains(build.Body+build.Description, "Claude") {
		t.Errorf("build names an agent product")
	}
}

// finalParagraph is the last blank-line-separated paragraph of s.
func finalParagraph(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "\n\n"); i >= 0 {
		return s[i+2:]
	}
	return s
}

func TestBuildCommitsOwnActorRecord(t *testing.T) {
	build, _ := Lookup("build")
	for _, want := range []string{"Record agent actor $RHIZA_ACTOR ", "git add $(", "never another actor's"} {
		if !strings.Contains(build.Body, want) {
			t.Errorf("build: lacks %q", want)
		}
	}
}

func TestBuildTrailersInOneFinalParagraph(t *testing.T) {
	build, _ := Lookup("build")
	i := strings.Index(build.Body, "\nStream CSV exports\n")
	if i < 0 {
		t.Fatal("build: no builder commit example")
	}
	example := build.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	last := finalParagraph(example)
	for _, tr := range []string{"Change: ", "Actor: ", "Leaf: export-csv", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("build example: final paragraph %q lacks %q", last, tr)
		}
	}
}

func TestGuidesFeedFindingsFromStdin(t *testing.T) {
	re := regexp.MustCompile(`rhiza run (findings|verdict)[^\n]*`)
	for _, name := range []string{"build", "review", "plan-review"} {
		g, ok := Lookup(name)
		if !ok {
			t.Fatalf("no guide %q", name)
		}
		body := g.Body
		hits := re.FindAllString(body, -1)
		for _, h := range hits {
			if strings.Contains(h, "<file>") {
				t.Errorf("%s: file form remains: %q", name, h)
			}
		}
		if strings.Count(body, "<<'EOF'") == 0 {
			t.Errorf("%s: no stdin heredoc form", name)
		}
	}
}

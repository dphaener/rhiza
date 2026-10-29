package guide

import (
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
			`"type":"verifies","to":"<leaf>"`, "rhiza run findings --leaf <leaf> <file>",
			`{"findings":[{"reason":`, "neither Evidence nor findings fails",
			"Evidence\nwritten by anyone else does not count",
		},
		"done": {"rhiza session done"},
	})
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

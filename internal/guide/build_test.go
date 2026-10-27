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

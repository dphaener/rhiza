package guide

import (
	"strings"
	"testing"
)

func TestBuilderText(t *testing.T) {
	b, ok := Lookup("builder")
	if !ok || b.Skill {
		t.Fatalf("builder: found %v, skill %v; want a topic", ok, b.Skill)
	}
	at := -1
	for _, h := range []string{"## Task\n", "## Expected outcome\n", "## Read first\n", "## Must do\n", "## Must not do\n", "## Context\n"} {
		i := strings.Index(b.Body, h)
		if i < 0 || i < at {
			t.Fatalf("builder: heading %q missing or out of order", strings.TrimSpace(h))
		}
		at = i
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "unstable-ci", "command-budget", "outside-work", "never-list"} {
		if !strings.Contains(b.Body, strings.TrimRight(shared[blk+".md"], "\n")) {
			t.Errorf("builder lacks the %s block", blk)
		}
	}
	checkPhrases(t, "builder", map[string][]string{
		"rules": {
			"The leaf's `acceptance` is the contract", "Read the slice, the leaf and the existing code and tests",
			"before you write anything", "raise a question before you start work, not midway",
		},
		"checklist": {
			"acceptance\ncommands", "`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code",
			"exactly one new commit",
		},
		"actor": {"Record agent actor $RHIZA_ACTOR ", "git add $(", "never another actor's"},
		"retry": {"RHIZA_ATTEMPT` > 1", "/rhiza-build-retry", "`rhiza guide build-retry`"},
	})
	i := strings.Index(b.Body, "\nStream CSV exports\n")
	if i < 0 {
		t.Fatal("builder: no commit example")
	}
	example := b.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	c, ch := strings.Index(example, "\nChecks:\n"), strings.Index(example, "\nChange: ")
	if c < 0 || ch < c {
		t.Errorf("builder example: Checks: block must sit above the trailers")
	}
	last := finalParagraph(example)
	if strings.Contains(last, "Checks:") {
		t.Errorf("builder example: Checks: inside the trailer paragraph")
	}
	for _, tr := range []string{"Change: ", "Actor: ", "Leaf: export-csv", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("builder example: final paragraph lacks %q", tr)
		}
	}
	for _, banned := range []string{"go test", "go build", "go vet", "make ", "testscript", "Claude"} {
		if strings.Contains(b.Body+b.Description, banned) {
			t.Errorf("builder names %q", banned)
		}
	}
}

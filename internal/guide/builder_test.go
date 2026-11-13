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
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(b.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("builder: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "unstable-ci", "command-budget", "outside-work", "never-list"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(b.Body, text) && !strings.Contains(b.Body, demote(text)) {
			t.Errorf("builder lacks the %s block", blk)
		}
	}
	must, never := sectionOf(b.Body, "## Must do"), sectionOf(b.Body, "## Must not do")
	for _, blk := range []string{"command-budget", "outside-work"} {
		if !strings.Contains(must, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("builder: the %s block is not inside Must do", blk)
		}
	}
	if !strings.Contains(never, demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("builder: the never-list block is not inside Must not do")
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

// sectionOf returns the text of body's section under heading, up to the
// next `## ` heading.
func sectionOf(body, heading string) string {
	i := strings.Index(body, "\n"+heading+"\n")
	if i < 0 {
		return ""
	}
	rest := body[i+len(heading)+2:]
	if j := strings.Index(rest, "\n## "); j >= 0 {
		rest = rest[:j]
	}
	return rest
}

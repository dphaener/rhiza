package guide

import (
	"strings"
	"testing"
)

func TestCIFixText(t *testing.T) {
	c, ok := Lookup("ci-fix")
	if !ok || c.Skill {
		t.Fatalf("ci-fix: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("ci-fix: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "unstable-ci", "command-budget", "never-list"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("ci-fix lacks the %s block", blk)
		}
		if strings.Contains(c.Body, "<!-- shared:") {
			t.Errorf("ci-fix: an unexpanded shared marker")
		}
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("ci-fix: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "ci-fix", map[string][]string{
		"findings":   {"rhiza ship check --json", "`tests`", "`findings`", "failing log lines"},
		"reproduce":  {"Reproduce each finding locally first", "scoped command", "Never push to see what CI does"},
		"regenerate": {"Stale generated files: regenerate them"},
		"flake":      {"flake", "quarantine it (skip it)", "rhiza request --kind bug", "named in a comment at the skip"},
		"checklist":  {"`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code", "exactly one new commit"},
		"trailer":    {"`Review: <run>`"},
	})
	if !strings.Contains(sectionOf(c.Body, "## Must do"), "Reproduce each finding locally first") ||
		!strings.Contains(sectionOf(c.Body, "## Must not do"), "Push, or rerun CI, to see") {
		t.Errorf("ci-fix: reproduce-first belongs in Must do, push-to-see in Must not do")
	}
	i := strings.Index(c.Body, "\nFix the CI findings of the export run\n")
	if i < 0 {
		t.Fatal("ci-fix: no commit example")
	}
	example := c.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	ck, ch := strings.Index(example, "\nChecks:\n"), strings.Index(example, "\nChange: ")
	if ck < 0 || ch < ck {
		t.Errorf("ci-fix example: Checks: block must sit above the trailers")
	}
	last := finalParagraph(example)
	if strings.Contains(last, "Checks:") {
		t.Errorf("ci-fix example: Checks: inside the trailer paragraph")
	}
	for _, tr := range []string{"Change: ", "Actor: ", "Review: ", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("ci-fix example: final paragraph lacks %q", tr)
		}
	}
	if strings.Contains(last, "Leaf: ") {
		t.Errorf("ci-fix example: a Leaf: trailer; CI findings name no leaf")
	}
}

// TestCIFixNamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestCIFixNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("ci-fix")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("ci-fix names %q", banned)
		}
	}
}

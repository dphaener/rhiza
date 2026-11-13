package guide

import (
	"strings"
	"testing"
)

func TestGateFixText(t *testing.T) {
	c, ok := Lookup("gate-fix")
	if !ok || c.Skill {
		t.Fatalf("gate-fix: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("gate-fix: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "unstable-ci", "command-budget", "never-list"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("gate-fix lacks the %s block", blk)
		}
		if strings.Contains(c.Body, "<!-- shared:") {
			t.Errorf("gate-fix: an unexpanded shared marker")
		}
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("gate-fix: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "gate-fix", map[string][]string{
		"settings":  {"`rhiza.ship.localGate`", "`rhiza.ship.prePushCheck`"},
		"read":      {"rhiza run status --json", "`Full log:`", "kept output"},
		"reproduce": {"Never rerun `rhiza ship gate`", "scoped command", "fail again or plainly pass"},
		"flake":     {"flake", "quarantine it (skip it)", "rhiza request --kind bug", "named in a comment at the skip"},
		"checklist": {"`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code", "exactly one new commit"},
		"trailer":   {"`Review: <run>`"},
	})
	if !strings.Contains(sectionOf(c.Body, "## Must do"), "Never rerun `rhiza ship gate`") ||
		!strings.Contains(sectionOf(c.Body, "## Must not do"), "Rerun `rhiza ship gate`, or push, to see") {
		t.Errorf("gate-fix: never-rerun belongs in Must do and Must not do")
	}
	i := strings.Index(c.Body, "\nFix the failed gate of the export run\n")
	if i < 0 {
		t.Fatal("gate-fix: no commit example")
	}
	example := c.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	ck, ch := strings.Index(example, "\nChecks:\n"), strings.Index(example, "\nChange: ")
	if ck < 0 || ch < ck {
		t.Errorf("gate-fix example: Checks: block must sit above the trailers")
	}
	last := finalParagraph(example)
	if strings.Contains(last, "Checks:") {
		t.Errorf("gate-fix example: Checks: inside the trailer paragraph")
	}
	for _, tr := range []string{"Change: ", "Actor: ", "Review: ", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("gate-fix example: final paragraph lacks %q", tr)
		}
	}
	if strings.Contains(last, "Leaf: ") {
		t.Errorf("gate-fix example: a Leaf: trailer; a gate failure names no leaf")
	}
}

// TestGateFixNamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestGateFixNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("gate-fix")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("gate-fix names %q", banned)
		}
	}
}

package guide

import (
	"strings"
	"testing"
)

func TestConflictFixText(t *testing.T) {
	c, ok := Lookup("conflict-fix")
	if !ok || c.Skill {
		t.Fatalf("conflict-fix: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("conflict-fix: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("conflict-fix lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("conflict-fix: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("conflict-fix: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "conflict-fix", map[string][]string{
		"code":    {"rhiza ship check --json", "`merge_conflict`", "`git add -A -- .rhiza/rendered .rhiza/log`", "`rhiza ship generate`", "`rhiza render`", "`Checks:` block"},
		"model":   {"rhiza ship check --model-conflict --json", "`rhiza set <node> --rebase`", "`Model-Conflict: <other change set>`"},
		"trailer": {"`Change:`", "`Actor:`", "`Review: <run>`"},
		"clash":   {"commit nothing", "one blocking question", "names both change sets"},
		"never":   {"End your session with the merge still in progress", "Push, rebase or force anything"},
	})
	i := strings.Index(c.Body, "\nMerge origin/main into the export branch\n")
	if i < 0 {
		t.Fatal("conflict-fix: no merge commit example")
	}
	example := c.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	if ck, ch := strings.Index(example, "\nChecks:\n"), strings.Index(example, "\nChange: "); ck < 0 || ch < ck {
		t.Errorf("conflict-fix example: Checks: block must sit above the trailers")
	}
	last := finalParagraph(example)
	for _, tr := range []string{"Change: ", "Actor: ", "Review: ", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("conflict-fix example: final paragraph lacks %q", tr)
		}
	}
}

// TestConflictFixNamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestConflictFixNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("conflict-fix")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("conflict-fix names %q", banned)
		}
	}
}

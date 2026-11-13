package guide

import (
	"strings"
	"testing"
)

func TestFixerText(t *testing.T) {
	f, ok := Lookup("fixer")
	if !ok || f.Skill {
		t.Fatalf("fixer: found %v, skill %v; want a topic", ok, f.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(f.Body, "\n") {
		if strings.HasPrefix(l, "## ") && !strings.HasPrefix(l, "## Severities") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("fixer: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "unstable-ci", "command-budget", "outside-work", "never-list", "severities"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(f.Body, text) && !strings.Contains(f.Body, demote(text)) {
			t.Errorf("fixer lacks the %s block", blk)
		}
	}
	if strings.Contains(f.Body, "<!-- shared:") {
		t.Errorf("fixer: an unexpanded shared marker")
	}
	must, never := sectionOf(f.Body, "## Must do"), sectionOf(f.Body, "## Must not do")
	for _, blk := range []string{"command-budget", "outside-work"} {
		if !strings.Contains(must, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("fixer: the %s block is not inside Must do", blk)
		}
	}
	if !strings.Contains(never, demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("fixer: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "fixer", map[string][]string{
		"only":      {"Fix only the findings your prompt names", "each exactly as written and in place"},
		"reproduce": {"Reproduce each finding first", "see it fail"},
		"checklist": {"`rhiza ship generate`", "`rhiza render`", "`Checks:` block", "exit code", "exactly one new commit"},
		"trailers":  {"`Leaf: <leaf>`", "`Review: <run>`", "every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line"},
		"unfixed":   {"A finding you could not fix is not dropped", "Drop a finding you could not fix: carry it to the end as unfixed", "Unfixed:"},
		"never":     {"Push, rebase, force or use `--no-verify`"},
	})
	if !strings.Contains(must, "Reproduce each finding first") || !strings.Contains(never, "Fix a finding your prompt does not name") {
		t.Errorf("fixer: reproduce-first belongs in Must do, fix-only-named in Must not do")
	}
	i := strings.Index(f.Body, "\nFix the findings of the export review\n")
	if i < 0 {
		t.Fatal("fixer: no commit example")
	}
	example := f.Body[i:]
	example = example[:strings.Index(example, "\n```")]
	ck, ch := strings.Index(example, "\nChecks:\n"), strings.Index(example, "\nChange: ")
	if ck < 0 || ch < ck {
		t.Errorf("fixer example: Checks: block must sit above the trailers")
	}
	last := finalParagraph(example)
	if strings.Contains(last, "Checks:") {
		t.Errorf("fixer example: Checks: inside the trailer paragraph")
	}
	for _, tr := range []string{"Change: ", "Actor: ", "Leaf: ", "Co-Authored-By: "} {
		if !strings.Contains(last, tr) {
			t.Errorf("fixer example: final paragraph lacks %q", tr)
		}
	}
}

// TestFixerNamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestFixerNamesNoRepoCommand(t *testing.T) {
	f, _ := Lookup("fixer")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude"} {
		if strings.Contains(f.Body+f.Description, banned) {
			t.Errorf("fixer names %q", banned)
		}
	}
}

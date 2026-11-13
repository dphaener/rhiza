package guide

import (
	"strings"
	"testing"
)

func TestCheckQAText(t *testing.T) {
	c, ok := Lookup("check-qa")
	if !ok || c.Skill {
		t.Fatalf("check-qa: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("check-qa: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list", "severity"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("check-qa lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("check-qa: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("check-qa: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "check-qa", map[string][]string{
		"isolation": {
			"You know nothing of how the work was built",
			"Never run `rhiza run status` or\nrhiza run show (or tail) on the build run",
			"its findings, its run notes or its\nsession transcripts",
			"the contract from the slice, the diff from the\nmerge-base and the model's Evidence",
		},
		"clean":     {"Start from a clean state"},
		"scratch":   {"to a scratch path outside `$PATH`", "Never install a\nbinary onto `$PATH`"},
		"scenarios": {"- Happy:", "- Failure:", "fakes alone do not count"},
		"explore":   {"Pairs first", "Then edge cases", "write no new test files"},
		"broken":    {"only for behaviour that is broken, never for a missing test or a style\npoint"},
		"evidence": {
			`"handle":"<leaf>-qa-<cycle>"`,
			`"type":"verifies"`,
			"rhiza apply <file>",
		},
		"record": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`,
			`{"findings":[]}`,
		},
	})
	if !strings.Contains(sectionOf(c.Body, "## Task"), "Never run `rhiza run status`") {
		t.Errorf("check-qa: the isolation rule belongs in Task")
	}
}

// TestCheckQANamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestCheckQANamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("check-qa")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude", "Co-Authored-By", "git commit"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("check-qa names %q", banned)
		}
	}
}

package guide

import (
	"strings"
	"testing"
)

func TestCheckQualityText(t *testing.T) {
	c, ok := Lookup("check-quality")
	if !ok || c.Skill {
		t.Fatalf("check-quality: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("check-quality: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list", "severity"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("check-quality lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("check-quality: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("check-quality: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "check-quality", map[string][]string{
		"isolation": {
			"You know nothing of how the work was built",
			"Never run `rhiza run status` or\nrhiza run show (or tail) on the build run",
			"its findings, its run notes or its\nsession transcripts",
			"the contract from the slice, the diff from the\nmerge-base and the model's Evidence",
		},
		"discovery": {"find the repo's test commands in its contributor guide", "discover them there and\nnever guess them"},
		"radius":    {"the impact radius", "- Direct:", "- Integration:", "- Downstream:"},
		"no-full":   {"Never silently\nfall back to the full suite"},
		"checklist": {
			"Swallowed errors", "Debug output", "Commented-out code", "Unused imports",
			"TODOs, stubs or placeholders", "Generic names", "Missing error handling",
		},
		"record": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`,
			`{"findings":[]}`,
		},
	})
	if !strings.Contains(sectionOf(c.Body, "## Task"), "Never run `rhiza run status`") {
		t.Errorf("check-quality: the isolation rule belongs in Task")
	}
	if n := strings.Count(sectionOf(c.Body, "## Must do"), "\n- "); n < 10 {
		t.Errorf("check-quality: Must do has %d bullets, want the radius and the seven checklist items", n)
	}
}

// TestCheckQualityNamesNoRepoCommand is the negative case: the topic serves
// any repo, so it names no Go, make or Rhiza-repo command and gives no commit
// or trailer instruction.
func TestCheckQualityNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("check-quality")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude", "Checks:", "trailer", "Co-Authored-By", "Leaf:", "git commit"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("check-quality names %q", banned)
		}
	}
}

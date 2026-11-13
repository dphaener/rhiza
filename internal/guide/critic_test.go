package guide

import (
	"strings"
	"testing"
)

func TestCriticText(t *testing.T) {
	c, ok := Lookup("critic")
	if !ok || c.Skill {
		t.Fatalf("critic: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") && !strings.HasPrefix(l, "## Severities") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("critic: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "command-budget", "never-list", "severity"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("critic lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("critic: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("critic: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "critic", map[string][]string{
		"isolation": {
			"You know nothing of how the work was built",
			"Never run `rhiza run status` or\nrhiza run show (or tail) on the build run",
			"its findings, its run notes or its\nsession transcripts",
			"the contract from the slice, the diff from the\nmerge-base and the model's Evidence",
			"The one exception is your own check's last round",
		},
		"verdicts": {"`confirmed`", "`downgraded`", "`dropped`", "severity too high", "lower severity"},
		"never":    {"Never add a finding, raise a severity or rewrite one", "Add a finding, raise a severity or rewrite a finding"},
		"record": {
			`rhiza run verdict --check "$RHIZA_CHECK" - <<'EOF'`,
			`{"verdicts":[`,
		},
	})
	if !strings.Contains(sectionOf(c.Body, "## Task"), "Never run `rhiza run status`") {
		t.Errorf("critic: the isolation rule belongs in Task")
	}
}

// TestCriticRecordsNoFindingAndNoCommit is the negative case: a critic adds no
// finding and commits nothing.
func TestCriticRecordsNoFindingAndNoCommit(t *testing.T) {
	c, _ := Lookup("critic")
	for _, banned := range []string{"rhiza run findings", "git commit", "Checks:", "Leaf:", "Review:", "Co-Authored-By", "go test", "make ", "-trimpath", "internal/", "cmd/rhiza", "Claude"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("critic names %q", banned)
		}
	}
}

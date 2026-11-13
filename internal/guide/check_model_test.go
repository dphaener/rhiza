package guide

import (
	"strings"
	"testing"
)

func TestCheckModelText(t *testing.T) {
	c, ok := Lookup("check-model")
	if !ok || c.Skill {
		t.Fatalf("check-model: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("check-model: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "command-budget", "outside-work", "never-list", "severity"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("check-model lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("check-model: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("check-model: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "check-model", map[string][]string{
		"isolation": {
			"You know nothing of how the work was built",
			"Never run `rhiza run status` or\nrhiza run show (or tail) on the build run",
			"its findings, its run notes or its\nsession transcripts",
			"the contract from the slice, the diff from the\nmerge-base and the model's Evidence",
		},
		"grading": {"Code contradicts a node", "Missing constraint", "is a `major`", "is a `minor`"},
		"outside": {
			"`rhiza change new \"<intent>\" --json`",
			"`rhiza link <new> discovered-from <this> --change <new> --json`",
			"never built in this run",
		},
		"record": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`,
			`{"findings":[]}`,
		},
	})
	if !strings.Contains(sectionOf(c.Body, "## Task"), "Never run `rhiza run status`") {
		t.Errorf("check-model: the isolation rule belongs in Task")
	}
}

// TestCheckModelNamesNoRepoCommand is the negative case: the topic serves any
// repo, so it names no Go, make or Rhiza-repo command.
func TestCheckModelNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("check-model")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude", "Co-Authored-By", "git commit"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("check-model names %q", banned)
		}
	}
}

package guide

import (
	"strings"
	"testing"
)

func TestCheckPlanText(t *testing.T) {
	c, ok := Lookup("check-plan")
	if !ok || c.Skill {
		t.Fatalf("check-plan: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("check-plan: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list", "severities"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("check-plan lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("check-plan: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("check-plan: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "check-plan", map[string][]string{
		"isolation": {
			"You know nothing of how the work was built",
			"Never run `rhiza run status` or\nrhiza run show (or tail) on the build run",
			"its findings, its run notes or its\nsession transcripts",
			"the contract from the slice, the diff from the\nmerge-base and the model's Evidence",
		},
		"criteria": {
			"Under-built: a leaf with nothing in the diff for it is a `blocker`",
			"an acceptance point\n  that is unmet or only partly met is a `major`",
			"Scope creep: behaviour outside the intent and the leaves is a `major`",
			"Unaccounted files", "is a\n  `major` if it changes behaviour and a `minor` otherwise",
		},
		"model": {
			"code against the slice's constraints and Decisions",
			"Code contradicts a node: code that does what a leaf, Decision or constraint forbids, or\n  the reverse, is a `major`",
			"Missing constraint: a rule the change introduced that no constraint states is a `minor`",
		},
		"line":  {"one line with its evidence", "re-diagnosing"},
		"fakes": {"Test with fakes only, never real Rex, agents or `gh`"},
		"record": {
			`rhiza run findings --check "$RHIZA_CHECK" - <<'EOF'`,
			`{"findings":[{"severity":"major","leaf":"<leaf>","place":"path:line","reason":`,
			`{"findings":[]}`,
		},
	})
	if !strings.Contains(sectionOf(c.Body, "## Task"), "Never run `rhiza run status`") {
		t.Errorf("check-plan: the isolation rule belongs in Task")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), "the code check does that") {
		t.Errorf("check-plan: Must not do does not leave style to the code check")
	}
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("check-plan costs %d tokens, cap %d", n, TopicCap)
	}
}

// TestCheckPlanNamesNoRepoCommand is the negative case: the topic serves
// any repo, so it names no Go, make or Rhiza-repo command and gives no commit
// or trailer instruction.
func TestCheckPlanNamesNoRepoCommand(t *testing.T) {
	c, _ := Lookup("check-plan")
	for _, banned := range []string{"go test", "go build", "go vet", "gofmt", "make ", "testscript", "-trimpath", "internal/", "cmd/rhiza", "tools/gate", "Claude", "Checks:", "trailer", "Co-Authored-By", "Leaf:", "git commit"} {
		if strings.Contains(c.Body+c.Description, banned) {
			t.Errorf("check-plan names %q", banned)
		}
	}
}

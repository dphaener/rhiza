package guide

import (
	"strings"
	"testing"
)

func TestVerifierText(t *testing.T) {
	v, ok := Lookup("verifier")
	if !ok || v.Skill {
		t.Fatalf("verifier: found %v, skill %v; want a topic", ok, v.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(v.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("verifier: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(v.Body, text) && !strings.Contains(v.Body, demote(text)) {
			t.Errorf("verifier lacks the %s block", blk)
		}
	}
	if !strings.Contains(sectionOf(v.Body, "## Must do"), demote(strings.TrimRight(shared["command-budget.md"], "\n"))) {
		t.Errorf("verifier: the command-budget block is not inside Must do")
	}
	if !strings.Contains(sectionOf(v.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("verifier: the never-list block is not inside Must not do")
	}
	checkPhrases(t, "verifier", map[string][]string{
		"checkout": {"detached checkout at the leaf's commit", "Never edit code", "fakes only,\nnever real Rex, agents or `gh`"},
		"reruns": {
			"Re-run every command in the commit's `Checks:` block", "every acceptance command",
			"A different exit code, or no `Checks:` block, is a finding naming the command",
		},
		"scope": {"each acceptance point is met and tested", "nothing\noutside the leaf changed", "no stubs or TODOs"},
		"evidence": {
			"rhiza apply <file> --change \"$RHIZA_CHANGE\" --json", `"kind":"evidence"`, `"type":"verifies"`,
			"result=PASS", "Evidence written by anyone else does not count",
		},
		"findings": {
			"rhiza run findings --leaf <leaf> - <<'EOF'", "Each finding names the file, command or acceptance point",
			"without re-diagnosing",
		},
	})
	for _, banned := range []string{"Checks:\n- <acceptance", "Co-Authored-By", "Actor: ", "Leaf: ", "Record agent actor", "git add", "ship generate", "one new commit"} {
		if strings.Contains(v.Body, banned) {
			t.Errorf("verifier carries builder instruction %q", banned)
		}
	}
}

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
	at := -1
	for _, h := range []string{"## Task\n", "## Expected outcome\n", "## Read first\n", "## Must do\n", "## Must not do\n", "## Context\n"} {
		i := strings.Index(v.Body, h)
		if i < 0 || i < at {
			t.Fatalf("verifier: heading %q missing or out of order", strings.TrimSpace(h))
		}
		at = i
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "test-scope", "command-budget", "never-list"} {
		if !strings.Contains(v.Body, strings.TrimRight(shared[blk+".md"], "\n")) {
			t.Errorf("verifier lacks the %s block", blk)
		}
	}
	checkPhrases(t, "verifier", map[string][]string{
		"checkout": {"detached checkout at the leaf's commit", "Never edit code", "fakes only"},
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

package guide

import (
	"strings"
	"testing"
)

func TestConduct(t *testing.T) {
	c, ok := Lookup("conduct")
	if !ok || c.Skill {
		t.Fatalf("conduct: found %v, skill %v; want a topic", ok, c.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(c.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("conduct: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	for _, blk := range []string{"session-done", "command-budget", "never-list", "contract-ref"} {
		text := strings.TrimRight(shared[blk+".md"], "\n")
		if !strings.Contains(c.Body, text) && !strings.Contains(c.Body, demote(text)) {
			t.Errorf("conduct lacks the %s block", blk)
		}
	}
	if strings.Contains(c.Body, "<!-- shared:") {
		t.Errorf("conduct: an unexpanded shared marker")
	}
	if !strings.Contains(sectionOf(c.Body, "## Must not do"), demote(strings.TrimRight(shared["never-list.md"], "\n"))) {
		t.Errorf("conduct: the never-list block is not inside Must not do")
	}
	must := sectionOf(c.Body, "## Must do")
	last := -1
	for _, phase := range []string{
		"`rhiza apply`", "--role builder", "--role verifier", "--role check-plan", "--role check-code",
		"--role check-qa", "--role critic", "--role fixer", "`rhiza run review-check`", "`rhiza land`",
		"`rhiza render`", "`rhiza session done`",
	} {
		i := strings.Index(must, phase)
		if i < 0 {
			t.Errorf("conduct: Must do lacks %s", phase)
			continue
		}
		if i < last {
			t.Errorf("conduct: %s is out of phase order", phase)
		}
		last = i
	}
	checkPhrases(t, "conduct", map[string][]string{
		"findings":  {"`rhiza run findings --check`"},
		"arbiter":   {"a question goes to the arbiter"},
		"critic":    {"only for a check whose latest round has a blocker or\n  major finding"},
		"never":     {"Merge or push", "Edit `.rhiza/log`", "Pass your own gate"},
		"contract":  {"rhiza guide contract"},
		"spawn-tag": {"--brief <file>"},
	})
	if n := Tokens(c.Body); n > TopicCap {
		t.Errorf("conduct costs %d tokens, cap %d", n, TopicCap)
	}
}

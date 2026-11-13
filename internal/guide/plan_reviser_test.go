package guide

import (
	"strings"
	"testing"
)

func TestPlanReviserText(t *testing.T) {
	p, ok := Lookup("plan-reviser")
	if !ok || p.Skill {
		t.Fatalf("plan-reviser: found %v, skill %v; want a topic", ok, p.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(p.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("plan-reviser: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	must, never, ctx := sectionOf(p.Body, "## Must do"), sectionOf(p.Body, "## Must not do"), sectionOf(p.Body, "## Context")
	for blk, sec := range map[string]string{"command-budget": must, "never-list": never} {
		if !strings.Contains(sec, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("plan-reviser: the %s block is not expanded inside its section", blk)
		}
	}
	if !strings.Contains(ctx, strings.TrimRight(shared["session-done.md"], "\n")) {
		t.Errorf("plan-reviser: the session-done block is not expanded inside Context")
	}
	for _, w := range []string{"`rhiza apply <file> --change <cs> --json`", "`rhiza apply <file> --record --json`", "split the leaf into leaves of one concern each", "add a leaf serving the intent clause", "Never invent an owner's words", "`quote`"} {
		if !strings.Contains(must, w) {
			t.Errorf("plan-reviser Must do: no %q", w)
		}
	}
	if !strings.Contains(sectionOf(p.Body, "## Read first"), "`confirmed` and major") {
		t.Errorf("plan-reviser Read first: no confirmed-major-only rule")
	}
	if !strings.Contains(never, "not confirmed and major") {
		t.Errorf("plan-reviser Must not do: no confirmed-major-only rule")
	}
	if strings.Contains(p.Body, "<!-- shared:") {
		t.Errorf("plan-reviser: a shared marker is left unexpanded")
	}
	// Negative: no findings or verdict heredoc.
	for _, banned := range []string{"rhiza run findings", "rhiza run verdict", "<<'EOF'", `"findings":`, `"verdicts":`} {
		if strings.Contains(p.Body, banned) {
			t.Errorf("plan-reviser carries %q", banned)
		}
	}
	if n := Tokens(p.Body); n > TopicCap {
		t.Errorf("plan-reviser costs %d tokens, cap %d", n, TopicCap)
	}
}

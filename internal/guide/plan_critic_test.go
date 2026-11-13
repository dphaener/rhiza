package guide

import (
	"strings"
	"testing"
)

func TestPlanCriticText(t *testing.T) {
	p, ok := Lookup("plan-critic")
	if !ok || p.Skill {
		t.Fatalf("plan-critic: found %v, skill %v; want a topic", ok, p.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(p.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("plan-critic: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	must, never, ctx := sectionOf(p.Body, "## Must do"), sectionOf(p.Body, "## Must not do"), sectionOf(p.Body, "## Context")
	for blk, sec := range map[string]string{"command-budget": must, "never-list": never, "severities": ctx} {
		if !strings.Contains(sec, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("plan-critic: the %s block is not expanded inside its section", blk)
		}
	}
	if !strings.Contains(ctx, strings.TrimRight(shared["session-done.md"], "\n")) {
		t.Errorf("plan-critic: the session-done block is not expanded inside Context")
	}
	if !strings.Contains(sectionOf(p.Body, "## Read first"), "`rhiza run status --json`") {
		t.Errorf("plan-critic Read first: no rhiza run status --json")
	}
	for _, w := range []string{"rhiza run verdict - --json <<'EOF'", "`confirmed`", "`downgraded`", "`dropped`", `"severity":"minor"`, `{"verdicts":[`} {
		if !strings.Contains(must, w) {
			t.Errorf("plan-critic Must do: no %q", w)
		}
	}
	for _, w := range []string{"Add a finding", "Raise a severity or rewrite a finding"} {
		if !strings.Contains(never, w) {
			t.Errorf("plan-critic Must not do: no %q", w)
		}
	}
	if strings.Contains(p.Body, "<!-- shared:") {
		t.Errorf("plan-critic: a shared marker is left unexpanded")
	}
	// Negative: no findings heredoc or apply instructions.
	for _, banned := range []string{"rhiza run findings", "rhiza apply", `"findings":`} {
		if strings.Contains(p.Body, banned) {
			t.Errorf("plan-critic carries %q", banned)
		}
	}
	if n := Tokens(p.Body); n > TopicCap {
		t.Errorf("plan-critic costs %d tokens, cap %d", n, TopicCap)
	}
}

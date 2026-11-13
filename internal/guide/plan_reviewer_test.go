package guide

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

func TestPlanReviewerText(t *testing.T) {
	p, ok := Lookup("plan-reviewer")
	if !ok || p.Skill {
		t.Fatalf("plan-reviewer: found %v, skill %v; want a topic", ok, p.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(p.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("plan-reviewer: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	must, never, ctx := sectionOf(p.Body, "## Must do"), sectionOf(p.Body, "## Must not do"), sectionOf(p.Body, "## Context")
	for blk, sec := range map[string]string{"command-budget": must, "never-list": never, "severities": ctx} {
		if !strings.Contains(sec, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("plan-reviewer: the %s block is not expanded inside its section", blk)
		}
	}
	if !strings.Contains(ctx, strings.TrimRight(shared["session-done.md"], "\n")) {
		t.Errorf("plan-reviewer: the session-done block is not expanded inside Context")
	}
	for _, k := range run.Kinds {
		if !strings.Contains(must, "`"+k+"`") {
			t.Errorf("plan-reviewer Must do: no kind %q", k)
		}
	}
	for _, w := range []string{"one that is vague or nothing could", "assumes something about the codebase", "`rhiza overlaps --change <cs> --json`", "its findings are not recorded", "rhiza run findings - --json <<'EOF'", `{"findings":[]}`, "an empty round"} {
		if !strings.Contains(must, w) {
			t.Errorf("plan-reviewer Must do: no %q", w)
		}
	}
	if strings.Contains(p.Body, "<!-- shared:") {
		t.Errorf("plan-reviewer: a shared marker is left unexpanded")
	}
	for _, banned := range []string{"rhiza run verdict", "rhiza apply", `"verdicts":`} {
		if strings.Contains(p.Body, banned) {
			t.Errorf("plan-reviewer carries %q", banned)
		}
	}
	if n := Tokens(p.Body); n > TopicCap {
		t.Errorf("plan-reviewer costs %d tokens, cap %d", n, TopicCap)
	}
}

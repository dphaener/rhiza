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
	if len(run.Kinds) != 8 || !strings.Contains(must, "of one of eight kinds") {
		t.Errorf("plan-reviewer Must do: want eight kinds, run.Kinds has %d", len(run.Kinds))
	}
	for _, k := range run.Kinds {
		if !strings.Contains(must, "`"+k+"`") {
			t.Errorf("plan-reviewer Must do: no kind %q", k)
		}
	}
	if o, n := strings.Index(must, "`oversized-leaf`"), strings.Index(must, "`no-failure-case`"); o < 0 || n < o {
		t.Errorf("plan-reviewer Must do: no-failure-case is not listed after oversized-leaf")
	}
	// missing-acceptance is narrowed: a missing acceptance is review-check's.
	if strings.Contains(must, "a leaf with no `acceptance`") {
		t.Errorf("plan-reviewer Must do: missing-acceptance still covers a leaf with no acceptance")
	}
	if !strings.Contains(never, "what `rhiza run review-check` already checks") {
		t.Errorf("plan-reviewer Must not do: no rhiza run review-check")
	}
	for _, w := range []string{"an acceptance that is vague or nothing could run", "runs\n     the repo's full suite, not the leaf's tests", "not a real failure path of the", "a `failure_waiver` with no real reason", "Never record what `rhiza run review-check` already checks", "`must_not`", "no `failure case` words", "duplicate handles", "assumes something about the codebase", "`rhiza overlaps --change <cs> --json`", "its findings\n   are not recorded", "rhiza run findings - --json <<'EOF'", `{"findings":[]}`, "an empty round"} {
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

// Failure case: a topic or description still counting the plan-review
// kinds as seven is stale now that run.Kinds has no-failure-case.
func TestNoTextSaysSevenKinds(t *testing.T) {
	for _, p := range List() {
		for where, s := range map[string]string{"description": p.Description, "body": p.Body} {
			if strings.Contains(s, "seven kinds") || strings.Contains(s, "the seven,") {
				t.Errorf("%s %s still says seven kinds", p.Name, where)
			}
		}
	}
}

package guide

import (
	"strings"
	"testing"
)

func TestPlannerText(t *testing.T) {
	p, ok := Lookup("planner")
	if !ok || p.Skill {
		t.Fatalf("planner: found %v, skill %v; want a topic", ok, p.Skill)
	}
	want := []string{"## Task", "## Expected outcome", "## Read first", "## Must do", "## Must not do", "## Context"}
	var got []string
	for _, l := range strings.Split(p.Body, "\n") {
		if strings.HasPrefix(l, "## ") {
			got = append(got, l)
		}
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("planner: sections %q, want exactly %q in order", got, want)
	}
	shared := sharedFiles(t)
	must, never, ctx := sectionOf(p.Body, "## Must do"), sectionOf(p.Body, "## Must not do"), sectionOf(p.Body, "## Context")
	for blk, sec := range map[string]string{"command-budget": must, "outside-work": must, "never-list": never} {
		if !strings.Contains(sec, demote(strings.TrimRight(shared[blk+".md"], "\n"))) {
			t.Errorf("planner: the %s block is not expanded inside its section", blk)
		}
	}
	if !strings.Contains(ctx, strings.TrimRight(shared["session-done.md"], "\n")) {
		t.Errorf("planner: the session-done block is not expanded inside Context")
	}

	read := sectionOf(p.Body, "## Read first")
	for _, w := range []string{"`rhiza guide plan-recipes`", "`rhiza guide plan-interview`"} {
		if !strings.Contains(read, w) {
			t.Errorf("planner Read first: no %s", w)
		}
	}
	// The research pass sits in Read first, before the interview step.
	research := strings.Index(read, "research pass")
	if research < 0 {
		t.Fatalf("planner Read first: no research pass")
	}
	pass := read[research:]
	for _, w := range []string{"before the interview", "rhiza why", "search the repo for the names the intent uses", "code and tests", "the files each leaf touches"} {
		if !strings.Contains(pass, w) {
			t.Errorf("planner research pass: no %q", w)
		}
	}
	if strings.Index(p.Body, "research pass") > strings.Index(p.Body, "Interview the owner") {
		t.Errorf("planner: the research pass comes after the interview step")
	}

	i := strings.Index(must, "Self-check")
	if i < 0 {
		t.Fatalf("planner Must do: no self-check")
	}
	check := must[i:]
	for _, w := range []string{"`rhiza slice <change> --json`", "every intent clause is\n   served by a leaf", "one concern", "provenance", "depends-on", "`rhiza check`"} {
		if !strings.Contains(check, w) {
			t.Errorf("planner self-check: no %q", w)
		}
	}
}

// The three plan topics are plain topics within TopicCap, shared blocks
// expanded, and name no repo-specific command.
func TestPlanTopicsPlain(t *testing.T) {
	for _, name := range []string{"planner", "plan-recipes", "plan-interview"} {
		tx, ok := Lookup(name)
		if !ok || tx.Skill {
			t.Errorf("%s: found %v, skill %v; want a topic", name, ok, tx.Skill)
			continue
		}
		if strings.Contains(tx.Body, "<!-- shared:") {
			t.Errorf("%s: a shared marker is left unexpanded", name)
		}
		if n := Tokens(tx.Body); n > TopicCap {
			t.Errorf("%s costs %d tokens, cap %d", name, n, TopicCap)
		}
		for _, banned := range []string{"go test", "go build", "go vet", "make ", "make\n", "testscript", "Claude"} {
			if strings.Contains(tx.Body+tx.Description, banned) {
				t.Errorf("%s names %q", name, banned)
			}
		}
	}
}

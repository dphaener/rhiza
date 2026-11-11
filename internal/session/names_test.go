package session

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
)

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Readable tabs":       "readable-tabs",
		"  The Code adapter ": "the-code-adapter",
		"A -- B__C!!":         "a-b-c",
		"Carving keeps IDs":   "carving-keeps-ids",
		"v2.0 launch":         "v2-0-launch",
		"":                    "session",
		"!!! ???":             "session",
		"ünï":                 "n",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSlugBranchWorktreeFromOneTitle(t *testing.T) {
	slug := Slug("Readable Tabs")
	if got := Branch(slug); got != "cs/readable-tabs" {
		t.Errorf("branch %q", got)
	}
	if got := Worktree("/code/rhiza", "rhiza", slug); got != "/code/rhiza.readable-tabs" {
		t.Errorf("worktree %q", got)
	}
	if got := Worktree("/code/rhiza.other", "rhiza", slug); got != "/code/rhiza.readable-tabs" {
		t.Errorf("worktree from sibling %q", got)
	}
}

func TestTabTitle(t *testing.T) {
	cases := []struct{ repo, title, step, want string }{
		{"rhiza", "Readable tabs", "build", "rhiza: Readable tabs"},
		{"rhiza", "Readable tabs", "", "rhiza: Readable tabs"},
		{"rhiza", "Readable tabs", "review", "rhiza: Readable tabs · review"},
		{"rhiza", "Carving keeps IDs", "review", "rhiza: Carving keeps IDs · review"},
	}
	for _, c := range cases {
		if got := TabTitle(c.repo, c.title, c.step); got != c.want {
			t.Errorf("TabTitle(%q,%q,%q) = %q, want %q", c.repo, c.title, c.step, got, c.want)
		}
	}
}

func TestRunTabTitle(t *testing.T) {
	cases := []struct{ step, label, want string }{
		{"build", "builder export-csv #2", "rhiza: Readable tabs · run R1 · build · builder export-csv #2"},
		{"review", "check alpha #1", "rhiza: Readable tabs · run R1 · review · check alpha #1"},
		{"interview", "", "rhiza: Readable tabs · run R1 · interview"},
	}
	for _, c := range cases {
		if got := RunTabTitle("rhiza", "Readable tabs", "R1", c.step, c.label); got != c.want {
			t.Errorf("RunTabTitle(%q,%q) = %q, want %q", c.step, c.label, got, c.want)
		}
	}
}

func TestTabTitlesFitRex(t *testing.T) {
	long := strings.Repeat("Inbox and daemon treat a landed plan run ", 3)
	got := RunTabTitle("rhiza", long, "01M6MHWK3K41VADZ6YD6M0B39J", "build/build", "verifier inbox-landed-plan-run-in-flight #1")
	if n := utf8.RuneCountInString(got); n > MaxTabTitle {
		t.Fatalf("RunTabTitle is %d characters, over %d: %q", n, MaxTabTitle, got)
	}
	if want := "… · run 01M6MHWK3K41VADZ6YD6M0B39J · build/build · verifier inbox-landed-plan-run-in-flight #1"; !strings.HasSuffix(got, want) {
		t.Errorf("RunTabTitle = %q, want the title shortened and the tail %q kept", got, want)
	}
	if !strings.HasPrefix(got, "rhiza: Inbox and daemon") {
		t.Errorf("RunTabTitle = %q, want the repo and the title's start kept", got)
	}
	if got := TabTitle("rhiza", long+long, "review"); utf8.RuneCountInString(got) > MaxTabTitle || !strings.HasSuffix(got, "… · review") {
		t.Errorf("TabTitle = %q, want at most %d characters ending in the step", got, MaxTabTitle)
	}
	tail := strings.Repeat("x", 200)
	if got := RunTabTitle("rhiza", "T", "R", "s", tail); utf8.RuneCountInString(got) != MaxTabTitle {
		t.Errorf("RunTabTitle with an overlong tail is %d characters, want %d", utf8.RuneCountInString(got), MaxTabTitle)
	}
}

func TestChangeSetSlug(t *testing.T) {
	mk := func(idStr, title string, state changeset.State) changeset.ChangeSet {
		cs, err := id.Parse(idStr)
		if err != nil {
			t.Fatal(err)
		}
		return changeset.ChangeSet{ID: cs, OpenEntry: cs, State: state, Open: changeset.Open{Title: title}}
	}
	a := mk("01M6T1ZR0000000000000000AA", "Fix the usage governor's", changeset.StateLanded)
	b := mk("01M6T2DC1234567890ABCDEFGH", "Fix the usage governor's", changeset.StateDraft)
	c := mk("01M6T2DC1234567890ABCDEFGJ", "Fix the usage governor's", changeset.StateDraft)
	d := mk("01M6T3000000000000000000AA", "Other work", changeset.StateDraft)
	sets := []changeset.ChangeSet{b, d, a, c}

	if got := ChangeSetSlug(sets, a.ID); got != "fix-the-usage-governor-s" {
		t.Errorf("earlier = %q", got)
	}
	if got, want := ChangeSetSlug(sets, b.ID), "fix-the-usage-governor-s-01m6t2dc1234"; got != want {
		t.Errorf("later = %q, want %q", got, want)
	}
	if got := ChangeSetSlug(sets, d.ID); got != "other-work" {
		t.Errorf("no clash = %q", got)
	}
	// IDs sharing a long prefix still get fixed-length suffixes, stable on repeat.
	if got := ChangeSetSlug(sets, c.ID); got != "fix-the-usage-governor-s-01m6t2dc1234" {
		t.Errorf("c = %q", got)
	}
	// Same after the earlier set is withdrawn, and on repeat.
	a.State = changeset.StateWithdrawn
	for i := 0; i < 2; i++ {
		if got := ChangeSetSlug([]changeset.ChangeSet{a, b, c}, b.ID); got != "fix-the-usage-governor-s-01m6t2dc1234" {
			t.Errorf("after abandon = %q", got)
		}
	}
	// Without the earlier set b is bare.
	if got := ChangeSetSlug([]changeset.ChangeSet{b, c}, b.ID); got != "fix-the-usage-governor-s" {
		t.Errorf("alone = %q", got)
	}
}

func TestTabTitleNestedStep(t *testing.T) {
	const step = "ship-change/build/build"
	if got, want := TabTitle("rhiza", "Fix it", step), "rhiza: Fix it · "+step; got != want {
		t.Errorf("TabTitle = %q, want %q", got, want)
	}
	if got, want := RunTabTitle("rhiza", "Fix it", "r1", step, "builder leaf #1"), "rhiza: Fix it · run r1 · "+step+" · builder leaf #1"; got != want {
		t.Errorf("RunTabTitle = %q, want %q", got, want)
	}
}

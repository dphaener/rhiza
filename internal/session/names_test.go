package session

import "testing"

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

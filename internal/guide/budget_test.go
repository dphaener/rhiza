package guide

import (
	"strings"
	"testing"
)

func TestBuildReviewCommandBudget(t *testing.T) {
	for _, name := range []string{"review"} {
		g, ok := Lookup(name)
		if !ok {
			t.Fatalf("no guide %q", name)
		}
		for _, p := range []string{
			"## Command budget", "Never sleep and never poll", "about four\nminutes",
			"runs once in the background and is\nwaited for once",
		} {
			if !strings.Contains(g.Body, p) {
				t.Errorf("%s: no %q", name, p)
			}
		}
		i := strings.Index(g.Body, "## Command budget")
		if i < 0 {
			continue
		}
		sec := g.Body[i:]
		if j := strings.Index(sec[3:], "\n## "); j >= 0 {
			sec = sec[:j+3]
		}
		for _, banned := range []string{"go ", "make", "rhiza", "testscript", "Claude"} {
			if strings.Contains(sec, banned) {
				t.Errorf("%s: budget rule names %q", name, banned)
			}
		}
	}
}

package guide

import (
	"strings"
	"testing"
)

func TestBuildRetryText(t *testing.T) {
	r, ok := Lookup("build-retry")
	if !ok || r.Skill {
		t.Fatalf("build-retry: found %v, skill %v; want a topic", ok, r.Skill)
	}
	checkPhrases(t, "build-retry", map[string][]string{
		"tree":     {"inspect the working tree", "earlier attempts' commits"},
		"findings": {"`rhiza run status --json`", "Address every finding exactly as written"},
		"commit":   {"one fix commit with the same `Leaf:` trailer", "Never rewrite history"},
	})
	if n := Tokens(r.Body); n > TopicCap/3 {
		t.Errorf("build-retry costs %d tokens, want at most %d (a third of TopicCap)", n, TopicCap/3)
	}
	b, _ := Lookup("builder")
	shared := sharedFiles(t)
	for name, text := range shared {
		text = strings.TrimRight(text, "\n")
		if !strings.Contains(b.Body, text) && !strings.Contains(b.Body, demote(text)) {
			continue // not one of the builder's blocks
		}
		if strings.Contains(r.Body, text) || strings.Contains(r.Body, demote(text)) {
			t.Errorf("build-retry repeats the builder's %s block", name)
		}
	}
}

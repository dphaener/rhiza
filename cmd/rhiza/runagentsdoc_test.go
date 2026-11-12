package main

import (
	"os"
	"strings"
	"testing"
)

func TestRunAgentsDocLessons(t *testing.T) {
	b, err := os.ReadFile("../../agents/run-agents.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, lesson := range []string{
		"share the\none last paragraph",
		"`Co-Authored-By:`",
		"already landed",
		"--allow-empty",
		"Merge `origin/main`",
		"Never rebase and\nforce-push",
		".git/rhiza/ship/<cs>.json",
	} {
		if !strings.Contains(doc, lesson) {
			t.Errorf("agents/run-agents.md lacks %q", lesson)
		}
	}
}

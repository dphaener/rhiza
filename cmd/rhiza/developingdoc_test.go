package main

import (
	"os"
	"strings"
	"testing"
)

func TestDevelopingDocLessons(t *testing.T) {
	b, err := os.ReadFile("../../agents/developing.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, lesson := range []string{
		"## Testscripts",
		"${WORK@R}",
		"-count=N",
		"`chmod` takes one path per call",
		"relative to the current `cd`",
		"Fake binaries",
		"session scratchpad",
		"Unset `RHIZA_ROOT`, `RHIZA_CHANGE` and `RHIZA_RUN`",
		"GOTMPDIR",
		"scratch worktree",
		"root daemon",
		"rhiza render",
	} {
		if !strings.Contains(doc, lesson) {
			t.Errorf("agents/developing.md lacks %q", lesson)
		}
	}
}

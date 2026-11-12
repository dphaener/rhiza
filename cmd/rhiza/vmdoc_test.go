package main

import (
	"os"
	"strings"
	"testing"
)

func TestVMDocLessons(t *testing.T) {
	b, err := os.ReadFile("../../agents/vm.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, lesson := range []string{
		"rhiza-root.service",
		"rhiza-release.timer",
		"rhiza.daemon.maxRuns",
		"The Rex package is pinned",
		"stray-logs",
		"injecting a delay",
	} {
		if !strings.Contains(doc, lesson) {
			t.Errorf("agents/vm.md lacks %q", lesson)
		}
	}
}

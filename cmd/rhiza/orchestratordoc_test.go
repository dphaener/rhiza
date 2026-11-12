package main

import (
	"os"
	"strings"
	"testing"
)

func TestOrchestratorDocRules(t *testing.T) {
	b, err := os.ReadFile("../../agents/orchestrator.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	for _, rule := range []string{
		"Never run `git fetch` in the main clone",
		"Never `pkill` by pattern",
		"Nothing is ever left uncommitted on main",
		"its own branch and PR",
		"rhiza.daemon.maxRuns",
		"merge `origin/main` into its branch",
		"Never rebase and force-push",
		"gh pr merge <n> --squash",
		"Quote the owner exactly",
		"Plan a milestone only with the owner",
	} {
		if !strings.Contains(doc, rule) {
			t.Errorf("agents/orchestrator.md lacks rule %q", rule)
		}
	}
}

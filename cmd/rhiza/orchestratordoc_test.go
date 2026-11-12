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

func TestOrchestratorConsultsOwner(t *testing.T) {
	read := func(p string) string {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	doc := read("../../agents/orchestrator.md")
	for _, rule := range []string{
		"Consult the owner, with a recommendation, before any judgment call",
		"has granted autonomy for that task or",
		"propose the next high-value work to the owner",
		"instead of starting it",
		"Queued bug requests still flow",
	} {
		if !strings.Contains(doc, rule) {
			t.Errorf("agents/orchestrator.md lacks rule %q", rule)
		}
	}
	if strings.Contains(read("../../AGENTS.md"), "Consult the owner") {
		t.Error("AGENTS.md carries the consult rule")
	}
}

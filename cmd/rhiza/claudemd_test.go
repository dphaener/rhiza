package main

import (
	"os"
	"strings"
	"testing"
)

// AGENTS.md tells agents in runs to test only what their change touches.
func TestClaudeMDScopedTests(t *testing.T) {
	b, err := os.ReadFile("../../AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := strings.Join(strings.Fields(string(b)), " ")
	for _, want := range []string{
		"go test -trimpath",
		"TestScripts/<name>",
		"make self-check",
		"never run `make check`, `make test` or a bare `go test ./...`",
		"While CI is unstable",
		"not explained by your change",
		"do not push or rerun CI to see",
		"Reproduce the failing test locally",
		"push only after it and the touched packages pass locally",
		"inject a delay at the suspected race",
		"Never use busy loops, CPU stress or parallel stress runs on the host",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("AGENTS.md lacks %q", want)
		}
	}
	if strings.Contains(doc, "Run `make check` before calling any change done.") {
		t.Error("AGENTS.md still has the unqualified make check rule")
	}
}

// CLAUDE.md is only the import of AGENTS.md, which holds the repo guidance.
func TestClaudeMDImportsAgentsMD(t *testing.T) {
	b, err := os.ReadFile("../../CLAUDE.md")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(b)); got != "@AGENTS.md" {
		t.Errorf("CLAUDE.md = %q, want only @AGENTS.md", got)
	}
	a, err := os.ReadFile("../../AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{
		"## Ignore the `~/code/.claude/rules` files",
		"## Stack and shape",
		"## Bug requests",
		"## Commands",
		"## Tests",
		"## Binaries",
	} {
		if !strings.Contains(string(a), "\n"+h+"\n") {
			t.Errorf("AGENTS.md lacks section %q", h)
		}
	}
}

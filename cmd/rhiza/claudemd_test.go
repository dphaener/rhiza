package main

import (
	"os"
	"strings"
	"testing"
)

// CLAUDE.md tells agents in runs to test only what their change touches.
func TestClaudeMDScopedTests(t *testing.T) {
	b, err := os.ReadFile("../../CLAUDE.md")
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
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("CLAUDE.md lacks %q", want)
		}
	}
	if strings.Contains(doc, "Run `make check` before calling any change done.") {
		t.Error("CLAUDE.md still has the unqualified make check rule")
	}
}

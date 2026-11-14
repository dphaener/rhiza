package main

import (
	"os"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// docs/primitives.md has a section for each of the eight primitives and
// names every step kind the workflow package has.
func TestDocsPrimitives(t *testing.T) {
	b, err := os.ReadFile("../../docs/primitives.md")
	if err != nil {
		t.Fatal(err)
	}
	doc := string(b)
	headings := map[string]bool{}
	for line := range strings.SplitSeq(doc, "\n") {
		if h, ok := strings.CutPrefix(line, "## "); ok {
			headings[strings.TrimSpace(h)] = true
		}
	}
	for _, p := range []string{"Agent", "Spawn", "Gate", "Push/PR", "Wait", "Merge", "Sign-off", "Hook"} {
		if !headings[p] {
			t.Errorf("docs/primitives.md has no ## %s section", p)
		}
	}
	for _, k := range workflow.Kinds {
		if !strings.Contains(doc, "`"+k+"`") {
			t.Errorf("docs/primitives.md does not name step kind `%s`", k)
		}
	}
	for _, s := range []string{"handback:", "handback_cap", "rhiza run spawn", "rhiza ship open-pr", "rhiza ship merge", "[[cleanup]]", "](reference/workflow-file.md)"} {
		if !strings.Contains(doc, s) {
			t.Errorf("docs/primitives.md does not mention %s", s)
		}
	}
	for _, f := range []string{"../../README.md", "../../docs/runs-and-daemon.md", "../../docs/custom-workflows.md"} {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), "primitives.md)") {
			t.Errorf("%s does not link docs/primitives.md", strings.TrimPrefix(f, "../../"))
		}
	}
}

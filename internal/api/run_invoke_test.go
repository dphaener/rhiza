package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A session's first prompt is its role's invoke rendered with the session's
// file paths and the run's values; the lines the daemon adds for the session
// follow it unchanged, and with none the prompt is the rendered invoke.
func TestInvokePrompt(t *testing.T) {
	vars := workflow.PromptVars{Change: "CS1", Run: "R1", Brief: "/s/brief.md", Slice: "/s/slice.json", Handback: "/s/handback.md", Worktree: "/wt"}
	const invoke = "/myskill:build --prompt {brief} {leaf} {attempt}"
	role := workflow.Role{Prompt: invoke}

	builder := run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 1}
	if got, want := leafPrompt(role, builder, vars), "/myskill:build --prompt /s/brief.md leaf-a 1"; got != want {
		t.Errorf("attempt 1 = %q, want %q", got, want)
	}
	builder.Attempt = 2
	got := leafPrompt(role, builder, vars)
	if !strings.HasPrefix(got, "/myskill:build --prompt /s/brief.md leaf-a 2\n\nThis is attempt 2 of 3 on leaf leaf-a") {
		t.Errorf("attempt 2 = %q", got)
	}

	fixer := workflow.Role{Prompt: "/myskill:fix {slice} {worktree} {leaf} {cycle}"}
	fs := []run.Finding{{ID: "f1", Check: "alpha", Severity: "major", Reason: "bad"}}
	answers := []run.ReviewAnswer{{Question: "Q1", Decision: "rename-it", Statement: "Rename it."}}
	leaf := fixerPrompt(fixer, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "leaf-a"}, Cycle: 2, Findings: fs, Answers: answers}, vars)
	cs := changeSetFixerPrompt(fixer, run.FixSpec{Cycle: 2, Findings: fs, Answers: answers}, vars)
	for name, c := range map[string]struct{ got, head string }{
		"leaf fixer":       {leaf, "/myskill:fix /s/slice.json /wt leaf-a 2\n\nConfirmed findings to fix in review cycle 2 on leaf leaf-a:"},
		"change set fixer": {cs, "/myskill:fix /s/slice.json /wt  2\n\nConfirmed findings to fix in review cycle 2, which name no leaf"},
	} {
		if !strings.HasPrefix(c.got, c.head) {
			t.Errorf("%s = %q, want it to start %q", name, c.got, c.head)
		}
		if !strings.Contains(c.got, "- Decision rename-it answers question Q1: Rename it.") {
			t.Errorf("%s lacks its answer line: %q", name, c.got)
		}
	}
}

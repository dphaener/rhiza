package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A sign-off's line names the owner by actor, an autonomous run's by its
// rule, and says changes were asked for; a run waits at its sign-off only
// while stopped there, unfinished.
func TestSignoffLine(t *testing.T) {
	for _, tt := range []struct {
		s    RunSignoff
		want string
	}{
		{RunSignoff{Step: "signoff", Visit: 1, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffApproved}, "signoff (visit 1): approved by ada"},
		{RunSignoff{Step: "signoff", Visit: 2, Actor: "ada", Rule: run.RuleOwner, Outcome: run.SignoffChanges}, "signoff (visit 2): changes asked by ada"},
		{RunSignoff{Step: "signoff", Visit: 1, Rule: run.RuleAutonomous, Outcome: run.SignoffApproved}, "signoff (visit 1): approved by rule autonomous"},
	} {
		if got := tt.s.line(); got != tt.want {
			t.Errorf("line = %q, want %q", got, tt.want)
		}
	}
	if !signoffWaiting(run.State{Stopped: run.SignoffWait}) {
		t.Error("a run stopped at its sign-off is not waiting")
	}
	if signoffWaiting(run.State{Stopped: run.SignoffWait, Ended: true}) || signoffWaiting(run.State{Stopped: "waiting: false"}) {
		t.Error("a finished run, or one stopped elsewhere, is waiting")
	}
}

// In an owner cycle the fixers' prompts say the owner asked for the
// changes, not that a review cycle confirmed them.
func TestOwnerCycleFixerPrompts(t *testing.T) {
	role := workflow.Role{Prompt: "fix {leaf}"}
	vars := workflow.PromptVars{Run: "R1"}
	fs := []run.Finding{{ID: "f1", Check: "alpha", Severity: "major", Place: "a.go:1", Reason: "say more"}}
	leaf := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "thing"}, Cycle: 4, Owner: true, Findings: fs}, vars)
	cs := changeSetFixerPrompt(role, run.FixSpec{Cycle: 4, Owner: true, Findings: fs}, vars)
	for _, tt := range []struct{ p, want string }{
		{leaf, "The owner requested these changes on leaf thing:\n\n- alpha f1 (major at a.go:1): say more\n"},
		{cs, "The owner requested these changes, which name no leaf of the change set:\n\n- alpha f1 (major at a.go:1): say more\n"},
	} {
		if !strings.Contains(tt.p, tt.want) || strings.Contains(tt.p, "Review cycle") {
			t.Errorf("prompt lacks %q or names a review cycle:\n%s", tt.want, tt.p)
		}
	}
	plain := fixerPrompt(role, run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "thing"}, Cycle: 2, Findings: fs}, vars)
	if !strings.Contains(plain, "Review cycle 2 confirmed these findings on leaf thing:") {
		t.Errorf("a fixer outside an owner cycle is not told its review cycle:\n%s", plain)
	}
}

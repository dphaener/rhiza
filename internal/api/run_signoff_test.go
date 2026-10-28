package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/run"
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

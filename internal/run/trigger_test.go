package run

import (
	"strings"
	"testing"
)

// A run with no change set (a scheduled run) fails at a step that needs one,
// naming the step, and never reaches the runner.
func TestRunWithoutChangeSetFailsLeavesStep(t *testing.T) {
	e, r, _, def := setup(t, leavesWF)
	if err := e.Store.Create("SCHED", Entry{Type: TypeStarted, At: t0, Step: def.Start, Started: &Started{
		Workflow: "wf", Definition: leavesWF, Trigger: TriggerSchedule, TriggerName: "intake",
	}}); err != nil {
		t.Fatal(err)
	}
	st, err := e.Run("SCHED")
	if err != nil {
		t.Fatal(err)
	}
	if st.Trigger != TriggerSchedule || st.TriggerName != "intake" || st.Change != "" {
		t.Errorf("state trigger = %q %q, change %q", st.Trigger, st.TriggerName, st.Change)
	}
	if st.Result != Failed || !strings.Contains(st.Detail, `step "build"`) || !strings.Contains(st.Detail, "needs a change set") {
		t.Errorf("result %q detail %q, want a failure naming step build", st.Result, st.Detail)
	}
	if len(r.leafSpecs) != 0 {
		t.Errorf("a leaf was started: %v", r.leafSpecs)
	}
}

package guide

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func TestWorkflowsText(t *testing.T) {
	wf, ok := Lookup("workflows")
	if !ok || wf.Skill {
		t.Fatalf("workflows: found %v, skill %v; want a topic", ok, wf.Skill)
	}
	// Every step kind the workflow package defines is named as a list item.
	for _, k := range []string{
		workflow.KindAgent, workflow.KindWait, workflow.KindCommand, workflow.KindArbiter,
		workflow.KindLeaves, workflow.KindChecks, workflow.KindSignoff, workflow.KindWorkflow,
	} {
		if !strings.Contains(wf.Body, "\n- `"+k+"`") {
			t.Errorf("workflows: names no step kind %q", k)
		}
	}
	checkPhrases(t, "workflows", map[string][]string{
		"layers":   {"$XDG_CONFIG_HOME/rhiza/workflows/<name>.toml", ".rhiza/workflows/<name>.toml", "built-in"},
		"targets":  {"`done` and `failed`", "max_visits"},
		"includes": {"skip = [steps]", "<workflow>/<step>"},
		"loop": {
			"rhiza workflow ls", "rhiza workflow new mine --from ship-change",
			"rhiza workflow describe mine", "rhiza workflow check mine",
			"rhiza run start --workflow mine --dry-run",
		},
		"outside input": {
			"## Outside input", "bypassPermissions", "rhiza.autonomy.plan", "rhiza request",
			"Rhiza cannot see", "permission_mode", "rhiza.autonomy.signoff",
			"does not gate intake", "rhiza.daemon.permissionMode", "outside Rhiza",
			"accept the plan", "plan workflow's steps",
		},
	})
	if strings.Contains(wf.Body+wf.Description, "Claude") {
		t.Errorf("workflows names an agent product")
	}
}

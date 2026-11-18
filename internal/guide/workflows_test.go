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
		workflow.KindGate,
	} {
		if !strings.Contains(wf.Body, "\n- `"+k+"`") {
			t.Errorf("workflows: names no step kind %q", k)
		}
	}
	if n := Tokens(wf.Body); n > TopicCap {
		t.Errorf("workflows costs %d tokens, cap %d", n, TopicCap)
	}
	checkPhrases(t, "workflows", map[string][]string{
		"workflow keys": {"`rhiza guide workflow-keys`"},
		"layers":        {"$XDG_CONFIG_HOME/rhiza/workflows/<name>.toml", ".rhiza/workflows/<name>.toml", "built-in"},
		"targets":       {"`done` and `failed`", "max_visits"},
		"includes":      {"skip = [steps]", "<workflow>/<step>"},
		"dispatch": {
			".rhiza/dispatch.toml", "[request]", "[ready]", "rhiza request --workflow",
			"feature", "bug", "milestone", "ship-change", "workflow-unresolved",
			"ship-change-classic", "conductor default", "handing back to `conduct`",
		},
		"loop": {
			"rhiza workflow ls", "rhiza workflow new mine --from ship-change",
			"rhiza workflow describe mine", "rhiza workflow check mine",
			"rhiza run start --workflow mine --dry-run",
		},
		"schedules": {
			"## Schedules", "[[schedules]]", "CRON_TZ", "catch_up", "overlap", "every", "cron",
			"main checkout", "no change set", "as it found it",
			"rhiza request", "rhiza health", "daemon log",
		},
		"outside input": {
			"## Outside input", "bypassPermissions", "rhiza.autonomy.plan", "rhiza request",
			"Rhiza cannot see", "permission_mode", "rhiza.autonomy.signoff",
			"does not gate intake", "rhiza.daemon.permissionMode", "outside Rhiza",
			"accept the plan", "plan workflow's steps", "Nothing is refused",
			"narrow `rhiza.daemon.permissionMode`",
		},
	})
	// ship-change's step list names every built-in step after conduct, in
	// file order.
	order := shipChange(t).Order
	if len(order) != 12 || order[0] != "conduct" {
		t.Fatalf("ship-change steps %q, want 12 starting at conduct", order)
	}
	var quoted []string
	for _, s := range order[1:] {
		quoted = append(quoted, "`"+s+"`")
	}
	steps := strings.Join(quoted[:len(quoted)-1], ", ") + " and " + quoted[len(quoted)-1] + " steps"
	if flat := strings.Join(strings.Fields(wf.Body), " "); !strings.Contains(flat, "one `conduct` agent step, then the "+steps) {
		t.Errorf("workflows: ship-change's step list is not %q", steps)
	}
	if !strings.Contains(wf.Body, "a\nfailed merge goes to `rebase`") {
		t.Errorf("workflows: does not say a failed merge goes to rebase")
	}
	if strings.Contains(wf.Body+wf.Description, "Claude") {
		t.Errorf("workflows names an agent product")
	}
}

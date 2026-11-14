package workflow

import (
	"strings"
	"testing"
)

const handbackInner = `
[step.conduct]
kind = "agent"
prompt = "conduct"
on.pass = "gate"
[step.gate]
kind = "command"
run = "gate"
on.pass = "done"
on.fail = "handback:conduct"
`

// A handback inside an included workflow keeps pointing at the included
// step, and done and failed targets map as for any included step.
func TestIncludeHandbackTarget(t *testing.T) {
	files := map[string]string{"inner": handbackInner}
	src := "[step.w]\nkind = \"workflow\"\nworkflow = \"inner\"\non.pass = \"after\"\non.fail = \"cleanup\"\n" +
		"[step.after]\nkind = \"command\"\nrun = \"true\"\n" +
		"[step.cleanup]\nkind = \"command\"\nrun = \"true\"\n"
	d, err := ParseWith("test", []byte(src), mapResolver(files))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.Steps["inner/conduct"]; !ok {
		t.Fatalf("flat steps lack inner/conduct: %v", d.Order)
	}
	if got := d.Steps["inner/gate"].On[OutcomeFail]; got != "handback:inner/conduct" {
		t.Errorf("inner/gate on.fail = %q, want handback:inner/conduct", got)
	}
	if next, _ := d.Next("inner/gate", OutcomePass); next != "after" {
		t.Errorf("inner/gate on pass goes to %q, want after", next)
	}

	// done and failed still map to the step's on.pass and on.fail.
	failing := strings.Replace(handbackInner, `on.fail = "handback:conduct"`, `on.fail = "failed"`, 1)
	d, err = ParseWith("test", []byte(src), mapResolver(map[string]string{"inner": failing}))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Steps["inner/gate"].On[OutcomeFail]; got != "cleanup" {
		t.Errorf("failed maps to %q, want cleanup", got)
	}

	bad := strings.Replace(handbackInner, "handback:conduct", "handback:nosuch", 1)
	_, err = ParseWith("test", []byte(src), mapResolver(map[string]string{"inner": bad}))
	if err == nil || !strings.Contains(err.Error(), "nosuch") || !strings.Contains(err.Error(), `"inner"`) {
		t.Errorf("handback:nosuch: error %v, want one naming nosuch and inner", err)
	}
}

const signoffHandbackInner = `
[step.conduct]
kind = "agent"
prompt = "conduct"
on.pass = "sign"
[step.sign]
kind = "signoff"
on.pass = "done"
on.changes = "handback:conduct"
`

// A signoff may hand back to an agent step in a workflow with includes:
// an included workflow's signoff keeps its hand-back through the include
// (handback:inner/conduct), and so does an including workflow's own. A
// signoff hand-back naming a missing step, an included step or a step that
// is no agent is refused.
func TestIncludeSignoffHandback(t *testing.T) {
	files := map[string]string{"inner": signoffHandbackInner}
	d, err := ParseWith("test", []byte("[step.w]\nkind = \"workflow\"\nworkflow = \"inner\"\non.pass = \"done\"\n"), mapResolver(files))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Steps["inner/sign"].On[OutcomeChanges]; got != "handback:inner/conduct" {
		t.Errorf("inner/sign on.changes = %q, want handback:inner/conduct", got)
	}

	outer := func(changes string) string {
		return `
[step.w]
kind = "workflow"
workflow = "inner"
on.pass = "fix"
[step.fix]
kind = "agent"
prompt = "fix"
on.pass = "gate"
[step.cmd]
kind = "command"
run = "true"
[step.gate]
kind = "signoff"
on.pass = "done"
on.changes = "` + changes + `"
`
	}
	d, err = ParseWith("test", []byte(outer("handback:fix")), mapResolver(files))
	if err != nil {
		t.Fatalf("an including workflow's signoff handing back to its agent step: %v", err)
	}
	if got := d.Steps["gate"].On[OutcomeChanges]; got != "handback:fix" {
		t.Errorf("gate on.changes = %q, want handback:fix", got)
	}
	for _, tt := range []struct{ changes, want string }{
		{"handback:nosuch", `hands back to "nosuch", which is not a step`},
		{"handback:w/conduct", `hands back to "w/conduct", which is not a step`},
		{"handback:cmd", `hands back to "cmd", which is a command step`},
	} {
		_, err := ParseWith("test", []byte(outer(tt.changes)), mapResolver(files))
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("on.changes %q: error %v, want one containing %q", tt.changes, err, tt.want)
		}
	}
}

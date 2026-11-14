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

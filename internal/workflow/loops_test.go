package workflow

import (
	"reflect"
	"strings"
	"testing"
)

func TestPlainRedLoops(t *testing.T) {
	const head = `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "test"

[step.test]
kind = "command"
run = "go test ./..."
max_visits = 3
on.pass = "%s"
`
	cases := []struct {
		name, src string
		want      []PlainRedLoop
	}{
		{"fail back to the agent step before it", head + "on.fail = \"build\"\n",
			[]PlainRedLoop{{Step: "test", Outcome: OutcomeFail, Target: "build", Agent: "build"}}},
		{"handback to the agent step", head + "on.fail = \"handback:build\"\n", nil},
		{"red route to a later step", head + "on.fail = \"report\"\n" + `
[step.report]
kind = "command"
run = "echo failed"
on.pass = "failed"
`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := []byte(strings.Replace(c.src, "%s", "done", 1))
			d, err := Parse("w", src)
			if err != nil {
				t.Fatal(err)
			}
			if got := d.PlainRedLoops(); !reflect.DeepEqual(got, c.want) {
				t.Errorf("PlainRedLoops() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestPlainRedLoopsNeedAnAgentStep(t *testing.T) {
	// fail goes back to a command step whose way back to test passes no
	// agent step: a command-only retry is not a plain red loop.
	src := `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "lint"

[step.lint]
kind = "command"
run = "lint"
on.pass = "test"

[step.test]
kind = "command"
run = "go test ./..."
on.pass = "done"
on.fail = "lint"
`
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.PlainRedLoops(); got != nil {
		t.Errorf("no agent on the way back: PlainRedLoops() = %+v, want none", got)
	}
}

func TestPlainRedLoopsThroughHandbackAreCapped(t *testing.T) {
	// fail goes to a command step whose only way back to the agent step is
	// a handback:build route: handback_cap counts that loop.
	src := `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "test"

[step.test]
kind = "command"
run = "go test ./..."
on.pass = "done"
on.fail = "triage"

[step.triage]
kind = "command"
run = "triage"
on.pass = "done"
on.fail = "handback:build"
`
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.PlainRedLoops(); got != nil {
		t.Errorf("way back through a handback route: PlainRedLoops() = %+v, want none", got)
	}
}

func TestPlainRedLoopsShipChangeHasNone(t *testing.T) {
	src, ok := Builtin("ship-change")
	if !ok {
		t.Fatal("no built-in ship-change")
	}
	d, err := Parse("ship-change", src)
	if err != nil {
		t.Fatal(err)
	}
	if got := d.PlainRedLoops(); got != nil {
		t.Errorf("ship-change: PlainRedLoops() = %+v, want none", got)
	}
}

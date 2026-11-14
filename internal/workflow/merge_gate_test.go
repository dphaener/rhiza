package workflow

import (
	"reflect"
	"testing"
)

func TestMergeWithoutGate(t *testing.T) {
	cases := []struct {
		name, src string
		want      []UngatedMerge
	}{
		{"merge reached only through a gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"

[step.gate]
kind = "gate"
commands = ["rhiza check"]
on.pass = "merge"
on.fail = "handback:build"

[step.merge]
kind = "command"
run = "rhiza ship check --no-rebase && rhiza ship merge"
on.pass = "done"
on.fail = "handback:build"
`, nil},
		{"merge reachable around the gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"
on.fail = "push"

[step.gate]
kind = "gate"
commands = ["rhiza check"]
on.pass = "merge"

[step.push]
kind = "command"
run = "rhiza ship open-pr"
on.pass = "merge"

[step.merge]
kind = "command"
run = "gh pr merge --squash"
on.pass = "done"
`, []UngatedMerge{{Step: "merge", Path: []string{"build", "push", "merge"}}}},
		{"no merge", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "test"

[step.test]
kind = "command"
run = "go test ./... && echo rhiza ship merged"
on.pass = "done"
`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := Parse("w", []byte(c.src))
			if err != nil {
				t.Fatal(err)
			}
			if got := d.UngatedMerges(); !reflect.DeepEqual(got, c.want) {
				t.Errorf("UngatedMerges() = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestMergeWithoutGateShipChangeHasNone(t *testing.T) {
	src, ok := Builtin("ship-change")
	if !ok {
		t.Fatal("no built-in ship-change")
	}
	d, err := Parse("ship-change", src)
	if err != nil {
		t.Fatal(err)
	}
	if !IsMergeStep(d.Steps["merge"]) {
		t.Fatal("ship-change's merge step is not a merge step")
	}
	if got := d.UngatedMerges(); got != nil {
		t.Errorf("ship-change: UngatedMerges() = %+v, want none", got)
	}
}

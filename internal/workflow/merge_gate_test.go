package workflow

import (
	"fmt"
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
		{"merge reached only through a command gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"

[step.gate]
kind = "command"
run = "rhiza ship gate"
on.pass = "merge"
on.fail = "handback:build"

[step.merge]
kind = "command"
run = "rhiza ship merge"
on.pass = "done"
`, nil},
		{"a --failed gate is no gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"

[step.gate]
kind = "command"
run = "rhiza run gate --failed"
on.pass = "merge"

[step.merge]
kind = "command"
run = "rhiza ship merge"
on.pass = "done"
`, []UngatedMerge{{Step: "merge", Path: []string{"build", "gate", "merge"}}}},
		{"an echoed gate is no gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"

[step.gate]
kind = "command"
run = "echo rhiza ship gate"
on.pass = "merge"

[step.merge]
kind = "command"
run = "rhiza ship merge"
on.pass = "done"
`, []UngatedMerge{{Step: "merge", Path: []string{"build", "gate", "merge"}}}},
		{"a gate after another command is a gate", `name = "w"
[step.build]
kind = "agent"
prompt = "build {change}"
on.pass = "gate"

[step.gate]
kind = "command"
run = "go vet ./... && rhiza run gate"
on.pass = "merge"

[step.merge]
kind = "command"
run = "rhiza ship merge"
on.pass = "done"
`, nil},
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

func TestMergeWithoutGateBuiltinsHaveNone(t *testing.T) {
	resolve := func(name string) ([]byte, error) {
		src, ok := Builtin(name)
		if !ok {
			return nil, fmt.Errorf("no built-in workflow %q", name)
		}
		return src, nil
	}
	for _, name := range BuiltinNames() {
		src, _ := Builtin(name)
		d, err := ParseWith(name, src, resolve)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := d.UngatedMerges(); got != nil {
			t.Errorf("%s: UngatedMerges() = %+v, want none", name, got)
		}
	}
}

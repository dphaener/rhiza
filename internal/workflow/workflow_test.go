package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestLoadValid(t *testing.T) {
	d, err := Load(filepath.Join("testdata", "valid.toml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if d.Name != "example" {
		t.Errorf("Name = %q, want example", d.Name)
	}
	if d.Start != "build" {
		t.Errorf("Start = %q, want build", d.Start)
	}
	if want := []string{"build", "test", "ci"}; !reflect.DeepEqual(d.Order, want) {
		t.Errorf("Order = %v, want %v", d.Order, want)
	}
	want := map[string]Step{
		"build": {
			Name: "build", Kind: KindAgent, Agent: "claude-code", PermissionMode: "manual",
			Prompt:   "Build change set {change} ({title}) in step {step} of run {run}",
			Interval: DefaultInterval, On: map[string]string{"pass": "test"},
		},
		"test": {
			Name: "test", Kind: KindCommand, Run: "make test", MaxVisits: 3,
			On: map[string]string{"pass": "ci", "fail": "build"},
		},
		"ci": {
			Name: "ci", Kind: KindWait, Until: "gh pr checks",
			Timeout: 30 * time.Minute, Interval: 30 * time.Second,
			On: map[string]string{"pass": "done", "timeout": "failed"},
		},
	}
	if !reflect.DeepEqual(d.Steps, want) {
		t.Errorf("Steps = %+v\nwant %+v", d.Steps, want)
	}
}

func TestLoadRefused(t *testing.T) {
	tests := []struct {
		file string
		want []string
	}{
		{"unknown-kind.toml", []string{`step "repeat": unknown kind "loop": want agent, wait, command, arbiter, leaves, checks, signoff or workflow`}},
		{"missing-target.toml", []string{`step "build": on.pass names "nope", which is not a step; want a step or done or failed`}},
		{"bad-outcome.toml", []string{`step "build": on.sucess: a command step has no outcome "sucess"; want on.pass or on.fail`}},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			path := filepath.Join("testdata", tt.file)
			_, err := Load(path)
			if err == nil {
				t.Fatal("Load succeeded, want an error")
			}
			for _, w := range append(tt.want, path) {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error %q does not contain %q", err, w)
				}
			}
		})
	}
}

func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "absent.toml")); err == nil {
		t.Fatal("Load of a missing file succeeded")
	}
}

func TestParseName(t *testing.T) {
	cmd := "[step.a]\nkind = \"command\"\nrun = \"true\"\n"
	d, err := Parse("fallback", []byte(cmd))
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "fallback" {
		t.Errorf("Name = %q, want fallback", d.Name)
	}
	d, err = Parse("fallback", []byte("name = \"own\"\n"+cmd))
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "own" {
		t.Errorf("Name = %q, want own", d.Name)
	}
	if _, err := Parse("", []byte(cmd)); err == nil || !strings.Contains(err.Error(), "workflow has no name") {
		t.Errorf("Parse with no name: err = %v, want workflow has no name", err)
	}
}

func TestParseOrderAndStart(t *testing.T) {
	src := `
[step.zeta]
kind = "command"
run = "true"
on.pass = "alpha"

[step.alpha]
kind = "command"
run = "true"
on.pass = "mid"

[step.mid]
kind = "wait"
until = "true"
timeout = "1m"
`
	d, err := Parse("order", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"zeta", "alpha", "mid"}; !reflect.DeepEqual(d.Order, want) {
		t.Errorf("Order = %v, want %v", d.Order, want)
	}
	if d.Start != "zeta" {
		t.Errorf("Start = %q, want zeta (the first step in the file)", d.Start)
	}
	if got := d.Steps["mid"].Interval; got != DefaultInterval {
		t.Errorf("wait Interval = %v, want default %v", got, DefaultInterval)
	}
	if got := d.Steps["zeta"].Interval; got != 0 {
		t.Errorf("command Interval = %v, want 0", got)
	}

	d, err = Parse("order", []byte("start = \"mid\"\n"+src))
	if err != nil {
		t.Fatal(err)
	}
	if d.Start != "mid" {
		t.Errorf("Start = %q, want mid", d.Start)
	}
}

func TestParseErrors(t *testing.T) {
	cmd := "[step.a]\nkind = \"command\"\nrun = \"true\"\n"
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"unknown kind", "[step.x]\nkind = \"loop\"\n",
			`step "x": unknown kind "loop": want agent, wait, command, arbiter, leaves, checks, signoff or workflow`},
		{"workflow without workflow", "[step.x]\nkind = \"workflow\"\n",
			`step "x": a workflow step needs workflow =`},
		{"workflow on command", cmd + "workflow = \"b\"\n",
			`step "a": workflow is only for workflow steps`},
		{"skip on command", cmd + "skip = [\"b\"]\n",
			`step "a": skip is only for workflow steps`},
		{"workflow without resolver", "[step.x]\nkind = \"workflow\"\nworkflow = \"b\"\n",
			`step "x": workflow "b" cannot be resolved here`},
		{"no kind", "[step.x]\nrun = \"true\"\n",
			`step "x": no kind: want agent, wait, command, arbiter, leaves, checks, signoff or workflow`},
		{"missing target", cmd + "on.pass = \"nope\"\n",
			`step "a": on.pass names "nope", which is not a step; want a step or done or failed`},
		{"unknown step key", cmd + "retries = 2\n",
			`step "a": unknown key "retries"`},
		{"unknown top key", "steps = 1\n" + cmd,
			`unknown key "steps"`},
		{"unknown table", cmd + "[other]\nx = 1\n",
			`unknown key "other"`},
		{"missing run", "[step.a]\nkind = \"command\"\n",
			`step "a": a command step needs run = "<shell command>"`},
		{"wait without timeout", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\n",
			`step "w": a wait step needs timeout`},
		{"wait without until", "[step.w]\nkind = \"wait\"\ntimeout = \"1m\"\n",
			`step "w": a wait step needs until`},
		{"wait zero timeout", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"0s\"\n",
			`step "w": timeout "0s" is not positive`},
		{"wait bad timeout", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"soon\"\n",
			`step "w": timeout "soon" is not a duration`},
		{"wait negative interval", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\ninterval = \"-1s\"\n",
			`step "w": interval "-1s" is not positive`},
		{"reserved done", "[step.done]\nkind = \"command\"\nrun = \"true\"\n",
			`step "done": "done" is a reserved target`},
		{"reserved failed", "[step.failed]\nkind = \"command\"\nrun = \"true\"\n",
			`step "failed": "failed" is a reserved target`},
		{"bad step name", "[step.Build]\nkind = \"command\"\nrun = \"true\"\n",
			`step "Build": want lowercase letters and digits`},
		{"double hyphen", "[step.a--b]\nkind = \"command\"\nrun = \"true\"\n",
			`step "a--b": want lowercase letters`},
		{"long step name", "[step.a-very-long-step-name-over-thirty]\nkind = \"command\"\nrun = \"true\"\n",
			`step "a-very-long-step-name-over-thirty": name is 33 characters; want at most 30`},
		{"bad outcome name", cmd + "on.Pass = \"done\"\n",
			`step "a": on.Pass: outcome "Pass"`},
		{"run on wait", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\nrun = \"x\"\n",
			`step "w": run is only for command steps`},
		{"until on command", cmd + "until = \"x\"\n",
			`step "a": until is only for wait and agent steps`},
		{"timeout on agent", "[step.g]\nkind = \"agent\"\ntimeout = \"1m\"\n",
			`step "g": timeout is only for wait steps`},
		{"interval on command", cmd + "interval = \"1s\"\n",
			`step "a": interval is only for wait and agent steps`},
		{"agent on command", cmd + "agent = \"x\"\n",
			`step "a": agent is only for agent and arbiter steps`},
		{"permission_mode on wait", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\npermission_mode = \"manual\"\n",
			`step "w": permission_mode is only for agent and arbiter steps`},
		{"negative max_visits", cmd + "max_visits = -1\n",
			`step "a": max_visits is -1; want 0 or more`},
		{"start not a step", "start = \"nope\"\n" + cmd,
			`start names "nope", which is not a step`},
		{"no steps", "name = \"empty\"\n",
			`workflow has no steps`},
		{"syntax", "[step.a\n", "toml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.src))
			if err == nil {
				t.Fatal("Parse succeeded, want an error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not contain %q", err, tt.want)
			}
		})
	}
}

func TestCap(t *testing.T) {
	for _, tt := range []struct{ max, want int }{{0, 1}, {1, 1}, {3, 3}} {
		if got := (Step{MaxVisits: tt.max}).Cap(); got != tt.want {
			t.Errorf("Cap with MaxVisits %d = %d, want %d", tt.max, got, tt.want)
		}
	}
}

func TestNext(t *testing.T) {
	d, err := Load(filepath.Join("testdata", "valid.toml"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		step, outcome string
		next          string
		end           End
	}{
		{"build", "pass", "test", EndNone}, // entry to a step
		{"test", "fail", "build", EndNone}, // entry loops back
		{"ci", "pass", "", EndSucceeded},   // entry to done
		{"ci", "timeout", "", EndFailed},   // entry to failed
		{"build", "fail", "", EndFailed},   // no entry, not pass
		{"test", "error", "", EndFailed},   // no entry, not pass
		{"nope", "pass", "", EndFailed},    // not a step
	}
	for _, tt := range tests {
		next, end := d.Next(tt.step, tt.outcome)
		if next != tt.next || end != tt.end {
			t.Errorf("Next(%q, %q) = (%q, %v), want (%q, %v)", tt.step, tt.outcome, next, end, tt.next, tt.end)
		}
	}

	// An outcome with no entry ends a pass as succeeded.
	bare, err := Parse("bare", []byte("[step.a]\nkind = \"command\"\nrun = \"true\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	if next, end := bare.Next("a", "pass"); next != "" || end != EndSucceeded {
		t.Errorf("Next(a, pass) with no entry = (%q, %v), want (\"\", succeeded)", next, end)
	}
	if next, end := bare.Next("a", "fail"); next != "" || end != EndFailed {
		t.Errorf("Next(a, fail) with no entry = (%q, %v), want (\"\", failed)", next, end)
	}
}

func TestEndString(t *testing.T) {
	for e, want := range map[End]string{EndNone: "", EndSucceeded: "succeeded", EndFailed: "failed"} {
		if got := e.String(); got != want {
			t.Errorf("End(%d).String() = %q, want %q", e, got, want)
		}
	}
}

func TestOutcomeNames(t *testing.T) {
	tests := []struct {
		name, src, want string
	}{
		{"wait has no fail", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\non.fail = \"done\"\n",
			`step "w": on.fail: a wait step has no outcome "fail"; want on.pass or on.timeout`},
		{"agent has no timeout", "[step.a]\nkind = \"agent\"\nprompt = \"go\"\non.timeout = \"done\"\n",
			`step "a": on.timeout: an agent step`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("o", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
	ok := "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\non.timeout = \"failed\"\non.pass = \"done\"\n"
	if _, err := Parse("o", []byte(ok)); err != nil {
		t.Errorf("wait with on.pass and on.timeout: %v", err)
	}
}

func TestAgentPrompt(t *testing.T) {
	refused := []struct {
		name, src, want string
	}{
		{"no prompt", "[step.a]\nkind = \"agent\"\n", `step "a": an agent step needs prompt = `},
		{"blank prompt", "[step.a]\nkind = \"agent\"\nprompt = \"  \"\n", `step "a": an agent step needs prompt = `},
		{"unknown placeholder", "[step.a]\nkind = \"agent\"\nprompt = \"do {chnage}\"\n",
			`step "a": prompt uses unknown placeholder {chnage}; want one of {change}, {slug}, {title}, {run}, {step}`},
		{"prompt on a command", "[step.c]\nkind = \"command\"\nrun = \"true\"\nprompt = \"x\"\n",
			`step "c": prompt is only for agent steps`},
	}
	for _, tt := range refused {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("p", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}

	d, err := Parse("p", []byte("[step.build]\nkind = \"agent\"\nprompt = \"/shipit:go {change} {slug} {title} {run} {step} {x\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	got := d.Steps["build"].RenderPrompt(PromptVars{Change: "01ABC", Slug: "runs", Title: "Runs", Run: "r1", Step: "build"})
	if want := "/shipit:go 01ABC runs Runs r1 build {x"; got != want {
		t.Errorf("RenderPrompt = %q, want %q", got, want)
	}
}

func TestBuiltinsParse(t *testing.T) {
	names := BuiltinNames()
	if len(names) == 0 {
		t.Fatal("no built-in workflows")
	}
	for _, n := range names {
		data, ok := Builtin(n)
		if !ok {
			t.Fatalf("Builtin(%q) not found", n)
		}
		def, err := ParseWith(n, data, builtinResolver)
		if err != nil {
			t.Errorf("built-in %q does not parse: %v", n, err)
			continue
		}
		if def.Name != n {
			t.Errorf("built-in %q is named %q", n, def.Name)
		}
	}
	if _, ok := Builtin("nope"); ok {
		t.Error("Builtin(nope) found")
	}
	if _, ok := Builtin("../x"); ok {
		t.Error("Builtin(../x) found")
	}
}

func builtinResolver(name string) ([]byte, error) {
	if data, ok := Builtin(name); ok {
		return data, nil
	}
	return nil, fmt.Errorf("no built-in workflow %q", name)
}

func TestShipChangeBuiltin(t *testing.T) {
	data, ok := Builtin("ship-change")
	if !ok {
		t.Fatal("no ship-change built-in")
	}
	def, err := ParseWith("ship-change", data, builtinResolver)
	if err != nil {
		t.Fatal(err)
	}
	for name, step := range def.Steps {
		if strings.Contains(step.Prompt, "shipit") || strings.Contains(step.Run, "shipit") {
			t.Errorf("step %s mentions shipit", name)
		}
		if strings.Contains(step.Until, "--pr-open") {
			t.Errorf("step %s waits for ship check --pr-open", name)
		}
	}
	if strings.Contains(string(data), "shipit") {
		t.Error("ship-change.toml mentions shipit")
	}
	for _, inner := range []string{"build", "review"} {
		if _, ok := def.Includes[inner]; !ok {
			t.Errorf("ship-change includes no %s workflow", inner)
		}
	}
	if _, ok := def.Steps["review/signoff"]; ok {
		t.Error("review's signoff step is not skipped")
	}
	want := []string{"build/build", "build/render", "build/check", "review/code", "review/review", "review/render", "render-pr", "open-pr", "wait-pr", "pr-green", "signoff", "rebase", "wait-rebased", "rebased-green", "land", "wait-landed", "landed-green",
		"rebase-merge", "wait-rebased-merge", "rebased-merge-green", "merge", "wait-main", "main-green", "close-out", "ship"}
	if !slices.Equal(def.Order, want) {
		t.Errorf("steps = %v; want %v", def.Order, want)
	}
}

// The built-in plan workflow: an interview driven by the rhiza-plan skill,
// then review rounds (a fresh reviewer told only the change set, a critic,
// rhiza run review-check, a reviser on a failed check), then the gate.
func TestPlanBuiltin(t *testing.T) {
	if !slices.Contains(BuiltinNames(), "plan") {
		t.Fatalf("BuiltinNames = %v; want it to list plan", BuiltinNames())
	}
	data, ok := Builtin("plan")
	if !ok {
		t.Fatal("no plan built-in")
	}
	def, err := Parse("plan", data)
	if err != nil {
		t.Fatal(err)
	}
	if def.Name != "plan" || def.Start != "interview" {
		t.Errorf("name %q, start %q; want plan, interview", def.Name, def.Start)
	}
	want := []string{"interview", "milestone-check", "review", "critic", "check", "revise", "gate", "open-pr", "wait-pr", "pr-green", "rebase", "wait-rebased", "rebased-green", "merge", "close-out"}
	if !slices.Equal(def.Order, want) {
		t.Fatalf("steps = %v; want %v", def.Order, want)
	}

	interview := def.Steps["interview"]
	if interview.Kind != KindAgent || !interview.WaitsForDone() || interview.Until == "" {
		t.Errorf("interview = %+v; want an agent step that waits for done", interview)
	}
	for _, w := range []string{"rhiza-plan", "{change}", "rhiza guide plan", "rhiza session done"} {
		if !strings.Contains(interview.Prompt, w) {
			t.Errorf("interview prompt lacks %q: %q", w, interview.Prompt)
		}
	}

	roles := map[string]string{"review": "reviewer", "critic": "critic", "revise": "reviser"}
	for name, role := range roles {
		s := def.Steps[name]
		if s.Kind != KindAgent || !s.WaitsForDone() || s.Until == "" {
			t.Errorf("%s = %+v; want an agent step that waits for done", name, s)
		}
		var ph []string
		for _, m := range placeholderPattern.FindAllStringSubmatch(s.Prompt, -1) {
			ph = append(ph, m[1])
		}
		if !slices.Equal(ph, []string{"change"}) {
			t.Errorf("%s prompt placeholders = %v; want only {change}: %q", name, ph, s.Prompt)
		}
		for _, w := range []string{"rhiza guide plan-review", role, "rhiza session done"} {
			if !strings.Contains(s.Prompt, w) {
				t.Errorf("%s prompt lacks %q: %q", name, w, s.Prompt)
			}
		}
		for _, bad := range []string{"plan file", "rhiza-plan", ".md"} {
			if strings.Contains(s.Prompt, bad) {
				t.Errorf("%s prompt names %q: %q", name, bad, s.Prompt)
			}
		}
	}

	check := def.Steps["check"]
	if check.Kind != KindCommand || check.Run != "rhiza run review-check" {
		t.Errorf("check = %+v; want command rhiza run review-check", check)
	}
	gate := def.Steps["gate"]
	if gate.Kind != KindCommand || gate.Run != "rhiza run gate" {
		t.Errorf("gate = %+v; want command rhiza run gate", gate)
	}

	trans := []struct {
		step, outcome, next string
		end                 End
	}{
		{"interview", OutcomePass, "milestone-check", EndNone},
		{"milestone-check", OutcomePass, "review", EndNone},
		{"milestone-check", OutcomeFail, "", EndFailed},
		{"interview", OutcomeFail, "", EndFailed},
		{"review", OutcomePass, "critic", EndNone},
		{"critic", OutcomePass, "check", EndNone},
		{"check", OutcomePass, "gate", EndNone},
		{"check", OutcomeFail, "revise", EndNone},
		{"revise", OutcomePass, "review", EndNone},
		{"revise", OutcomeFail, "", EndFailed},
		{"gate", OutcomePass, "open-pr", EndNone},
		{"gate", OutcomeFail, "", EndFailed},
		{"open-pr", OutcomePass, "wait-pr", EndNone},
		{"open-pr", OutcomeFail, "", EndFailed},
		{"wait-pr", OutcomePass, "pr-green", EndNone},
		{"pr-green", OutcomePass, "rebase", EndNone},
		{"rebase", OutcomePass, "wait-rebased", EndNone},
		{"rebase", OutcomeFail, "", EndFailed},
		{"wait-rebased", OutcomePass, "rebased-green", EndNone},
		{"rebased-green", OutcomePass, "merge", EndNone},
		{"merge", OutcomePass, "close-out", EndNone},
		{"close-out", OutcomePass, "", EndSucceeded},
	}
	for _, tt := range trans {
		next, end := def.Next(tt.step, tt.outcome)
		if next != tt.next || end != tt.end {
			t.Errorf("Next(%q, %q) = (%q, %v); want (%q, %v)", tt.step, tt.outcome, next, end, tt.next, tt.end)
		}
	}

	caps := map[string]int{"interview": 1, "review": 3, "critic": 3, "check": 3, "revise": 2, "gate": 1, "open-pr": 3, "wait-pr": 3, "pr-green": 3}
	for name, c := range caps {
		if got := def.Steps[name].Cap(); got != c {
			t.Errorf("%s cap = %d; want %d", name, got, c)
		}
	}
}

// The built-in build workflow: a leaves step (builder models sonnet,
// opus, opus; verifier opus; split auto) whose prompts follow rhiza guide
// build, then rhiza run build-check, ending done with no land step.
func TestBuildBuiltin(t *testing.T) {
	if !slices.Contains(BuiltinNames(), "build") {
		t.Fatalf("BuiltinNames = %v; want it to list build", BuiltinNames())
	}
	data, ok := Builtin("build")
	if !ok {
		t.Fatal("no build built-in")
	}
	def, err := Parse("build", data)
	if err != nil {
		t.Fatal(err)
	}
	if def.Name != "build" || def.Start != "build" {
		t.Errorf("name %q, start %q; want build, build", def.Name, def.Start)
	}
	if want := []string{"build", "render", "check"}; !slices.Equal(def.Order, want) {
		t.Fatalf("steps = %v; want %v", def.Order, want)
	}
	build := def.Steps["build"]
	if build.Kind != KindLeaves || build.Split != SplitAuto || build.Builder == nil || build.Verifier == nil {
		t.Fatalf("build = %+v; want a leaves step with split auto and both roles", build)
	}
	if want := []string{"sonnet", "opus", "opus"}; !slices.Equal(build.Builder.Models, want) {
		t.Errorf("builder models = %v; want %v", build.Builder.Models, want)
	}
	if want := []string{"opus"}; !slices.Equal(build.Verifier.Models, want) {
		t.Errorf("verifier models = %v; want %v", build.Verifier.Models, want)
	}
	for role, r := range map[string]*Role{"builder": build.Builder, "verifier": build.Verifier} {
		for _, w := range []string{"rhiza guide build", "{leaf}", "{change}", role, "rhiza session done"} {
			if !strings.Contains(r.Prompt, w) {
				t.Errorf("%s prompt lacks %q: %q", role, w, r.Prompt)
			}
		}
		if strings.Contains(strings.ToLower(r.Prompt), "claude") {
			t.Errorf("%s prompt names an agent product: %q", role, r.Prompt)
		}
	}
	for _, w := range []string{"{attempt}", "Leaf: {leaf}"} {
		if !strings.Contains(build.Builder.Prompt, w) {
			t.Errorf("builder prompt lacks %q", w)
		}
	}
	for _, w := range []string{"$RHIZA_COMMIT", "rhiza run findings --leaf {leaf}", "Never edit code"} {
		if !strings.Contains(build.Verifier.Prompt, w) {
			t.Errorf("verifier prompt lacks %q", w)
		}
	}
	check := def.Steps["check"]
	if check.Kind != KindCommand || check.Run != "rhiza run build-check" {
		t.Errorf("check = %+v; want command rhiza run build-check", check)
	}
	for _, s := range def.Steps {
		if strings.Contains(s.Run, "land") {
			t.Errorf("step %s lands: %q", s.Name, s.Run)
		}
	}
	trans := []struct {
		step, outcome, next string
		end                 End
	}{
		{"build", OutcomePass, "render", EndNone},
		{"render", OutcomePass, "check", EndNone},
		{"render", OutcomeFail, "", EndFailed},
		{"build", OutcomeFail, "", EndFailed},
		{"check", OutcomePass, "", EndSucceeded},
		{"check", OutcomeFail, "", EndFailed},
	}
	for _, tt := range trans {
		next, end := def.Next(tt.step, tt.outcome)
		if next != tt.next || end != tt.end {
			t.Errorf("Next(%q, %q) = (%q, %v); want (%q, %v)", tt.step, tt.outcome, next, end, tt.next, tt.end)
		}
	}
}

// An agent step may have until; its interval then defaults to
// AgentUntilInterval, and a set interval still wins.
func TestAgentStepUntil(t *testing.T) {
	src := "[step.g]\nkind = \"agent\"\nprompt = \"go\"\nuntil = \"true\"\n"
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if g := d.Steps["g"]; g.Until != "true" || g.Interval != AgentUntilInterval {
		t.Errorf("agent step with until = %+v; want until \"true\" every %v", g, AgentUntilInterval)
	}
	d, err = Parse("w", []byte(src+"interval = \"5s\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Steps["g"].Interval; got != 5*time.Second {
		t.Errorf("interval = %v, want the file's 5s", got)
	}
	d, err = Parse("w", []byte("[step.g]\nkind = \"agent\"\nprompt = \"go\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Steps["g"].Interval; got != DefaultInterval {
		t.Errorf("agent step without until: interval = %v, want %v", got, DefaultInterval)
	}
}

// until_when says when an agent step checks its until: idle (the default)
// or done. It is only for agent steps with an until.
func TestAgentStepUntilWhen(t *testing.T) {
	src := "[step.g]\nkind = \"agent\"\nprompt = \"go\"\nuntil = \"true\"\n"
	for _, tt := range []struct {
		line      string
		when      string
		waitsDone bool
	}{
		{"", "", false},
		{"until_when = \"idle\"\n", UntilWhenIdle, false},
		{"until_when = \"done\"\n", UntilWhenDone, true},
	} {
		d, err := Parse("w", []byte(src+tt.line))
		if err != nil {
			t.Fatalf("%q: %v", tt.line, err)
		}
		g := d.Steps["g"]
		if g.UntilWhen != tt.when || g.WaitsForDone() != tt.waitsDone {
			t.Errorf("%q: until_when %q, waits for done %v", tt.line, g.UntilWhen, g.WaitsForDone())
		}
	}
	for _, tt := range []struct {
		name, src, want string
	}{
		{"unknown value", src + "until_when = \"soon\"\n",
			`step "g": until_when "soon": want "idle" (check until whenever the session is idle) or "done"`},
		{"no until", "[step.g]\nkind = \"agent\"\nprompt = \"go\"\nuntil_when = \"done\"\n",
			`step "g": until_when needs until = "<shell command>" on the same step`},
		{"wait step", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\nuntil_when = \"done\"\n",
			`step "w": until_when is only for agent steps; a wait step cannot have it`},
		{"command step", "[step.a]\nkind = \"command\"\nrun = \"true\"\nuntil_when = \"idle\"\n",
			`step "a": until_when is only for agent steps; a command step cannot have it`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
}

// arbiterBase is a workflow whose agent step starts it, for the arbiter
// step tests to add an arbiter step to.
const arbiterBase = "[step.build]\nkind = \"agent\"\nprompt = \"go\"\n"

// A file may declare one arbiter step, with agent, permission_mode and
// interval; it defaults to DefaultInterval and produces pass or fail.
func TestArbiterStepLoads(t *testing.T) {
	src := arbiterBase + "[step.decide]\nkind = \"arbiter\"\nagent = \"x\"\npermission_mode = \"manual\"\ninterval = \"5s\"\n"
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := Step{Name: "decide", Kind: KindArbiter, Agent: "x", PermissionMode: "manual", Interval: 5 * time.Second, On: map[string]string{}}
	if got := d.Steps["decide"]; !reflect.DeepEqual(got, want) {
		t.Errorf("arbiter step = %+v\nwant %+v", got, want)
	}
	if got := d.Arbiter(); !reflect.DeepEqual(got, want) {
		t.Errorf("Arbiter() = %+v, want the file's step", got)
	}
	if d.Start != "build" {
		t.Errorf("start = %q, want build", d.Start)
	}
	if want := []string{OutcomePass, OutcomeFail}; !slices.Equal(kindOutcomes[KindArbiter], want) {
		t.Errorf("arbiter outcomes = %v, want %v", kindOutcomes[KindArbiter], want)
	}

	d, err = Parse("w", []byte(arbiterBase+"[step.decide]\nkind = \"arbiter\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Arbiter().Interval; got != DefaultInterval {
		t.Errorf("arbiter interval = %v, want %v", got, DefaultInterval)
	}
}

// Without an arbiter step, Arbiter gives the built-in one: the run's agent
// and mode, every DefaultInterval.
func TestArbiterBuiltin(t *testing.T) {
	d, err := Parse("w", []byte(arbiterBase))
	if err != nil {
		t.Fatal(err)
	}
	want := Step{Name: ArbiterName, Kind: KindArbiter, Interval: DefaultInterval, On: map[string]string{}}
	if got := d.Arbiter(); !reflect.DeepEqual(got, want) {
		t.Errorf("Arbiter() = %+v, want %+v", got, want)
	}
	if _, ok := d.Steps[ArbiterName]; ok {
		t.Error("the built-in arbiter became a step of the file")
	}
}

// An arbiter step refuses a prompt and every key but agent,
// permission_mode and interval; a file has at most one, it is never the
// start step and no on. goes to it.
func TestArbiterStepRefused(t *testing.T) {
	arb := "[step.decide]\nkind = \"arbiter\"\n"
	for _, tt := range []struct {
		name, src, want string
	}{
		{"prompt", arbiterBase + arb + "prompt = \"decide {change}\"\n",
			`step "decide": an arbiter step cannot have prompt: its prompt is fixed`},
		{"run", arbiterBase + arb + "run = \"true\"\n", `step "decide": an arbiter step cannot have run`},
		{"until", arbiterBase + arb + "until = \"true\"\n", `step "decide": an arbiter step cannot have until`},
		{"until_when", arbiterBase + arb + "until_when = \"done\"\n", `step "decide": an arbiter step cannot have until_when`},
		{"timeout", arbiterBase + arb + "timeout = \"1m\"\n", `step "decide": an arbiter step cannot have timeout`},
		{"max_visits", arbiterBase + arb + "max_visits = 2\n", `step "decide": an arbiter step cannot have max_visits`},
		{"on", arbiterBase + arb + "on.pass = \"build\"\n", `step "decide": on.pass: an arbiter step cannot have on.`},
		{"bad interval", arbiterBase + arb + "interval = \"soon\"\n", `step "decide": interval "soon" is not a duration`},
		{"two arbiters", arbiterBase + arb + "[step.again]\nkind = \"arbiter\"\n",
			`steps "decide" and "again" are both arbiter steps; a workflow has at most one`},
		{"first step", arb + arbiterBase,
			`step "decide" is an arbiter step and cannot be the start step`},
		{"start", "start = \"decide\"\n" + arbiterBase + arb,
			`step "decide" is an arbiter step and cannot be the start step`},
		{"target", "[step.build]\nkind = \"agent\"\nprompt = \"go\"\non.fail = \"decide\"\n" + arb,
			`step "build": on.fail names "decide", an arbiter step; no step goes to it`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
}

const leavesBase = "[step.build]\nkind = \"leaves\"\n"
const leavesRoles = "[step.build.builder]\nprompt = \"build {leaf} of {change}, attempt {attempt}\"\n" +
	"[step.build.verifier]\nprompt = \"verify {leaf}\"\n"

// A leaves step loads with both sub-tables, its split, models per role
// and the leaves prompt placeholders.
func TestLeavesStepLoads(t *testing.T) {
	src := "[step.build]\nkind = \"leaves\"\nsplit = \"below\"\ninterval = \"5s\"\non.fail = \"failed\"\n" +
		"[step.build.builder]\nprompt = \"build {leaf} of {change}, attempt {attempt}\"\nagent = \"claude-code\"\n" +
		"permission_mode = \"acceptEdits\"\nmodels = [\"sonnet\", \"opus\"]\n" +
		"[step.build.verifier]\nprompt = \"verify {leaf}\"\nmodels = [\"opus\"]\n"
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := Step{
		Name: "build", Kind: KindLeaves, Split: SplitBelow, Interval: 5 * time.Second,
		On: map[string]string{"fail": "failed"},
		Builder: &Role{Prompt: "build {leaf} of {change}, attempt {attempt}", Agent: "claude-code",
			PermissionMode: "acceptEdits", Models: []string{"sonnet", "opus"}},
		Verifier: &Role{Prompt: "verify {leaf}", Models: []string{"opus"}},
	}
	if got := d.Steps["build"]; !reflect.DeepEqual(got, want) {
		t.Errorf("step = %+v\nwant %+v", got, want)
	}
	got := d.Steps["build"].Builder.RenderPrompt(LeafVars{PromptVars: PromptVars{Change: "01ABC"}, Leaf: "leaf-a", Attempt: 2})
	if want := "build leaf-a of 01ABC, attempt 2"; got != want {
		t.Errorf("RenderPrompt = %q, want %q", got, want)
	}

	d, err = Parse("w", []byte(leavesBase+leavesRoles))
	if err != nil {
		t.Fatal(err)
	}
	if s := d.Steps["build"]; s.Split != SplitAuto || s.Interval != DefaultInterval {
		t.Errorf("defaults: split %q, interval %v; want %q, %v", s.Split, s.Interval, SplitAuto, DefaultInterval)
	}

	d, err = Parse("w", []byte("[step.a]\nkind = \"agent\"\nprompt = \"go\"\nmodels = [\"opus\"]\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Steps["a"].Models; !reflect.DeepEqual(got, []string{"opus"}) {
		t.Errorf("agent step Models = %v, want [opus]", got)
	}
}

// Role.Model takes the attempt's model, keeps the last for later
// attempts, and is empty without models.
func TestLeavesRoleModel(t *testing.T) {
	r := Role{Models: []string{"sonnet", "opus", "opus-max"}}
	for attempt, want := range map[int]string{0: "sonnet", 1: "sonnet", 2: "opus", 3: "opus-max", 7: "opus-max"} {
		if got := r.Model(attempt); got != want {
			t.Errorf("Model(%d) = %q, want %q", attempt, got, want)
		}
	}
	if got := (Role{}).Model(1); got != "" {
		t.Errorf("Model without models = %q, want empty", got)
	}
}

// Step.Model takes the visit's model as Role.Model takes the attempt's.
func TestStepModel(t *testing.T) {
	s := Step{Models: []string{"sonnet", "opus"}}
	for visit, want := range map[int]string{1: "sonnet", 2: "opus", 3: "opus"} {
		if got := s.Model(visit); got != want {
			t.Errorf("Model(%d) = %q, want %q", visit, got, want)
		}
	}
	if got := (Step{Models: []string{"opus"}}).Model(4); got != "opus" {
		t.Errorf("one model, visit 4 = %q, want opus", got)
	}
	if got := (Step{}).Model(1); got != "" {
		t.Errorf("Model without models = %q, want empty", got)
	}
}

// A leaves step refuses a missing role, an unknown sub-table key, models
// and split where they do not belong, a bad split, and the session keys
// that belong in its sub-tables.
func TestLeavesStepRefused(t *testing.T) {
	builder := "[step.build.builder]\nprompt = \"build {leaf}\"\n"
	verifier := "[step.build.verifier]\nprompt = \"verify {leaf}\"\n"
	cmd := "[step.c]\nkind = \"command\"\nrun = \"true\"\n"
	wait := "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\n"
	agent := "[step.a]\nkind = \"agent\"\nprompt = \"go\"\n"
	for _, tt := range []struct {
		name, src, want string
	}{
		{"missing builder", leavesBase + verifier,
			`step "build": a leaves step needs a [step.build.builder] table with prompt = `},
		{"missing verifier", leavesBase + builder,
			`step "build": a leaves step needs a [step.build.verifier] table with prompt = `},
		{"blank builder prompt", leavesBase + "[step.build.builder]\nprompt = \" \"\n" + verifier,
			`step "build": builder needs prompt = "<what the builder should do>"; it may use {change}, {slug}, {title}, {run}, {step}, {leaf}, {attempt}`},
		{"unknown placeholder", leavesBase + builder + "[step.build.verifier]\nprompt = \"verify {leef}\"\n",
			`step "build": verifier: prompt uses unknown placeholder {leef}`},
		{"unknown sub-table key", leavesBase + builder + verifier + "retries = 3\n",
			`step "build": verifier: unknown key "retries"; want prompt, agent, permission_mode or models`},
		{"models on a wait step", wait + "models = [\"opus\"]\n",
			`step "w": models is only for agent steps and a leaves step's builder and verifier; a wait step cannot have it`},
		{"models on a command step", cmd + "models = [\"opus\"]\n",
			`step "c": models is only for agent steps and a leaves step's builder and verifier; a command step cannot have it`},
		{"models on a leaves step", leavesBase + "models = [\"opus\"]\n" + builder + verifier,
			`step "build": models is only for agent steps and a leaves step's builder and verifier; set it under [step.build.builder] or [step.build.verifier]`},
		{"blank model", leavesBase + builder + verifier + "models = [\"opus\", \"\"]\n",
			`step "build": verifier: models entry 2 is blank`},
		{"split on an agent step", agent + "split = \"right\"\n",
			`step "a": split is only for leaves and checks steps; an agent step cannot have it`},
		{"bad split", leavesBase + "split = \"diagonal\"\n" + builder + verifier,
			`step "build": split "diagonal": want "auto" (by the tab's shape), "right" or "below"`},
		{"builder on an agent step", agent + "[step.a.builder]\nprompt = \"x\"\n",
			`step "a": builder is only for leaves steps; an agent step cannot have it`},
		{"prompt on a leaves step", leavesBase + "prompt = \"x\"\n" + builder + verifier,
			`step "build": a leaves step cannot have prompt: each role has its own`},
		{"agent on a leaves step", leavesBase + "agent = \"x\"\n" + builder + verifier,
			`step "build": a leaves step cannot have agent: each role has its own`},
		{"permission_mode on a leaves step", leavesBase + "permission_mode = \"x\"\n" + builder + verifier,
			`step "build": a leaves step cannot have permission_mode: each role has its own`},
		{"run on a leaves step", leavesBase + "run = \"x\"\n" + builder + verifier,
			`step "build": a leaves step cannot have run`},
		{"until on a leaves step", leavesBase + "until = \"x\"\n" + builder + verifier,
			`step "build": a leaves step cannot have until`},
		{"until_when on a leaves step", leavesBase + "until_when = \"done\"\n" + builder + verifier,
			`step "build": a leaves step cannot have until_when`},
		{"leaf placeholder on an agent step", "[step.a]\nkind = \"agent\"\nprompt = \"do {leaf}\"\n",
			`step "a": prompt uses unknown placeholder {leaf}`},
		{"bad outcome", leavesBase + "on.timeout = \"done\"\n" + builder + verifier,
			`step "build": on.timeout: a leaves step has no outcome "timeout"; want on.pass or on.fail`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
}

// checksBase is a checks step with two checks; checksRoles its critic,
// fixer and verifier.
const checksBase = "[step.review]\nkind = \"checks\"\n" + checksList

// checksList are checksBase's two checks.
const checksList = "[[step.review.check]]\nname = \"plan-compliance\"\nprompt = \"check {check} of {change}, cycle {cycle}\"\n" +
	"[[step.review.check]]\nname = \"code-quality\"\nprompt = \"check {check}\"\nagent = \"claude-code\"\npermission_mode = \"acceptEdits\"\nmodels = [\"sonnet\", \"opus\"]\n"

const checksRoles = "[step.review.critic]\nprompt = \"criticise {check} in cycle {cycle}\"\n" +
	"[step.review.fixer]\nprompt = \"fix {leaf} in cycle {cycle}\"\nmodels = [\"opus\"]\n" +
	"[step.review.verifier]\nprompt = \"verify {leaf}, attempt {attempt}\"\n"

// A checks step loads its checks in file order, each a role of its own,
// and its critic, fixer and verifier; it takes split, interval and
// max_visits, and renders the check, critic and fixer placeholders.
func TestChecksStepLoads(t *testing.T) {
	d, err := Parse("w", []byte("[step.review]\nkind = \"checks\"\nsplit = \"right\"\ninterval = \"5s\"\nmax_visits = 3\non.fail = \"failed\"\n"+checksList+checksRoles))
	if err != nil {
		t.Fatal(err)
	}
	want := Step{
		Name: "review", Kind: KindChecks, Split: SplitRight, Interval: 5 * time.Second, MaxVisits: 3,
		On: map[string]string{"fail": "failed"},
		Checks: []Check{
			{Name: "plan-compliance", Role: Role{Prompt: "check {check} of {change}, cycle {cycle}"}},
			{Name: "code-quality", Role: Role{Prompt: "check {check}", Agent: "claude-code", PermissionMode: "acceptEdits", Models: []string{"sonnet", "opus"}}},
		},
		Critic:   &Role{Prompt: "criticise {check} in cycle {cycle}"},
		Fixer:    &Role{Prompt: "fix {leaf} in cycle {cycle}", Models: []string{"opus"}},
		Verifier: &Role{Prompt: "verify {leaf}, attempt {attempt}"},
	}
	s := d.Steps["review"]
	if !reflect.DeepEqual(s, want) {
		t.Errorf("step = %+v\nwant %+v", s, want)
	}
	if c, ok := s.Check("code-quality"); !ok || c.Agent != "claude-code" {
		t.Errorf("Check(code-quality) = %+v, %v", c, ok)
	}
	if _, ok := s.Check("nope"); ok {
		t.Errorf("Check(nope) found")
	}
	vars := PromptVars{Change: "01ABC"}
	if got := s.Checks[0].RenderCheckPrompt(CheckVars{PromptVars: vars, Check: "plan-compliance", Cycle: 2}); got != "check plan-compliance of 01ABC, cycle 2" {
		t.Errorf("check prompt = %q", got)
	}
	if got := s.Critic.RenderCheckPrompt(CheckVars{PromptVars: vars, Check: "code-quality", Cycle: 1}); got != "criticise code-quality in cycle 1" {
		t.Errorf("critic prompt = %q", got)
	}
	if got := s.Fixer.RenderFixerPrompt(FixerVars{PromptVars: vars, Leaf: "leaf-a", Cycle: 2}); got != "fix leaf-a in cycle 2" {
		t.Errorf("fixer prompt = %q", got)
	}
	if s.Role(RoleCritic) != s.Critic || s.Role(RoleFixer) != s.Fixer || s.Role(RoleVerifier) != s.Verifier {
		t.Errorf("Role does not give the critic, fixer and verifier")
	}

	d, err = Parse("w", []byte(checksBase+checksRoles))
	if err != nil {
		t.Fatal(err)
	}
	if s := d.Steps["review"]; s.Split != SplitAuto || s.Interval != DefaultInterval {
		t.Errorf("defaults: split %q, interval %v; want %q, %v", s.Split, s.Interval, SplitAuto, DefaultInterval)
	}
}

// A checks step refuses no checks, a duplicate or invalid check name, a
// missing or blank role, an unknown key in a check or a role, unknown
// placeholders, and the keys that belong elsewhere; other kinds refuse
// its sub-tables.
func TestChecksStepRefused(t *testing.T) {
	head := "[step.review]\nkind = \"checks\"\n"
	check := func(name string) string {
		return "[[step.review.check]]\nname = \"" + name + "\"\nprompt = \"check {check}\"\n"
	}
	critic := "[step.review.critic]\nprompt = \"criticise {check}\"\n"
	fixer := "[step.review.fixer]\nprompt = \"fix {leaf}\"\n"
	verifier := "[step.review.verifier]\nprompt = \"verify {leaf}\"\n"
	roles := critic + fixer + verifier
	agent := "[step.a]\nkind = \"agent\"\nprompt = \"go\"\n"
	for _, tt := range []struct {
		name, src, want string
	}{
		{"no checks", head + roles,
			`step "review": a checks step needs at least one [[step.review.check]] table with name = "..." and prompt = "..."`},
		{"duplicate check name", head + check("lint") + check("lint") + roles,
			`step "review": check "lint" is declared twice; give each check its own name`},
		{"invalid check name", head + check("Lint Me") + roles,
			`step "review": check "Lint Me": want lowercase letters and digits in words joined by single hyphens`},
		{"check with no name", head + "[[step.review.check]]\nprompt = \"x\"\n" + roles,
			`step "review": check 1 has no name: set name = "..."`},
		{"blank check prompt", head + "[[step.review.check]]\nname = \"lint\"\nprompt = \" \"\n" + roles,
			`step "review": check "lint" needs prompt = "<what the check should do>"; it may use {change}, {slug}, {title}, {run}, {step}, {check}, {cycle}`},
		{"unknown check key", head + check("lint") + "retries = 3\n" + roles,
			`step "review": check: unknown key "retries"; want name, prompt, agent, permission_mode or models`},
		{"unknown critic key", head + check("lint") + critic + "retries = 3\n" + fixer + verifier,
			`step "review": critic: unknown key "retries"; want prompt, agent, permission_mode or models`},
		{"unknown fixer key", head + check("lint") + critic + fixer + "retries = 3\n" + verifier,
			`step "review": fixer: unknown key "retries"; want prompt, agent, permission_mode or models`},
		{"missing critic", head + check("lint") + fixer + verifier,
			`step "review": a checks step needs a [step.review.critic] table with prompt = `},
		{"missing fixer", head + check("lint") + critic + verifier,
			`step "review": a checks step needs a [step.review.fixer] table with prompt = `},
		{"missing verifier", head + check("lint") + critic + fixer,
			`step "review": a checks step needs a [step.review.verifier] table with prompt = `},
		{"blank critic prompt", head + check("lint") + "[step.review.critic]\nprompt = \"\"\n" + fixer + verifier,
			`step "review": critic needs prompt = "<what the critic should do>"`},
		{"unknown check placeholder", head + "[[step.review.check]]\nname = \"lint\"\nprompt = \"do {leaf}\"\n" + roles,
			`step "review": check "lint": prompt uses unknown placeholder {leaf}; want one of {change}, {slug}, {title}, {run}, {step}, {check}, {cycle}`},
		{"check placeholder on the fixer", head + check("lint") + critic + "[step.review.fixer]\nprompt = \"fix {check}\"\n" + verifier,
			`step "review": fixer: prompt uses unknown placeholder {check}; want one of {change}, {slug}, {title}, {run}, {step}, {leaf}, {cycle}`},
		{"cycle placeholder on the verifier", head + check("lint") + critic + fixer + "[step.review.verifier]\nprompt = \"verify {cycle}\"\n",
			`step "review": verifier: prompt uses unknown placeholder {cycle}`},
		{"blank check model", head + check("lint") + "models = [\"\"]\n" + roles,
			`step "review": check "lint": models entry 1 is blank`},
		{"builder on a checks step", head + check("lint") + roles + "[step.review.builder]\nprompt = \"x\"\n",
			`step "review": a checks step cannot have builder: a fixer fixes what its checks find; remove it`},
		{"models on a checks step", head + "models = [\"opus\"]\n" + check("lint") + roles,
			`step "review": a checks step cannot have models: each check and role has its own`},
		{"prompt on a checks step", head + "prompt = \"x\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have prompt: each check and role has its own`},
		{"agent on a checks step", head + "agent = \"x\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have agent`},
		{"permission_mode on a checks step", head + "permission_mode = \"x\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have permission_mode`},
		{"run on a checks step", head + "run = \"x\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have run`},
		{"until on a checks step", head + "until = \"x\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have until`},
		{"until_when on a checks step", head + "until_when = \"done\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have until_when`},
		{"timeout on a checks step", head + "timeout = \"1m\"\n" + check("lint") + roles,
			`step "review": a checks step cannot have timeout`},
		{"bad split", head + "split = \"diagonal\"\n" + check("lint") + roles,
			`step "review": split "diagonal": want "auto" (by the tab's shape), "right" or "below"`},
		{"bad outcome", head + "on.timeout = \"done\"\n" + check("lint") + roles,
			`step "review": on.timeout: a checks step has no outcome "timeout"; want on.pass or on.fail`},
		{"check on an agent step", agent + "[[step.a.check]]\nname = \"lint\"\nprompt = \"x\"\n",
			`step "a": check is only for checks steps; an agent step cannot have it`},
		{"critic on an agent step", agent + "[step.a.critic]\nprompt = \"x\"\n",
			`step "a": critic is only for checks steps; an agent step cannot have it`},
		{"fixer on a leaves step", leavesBase + leavesRoles + "[step.build.fixer]\nprompt = \"x\"\n",
			`step "build": fixer is only for checks steps; a leaves step cannot have it`},
		{"verifier on an agent step", agent + "[step.a.verifier]\nprompt = \"x\"\n",
			`step "a": verifier is only for leaves and checks steps; an agent step cannot have it`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
}

// signoffBase is a checks step and a signoff step going back to it.
const signoffBase = "[step.review]\nkind = \"checks\"\non.pass = \"signoff\"\n" + checksList + checksRoles +
	"[step.signoff]\nkind = \"signoff\"\nmax_visits = 3\non.pass = \"done\"\non.changes = \"review\"\n"

// A signoff step has outcomes pass and changes, takes max_visits and
// on., and its on.changes names a checks step.
func TestSignoffStepLoads(t *testing.T) {
	d, err := Parse("w", []byte(signoffBase))
	if err != nil {
		t.Fatal(err)
	}
	want := Step{Name: "signoff", Kind: KindSignoff, MaxVisits: 3, On: map[string]string{"pass": "done", "changes": "review"}}
	if s := d.Steps["signoff"]; !reflect.DeepEqual(s, want) {
		t.Errorf("step = %+v\nwant %+v", s, want)
	}
	if next, end := d.Next("signoff", OutcomeChanges); next != "review" || end != EndNone {
		t.Errorf("Next(signoff, changes) = %q, %v", next, end)
	}
}

// A signoff step refuses every session and condition key, sub-tables, an
// outcome other than pass or changes, and a missing on.changes or one not
// naming a checks step.
func TestSignoffStepRefused(t *testing.T) {
	review := checksBase + checksRoles
	head := review + "[step.signoff]\nkind = \"signoff\"\n"
	changes := "on.changes = \"review\"\n"
	for _, tt := range []struct {
		name, src, want string
	}{
		{"no on.changes", head + "on.pass = \"done\"\n",
			`step "signoff": a signoff step needs on.changes = "<checks step>"`},
		{"on.changes to an agent step", review + "[step.a]\nkind = \"agent\"\nprompt = \"go\"\n[step.signoff]\nkind = \"signoff\"\non.changes = \"a\"\n",
			`step "signoff": on.changes names "a", which is not a checks step`},
		{"on.changes to done", head + "on.changes = \"done\"\n",
			`step "signoff": on.changes names "done", which is not a checks step`},
		{"on.fail", head + changes + "on.fail = \"failed\"\n",
			`step "signoff": on.fail: a signoff step has no outcome "fail"; want on.pass or on.changes`},
		{"prompt", head + changes + "prompt = \"x\"\n", `step "signoff": a signoff step cannot have prompt: it runs nothing`},
		{"run", head + changes + "run = \"x\"\n", `step "signoff": a signoff step cannot have run`},
		{"until", head + changes + "until = \"x\"\n", `step "signoff": a signoff step cannot have until`},
		{"until_when", head + changes + "until_when = \"done\"\n", `step "signoff": a signoff step cannot have until_when`},
		{"timeout", head + changes + "timeout = \"1m\"\n", `step "signoff": a signoff step cannot have timeout`},
		{"interval", head + changes + "interval = \"1m\"\n", `step "signoff": a signoff step cannot have interval`},
		{"agent", head + changes + "agent = \"x\"\n", `step "signoff": a signoff step cannot have agent`},
		{"permission_mode", head + changes + "permission_mode = \"x\"\n", `step "signoff": a signoff step cannot have permission_mode`},
		{"models", head + changes + "models = [\"opus\"]\n", `step "signoff": a signoff step cannot have models`},
		{"split", head + changes + "split = \"right\"\n", `step "signoff": split is only for leaves and checks steps; a signoff step cannot have it`},
		{"critic", head + changes + "[step.signoff.critic]\nprompt = \"x\"\n", `step "signoff": critic is only for checks steps; a signoff step cannot have it`},
		{"check", head + changes + "[[step.signoff.check]]\nname = \"x\"\nprompt = \"x\"\n", `step "signoff": check is only for checks steps; a signoff step cannot have it`},
		{"verifier", head + changes + "[step.signoff.verifier]\nprompt = \"x\"\n", `step "signoff": verifier is only for leaves and checks steps; a signoff step cannot have it`},
		{"max_visits", head + changes + "max_visits = -1\n", `step "signoff": max_visits is -1`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tt.src))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v; want %q", err, tt.want)
			}
		})
	}
}

// The built-in review workflow: code-check (no on.fail), then a checks
// step with compliance, quality, qa, scope and model checks, its critic,
// fixer and verifier, all on opus, each prompt following a guide and
// ending in rhiza session done, then the sign-off, which goes back to the
// checks on changes and ends done without landing.
func TestReviewBuiltin(t *testing.T) {
	if !slices.Contains(BuiltinNames(), "review") {
		t.Fatalf("BuiltinNames = %v; want it to list review", BuiltinNames())
	}
	data, ok := Builtin("review")
	if !ok {
		t.Fatal("no review built-in")
	}
	def, err := Parse("review", data)
	if err != nil {
		t.Fatal(err)
	}
	if def.Name != "review" || def.Start != "code" {
		t.Errorf("name %q, start %q; want review, code", def.Name, def.Start)
	}
	if want := []string{"code", "review", "render", "signoff"}; !slices.Equal(def.Order, want) {
		t.Fatalf("steps = %v; want %v", def.Order, want)
	}
	code := def.Steps["code"]
	if code.Kind != KindCommand || code.Run != "rhiza run code-check" {
		t.Errorf("code = %+v; want command rhiza run code-check", code)
	}
	if _, ok := code.On[OutcomeFail]; ok {
		t.Errorf("code has on.fail %q; a failed code check ends the run", code.On[OutcomeFail])
	}
	review := def.Steps["review"]
	if review.Kind != KindChecks || review.Split != SplitAuto || review.MaxVisits < 2 {
		t.Fatalf("review = %+v; want a checks step, split auto, max_visits for request-changes", review)
	}
	var names []string
	for _, c := range review.Checks {
		names = append(names, c.Name)
		if !slices.Equal(c.Models, []string{"opus"}) {
			t.Errorf("check %s models = %v; want [opus]", c.Name, c.Models)
		}
		var ph []string
		for _, m := range placeholderPattern.FindAllStringSubmatch(c.Prompt, -1) {
			ph = append(ph, m[1])
		}
		for _, p := range ph {
			if !slices.Contains([]string{"check", "change", "cycle"}, p) {
				t.Errorf("check %s prompt uses {%s}; want only {check}, {change}, {cycle}", c.Name, p)
			}
		}
		for _, w := range []string{"/rhiza-review", "{check}", "{change}", "rhiza guide review", "rhiza run findings --check", "rhiza session done"} {
			if !strings.Contains(c.Prompt, w) {
				t.Errorf("check %s prompt lacks %q: %q", c.Name, w, c.Prompt)
			}
		}
	}
	if want := []string{"compliance", "quality", "qa", "scope", "model"}; !slices.Equal(names, want) {
		t.Errorf("checks = %v; want %v", names, want)
	}
	roles := map[string]struct {
		r    *Role
		want []string
	}{
		"critic":   {review.Critic, []string{"/rhiza-review", "{check}", "{change}", "rhiza guide review", "rhiza run verdict --check", "never add a finding", "rhiza session done"}},
		"fixer":    {review.Fixer, []string{"/rhiza-review", "{change}", "{cycle}", "rhiza guide review", "rhiza session done"}},
		"verifier": {review.Verifier, []string{"/rhiza-build", "{leaf}", "{change}", "$RHIZA_COMMIT", "rhiza guide build", "Never edit code", "rhiza run findings --leaf {leaf}", "rhiza session done"}},
	}
	for name, tt := range roles {
		if tt.r == nil {
			t.Fatalf("review has no %s", name)
		}
		if !slices.Equal(tt.r.Models, []string{"opus"}) {
			t.Errorf("%s models = %v; want [opus]", name, tt.r.Models)
		}
		for _, w := range tt.want {
			if !strings.Contains(tt.r.Prompt, w) {
				t.Errorf("%s prompt lacks %q: %q", name, w, tt.r.Prompt)
			}
		}
	}
	// The fixer's prompt reads the same for a leaf's fixer and the change
	// set's ({leaf} empty): the appended findings name the leaf or the
	// change set.
	for _, l := range []string{"leaf-a", ""} {
		got := review.Fixer.RenderFixerPrompt(FixerVars{PromptVars: PromptVars{Change: "C1"}, Leaf: l, Cycle: 2})
		if strings.Contains(got, "  ") || strings.Contains(got, "()") || !strings.Contains(got, "change set C1 in review cycle 2:") {
			t.Errorf("fixer prompt with leaf %q = %q", l, got)
		}
	}
	signoff := def.Steps["signoff"]
	if signoff.Kind != KindSignoff || signoff.MaxVisits < 2 {
		t.Errorf("signoff = %+v; want a signoff step with max_visits for request-changes", signoff)
	}
	for _, s := range def.Steps {
		if strings.Contains(s.Run, "land") {
			t.Errorf("step %s lands: %q", s.Name, s.Run)
		}
	}
	prompts := []string{review.Critic.Prompt, review.Fixer.Prompt, review.Verifier.Prompt}
	for _, c := range review.Checks {
		prompts = append(prompts, c.Prompt)
	}
	for _, p := range prompts {
		if strings.Contains(strings.ToLower(p), "claude") {
			t.Errorf("a prompt names an agent product: %q", p)
		}
	}
	trans := []struct {
		step, outcome, next string
		end                 End
	}{
		{"code", OutcomePass, "review", EndNone},
		{"code", OutcomeFail, "", EndFailed},
		{"review", OutcomePass, "render", EndNone},
		{"render", OutcomePass, "signoff", EndNone},
		{"render", OutcomeFail, "", EndFailed},
		{"review", OutcomeFail, "", EndFailed},
		{"signoff", OutcomePass, "", EndSucceeded},
		{"signoff", OutcomeChanges, "review", EndNone},
	}
	for _, tt := range trans {
		next, end := def.Next(tt.step, tt.outcome)
		if next != tt.next || end != tt.end {
			t.Errorf("Next(%q, %q) = (%q, %v); want (%q, %v)", tt.step, tt.outcome, next, end, tt.next, tt.end)
		}
	}
}

func mapResolver(files map[string]string) Resolver {
	return func(name string) ([]byte, error) {
		src, ok := files[name]
		if !ok {
			return nil, fmt.Errorf("no workflow %q", name)
		}
		return []byte(src), nil
	}
}

const innerThree = `
[step.one]
kind = "command"
run = "true"
on.pass = "two"
[step.two]
kind = "command"
run = "true"
on.pass = "three"
[step.three]
kind = "command"
run = "true"
`

func TestWorkflowStepErrors(t *testing.T) {
	outer := func(extra string) string {
		return "[step.w]\nkind = \"workflow\"\nworkflow = \"inner\"\n" + extra
	}
	tests := []struct {
		name, src string
		files     map[string]string
		want      string
	}{
		{"unknown", outer(""), map[string]string{}, `step "w": workflow "inner": no workflow "inner"`},
		{"skip lacks", outer("skip = [\"nope\"]\n"), map[string]string{"inner": innerThree}, `skip names "nope", which workflow "inner" does not have`},
		{"self", outer(""), map[string]string{"inner": "[step.w]\nkind = \"workflow\"\nworkflow = \"test\"\n", "test": ""}, `includes itself: test -> inner -> test`},
		{"directly", outer(""), map[string]string{"inner": "[step.w]\nkind = \"workflow\"\nworkflow = \"inner\"\n"}, `includes itself: inner -> inner`},
		{"bad name", "[step.w]\nkind = \"workflow\"\nworkflow = \"a/b\"\n", nil, `want a workflow's name`},
		{"with run", outer("run = \"x\"\n"), nil, `a workflow step cannot have run`},
		{"bad outcome", outer("on.timeout = \"done\"\n"), map[string]string{"inner": innerThree}, `has no outcome "timeout"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseWith("test", []byte(tt.src), mapResolver(tt.files))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v, want one containing %q", err, tt.want)
			}
		})
	}
}

func TestWorkflowStepExpands(t *testing.T) {
	src := `
[step.pre]
kind = "command"
run = "true"
on.pass = "w"
[step.w]
kind = "workflow"
workflow = "inner"
skip = ["two"]
on.pass = "post"
on.fail = "cleanup"
[step.post]
kind = "command"
run = "true"
[step.cleanup]
kind = "command"
run = "true"
`
	d, err := ParseWith("test", []byte(src), mapResolver(map[string]string{"inner": innerThree}))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(d.Order, ","), "pre,inner/one,inner/three,post,cleanup"; got != want {
		t.Errorf("order %s, want %s", got, want)
	}
	for step, want := range map[string]string{"pre": "inner/one", "inner/one": "inner/three", "inner/three": "post"} {
		if next, _ := d.Next(step, OutcomePass); next != want {
			t.Errorf("%s on pass goes to %q, want %q", step, next, want)
		}
	}
	if next, _ := d.Next("inner/one", OutcomeFail); next != "cleanup" {
		t.Errorf("inner/one on fail goes to %q, want cleanup", next)
	}
	if _, ok := d.Steps["inner/two"]; ok {
		t.Error("skipped step inner/two is in the definition")
	}
	if _, ok := d.Includes["inner"]; !ok {
		t.Error("Includes lacks inner")
	}
	again, err := ParseWith("test", []byte(src), mapResolver(map[string]string{"inner": d.Includes["inner"]}))
	if err != nil || strings.Join(again.Order, ",") != strings.Join(d.Order, ",") {
		t.Errorf("reparse from Includes: %v, %v", again, err)
	}
}

func TestWorkflowStepAsLastAndStart(t *testing.T) {
	src := "[step.w]\nkind = \"workflow\"\nworkflow = \"inner\"\n"
	d, err := ParseWith("test", []byte(src), mapResolver(map[string]string{"inner": innerThree}))
	if err != nil {
		t.Fatal(err)
	}
	if d.Start != "inner/one" {
		t.Errorf("start %q, want inner/one", d.Start)
	}
	if _, end := d.Next("inner/three", OutcomePass); end != EndSucceeded {
		t.Errorf("last inner step on pass ends %v, want succeeded", end)
	}
	if _, end := d.Next("inner/one", OutcomeFail); end != EndFailed {
		t.Errorf("inner step on fail ends %v, want failed", end)
	}
}

func TestWorkflowStepNested(t *testing.T) {
	files := map[string]string{
		"mid":   "[step.m]\nkind = \"workflow\"\nworkflow = \"inner\"\non.pass = \"tail\"\n[step.tail]\nkind = \"command\"\nrun = \"true\"\n",
		"inner": innerThree,
	}
	src := "[step.w]\nkind = \"workflow\"\nworkflow = \"mid\"\nskip = [\"m\"]\n"
	d, err := ParseWith("test", []byte(src), mapResolver(files))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(d.Order, ","); got != "mid/tail" {
		t.Errorf("skipping the nested workflow step leaves %s, want mid/tail", got)
	}
}

const innerChecks = `
[step.rev]
kind = "checks"
on.pass = "sign"
[[step.rev.check]]
name = "alpha"
prompt = "check {check}"
[step.rev.critic]
prompt = "critic"
[step.rev.fixer]
prompt = "fix"
[step.rev.verifier]
prompt = "verify"
[step.sign]
kind = "signoff"
on.pass = "done"
on.changes = "rev"
`

func TestSignoffChangesIntoIncludedChecks(t *testing.T) {
	files := map[string]string{"inner": innerChecks}
	src := func(changes string) string {
		return `
[step.w]
kind = "workflow"
workflow = "inner"
skip = ["sign"]
on.pass = "gate"
[step.gate]
kind = "signoff"
on.pass = "done"
on.changes = "` + changes + `"
`
	}
	d, err := ParseWith("test", []byte(src("w/rev")), mapResolver(files))
	if err != nil {
		t.Fatal(err)
	}
	if next, _ := d.Next("gate", OutcomeChanges); next != "inner/rev" {
		t.Errorf("gate on changes goes to %q, want inner/rev", next)
	}
	for _, tt := range []struct{ changes, want string }{
		{"w/nope", `workflow "inner" has no step "nope" left to run`},
		{"w/sign", `workflow "inner" has no step "sign" left to run`},
		{"w", `not a checks step`},
		{"x/rev", `which is not a step`},
	} {
		_, err := ParseWith("test", []byte(src(tt.changes)), mapResolver(files))
		if err == nil || !strings.Contains(err.Error(), tt.want) {
			t.Errorf("on.changes %q: error %v, want one containing %q", tt.changes, err, tt.want)
		}
	}
}

// The render step of the build and the review workflow, run against a
// stub rhiza in a scratch git repo: a branch that is already rendered
// and has no uncommitted .rhiza/log entry gets no commit; one left stale
// gets exactly one.
func TestRenderStepAddsNoCommitWhenClean(t *testing.T) {
	for _, name := range []string{"build", "review"} {
		t.Run(name, func(t *testing.T) {
			data, ok := Builtin(name)
			if !ok {
				t.Fatalf("no %s built-in", name)
			}
			def, err := Parse(name, data)
			if err != nil {
				t.Fatal(err)
			}
			run := def.Steps["render"].Run
			if run == "" {
				t.Fatalf("%s has no render command", name)
			}
			dir := t.TempDir()
			bin := filepath.Join(dir, "bin")
			repo := filepath.Join(dir, "repo")
			for _, d := range []string{bin, repo} {
				if err := os.MkdirAll(d, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			// The stub's render rewrites the rendered page with the
			// content of $STUB_PAGE, as a real render does when stale.
			stub := "#!/bin/sh\n[ \"$1\" = render ] || exit 2\nmkdir -p .rhiza/rendered\nprintf '%s\\n' \"$STUB_PAGE\" > .rhiza/rendered/index.md\n"
			if err := os.WriteFile(filepath.Join(bin, "rhiza"), []byte(stub), 0o755); err != nil {
				t.Fatal(err)
			}
			env := append(os.Environ(),
				"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
				"RHIZA_CHANGE=01M61GW7K1PDG7HMFXMN3YNA7B",
				"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
				"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
				"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
			sh := func(script, page string) string {
				t.Helper()
				cmd := exec.Command("sh", "-c", script)
				cmd.Dir = repo
				cmd.Env = append(env, "STUB_PAGE="+page)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("%s: %v\n%s", script, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			sh("git init -q && mkdir -p .rhiza/log && echo x > .rhiza/log/a.json", "v1")
			sh(run, "v1")
			head := sh("git rev-parse HEAD", "")
			if n := sh("git rev-list --count HEAD", ""); n != "1" {
				t.Fatalf("stale render made %s commits; want 1", n)
			}

			sh(run, "v1")
			if got := sh("git rev-parse HEAD", ""); got != head {
				t.Errorf("render on a rendered branch moved HEAD %s to %s", head, got)
			}
			if st := sh("git status --porcelain -- .rhiza", ""); st != "" {
				t.Errorf("render left %q uncommitted", st)
			}
		})
	}
}

// TestBuiltinRolePromptsScopeTests: every builder, verifier, check, critic
// and fixer prompt of the build and review workflows tells the agent to
// run only the tests its change touches, and names no test command.
func TestBuiltinRolePromptsScopeTests(t *testing.T) {
	rule := "Run only the tests covering what your change touches, never the repo's full suite"
	var prompts = map[string]string{}
	for _, name := range []string{"build", "review"} {
		data, ok := Builtin(name)
		if !ok {
			t.Fatalf("no %s built-in", name)
		}
		def, err := Parse(name, data)
		if err != nil {
			t.Fatal(err)
		}
		for sn, st := range def.Steps {
			for rn, r := range map[string]*Role{"builder": st.Builder, "verifier": st.Verifier, "critic": st.Critic, "fixer": st.Fixer} {
				if r != nil {
					prompts[name+"/"+sn+"/"+rn] = r.Prompt
				}
			}
			for _, c := range st.Checks {
				prompts[name+"/"+sn+"/check "+c.Name] = c.Prompt
			}
		}
	}
	if len(prompts) < 10 {
		t.Fatalf("found %d role prompts; want builder, verifiers, 5 checks, critic, fixer", len(prompts))
	}
	for where, p := range prompts {
		if !strings.Contains(p, rule) {
			t.Errorf("%s prompt lacks the rule: %q", where, p)
		}
		for _, bad := range []string{"go test", "make check", "make test", "testscript"} {
			if strings.Contains(p, bad) {
				t.Errorf("%s prompt names %q", where, bad)
			}
		}
	}
}

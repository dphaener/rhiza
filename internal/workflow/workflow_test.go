package workflow

import (
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
		{"unknown-kind.toml", []string{`step "repeat": unknown kind "loop": want agent, wait, command, arbiter or leaves`}},
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
			`step "x": unknown kind "loop": want agent, wait, command, arbiter or leaves`},
		{"workflow kind is later", "[step.x]\nkind = \"workflow\"\n",
			`step "x": unknown kind "workflow"`},
		{"no kind", "[step.x]\nrun = \"true\"\n",
			`step "x": no kind: want agent, wait, command, arbiter or leaves`},
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
		def, err := Parse(n, data)
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

func TestShipChangeBuiltin(t *testing.T) {
	data, ok := Builtin("ship-change")
	if !ok {
		t.Fatal("no ship-change built-in")
	}
	def, err := Parse("ship-change", data)
	if err != nil {
		t.Fatal(err)
	}
	build := def.Steps["build"]
	for _, want := range []string{"shipit", "rhiza slice --change {change}"} {
		if !strings.Contains(build.Prompt, want) {
			t.Errorf("build prompt lacks %q", want)
		}
	}
	if !build.WaitsForDone() {
		t.Errorf("build until_when = %q; want %q", build.UntilWhen, UntilWhenDone)
	}
	if !strings.Contains(build.Prompt, "rhiza session done") {
		t.Error("build prompt does not tell the agent to run rhiza session done")
	}
	want := []string{"build", "wait-pr", "pr-green", "rebase", "wait-rebased", "rebased-green", "land", "wait-landed", "landed-green",
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
	want := []string{"interview", "review", "critic", "check", "revise", "gate"}
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
		{"interview", OutcomePass, "review", EndNone},
		{"interview", OutcomeFail, "", EndFailed},
		{"review", OutcomePass, "critic", EndNone},
		{"critic", OutcomePass, "check", EndNone},
		{"check", OutcomePass, "gate", EndNone},
		{"check", OutcomeFail, "revise", EndNone},
		{"revise", OutcomePass, "review", EndNone},
		{"revise", OutcomeFail, "", EndFailed},
		{"gate", OutcomePass, "", EndSucceeded},
		{"gate", OutcomeFail, "", EndFailed},
	}
	for _, tt := range trans {
		next, end := def.Next(tt.step, tt.outcome)
		if next != tt.next || end != tt.end {
			t.Errorf("Next(%q, %q) = (%q, %v); want (%q, %v)", tt.step, tt.outcome, next, end, tt.next, tt.end)
		}
	}

	caps := map[string]int{"interview": 1, "review": 3, "critic": 3, "check": 3, "revise": 2, "gate": 1}
	for name, c := range caps {
		if got := def.Steps[name].Cap(); got != c {
			t.Errorf("%s cap = %d; want %d", name, got, c)
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
			`step "a": split is only for leaves steps; an agent step cannot have it`},
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

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
		{"unknown-kind.toml", []string{`step "repeat": unknown kind "loop": want agent, wait or command`}},
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
			`step "x": unknown kind "loop": want agent, wait or command`},
		{"workflow kind is later", "[step.x]\nkind = \"workflow\"\n",
			`step "x": unknown kind "workflow"`},
		{"no kind", "[step.x]\nrun = \"true\"\n",
			`step "x": no kind: want agent, wait or command`},
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
			`step "a": agent is only for agent steps`},
		{"permission_mode on wait", "[step.w]\nkind = \"wait\"\nuntil = \"true\"\ntimeout = \"1m\"\npermission_mode = \"manual\"\n",
			`step "w": permission_mode is only for agent steps`},
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

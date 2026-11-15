package api

import (
	"errors"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/workflow"
)

// projectAgentCtx is a test context whose git config holds
// rhiza.project.agent = agent (unset when empty) and records what git
// config writes.
func projectAgentCtx(t *testing.T, agent string) (*Ctx, map[string]string) {
	t.Helper()
	conf := map[string]string{}
	if agent != "" {
		conf[ProjectAgentKey] = agent
	}
	home := t.TempDir()
	ctx := testCtx(t.TempDir())
	ctx.Home = home
	ctx.Getenv = func(k string) string {
		if k == "HOME" {
			return home
		}
		return ""
	}
	ctx.Git = func(args ...string) (string, error) {
		switch {
		case len(args) == 3 && args[0] == "config" && args[1] == "--get":
			if v, ok := conf[args[2]]; ok {
				return v + "\n", nil
			}
			return "", errUnset
		case len(args) == 3 && args[0] == "config":
			conf[args[1]] = args[2]
			return "", nil
		}
		return "", nil
	}
	return ctx, conf
}

var errUnset = errors.New("unset")

func TestProjectAgentStart(t *testing.T) {
	t.Parallel()
	def := adapter.Default().Name()
	cases := []struct {
		name, project, flag, step, want string
	}{
		{"Unset", "", "", "", def},
		{"Set", askAgent, "", "", askAgent},
		{"ExplicitWins", askAgent, def, "", def},
		{"StepWins", askAgent, "", def, def},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, _ := projectAgentCtx(t, c.project)
			const mode = "manual"
			// session start
			if c.step == "" {
				a, err := startAdapter(ctx, SessionStartRequest{Agent: c.flag, PermissionMode: mode})
				if err != nil {
					t.Fatalf("session start: %v", err)
				}
				if a.Name() != c.want {
					t.Errorf("session start agent = %s, want %s", a.Name(), c.want)
				}
			}
			// run start's agent steps: the step's adapter, else the run's
			// --agent, else the project's.
			req := agentStepRequest(workflow.Step{Name: "build", Agent: c.step}, mode, c.flag)
			a, err := startAdapter(ctx, req)
			if err != nil {
				t.Fatalf("run step: %v", err)
			}
			if a.Name() != c.want {
				t.Errorf("run step agent = %s, want %s", a.Name(), c.want)
			}
			if err := checkAgentStep(ctx, workflow.Step{Name: "build", Agent: c.step}, "an agent step", mode, c.flag, "next"); err != nil {
				t.Errorf("checkAgentStep: %v", err)
			}
		})
	}
}

// TestProjectAgentSkillCheck: skill check takes the project's agent. The
// ask-fake adapter reads no slash commands, so with it the check passes
// the skill check and stops later; the default adapter refuses the
// uninstalled skill by its name.
func TestProjectAgentSkillCheck(t *testing.T) {
	t.Parallel()
	def := adapter.Default().Name()
	missing := "skill not installed"
	cases := []struct {
		name, project, flag, want string
	}{
		{"Unset", "", "", missing},
		{"Set", askAgent, "", "this caller runs no programs"},
		{"ExplicitWins", askAgent, def, missing},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, _ := projectAgentCtx(t, c.project)
			req := `{"invoke":"/rhiza-builder Build {leaf}","permission_mode":"manual"`
			if c.flag != "" {
				req += `,"agent":"` + c.flag + `"`
			}
			_, err := callOp(t, ctx, "skill.check", req+"}")
			if err == nil || !strings.Contains(AsError(err).Message, c.want) {
				t.Fatalf("err = %v, want %q", err, c.want)
			}
		})
	}
}

func TestProjectAgentDaemonSet(t *testing.T) {
	t.Parallel()
	ctx, conf := projectAgentCtx(t, "")
	_, err := callOp(t, ctx, "daemon.set", `{"key":"rhiza.project.agent","value":"nope"}`)
	if err == nil {
		t.Fatal("daemon set nope: no error")
	}
	msg := AsError(err).Message
	if AsError(err).Code != CodeUsage || !strings.Contains(msg, "want one of "+strings.Join(adapter.Names(), ", ")) {
		t.Fatalf("daemon set nope: %v; want a usage error listing %v", err, adapter.Names())
	}
	if _, ok := conf[ProjectAgentKey]; ok {
		t.Fatalf("a refused value was written: %v", conf)
	}
	resp, err := callOp(t, ctx, "daemon.set", `{"key":"rhiza.project.agent","value":"`+askAgent+`"}`)
	if err != nil {
		t.Fatalf("daemon set %s: %v", askAgent, err)
	}
	if r := resp.(DaemonSetResponse); r.Scope != ScopeProject || conf[ProjectAgentKey] != askAgent {
		t.Fatalf("daemon set = %+v, config %v", r, conf)
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, s := range rows {
		if s.Key == ProjectAgentKey {
			found = true
			if s.Value != askAgent || s.Source != SettingConfig || s.Scope != ScopeProject || s.Default != "" {
				t.Errorf("settings row = %+v", s)
			}
		}
	}
	if !found {
		t.Errorf("daemon settings lists no %s", ProjectAgentKey)
	}
}

// TestProjectAgentUnknown: a name in git config the registry does not
// know refuses the start, naming the key and the registered names.
func TestProjectAgentUnknown(t *testing.T) {
	t.Parallel()
	ctx, _ := projectAgentCtx(t, "nope")
	_, err := startAdapter(ctx, SessionStartRequest{PermissionMode: "manual"})
	if err == nil || !strings.Contains(AsError(err).Message, `git config rhiza.project.agent is "nope"`) ||
		!strings.Contains(AsError(err).Message, strings.Join(adapter.Names(), ", ")) {
		t.Fatalf("err = %v", err)
	}
	// An explicit --agent does not read the setting.
	if _, err := startAdapter(ctx, SessionStartRequest{Agent: adapter.Default().Name(), PermissionMode: "manual"}); err != nil {
		t.Fatalf("explicit agent: %v", err)
	}
}

package api

// This file works out which workflows the daemon would run from a queued
// request with no person in between, and which of their agent sessions run
// in a bypass mode: the workflows outside text (a request's text and its
// --ref) reaches with no permission checks. rhiza health and rhiza workflow
// check warn from it; nothing here refuses anything.

import (
	"fmt"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/workflow"
)

// The triggers of an exposed workflow.
const (
	// TriggerRequest is the workflow the daemon starts from every queued
	// request.
	TriggerRequest = "request"
	// TriggerReady is the workflow the daemon starts on a change set
	// that is ready to ship.
	TriggerReady = "ready"
)

// ExposedWorkflow is a workflow outside text reaches with no permission
// checks.
type ExposedWorkflow struct {
	Workflow string   `json:"workflow" desc:"The workflow's name"`
	Trigger  string   `json:"trigger" enum:"request,ready" desc:"What starts it: request (a queued request) or ready (a ready change set)"`
	Steps    []string `json:"steps" desc:"The steps and roles whose sessions run in a bypass mode, in step order, such as \"interview\" or \"build/build builder\""`
}

// requestExposedWorkflows lists the workflows the daemon would run from a
// queued request with no person in between, each with the steps and roles
// whose effective mode its adapter reports as a bypass mode. settings are
// the daemon's settings (settingsOf). The request workflow (plan) counts
// whatever rhiza.autonomy.plan says: its interview agent starts at once on
// the request's text. The ready workflow (ship-change) counts only when
// rhiza.autonomy.plan, accept and ship are all auto, as only then does
// request text reach its agents unseen. A workflow that cannot be loaded
// is left out: rhiza workflow check reports it.
func requestExposedWorkflows(ctx *Ctx, settings []Setting) []ExposedWorkflow {
	mode := settingValue(settings, "rhiza.daemon.permissionMode")
	picks := []struct{ name, trigger string }{{planWorkflow, TriggerRequest}}
	if settingValue(settings, "rhiza.autonomy.plan") == AutonomyAuto &&
		settingValue(settings, "rhiza.autonomy.accept") == AutonomyAuto &&
		settingValue(settings, "rhiza.autonomy.ship") == AutonomyAuto {
		picks = append(picks, struct{ name, trigger string }{shipChangeWorkflow, TriggerReady})
	}
	var out []ExposedWorkflow
	for _, p := range picks {
		_, def, _, err := resolveWorkflow(ctx.Root, p.name)
		if err != nil {
			continue
		}
		if steps := bypassSteps(def, mode); len(steps) > 0 {
			out = append(out, ExposedWorkflow{Workflow: p.name, Trigger: p.trigger, Steps: steps})
		}
	}
	return out
}

// bypassSteps names every session def starts, in step order, whose
// effective permission mode (its own, else the run's mode) is a bypass mode
// for its adapter (its own, else the default). The arbiter counts when the
// workflow has an agent session, as an autonomous run launches it then.
func bypassSteps(def *workflow.Definition, mode string) []string {
	var out []string
	agents := false
	check := func(name, agentName, stepMode string) {
		agents = true
		if stepMode == "" {
			stepMode = mode
		}
		a := adapter.Default()
		if agentName != "" {
			a, _ = adapter.Lookup(agentName)
		}
		if a != nil && stepMode != "" && a.Bypass(stepMode) {
			out = append(out, name)
		}
	}
	role := func(step, label string, r *workflow.Role) {
		if r != nil {
			check(fmt.Sprintf("%s %s", step, label), r.Agent, r.PermissionMode)
		}
	}
	for _, n := range def.Order {
		s := def.Steps[n]
		switch s.Kind {
		case workflow.KindAgent:
			check(n, s.Agent, s.PermissionMode)
		case workflow.KindLeaves:
			role(n, workflow.RoleBuilder, s.Builder)
			role(n, workflow.RoleVerifier, s.Verifier)
		case workflow.KindChecks:
			for _, c := range s.Checks {
				check(fmt.Sprintf("%s check %s", n, c.Name), c.Agent, c.PermissionMode)
			}
			role(n, workflow.RoleCritic, s.Critic)
			role(n, workflow.RoleFixer, s.Fixer)
			role(n, workflow.RoleVerifier, s.Verifier)
		}
	}
	if agents {
		a := def.Arbiter()
		check(a.Name, a.Agent, a.PermissionMode)
	}
	return out
}

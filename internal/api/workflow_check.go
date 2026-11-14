package api

// This file holds workflow.check: a workflow parsed and validated as run
// start would, includes expanded, every problem reported with its step. It
// starts nothing and writes nothing.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "workflow.check",
		Summary:  "Validate a workflow as run start would, with includes expanded, and report every problem with the step it is in; starts nothing",
		Request:  reflect.TypeFor[WorkflowCheckRequest](),
		Response: reflect.TypeFor[WorkflowCheckResponse](),
		Run:      runWorkflowCheck,
	})
}

// WorkflowCheckRequest names the workflow.
type WorkflowCheckRequest struct {
	Workflow string `json:"workflow" arg:"0" desc:"A workflow name (personal, repo, then built-in) or a .toml file path"`
}

// WorkflowProblem is one thing wrong with a workflow.
type WorkflowProblem struct {
	Step    string `json:"step" desc:"The step the problem is in; empty for a problem of the whole file"`
	Message string `json:"message" desc:"What is wrong and what is wanted"`
}

// WorkflowCheckResponse is the verdict.
type WorkflowCheckResponse struct {
	Name     string            `json:"name" desc:"The workflow's name"`
	Path     string            `json:"path" desc:"The file it resolved from, or builtin:<name>"`
	Valid    bool              `json:"valid" desc:"True when run start would accept the workflow"`
	Problems []WorkflowProblem `json:"problems" desc:"Every problem found; empty when valid"`
	Warnings []string          `json:"warnings,omitempty" desc:"Risks worth knowing about that do not make the workflow invalid, such as outside text reaching its agents with no permission checks; never change valid or the exit code"`
}

// Human writes the verdict as text.
func (r WorkflowCheckResponse) Human(w io.Writer) {
	defer func() {
		for _, warn := range r.Warnings {
			fmt.Fprintf(w, "  warning: %s\n", warn)
		}
	}()
	if r.Valid {
		fmt.Fprintf(w, "workflow %s is valid (%s)\n", r.Name, r.Path)
		return
	}
	fmt.Fprintf(w, "workflow %s has %d problem(s) (%s)\n", r.Name, len(r.Problems), r.Path)
	for _, p := range r.Problems {
		if p.Step != "" {
			fmt.Fprintf(w, "  step %s: %s\n", p.Step, p.Message)
		} else {
			fmt.Fprintf(w, "  %s\n", p.Message)
		}
	}
}

// ExitCode is 1 when the workflow has a problem, else 0.
func (r WorkflowCheckResponse) ExitCode() int {
	if r.Valid {
		return 0
	}
	return 1
}

func runWorkflowCheck(ctx *Ctx, req any) (any, error) {
	r := req.(WorkflowCheckRequest)
	ref := strings.TrimSpace(r.Workflow)
	if ref == "" {
		ref = defaultWorkflow
	}
	var (
		path string
		data []byte
		name string
	)
	if isWorkflowPath(ref) {
		p, err := workflowFile(ctx.Root, ref)
		if err != nil {
			return nil, err
		}
		path = p
		if data, err = os.ReadFile(path); err != nil {
			return nil, Usagef("no workflow %q: %v", ref, err).WithNext(runStartNext)
		}
		name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	} else {
		p, d, ok, err := lookupWorkflow(ctx.Root, ref)
		if err != nil {
			return nil, Usagef("workflow %q: %v", ref, err).WithNext(runStartNext)
		}
		if !ok {
			return nil, Usagef("no workflow %q: there is no such file in %s or your personal workflows, and no built-in of that name (built-in: %s)",
				ref, workflowsDir, strings.Join(workflow.BuiltinNames(), ", ")).WithNext(runStartNext)
		}
		path, data, name = p, d, ref
	}
	out := WorkflowCheckResponse{Name: name, Path: path, Problems: []WorkflowProblem{}}
	name, problems := workflow.Validate(name, data, workflowResolver(ctx.Root))
	out.Name = name
	for _, p := range problems {
		out.Problems = append(out.Problems, WorkflowProblem{Step: p.Step, Message: p.Message})
	}
	if len(out.Problems) == 0 {
		def, err := workflow.ParseWith(name, data, workflowResolver(ctx.Root))
		if err != nil {
			return nil, err
		}
		out.Problems = append(out.Problems, agentProblems(ctx, def)...)
	}
	out.Valid = len(out.Problems) == 0
	if !isWorkflowPath(ref) {
		out.Warnings = exposureWarnings(ctx, name)
	}
	out.Warnings = append(out.Warnings, placeholderWarnings(name, data, ctx)...)
	out.Warnings = append(out.Warnings, invokeWarnings(name, data, ctx)...)
	out.Warnings = append(out.Warnings, loopWarnings(name, data, ctx)...)
	out.Warnings = append(out.Warnings, mergeGateWarnings(name, data, ctx)...)
	return out, nil
}

// mergeGateWarnings warns for each merge step (workflow.IsMergeStep) the
// run can reach from its start step without entering a gate step
// (workflow.UngatedMerges), naming the step and one such path. It never
// makes the workflow invalid.
func mergeGateWarnings(name string, data []byte, ctx *Ctx) []string {
	def, err := workflow.ParseWith(name, data, workflowResolver(ctx.Root))
	if err != nil {
		return nil
	}
	var out []string
	for _, m := range def.UngatedMerges() {
		out = append(out, fmt.Sprintf("step %s merges with no gate step before it: the run can reach it along %s; route every path to it through a gate step",
			m.Step, strings.Join(m.Path, " -> ")))
	}
	return out
}

// loopWarnings warns for each plain red loop (workflow.PlainRedLoops): a
// red outcome naming by plain step name a step that leads back through an
// agent step, which handback_cap does not count. It names the step, the
// outcome and its target and suggests handback:<agent step>. It never
// makes the workflow invalid.
func loopWarnings(name string, data []byte, ctx *Ctx) []string {
	def, err := workflow.ParseWith(name, data, workflowResolver(ctx.Root))
	if err != nil {
		return nil
	}
	var out []string
	for _, l := range def.PlainRedLoops() {
		out = append(out, fmt.Sprintf("step %s: on.%s = %q loops back through agent step %s uncapped: handback_cap does not count it and no handback.md reaches the agent; use on.%s = \"%s%s\"",
			l.Step, l.Outcome, l.Target, l.Agent, l.Outcome, workflow.HandbackPrefix, l.Agent))
	}
	return out
}

// placeholderWarnings warns for each command step still running the
// skeleton's placeholder test command. It never makes the workflow invalid.
func placeholderWarnings(name string, data []byte, ctx *Ctx) []string {
	def, err := workflow.ParseWith(name, data, workflowResolver(ctx.Root))
	if err != nil {
		return nil
	}
	var out []string
	for _, n := range def.Order {
		s := def.Steps[n]
		if s.Kind == workflow.KindCommand && strings.Contains(s.Run, placeholderMarker) {
			out = append(out, fmt.Sprintf("step %s still runs the placeholder command from workflow new: replace it with your project's test command", s.Name))
		}
	}
	return out
}

// invokeWarnings warns for each agent step and role whose invoke starts
// with a slash command that its agent's adapter (the default for none)
// finds installed neither in the agent's config folder nor in the repo,
// naming the step, the role and the command, and for a Rhiza skill the
// rhiza guide install command that installs it; when Rhiza's own skills
// cannot be listed it warns so and names no install command. It never makes the
// workflow invalid; an agent no adapter knows is agentProblems' to report.
func invokeWarnings(name string, data []byte, ctx *Ctx) []string {
	def, err := workflow.ParseWith(name, data, workflowResolver(ctx.Root))
	if err != nil {
		return nil
	}
	var out []string
	files, err := pluginFiles(ctx.Version)
	if err != nil {
		// Without them no warning can tell a Rhiza skill apart, so none
		// names the install command.
		out = append(out, fmt.Sprintf("cannot list Rhiza's skills (%v): no invoke warning names rhiza guide install", err))
	}
	warn := func(step, role, agent, prompt string) {
		a := adapter.Default()
		if agent != "" {
			var ok bool
			if a, ok = adapter.Lookup(agent); !ok {
				return
			}
		}
		skill, guideAgent, missing := missingStartSkill(ctx, a, strings.TrimSpace(prompt), ctx.Root)
		if !missing {
			return
		}
		msg := fmt.Sprintf("step %s: the %s's invoke starts with /%s, which %s has not installed", step, role, skill, guideAgent)
		if isRhizaSkill(files, skill) {
			msg += "; run rhiza guide install --agent " + guideAgent
		}
		out = append(out, msg)
	}
	for _, n := range def.Order {
		s := def.Steps[n]
		switch s.Kind {
		case workflow.KindAgent:
			warn(s.Name, "agent step", s.Agent, s.Prompt)
		case workflow.KindArbiter:
			warn(s.Name, "arbiter step", s.Agent, s.Prompt)
		}
		for _, r := range []string{workflow.RoleBuilder, workflow.RoleVerifier, workflow.RoleCritic, workflow.RoleFixer} {
			if role := s.Role(r); role != nil {
				warn(s.Name, r, role.Agent, role.Prompt)
			}
		}
		for _, c := range s.Checks {
			warn(s.Name, "check "+c.Name, c.Agent, c.Prompt)
		}
	}
	return out
}

// isRhizaSkill says the slash command name is one of the skills rhiza
// guide install installs, plugin-qualified ("rhiza:<skill>") or not.
func isRhizaSkill(files map[string][]byte, name string) bool {
	name = strings.TrimPrefix(name, pluginName+":")
	if name == pluginName {
		return true
	}
	_, ok := files["skills/"+name+"/SKILL.md"]
	return ok
}

// exposureWarnings is the warning for the workflow called name when the
// request-exposed list flags it, else none. Settings that cannot be read
// give none: rhiza health reports them.
func exposureWarnings(ctx *Ctx, name string) []string {
	settings, err := settingsOf(ctx, nil)
	if err != nil {
		return nil
	}
	for _, e := range requestExposedWorkflows(ctx, settings) {
		if e.Workflow == name {
			return []string{ExposureDetail(e) + ": outside text reaches agents with no permission checks"}
		}
	}
	return nil
}

// agentProblems runs run start's adapter and permission mode checks
// (checkAgentSteps) on every agent and permission_mode the file sets
// explicitly, reporting each one that fails. A key left unset takes the
// run's --agent or --permission-mode, which only run start knows.
func agentProblems(ctx *Ctx, def *workflow.Definition) []WorkflowProblem {
	var out []WorkflowProblem
	check := func(step, what, agent, mode string) {
		if agent == "" && mode == "" {
			return
		}
		if mode == "" {
			if _, ok := adapter.Lookup(agent); !ok {
				out = append(out, WorkflowProblem{Step: step, Message: fmt.Sprintf("%s: unknown agent %q: want one of %s",
					what, agent, strings.Join(adapter.Names(), ", "))})
			}
			return
		}
		if _, err := startAdapter(ctx, SessionStartRequest{Agent: agent, PermissionMode: mode, Step: step}); err != nil {
			out = append(out, WorkflowProblem{Step: step, Message: what + ": " + AsError(err).Message})
		}
	}
	role := func(s workflow.Step, name, what string) {
		if r := s.Role(name); r != nil {
			check(s.Name, what+" "+name, r.Agent, r.PermissionMode)
		}
	}
	for _, n := range def.Order {
		s := def.Steps[n]
		switch s.Kind {
		case workflow.KindAgent:
			check(s.Name, "the agent step", s.Agent, s.PermissionMode)
		case workflow.KindArbiter:
			check(s.Name, "the arbiter step", s.Agent, s.PermissionMode)
		case workflow.KindLeaves:
			for _, r := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
				role(s, r, "the leaves step's")
			}
		case workflow.KindChecks:
			for _, c := range s.Checks {
				check(s.Name, "check "+c.Name, c.Agent, c.PermissionMode)
			}
			for _, r := range []string{workflow.RoleCritic, workflow.RoleFixer, workflow.RoleVerifier} {
				role(s, r, "the checks step's")
			}
		}
	}
	return out
}

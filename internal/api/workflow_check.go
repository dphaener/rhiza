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
	return out, nil
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

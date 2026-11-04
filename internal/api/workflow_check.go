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
}

// Human writes the verdict as text.
func (r WorkflowCheckResponse) Human(w io.Writer) {
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
	out.Valid = len(out.Problems) == 0
	return out, nil
}

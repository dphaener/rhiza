package api

// This file holds workflow.describe: a workflow resolved as run start
// resolves it, shown without starting anything. It writes nothing.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "workflow.describe",
		Summary:  "Show a workflow as run start resolves it: where it comes from, the layers it shadows, its start step, each step with its targets and roles, and the workflows it includes; starts nothing",
		Request:  reflect.TypeFor[WorkflowDescribeRequest](),
		Response: reflect.TypeFor[WorkflowDescribeResponse](),
		Run:      runWorkflowDescribe,
	})
}

// WorkflowDescribeRequest names the workflow.
type WorkflowDescribeRequest struct {
	Workflow string `json:"workflow" arg:"0" desc:"A workflow name (personal, repo, then built-in) or a .toml file path"`
}

// WorkflowRoleView is one session kind of a step.
type WorkflowRoleView struct {
	Role   string   `json:"role" desc:"The role: agent, builder, verifier, critic, fixer or a check's name"`
	Models []string `json:"models" desc:"The model ladder by attempt or cycle; empty means the agent's default"`
	Prompt string   `json:"prompt" desc:"The role's first prompt template"`
}

// WorkflowStepView is one step in order.
type WorkflowStepView struct {
	Name  string             `json:"name" desc:"The step's name; <workflow>/<step> for a step of an included workflow"`
	Kind  string             `json:"kind" desc:"agent, wait, command, arbiter, leaves, checks or signoff"`
	On    map[string]string  `json:"on" desc:"The next step (or done or failed) by outcome"`
	Roles []WorkflowRoleView `json:"roles" desc:"The step's sessions, with their model ladders and prompts"`
}

// WorkflowDescribeResponse is the described workflow.
type WorkflowDescribeResponse struct {
	Name     string             `json:"name" desc:"The workflow's name"`
	Path     string             `json:"path" desc:"The file it resolved from, or builtin:<name>"`
	Layer    string             `json:"layer" desc:"personal, repo, builtin or file"`
	Shadows  []WorkflowLayer    `json:"shadows" desc:"The lower layers the winner hides, highest first; empty when none"`
	Start    string             `json:"start" desc:"The step a run enters first"`
	Steps    []WorkflowStepView `json:"steps" desc:"Every step in order"`
	Includes []string           `json:"includes" desc:"The workflows it includes, directly or through another, sorted"`
}

// Human writes the workflow as text.
func (r WorkflowDescribeResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "%s (%s) %s\n", r.Name, r.Layer, r.Path)
	for _, s := range r.Shadows {
		fmt.Fprintf(w, "  shadows %s %s\n", s.Layer, s.Path)
	}
	fmt.Fprintf(w, "start %s\n", r.Start)
	if len(r.Includes) > 0 {
		fmt.Fprintf(w, "includes %s\n", strings.Join(r.Includes, ", "))
	}
	for _, s := range r.Steps {
		fmt.Fprintf(w, "%s %s\n", s.Name, s.Kind)
		keys := make([]string, 0, len(s.On))
		for k := range s.On {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  on.%s -> %s\n", k, s.On[k])
		}
		for _, ro := range s.Roles {
			fmt.Fprintf(w, "  role %s models [%s]\n", ro.Role, strings.Join(ro.Models, ", "))
			for _, line := range strings.Split(strings.TrimRight(ro.Prompt, "\n"), "\n") {
				fmt.Fprintf(w, "    %s\n", line)
			}
		}
	}
}

func roleView(name string, r workflow.Role) WorkflowRoleView {
	return WorkflowRoleView{Role: name, Models: append([]string{}, r.Models...), Prompt: r.Prompt}
}

func stepView(s workflow.Step) WorkflowStepView {
	v := WorkflowStepView{Name: s.Name, Kind: s.Kind, On: map[string]string{}, Roles: []WorkflowRoleView{}}
	for k, t := range s.On {
		v.On[k] = t
	}
	if s.Prompt != "" || len(s.Models) > 0 {
		v.Roles = append(v.Roles, roleView("agent", workflow.Role{Prompt: s.Prompt, Models: s.Models}))
	}
	for _, r := range []struct {
		name string
		role *workflow.Role
	}{{"builder", s.Builder}, {"verifier", s.Verifier}, {"critic", s.Critic}, {"fixer", s.Fixer}} {
		if r.role != nil {
			v.Roles = append(v.Roles, roleView(r.name, *r.role))
		}
	}
	for _, c := range s.Checks {
		v.Roles = append(v.Roles, roleView(c.Name, c.Role))
	}
	return v
}

// lowerLayers lists the layers of a bare name below the one that won.
func lowerLayers(root, name, winner string) []WorkflowLayer {
	out := []WorkflowLayer{}
	passed := false
	check := func(layer, p string) {
		if p == winner {
			passed = true
			return
		}
		if !passed {
			return
		}
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			out = append(out, WorkflowLayer{Layer: layer, Path: p})
		}
	}
	if d := personalWorkflowsDir(); d != "" {
		check(layerPersonal, filepath.Join(d, name+".toml"))
	}
	check(layerRepo, filepath.Join(root, filepath.FromSlash(workflowsDir), name+".toml"))
	if _, ok := workflow.Builtin(name); ok && winner != builtinPrefix+name {
		out = append(out, WorkflowLayer{Layer: layerBuiltin, Path: builtinPrefix + name})
	}
	return out
}

func runWorkflowDescribe(ctx *Ctx, req any) (any, error) {
	r := req.(WorkflowDescribeRequest)
	ref := strings.TrimSpace(r.Workflow)
	path, def, _, err := resolveWorkflow(ctx.Root, ref)
	if err != nil {
		return nil, err
	}
	out := WorkflowDescribeResponse{Name: def.Name, Path: path, Start: def.Start, Shadows: []WorkflowLayer{}, Steps: []WorkflowStepView{}, Includes: []string{}}
	switch {
	case isWorkflowPath(ref):
		out.Layer = "file"
	case strings.HasPrefix(path, builtinPrefix):
		out.Layer = layerBuiltin
	case path == filepath.Join(ctx.Root, filepath.FromSlash(workflowsDir), def.Name+".toml"):
		out.Layer = layerRepo
	default:
		out.Layer = layerPersonal
	}
	if out.Layer != "file" {
		out.Shadows = lowerLayers(ctx.Root, def.Name, path)
	}
	for _, n := range def.Order {
		out.Steps = append(out.Steps, stepView(def.Steps[n]))
	}
	for n := range def.Includes {
		out.Includes = append(out.Includes, n)
	}
	sort.Strings(out.Includes)
	return out, nil
}

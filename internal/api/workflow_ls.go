package api

// This file holds workflow.ls: every workflow a bare name can reach, each
// with the layer and path that wins and the lower layers it shadows. It
// reads the workflow folders only and never writes.

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
		Name:     "workflow.ls",
		Summary:  "List every workflow a bare name can reach (personal, repo, built-in), each with the layer and path it comes from and the layers it shadows",
		Request:  reflect.TypeFor[WorkflowLsRequest](),
		Response: reflect.TypeFor[WorkflowLsResponse](),
		Run:      runWorkflowLs,
	})
}

// Workflow layers, highest first.
const (
	layerPersonal = "personal"
	layerRepo     = "repo"
	layerBuiltin  = "builtin"
)

// WorkflowLsRequest takes no arguments.
type WorkflowLsRequest struct{}

// WorkflowLayer is one place a workflow name is found.
type WorkflowLayer struct {
	Layer string `json:"layer" desc:"personal, repo or builtin"`
	Path  string `json:"path" desc:"The file, or builtin:<name> for a built-in"`
}

// WorkflowLsItem is one workflow name.
type WorkflowLsItem struct {
	Name    string          `json:"name" desc:"The workflow's bare name"`
	Layer   string          `json:"layer" desc:"The layer a bare name resolves to: personal, repo or builtin"`
	Path    string          `json:"path" desc:"The winning file, or builtin:<name>"`
	Shadows []WorkflowLayer `json:"shadows" desc:"The lower layers the winner hides, highest first; empty when none"`
}

// WorkflowLsResponse is the listing.
type WorkflowLsResponse struct {
	Workflows []WorkflowLsItem `json:"workflows" desc:"Every reachable workflow, sorted by name"`
}

// Human writes the listing as text.
func (r WorkflowLsResponse) Human(w io.Writer) {
	if len(r.Workflows) == 0 {
		fmt.Fprintln(w, "No workflows")
		return
	}
	for _, it := range r.Workflows {
		fmt.Fprintf(w, "%s %s %s\n", it.Name, it.Layer, it.Path)
		for _, s := range it.Shadows {
			fmt.Fprintf(w, "  shadows %s %s\n", s.Layer, s.Path)
		}
	}
}

// workflowNames lists the names of the .toml files in dir, sorted; a
// missing dir has none.
func workflowNames(dir string) []string {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if n, ok := strings.CutSuffix(e.Name(), ".toml"); ok && n != "" {
			names = append(names, n)
		}
	}
	return names
}

// runWorkflowLs finds each name in the three layers lookupWorkflow reads.
func runWorkflowLs(ctx *Ctx, req any) (any, error) {
	personal := personalWorkflowsDir()
	repo := filepath.Join(ctx.Root, filepath.FromSlash(workflowsDir))
	layers := []struct{ name, dir string }{{layerPersonal, personal}, {layerRepo, repo}}
	seen := map[string]bool{}
	for _, l := range layers {
		for _, n := range workflowNames(l.dir) {
			seen[n] = true
		}
	}
	for _, n := range workflow.BuiltinNames() {
		seen[n] = true
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	out := []WorkflowLsItem{}
	for _, n := range names {
		var found []WorkflowLayer
		for _, l := range layers {
			if l.dir == "" {
				continue
			}
			p := filepath.Join(l.dir, n+".toml")
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				found = append(found, WorkflowLayer{Layer: l.name, Path: p})
			}
		}
		if _, ok := workflow.Builtin(n); ok {
			found = append(found, WorkflowLayer{Layer: layerBuiltin, Path: builtinPrefix + n})
		}
		out = append(out, WorkflowLsItem{
			Name: n, Layer: found[0].Layer, Path: found[0].Path, Shadows: append([]WorkflowLayer{}, found[1:]...),
		})
	}
	return WorkflowLsResponse{Workflows: out}, nil
}

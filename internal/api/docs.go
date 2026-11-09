package api

// This file holds the docs operation: it writes the generated reference
// from the operation registry.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/dphaener/rhiza/internal/docs"
	"github.com/dphaener/rhiza/internal/workflow"
)

// docsDir is the generated reference folder, relative to the project root.
const docsDir = "docs/reference"

// DocsRequest names where the reference goes.
type DocsRequest struct {
	Out string `json:"out,omitempty" desc:"The folder to write into; default docs/reference in the project root"`
}

// DocsResponse says what docs wrote.
type DocsResponse struct {
	Written   []string `json:"written" desc:"The files written: new, or different from what was on disk"`
	Unchanged int      `json:"unchanged" desc:"How many files were already current"`
}

func init() {
	register(Operation{
		Name:     "docs",
		Summary:  "Write the generated reference (CLI commands, settings, workflow file) under docs/reference",
		Request:  reflect.TypeFor[DocsRequest](),
		Response: reflect.TypeFor[DocsResponse](),
		Run:      runDocs,
	})
}

func runDocs(ctx *Ctx, req any) (any, error) {
	r := req.(DocsRequest)
	if ctx.DryRun {
		return nil, Usagef("docs does not take --dry-run: it writes no log entry").WithNext("rhiza docs")
	}
	dir := r.Out
	if dir == "" {
		dir = filepath.Join(ctx.Root, filepath.FromSlash(docsDir))
	} else if !filepath.IsAbs(dir) {
		dir = filepath.Join(ctx.Root, dir)
	}
	files := map[string][]byte{"cli.md": docs.CLI(docsCommands()), "config.md": docs.Config(docsSettings()), "workflow-file.md": docs.Workflow(docsWorkflow())}
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	resp := DocsResponse{Written: []string{}}
	for _, n := range names {
		dst := filepath.Join(dir, n)
		if have, err := os.ReadFile(dst); err == nil && string(have) == string(files[n]) {
			resp.Unchanged++
			continue
		}
		if err := writeAtomic(dst, files[n]); err != nil {
			return nil, Errorf("docs: %v", err)
		}
		shown := dst
		if rel, err := filepath.Rel(ctx.Root, dst); err == nil {
			shown = rel
		}
		resp.Written = append(resp.Written, filepath.ToSlash(shown))
	}
	return resp, nil
}

// docsCommands turns the registry into the renderer's commands.
func docsCommands() []docs.Command {
	exiter := reflect.TypeFor[Exiter]()
	var out []docs.Command
	for _, o := range Operations() {
		c := docs.Command{
			Name: CommandPath(o.Name), Summary: o.Summary, Usage: UsageLine(o),
			Writes: o.Write, SetsExit: o.Response.Implements(exiter),
		}
		if o.Input != nil {
			c.Input = o.Input.Summary
		}
		for _, p := range Params(o.Request) {
			name := "--" + p.Flag
			if p.Arg >= 0 {
				name = "<" + p.Name + ">"
			}
			var notes []string
			if p.Required && p.Kind != ParamBool {
				notes = append(notes, "required")
			}
			if p.Kind == ParamList && p.Arg < 0 {
				notes = append(notes, "repeatable")
			}
			if len(p.Enum) > 0 {
				notes = append(notes, "one of "+joinComma(p.Enum))
			}
			c.Params = append(c.Params, docs.Param{Name: name, Desc: p.Desc, Notes: notes})
		}
		c.Output = outputFields(o.Response)
		out = append(out, c)
	}
	return out
}

func joinComma(s []string) string {
	r := ""
	for i, v := range s {
		if i > 0 {
			r += ", "
		}
		r += v
	}
	return r
}

// outputFields lists a response type's top-level JSON fields, sorted.
func outputFields(t reflect.Type) []docs.Field {
	props, _ := Schema(t)["properties"].(map[string]any)
	var out []docs.Field
	for name, v := range props {
		m, _ := v.(map[string]any)
		typ, _ := m["type"].(string)
		if typ == "" {
			typ = "any"
		}
		if typ == "array" {
			if it, ok := m["items"].(map[string]any); ok {
				if s, _ := it["type"].(string); s != "" {
					typ = "array of " + s
				}
			}
		}
		desc, _ := m["description"].(string)
		out = append(out, docs.Field{Name: name, Type: typ, Desc: desc})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Human lists the files written.
func (r DocsResponse) Human(w io.Writer) {
	for _, p := range r.Written {
		fmt.Fprintf(w, "- wrote `%s`\n", p)
	}
	fmt.Fprintf(w, "%d unchanged\n", r.Unchanged)
}

// docsSettings turns the settings definitions into the renderer's settings.
func docsSettings() []docs.Setting {
	out := make([]docs.Setting, 0, len(settingDefs))
	for _, d := range settingDefs {
		out = append(out, docs.Setting{Key: d.key, Type: d.typ, Default: d.def, Scope: d.scope(), Desc: d.desc, Flag: d.flag})
	}
	return out
}

// docsWorkflow turns the workflow package's format description into the
// renderer's.
func docsWorkflow() docs.WorkflowFormat {
	conv := func(ks []workflow.KeyDoc) []docs.WorkflowKey {
		var out []docs.WorkflowKey
		for _, k := range ks {
			out = append(out, docs.WorkflowKey{Key: k.Key, Desc: k.Desc, Required: k.Required})
		}
		return out
	}
	f := docs.WorkflowFormat{File: conv(workflow.FileKeys), Shared: conv(workflow.SharedKeys), RoleKeys: conv(workflow.RoleKeys())}
	for _, k := range workflow.KindDocs() {
		f.Kinds = append(f.Kinds, docs.WorkflowKind{Kind: k.Kind, Summary: k.Summary, Keys: conv(k.Keys), Tables: conv(k.Tables), Outcomes: k.Outcomes})
	}
	for _, p := range workflow.PromptPlaceholders() {
		f.Placeholders = append(f.Placeholders, docs.Field{Name: p.Prompt, Desc: "{" + strings.Join(p.Names, "}, {") + "}"})
	}
	return f
}

package api

// This file holds rhiza project add and rhiza project ls: the projects a
// Rhiza root registers in its rhiza.toml (see internal/rootcfg).

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "project.add",
		Summary:  "Register an existing git checkout directly under the root as a project",
		Request:  reflect.TypeFor[ProjectAddRequest](),
		Response: reflect.TypeFor[ProjectAddResponse](),
		Run:      runProjectAdd,
	})
	register(Operation{
		Name:     "project.ls",
		ReadOnly: true,
		Summary:  "List the projects registered in the root",
		Request:  reflect.TypeFor[ProjectLsRequest](),
		Response: reflect.TypeFor[ProjectLsResponse](),
		Run:      runProjectLs,
	})
}

// ProjectAddRequest names the checkout.
type ProjectAddRequest struct {
	Path string `json:"path" arg:"0" desc:"The git checkout to register: a folder directly under the root"`
	Name string `json:"name,omitempty" flag:"name" desc:"The project's name (default: the folder's name)"`
}

// ProjectAddResponse reports the registered project.
type ProjectAddResponse struct {
	Name string `json:"name" desc:"The project's name"`
	Path string `json:"path" desc:"The project's folder, relative to the root"`
	Repo string `json:"repo,omitempty" desc:"The project's repo ID; empty when it has no .rhiza/"`
}

// Human writes the project as text.
func (r ProjectAddResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Registered project %s at %s.\nnext: rhiza project ls\n", r.Name, r.Path)
}

func runProjectAdd(ctx *Ctx, req any) (any, error) {
	r := req.(ProjectAddRequest)
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return nil, Usagef("%s is not inside a Rhiza root", ctx.Wd).WithNext("rhiza root init")
	}
	if r.Path == "" {
		return nil, Usagef("no path to register").WithNext("rhiza project add <path>")
	}
	abs := r.Path
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(ctx.Wd, abs)
	}
	abs = filepath.Clean(abs)
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		abs = real
	}
	if real, err := filepath.EvalSymlinks(root); err == nil {
		root = real
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil || rel == "." || rel == ".." || !filepath.IsLocal(rel) {
		return nil, Usagef("%s is outside the root %s", abs, root).
			WithNext("move the checkout directly under the root, or clone it there")
	}
	if rel != filepath.Base(rel) {
		return nil, Usagef("%s is nested under the root %s; a project sits directly under it", abs, root).
			WithNext("move the checkout directly under the root")
	}
	if _, err := os.Stat(filepath.Join(abs, ".git")); err != nil {
		return nil, Usagef("%s is not a git checkout", abs).WithNext("git init " + abs)
	}
	name := r.Name
	if name == "" {
		name = rel
	}
	if name == "" || name != filepath.Base(name) || name == "." || name == ".." {
		return nil, Usagef("%q is not a usable project name", name).WithNext("rhiza project add " + r.Path + " --name <name>")
	}
	cfg, err := rootcfg.Load(root)
	if err != nil {
		return nil, err
	}
	taken := Usagef("a project named %s is already registered", name).
		WithNext("rhiza project add " + r.Path + " --name <other name>")
	if _, ok := cfg.Projects[name]; ok {
		return nil, taken
	}
	p := rootcfg.Project{Path: rel}
	switch s, err := store.Open(abs); {
	case err == nil:
		p.Repo = s.Repo().ID.String()
	case errors.Is(err, store.ErrNotFound):
	default:
		return nil, err
	}
	if ctx.DryRun {
		return ProjectAddResponse{Name: name, Path: p.Path, Repo: p.Repo}, nil
	}
	if err := rootcfg.Add(root, name, p); err != nil {
		if errors.Is(err, rootcfg.ErrTaken) {
			return nil, taken
		}
		return nil, err
	}
	return ProjectAddResponse{Name: name, Path: p.Path, Repo: p.Repo}, nil
}

// ProjectLsRequest takes nothing.
type ProjectLsRequest struct{}

// ProjectInfo is one registered project.
type ProjectInfo struct {
	Name        string `json:"name" desc:"The project's name"`
	Path        string `json:"path" desc:"The project's folder, relative to the root"`
	Repo        string `json:"repo,omitempty" desc:"The project's repo ID; empty when it has no .rhiza/"`
	Initialised bool   `json:"initialised" desc:"Whether the project's folder holds .rhiza/"`
}

// ProjectLsResponse lists the projects.
type ProjectLsResponse struct {
	Projects []ProjectInfo `json:"projects" desc:"The registered projects, by name"`
}

// Human writes the projects as text.
func (r ProjectLsResponse) Human(w io.Writer) {
	if len(r.Projects) == 0 {
		fmt.Fprintln(w, "no projects\nnext: rhiza project add <path>")
		return
	}
	for _, p := range r.Projects {
		state := "initialised"
		if !p.Initialised {
			state = "not initialised"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.Path, state)
	}
}

func runProjectLs(ctx *Ctx, req any) (any, error) {
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return nil, Usagef("%s is not inside a Rhiza root", ctx.Wd).WithNext("rhiza root init")
	}
	cfg, err := rootcfg.Load(root)
	if err != nil {
		return nil, err
	}
	out := ProjectLsResponse{Projects: []ProjectInfo{}}
	for _, n := range cfg.Names() {
		p := cfg.Projects[n]
		_, serr := os.Stat(filepath.Join(cfg.Dir(root, n), store.Dir))
		out.Projects = append(out.Projects, ProjectInfo{Name: n, Path: p.Path, Repo: p.Repo, Initialised: serr == nil})
	}
	return out, nil
}

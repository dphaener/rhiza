package api

// This file holds rhiza root init and rhiza root show: the Rhiza root is a
// directory holding a rhiza.toml, found by walking up from the working
// directory (see internal/rootcfg).

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

func init() {
	register(Operation{
		Name:     "root.init",
		Summary:  "Make a directory a Rhiza root by writing its rhiza.toml",
		Request:  reflect.TypeFor[RootInitRequest](),
		Response: reflect.TypeFor[RootInitResponse](),
		Run:      runRootInit,
	})
	register(Operation{
		Name:     "root.show",
		ReadOnly: true,
		Summary:  "Show the Rhiza root, its config file and the project the working directory belongs to",
		Request:  reflect.TypeFor[RootShowRequest](),
		Response: reflect.TypeFor[RootShowResponse](),
		Run:      runRootShow,
	})
}

// RootInitRequest names the directory.
type RootInitRequest struct {
	Dir string `json:"dir,omitempty" arg:"0" desc:"The directory to make a root; created if missing (default: the working directory)"`
}

// RootInitResponse reports the new root.
type RootInitResponse struct {
	Root   string `json:"root" desc:"The root directory"`
	Config string `json:"config" desc:"The root's config file"`
}

// Human writes the root as text.
func (r RootInitResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Initialised a Rhiza root in %s.\nnext: rhiza root show\n", r.Root)
}

func runRootInit(ctx *Ctx, req any) (any, error) {
	r := req.(RootInitRequest)
	dir := r.Dir
	switch {
	case dir == "":
		dir = ctx.Wd
	case !filepath.IsAbs(dir) && ctx.Wd != "":
		dir = filepath.Join(ctx.Wd, dir)
	}
	if dir == "" {
		return nil, Usagef("no directory to make a root")
	}
	if ctx.DryRun {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		return RootInitResponse{Root: abs, Config: filepath.Join(abs, rootcfg.File)}, nil
	}
	path, err := rootcfg.Init(dir)
	if err == rootcfg.ErrExists {
		return nil, Usagef("%s is already a Rhiza root", filepath.Dir(path)).WithNext("rhiza root show")
	}
	if err != nil {
		return nil, err
	}
	return RootInitResponse{Root: filepath.Dir(path), Config: path}, nil
}

// RootShowRequest takes nothing.
type RootShowRequest struct{}

// RootShowResponse reports the root and project of the working directory.
type RootShowResponse struct {
	Root    string `json:"root" desc:"The root directory: the nearest ancestor holding a rhiza.toml"`
	Config  string `json:"config" desc:"The root's config file"`
	Project string `json:"project" desc:"The registered project the working directory belongs to, by its folder or a <project>.<slug> sibling worktree; empty at the root itself"`
}

// Human writes the root as text.
func (r RootShowResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "root: %s\nconfig: %s\n", r.Root, r.Config)
	if r.Project == "" {
		fmt.Fprintln(w, "project: none")
	} else {
		fmt.Fprintf(w, "project: %s\n", r.Project)
	}
}

func runRootShow(ctx *Ctx, req any) (any, error) {
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		return nil, Usagef("%s is not inside a Rhiza root", ctx.Wd).WithNext("rhiza root init")
	}
	cfg, err := rootcfg.Load(root)
	if err != nil {
		return nil, err
	}
	project, _ := cfg.ProjectOf(root, ctx.Wd)
	return RootShowResponse{Root: root, Config: filepath.Join(root, rootcfg.File), Project: project}, nil
}

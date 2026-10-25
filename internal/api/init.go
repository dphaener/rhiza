package api

// This file holds the init operation: it creates the repository, the
// .genesis/.gitignore that keeps local/ (the current change set) out of git,
// and .genesis/log/.gitkeep, since git keeps no empty folder and a clone of a
// fresh repository would otherwise have no log at all.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/store"
)

// gitignore is the content of .genesis/.gitignore: local/ holds per-worktree
// state that must never be committed.
const gitignore = localDir + "/\n"

// InitRequest takes nothing: init works in the caller's root.
type InitRequest struct{}

// InitResponse names the new repository.
type InitResponse struct {
	Repo id.ID  `json:"repo" desc:"The new repository's ID"`
	Root string `json:"root" desc:"The absolute path of the project root holding .genesis"`
}

func init() {
	register(Operation{
		Name:     "init",
		Summary:  "Create a Genesis repository in the current folder",
		Request:  reflect.TypeFor[InitRequest](),
		Response: reflect.TypeFor[InitResponse](),
		Run:      runInit,
	})
}

// runInit creates the repository in ctx.Root. An existing one is a
// conflict. Init writes no log entry, so --dry-run has nothing to preview
// and is a usage error rather than a silent no-op. The log starts empty and
// git keeps no empty folder, so init writes an empty log/.gitkeep for the
// log to survive a clone; the log reader skips dot-files.
func runInit(ctx *Ctx, _ any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("init does not take --dry-run: it writes no log entry").WithNext("genesis init")
	}
	root, err := filepath.Abs(ctx.Root)
	if err != nil {
		return nil, Errorf("project root: %v", err)
	}
	st, err := store.Init(root, ctx.IDs)
	if err != nil {
		if errors.Is(err, store.ErrExists) {
			return nil, Conflictf("a Genesis repository already exists in %s", root).WithNext("genesis change ls")
		}
		return nil, Errorf("init: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, store.Dir, ".gitignore"), []byte(gitignore), 0o644); err != nil {
		return nil, Errorf("write .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, store.Dir, store.LogDir, ".gitkeep"), nil, 0o644); err != nil {
		return nil, Errorf("write log/.gitkeep: %v", err)
	}
	return InitResponse{Repo: st.Repo().ID, Root: root}, nil
}

// Human reports the new repository and the first command to run.
func (r InitResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Initialised a Genesis repository in %s (repo %s).\n\n", r.Root, r.Repo)
	fmt.Fprintf(w, "next: %s\n", newEntryNext)
}

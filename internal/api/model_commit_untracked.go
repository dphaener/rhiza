package api

// rhiza model commit-untracked is the catch-up for model writes that never
// reached a model branch: it finds the untracked .rhiza/log entries of the
// main checkout that no pending model branch or live run holds (the ones
// rhiza check warns about) and commits them onto one model branch per change
// set, each with a model-only PR, as a main-checkout write does. The daemon
// runs it on start and on each tick.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

// ModelCommitUntrackedRequest takes nothing: it works in the caller's root.
type ModelCommitUntrackedRequest struct{}

// ModelBranchShipped is one model branch's worth of entries.
type ModelBranchShipped struct {
	ChangeSet string   `json:"change_set" desc:"The change set the entries belong to"`
	Entries   []string `json:"entries" desc:"The entry files, relative to the project root, committed onto the change set's model branch"`
	Error     string   `json:"error,omitempty" desc:"Why the entries are still untracked, when shipping them failed"`
}

// ModelCommitUntrackedResponse lists what was shipped, per change set.
type ModelCommitUntrackedResponse struct {
	Shipped []ModelBranchShipped `json:"shipped" desc:"One row per change set with orphan entries"`
}

// Human prints one line per change set.
func (r ModelCommitUntrackedResponse) Human(w io.Writer) {
	if len(r.Shipped) == 0 {
		fmt.Fprintln(w, "no untracked log entries to commit")
		return
	}
	for _, s := range r.Shipped {
		if s.Error != "" {
			fmt.Fprintf(w, "change set %s: %d entries not committed: %s\n", s.ChangeSet, len(s.Entries), s.Error)
			continue
		}
		fmt.Fprintf(w, "change set %s: committed %d entries on a model branch\n", s.ChangeSet, len(s.Entries))
	}
}

func init() {
	register(Operation{
		Name:     "model.commit-untracked",
		Summary:  "Commit the untracked .rhiza/log entries of a main checkout that no model branch or live run holds onto model branches, one per change set, each with a model-only PR",
		Request:  reflect.TypeFor[ModelCommitUntrackedRequest](),
		Response: reflect.TypeFor[ModelCommitUntrackedResponse](),
		Run:      runModelCommitUntracked,
	})
}

func runModelCommitUntracked(ctx *Ctx, _ any) (any, error) {
	resp := ModelCommitUntrackedResponse{Shipped: []ModelBranchShipped{}}
	orphans := orphanLogEntries(ctx)
	if len(orphans) == 0 {
		return resp, nil
	}
	top, _, ok := gitDirsOf(ctx.Root)
	if !ok {
		return resp, nil
	}
	by := map[string][]string{}
	for _, f := range orphans {
		if parts := strings.Split(f, "/"); len(parts) == 4 {
			by[parts[2]] = append(by[parts[2]], f)
		}
	}
	var sets []string
	for cs := range by {
		sets = append(sets, cs)
	}
	slices.Sort(sets)
	var failed []string
	for _, cs := range sets {
		files := by[cs] // orphans are sorted
		row := ModelBranchShipped{ChangeSet: cs, Entries: files}
		if !ctx.DryRun {
			if ctx.Exec == nil {
				return nil, Errorf("committing entries needs git and gh, and no programs can be run here")
			}
			tag := strings.TrimSuffix(filepath.Base(files[0]), ".json")
			err := shipModelWrite(ctx, top, cs, tag, func(root string) error {
				for _, f := range files {
					dst := filepath.Join(root, filepath.FromSlash(f))
					if _, err := os.Stat(dst); err == nil {
						continue
					}
					b, err := os.ReadFile(filepath.Join(ctx.Root, filepath.FromSlash(f)))
					if err != nil {
						return err
					}
					if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
						return err
					}
					if err := os.WriteFile(dst, b, 0o644); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				row.Error = err.Error()
				failed = append(failed, cs+": "+err.Error())
			}
		}
		resp.Shipped = append(resp.Shipped, row)
	}
	if len(failed) > 0 {
		return resp, Errorf("could not commit entries: %s", strings.Join(failed, "; ")).WithNext("rhiza model commit-untracked")
	}
	return resp, nil
}

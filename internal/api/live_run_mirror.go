package api

// A model write made outside a run for a change set that has a live run is
// also written into that run's worktree, so the run commits it on its branch
// at its next commit and the main checkout's copy is one the run holds, which
// close-out drops. A change set with no live run, a write made inside the
// run's own worktree and a dry run copy nothing.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
)

// liveRunMirror returns the function that copies a written entry into the
// worktree of its change set's live run. It reads the git directory from
// the file system, so a write runs no git command.
func liveRunMirror(ctx *Ctx) func(entry.Entry) error {
	return func(e entry.Entry) error {
		top, common, ok := gitDirsOf(ctx.Root)
		if !ok {
			return nil
		}
		j, ok, err := run.NewStore(run.Dir(common)).LiveOf(e.ChangeSet.String(), "")
		if err != nil || !ok {
			return nil
		}
		wt := j.State().Worktree
		if fi, err := os.Stat(wt); wt == "" || err != nil || !fi.IsDir() {
			return nil
		}
		root := wt
		if rel, err := filepath.Rel(top, ctx.Root); err == nil && rel != "." {
			root = filepath.Join(wt, rel)
		}
		if samePath(root, ctx.Root) {
			return nil
		}
		st, err := store.Open(root)
		if err != nil {
			return nil
		}
		if err := st.Write(e); err != nil && !errors.Is(err, store.ErrEntryExists) {
			return Errorf("copy %s into run %s's worktree: %v", e.ID, j.ID, err)
		}
		return nil
	}
}

// gitDirsOf finds the checkout's top folder and the git common dir above
// dir: a .git folder is the common dir itself, a .git file of a linked
// worktree names <common>/worktrees/<name>.
func gitDirsOf(dir string) (top, common string, ok bool) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", "", false
	}
	for {
		g := filepath.Join(d, ".git")
		if fi, err := os.Lstat(g); err == nil {
			if fi.IsDir() {
				return d, g, true
			}
			b, err := os.ReadFile(g)
			if err != nil {
				return "", "", false
			}
			p, found := strings.CutPrefix(strings.TrimSpace(string(b)), "gitdir:")
			if !found {
				return "", "", false
			}
			p = strings.TrimSpace(p)
			if !filepath.IsAbs(p) {
				p = filepath.Join(d, p)
			}
			p = filepath.Clean(p)
			if filepath.Base(filepath.Dir(p)) == "worktrees" {
				return d, filepath.Dir(filepath.Dir(p)), true
			}
			return "", "", false
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", "", false
		}
		d = parent
	}
}

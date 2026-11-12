package api

// A model write made outside a run for a change set that has a live run is
// also written into that run's worktree, so the run commits it on its branch
// at its next commit and the main checkout's copy is one the run holds, which
// close-out drops. A change set with no live run, a write made inside the
// run's own worktree and a dry run copy nothing.
//
// The copy takes every entry of the change set the main checkout holds and
// the worktree lacks, not the new entry alone, so the worktree never holds
// an entry whose predecessors (an open or propose still uncommitted in the
// main checkout) are missing. The run's own commits (commitOnly) then sweep
// in the copies still untracked there (mirroredPaths).

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
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
		if _, err := store.Open(root); err != nil {
			return nil
		}
		cs := e.ChangeSet.String()
		if err := copyMissingEntries(filepath.Join(ctx.Root, store.Dir, store.LogDir, cs), filepath.Join(root, store.Dir, store.LogDir, cs)); err != nil {
			return Errorf("copy %s into run %s's worktree: %v", e.ID, j.ID, err)
		}
		return nil
	}
}

// copyMissingEntries copies each entry file in from that to lacks, byte for
// byte, through a temporary file renamed into place. Entry files are never
// rewritten, so one already in to is left as it is.
func copyMissingEntries(from, to string) error {
	es, err := os.ReadDir(from)
	if err != nil {
		return err
	}
	for _, de := range es {
		name := de.Name()
		if de.IsDir() || filepath.Ext(name) != ".json" || strings.HasPrefix(name, ".") {
			continue
		}
		dst := filepath.Join(to, name)
		if _, err := os.Lstat(dst); err == nil {
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		b, err := os.ReadFile(filepath.Join(from, name))
		if err != nil {
			return err
		}
		if err := os.MkdirAll(to, 0o755); err != nil {
			return err
		}
		tmp, err := os.CreateTemp(to, ".mirror-*")
		if err != nil {
			return err
		}
		_, werr := tmp.Write(b)
		cerr := tmp.Close()
		if werr == nil {
			werr = cerr
		}
		if werr == nil {
			werr = os.Rename(tmp.Name(), dst)
		}
		if werr != nil {
			os.Remove(tmp.Name())
			return werr
		}
	}
	return nil
}

// mirroredPaths are the entry files of the run's change set, as paths from
// the worktree's top, that the worktree holds untracked and the main
// checkout holds with the same bytes: the copies liveRunMirror made, which
// the run's next commit takes. A runner with no main checkout has none.
func (r *stepRunner) mirroredPaths() ([]string, error) {
	if r.main == "" || samePath(r.main, r.worktree) {
		return nil, nil
	}
	dir := filepath.ToSlash(filepath.Join(r.prefix(), store.Dir, store.LogDir, r.change))
	out, err := r.worktreeGit("ls-files", "--others", "--exclude-standard", "--full-name", "-z", "--", dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, p := range strings.Split(out, "\x00") {
		if p == "" || path.Ext(p) != ".json" {
			continue
		}
		mine, err := os.ReadFile(filepath.Join(r.worktree, filepath.FromSlash(p)))
		if err != nil {
			continue
		}
		theirs, err := os.ReadFile(filepath.Join(r.main, filepath.FromSlash(p)))
		if err != nil || !bytes.Equal(mine, theirs) {
			continue
		}
		paths = append(paths, p)
	}
	return paths, nil
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

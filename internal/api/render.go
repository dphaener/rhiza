package api

// This file holds the render operation: it writes the as-is model as
// markdown under .genesis/rendered, the files check compares against. It
// reads the log through the store alone, never the workspace: render writes
// no log entry and needs no change set or actor.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	gmodel "github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/render"
	"github.com/dphaener/genesis/internal/store"
)

// renderedDir is the rendered folder, relative to the project root and
// written with slashes, as response paths are.
const renderedDir = store.Dir + "/rendered"

// RenderRequest takes nothing: render works in the caller's root.
type RenderRequest struct{}

// RenderResponse says what render changed. Paths are relative to the
// project root, with slashes, sorted.
type RenderResponse struct {
	Written   []string `json:"written" desc:"The files written: new, or different from what was on disk"`
	Removed   []string `json:"removed" desc:"The files render does not write, and then the folders left empty, that were removed"`
	Unchanged int      `json:"unchanged" desc:"How many rendered files were already current"`
}

func init() {
	register(Operation{
		Name:     "render",
		Summary:  "Write the as-is model as markdown under .genesis/rendered",
		Request:  reflect.TypeFor[RenderRequest](),
		Response: reflect.TypeFor[RenderResponse](),
		Run:      runRender,
	})
}

// openStore opens the repository in root. No repository is not_found with
// next "genesis init"; anything else wrong with it is an error.
func openStore(root string) (*store.Store, error) {
	st, err := store.Open(root)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, NotFoundf("no Genesis repository in %s", root).WithNext("genesis init")
		}
		return nil, Errorf("%s", err.Error())
	}
	return st, nil
}

// runRender renders the as-is view of the log and makes the rendered folder
// hold exactly those files: it removes the files render does not write,
// then the folders left empty (never the rendered folder itself), and
// writes each file whose bytes differ from the disk's, atomically. Render
// writes no log entry, so --dry-run is a usage error, as for init.
func runRender(ctx *Ctx, _ any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("render does not take --dry-run: it writes no log entry").WithNext("genesis render")
	}
	st, err := openStore(ctx.Root)
	if err != nil {
		return nil, err
	}
	sets, _ := changeset.FoldLog(st.ReadLog())
	m, _ := gmodel.AsIs(sets)
	want := render.Render(m, sets)

	dir := filepath.Join(ctx.Root, filepath.FromSlash(renderedDir))
	resp := RenderResponse{Written: []string{}, Removed: []string{}}
	removed, err := pruneRendered(dir, want)
	if err != nil {
		return nil, Errorf("render: %v", err)
	}
	resp.Removed = append(resp.Removed, removed...)
	for _, p := range slices.Sorted(maps.Keys(want)) {
		dst := filepath.Join(dir, filepath.FromSlash(p))
		if have, err := os.ReadFile(dst); err == nil && bytes.Equal(have, want[p]) {
			resp.Unchanged++
			continue
		}
		if err := writeAtomic(dst, want[p]); err != nil {
			return nil, Errorf("render: %v", err)
		}
		resp.Written = append(resp.Written, renderedDir+"/"+p)
	}
	return resp, nil
}

// pruneRendered removes from dir every file not in want, then every folder
// left empty that no file of want sits under. It returns the removed paths,
// relative to the project root and sorted; a missing dir removes nothing.
func pruneRendered(dir string, want map[string][]byte) ([]string, error) {
	keep := map[string]bool{}
	for p := range want {
		for d := path.Dir(p); d != "."; d = path.Dir(d) {
			keep[d] = true
		}
	}
	var removed, dirs []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == dir && errors.Is(err, fs.ErrNotExist) {
				return fs.SkipAll
			}
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		switch {
		case rel == ".":
		case d.IsDir():
			if !keep[rel] {
				dirs = append(dirs, rel)
			}
		default:
			if _, ok := want[rel]; !ok {
				if err := os.Remove(p); err != nil {
					return err
				}
				removed = append(removed, renderedDir+"/"+rel)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Deepest first, so a folder's subfolders go before it.
	slices.SortFunc(dirs, func(a, b string) int { return strings.Count(b, "/") - strings.Count(a, "/") })
	for _, d := range dirs {
		p := filepath.Join(dir, filepath.FromSlash(d))
		if es, err := os.ReadDir(p); err != nil || len(es) > 0 {
			continue
		}
		if err := os.Remove(p); err != nil {
			return nil, err
		}
		removed = append(removed, renderedDir+"/"+d)
	}
	slices.Sort(removed)
	return removed, nil
}

// writeAtomic writes data to dst through a temp file in dst's folder and a
// rename, creating the folder as needed: a reader sees the old bytes or the
// new, never a part. A temp file left by a crash is removed by the next
// render, since render does not write it.
func writeAtomic(dst string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(dst), ".render-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(tmp, 0o644)
	}
	if err == nil {
		err = os.Rename(tmp, dst)
	}
	if err != nil {
		os.Remove(tmp)
	}
	return err
}

// Human reports the counts, then each path written and removed.
func (r RenderResponse) Human(w io.Writer) {
	n := len(r.Written) + r.Unchanged
	files := "files"
	if n == 1 {
		files = "file"
	}
	fmt.Fprintf(w, "Rendered %d %s (%d written, %d removed, %d unchanged).\n",
		n, files, len(r.Written), len(r.Removed), r.Unchanged)
	if len(r.Written)+len(r.Removed) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, p := range r.Written {
		fmt.Fprintf(w, "- wrote %s\n", p)
	}
	for _, p := range r.Removed {
		fmt.Fprintf(w, "- removed %s\n", p)
	}
}

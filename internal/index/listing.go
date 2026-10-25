package index

import (
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/dphaener/genesis/internal/store"
)

// listWorkers is how many change set folders listLog reads at once.
const listWorkers = 4

// listLog is store.ListLog over fsys with the change set folders read
// concurrently. Reading ~550 folders one after another is most of a fresh
// open's cost; the listing it returns is the one store.ListLog returns.
func listLog(fsys fs.FS) store.Listing {
	return store.ListLog(&prefetchFS{FS: fsys})
}

// dirResult is one ReadDir's outcome.
type dirResult struct {
	des []fs.DirEntry
	err error
}

// prefetchFS answers ReadDir of the log folder by also reading, in
// parallel, every folder in it that store.ListLog will descend into, and
// answers those later ReadDir calls from what it read. Everything else
// goes straight to the wrapped FS.
type prefetchFS struct {
	fs.FS
	dirs map[string]dirResult
}

func (p *prefetchFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(p.FS, name)
}

func (p *prefetchFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if r, ok := p.dirs[name]; ok {
		delete(p.dirs, name)
		return r.des, r.err
	}
	des, err := fs.ReadDir(p.FS, name)
	if name != store.LogDir || err != nil {
		return des, err
	}
	var todo []string
	for _, de := range des {
		if de.IsDir() && !strings.HasPrefix(de.Name(), ".") {
			todo = append(todo, path.Join(name, de.Name()))
		}
	}
	results := make([]dirResult, len(todo))
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(listWorkers, len(todo)) {
		wg.Go(func() {
			for i := range next {
				d, err := fs.ReadDir(p.FS, todo[i])
				results[i] = dirResult{d, err}
			}
		})
	}
	for i := range todo {
		next <- i
	}
	close(next)
	wg.Wait()
	p.dirs = make(map[string]dirResult, len(todo))
	for i, dir := range todo {
		p.dirs[dir] = results[i]
	}
	return des, nil
}

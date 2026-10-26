package index

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io/fs"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

// racyWindow is how old every stat'ed modification time must be, against
// the moment the stats were taken, for a stat fingerprint to be stored.
// A folder changed within one tick of the file system's clock can keep
// the modification time it had: HFS+ keeps whole seconds, FAT two, and
// ext4 stamps with the kernel's coarse clock, which lags the wall clock.
// A time older than the window is past any such tick, so a later change
// gives the folder a later time and the fingerprint differs. Two seconds
// covers the coarsest of these.
const racyWindow = 2 * time.Second

// osKey is what the os adds to a folder's stat: its change time, which no
// one can set back, and its inode and device, which change when the folder
// is replaced.
type osKey struct {
	ctime    int64
	ino, dev uint64
}

// statPrint is a stat fingerprint of the log: the names in log/ and, for
// each folder store.ListLog descends into, its modification time, mode
// and os key.
type statPrint struct {
	// fp is the fingerprint, "" when a stat failed or a time is unknown.
	fp string
	// racy is true when a folder's modification time is inside racyWindow
	// of the moment the stats were taken (or after it).
	racy bool
}

// storable reports whether p may be stored as the stat fingerprint of a
// listing taken after it.
func (p statPrint) storable() bool { return p.fp != "" && !p.racy }

// statWorkers is how many goroutines takeStats stats the folders with. A
// stat is a system call that blocks its goroutine, not CPU work, so on a
// few cores several at once still finish sooner than one after another:
// measured on ~550 folders, four workers take about half the time of one
// at 11 cores and about three quarters at two; more than four gain
// nothing.
const statWorkers = 4

// statResult is one folder's stat.
type statResult struct {
	fi  fs.FileInfo
	err error
}

// takeStats reads log/ and stats each folder in it that store.ListLog
// would descend into, as prefetchFS selects them (folders whose names do
// not start with "."). now is the moment the stats are taken. Listing a
// change set folder is what costs; a stat is a fraction of it, and when
// no folder's stat changed no folder's names did. The stats are taken by
// statWorkers at once and hashed in the order log/ lists them, so the
// fingerprint does not depend on which finished first.
func takeStats(fsys fs.FS, now time.Time) statPrint {
	des, err := fs.ReadDir(fsys, store.LogDir)
	if err != nil {
		return statPrint{}
	}
	stats := make([]statResult, len(des))
	var next atomic.Int64
	var wg sync.WaitGroup
	for range min(statWorkers, len(des)) {
		wg.Go(func() {
			for k := int(next.Add(1)) - 1; k < len(des); k = int(next.Add(1)) - 1 {
				if de := des[k]; de.IsDir() && !strings.HasPrefix(de.Name(), ".") {
					stats[k].fi, stats[k].err = fs.Stat(fsys, path.Join(store.LogDir, de.Name()))
				}
			}
		})
	}
	wg.Wait()
	cutoff := now.Add(-racyWindow)
	h := sha256.New()
	var buf [8]byte
	num := func(v uint64) {
		binary.BigEndian.PutUint64(buf[:], v)
		h.Write(buf[:])
	}
	var p statPrint
	for k, de := range des {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			// store.ListLog skips these names, whatever they hold.
			continue
		}
		h.Write([]byte(name))
		h.Write([]byte{0})
		if !de.IsDir() {
			h.Write([]byte{'f'})
			continue
		}
		h.Write([]byte{'d'})
		fi, err := stats[k].fi, stats[k].err
		if err != nil {
			return statPrint{}
		}
		mt := fi.ModTime()
		if mt.IsZero() {
			// A file system that keeps no times cannot show a change.
			return statPrint{}
		}
		if !mt.Before(cutoff) {
			p.racy = true
		}
		num(uint64(mt.UnixNano()))
		num(uint64(fi.Mode()))
		if k, ok := statKey(fi); ok {
			num(uint64(k.ctime))
			num(k.ino)
			num(k.dev)
		}
	}
	p.fp = hex.EncodeToString(h.Sum(nil))
	return p
}

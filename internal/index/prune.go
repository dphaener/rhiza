package index

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// metaRoot is the meta key holding the resolved path of the worktree the
// index was opened for. The file name is only a hash of it, so this is how
// a sweep learns the worktree is gone.
const metaRoot = "root"

// recordRoot stores the resolved path of root in the index at dbPath.
// Errors are skipped: the record only serves the orphan sweep.
func recordRoot(x *Index, root string) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return
	}
	if v, err := metaOf(x.db, metaRoot); err == nil && v == resolved {
		return
	}
	_, _ = x.db.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaRoot, resolved)
}

// Pruned is what PruneOrphans removed.
type Pruned struct {
	Files int
	Bytes int64
}

// PruneOrphans removes the index databases in dir whose recorded worktree
// path no longer exists, with their WAL and shared-memory files. A database
// with no recorded path (made before paths were recorded) or one that cannot
// be read stays. Errors are skipped: this is housekeeping.
func PruneOrphans(dir string) Pruned {
	var p Pruned
	entries, err := os.ReadDir(dir)
	if err != nil {
		return p
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		root := recordedRoot(path)
		if root == "" {
			continue
		}
		if _, err := os.Lstat(root); !errors.Is(err, fs.ErrNotExist) {
			continue
		}
		var size int64
		for _, f := range []string{path, path + "-wal", path + "-shm"} {
			if info, err := os.Stat(f); err == nil {
				size += info.Size()
			}
		}
		if removeDB(path) != nil {
			continue
		}
		p.Files++
		p.Bytes += size
	}
	return p
}

// IndexDir is the folder holding every worktree's index database.
func IndexDir() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "rhiza", "index"), nil
}

// recordedRoot reads the worktree path stored in the database at path, ""
// when there is none or the file cannot be read. It never creates or
// changes the file.
func recordedRoot(path string) string {
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro&_pragma=busy_timeout(2000)")
	if err != nil {
		return ""
	}
	defer db.Close()
	v, err := metaOf(db, metaRoot)
	if err != nil {
		return ""
	}
	return v
}

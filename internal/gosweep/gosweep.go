// Package gosweep trims the Go build cache and the go-build* temp folders a
// killed build leaves behind. Go itself only trims cache entries unused for
// five days, which one evening of runs can outgrow the disk long before.
package gosweep

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// CacheMaxAge is how long a cache entry may sit unused: go refreshes an
	// entry's mtime when it uses it.
	CacheMaxAge = 3 * time.Hour
	// TempMaxAge is how old a go-build* temp folder must be to go.
	TempMaxAge = 2 * time.Hour
	// TempPrefix starts the name of every temp folder go makes.
	TempPrefix = "go-build"
)

// Freed is what a sweep removed.
type Freed struct {
	CacheEntries int
	CacheBytes   int64
	TempFolders  int
	TempBytes    int64
}

// Bytes is everything freed.
func (f Freed) Bytes() int64 { return f.CacheBytes + f.TempBytes }

// Any reports whether the sweep freed anything.
func (f Freed) Any() bool { return f.CacheEntries > 0 || f.TempFolders > 0 }

// CacheDir is the Go build cache: $GOCACHE, or os.UserCacheDir()/go-build
// when it is unset. It is empty when the cache is off or cannot be placed.
func CacheDir() string {
	if v := os.Getenv("GOCACHE"); v != "" {
		if v == "off" {
			return ""
		}
		return v
	}
	d, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(d, "go-build")
}

// Sweep trims the default cache and os.TempDir() as of now.
func Sweep(now time.Time) Freed {
	return SweepDirs(CacheDir(), os.TempDir(), now)
}

// SweepDirs removes the entries of cacheDir (empty means none) not modified
// since now-CacheMaxAge, and the go-build* folders directly under tempDir not
// modified since now-TempMaxAge. Errors are skipped: this is housekeeping.
func SweepDirs(cacheDir, tempDir string, now time.Time) Freed {
	var f Freed
	if cacheDir != "" {
		f.CacheEntries, f.CacheBytes = sweepCache(cacheDir, now.Add(-CacheMaxAge))
	}
	if tempDir != "" {
		f.TempFolders, f.TempBytes = sweepTemp(tempDir, now.Add(-TempMaxAge))
	}
	return f
}

// isShard reports whether name is one of the cache's two-hex-digit folders.
func isShard(name string) bool {
	if len(name) != 2 {
		return false
	}
	for _, c := range name {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}

func sweepCache(dir string, cutoff time.Time) (n int, bytes int64) {
	shards, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}
	for _, s := range shards {
		if !s.IsDir() || !isShard(s.Name()) {
			continue
		}
		sd := filepath.Join(dir, s.Name())
		files, err := os.ReadDir(sd)
		if err != nil {
			continue
		}
		for _, e := range files {
			if e.IsDir() {
				continue
			}
			info, err := e.Info()
			if err != nil || !info.ModTime().Before(cutoff) {
				continue
			}
			if os.Remove(filepath.Join(sd, e.Name())) == nil {
				n++
				bytes += info.Size()
			}
		}
	}
	return n, bytes
}

func sweepTemp(dir string, cutoff time.Time) (n int, bytes int64) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), TempPrefix) {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.ModTime().Before(cutoff) {
			continue
		}
		p := filepath.Join(dir, e.Name())
		size := treeSize(p)
		if os.RemoveAll(p) == nil {
			n++
			bytes += size
		}
	}
	return n, bytes
}

func treeSize(root string) int64 {
	var total int64
	filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			total += info.Size()
		}
		return nil
	})
	return total
}

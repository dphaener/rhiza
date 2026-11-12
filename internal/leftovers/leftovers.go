// Package leftovers prunes what ended runs left on disk: the scratch folders
// of agent sessions that belong to no live run.
package leftovers

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// ScratchMaxAge is how old a session's scratch folder must be to go.
const ScratchMaxAge = 24 * time.Hour

// Session is one recorded agent session.
type Session struct {
	ID    string
	Agent string
	Run   string
}

// Freed is what a sweep removed.
type Freed struct {
	Folders int
	Bytes   int64
}

// Scratch removes the scratch folders of the sessions whose run is not live
// (live reports it) and whose folder was last modified before now-ScratchMaxAge.
// globs gives the patterns of one session's folders. A session of a live run
// is never touched, nor is a fresh folder. Errors are skipped: this is
// housekeeping.
func Scratch(sessions []Session, live func(run string) bool, globs func(Session) []string, now time.Time) Freed {
	var f Freed
	cutoff := now.Add(-ScratchMaxAge)
	for _, s := range sessions {
		if s.Run != "" && live(s.Run) {
			continue
		}
		for _, g := range globs(s) {
			matches, err := filepath.Glob(g)
			if err != nil {
				continue
			}
			for _, p := range matches {
				info, err := os.Stat(p)
				if err != nil || !info.IsDir() || !info.ModTime().Before(cutoff) {
					continue
				}
				size := dirSize(p)
				if os.RemoveAll(p) != nil {
					continue
				}
				f.Folders++
				f.Bytes += size
			}
		}
	}
	return f
}

func dirSize(root string) int64 {
	var n int64
	_ = filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if info, err := d.Info(); err == nil {
			n += info.Size()
		}
		return nil
	})
	return n
}

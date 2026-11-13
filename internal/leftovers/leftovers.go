// Package leftovers prunes what ended runs left on disk: the scratch folders
// of agent sessions that belong to no live run.
package leftovers

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// ReadDir lists one directory, as os.ReadDir does. A sweep takes it so a test
// can count the reads.
type ReadDir func(name string) ([]os.DirEntry, error)

// SafeID reports whether id can stand as one folder name under a scratch
// parent: not empty, no separator, no parent reference and no pattern syntax.
func SafeID(id string) bool {
	return id != "" && !strings.ContainsAny(id, `/\*?[`) && !strings.Contains(id, "..")
}

// Scratch removes the scratch folders of the sessions whose run is not live
// (live reports it) and whose folder was last modified before now-ScratchMaxAge.
// parents gives, for an agent, the glob patterns of the folders whose direct
// children are session folders named by session ID. Each parent is listed once
// (through read, os.ReadDir when nil), however many sessions there are, and
// its children are matched against the orphaned sessions' IDs. A session of a
// live run is never touched, nor is a fresh folder. Errors are skipped: this
// is housekeeping.
func Scratch(sessions []Session, live func(run string) bool, parents func(agent string) []string, read ReadDir, now time.Time) Freed {
	if read == nil {
		read = os.ReadDir
	}
	var f Freed
	cutoff := now.Add(-ScratchMaxAge)
	byAgent := map[string]map[string]bool{}
	for _, s := range sessions {
		if !SafeID(s.ID) || (s.Run != "" && live(s.Run)) {
			continue
		}
		if byAgent[s.Agent] == nil {
			byAgent[s.Agent] = map[string]bool{}
		}
		byAgent[s.Agent][s.ID] = true
	}
	agents := make([]string, 0, len(byAgent))
	for a := range byAgent {
		agents = append(agents, a)
	}
	sort.Strings(agents)
	// ids per parent folder, merged across agents that share one.
	ids := map[string]map[string]bool{}
	var dirs []string
	for _, a := range agents {
		for _, g := range parents(a) {
			matches, err := filepath.Glob(g)
			if err != nil {
				continue
			}
			for _, m := range matches {
				if ids[m] == nil {
					ids[m] = map[string]bool{}
					dirs = append(dirs, m)
				}
				for id := range byAgent[a] {
					ids[m][id] = true
				}
			}
		}
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		entries, err := read(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() || !ids[dir][e.Name()] {
				continue
			}
			info, err := e.Info()
			if err != nil || !info.ModTime().Before(cutoff) {
				continue
			}
			p := filepath.Join(dir, e.Name())
			size := dirSize(p)
			if os.RemoveAll(p) != nil {
				continue
			}
			f.Folders++
			f.Bytes += size
		}
	}
	return f
}

// RemoveSession removes the folder of session id under each of the scratch
// parents (glob patterns). An unsafe ID removes nothing. Errors are skipped.
func RemoveSession(parents []string, id string) {
	if !SafeID(id) {
		return
	}
	for _, g := range parents {
		matches, err := filepath.Glob(g)
		if err != nil {
			continue
		}
		for _, m := range matches {
			_ = os.RemoveAll(filepath.Join(m, id))
		}
	}
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

// Package rootguard lets a package's tests prove they left the state folder
// of a Rhiza root alone. A checkout under ~/code sits under a root, and a
// test that reaches that root (a fake usage-limit hit, a stray file) writes
// into the machine's real .rhiza-root/. TestMain wraps m.Run in Main, which
// snapshots the folder before the run and, after it, reports each write it
// attributes to the tests. The live root daemon and run drivers write to
// the same folder at the same time, so files they own are never reported.
package rootguard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

// StateDir is the name of a root's state folder.
const StateDir = ".rhiza-root"

// limitsFile is the usage-limit log; its lines are judged by hit_at.
const limitsFile = "usage.limits.jsonl"

// owned are the files the root daemon and run drivers write. A change to
// one is never attributed to a test run. The limits file is among them,
// but its appended lines are still judged by hit_at.
var owned = map[string]bool{
	"daemon.lock":       true,
	"daemon.home":       true,
	"daemon.state.json": true,
	"requests.json":     true,
	"requests.lock":     true,
	"starts.json":       true,
	"starts.lock":       true,
	"continues.json":    true,
	"continues.lock":    true,
	"main-stale.json":   true,
	limitsFile:          true,
}

type fileState struct {
	size int64
	mod  time.Time
}

// Snapshot is the state of a root's .rhiza-root/ at one moment.
type Snapshot struct {
	dir   string // the state folder; empty when no root encloses the dir
	files map[string]fileState
	lines int // lines in usage.limits.jsonl
}

// Take snapshots the state folder of the root enclosing dir. With no
// enclosing root, or no state folder, the snapshot is empty and Check
// reports nothing unless a folder appears.
func Take(dir string) *Snapshot {
	s := &Snapshot{files: map[string]fileState{}}
	root, ok := rootcfg.Find(dir)
	if !ok {
		return s
	}
	s.dir = filepath.Join(root, StateDir)
	s.files = scan(s.dir)
	s.lines = len(readLines(filepath.Join(s.dir, limitsFile)))
	return s
}

func scan(dir string) map[string]fileState {
	out := map[string]fileState{}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return nil
		}
		out[filepath.ToSlash(rel)] = fileState{fi.Size(), fi.ModTime()}
		return nil
	})
	return out
}

func readLines(path string) [][]byte {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines [][]byte
	for _, l := range bytes.Split(b, []byte("\n")) {
		if len(bytes.TrimSpace(l)) > 0 {
			lines = append(lines, l)
		}
	}
	return lines
}

// Check compares the folder now with the snapshot and returns one message
// per write attributed to a test run that lasted from start to end: a file
// that appeared or changed and no daemon or driver owns, and a line
// appended to usage.limits.jsonl whose hit_at lies outside [start, end]
// (tests stamp hits with fixed clocks, real drivers with the current time).
func (s *Snapshot) Check(start, end time.Time) []string {
	if s.dir == "" {
		return nil
	}
	var out []string
	now := scan(s.dir)
	names := make([]string, 0, len(now))
	for n := range now {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if owned[n] {
			continue
		}
		old, had := s.files[n]
		switch {
		case !had:
			out = append(out, fmt.Sprintf("%s/%s: file added", StateDir, n))
		case old != now[n]:
			out = append(out, fmt.Sprintf("%s/%s: file changed", StateDir, n))
		}
	}
	lines := readLines(filepath.Join(s.dir, limitsFile))
	if len(lines) > s.lines {
		for i, l := range lines[s.lines:] {
			var hit struct {
				HitAt time.Time `json:"hit_at"`
			}
			if err := json.Unmarshal(l, &hit); err != nil {
				out = append(out, fmt.Sprintf("%s/%s: line %d appended, unreadable: %v", StateDir, limitsFile, s.lines+i+1, err))
				continue
			}
			if hit.HitAt.Before(start) || hit.HitAt.After(end) {
				out = append(out, fmt.Sprintf("%s/%s: line %d appended with hit_at %s, outside the test run (%s to %s)",
					StateDir, limitsFile, s.lines+i+1, hit.HitAt.Format(time.RFC3339), start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339)))
			}
		}
	}
	return out
}

// Main runs the tests with run (m.Run, or testscript.RunMain) guarded: it
// snapshots the state folder of the root enclosing the working directory,
// runs, and returns a non-zero code, printing each attributed write to
// stderr, when the run wrote into it. Use it as os.Exit(rootguard.Main(run)).
func Main(run func() int) int {
	dir, _ := os.Getwd()
	snap := Take(dir)
	start := time.Now()
	code := run()
	end := time.Now()
	if msgs := snap.Check(start, end); len(msgs) > 0 {
		fmt.Fprintln(os.Stderr, "FAIL: the test run wrote into the state folder of the root enclosing it:")
		for _, msg := range msgs {
			fmt.Fprintln(os.Stderr, "  "+msg)
		}
		if code == 0 {
			code = 1
		}
	}
	return code
}

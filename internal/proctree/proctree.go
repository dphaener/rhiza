// Package proctree reads the host's process table and gives an agent
// session's process tree: the `rhiza session run <session-id>` process and
// every descendant, each with its pid, parent, resident memory and command
// line, plus the count and the summed resident memory.
//
// The table comes from one ps call, never /proc, so that Linux and macOS
// read the same: both print rss in KiB. Parsing is pure (Parse,
// SessionTree); only Read and ReadTable run ps.
package proctree

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// Process is one row of the process table.
type Process struct {
	PID  int    `json:"pid"`
	PPID int    `json:"ppid"`
	RSS  int64  `json:"rss"` // resident memory in bytes
	Args string `json:"args"`
}

// Tree is a session's process tree: its root process (the
// `rhiza session run <session-id>` one) and every descendant, root first,
// then by pid.
type Tree struct {
	Session   string    `json:"session"`
	Root      int       `json:"root"`
	Processes []Process `json:"processes"`
	Count     int       `json:"count"`
	RSS       int64     `json:"rss"` // summed resident memory in bytes
}

// ErrNotFound is returned when no `rhiza session run <session-id>` process
// is in the table: the session has no tree to read, which is not a reading
// of zero processes.
var ErrNotFound = errors.New("session process not found")

// psArgs asks for every process, unlimited width (BSD ps truncates args to
// the terminal's otherwise), with no header: pid, parent pid, resident
// memory in KiB and the full command line.
var psArgs = []string{"-A", "-ww", "-o", "pid=", "-o", "ppid=", "-o", "rss=", "-o", "args="}

// runPS runs ps; tests replace it to fail without running ps.
var runPS = func() ([]byte, error) {
	return exec.Command("ps", psArgs...).Output()
}

// ReadTable reads the host's process table with one ps call.
func ReadTable() ([]Process, error) {
	out, err := runPS()
	if err != nil {
		return nil, fmt.Errorf("ps: %w", err)
	}
	return Parse(string(out)), nil
}

// Read reads the host's process table and returns session's tree, or
// ErrNotFound when no process runs `rhiza session run <session>`.
func Read(session string) (Tree, error) {
	procs, err := ReadTable()
	if err != nil {
		return Tree{}, err
	}
	return SessionTree(procs, session)
}

// Parse reads ps output, one process per line as pid, ppid, rss in KiB and
// the command line. Lines that do not hold all four, or whose numbers do
// not parse (a header line among them), are skipped.
func Parse(ps string) []Process {
	var procs []Process
	for line := range strings.Lines(ps) {
		p, ok := parseLine(line)
		if ok {
			procs = append(procs, p)
		}
	}
	return procs
}

func parseLine(line string) (Process, bool) {
	var nums [3]int64
	rest := line
	for i := range nums {
		rest = strings.TrimLeft(rest, " \t")
		end := strings.IndexAny(rest, " \t\r\n")
		if end < 0 {
			return Process{}, false
		}
		n, err := strconv.ParseInt(rest[:end], 10, 64)
		if err != nil || n < 0 {
			return Process{}, false
		}
		nums[i] = n
		rest = rest[end:]
	}
	args := strings.TrimSpace(rest)
	if args == "" {
		return Process{}, false
	}
	return Process{PID: int(nums[0]), PPID: int(nums[1]), RSS: nums[2] * 1024, Args: args}, true
}

// IsSessionRoot says whether args is the command line of
// `<rhiza> session run <session>`: the program's first three arguments are
// session, run and the session ID, matched whole (an ID that only shares a
// prefix is another session). A command that only passes that line on (a
// Rex tab opening it) is not the root.
func IsSessionRoot(args, session string) bool {
	f := strings.Fields(args)
	return session != "" && len(f) >= 4 && f[1] == "session" && f[2] == "run" && f[3] == session
}

// SessionTree returns session's tree from a process table: the process
// running `rhiza session run <session>` and every descendant, or
// ErrNotFound when no process runs it.
func SessionTree(procs []Process, session string) (Tree, error) {
	children := map[int][]Process{}
	var roots []Process
	for _, p := range procs {
		children[p.PPID] = append(children[p.PPID], p)
		if IsSessionRoot(p.Args, session) {
			roots = append(roots, p)
		}
	}
	if len(roots) == 0 {
		return Tree{}, fmt.Errorf("%w: rhiza session run %s", ErrNotFound, session)
	}
	// A root nested under another (one session run starting the same
	// session again) belongs to the outer tree; the outermost, lowest pid
	// first, is the tree's root.
	slices.SortFunc(roots, func(a, b Process) int { return a.PID - b.PID })
	t := Tree{Session: session, Root: roots[0].PID}
	seen := map[int]bool{}
	var walk func(p Process)
	walk = func(p Process) {
		if seen[p.PID] {
			return
		}
		seen[p.PID] = true
		t.Processes = append(t.Processes, p)
		for _, c := range children[p.PID] {
			if c.PID != p.PID {
				walk(c)
			}
		}
	}
	for _, r := range roots {
		walk(r)
	}
	slices.SortStableFunc(t.Processes[1:], func(a, b Process) int { return a.PID - b.PID })
	t.Count = len(t.Processes)
	for _, p := range t.Processes {
		t.RSS += p.RSS
	}
	return t, nil
}

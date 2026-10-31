// Package quiettest lets a timing test wait until go test runs nothing else.
// go test runs packages at once, and a read timed beside another package's
// build or test measures that contention, not the read.
package quiettest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// quietWait is the most Wait waits for the other packages' builds
// and tests to finish before it times the reads. go test runs -p packages
// at once (GOMAXPROCS by default), and budgetLock orders only the three
// budget tests: on GitHub's two-core runner the rest of gen, index and
// slice ran beside TestReadBudget and every read measured about 3x slower
// (show 1.3 -> 4.2 ms, roadmap 8.9 -> 28.7 ms median). The reads do the
// same work either way; the wait decides only when they are timed.
const quietWait = 3 * time.Minute

// quietInterval is how long Wait pauses between looks, and
// quietChecks how many looks in a row must find nothing, two seconds in
// all: between one package's test ending and the next one's starting, go
// test runs nothing, and on a loaded runner that gap reached half a second
// (a run that declared itself alone after two looks timed its reads beside
// slice.test, which started right after).
const (
	quietInterval = 250 * time.Millisecond
	quietChecks   = 8
)

// otherTestProcesses returns the names of the processes go test is running
// beside this one: the other children of this process's parent that are a
// package's test binary or a toolchain step. ps is asked, not /proc, so
// that Linux and darwin read the same.
func otherTestProcesses() ([]string, error) {
	out, err := exec.Command("ps", "-A", "-ww", "-o", "pid=", "-o", "ppid=", "-o", "args=").Output()
	if err != nil {
		return nil, fmt.Errorf("ps: %w", err)
	}
	return testProcesses(string(out), os.Getpid(), os.Getppid()), nil
}

// testProcesses reads ps output, one process per line as pid, ppid and
// command line, and returns the command names of the test-run processes
// (see isTestProcess) whose parent is ppid, other than pid itself.
func testProcesses(ps string, pid, ppid int) []string {
	var names []string
	for line := range strings.Lines(ps) {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		p, err := strconv.Atoi(fields[0])
		if err != nil || p == pid {
			continue
		}
		pp, err := strconv.Atoi(fields[1])
		if err != nil || pp != ppid {
			continue
		}
		if name := filepath.Base(fields[2]); isTestProcess(name) {
			names = append(names, name)
		}
	}
	return names
}

// isTestProcess says whether a command name is one go test runs for a
// package: its test binary, named after the package with a .test suffix,
// or a step of building and vetting it.
func isTestProcess(name string) bool {
	switch name {
	case "compile", "link", "asm", "cgo", "vet":
		return true
	}
	return strings.HasSuffix(name, ".test")
}

// Quiet is what Wait found: how long it waited, the processes it
// last saw running (nil when it never saw any), what was still running
// when quietWait ran out (nil once the run went quiet), and the error if
// ps could not be asked.
type Quiet struct {
	waited  time.Duration
	last    []string
	running []string
	err     error
}

// Wait waits until go test runs no other package's build or test
// beside this process, for at most quietWait. It must run before
// budgetLock is taken: the budget tests of index and slice wait on that
// lock, so holding it here would keep the run from ever going quiet.
func Wait() Quiet {
	start := time.Now()
	var q Quiet
	checks := 0
	for {
		running, err := otherTestProcesses()
		if err != nil {
			q.err = err
			break
		}
		if len(running) == 0 {
			if checks++; checks == quietChecks {
				break
			}
		} else {
			checks = 0
			q.last = running
			if time.Since(start) >= quietWait {
				q.running = running
				break
			}
		}
		time.Sleep(quietInterval)
	}
	q.waited = time.Since(start)
	return q
}

// String says how the reads were timed, for the timing log's header:
// alone, after waiting for the other packages (naming the last seen, so a
// read that still looks slow can be told from a quiet one), or beside
// what was still running when quietWait ran out.
func (q Quiet) String() string {
	switch {
	case q.err != nil:
		return fmt.Sprintf("other packages unknown (%v)", q.err)
	case len(q.running) > 0:
		return fmt.Sprintf("beside %s after %v", strings.Join(q.running, " "), q.waited.Round(time.Second))
	case len(q.last) > 0:
		return fmt.Sprintf("alone after %v (last beside %s)", q.waited.Round(100*time.Millisecond), strings.Join(q.last, " "))
	default:
		return fmt.Sprintf("alone after %v", q.waited.Round(100*time.Millisecond))
	}
}

// Late says which other packages' processes were running when the
// reads named by read and view had just been timed, for a note at the end
// of the timing log; "" when none were. A package that go test started
// after the wait found the run quiet shows up here, so a slow read is
// never silently taken for a quiet one.
func Late(read, view string) string {
	running, err := otherTestProcesses()
	if err != nil || len(running) == 0 {
		return ""
	}
	return fmt.Sprintf("note: %s ran beside %s %s and the reads after it", strings.Join(running, " "), read, view)
}

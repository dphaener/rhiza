package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/gen"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
)

// readBudget is the most a read's median may take, fresh workspace and
// index open to answer: milestone 1's target for every read on the
// 50,000-node model (D14, as revised).
const readBudget = 50 * time.Millisecond

// readBudgetRuns is how many fresh runs TestReadBudget times per read.
const readBudgetRuns = 20

// generated is gen.Default() written once, with its index built, shared by
// TestReadBudget and the benchmarks. It lives under one os.MkdirTemp folder
// that TestMain removes.
type generated struct {
	root      string
	indexPath string
	summary   gen.Summary
	// subtree is the deep subtree tree reads (D14): the deep sample's
	// grandparent, two levels above the deepest.
	subtree id.ID
	err     error
}

var (
	genOnce sync.Once
	genDir  string
	genned  generated
)

func TestMain(m *testing.M) {
	code := m.Run()
	if genDir != "" {
		os.RemoveAll(genDir)
	}
	os.Exit(code)
}

// defaultModel returns the shared generated model, writing and indexing it
// on first use.
func defaultModel(t testing.TB) generated {
	t.Helper()
	genOnce.Do(func() {
		dir, err := os.MkdirTemp("", "rhiza-api-budget-")
		if err != nil {
			genned.err = err
			return
		}
		genDir = dir
		root := filepath.Join(dir, "repo")
		if err := os.Mkdir(root, 0o755); err != nil {
			genned.err = err
			return
		}
		s, err := gen.Write(root, gen.Default())
		if err != nil {
			genned.err = err
			return
		}
		indexPath := filepath.Join(dir, "index.db")
		x, err := index.Open(root, indexPath)
		if err != nil {
			genned.err = err
			return
		}
		as, err := x.Ancestors(index.AsIs, s.Deep.ID)
		if err == nil && len(as) < 2 {
			err = fmt.Errorf("deep node %s has %d ancestors, want 2 or more", s.Deep.ID, len(as))
		}
		if cerr := x.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			genned.err = err
			return
		}
		genned = generated{root: root, indexPath: indexPath, summary: s, subtree: as[len(as)-2].ID}
	})
	if genned.err != nil {
		t.Fatalf("generate default model: %v", genned.err)
	}
	return genned
}

// budgetLock holds an exclusive lock on a file shared by every test that
// builds the 50,000-node model (the budget tests of internal/api,
// internal/index and internal/slice), until t ends. Each takes it before
// it writes and indexes its model: go test runs packages at once, and a
// budget timed while another package builds or reads its own model would
// measure that contention, not the read.
func budgetLock(t *testing.T) {
	t.Helper()
	f, err := os.OpenFile(filepath.Join(os.TempDir(), "rhiza-budget.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	})
}

// quietWait is the most TestReadBudget waits for the other packages' builds
// and tests to finish before it times the reads. go test runs -p packages
// at once (GOMAXPROCS by default), and budgetLock orders only the three
// budget tests: on GitHub's two-core runner the rest of gen, index and
// slice ran beside TestReadBudget and every read measured about 3x slower
// (show 1.3 -> 4.2 ms, roadmap 8.9 -> 28.7 ms median). The reads do the
// same work either way; the wait decides only when they are timed.
const quietWait = 3 * time.Minute

// quietInterval is how long waitForQuiet pauses between looks, and
// quietChecks how many looks in a row must find nothing: go test starts
// the next package's build or test within milliseconds of the last one
// ending, and a single look could land in that gap.
const (
	quietInterval = 250 * time.Millisecond
	quietChecks   = 2
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

// waitForQuiet waits until go test runs no other package's build or test
// beside this process, for at most quietWait, and returns how long it
// waited and what was still running when it stopped (nil once the run was
// quiet). It must run before budgetLock is taken: the budget tests of
// index and slice wait on that lock, so holding it here would keep the
// run from ever going quiet.
func waitForQuiet() (waited time.Duration, running []string, err error) {
	start := time.Now()
	quiet := 0
	for {
		running, err = otherTestProcesses()
		if err != nil {
			return time.Since(start), nil, err
		}
		if len(running) == 0 {
			if quiet++; quiet == quietChecks {
				return time.Since(start), nil, nil
			}
		} else {
			quiet = 0
			if time.Since(start) >= quietWait {
				return time.Since(start), running, nil
			}
		}
		time.Sleep(quietInterval)
	}
}

// quietReport says how the reads were timed, for the timing log's header:
// alone, after waiting for the other packages, or beside what was still
// running when quietWait ran out.
func quietReport(waited time.Duration, running []string, err error) string {
	switch {
	case err != nil:
		return fmt.Sprintf("other packages unknown (%v)", err)
	case len(running) > 0:
		return fmt.Sprintf("beside %s after %v", strings.Join(running, " "), waited.Round(time.Second))
	default:
		return fmt.Sprintf("alone after %v", waited.Round(100*time.Millisecond))
	}
}

// TestTestProcesses reads ps output as Linux (procps) and darwin print it
// for the same go test run: this process is pid 300 under the go command,
// pid 200.
func TestTestProcesses(t *testing.T) {
	const ps = `    1     0 /sbin/init splash
  100     1 zsh
  200   100 go test ./...
  300   200 /tmp/go-build1/b001/api.test -test.paniconexit0 -test.timeout=10m0s
  301   200 /tmp/go-build1/b002/slice.test -test.paniconexit0 -test.timeout=10m0s
  302   200 /usr/local/go/pkg/tool/linux_amd64/compile -o /tmp/go-build1/b003/_pkg_.a -trimpath /tmp/go-build1/b003=> -p index
  303   200 /usr/local/go/pkg/tool/linux_amd64/vet -atomic -bool -buildtags
  304   200 /usr/local/go/pkg/tool/linux_amd64/link -o /tmp/go-build1/b004/gen.test
  305   301 /tmp/go-build1/b002/slice.test -test.run=^TestChild$
  306   100 tail -5
  307   200 git status
  308   200 /private/var/folders/x/T/go-build2/b005/changeset.test -test.paniconexit0
bad line
`
	got := testProcesses(ps, 300, 200)
	want := []string{"slice.test", "compile", "vet", "link", "changeset.test"}
	if !slices.Equal(got, want) {
		t.Errorf("testProcesses = %q, want %q", got, want)
	}
	if got := testProcesses(ps, 300, 999); got != nil {
		t.Errorf("testProcesses under another parent = %q, want none", got)
	}
	if got := testProcesses("", 300, 200); got != nil {
		t.Errorf("testProcesses of nothing = %q, want none", got)
	}
}

func TestQuietReport(t *testing.T) {
	cases := []struct {
		waited  time.Duration
		running []string
		err     error
		want    string
	}{
		{250 * time.Millisecond, nil, nil, "alone after 300ms"},
		{41340 * time.Millisecond, nil, nil, "alone after 41.3s"},
		{3 * time.Minute, []string{"index.test", "slice.test"}, nil, "beside index.test slice.test after 3m0s"},
		{0, nil, fmt.Errorf("ps: exec: no ps"), "other packages unknown (ps: exec: no ps)"},
	}
	for _, c := range cases {
		if got := quietReport(c.waited, c.running, c.err); got != c.want {
			t.Errorf("quietReport(%v, %q, %v) = %q, want %q", c.waited, c.running, c.err, got, c.want)
		}
	}
}

// budgetRead is one read the budget times: an operation and its request.
type budgetRead struct {
	name string
	op   string
	req  any
}

// answerSize is how much a read's response holds: the nodes of show, ls
// and tree, the change sets of why, log and ready, the groups of roadmap.
// The budget wants it positive, so that no read is timed answering
// nothing.
func answerSize(resp any) int {
	switch r := resp.(type) {
	case ShowResponse:
		if r.Node != nil {
			return 1
		}
	case LsResponse:
		return len(r.Nodes)
	case TreeResponse:
		return len(r.Nodes)
	case WhyResponse:
		return len(r.ChangeSets)
	case LogResponse:
		return len(r.ChangeSets)
	case RoadmapResponse:
		return len(r.Groups)
	case ReadyResponse:
		return len(r.ChangeSets)
	}
	return 0
}

// budgetReads is every budgeted read of g in view v (D14): show of a deep
// node by handle, the first page of ls, a deep subtree, why and log of a
// deep node, roadmap and ready.
func budgetReads(g generated, v string) []budgetRead {
	s := g.summary
	return []budgetRead{
		{"show", "show", ShowRequest{Ref: s.Deep.Handle, View: v}},
		{"ls", "ls", LsRequest{View: v}},
		{"tree", "tree", TreeRequest{Ref: g.subtree.String(), View: v}},
		{"why", "why", WhyRequest{Ref: s.Deep.ID.String(), View: v}},
		{"log", "log", LogRequest{Ref: s.Deep.ID.String(), View: v}},
		{"roadmap", "roadmap", RoadmapRequest{View: v}},
		{"ready", "ready", ReadyRequest{View: v}},
	}
}

// runRead runs rd once as the CLI does: a fresh Ctx, so a fresh workspace
// and index open.
func runRead(g generated, rd budgetRead) (any, error) {
	ctx := testCtx(g.root)
	ctx.IndexPath = g.indexPath
	o, _ := Lookup(rd.op)
	return o.Run(ctx, rd.req)
}

// openTimingLog returns a function that appends a line to the file named by
// RHIZA_BUDGET_LOG, so that CI can print the timings after the run, passing
// or failing. The file is opened in append mode, starts its block with a
// header naming GOMAXPROCS, the CPU count and quiet (what waitForQuiet
// reported), and is closed when t ends. An open or write error fails t.
// With the variable unset it returns a no-op.
func openTimingLog(t *testing.T, quiet string) func(line string) {
	t.Helper()
	path := os.Getenv("RHIZA_BUDGET_LOG")
	if path == "" {
		return func(string) {}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("RHIZA_BUDGET_LOG: %v", err)
	}
	t.Cleanup(func() {
		if err := f.Close(); err != nil {
			t.Errorf("RHIZA_BUDGET_LOG %s: %v", path, err)
		}
	})
	write := func(line string) {
		t.Helper()
		if _, err := fmt.Fprintln(f, line); err != nil {
			t.Fatalf("RHIZA_BUDGET_LOG %s: %v", path, err)
		}
	}
	write(fmt.Sprintf("read budget: GOMAXPROCS=%d NumCPU=%d, %s", runtime.GOMAXPROCS(0), runtime.NumCPU(), quiet))
	return write
}

func TestReadBudget(t *testing.T) {
	if testing.Short() {
		t.Skip("the read budget generates the 50,000-node model; skipped under -short")
	}
	quiet := quietReport(waitForQuiet())
	t.Logf("other packages: %s", quiet)
	budgetLock(t)
	timingLog := openTimingLog(t, quiet)
	g := defaultModel(t)
	for _, v := range []string{viewAsIs, viewToBe} {
		for _, rd := range budgetReads(g, v) {
			var times []time.Duration
			var size int
			for range readBudgetRuns {
				start := time.Now()
				resp, err := runRead(g, rd)
				times = append(times, time.Since(start))
				if err != nil {
					t.Fatalf("%s %s: %v", rd.name, v, err)
				}
				if size = answerSize(resp); size == 0 {
					t.Fatalf("%s %s answered nothing: %+v", rd.name, v, resp)
				}
			}
			slices.Sort(times)
			median := (times[readBudgetRuns/2-1] + times[readBudgetRuns/2]) / 2
			line := fmt.Sprintf("%-7s %-5s median %v (min %v, max %v), size %d", rd.name, v, median, times[0], times[len(times)-1], size)
			t.Log(line)
			timingLog(line)
			if median >= readBudget {
				t.Errorf("%s %s median %v, want < %v", rd.name, v, median, readBudget)
			}
		}
	}
}

// benchRead runs the read named name in to-be b.N times.
func benchRead(b *testing.B, name string) {
	g := defaultModel(b)
	i := slices.IndexFunc(budgetReads(g, viewToBe), func(rd budgetRead) bool { return rd.name == name })
	rd := budgetReads(g, viewToBe)[i]
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadShow(b *testing.B)    { benchRead(b, "show") }
func BenchmarkReadLs(b *testing.B)      { benchRead(b, "ls") }
func BenchmarkReadTree(b *testing.B)    { benchRead(b, "tree") }
func BenchmarkReadWhy(b *testing.B)     { benchRead(b, "why") }
func BenchmarkReadLog(b *testing.B)     { benchRead(b, "log") }
func BenchmarkReadRoadmap(b *testing.B) { benchRead(b, "roadmap") }
func BenchmarkReadReady(b *testing.B)   { benchRead(b, "ready") }

// BenchmarkTreeWhole shows the whole capability tree: not budgeted (D14).
func BenchmarkTreeWhole(b *testing.B) {
	g := defaultModel(b)
	rd := budgetRead{"tree", "tree", TreeRequest{View: viewToBe}}
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkShowWhatIf shows a deep node in the what-if of the draft change
// set: not budgeted (D14).
func BenchmarkShowWhatIf(b *testing.B) {
	g := defaultModel(b)
	s := g.summary
	rd := budgetRead{"show", "show", ShowRequest{Ref: s.Deep.Handle, View: viewWhatIfPref + s.Draft.String()}}
	for b.Loop() {
		if _, err := runRead(g, rd); err != nil {
			b.Fatal(err)
		}
	}
}

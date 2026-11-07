package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/rootguard"
	runpkg "github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(rootguard.Main(func() int {
		return testscript.RunMain(m, map[string]func() int{
			"rhiza": func() int {
				return profiled(os.Getenv("RHIZA_CPUPROFILE"), func() int {
					return run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Getenv)
				})
			},
		})
	}))
}

// -update rewrites a script's archive file when a `cmp stdout <file>` fails,
// e.g. describe-golden.txtar's golden after a request or response type
// changes: go test ./cmd/rhiza -run TestScripts/describe-golden -update
var update = flag.Bool("update", false, "rewrite the archive files that cmp compares against")

func TestScripts(t *testing.T) {
	dir := "testdata/script"
	files, err := shardFiles(dir, os.Getenv("RHIZA_TEST_SHARD"))
	if err != nil {
		t.Fatal(err)
	}
	// testscript reads Files only when Dir is empty.
	if files != nil {
		dir = ""
	}
	testscript.Run(t, testscript.Params{
		Dir:           dir,
		Files:         files,
		Setup:         setup,
		Cmds:          map[string]func(ts *testscript.TestScript, neg bool, args []string){"fold-ok": foldOK, "chtimes": chtimes, "health-server": healthServer, "run-lock": runLock, "root-daemon-lock": rootDaemonLock},
		UpdateScripts: *update,
	})
}

// shardFiles lists the scripts in dir that shard (RHIZA_TEST_SHARD, "i/n",
// 1-based) selects, as paths for testscript's Files. An empty shard returns
// nil, and testscript then runs every script in Dir.
func shardFiles(dir, shard string) ([]string, error) {
	if shard == "" {
		return nil, nil
	}
	names, err := filepath.Glob(filepath.Join(dir, "*.txtar"))
	if err != nil {
		return nil, err
	}
	return shardNames(names, shard)
}

// shardNames returns the names whose index in the sorted list, mod n, is
// i-1, for shard "i/n": round-robin, so the shards are disjoint and together
// cover every name.
func shardNames(names []string, shard string) ([]string, error) {
	var i, n int
	var rest string
	k, err := fmt.Sscanf(shard, "%d/%d%s", &i, &n, &rest)
	_ = err
	if k != 2 || n < 1 || i < 1 || i > n {
		return nil, fmt.Errorf("RHIZA_TEST_SHARD=%q: want i/n with 1 <= i <= n", shard)
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	out := []string{}
	for idx, name := range sorted {
		if idx%n == i-1 {
			out = append(out, name)
		}
	}
	return out, nil
}

// Every script starts with a fixed clock and ID seed, so its output bytes
// are stable. A script that wants the real clock or random IDs clears them
// with `env RHIZA_TEST_NOW=` and `env RHIZA_TEST_SEED=`.
const (
	testNow  = "2026-10-24T12:00:00Z"
	testSeed = "1"
	// testFreeBytes is the free disk every script sees (1 TiB), so the
	// daemon's minFree guard never reads the host's disk. A script that
	// wants another figure sets RHIZA_TEST_FREE_BYTES itself.
	testFreeBytes = "1099511627776"
)

func setup(env *testscript.Env) error {
	env.Setenv("RHIZA_TEST_NOW", testNow)
	env.Setenv("RHIZA_TEST_SEED", testSeed)
	env.Setenv("RHIZA_TEST_FREE_BYTES", testFreeBytes)
	// RHIZA_TEST_BIN is the directory holding the in-process rhiza, so a
	// script can set PATH to it alone and hide the host's programs (e.g.
	// systemctl) from rhiza.
	if bin, err := exec.LookPath("rhiza"); err == nil {
		env.Setenv("RHIZA_TEST_BIN", filepath.Dir(bin))
	}
	// Every script's $WORK is its own Rhiza root with a project named repo
	// at $WORK/repo (written directly, not through rhiza project add).
	// $NOROOT is a folder outside any root, for the outside-a-root cases.
	err := os.WriteFile(filepath.Join(env.WorkDir, "rhiza.toml"),
		[]byte("[projects.repo]\npath = \"repo\"\n"), 0o644)
	if err != nil {
		return err
	}
	noroot, err := os.MkdirTemp("", "rhiza-noroot-")
	if err != nil {
		return err
	}
	env.Defer(func() { os.RemoveAll(noroot) })
	env.Setenv("NOROOT", noroot)
	return hermeticGit(env)
}

// foldOK is the script command `fold-ok [dir]`: it reads the .rhiza folder
// dir (default $WORK/.rhiza) with store.ReadLog, folds it with
// changeset.FoldLog and prints every log and fold problem to stdout, one per
// line. It fails if there is any; `! fold-ok` fails if there is none.
// chtimes sets the modification time of files: chtimes <RFC 3339 time> file...
func chtimes(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) < 2 {
		ts.Fatalf("usage: chtimes time file...")
	}
	t, err := time.Parse(time.RFC3339, args[0])
	if err != nil {
		ts.Fatalf("chtimes: %v", err)
	}
	for _, f := range args[1:] {
		ts.Check(os.Chtimes(ts.MkAbs(f), t, t))
	}
}

// runLock is the script command `run-lock dir id`: it holds the run lock of
// run id in the runs directory dir until the script ends, so a fixture
// journal there counts as a live run.
func runLock(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 2 {
		ts.Fatalf("usage: run-lock dir id")
	}
	unlock, ok, err := runpkg.NewStore(ts.MkAbs(args[0])).TryLock(args[1])
	if err != nil || !ok {
		ts.Fatalf("run-lock %s: held=%v err=%v", args[1], !ok, err)
	}
	ts.Defer(unlock)
}

// rootDaemonLock is the script command `root-daemon-lock <root>`: it takes
// the root daemon's lock in the test process, so every rhiza the script
// runs sees a root daemon running, until the script ends.
func rootDaemonLock(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 1 {
		ts.Fatalf("usage: root-daemon-lock root")
	}
	l, err := daemon.TryLockRoot(ts.MkAbs(args[0]), time.Now())
	if err != nil {
		ts.Fatalf("root-daemon-lock: %v", err)
	}
	ts.Defer(l.Unlock)
}

// healthServer is the script command `health-server <var> <log> <status>...`:
// it starts a fake deploy health server on loopback, in the test process, so
// it outlives each rhiza the script runs, and sets <var> to its URL. The
// server answers the statuses in order, the last one from then on, and
// appends each status it answers to <log>, one per line. It closes when the
// script ends.
func healthServer(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) < 3 {
		ts.Fatalf("usage: health-server var log status...")
	}
	var statuses []int
	for _, a := range args[2:] {
		n, err := strconv.Atoi(a)
		if err != nil {
			ts.Fatalf("health-server: status %q: %v", a, err)
		}
		statuses = append(statuses, n)
	}
	log := ts.MkAbs(args[1])
	var mu sync.Mutex
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		code := statuses[min(hits, len(statuses)-1)]
		hits++
		if f, err := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			fmt.Fprintln(f, code)
			f.Close()
		}
		mu.Unlock()
		w.WriteHeader(code)
	}))
	ts.Defer(srv.Close)
	ts.Setenv(args[0], srv.URL+"/health")
}

func foldOK(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) > 1 {
		ts.Fatalf("usage: fold-ok [dir]")
	}
	dir := ".rhiza"
	if len(args) == 1 {
		dir = args[0]
	}
	dir = ts.MkAbs(dir)
	var problems []string
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		problems = append(problems, fmt.Sprintf("%s is not a folder", dir))
	}
	l := store.ReadLog(os.DirFS(dir))
	_, foldProblems := changeset.FoldLog(l)
	for _, p := range l.Problems {
		problems = append(problems, fmt.Sprintf("log: %s: %v", p.Path, p.Err))
	}
	for _, p := range foldProblems {
		problems = append(problems, "fold: "+p.String())
	}
	for _, p := range problems {
		fmt.Fprintln(ts.Stdout(), p)
	}
	switch {
	case neg && len(problems) == 0:
		ts.Fatalf("fold-ok: %d entries, no problems", len(l.Entries))
	case !neg && len(problems) > 0:
		ts.Fatalf("fold-ok: %d problems", len(problems))
	}
}

// testGitConfig is the global git config every script runs with.
const testGitConfig = `[user]
	name = Rhiza Test
	email = test@rhiza.invalid
[init]
	defaultBranch = main
[commit]
	gpgsign = false
`

// hermeticGit points git at a sandboxed HOME and global config so scripts
// never read the developer's real configuration. XDG_CACHE_HOME goes there
// too, so os.UserCacheDir, and with it every script's index, is under the
// script's work dir on linux as on darwin (where it is $HOME/Library/Caches)
// and never in the real user cache. A missing git binary is
// left to fail the scripts that use it rather than skipping them.
func hermeticGit(env *testscript.Env) error {
	home := filepath.Join(env.WorkDir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		return err
	}
	gitconfig := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(gitconfig, []byte(testGitConfig), 0o644); err != nil {
		return err
	}
	env.Setenv("HOME", home)
	env.Setenv("XDG_CACHE_HOME", filepath.Join(home, ".cache"))
	env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	env.Setenv("GIT_CONFIG_GLOBAL", gitconfig)
	return nil
}

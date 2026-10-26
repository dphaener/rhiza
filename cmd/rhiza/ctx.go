package main

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
)

// newCtx fills an api.Ctx from the global flags and the environment. Flags
// win over RHIZA_CHANGE, RHIZA_ACTOR and RHIZA_FOR; the workspace
// resolves the raw refs. The index path is index.DefaultPath of the root,
// under the user cache dir, when a root was found; otherwise, or when the
// cache dir is unknown, it is empty and reads fail. RHIZA_TEST_NOW
// (RFC 3339) and RHIZA_TEST_SEED fix the clock and ID seed for tests.
// Each process starts afresh, so both are offset by the repo's entry files,
// plus one for repo.json, plus the git repository's session records: every
// run draws fresh, deterministic IDs, and a dry run predicts the next real
// run.
func newCtx(g globals, stdin io.Reader, getenv func(string) string) (*api.Ctx, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working directory: %w", err)
	}
	ctx := &api.Ctx{
		Root:      findRoot(wd),
		ChangeRef: pick(g.change, g.changeSet, getenv("RHIZA_CHANGE")),
		ActorRef:  pick(g.actor, g.actorSet, getenv("RHIZA_ACTOR")),
		ForRef:    pick(g.forRef, g.forSet, getenv("RHIZA_FOR")),
		Record:    g.record,
		DryRun:    g.dryRun,
		Now:       time.Now,
		IDs:       id.Default(),
		GitConfig: gitConfig,
		Version:   versionString(),
		Stdin:     stdin,
		Home:      getenv("HOME"),
		Exec:      execIn,
		Self:      selfPath(),
		Start:     startProgram,
		// REX_SESSION and REX_BLOCK are the Rex session and block this
		// terminal runs in; Rex sets them in every shell it starts.
		RexSession: getenv("REX_SESSION"),
		RexBlock:   getenv("REX_BLOCK"),
		Stderr:     os.Stderr,
	}
	ctx.Git = gitIn(ctx.Root)
	if _, ok := api.FindRoot(wd); ok {
		if p, err := index.DefaultPath(ctx.Root); err == nil {
			ctx.IndexPath = p
		}
	}
	nowEnv, seedEnv := getenv("RHIZA_TEST_NOW"), getenv("RHIZA_TEST_SEED")
	if nowEnv == "" && seedEnv == "" {
		return ctx, nil
	}
	n := testOffset(ctx.Root)
	if nowEnv != "" {
		t, err := time.Parse(time.RFC3339, nowEnv)
		if err != nil {
			return nil, api.Usagef("RHIZA_TEST_NOW: %v", err)
		}
		t = t.Add(time.Duration(n) * time.Second)
		ctx.Now = func() time.Time { return t }
	}
	if seedEnv == "" {
		ctx.IDs = id.NewGenerator(ctx.Now, crand.Reader)
		return ctx, nil
	}
	seed, err := strconv.ParseInt(seedEnv, 10, 64)
	if err != nil {
		return nil, api.Usagef("RHIZA_TEST_SEED: not an integer: %q", seedEnv)
	}
	ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(seed+int64(n))))
	return ctx, nil
}

// testOffset is N for test mode: the entry files in the log, plus one if
// .rhiza/repo.json exists, plus the session records. A session start writes
// its entries into another worktree's log, so without the records two
// launches in a row would draw the same IDs.
func testOffset(root string) int {
	n := countEntryFiles(root) + countSessionFiles(root)
	if fi, err := os.Stat(filepath.Join(root, store.Dir, store.RepoFile)); err == nil && !fi.IsDir() {
		n++
	}
	return n
}

// countEntryFiles counts the *.json files in root/.rhiza/log/*/, skipping
// names that start with "."; 0 if there is no log.
func countEntryFiles(root string) int {
	matches, _ := filepath.Glob(filepath.Join(root, store.Dir, store.LogDir, "*", "*.json"))
	n := 0
	for _, m := range matches {
		if !strings.HasPrefix(filepath.Base(m), ".") && !strings.HasPrefix(filepath.Base(filepath.Dir(m)), ".") {
			n++
		}
	}
	return n
}

// countSessionFiles counts the *.json files in the session folder of
// root's git repository, skipping names that start with "."; 0 if there is
// no git repository or no session.
func countSessionFiles(root string) int {
	out, err := gitIn(root)("rev-parse", "--git-common-dir")
	if err != nil {
		return 0
	}
	dir := strings.TrimSpace(out)
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(root, dir)
	}
	matches, _ := filepath.Glob(filepath.Join(session.Dir(dir), "*.json"))
	n := 0
	for _, m := range matches {
		if !strings.HasPrefix(filepath.Base(m), ".") {
			n++
		}
	}
	return n
}

func pick(flag string, set bool, env string) string {
	if set {
		return flag
	}
	return env
}

// findRoot is api.FindRoot, falling back to dir when no ancestor holds
// .rhiza: init works there, and operations needing a repo report
// not_found themselves.
func findRoot(dir string) string {
	if root, ok := api.FindRoot(dir); ok {
		return root
	}
	return dir
}

// gitIn returns an api.Ctx.Git that runs `git -C root <args...>` and returns
// its standard output. A failure carries git's trimmed standard error.
func gitIn(root string) func(args ...string) (string, error) {
	return func(args ...string) (string, error) {
		out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
		if err != nil {
			var ee *exec.ExitError
			if errors.As(err, &ee) && len(bytes.TrimSpace(ee.Stderr)) > 0 {
				return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), bytes.TrimSpace(ee.Stderr))
			}
			return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
		return string(out), nil
	}
}

// execIn is api.Ctx.Exec: it runs name with args in dir, with env (nil for
// this process's environment), and returns its standard output. A failure
// carries the program's trimmed standard error.
func execIn(name string, args []string, dir string, env []string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = env
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if msg := bytes.TrimSpace(stderr.Bytes()); len(msg) > 0 {
			return out, fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), msg)
		}
		return out, fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return out, nil
}

// streamGrace is how long a killed program's output stays open for a child
// of it that still holds it.
const streamGrace = 2 * time.Second

// startProgram is api.Ctx.Start, over os/exec.
func startProgram(p api.Program) (api.Process, error) {
	cmd := exec.Command(p.Name, p.Args...)
	cmd.Dir = p.Dir
	set := p.Env
	if p.Dir != "" {
		// As a shell would: the inherited PWD names another folder.
		if abs, err := filepath.Abs(p.Dir); err == nil {
			set = append([]string{"PWD=" + abs}, set...)
		}
	}
	cmd.Env = programEnv(os.Environ(), set, p.Unset)
	proc := &process{cmd: cmd}
	if p.Terminal {
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		// The terminal sends its interrupt, quit and hangup to the whole
		// foreground process group. Catching them (not ignoring them,
		// which the program would inherit) leaves them to the program,
		// and keeps this process alive to see how it ends.
		proc.sigs = make(chan os.Signal, 1)
		signal.Notify(proc.sigs, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
		if err := cmd.Start(); err != nil {
			signal.Stop(proc.sigs)
			return nil, err
		}
		return proc, nil
	}
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	cmd.Stdout = w
	if err := cmd.Start(); err != nil {
		r.Close()
		w.Close()
		return nil, err
	}
	w.Close()
	proc.out = r
	return proc, nil
}

// programEnv is env without the variables named in unset or set, followed
// by set.
func programEnv(env, set, unset []string) []string {
	drop := map[string]bool{}
	for _, n := range unset {
		drop[n] = true
	}
	for _, kv := range set {
		n, _, _ := strings.Cut(kv, "=")
		drop[n] = true
	}
	out := make([]string, 0, len(env)+len(set))
	for _, kv := range env {
		if n, _, _ := strings.Cut(kv, "="); !drop[n] {
			out = append(out, kv)
		}
	}
	return append(out, set...)
}

// process is an api.Process over an *exec.Cmd.
type process struct {
	cmd  *exec.Cmd
	out  *os.File
	sigs chan os.Signal
}

func (p *process) Output() io.Reader {
	if p.out == nil {
		return nil
	}
	return p.out
}

func (p *process) Wait() (int, error) {
	err := p.cmd.Wait()
	if p.sigs != nil {
		signal.Stop(p.sigs)
	}
	if p.out != nil {
		p.out.Close()
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return 128 + int(ws.Signal()), nil
		}
		return ee.ExitCode(), nil
	}
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func (p *process) Kill() {
	p.cmd.Process.Kill()
	if p.out != nil {
		time.AfterFunc(streamGrace, func() { p.out.Close() })
	}
}

// selfPath is the absolute path of the running rhiza: os.Args[0], looked
// up on PATH when it has no path separator (as a shell runs it), then made
// absolute. Empty if neither works.
func selfPath() string {
	p := os.Args[0]
	if !strings.ContainsRune(p, filepath.Separator) {
		lp, err := exec.LookPath(p)
		if err != nil {
			return ""
		}
		p = lp
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return ""
	}
	return abs
}

// gitConfig returns the trimmed value of a git config key, or an error if
// it is unset or empty.
func gitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", "--get", key).Output()
	v := strings.TrimSpace(string(out))
	if err != nil || v == "" {
		return "", fmt.Errorf("git config %s is not set", key)
	}
	return v, nil
}

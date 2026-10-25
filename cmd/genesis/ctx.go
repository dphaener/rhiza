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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/genesis/internal/api"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/store"
)

// newCtx fills an api.Ctx from the global flags and the environment. Flags
// win over GENESIS_CHANGE, GENESIS_ACTOR and GENESIS_FOR; the workspace
// resolves the raw refs. GENESIS_TEST_NOW (RFC 3339) and GENESIS_TEST_SEED
// fix the clock and ID seed for tests. Each process starts afresh, so both
// are offset by the repo's entry files plus one for repo.json: every run
// draws fresh, deterministic IDs, and a dry run predicts the next real run.
func newCtx(g globals, stdin io.Reader, getenv func(string) string) (*api.Ctx, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working directory: %w", err)
	}
	ctx := &api.Ctx{
		Root:      findRoot(wd),
		ChangeRef: pick(g.change, g.changeSet, getenv("GENESIS_CHANGE")),
		ActorRef:  pick(g.actor, g.actorSet, getenv("GENESIS_ACTOR")),
		ForRef:    pick(g.forRef, g.forSet, getenv("GENESIS_FOR")),
		Record:    g.record,
		DryRun:    g.dryRun,
		Now:       time.Now,
		IDs:       id.Default(),
		GitConfig: gitConfig,
		Version:   versionString(),
		Stdin:     stdin,
	}
	ctx.Git = gitIn(ctx.Root)
	if home, err := os.UserHomeDir(); err == nil {
		ctx.Home = home
	}
	nowEnv, seedEnv := getenv("GENESIS_TEST_NOW"), getenv("GENESIS_TEST_SEED")
	if nowEnv == "" && seedEnv == "" {
		return ctx, nil
	}
	n := testOffset(ctx.Root)
	if nowEnv != "" {
		t, err := time.Parse(time.RFC3339, nowEnv)
		if err != nil {
			return nil, api.Usagef("GENESIS_TEST_NOW: %v", err)
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
		return nil, api.Usagef("GENESIS_TEST_SEED: not an integer: %q", seedEnv)
	}
	ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(seed+int64(n))))
	return ctx, nil
}

// testOffset is N for test mode: the entry files in the log, plus one if
// .genesis/repo.json exists.
func testOffset(root string) int {
	n := countEntryFiles(root)
	if fi, err := os.Stat(filepath.Join(root, store.Dir, store.RepoFile)); err == nil && !fi.IsDir() {
		n++
	}
	return n
}

// countEntryFiles counts the *.json files in root/.genesis/log/*/, skipping
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

func pick(flag string, set bool, env string) string {
	if set {
		return flag
	}
	return env
}

// findRoot is api.FindRoot, falling back to dir when no ancestor holds
// .genesis: init works there, and operations needing a repo report
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

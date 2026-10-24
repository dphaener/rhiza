package main

import (
	crand "crypto/rand"
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
)

// newCtx fills an api.Ctx from the global flags and the environment. Flags
// win over GENESIS_CHANGE, GENESIS_ACTOR and GENESIS_FOR. The refs stay
// raw; the workspace resolves them.
//
// GENESIS_TEST_NOW (RFC 3339) fixes the clock and GENESIS_TEST_SEED (an
// integer) seeds the ID entropy, so tests get stable output. They are for
// tests only. Every genesis process starts afresh, so a plain fixed clock
// and seed would make each run generate the same IDs as the last and a
// second write would collide. In test mode both are therefore offset by N,
// the number of entry files already in the repo's log: the clock reads
// GENESIS_TEST_NOW + N seconds and the seed is GENESIS_TEST_SEED + N. That
// stays deterministic, and every write changes N for the next run.
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
		Stdin:     stdin,
	}
	nowEnv, seedEnv := getenv("GENESIS_TEST_NOW"), getenv("GENESIS_TEST_SEED")
	if nowEnv == "" && seedEnv == "" {
		return ctx, nil
	}
	n := countEntryFiles(ctx.Root)
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

// countEntryFiles counts the *.json files in root/.genesis/log/*/, skipping
// names that start with "."; 0 if there is no log.
func countEntryFiles(root string) int {
	matches, _ := filepath.Glob(filepath.Join(root, ".genesis", "log", "*", "*.json"))
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

// findRoot walks up from dir to the first folder holding .genesis. With
// none, it returns dir: init works there, and operations needing a repo
// report not_found themselves.
func findRoot(dir string) string {
	for d := dir; ; {
		if fi, err := os.Stat(filepath.Join(d, ".genesis")); err == nil && fi.IsDir() {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return dir
		}
		d = parent
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

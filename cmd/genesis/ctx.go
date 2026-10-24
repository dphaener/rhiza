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
// tests only.
func newCtx(g globals, stdin io.Reader, getenv func(string) string) (*api.Ctx, error) {
	ctx := &api.Ctx{
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
	fixed := false
	if s := getenv("GENESIS_TEST_NOW"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, api.Usagef("GENESIS_TEST_NOW: %v", err)
		}
		ctx.Now = func() time.Time { return t }
		fixed = true
	}
	if s := getenv("GENESIS_TEST_SEED"); s != "" {
		seed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, api.Usagef("GENESIS_TEST_SEED: not an integer: %q", s)
		}
		ctx.IDs = id.NewGenerator(ctx.Now, rand.New(rand.NewSource(seed)))
	} else if fixed {
		ctx.IDs = id.NewGenerator(ctx.Now, crand.Reader)
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("working directory: %w", err)
	}
	ctx.Root = findRoot(wd)
	return ctx, nil
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

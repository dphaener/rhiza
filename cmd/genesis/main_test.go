package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/store"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"genesis": main,
	})
}

// -update rewrites a script's archive file when a `cmp stdout <file>` fails,
// e.g. describe-golden.txtar's golden after a request or response type
// changes: go test ./cmd/genesis -run TestScripts/describe-golden -update
var update = flag.Bool("update", false, "rewrite the archive files that cmp compares against")

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:           "testdata/script",
		Setup:         setup,
		Cmds:          map[string]func(ts *testscript.TestScript, neg bool, args []string){"fold-ok": foldOK},
		UpdateScripts: *update,
	})
}

// Every script starts with a fixed clock and ID seed, so its output bytes
// are stable. A script that wants the real clock or random IDs clears them
// with `env GENESIS_TEST_NOW=` and `env GENESIS_TEST_SEED=`.
const (
	testNow  = "2026-10-24T12:00:00Z"
	testSeed = "1"
)

func setup(env *testscript.Env) error {
	env.Setenv("GENESIS_TEST_NOW", testNow)
	env.Setenv("GENESIS_TEST_SEED", testSeed)
	return hermeticGit(env)
}

// foldOK is the script command `fold-ok [dir]`: it reads the .genesis folder
// dir (default $WORK/.genesis) with store.ReadLog, folds it with
// changeset.FoldLog and prints every log and fold problem to stdout, one per
// line. It fails if there is any; `! fold-ok` fails if there is none.
func foldOK(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) > 1 {
		ts.Fatalf("usage: fold-ok [dir]")
	}
	dir := ".genesis"
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
	name = Genesis Test
	email = test@genesis.invalid
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

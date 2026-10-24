package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"genesis": main,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:   "testdata/script",
		Setup: hermeticGit,
	})
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
// never read the developer's real configuration. A missing git binary is
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
	env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	env.Setenv("GIT_CONFIG_GLOBAL", gitconfig)
	return nil
}

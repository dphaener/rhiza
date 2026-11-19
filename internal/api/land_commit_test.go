package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// landCommitFixture is an accepted change set in a real git repository
// with everything committed but the change set's own entries, a file
// staged outside .rhiza, and a head build (SHIP_HEAD_RHIZA) that is not the
// running binary. Its ctx.Exec runs git hermetically and fakes the head:
// it records the head's argv and runs the op in-process with an
// environment naming no head, or fails land when failLand is set.
func landCommitFixture(t *testing.T, failLand bool) (ctx *Ctx, git func(...string) string, argv *[][]string) {
	t.Helper()
	ctx = initCtx(t)
	newChange(t, ctx, `{"intent":"Land and commit"}`)
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"landing"}`)
	for _, v := range []string{"propose", "accept"} {
		if _, err := callOp(t, ctx, v, `{}`); err != nil {
			t.Fatalf("%s: %v", v, err)
		}
	}
	gitExec := hermeticExec(t, t.TempDir())
	git = func(args ...string) string {
		t.Helper()
		out, err := gitExec("git", args, ctx.Root, nil)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	ctx.Git = func(args ...string) (string, error) {
		out, err := gitExec("git", args, ctx.Root, nil)
		return string(out), err
	}
	git("init", "-q")
	git("add", "-A", "--", ":!.rhiza/log")
	git("commit", "-q", "-m", "base")
	if err := os.WriteFile(filepath.Join(ctx.Root, "code.go"), []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "code.go")

	bin := filepath.Join(t.TempDir(), "rhiza-head")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx.Self = filepath.Join(t.TempDir(), "rhiza-installed")
	parent := ctx.Getenv
	ctx.Getenv = func(k string) string {
		if k == envShipRhiza {
			return bin
		}
		return parent(k)
	}
	argv = new([][]string)
	ctx.Exec = func(name string, args []string, dir string, _ []string) ([]byte, error) {
		if name != bin {
			return gitExec(name, args, dir, nil)
		}
		*argv = append(*argv, append([]string{}, args...))
		if failLand && args[0] == "land" {
			return []byte(`{"error":{"code":"conflict","message":"head refused"}}`), errors.New("exit status 3")
		}
		child := *ctx
		child.Getenv = parent
		resp, err := callOp(t, &child, args[0], `{}`)
		if err != nil {
			return nil, err
		}
		return json.Marshal(resp)
	}
	return ctx, git, argv
}

// With a head build named, land --commit lands and renders through the
// head, never forwarding --commit, then commits only .rhiza itself with
// the Change: and Actor: trailers; a failed head land commits nothing.
func TestLandCommitViaHead(t *testing.T) {
	t.Parallel()
	t.Run("lands and renders through the head, then commits .rhiza", func(t *testing.T) {
		ctx, git, argv := landCommitFixture(t, false)
		cs := currentOf(t, ctx)
		before := git("rev-parse", "HEAD")
		out, err := callOp(t, ctx, "land", `{"commit":true}`)
		if err != nil {
			t.Fatal(err)
		}
		resp := out.(LifecycleResponse)
		if len(*argv) != 2 || (*argv)[0][0] != "land" || (*argv)[1][0] != "render" {
			t.Fatalf("head argv = %v, want land then render", *argv)
		}
		for _, a := range *argv {
			if slices.Contains(a, "--commit") {
				t.Errorf("head argv %v forwards --commit", a)
			}
		}
		head := git("rev-parse", "HEAD")
		if head == before || resp.Commit != head || git("rev-parse", "HEAD~1") != before {
			t.Fatalf("HEAD %s (was %s), resp.Commit %s: want one new commit", head, before, resp.Commit)
		}
		files := strings.Fields(git("show", "--name-only", "--format=", "HEAD"))
		if len(files) == 0 {
			t.Fatal("the commit has no files")
		}
		for _, f := range files {
			if !strings.HasPrefix(f, ".rhiza/") {
				t.Errorf("the commit has %s, outside .rhiza", f)
			}
		}
		if !slices.ContainsFunc(files, func(f string) bool { return strings.HasPrefix(f, ".rhiza/rendered/") }) {
			t.Errorf("the commit has no rendered file: %v", files)
		}
		w := ctxWorkspace(t, ctx)
		actor, err := shipActorHandle(w, cs, ctx)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := git("log", "-1", "--format=%s"), landCommitMessage; got != want {
			t.Errorf("subject = %q, want %q", got, want)
		}
		if got, want := git("log", "-1", "--format=%(trailers)"), "Change: "+cs.String()+"\nActor: "+actor; got != want {
			t.Errorf("trailers = %q, want %q", got, want)
		}
		if s := git("status", "--porcelain", "--", ".rhiza"); s != "" {
			t.Errorf(".rhiza is not clean:\n%s", s)
		}
		if s := git("diff", "--cached", "--name-only"); s != "code.go" {
			t.Errorf("staged = %q, want code.go still staged", s)
		}
	})

	t.Run("a failed head land commits nothing", func(t *testing.T) {
		ctx, git, argv := landCommitFixture(t, true)
		before := git("rev-parse", "HEAD")
		if _, err := callOp(t, ctx, "land", `{"commit":true}`); err == nil {
			t.Fatal("land --commit succeeded on a failed head land")
		}
		if len(*argv) != 1 || (*argv)[0][0] != "land" {
			t.Errorf("head argv = %v, want land only", *argv)
		}
		if head := git("rev-parse", "HEAD"); head != before {
			t.Errorf("HEAD moved to %s", head)
		}
		if s := git("diff", "--cached", "--name-only"); s != "code.go" {
			t.Errorf("staged = %q, want code.go only", s)
		}
	})
}

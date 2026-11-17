package api

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/store"
)

// testCtx is a Ctx over root with a fixed clock, seeded IDs, a git config
// whose user.email is test@example.com, an Environ holding the process's
// PATH (children are found on Environ's PATH alone) and, as in production, an index:
// IndexPath is a fresh file under t's own t.TempDir(), so no two tests (or
// two Ctxs of one test) share an index. A test that exercises the no-index
// fallback on purpose clears IndexPath, with a comment saying so.
func testCtx(t testing.TB, root string) *Ctx {
	t.Helper()
	ctx := testCtxNoIndex(root)
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	guardProjections(t, ctx)
	return ctx
}

// guardProjections fails t, once it and its cleanups registered after this
// one are done, when a projection opened on ctx (or a copy of it) is still
// open, or a throwaway index's folder is left in ctx's TMPDIR: a fresh
// folder under t's own t.TempDir(), so no other test's land there.
func guardProjections(t testing.TB, ctx *Ctx) {
	t.Helper()
	ctx.liveProjections = &liveProjections{}
	tmp := t.TempDir()
	inner := ctx.Getenv
	ctx.Getenv = func(k string) string {
		if k == "TMPDIR" {
			return tmp
		}
		return inner(k)
	}
	t.Cleanup(func() {
		for _, s := range ctx.liveProjections.stacks() {
			t.Errorf("projection left open, opened at:\n%s", s)
		}
		left, err := filepath.Glob(filepath.Join(tmp, "rhiza-projection-*"))
		if err != nil {
			t.Error(err)
		}
		if len(left) > 0 {
			t.Errorf("throwaway indexes left: %v", left)
		}
	})
}

// testCtxNoIndex is testCtx with no index (IndexPath empty), for code with
// no test to put one under: a snapshot's build and copy (snapshot.copy gives
// the copy's Ctx its test's index) and the budget reads (runRead sets the
// generated repository's own index).
func testCtxNoIndex(root string) *Ctx {
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	return &Ctx{
		Root:      root,
		Now:       clock,
		IDs:       id.NewGenerator(clock, rand.New(rand.NewSource(3))),
		GitConfig: gitConfig(map[string]string{"user.email": "test@example.com", "user.name": "Test User"}),
		Environ:   []string{"PATH=" + os.Getenv("PATH")},
		Getenv:    testCache,
	}
}

// testCache is a test Ctx's getenv: the test binary's user cache folder
// (XDG_CACHE_HOME, set once in TestMain) and nothing else, so an index or
// run folder lands there, as it did when a nil getenv read the process's.
func testCache(key string) string {
	if key == "XDG_CACHE_HOME" {
		return os.Getenv("XDG_CACHE_HOME")
	}
	return ""
}

// callOp decodes reqJSON for the named operation and runs it on ctx.
func callOp(t *testing.T, ctx *Ctx, name, reqJSON string) (any, error) {
	t.Helper()
	o, ok := Lookup(name)
	if !ok {
		t.Fatalf("%s is not registered", name)
	}
	req, err := DecodeRequest(o, []byte(reqJSON))
	if err != nil {
		t.Fatalf("DecodeRequest(%s, %s): %v", name, reqJSON, err)
	}
	return o.Run(ctx, req)
}

// human renders a response's markdown.
func human(t *testing.T, resp any) string {
	t.Helper()
	h, ok := resp.(Humaner)
	if !ok {
		t.Fatalf("%T is not a Humaner", resp)
	}
	var buf bytes.Buffer
	h.Human(&buf)
	return buf.String()
}

func TestInit(t *testing.T) {
	t.Run("Fresh", func(t *testing.T) {
		root := t.TempDir()
		resp, err := callOp(t, testCtx(t, root), "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		r := resp.(InitResponse)
		st, err := store.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		if r.Repo != st.Repo().ID || r.Repo.IsZero() {
			t.Errorf("repo = %s, want %s", r.Repo, st.Repo().ID)
		}
		if r.Root != root {
			t.Errorf("root = %q, want %q", r.Root, root)
		}
		data, err := os.ReadFile(filepath.Join(root, store.Dir, ".gitignore"))
		if err != nil || string(data) != "local/\n" {
			t.Errorf(".gitignore = %q, %v; want \"local/\\n\"", data, err)
		}
		if data, err := os.ReadFile(filepath.Join(root, store.Dir, store.LogDir, ".gitkeep")); err != nil || len(data) != 0 {
			t.Errorf("log/.gitkeep = %q, %v; want an empty file", data, err)
		}
		out := human(t, resp)
		for _, s := range []string{root, r.Repo.String(), `rhiza change new "<intent>"`} {
			if !strings.Contains(out, s) {
				t.Errorf("human output lacks %q:\n%s", s, out)
			}
		}
	})

	t.Run("Exists", func(t *testing.T) {
		root := t.TempDir()
		if _, err := callOp(t, testCtx(t, root), "init", `{}`); err != nil {
			t.Fatal(err)
		}
		_, err := callOp(t, testCtx(t, root), "init", `{}`)
		if ExitCode(err) != 3 {
			t.Fatalf("second init: %v, exit %d; want 3", err, ExitCode(err))
		}
		if !strings.Contains(err.Error(), "already") {
			t.Errorf("message = %q", err)
		}
	})

	t.Run("RelativeRootIsMadeAbsolute", func(t *testing.T) {
		// A relative root is made absolute against the Ctx's working
		// directory, not the process's.
		wd := t.TempDir()
		if err := os.Mkdir(filepath.Join(wd, "proj"), 0o755); err != nil {
			t.Fatal(err)
		}
		ctx := testCtx(t, "proj")
		ctx.Wd = wd
		resp, err := callOp(t, ctx, "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := resp.(InitResponse).Root, filepath.Join(wd, "proj"); got != want {
			t.Errorf("root = %q, want %q", got, want)
		}
	})

	t.Run("DryRunIsUsage", func(t *testing.T) {
		root := t.TempDir()
		ctx := testCtx(t, root)
		ctx.DryRun = true
		_, err := callOp(t, ctx, "init", `{}`)
		if ExitCode(err) != 2 {
			t.Fatalf("dry-run init: %v, exit %d; want 2", err, ExitCode(err))
		}
		if _, err := os.Stat(filepath.Join(root, store.Dir)); !os.IsNotExist(err) {
			t.Errorf("dry-run init created .rhiza: %v", err)
		}
	})
}

func TestInitAgent(t *testing.T) {
	t.Parallel()
	def := adapter.Default().Name()
	t.Run("FlagWritten", func(t *testing.T) {
		ctx, conf := projectAgentCtx(t, "")
		resp, err := callOp(t, ctx, "init", `{"agent":"`+askAgent+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		r := resp.(InitResponse)
		if r.Agent != askAgent || !r.Saved || conf[ProjectAgentKey] != askAgent {
			t.Errorf("agent %q saved %v, git config %q; want %s written", r.Agent, r.Saved, conf[ProjectAgentKey], askAgent)
		}
		if out := human(t, resp); !strings.Contains(out, "Agent: "+askAgent+", written to "+ProjectAgentKey) {
			t.Errorf("human output lacks the agent:\n%s", out)
		}
	})

	t.Run("UnknownFlagRefused", func(t *testing.T) {
		ctx, conf := projectAgentCtx(t, "")
		_, err := callOp(t, ctx, "init", `{"agent":"nope"}`)
		if ExitCode(err) != 2 {
			t.Fatalf("init --agent nope: %v, exit %d; want 2", err, ExitCode(err))
		}
		for _, n := range adapter.Names() {
			if !strings.Contains(err.Error(), n) {
				t.Errorf("message %q lacks registered adapter %s", err, n)
			}
		}
		if _, err := os.Stat(filepath.Join(ctx.Root, store.Dir)); !os.IsNotExist(err) {
			t.Errorf("refused init created .rhiza: %v", err)
		}
		if _, ok := conf[ProjectAgentKey]; ok {
			t.Errorf("refused init wrote %s", ProjectAgentKey)
		}
	})

	t.Run("NoTerminalTakesDefault", func(t *testing.T) {
		ctx, conf := projectAgentCtx(t, "")
		// Input with no terminal is never read as an answer.
		ctx.Stdin = strings.NewReader(askAgent + "\n")
		resp, err := callOp(t, ctx, "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(InitResponse); r.Agent != def || conf[ProjectAgentKey] != def {
			t.Errorf("agent %q, git config %q; want the default %s", r.Agent, conf[ProjectAgentKey], def)
		}
	})

	prompt := func(t *testing.T, input string) (InitResponse, string, error) {
		t.Helper()
		ctx, conf := projectAgentCtx(t, "")
		var term bytes.Buffer
		ctx.Terminal, ctx.Stdin = &term, strings.NewReader(input)
		resp, err := callOp(t, ctx, "init", `{}`)
		if err != nil {
			return InitResponse{}, term.String(), err
		}
		r := resp.(InitResponse)
		if conf[ProjectAgentKey] != r.Agent {
			t.Errorf("git config %q, want %q", conf[ProjectAgentKey], r.Agent)
		}
		return r, term.String(), nil
	}

	t.Run("PromptEmptyTakesDefault", func(t *testing.T) {
		r, term, err := prompt(t, "\n")
		if err != nil {
			t.Fatal(err)
		}
		if r.Agent != def {
			t.Errorf("agent %q, want the default %s", r.Agent, def)
		}
		for _, want := range []string{"Which agent does this project use?", "  " + def + " (default)\n", "  " + askAgent + "\n", "Agent [" + def + "]: "} {
			if !strings.Contains(term, want) {
				t.Errorf("prompt lacks %q:\n%s", want, term)
			}
		}
	})

	t.Run("PromptReasksUnknown", func(t *testing.T) {
		r, term, err := prompt(t, "nope\n"+askAgent+"\n")
		if err != nil {
			t.Fatal(err)
		}
		if r.Agent != askAgent {
			t.Errorf("agent %q, want %s", r.Agent, askAgent)
		}
		if !strings.Contains(term, "nope is not a registered agent adapter") || strings.Count(term, "Agent [") != 2 {
			t.Errorf("want one refusal and a second ask:\n%s", term)
		}
	})

	t.Run("PromptChosenStored", func(t *testing.T) {
		r, _, err := prompt(t, " "+askAgent+" \n")
		if err != nil {
			t.Fatal(err)
		}
		if r.Agent != askAgent || !r.Saved {
			t.Errorf("agent %q saved %v, want %s saved", r.Agent, r.Saved, askAgent)
		}
	})

	t.Run("PromptUnknownAtEndRefused", func(t *testing.T) {
		_, _, err := prompt(t, "nope")
		if ExitCode(err) != 2 {
			t.Fatalf("init: %v, exit %d; want 2", err, ExitCode(err))
		}
	})

	t.Run("JSONDoesNotAsk", func(t *testing.T) {
		ctx, _ := projectAgentCtx(t, "")
		var term bytes.Buffer
		ctx.Terminal, ctx.Stdin, ctx.JSON = &term, strings.NewReader(askAgent+"\n"), true
		resp, err := callOp(t, ctx, "init", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(InitResponse); r.Agent != def || term.Len() != 0 {
			t.Errorf("agent %q, prompt %q; want the default and no prompt", r.Agent, term.String())
		}
	})

	t.Run("NoGitWritesNothing", func(t *testing.T) {
		ctx, conf := projectAgentCtx(t, "")
		inner := ctx.Git
		ctx.Git = func(args ...string) (string, error) {
			if args[0] == "rev-parse" {
				return "", errors.New("not a git repository")
			}
			return inner(args...)
		}
		resp, err := callOp(t, ctx, "init", `{"agent":"`+askAgent+`"}`)
		if err != nil {
			t.Fatal(err)
		}
		if r := resp.(InitResponse); r.Saved || len(conf) != 0 {
			t.Errorf("saved %v, git config %v; want nothing written", r.Saved, conf)
		}
		if out := human(t, resp); !strings.Contains(out, "rhiza daemon set "+ProjectAgentKey+" "+askAgent) {
			t.Errorf("human output lacks how to set it later:\n%s", out)
		}
	})
}

// underRoot makes ctx's repo folder a Rhiza root too (a rhiza.toml beside
// its .rhiza) registering the repo as a project, and sets its working directory there, as the daemon needs.
func underRoot(t *testing.T, ctx *Ctx) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(ctx.Root, "rhiza.toml"), []byte("[projects.repo]\npath = \".\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx.Wd = ctx.Root
}

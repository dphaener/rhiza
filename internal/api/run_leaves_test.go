package api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// hermeticExec is a Ctx.Exec that runs programs with git pointed at a
// sandboxed HOME and global config under dir, never the developer's.
func hermeticExec(t *testing.T, dir string) func(string, []string, string, []string) ([]byte, error) {
	t.Helper()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(cfg, []byte("[user]\n\tname = Rhiza Test\n\temail = test@rhiza.invalid\n[init]\n\tdefaultBranch = main\n[commit]\n\tgpgsign = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "HOME="+home, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+cfg)
	return func(name string, args []string, d string, _ []string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		cmd.Dir, cmd.Env = d, env
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			return out, fmt.Errorf("%s %s: %v: %s", name, strings.Join(args, " "), err, stderr.String())
		}
		return out, nil
	}
}

// leafCommitRunner is a stepRunner whose worktree is a real git repository
// with one commit, and a function that runs git there.
func leafCommitRunner(t *testing.T) (*stepRunner, func(...string) string) {
	t.Helper()
	tmp := t.TempDir()
	wt := filepath.Join(tmp, "wt")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(wt)
	ctx.Exec = hermeticExec(t, tmp)
	git := func(args ...string) string {
		t.Helper()
		out, err := ctx.Exec("git", args, wt, nil)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	git("init", "-q")
	git("commit", "-q", "--allow-empty", "-m", "base")
	return &stepRunner{ctx: ctx, worktree: wt}, git
}

// A builder's leaf commit is the newest commit since its start HEAD whose
// Leaf: trailer names its leaf. With no new commit, or new commits naming
// only another leaf or none, there is none, and the detail says why.
func TestLeafCommit(t *testing.T) {
	leaf := run.Leaf{ID: "L1", Handle: "leaf-a"}
	t.Run("none", func(t *testing.T) {
		r, _ := leafCommitRunner(t)
		head, err := r.Head()
		if err != nil {
			t.Fatal(err)
		}
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != "" || detail != "no new commit since "+head[:12] {
			t.Errorf("LeafCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("other leaf", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "work", "-m", "Change: CS1\nLeaf: leaf-b")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != "" || detail != "new commits name another leaf: leaf-b" {
			t.Errorf("LeafCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("no trailer", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "work")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != "" || !strings.Contains(detail, "none with a Leaf: trailer") {
			t.Errorf("LeafCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("own leaf", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "work", "-m", "Leaf: leaf-a")
		want := git("rev-parse", "HEAD")
		git("commit", "-q", "--allow-empty", "-m", "unrelated")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != want || detail != "" {
			t.Errorf("LeafCommit = %q, %q, %v; want %s", c, detail, err, want)
		}
	})
	t.Run("before start does not count", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		git("commit", "-q", "--allow-empty", "-m", "earlier", "-m", "Leaf: leaf-a")
		head, _ := r.Head()
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != "" || !strings.HasPrefix(detail, "no new commit since") {
			t.Errorf("LeafCommit = %q, %q, %v", c, detail, err)
		}
	})
}

// gitLeafRunner is the stepRunner on a real git worktree, its sessions
// faked: each builder, as it starts, makes the commits commits names for
// its leaf, and every session exits 0 at its first poll.
type gitLeafRunner struct {
	*stepRunner
	git     func(...string) string
	leaves  []run.Leaf
	commits map[string][]string // per leaf handle, the commits' messages
	specs   []run.LeafSpec
}

func (g *gitLeafRunner) Leaves(string) ([]run.Leaf, error) { return g.leaves, nil }

func (g *gitLeafRunner) StartLeaf(spec run.LeafSpec) (string, error) {
	g.specs = append(g.specs, spec)
	if spec.Role == "builder" {
		for _, m := range g.commits[spec.Leaf.Handle] {
			g.git("commit", "-q", "--allow-empty", "-m", m)
		}
	}
	return spec.Leaf.Handle + "-" + spec.Role, nil
}

func (g *gitLeafRunner) AgentStatus(string) (run.AgentStatus, error) {
	return run.AgentStatus{Done: true, Outcome: "pass", Detail: "exit 0"}, nil
}

func (g *gitLeafRunner) AgentEnded(workflow.Step, string) error { return nil }

type noSleep struct{}

func (noSleep) Now() time.Time      { return time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC) }
func (noSleep) Sleep(time.Duration) {}

// The engine, against a real git repository, fails a builder that ends
// with no new commit, or with one naming another leaf, with that detail,
// and follows a builder whose commit names its leaf with a verifier given
// that commit, on its leaf-session and in its spec.
func TestLeafCommitEngine(t *testing.T) {
	sr, git := leafCommitRunner(t)
	g := &gitLeafRunner{stepRunner: sr, git: git,
		leaves: []run.Leaf{{ID: "L1", Handle: "leaf-none"}, {ID: "L2", Handle: "leaf-other"}, {ID: "L3", Handle: "leaf-own"}},
		commits: map[string][]string{
			"leaf-other": {"work\n\nLeaf: leaf-x"},
			"leaf-own":   {"work\n\nLeaf: leaf-own"},
		},
	}
	src := "[step.build]\nkind = \"leaves\"\n[step.build.builder]\nprompt = \"b\"\n[step.build.verifier]\nprompt = \"v\"\n"
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	store := run.NewStore(filepath.Join(t.TempDir(), "runs"))
	now := noSleep{}.Now()
	if err := run.Start(store, "R1", run.Started{Workflow: "wf", Definition: src, Change: "CS1"}, now, def); err != nil {
		t.Fatal(err)
	}
	e := &run.Engine{Store: store, Runner: g, Clock: noSleep{}, Now: func() time.Time { return now }}
	st, err := e.Run("R1")
	if err != nil {
		t.Fatal(err)
	}
	byHandle := map[string]run.LeafState{}
	for _, l := range st.Leaves {
		byHandle[l.Handle] = l
	}
	start, _ := byHandle["leaf-none"].Session("builder", 1)
	if j, _ := byHandle["leaf-none"].Judged("builder", 1); j.Outcome != "fail" || j.Detail != "no new commit since "+shortRev(start.Commit) {
		t.Errorf("leaf-none builder = %+v", j)
	}
	if j, _ := byHandle["leaf-other"].Judged("builder", 1); j.Outcome != "fail" || j.Detail != "new commits name another leaf: leaf-x" {
		t.Errorf("leaf-other builder = %+v", j)
	}
	for _, h := range []string{"leaf-none", "leaf-other"} {
		if l := byHandle[h]; l.Outcome != run.LeafFailed {
			t.Errorf("%s = %+v, want failed", h, l)
		} else if _, ok := l.Session("verifier", 1); ok {
			t.Errorf("%s started a verifier", h)
		}
	}
	own := git("rev-parse", "HEAD")
	l := byHandle["leaf-own"]
	if j, _ := l.Judged("builder", 1); j.Outcome != "pass" || j.Commit != own {
		t.Errorf("leaf-own builder = %+v, want pass with %s", j, own)
	}
	if v, ok := l.Session("verifier", 1); !ok || v.Commit != own {
		t.Errorf("leaf-own verifier session = %+v, want commit %s", v, own)
	}
	if l.Outcome != run.LeafPassed {
		t.Errorf("leaf-own = %+v", l)
	}
	spec := g.specs[len(g.specs)-1]
	if spec.Role != "verifier" || spec.Commit != own || !slices.Contains(spec.Env, run.EnvCommit+"="+own) {
		t.Errorf("verifier spec = %+v", spec)
	}
}

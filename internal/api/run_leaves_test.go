package api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
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
	return func(name string, args []string, d string, e []string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		cmd.Dir, cmd.Env = d, env
		if e != nil {
			cmd.Env = append(append([]string{}, e...), env[len(env)-3:]...)
		}
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

// Evidence finds every verifier's own Evidence: its sessions are fakes
// with no record to read.
func (g *gitLeafRunner) Evidence(run.Leaf, string) (bool, string, error) { return true, "", nil }

// CommitEvidence commits nothing: there is no Evidence entry to commit.
func (g *gitLeafRunner) CommitEvidence(leaf run.Leaf, _ string) (string, error) {
	return "E-" + leaf.Handle, nil
}

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

// evidenceCommitRunner is a stepRunner on a real git worktree holding a
// Rhiza repository, committed, whose change set has leaf "leaf-a", and a
// verifier session "S-v" whose actor is ada, recorded and committed too.
// It returns the runner, a git runner there, the context (to write as
// ada), the leaf and ada.
func evidenceCommitRunner(t *testing.T) (*stepRunner, func(...string) string, *Ctx, run.Leaf, id.ID) {
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
	ctx.Git = func(args ...string) (string, error) {
		out, err := ctx.Exec("git", args, wt, nil)
		return string(out), err
	}
	git("init", "-q")
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	cs := newChange(t, ctx, `{"intent":"Build"}`).ChangeSet
	leaf := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"leaf-a"}`).Nodes[0]
	ada := recordAda(t, ctx, "ada@example.com")
	git("add", "-A")
	git("commit", "-q", "-m", "model")
	sessions := session.NewStore(filepath.Join(tmp, "sessions"))
	if err := sessions.Write(session.Record{ID: "S-v", Actor: ada.String(), Change: cs.String()}); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: ctx, change: cs.String(), run: "R1", worktree: wt, sessions: sessions}
	return r, git, ctx, run.Leaf{ID: leaf.String(), Handle: "leaf-a"}, ada
}

// After a verifier passes, its Evidence entry files, and only they, are
// committed on HEAD as one chore(rhiza) commit naming the leaf, with the
// run's Change: and the verifier's Actor: trailers and no Leaf: trailer.
// The builder's staged, unstaged and untracked edits stay exactly as they
// were, and the committed entries show clean. A second call (a resume)
// makes no new commit and returns the same one.
func TestEvidenceCommit(t *testing.T) {
	r, git, ctx, leaf, ada := evidenceCommitRunner(t)
	if err := os.WriteFile(filepath.Join(r.worktree, "notes.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "notes.txt")
	git("commit", "-q", "-m", "notes")
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(r.worktree, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("notes.txt", "base\nstaged\n")
	write("new.txt", "staged new file\n")
	git("add", "notes.txt", "new.txt")
	write("notes.txt", "base\nstaged\nunstaged\n")
	write("scratch.txt", "untracked\n")

	// The verifier writes its Evidence as its own actor: an add and a link,
	// two entries.
	v := *ctx
	v.Actor = ada
	ev := mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-ok","fields":["statement=leaf-a works","result=PASS"]}`)
	ln := mustWrite(t, &v, "link", `{"from":"leaf-a-ok","type":"verifies","to":"leaf-a"}`)
	want := []string{
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ev.Entry),
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ln.Entry),
	}
	slices.Sort(want)

	cached, diff := git("diff", "--cached"), git("diff")
	status := git("status", "--porcelain", "--untracked-files=all", "--", "notes.txt", "new.txt", "scratch.txt")
	old := git("rev-parse", "HEAD")

	c, err := r.CommitEvidence(leaf, "S-v")
	if err != nil {
		t.Fatal(err)
	}
	if head := git("rev-parse", "HEAD"); head != c {
		t.Fatalf("HEAD = %s, want the Evidence commit %s", head, c)
	}
	if parents := git("rev-list", "--parents", "-n1", c); parents != c+" "+old {
		t.Errorf("commit parents = %q, want one new commit on %s", parents, old)
	}
	show := git("show", "--stat", "--format=%s%n%n%b", c)
	t.Logf("git show --stat:\n%s", show)
	files := strings.Split(git("show", "--name-only", "--format=", c), "\n")
	if !reflect.DeepEqual(files, want) {
		t.Errorf("committed files = %v, want only the Evidence entries %v", files, want)
	}
	msg := git("log", "-1", "--format=%B", c)
	if subject := strings.SplitN(msg, "\n", 2)[0]; subject != "chore(rhiza): evidence for leaf-a" {
		t.Errorf("subject = %q", subject)
	}
	trailers := git("log", "-1", "--format=%(trailers:unfold)", c)
	if want := "Change: " + r.change + "\nActor: " + nodeName(ctxWorkspace(t, ctx).model(id.ID{}), ada); trailers != want {
		t.Errorf("trailers = %q, want %q", trailers, want)
	}
	if strings.Contains(msg, "Leaf:") {
		t.Errorf("message has a Leaf: trailer:\n%s", msg)
	}
	if a, b := git("diff", "--cached"), git("diff"); a != cached || b != diff {
		t.Errorf("builder's edits changed:\ncached before:\n%s\nafter:\n%s\ndiff before:\n%s\nafter:\n%s", cached, a, diff, b)
	}
	if s := git("status", "--porcelain", "--untracked-files=all", "--", "notes.txt", "new.txt", "scratch.txt"); s != status {
		t.Errorf("builder's status = %q, want %q", s, status)
	}
	if s := git("status", "--porcelain", "--untracked-files=all", "--", ".rhiza"); s != "" {
		t.Errorf("the .rhiza folder is not clean after the commit:\n%s", s)
	}
	t.Logf("git diff --cached (before = after):\n%s\ngit diff (before = after):\n%s\nstatus (before = after):\n%s", cached, diff, status)

	again, err := r.CommitEvidence(leaf, "S-v")
	if err != nil || again != c {
		t.Errorf("second CommitEvidence = %q, %v; want %s", again, err, c)
	}
	if head := git("rev-parse", "HEAD"); head != c {
		t.Errorf("second call moved HEAD to %s", head)
	}
}

// With no Evidence of the verifier's own actor, nothing is committed.
func TestEvidenceCommitNone(t *testing.T) {
	r, git, _, leaf, _ := evidenceCommitRunner(t)
	old := git("rev-parse", "HEAD")
	if _, err := r.CommitEvidence(leaf, "S-v"); err == nil || !strings.Contains(err.Error(), "no Evidence by") {
		t.Errorf("CommitEvidence err = %v", err)
	}
	if head := git("rev-parse", "HEAD"); head != old {
		t.Errorf("HEAD moved to %s", head)
	}
}

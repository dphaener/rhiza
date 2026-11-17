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
	ctx := testCtx(t, wt)
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
	r := &stepRunner{ctx: ctx, worktree: wt}
	t.Cleanup(r.closeModel)
	return r, git
}

// A builder's leaf commit is the newest commit since its start HEAD whose
// Leaf: trailer names its leaf. With no new commit, or new commits naming
// only another leaf or none, there is none, and the detail says why.
func TestLeafCommit(t *testing.T) {
	t.Parallel()
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
	t.Run("Co-Authored-By in a later paragraph", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "work", "-m", "Change: CS1\nLeaf: leaf-a",
			"-m", "Co-Authored-By: Claude <noreply@anthropic.com>")
		want := git("rev-parse", "HEAD")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != want || detail != "" {
			t.Errorf("LeafCommit = %q, %q, %v; want %s", c, detail, err, want)
		}
	})
	t.Run("a fixer's commit after the builder's, Co-Authored-By in a later paragraph", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		git("commit", "-q", "--allow-empty", "-m", "work", "-m", "Change: CS1\nLeaf: leaf-a\n\nCo-Authored-By: Claude <noreply@anthropic.com>")
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "fix", "-m", "Change: CS1\nLeaf: leaf-a",
			"-m", "Co-Authored-By: Claude <noreply@anthropic.com>")
		want := git("rev-parse", "HEAD")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != want || detail != "" {
			t.Errorf("LeafCommit = %q, %q, %v; want %s", c, detail, err, want)
		}
	})
	t.Run("another leaf with Co-Authored-By in a later paragraph", func(t *testing.T) {
		r, git := leafCommitRunner(t)
		head, _ := r.Head()
		git("commit", "-q", "--allow-empty", "-m", "work", "-m", "Leaf: leaf-b",
			"-m", "Co-Authored-By: Claude <noreply@anthropic.com>")
		c, detail, err := r.LeafCommit(leaf, head)
		if err != nil || c != "" || detail != "new commits name another leaf: leaf-b" {
			t.Errorf("LeafCommit = %q, %q, %v", c, detail, err)
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

// ClosePane closes nothing: its sessions are fakes with no pane.
func (g *gitLeafRunner) ClosePane(string) error { return nil }

type noSleep struct{}

func (noSleep) Now() time.Time      { return time.Date(2026, 10, 27, 0, 0, 0, 0, time.UTC) }
func (noSleep) Sleep(time.Duration) {}

// The engine, against a real git repository, fails a builder that ends
// with no new commit, or with one naming another leaf, with that detail,
// and follows a builder whose commit names its leaf with a verifier given
// that commit, on its leaf-session and in its spec.
func TestLeafCommitEngine(t *testing.T) {
	t.Parallel()
	sr, git := leafCommitRunner(t)
	g := &gitLeafRunner{stepRunner: sr, git: git,
		leaves: []run.Leaf{{ID: "L1", Handle: "leaf-none"}, {ID: "L2", Handle: "leaf-other"}, {ID: "L3", Handle: "leaf-own"},
			{ID: "L4", Handle: "leaf-split-other"}, {ID: "L5", Handle: "leaf-split"}},
		commits: map[string][]string{
			"leaf-other": {"work\n\nLeaf: leaf-x"},
			"leaf-own":   {"work\n\nLeaf: leaf-own"},
			// Co-Authored-By in a paragraph of its own after the trailers.
			"leaf-split-other": {"work\n\nChange: CS1\nLeaf: leaf-x\n\nCo-Authored-By: Claude <noreply@anthropic.com>"},
			"leaf-split":       {"work\n\nChange: CS1\nLeaf: leaf-split\n\nCo-Authored-By: Claude <noreply@anthropic.com>"},
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
	split := git("log", "-1", "--format=%H", "--grep=^Leaf: leaf-split$")
	own := git("log", "-1", "--format=%H", "--grep=^Leaf: leaf-own$")
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
	if j, _ := byHandle["leaf-split-other"].Judged("builder", 1); j.Outcome != "fail" || j.Detail != "new commits name another leaf: leaf-x" {
		t.Errorf("leaf-split-other builder = %+v", j)
	}
	if l := byHandle["leaf-split"]; l.Outcome != run.LeafPassed {
		t.Errorf("leaf-split = %+v, want passed", l)
	} else if j, _ := l.Judged("builder", 1); j.Outcome != "pass" || j.Commit != split {
		t.Errorf("leaf-split builder = %+v, want pass with %s", j, split)
	}
	spec := g.specs[len(g.specs)-1]
	if spec.Role != "verifier" || spec.Commit != split || !slices.Contains(spec.Env, run.EnvCommit+"="+split) {
		t.Errorf("verifier spec = %+v", spec)
	}
}

// evidenceCommitRunner is a stepRunner on a real git worktree holding a
// Rhiza repository, committed, whose change set has leaf "leaf-a", and a
// verifier session "S-v" whose actor is ada, recorded in its own record
// change set whose entries are left uncommitted.
// It returns the runner, a git runner there, the context (to write as
// ada), the leaf and ada.
func evidenceCommitRunner(t *testing.T) (*stepRunner, func(...string) string, *Ctx, run.Leaf, id.ID) {
	t.Helper()
	tmp := t.TempDir()
	wt := filepath.Join(tmp, "wt")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(t, wt)
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
	git("add", "-A")
	git("commit", "-q", "-m", "model")
	// The verifier's actor is recorded as session start records it, its
	// record change set left uncommitted in the worktree.
	ada := recordAda(t, ctx, "ada@example.com")
	sessions := session.NewStore(filepath.Join(tmp, "sessions"))
	if err := sessions.Write(session.Record{ID: "S-v", Actor: ada.String(), Change: cs.String()}); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: ctx, change: cs.String(), run: "R1", worktree: wt, sessions: sessions}
	t.Cleanup(r.closeModel)
	return r, git, ctx, run.Leaf{ID: leaf.String(), Handle: "leaf-a"}, ada
}

// After a verifier passes, its Evidence entry files, and only they, are
// committed on HEAD as one chore(rhiza) commit naming the leaf, with the
// run's Change: and the verifier's Actor: trailers and no Leaf: trailer.
// The builder's staged, unstaged and untracked edits stay exactly as they
// were, and the committed entries show clean. A second call (a resume)
// makes no new commit and returns the same one.
func TestEvidenceCommit(t *testing.T) {
	t.Parallel()
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
	// Another session's write, untracked in the run's change set, which
	// the Evidence commit must leave alone.
	other := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"leaf-b"}`)
	otherPath := fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, other.Entry)
	want := []string{
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ev.Entry),
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ln.Entry),
	}
	var record []string
	for _, f := range strings.Split(git("ls-files", "--others", "--", ".rhiza/log"), "\n") {
		if f != "" && !slices.Contains(want, f) && f != otherPath {
			record = append(record, f)
		}
	}
	if len(record) == 0 {
		t.Fatal("the verifier's actor record left no uncommitted entries")
	}
	for _, f := range record {
		if strings.Contains(f, r.change) {
			t.Fatalf("actor record entry %s is in the run's change set", f)
		}
	}
	want = append(want, record...)
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
	t.Logf("committed files:\n%s", strings.Join(files, "\n"))
	if !reflect.DeepEqual(files, want) {
		t.Errorf("committed files = %v, want only the Evidence and actor record entries %v", files, want)
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
	if s := git("status", "--porcelain", "--untracked-files=all", "--", ".rhiza"); s != "?? "+otherPath {
		t.Errorf(".rhiza status after the commit = %q, want only the other write %s untracked", s, otherPath)
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

// raceEvidence writes the verifier's Evidence for leaf (an add and a
// link) and returns the paths of those two entries.
func raceEvidence(t *testing.T, r *stepRunner, ctx *Ctx, ada id.ID) []string {
	t.Helper()
	v := *ctx
	v.Actor = ada
	ev := mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-ok","fields":["statement=leaf-a works","result=PASS"]}`)
	ln := mustWrite(t, &v, "link", `{"from":"leaf-a-ok","type":"verifies","to":"leaf-a"}`)
	return []string{
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ev.Entry),
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ln.Entry),
	}
}

// setCommitOnlyHeadRead sets the commitOnly hook for this test alone. Its
// callers do not run in parallel, so no other test sees it.
func setCommitOnlyHeadRead(t *testing.T, hook func(worktree string)) {
	t.Helper()
	commitOnlyHeadRead = hook
	t.Cleanup(func() { commitOnlyHeadRead = nil })
}

// Another session committing on the run worktree's HEAD between
// commitOnly's read of HEAD and its update-ref does not fail the Evidence
// commit: it is built again on the new HEAD, so both commits are on the
// branch, the racing commit's files are intact and the real index shows
// the Evidence entries committed.
func TestCommitOnlyHeadMoved(t *testing.T) {
	r, git, ctx, leaf, ada := evidenceCommitRunner(t)
	evidence := raceEvidence(t, r, ctx, ada)
	old := git("rev-parse", "HEAD")
	var racer string
	calls := 0
	setCommitOnlyHeadRead(t, func(wt string) {
		calls++
		if calls > 1 {
			return
		}
		if err := os.WriteFile(filepath.Join(wt, "racer.txt"), []byte("racing commit\n"), 0o644); err != nil {
			t.Error(err)
			return
		}
		git("add", "racer.txt")
		git("commit", "-q", "-m", "racing commit")
		racer = git("rev-parse", "HEAD")
	})

	c, err := r.CommitEvidence(leaf, "S-v")
	if err != nil {
		t.Fatalf("CommitEvidence with HEAD moved: %v", err)
	}
	if calls != 2 {
		t.Errorf("commitOnly built its commit %d times, want 2", calls)
	}
	if racer == "" || racer == old {
		t.Fatalf("the hook made no racing commit (HEAD %s)", racer)
	}
	if head := git("rev-parse", "HEAD"); head != c {
		t.Fatalf("HEAD = %s, want the Evidence commit %s", head, c)
	}
	if parents := git("rev-list", "--parents", "-n1", c); parents != c+" "+racer {
		t.Errorf("Evidence commit parents = %q, want the racing commit %s", parents, racer)
	}
	if log := git("rev-list", old+"..HEAD"); log != c+"\n"+racer {
		t.Errorf("commits on the branch since %s = %q, want the Evidence commit and the racing one", old, log)
	}
	if body := git("show", c+":racer.txt"); body != "racing commit" {
		t.Errorf("racer.txt in the Evidence commit = %q", body)
	}
	if files := git("show", "--name-only", "--format=", racer); files != "racer.txt" {
		t.Errorf("racing commit's files = %q, want racer.txt", files)
	}
	files := strings.Split(git("show", "--name-only", "--format=", c), "\n")
	for _, p := range evidence {
		if !slices.Contains(files, p) {
			t.Errorf("Evidence commit's files %v lack %s", files, p)
		}
	}
	if slices.Contains(files, "racer.txt") {
		t.Errorf("Evidence commit rewrote racer.txt: %v", files)
	}
	if s := git(append([]string{"status", "--porcelain", "--untracked-files=all", "--", "racer.txt"}, evidence...)...); s != "" {
		t.Errorf("status of the Evidence entries and racer.txt = %q, want them committed", s)
	}
	if d := git("diff", "--cached", "--name-only"); d != "" {
		t.Errorf("real index differs from HEAD: %q", d)
	}
}

// A HEAD that moves on every try fails the Evidence commit after
// commitOnlyTries tries with update-ref's error, HEAD left at the last
// racing commit.
func TestCommitOnlyHeadMovedEveryTry(t *testing.T) {
	r, git, ctx, leaf, ada := evidenceCommitRunner(t)
	raceEvidence(t, r, ctx, ada)
	calls := 0
	setCommitOnlyHeadRead(t, func(string) {
		calls++
		git("commit", "-q", "--allow-empty", "-m", fmt.Sprintf("racing commit %d", calls))
	})
	_, err := r.CommitEvidence(leaf, "S-v")
	if err == nil {
		t.Fatal("CommitEvidence succeeded with HEAD moving on every try")
	}
	if calls != commitOnlyTries || commitOnlyTries != 5 {
		t.Errorf("tries = %d (commitOnlyTries %d), want 5", calls, commitOnlyTries)
	}
	if !strings.Contains(err.Error(), "update-ref") {
		t.Errorf("err = %v, want update-ref's error", err)
	}
	if subject := git("log", "-1", "--format=%s"); subject != "racing commit 5" {
		t.Errorf("HEAD's subject = %q, want the last racing commit", subject)
	}
}

// With no Evidence of the verifier's own actor, nothing is committed.
func TestEvidenceCommitNone(t *testing.T) {
	t.Parallel()
	r, git, _, leaf, _ := evidenceCommitRunner(t)
	old := git("rev-parse", "HEAD")
	if _, err := r.CommitEvidence(leaf, "S-v"); err == nil || !strings.Contains(err.Error(), "no Evidence by") {
		t.Errorf("CommitEvidence err = %v", err)
	}
	if head := git("rev-parse", "HEAD"); head != old {
		t.Errorf("HEAD moved to %s", head)
	}
}

// A builder's prompt is its role's, rendered with the leaf and attempt; a
// retry builder's also names the leaf and its attempt and points it at
// the earlier attempts' findings. A verifier's gets no retry line.
func TestLeafPrompt(t *testing.T) {
	t.Parallel()
	role := workflow.Role{Prompt: "build {leaf} attempt {attempt} of {change}"}
	vars := workflow.PromptVars{Change: "CS1"}
	spec := run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 1}
	if got := leafPrompt(role, spec, vars); got != "build leaf-a attempt 1 of CS1" {
		t.Errorf("attempt 1 = %q", got)
	}
	spec.Attempt = 2
	got := leafPrompt(role, spec, vars)
	for _, want := range []string{"build leaf-a attempt 2 of CS1\n\n", "This is attempt 2 of 3 on leaf leaf-a", "rhiza run status --json", "Leaf: leaf-a"} {
		if !strings.Contains(got, want) {
			t.Errorf("attempt 2 = %q, lacks %q", got, want)
		}
	}
	spec.Role = workflow.RoleVerifier
	if got := leafPrompt(role, spec, vars); got != "build leaf-a attempt 2 of CS1" {
		t.Errorf("verifier attempt 2 = %q", got)
	}

	// A fresh round's first builder gets the answering Decision and no
	// retry line; its second gets both, counted within the round.
	spec = run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 4, Round: 2, RoundAttempt: 1,
		Question: "Q1", Decision: "use-a-map", Answer: "Use a map, not a list."}
	got = leafPrompt(role, spec, vars)
	want := "build leaf-a attempt 4 of CS1\n\nLeaf leaf-a failed every attempt of its last round and raised question Q1. " +
		"Decision use-a-map answers it:\n\nUse a map, not a list.\n\nBuild the leaf as it says; this round has 3 fresh attempts."
	if got != want {
		t.Errorf("round 2 attempt 1 =\n%q\nwant\n%q", got, want)
	}
	spec.Attempt, spec.RoundAttempt = 5, 2
	got = leafPrompt(role, spec, vars)
	for _, want := range []string{"Decision use-a-map answers it:\n\nUse a map, not a list.", "This is attempt 2 of 3 on leaf leaf-a"} {
		if !strings.Contains(got, want) {
			t.Errorf("round 2 attempt 2 = %q, lacks %q", got, want)
		}
	}
	spec.Role = workflow.RoleVerifier
	if got := leafPrompt(role, spec, vars); got != "build leaf-a attempt 5 of CS1" {
		t.Errorf("verifier in round 2 = %q", got)
	}

	// The retry's round names the retry and no Decision.
	spec = run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 4, Round: 2, RoundAttempt: 1,
		Question: "Q1", Retry: true}
	got = leafPrompt(role, spec, vars)
	for _, want := range []string{"A person retried the run", "fresh round of 3 attempts", "no Decision answers its failure"} {
		if !strings.Contains(got, want) {
			t.Errorf("retry round = %q, lacks %q", got, want)
		}
	}
	if strings.Contains(got, "Decision use-a-map") || strings.Contains(got, "answers it:") || strings.Contains(got, "raised question") {
		t.Errorf("retry round names a Decision: %q", got)
	}
}

// A retry builder's prompt quotes the previous attempt's findings byte
// for byte, a multi-line one and one past 200 characters among them, and
// tells it to inspect the working tree first and address every finding as
// written; with no findings it quotes the attempt's failure instead, and a
// first attempt gets no findings section.
func TestLeafPromptQuotesFindings(t *testing.T) {
	t.Parallel()
	role := workflow.Role{Prompt: "build {leaf}"}
	vars := workflow.PromptVars{Change: "CS1"}
	multi := "the exit code is wrong:\n  got 0\n  want 2\n"
	long := "  leading space kept; " + strings.Repeat("this reason runs well past two hundred characters, ", 6) + "and ends here.  "
	spec := run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 2, RoundAttempt: 2,
		Findings: []run.Finding{
			{ID: "f1", Check: "verifier", Severity: "major", Place: "cmd/x.go:12", Reason: multi},
			{ID: "f2", Reason: long},
		},
		Failure: "verifier failed: verifier found 2 problem(s)"}
	got := leafPrompt(role, spec, vars)
	for _, want := range []string{
		"Finding f1 (verifier, major at cmd/x.go:12):\n" + multi,
		"Finding f2:\n" + long,
		"inspect the working tree and the earlier attempts' commits",
		"address every finding exactly as written, in one fix commit with the trailer `Leaf: leaf-a`",
		"`rhiza run status --json`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("retry prompt lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "it failed with") {
		t.Errorf("retry prompt with findings quotes the failure:\n%s", got)
	}

	// No findings journalled: the failure detail is quoted, no list.
	spec.Findings = nil
	spec.Failure = "verifier failed: no Evidence verifies the leaf"
	got = leafPrompt(role, spec, vars)
	for _, want := range []string{
		"journalled no findings; it failed with:\n\nverifier failed: no Evidence verifies the leaf\n\n",
		"inspect the working tree",
		"`rhiza run status --json`",
		"`Leaf: leaf-a`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("no-findings prompt lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Finding ") || strings.Contains(got, "these findings") {
		t.Errorf("no-findings prompt has a findings list:\n%s", got)
	}

	// A first attempt has no findings section.
	spec = run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-a"}, Attempt: 1, RoundAttempt: 1}
	if got := leafPrompt(role, spec, vars); got != "build leaf-a" {
		t.Errorf("first attempt = %q", got)
	}
}

// The run's notes, journaled and read back oldest first, end a later
// leaf's builder and verifier prompts under their heading; a fixer's
// prompt carries none, and a run with no notes adds nothing.
func TestLeafPromptNotes(t *testing.T) {
	t.Parallel()
	store := run.NewStore(t.TempDir())
	t0 := time.Date(2026, 11, 17, 5, 0, 0, 0, time.UTC)
	if err := store.Create("R1", run.Entry{Type: run.TypeStarted, Step: "build", At: t0,
		Started: &run.Started{Workflow: "w", Change: "CS1"}}); err != nil {
		t.Fatal(err)
	}
	for i, text := range []string{"first: make gate needs RHIZA_GATE_BASE", "second: the fake rex ignores -b"} {
		if err := store.Append("R1", run.Entry{Type: run.TypeNote, At: t0.Add(time.Duration(i+1) * time.Minute), Step: "build", Role: "builder", Leaf: "leaf-a", Text: text}); err != nil {
			t.Fatal(err)
		}
	}
	j, err := store.Read("R1")
	if err != nil {
		t.Fatal(err)
	}
	notes := j.State().Notes
	role := workflow.Role{Prompt: "work on {leaf}"}
	vars := workflow.PromptVars{Change: "CS1"}
	want := "work on leaf-b\n\nNotes from earlier sessions of this run:\n" +
		"- first: make gate needs RHIZA_GATE_BASE\n- second: the fake rex ignores -b"
	for _, r := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
		spec := run.LeafSpec{Role: r, Leaf: run.Leaf{Handle: "leaf-b"}, Attempt: 1}
		if got := withRunNotes(leafPrompt(role, spec, vars), r, notes); got != want {
			t.Errorf("%s prompt =\n%q\nwant\n%q", r, got, want)
		}
		if got := withRunNotes(leafPrompt(role, spec, vars), r, nil); got != "work on leaf-b" || strings.Contains(got, run.NotesHeading) {
			t.Errorf("%s prompt with no notes = %q", r, got)
		}
	}
	// A retry builder's notes follow its retry lines.
	spec := run.LeafSpec{Role: workflow.RoleBuilder, Leaf: run.Leaf{Handle: "leaf-b"}, Attempt: 2, RoundAttempt: 2, Failure: "verifier failed"}
	if got := withRunNotes(leafPrompt(role, spec, vars), workflow.RoleBuilder, notes); !strings.HasSuffix(got, "\n\n"+run.NotesSection(notes)) || !strings.Contains(got, "This is attempt 2") {
		t.Errorf("retry builder prompt = %q", got)
	}
	fixer := workflow.Role{Prompt: "fix {leaf}"}
	fspec := run.LeafSpec{Role: workflow.RoleFixer, Leaf: run.Leaf{Handle: "leaf-b"}, Attempt: 1, Cycle: 1}
	if got := withRunNotes(leafPrompt(fixer, fspec, vars), workflow.RoleFixer, notes); strings.Contains(got, run.NotesHeading) || strings.Contains(got, "first:") {
		t.Errorf("fixer prompt carries notes: %q", got)
	}
}

// A round the retry gave shows as the retry's in run status, with no
// decision.
func TestRunInfoRetryRound(t *testing.T) {
	t.Parallel()
	st := run.State{Leaves: []run.LeafState{{Step: "build", Leaf: "LA", Handle: "leaf-a", Rounds: []run.LeafRound{
		{Round: 2, After: 3, Question: "Q1", QuestionHandle: "q-one", Decision: "use-a-map"},
		{Round: 3, After: 6, Question: "Q2", Retry: true},
	}}}}
	info := runInfo(&Ctx{}, st, false, nil)
	if len(info.Leaves) != 1 || len(info.Leaves[0].Rounds) != 2 {
		t.Fatalf("leaves %+v", info.Leaves)
	}
	rs := info.Leaves[0].Rounds
	if rs[0].Retry || rs[0].Decision != "use-a-map" || !rs[1].Retry || rs[1].Decision != "" {
		t.Errorf("rounds %+v", rs)
	}
	var b strings.Builder
	info.Human(&b)
	for _, want := range []string{
		"\n  - round 2 from attempt 4: decision use-a-map answers question q-one (Q1)\n",
		"\n  - round 3 from attempt 7: given by the retry\n",
	} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("run status text lacks %q:\n%s", want, b.String())
		}
	}
	if strings.Contains(b.String(), "decision  ") {
		t.Errorf("run status text names an empty decision:\n%s", b.String())
	}
}

// AddLeafQuestion writes one question node into the run's change set,
// in the worktree where its sessions write, with the text as its
// statement and a blocks link to the leaf; a second gets the next handle.
// LeafAnswer finds the owner's Decision answering it, whatever its
// wording, and nothing before one does.
func TestAddLeafQuestion(t *testing.T) {
	t.Parallel()
	_, r, wctx := arbiterFixture(t)
	leaf := mustWrite(t, wctx, "add", `{"kind":"capability","handle":"leaf-a"}`).Nodes[0]
	cs := mustParseID(r.change)
	q, err := r.AddLeafQuestion(run.Leaf{ID: leaf.String(), Handle: "leaf-a"}, "Leaf leaf-a failed all 3 attempts of round 1. How?")
	if err != nil {
		t.Fatal(err)
	}
	if q.Handle != "leaf-failed-01runarb-1" {
		t.Errorf("handle %q", q.Handle)
	}
	ww, err := r.worktreeWorkspace(false)
	if err != nil {
		t.Fatal(err)
	}
	defer ww.closeProjection()
	m := ww.model(cs)
	n, ok := m.Node(mustParseID(q.ID))
	if !ok || n.Kind != "question" || n.Handle != q.Handle || fieldText(n.Fields["statement"]) != "Leaf leaf-a failed all 3 attempts of round 1. How?" {
		t.Fatalf("question node %+v", n)
	}
	blocks := m.LinksTo(leaf, "blocks")
	if len(blocks) != 1 || blocks[0].From.String() != q.ID {
		t.Errorf("blocks links to the leaf: %+v", blocks)
	}
	c, _ := ww.changeSet(cs)
	created := false
	for _, o := range c.Ops {
		created = created || (o.Type == "node.create" && o.Handle == q.Handle)
	}
	if !created {
		t.Errorf("the question is not created in the run's change set %s", cs)
	}
	q2, err := r.AddLeafQuestion(run.Leaf{ID: leaf.String(), Handle: "leaf-a"}, "Again?")
	if err != nil || q2.Handle != "leaf-failed-01runarb-2" {
		t.Errorf("second question %+v, %v", q2, err)
	}

	if d, s, ok, err := r.LeafAnswer(q.ID); err != nil || ok || d != "" || s != "" {
		t.Errorf("unanswered: %q %q %v %v", d, s, ok, err)
	}
	mustWrite(t, wctx, "add", `{"kind":"decision","handle":"use-a-map","fields":["statement=Use a map."]}`)
	mustWrite(t, wctx, "link", fmt.Sprintf(`{"from":"use-a-map","type":"answers","to":%q}`, q.ID))
	if d, s, ok, err := r.LeafAnswer(q.ID); err != nil || !ok || d != "use-a-map" || s != "Use a map." {
		t.Errorf("answered: %q %q %v %v", d, s, ok, err)
	}
	if _, err := r.AddLeafQuestion(run.Leaf{ID: id.ID{}.String(), Handle: "gone"}, "?"); err == nil {
		t.Error("a question on a leaf that is no node was added")
	}
}

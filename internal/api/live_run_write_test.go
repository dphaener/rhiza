package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/dphaener/rhiza/internal/workflow"
)

func startLiveRun(t *testing.T, f *sessionFixture, cs string) {
	t.Helper()
	startLiveRunAt(t, f.common, f.worktree, cs)
}

// startLiveRunAt starts run R1 on change set cs in worktree wt, its journal
// under git common dir common.
func startLiveRunAt(t *testing.T, common, wt, cs string) {
	t.Helper()
	src := "[step.s]\nkind = \"leaves\"\n[step.s.builder]\nprompt = \"b\"\n[step.s.verifier]\nprompt = \"v\"\n"
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)
	if err := run.Start(run.NewStore(run.Dir(common)), "R1", run.Started{Workflow: "wf", Definition: src, Change: cs, Worktree: wt}, now, def); err != nil {
		t.Fatal(err)
	}
}

func wtLogFiles(t *testing.T, root, cs string) []string {
	t.Helper()
	es, err := os.ReadDir(filepath.Join(root, ".rhiza", "log", cs))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	var out []string
	for _, e := range es {
		out = append(out, e.Name())
	}
	return out
}

// An accept made in the main checkout while the change set's run is live
// also lands in the run's worktree, byte for byte.
func TestLiveRunWriteReachesWorktree(t *testing.T) {
	f := newSessionFixture(t, "live")
	cs := newChange(t, f.ctx, `{"intent":"Live"}`).ChangeSet.String()
	if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
		t.Fatal(err)
	}
	startLiveRun(t, f, cs)
	resp, err := callOp(t, f.ctx, "accept", `{"reason":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	name := resp.(LifecycleResponse).Entry.String() + ".json"
	a, err := os.ReadFile(filepath.Join(f.ctx.Root, ".rhiza", "log", cs, name))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(f.worktree, ".rhiza", "log", cs, name))
	if err != nil || string(a) != string(b) {
		t.Fatalf("worktree copy of the accept: %v, equal=%v", err, string(a) == string(b))
	}
	// The open and propose, uncommitted in the main checkout, come along,
	// so the worktree's log holds the whole change set.
	if got, want := wtLogFiles(t, f.worktree, cs), wtLogFiles(t, f.ctx.Root, cs); !slices.Equal(got, want) || len(got) != 3 {
		t.Errorf("worktree holds %v, want the main checkout's %v", got, want)
	}
	assertFoldsAccepted(t, f.worktree, cs)
}

// assertFoldsAccepted checks that the log under root folds with no
// problems and holds change set cs, accepted, with every entry applied.
func assertFoldsAccepted(t *testing.T, root, cs string) {
	t.Helper()
	l := store.ReadLog(os.DirFS(filepath.Join(root, store.Dir)))
	sets, problems := changeset.FoldLog(l)
	if len(l.Problems) != 0 || len(problems) != 0 {
		t.Errorf("worktree log problems: %v, fold problems: %v", l.Problems, problems)
	}
	for _, c := range sets {
		if c.ID.String() == cs {
			if c.State != changeset.StateAccepted || len(c.Applied) != len(wtLogFiles(t, root, cs)) {
				t.Errorf("worktree folds %s as %s with %d of %d entries applied", cs, c.State, len(c.Applied), len(wtLogFiles(t, root, cs)))
			}
			return
		}
	}
	t.Errorf("worktree log has no change set %s", cs)
}

// The copies of main-checkout writes go into the run's next commit on its
// branch: here a check session's record commit, which would otherwise hold
// only that session's actor. Run with real git, the main checkout and the
// run's worktree linked by git worktree add.
func TestLiveRunWriteCommittedByRunsNextCommit(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "repo")
	wt := filepath.Join(tmp, "repo.live")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	exec := hermeticExec(t, tmp)
	git := func(dir string, args ...string) string {
		t.Helper()
		out, err := exec("git", args, dir, nil)
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	ctxAt := func(dir string) *Ctx {
		c := testCtx(dir)
		c.Exec = exec
		c.Git = func(args ...string) (string, error) {
			out, err := exec("git", args, dir, nil)
			return string(out), err
		}
		return c
	}
	main := ctxAt(root)
	git(root, "init", "-q", "-b", "main")
	if _, err := callOp(t, main, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "init")
	git(root, "worktree", "add", "-q", "-b", "cs/live", wt)

	// The change set is opened and proposed in the main checkout and left
	// uncommitted there; then its run starts and the owner accepts it.
	cs := newChange(t, main, `{"intent":"Live"}`).ChangeSet.String()
	if _, err := callOp(t, main, "propose", `{}`); err != nil {
		t.Fatal(err)
	}
	startLiveRunAt(t, filepath.Join(root, ".git"), wt, cs)
	resp, err := callOp(t, main, "accept", `{"reason":"ok"}`)
	if err != nil {
		t.Fatal(err)
	}
	accept := filepath.ToSlash(filepath.Join(".rhiza", "log", cs, resp.(LifecycleResponse).Entry.String()+".json"))
	assertFoldsAccepted(t, wt, cs)

	wctx := ctxAt(wt)
	ada := recordAda(t, wctx, "ada@example.com")
	sessions := session.NewStore(filepath.Join(tmp, "sessions"))
	if err := sessions.Write(session.Record{ID: "S-c", Actor: ada.String(), Change: cs, Check: "quality", Role: "critic", Cycle: 1}); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: wctx, change: cs, run: "R1", worktree: wt, main: root, sessions: sessions}
	c, err := r.CommitCheckWrites("S-c")
	if err != nil || c == "" {
		t.Fatalf("CommitCheckWrites = %q, %v", c, err)
	}
	if head := git(wt, "rev-parse", "HEAD"); head != c {
		t.Fatalf("HEAD %s, want the run's commit %s", head, c)
	}
	files := strings.Split(git(wt, "show", "--name-only", "--format=", c), "\n")
	for _, name := range wtLogFiles(t, root, cs) {
		if p := ".rhiza/log/" + cs + "/" + name; !slices.Contains(files, p) {
			t.Errorf("run's commit lacks %s; it holds %v", p, files)
		}
	}
	if !slices.Contains(files, accept) {
		t.Errorf("run's commit lacks the accept %s", accept)
	}
	if s := git(wt, "status", "--porcelain", "--untracked-files=all", "--", ".rhiza/log/"+cs); s != "" {
		t.Errorf("worktree still holds uncommitted entries of %s: %q", cs, s)
	}
	// The main checkout's copies are byte for byte the committed ones, so
	// the run's branch holds them.
	for _, name := range wtLogFiles(t, root, cs) {
		p := ".rhiza/log/" + cs + "/" + name
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
		if err != nil {
			t.Fatal(err)
		}
		if got := git(wt, "show", "HEAD:"+p); got != strings.TrimSpace(string(b)) {
			t.Errorf("branch's %s differs from the main checkout's", p)
		}
	}
}

// A change set with no live run, and a dry run, copy nothing.
func TestLiveRunWriteNoRunNoCopy(t *testing.T) {
	f := newSessionFixture(t, "none")
	cs := newChange(t, f.ctx, `{"intent":"Quiet"}`).ChangeSet.String()
	if _, err := callOp(t, f.ctx, "propose", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := wtLogFiles(t, f.worktree, cs); len(got) != 0 {
		t.Errorf("no live run, worktree holds %v", got)
	}
	startLiveRun(t, f, cs)
	f.ctx.DryRun = true
	if _, err := callOp(t, f.ctx, "accept", `{}`); err != nil {
		t.Fatal(err)
	}
	if got := wtLogFiles(t, f.worktree, cs); len(got) != 0 {
		t.Errorf("dry run, worktree holds %v", got)
	}
}

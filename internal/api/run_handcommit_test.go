package api

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// handCommitBranch is a stepRunner over a real repository on branch
// cs/x, branched from main, with the run's state (run R1, Base the
// branch's first commit) and a git function.
func handCommitBranch(t *testing.T) (*stepRunner, run.State, func(...string) string) {
	t.Helper()
	r, git := leafCommitRunner(t)
	r.run = "R1"
	r.sessions = session.NewStore(filepath.Join(t.TempDir(), "sessions"))
	if err := r.sessions.Write(session.Record{ID: "S1", Run: "R1", Actor: "claude-code-s1", StartedAt: time.Unix(0, 0)}); err != nil {
		t.Fatal(err)
	}
	if err := r.sessions.Write(session.Record{ID: "S2", Run: "R2", Actor: "claude-code-s2", StartedAt: time.Unix(0, 0)}); err != nil {
		t.Fatal(err)
	}
	r.ctx.Git = func(args ...string) (string, error) {
		out, err := r.ctx.Exec("git", append([]string{"-C", r.worktree}, args...), r.worktree, nil)
		return string(out), err
	}
	git("checkout", "-q", "-b", "cs/x")
	st := run.State{ID: "R1", Change: "CS1", Branch: "cs/x", Base: branchHead(r.ctx, "cs/x")}
	if st.Base == "" {
		t.Fatal("branchHead found no commit on cs/x")
	}
	return r, st, git
}

func scanHand(t *testing.T, r *stepRunner, st run.State, leaves []run.Leaf) []string {
	t.Helper()
	found, err := handCommits(r.ctx.Git, st, leaves, r.sessions)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, iv := range found {
		out = append(out, iv.Kind+" "+iv.Commit+" "+iv.Actor)
	}
	return out
}

// The run's own commits (a Leaf: commit of its change set, a Review: <run>
// commit, an Actor: commit of a run session and a merge of main) are no
// hand commits; a plain commit made by hand is one, its actor the author;
// a commit whose Review: names another run, or whose Actor: is another
// run's session, is one too, its actor the trailer's.
func TestHandCommitsByTrailer(t *testing.T) {
	t.Parallel()
	r, st, git := handCommitBranch(t)
	leaves := []run.Leaf{{ID: "L1", Handle: "leaf-a"}}
	git("commit", "-q", "--allow-empty", "-m", "build leaf", "-m", "Change: CS1\nActor: claude-code-b1\nLeaf: leaf-a")
	git("commit", "-q", "--allow-empty", "-m", "fix review", "-m", "Change: CS1\nReview: R1")
	git("commit", "-q", "--allow-empty", "-m", "plan", "-m", "Change: CS1\nActor: claude-code-s1")
	git("checkout", "-q", "main")
	git("commit", "-q", "--allow-empty", "-m", "main moves on")
	git("commit", "-q", "--allow-empty", "-m", "main again", "-m", "Actor: someone-on-main")
	git("checkout", "-q", "cs/x")
	git("merge", "-q", "--no-ff", "--no-edit", "main")
	if got := scanHand(t, r, st, leaves); len(got) != 0 {
		t.Fatalf("hand commits = %q, want none", got)
	}

	git("commit", "-q", "--allow-empty", "-m", "fix by hand")
	plain := git("rev-parse", "HEAD")
	want := []string{"hand-commit " + plain + " Rhiza Test"}
	if got := scanHand(t, r, st, leaves); !reflect.DeepEqual(got, want) {
		t.Fatalf("hand commits = %q, want %q", got, want)
	}

	// Failure case: trailers that name another run, its session's actor or
	// another change set's leaf do not make a commit the run's.
	git("commit", "-q", "--allow-empty", "-m", "other review", "-m", "Change: CS1\nReview: R2")
	other := git("rev-parse", "HEAD")
	git("commit", "-q", "--allow-empty", "-m", "other session", "-m", "Change: CS1\nActor: claude-code-s2\nLeaf: leaf-z")
	stranger := git("rev-parse", "HEAD")
	want = append(want, "hand-commit "+other+" Rhiza Test", "hand-commit "+stranger+" claude-code-s2")
	if got := scanHand(t, r, st, leaves); !reflect.DeepEqual(got, want) {
		t.Fatalf("hand commits = %q, want %q", got, want)
	}

	// A merge of a branch that is not the base branch is a hand commit.
	git("checkout", "-q", "-b", "side", "main")
	git("commit", "-q", "--allow-empty", "-m", "side work")
	git("checkout", "-q", "cs/x")
	git("merge", "-q", "--no-ff", "--no-edit", "side")
	merge := git("rev-parse", "HEAD")
	if got := scanHand(t, r, st, leaves); len(got) != 4 || got[3] != "hand-commit "+merge+" Rhiza Test" {
		t.Fatalf("hand commits = %q, want the side merge %s last", got, merge)
	}
}

// Commits on another branch, or before the run's Base, are never scanned;
// a run with no recorded Base or branch scans nothing.
func TestHandCommitsOnlyRunBranch(t *testing.T) {
	t.Parallel()
	r, st, git := handCommitBranch(t)
	git("checkout", "-q", "-b", "elsewhere")
	git("commit", "-q", "--allow-empty", "-m", "elsewhere by hand")
	git("checkout", "-q", "cs/x")
	if got := scanHand(t, r, st, nil); len(got) != 0 {
		t.Fatalf("hand commits = %q, want none", got)
	}
	git("commit", "-q", "--allow-empty", "-m", "by hand")
	for _, s := range []run.State{{ID: "R1", Branch: "cs/x"}, {ID: "R1", Base: st.Base}} {
		found, err := r.HandCommits(s)
		if err != nil || len(found) != 0 {
			t.Errorf("HandCommits(%+v) = %v, %v; want none", s, found, err)
		}
	}
}

// handEngineRunner runs every command step as passing and scans hand
// commits with the stepRunner.
type handEngineRunner struct {
	run.Runner
	s *stepRunner
}

func (h handEngineRunner) Command(workflow.Step, []string) (int, string, error) { return 0, "", nil }
func (h handEngineRunner) HandCommits(st run.State) ([]run.Intervention, error) {
	return h.s.HandCommits(st)
}

const handCommitWF = `
name = "wf"
[step.one]
kind = "command"
run = "true"
on.pass = "two"

[step.two]
kind = "command"
run = "true"
`

// The engine journals one hand-commit intervention naming a plain commit
// made by hand, at the step it scanned on; the later scans (entering step
// two, before the finish) add none, and the run succeeds unchanged.
func TestHandCommitJournalledOnce(t *testing.T) {
	t.Parallel()
	r, st, git := handCommitBranch(t)
	git("commit", "-q", "--allow-empty", "-m", "fix review", "-m", "Review: R1")
	git("commit", "-q", "--allow-empty", "-m", "fix by hand")
	plain := git("rev-parse", "HEAD")
	def, err := workflow.Parse("wf", []byte(handCommitWF))
	if err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(t.TempDir())
	cache := t.TempDir()
	rs.Getenv = func(k string) string {
		if k == "XDG_CACHE_HOME" {
			return cache
		}
		return ""
	}
	// No change set: the run's leaves are none, read from no model.
	if err := run.Start(rs, "R1", run.Started{Workflow: "wf", Definition: handCommitWF, Branch: st.Branch, Base: st.Base}, time.Unix(0, 0), def); err != nil {
		t.Fatal(err)
	}
	e := &run.Engine{Store: rs, Runner: handEngineRunner{s: r}, Clock: run.RealClock{}, Now: func() time.Time { return time.Unix(10, 0) }}
	got, err := e.Run("R1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Result != run.Succeeded {
		t.Fatalf("result %q detail %q", got.Result, got.Detail)
	}
	j, err := rs.Read("R1")
	if err != nil {
		t.Fatal(err)
	}
	ivs := j.State().Interventions
	if len(ivs) != 1 || ivs[0].Kind != run.InterventionHandCommit || ivs[0].Commit != plain || ivs[0].Actor != "Rhiza Test" || ivs[0].Step != "one" ||
		!strings.Contains(ivs[0].Detail, "fix by hand") {
		t.Fatalf("interventions = %+v, want one hand-commit of %s at step one", ivs, plain)
	}
}

// A branch still at Base that merges main fast-forwards to it, as ship's
// git merge --no-edit does: main's commits, now on the run's first-parent
// line, are no hand commits. A commit made by hand after them still is.
func TestHandCommitsFastForwardedBranch(t *testing.T) {
	t.Parallel()
	r, st, git := handCommitBranch(t)
	git("checkout", "-q", "main")
	git("commit", "-q", "--allow-empty", "-m", "main moves on")
	git("commit", "-q", "--allow-empty", "-m", "main again", "-m", "Actor: someone-on-main")
	git("checkout", "-q", "cs/x")
	git("merge", "-q", "--no-edit", "main")
	if got := scanHand(t, r, st, nil); len(got) != 0 {
		t.Fatalf("hand commits = %q, want none", got)
	}
	git("commit", "-q", "--allow-empty", "-m", "fix by hand")
	plain := git("rev-parse", "HEAD")
	want := []string{"hand-commit " + plain + " Rhiza Test"}
	if got := scanHand(t, r, st, nil); !reflect.DeepEqual(got, want) {
		t.Fatalf("hand commits = %q, want %q", got, want)
	}
}

// A person rebasing the run's branch onto newer main puts main's commits on
// its first-parent line: they are no hand commits, while the run's own
// rebased commit stays the run's and a hand commit stays one.
func TestHandCommitsRebasedBranch(t *testing.T) {
	t.Parallel()
	r, st, git := handCommitBranch(t)
	leaves := []run.Leaf{{ID: "L1", Handle: "leaf-a"}}
	git("commit", "-q", "--allow-empty", "-m", "build leaf", "-m", "Change: CS1\nActor: claude-code-b1\nLeaf: leaf-a")
	git("commit", "-q", "--allow-empty", "-m", "fix by hand")
	git("checkout", "-q", "main")
	git("commit", "-q", "--allow-empty", "-m", "main moves on")
	git("commit", "-q", "--allow-empty", "-m", "main again", "-m", "Actor: someone-on-main")
	git("checkout", "-q", "cs/x")
	git("rebase", "-q", "--keep-empty", "main")
	plain := git("rev-parse", "HEAD")
	want := []string{"hand-commit " + plain + " Rhiza Test"}
	if got := scanHand(t, r, st, leaves); !reflect.DeepEqual(got, want) {
		t.Fatalf("hand commits = %q, want %q", got, want)
	}
}

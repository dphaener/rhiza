package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
)

// otherRoot is a plan run's worktree, a repository other than the caller's
// root: its plan change set landed an intent, a planned change set asks
// two questions, and run ROTHER is stopped at its gate on both. The run
// lives in the git common folder common.
type otherRoot struct {
	root    string
	own, cs id.ID
	common  string
}

// otherRootFixture builds the worktree of an otherRoot and a caller's Ctx
// rooted in a different repository, whose user cache folder is the test's
// HOME/.cache (home empty: no user cache folder at all). The caller's own
// index path is its own root's, never the worktree's.
func otherRootFixture(t *testing.T, home string) (*Ctx, otherRoot) {
	t.Helper()
	wt := initCtx(t)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, wt, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	o := otherRoot{root: wt.Root, common: t.TempDir()}
	o.own = call("change.new", `{"intent":"Plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"other-plan-intent","fields":["statement=Plan the other root"]}`)
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	call("land", `{}`)
	o.cs = call("change.new", `{"intent":"Planned"}`).(ChangeNewResponse).ChangeSet
	call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"other-plan-intent"}`, o.cs.String()))
	call("add", `{"kind":"question","handle":"other-q","fields":["statement=Statement of other-q"]}`)
	call("add", `{"kind":"question","handle":"other-r","fields":["statement=Statement of other-r"]}`)
	call("propose", `{}`)

	runs := run.NewStore(run.Dir(o.common))
	at := wt.Now()
	if err := runs.Create("ROTHER", run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan", Change: o.own.String(), Worktree: o.root}}); err != nil {
		t.Fatal(err)
	}
	if err := runs.Append("ROTHER", run.Entry{Type: run.TypeStopped, At: at, Step: "gate", Detail: run.OpenStop + "other-q, other-r; answer each with a Decision, then rhiza run resume ROTHER"}); err != nil {
		t.Fatal(err)
	}

	caller := initCtx(t)
	// Its own IDs: the two copies of the init snapshot draw the same ones.
	seededIDs(caller, 4242)
	caller.Getenv = func(k string) string {
		if k == "HOME" {
			return home
		}
		return ""
	}
	caller.Git = func(...string) (string, error) { return o.common + "\n", nil }
	caller.Exec = func(string, []string, string, []string) ([]byte, error) {
		return []byte(o.common + "\n"), nil
	}
	return caller, o
}

// otherIndex is the worktree's default index file under home's user cache
// folder.
func otherIndex(t *testing.T, home, root string) string {
	t.Helper()
	p, err := index.DefaultPath(root, func(k string) string {
		if k == "HOME" {
			return home
		}
		return ""
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(p, filepath.Join(home, ".cache")+string(filepath.Separator)) {
		t.Fatalf("index path %s is not under %s's user cache folder", p, home)
	}
	return p
}

// checkOtherIndex requires the worktree's own index file under home's user
// cache folder (home set), or no user cache folder made at all (home
// empty: the read answered through a throwaway index).
func checkOtherIndex(t *testing.T, caller *Ctx, home string, o otherRoot) {
	t.Helper()
	if home == "" {
		if _, err := index.DefaultPath(o.root, caller.Getenv); err == nil {
			t.Fatal("the worktree has a default index path with no user cache folder")
		}
		return
	}
	if _, err := os.Stat(otherIndex(t, home, o.root)); err != nil {
		t.Errorf("the read did not open the worktree's own index: %v", err)
	}
}

// otherRootCases runs f with the user cache folder under a fresh HOME (the
// worktree's own index) and with none (a throwaway index).
func otherRootCases(t *testing.T, f func(t *testing.T, home string)) {
	t.Run("own-index", func(t *testing.T) {
		t.Parallel()
		f(t, t.TempDir())
	})
	t.Run("throwaway", func(t *testing.T) {
		t.Parallel()
		f(t, "")
	})
}

// The daemon's context for a project (projectCtx) reads the project's own
// index, not the caller's, and answers as a read in that project does.
func TestOtherRootDaemonProject(t *testing.T) {
	t.Parallel()
	otherRootCases(t, func(t *testing.T, home string) {
		caller, o := otherRootFixture(t, home)
		caller.Exec = nil
		pc := projectCtx(caller, o.root)
		if pc.Root != o.root {
			t.Fatalf("projectCtx root %s, want %s", pc.Root, o.root)
		}
		resp, err := callOp(t, pc, "ls", `{"kinds":["intent"]}`)
		if err != nil {
			t.Fatal(err)
		}
		ls := resp.(LsResponse)
		if len(ls.Nodes) != 1 || ls.Nodes[0].Handle != "other-plan-intent" {
			t.Fatalf("ls in the project: %+v, want other-plan-intent", ls.Nodes)
		}
		checkOtherIndex(t, caller, home, o)
	})
}

// The gate read of the inbox (inboxGateRead) reads the run's worktree
// through its own index and answers each question's statement.
func TestOtherRootInboxGateRead(t *testing.T) {
	t.Parallel()
	otherRootCases(t, func(t *testing.T, home string) {
		caller, o := otherRootFixture(t, home)
		st := run.State{Change: o.own.String(), Worktree: o.root, Stopped: run.OpenStop + "other-q, other-r"}
		got := inboxGateQuestions(caller, st)
		if len(got) != 2 {
			t.Fatalf("got %+v, want other-q and other-r", got)
		}
		for i, h := range []string{"other-q", "other-r"} {
			if got[i].Handle != h || got[i].Statement != "Statement of "+h {
				t.Errorf("question %d: %+v, want %s with its statement", i, got[i], h)
			}
		}
		checkOtherIndex(t, caller, home, o)
	})
}

// question.answer --run reads and writes the run's worktree through its
// own index, and the web's gateAnswered reads the answer from it.
func TestOtherRootGateAnswer(t *testing.T) {
	t.Parallel()
	otherRootCases(t, func(t *testing.T, home string) {
		caller, o := otherRootFixture(t, home)
		if done, err := gateAnswered(caller, "ROTHER", "other-q"); err != nil || done {
			t.Fatalf("gateAnswered before the answer = %v, %v; want false", done, err)
		}
		if _, err := callOp(t, caller, "question.answer", `{"question":"other-q","text":"Yes, do it.","run":"ROTHER"}`); err != nil {
			t.Fatal(err)
		}
		w, err := openWorkspace(o.root, caller.Now, caller.IDs, false)
		if err != nil {
			t.Fatal(err)
		}
		defer w.closeProjection()
		if d, err := w.model(o.cs).Resolve("answer-other-q"); err != nil || d.Kind != op.KindDecision {
			t.Errorf("no answering Decision in %s: %+v, %v", o.cs, d, err)
		}
		if _, err := ctxWorkspace(t, caller).model(o.cs).Resolve("answer-other-q"); err == nil {
			t.Error("the answer was written into the caller's root")
		}
		if done, err := gateAnswered(caller, "ROTHER", "other-q"); err != nil || !done {
			t.Errorf("gateAnswered after the answer = %v, %v; want true", done, err)
		}
		if done, err := gateAnswered(caller, "ROTHER", "other-r"); err != nil || done {
			t.Errorf("gateAnswered on the open question = %v, %v; want false", done, err)
		}
		checkOtherIndex(t, caller, home, o)
	})
}

// Another root whose log does not read (no .rhiza) answers as before: the
// inbox's gate read lists the handles with empty statements, a gate answer
// and gateAnswered fail with no Rhiza repository, and a project read fails
// the same way.
func TestOtherRootUnreadable(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	caller, o := otherRootFixture(t, home)
	if err := os.RemoveAll(filepath.Join(o.root, ".rhiza")); err != nil {
		t.Fatal(err)
	}
	st := run.State{Change: o.own.String(), Worktree: o.root, Stopped: run.OpenStop + "other-q, other-r"}
	got := inboxGateQuestions(caller, st)
	if len(got) != 2 || got[0].Statement != "" || got[1].Statement != "" {
		t.Errorf("inbox gate read of an unreadable worktree: %+v, want handles with empty statements", got)
	}
	_, err := callOp(t, caller, "question.answer", `{"question":"other-q","text":"Yes.","run":"ROTHER"}`)
	if err == nil || !strings.Contains(err.Error(), "no Rhiza repository") {
		t.Errorf("gate answer in an unreadable worktree: %v, want no Rhiza repository", err)
	}
	if _, err := gateAnswered(caller, "ROTHER", "other-q"); err == nil || !strings.Contains(err.Error(), "no Rhiza repository") {
		t.Errorf("gateAnswered in an unreadable worktree: %v, want no Rhiza repository", err)
	}
	caller.Exec = nil
	if _, err := callOp(t, projectCtx(caller, o.root), "ls", `{"kinds":["intent"]}`); err == nil || !strings.Contains(err.Error(), "no Rhiza repository") {
		t.Errorf("ls in an unreadable project: %v, want no Rhiza repository", err)
	}
}

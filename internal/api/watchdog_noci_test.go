package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
)

// TestWatchdogPRNoCI drives the no-CI recovery with fake git remotes (local
// bare repos) and a fake gh read: a clean merge is pushed without force and
// recorded in the ship state, a conflicting one is aborted, the branch left
// unchanged, and the run resumed or an inbox item held.
func TestWatchdogPRNoCI(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, "gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	run1 := func(dir string, args ...string) string {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	ctx := &Ctx{Exec: func(name string, args []string, dir string, env []string) ([]byte, error) {
		c := exec.Command(name, args...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		return c.Output()
	}}
	write := func(dir, file, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, file), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// setup makes origin, a main clone and a branch worktree, and moves
	// main on by one commit that touches file with text.
	setup := func(file, text string) shipTarget {
		root := t.TempDir()
		origin := filepath.Join(root, "origin.git")
		run1(root, "init", "-q", "--bare", "-b", "main", origin)
		mainDir := filepath.Join(root, "main")
		run1(root, "clone", "-q", origin, mainDir)
		write(mainDir, "a.txt", "one\n")
		write(mainDir, "b.txt", "one\n")
		run1(mainDir, "add", ".")
		run1(mainDir, "commit", "-q", "-m", "base")
		run1(mainDir, "push", "-q", "origin", "HEAD:main")
		wt := filepath.Join(root, "wt")
		run1(mainDir, "worktree", "add", "-q", "-b", "cs/x", wt)
		write(wt, "a.txt", "branch\n")
		run1(wt, "commit", "-q", "-am", "branch work")
		run1(wt, "push", "-q", "origin", "cs/x")
		write(mainDir, file, text)
		run1(mainDir, "commit", "-q", "-am", "main moves")
		run1(mainDir, "push", "-q", "origin", "HEAD:main")
		cs, _ := id.Parse("01M6WHRFPM0C0X3H6VBFT4R19Y")
		return shipTarget{cs: cs, branch: "cs/x", worktree: wt, main: mainDir, common: filepath.Join(mainDir, ".git")}
	}

	// A clean merge is pushed (never forced) and recorded.
	tc := setup("b.txt", "main\n")
	before := run1(tc.worktree, "rev-parse", "HEAD")
	res, err := watchdogMergeMain(ctx, tc)
	if err != nil || res.Pushed == "" || len(res.Conflicts) > 0 {
		t.Fatalf("clean merge: %+v, %v", res, err)
	}
	if got := run1(tc.main, "ls-remote", "origin", "refs/heads/cs/x"); !strings.HasPrefix(got, res.Pushed) {
		t.Errorf("origin's branch %q, want %s", got, res.Pushed)
	}
	if run1(tc.worktree, "merge-base", "--is-ancestor", before, res.Pushed); run1(tc.worktree, "rev-list", "--merges", "--count", before+".."+res.Pushed) != "1" {
		t.Error("want a merge commit on top of the branch's old head")
	}
	st, err := readShipState(tc.statePath())
	if err != nil || st.Pushed != res.Pushed || st.Head != "" {
		t.Errorf("ship state %+v, %v", st, err)
	}
	if again, err := watchdogMergeMain(ctx, tc); err != nil || again.Pushed != "" {
		t.Errorf("second merge %+v, %v: want nothing to do", again, err)
	}

	// A conflicting merge is aborted and the branch is unchanged.
	tx := setup("a.txt", "main\n")
	before = run1(tx.worktree, "rev-parse", "HEAD")
	remoteBefore := run1(tx.main, "ls-remote", "origin", "refs/heads/cs/x")
	res, err = watchdogMergeMain(ctx, tx)
	if err != nil || res.Pushed != "" || strings.Join(res.Conflicts, ",") != "a.txt" {
		t.Fatalf("conflict: %+v, %v", res, err)
	}
	if got := run1(tx.worktree, "rev-parse", "HEAD"); got != before {
		t.Errorf("branch moved to %s, want %s", got, before)
	}
	if run1(tx.worktree, "status", "--porcelain") != "" {
		t.Error("worktree not clean after the abort")
	}
	if got := run1(tx.main, "ls-remote", "origin", "refs/heads/cs/x"); got != remoteBefore {
		t.Errorf("origin's branch changed: %q", got)
	}
	if _, err := os.Stat(tx.statePath()); err == nil {
		t.Error("ship state written for a conflict")
	}

	// The pass: which runs it takes, and what it does with each outcome.
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	failed := func(id, def string) run.Journal {
		return run.Journal{ID: id, Entries: []run.Entry{
			{Type: run.TypeStarted, At: t0, Started: &run.Started{Autonomous: true, Definition: def}},
			{Type: run.TypeEntered, At: t0, Step: "wait", Visit: 1},
			{Type: run.TypeEnded, At: t0, Step: "wait", Visit: 1, Outcome: "fail", Detail: "no checks"},
			{Type: run.TypeFinished, At: t0, Step: "wait", End: "failed", Detail: "step wait failed"},
		}}
	}
	const handled = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check --settled\"\n"
	const unhandled = "name = \"w\"\n[step.wait]\nkind = \"command\"\nrun = \"rhiza ship check\"\n"
	pr := watchdogPR{Head: "aaaa111", Mergeable: "CONFLICTING"}
	var merges, resumed, log []string
	out := watchdogMerge{Pushed: "bbbb2222"}
	d := watchdogDeps{
		ci:     func(run.State) (watchdogPR, error) { return pr, nil },
		merge:  func(st run.State) (watchdogMerge, error) { merges = append(merges, st.ID); return out, nil },
		resume: func(id string) error { resumed = append(resumed, id); return nil },
		log:    func(f string, a ...any) { log = append(log, fmt.Sprintf(f, a...)) },
	}
	pol := WatchdogPolicy{Enabled: true, RetryBudget: 1, AutonomousOnly: true}
	nolive := func(string) bool { return false }
	pass := func(js ...run.Journal) []daemon.WatchdogHold {
		t.Helper()
		_, holds, err := watchdogPRNoCI(pol, js, nolive, nil, d)
		if err != nil {
			t.Fatal(err)
		}
		return holds
	}

	if h := pass(failed("clean", handled)); len(h) != 0 || strings.Join(merges, ",") != "clean" || strings.Join(resumed, ",") != "clean" {
		t.Fatalf("clean: holds %v merges %v resumed %v", h, merges, resumed)
	}
	// With CI present or the PR mergeable, nothing is done.
	merges, resumed = nil, nil
	pr.Checks = []ShipCheck{{Name: "test", Status: "IN_PROGRESS"}}
	pass(failed("ci", handled))
	pr.Checks, pr.Mergeable = nil, "MERGEABLE"
	pass(failed("ok", handled))
	if len(merges) != 0 {
		t.Fatalf("merged %v for a PR with CI or no conflict", merges)
	}
	pr.Mergeable = "CONFLICTING"
	out = watchdogMerge{Conflicts: []string{"a.txt"}}
	if h := pass(failed("c1", handled)); len(h) != 0 || strings.Join(resumed, ",") != "c1" {
		t.Fatalf("conflict, handled: holds %v resumed %v", h, resumed)
	}
	resumed = nil
	h := pass(failed("c2", unhandled))
	if len(resumed) != 0 || len(h) != 1 || h[0].Run != "c2" || !strings.Contains(h[0].Detail, "a.txt") {
		t.Fatalf("conflict, no handling: holds %+v resumed %v", h, resumed)
	}
	if items := watchdogHoldItems(h); len(items) != 1 || items[0].Kind != InboxStalled {
		t.Errorf("items %+v", items)
	}
	if len(log) == 0 {
		t.Error("no log lines")
	}

	// A head left to the owner keeps its hold on later polls without
	// another merge or log line; a new head is merged again.
	c3 := failed("c3", unhandled)
	mem, h, err := watchdogPRNoCI(pol, []run.Journal{c3}, nolive, nil, d)
	if err != nil || len(h) != 1 {
		t.Fatalf("c3 first poll: holds %+v, %v", h, err)
	}
	merges, log = nil, nil
	for range 3 {
		if mem, h, err = watchdogPRNoCI(pol, []run.Journal{c3}, nolive, mem, d); err != nil || len(h) != 1 || !strings.Contains(h[0].Detail, "a.txt") {
			t.Fatalf("c3 later poll: holds %+v, %v", h, err)
		}
	}
	if len(merges) != 0 || len(log) != 0 {
		t.Fatalf("c3 later polls merged %v, logged %q", merges, log)
	}
	pr.Head = "cccc333"
	if mem, _, _ = watchdogPRNoCI(pol, []run.Journal{c3}, nolive, mem, d); len(merges) != 1 {
		t.Fatalf("new head not merged: %v", merges)
	}

	// With the budget spent it holds once without merging.
	merges, log = nil, nil
	spent := []daemon.WatchdogRun{{Run: "c4", Spent: 1}}
	for range 3 {
		if spent, h, err = watchdogPRNoCI(pol, []run.Journal{failed("c4", handled)}, nolive, spent, d); err != nil || len(h) != 1 {
			t.Fatalf("spent: holds %+v, %v", h, err)
		}
	}
	if len(merges) != 0 || len(log) != 1 {
		t.Fatalf("spent budget merged %v, logged %q", merges, log)
	}
}

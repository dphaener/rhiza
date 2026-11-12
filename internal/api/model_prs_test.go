package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// modelPRFixture is modelBranchFixture with every orphan entry shipped on a
// model branch, and a gh whose pull requests report ci (green, red or
// conflict) and whose pr merge squash-merges the branch into origin/main.
func modelPRFixture(t *testing.T, ci string) *modelFixture {
	t.Helper()
	f := modelBranchFixture(t, true)
	f.ctx.Git = func(args ...string) (string, error) {
		out, err := f.gh.exec("git", args, f.root, nil)
		return string(out), err
	}
	daemonCommitUntracked(f.ctx)
	if len(f.branches(t)) == 0 {
		t.Fatal("no model branches")
	}
	base := f.gh.run
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		k := strings.Join(args, " ")
		if name == "gh" && strings.HasPrefix(k, "pr view ") {
			var checks, mergeable = `[{"name":"ci","status":"COMPLETED","conclusion":"SUCCESS"}]`, "MERGEABLE"
			switch ci {
			case "red":
				checks = `[{"name":"ci","status":"COMPLETED","conclusion":"FAILURE"}]`
			case "conflict":
				mergeable = "CONFLICTING"
			}
			return []byte(`{"number":7,"url":"https://github.com/o/r/pull/7","state":"OPEN","body":"","headRefOid":"x","mergeable":"` + mergeable + `","statusCheckRollup":` + checks + `}`), nil
		}
		if name == "gh" && strings.HasPrefix(k, "pr merge ") {
			f.gh.calls++
			merges++
			scratch := filepath.Join(t.TempDir(), "m")
			f.git(filepath.Dir(scratch), "clone", "-q", f.origin, scratch)
			f.git(scratch, "merge", "--squash", "origin/"+args[2])
			f.git(scratch, "commit", "-q", "-m", "squash "+args[2])
			f.git(scratch, "push", "-q", "origin", "HEAD:main")
			return nil, nil
		}
		return base(name, args, dir, env)
	}
	return f
}

var merges int

func untrackedLog(t *testing.T, f *modelFixture) string {
	t.Helper()
	return f.git(f.root, "ls-files", "--others", "--exclude-standard", ".rhiza/log")
}

// TestDaemonModelPR: with no run, a green model PR is squash-merged, its
// branch deleted and main fast-forwarded with no untracked copy left; a red
// or conflicting one is not merged and is an inbox item naming the change
// set.
func TestDaemonModelPR(t *testing.T) {
	t.Run("Green", func(t *testing.T) {
		merges = 0
		f := modelPRFixture(t, "green")
		if untrackedLog(t, f) == "" {
			t.Fatal("fixture has no untracked entries")
		}
		daemonMergeModelPRs(f.ctx)
		if f.stderr.Len() != 0 {
			t.Fatalf("stderr: %s", f.stderr)
		}
		if merges == 0 {
			t.Fatal("nothing merged")
		}
		if bs := f.branches(t); len(bs) != 0 {
			t.Errorf("branches left: %v", bs)
		}
		if u := untrackedLog(t, f); u != "" {
			t.Errorf("untracked copies left:\n%s", u)
		}
		if h, o := f.git(f.root, "rev-parse", "HEAD"), f.git(f.root, "rev-parse", "origin/main"); h != o {
			t.Errorf("main %s is not origin/main %s", h, o)
		}
		if _, err := os.Stat(filepath.Join(f.root, ".rhiza/log", f.cs)); err != nil {
			t.Errorf("change set entries not in main: %v", err)
		}
		if items := inboxModelPRs(f.ctx); len(items) != 0 {
			t.Errorf("inbox items: %+v", items)
		}
	})
	for _, ci := range []string{"red", "conflict"} {
		t.Run(ci, func(t *testing.T) {
			merges = 0
			f := modelPRFixture(t, ci)
			n := len(f.branches(t))
			daemonMergeModelPRs(f.ctx)
			if merges != 0 || len(f.branches(t)) != n {
				t.Errorf("merged %d, branches %v", merges, f.branches(t))
			}
			var found bool
			for _, it := range inboxModelPRs(f.ctx) {
				if it.Kind == InboxModelPR && it.Change == f.cs && it.Detail != "" {
					found = true
				}
			}
			if !found {
				t.Errorf("no inbox item for %s: %+v", f.cs, inboxModelPRs(f.ctx))
			}
		})
	}
	t.Run("GhFailureNeverFails", func(t *testing.T) {
		f := modelPRFixture(t, "green")
		inner := f.ctx.Exec
		f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
			if name == "gh" {
				return nil, errors.New("gh: not authenticated")
			}
			return inner(name, args, dir, env)
		}
		daemonMergeModelPRs(f.ctx)
	})
}

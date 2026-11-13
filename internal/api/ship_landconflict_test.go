package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// landHeadFails makes the fixture's head build answer land with exit
// status 3 and a --json error body carrying message.
func landHeadFails(t *testing.T, f *sessionFixture, message, status string) {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "rhiza")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	f.ctx.Getenv = func(k string) string {
		if k == envShipRhiza {
			return bin
		}
		return ""
	}
	body := `{"error":{"code":"conflict","message":"` + message + `"}}`
	inner := f.ctx.Exec
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == bin {
			return []byte(body), errors.New(status)
		}
		return inner(name, args, dir, env)
	}
}

// A land the head build refuses with exit 3 because a node changed on main
// merges origin/main into the branch, leaves the conflict in progress and
// recorded for the conflict fixer, and reports a conflict naming the node.
func TestShipLandNodeChangedMergesMain(t *testing.T) {
	f, state, _ := modelConflictFixture(t)
	landHeadFails(t, f, "cannot land change set X: node N changed", "exit status 3")
	for k, v := range map[string]string{
		"fetch -q origin": "",
		"rev-parse --verify --quiet refs/remotes/origin/main": "abc123\n",
		"rev-parse --show-prefix":                             "\n",
		"diff --name-only --diff-filter=U":                    "main.go\n",
	} {
		f.git.out[k] = v
	}
	_, err := callOp(t, f.ctx, "ship.land", `{}`)
	if err == nil {
		t.Fatal("land succeeded")
	}
	e := AsError(err)
	if e.Code != CodeConflict || !strings.Contains(e.Message, "node N changed") || !strings.Contains(e.Message, "main.go") {
		t.Errorf("error = %+v, want a conflict naming node N and the conflicted file", e)
	}
	st, rerr := readShipState(state)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if st.MergeConflict == nil || len(st.MergeConflict.Files) != 1 || st.MergeConflict.Files[0] != "main.go" {
		t.Errorf("merge conflict recorded = %+v", st.MergeConflict)
	}
	if st.ModelConflict != nil {
		t.Errorf("model conflict recorded = %+v, want none", st.ModelConflict)
	}
}

// Any other head-build failure of land still fails with the head build's
// detail, merges nothing and records no conflict.
func TestShipLandOtherHeadFailureStaysPlain(t *testing.T) {
	for _, c := range []struct{ name, msg, status string }{
		{"other message", "cannot land change set X: it depends on Y (draft)", "exit status 3"},
		{"other exit", "cannot land change set X: node N changed", "exit status 1"},
	} {
		t.Run(c.name, func(t *testing.T) {
			f, state, _ := modelConflictFixture(t)
			landHeadFails(t, f, c.msg, c.status)
			_, err := callOp(t, f.ctx, "ship.land", `{}`)
			if err == nil {
				t.Fatal("land succeeded")
			}
			if !strings.Contains(err.Error(), c.msg) {
				t.Errorf("error %q lacks the head build's detail %q", err, c.msg)
			}
			for _, call := range f.git.calls {
				if strings.HasPrefix(call, "merge") {
					t.Errorf("git %s ran", call)
				}
			}
			if st, _ := readShipState(state); st.MergeConflict != nil {
				t.Errorf("merge conflict recorded = %+v", st.MergeConflict)
			}
		})
	}
}

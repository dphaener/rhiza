package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// scriptedStart runs each started program as the next scripted result.
func scriptedStart(results []fakeProcess, n *int) func(Program) (Process, error) {
	return func(Program) (Process, error) {
		p := results[*n]
		*n++
		return &p, nil
	}
}

func TestCommandRerunsOnTransient(t *testing.T) {
	fail := func(s string) fakeProcess { return fakeProcess{out: strings.NewReader(s), code: 1} }
	ok := func() fakeProcess { return fakeProcess{out: strings.NewReader("ok\n"), code: 0} }
	transient := func(out string) bool { return strings.Contains(out, "flaky") }
	for _, tc := range []struct {
		name     string
		classify func(string) bool
		results  []fakeProcess
		runs     int
		code     int
		notes    int
	}{
		{"transient then passes", transient, []fakeProcess{fail("flaky\n"), ok()}, 2, 0, 1},
		{"transient twice stands", transient, []fakeProcess{fail("flaky\n"), fail("flaky again\n")}, 2, 1, 1},
		{"real not rerun", transient, []fakeProcess{fail("compile error\n")}, 1, 1, 0},
		{"no classifier not rerun", nil, []fakeProcess{fail("flaky\n")}, 1, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testCtx(t.TempDir())
			var stderr bytes.Buffer
			ctx.Stderr = &stderr
			ctx.Classify = tc.classify
			n := 0
			ctx.Start = scriptedStart(tc.results, &n)
			r := &stepRunner{ctx: ctx, worktree: ctx.Root, main: ctx.Root}
			code, _, err := r.Command(workflow.Step{Run: "go build ./..."}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if code != tc.code || n != tc.runs {
				t.Errorf("code %d after %d runs, want %d after %d", code, n, tc.code, tc.runs)
			}
			if got := strings.Count(stderr.String(), "running it once more"); got != tc.notes {
				t.Errorf("%d retry lines, want %d: %q", got, tc.notes, stderr.String())
			}
		})
	}
}

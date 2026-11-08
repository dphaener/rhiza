package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func TestVanishedGoWork(t *testing.T) {
	for _, tc := range []struct {
		name, out string
		want      bool
	}{
		{"import", "x.go:3:8: could not import io (open $WORK/b033/_pkg_.a: no such file or directory)", true},
		{"create", "can't create $WORK/b094/_pkg_.a: open $WORK/b094/_pkg_.a: no such file or directory", true},
		{"compile error", "./x.go:3:2: undefined: foo", false},
		{"missing file outside work", "open /home/x/y.txt: no such file or directory", false},
		{"work other error", "open $WORK/b033/_pkg_.a: permission denied", false},
	} {
		if got := vanishedGoWork(tc.out); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// scriptedStart runs each started program as the next scripted result.
func scriptedStart(results []fakeProcess, n *int) func(Program) (Process, error) {
	return func(Program) (Process, error) {
		p := results[*n]
		*n++
		return &p, nil
	}
}

func TestCommandRetriesVanishedGoWork(t *testing.T) {
	sig := "could not import io (open $WORK/b033/_pkg_.a: no such file or directory)\n"
	for _, tc := range []struct {
		name    string
		results []fakeProcess
		runs    int
		code    int
		notes   int
	}{
		{"once then passes", []fakeProcess{{out: strings.NewReader(sig), code: 1}, {out: strings.NewReader("ok\n"), code: 0}}, 2, 0, 1},
		{"twice stands", []fakeProcess{{out: strings.NewReader(sig), code: 1}, {out: strings.NewReader(sig + "again\n"), code: 2}}, 2, 2, 1},
		{"other failure", []fakeProcess{{out: strings.NewReader("compile error\n"), code: 1}}, 1, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testCtx(t.TempDir())
			var stderr bytes.Buffer
			ctx.Stderr = &stderr
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

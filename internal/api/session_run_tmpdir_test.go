package api

import (
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// tmpdirSteps is a workflow body with one leaves step and no [env].
const tmpdirSteps = "[step.s]\nkind = \"leaves\"\n[step.s.builder]\nprompt = \"b\"\n[step.s.verifier]\nprompt = \"v\"\n"

// startTmpdirRun starts run id on f's change set with the workflow src and
// returns the run's temp folder, removed when the test ends.
func startTmpdirRun(t *testing.T, f *sessionFixture, change, id, src string) string {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	rs.Getenv = f.ctx.Getenv
	now := time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC)
	if err := run.Start(rs, id, run.Started{Workflow: "wf", Definition: src, Change: change, Worktree: f.worktree}, now, def); err != nil {
		t.Fatal(err)
	}
	tmp := rs.TmpDir(id)
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Fatalf("run temp folder exists before the session runs: %v", err)
	}
	return tmp
}

// envValues is every value env gives key, in order.
func envValues(env []string, key string) []string {
	var out []string
	for _, kv := range env {
		if k, v, ok := strings.Cut(kv, "="); ok && k == key {
			out = append(out, v)
		}
	}
	return out
}

// A session of a run whose workflow has no [env], or an [env] that names
// no {run_tmp}, gets TMPDIR=<run tmp> as the run's steps do: a dry run
// shows it and makes nothing, a real start makes the folder first.
//
// Serial: each run's temp folder lives under the cache folder every test
// in the binary shares.
func TestSessionRunTmpdirDefault(t *testing.T) {
	for _, tc := range []struct{ name, id, src string }{
		{"no env", "RT-noenv", "name = \"wf\"\n" + tmpdirSteps},
		{"env without run_tmp", "RT-plainenv", "name = \"wf\"\n[env]\nGOFLAGS = \"-count=1\"\n" + tmpdirSteps},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, s, fs := runFixture(t, `{"permission_mode":"manual","run":"`+tc.id+`"}`)
			tmp := startTmpdirRun(t, f, s.Change, tc.id, tc.src)

			dry := *f.ctx
			dry.DryRun = true
			resp, err := callOp(t, &dry, "session.run", `{"id":"`+s.ID+`"}`)
			if err != nil {
				t.Fatal(err)
			}
			if got := envValues(resp.(SessionRunResponse).Env, "TMPDIR"); !slices.Equal(got, []string{tmp}) {
				t.Errorf("dry run TMPDIR %v, want [%s]", got, tmp)
			}
			if _, err := os.Stat(tmp); !os.IsNotExist(err) {
				t.Errorf("a dry run made the run temp folder (stat err %v)", err)
			}

			if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
				t.Fatal(err)
			}
			agent := fs.programs[len(fs.programs)-1]
			if got := envValues(agent.Env, "TMPDIR"); !slices.Equal(got, []string{tmp}) {
				t.Errorf("agent TMPDIR %v, want [%s]", got, tmp)
			}
			if fi, err := os.Stat(tmp); err != nil || !fi.IsDir() {
				t.Errorf("run temp folder %s not made: %v", tmp, err)
			}
		})
	}
}

// A workflow whose [env] sets TMPDIR gives the session that value only,
// and no run temp folder is made for it.
//
// Serial: see TestSessionRunTmpdirDefault.
func TestSessionRunTmpdirOwn(t *testing.T) {
	f, s, fs := runFixture(t, `{"permission_mode":"manual","run":"RT-own"}`)
	tmp := startTmpdirRun(t, f, s.Change, "RT-own", "name = \"wf\"\n[env]\nTMPDIR = \"/x\"\n"+tmpdirSteps)

	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	agent := fs.programs[len(fs.programs)-1]
	if got := envValues(agent.Env, "TMPDIR"); !slices.Equal(got, []string{"/x"}) {
		t.Errorf("agent TMPDIR %v, want [/x]", got)
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Errorf("run temp folder %s made for a workflow that sets TMPDIR (stat err %v)", tmp, err)
	}
}

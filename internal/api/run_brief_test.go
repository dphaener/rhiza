package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// envOf is the value launchEnv sets for name, and whether it sets it.
func envOf(rec session.Record, name string) (string, bool) {
	set, _ := launchEnv(rec)
	for _, kv := range set {
		if v, ok := strings.CutPrefix(kv, name+"="); ok {
			return v, true
		}
	}
	return "", false
}

// An agent step's, a builder's and a checks step's verifier's session each
// get their own directory with brief.md and slice.json, slice.json the
// bytes `rhiza slice --json` prints, RHIZA_BRIEF and RHIZA_SLICE set, and
// RHIZA_HANDBACK only once handback.md exists; a session started outside a
// run gets none.
func TestSessionBrief(t *testing.T) {
	f, r := paneRunner(t)
	// The run's worktree holds the change set, as a real one does.
	csID, err := id.Parse(r.change)
	if err != nil {
		t.Fatal(err)
	}
	if err := copyChangeSetLog(f.ctx.Root, f.worktree, csID); err != nil {
		t.Fatal(err)
	}
	// slice is what rhiza slice --json prints for the change set now.
	slice := func() []byte {
		sc := *testCtx(f.worktree)
		sc.ChangeRef = r.change
		sl, err := callOp(t, &sc, "slice", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		b, err := canon.Marshal(sl)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	leaf := run.Leaf{Handle: "leaf-a"}
	starts := map[string]func() (string, error){
		"agent": func() (string, error) {
			return r.StartAgent(workflow.Step{Name: "plan", Prompt: "plan"}, "", nil)
		},
		"builder": func() (string, error) {
			return r.StartLeaf(run.LeafSpec{Step: workflow.Step{Name: "build", Builder: &workflow.Role{Prompt: "b"}},
				Role: workflow.RoleBuilder, Leaf: leaf, Attempt: 1})
		},
		"verifier": func() (string, error) {
			return r.StartCheck(run.CheckSpec{Step: workflow.Step{Name: "review", Verifier: &workflow.Role{Prompt: "v"}},
				Role: workflow.RoleVerifier, Check: workflow.Check{Name: "alpha"}, Cycle: 1})
		},
	}
	dirs := map[string]bool{}
	for name, start := range starts {
		want := slice()
		sid, err := start()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		rec, err := r.sessions.Read(sid)
		if err != nil {
			t.Fatal(err)
		}
		brief, ok := envOf(rec, envBrief)
		if !ok || filepath.Base(brief) != "brief.md" {
			t.Fatalf("%s: RHIZA_BRIEF = %q, %v", name, brief, ok)
		}
		dir := filepath.Dir(brief)
		if dirs[dir] || !strings.HasPrefix(dir, filepath.Join(f.common, "rhiza", "runs", "R1")) {
			t.Errorf("%s: directory %s is shared or outside the run's", name, dir)
		}
		dirs[dir] = true
		if sp, _ := envOf(rec, envSlice); sp != filepath.Join(dir, "slice.json") {
			t.Errorf("%s: RHIZA_SLICE = %q", name, sp)
		}
		got, err := os.ReadFile(filepath.Join(dir, "slice.json"))
		if err != nil || string(got) != string(want) {
			t.Errorf("%s: slice.json differs from rhiza slice --json (%v)\n got %s\nwant %s", name, err, got, want)
		}
		md, err := os.ReadFile(brief)
		if err != nil || !strings.Contains(string(md), "Readable tabs") || !strings.HasSuffix(string(md), run.ContractLastLine+"\n") {
			t.Errorf("%s: brief.md %q (%v)", name, md, err)
		}
		if _, ok := envOf(rec, envHandback); ok {
			t.Errorf("%s: RHIZA_HANDBACK set with no handback.md", name)
		}
		if err := os.WriteFile(filepath.Join(dir, "handback.md"), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if hb, ok := envOf(rec, envHandback); !ok || hb != filepath.Join(dir, "handback.md") {
			t.Errorf("%s: RHIZA_HANDBACK = %q, %v with handback.md", name, hb, ok)
		}
	}
	if rec := (session.Record{ID: "S-out", Change: r.change}); true {
		for _, n := range []string{envBrief, envSlice, envHandback} {
			if _, ok := envOf(rec, n); ok {
				t.Errorf("a session outside a run got %s", n)
			}
		}
	}
}

// An unwritable runs directory fails the start with its path and opens no
// tab.
func TestSessionBriefUnwritable(t *testing.T) {
	f, r := paneRunner(t)
	if err := os.MkdirAll(filepath.Join(f.common, "rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.common, "rhiza", "runs"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
	_, err := r.StartAgent(workflow.Step{Name: "plan", Prompt: "plan"}, "", nil)
	if err == nil || !strings.Contains(err.Error(), filepath.Join(f.common, "rhiza", "runs")) {
		t.Fatalf("error = %v, want one naming the runs directory", err)
	}
	if len(f.rex.calls) != 0 {
		t.Errorf("rex called: %q", f.rex.calls)
	}
}

package api

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// fixLeafWF is a review over two checks whose fixer prompt names {leaf}
// and {cycle}.
const fixLeafWF = `[step.review]
kind = "checks"
[[step.review.check]]
name = "alpha"
prompt = "check {check}"
[[step.review.check]]
name = "beta"
prompt = "check {check}"
[step.review.critic]
prompt = "criticise {check}"
[step.review.fixer]
prompt = "fix leaf {leaf} in cycle {cycle} of {change}"
[step.review.verifier]
prompt = "verify {leaf}"
`

// fixLeafRunner is the stepRunner on a real git worktree holding a Rhiza
// repository, with leaf "leaf-a", its sessions faked: in cycle 1 check
// alpha finds a major on leaf-a, which alpha's critic confirms, and beta
// finds nothing; in cycle 2 alpha finds nothing. The fixer, as it starts,
// commits a fix naming leaf-a unless noCommit; the verifier writes
// Evidence on leaf-a as ada, whose session S-v is its own. Every session
// exits 0 at its first poll. The real git machinery judges the fixer's
// commit and the verifier's Evidence, and opens their checkouts.
type fixLeafRunner struct {
	*stepRunner
	t        *testing.T
	git      func(...string) string
	ctx      *Ctx
	store    *run.Store
	leaf     run.Leaf
	ada      id.ID
	noCommit bool
	checks   []string // "check role cycle" per StartCheck call
	specs    []run.LeafSpec
	prompts  []string // each leaf session's rendered prompt
	heads    []string // each verifier checkout's HEAD as it started
	fix      string   // the fixer's commit
}

func (g *fixLeafRunner) Leaves(string) ([]run.Leaf, error) { return []run.Leaf{g.leaf}, nil }

func (g *fixLeafRunner) StartCheck(spec run.CheckSpec) (string, error) {
	g.checks = append(g.checks, spec.Check.Name+" "+spec.Role+" "+strconv.Itoa(spec.Cycle))
	e := run.Entry{Step: spec.Step.Name, Check: spec.Check.Name, Cycle: spec.Cycle, At: noSleep{}.Now()}
	switch {
	case spec.Role == workflow.RoleCheck && spec.Check.Name == "alpha" && spec.Cycle == 1:
		e.Type, e.Findings = run.TypeFindings, []run.Finding{{
			ID: "f1", Check: "alpha", Severity: run.SeverityMajor, Leaf: g.leaf.ID, LeafHandle: g.leaf.Handle,
			Place: "fix.go:1", Reason: "the map is written before it is made",
		}}
	case spec.Role == workflow.RoleCheck:
		e.Type, e.Findings = run.TypeFindings, []run.Finding{}
	default:
		e.Type, e.Verdicts = run.TypeVerdict, []run.Verdict{{Finding: "f1", Verdict: run.Confirmed, Reason: "it panics"}}
	}
	if err := g.store.Append("R1", e); err != nil {
		return "", err
	}
	return "S-" + spec.Check.Name + "-" + spec.Role + "-" + strconv.Itoa(spec.Cycle), nil
}

func (g *fixLeafRunner) StartLeaf(spec run.LeafSpec) (string, error) {
	t := g.t
	g.specs = append(g.specs, spec)
	g.prompts = append(g.prompts, leafPrompt(*spec.Step.Role(spec.Role), spec, workflow.PromptVars{Change: g.change}))
	switch spec.Role {
	case workflow.RoleFixer:
		if g.noCommit {
			return "S-fixer", nil
		}
		if err := os.WriteFile(filepath.Join(g.worktree, "fix.go"), []byte("package fix\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		g.git("add", "fix.go")
		g.git("commit", "-q", "-m", "fix: make the map first", "-m", "Change: "+g.change+"\nLeaf: "+g.leaf.Handle)
		g.fix = g.git("rev-parse", "HEAD")
		return "S-fixer", nil
	case workflow.RoleVerifier:
		g.heads = append(g.heads, g.git("-C", spec.Dir, "rev-parse", "HEAD"))
		v := *g.ctx
		v.Actor = g.ada
		mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-fixed","fields":["statement=leaf-a's fix works","result=PASS"]}`)
		mustWrite(t, &v, "link", `{"from":"leaf-a-fixed","type":"verifies","to":"leaf-a"}`)
		return "S-v", nil
	}
	t.Fatalf("unexpected leaf session %+v", spec)
	return "", nil
}

func (g *fixLeafRunner) AgentStatus(string) (run.AgentStatus, error) {
	return run.AgentStatus{Done: true, Outcome: "pass", Detail: "exit 0"}, nil
}

func (g *fixLeafRunner) AgentEnded(workflow.Step, string) error { return nil }

// ClosePane closes nothing: its sessions are fakes with no pane.
func (g *fixLeafRunner) ClosePane(string) error { return nil }

// fixLeafRun runs the review of fixLeafWF with g's sessions and returns
// the folded journal and its entries.
func fixLeafRun(t *testing.T, noCommit bool) (*fixLeafRunner, run.State, []run.Entry) {
	t.Helper()
	sr, git, ctx, leaf, ada := evidenceCommitRunner(t)
	def, err := workflow.Parse("review", []byte(fixLeafWF))
	if err != nil {
		t.Fatal(err)
	}
	store := run.NewStore(filepath.Join(t.TempDir(), "runs"))
	now := noSleep{}.Now()
	if err := run.Start(store, "R1", run.Started{Workflow: "review", Definition: fixLeafWF, Change: sr.change}, now, def); err != nil {
		t.Fatal(err)
	}
	g := &fixLeafRunner{stepRunner: sr, t: t, git: git, ctx: ctx, store: store, leaf: leaf, ada: ada, noCommit: noCommit}
	e := &run.Engine{Store: store, Runner: g, Clock: noSleep{}, Now: func() time.Time { return now }}
	st, err := e.Run("R1")
	if err != nil {
		t.Fatal(err)
	}
	j, err := store.Read("R1")
	if err != nil {
		t.Fatal(err)
	}
	return g, st, j.Entries
}

// The acceptance of leaf review-fix-leaf, against a real git repository:
// a confirmed finding naming leaf A goes to a fixer whose prompt names A
// and the finding; its new commit carries `Leaf: A`; A's verifier starts
// in a checkout of that commit, and its Evidence, written by the
// verifier's actor, is committed alone; then only the failed check runs
// again, and the step passes.
func TestFixLeafEngine(t *testing.T) {
	g, st, entries := fixLeafRun(t, false)
	if len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomePass {
		t.Fatalf("steps = %+v", st.Steps)
	}
	if len(g.specs) != 2 || g.specs[0].Role != workflow.RoleFixer || g.specs[1].Role != workflow.RoleVerifier {
		t.Fatalf("leaf sessions = %+v", g.specs)
	}
	fixer := g.prompts[0]
	t.Logf("fixer prompt:\n%s", fixer)
	for _, want := range []string{"fix leaf leaf-a in cycle 2 of " + g.change, "alpha f1", "major", "fix.go:1",
		"the map is written before it is made", "`Leaf: leaf-a`"} {
		if !strings.Contains(fixer, want) {
			t.Errorf("fixer prompt lacks %q:\n%s", want, fixer)
		}
	}
	if g.specs[0].Dir != "" {
		t.Errorf("the fixer works in %q, want the run's worktree", g.specs[0].Dir)
	}

	if g.fix == "" {
		t.Fatal("the fixer made no commit")
	}
	if tr := g.git("log", "-1", "--format=%(trailers:key=Leaf,valueonly)", g.fix); tr != "leaf-a" {
		t.Errorf("fix commit's Leaf: trailer = %q", tr)
	}
	l, ok := st.Leaf("review", g.leaf.ID)
	if !ok {
		t.Fatalf("no leaf state for leaf-a on the checks step: %+v", st.Leaves)
	}
	if j, _ := l.Judged(workflow.RoleFixer, 2); j.Outcome != workflow.OutcomePass || j.Commit != g.fix {
		t.Errorf("fixer judgement = %+v, want pass on %s", j, g.fix)
	}
	if v := g.specs[1]; v.Commit != g.fix || !reflect.DeepEqual(g.heads, []string{g.fix}) {
		t.Errorf("verifier started on %s with checkout HEAD %v, want the fix commit %s", v.Commit, g.heads, g.fix)
	}
	if _, err := os.Stat(g.specs[1].Dir); !os.IsNotExist(err) {
		t.Errorf("the verifier's checkout %s is still there: %v", g.specs[1].Dir, err)
	}
	vj, _ := l.Judged(workflow.RoleVerifier, 2)
	if vj.Outcome != workflow.OutcomePass || vj.Detail != "Evidence written by the verifier" {
		t.Fatalf("verifier judgement = %+v", vj)
	}
	actor := nodeName(ctxWorkspace(t, g.ctx).model(id.ID{}), g.ada)
	if tr := g.git("log", "-1", "--format=%(trailers:key=Actor,valueonly)", vj.Commit); tr != actor {
		t.Errorf("Evidence commit %s's Actor: trailer = %q, want the verifier's %q", vj.Commit, tr, actor)
	}
	if head := g.git("rev-parse", "HEAD"); head != vj.Commit {
		t.Errorf("HEAD = %s, want the Evidence commit %s", head, vj.Commit)
	}

	if want := []string{"alpha check 1", "beta check 1", "alpha critic 1", "alpha check 2"}; !reflect.DeepEqual(g.checks, want) {
		t.Errorf("check sessions = %v, want %v", g.checks, want)
	}
	judged, recheck := -1, -1
	for i, e := range entries {
		if e.Type == run.TypeLeafAttempt && e.Role == workflow.RoleVerifier {
			judged = i
		}
		if e.Type == run.TypeCheckSession && e.Cycle == 2 && recheck < 0 {
			recheck = i
		}
		if e.Type == run.TypeLeaf {
			t.Errorf("a leaf outcome was journalled: %+v", e)
		}
	}
	if judged < 0 || recheck < judged {
		t.Errorf("verifier judged at entry %d, cycle 2's check started at %d", judged, recheck)
	}
	c, _ := st.Check("review", "alpha")
	if o, _ := c.Outcome(2); o.Outcome != run.CheckPassed {
		t.Errorf("alpha's cycle 2 outcome = %+v", o)
	}
}

// A fixer that ends with no new commit naming its leaf fails the cycle:
// no verifier starts and the failed check runs no session in it.
func TestFixLeafEngineNoCommit(t *testing.T) {
	g, st, _ := fixLeafRun(t, true)
	if len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomeFail {
		t.Fatalf("steps = %+v", st.Steps)
	}
	for _, s := range g.specs {
		if s.Role != workflow.RoleFixer {
			t.Errorf("a %s started after the fixer made no commit", s.Role)
		}
	}
	if want := []string{"alpha check 1", "beta check 1", "alpha critic 1"}; !reflect.DeepEqual(g.checks, want) {
		t.Errorf("check sessions = %v, want %v", g.checks, want)
	}
	l, _ := st.Leaf("review", g.leaf.ID)
	j, _ := l.Judged(workflow.RoleFixer, 2)
	if j.Outcome != workflow.OutcomeFail || !strings.HasPrefix(j.Detail, "no new commit since ") {
		t.Errorf("fixer judgement = %+v", j)
	}
	c, _ := st.Check("review", "alpha")
	if o, _ := c.Outcome(2); o.Outcome != run.CheckFailed || o.Detail != "cycle 2: fix failed: leaf leaf-a: fixer failed: "+j.Detail {
		t.Errorf("alpha's cycle 2 outcome = %+v", o)
	}
}

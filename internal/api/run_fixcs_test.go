package api

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// A change-set fixer's commit is the newest commit since its start HEAD
// whose Review: trailer names the run. A Review: trailer naming another
// run, or other trailers alone, do not count, and the detail says why.
func TestReviewCommit(t *testing.T) {
	setup := func(t *testing.T) (*stepRunner, func(...string) string, string) {
		r, git := leafCommitRunner(t)
		r.run = "R1"
		head, err := r.Head()
		if err != nil {
			t.Fatal(err)
		}
		return r, git, head
	}
	t.Run("none", func(t *testing.T) {
		r, _, head := setup(t)
		c, detail, err := r.ReviewCommit(head)
		if err != nil || c != "" || detail != "no new commit since "+head[:12] {
			t.Errorf("ReviewCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("other run", func(t *testing.T) {
		r, git, head := setup(t)
		git("commit", "-q", "--allow-empty", "-m", "fix", "-m", "Change: CS1\nReview: R2")
		c, detail, err := r.ReviewCommit(head)
		if err != nil || c != "" || detail != "new commits' Review: trailers name another run: R2" {
			t.Errorf("ReviewCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("other trailers", func(t *testing.T) {
		r, git, head := setup(t)
		git("commit", "-q", "--allow-empty", "-m", "fix", "-m", "Change: CS1\nLeaf: R1\nActor: R1")
		c, detail, err := r.ReviewCommit(head)
		if err != nil || c != "" || detail != "1 new commit(s) since "+head[:12]+", none with a Review: trailer" {
			t.Errorf("ReviewCommit = %q, %q, %v", c, detail, err)
		}
	})
	t.Run("found", func(t *testing.T) {
		r, git, head := setup(t)
		git("commit", "-q", "--allow-empty", "-m", "fix", "-m", "Change: CS1\nreview:  R1 ")
		want := git("rev-parse", "HEAD")
		git("commit", "-q", "--allow-empty", "-m", "later", "-m", "Review: R2")
		c, detail, err := r.ReviewCommit(head)
		if err != nil || c != want || detail != "" {
			t.Errorf("ReviewCommit = %q, %q, %v; want %s", c, detail, err, want)
		}
	})
	t.Run("only since", func(t *testing.T) {
		r, git, _ := setup(t)
		git("commit", "-q", "--allow-empty", "-m", "old fix", "-m", "Review: R1")
		head, _ := r.Head()
		c, _, err := r.ReviewCommit(head)
		if err != nil || c != "" {
			t.Errorf("ReviewCommit = %q, %v; an old commit counted", c, err)
		}
	})
}

// A change-set fixer's session gets RHIZA_CYCLE and RHIZA_ROLE, and no
// RHIZA_LEAF, RHIZA_ATTEMPT or RHIZA_CHECK.
func TestFixChangeSetEnv(t *testing.T) {
	set, unset := launchEnv(session.Record{ID: "S1", Change: "CS1", Run: "R1", Cycle: 2, Role: workflow.RoleFixer})
	for _, v := range []string{"RHIZA_CYCLE=2", "RHIZA_ROLE=fixer", "RHIZA_RUN=R1"} {
		if !slices.Contains(set, v) {
			t.Errorf("set %v lacks %s", set, v)
		}
	}
	for _, v := range []string{"RHIZA_LEAF", "RHIZA_ATTEMPT", "RHIZA_CHECK"} {
		if !slices.Contains(unset, v) {
			t.Errorf("unset %v lacks %s", unset, v)
		}
	}
}

// fixCSRunner is the stepRunner on a real git worktree holding a Rhiza
// repository, its sessions faked: in cycle 1 check alpha finds a major
// and beta a blocker, neither naming a leaf, which their critics confirm;
// in cycle 2 both find nothing. The change set's fixer, as it starts,
// commits a fix with the trailer `Review: <review>`. Every session exits
// 0 at its first poll.
type fixCSRunner struct {
	*stepRunner
	t       *testing.T
	git     func(...string) string
	store   *run.Store
	review  string // the run its fix commit's Review: trailer names
	checks  []string
	specs   []run.FixSpec
	prompts []string
	fix     string // the fixer's commit
}

func (g *fixCSRunner) Leaves(string) ([]run.Leaf, error) { return nil, nil }

func (g *fixCSRunner) StartCheck(spec run.CheckSpec) (string, error) {
	g.checks = append(g.checks, spec.Check.Name+" "+spec.Role+" "+strconv.Itoa(spec.Cycle))
	e := run.Entry{Step: spec.Step.Name, Check: spec.Check.Name, Cycle: spec.Cycle, At: noSleep{}.Now()}
	switch {
	case spec.Role == workflow.RoleCheck && spec.Cycle == 1 && spec.Check.Name == "alpha":
		e.Type, e.Findings = run.TypeFindings, []run.Finding{{
			ID: "f1", Check: "alpha", Severity: run.SeverityMajor, Place: "README.md:3", Reason: "the install line is wrong",
		}}
	case spec.Role == workflow.RoleCheck && spec.Cycle == 1:
		e.Type, e.Findings = run.TypeFindings, []run.Finding{{
			ID: "f1", Check: "beta", Severity: run.SeverityBlocker, Place: "go.mod:1", Reason: "the module path is stale",
		}}
	case spec.Role == workflow.RoleCheck:
		e.Type, e.Findings = run.TypeFindings, []run.Finding{}
	default:
		e.Type, e.Verdicts = run.TypeVerdict, []run.Verdict{{Finding: "f1", Verdict: run.Confirmed, Reason: "real"}}
	}
	if err := g.store.Append("R1", e); err != nil {
		return "", err
	}
	return "S-" + spec.Check.Name + "-" + spec.Role + "-" + strconv.Itoa(spec.Cycle), nil
}

func (g *fixCSRunner) StartLeaf(spec run.LeafSpec) (string, error) {
	g.t.Fatalf("a leaf session started: %+v", spec)
	return "", nil
}

func (g *fixCSRunner) StartFix(spec run.FixSpec) (string, error) {
	g.specs = append(g.specs, spec)
	g.prompts = append(g.prompts, changeSetFixerPrompt(*spec.Step.Role(workflow.RoleFixer), spec, workflow.PromptVars{Change: g.change, Run: g.run}))
	if err := os.WriteFile(filepath.Join(g.worktree, "README.md"), []byte("fixed in cycle "+strconv.Itoa(spec.Cycle)+" attempt "+strconv.Itoa(spec.Attempt)+"\n"), 0o644); err != nil {
		g.t.Fatal(err)
	}
	g.git("add", "README.md")
	g.git("commit", "-q", "-m", "fix: the review's findings", "-m", "Change: "+g.change+"\nReview: "+g.review)
	g.fix = g.git("rev-parse", "HEAD")
	return "S-fixer-" + strconv.Itoa(spec.Cycle) + map[bool]string{true: "", false: "-r" + strconv.Itoa(spec.Attempt)}[spec.Attempt == 1], nil
}

func (g *fixCSRunner) AgentStatus(string) (run.AgentStatus, error) {
	return run.AgentStatus{Done: true, Outcome: "pass", Detail: "exit 0"}, nil
}

func (g *fixCSRunner) AgentEnded(workflow.Step, string) error { return nil }

func (g *fixCSRunner) ClosePane(string) error { return nil }

// fixCSRun runs the review of fixLeafWF with g's sessions, the fixer's
// commit naming run review, and returns the folded journal and its
// entries.
func fixCSRun(t *testing.T, review string) (*fixCSRunner, run.State, []run.Entry) {
	t.Helper()
	sr, git, _, _, _ := evidenceCommitRunner(t)
	def, err := workflow.Parse("review", []byte(fixLeafWF))
	if err != nil {
		t.Fatal(err)
	}
	store := run.NewStore(filepath.Join(t.TempDir(), "runs"))
	now := noSleep{}.Now()
	if err := run.Start(store, "R1", run.Started{Workflow: "review", Definition: fixLeafWF, Change: sr.change}, now, def); err != nil {
		t.Fatal(err)
	}
	g := &fixCSRunner{stepRunner: sr, t: t, git: git, store: store, review: review}
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

// The acceptance of leaf review-fix-change-set, against a real git
// repository: two confirmed findings naming no leaf go to a single fixer
// session, whose prompt names both; its new commit carries a Review:
// trailer naming the run; and the failed checks re-run after it.
func TestFixChangeSetEngine(t *testing.T) {
	g, st, entries := fixCSRun(t, "R1")
	if len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomePass {
		t.Fatalf("steps = %+v", st.Steps)
	}
	if len(g.specs) != 1 {
		t.Fatalf("change-set fixer sessions = %+v, want one", g.specs)
	}
	p := g.prompts[0]
	t.Logf("fixer prompt:\n%s", p)
	for _, want := range []string{"fix leaf  in cycle 2 of " + g.change, "alpha f1 (major at README.md:3): the install line is wrong",
		"beta f1 (blocker at go.mod:1): the module path is stale", "`Review: R1`"} {
		if !strings.Contains(p, want) {
			t.Errorf("fixer prompt lacks %q:\n%s", want, p)
		}
	}
	if tr := g.git("log", "-1", "--format=%(trailers:key=Review,valueonly)", g.fix); tr != "R1" {
		t.Errorf("fix commit's Review: trailer = %q", tr)
	}
	fix, ok := st.ReviewFix("review", 2)
	if !ok || fix.Outcome != workflow.OutcomePass || fix.Commit != g.fix {
		t.Errorf("fixer = %+v, want pass on %s", fix, g.fix)
	}
	if want := []string{"alpha check 1", "beta check 1", "alpha critic 1", "beta critic 1", "alpha check 2", "beta check 2"}; !reflect.DeepEqual(g.checks, want) {
		t.Errorf("check sessions = %v, want %v", g.checks, want)
	}
	judged, recheck := -1, -1
	for i, e := range entries {
		if e.Type == run.TypeReviewFixJudged {
			judged = i
		}
		if e.Type == run.TypeCheckSession && e.Cycle == 2 && recheck < 0 {
			recheck = i
		}
	}
	if judged < 0 || recheck < judged {
		t.Errorf("fixer judged at entry %d, cycle 2's check started at %d", judged, recheck)
	}

	var b strings.Builder
	runInfo(st, false, nil).Human(&b)
	if !strings.Contains(b.String(), "## Change-set fixes") || !strings.Contains(b.String(), "cycle 2 (step review): fixer session S-fixer-2: pass, commit "+g.fix[:12]) {
		t.Errorf("run status lacks the change-set fixer:\n%s", b.String())
	}
	if info := runInfo(st, false, nil); len(info.Fixes) != 1 || info.Fixes[0].Commit != g.fix {
		t.Errorf("RunInfo.Fixes = %+v", info.Fixes)
	}
}

// A change-set fixer ending without a commit whose Review: trailer names
// the run (here it names another) fails the cycle: the failed checks run
// no session in it.
func TestFixChangeSetEngineNoCommit(t *testing.T) {
	g, st, _ := fixCSRun(t, "R2")
	if len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomeFail {
		t.Fatalf("steps = %+v", st.Steps)
	}
	if want := []string{"alpha check 1", "beta check 1", "alpha critic 1", "beta critic 1"}; !reflect.DeepEqual(g.checks, want) {
		t.Errorf("check sessions = %v, want %v", g.checks, want)
	}
	fix, _ := st.ReviewFix("review", 2)
	if fix.Outcome != workflow.OutcomeFail || fix.Detail != "new commits' Review: trailers name another run: R2" {
		t.Errorf("fixer = %+v", fix)
	}
	c, _ := st.Check("review", "alpha")
	if o, _ := c.Outcome(2); o.Outcome != run.CheckFailed || o.Detail != "cycle 2: fix failed: change set: fixer failed: "+fix.Detail {
		t.Errorf("alpha's cycle 2 outcome = %+v", o)
	}
}

package api

// This file holds skill.check: an invoke run as a role's session against a
// fixture change set in a scratch repo outside every git repository, and
// each gate of that role's contract reported as pass or fail. The current
// repo is never read or written: everything happens in the scratch repo,
// through the running rhiza, and the scratch folder is removed afterwards
// unless kept.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "skill.check",
		Summary:  "Run an invoke as a role's session against a one-leaf fixture change set in a scratch repo outside every git repository, and report each gate of the role's contract as pass or fail; exits 1 on any fail and never touches the current repo",
		Request:  reflect.TypeFor[SkillCheckRequest](),
		Response: reflect.TypeFor[SkillCheckResponse](),
		Run:      runSkillCheck,
	})
}

// SkillCheckRequest names the invoke, the role it plays and its agent.
type SkillCheckRequest struct {
	Invoke         string `json:"invoke" arg:"0" desc:"The text the role's session starts with, as a workflow's invoke: a template over {change}, {leaf}, {attempt} and the other placeholders, such as \"/my-builder Build {leaf}\""`
	Role           string `json:"role,omitempty" desc:"The role the invoke plays: builder; default builder"`
	Agent          string `json:"agent,omitempty" desc:"The agent adapter to launch; default the built-in default adapter"`
	PermissionMode string `json:"permission_mode" desc:"The agent's permission mode, always explicit: one the agent accepts; modes that bypass its permission checks are refused"`
	Keep           bool   `json:"keep,omitempty" desc:"Keep the scratch repo instead of removing it, and name it in the report"`
}

// Skill check's gate names.
const (
	gateSessionDone = "session-done"
	gateTrailers    = "trailers"
	gateCheck       = "check"
)

// SkillGate is one gate of the role's contract and how the session did.
type SkillGate struct {
	Name   string `json:"name" desc:"The gate: session-done, trailers or check"`
	Result string `json:"result" desc:"pass or fail"`
	Detail string `json:"detail" desc:"What the gate wants and what the session did"`
}

// SkillCheckResponse is the gate report.
type SkillCheckResponse struct {
	Invoke  string      `json:"invoke" desc:"The invoke checked"`
	Role    string      `json:"role" desc:"The role it played"`
	Agent   string      `json:"agent" desc:"The agent adapter that ran it"`
	Passed  bool        `json:"passed" desc:"True when every gate passes"`
	Gates   []SkillGate `json:"gates" desc:"Each gate of the role's contract, in order"`
	Run     string      `json:"run" desc:"How the scratch run ended: its result and detail"`
	Scratch string      `json:"scratch,omitempty" desc:"The scratch repo, with --keep; empty once it is removed"`
}

// Human writes the report, a line per gate.
func (r SkillCheckResponse) Human(w io.Writer) {
	pass := 0
	for _, g := range r.Gates {
		if g.Result == workflow.OutcomePass {
			pass++
		}
	}
	verdict := "passes"
	if !r.Passed {
		verdict = "fails"
	}
	fmt.Fprintf(w, "skill check %s as %s %s: %d of %d gates pass\n", r.Invoke, r.Role, verdict, pass, len(r.Gates))
	for _, g := range r.Gates {
		fmt.Fprintf(w, "  %-4s %-12s %s\n", g.Result, g.Name, g.Detail)
	}
	fmt.Fprintf(w, "  run: %s\n", r.Run)
	if r.Scratch != "" {
		fmt.Fprintf(w, "  scratch repo kept: %s\n", r.Scratch)
	}
}

// ExitCode is 1 when a gate fails, else 0.
func (r SkillCheckResponse) ExitCode() int {
	if r.Passed {
		return 0
	}
	return 1
}

// The fixture: the change set, its one leaf and the one-step workflow the
// invoke runs in.
const (
	skillCheckChange   = "Skill check fixture"
	skillCheckLeaf     = "greeting-file"
	skillCheckWorkflow = "skill-check"
	skillCheckStep     = "build"
	// skillCheckVerifier is the fixture workflow's verifier invoke: a
	// leaves step needs one, and it runs only after the builder passes.
	skillCheckVerifier = "Verify leaf {leaf} of change set {change}: follow `rhiza guide verifier`."
)

// skillCheckRoot is the scratch folder's root config: the scratch repo
// is its one project, and no failure classifier runs.
const skillCheckRoot = "[projects.repo]\npath = \"repo\"\n\n[daemon]\n\"rhiza.daemon.failureClassifier\" = \"\"\n"

// skillCheckUnset are the variables the scratch repo's rhiza must not
// inherit from the caller: a session's or run's marks (RHIZA_ROOT would
// send every write to the caller's repo), Rex's block (its status reports
// belong to the caller) and git's repository overrides.
var skillCheckUnset = []string{
	envChange, envActor, envFor, envSession, envRun, envAutonomous, envMilestone, envMilestoneNumber,
	envLeaf, envAttempt, envRole, envCommit, envCheck, envCycle, envRoot,
	envBrief, envSlice, envHandback, "RHIZA_STEP", "RHIZA_RUN_LOG",
	"REX_BLOCK", "GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR",
}

// runSkillCheck refuses an unknown role or an invoke whose skill the agent
// has not installed, then builds the scratch repo, runs the invoke there
// and grades the gates. The scratch folder is removed unless r.Keep.
func runSkillCheck(ctx *Ctx, req any) (any, error) {
	r := req.(SkillCheckRequest)
	r.Invoke = strings.TrimSpace(r.Invoke)
	if r.Invoke == "" {
		return nil, Usagef("skill check needs an invoke: the text the role's session starts with, such as \"/my-builder Build {leaf}\"")
	}
	if r.Role == "" {
		r.Role = workflow.RoleBuilder
	}
	if r.Role != workflow.RoleBuilder {
		return nil, Usagef("unknown role %q: want builder", r.Role)
	}
	a, err := startAdapter(ctx, SessionStartRequest{Agent: r.Agent, PermissionMode: r.PermissionMode})
	if err != nil {
		return nil, err
	}
	// The scratch repo has no skills of its own, so only the agent's
	// config folder counts.
	if err := checkStartSkill(ctx, a, r.Invoke, ""); err != nil {
		return nil, err
	}
	if ctx.Self == "" || ctx.Exec == nil {
		return nil, Errorf("skill check runs the invoke through rhiza: the path of the running rhiza is unknown, or this caller runs no programs")
	}
	if ctx.DryRun {
		return nil, Usagef("skill check has no dry run: it writes only to its own scratch repo")
	}
	dir, err := os.MkdirTemp("", "rhiza-skill-check-*")
	if err != nil {
		return nil, Errorf("scratch folder: %v", err)
	}
	if dir, err = filepath.EvalSymlinks(dir); err != nil {
		os.RemoveAll(dir)
		return nil, Errorf("scratch folder: %v", err)
	}
	sc := &scratch{ctx: ctx, dir: dir, repo: filepath.Join(dir, "repo")}
	out := SkillCheckResponse{Invoke: r.Invoke, Role: r.Role, Agent: a.Name()}
	if r.Keep {
		out.Scratch = sc.repo
	} else {
		defer os.RemoveAll(dir)
	}
	if top, err := sc.git(dir, "rev-parse", "--show-toplevel"); err == nil {
		return nil, Usagef("the scratch folder %s is inside the git repository %s: set TMPDIR to a folder outside every git repository", dir, strings.TrimSpace(top))
	}
	if err := sc.fixture(r.Invoke); err != nil {
		return nil, err
	}
	args := []string{"run", "start", "--foreground", "--workflow", skillCheckWorkflow, "--permission-mode", r.PermissionMode}
	if r.Agent != "" {
		args = append(args, "--agent", r.Agent)
	}
	_, runErr := sc.rhiza(sc.repo, args...)
	gates, ran, err := sc.grade(runErr)
	if err != nil {
		return nil, err
	}
	out.Gates, out.Run = gates, ran
	out.Passed = !slices.ContainsFunc(gates, func(g SkillGate) bool { return g.Result != workflow.OutcomePass })
	return out, nil
}

// scratch is the scratch folder and the repo in it.
type scratch struct {
	ctx  *Ctx
	dir  string
	repo string
}

// rhiza runs the running rhiza with args in dir, its environment the
// caller's without skillCheckUnset, and returns its standard output.
func (s *scratch) rhiza(dir string, args ...string) (string, error) {
	return s.exec(dir, s.ctx.Self, args...)
}

// git runs git with args in dir as rhiza does.
func (s *scratch) git(dir string, args ...string) (string, error) {
	return s.exec(dir, "git", args...)
}

func (s *scratch) exec(dir, name string, args ...string) (string, error) {
	full := make([]string, 0, 2*len(skillCheckUnset)+1+len(args))
	for _, v := range skillCheckUnset {
		full = append(full, "-u", v)
	}
	full = append(append(full, name), args...)
	out, err := s.ctx.Exec("env", full, dir, nil)
	if err != nil {
		// Name the command as run, without env's unsets.
		msg := strings.TrimPrefix(err.Error(), "env "+strings.Join(full, " ")+": ")
		return string(out), fmt.Errorf("%s %s: %s", filepath.Base(name), strings.Join(args, " "), msg)
	}
	return string(out), nil
}

// fixture builds the scratch repo: git and rhiza initialized, a change
// set with one leaf capability, and the one-step workflow whose builder
// runs invoke, all committed.
func (s *scratch) fixture(invoke string) error {
	if err := os.Mkdir(s.repo, 0o755); err != nil {
		return Errorf("scratch repo: %v", err)
	}
	wf, err := skillCheckDefinition(invoke)
	if err != nil {
		return err
	}
	steps := []func() error{
		// Runs need a project registered in a Rhiza root: the scratch
		// folder is its own, so the nearest root is never the caller's.
		func() error {
			return os.WriteFile(filepath.Join(s.dir, rootcfg.File), []byte(skillCheckRoot), 0o644)
		},
		func() error { _, err := s.git(s.repo, "init", "-q"); return err },
		func() error { _, err := s.rhiza(s.repo, "init"); return err },
		func() error { return s.commit("init") },
		func() error {
			_, err := s.rhiza(s.repo, "change", "new", skillCheckChange, "--title", skillCheckChange)
			return err
		},
		func() error {
			_, err := s.rhiza(s.repo, "add", "capability", skillCheckLeaf,
				"--field", "statement=The file GREETING at the repo's top holds the line hello.",
				"--field", "acceptance=test \"$(cat GREETING)\" = hello")
			return err
		},
		func() error {
			p := filepath.Join(s.repo, workflowsDir, skillCheckWorkflow+".toml")
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return err
			}
			return os.WriteFile(p, wf, 0o644)
		},
		func() error { _, err := s.rhiza(s.repo, "render"); return err },
		func() error { return s.commit("fixture") },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return Errorf("build the scratch repo %s: %v", s.repo, err)
		}
	}
	return nil
}

// commit commits everything in the scratch repo.
func (s *scratch) commit(msg string) error {
	if _, err := s.git(s.repo, "add", "-A"); err != nil {
		return err
	}
	_, err := s.git(s.repo, "commit", "-qm", msg)
	return err
}

// skillCheckDefinition is the fixture's one-step workflow: a leaves step
// whose builder runs invoke.
func skillCheckDefinition(invoke string) ([]byte, error) {
	doc := map[string]any{
		"name":  skillCheckWorkflow,
		"start": skillCheckStep,
		"step": map[string]any{
			skillCheckStep: map[string]any{
				"kind":     workflow.KindLeaves,
				"on":       map[string]string{workflow.OutcomePass: "done"},
				"builder":  map[string]string{"invoke": invoke},
				"verifier": map[string]string{"invoke": skillCheckVerifier},
			},
		},
	}
	var b bytes.Buffer
	if err := toml.NewEncoder(&b).Encode(doc); err != nil {
		return nil, Errorf("fixture workflow: %v", err)
	}
	return b.Bytes(), nil
}

// grade reads the scratch run's journal and grades the builder's first
// attempt: its session said done, a commit it made carries the Change:,
// Actor: and Leaf: trailers, and rhiza check is clean at its newest
// commit. runErr is run start's failure, if any: with no run journalled
// it is the error. A run that never started a builder fails every gate.
func (s *scratch) grade(runErr error) ([]SkillGate, string, error) {
	common := filepath.Join(s.repo, ".git")
	js, err := run.NewStore(run.Dir(common)).List()
	if err != nil {
		return nil, "", Errorf("read the scratch run: %v", err)
	}
	if len(js) == 0 {
		why := "no run started"
		if runErr != nil {
			why = runErr.Error()
		}
		return nil, "", Errorf("the scratch run did not start: %s", why)
	}
	st := js[len(js)-1].State()
	ran := st.Result
	if ran == "" {
		ran = "unfinished"
	}
	if st.Detail != "" {
		ran += ": " + st.Detail
	}
	var leaf run.LeafState
	for _, l := range st.Leaves {
		if l.Handle == skillCheckLeaf {
			leaf = l
		}
	}
	b, started := leaf.Session(workflow.RoleBuilder, 1)
	if !started {
		why := "the builder never started (" + ran + ")"
		return []SkillGate{
			{Name: gateSessionDone, Result: workflow.OutcomeFail, Detail: why},
			{Name: gateTrailers, Result: workflow.OutcomeFail, Detail: why},
			{Name: gateCheck, Result: workflow.OutcomeFail, Detail: why},
		}, ran, nil
	}
	commits, why := s.builderCommits(st, leaf, b)
	return []SkillGate{
		s.sessionDoneGate(common, b.Session),
		trailersGate(commits, why, st.Change),
		s.checkGate(commits, why),
	}, ran, nil
}

// builderCommits are the commits the builder's first attempt made, newest
// first: from the commit it started at to the commit its verifier checks,
// else the next attempt's start, else the run's branch head. With none,
// why says so.
func (s *scratch) builderCommits(st run.State, leaf run.LeafState, b run.LeafSession) (commits []check.Commit, why string) {
	if b.Commit == "" || st.Worktree == "" {
		return nil, "the run recorded no commit the builder started from"
	}
	head := "HEAD"
	if v, ok := leaf.Session(workflow.RoleVerifier, 1); ok && v.Commit != "" {
		head = v.Commit
	} else if next, ok := leaf.Session(workflow.RoleBuilder, 2); ok && next.Commit != "" {
		head = next.Commit
	}
	out, err := s.git(st.Worktree, "log", "--format="+trailerLogFormat, b.Commit+".."+head)
	if err != nil {
		return nil, fmt.Sprintf("read the builder's commits: %v", err)
	}
	if commits = parseTrailerLog(out); len(commits) == 0 {
		return nil, "the builder made no commit"
	}
	return commits, ""
}

// sessionDoneGate passes when the builder's session ran rhiza session done
// before it exited.
func (s *scratch) sessionDoneGate(common, sid string) SkillGate {
	g := SkillGate{Name: gateSessionDone}
	rec, err := session.NewStore(session.Dir(common)).Read(sid)
	switch {
	case err != nil:
		g.Result, g.Detail = workflow.OutcomeFail, fmt.Sprintf("read the builder's session %s: %v", sid, err)
	case rec.Status.SaidDone():
		g.Result, g.Detail = workflow.OutcomePass, "the builder ran rhiza session done"
	default:
		g.Result, g.Detail = workflow.OutcomeFail, "the builder never ran rhiza session done before it exited"
	}
	return g
}

// trailersGate passes when one of the builder's commits carries a Change:
// trailer naming the run's change set, an Actor: trailer and a Leaf:
// trailer naming the leaf.
func trailersGate(commits []check.Commit, why, change string) SkillGate {
	g := SkillGate{Name: gateTrailers, Result: workflow.OutcomeFail, Detail: why}
	if len(commits) == 0 {
		return g
	}
	for _, c := range commits {
		if slices.Contains(c.Changes, change) && len(c.Actors) > 0 && slices.Contains(c.Leaves, skillCheckLeaf) {
			g.Result = workflow.OutcomePass
			g.Detail = fmt.Sprintf("commit %s carries Change:, Actor: and Leaf: trailers", shortRev(c.Hash))
			return g
		}
	}
	g.Detail = fmt.Sprintf("none of the builder's %d commit(s) carries all of Change: %s, Actor: and Leaf: %s", len(commits), change, skillCheckLeaf)
	return g
}

// checkGate passes when rhiza check is clean in a detached checkout of
// the builder's newest commit: what it committed, not what it left
// untracked nor what the verifier wrote after it.
func (s *scratch) checkGate(commits []check.Commit, why string) SkillGate {
	g := SkillGate{Name: gateCheck, Result: workflow.OutcomeFail, Detail: why}
	if len(commits) == 0 {
		return g
	}
	at := commits[0].Hash
	dir := filepath.Join(s.dir, "check")
	if _, err := s.git(s.repo, "worktree", "add", "-q", "--detach", dir, at); err != nil {
		g.Detail = fmt.Sprintf("check out the builder's commit %s: %v", shortRev(at), err)
		return g
	}
	if out, err := s.rhiza(dir, "check"); err != nil {
		why := strings.Join(strings.Fields(strings.TrimSpace(out)), " ")
		if why == "" {
			why = err.Error()
		}
		g.Detail = fmt.Sprintf("rhiza check at the builder's commit %s is not clean: %s", shortRev(at), why)
		return g
	}
	g.Result, g.Detail = workflow.OutcomePass, fmt.Sprintf("rhiza check is clean at the builder's commit %s", shortRev(at))
	return g
}

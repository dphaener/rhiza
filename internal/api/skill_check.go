package api

// This file holds skill.check: an invoke run as a role's session against a
// fixture change set in a scratch repo outside every git repository, and
// each gate of that role's contract reported as pass or fail. The current
// repo is never read or written: everything happens in the scratch repo,
// through the running rhiza, and the scratch folder is removed afterwards
// unless kept.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/id"
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
	Role           string `json:"role,omitempty" desc:"The role the invoke plays: builder, verifier or check; default builder"`
	Agent          string `json:"agent,omitempty" desc:"The agent adapter to launch; default the project's rhiza.project.agent, else the built-in default adapter"`
	PermissionMode string `json:"permission_mode" desc:"The agent's permission mode, always explicit: one the agent accepts; modes that bypass its permission checks are refused"`
	Keep           bool   `json:"keep,omitempty" desc:"Keep the scratch repo instead of removing it, and name it in the report"`
}

// Skill check's gate names.
const (
	gateSessionDone = "session-done"
	gateTrailers    = "trailers"
	gateCheck       = "check"
	gateEvidence    = "evidence"
	gateFindings    = "findings"
)

// skillCheckRoles are the roles skill check runs an invoke as.
var skillCheckRoles = []string{workflow.RoleBuilder, workflow.RoleVerifier, workflow.RoleCheck}

// SkillGate is one gate of the role's contract and how the session did.
type SkillGate struct {
	Name   string `json:"name" desc:"The gate: session-done, then the role's own (a builder's trailers, a verifier's evidence, a check's findings), then check"`
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
	// skillCheckChecks and skillCheckCheck are the check fixture's checks
	// step and its one check.
	skillCheckChecks = "review"
	skillCheckCheck  = "fixture"
	// skillCheckBuilder is the actor of the fixture's builder commit, the
	// one the verifier and check roles start from, and skillCheckBuilt
	// the session the verifier fixture's journal says made it.
	skillCheckBuilder = "fixture-builder"
	skillCheckBuilt   = "skill-check-fixture-builder"
	// skillCheckUnused starts the invoke of every fixture role but the
	// one checked: a slash command no agent has installed, so its session
	// start is refused before an agent launches and the run ends there,
	// rather than a fixture role spending a real session.
	skillCheckUnused = "/rhiza-skill-check-unused The fixture's "
	// skillCheckArbiter is the fixture's arbiter step and skillCheckNoMode
	// its permission mode, one no adapter accepts: the question a failed
	// leaf or review raises gets no advisor session, whose start is
	// refused and only logged. The run is never autonomous, so no arbiter
	// runs.
	skillCheckArbiter = "arbiter"
	skillCheckNoMode  = "skill-check-none"
)

// envSkillCheckOnce, set in the builder and check fixtures' scratch rhiza,
// refuses a leaf builder's start past its round's first attempt and a
// check's past its first cycle: a failed builder or check ends the
// scratch run instead of launching the agent again, as only the first
// attempt or cycle is graded.
const envSkillCheckOnce = "RHIZA_SKILL_CHECK_ONCE"

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
	envBrief, envSlice, envHandback, envSkillCheckOnce, "RHIZA_STEP", "RHIZA_RUN_LOG",
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
	if !slices.Contains(skillCheckRoles, r.Role) {
		return nil, Usagef("unknown role %q: want builder, verifier or check", r.Role)
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
	tmp, err := os.MkdirTemp(tempDir(ctx), "rhiza-skill-check-*")
	if err != nil {
		return nil, Errorf("scratch folder: %v", err)
	}
	dir, err := filepath.EvalSymlinks(tmp)
	if err != nil {
		os.RemoveAll(tmp)
		return nil, Errorf("scratch folder: %v", err)
	}
	sc := &scratch{ctx: ctx, dir: dir, repo: filepath.Join(dir, "repo"), role: r.Role}
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
	// The scratch repo has no project setting of its own, so it is told the
	// adapter this project resolved.
	args = append(args, "--agent", a.Name())
	_, runErr := sc.rhiza(sc.repo, args...)
	if r.Role == workflow.RoleVerifier {
		// The fixture's builder is refused at its start, so the run fails
		// there; journal its attempt as passed at the fixture's builder
		// commit and retry the step, which starts the verifier on it.
		if err := sc.seedBuilder(runErr); err != nil {
			return nil, Errorf("the scratch run: %v", err)
		}
		_, runErr = sc.rhiza(sc.repo, "run", "resume", "--retry", "--foreground", sc.run)
	}
	gates, ran, err := sc.grade(runErr)
	if err != nil {
		return nil, err
	}
	out.Gates, out.Run = gates, ran
	out.Passed = !slices.ContainsFunc(gates, func(g SkillGate) bool { return g.Result != workflow.OutcomePass })
	return out, nil
}

// scratch is the scratch folder and the repo in it, the role checked,
// and what the fixture made: its change set, its leaf, the commit the
// builder started from and, for the verifier and check roles, the
// fixture's builder commit; run is the scratch run once seeded.
type scratch struct {
	ctx    *Ctx
	dir    string
	repo   string
	role   string
	change string
	leaf   string
	base   string
	built  string
	run    string
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
	full := make([]string, 0, 2*len(skillCheckUnset)+2+len(args))
	for _, v := range skillCheckUnset {
		full = append(full, "-u", v)
	}
	if s.role == workflow.RoleBuilder || s.role == workflow.RoleCheck {
		full = append(full, envSkillCheckOnce+"=1")
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
// set with one leaf capability, and the one-step workflow whose role
// s.role runs invoke, all committed. For the verifier and check roles the
// leaf is then built: a commit writing the file, its Change:, Actor: and
// Leaf: trailers naming the change set, the fixture's recorded builder
// actor and the leaf.
func (s *scratch) fixture(invoke string) error {
	if err := os.Mkdir(s.repo, 0o755); err != nil {
		return Errorf("scratch repo: %v", err)
	}
	wf, err := skillCheckDefinition(s.role, invoke)
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
		func() error {
			if s.role == workflow.RoleBuilder {
				return nil
			}
			if _, err := s.rhiza(s.repo, "add", "actor", skillCheckBuilder, "--record"); err != nil {
				return err
			}
			_, err := s.rhiza(s.repo, "render")
			return err
		},
		func() error { return s.commit("init") },
		func() error {
			out, err := s.rhiza(s.repo, "change", "new", skillCheckChange, "--title", skillCheckChange, "--json")
			if err != nil {
				return err
			}
			var cs struct {
				ChangeSet string `json:"change_set"`
			}
			if err := json.Unmarshal([]byte(out), &cs); err != nil || cs.ChangeSet == "" {
				return fmt.Errorf("rhiza change new printed no change set: %s", out)
			}
			s.change = cs.ChangeSet
			return nil
		},
		func() error {
			out, err := s.rhiza(s.repo, "add", "capability", skillCheckLeaf,
				"--field", "statement=The file GREETING at the repo's top holds the line hello.",
				"--field", "acceptance=test \"$(cat GREETING)\" = hello", "--json")
			if err != nil {
				return err
			}
			var n struct {
				Nodes []string `json:"nodes"`
			}
			if err := json.Unmarshal([]byte(out), &n); err != nil || len(n.Nodes) != 1 {
				return fmt.Errorf("rhiza add capability printed no node: %s", out)
			}
			s.leaf = n.Nodes[0]
			return nil
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
		func() (err error) { s.base, err = s.head(); return err },
		func() error {
			if s.role == workflow.RoleBuilder {
				return nil
			}
			if err := os.WriteFile(filepath.Join(s.repo, "GREETING"), []byte("hello\n"), 0o644); err != nil {
				return err
			}
			if err := s.commit("build " + skillCheckLeaf + "\n\nChange: " + s.change + "\nActor: " + skillCheckBuilder + "\nLeaf: " + skillCheckLeaf); err != nil {
				return err
			}
			s.built, err = s.head()
			return err
		},
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return Errorf("build the scratch repo %s: %v", s.repo, err)
		}
	}
	return nil
}

// head is the scratch repo's HEAD commit.
func (s *scratch) head() (string, error) {
	out, err := s.git(s.repo, "rev-parse", "HEAD")
	return strings.TrimSpace(out), err
}

// seedBuilder journals, on the verifier fixture's failed run, its
// builder's first attempt: a session started at the fixture commit and
// judged pass at the fixture's builder commit, as a builder that built
// the leaf leaves it. runErr is run start's failure: the fixture's
// builder is refused at its start, so the run must have failed at its
// one step, and anything else is that failure.
func (s *scratch) seedBuilder(runErr error) error {
	rs := run.NewStore(run.Dir(filepath.Join(s.repo, ".git")))
	js, err := rs.List()
	if err != nil {
		return fmt.Errorf("read the scratch run: %v", err)
	}
	if len(js) == 0 {
		if runErr == nil {
			runErr = fmt.Errorf("no run started")
		}
		return runErr
	}
	j := js[len(js)-1]
	st := j.State()
	if step, err := st.FailedStep(); err != nil || step != skillCheckStep || st.Visits[skillCheckStep] != 1 {
		return fmt.Errorf("the run did not fail at its fixture builder's start: %s: %s", st.Result, st.Detail)
	}
	s.run = j.ID
	now := s.ctx.Now()
	for _, e := range []run.Entry{
		{Type: run.TypeLeafSession, At: now, Step: skillCheckStep, Leaf: s.leaf, Handle: skillCheckLeaf,
			Role: workflow.RoleBuilder, Attempt: 1, Session: skillCheckBuilt, Commit: s.base},
		{Type: run.TypeLeafAttempt, At: now, Step: skillCheckStep, Leaf: s.leaf, Handle: skillCheckLeaf,
			Role: workflow.RoleBuilder, Attempt: 1, Outcome: workflow.OutcomePass, Commit: s.built,
			Detail: "the fixture's builder commit"},
	} {
		if err := rs.Append(j.ID, e); err != nil {
			return fmt.Errorf("journal the fixture's builder: %v", err)
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

// skillCheckDefinition is the fixture's one-step workflow for role: a
// leaves step whose builder runs invoke, its verifier never started (the
// scratch rhiza runs the builder once, envSkillCheckOnce); a leaves step whose verifier
// runs invoke, its builder never started; or a checks step with one
// check that runs invoke, its critic, fixer and verifier never started
// (the scratch rhiza runs the check in its first cycle only,
// envSkillCheckOnce). Its arbiter step launches no advisor.
func skillCheckDefinition(role, invoke string) ([]byte, error) {
	unused := func(r string) map[string]string {
		return map[string]string{"invoke": skillCheckUnused + r + " never runs."}
	}
	step, body := skillCheckStep, map[string]any{
		"kind":     workflow.KindLeaves,
		"on":       map[string]string{workflow.OutcomePass: "done"},
		"builder":  map[string]string{"invoke": invoke},
		"verifier": unused(workflow.RoleVerifier),
	}
	switch role {
	case workflow.RoleVerifier:
		body["builder"], body["verifier"] = unused(workflow.RoleBuilder), map[string]string{"invoke": invoke}
	case workflow.RoleCheck:
		step, body = skillCheckChecks, map[string]any{
			"kind":                workflow.KindChecks,
			"on":                  map[string]string{workflow.OutcomePass: "done"},
			"check":               []map[string]string{{"name": skillCheckCheck, "invoke": invoke}},
			workflow.RoleCritic:   unused(workflow.RoleCritic),
			workflow.RoleFixer:    unused(workflow.RoleFixer),
			workflow.RoleVerifier: unused(workflow.RoleVerifier),
		}
	}
	doc := map[string]any{
		"name":  skillCheckWorkflow,
		"start": step,
		"step": map[string]any{
			step:              body,
			skillCheckArbiter: map[string]string{"kind": workflow.KindArbiter, "permission_mode": skillCheckNoMode},
		},
	}
	var b bytes.Buffer
	if err := toml.NewEncoder(&b).Encode(doc); err != nil {
		return nil, Errorf("fixture workflow: %v", err)
	}
	return b.Bytes(), nil
}

// grade reads the scratch run's journal and grades the first session of
// the role checked (gradeBuilder, gradeVerifier, gradeCheck). runErr is
// the scratch run's failure, if any: with no run journalled it is the
// error.
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
	ran := runOutcome(st, s.role)
	var leaf run.LeafState
	for _, l := range st.Leaves {
		if l.Handle == skillCheckLeaf {
			leaf = l
		}
	}
	switch s.role {
	case workflow.RoleVerifier:
		return s.gradeVerifier(common, st, leaf, ran), ran, nil
	case workflow.RoleCheck:
		return s.gradeCheck(common, st, ran), ran, nil
	}
	return s.gradeBuilder(common, st, leaf, ran), ran, nil
}

// runOutcome is how the scratch run ended: its result and detail, except
// that the fixture's own stop at a role it never runs (skillCheckUnused
// refused at its start) is reported as skill check's stop, not a failure.
func runOutcome(st run.State, role string) string {
	if st.Result == run.Failed && strings.Contains(st.Detail, "start command "+strings.Fields(skillCheckUnused)[0]+" ") {
		return "stopped by skill check after the " + role + ": the fixture starts no other role's session"
	}
	ran := st.Result
	if ran == "" {
		ran = "unfinished"
	}
	if st.Detail != "" {
		ran += ": " + st.Detail
	}
	return ran
}

// notStarted fails each of gates: the role's session never started.
func notStarted(role, ran string, gates ...string) []SkillGate {
	why := "the " + role + " never started (" + ran + ")"
	out := make([]SkillGate, len(gates))
	for i, g := range gates {
		out[i] = SkillGate{Name: g, Result: workflow.OutcomeFail, Detail: why}
	}
	return out
}

// gradeBuilder grades the builder's first attempt: its session said done,
// a commit it made carries the Change:, Actor: and Leaf: trailers, and
// rhiza check is clean at its newest commit. A run that never started a
// builder fails every gate.
func (s *scratch) gradeBuilder(common string, st run.State, leaf run.LeafState, ran string) []SkillGate {
	b, started := leaf.Session(workflow.RoleBuilder, 1)
	if !started {
		return notStarted(workflow.RoleBuilder, ran, gateSessionDone, gateTrailers, gateCheck)
	}
	commits, why := s.builderCommits(st, leaf, b)
	return []SkillGate{
		s.sessionDoneGate(common, workflow.RoleBuilder, b.Session),
		trailersGate(commits, why, st.Change),
		s.checkGate(commits, why),
	}
}

// gradeVerifier grades the verifier's session on the fixture's builder
// commit: it said done, live Evidence its actor wrote verifies the leaf,
// and rhiza check is clean in the run's worktree once it ended, its
// Evidence committed there. A run that never started the verifier fails
// every gate.
func (s *scratch) gradeVerifier(common string, st run.State, leaf run.LeafState, ran string) []SkillGate {
	v, started := leaf.Session(workflow.RoleVerifier, 1)
	if !started {
		return notStarted(workflow.RoleVerifier, ran, gateSessionDone, gateEvidence, gateCheck)
	}
	return []SkillGate{
		s.sessionDoneGate(common, workflow.RoleVerifier, v.Session),
		s.evidenceGate(common, st, v.Session),
		s.worktreeCheckGate(st, workflow.RoleVerifier),
	}
}

// gradeCheck grades the check's first cycle: its session said done, it
// recorded a round with rhiza run findings for its check, and rhiza check
// is clean in the run's worktree once it ended, its writes committed
// there. A run that never started the check fails every gate.
func (s *scratch) gradeCheck(common string, st run.State, ran string) []SkillGate {
	c, _ := st.Check(skillCheckChecks, skillCheckCheck)
	cs, started := c.Session(workflow.RoleCheck, 1)
	if !started {
		return notStarted(workflow.RoleCheck, ran, gateSessionDone, gateFindings, gateCheck)
	}
	return []SkillGate{
		s.sessionDoneGate(common, workflow.RoleCheck, cs.Session),
		findingsGate(c),
		s.worktreeCheckGate(st, workflow.RoleCheck),
	}
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

// sessionDoneGate passes when role's session sid ran rhiza session done
// before it exited.
func (s *scratch) sessionDoneGate(common, role, sid string) SkillGate {
	g := SkillGate{Name: gateSessionDone}
	rec, err := session.NewStore(session.Dir(common)).Read(sid)
	switch {
	case err != nil:
		g.Result, g.Detail = workflow.OutcomeFail, fmt.Sprintf("read the %s's session %s: %v", role, sid, err)
	case rec.Status.SaidDone():
		g.Result, g.Detail = workflow.OutcomePass, "the "+role+" ran rhiza session done"
	default:
		g.Result, g.Detail = workflow.OutcomeFail, "the "+role+" never ran rhiza session done before it exited"
	}
	return g
}

// evidenceGate passes when live Evidence created by the actor of verifier
// session sid verifies the fixture's leaf, as the leaves step judges a
// verifier: read in the run's worktree, where the verifier's writes land.
func (s *scratch) evidenceGate(common string, st run.State, sid string) SkillGate {
	g := SkillGate{Name: gateEvidence, Result: workflow.OutcomeFail}
	ok, detail, err := s.verifierEvidence(common, st, sid)
	switch {
	case err != nil:
		g.Detail = fmt.Sprintf("read the verifier's Evidence: %v", err)
	case ok:
		g.Result, g.Detail = workflow.OutcomePass, "Evidence the verifier wrote verifies "+skillCheckLeaf
	case detail != "":
		g.Detail = "no Evidence the verifier wrote verifies " + skillCheckLeaf + ": " + detail
	default:
		g.Detail = "no Evidence verifies " + skillCheckLeaf + ": the verifier wrote none"
	}
	return g
}

// verifierEvidence is leafEvidence over the run change set's view in the
// run's worktree, for the actor of verifier session sid.
func (s *scratch) verifierEvidence(common string, st run.State, sid string) (bool, string, error) {
	rec, err := session.NewStore(session.Dir(common)).Read(sid)
	if err != nil {
		return false, "", err
	}
	cs, err := id.Parse(st.Change)
	if err != nil {
		return false, "", err
	}
	lid, err := id.Parse(s.leaf)
	if err != nil {
		return false, "", err
	}
	w, err := openRootWorkspace(s.ctx, st.Worktree, true)
	if err != nil {
		return false, "", err
	}
	defer w.closeProjection()
	m := w.model(cs)
	actor, err := liveActor(m, rec.Actor)
	if err != nil {
		return false, "", fmt.Errorf("its actor %s: %s", rec.Actor, AsError(err).Message)
	}
	ok, detail := leafEvidence(m, w.creator(m), lid, actor)
	return ok, detail, nil
}

// findingsGate passes when check c recorded a findings round in cycle 1,
// with rhiza run findings --check.
func findingsGate(c run.CheckState) SkillGate {
	g := SkillGate{Name: gateFindings, Result: workflow.OutcomeFail,
		Detail: fmt.Sprintf("the check recorded no round with rhiza run findings --check %s in cycle 1", skillCheckCheck)}
	if r, ok := c.RoundFor(1); ok {
		g.Result = workflow.OutcomePass
		g.Detail = fmt.Sprintf("the check recorded round %d with %d finding(s) for check %s", r.Number, len(r.Findings), skillCheckCheck)
	}
	return g
}

// worktreeCheckGate passes when rhiza check is clean in the run's
// worktree once role's session ended: the fixture's builder commit and
// what the role wrote there. A verifier or check never renders, so the
// worktree is rendered first, as a review's render step does after them.
func (s *scratch) worktreeCheckGate(st run.State, role string) SkillGate {
	g := SkillGate{Name: gateCheck, Result: workflow.OutcomeFail}
	if st.Worktree == "" {
		g.Detail = "the run recorded no worktree"
		return g
	}
	if _, err := s.rhiza(st.Worktree, "render"); err != nil {
		g.Detail = fmt.Sprintf("rhiza render in the run's worktree after the %s: %v", role, err)
		return g
	}
	if out, err := s.rhiza(st.Worktree, "check"); err != nil {
		why := strings.Join(strings.Fields(strings.TrimSpace(out)), " ")
		if why == "" {
			why = err.Error()
		}
		g.Detail = fmt.Sprintf("rhiza check in the run's worktree after the %s is not clean: %s", role, why)
		return g
	}
	g.Result, g.Detail = workflow.OutcomePass, "rhiza check is clean in the run's worktree after the "+role
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

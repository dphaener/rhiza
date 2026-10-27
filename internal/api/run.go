package api

// This file holds the run operations: run.start runs a workflow's steps on
// a change set in the foreground, run.status reports a run from its
// journal, and run.resume continues a stopped or interrupted run, or with
// --retry restarts a failed one at the step that failed. A run's
// journal lives in <git common dir>/rhiza/runs, beside the session
// records, and is never committed: it is work-layer state, not the model.
//
// The engine (internal/run) decides what runs next; the stepRunner here
// runs it: command and wait steps through `sh -c` in the run's worktree,
// agent steps through session.start's path, so nothing here names an
// agent.

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "run.start",
		Summary:  "Run a workflow's steps on a change set in its worktree, in the foreground",
		Write:    true,
		Request:  reflect.TypeFor[RunStartRequest](),
		Response: reflect.TypeFor[RunResponse](),
		Run:      runRunStart,
	})
	register(Operation{
		Name:     "run.status",
		Summary:  "Show a run's state and the steps it has done, from its journal",
		Request:  reflect.TypeFor[RunStatusRequest](),
		Response: reflect.TypeFor[RunInfo](),
		Run:      runRunStatus,
	})
	register(Operation{
		Name:     "run.resume",
		Summary:  "Continue a stopped or interrupted run where its journal left it, or retry a failed run at its failed step",
		Write:    true,
		Request:  reflect.TypeFor[RunResumeRequest](),
		Response: reflect.TypeFor[RunResponse](),
		Run:      runRunResume,
	})
}

// Next commands of the run operations.
const (
	runStartNext  = "rhiza describe run.start"
	runStatusNext = "rhiza run status"
)

// Run states, as RunInfo.State reports them.
const (
	runRunning = "running"
	runStopped = "stopped"
	runBlocked = "blocked"
)

// runClock is the clock run.start and run.resume poll by; tests replace
// it to drive a session's status between polls.
var runClock run.Clock = run.RealClock{}

// workflowsDir is where named workflows live, under the project root.
const workflowsDir = store.Dir + "/workflows"

// maxRunIDDraws bounds the IDs run.start draws looking for one no run has.
const maxRunIDDraws = 100

// RunStartRequest starts a run of a workflow on a change set.
type RunStartRequest struct {
	Change         string `json:"change,omitempty" arg:"0" desc:"The change set: a full ID or a unique prefix of 4+ characters; default --change, else the current one"`
	Workflow       string `json:"workflow,omitempty" desc:"The workflow: a .toml file (a path with a / or ending in .toml), else a name: .rhiza/workflows/<name>.toml, else a built-in; default ship-change"`
	PermissionMode string `json:"permission_mode,omitempty" desc:"The permission mode of agent steps that set none; modes that bypass the agent's permission checks are refused"`
	Agent          string `json:"agent,omitempty" desc:"The agent adapter of agent steps that name none; default the built-in default adapter"`
	NoWait         bool   `json:"no_wait,omitempty" desc:"Stop at a wait step whose condition does not hold yet, instead of waiting; resume the run later"`
	Autonomous     bool   `json:"autonomous,omitempty" desc:"Run without the owner: skip the owner's gates and questions, recording each; never changes agent permissions. Resume keeps it"`
}

// RunStatusRequest names a run, or none for the latest.
type RunStatusRequest struct {
	Run string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
}

// RunResumeRequest names the run to continue.
type RunResumeRequest struct {
	Run    string `json:"run" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters"`
	NoWait bool   `json:"no_wait,omitempty" desc:"Stop at a wait step whose condition does not hold yet, instead of waiting"`
	Retry  bool   `json:"retry,omitempty" flag:"retry" desc:"Restart a failed run at the step that failed, once its cause is fixed; the retry is journalled. A run that succeeded is refused"`
}

// RunResponse is a run after run.start or run.resume ran it, or in a dry
// run the one that would run.
type RunResponse struct {
	Run    RunInfo `json:"run" desc:"The run as its journal stands; the command exits 1 when the run failed"`
	DryRun bool    `json:"dry_run,omitempty" desc:"True if nothing was written or run"`
	Retry  string  `json:"retry,omitempty" desc:"The step a resume --retry restarted the failed run at, or in a dry run would; absent otherwise"`
}

// RunInfo is a run folded from its journal.
type RunInfo struct {
	ID         string         `json:"id" desc:"The run ID"`
	Workflow   string         `json:"workflow" desc:"The workflow's name"`
	Change     string         `json:"change" desc:"The change set it runs on"`
	Title      string         `json:"title" desc:"The change set's title"`
	Worktree   string         `json:"worktree" desc:"The worktree every step runs in"`
	Autonomous bool           `json:"autonomous,omitempty" desc:"True for a run started with --autonomous: it goes on without the owner; absent otherwise"`
	State      string         `json:"state" desc:"running (another process runs it now), blocked (its agent step waits on the owner), stopped (resumable), succeeded or failed; empty in a dry run"`
	Current    string         `json:"current,omitempty" desc:"The step being run or to run next; empty once the run is finished"`
	Detail     string         `json:"detail,omitempty" desc:"Why it stopped or failed"`
	Steps      []RunStep      `json:"steps" desc:"The steps done, in order"`
	StartedAt  string         `json:"started_at" desc:"When it started, RFC 3339 UTC"`
	Blocked    *RunBlocked    `json:"blocked,omitempty" desc:"What the run waits on while its agent step's session is blocked; absent otherwise"`
	Amendments []RunAmendment `json:"amendments" desc:"The amendments made to its agent steps' briefs or plans, in order"`
	Retries    []RunRetry     `json:"retries,omitempty" desc:"The retries of the run after it failed, in order; absent when it was never retried"`
	Accepts    []RunAccept    `json:"accepts,omitempty" desc:"The change sets an autonomous run accepted or rejected in the owner's place under rule autonomous, in order; absent when there were none"`
	Questions  []RunQuestion  `json:"questions,omitempty" desc:"The questions an autonomous run's agent steps asked, each with its arbiter's verdict, in order; absent when there were none"`
}

// RunAccept is one accepted or rejected entry journalled in an autonomous
// run: it moved a change set past the owner's gate under rule autonomous.
type RunAccept struct {
	Change string `json:"change" desc:"The change set accepted or rejected"`
	Entry  string `json:"entry" desc:"The model log entry ID of the accept or reject"`
	Verb   string `json:"verb" desc:"accepted or rejected"`
	At     string `json:"at" desc:"When it was journalled, RFC 3339 UTC"`
}

// RunQuestion is one question an autonomous run's agent step asked, which
// the engine put to the arbiter as a question node.
type RunQuestion struct {
	Question string `json:"question" desc:"The question node's ID"`
	Handle   string `json:"handle" desc:"The question node's handle"`
	Text     string `json:"text" desc:"What the session asked"`
	Step     string `json:"step" desc:"The agent step that asked"`
	Outcome  string `json:"outcome,omitempty" desc:"The arbiter's verdict: pass, fail or error; empty while it is open"`
	Decision string `json:"decision,omitempty" desc:"The handle of the Decision that answers it, once the arbiter passed; empty otherwise"`
	Answered bool   `json:"answered,omitempty" desc:"True once the answer was delivered to the asking session"`
	At       string `json:"at" desc:"When it was asked, RFC 3339 UTC"`
}

// RunRetry is one retry journalled in a run: the run had failed and
// resume --retry restarted it at a step.
type RunRetry struct {
	Step   string `json:"step" desc:"The step the run was restarted at"`
	Detail string `json:"detail,omitempty" desc:"Why the run had failed"`
	At     string `json:"at" desc:"When it was retried, RFC 3339 UTC"`
}

// RunAmendment is one amendment journaled in a run: a step's brief or
// plan file replaced while the step ran.
type RunAmendment struct {
	Step    string `json:"step" desc:"The amended step"`
	Session string `json:"session" desc:"The step's session, which was sent a pointer to the file"`
	Actor   string `json:"actor" desc:"The sender's actor handle"`
	File    string `json:"file" desc:"The amended file, relative to the run's worktree"`
	At      string `json:"at" desc:"When it was amended, RFC 3339 UTC"`
}

// RunBlocked is what a blocked run waits on, read from its agent step's
// session record when the run is reported.
type RunBlocked struct {
	Step     string     `json:"step" desc:"The agent step that waits"`
	Session  string     `json:"session" desc:"The step's session"`
	Kind     string     `json:"kind,omitempty" desc:"What kind of block: permission, question or auth; empty when the session gave none"`
	Question string     `json:"question,omitempty" desc:"What the session asks the owner, as it reported it; empty when it gave nothing"`
	Tab      SessionTab `json:"tab" desc:"The terminal tab the session runs in, where the owner answers"`
}

// RunStep is one step a run has done.
type RunStep struct {
	Step    string `json:"step" desc:"The step's name"`
	Visit   int    `json:"visit" desc:"Which visit of the step it was, from 1"`
	Outcome string `json:"outcome" desc:"How it ended: pass, fail, timeout, or error when it could not run"`
	Next    string `json:"next,omitempty" desc:"The step the outcome led to; empty when it ended the run"`
	Detail  string `json:"detail,omitempty" desc:"What the step reported, such as its exit code"`
	Session string `json:"session,omitempty" desc:"An agent step's session"`
	At      string `json:"at" desc:"When it ended, RFC 3339 UTC"`
	Retry   bool   `json:"retry,omitempty" desc:"True for the first visit of the step after a retry restarted the run at it"`
}

// ExitCode is 1 when the run failed, else 0.
func (r RunResponse) ExitCode() int {
	if r.Run.State == run.Failed {
		return 1
	}
	return 0
}

// runRunStart resolves and checks everything (the workflow, its agent
// steps' adapters and modes, the change set and its worktree) before
// writing anything; then creates or reuses the worktree, starts the run's
// journal and runs it in the foreground. A dry run writes and runs nothing.
func runRunStart(ctx *Ctx, req any) (any, error) {
	r := req.(RunStartRequest)
	path, def, data, err := resolveWorkflow(ctx.Root, r.Workflow)
	if err != nil {
		return nil, err
	}
	if err := checkAgentSteps(def, r.PermissionMode, r.Agent, r.Autonomous, runStartNext); err != nil {
		return nil, err
	}
	t, err := resolveChangeWorktree(ctx, r.Change, "no run can work on it")
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(t.common))
	rid, err := newRunID(ctx, rs)
	if err != nil {
		return nil, err
	}
	started := run.Started{
		Workflow: def.Name, WorkflowPath: path, Definition: string(data),
		Change: t.cs.String(), Title: t.title, Worktree: t.worktree, Branch: t.branch,
		PermissionMode: r.PermissionMode, Agent: r.Agent, Autonomous: r.Autonomous,
	}
	if ctx.DryRun {
		return RunResponse{Run: RunInfo{
			ID: rid, Workflow: def.Name, Change: started.Change, Title: t.title, Worktree: t.worktree,
			Autonomous: r.Autonomous, Current: def.Start, Steps: []RunStep{}, StartedAt: rfc3339(ctx.Now()), Amendments: []RunAmendment{},
		}, DryRun: true}, nil
	}
	if err := canRun(ctx, def); err != nil {
		return nil, err
	}
	if _, err := t.open(ctx); err != nil {
		return nil, err
	}
	if err := run.Start(rs, rid, started, ctx.Now(), def); err != nil {
		return nil, Errorf("start run %s: %v", rid, err)
	}
	return runEngine(ctx, rs, t.common, rid, started.Change, t.worktree, r.PermissionMode, r.Agent, started.Autonomous, r.NoWait)
}

// runRunResume continues a run that is neither finished nor running in
// another process, with the definition its journal started with. With
// Retry it also restarts a failed run at its failed step: the engine
// appends the retried entry under the run's lock, then runs it. Without
// Retry a finished run is refused, as is a succeeded one with it.
func runRunResume(ctx *Ctx, req any) (any, error) {
	r := req.(RunResumeRequest)
	ref := strings.TrimSpace(r.Run)
	if ref == "" {
		return nil, Usagef("run resume needs a run ID").WithNext(runStatusNext)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	j, err := readRun(rs, ref)
	if err != nil {
		return nil, err
	}
	st := j.State()
	retryStep := ""
	switch {
	case st.Ended && !r.Retry:
		next := "rhiza run status " + j.ID
		if st.Result == run.Failed {
			next = "rhiza run resume --retry " + j.ID
		}
		return nil, Conflictf("run %s is finished: it %s", j.ID, st.Result).WithNext(next)
	case st.Ended:
		step, err := st.FailedStep()
		if err != nil {
			return nil, Conflictf("run %s %s: only a failed run can be retried", j.ID, st.Result).WithNext("rhiza run status " + j.ID)
		}
		retryStep = step
	}
	if rs.Live(j.ID) {
		return nil, Conflictf("run %s is running in another process", j.ID).WithNext("rhiza run status " + j.ID)
	}
	def, err := workflow.Parse(st.Workflow, []byte(st.Definition))
	if err != nil {
		return nil, Errorf("run %s: its workflow no longer parses: %v", j.ID, err)
	}
	// Start checked the agent steps; check again, since the journal is a
	// file anyone can edit.
	if err := checkAgentSteps(def, st.PermissionMode, st.Agent, st.Autonomous, "rhiza run status "+j.ID); err != nil {
		return nil, err
	}
	if retryStep != "" {
		if _, ok := def.Steps[retryStep]; !ok {
			return nil, Errorf("run %s: step %q is not in workflow %q", j.ID, retryStep, def.Name).WithNext("rhiza run status " + j.ID)
		}
	}
	if ctx.DryRun {
		if retryStep != "" {
			// Fold the retried entry the engine would append, without
			// writing it, to show the run as the retry would leave it.
			j.Entries = append(j.Entries, run.Entry{Type: run.TypeRetried, At: ctx.Now(), Step: retryStep, Detail: st.Detail})
			st = j.State()
		}
		return RunResponse{Run: runInfo(st, false, session.NewStore(session.Dir(common))), DryRun: true, Retry: retryStep}, nil
	}
	if err := canRun(ctx, def); err != nil {
		return nil, err
	}
	resp, err := runEngineRetry(ctx, rs, common, j.ID, st.Change, st.Worktree, st.PermissionMode, st.Agent, st.Autonomous, r.NoWait, r.Retry)
	if err != nil {
		return nil, err
	}
	resp.Retry = retryStep
	return resp, nil
}

// runRunStatus reports a run from its journal alone.
func runRunStatus(ctx *Ctx, req any) (any, error) {
	r := req.(RunStatusRequest)
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	var j run.Journal
	if ref := strings.TrimSpace(r.Run); ref != "" {
		if j, err = readRun(rs, ref); err != nil {
			return nil, err
		}
	} else if j, err = latestRun(ctx, rs); err != nil {
		return nil, err
	}
	return runInfo(j.State(), rs.Live(j.ID), session.NewStore(session.Dir(common))), nil
}

// latestRun is the latest started run of the effective change set, else
// the latest run.
func latestRun(ctx *Ctx, rs *run.Store) (run.Journal, error) {
	runs, err := rs.List()
	if err != nil {
		return run.Journal{}, Errorf("read runs: %v", err)
	}
	if len(runs) == 0 {
		return run.Journal{}, NotFoundf("no runs yet").WithNext("rhiza run start --workflow <name>")
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return run.Journal{}, err
	}
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return run.Journal{}, err
	}
	if !cs.IsZero() {
		for i := len(runs) - 1; i >= 0; i-- {
			if runs[i].State().Change == cs.String() {
				return runs[i], nil
			}
		}
	}
	return runs[len(runs)-1], nil
}

// readRun reads the run ref names, a full ID or a unique prefix, mapping
// the store's errors to not_found and usage.
func readRun(rs *run.Store, ref string) (run.Journal, error) {
	j, err := rs.Read(ref)
	var nf *run.NotFoundError
	var amb *run.AmbiguousError
	switch {
	case errors.As(err, &nf):
		return j, NotFoundf("no run matches %q", ref).WithNext(runStatusNext)
	case errors.As(err, &amb):
		return j, Usagef("%v", amb).WithNext(runStatusNext)
	case err != nil:
		return j, Errorf("read run %s: %v", ref, err)
	}
	return j, nil
}

// defaultWorkflow is the workflow a run starts when none is named.
const defaultWorkflow = "ship-change"

// builtinPrefix starts the WorkflowPath of a run of a built-in workflow.
const builtinPrefix = "builtin:"

// workflowFile is the file ref names: a path when it has a separator or
// ends in .toml (made absolute against the working directory, as apply
// reads its file), else <root>/.rhiza/workflows/<ref>.toml. An empty ref
// names the default workflow.
func workflowFile(root, ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		ref = defaultWorkflow
	}
	if strings.ContainsRune(ref, '/') || strings.ContainsRune(ref, filepath.Separator) || strings.HasSuffix(ref, ".toml") {
		abs, err := filepath.Abs(ref)
		if err != nil {
			return "", Errorf("workflow %s: %v", ref, err)
		}
		return abs, nil
	}
	if ref == "." || ref == ".." {
		return "", Usagef("workflow %q is not a workflow name", ref).WithNext(runStartNext)
	}
	return filepath.Join(root, filepath.FromSlash(workflowsDir), ref+".toml"), nil
}

// isWorkflowPath reports whether ref names a file rather than a workflow.
func isWorkflowPath(ref string) bool {
	return strings.ContainsRune(ref, '/') || strings.ContainsRune(ref, filepath.Separator) || strings.HasSuffix(ref, ".toml")
}

// resolveWorkflow finds the workflow ref names and returns the path to
// record in the journal (the file, or builtin:<name> for an embedded
// one), its definition and its source. A bare name resolves to the repo's
// .rhiza/workflows/<name>.toml when that exists, so a repo replaces a
// built-in by name, else to the built-in of that name.
func resolveWorkflow(root, ref string) (string, *workflow.Definition, []byte, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		ref = defaultWorkflow
	}
	path, err := workflowFile(root, ref)
	if err != nil {
		return "", nil, nil, err
	}
	if !isWorkflowPath(ref) {
		if _, statErr := os.Stat(path); errors.Is(statErr, fs.ErrNotExist) {
			data, ok := workflow.Builtin(ref)
			if !ok {
				return "", nil, nil, Usagef("no workflow %q: %s does not exist and there is no built-in of that name (built-in: %s)",
					ref, path, strings.Join(workflow.BuiltinNames(), ", ")).WithNext(runStartNext)
			}
			def, err := workflow.Parse(ref, data)
			if err != nil {
				return "", nil, nil, Errorf("built-in workflow %s: %v", ref, err)
			}
			return builtinPrefix + ref, def, data, nil
		}
	}
	def, data, err := loadWorkflow(path, ref)
	if err != nil {
		return "", nil, nil, err
	}
	return path, def, data, nil
}

// loadWorkflow reads and validates the workflow at path. A missing file
// and an invalid workflow are usage errors carrying the reason.
func loadWorkflow(path, ref string) (*workflow.Definition, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, Usagef("no workflow %q: %s does not exist", ref, path).WithNext(runStartNext)
		}
		return nil, nil, Usagef("workflow %q: %v", ref, err).WithNext(runStartNext)
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	def, err := workflow.Parse(name, data)
	if err != nil {
		return nil, nil, Usagef("workflow %s: %v", path, err).WithNext(runStartNext)
	}
	return def, data, nil
}

// checkAgentSteps checks every agent step's effective adapter and
// permission mode (the step's, else mode and agent) as session.start
// checks its own: a mode is required, it must be the adapter's, and a mode
// that bypasses permission checks is refused. An autonomous run's arbiter
// (the file's arbiter step or the built-in one) is checked the same way
// when the workflow has an agent step: only an autonomous run launches it,
// and only on an agent step's question.
func checkAgentSteps(def *workflow.Definition, mode, agent string, autonomous bool, next string) error {
	agents := false
	for _, name := range def.Order {
		s := def.Steps[name]
		if s.Kind != workflow.KindAgent {
			continue
		}
		agents = true
		if err := checkAgentStep(s, "an agent step", mode, agent, next); err != nil {
			return err
		}
	}
	if autonomous && agents {
		return checkAgentStep(def.Arbiter(), "the arbiter step", mode, agent, next)
	}
	return nil
}

// checkAgentStep checks one step's effective adapter and permission mode
// for checkAgentSteps; what names the step's role in the message.
func checkAgentStep(s workflow.Step, what, mode, agent, next string) error {
	req := agentStepRequest(s, mode, agent)
	if req.PermissionMode == "" {
		return Usagef("step %q is %s with no permission mode: give --permission-mode or set permission_mode on the step", s.Name, what).
			WithNext(next)
	}
	if _, err := startAdapter(req); err != nil {
		e := AsError(err)
		return &Error{Code: e.Code, Message: fmt.Sprintf("step %q: %s", s.Name, e.Message), Next: next}
	}
	return nil
}

// agentStepRequest is the session.start request of an agent step: the
// step's adapter and mode, else the run's.
func agentStepRequest(s workflow.Step, mode, agent string) SessionStartRequest {
	r := SessionStartRequest{PermissionMode: s.PermissionMode, Agent: s.Agent, Step: s.Name}
	if r.PermissionMode == "" {
		r.PermissionMode = mode
	}
	if r.Agent == "" {
		r.Agent = agent
	}
	return r
}

// canRun refuses a caller that cannot run def's steps, before anything is
// written.
func canRun(ctx *Ctx, def *workflow.Definition) error {
	if ctx.Start == nil {
		return Errorf("cannot run a workflow: this caller starts no programs")
	}
	for _, s := range def.Steps {
		if s.Kind != workflow.KindAgent {
			continue
		}
		switch {
		case ctx.Self == "":
			return Errorf("cannot run agent step %q: the path of the running rhiza is unknown", s.Name)
		case ctx.Exec == nil:
			return Errorf("cannot run agent step %q: this caller runs no programs", s.Name)
		}
	}
	return nil
}

// newRunID draws a ULID from ctx.IDs that no run in store has yet. Test
// mode fixes the seed per process, so a later process can draw an ID an
// earlier one took; the next draw differs.
func newRunID(ctx *Ctx, rs *run.Store) (string, error) {
	for range maxRunIDDraws {
		u, err := ctx.IDs.New()
		if err != nil {
			return "", Errorf("%s", err.Error())
		}
		rid := u.String()
		_, err = rs.Read(rid)
		var nf *run.NotFoundError
		if errors.As(err, &nf) {
			return rid, nil
		}
		if err != nil {
			return "", Errorf("read run %s: %v", rid, err)
		}
	}
	return "", Errorf("no free run ID after %d draws", maxRunIDDraws)
}

// runEngine runs the run in the foreground until it finishes or stops,
// and responds with its state. Autonomous is the journal's: every agent
// the run launches is marked with it.
func runEngine(ctx *Ctx, rs *run.Store, common, rid, change, worktree, mode, agent string, autonomous, noWait bool) (any, error) {
	resp, err := runEngineRetry(ctx, rs, common, rid, change, worktree, mode, agent, autonomous, noWait, false)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// runEngineRetry is runEngine; with retry it runs the engine's Retry, which
// restarts a failed run at its failed step.
func runEngineRetry(ctx *Ctx, rs *run.Store, common, rid, change, worktree, mode, agent string, autonomous, noWait, retry bool) (RunResponse, error) {
	sessions := session.NewStore(session.Dir(common))
	e := &run.Engine{
		Store: rs,
		Runner: &stepRunner{
			ctx: ctx, change: change, run: rid, worktree: worktree, mode: mode, agent: agent,
			autonomous: autonomous, sessions: sessions,
		},
		Clock:  runClock,
		Now:    ctx.Now,
		NoWait: noWait,
	}
	if ctx.Stderr != nil {
		e.Log = func(line string) { fmt.Fprintf(ctx.Stderr, "rhiza: %s\n", line) }
	}
	if ctx.Rex != nil {
		e.Notify = func(n run.Notice) { ctx.Rex.Write(rexNotice(n)) }
	}
	run1 := e.Run
	if retry {
		run1 = e.Retry
	}
	st, err := run1(rid)
	switch {
	case errors.Is(err, run.ErrLive), errors.Is(err, run.ErrFinished), errors.Is(err, run.ErrNotFailed):
		return RunResponse{}, Conflictf("%v", err).WithNext("rhiza run status " + rid)
	case err != nil:
		return RunResponse{}, Errorf("run %s: %v", rid, err).WithNext("rhiza run status " + rid)
	}
	return RunResponse{Run: runInfo(st, false, sessions)}, nil
}

// stepRunner is the engine's run.Runner: command and wait steps run
// through `sh -c` in the run's worktree, agent steps through session
// start on the run's change set.
type stepRunner struct {
	ctx      *Ctx
	change   string
	run      string
	worktree string
	mode     string
	agent    string
	// autonomous marks every session the run launches, so the agent's
	// environment tells it the owner is not there.
	autonomous bool
	sessions   *session.Store
}

// Command runs the step's command, its output copied to ctx.Stderr.
func (r *stepRunner) Command(step workflow.Step, env []string) (int, string, error) {
	code, err := r.sh(step.Run, env, r.ctx.Stderr)
	return code, "", err
}

// Holds runs the wait step's condition once; exit 0 means it holds.
func (r *stepRunner) Holds(step workflow.Step, env []string) (bool, error) {
	code, err := r.sh(step.Until, env, nil)
	return code == 0, err
}

// dir is where command and wait steps run: the worktree, or the main
// checkout once the worktree is gone. The close-out step removes the
// worktree, so the steps after it run in the main checkout.
func (r *stepRunner) dir() string {
	if _, err := os.Stat(r.worktree); err == nil {
		return r.worktree
	}
	if out, err := r.ctx.Git("worktree", "list", "--porcelain"); err == nil {
		for _, line := range strings.Split(out, "\n") {
			if p, ok := strings.CutPrefix(line, "worktree "); ok {
				return p
			}
		}
	}
	return r.ctx.Root
}

// sh runs cmd with `sh -c` in the worktree, env laid over the caller's,
// copying its standard output to out (nil discards it), and returns its
// exit code. Only failing to run it is an error.
func (r *stepRunner) sh(cmd string, env []string, out io.Writer) (int, error) {
	p, err := r.ctx.Start(Program{Name: "sh", Args: []string{"-c", cmd}, Dir: r.dir(), Env: env})
	if err != nil {
		return 0, fmt.Errorf("start sh -c %q: %v", cmd, err)
	}
	if o := p.Output(); o != nil {
		if out == nil {
			out = io.Discard
		}
		io.Copy(out, o)
	}
	code, err := p.Wait()
	if err != nil {
		return 0, fmt.Errorf("wait for sh -c %q: %v", cmd, err)
	}
	return code, nil
}

// StartAgent starts the step's session as session.start does, on the
// run's change set, and returns its ID. An autonomous run's session is
// recorded autonomous.
func (r *stepRunner) StartAgent(step workflow.Step, _ []string) (string, error) {
	req := agentStepRequest(step, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	t, err := resolveChangeWorktree(r.ctx, r.change, "no agent can work on it")
	if err != nil {
		return "", err
	}
	req.Prompt = step.RenderPrompt(workflow.PromptVars{
		Change: t.cs.String(), Slug: session.Slug(t.title), Title: t.title, Run: r.run, Step: step.Name,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// arbiterPrompt is the arbiter's whole prompt, over the question node's
// ID: it names nothing else, so the arbiter reads the question and its
// context from the model, as the guide topic tells it.
const arbiterPrompt = "Decide question %s: follow rhiza guide arbiter."

// arbiterWording is the wording field of a Decision an arbiter writes: it
// assumed the answer in the owner's place.
const arbiterWording = "(assumed, would have escalated)"

// StartArbiter starts the arbiter step's session as StartAgent starts an
// agent step's (the step's adapter and mode, else the run's; autonomous as
// the run is), on the run's change set, with a prompt that names only
// question and the arbiter guide topic.
func (r *stepRunner) StartArbiter(step workflow.Step, question string, _ []string) (string, error) {
	req := agentStepRequest(step, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	req.Prompt = fmt.Sprintf(arbiterPrompt, question)
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// Answer reads the model as the run's change set sees it in the run's
// worktree, where its sessions write, and looks for a live Decision with a
// live answers link to question node question, wording arbiterWording and
// a reasoning. The first such Decision, by link ID, is the answer; with
// none, Why names what the nearest one lacks.
func (r *stepRunner) Answer(question string) (run.Answer, error) {
	q, err := id.Parse(question)
	if err != nil {
		return run.Answer{}, fmt.Errorf("question %q is not a node ID: %v", question, err)
	}
	cs, err := id.Parse(r.change)
	if err != nil {
		return run.Answer{}, fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	prefix, err := r.ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return run.Answer{}, err
	}
	ww, err := openWorkspace(filepath.Join(r.worktree, strings.TrimSpace(prefix)), r.ctx.Now, r.ctx.IDs, true)
	if err != nil {
		return run.Answer{}, err
	}
	return arbiterAnswer(ww.model(cs), q), nil
}

// questionHandlePrefix starts the handle of every question node a run adds
// for its arbiter: question-<sid8>-<n>.
const questionHandlePrefix = "question-"

// AddQuestion writes what session sid asked the owner, text, as one
// question node question-<sid8>-<n> (n the next number free in the change
// set's view) with fields statement and session, into the run's change
// set, in the workspace of the session's worktree, as the session's actor
// for its for-actor, in one entry. It refuses a session of another change
// set, and a change set that is finished.
func (r *stepRunner) AddQuestion(sid, text string) (run.Question, error) {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return run.Question{}, err
	}
	if rec.Change != r.change {
		return run.Question{}, Errorf("session %s is on change set %s, not the run's %s", rec.ID, rec.Change, r.change)
	}
	if len(rec.ID) < 8 {
		return run.Question{}, Errorf("session ID %q is shorter than 8 characters", rec.ID)
	}
	ww, cs, err := sessionWorkspace(r.ctx, rec)
	if err != nil {
		return run.Question{}, err
	}
	c, _ := ww.changeSet(cs)
	if c.State.Terminal() {
		return run.Question{}, Conflictf("change set %s is %s; no question can be added to it", cs, c.State).
			WithNext("rhiza change show " + cs.String())
	}
	m := ww.model(cs)
	prefix := questionHandlePrefix + strings.ToLower(rec.ID[:8]) + "-"
	handle := ""
	for n := 1; handle == ""; n++ {
		if h := prefix + strconv.Itoa(n); !answerProposed(m, c.Ops, h) {
			handle = h
		}
	}
	actor, err := liveActor(m, rec.Actor)
	if err != nil {
		return run.Question{}, Errorf("session %s's actor %s: %s", rec.ID, rec.Actor, AsError(err).Message)
	}
	forID, err := sessionForActor(r.ctx, m, rec.For)
	if err != nil {
		return run.Question{}, Errorf("session %s's for-actor %s: %s", rec.ID, rec.For, AsError(err).Message)
	}
	items := []applyItem{{op: "add", build: func(b writeBase) (edit, error) {
		return buildAdd(b, AddRequest{Kind: string(op.KindQuestion), Handle: handle,
			Fields: []string{"statement=" + text, "session=" + rec.ID}})
	}}}
	e, err := buildApply(ww, cs, r.ctx.IDs, items)
	if err == nil {
		err = e.check("add")
	}
	if err != nil {
		return run.Question{}, err
	}
	if len(e.nodes) != 1 {
		return run.Question{}, Errorf("adding question %s touched %d nodes, not 1", handle, len(e.nodes))
	}
	if _, err := ww.appendEntry(cs, changeset.Ops{Ops: e.ops, Drop: e.drop}, actor, forID); err != nil {
		return run.Question{}, err
	}
	return run.Question{ID: e.nodes[0].String(), Handle: handle}, nil
}

// deliverText is the message that gives the asking session the arbiter's
// answer: the question's ID, the answering Decision's handle and its
// statement.
func deliverText(question, decision, statement string) string {
	return fmt.Sprintf("The arbiter answered your question %s with Decision %s:\n\n%s\n\nCarry on with your work.",
		question, decision, strings.TrimSpace(statement))
}

// Deliver sends session sid the arbiter's answer to question as a session
// message from the run's caller, as run amend sends its pointer: the file
// is written, and the pointer typed once the session waits at its prompt.
// Then it clears the session's asked mark, so the session is no longer
// blocked on the question.
func (r *stepRunner) Deliver(sid, question, decision, statement string) error {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return err
	}
	if err := messageable(rec); err != nil {
		return err
	}
	from, err := messageSender(r.ctx)
	if err != nil {
		return err
	}
	dir, rel, err := messageFolder(r.ctx, rec)
	if err != nil {
		return err
	}
	text := deliverText(question, decision, statement)
	name, err := writeMessage(dir, messageContent(from, rec, rfc3339(r.ctx.Now()), text))
	if err != nil {
		return err
	}
	if _, err := typePointer(r.ctx, r.sessions, rec.ID, messagePointer(from.Handle, rel+"/"+name), defaultMessageTimeout); err != nil {
		return withNote(err, "the answer is in "+filepath.Join(dir, name))
	}
	if _, err := r.sessions.Update(rec.ID, func(x *session.Record) { x.Status.Asked = nil }); err != nil {
		return Errorf("clear session %s's asked mark: %v", rec.ID, err)
	}
	return nil
}

// arbiterAnswer is Answer's verdict over model m on question q.
func arbiterAnswer(m Model, q id.ID) run.Answer {
	if n, ok := m.Node(q); !ok || n.Retired {
		return run.Answer{Why: fmt.Sprintf("question %s is not a live node", q)}
	}
	why := fmt.Sprintf("no live Decision answers question %s", q)
	for _, l := range m.LinksTo(q, op.LinkAnswers) {
		d, ok := m.Node(l.From)
		if !ok || d.Retired || d.Kind != op.KindDecision {
			continue
		}
		switch {
		case fieldText(d.Fields["wording"]) != arbiterWording:
			why = fmt.Sprintf("decision %s answers question %s but its wording is %q, not %q", d.Handle, q, fieldText(d.Fields["wording"]), arbiterWording)
		case strings.TrimSpace(fieldText(d.Fields["reasoning"])) == "":
			why = fmt.Sprintf("decision %s answers question %s but has no reasoning", d.Handle, q)
		default:
			return run.Answer{OK: true, Decision: d.Handle, Statement: fieldText(d.Fields["statement"])}
		}
	}
	return run.Answer{Why: why}
}

// AgentStatus reads the session's record: it is done once its process
// has exited (exit-is-final), with pass for exit 0, else fail. Until then
// it is blocked while its resolved report is blocked, on that report's
// kind and detail, in the record's tab, idle while it resolves to idle,
// and said done once the record has the agent's done mark (session.done).
func (r *stepRunner) AgentStatus(sid string) (run.AgentStatus, error) {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return run.AgentStatus{}, err
	}
	if !rec.Status.Exited() {
		rep, _ := session.Resolve(rec)
		return run.AgentStatus{
			Blocked: recordBlock(rec), Idle: rep.State == session.StateIdle, SaidDone: rec.Status.SaidDone(),
		}, nil
	}
	x := rec.Status.Exit
	switch {
	case x.Code != nil && *x.Code == 0:
		return run.AgentStatus{Done: true, Outcome: workflow.OutcomePass, Detail: "exit 0"}, nil
	case x.Code != nil:
		return run.AgentStatus{Done: true, Outcome: workflow.OutcomeFail, Detail: fmt.Sprintf("exit %d", *x.Code)}, nil
	}
	detail := x.State
	if x.Detail != "" {
		detail += ": " + x.Detail
	}
	return run.AgentStatus{Done: true, Outcome: workflow.OutcomeFail, Detail: detail}, nil
}

// recordBlock is what a live session waits on when its resolved report
// is blocked, else nil. A live session with an asked mark (a question its
// hook denied in an autonomous run) waits on that question whatever its
// reports say, until the answer is delivered.
func recordBlock(rec session.Record) *run.Block {
	if rec.Status.Exited() {
		return nil
	}
	rep, _ := session.Resolve(rec)
	if a := rec.Status.Asked; a != nil {
		rep = *a
	} else if rep.State != session.StateBlocked {
		return nil
	}
	return &run.Block{
		Kind: rep.Kind, Question: rep.Detail,
		TabSession: rec.Tab.Session, TabWindow: rec.Tab.Window, TabBlock: rec.Tab.Block, TabTitle: rec.Tab.Title,
	}
}

// runBlockedNow is what an unfinished run waits on now: its current
// step's session, read from its record, when that resolves to blocked.
// A record that cannot be read leaves the run not blocked.
func runBlockedNow(st run.State, sessions *session.Store) *RunBlocked {
	sid := st.Session[st.Current]
	if st.Ended || sid == "" || sessions == nil {
		return nil
	}
	rec, err := sessions.Read(sid)
	if err != nil {
		return nil
	}
	b := recordBlock(rec)
	if b == nil {
		return nil
	}
	return &RunBlocked{
		Step: st.Current, Session: sid, Kind: b.Kind, Question: b.Question,
		Tab: SessionTab{Session: b.TabSession, Window: b.TabWindow, Block: b.TabBlock, Title: b.TabTitle},
	}
}

// Limits of a Program Status report's fields, in decoded bytes.
const (
	rexTitleMax = 192
	rexMsgMax   = 2048
)

// rexNotice is the OSC 7501 root-record report of a notice: blocked, with
// the kind, the agent's tab as title and a message naming the change set,
// the step, the tab and the question; or working once the block ends.
func rexNotice(n run.Notice) []byte {
	if !n.Blocked || n.Block == nil {
		return rexReport("working", "", "", "")
	}
	b := n.Block
	what := n.Title
	if what == "" {
		what = n.Change
	}
	tab := b.TabTitle
	if tab == "" {
		tab = n.Session
	}
	msg := fmt.Sprintf("%s: %s waits on you in tab %s", what, n.Step, tab)
	if q := strings.TrimSpace(b.Question); q != "" {
		msg += ": " + q
	}
	return rexReport("blocked", b.Kind, b.TabTitle, msg)
}

// rexReport encodes one OSC 7501 root-record report from app rhiza. The
// kind is kept only on a blocked state and only when Rex knows it; title
// and msg have their control characters replaced with spaces, are capped
// on a UTF-8 boundary and are base64-encoded, and are left out when empty.
func rexReport(state, kind, title, msg string) []byte {
	var b strings.Builder
	b.WriteString("\x1b]7501;state=")
	b.WriteString(state)
	if state == session.StateBlocked {
		switch kind {
		case session.KindPermission, session.KindQuestion, session.KindAuth:
			b.WriteString(":kind=" + kind)
		}
	}
	b.WriteString(":app=rhiza")
	if t := rexField(title, rexTitleMax); t != "" {
		b.WriteString(":title=" + base64.StdEncoding.EncodeToString([]byte(t)))
	}
	if m := rexField(msg, rexMsgMax); m != "" {
		b.WriteString(":msg=" + base64.StdEncoding.EncodeToString([]byte(m)))
	}
	b.WriteString("\x1b\\")
	return []byte(b.String())
}

// rexField is s as valid UTF-8 with every control character replaced by
// a space, trimmed and cut to at most limit bytes without splitting a rune.
func rexField(s string, limit int) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	if len(s) <= limit {
		return s
	}
	cut := limit
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return strings.TrimSpace(s[:cut])
}

// runInfo is the response form of a folded journal; live says another
// process holds the run. An unfinished run whose current step's session
// is blocked reads blocked, with what it waits on from sessions.
func runInfo(st run.State, live bool, sessions *session.Store) RunInfo {
	info := RunInfo{
		ID: st.ID, Workflow: st.Workflow, Change: st.Change, Title: st.Title, Worktree: st.Worktree,
		Autonomous: st.Autonomous, Current: st.Current, Steps: []RunStep{}, StartedAt: rfc3339(st.StartedAt),
		Amendments: []RunAmendment{},
	}
	switch {
	case st.Ended:
		info.State, info.Detail = st.Result, st.Detail
	case live:
		info.State = runRunning
	default:
		info.State, info.Detail = runStopped, st.Stopped
	}
	if b := runBlockedNow(st, sessions); b != nil {
		info.State, info.Blocked = runBlocked, b
	}
	for _, s := range st.Steps {
		info.Steps = append(info.Steps, RunStep{
			Step: s.Step, Visit: s.Visit, Outcome: s.Outcome, Next: s.Next, Detail: s.Detail,
			Session: s.Session, At: rfc3339(s.At), Retry: s.Retried,
		})
	}
	for _, rt := range st.Retries {
		info.Retries = append(info.Retries, RunRetry{Step: rt.Step, Detail: rt.Detail, At: rfc3339(rt.At)})
	}
	for _, a := range st.Amendments {
		info.Amendments = append(info.Amendments, RunAmendment{
			Step: a.Step, Session: a.Session, Actor: a.Actor, File: a.File, At: rfc3339(a.At),
		})
	}
	for _, a := range st.Accepts {
		info.Accepts = append(info.Accepts, RunAccept{Change: a.Change, Entry: a.LogEntry, Verb: a.Verb, At: rfc3339(a.At)})
	}
	for _, q := range st.Asked {
		rq := RunQuestion{
			Question: q.Question, Handle: q.Handle, Text: q.Text, Step: q.Step,
			Decision: q.Answered, Answered: q.Answered != "", At: rfc3339(q.At),
		}
		if a, ok := st.Arbitration(q.Question); ok {
			rq.Outcome = a.Outcome
			if rq.Decision == "" && a.Outcome == workflow.OutcomePass {
				rq.Decision = a.Decision
			}
		}
		info.Questions = append(info.Questions, rq)
	}
	return info
}

// Human reports the run, or the plan of a dry run.
func (r RunResponse) Human(w io.Writer) {
	if r.DryRun && r.Retry != "" {
		fmt.Fprintf(w, "Run %s would be retried at step %s\n\n", r.Run.ID, r.Retry)
		runLines(w, r.Run)
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or run.\n")
		return
	}
	if r.DryRun {
		fmt.Fprintf(w, "Run %s would start\n\n", r.Run.ID)
		runLines(w, r.Run)
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or run.\n")
		return
	}
	r.Run.Human(w)
}

// Human shows the run and the steps it has done.
func (r RunInfo) Human(w io.Writer) {
	fmt.Fprintf(w, "Run %s %s\n\n", r.ID, r.State)
	runLines(w, r)
	if len(r.Steps) > 0 {
		fmt.Fprintf(w, "\n## Steps\n\n")
		for _, s := range r.Steps {
			to := "end"
			if s.Next != "" {
				to = s.Next
			}
			retry := ""
			if s.Retry {
				retry = ", retried"
			}
			fmt.Fprintf(w, "- %s (visit %d%s): %s, then %s", s.Step, s.Visit, retry, s.Outcome, to)
			if s.Session != "" {
				fmt.Fprintf(w, ", session %s", s.Session)
			}
			if s.Detail != "" {
				fmt.Fprintf(w, ": %s", s.Detail)
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.Accepts) > 0 {
		fmt.Fprintf(w, "\n## Accepted in the owner's place\n\n")
		for _, a := range r.Accepts {
			fmt.Fprintf(w, "- %s %s change set %s (log entry %s)\n", a.At, a.Verb, a.Change, a.Entry)
		}
	}
	if len(r.Questions) > 0 {
		fmt.Fprintf(w, "\n## Questions the arbiter decided\n\n")
		for _, q := range r.Questions {
			fmt.Fprintf(w, "- %s (%s): %s", q.Handle, q.Question, q.Text)
			switch {
			case q.Decision != "":
				fmt.Fprintf(w, " — decision %s", q.Decision)
				if !q.Answered {
					fmt.Fprintf(w, ", not yet delivered")
				}
			case q.Outcome != "":
				fmt.Fprintf(w, " — arbiter failed: %s", q.Outcome)
			default:
				fmt.Fprintf(w, " — open")
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.Retries) > 0 {
		fmt.Fprintf(w, "\n## Retries\n\n")
		for _, rt := range r.Retries {
			fmt.Fprintf(w, "- %s retried at %s", rt.At, rt.Step)
			if rt.Detail != "" {
				fmt.Fprintf(w, ", after: %s", rt.Detail)
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.Amendments) > 0 {
		fmt.Fprintf(w, "\n## Amendments\n\n")
		for _, a := range r.Amendments {
			fmt.Fprintf(w, "- %s %s amended %s (step %s)\n", a.At, a.Actor, a.File, a.Step)
		}
	}
	switch r.State {
	case runStopped:
		fmt.Fprintf(w, "\nnext: rhiza run resume %s\n", r.ID)
	case runRunning:
		fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.ID)
	case runBlocked:
		if r.Blocked != nil && r.Blocked.Tab.Block != "" {
			fmt.Fprintf(w, "\nnext: rex focus -b %s\n", r.Blocked.Tab.Block)
		} else if r.Blocked != nil {
			fmt.Fprintf(w, "\nnext: rhiza session show %s\n", r.Blocked.Session)
		}
	}
}

// runLines are a run's workflow, change set, worktree, whether it is
// autonomous, current step and detail, one per line.
func runLines(w io.Writer, r RunInfo) {
	fmt.Fprintf(w, "- Workflow: %s\n", r.Workflow)
	fmt.Fprintf(w, "- Change set: %s: %s\n", r.Change, r.Title)
	fmt.Fprintf(w, "- Worktree: %s\n", r.Worktree)
	if r.Autonomous {
		fmt.Fprintf(w, "- Autonomous: yes\n")
	}
	fmt.Fprintf(w, "- Started: %s\n", r.StartedAt)
	if r.Current != "" {
		fmt.Fprintf(w, "- Current step: %s\n", r.Current)
	}
	if r.Detail != "" {
		fmt.Fprintf(w, "- Detail: %s\n", r.Detail)
	}
	if b := r.Blocked; b != nil {
		tab := b.Tab.Title
		if tab == "" {
			tab = b.Session
		}
		fmt.Fprintf(w, "- Blocked: %s waits on the owner in tab %s", b.Step, tab)
		if b.Tab.Window != "" {
			fmt.Fprintf(w, " (window %s)", b.Tab.Window)
		}
		fmt.Fprintf(w, "\n")
		if b.Question != "" {
			fmt.Fprintf(w, "- Question: %s\n", b.Question)
		}
	}
}

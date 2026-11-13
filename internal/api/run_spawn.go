package api

// This file holds run.spawn: a session in a run starts a child session in
// the run's worktree and waits for it. The command runs in the caller's
// process, never the daemon's; the child is started as any run session is
// (through the run's agent adapter and permission mode), is journaled as a
// child of the calling step, and counts against no rhiza.daemon.maxRuns: only
// a usage hold or the machine load hold makes the spawn wait.

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "run.spawn",
		Summary:  "Start a child session of the run RHIZA_RUN names in the run's worktree, with the brief file given as its brief, wait for it (journaled as held while a usage or machine hold is on) and print its outcome",
		Write:    true,
		Request:  reflect.TypeFor[RunSpawnRequest](),
		Response: reflect.TypeFor[RunSpawnResponse](),
		Run:      runRunSpawn,
	})
}

const runSpawnNext = "rhiza describe run.spawn"

// spawnPoll is how often a spawn re-reads the holds and the child's record.
const spawnPoll = 2 * time.Second

// RunSpawnRequest names the child's role and brief.
type RunSpawnRequest struct {
	Role  string `json:"role" flag:"role" desc:"The child's role, such as builder or verifier: the guide it follows"`
	Leaf  string `json:"leaf,omitempty" flag:"leaf" desc:"A leaf of the run's change set the child works on, by handle; default none"`
	Brief string `json:"brief" flag:"brief" desc:"A file holding the caller's brief for the child: its brief.md holds the contract section, then this file"`
	Model string `json:"model,omitempty" flag:"model" desc:"The model the child's agent starts on; default the agent's"`
}

// RunSpawnResponse is the child's outcome.
type RunSpawnResponse struct {
	Run     string `json:"run" desc:"The run's full ID"`
	Step    string `json:"step" desc:"The calling step the child is a child of"`
	Parent  string `json:"parent,omitempty" desc:"The calling session"`
	Session string `json:"session,omitempty" desc:"The child's session; absent in a dry run"`
	Role    string `json:"role" desc:"The child's role"`
	Leaf    string `json:"leaf,omitempty" desc:"The leaf it worked on"`
	Brief   string `json:"brief,omitempty" desc:"The child's brief.md, an absolute path"`
	Held    bool   `json:"held,omitempty" desc:"True if the spawn waited on a usage or machine hold before it started the child"`
	Outcome string `json:"outcome,omitempty" enum:"pass,fail" desc:"pass when the child's agent exited 0, else fail"`
	Detail  string `json:"detail,omitempty" desc:"The child's exit, or why it failed"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was started or journaled"`
}

// ExitCode is 1 when the child failed, else 0.
func (r RunSpawnResponse) ExitCode() int {
	if r.Outcome == workflow.OutcomeFail {
		return 1
	}
	return 0
}

// Human prints the outcome in one line.
func (r RunSpawnResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would spawn a %s session under step %s of run %s\n", r.Role, r.Step, r.Run)
		return
	}
	fmt.Fprintf(w, "Spawned %s session %s under step %s of run %s: %s", r.Role, r.Session, r.Step, r.Run, r.Outcome)
	if r.Detail != "" {
		fmt.Fprintf(w, " (%s)", r.Detail)
	}
	fmt.Fprintln(w)
}

// spawner is one spawn's working state. start launches the child: the run's
// session.start path, which tests replace.
type spawner struct {
	ctx   *Ctx
	start func(SessionStartRequest) (string, error)
}

func runRunSpawn(ctx *Ctx, req any) (any, error) {
	return spawner{ctx: ctx, start: func(r SessionStartRequest) (string, error) {
		resp, err := runSessionStart(ctx, r)
		if err != nil {
			return "", err
		}
		return resp.(SessionStartResponse).Session.ID, nil
	}}.spawn(req.(RunSpawnRequest))
}

func (sp spawner) spawn(r RunSpawnRequest) (RunSpawnResponse, error) {
	ctx := sp.ctx
	rr, err := openReviewRun(ctx, "run spawn", runSpawnNext)
	if err != nil {
		return RunSpawnResponse{}, err
	}
	role := strings.TrimSpace(r.Role)
	if role == "" {
		return RunSpawnResponse{}, Usagef("run spawn needs --role: the child's role, such as builder").WithNext(runSpawnNext)
	}
	if strings.TrimSpace(r.Brief) == "" {
		return RunSpawnResponse{}, Usagef("run spawn needs --brief: a file holding the child's brief").WithNext(runSpawnNext)
	}
	path := r.Brief
	if !filepath.IsAbs(path) && ctx.Wd != "" {
		path = filepath.Join(ctx.Wd, path)
	}
	brief, err := os.ReadFile(path)
	if err != nil {
		return RunSpawnResponse{}, Usagef("run spawn cannot read --brief %s: %v", r.Brief, err).WithNext(runSpawnNext)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return RunSpawnResponse{}, err
	}
	getenv := ctxGetenv(ctx)
	parent := strings.TrimSpace(getenv(envSession))
	step := strings.TrimSpace(getenv("RHIZA_STEP"))
	sessions := session.NewStore(session.Dir(common))
	if parent != "" {
		if rec, err := sessions.Read(parent); err == nil {
			parent = rec.ID
			if rec.Step != "" {
				step = rec.Step
			}
		}
	}
	if step == "" {
		step = rr.state.Current
	}
	resp := RunSpawnResponse{Run: rr.id, Step: step, Parent: parent, Role: role, Leaf: strings.TrimSpace(r.Leaf)}
	if ctx.DryRun {
		resp.DryRun = true
		return resp, nil
	}

	// Wait while a hold is on, journaled once per wait.
	rs := rr.store
	for held := false; ; held = true {
		why := spawnHoldDetail(ctx, common)
		if why == "" {
			break
		}
		if !held {
			resp.Held = true
			if err := rs.Append(rr.id, run.Entry{Type: run.TypeSpawnHeld, At: ctx.Now(), Step: step, Parent: parent, Detail: why, Role: role}); err != nil {
				return resp, Errorf("journal the held spawn in run %s: %v", rr.id, err)
			}
		}
		spawnSleep(ctx)
	}

	req := SessionStartRequest{
		Change: rr.state.Change, PermissionMode: rr.state.PermissionMode, Agent: rr.state.Agent, Step: step, Run: rr.id,
		autonomous: rr.state.Autonomous, landedOK: true,
		role: role, leaf: resp.Leaf, model: strings.TrimSpace(r.Model), parent: parent,
	}
	if resp.Leaf != "" {
		req.attempt = 1
	}
	req.Prompt = fmt.Sprintf("You are a %s session spawned by session %s of run %s. Read your brief at $RHIZA_BRIEF, then follow `rhiza guide %s`.", role, parent, rr.id, role)
	req.briefFor = func(sid string) (string, error) {
		dir := run.KeyedSessionDir(common, rr.id, "spawn", sid)
		resp.Brief = filepath.Join(dir, briefFile)
		return dir, run.WriteBriefFile(dir, run.RenderSpawnBrief(string(brief)))
	}
	sid, err := sp.start(req)
	if err != nil {
		return resp, err
	}
	resp.Session = sid
	if err := rs.Append(rr.id, run.Entry{Type: run.TypeSpawn, At: ctx.Now(), Step: step, Parent: parent, Session: sid, Role: role, Leaf: resp.Leaf, Model: req.model}); err != nil {
		return resp, Errorf("journal the child session %s in run %s: %v", sid, rr.id, err)
	}
	for {
		rec, err := readSession(sessions, sid)
		if err != nil {
			return resp, err
		}
		if outcome, detail, done := childOutcome(rec); done {
			resp.Outcome, resp.Detail = outcome, detail
			break
		}
		spawnSleep(ctx)
	}
	if _, err := rs.AppendSpawnEnded(rr.id, run.Entry{At: ctx.Now(), Step: step, Parent: parent, Session: sid, Outcome: resp.Outcome, Detail: resp.Detail}); err != nil {
		return resp, Errorf("journal the end of child session %s in run %s: %v", sid, rr.id, err)
	}
	return resp, nil
}

// spawnSleep waits one poll.
func spawnSleep(ctx *Ctx) {
	if ctx.Sleep != nil {
		ctx.Sleep(spawnPoll)
		return
	}
	time.Sleep(spawnPoll)
}

// childOutcome is a child session's outcome once its process has exited:
// pass for exit 0, else fail with the exit or the reason it ended.
func childOutcome(rec session.Record) (outcome, detail string, done bool) {
	if !rec.Status.Exited() {
		return "", "", false
	}
	x := rec.Status.Exit
	switch {
	case x.Code != nil && *x.Code == 0:
		return workflow.OutcomePass, "exit 0", true
	case x.Code != nil:
		return workflow.OutcomeFail, fmt.Sprintf("exit %d", *x.Code), true
	}
	detail = x.State
	if x.Detail != "" {
		detail += ": " + x.Detail
	}
	return workflow.OutcomeFail, detail, true
}

package api

// This file holds the start of a scheduled run: run start --schedule <name>.
// A scheduled run has no change set, no worktree and no branch; its steps run
// in the project's main checkout.

import (
	"bytes"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// PickSchedule is the Pick of a scheduled run.
const PickSchedule = "schedule"

// scheduleCommandStep is the one step of a command schedule's implicit
// workflow.
const scheduleCommandStep = "command"

// scheduleWorkflow is the schedule's workflow: the one it names, resolved as
// run start resolves a name, or its command wrapped as an implicit one-step
// workflow with a single command step.
func scheduleWorkflow(root string, sc dispatch.Schedule) (string, *workflow.Definition, []byte, error) {
	if sc.Workflow != "" {
		return resolveWorkflow(root, sc.Workflow)
	}
	var buf bytes.Buffer
	doc := map[string]any{"step": map[string]any{scheduleCommandStep: map[string]any{"kind": workflow.KindCommand, "run": sc.Command}}}
	if err := toml.NewEncoder(&buf).Encode(doc); err != nil {
		return "", nil, nil, Errorf("schedule %q: wrap its command: %v", sc.Name, err)
	}
	def, err := workflow.Parse(sc.Name, buf.Bytes())
	if err != nil {
		return "", nil, nil, Errorf("schedule %q: wrap its command: %v", sc.Name, err)
	}
	return "schedule:" + sc.Name, def, buf.Bytes(), nil
}

func runScheduledStart(ctx *Ctx, r RunStartRequest) (any, error) {
	name := strings.TrimSpace(r.Schedule)
	for _, c := range []struct {
		set  bool
		flag string
	}{
		{strings.TrimSpace(r.Change) != "", "a change set"}, {strings.TrimSpace(r.Workflow) != "", "--workflow"},
		{strings.TrimSpace(r.New) != "", "--new"}, {r.Milestone, "--milestone"},
		{strings.TrimSpace(r.Ref) != "", "--ref"}, {strings.TrimSpace(r.Request) != "", "--request"},
	} {
		if c.set {
			return nil, Usagef("--schedule starts a run with no change set: it does not take %s", c.flag).WithNext(runStartNext)
		}
	}
	rules, err := dispatch.Load(ctx.Root)
	if err != nil {
		return nil, Usagef("%s: %v", dispatch.File, err).WithNext(runStartNext)
	}
	var sc dispatch.Schedule
	found := false
	for _, s := range rules.Schedules {
		if s.Name == name {
			sc, found = s, true
		}
	}
	if !found {
		return nil, Usagef("no schedule %q in %s", name, dispatch.File).WithNext(runStartNext)
	}
	path, def, data, err := scheduleWorkflow(ctx.Root, sc)
	if err != nil {
		return nil, err
	}
	// Nobody is there to answer: a scheduled run is always autonomous.
	r.Autonomous = true
	if err := checkAgentSteps(ctx, def, r.PermissionMode, r.Agent, true, runStartNext); err != nil {
		return nil, err
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	rid, err := newRunID(ctx, rs)
	if err != nil {
		return nil, err
	}
	started := run.Started{
		Workflow: def.Name, WorkflowPath: path, Pick: PickSchedule, Definition: string(data), Includes: def.Includes,
		PermissionMode: r.PermissionMode, Agent: r.Agent, Autonomous: true,
		Trigger: run.TriggerSchedule, TriggerName: sc.Name,
	}
	if ctx.DryRun {
		return RunResponse{Run: RunInfo{
			ID: rid, Workflow: def.Name, Pick: PickSchedule, Trigger: started.Trigger, TriggerName: sc.Name,
			Autonomous: true, Current: def.Start, Steps: []RunStep{}, StartedAt: rfc3339(ctx.Now()), Amendments: []RunAmendment{},
			Messages: []RunMessage{}, Sessions: []RunSession{},
		}, DryRun: true}, nil
	}
	if ctx.Getenv != nil {
		if sid := strings.TrimSpace(ctx.Getenv(envSession)); sid != "" {
			return nil, Usagef("run start refused: this is agent session %s, and no agent session starts a run", sid).
				WithNext("ask the owner to run it outside an agent session")
		}
	}
	if err := canRun(ctx, def); err != nil {
		return nil, err
	}
	claimStartLog(ctx, common, rid)
	if err := startRun(rs, rid, started, ctx.Now(), def); err != nil {
		return nil, err
	}
	return runEngine(ctx, rs, common, rid, "", "", r.PermissionMode, r.Agent, true, r.NoWait)
}

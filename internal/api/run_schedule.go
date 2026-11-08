package api

// This file holds the start of a scheduled run: run start --schedule <name>.
// A scheduled run has no change set, no worktree and no branch; its steps run
// in the project's main checkout.

import (
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
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
	main := filepath.Dir(common)
	before, snapErr := snapshotCheckout(ctx, main)
	resp, err := runEngine(ctx, rs, common, rid, "", "", r.PermissionMode, r.Agent, true, r.NoWait)
	if err != nil || snapErr != nil {
		return resp, err
	}
	rr, ok := resp.(RunResponse)
	if !ok || rr.Run.State == runRunning || rr.Run.State == runStopped || rr.Run.State == runBlocked {
		// The run goes on elsewhere: its end is not here.
		return resp, nil
	}
	after, err := snapshotCheckout(ctx, main)
	if err != nil {
		return resp, nil
	}
	if w := before.changes(after); w != "" {
		// A warning only: nothing is reverted or blocked.
		if err := rs.Append(rid, run.Entry{Type: run.TypeWarning, At: ctx.Now(), Detail: w}); err == nil {
			rr.Run.Warnings = append(rr.Run.Warnings, w)
			if ctx.Stderr != nil {
				fmt.Fprintf(ctx.Stderr, "rhiza: warning: scheduled run %s (schedule %s) %s\n", rid, sc.Name, w)
			}
		}
	}
	return rr, nil
}

// checkoutSnapshot is the main checkout's branch, HEAD and uncommitted
// tracked changes (path to its content's hash, "deleted" when gone) at one
// moment.
type checkoutSnapshot struct {
	branch, head string
	files        map[string]string
}

func snapshotCheckout(ctx *Ctx, dir string) (checkoutSnapshot, error) {
	git := func(args ...string) (string, error) {
		out, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
		return strings.TrimSpace(string(out)), err
	}
	var s checkoutSnapshot
	var err error
	if s.head, err = git("rev-parse", "HEAD"); err != nil {
		return s, err
	}
	// Detached HEAD reads as HEAD itself, which never equals a branch name.
	if s.branch, err = git("rev-parse", "--abbrev-ref", "HEAD"); err != nil {
		return s, err
	}
	names, err := git("diff", "HEAD", "--name-only")
	if err != nil {
		return s, err
	}
	s.files = map[string]string{}
	for _, f := range strings.Split(names, "\n") {
		if f = strings.TrimSpace(f); f == "" {
			continue
		}
		h, err := git("hash-object", "--", f)
		if err != nil {
			h = "deleted"
		}
		s.files[f] = h
	}
	return s, nil
}

// changes names what the run changed between s (its start) and after (its
// end): the branch, the HEAD and each tracked file whose uncommitted state
// differs; empty when nothing changed.
func (s checkoutSnapshot) changes(after checkoutSnapshot) string {
	var parts []string
	if s.branch != after.branch {
		parts = append(parts, fmt.Sprintf("switched the main checkout's branch from %s to %s", s.branch, after.branch))
	}
	if s.head != after.head {
		parts = append(parts, fmt.Sprintf("moved the main checkout's HEAD from %s to %s", s.head, after.head))
	}
	var files []string
	for f, h := range after.files {
		if old, ok := s.files[f]; !ok || old != h {
			files = append(files, f)
		}
	}
	for f := range s.files {
		if _, ok := after.files[f]; !ok {
			files = append(files, f)
		}
	}
	if len(files) > 0 {
		slices.Sort(files)
		parts = append(parts, "changed uncommitted tracked files in the main checkout: "+strings.Join(files, ", "))
	}
	return strings.Join(parts, "; ")
}

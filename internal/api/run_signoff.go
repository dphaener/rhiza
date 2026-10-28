package api

// This file holds the owner's sign-off of a review: run.approve and
// run.request-changes. A review's signoff step stops the run, waiting for
// the owner; these ops take the run by ref or the latest (the owner has
// no RHIZA_RUN), journal the sign-off and continue the run in-process, as
// run.resume does. approve ends the step pass, so the run goes on to its
// end; it never lands the change set. request-changes journals the
// owner's findings, each naming a check of the checks step the signoff
// step's on.changes names, as already confirmed: the step ends changes,
// and that checks step fixes them in one owner cycle outside the cap,
// re-runs the checks named, then the run waits at the sign-off again.

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "run.approve",
		Summary:  "Sign off a review run waiting at its signoff step as the owner: journal the approval and continue the run, which ends done; it never lands the change set",
		Write:    true,
		Request:  reflect.TypeFor[RunApproveRequest](),
		Response: reflect.TypeFor[RunSignoffResponse](),
		Run:      runRunApprove,
	})
	register(Operation{
		Name:     "run.request-changes",
		Summary:  "Ask for changes to a review run waiting at its signoff step: journal the owner's findings as confirmed and continue the run, which fixes them in one cycle outside the cap, re-runs the checks they name and waits for the sign-off again",
		Write:    true,
		Request:  reflect.TypeFor[RunRequestChangesRequest](),
		Response: reflect.TypeFor[RunSignoffResponse](),
		Input: &Input{
			Summary: `{"findings":[{"check":...,"severity":...,"leaf":...,"place":...,"reason":...},...]}, at least one; check names a check of the checks step the sign-off goes back to; severity is one of ` +
				strings.Join(run.Severities, ", ") + `, leaf optional`,
			Schema: changesFileSchema,
		},
		Run: runRunRequestChanges,
	})
}

const runRequestChangesNext = "rhiza describe run.request-changes"

// RunApproveRequest names the run to approve.
type RunApproveRequest struct {
	Run string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
}

// RunRequestChangesRequest names the file of the owner's findings and the
// run.
type RunRequestChangesRequest struct {
	File string `json:"file" arg:"0" desc:"A JSON file of the owner's findings, or - for stdin (see the input schema)"`
	Run  string `json:"run,omitempty" arg:"1" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
}

// RunSignoff is one sign-off of a review's signoff step.
type RunSignoff struct {
	Step     string       `json:"step" desc:"The signoff step"`
	Visit    int          `json:"visit" desc:"Which visit of the step it signed off, from 1"`
	Actor    string       `json:"actor,omitempty" desc:"Who signed off, by actor handle; absent under rule autonomous"`
	Rule     string       `json:"rule" desc:"owner, or autonomous when an autonomous run signed off in the owner's place"`
	Outcome  string       `json:"outcome" desc:"approved, or changes when the owner asked for changes"`
	Findings []RunFinding `json:"findings,omitempty" desc:"With changes: the owner's findings, each naming a check, which count as confirmed"`
	At       string       `json:"at" desc:"When, RFC 3339 UTC"`
}

// RunSignoffResponse is the sign-off journalled and the run as it stands
// once continued.
type RunSignoffResponse struct {
	Signoff RunSignoff `json:"signoff" desc:"The sign-off journalled, or in a dry run the one that would be"`
	Changes string     `json:"changes,omitempty" desc:"With changes: the checks step that fixes them"`
	Run     RunInfo    `json:"run" desc:"The run after it was continued; in a dry run, as it stands. The command exits 1 when the run failed"`
	DryRun  bool       `json:"dry_run,omitempty" desc:"True if nothing was journalled or run"`
}

// ExitCode is 1 when the run failed, else 0.
func (r RunSignoffResponse) ExitCode() int {
	if r.Run.State == run.Failed {
		return 1
	}
	return 0
}

// signoffRun is a run waiting at its signoff step for the owner.
type signoffRun struct {
	store  *run.Store
	common string
	j      run.Journal
	state  run.State
	def    *workflow.Definition
	step   workflow.Step
	cmd    string
}

// openSignoffRun reads the run ref names, or the latest, and refuses it,
// as usage, unless it is unfinished, not running, and stopped at a
// signoff step whose visit has no sign-off yet.
func openSignoffRun(ctx *Ctx, cmd, ref string) (signoffRun, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return signoffRun{}, err
	}
	rs := run.NewStore(run.Dir(common))
	var j run.Journal
	if ref = strings.TrimSpace(ref); ref != "" {
		j, err = readRun(rs, ref)
	} else {
		j, err = latestRun(ctx, rs)
	}
	if err != nil {
		if e := AsError(err); e.Code == CodeNotFound {
			return signoffRun{}, Usagef("%s works on a review run waiting at its sign-off: %s", cmd, e.Message).WithNext(e.Next)
		}
		return signoffRun{}, err
	}
	st := j.State()
	if rs.Live(j.ID) {
		return signoffRun{}, Usagef("%s works on a review run waiting at its sign-off, and run %s is running in another process", cmd, j.ID).WithNext("rhiza run status " + j.ID)
	}
	def, err := workflow.Parse(st.Workflow, []byte(st.Definition))
	if err != nil {
		return signoffRun{}, Errorf("run %s: its workflow no longer parses: %v", j.ID, err)
	}
	step, err := signoffWaitingStep(cmd, j.ID, st, def)
	if err != nil {
		return signoffRun{}, err
	}
	return signoffRun{store: rs, common: common, j: j, state: st, def: def, step: step, cmd: cmd}, nil
}

// signoffWaitingStep is the signoff step run id waits at, in state st, or
// a usage refusal: the run is finished, not stopped at a signoff step, or
// that step's visit is already signed off.
func signoffWaitingStep(cmd, id string, st run.State, def *workflow.Definition) (workflow.Step, error) {
	status := "rhiza run status " + id
	if st.Ended {
		return workflow.Step{}, Usagef("%s works on a review run waiting at its sign-off, and run %s is finished: it %s", cmd, id, st.Result).WithNext(status)
	}
	step, ok := def.Steps[st.Current]
	if !ok || step.Kind != workflow.KindSignoff || st.CurrentVisit == 0 || st.Stopped == "" {
		return workflow.Step{}, Usagef("%s works on a review run waiting at its sign-off, and run %s is not: it is %s",
			cmd, id, runWhere(st)).WithNext(status)
	}
	if s, ok := st.SignoffFor(step.Name, st.CurrentVisit); ok {
		return workflow.Step{}, Usagef("%s: visit %d of step %s of run %s is already signed off (%s); rhiza run resume continues it",
			cmd, st.CurrentVisit, step.Name, id, s.Outcome).WithNext("rhiza run resume " + id)
	}
	return step, nil
}

// runWhere says where an unfinished run is, for a refusal.
func runWhere(st run.State) string {
	switch {
	case st.Current == "":
		return "between steps"
	case st.CurrentVisit == 0:
		return "about to enter step " + st.Current
	case st.Stopped != "":
		return "stopped at step " + st.Current + ": " + st.Stopped
	}
	return "at step " + st.Current
}

// signoff journals e on the run and continues it in-process, as run
// resume does; a dry run journals and runs nothing. e is journalled only
// under the run's lock, after the checks openSignoffRun made are made
// again on the journal read under it: a sign-off that lost a race to
// another process exits 3 (the run is running) or 2 (it was signed off
// or moved on) and journals nothing.
func (sr signoffRun) signoff(ctx *Ctx, e run.Entry, changes string) (RunSignoffResponse, error) {
	resp := RunSignoffResponse{Signoff: runSignoff(run.Signoff{
		Step: e.Step, Visit: e.Visit, Actor: e.Actor, Rule: e.Rule, Outcome: e.Outcome, Findings: e.Findings, At: e.At,
	}), Changes: changes, DryRun: ctx.DryRun}
	sessions := session.NewStore(session.Dir(sr.common))
	if ctx.DryRun {
		resp.Run = runInfo(sr.state, false, sessions)
		return resp, nil
	}
	st := sr.state
	if err := checkAgentSteps(sr.def, st.PermissionMode, st.Agent, st.Autonomous, "rhiza run status "+sr.j.ID); err != nil {
		return RunSignoffResponse{}, err
	}
	if err := canRun(ctx, sr.def); err != nil {
		return RunSignoffResponse{}, err
	}
	before := func(now run.State) ([]run.Entry, error) {
		step, err := signoffWaitingStep(sr.cmd, sr.j.ID, now, sr.def)
		if err != nil {
			return nil, err
		}
		if step.Name != e.Step || now.CurrentVisit != e.Visit {
			return nil, Usagef("%s works on a review run waiting at its sign-off, and run %s moved on: it is %s",
				sr.cmd, sr.j.ID, runWhere(now)).WithNext("rhiza run status " + sr.j.ID)
		}
		return []run.Entry{e}, nil
	}
	r, err := runEngineRetry(ctx, sr.store, sr.common, sr.j.ID, st.Change, st.Worktree, st.PermissionMode, st.Agent, st.Autonomous, false, false, before)
	if err != nil {
		return RunSignoffResponse{}, err
	}
	resp.Run = r.Run
	return resp, nil
}

// runRunApprove journals the owner's approval of the run's sign-off
// (rule owner, the caller's actor) and continues the run: its signoff
// step passes.
func runRunApprove(ctx *Ctx, req any) (any, error) {
	r := req.(RunApproveRequest)
	sr, err := openSignoffRun(ctx, "run approve", r.Run)
	if err != nil {
		return nil, err
	}
	from, err := messageSender(ctx)
	if err != nil {
		return nil, err
	}
	return sr.signoff(ctx, run.Entry{
		Type: run.TypeSignoff, At: ctx.Now(), Step: sr.step.Name, Visit: sr.state.CurrentVisit,
		Actor: from.Handle, Rule: run.RuleOwner, Outcome: run.SignoffApproved,
	}, "")
}

// changesFile is run.request-changes' input.
type changesFile struct {
	Findings []struct {
		Check    string `json:"check"`
		Severity string `json:"severity"`
		Leaf     string `json:"leaf"`
		Place    string `json:"place"`
		Reason   string `json:"reason"`
	} `json:"findings"`
}

// runRunRequestChanges checks the run, the file and each finding (a check
// of the checks step on.changes names, a known severity, a place, a
// reason and, when it names one, a leaf of the run's change set), and
// that another visit of the checks step and of the signoff step stays
// within its max_visits. Only then does it journal the owner's request
// (rule owner, the caller's actor, the findings numbered f1, f2, ...) and
// continue the run: its signoff step ends changes.
func runRunRequestChanges(ctx *Ctx, req any) (any, error) {
	r := req.(RunRequestChangesRequest)
	const cmd = "run request-changes"
	next := runRequestChangesNext
	if strings.TrimSpace(r.File) == "" {
		return nil, Usagef("%s needs a file of findings, or - for stdin", cmd).WithNext(next)
	}
	sr, err := openSignoffRun(ctx, cmd, r.Run)
	if err != nil {
		return nil, err
	}
	target := sr.step.On[workflow.OutcomeChanges]
	checks, ok := sr.def.Steps[target]
	if !ok || checks.Kind != workflow.KindChecks {
		return nil, Errorf("run %s: step %s's on.changes names %q, which is not a checks step", sr.j.ID, sr.step.Name, target)
	}
	for _, s := range []workflow.Step{checks, sr.step} {
		if n := sr.state.CapVisits(s.Name); n+1 > s.Cap() {
			return nil, Usagef("%s: another visit of step %s of run %s would go over its max_visits %d, so no more changes can be fixed in this run; rhiza run approve signs it off",
				cmd, s.Name, sr.j.ID, s.Cap()).WithNext("rhiza run approve " + sr.j.ID)
		}
	}
	var f changesFile
	if err := readReviewFile(ctx, cmd, next, r.File, &f); err != nil {
		return nil, err
	}
	if len(f.Findings) == 0 {
		return nil, Usagef(`%s file: want {"findings":[...]} with at least one finding; to sign off as it is, rhiza run approve`, cmd).WithNext(next)
	}
	rr := reviewRun{store: sr.store, id: sr.j.ID, state: sr.state}
	if rr.change, err = id.Parse(sr.state.Change); err != nil {
		return nil, Errorf("run %s names change set %q: %v", sr.j.ID, sr.state.Change, err)
	}
	leafOf := changeLeafResolver(ctx, rr)
	found := make([]run.Finding, 0, len(f.Findings))
	for i, in := range f.Findings {
		n := i + 1
		check := strings.TrimSpace(in.Check)
		if _, ok := checks.Check(check); !ok {
			return nil, Usagef("finding %d: check %q is not a check of step %s: its checks are %s",
				n, in.Check, checks.Name, strings.Join(checkNames(checks), ", ")).WithNext(next)
		}
		if !run.ValidSeverity(in.Severity) {
			return nil, Usagef("finding %d: severity %q is not one of %s", n, in.Severity, strings.Join(run.Severities, ", ")).WithNext(next)
		}
		if strings.TrimSpace(in.Place) == "" {
			return nil, Usagef("finding %d: no place: give file:line or a node", n).WithNext(next)
		}
		if strings.TrimSpace(in.Reason) == "" {
			return nil, Usagef("finding %d: no reason", n).WithNext(next)
		}
		fd := run.Finding{ID: "f" + strconv.Itoa(n), Check: check, Severity: in.Severity, Place: in.Place, Reason: in.Reason}
		if ref := strings.TrimSpace(in.Leaf); ref != "" {
			node, err := leafOf(ref)
			if err != nil {
				if e := AsError(err); e.Code != CodeUsage {
					return nil, err
				}
				return nil, Usagef("finding %d: leaf %q: %s", n, ref, AsError(err).Message).WithNext("rhiza change show " + rr.change.String())
			}
			fd.Leaf, fd.LeafHandle = node.ID.String(), node.Handle
		}
		found = append(found, fd)
	}
	from, err := messageSender(ctx)
	if err != nil {
		return nil, err
	}
	return sr.signoff(ctx, run.Entry{
		Type: run.TypeSignoff, At: ctx.Now(), Step: sr.step.Name, Visit: sr.state.CurrentVisit,
		Actor: from.Handle, Rule: run.RuleOwner, Outcome: run.SignoffChanges, Findings: found,
	}, checks.Name)
}

// runSignoff is the response form of sign-off s.
func runSignoff(s run.Signoff) RunSignoff {
	out := RunSignoff{Step: s.Step, Visit: s.Visit, Actor: s.Actor, Rule: s.Rule, Outcome: s.Outcome, At: rfc3339(s.At)}
	for _, f := range s.Findings {
		out.Findings = append(out.Findings, runFinding(f, run.Round{}))
	}
	return out
}

// line is the sign-off on one line, as run status lists it.
func (s RunSignoff) line() string {
	who := s.Actor
	if s.Rule != run.RuleOwner || who == "" {
		who = "rule " + s.Rule
	}
	if s.Outcome == run.SignoffChanges {
		return fmt.Sprintf("%s (visit %d): changes asked by %s", s.Step, s.Visit, who)
	}
	return fmt.Sprintf("%s (visit %d): %s by %s", s.Step, s.Visit, s.Outcome, who)
}

// Human reports the sign-off, then the run.
func (r RunSignoffResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Sign-off: %s\n", r.Signoff.line())
	for _, f := range r.Signoff.Findings {
		fmt.Fprintf(w, "- %s %s\n", f.Check, findingLine(f))
	}
	if r.Changes != "" {
		fmt.Fprintf(w, "Step %s fixes them, then the run waits for the sign-off again.\n", r.Changes)
	}
	fmt.Fprintf(w, "\n")
	if r.DryRun {
		r.Run.Human(w)
		fmt.Fprintf(w, "\nThis was a dry run: nothing was journalled or run.\n")
		return
	}
	r.Run.Human(w)
}

// changesFileSchema is the JSON Schema of run.request-changes' input.
func changesFileSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"findings": map[string]any{
				"type":        "array",
				"minItems":    1,
				"description": "The owner's findings, in order, numbered f1, f2, ...; they count as confirmed",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"check":    map[string]any{"type": "string", "description": "A check of the checks step the sign-off goes back to: it runs again once the findings are fixed"},
						"severity": map[string]any{"enum": stringsAny(run.Severities), "description": "How serious it is"},
						"leaf":     map[string]any{"type": "string", "description": "The leaf of the run's change set it is about, by handle or ID; without one, the change set's fixer takes it"},
						"place":    map[string]any{"type": "string", "minLength": 1, "description": "Where: file:line or a node"},
						"reason":   map[string]any{"type": "string", "minLength": 1, "description": "What to change, and why"},
					},
					"required":             []string{"check", "place", "reason", "severity"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"findings"},
		"additionalProperties": false,
	}
}

// signoffWaiting says whether the run waits at a signoff step for the
// owner.
func signoffWaiting(st run.State) bool {
	return !st.Ended && st.Stopped == run.SignoffWait
}

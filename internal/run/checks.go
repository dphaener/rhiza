package run

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

// CheckRunner is the check-level surface a Runner gives a checks step. A
// Runner without it cannot run one: the step ends with an error outcome.
// It embeds LeafRunner, whose leaf sessions a checks step's fixers and
// verifiers reuse, and whose ClosePane closes a check session's pane too.
// A check session is polled with the Runner's AgentStatus and its end
// reported with AgentEnded, as an agent step's is.
type CheckRunner interface {
	LeafRunner
	// StartCheck launches the session spec names and returns its ID.
	StartCheck(spec CheckSpec) (session string, err error)
}

// CheckSpec is one check session to launch: role Role (check) of checks
// step Step for Check, in cycle Cycle (from 1). Model is the model to
// run, "" for the agent's default, from the check's ladder at Cycle. Env
// holds the session's RHIZA_* variables, the step's plus RHIZA_CHECK,
// RHIZA_CYCLE and RHIZA_ROLE. Dir is where the session works, "" for the
// run's worktree. Anchor is the check session whose pane the session's
// pane splits, the split's direction following Step's split; "" opens the
// change set's tab anew. Its prompt names only the change set and the
// check: the Runner renders Check's prompt, nothing more.
type CheckSpec struct {
	Step   workflow.Step
	Role   string
	Check  workflow.Check
	Cycle  int
	Model  string
	Dir    string
	Env    []string
	Anchor string
}

// The environment variables a check session gets beyond a step's; it
// gets EnvRole too.
const (
	EnvCheck = "RHIZA_CHECK"
	EnvCycle = "RHIZA_CYCLE"
)

// MaxFixCycles is how many further cycles a checks step gives the checks
// that failed its first: a check failing all 1 + MaxFixCycles cycles
// fails the step.
const MaxFixCycles = 2

// checkActionKind is what the conductor does next for a checks step.
type checkActionKind int

const (
	// checkCycle journals the step's next cycle over checks.
	checkCycle checkActionKind = iota + 1
	// checkStart launches the check's session in the cycle.
	checkStart
	// checkPoll polls the check's live session in the cycle.
	checkPoll
	// checkEnd journals the check's outcome in the cycle.
	checkEnd
	// checksDone ends the step with outcome and detail.
	checksDone
)

// checkAction is one move of the checks conductor, decided from the
// folded journal alone.
type checkAction struct {
	kind    checkActionKind
	check   string
	role    string
	cycle   int
	checks  []string
	session string
	outcome string
	detail  string
}

// nextCheckActions is the pure decision of a checks step: from its
// checks, its cycles in order and its folded check states by name, the
// moves to make now. Cycle 1 runs every check; each later cycle only
// those the cycle before failed, up to MaxFixCycles further cycles. In a
// cycle every check's session starts before any is polled, so they all
// run in parallel: the moves list each start before any poll or end. A
// check judged is ended by checkVerdict. Once each check of the cycle has
// an outcome the next cycle begins over the failed ones; with none
// failed the step passes, and with no cycle left it fails.
func nextCheckActions(step workflow.Step, cycles []CheckCycle, states map[string]CheckState) []checkAction {
	if len(cycles) == 0 {
		var all []string
		for _, c := range step.Checks {
			all = append(all, c.Name)
		}
		return []checkAction{{kind: checkCycle, cycle: 1, checks: all}}
	}
	cur := cycles[len(cycles)-1]
	var starts, rest []checkAction
	var failed []string
	for _, name := range cur.Checks {
		c := states[name]
		if o, ok := c.Outcome(cur.Cycle); ok {
			if o.Outcome != CheckPassed {
				failed = append(failed, name)
			}
			continue
		}
		s, started := c.Session(workflow.RoleCheck, cur.Cycle)
		if !started {
			starts = append(starts, checkAction{kind: checkStart, check: name, role: workflow.RoleCheck, cycle: cur.Cycle})
			continue
		}
		if outcome, detail, ok := checkVerdict(c, cur.Cycle); ok {
			rest = append(rest, checkAction{kind: checkEnd, check: name, cycle: cur.Cycle, outcome: outcome, detail: detail})
			continue
		}
		rest = append(rest, checkAction{kind: checkPoll, check: name, role: workflow.RoleCheck, cycle: cur.Cycle, session: s.Session})
	}
	if acts := append(starts, rest...); len(acts) > 0 {
		return acts
	}
	if len(failed) > 0 && cur.Cycle <= MaxFixCycles {
		return []checkAction{{kind: checkCycle, cycle: cur.Cycle + 1, checks: failed}}
	}
	passed := 0
	var lost []string
	for _, c := range step.Checks {
		if o, ok := states[c.Name].Latest(); ok && o.Outcome == CheckPassed {
			passed++
		} else {
			lost = append(lost, c.Name)
		}
	}
	if len(lost) == 0 {
		return []checkAction{{kind: checksDone, outcome: workflow.OutcomePass,
			detail: fmt.Sprintf("%d of %d checks passed", passed, len(step.Checks))}}
	}
	return []checkAction{{kind: checksDone, outcome: workflow.OutcomeFail,
		detail: fmt.Sprintf("%d of %d checks passed after %d cycles; failed: %s", passed, len(step.Checks), cur.Cycle, strings.Join(lost, ", "))}}
}

// checkVerdict is the check's outcome in cycle once its session's work is
// judged: passed when the session passed, else failed; ok false while the
// session is unjudged. It is the seam where a check's findings and its
// critic's verdicts will weigh in.
func checkVerdict(c CheckState, cycle int) (outcome, detail string, ok bool) {
	j, judged := c.Judgement(workflow.RoleCheck, cycle)
	if !judged {
		return "", "", false
	}
	if j.Outcome == workflow.OutcomePass {
		return CheckPassed, fmt.Sprintf("cycle %d: check passed", cycle), true
	}
	return CheckFailed, fmt.Sprintf("cycle %d: check failed: %s", cycle, j.Detail), true
}

// checkStates are step's folded checks, by name.
func (r *running) checkStates(step string) map[string]CheckState {
	out := map[string]CheckState{}
	for _, c := range r.j.State().Checks {
		if c.Step == step {
			out[c.Check] = c
		}
	}
	return out
}

// stepCycles are step's cycles, in order.
func (r *running) stepCycles(step string) []CheckCycle {
	var out []CheckCycle
	for _, c := range r.j.State().Cycles {
		if c.Step == step {
			out = append(out, c)
		}
	}
	return out
}

// checks conducts a checks step over its checks: one loop on this
// goroutine that acts on nextCheckActions, journalling each change, and
// sleeps the step's interval when a turn started, judged and ended
// nothing. Each turn polls every live check session. A session in the
// journal is rejoined, never started again. A session is judged once its
// process exits (pass on exit 0) or its agent says done (pass); idle
// alone is never enough. Blocks, notices and an autonomous run's
// arbitration are handled as on an agent step. Runner failures are an
// error outcome; the error return is for journal failures.
//
// Each session runs in a pane of the change set's tab, placed by
// checkPane. When the step ends, whatever its outcome, the pane of every
// judged session still open is closed; a live session's pane stays, as a
// resume rejoins it.
func (r *running) checks(step workflow.Step, env []string) (result, error) {
	cr, ok := r.e.Runner.(CheckRunner)
	if !ok {
		return errResult(fmt.Errorf("step %q is a checks step, and this runner cannot run checks", step.Name)), nil
	}
	res, err := r.conductChecks(cr, step, env)
	if err != nil {
		return res, err
	}
	if err := r.closeEndedCheckPanes(cr, step); err != nil {
		return result{}, err
	}
	return res, nil
}

// conductChecks is the checks step's loop, as checks describes.
func (r *running) conductChecks(cr CheckRunner, step workflow.Step, env []string) (result, error) {
	interval := step.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	for {
		progressed := false
		for _, a := range nextCheckActions(step, r.stepCycles(step.Name), r.checkStates(step.Name)) {
			switch a.kind {
			case checksDone:
				return result{outcome: a.outcome, detail: a.detail}, nil
			case checkCycle:
				if err := r.append(Entry{Type: TypeCheckCycle, Step: step.Name, Cycle: a.cycle, Checks: a.checks}); err != nil {
					return result{}, err
				}
				r.log("step %s: cycle %d: %s", step.Name, a.cycle, strings.Join(a.checks, ", "))
				progressed = true
			case checkStart:
				if res, failed, err := r.startCheck(cr, step, env, a); err != nil || failed {
					return res, err
				}
				progressed = true
			case checkEnd:
				if err := r.append(Entry{Type: TypeCheck, Step: step.Name, Check: a.check, Cycle: a.cycle, Outcome: a.outcome, Detail: a.detail}); err != nil {
					return result{}, err
				}
				r.log("step %s: check %s cycle %d: %s", step.Name, a.check, a.cycle, a.outcome)
				progressed = true
			case checkPoll:
				judged, res, ended, err := r.pollCheck(step, a)
				if err != nil || ended {
					return res, err
				}
				progressed = progressed || judged
			}
		}
		if !progressed {
			r.clock().Sleep(interval)
		}
	}
}

// startCheck launches a's session, on the check's model for the cycle, in
// a pane checkPane places, and journals its check-session entry. A launch
// that fails ends the step (failed true) with an error outcome.
func (r *running) startCheck(cr CheckRunner, step workflow.Step, env []string, a checkAction) (result, bool, error) {
	check, ok := step.Check(a.check)
	if !ok {
		return errResult(fmt.Errorf("check %s is not in step %q", a.check, step.Name)), true, nil
	}
	spec := CheckSpec{
		Step: step, Role: a.role, Check: check, Cycle: a.cycle, Model: check.Model(a.cycle),
		Env: append(append([]string{}, env...),
			EnvCheck+"="+a.check, EnvCycle+"="+strconv.Itoa(a.cycle), EnvRole+"="+a.role),
		Anchor: checkPane(r.checkStates(step.Name)),
	}
	s, err := cr.StartCheck(spec)
	if err != nil {
		return errResult(fmt.Errorf("check %s: start its session: %v", a.check, err)), true, nil
	}
	if err := r.append(Entry{
		Type: TypeCheckSession, Step: step.Name, Check: a.check, Role: a.role, Cycle: a.cycle,
		Session: s, Model: spec.Model, Detail: spec.Dir,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: check %s cycle %d: %s session %s", step.Name, a.check, a.cycle, a.role, s)
	return result{}, false, nil
}

// pollCheck polls a's session once. When it has exited or its agent said
// done, it journals the session's judgement (judged true) and reports
// its end to the Runner. A live session's block is journalled and
// notified as an agent step's, and in an autonomous run its question goes
// to the arbiter; a failed verdict ends the step (ended true) with res.
func (r *running) pollCheck(step workflow.Step, a checkAction) (judged bool, res result, ended bool, err error) {
	status, err := r.e.Runner.AgentStatus(a.session)
	if err != nil {
		if err := r.unblock(step.Name, a.session); err != nil {
			return false, result{}, false, err
		}
		return false, errResult(fmt.Errorf("check %s: %s session %s: %v", a.check, a.role, a.session, err)), true, nil
	}
	if status.Done || status.SaidDone {
		if err := r.unblock(step.Name, a.session); err != nil {
			return false, result{}, false, err
		}
		outcome, detail := workflow.OutcomePass, "said done"
		if status.Done {
			outcome, detail = status.Outcome, status.Detail
			if outcome == "" {
				outcome = workflow.OutcomeFail
			}
		}
		if err := r.append(Entry{
			Type: TypeCheckJudged, Step: step.Name, Check: a.check, Role: a.role, Cycle: a.cycle, Outcome: outcome, Detail: detail,
		}); err != nil {
			return false, result{}, false, err
		}
		r.log("step %s: check %s cycle %d: %s %s", step.Name, a.check, a.cycle, a.role, outcome)
		if err := r.e.Runner.AgentEnded(step, a.session); err != nil {
			r.log("step %s: %v", step.Name, err)
		}
		return true, result{}, false, nil
	}
	if status.Blocked != nil {
		err = r.block(step.Name, a.session, status.Blocked)
	} else {
		err = r.unblock(step.Name, a.session)
	}
	if err != nil {
		return false, result{}, false, err
	}
	if r.arbitrates(status.Blocked) {
		res, ended, err := r.ask(step, a.session, status.Blocked)
		if err != nil || ended {
			return false, res, ended, err
		}
	}
	return false, result{}, false, nil
}

// openCheckPane is check's session whose pane is still open, and
// whether the session was judged.
type openCheckPane struct {
	check  string
	s      CheckSession
	judged bool
}

// openCheckPanes are the step's check sessions whose panes are open (no
// pane-closed entry names them), in the order they started.
func openCheckPanes(states map[string]CheckState) []openCheckPane {
	var out []openCheckPane
	for _, c := range states {
		for _, s := range c.Sessions {
			if s.PaneClosed {
				continue
			}
			_, judged := c.Judgement(s.Role, s.Cycle)
			out = append(out, openCheckPane{check: c.Check, s: s, judged: judged})
		}
	}
	slices.SortFunc(out, func(a, b openCheckPane) int { return a.s.Seq - b.s.Seq })
	return out
}

// checkPane is the pure choice of where a new check session's pane opens,
// from the step's check states: the most recent open check pane, which
// it splits, or "" for a new tab when no pane of the step is open. So a
// cycle's checks share the change set's tab, each in its own pane.
func checkPane(states map[string]CheckState) string {
	open := openCheckPanes(states)
	if n := len(open); n > 0 {
		return open[n-1].s.Session
	}
	return ""
}

// closeEndedCheckPanes closes the pane of every judged check session of
// step still open, in the order they started, as closePane does a leaf
// session's: a failure is logged and leaves the pane open in the
// journal. The error return is for journal failures.
func (r *running) closeEndedCheckPanes(cr CheckRunner, step workflow.Step) error {
	for _, p := range openCheckPanes(r.checkStates(step.Name)) {
		if !p.judged {
			continue
		}
		if err := cr.ClosePane(p.s.Session); err != nil {
			r.log("step %s: close the pane of session %s: %v", step.Name, p.s.Session, err)
			continue
		}
		if err := r.append(Entry{Type: TypePaneClosed, Step: step.Name, Check: p.check, Session: p.s.Session}); err != nil {
			return err
		}
	}
	return nil
}

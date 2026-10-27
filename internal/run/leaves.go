package run

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

// LeafRunner is the leaf-level surface a Runner gives a leaves step. A
// Runner without it cannot run one: the step ends with an error outcome.
// A leaf session is polled with the Runner's AgentStatus and its end
// reported with AgentEnded, as an agent step's is.
type LeafRunner interface {
	// Leaves are change set change's leaves, in its leaf order: each
	// depends only on leaves before it.
	Leaves(change string) ([]Leaf, error)
	// StartLeaf launches the session spec names and returns its ID.
	StartLeaf(spec LeafSpec) (session string, err error)
}

// Leaf is one leaf of a change set: its capability node's ID and handle,
// and the IDs of the leaves among the change set's that it depends on.
type Leaf struct {
	ID        string
	Handle    string
	DependsOn []string
}

// LeafSpec is one leaf session to launch: role Role (builder or verifier)
// of leaves step Step on Leaf, for attempt Attempt (from 1). Model is the
// model to run, "" for the agent's default. Env holds the session's
// RHIZA_* variables, the step's plus RHIZA_LEAF and RHIZA_ATTEMPT.
type LeafSpec struct {
	Step    workflow.Step
	Role    string
	Leaf    Leaf
	Attempt int
	Model   string
	Env     []string
}

// The environment variables a leaf session gets beyond a step's.
const (
	EnvLeaf    = "RHIZA_LEAF"
	EnvAttempt = "RHIZA_ATTEMPT"
)

// leafActionKind is what the conductor does next for a leaves step.
type leafActionKind int

const (
	// leafStart launches role's session on the leaf's attempt.
	leafStart leafActionKind = iota + 1
	// leafPoll polls the live session of role on the leaf's attempt.
	leafPoll
	// leafEnd journals the leaf's outcome.
	leafEnd
	// leavesDone ends the step with outcome and detail.
	leavesDone
)

// leafAction is one move of the conductor, decided from the folded
// journal alone.
type leafAction struct {
	kind    leafActionKind
	leaf    Leaf
	role    string
	attempt int
	session string
	outcome string
	detail  string
}

// nextLeafActions is the pure decision of a leaves step: from the change
// set's leaves, in leaf order, and the step's folded leaf states, the moves
// to make now. Leaves go one at a time, in order. A leaf's attempt starts
// with its builder; a builder judged pass is followed by a verifier, whose
// pass passes the leaf. A fail of either fails the leaf. Once every leaf
// has an outcome the step ends: pass when every leaf passed, else fail.
func nextLeafActions(leaves []Leaf, states map[string]LeafState) []leafAction {
	passed, failed := 0, []string{}
	for _, lf := range leaves {
		l := states[lf.ID]
		switch l.Outcome {
		case LeafPassed:
			passed++
			continue
		case LeafFailed, LeafSkipped:
			failed = append(failed, lf.Handle)
			continue
		}
		return []leafAction{nextForLeaf(lf, l)}
	}
	if len(failed) == 0 {
		return []leafAction{{kind: leavesDone, outcome: workflow.OutcomePass, detail: fmt.Sprintf("%d of %d leaves passed", passed, len(leaves))}}
	}
	return []leafAction{{kind: leavesDone, outcome: workflow.OutcomeFail,
		detail: fmt.Sprintf("%d of %d leaves passed; not passed: %s", passed, len(leaves), strings.Join(failed, ", "))}}
}

// nextForLeaf is the next move on an open leaf l.
func nextForLeaf(lf Leaf, l LeafState) leafAction {
	attempt := max(1, l.Attempt())
	for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
		s, started := l.Session(role, attempt)
		if !started {
			return leafAction{kind: leafStart, leaf: lf, role: role, attempt: attempt}
		}
		j, judged := l.Judged(role, attempt)
		if !judged {
			return leafAction{kind: leafPoll, leaf: lf, role: role, attempt: attempt, session: s.Session}
		}
		if j.Outcome != workflow.OutcomePass {
			return leafAction{kind: leafEnd, leaf: lf, outcome: LeafFailed,
				detail: fmt.Sprintf("attempt %d: %s failed: %s", attempt, role, j.Detail)}
		}
	}
	return leafAction{kind: leafEnd, leaf: lf, outcome: LeafPassed, detail: fmt.Sprintf("attempt %d: verifier passed", attempt)}
}

// leafStates are step's folded leaves, by leaf ID.
func (r *running) leafStates(step string) map[string]LeafState {
	out := map[string]LeafState{}
	for _, l := range r.j.State().Leaves {
		if l.Step == step {
			out[l.Leaf] = l
		}
	}
	return out
}

// leaves conducts a leaves step over the run change set's leaves: one loop
// on this goroutine that acts on nextLeafActions, journalling each change,
// and sleeps the step's interval when a turn only polled. A session in the
// journal is rejoined, never started again. A session's attempt is judged
// once its process exits (pass on exit 0) or its agent says done (pass);
// idle alone is never enough. Blocks, notices and an autonomous run's
// arbitration are handled as on an agent step. Runner failures are an
// error outcome; the error return is for journal failures.
func (r *running) leaves(step workflow.Step, env []string) (result, error) {
	lr, ok := r.e.Runner.(LeafRunner)
	if !ok {
		return errResult(fmt.Errorf("step %q is a leaves step, and this runner cannot run leaves", step.Name)), nil
	}
	leaves, err := lr.Leaves(r.j.State().Change)
	if err != nil {
		return errResult(err), nil
	}
	interval := step.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	for {
		waiting := false
		for _, a := range nextLeafActions(leaves, r.leafStates(step.Name)) {
			switch a.kind {
			case leavesDone:
				return result{outcome: a.outcome, detail: a.detail}, nil
			case leafStart:
				if res, failed, err := r.startLeaf(lr, step, env, a); err != nil || failed {
					return res, err
				}
			case leafEnd:
				if err := r.append(Entry{
					Type: TypeLeaf, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle, Outcome: a.outcome, Detail: a.detail,
				}); err != nil {
					return result{}, err
				}
				r.log("step %s: leaf %s %s", step.Name, a.leaf.Handle, a.outcome)
			case leafPoll:
				judged, res, ended, err := r.pollLeaf(step, a)
				if err != nil || ended {
					return res, err
				}
				waiting = waiting || !judged
			}
		}
		if waiting {
			r.clock().Sleep(interval)
		}
	}
}

// startLeaf launches a's session and journals its leaf-session entry. A
// launch that fails ends the step (failed true) with an error outcome.
func (r *running) startLeaf(lr LeafRunner, step workflow.Step, env []string, a leafAction) (result, bool, error) {
	spec := LeafSpec{
		Step: step, Role: a.role, Leaf: a.leaf, Attempt: a.attempt,
		Env: append(append([]string{}, env...), EnvLeaf+"="+a.leaf.Handle, EnvAttempt+"="+strconv.Itoa(a.attempt)),
	}
	s, err := lr.StartLeaf(spec)
	if err != nil {
		return errResult(fmt.Errorf("leaf %s: start its %s: %v", a.leaf.Handle, a.role, err)), true, nil
	}
	if err := r.append(Entry{
		Type: TypeLeafSession, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle,
		Role: a.role, Attempt: a.attempt, Session: s, Model: spec.Model,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: leaf %s attempt %d: %s session %s", step.Name, a.leaf.Handle, a.attempt, a.role, s)
	return result{}, false, nil
}

// pollLeaf polls a's session once. When it has exited or its agent said
// done, it journals the attempt's judgement (judged true) and reports the
// session's end to the Runner. A live session's block is journalled and
// notified as an agent step's, and in an autonomous run its question goes
// to the arbiter; a failed verdict ends the step (ended true) with res.
func (r *running) pollLeaf(step workflow.Step, a leafAction) (judged bool, res result, ended bool, err error) {
	status, err := r.e.Runner.AgentStatus(a.session)
	if err != nil {
		if err := r.unblock(step.Name, a.session); err != nil {
			return false, result{}, false, err
		}
		return false, errResult(fmt.Errorf("leaf %s: %s session %s: %v", a.leaf.Handle, a.role, a.session, err)), true, nil
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
			Type: TypeLeafAttempt, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle,
			Role: a.role, Attempt: a.attempt, Outcome: outcome, Detail: detail,
		}); err != nil {
			return false, result{}, false, err
		}
		r.log("step %s: leaf %s attempt %d: %s %s", step.Name, a.leaf.Handle, a.attempt, a.role, outcome)
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

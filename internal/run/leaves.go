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
	// Head is the commit at HEAD of the worktree builders work in.
	Head() (string, error)
	// LeafCommit is the newest commit in since..HEAD of that worktree
	// whose Leaf: trailer names leaf's handle. When there is none, commit
	// is "" and detail says why: no new commit, or the new commits name
	// other leaves or none.
	LeafCommit(leaf Leaf, since string) (commit, detail string, err error)
	// Evidence reports whether leaf has live Evidence that verifies it,
	// created by the actor of verifier session session. When not, detail
	// says why: "" when no Evidence verifies the leaf at all, else which
	// other actor wrote the Evidence that does.
	Evidence(leaf Leaf, session string) (ok bool, detail string, err error)
	// CommitEvidence commits, on HEAD of the builders' worktree, only the
	// log entry files holding the Evidence of leaf that verifier session
	// session's actor wrote, as its own commit, and returns it. The
	// builder's staged and unstaged edits stay as they were. Called again
	// once that commit exists, it returns it and commits nothing.
	CommitEvidence(leaf Leaf, session string) (commit string, err error)
	// OpenCheckout adds a throwaway detached worktree at commit, where
	// leaf's verifier on attempt works, and returns its path. It never
	// touches the builders' worktree.
	OpenCheckout(leaf Leaf, attempt int, commit string) (dir string, err error)
	// RemoveCheckout removes a checkout OpenCheckout added; one already
	// gone is no error.
	RemoveCheckout(dir string) error
	// AddLeafQuestion adds a question node, statement text, to the run's
	// change set, linked blocks to leaf, and returns it. It refuses when
	// the change set is finished.
	AddLeafQuestion(leaf Leaf, text string) (Question, error)
	// LeafAnswer reports whether a live Decision answers question node
	// question (a live answers link to it), and if so that Decision's
	// handle and statement. Any Decision does: the owner's needs no
	// arbiter wording.
	LeafAnswer(question string) (decision, statement string, ok bool, err error)
}

// Leaf is one leaf of a change set: its capability node's ID and handle,
// and the IDs of the leaves among the change set's that it depends on.
type Leaf struct {
	ID        string
	Handle    string
	DependsOn []string
}

// LeafSpec is one leaf session to launch: role Role (builder or verifier)
// of leaves step Step on Leaf, for attempt Attempt (from 1, counted across
// rounds), which is attempt RoundAttempt (from 1 to MaxAttempts) of the
// leaf's round Round (from 1). Model is the
// model to run, "" for the agent's default, from the role's ladder at
// RoundAttempt. In a fresh round (Round past 1), Question is the question
// node (ID) the leaf's failure raised, Decision the handle of the Decision
// that answers it and Answer its statement. Commit is, for a verifier,
// the builder's commit it checks; "" for a builder. Env holds the
// session's RHIZA_* variables, the step's plus RHIZA_LEAF, RHIZA_ATTEMPT,
// RHIZA_ROLE and, for a verifier, RHIZA_COMMIT. Dir is, for a verifier,
// the checkout of Commit it works in; "" for a builder.
type LeafSpec struct {
	Step         workflow.Step
	Role         string
	Leaf         Leaf
	Attempt      int
	Round        int
	RoundAttempt int
	Model        string
	Commit       string
	Dir          string
	Env          []string
	Question     string
	Decision     string
	Answer       string
}

// The environment variables a leaf session gets beyond a step's.
const (
	EnvLeaf    = "RHIZA_LEAF"
	EnvAttempt = "RHIZA_ATTEMPT"
	// EnvRole is the session's role: builder or verifier.
	EnvRole = "RHIZA_ROLE"
	// EnvCommit is the commit a verifier checks.
	EnvCommit = "RHIZA_COMMIT"
)

// MaxAttempts is how many attempts a leaf gets in a round: each failed
// attempt but the last goes to a fresh builder session for the next
// attempt, on the next model of the builder's ladder; once the last fails
// the leaf fails. A fresh round (a leaf-round entry) gives it MaxAttempts
// more, the ladder starting over.
const MaxAttempts = 3

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
	// commit is, for a verifier's start, the builder's commit it checks,
	// and for a passed leaf's end its Evidence commit.
	commit  string
	outcome string
	detail  string
}

// nextLeafActions is the pure decision of a leaves step: from the change
// set's leaves, in leaf order, and the step's folded leaf states, the moves
// to make now. A leaf's attempt starts with its builder; a builder judged
// pass is followed by its verifier, whose pass passes the leaf. A fail of
// either goes to a fresh builder on the next attempt, up to MaxAttempts;
// the last one's fail fails the leaf. Builders run one at a time: while
// none is live, a retrying leaf's next builder starts first, in leaf
// order, else the next builder starts on the first leaf in leaf order not
// yet started whose depends-on leaves have all passed. A built leaf's
// verifier runs
// alongside the next builder. A leaf whose dependency failed or was
// skipped is skipped. Once every leaf has an outcome the step ends: pass
// when every leaf passed, else fail.
func nextLeafActions(leaves []Leaf, states map[string]LeafState) []leafAction {
	known := map[string]bool{}
	for _, lf := range leaves {
		known[lf.ID] = true
	}
	var acts []leafAction
	var ready []Leaf
	var retries []leafAction
	builderLive := false
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
		if _, started := l.Session(workflow.RoleBuilder, max(1, l.Attempt())); started {
			a := nextForLeaf(lf, l)
			if a.kind == leafStart && a.role == workflow.RoleBuilder {
				retries = append(retries, a)
				continue
			}
			builderLive = builderLive || (a.kind == leafPoll && a.role == workflow.RoleBuilder)
			acts = append(acts, a)
			continue
		}
		if dep, ok := lostDependency(lf, states, known, leaves); ok {
			acts = append(acts, leafAction{kind: leafEnd, leaf: lf, outcome: LeafSkipped,
				detail: fmt.Sprintf("its dependency %s did not pass", dep)})
			continue
		}
		if dependenciesPassed(lf, states, known) {
			ready = append(ready, lf)
		}
	}
	switch {
	case builderLive:
	case len(retries) > 0:
		acts = append(acts, retries[0])
	case len(ready) > 0:
		acts = append(acts, leafAction{kind: leafStart, leaf: ready[0], role: workflow.RoleBuilder, attempt: 1})
	}
	if len(acts) > 0 || len(retries) > 0 {
		return acts
	}
	open := len(leaves) - passed - len(failed)
	if open > 0 {
		return []leafAction{{kind: leavesDone, outcome: workflow.OutcomeFail,
			detail: fmt.Sprintf("%d of %d leaves passed; %d cannot start", passed, len(leaves), open)}}
	}
	if len(failed) == 0 {
		return []leafAction{{kind: leavesDone, outcome: workflow.OutcomePass, detail: fmt.Sprintf("%d of %d leaves passed", passed, len(leaves))}}
	}
	return []leafAction{{kind: leavesDone, outcome: workflow.OutcomeFail,
		detail: fmt.Sprintf("%d of %d leaves passed; not passed: %s", passed, len(leaves), strings.Join(failed, ", "))}}
}

// dependenciesPassed reports whether every leaf lf depends on, among the
// step's leaves, has passed.
func dependenciesPassed(lf Leaf, states map[string]LeafState, known map[string]bool) bool {
	for _, d := range lf.DependsOn {
		if known[d] && states[d].Outcome != LeafPassed {
			return false
		}
	}
	return true
}

// lostDependency is the handle of the first leaf lf depends on that failed
// or was skipped, if any.
func lostDependency(lf Leaf, states map[string]LeafState, known map[string]bool, leaves []Leaf) (string, bool) {
	for _, d := range lf.DependsOn {
		if !known[d] {
			continue
		}
		if o := states[d].Outcome; o == LeafFailed || o == LeafSkipped {
			for _, l := range leaves {
				if l.ID == d {
					return l.Handle, true
				}
			}
		}
	}
	return "", false
}

// nextForLeaf is the next move on an open leaf l whose builder started. A
// verifier starts on the commit its builder's passed attempt recorded; a
// passed leaf ends with the Evidence commit its verifier's attempt
// recorded. A failed attempt before the round's MaxAttempts-th (a fresh
// round's first move is this, after the last round's failed attempt) starts a fresh builder on
// the next attempt; the caller decides when, as builders run one at a
// time. The last attempt's fail fails the leaf.
func nextForLeaf(lf Leaf, l LeafState) leafAction {
	attempt := max(1, l.Attempt())
	for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
		s, started := l.Session(role, attempt)
		if !started {
			a := leafAction{kind: leafStart, leaf: lf, role: role, attempt: attempt}
			if b, ok := l.Judged(workflow.RoleBuilder, attempt); ok && role == workflow.RoleVerifier {
				a.commit = b.Commit
			}
			return a
		}
		j, judged := l.Judged(role, attempt)
		if !judged {
			return leafAction{kind: leafPoll, leaf: lf, role: role, attempt: attempt, session: s.Session}
		}
		if j.Outcome != workflow.OutcomePass {
			if attempt-l.RoundBase() < MaxAttempts {
				return leafAction{kind: leafStart, leaf: lf, role: workflow.RoleBuilder, attempt: attempt + 1}
			}
			return leafAction{kind: leafEnd, leaf: lf, outcome: LeafFailed,
				detail: fmt.Sprintf("attempt %d: %s failed: %s", attempt, role, j.Detail)}
		}
	}
	v, _ := l.Judged(workflow.RoleVerifier, attempt)
	return leafAction{kind: leafEnd, leaf: lf, outcome: LeafPassed, commit: v.Commit,
		detail: fmt.Sprintf("attempt %d: verifier passed", attempt)}
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
// and sleeps the step's interval when a turn started, judged and ended
// nothing. Each turn polls every live session: the builder and any
// verifiers alongside it. A session in the
// journal is rejoined, never started again. A session's attempt is judged
// once its process exits (pass on exit 0) or its agent says done (pass);
// idle alone is never enough. Blocks, notices and an autonomous run's
// arbitration are handled as on an agent step. Runner failures are an
// error outcome; the error return is for journal failures.
//
// A leaf that fails raises a question (AddLeafQuestion) linked blocks to
// it; its dependents are skipped and the other leaves carry on, so the
// step ends fail once nothing more can proceed. On entry (a resume, or a
// retry, which keeps the step's leaves) a failed leaf whose question a
// Decision now answers gets a fresh round (reviveLeaves). In an autonomous
// run the arbiter answers a leaf's first failure's question at once.
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
	r.removeEndedCheckouts(lr, step)
	if res, ended, err := r.reviveLeaves(lr, step); err != nil || ended {
		return res, err
	}
	for {
		progressed := false
		for _, a := range nextLeafActions(leaves, r.leafStates(step.Name)) {
			switch a.kind {
			case leavesDone:
				return result{outcome: a.outcome, detail: a.detail}, nil
			case leafStart:
				if res, failed, err := r.startLeaf(lr, step, env, a); err != nil || failed {
					return res, err
				}
				progressed = true
			case leafEnd:
				if res, ended, err := r.endLeaf(lr, step, a); err != nil || ended {
					return res, err
				}
				progressed = true
			case leafPoll:
				judged, res, ended, err := r.pollLeaf(lr, step, a)
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

// startLeaf launches a's session and journals its leaf-session entry. Its
// model comes from the role's ladder at the attempt within the leaf's
// round, and a fresh round's builder is given the question the leaf's
// failure raised and the answering Decision, its statement read anew
// (LeafAnswer), so the journal holds no model text. A
// builder's with the worktree's HEAD before it started, which its leaf
// commit must follow, and a verifier's with the commit it checks and, as
// Detail, the checkout of it opened for the verifier to work in. A launch
// that fails ends the step (failed true) with an error outcome, its
// checkout removed.
func (r *running) startLeaf(lr LeafRunner, step workflow.Step, env []string, a leafAction) (result, bool, error) {
	l := r.leafStates(step.Name)[a.leaf.ID]
	spec := LeafSpec{
		Step: step, Role: a.role, Leaf: a.leaf, Attempt: a.attempt, Commit: a.commit,
		Round: l.Round(), RoundAttempt: a.attempt - l.RoundBase(),
		Env: append(append([]string{}, env...),
			EnvLeaf+"="+a.leaf.Handle, EnvAttempt+"="+strconv.Itoa(a.attempt), EnvRole+"="+a.role),
	}
	if a.commit != "" {
		spec.Env = append(spec.Env, EnvCommit+"="+a.commit)
	}
	if role := step.Role(a.role); role != nil {
		spec.Model = role.Model(spec.RoundAttempt)
	}
	if rd, ok := l.LastRound(); ok && a.role == workflow.RoleBuilder {
		spec.Question, spec.Decision = rd.Question, rd.Decision
		d, statement, ok, err := lr.LeafAnswer(rd.Question)
		if err != nil {
			return errResult(fmt.Errorf("leaf %s: read the answer to question %s: %v", a.leaf.Handle, rd.Question, err)), true, nil
		}
		if ok && d == rd.Decision {
			spec.Answer = statement
		}
	}
	commit := a.commit
	if a.role == workflow.RoleBuilder {
		head, err := lr.Head()
		if err != nil {
			return errResult(fmt.Errorf("leaf %s: read the worktree's HEAD: %v", a.leaf.Handle, err)), true, nil
		}
		commit = head
	}
	if a.role == workflow.RoleVerifier {
		dir, err := lr.OpenCheckout(a.leaf, a.attempt, a.commit)
		if err != nil {
			return errResult(fmt.Errorf("leaf %s: check out %s for its verifier: %v", a.leaf.Handle, a.commit, err)), true, nil
		}
		spec.Dir = dir
	}
	s, err := lr.StartLeaf(spec)
	if err != nil {
		if spec.Dir != "" {
			r.removeCheckout(lr, step, spec.Dir)
		}
		return errResult(fmt.Errorf("leaf %s: start its %s: %v", a.leaf.Handle, a.role, err)), true, nil
	}
	if err := r.append(Entry{
		Type: TypeLeafSession, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle,
		Role: a.role, Attempt: a.attempt, Session: s, Model: spec.Model, Commit: commit, Detail: spec.Dir,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: leaf %s attempt %d: %s session %s", step.Name, a.leaf.Handle, a.attempt, a.role, s)
	return result{}, false, nil
}

// pollLeaf polls a's session once. When it has exited or its agent said
// done, it journals the attempt's judgement (judged true), removes a
// verifier's checkout whatever the judgement, and reports the session's
// end to the Runner. A builder judged pass passes only with a
// new commit since its start naming its leaf (LeafCommit), journalled on
// the attempt; without one its attempt fails with LeafCommit's detail. A
// verifier judged pass passes only as judgeVerifier finds, its Evidence
// commit journalled on the attempt. A live session's block is journalled and
// notified as an agent step's, and in an autonomous run its question goes
// to the arbiter; a failed verdict ends the step (ended true) with res.
func (r *running) pollLeaf(lr LeafRunner, step workflow.Step, a leafAction) (judged bool, res result, ended bool, err error) {
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
		commit := ""
		if a.role == workflow.RoleVerifier && outcome == workflow.OutcomePass {
			if outcome, detail, commit, err = r.judgeVerifier(lr, step, a); err != nil {
				return false, result{}, false, err
			}
		}
		if a.role == workflow.RoleBuilder && outcome == workflow.OutcomePass {
			s, _ := r.leafStates(step.Name)[a.leaf.ID].Session(a.role, a.attempt)
			c, why, err := lr.LeafCommit(a.leaf, s.Commit)
			switch {
			case err != nil:
				outcome, detail = workflow.OutcomeFail, fmt.Sprintf("find its leaf commit: %v", err)
			case c == "":
				outcome, detail = workflow.OutcomeFail, why
			default:
				commit = c
			}
		}
		if err := r.append(Entry{
			Type: TypeLeafAttempt, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle,
			Role: a.role, Attempt: a.attempt, Outcome: outcome, Commit: commit, Detail: detail,
		}); err != nil {
			return false, result{}, false, err
		}
		r.log("step %s: leaf %s attempt %d: %s %s", step.Name, a.leaf.Handle, a.attempt, a.role, outcome)
		if s, ok := r.leafStates(step.Name)[a.leaf.ID].Session(a.role, a.attempt); ok && a.role == workflow.RoleVerifier && s.Detail != "" {
			r.removeCheckout(lr, step, s.Detail)
		}
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

// judgeVerifier judges a verifier session of a that ended well: it fails
// the attempt when the verifier journalled findings on it, else passes
// only when the leaf has Evidence the verifier's own actor wrote
// (Evidence), and that Evidence's entries are committed alone
// (CommitEvidence), whose commit it returns. A commit that fails fails
// the attempt. Findings win over Evidence. The verifier's `rhiza run
// findings` appended to the journal without the run's lock, so the
// journal is read again first. The error return is for journal failures.
func (r *running) judgeVerifier(lr LeafRunner, step workflow.Step, a leafAction) (outcome, detail, commit string, err error) {
	j, err := r.e.Store.Read(r.j.ID)
	if err != nil {
		return "", "", "", err
	}
	r.j = j
	if n := len(r.leafStates(step.Name)[a.leaf.ID].AttemptFindings(a.attempt)); n > 0 {
		return workflow.OutcomeFail, fmt.Sprintf("verifier found %d problem(s)", n), "", nil
	}
	ok, why, err := lr.Evidence(a.leaf, a.session)
	switch {
	case err != nil:
		return workflow.OutcomeFail, fmt.Sprintf("read its Evidence: %v", err), "", nil
	case ok:
		c, err := lr.CommitEvidence(a.leaf, a.session)
		if err != nil {
			return workflow.OutcomeFail, fmt.Sprintf("commit its Evidence: %v", err), "", nil
		}
		return workflow.OutcomePass, "Evidence written by the verifier", c, nil
	case why == "":
		return workflow.OutcomeFail, "verifier ended with neither Evidence nor findings", "", nil
	}
	return workflow.OutcomeFail, "verifier ended with no Evidence of its own and no findings: " + why, "", nil
}

// removeEndedCheckouts removes the checkout of every verifier of step
// whose attempt is judged: on a resume, one whose removal a killed run
// never reached. A live verifier's checkout stays.
func (r *running) removeEndedCheckouts(lr LeafRunner, step workflow.Step) {
	for _, l := range r.j.State().Leaves {
		if l.Step != step.Name {
			continue
		}
		for _, s := range l.Sessions {
			if s.Role != workflow.RoleVerifier || s.Detail == "" {
				continue
			}
			if _, judged := l.Judged(s.Role, s.Attempt); judged {
				r.removeCheckout(lr, step, s.Detail)
			}
		}
	}
}

// removeCheckout removes a verifier's checkout dir. A failure is logged,
// never fatal: a leftover checkout wastes space but changes no outcome.
func (r *running) removeCheckout(lr LeafRunner, step workflow.Step, dir string) {
	if err := lr.RemoveCheckout(dir); err != nil {
		r.log("step %s: remove checkout %s: %v", step.Name, dir, err)
	}
}

// endLeaf journals a's leaf outcome. A failed leaf first raises a question
// (AddLeafQuestion) naming its round's attempts and findings, linked
// blocks to it, journalled on the leaf entry; one that cannot be added is
// named in the entry's detail, and the run carries on. In an autonomous
// run, the question of a leaf's first round goes to the arbiter at once
// (reviveLeaf). An ended step (ended true) carries res.
func (r *running) endLeaf(lr LeafRunner, step workflow.Step, a leafAction) (res result, ended bool, err error) {
	en := Entry{Type: TypeLeaf, Step: step.Name, Leaf: a.leaf.ID, Handle: a.leaf.Handle, Outcome: a.outcome, Commit: a.commit, Detail: a.detail}
	if a.outcome == LeafFailed {
		l := r.leafStates(step.Name)[a.leaf.ID]
		q, err := lr.AddLeafQuestion(a.leaf, leafQuestionText(l))
		if err != nil {
			en.Detail += fmt.Sprintf("; its question could not be added: %v", err)
		} else {
			en.Question, en.QuestionHandle = q.ID, q.Handle
		}
	}
	if err := r.append(en); err != nil {
		return result{}, false, err
	}
	if en.Question != "" {
		r.log("step %s: leaf %s %s; question %s (%s)", step.Name, a.leaf.Handle, a.outcome, en.QuestionHandle, en.Question)
	} else {
		r.log("step %s: leaf %s %s", step.Name, a.leaf.Handle, a.outcome)
	}
	l := r.leafStates(step.Name)[a.leaf.ID]
	if en.Question != "" && r.j.State().Autonomous && l.Round() == 1 {
		return r.reviveLeaf(lr, step, l)
	}
	return result{}, false, nil
}

// reviveLeaves runs reviveLeaf on each of step's failed leaves, in the
// order the journal first names them.
func (r *running) reviveLeaves(lr LeafRunner, step workflow.Step) (result, bool, error) {
	var ids []string
	for _, l := range r.j.State().Leaves {
		if l.Step == step.Name && l.Outcome == LeafFailed && l.Question != "" {
			ids = append(ids, l.Leaf)
		}
	}
	for _, leaf := range ids {
		if res, ended, err := r.reviveLeaf(lr, step, r.leafStates(step.Name)[leaf]); err != nil || ended {
			return res, ended, err
		}
	}
	return result{}, false, nil
}

// reviveLeaf gives failed leaf l a fresh round (a leaf-round entry) once a
// Decision answers the question its failure raised. In an autonomous run
// the question of the leaf's first round is the arbiter's: it is
// arbitrated, or its open arbitration rejoined, and only a pass revives
// the leaf; any other verdict leaves it failed. So an autonomous run
// arbitrates a leaf once: a failure in its fresh round is never
// arbitrated again. Otherwise the owner's answer is read (LeafAnswer);
// with none the leaf stays failed. A LeafAnswer failure ends the step
// (ended true) with an error outcome. The error return is for journal
// failures.
func (r *running) reviveLeaf(lr LeafRunner, step workflow.Step, l LeafState) (result, bool, error) {
	if l.Outcome != LeafFailed || l.Question == "" {
		return result{}, false, nil
	}
	var decision string
	if r.j.State().Autonomous && l.Round() == 1 {
		a, ok := r.j.State().Arbitration(l.Question)
		if !ok || a.Open() {
			if _, err := r.arbitrate(step, l.Question); err != nil {
				return result{}, false, err
			}
			a, _ = r.j.State().Arbitration(l.Question)
		}
		if a.Outcome != workflow.OutcomePass {
			r.log("step %s: leaf %s stays failed: the arbiter did not answer question %s", step.Name, l.Handle, l.QuestionHandle)
			return result{}, false, nil
		}
		decision = a.Decision
	} else {
		d, _, ok, err := lr.LeafAnswer(l.Question)
		if err != nil {
			return errResult(fmt.Errorf("leaf %s: read the answer to question %s: %v", l.Handle, l.Question, err)), true, nil
		}
		if !ok {
			return result{}, false, nil
		}
		decision = d
	}
	if err := r.append(Entry{
		Type: TypeLeafRound, Step: step.Name, Leaf: l.Leaf, Handle: l.Handle,
		Question: l.Question, QuestionHandle: l.QuestionHandle, Decision: decision,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: leaf %s gets round %d: decision %s answers question %s", step.Name, l.Handle, l.Round()+1, decision, l.QuestionHandle)
	return result{}, false, nil
}

// leafQuestionText is the statement of the question failed leaf l raises:
// its handle, round, and each attempt of the round with the role that
// failed it, why, and its verifier's findings, each kept short.
func leafQuestionText(l LeafState) string {
	var parts []string
	for n := l.RoundBase() + 1; n <= l.Attempt(); n++ {
		part := fmt.Sprintf("attempt %d", n)
		for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
			if j, ok := l.Judged(role, n); ok && j.Outcome != workflow.OutcomePass {
				part += fmt.Sprintf(" %s failed: %s", role, clip(j.Detail))
			}
		}
		if fs := l.AttemptFindings(n); len(fs) > 0 {
			var rs []string
			for _, f := range fs {
				rs = append(rs, f.ID+" "+clip(f.Reason))
			}
			part += " (" + strings.Join(rs, "; ") + ")"
		}
		parts = append(parts, part)
	}
	return fmt.Sprintf("Leaf %s failed all %d attempts of round %d: %s. How should it be built, or should it change?",
		l.Handle, MaxAttempts, l.Round(), strings.Join(parts, "; "))
}

// clipRunes is how much of one detail or finding a leaf's question keeps.
const clipRunes = 160

// clip is s on one line, without a trailing full stop, cut to clipRunes.
func clip(s string) string {
	s = strings.TrimRight(strings.Join(strings.Fields(s), " "), ".")
	if r := []rune(s); len(r) > clipRunes {
		return string(r[:clipRunes-1]) + "…"
	}
	return s
}

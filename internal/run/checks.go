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
	// OpenCheckCheckout adds a throwaway detached worktree at commit,
	// where check's session in role on cycle works, and returns its path.
	// Its path is check's, role's and cycle's alone, apart from every
	// leaf verifier's checkout. It never touches the run's worktree; the
	// LeafRunner's RemoveCheckout removes it.
	OpenCheckCheckout(check, role string, cycle int, commit string) (dir string, err error)
	// StartFix launches the change set's fixer spec names, in the run's
	// worktree, and returns its session's ID.
	StartFix(spec FixSpec) (session string, err error)
	// ReviewCommit is the newest commit in since..HEAD of the run's
	// worktree whose Review: trailer names the run. With none, commit is
	// "" and detail says why.
	ReviewCommit(since string) (commit, detail string, err error)
	// AddReviewQuestion adds a question node, statement text, to the
	// run's change set, linked to no leaf, and returns it: the question a
	// capped review raises on a finding naming no leaf of the change set,
	// or on a check's failure with no finding. Once the change set has
	// landed, it goes in a record change set of its own instead; it
	// refuses one otherwise finished.
	AddReviewQuestion(text string) (Question, error)
	// CommitCheckWrites commits, on HEAD of the run's worktree, only the
	// log entry files that check or critic session session's actor wrote
	// there and no commit holds yet (a QA check's Evidence, a model
	// check's draft change set), as its own commit with no Leaf: trailer,
	// and returns it; "" when the actor wrote none. Other uncommitted
	// edits stay as they were.
	CommitCheckWrites(session string) (commit string, err error)
}

// ReviewAnswer is one answer an answered cycle's fixers are given: the
// Decision (handle) answering question node Question (handle
// QuestionHandle) that the capped cycle raised, and its statement, read
// anew from the model ("" when it no longer answers the question).
type ReviewAnswer struct {
	Question       string
	QuestionHandle string
	Decision       string
	Statement      string
}

// FixSpec is the change set's fixer to launch: role fixer of checks step
// Step in cycle Cycle (from 2), for Findings, the confirmed blocker and
// major findings of the cycle's checks that name no leaf of the change
// set, each with its check. Model is the fixer's model at Cycle, "" for
// the agent's default. Env holds the session's RHIZA_* variables, the
// step's plus RHIZA_CYCLE and RHIZA_ROLE (no RHIZA_LEAF, no RHIZA_CHECK).
// It works in the run's worktree, where its fix commit must carry the
// trailer `Review: <run>`. Anchor is the session whose pane its pane
// splits; "" opens the change set's tab anew. The Runner renders the
// fixer's prompt with {leaf} empty and names the findings.
type FixSpec struct {
	Step  workflow.Step
	Cycle int
	// Attempt is the change set's fixer's attempt in the cycle, from 1:
	// a failed fixer is retried within its cycle.
	Attempt  int
	Model    string
	Env      []string
	Anchor   string
	Findings []Finding
	// Answers are, in an answered cycle, the Decisions answering the
	// questions the capped cycle raised; nil otherwise.
	Answers []ReviewAnswer
	// Owner is true in an owner cycle: the findings are the owner's
	// request for changes, not ones a review cycle confirmed.
	Owner bool
}

// CheckSpec is one check session to launch: role Role (check) of checks
// step Step for Check, in cycle Cycle (from 1). Model is the model to
// run, "" for the agent's default, from the check's ladder at Cycle. Env
// holds the session's RHIZA_* variables, the step's plus RHIZA_CHECK,
// RHIZA_CYCLE and RHIZA_ROLE. Dir is where the session works, "" for the
// run's worktree: the engine gives each session a throwaway checkout of
// the run worktree's HEAD (OpenCheckCheckout). Anchor is the check session whose pane the session's
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
	// fixStart launches the fix phase's session fix names: a leaf's fixer
	// or its verifier, on the cycle.
	fixStart
	// fixPoll polls the fix phase's live session fix names.
	fixPoll
	// reviewFixStart launches the change set's fixer on the cycle, for
	// findings, on its attempt.
	reviewFixStart
	// reviewFixPoll polls the change set's live fixer session.
	reviewFixPoll
)

// checkAction is one move of the checks conductor, decided from the
// folded journal alone.
type checkAction struct {
	kind  checkActionKind
	check string
	role  string
	cycle int
	// attempt is a reviewFixStart's or reviewFixPoll's attempt of the
	// change set's fixer in the cycle, from 1.
	attempt int
	checks  []string
	session string
	outcome string
	detail  string
	// fix is a fixStart's or fixPoll's leaf session: role fixer or
	// verifier, its attempt the cycle.
	fix leafAction
	// findings are a reviewFixStart's findings: those naming no leaf.
	findings []Finding
}

// nextCheckActions is the pure decision of a checks step: from its
// checks, its cycles in order and its folded check states by name, the
// moves to make now. Cycle 1 runs every check; each later cycle only
// those the cycle before failed, up to MaxFixCycles further cycles. In a
// cycle every check's session starts before any is polled, so they all
// run in parallel: the moves list each start before any poll or end.
//
// After the cycle, once every check session of it is judged, each check
// whose session passed with a blocker or major finding in its round gets
// a critic session (needsCritic), all of them started before any is
// polled, so the critics run in parallel too. A check judged, and its
// critic's session if it needs one, is ended by checkVerdict. Once each
// check of the cycle has an outcome the next cycle begins over the failed
// ones; with none failed the step passes, and with no cycle left it
// fails.
//
// A cycle past the first opens with its fix phase (fixActions), before
// any of its checks starts: the confirmed findings of its checks that
// name a leaf (phaseFindings, by leaf through leafFindings) go to that
// leaf's fixer, then its verifier, leaves being the change set's in leaf
// order and fixes the step's folded leaf sessions by leaf ID. A fix phase that fails ends each check
// of the cycle failed ("cycle N: fix failed: ...") without a session, so
// the cycle is spent: the next one, if any is left, tries the fix again.
// A leaf's fixer or verifier judged fail, and a change-set fixer judged
// fail, are first retried within the cycle, up to MaxFixCycles more
// attempts (fixActions, changeSetFix), so only one failing all
// 1 + MaxFixCycles of them fails the phase.
// The confirmed findings naming no leaf of leaves (changeSetFindings) go
// to one fixer for the change set, after every leaf's fixer
// (changeSetFix); fixed holds the step's change-set fixers.
//
// An owner cycle (CheckCycle.Owner) fixes the owner's findings instead
// (phaseFindings), then runs the checks they name. Like an answered
// cycle it is a single cycle, outside the cap: a check failing it, or a
// failure of any cycle that is not within the cap, ends the step fail,
// which the conductor turns into questions (capped).
func nextCheckActions(step workflow.Step, cycles []CheckCycle, states map[string]CheckState, leaves []Leaf, fixes map[string]LeafState, fixed []ReviewFix) []checkAction {
	if len(cycles) == 0 {
		var all []string
		for _, c := range step.Checks {
			all = append(all, c.Name)
		}
		return []checkAction{{kind: checkCycle, cycle: 1, checks: all}}
	}
	cur := cycles[len(cycles)-1]
	if cur.Cycle > 1 {
		found := phaseFindings(cur, states)
		acts, failure := fixActions(cur.Cycle, leaves, leafFindings(found), fixes)
		acts, failure = changeSetFix(cur.Cycle, noLeafFindings(found, leaves), fixed, acts, failure)
		if len(acts) > 0 {
			return acts
		}
		if failure != "" {
			for _, name := range cur.Checks {
				if _, ok := states[name].Outcome(cur.Cycle); !ok {
					acts = append(acts, checkAction{kind: checkEnd, check: name, cycle: cur.Cycle, outcome: CheckFailed,
						detail: fmt.Sprintf("cycle %d: fix failed: %s", cur.Cycle, failure)})
				}
			}
			if len(acts) > 0 {
				return acts
			}
		}
	}
	var starts, rest []checkAction
	var failed, critics []string
	allJudged := true
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
			allJudged = false
			starts = append(starts, checkAction{kind: checkStart, check: name, role: workflow.RoleCheck, cycle: cur.Cycle})
			continue
		}
		if _, judged := c.Judgement(workflow.RoleCheck, cur.Cycle); !judged {
			allJudged = false
			rest = append(rest, checkAction{kind: checkPoll, check: name, role: workflow.RoleCheck, cycle: cur.Cycle, session: s.Session})
			continue
		}
		if needsCritic(c, cur.Cycle) {
			critics = append(critics, name)
			continue
		}
		outcome, detail, _ := checkVerdict(c, cur.Cycle)
		rest = append(rest, checkAction{kind: checkEnd, check: name, cycle: cur.Cycle, outcome: outcome, detail: detail})
	}
	if allJudged {
		for _, name := range critics {
			c := states[name]
			s, started := c.Session(workflow.RoleCritic, cur.Cycle)
			if !started {
				starts = append(starts, checkAction{kind: checkStart, check: name, role: workflow.RoleCritic, cycle: cur.Cycle})
				continue
			}
			if outcome, detail, ok := checkVerdict(c, cur.Cycle); ok {
				rest = append(rest, checkAction{kind: checkEnd, check: name, cycle: cur.Cycle, outcome: outcome, detail: detail})
				continue
			}
			rest = append(rest, checkAction{kind: checkPoll, check: name, role: workflow.RoleCritic, cycle: cur.Cycle, session: s.Session})
		}
	}
	if acts := append(starts, rest...); len(acts) > 0 {
		return acts
	}
	if len(failed) > 0 && cur.Cycle <= MaxFixCycles && cur.Kind == "" {
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

// phaseFindings are the findings cur's fix phase fixes, in order: an
// owner cycle's own (the owner's, confirmed by asking), else the
// confirmed findings of its checks (confirmedFindings).
func phaseFindings(cur CheckCycle, states map[string]CheckState) []Finding {
	if cur.Owner() {
		return cur.Findings
	}
	var out []Finding
	confirmedFindings(cur.Cycle, cur.Checks, states, func(f Finding) { out = append(out, f) })
	return out
}

// leafFindings are those of found that name a leaf, by leaf ID.
func leafFindings(found []Finding) map[string][]Finding {
	out := map[string][]Finding{}
	for _, f := range found {
		if f.Leaf != "" {
			out[f.Leaf] = append(out[f.Leaf], f)
		}
	}
	return out
}

// noLeafFindings are those of found, in order, that name no leaf, or one
// not among leaves: the change set's fixer's.
func noLeafFindings(found []Finding, leaves []Leaf) []Finding {
	var out []Finding
	for _, f := range found {
		if f.Leaf == "" || !slices.ContainsFunc(leaves, func(l Leaf) bool { return l.ID == f.Leaf }) {
			out = append(out, f)
		}
	}
	return out
}

// changeSetFindings are the findings cycle's change-set fixer fixes: of
// the confirmed findings phaseFindings reads, in check order, those that
// name no leaf, or name one not among leaves.
func changeSetFindings(cycle int, checks []string, states map[string]CheckState, leaves []Leaf) []Finding {
	return noLeafFindings(phaseFindings(CheckCycle{Cycle: cycle, Checks: checks}, states), leaves)
}

// confirmedFindings calls each with every confirmed finding cycle's fix
// phase fixes, its check's name set: for each of checks (those the cycle
// runs, each failed the cycle before), in order, the standing blocker and
// major findings of the round of the last earlier cycle the check's
// session ran in, provided its critic was judged pass there, so they are
// confirmed (R8). A cycle a failed fix spent ran no session, so the next
// one fixes the same findings again; a check whose session or critic
// failed has none to fix and only runs again. Those naming a leaf go to
// its fixer (leafFindings); the rest go to the change set's fixer
// (changeSetFindings), as do those naming a leaf the change set does not
// have.
func confirmedFindings(cycle int, checks []string, states map[string]CheckState, each func(Finding)) {
	if cycle < 2 {
		return
	}
	for _, name := range checks {
		c := states[name]
		last := 0
		for _, s := range c.Sessions {
			if s.Role == workflow.RoleCheck && s.Cycle < cycle {
				last = max(last, s.Cycle)
			}
		}
		if j, ok := c.Judgement(workflow.RoleCritic, last); !ok || j.Outcome != workflow.OutcomePass {
			continue
		}
		r, ok := c.RoundFor(last)
		if !ok {
			continue
		}
		for _, f := range r.Standing() {
			f.Check = name
			each(f)
		}
	}
}

// changeSetFix adds the change set's fixer of cycle to the leaves' fix
// phase, acts and failure as fixActions gave them, and returns the
// phase's moves and failure. With findings to fix it starts once every
// leaf's fixer was judged pass (no fixer move in acts, no failure), as
// fixers share the run's worktree, beside a fixed leaf's verifier if one
// runs; while live it is polled, even after a failure; judged fail it
// fails the phase ("change set: fixer failed: ..."), unless a leaf's
// failure came first. fixed are the step's change-set fixers.
func changeSetFix(cycle int, findings []Finding, fixed []ReviewFix, acts []checkAction, failure string) ([]checkAction, string) {
	if len(findings) == 0 {
		return acts, failure
	}
	var fix ReviewFix
	started := false
	for _, f := range fixed {
		if f.Cycle == cycle {
			fix, started = f, true
		}
	}
	attempt := 1
	if started && fix.Judged && fix.Outcome != workflow.OutcomePass && failure == "" && fix.Attempt <= MaxFixCycles {
		// A failed fixer is retried within its cycle, once any leaf
		// fixer's failure is out of the way.
		started, attempt = false, fix.Attempt+1
	}
	var mine checkAction
	switch {
	case !started:
		busy := slices.ContainsFunc(acts, func(a checkAction) bool { return a.fix.role == workflow.RoleFixer })
		if busy || failure != "" {
			return acts, failure
		}
		mine = checkAction{kind: reviewFixStart, cycle: cycle, attempt: attempt, role: workflow.RoleFixer, findings: findings}
	case !fix.Judged:
		mine = checkAction{kind: reviewFixPoll, cycle: cycle, attempt: fix.Attempt, role: workflow.RoleFixer, session: fix.Session}
	default:
		if fix.Outcome != workflow.OutcomePass && failure == "" {
			failure = "change set: fixer failed: " + fix.Detail
		}
		return acts, failure
	}
	// Starts before polls: the leaves' starts lead acts.
	i := slices.IndexFunc(acts, func(a checkAction) bool { return a.kind != fixStart })
	if i < 0 || mine.kind == reviewFixPoll {
		i = len(acts)
	}
	return slices.Insert(acts, i, mine), failure
}

// fixActions is the pure decision of cycle's fix phase, from the change
// set's leaves in leaf order, the findings to fix by leaf ID, and the
// step's folded leaf sessions by leaf ID: each leaf with findings gets a
// fixer, then, on the commit the fixer's judgement found, a verifier, both
// on attempt cycle. Fixers run one at a time, in leaf order, as they
// share the run's worktree; a fixed leaf's verifier runs alongside the
// next fixer. A leaf whose verifier failed waits for any other leaf's
// live fixer before its retry fixer starts. The moves list each start before any poll. failure says
// why the phase failed, once a fixer or verifier was judged fail: then
// nothing more starts, and the phase is over once its live sessions are
// judged. No moves and no failure: the phase is done, or had nothing to
// fix.
func fixActions(cycle int, leaves []Leaf, found map[string][]Finding, fixes map[string]LeafState) (acts []checkAction, failure string) {
	var starts, polls []checkAction
	// The fixer slot is taken by any leaf's live fixer, not only an
	// earlier leaf's: a retry can come up on a leaf the walk has passed.
	fixerFree := !slices.ContainsFunc(leaves, func(lf Leaf) bool {
		return len(found[lf.ID]) > 0 && liveFixer(cycle, fixes[lf.ID])
	})
	fail := func(lf Leaf, role, detail string) {
		if failure == "" {
			failure = fmt.Sprintf("leaf %s: %s failed: %s", lf.Handle, role, detail)
		}
	}
	for _, lf := range leaves {
		fs := found[lf.ID]
		if len(fs) == 0 {
			continue
		}
		l := fixes[lf.ID]
		// A failed fixer or verifier is retried as a fresh pair on the
		// next attempt of the cycle, up to MaxFixCycles more of them.
	attempts:
		for try := 1; try <= 1+MaxFixCycles; try++ {
			attempt := leafFixAttempt(cycle, try)
			retry := try <= MaxFixCycles
			s, started := l.Session(workflow.RoleFixer, attempt)
			if !started {
				if fixerFree && failure == "" {
					starts = append(starts, checkAction{kind: fixStart, cycle: cycle,
						fix: leafAction{kind: leafStart, leaf: lf, role: workflow.RoleFixer, attempt: attempt, cycle: cycle, findings: fs}})
				}
				fixerFree = false
				break
			}
			j, judged := l.Judged(workflow.RoleFixer, attempt)
			if !judged {
				polls = append(polls, checkAction{kind: fixPoll, cycle: cycle,
					fix: leafAction{kind: leafPoll, leaf: lf, role: workflow.RoleFixer, attempt: attempt, session: s.Session}})
				fixerFree = false
				break
			}
			if j.Outcome != workflow.OutcomePass {
				if !retry {
					fail(lf, workflow.RoleFixer, j.Detail)
				}
				continue attempts
			}
			v, started := l.Session(workflow.RoleVerifier, attempt)
			if !started {
				starts = append(starts, checkAction{kind: fixStart, cycle: cycle,
					fix: leafAction{kind: leafStart, leaf: lf, role: workflow.RoleVerifier, attempt: attempt, cycle: cycle, commit: j.Commit}})
				break
			}
			vj, judged := l.Judged(workflow.RoleVerifier, attempt)
			if !judged {
				polls = append(polls, checkAction{kind: fixPoll, cycle: cycle,
					fix: leafAction{kind: leafPoll, leaf: lf, role: workflow.RoleVerifier, attempt: attempt, session: v.Session}})
				break
			}
			if vj.Outcome != workflow.OutcomePass {
				if !retry {
					fail(lf, workflow.RoleVerifier, vj.Detail)
				}
				continue attempts
			}
			break
		}
	}
	if failure != "" {
		starts = nil
	}
	return append(starts, polls...), failure
}

// liveFixer says whether one of l's fixers of cycle is started and not
// yet judged.
func liveFixer(cycle int, l LeafState) bool {
	for try := 1; try <= 1+MaxFixCycles; try++ {
		attempt := leafFixAttempt(cycle, try)
		if _, started := l.Session(workflow.RoleFixer, attempt); started {
			if _, judged := l.Judged(workflow.RoleFixer, attempt); !judged {
				return true
			}
		}
	}
	return false
}

// leafFixAttempt is the attempt number, unique among the leaf's sessions
// of the step, of try (from 1) of its fix in cycle (from 2): cycle 2's
// first try is attempt 2, and each cycle gets 1 + MaxFixCycles of them.
func leafFixAttempt(cycle, try int) int {
	return (cycle-2)*(1+MaxFixCycles) + try + 1
}

// needsCritic says whether the check's session of cycle passed having
// recorded, in its round for the cycle, a blocker or major finding: only
// then does a critic weigh in. Minor and nit findings never call one.
func needsCritic(c CheckState, cycle int) bool {
	j, judged := c.Judgement(workflow.RoleCheck, cycle)
	if !judged || j.Outcome != workflow.OutcomePass {
		return false
	}
	r, ok := c.RoundFor(cycle)
	return ok && len(r.Serious()) > 0
}

// checkVerdict is the check's outcome in cycle once its session's work,
// and its critic's if it needs one, is judged. A session that failed
// fails its check; one that passed with no round for the cycle fails it
// too ("no findings round"), so it gets no critic and re-runs next cycle.
// A round with no blocker or major finding passes, with no critic. Else
// the critic decides: a critic session that failed fails the check
// ("critic failed"), and otherwise the check fails iff a serious finding
// stands (Round.Standing: confirmed, or given no verdict). A critic can
// only confirm or drop the round's findings, never add one. ok is false
// while a session the outcome waits on is unjudged.
func checkVerdict(c CheckState, cycle int) (outcome, detail string, ok bool) {
	j, judged := c.Judgement(workflow.RoleCheck, cycle)
	if !judged {
		return "", "", false
	}
	if j.Outcome != workflow.OutcomePass {
		return CheckFailed, fmt.Sprintf("cycle %d: check failed: %s", cycle, j.Detail), true
	}
	r, ok := c.RoundFor(cycle)
	if !ok {
		return CheckFailed, fmt.Sprintf("cycle %d: no findings round", cycle), true
	}
	serious := r.Serious()
	if len(serious) == 0 {
		return CheckPassed, fmt.Sprintf("cycle %d: check passed", cycle), true
	}
	cj, judged := c.Judgement(workflow.RoleCritic, cycle)
	if !judged {
		return "", "", false
	}
	if cj.Outcome != workflow.OutcomePass {
		return CheckFailed, fmt.Sprintf("cycle %d: critic failed: %s", cycle, cj.Detail), true
	}
	if standing := r.Standing(); len(standing) > 0 {
		var ids []string
		for _, f := range standing {
			ids = append(ids, f.ID+" ("+f.Severity+")")
		}
		return CheckFailed, fmt.Sprintf("cycle %d: confirmed %s", cycle, strings.Join(ids, ", ")), true
	}
	var ids []string
	for _, f := range serious {
		ids = append(ids, f.ID)
	}
	return CheckPassed, fmt.Sprintf("cycle %d: critic dropped %s", cycle, strings.Join(ids, ", ")), true
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
	r.removeCheckCheckouts(cr, step, false)
	r.removeEndedCheckouts(cr, step)
	res, err := r.conductChecks(cr, step, env)
	if err != nil {
		return res, err
	}
	if err := r.closeEndedCheckPanes(cr, step); err != nil {
		return result{}, err
	}
	if err := r.closeEndedPanes(cr, step); err != nil {
		return result{}, err
	}
	if err := r.closeEndedFixPanes(cr, step); err != nil {
		return result{}, err
	}
	r.removeCheckCheckouts(cr, step, true)
	r.removeExitedCheckouts(cr, step)
	return res, nil
}

// openCheckCheckout opens the throwaway checkout a's session works in: a
// detached worktree at the run worktree's HEAD, named for a's check, role
// and cycle. Check and critic sessions alike work in one.
func (r *running) openCheckCheckout(cr CheckRunner, a checkAction) (string, error) {
	head, err := cr.Head()
	if err != nil {
		return "", fmt.Errorf("read the worktree's HEAD: %v", err)
	}
	dir, err := cr.OpenCheckCheckout(a.check, a.role, a.cycle, head)
	if err != nil {
		return "", fmt.Errorf("check out %s for its %s: %v", head, a.role, err)
	}
	return dir, nil
}

// removeCheckCheckout removes the checkout of check's session in role on
// cycle, if the journal names one. A failure is logged, never fatal.
func (r *running) removeCheckCheckout(cr CheckRunner, step workflow.Step, check, role string, cycle int) {
	if s, ok := r.checkStates(step.Name)[check].Session(role, cycle); ok && s.Detail != "" {
		r.removeCheckout(cr, step, s.Detail)
	}
}

// removeCheckCheckouts removes the checkout of every judged session of
// step, whatever its role: on a resume (atEnd false), one whose removal a
// killed run never reached. As the step ends (atEnd true) it removes
// instead that of each session not yet judged whose process has exited,
// one the step ended on before judging it. A live session keeps its
// checkout.
func (r *running) removeCheckCheckouts(cr CheckRunner, step workflow.Step, atEnd bool) {
	for _, c := range r.j.State().Checks {
		if c.Step != step.Name {
			continue
		}
		for _, s := range c.Sessions {
			if s.Detail == "" {
				continue
			}
			_, judged := c.Judgement(s.Role, s.Cycle)
			switch {
			case judged && !atEnd:
				r.removeCheckout(cr, step, s.Detail)
			case !judged && atEnd:
				if st, err := r.e.Runner.AgentStatus(s.Session); err == nil && st.Done {
					r.removeCheckout(cr, step, s.Detail)
				}
			}
		}
	}
}

// conductChecks is the checks step's loop, as checks describes.
func (r *running) conductChecks(cr CheckRunner, step workflow.Step, env []string) (result, error) {
	interval := step.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	if err := r.ownerCycle(step); err != nil {
		return result{}, err
	}
	var leaves []Leaf
	for {
		cycles, states := r.stepCycles(step.Name), r.checkStates(step.Name)
		if leaves == nil && needsLeaves(cycles, states) {
			l, err := cr.Leaves(r.j.State().Change)
			if err != nil {
				return errResult(fmt.Errorf("read the change set's leaves for its fixers: %v", err)), nil
			}
			leaves = l
			if r.leafList == nil {
				r.leafList = map[string][]Leaf{}
			}
			r.leafList[step.Name] = l
		}
		progressed := false
		for _, a := range nextCheckActions(step, cycles, states, leaves, r.leafStates(step.Name), r.reviewFixes(step.Name)) {
			switch a.kind {
			case checksDone:
				if a.outcome != workflow.OutcomeFail {
					return result{outcome: a.outcome, detail: a.detail}, nil
				}
				res, again, err := r.capped(cr, step, a, &leaves)
				if err != nil || !again {
					return res, err
				}
				progressed = true
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
				judged, res, ended, err := r.pollCheck(cr, step, a)
				if err != nil || ended {
					return res, err
				}
				progressed = progressed || judged
			case fixStart:
				if res, failed, err := r.startLeaf(cr, step, env, a.fix); err != nil || failed {
					return res, err
				}
				progressed = true
			case fixPoll:
				judged, res, ended, err := r.pollLeaf(cr, step, a.fix)
				if err != nil || ended {
					return res, err
				}
				progressed = progressed || judged
			case reviewFixStart:
				if res, failed, err := r.startReviewFix(cr, step, env, a); err != nil || failed {
					return res, err
				}
				progressed = true
			case reviewFixPoll:
				judged, res, ended, err := r.pollReviewFix(cr, step, a)
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

// needsLeaves says whether the latest of a checks step's cycles has
// findings for leaf fixers, so the conductor needs the change set's
// leaves: read only then, as a review without them never needs them.
func needsLeaves(cycles []CheckCycle, states map[string]CheckState) bool {
	if len(cycles) == 0 {
		return false
	}
	cur := cycles[len(cycles)-1]
	return cur.Cycle > 1 && len(leafFindings(phaseFindings(cur, states))) > 0
}

// reviewFixes are step's change-set fixers, in order.
func (r *running) reviewFixes(step string) []ReviewFix {
	var out []ReviewFix
	for _, f := range r.j.State().Fixes {
		if f.Step == step {
			out = append(out, f)
		}
	}
	return out
}

// startReviewFix launches the change set's fixer a names, on the fixer's
// model for the cycle, in the run's worktree, its pane where stepPane
// says, and journals its review-fix entry with the worktree's HEAD before
// it started, which its Review: commit must follow. A launch that fails
// ends the step (failed true) with an error outcome.
func (r *running) startReviewFix(cr CheckRunner, step workflow.Step, env []string, a checkAction) (result, bool, error) {
	role := step.Role(workflow.RoleFixer)
	if role == nil {
		return errResult(fmt.Errorf("step %q has no fixer", step.Name)), true, nil
	}
	head, err := cr.Head()
	if err != nil {
		return errResult(fmt.Errorf("change set's fixer: read the worktree's HEAD: %v", err)), true, nil
	}
	spec := FixSpec{
		Step: step, Cycle: a.cycle, Attempt: a.attempt, Model: role.Model(a.cycle), Findings: a.findings,
		Env: append(append([]string{}, env...),
			EnvCycle+"="+strconv.Itoa(a.cycle), EnvRole+"="+workflow.RoleFixer),
		Anchor: stepPane(r.checkStates(step.Name), r.leafStates(step.Name), r.reviewFixes(step.Name)),
	}
	answers, err := r.cycleAnswers(cr, step, a.cycle)
	if err != nil {
		return errResult(fmt.Errorf("change set's fixer: %v", err)), true, nil
	}
	spec.Answers = answers
	spec.Owner = r.j.State().OwnerCycle(step.Name, a.cycle)
	s, err := cr.StartFix(spec)
	if err != nil {
		return errResult(fmt.Errorf("start the change set's fixer: %v", err)), true, nil
	}
	if err := r.append(Entry{
		Type: TypeReviewFix, Step: step.Name, Role: workflow.RoleFixer, Cycle: a.cycle, Attempt: a.attempt, Session: s, Model: spec.Model, Commit: head,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: cycle %d: change set's fixer session %s", step.Name, a.cycle, s)
	return result{}, false, nil
}

// pollReviewFix polls the change set's fixer a names once. When it has
// exited or its agent said done, it journals its judgement (judged true)
// and reports its end to the Runner: judged pass, it passes only with a
// new commit since its start whose Review: trailer names the run
// (ReviewCommit), journalled on the judgement; without one it fails with
// ReviewCommit's detail. A live session's block is journalled and
// notified as an agent step's, and in an autonomous run its question goes
// to the arbiter; a failed verdict ends the step (ended true) with res.
func (r *running) pollReviewFix(cr CheckRunner, step workflow.Step, a checkAction) (judged bool, res result, ended bool, err error) {
	status, err := r.e.Runner.AgentStatus(a.session)
	if err != nil {
		if err := r.unblock(step.Name, a.session); err != nil {
			return false, result{}, false, err
		}
		return false, errResult(fmt.Errorf("change set's fixer session %s: %v", a.session, err)), true, nil
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
		if outcome == workflow.OutcomePass {
			var since string
			for _, f := range r.reviewFixes(step.Name) {
				if f.Cycle == a.cycle && f.Attempt == a.attempt {
					since = f.Since
				}
			}
			c, why, err := cr.ReviewCommit(since)
			switch {
			case err != nil:
				outcome, detail = workflow.OutcomeFail, fmt.Sprintf("find its Review: commit: %v", err)
			case c == "":
				outcome, detail = workflow.OutcomeFail, why
			default:
				commit = c
			}
		}
		if err := r.append(Entry{
			Type: TypeReviewFixJudged, Step: step.Name, Role: workflow.RoleFixer, Cycle: a.cycle, Attempt: a.attempt, Outcome: outcome, Commit: commit, Detail: detail,
		}); err != nil {
			return false, result{}, false, err
		}
		r.log("step %s: cycle %d: change set's fixer %s", step.Name, a.cycle, outcome)
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

// closeEndedFixPanes closes the pane of every judged change-set fixer of
// step still open, as closeEndedCheckPanes does a check session's: a
// failure is logged and leaves the pane open in the journal. The error
// return is for journal failures.
func (r *running) closeEndedFixPanes(cr CheckRunner, step workflow.Step) error {
	for _, f := range r.reviewFixes(step.Name) {
		if !f.Judged || f.PaneClosed {
			continue
		}
		if err := cr.ClosePane(f.Session); err != nil {
			r.log("step %s: close the pane of session %s: %v", step.Name, f.Session, err)
			continue
		}
		if err := r.append(Entry{Type: TypePaneClosed, Step: step.Name, Role: workflow.RoleFixer, Cycle: f.Cycle, Attempt: f.Attempt, Session: f.Session}); err != nil {
			return err
		}
	}
	return nil
}

// startCheck launches a's session, on the check's model for the cycle, in
// a pane checkPane places and in a checkout openCheckCheckout opens, and
// journals its check-session entry, the checkout as its detail. A launch
// that fails ends the step (failed true) with an error outcome.
func (r *running) startCheck(cr CheckRunner, step workflow.Step, env []string, a checkAction) (result, bool, error) {
	check, ok := step.Check(a.check)
	if !ok {
		return errResult(fmt.Errorf("check %s is not in step %q", a.check, step.Name)), true, nil
	}
	model := check.Model(a.cycle)
	if a.role != workflow.RoleCheck {
		role := step.Role(a.role)
		if role == nil {
			return errResult(fmt.Errorf("check %s: step %q has no %s", a.check, step.Name, a.role)), true, nil
		}
		model = role.Model(a.cycle)
	}
	dir, err := r.openCheckCheckout(cr, a)
	if err != nil {
		return errResult(fmt.Errorf("check %s: %v", a.check, err)), true, nil
	}
	spec := CheckSpec{
		Step: step, Role: a.role, Check: check, Cycle: a.cycle, Model: model, Dir: dir,
		Env: append(append([]string{}, env...),
			EnvCheck+"="+a.check, EnvCycle+"="+strconv.Itoa(a.cycle), EnvRole+"="+a.role),
		Anchor: stepPane(r.checkStates(step.Name), r.leafStates(step.Name), r.reviewFixes(step.Name)),
	}
	s, err := cr.StartCheck(spec)
	if err != nil {
		r.removeCheckout(cr, step, dir)
		return errResult(fmt.Errorf("check %s: start its %s: %v", a.check, sessionWord(a.role), err)), true, nil
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

// sessionWord names a session of role in an error: "session" for a
// check's own, "critic session" for its critic's.
func sessionWord(role string) string {
	if role == workflow.RoleCheck {
		return "session"
	}
	return role + " session"
}

// pollCheck polls a's session once. When it has exited or its agent said
// done, it reads the journal again (the session's `rhiza run findings`
// appended to it without the run's lock), journals the session's
// judgement (judged true), removes its checkout whatever the judgement,
// and reports its end to the Runner. A live session's block is journalled and
// notified as an agent step's, and in an autonomous run its question goes
// to the arbiter; a failed verdict ends the step (ended true) with res.
func (r *running) pollCheck(cr CheckRunner, step workflow.Step, a checkAction) (judged bool, res result, ended bool, err error) {
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
		j, err := r.e.Store.Read(r.j.ID)
		if err != nil {
			return false, result{}, false, err
		}
		r.j = j
		if err := r.append(Entry{
			Type: TypeCheckJudged, Step: step.Name, Check: a.check, Role: a.role, Cycle: a.cycle, Outcome: outcome, Detail: detail,
		}); err != nil {
			return false, result{}, false, err
		}
		r.log("step %s: check %s cycle %d: %s %s", step.Name, a.check, a.cycle, a.role, outcome)
		// The session's model writes (QA Evidence, a draft change set)
		// are committed now, so no later fixer sweeps them into its
		// commit. A failure is only logged: the judgement stands, and the
		// writes stay uncommitted in the run's worktree.
		if c, err := cr.CommitCheckWrites(a.session); err != nil {
			r.log("step %s: check %s cycle %d: commit the %s's writes: %v", step.Name, a.check, a.cycle, a.role, err)
		} else if c != "" {
			r.log("step %s: check %s cycle %d: committed the %s's writes as %s", step.Name, a.check, a.cycle, a.role, c)
		}
		r.removeCheckCheckout(cr, step, a.check, a.role, a.cycle)
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

// stepPane is the pure choice of where a checks step's new session's
// pane opens, a check's or a fixer's or verifier's alike: the most recent
// open pane of the step's check sessions (checkPane), its leaf sessions
// and its change-set fixers, or "" for a new tab when none is open.
func stepPane(states map[string]CheckState, fixes map[string]LeafState, fixed []ReviewFix) string {
	anchor, seq := "", -1
	if open := openCheckPanes(states); len(open) > 0 {
		p := open[len(open)-1]
		anchor, seq = p.s.Session, p.s.Seq
	}
	if open := openPanes(fixes); len(open) > 0 {
		if p := open[len(open)-1]; p.s.Seq > seq {
			anchor, seq = p.s.Session, p.s.Seq
		}
	}
	for _, f := range fixed {
		if !f.PaneClosed && f.Seq > seq {
			anchor, seq = f.Session, f.Seq
		}
	}
	return anchor
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

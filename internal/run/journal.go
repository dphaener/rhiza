// Package run keeps the journal of a workflow run: an append-only JSONL file
// per run, folded to a state. It is work-layer state and lives outside the
// model's log. It runs no program and names no agent.
package run

import (
	"slices"
	"strings"
	"time"
)

// Entry types.
const (
	TypeStarted    = "started"
	TypeEntered    = "entered"
	TypeSession    = "session"
	TypeEnded      = "ended"
	TypeStopped    = "stopped"
	TypeFinished   = "finished"
	TypeBlocked    = "blocked"
	TypeUnblocked  = "unblocked"
	TypeAmended    = "amended"
	TypeRetried    = "retried"
	TypeAccepted   = "accepted"
	TypeRejected   = "rejected"
	TypeArbiter    = "arbiter"
	TypeArbitrated = "arbitrated"
	TypeAsked      = "asked"
	// TypeAdvisor and TypeRecommended are a normal run's advice on a
	// question: an advisor session was launched, and it recommended an
	// answer.
	TypeAdvisor     = "advisor"
	TypeRecommended = "recommended"
	TypeAnswered    = "answered"
	TypeMessage     = "message"
	TypeTokens      = "tokens"
	// TypeFreshSession records that a session started over as a fresh
	// session after a usage-limit pause, its transcript gone.
	TypeFreshSession = "fresh-session"
	TypeFindings     = "findings"
	TypeVerdict      = "verdict"
	TypeGated        = "gated"
	// TypeLeafSession, TypeLeafAttempt and TypeLeaf are a leaves step's
	// entries: a leaf's session started, a session's work judged, and a
	// leaf's outcome.
	TypeLeafSession = "leaf-session"
	TypeLeafAttempt = "leaf-attempt"
	TypeLeaf        = "leaf"
	// TypeLeafRound gives a failed leaf a fresh round of attempts, once a
	// Decision answers the question its failure raised.
	TypeLeafRound = "leaf-round"
	// TypePaneClosed records that a leaf or check session's pane was
	// closed.
	TypePaneClosed = "pane-closed"
	// TypeCheckCycle, TypeCheckSession, TypeCheckJudged and TypeCheck are
	// a checks step's entries: a cycle of its checks began, a check's
	// session started, a session's work judged, and a check's outcome in
	// a cycle.
	TypeCheckCycle   = "check-cycle"
	TypeCheckSession = "check-session"
	TypeCheckJudged  = "check-judged"
	TypeCheck        = "check"
	// TypeReviewFix and TypeReviewFixJudged are a checks step's change-set
	// fixer's entries: its session started in a cycle (Commit the
	// worktree's HEAD then), and its work judged (Commit, when it passed,
	// the new commit whose Review: trailer names the run). Neither names
	// a leaf or a check.
	TypeReviewFix       = "review-fix"
	TypeReviewFixJudged = "review-fix-judged"
	// TypeReviewQuestion is a question a checks step raised when a check
	// still failed with no fix cycle left: one per surviving confirmed
	// finding of the check, or one for the check's failure when it had
	// none.
	TypeReviewQuestion = "review-question"
	// TypeCodeCheck is a review's code check: rhiza run code-check ran
	// rhiza check and build-check on the run's change set.
	TypeCodeCheck = "code-check"
	// TypeSignoff is a signoff step's sign-off: the owner approved the
	// review or asked for changes (rule owner), or an autonomous run
	// passed it in the owner's place (rule autonomous).
	TypeSignoff = "signoff"
	// TypeProcess records the process that runs the run, each time it
	// starts or is resumed.
	TypeProcess = "process"
	// TypeDaemonResumed records that the daemon resumed the run, which it
	// does once.
	TypeDaemonResumed = "daemon-resumed"
)

// Sign-off outcomes, a signoff entry's Outcome.
const (
	SignoffApproved = "approved"
	SignoffChanges  = "changes"
)

// Sign-off rules, a signoff entry's Rule: who signed off.
const (
	RuleOwner      = "owner"
	RuleAutonomous = "autonomous"
)

// CycleOwner is the Kind of a check-cycle entry that runs the checks the
// owner's request for changes names, after a fix phase for the owner's
// findings: one cycle outside the cap, granted by the sign-off.
const CycleOwner = "owner"

// Code check outcomes, a code-check entry's Outcome.
const (
	CodeCheckPass = "pass"
	CodeCheckFail = "fail"
)

// CycleAnswered is the Kind of a check-cycle entry that gives a capped
// checks step one fresh cycle, once Decisions answer every question its
// last cycle raised.
const CycleAnswered = "answered"

// Check outcomes, a check entry's Outcome.
const (
	CheckPassed = "passed"
	CheckFailed = "failed"
)

// Leaf outcomes, a leaf entry's Outcome.
const (
	LeafPassed  = "passed"
	LeafFailed  = "failed"
	LeafSkipped = "skipped"
)

// Finding kinds: what a plan reviewer's finding says is wrong with a node.
const (
	KindMissingAcceptance = "missing-acceptance"
	KindMissingProvenance = "missing-provenance"
	KindUnrealDependsOn   = "unreal-depends-on"
	KindOutsideIntent     = "outside-intent"
	KindContradiction     = "contradiction"
)

// Kinds are the finding kinds, in the order they are documented.
var Kinds = []string{KindMissingAcceptance, KindMissingProvenance, KindUnrealDependsOn, KindOutsideIntent, KindContradiction}

// ValidKind reports whether kind is one of Kinds.
func ValidKind(kind string) bool {
	for _, k := range Kinds {
		if k == kind {
			return true
		}
	}
	return false
}

// Severities of a review finding, from the most to the least serious.
const (
	SeverityBlocker = "blocker"
	SeverityMajor   = "major"
	SeverityMinor   = "minor"
	SeverityNit     = "nit"
)

// Severities are a review finding's severities, most serious first.
var Severities = []string{SeverityBlocker, SeverityMajor, SeverityMinor, SeverityNit}

// ValidSeverity reports whether sev is one of Severities.
func ValidSeverity(sev string) bool {
	for _, s := range Severities {
		if s == sev {
			return true
		}
	}
	return false
}

// Serious reports whether a finding of severity sev is serious: blocker or
// major, the findings a critic judges.
func Serious(sev string) bool { return sev == SeverityBlocker || sev == SeverityMajor }

// Verdicts a critic gives a finding.
const (
	Confirmed = "confirmed"
	Dropped   = "dropped"
)

// Run results, used in Entry.End and State.Result.
const (
	Succeeded = "succeeded"
	Failed    = "failed"
)

// Started is what a run needs to be resumed: the definition as it was
// when the run began, and where and how it runs.
type Started struct {
	Workflow     string `json:"workflow"`
	WorkflowPath string `json:"workflow_path,omitempty"`
	Definition   string `json:"definition"`
	// Includes holds the source of each workflow the definition's
	// workflow steps include, by name, as the run started with them.
	Includes       map[string]string `json:"includes,omitempty"`
	Change         string            `json:"change"`
	Title          string            `json:"title,omitempty"`
	Worktree       string            `json:"worktree,omitempty"`
	Branch         string            `json:"branch,omitempty"`
	PermissionMode string            `json:"permission_mode,omitempty"`
	Agent          string            `json:"agent,omitempty"`
	// Autonomous is true for a run started without the owner: it skips
	// the owner's gates and questions, recording each. It never changes
	// an agent's permissions.
	Autonomous bool `json:"autonomous,omitempty"`
	// Milestone is the ID of the milestone intent a plan run's request
	// names, handed to the planner as RHIZA_MILESTONE; empty for none.
	Milestone string `json:"milestone,omitempty"`
	// MilestonePlan is true for the plan run of a milestone request: its
	// interview session waits on the owner from launch.
	MilestonePlan bool `json:"milestone_plan,omitempty"`
}

// Block is what a blocked agent step waits on: the kind of block and the
// question, as the session's blocked report gave them, and the tab its
// session runs in. The tab is plain strings so this package stays free of
// the session records it comes from.
type Block struct {
	Kind       string `json:"kind,omitempty"`
	Question   string `json:"question,omitempty"`
	TabSession string `json:"tab_session,omitempty"`
	TabWindow  string `json:"tab_window,omitempty"`
	TabBlock   string `json:"tab_block,omitempty"`
	TabTitle   string `json:"tab_title,omitempty"`
}

// Entry is one line of a journal.
//
//	started:   Step is the start step, Started holds the run's setup
//	entered:   Step, Visit
//	session:   Step, Session (an agent step's launched session)
//	blocked:   Step, Session, Block (the session began waiting on the owner)
//	unblocked: Step, Session (the block ended: answered, or the session ended)
//	ended:     Step, Outcome, and Next (a step) or End (the run ends)
//	stopped:   Step, Detail (why it stopped, e.g. waiting)
//	finished:  End (succeeded or failed), Detail
//	amended:   Step, Session, Actor, File (the owner or an agent amended the step's brief)
//	retried:   Step, Detail (a failed run reopened at Step; Detail is why it had failed)
//	accepted:  Change, LogEntry (the run accepted a change set under rule autonomous)
//	rejected:  Change, LogEntry (the run rejected a change set under rule autonomous)
//	arbiter:   Step, Session, Question (the arbiter's session launched on question
//	           node Question, asked from agent step Step)
//	arbitrated: Step, Question, Outcome (pass, fail or error), Detail, Decision (the
//	           arbiter's verdict on Question; Decision is the answering Decision's handle)
//	asked:     Step, Session, Question, Handle, Detail (an autonomous run's agent step's
//	           session asked the owner Detail; the engine added question node Question,
//	           handle Handle, for the arbiter)
//	advisor:   Step, Session, Asker (the advisor's session launched on the question
//	           Asker, a session of agent step Step, blocked on the owner)
//	recommended: Step, Session, Asker, Text, Detail (advisor session Session
//	           recommended Text to the question Asker asked, for the reasons
//	           in Detail, restating it as Restated when it has; it answers
//	           nothing)
//	message:   Step, Session, Actor (the sender), To, Text, Summary, Structured (a session sent another a message)
//	tokens:    Step, Session, Tokens (a reading of the session's token use; the latest counts)
//	fresh-session: Step, Session, Detail (the session's transcript was gone
//	           after a usage-limit pause, so its step started over as a
//	           fresh session instead of resuming)
//	findings:  Findings (a plan reviewer's findings, maybe none; each findings entry
//	           opens the next review round). With Leaf set it is instead a
//	           verifier's findings on leaf Leaf (handle Handle) of leaves step
//	           Step, for its attempt Attempt; plan review rounds skip it
//	verdict:   Verdicts (a critic's verdicts on findings of the latest round)
//	           With Check set, findings and verdict entries are instead a
//	           review check's: findings opens the next round of check Check
//	           of checks step Step, in cycle Cycle (when 0, the cycle of the
//	           check's latest session), and verdict judges that check's
//	           latest round; plan review rounds skip both
//	gated:     Changes (run gate ended the plan: its change sets, by ID; the
//	           latest gated entry wins)
//	leaf-session: Step, Leaf, Handle, Role (builder or verifier), Attempt, Session,
//	           Model, Commit (a builder's: the worktree's HEAD when it started; a
//	           verifier's: the commit it checks), Detail (a leaves step started a
//	           session on leaf node Leaf, handle Handle); on a checks step, Cycle
//	           is the review cycle of the fix it belongs to
//	leaf-attempt: Step, Leaf, Handle, Role, Attempt, Outcome (pass or fail), Commit
//	           (a passed builder's leaf commit; a passed verifier's Evidence
//	           commit), Detail (the role's session of the attempt was judged)
//	leaf:      Step, Leaf, Handle, Outcome (passed, failed or skipped), Question and
//	           QuestionHandle (for failed: the question node its failure raised),
//	           Commit (a passed leaf's Evidence commit), Detail (the leaf's outcome)
//	leaf-round: Step, Leaf, Handle, Question, QuestionHandle, Decision (failed leaf
//	           Leaf gets a fresh round of MaxAttempts attempts: Decision, by
//	           handle, answers Question, which its failure raised; the leaf is
//	           open again, and the step's skipped leaves with it)
//	pane-closed: Step, Leaf, Handle, Session (the pane of leaf session Session was
//	           closed: its session passed, the next builder took its place, or
//	           its leaves step ended). With Check set instead of Leaf, the pane
//	           of check Check's session Session, closed as its checks step ended
//	check-cycle: Step, Cycle, Checks (checks step Step began cycle Cycle, from
//	           1, over the checks named: every check in cycle 1, then those that
//	           failed the cycle before)
//	check-session: Step, Check, Role (check), Cycle, Session, Model, Detail (a
//	           checks step started a session of role Role for check Check in
//	           cycle Cycle)
//	check-judged: Step, Check, Role, Cycle, Outcome (pass or fail), Detail (the
//	           role's session of the check's cycle was judged)
//	check:     Step, Check, Cycle, Outcome (passed or failed), Detail (the check's
//	           outcome in the cycle)
//	           A check-cycle with Kind answered also has Questions and
//	           Decisions: the questions the cycle before raised, by ID, and
//	           the Decisions answering them, by handle, in the same order;
//	           it is one fresh cycle past the cap
//	review-question: Step, Check, Cycle, Finding, Leaf, Handle, Question,
//	           QuestionHandle (check Check still failed cycle Cycle with no
//	           fix cycle left; question node Question was raised for its
//	           surviving confirmed finding Finding, on leaf node Leaf
//	           (handle Handle) when it names one of the change set's leaves,
//	           linked blocks to it; Finding is empty for a question on the
//	           check's failure itself)
//	code-check: Step, Outcome (pass or fail), Detail (rhiza run code-check
//	           ran rhiza check and build-check on the run's change set while
//	           the run was at step Step; Detail names each cause of a fail)
//	process:   Process (the process that began running the run: when it
//	           started or was resumed; the latest process entry wins)
//	daemon-resumed: Detail (the daemon opened a tab running rhiza run
//	           resume on the run; it does so once)
//	signoff:   Step, Visit, Actor, Rule (owner or autonomous), Outcome (approved
//	           or changes), Findings (signoff step Step's visit Visit was
//	           signed off: approved, or with changes, the owner's findings,
//	           each naming a check, which count as confirmed)
//	           A check-cycle with Kind owner also has Findings: the owner's,
//	           from the signoff entry it consumes; it runs the checks they
//	           name, one cycle outside the cap
type Entry struct {
	Type    string    `json:"type"`
	At      time.Time `json:"at"`
	Step    string    `json:"step,omitempty"`
	Visit   int       `json:"visit,omitempty"`
	Outcome string    `json:"outcome,omitempty"`
	Next    string    `json:"next,omitempty"`
	Detail  string    `json:"detail,omitempty"`
	Session string    `json:"session,omitempty"`
	Started *Started  `json:"started,omitempty"`
	End     string    `json:"end,omitempty"`
	Block   *Block    `json:"block,omitempty"`
	Actor   string    `json:"actor,omitempty"`
	File    string    `json:"file,omitempty"`
	// Change and LogEntry are an accepted or rejected entry's change set
	// and the model log entry that moved it.
	Change   string `json:"change,omitempty"`
	LogEntry string `json:"log_entry,omitempty"`
	// Question is an arbiter or arbitrated entry's question node ID.
	Question string `json:"question,omitempty"`
	// Decision is an arbitrated entry's answering Decision's handle, when
	// the arbiter passed, and a leaf-round entry's.
	Decision string `json:"decision,omitempty"`
	// Handle is an asked entry's question node handle.
	Handle string `json:"handle,omitempty"`
	// Restated is a recommended entry's plain-language restatement of the
	// question.
	Restated string `json:"restated,omitempty"`
	// Asker is an advisor or recommended entry's asking session, or the
	// question node's ID when the engine raised the question itself.
	Asker      string    `json:"asker,omitempty"`
	To         string    `json:"to,omitempty"`
	Text       string    `json:"text,omitempty"`
	Summary    string    `json:"summary,omitempty"`
	Structured bool      `json:"structured,omitempty"`
	Tokens     *TokenUse `json:"tokens,omitempty"`
	// Findings are a findings entry's findings; Verdicts a verdict entry's.
	Findings []Finding `json:"findings,omitempty"`
	Verdicts []Verdict `json:"verdicts,omitempty"`
	// Changes are a gated entry's plan change sets, by ID.
	Changes []string `json:"changes,omitempty"`
	// Leaf is a leaf entry's leaf node ID, and Role its session's role:
	// builder or verifier. Attempt counts a leaf's attempts from 1. Model
	// is the model a leaf session runs. Commit is, on a builder's
	// leaf-session, the worktree's HEAD when it started; on a verifier's,
	// the commit it checks; on a passed builder's leaf-attempt, the commit
	// naming its leaf; on a passed verifier's leaf-attempt and a passed
	// leaf entry, the commit of its Evidence entries alone.
	Leaf    string `json:"leaf,omitempty"`
	Role    string `json:"role,omitempty"`
	Attempt int    `json:"attempt,omitempty"`
	Model   string `json:"model,omitempty"`
	Commit  string `json:"commit,omitempty"`
	// QuestionHandle is a failed leaf's or a leaf-round entry's question
	// node handle (Question is its ID).
	QuestionHandle string `json:"question_handle,omitempty"`
	// Check is a checks step entry's check name, and Cycle its cycle,
	// from 1. Checks are a check-cycle entry's checks, by name.
	Check  string   `json:"check,omitempty"`
	Cycle  int      `json:"cycle,omitempty"`
	Checks []string `json:"checks,omitempty"`
	// Kind is a check-cycle entry's kind: "" for a cycle within the cap,
	// CycleAnswered for one past it. Questions and Decisions are an
	// answered cycle's questions (IDs) and answering Decisions (handles).
	Kind      string   `json:"kind,omitempty"`
	Questions []string `json:"questions,omitempty"`
	Decisions []string `json:"decisions,omitempty"`
	// Finding is a review-question entry's finding ID.
	Finding string `json:"finding,omitempty"`
	// Rule is a signoff entry's rule: owner or autonomous.
	Rule string `json:"rule,omitempty"`
	// Process is a process entry's process.
	Process *Process `json:"process,omitempty"`
}

// Finding is one finding: of a plan review round, node Node (its handle
// Handle when it was recorded) has a problem of kind Kind, for Reason; of
// a leaf's verifier, only ID and Reason are set; of a review check, check
// Check found a problem of severity Severity at Place (file:line or a
// node), about leaf node Leaf (handle LeafHandle) when it names one. ID
// is f1, f2, ... in the round, or in the leaf's attempt.
//
// Every field but ID and Reason is optional, so a finding already
// journalled keeps its shape.
type Finding struct {
	ID         string `json:"id"`
	Node       string `json:"node,omitempty"`
	Handle     string `json:"handle,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Check      string `json:"check,omitempty"`
	Severity   string `json:"severity,omitempty"`
	Leaf       string `json:"leaf,omitempty"`
	LeafHandle string `json:"leaf_handle,omitempty"`
	Place      string `json:"place,omitempty"`
	Reason     string `json:"reason"`
}

// Verdict is a critic's verdict on the finding with ID Finding: confirmed
// or dropped, for Reason.
type Verdict struct {
	Finding string `json:"finding"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

// Round is one plan review round: a findings entry, numbered from 1 in
// journal order, and the verdicts journalled after it and before the next.
// A review check's rounds are the same, numbered per check, with Cycle the
// check's cycle the round was recorded in; a plan review round's is 0.
type Round struct {
	Number   int
	Cycle    int
	At       time.Time
	Findings []Finding
	Verdicts []Verdict
}

// Finding is the round's finding with ID id.
func (r Round) Finding(id string) (Finding, bool) {
	for _, f := range r.Findings {
		if f.ID == id {
			return f, true
		}
	}
	return Finding{}, false
}

// Verdict is the latest verdict on the finding with ID id.
func (r Round) Verdict(id string) (Verdict, bool) {
	for i := len(r.Verdicts) - 1; i >= 0; i-- {
		if r.Verdicts[i].Finding == id {
			return r.Verdicts[i], true
		}
	}
	return Verdict{}, false
}

// Confirmed are the round's findings whose latest verdict is confirmed, in
// order.
func (r Round) Confirmed() []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if v, ok := r.Verdict(f.ID); ok && v.Verdict == Confirmed {
			out = append(out, f)
		}
	}
	return out
}

// Serious are the round's blocker and major findings, in order: those a
// review check's critic weighs.
func (r Round) Serious() []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if Serious(f.Severity) {
			out = append(out, f)
		}
	}
	return out
}

// Standing are the round's serious findings no verdict dropped, in order:
// one with no verdict stands, as confirmed (R8). A review check fails
// while any stands.
func (r Round) Standing() []Finding {
	var out []Finding
	for _, f := range r.Serious() {
		if v, ok := r.Verdict(f.ID); !ok || v.Verdict != Dropped {
			out = append(out, f)
		}
	}
	return out
}

// Asked is one asked entry: in an autonomous run, agent step Step's
// session Session asked the owner Text, and the engine added question node
// Question (handle Handle) to the run's change set for the arbiter.
// Answered is the handle of the arbiter's Decision once the engine has
// delivered it to the asking session (an answered entry), else "".
type Asked struct {
	Step     string
	Session  string
	Question string
	Handle   string
	Text     string
	At       time.Time
	Answered string
}

// TokenCounts are the tokens a set of model calls used, by kind.
type TokenCounts struct {
	Input         int64 `json:"input"`
	Output        int64 `json:"output"`
	CacheCreation int64 `json:"cache_creation"`
	CacheRead     int64 `json:"cache_read"`
}

// TokenUse is a reading of a session's token use. Session counts the
// session's own calls; Helpers counts the calls its helper agents logged in
// its transcript.
type TokenUse struct {
	Session TokenCounts `json:"session"`
	Helpers TokenCounts `json:"helpers"`
}

// Message is one message entry: a session sent another a message while its
// step ran. Actor is the sender's actor handle; Structured is true when the
// message was a structured one rather than plain text.
type Message struct {
	Step       string
	Session    string
	Actor      string
	To         string
	Text       string
	Summary    string
	Structured bool
	At         time.Time
}

// RunSession is one session the run launched or read token use for, with
// the step it first appeared under.
type RunSession struct {
	Step    string
	Session string
}

// Amendment is one amended entry: a step's brief or plan file replaced
// while the step ran. File is relative to the run's worktree and Actor is
// the sender's actor handle.
type Amendment struct {
	Step    string
	Session string
	Actor   string
	File    string
	At      time.Time
}

// Accept is one accepted or rejected entry: the run moved change set
// Change past the owner's gate under rule autonomous, with model log entry
// LogEntry. Verb is the entry's type, accepted or rejected.
type Accept struct {
	Change   string
	LogEntry string
	Verb     string
	At       time.Time
}

// Arbitration is one arbiter entry, completed by its arbitrated entry: the
// arbiter's session Session decided question node Question, asked from
// agent step Step. Outcome is empty while it is open (no arbitrated entry
// yet); then it is pass, fail or error, with Detail, and Decision is the
// answering Decision's handle on a pass. An arbitrated entry with no open
// arbiter entry (the arbiter could not be launched) has no Session. At is
// when it began.
type Arbitration struct {
	Step     string
	Session  string
	Question string
	Outcome  string
	Detail   string
	Decision string
	At       time.Time
}

// Open reports whether the arbitration has no verdict yet.
func (a Arbitration) Open() bool { return a.Outcome == "" }

// LeafSession is one leaf-session entry: a session of role Role (builder
// or verifier) on the leaf's attempt Attempt. A builder's Commit is the
// worktree's HEAD when it started; a verifier's is the commit it checks.
// On a checks step, Cycle is the review cycle of the fix it belongs to,
// whose tries have attempts of their own (leafFixAttempt); 0 elsewhere.
//
// Seq is the position of its entry in the journal, which orders a step's
// sessions across leaves. PaneClosed is true once a pane-closed entry
// names the session: its pane is gone from the change set's tab.
type LeafSession struct {
	Role       string
	Attempt    int
	Cycle      int
	Session    string
	Model      string
	Commit     string
	Detail     string
	At         time.Time
	Seq        int
	PaneClosed bool
}

// LeafAttempt is one leaf-attempt entry: role Role's session of attempt
// Attempt was judged pass or fail (Outcome), for Detail. A passed
// builder's Commit is the new commit naming its leaf; a passed verifier's
// is the commit of its Evidence entries.
type LeafAttempt struct {
	Role    string
	Attempt int
	Outcome string
	Commit  string
	Detail  string
	At      time.Time
}

// LeafFinding is one finding a verifier journalled on a leaf's attempt
// Attempt with `rhiza run findings --leaf`.
type LeafFinding struct {
	Attempt int
	Finding Finding
	At      time.Time
}

// LeafRound is one leaf-round entry: the leaf's round Round (from 2)
// began after attempt After, because Decision (by handle) answers question
// Question (ID; handle QuestionHandle) its failure in the round before
// raised.
type LeafRound struct {
	Round          int
	After          int
	Question       string
	QuestionHandle string
	Decision       string
	At             time.Time
}

// LeafState is one leaf of a leaves step, folded from its entries: its
// sessions, judged attempts and findings in journal order, and its outcome
// (passed, failed or skipped) once a leaf entry has one; empty while it is
// open. Attempts are numbered across rounds: round 2 of a leaf that used
// its 3 attempts starts at attempt 4.
type LeafState struct {
	Step     string
	Leaf     string
	Handle   string
	Sessions []LeafSession
	Attempts []LeafAttempt
	Findings []LeafFinding
	// Rounds are the leaf's fresh rounds, in order; none in round 1.
	Rounds  []LeafRound
	Outcome string
	// Question and QuestionHandle are the question node a failed leaf's
	// failure raised, if any.
	Question       string
	QuestionHandle string
	// Commit is a passed leaf's Evidence commit.
	Commit string
	Detail string
}

// Round is the leaf's current round, from 1.
func (l LeafState) Round() int { return len(l.Rounds) + 1 }

// RoundBase is how many attempts came before the current round: 0 in
// round 1.
func (l LeafState) RoundBase() int {
	if n := len(l.Rounds); n > 0 {
		return l.Rounds[n-1].After
	}
	return 0
}

// LastRound is the leaf's latest fresh round, if it has one.
func (l LeafState) LastRound() (LeafRound, bool) {
	if n := len(l.Rounds); n > 0 {
		return l.Rounds[n-1], true
	}
	return LeafRound{}, false
}

// AttemptFindings are the findings journalled on the leaf's attempt, in
// order.
func (l LeafState) AttemptFindings(attempt int) []Finding {
	var out []Finding
	for _, f := range l.Findings {
		if f.Attempt == attempt {
			out = append(out, f.Finding)
		}
	}
	return out
}

// Session is the leaf's session of role on attempt, if it has one.
func (l LeafState) Session(role string, attempt int) (LeafSession, bool) {
	for i := len(l.Sessions) - 1; i >= 0; i-- {
		if s := l.Sessions[i]; s.Role == role && s.Attempt == attempt {
			return s, true
		}
	}
	return LeafSession{}, false
}

// CycleAttempts are the attempts of role's sessions on the leaf in review
// cycle cycle of a checks step, in journal order: for its fixer, one per
// try of the cycle's fix.
func (l LeafState) CycleAttempts(role string, cycle int) []int {
	var out []int
	for _, s := range l.Sessions {
		if s.Role == role && s.Cycle == cycle {
			out = append(out, s.Attempt)
		}
	}
	return out
}

// Judged is the judgement of role's session on attempt, if it has one.
func (l LeafState) Judged(role string, attempt int) (LeafAttempt, bool) {
	for i := len(l.Attempts) - 1; i >= 0; i-- {
		if a := l.Attempts[i]; a.Role == role && a.Attempt == attempt {
			return a, true
		}
	}
	return LeafAttempt{}, false
}

// Attempt is the leaf's latest attempt: the highest one any session has,
// or 0 before its first session.
func (l LeafState) Attempt() int {
	n := 0
	for _, s := range l.Sessions {
		n = max(n, s.Attempt)
	}
	return n
}

// CheckCycle is one check-cycle entry: checks step Step's cycle Cycle,
// from 1, runs the checks named Checks, in the step's order. Kind is
// CycleAnswered for a fresh cycle past the cap, which Decisions (handles)
// granted by answering Questions (IDs), in the same order; CycleOwner
// for one outside the cap, for the owner's Findings.
type CheckCycle struct {
	Step      string
	Cycle     int
	Checks    []string
	Kind      string
	Questions []string
	Decisions []string
	Findings  []Finding
	At        time.Time
}

// Answered says whether the cycle is a fresh one past the cap, granted by
// answers to the questions the cycle before raised.
func (c CheckCycle) Answered() bool { return c.Kind == CycleAnswered }

// Owner says whether the cycle is the owner's: it fixes the findings of
// the owner's request for changes, then runs the checks they name.
func (c CheckCycle) Owner() bool { return c.Kind == CycleOwner }

// Signoff is one signoff entry: signoff step Step's visit Visit was
// signed off by Actor under Rule (owner or autonomous), approved or with
// changes (Outcome), the owner's Findings for changes. A request for
// changes is Consumed once a checks step's owner cycle has taken its
// findings.
type Signoff struct {
	Step     string
	Visit    int
	Actor    string
	Rule     string
	Outcome  string
	Findings []Finding
	Consumed bool
	At       time.Time
}

// ReviewQuestion is one review-question entry: check Check of checks step
// Step still failed cycle Cycle with no fix cycle left, and question node
// Question (handle QuestionHandle) was raised for its surviving confirmed
// finding Finding, on leaf node Leaf (handle Handle) when it names one;
// Finding is "" for a question on the check's failure itself.
type ReviewQuestion struct {
	Step           string
	Check          string
	Cycle          int
	Finding        string
	Leaf           string
	Handle         string
	Question       string
	QuestionHandle string
	// Text is the question as raised.
	Text string
	At   time.Time
}

// CheckSession is one check-session entry: a session of role Role (check)
// on the check's cycle Cycle. Seq is the position of its entry in the
// journal, which orders a step's sessions across checks. PaneClosed is
// true once a pane-closed entry names the session.
type CheckSession struct {
	Role       string
	Cycle      int
	Session    string
	Model      string
	Detail     string
	At         time.Time
	Seq        int
	PaneClosed bool
}

// CheckJudged is one check-judged entry: role Role's session of cycle
// Cycle was judged pass or fail (Outcome), for Detail.
type CheckJudged struct {
	Role    string
	Cycle   int
	Outcome string
	Detail  string
	At      time.Time
}

// CheckOutcome is one check entry: the check passed or failed (Outcome)
// cycle Cycle, for Detail.
type CheckOutcome struct {
	Cycle   int
	Outcome string
	Detail  string
	At      time.Time
}

// ReviewFix is a checks step's change-set fixer in cycle Cycle, on its
// Attempt (from 1; a fixer that failed is retried within its cycle), folded
// from its review-fix entry and, once judged, its review-fix-judged
// entry: the fixer for the confirmed findings that name no leaf of the
// change set. Since is the worktree's HEAD when it started, which its
// commit must follow. Judged is true once its work was judged pass or
// fail (Outcome), for Detail; a passed fixer's Commit is the new commit
// whose Review: trailer names the run. Seq is the position of its
// review-fix entry in the journal, and PaneClosed true once a
// pane-closed entry names its session.
type ReviewFix struct {
	Step       string
	Cycle      int
	Attempt    int
	Session    string
	Model      string
	Since      string
	At         time.Time
	Seq        int
	PaneClosed bool
	Judged     bool
	Outcome    string
	Commit     string
	Detail     string
	JudgedAt   time.Time
}

// CheckState is one check of a checks step, folded from its entries: its
// sessions, judgements and outcomes, one per cycle it ran, in journal
// order.
type CheckState struct {
	Step     string
	Check    string
	Sessions []CheckSession
	Judged   []CheckJudged
	Outcomes []CheckOutcome
	// Rounds are the check's findings rounds, from findings entries naming
	// it, each with the verdicts journalled on it.
	Rounds []Round
}

// LatestRound is the check's latest findings round, if it has one.
func (c CheckState) LatestRound() (Round, bool) {
	if n := len(c.Rounds); n > 0 {
		return c.Rounds[n-1], true
	}
	return Round{}, false
}

// RoundFor is the check's latest findings round recorded in cycle, if it
// has one: none means the check's session of the cycle recorded none.
func (c CheckState) RoundFor(cycle int) (Round, bool) {
	for i := len(c.Rounds) - 1; i >= 0; i-- {
		if c.Rounds[i].Cycle == cycle {
			return c.Rounds[i], true
		}
	}
	return Round{}, false
}

// latestSessionCycle is the cycle of the check's latest session, 0 with
// none.
func (c CheckState) latestSessionCycle() int {
	if n := len(c.Sessions); n > 0 {
		return c.Sessions[n-1].Cycle
	}
	return 0
}

// Session is the check's session of role on cycle, if it has one.
func (c CheckState) Session(role string, cycle int) (CheckSession, bool) {
	for i := len(c.Sessions) - 1; i >= 0; i-- {
		if s := c.Sessions[i]; s.Role == role && s.Cycle == cycle {
			return s, true
		}
	}
	return CheckSession{}, false
}

// Judgement is the judgement of role's session on cycle, if it has one.
func (c CheckState) Judgement(role string, cycle int) (CheckJudged, bool) {
	for i := len(c.Judged) - 1; i >= 0; i-- {
		if j := c.Judged[i]; j.Role == role && j.Cycle == cycle {
			return j, true
		}
	}
	return CheckJudged{}, false
}

// Outcome is the check's outcome in cycle, if it has one.
func (c CheckState) Outcome(cycle int) (CheckOutcome, bool) {
	for i := len(c.Outcomes) - 1; i >= 0; i-- {
		if o := c.Outcomes[i]; o.Cycle == cycle {
			return o, true
		}
	}
	return CheckOutcome{}, false
}

// CheckOutcome is check c's current outcome in st: its latest cycle's, or
// none while that cycle is still open, even when an earlier cycle was
// judged. Open is the cycle number that is open then (0 when none is),
// and Prior the outcome of the cycle judged before it, if any.
func (s State) CheckOutcome(c CheckState) (o CheckOutcome, ok bool) {
	open, _, _ := s.CheckOpen(c)
	if open > 0 {
		return CheckOutcome{}, false
	}
	return c.Latest()
}

// CheckOpen reports whether c's latest cycle in st is open: its cycle
// number (0 when it is not), and the outcome of the newest earlier cycle
// that was judged, if any.
func (s State) CheckOpen(c CheckState) (cycle int, prior CheckOutcome, hasPrior bool) {
	latest := 0
	for _, cy := range s.Cycles {
		if cy.Step == c.Step && slices.Contains(cy.Checks, c.Check) && cy.Cycle > latest {
			latest = cy.Cycle
		}
	}
	if latest == 0 {
		return 0, CheckOutcome{}, false
	}
	if _, ok := c.Outcome(latest); ok {
		return 0, CheckOutcome{}, false
	}
	for _, o := range c.Outcomes {
		if o.Cycle < latest && (!hasPrior || o.Cycle >= prior.Cycle) {
			prior, hasPrior = o, true
		}
	}
	return latest, prior, hasPrior
}

// Latest is the check's latest outcome, if it has one.
func (c CheckState) Latest() (CheckOutcome, bool) {
	if n := len(c.Outcomes); n > 0 {
		return c.Outcomes[n-1], true
	}
	return CheckOutcome{}, false
}

// Journal is a run's entries in order.
type Journal struct {
	ID      string
	Entries []Entry
}

// StepDone is one ended step.
type StepDone struct {
	Step    string
	Visit   int
	Outcome string
	Next    string
	Detail  string
	Session string
	At      time.Time
	// Retried is true for the first visit of Step after a retried entry
	// that reopened the run at it.
	Retried bool
}

// Retry is one retried entry: the run had failed (Detail says why) and
// was reopened at Step.
type Retry struct {
	Step   string
	Detail string
	At     time.Time
}

// State is a journal folded.
type State struct {
	ID             string
	Workflow       string
	Change         string
	Title          string
	Worktree       string
	PermissionMode string
	Agent          string
	Definition     string
	Includes       map[string]string
	StartedAt      time.Time
	// Autonomous is the started entry's Autonomous: the run goes on
	// without the owner, and every agent it launches is told so.
	Autonomous bool
	// Milestone is the started entry's Milestone.
	Milestone string
	// MilestonePlan is the started entry's MilestonePlan.
	MilestonePlan bool
	// Current is the step entered and not ended, or the next step to
	// enter. Empty when the run is finished.
	Current      string
	CurrentVisit int
	Visits       map[string]int
	Steps        []StepDone
	// Session holds the latest session per unfinished step.
	Session map[string]string
	Ended   bool
	// Result is succeeded or failed; empty while the run is not finished.
	Result string
	Detail string
	// Stopped is the detail of the last stopped entry after the last
	// entered or ended entry, if any.
	Stopped string
	// Blocked is the latest open block of Blocks: what the run waits on,
	// from the last blocked entry still open; nil when no block is open.
	// BlockedStep and BlockedSession name the step and session it belongs
	// to.
	Blocked        *Block
	BlockedStep    string
	BlockedSession string
	// Blocks are the open blocks, one per session, in the order of their
	// blocked entries. A session's block is open from its blocked entry
	// until an unblocked entry for the session, an ended entry of its step,
	// or a finished or retried entry. A leaves step polls several sessions
	// at once, so several can be blocked together.
	Blocks []SessionBlock
	// Amendments are the amended entries in journal order. An amendment
	// changes nothing else in the fold.
	Amendments []Amendment
	// FinishedStep is the step the finished entry names: the step whose
	// end finished the run, or the one that went over its max_visits.
	FinishedStep string
	// Retries are the retried entries in journal order.
	Retries []Retry
	// Accepts are the accepted and rejected entries in journal order. They
	// change nothing else in the fold.
	Accepts []Accept
	// Arbitrations are the arbiter entries in journal order, each completed
	// by its arbitrated entry. They change nothing else in the fold.
	Arbitrations []Arbitration
	// Asked are the asked entries in journal order. They change nothing
	// else in the fold.
	Asked []Asked
	// Advice are the advisor entries in journal order, each completed by
	// its recommended entry. They change nothing else in the fold.
	Advice []Advice
	// RetryBase holds each step's Visits as they stood at the run's last
	// retried entry. Only visits past it count against max_visits, so a
	// retry gives every step its full cap again; visit numbers still
	// count from the run's start. Nil until a retry.
	RetryBase map[string]int
	// retryNext is the step the last retried entry reopened the run at,
	// until its next entered entry marks that visit as the retry.
	retryNext  string
	retryVisit map[string]int
	// Messages are the message entries in journal order. A message changes
	// nothing else in the fold.
	Messages []Message
	// Sessions holds one entry per distinct session ID, in the order of its
	// first session or tokens entry, with the step from that entry.
	Sessions []RunSession
	// Tokens holds each session's latest token use, keyed by session ID:
	// the last tokens entry wins. A tokens entry changes nothing else in
	// the fold.
	Tokens map[string]TokenUse
	// Rounds are the plan review rounds in journal order. A verdict entry
	// before any findings entry belongs to no round and is dropped. Rounds
	// change nothing else in the fold.
	Rounds []Round
	// Gated are the plan's change sets the latest gated entry names, by
	// ID; nil before run gate. They change nothing else in the fold.
	Gated []string
	// GatedAt is when the latest gated entry was journalled.
	GatedAt time.Time
	// Leaves are the leaves of the leaves steps, in the order their first
	// entry appears. A new visit of a leaves step starts its leaves anew:
	// its entered entry drops the step's earlier leaves, except the visit
	// a retried entry reopened the run at, which keeps them, so passed
	// leaves stay passed. A checks step's leaves (its fixers and their
	// verifiers) survive a new visit, as its Checks do. A leaf session is
	// also in Sessions; nothing else in the fold changes.
	Leaves []LeafState
	// Checks are the checks of the checks steps, in the order their first
	// entry appears, and Cycles the checks steps' cycles in journal order.
	// Unlike Leaves they survive a new visit of their step: a review
	// re-entered after the owner asks for changes keeps its history. A
	// check session is also in Sessions; nothing else in the fold changes.
	Checks []CheckState
	Cycles []CheckCycle
	// Fixes are the checks steps' change-set fixers, one per cycle that
	// ran one, in journal order. Like Checks they survive a new visit of
	// their step. A fixer session is also in Sessions; nothing else in the
	// fold changes.
	Fixes []ReviewFix
	// ReviewQuestions are the checks steps' review-question entries, in
	// journal order. Like Checks they survive a new visit of their step.
	ReviewQuestions []ReviewQuestion
	// CodeChecks are the code-check entries, in journal order. Like
	// Checks they survive a new visit of their step.
	CodeChecks []CodeCheck
	// Signoffs are the signoff entries, in journal order.
	Signoffs []Signoff
	// Process is the latest process entry's process; nil when the journal
	// has none (a run journalled before processes were recorded).
	Process *Process
	// ProcessAt is when that process entry was journalled.
	ProcessAt time.Time
	// DaemonResumed is true once the daemon resumed the run, and
	// DaemonResumedAt is when.
	DaemonResumed   bool
	DaemonResumedAt time.Time
	// StoppedAt is when the last stopped entry was journalled, and EndedAt
	// when the finished entry was; zero while there is none.
	StoppedAt time.Time
	EndedAt   time.Time
	// checksSteps are the steps with a check-cycle entry: checks steps,
	// whose leaves (fixers and verifiers) survive a new visit.
	checksSteps map[string]bool
}

// SignoffFor is the latest sign-off of signoff step step's visit, if it
// has one.
func (st State) SignoffFor(step string, visit int) (Signoff, bool) {
	for i := len(st.Signoffs) - 1; i >= 0; i-- {
		if s := st.Signoffs[i]; s.Step == step && s.Visit == visit {
			return s, true
		}
	}
	return Signoff{}, false
}

// PendingChanges is the latest request for changes no owner cycle has
// consumed yet, if there is one.
func (st State) PendingChanges() (Signoff, bool) {
	for i := len(st.Signoffs) - 1; i >= 0; i-- {
		if s := st.Signoffs[i]; s.Outcome == SignoffChanges && !s.Consumed {
			return s, true
		}
	}
	return Signoff{}, false
}

// CodeCheck is one code-check entry: rhiza run code-check's outcome
// (pass or fail) while the run was at step Step, and Detail, the causes
// of a fail.
type CodeCheck struct {
	Step    string
	Outcome string
	Detail  string
	At      time.Time
}

// OwnerCycle says whether checks step step's cycle is an owner cycle:
// one that fixes the owner's request for changes.
func (st State) OwnerCycle(step string, cycle int) bool {
	for _, c := range st.Cycles {
		if c.Step == step && c.Cycle == cycle {
			return c.Owner()
		}
	}
	return false
}

// LatestCodeCheck is the journal's latest code check, if it has one.
func (st State) LatestCodeCheck() (CodeCheck, bool) {
	if len(st.CodeChecks) == 0 {
		return CodeCheck{}, false
	}
	return st.CodeChecks[len(st.CodeChecks)-1], true
}

// ReviewQuestionsFor are checks step step's questions raised when cycle
// ended with no fix cycle left, in journal order.
func (st State) ReviewQuestionsFor(step string, cycle int) []ReviewQuestion {
	var out []ReviewQuestion
	for _, q := range st.ReviewQuestions {
		if q.Step == step && q.Cycle == cycle {
			out = append(out, q)
		}
	}
	return out
}

// ReviewFix is checks step step's change-set fixer in cycle, its latest
// attempt, if the journal has one.
func (st State) ReviewFix(step string, cycle int) (ReviewFix, bool) {
	var out ReviewFix
	found := false
	for _, f := range st.Fixes {
		if f.Step == step && f.Cycle == cycle {
			out, found = f, true
		}
	}
	return out, found
}

// Check is step's check named name, if the journal has one.
func (st State) Check(step, name string) (CheckState, bool) {
	for _, c := range st.Checks {
		if c.Step == step && c.Check == name {
			return c, true
		}
	}
	return CheckState{}, false
}

// LatestCycle is checks step step's latest cycle, if it has one.
func (st State) LatestCycle(step string) (CheckCycle, bool) {
	for i := len(st.Cycles) - 1; i >= 0; i-- {
		if c := st.Cycles[i]; c.Step == step {
			return c, true
		}
	}
	return CheckCycle{}, false
}

// check is step's check state for e, added in first-seen order when new.
func (st *State) check(e Entry) *CheckState {
	for i := range st.Checks {
		if c := &st.Checks[i]; c.Step == e.Step && c.Check == e.Check {
			return c
		}
	}
	st.Checks = append(st.Checks, CheckState{Step: e.Step, Check: e.Check})
	return &st.Checks[len(st.Checks)-1]
}

// closeFixPane marks the pane of the change-set fixer session pane-closed
// entry e names closed, and says whether e named one.
func (st *State) closeFixPane(e Entry) bool {
	for i := range st.Fixes {
		if f := &st.Fixes[i]; f.Step == e.Step && f.Session == e.Session {
			f.PaneClosed = true
			return true
		}
	}
	return false
}

// Leaf is step's leaf with node ID leaf, if the journal has one.
func (st State) Leaf(step, leaf string) (LeafState, bool) {
	for _, l := range st.Leaves {
		if l.Step == step && l.Leaf == leaf {
			return l, true
		}
	}
	return LeafState{}, false
}

// leaf is step's leaf state for e, added in first-seen order when new;
// it takes e's handle when it had none.
func (st *State) leaf(e Entry) *LeafState {
	for i := range st.Leaves {
		if l := &st.Leaves[i]; l.Step == e.Step && l.Leaf == e.Leaf {
			if l.Handle == "" {
				l.Handle = e.Handle
			}
			return l
		}
	}
	st.Leaves = append(st.Leaves, LeafState{Step: e.Step, Leaf: e.Leaf, Handle: e.Handle})
	return &st.Leaves[len(st.Leaves)-1]
}

// State folds the journal.
func (j Journal) State() State {
	st := State{ID: j.ID, Visits: map[string]int{}, Session: map[string]string{}, Tokens: map[string]TokenUse{}}
	for seq, e := range j.Entries {
		switch e.Type {
		case TypeStarted:
			st.StartedAt = e.At
			st.Current = e.Step
			if s := e.Started; s != nil {
				st.Workflow = s.Workflow
				st.Change = s.Change
				st.Title = s.Title
				st.Worktree = s.Worktree
				st.PermissionMode = s.PermissionMode
				st.Agent = s.Agent
				st.Definition = s.Definition
				st.Includes = s.Includes
				st.Autonomous = s.Autonomous
				st.Milestone = s.Milestone
				st.MilestonePlan = s.MilestonePlan
			}
		case TypeProcess:
			if e.Process != nil {
				p := *e.Process
				st.Process = &p
				st.ProcessAt = e.At
			}
		case TypeDaemonResumed:
			st.DaemonResumed, st.DaemonResumedAt = true, e.At
		case TypeEntered:
			st.Visits[e.Step]++
			st.Current = e.Step
			st.CurrentVisit = e.Visit
			if st.CurrentVisit == 0 {
				st.CurrentVisit = st.Visits[e.Step]
			}
			st.Stopped = ""
			if st.retryNext == e.Step {
				if st.retryVisit == nil {
					st.retryVisit = map[string]int{}
				}
				st.retryVisit[e.Step] = st.CurrentVisit
			}
			retried := st.retryNext == e.Step
			st.retryNext = ""
			if !retried && !st.checksSteps[e.Step] {
				st.Leaves = slices.DeleteFunc(st.Leaves, func(l LeafState) bool { return l.Step == e.Step })
			}
		case TypeLeafSession:
			l := st.leaf(e)
			l.Sessions = append(l.Sessions, LeafSession{
				Role: e.Role, Attempt: e.Attempt, Cycle: e.Cycle, Session: e.Session, Model: e.Model, Commit: e.Commit, Detail: e.Detail, At: e.At,
				Seq: seq,
			})
			st.addSession(e.Step, e.Session)
		case TypeLeafAttempt:
			l := st.leaf(e)
			l.Attempts = append(l.Attempts, LeafAttempt{Role: e.Role, Attempt: e.Attempt, Outcome: e.Outcome, Commit: e.Commit, Detail: e.Detail, At: e.At})
		case TypeLeaf:
			l := st.leaf(e)
			l.Outcome, l.Question, l.QuestionHandle, l.Commit, l.Detail = e.Outcome, e.Question, e.QuestionHandle, e.Commit, e.Detail
		case TypeCheckCycle:
			st.Cycles = append(st.Cycles, CheckCycle{
				Step: e.Step, Cycle: e.Cycle, Checks: e.Checks, Kind: e.Kind, Questions: e.Questions, Decisions: e.Decisions,
				Findings: e.Findings, At: e.At,
			})
			if st.checksSteps == nil {
				st.checksSteps = map[string]bool{}
			}
			st.checksSteps[e.Step] = true
			if e.Kind == CycleOwner {
				for i := range st.Signoffs {
					if st.Signoffs[i].Outcome == SignoffChanges {
						st.Signoffs[i].Consumed = true
					}
				}
			}
		case TypeSignoff:
			st.Signoffs = append(st.Signoffs, Signoff{
				Step: e.Step, Visit: e.Visit, Actor: e.Actor, Rule: e.Rule, Outcome: e.Outcome, Findings: e.Findings, At: e.At,
			})
		case TypeReviewQuestion:
			st.ReviewQuestions = append(st.ReviewQuestions, ReviewQuestion{
				Step: e.Step, Check: e.Check, Cycle: e.Cycle, Finding: e.Finding, Leaf: e.Leaf, Handle: e.Handle,
				Question: e.Question, QuestionHandle: e.QuestionHandle, Text: e.Detail, At: e.At,
			})
		case TypeCodeCheck:
			st.CodeChecks = append(st.CodeChecks, CodeCheck{Step: e.Step, Outcome: e.Outcome, Detail: e.Detail, At: e.At})
		case TypeCheckSession:
			c := st.check(e)
			c.Sessions = append(c.Sessions, CheckSession{
				Role: e.Role, Cycle: e.Cycle, Session: e.Session, Model: e.Model, Detail: e.Detail, At: e.At, Seq: seq,
			})
			st.addSession(e.Step, e.Session)
		case TypeCheckJudged:
			c := st.check(e)
			c.Judged = append(c.Judged, CheckJudged{Role: e.Role, Cycle: e.Cycle, Outcome: e.Outcome, Detail: e.Detail, At: e.At})
		case TypeCheck:
			c := st.check(e)
			c.Outcomes = append(c.Outcomes, CheckOutcome{Cycle: e.Cycle, Outcome: e.Outcome, Detail: e.Detail, At: e.At})
		case TypeReviewFix:
			st.Fixes = append(st.Fixes, ReviewFix{
				Step: e.Step, Cycle: e.Cycle, Attempt: max(e.Attempt, 1), Session: e.Session, Model: e.Model, Since: e.Commit, At: e.At, Seq: seq,
			})
			st.addSession(e.Step, e.Session)
		case TypeReviewFixJudged:
			for i := range st.Fixes {
				if f := &st.Fixes[i]; f.Step == e.Step && f.Cycle == e.Cycle && f.Attempt == max(e.Attempt, 1) {
					f.Judged, f.Outcome, f.Commit, f.Detail, f.JudgedAt = true, e.Outcome, e.Commit, e.Detail, e.At
				}
			}
		case TypePaneClosed:
			if e.Leaf == "" && e.Check == "" && st.closeFixPane(e) {
				break
			}
			if e.Check != "" {
				c := st.check(e)
				for i := range c.Sessions {
					if c.Sessions[i].Session == e.Session {
						c.Sessions[i].PaneClosed = true
					}
				}
				break
			}
			l := st.leaf(e)
			for i := range l.Sessions {
				if l.Sessions[i].Session == e.Session {
					l.Sessions[i].PaneClosed = true
				}
			}
		case TypeLeafRound:
			l := st.leaf(e)
			l.Rounds = append(l.Rounds, LeafRound{
				Round: l.Round() + 1, After: l.Attempt(), Question: e.Question, QuestionHandle: e.QuestionHandle, Decision: e.Decision, At: e.At,
			})
			l.Outcome, l.Question, l.QuestionHandle, l.Commit, l.Detail = "", "", "", "", ""
			// Leaves skipped for a failed dependency may start again; the
			// conductor skips anew any whose dependency is still lost.
			for i := range st.Leaves {
				if o := &st.Leaves[i]; o.Step == e.Step && o.Outcome == LeafSkipped {
					o.Outcome, o.Detail = "", ""
				}
			}
		case TypeSession:
			st.Session[e.Step] = e.Session
			st.addSession(e.Step, e.Session)
		case TypeBlocked:
			b := Block{}
			if e.Block != nil {
				b = *e.Block
			}
			st.setBlocks(append(st.blocksBut(func(o SessionBlock) bool { return o.Session == e.Session }),
				SessionBlock{Step: e.Step, Session: e.Session, Block: b, At: e.At}))
		case TypeUnblocked:
			// An unblocked entry with no session ends every block.
			st.setBlocks(st.blocksBut(func(o SessionBlock) bool { return e.Session == "" || o.Session == e.Session }))
		case TypeEnded:
			visit := e.Visit
			if visit == 0 {
				visit = st.Visits[e.Step]
			}
			st.Steps = append(st.Steps, StepDone{
				Step: e.Step, Visit: visit, Outcome: e.Outcome, Next: e.Next,
				Detail: e.Detail, Session: st.Session[e.Step], At: e.At,
				Retried: visit > 0 && st.retryVisit[e.Step] == visit,
			})
			delete(st.Session, e.Step)
			st.setBlocks(st.blocksBut(func(o SessionBlock) bool { return o.Step == e.Step }))
			st.Current = e.Next
			st.CurrentVisit = 0
			st.Stopped = ""
		case TypeAmended:
			st.Amendments = append(st.Amendments, Amendment{
				Step: e.Step, Session: e.Session, Actor: e.Actor, File: e.File, At: e.At,
			})
		case TypeAccepted, TypeRejected:
			st.Accepts = append(st.Accepts, Accept{Change: e.Change, LogEntry: e.LogEntry, Verb: e.Type, At: e.At})
		case TypeArbiter:
			st.Arbitrations = append(st.Arbitrations, Arbitration{
				Step: e.Step, Session: e.Session, Question: e.Question, At: e.At,
			})
		case TypeArbitrated:
			st.arbitrated(e)
		case TypeAdvisor:
			st.Advice = append(st.Advice, Advice{Step: e.Step, Advisor: e.Session, Asker: e.Asker, At: e.At})
		case TypeRecommended:
			st.recommended(e)
		case TypeAsked:
			st.Asked = append(st.Asked, Asked{
				Step: e.Step, Session: e.Session, Question: e.Question, Handle: e.Handle, Text: e.Detail, At: e.At,
			})
			// Only a session's question is its block's: a command step's
			// asked entry has no session.
			if e.Session != "" {
				for i := range st.Blocks {
					if st.Blocks[i].Session == e.Session {
						st.Blocks[i].asked = len(st.Asked)
					}
				}
			}
		case TypeAnswered:
			for i := len(st.Asked) - 1; i >= 0; i-- {
				if st.Asked[i].Question == e.Question {
					st.Asked[i].Answered = e.Decision
					break
				}
			}
		case TypeMessage:
			st.Messages = append(st.Messages, Message{
				Step: e.Step, Session: e.Session, Actor: e.Actor, To: e.To,
				Text: e.Text, Summary: e.Summary, Structured: e.Structured, At: e.At,
			})
		case TypeTokens:
			st.addSession(e.Step, e.Session)
			if e.Tokens != nil {
				st.Tokens[e.Session] = *e.Tokens
			}
		case TypeGated:
			st.Gated, st.GatedAt = e.Changes, e.At
		case TypeFindings:
			if e.Check != "" {
				c := st.check(e)
				cycle := e.Cycle
				if cycle == 0 {
					cycle = c.latestSessionCycle()
				}
				c.Rounds = append(c.Rounds, Round{Number: len(c.Rounds) + 1, Cycle: cycle, At: e.At, Findings: e.Findings})
				break
			}
			if e.Leaf != "" {
				l := st.leaf(e)
				for _, f := range e.Findings {
					l.Findings = append(l.Findings, LeafFinding{Attempt: e.Attempt, Finding: f, At: e.At})
				}
				break
			}
			st.Rounds = append(st.Rounds, Round{Number: len(st.Rounds) + 1, At: e.At, Findings: e.Findings})
		case TypeVerdict:
			if e.Check != "" {
				c := st.check(e)
				if n := len(c.Rounds); n > 0 {
					c.Rounds[n-1].Verdicts = append(c.Rounds[n-1].Verdicts, e.Verdicts...)
				}
				break
			}
			if n := len(st.Rounds); n > 0 {
				st.Rounds[n-1].Verdicts = append(st.Rounds[n-1].Verdicts, e.Verdicts...)
			}
		case TypeStopped:
			st.Stopped, st.StoppedAt = e.Detail, e.At
		case TypeFinished:
			st.clearBlock()
			st.Ended, st.EndedAt = true, e.At
			st.Result = e.End
			st.Detail = e.Detail
			st.Current = ""
			st.CurrentVisit = 0
			st.FinishedStep = e.Step
			st.retryNext = ""
		case TypeRetried:
			// A retry reopens a finished run at e.Step with no open visit,
			// so the engine enters it anew. Nothing before it is undone:
			// the failed steps and the finished entry stay in the history.
			st.Retries = append(st.Retries, Retry{Step: e.Step, Detail: e.Detail, At: e.At})
			st.clearBlock()
			st.Ended = false
			st.Result = ""
			st.Detail = ""
			st.FinishedStep = ""
			st.Stopped = ""
			st.Current = e.Step
			st.CurrentVisit = 0
			delete(st.Session, e.Step)
			st.RetryBase = make(map[string]int, len(st.Visits))
			for step, n := range st.Visits {
				st.RetryBase[step] = n
			}
			st.retryNext = e.Step
		}
	}
	return st
}

// arbitrated completes the last open arbitration of e's question with e's
// verdict, or, when there is none, adds a completed one without a session.
func (st *State) arbitrated(e Entry) {
	for i := len(st.Arbitrations) - 1; i >= 0; i-- {
		a := &st.Arbitrations[i]
		if a.Open() && a.Question == e.Question {
			a.Outcome, a.Detail, a.Decision = e.Outcome, e.Detail, e.Decision
			return
		}
	}
	st.Arbitrations = append(st.Arbitrations, Arbitration{
		Step: e.Step, Question: e.Question, Outcome: e.Outcome, Detail: e.Detail, Decision: e.Decision, At: e.At,
	})
}

// LatestRound is the run's latest plan review round, if it has one.
func (st State) LatestRound() (Round, bool) {
	if n := len(st.Rounds); n > 0 {
		return st.Rounds[n-1], true
	}
	return Round{}, false
}

// OpenArbitration is the open arbitration of question, the latest when
// there are several: the arbiter launched on it and not yet judged, which a
// resume rejoins.
func (st State) OpenArbitration(question string) (Arbitration, bool) {
	for i := len(st.Arbitrations) - 1; i >= 0; i-- {
		if a := st.Arbitrations[i]; a.Open() && a.Question == question {
			return a, true
		}
	}
	return Arbitration{}, false
}

// Arbitration is the latest arbitration of question, open or judged.
func (st State) Arbitration(question string) (Arbitration, bool) {
	for i := len(st.Arbitrations) - 1; i >= 0; i-- {
		if a := st.Arbitrations[i]; a.Question == question {
			return a, true
		}
	}
	return Arbitration{}, false
}

// Advice is one advisor entry, completed by its recommended entry: the
// advisor's session Advisor was launched on the question session Asker, of
// agent step Step, is blocked on. Recommendation is empty while the
// advisor has recommended nothing (or never will: it failed); then it is
// what it recommends, with Reasons. At is when the advisor was launched.
type Advice struct {
	Step           string
	Advisor        string
	Asker          string
	Recommendation string
	Reasons        string
	Restatement    string
	At             time.Time
}

// recommended completes the latest advisor entry of e's advisor session
// with its recommendation. A recommended entry with no advisor entry
// (the journal does not show the launch) is ignored.
func (st *State) recommended(e Entry) {
	for i := len(st.Advice) - 1; i >= 0; i-- {
		if a := &st.Advice[i]; a.Advisor == e.Session && a.Asker == e.Asker {
			a.Recommendation, a.Reasons, a.Restatement = e.Text, e.Detail, e.Restated
			return
		}
	}
}

// AdviceFor is the advice on the question block sb is: the latest
// recommendation, among the advisors launched on its session since the
// block began, that has one.
func (st State) AdviceFor(sb SessionBlock) (Advice, bool) {
	for i := len(st.Advice) - 1; i >= 0; i-- {
		if a := st.Advice[i]; a.Asker == sb.Session && !a.At.Before(sb.At) && a.Recommendation != "" {
			return a, true
		}
	}
	return Advice{}, false
}

// SessionBlock is one session's open block: what session, of step Step,
// waits on.
type SessionBlock struct {
	Step    string
	Session string
	Block   Block
	// At is when its blocked entry was journalled.
	At time.Time
	// asked is 1 + the index in Asked of the asked entry for this block,
	// 0 when the block has none.
	asked int
}

// BlockAsked is the asked entry of session's open block: the question
// the engine raised for it, so a session still blocked is never asked
// twice for one block. It has none when the session is not blocked, or
// its block has not been asked yet.
func (st State) BlockAsked(session string) (Asked, bool) {
	b, ok := st.BlockOf(session)
	if !ok || b.asked == 0 {
		return Asked{}, false
	}
	return st.Asked[b.asked-1], true
}

// BlockOf is session's open block, if it has one.
func (st State) BlockOf(session string) (SessionBlock, bool) {
	for _, b := range st.Blocks {
		if b.Session == session {
			return b, true
		}
	}
	return SessionBlock{}, false
}

// CapVisits is how many visits of step count against its max_visits: those
// since the run's last retried entry, or all of them when it has none.
func (st State) CapVisits(step string) int {
	return st.Visits[step] - st.RetryBase[step]
}

// FailedStep is the step a failed run would be retried at: the step its
// finished entry names, which is the step whose end failed the run or the
// one that went over its max_visits. Journals whose finished entry names no
// step fall back to the last ended step. A run that is not finished, or
// that succeeded, has none and returns ErrNotFailed.
func (st State) FailedStep() (string, error) {
	if !st.Ended || st.Result != Failed {
		return "", ErrNotFailed
	}
	if st.FinishedStep != "" {
		return st.FinishedStep, nil
	}
	if n := len(st.Steps); n > 0 {
		return st.Steps[n-1].Step, nil
	}
	return "", ErrNotFailed
}

func (st *State) clearBlock() {
	st.setBlocks(nil)
}

// blocksBut is a new slice of the open blocks drop does not match.
func (st *State) blocksBut(drop func(SessionBlock) bool) []SessionBlock {
	var out []SessionBlock
	for _, b := range st.Blocks {
		if !drop(b) {
			out = append(out, b)
		}
	}
	return out
}

// setBlocks sets the open blocks, and Blocked, BlockedStep and
// BlockedSession to the latest of them.
func (st *State) setBlocks(bs []SessionBlock) {
	st.Blocks = bs
	st.Blocked, st.BlockedStep, st.BlockedSession = nil, "", ""
	if n := len(bs); n > 0 {
		b := bs[n-1].Block
		st.Blocked, st.BlockedStep, st.BlockedSession = &b, bs[n-1].Step, bs[n-1].Session
	}
}

// addSession records session under step in Sessions if it is not there yet.
func (st *State) addSession(step, session string) {
	for _, s := range st.Sessions {
		if s.Session == session {
			return
		}
	}
	st.Sessions = append(st.Sessions, RunSession{Step: step, Session: session})
}

// AbandonedPrefix starts the Detail of a run that rhiza run abandon failed.
const AbandonedPrefix = "abandoned by "

// Abandoned reports whether the run finished failed because it was
// abandoned, which a resume (a retry included) refuses.
func (st State) Abandoned() bool {
	return st.Ended && st.Result == Failed && strings.HasPrefix(st.Detail, AbandonedPrefix)
}

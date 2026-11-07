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
	"bytes"
	"cmp"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
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
		Summary:  "Queue a run of a workflow on a change set for the daemon, which runs its steps in its worktree headless; --foreground drives it in this process",
		Write:    true,
		Request:  reflect.TypeFor[RunStartRequest](),
		Response: reflect.TypeFor[RunResponse](),
		Run:      runRunStart,
	})
	register(Operation{
		Name:     "run.status",
		Summary:  "Show a run's state and the steps it has done, from its journal",
		ReadOnly: true,
		Request:  reflect.TypeFor[RunStatusRequest](),
		Response: reflect.TypeFor[RunInfo](),
		Run:      runRunStatus,
	})
	register(Operation{
		Name:     "run.resume",
		Summary:  "Queue the continuation of a stopped or interrupted run for the daemon, which drives it headless where its journal left it, or retries a failed run at its failed step; --foreground drives it in this process",
		Write:    true,
		Request:  reflect.TypeFor[RunResumeRequest](),
		Response: reflect.TypeFor[RunResponse](),
		Run:      runRunResume,
	})
	register(Operation{
		Name:     "run.abandon",
		Summary:  "Give up on an unfinished or failed run: journal it failed as abandoned, so the inbox stops listing it and a new run can start on its change set",
		Write:    true,
		Request:  reflect.TypeFor[RunAbandonRequest](),
		Response: reflect.TypeFor[RunAbandonResponse](),
		Run:      runRunAbandon,
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
	Workflow       string `json:"workflow,omitempty" desc:"The workflow: a .toml file (a path with a / or ending in .toml), else a name: your personal $XDG_CONFIG_HOME/rhiza/workflows/<name>.toml, else .rhiza/workflows/<name>.toml, else a built-in; default ship-change"`
	PermissionMode string `json:"permission_mode,omitempty" desc:"The permission mode of agent steps that set none; modes that bypass the agent's permission checks are refused"`
	Agent          string `json:"agent,omitempty" desc:"The agent adapter of agent steps that name none; default the built-in default adapter"`
	NoWait         bool   `json:"no_wait,omitempty" desc:"Stop at a wait step whose condition does not hold yet, instead of waiting; resume the run later"`
	Autonomous     bool   `json:"autonomous,omitempty" desc:"Run without the owner: skip the owner's gates and questions, recording each; never changes agent permissions. Resume keeps it"`
	New            string `json:"new,omitempty" desc:"With --workflow plan and no change set: open a draft change set with this one-line intent and start the plan run on it"`
	Milestone      bool   `json:"milestone,omitempty" desc:"With --new: open a record change set, to plan a milestone in; without it an intention"`
	Request        string `json:"request,omitempty" desc:"With --new: the plan request this run plans; the run takes the ID the daemon drew for it. The daemon passes it"`
	Foreground     bool   `json:"foreground,omitempty" desc:"Drive the run in this process, rather than hand it to the daemon. The daemon passes it to the detached process it starts"`
}

// RunStatusRequest names a run, or none for the latest.
type RunStatusRequest struct {
	Run string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
}

// RunResumeRequest names the run to continue.
type RunResumeRequest struct {
	Run        string `json:"run" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters"`
	NoWait     bool   `json:"no_wait,omitempty" desc:"Stop at a wait step whose condition does not hold yet, instead of waiting"`
	Retry      bool   `json:"retry,omitempty" flag:"retry" desc:"Restart a failed run at the step that failed, once its cause is fixed; the retry is journalled. A run that succeeded is refused"`
	Foreground bool   `json:"foreground,omitempty" desc:"Drive the run in this process, rather than hand it to the daemon. The daemon passes it to the detached process it starts"`
}

// RunAbandonRequest names the run to give up on.
type RunAbandonRequest struct {
	Run string `json:"run" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters"`
}

// RunAbandonResponse is the run as abandoning left it, or in a dry run
// would leave it. The command exits 0: the run is failed by the owner's
// choice.
type RunAbandonResponse struct {
	Run    RunInfo `json:"run" desc:"The run as its journal stands, failed with the abandoned detail"`
	DryRun bool    `json:"dry_run,omitempty" desc:"True if nothing was written"`
	// Cleanup reports the removal of the run's worktree and local branch.
	Cleanup []RunAbandonStep `json:"cleanup,omitempty" desc:"The clean-up of the run's worktree and local branch, in order; absent in a dry run or when the run has no worktree"`
}

// RunAbandonStep is one clean-up step of an abandon.
type RunAbandonStep struct {
	Step   string `json:"step" enum:"remove-worktree,delete-branch" desc:"remove-worktree or delete-branch"`
	Status string `json:"status" enum:"done,already-gone,kept" desc:"done, already-gone, or kept (the worktree is dirty, or git refused)"`
	Detail string `json:"detail,omitempty" desc:"Why it was kept"`
}

// RunResponse is a run after run.start or run.resume ran it, or in a dry
// run the one that would run.
type RunResponse struct {
	Run    RunInfo `json:"run" desc:"The run as its journal stands; the command exits 1 when the run failed"`
	DryRun bool    `json:"dry_run,omitempty" desc:"True if nothing was written or run"`
	Retry  string  `json:"retry,omitempty" desc:"The step a resume --retry restarted the failed run at, or in a dry run would; absent otherwise"`
	Queued bool    `json:"queued,omitempty" desc:"True when run start or run resume queued the run for the daemon rather than driving it: a started run has no ID yet, and the daemon drives it headless; absent otherwise"`
	Daemon string  `json:"daemon,omitempty" desc:"With queued: the daemon's process, as host:pid; empty when no daemon runs and the run waits for one"`
}

// RunInfo is a run folded from its journal.
type RunInfo struct {
	ID             string          `json:"id" desc:"The run ID"`
	Workflow       string          `json:"workflow" desc:"The workflow's name"`
	Change         string          `json:"change" desc:"The change set it runs on"`
	Title          string          `json:"title" desc:"The change set's title"`
	Worktree       string          `json:"worktree" desc:"The worktree every step runs in"`
	Autonomous     bool            `json:"autonomous,omitempty" desc:"True for a run started with --autonomous: it goes on without the owner; absent otherwise"`
	State          string          `json:"state" desc:"running (another process runs it now), blocked (its agent step waits on the owner), stopped (resumable), succeeded or failed; empty in a dry run"`
	Current        string          `json:"current,omitempty" desc:"The step being run or to run next; empty once the run is finished"`
	Detail         string          `json:"detail,omitempty" desc:"Why it stopped or failed"`
	Steps          []RunStep       `json:"steps" desc:"The steps done, in order"`
	StartedAt      string          `json:"started_at" desc:"When it started, RFC 3339 UTC"`
	Blocked        *RunBlocked     `json:"blocked,omitempty" desc:"What the run waits on while its agent step's session is blocked; absent otherwise"`
	Stalled        *RunStalled     `json:"stalled,omitempty" desc:"A session of a live run's current step (its agent session, or a leaf's, check's or fixer's latest session) while it is stalled: working, or idle in an autonomous run, without saying it is done and silent for more than 20 minutes, not paused at a usage limit; absent otherwise"`
	Amendments     []RunAmendment  `json:"amendments" desc:"The amendments made to its agent steps' briefs or plans, in order"`
	Messages       []RunMessage    `json:"messages" desc:"The messages its sessions sent each other, in order"`
	Sessions       []RunSession    `json:"sessions" desc:"The sessions it launched or read token use for, in order of first appearance, each with its latest reading"`
	Totals         SessionTokenUse `json:"totals" desc:"The sums of the sessions with recorded token use, the session's own and helpers' kept apart"`
	Retries        []RunRetry      `json:"retries,omitempty" desc:"The retries of the run after it failed, in order; absent when it was never retried"`
	Accepts        []RunAccept     `json:"accepts,omitempty" desc:"The change sets an autonomous run accepted or rejected in the owner's place under rule autonomous, in order; absent when there were none"`
	Questions      []RunQuestion   `json:"questions,omitempty" desc:"The questions an autonomous run's agent steps asked, each with its arbiter's verdict, in order; absent when there were none"`
	Rounds         []RunRound      `json:"rounds,omitempty" desc:"The plan review rounds journalled by run findings and run verdict, in order; absent when there were none"`
	Waiting        []RunWaiting    `json:"waiting,omitempty" desc:"The plan's change sets run gate proposed that still wait proposed for the owner, by ID; absent when none does"`
	Leaves         []RunLeaf       `json:"leaves,omitempty" desc:"The leaves its leaves steps conducted, and those its checks steps' fixers fixed, in the order they began, each with its sessions, judged attempts and outcome; absent when there were none"`
	Checks         []RunCheck      `json:"checks,omitempty" desc:"The checks its checks steps ran, in the order they began, each with its cycles and outcome; absent when there were none"`
	Fixes          []RunFix        `json:"fixes,omitempty" desc:"The change-set fixers its checks steps ran, for confirmed findings naming no leaf, one per cycle, in order; absent when there were none"`
	CodeChecks     []RunCodeCheck  `json:"code_checks,omitempty" desc:"The code checks rhiza run code-check journalled, in order; absent when there were none"`
	Signoffs       []RunSignoff    `json:"signoffs,omitempty" desc:"The sign-offs of its signoff steps, in order: the owner's approvals and requests for changes, and an autonomous run's under rule autonomous; absent when there were none"`
	PR             string          `json:"pr,omitempty" desc:"While the run waits at a signoff step: the URL of the change set's pull request, when ship open-pr recorded one; absent otherwise"`
	SignoffWaiting bool            `json:"signoff_waiting,omitempty" desc:"True while the run is stopped at a signoff step, waiting for the owner's rhiza run approve or rhiza run request-changes, with rhiza run summary rendering its verdicts; absent otherwise"`
}

// RunCodeCheck is one code check: rhiza check and build-check on the
// run's change set.
type RunCodeCheck struct {
	Step    string `json:"step,omitempty" desc:"The step the run was at"`
	Outcome string `json:"outcome" enum:"pass,fail" desc:"pass or fail"`
	Detail  string `json:"detail,omitempty" desc:"A fail's causes, separated by semicolons"`
	At      string `json:"at" desc:"When it ran, RFC 3339 UTC"`
}

// RunFix is one change-set fixer a checks step ran in a cycle.
type RunFix struct {
	Step    string `json:"step" desc:"The checks step"`
	Cycle   int    `json:"cycle" desc:"The cycle whose fix phase ran it, from 2"`
	Session string `json:"session" desc:"The fixer session's ID"`
	Model   string `json:"model,omitempty" desc:"The model it runs; empty for the agent's default"`
	Since   string `json:"since" desc:"The worktree's HEAD when it started, which its commit must follow"`
	Outcome string `json:"outcome,omitempty" desc:"pass or fail once its work was judged; empty while it runs"`
	Commit  string `json:"commit,omitempty" desc:"A passed fixer's new commit whose Review: trailer names the run"`
	Detail  string `json:"detail,omitempty" desc:"Why it was judged so"`
	At      string `json:"at" desc:"When it started, RFC 3339 UTC"`
}

// RunCheck is one check a checks step ran.
type RunCheck struct {
	Step         string             `json:"step" desc:"The checks step"`
	Check        string             `json:"check" desc:"The check's name"`
	Cycles       []RunCheckCycle    `json:"cycles" desc:"The cycles it ran in, in order"`
	Rounds       []RunCheckRound    `json:"rounds,omitempty" desc:"Its findings rounds, journalled by run findings --check and judged by run verdict --check, in order; absent when there were none"`
	Outcome      string             `json:"outcome,omitempty" desc:"passed or failed: its latest cycle's outcome; empty while that cycle is open, even when an earlier cycle was judged"`
	OpenCycle    int                `json:"open_cycle,omitempty" desc:"The cycle that is open, when its latest cycle is open after an earlier one was judged; absent otherwise"`
	PriorCycle   int                `json:"prior_cycle,omitempty" desc:"The newest earlier cycle that was judged, when open_cycle is set"`
	PriorOutcome string             `json:"prior_outcome,omitempty" desc:"That cycle's outcome, passed or failed, when open_cycle is set"`
	Questions    []RunCheckQuestion `json:"questions,omitempty" desc:"The questions it raised when it still failed with no fix cycle left, in order; absent when none"`
}

// RunCheckQuestion is one question a capped review raised for a check.
type RunCheckQuestion struct {
	Cycle          int    `json:"cycle" desc:"The cycle that ended with no fix cycle left"`
	Finding        string `json:"finding,omitempty" desc:"The surviving confirmed finding it is about; empty for one on the check's failure itself"`
	Leaf           string `json:"leaf,omitempty" desc:"The leaf node the finding names, which the question blocks; empty for a question in the change set at large"`
	LeafHandle     string `json:"leaf_handle,omitempty" desc:"That leaf's handle"`
	Question       string `json:"question" desc:"The question node's ID"`
	QuestionHandle string `json:"question_handle" desc:"The question node's handle"`
}

// RunCheckRound is one findings round of a review check.
type RunCheckRound struct {
	Round     int          `json:"round" desc:"The round's number in the check, from 1"`
	Cycle     int          `json:"cycle" desc:"The check's cycle it was recorded in, from 1"`
	At        string       `json:"at" desc:"When its findings were journalled, RFC 3339 UTC"`
	Findings  []RunFinding `json:"findings" desc:"The check's findings, in order, each with its latest verdict; empty for a clean check"`
	Verdicts  []RunVerdict `json:"verdicts,omitempty" desc:"The verdicts journalled on the round, in order; absent before any"`
	Confirmed int          `json:"confirmed" desc:"How many findings the latest verdict confirms"`
}

// RunCheckCycle is one cycle a check ran in.
type RunCheckCycle struct {
	Cycle     int               `json:"cycle" desc:"The cycle, from 1"`
	Sessions  []RunCheckSession `json:"sessions" desc:"The check's sessions in the cycle, in the order they started"`
	Outcome   string            `json:"outcome,omitempty" desc:"passed or failed; empty while the cycle is open for the check"`
	Detail    string            `json:"detail,omitempty" desc:"Why the check ended the cycle as it did"`
	At        string            `json:"at,omitempty" desc:"When the check's outcome in the cycle was journalled, RFC 3339 UTC"`
	Decisions []string          `json:"decisions,omitempty" desc:"For a fresh cycle past the cap, the Decisions (handles) answering the questions the cycle before raised; absent otherwise"`
	Kind      string            `json:"kind,omitempty" desc:"answered for a fresh cycle past the cap, owner for the cycle fixing the owner's request for changes; absent for a cycle within the cap"`
}

// RunCheckSession is one session a checks step started for a check.
type RunCheckSession struct {
	Role    string `json:"role" desc:"check"`
	Session string `json:"session" desc:"The session's ID"`
	Model   string `json:"model,omitempty" desc:"The model it runs; empty for the agent's default"`
	Outcome string `json:"outcome,omitempty" desc:"pass or fail once its work was judged; empty while it runs"`
	Detail  string `json:"detail,omitempty" desc:"Why it was judged so: the session's exit, or that its agent said done"`
	At      string `json:"at" desc:"When it started, RFC 3339 UTC"`
}

// RunLeaf is one leaf a leaves step conducted.
type RunLeaf struct {
	Step           string           `json:"step" desc:"The leaves step, or the checks step whose fixer fixed it"`
	Leaf           string           `json:"leaf" desc:"The leaf capability's ID"`
	Handle         string           `json:"handle" desc:"The leaf capability's handle"`
	Sessions       []RunLeafSession `json:"sessions" desc:"The leaf's sessions, in the order they started"`
	Attempts       []RunLeafAttempt `json:"attempts" desc:"The judgements of its sessions' work, in order"`
	Findings       []RunLeafFinding `json:"findings,omitempty" desc:"The findings its verifiers journalled with run findings --leaf, in order; absent when there were none"`
	Outcome        string           `json:"outcome,omitempty" desc:"passed, failed or skipped; empty while the leaf is open, and always on a checks step"`
	Question       string           `json:"question,omitempty" desc:"The question node a failed leaf raised, by ID, if any"`
	QuestionHandle string           `json:"question_handle,omitempty" desc:"That question node's handle"`
	Round          int              `json:"round" desc:"The leaf's current round of attempts, from 1"`
	Rounds         []RunLeafRound   `json:"rounds,omitempty" desc:"Its fresh rounds, each given once a Decision answered the question its failure raised; absent in round 1"`
	Commit         string           `json:"commit,omitempty" desc:"A passed leaf's Evidence commit: its verifier's Evidence entries, committed alone"`
	Detail         string           `json:"detail,omitempty" desc:"Why the leaf ended as it did"`
}

// RunLeafRound is one fresh round a failed leaf was given.
type RunLeafRound struct {
	Round          int    `json:"round" desc:"The round, from 2"`
	After          int    `json:"after" desc:"The attempts before it: its first attempt is after + 1"`
	Question       string `json:"question" desc:"The question node the leaf's failure raised, by ID"`
	QuestionHandle string `json:"question_handle,omitempty" desc:"That question node's handle"`
	Decision       string `json:"decision" desc:"The handle of the Decision that answers it"`
	At             string `json:"at" desc:"When the round began, RFC 3339 UTC"`
}

// RunLeafSession is one session a leaves step started on a leaf.
type RunLeafSession struct {
	Role    string `json:"role" desc:"builder or verifier; on a checks step fixer or verifier"`
	Attempt int    `json:"attempt" desc:"The leaf's attempt it works on, from 1; on a checks step, one per try of a cycle's fix: cycle 2's first try is 2, and each cycle has 3 tries"`
	Cycle   int    `json:"cycle,omitempty" desc:"On a checks step, the review cycle of the fix it belongs to"`
	Session string `json:"session" desc:"The session's ID"`
	Model   string `json:"model,omitempty" desc:"The model it runs; empty for the agent's default"`
	Commit  string `json:"commit,omitempty" desc:"A builder's or fixer's: the worktree's HEAD when it started; a verifier's: the commit it checks"`
	At      string `json:"at" desc:"When it started, RFC 3339 UTC"`
}

// RunLeafAttempt is one judgement of a leaf session's work.
type RunLeafAttempt struct {
	Role    string `json:"role" desc:"builder or verifier; on a checks step fixer or verifier"`
	Attempt int    `json:"attempt" desc:"The leaf's attempt, from 1; on a checks step, its session's attempt, one per try of a cycle's fix"`
	Cycle   int    `json:"cycle,omitempty" desc:"On a checks step, the review cycle of the fix it belongs to"`
	Outcome string `json:"outcome" desc:"pass or fail"`
	Commit  string `json:"commit,omitempty" desc:"A passed builder's or fixer's commit naming its leaf; a passed verifier's Evidence commit"`
	Detail  string `json:"detail,omitempty" desc:"Why: the session's exit, that its agent said done, or that a builder left no commit naming its leaf"`
	At      string `json:"at" desc:"When it was judged, RFC 3339 UTC"`
}

// RunLeafFinding is one finding a verifier journalled on a leaf's attempt.
type RunLeafFinding struct {
	Attempt int    `json:"attempt" desc:"The leaf's attempt it is on, from 1"`
	ID      string `json:"id" desc:"The finding's ID in the attempt: f1, f2, ..."`
	Reason  string `json:"reason" desc:"What the verifier found wrong"`
	At      string `json:"at" desc:"When it was journalled, RFC 3339 UTC"`
}

// RunWaiting is one of the plan's change sets run gate named that waits
// proposed for the owner's accept.
type RunWaiting struct {
	Change string `json:"change" desc:"The change set's ID"`
	Title  string `json:"title" desc:"Its title"`
	Next   string `json:"next" desc:"The owner's command to accept it: rhiza accept --rule owner --change <id>"`
}

// RunRound is one plan review round: a reviewer's findings, each with the
// critic's latest verdict.
type RunRound struct {
	Round     int          `json:"round" desc:"The round's number, from 1"`
	At        string       `json:"at" desc:"When its findings were journalled, RFC 3339 UTC"`
	Findings  []RunFinding `json:"findings" desc:"The reviewer's findings, in order; empty for a clean review"`
	Confirmed int          `json:"confirmed" desc:"How many findings the critic's latest verdict confirms"`
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

// RunMessage is one message journaled in a run: a SendMessage one of its
// sessions sent.
type RunMessage struct {
	At         string `json:"at" desc:"When it was sent, RFC 3339 UTC"`
	Step       string `json:"step" desc:"The sender's step"`
	Session    string `json:"session" desc:"The sender's session"`
	Actor      string `json:"actor" desc:"The sender's actor handle"`
	To         string `json:"to" desc:"The recipient, as the sender named it"`
	Text       string `json:"text" desc:"The message; a structured one is its raw JSON"`
	Summary    string `json:"summary,omitempty" desc:"The sender's summary of it; absent when it gave none"`
	Structured bool   `json:"structured,omitempty" desc:"True if the message was not a string, so text holds its raw JSON"`
}

// RunSession is one session of a run, under the step it first appeared
// in, with the latest reading of its token use.
type RunSession struct {
	Step     string           `json:"step" desc:"The step the session first appeared under"`
	Session  string           `json:"session" desc:"The session ID"`
	Tokens   *SessionTokenUse `json:"tokens,omitempty" desc:"The latest reading of its token use; absent when none could be read, and then left out of the totals"`
	NoTokens string           `json:"no_tokens,omitempty" desc:"Why no token use is shown, when tokens is absent: still running, or the transcript is unreadable with its warning"`
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

// RunStalled is a session of an unfinished run that works but has not
// reported for more than 20 minutes.
type RunStalled struct {
	Step    string `json:"step" desc:"The step whose session is stalled"`
	Session string `json:"session" desc:"The stalled session"`
	Since   string `json:"since" desc:"When the session last reported, RFC 3339 UTC"`
}

// runStalledNow is the first stalled session at now of the current step of
// a live, unfinished run, of those runStalledAll lists.
func runStalledNow(st run.State, live bool, sessions *session.Store, now time.Time) *RunStalled {
	if all := runStalledAll(st, live, sessions, now); len(all) > 0 {
		return &all[0]
	}
	return nil
}

// runStalledAll is every stalled session at now of the current step of a
// live, unfinished run: an agent step's session, then the latest session of
// each open leaf, of each check and of each unjudged change-set fixer of a
// leaves or checks step. Only a live run's sessions stall: a stopped run has
// no process holding it. A record that cannot be read is not stalled.
func runStalledAll(st run.State, live bool, sessions *session.Store, now time.Time) []RunStalled {
	if !live || st.Ended || sessions == nil {
		return nil
	}
	var out []RunStalled
	for _, sid := range stepSessions(st) {
		rec, err := sessions.Read(sid)
		if err != nil || !session.Stalled(rec, now) {
			continue
		}
		rep, _ := session.Resolve(rec)
		out = append(out, RunStalled{Step: st.Current, Session: sid, Since: rfc3339(rep.At)})
	}
	return out
}

// stepSessions are the sessions of the run's current step that can still
// be working, in the order runStalledAll checks them. A session whose pane
// is closed is gone.
func stepSessions(st run.State) []string {
	var out []string
	if sid := st.Session[st.Current]; sid != "" {
		out = append(out, sid)
	}
	// add takes the latest of a leaf's or a check's sessions, by journal
	// position.
	type started struct {
		seq    int
		id     string
		closed bool
	}
	add := func(all []started) {
		if len(all) == 0 {
			return
		}
		last := slices.MaxFunc(all, func(a, b started) int { return cmp.Compare(a.seq, b.seq) })
		if !last.closed && last.id != "" {
			out = append(out, last.id)
		}
	}
	for _, l := range st.Leaves {
		if l.Step != st.Current || l.Outcome != "" {
			continue
		}
		var all []started
		for _, s := range l.Sessions {
			all = append(all, started{s.Seq, s.Session, s.PaneClosed})
		}
		add(all)
	}
	for _, c := range st.Checks {
		if c.Step != st.Current {
			continue
		}
		var all []started
		for _, s := range c.Sessions {
			all = append(all, started{s.Seq, s.Session, s.PaneClosed})
		}
		add(all)
	}
	for _, f := range st.Fixes {
		if f.Step == st.Current && !f.Judged && !f.PaneClosed && f.Session != "" {
			out = append(out, f.Session)
		}
	}
	return out
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
	if err := requireProject(ctx); err != nil {
		return nil, err
	}
	if r.Milestone && strings.TrimSpace(r.New) == "" {
		return nil, Usagef("--milestone goes with --new").WithNext(runStartNext)
	}
	if strings.TrimSpace(r.Request) != "" && strings.TrimSpace(r.New) == "" {
		return nil, Usagef("--request goes with --new").WithNext(runStartNext)
	}
	path, def, data, err := resolveWorkflow(ctx.Root, r.Workflow)
	if err != nil {
		return nil, err
	}
	if r.Autonomous && def.Name == planWorkflow && r.Milestone {
		return nil, milestoneNeedsOwner(runStartNext)
	}
	if err := checkAgentSteps(ctx, def, r.PermissionMode, r.Agent, r.Autonomous, runStartNext); err != nil {
		return nil, err
	}
	// A request's run takes the ID the daemon drew for it: check it before
	// a draft is opened, so a refused start leaves none.
	var reqRun, reqMilestone string
	if q := strings.TrimSpace(r.Request); q != "" {
		common, err := gitCommonDir(ctx)
		if err != nil {
			return nil, err
		}
		req, err := requestRun(common, q, run.NewStore(run.Dir(common)))
		if err != nil {
			return nil, err
		}
		// A milestone request's plan plans the milestone, --milestone or not.
		if r.Autonomous && def.Name == planWorkflow && req.Kind == daemon.KindMilestone {
			return nil, milestoneNeedsOwner(runStartNext)
		}
		reqRun, reqMilestone = req.Run, req.Milestone
	}
	var t changeWorktree
	if strings.TrimSpace(r.New) != "" {
		t, err = newPlanChange(ctx, r, def)
	} else {
		t, err = resolveChangeWorktree(ctx, r.Change, "no run can work on it")
	}
	if err != nil {
		return nil, err
	}
	if r.Autonomous && def.Name == planWorkflow && planIsMilestone(ctx, t.cs.String()) {
		return nil, milestoneNeedsOwner(runStartNext)
	}
	if def.Name == shipChangeWorkflow {
		if err := checkShipLeaves(t); err != nil {
			return nil, err
		}
	}
	rs := run.NewStore(run.Dir(t.common))
	rid := reqRun
	if rid == "" {
		if rid, err = newRunID(ctx, rs); err != nil {
			return nil, err
		}
	}
	started := run.Started{
		Workflow: def.Name, WorkflowPath: path, Definition: string(data), Includes: def.Includes,
		Change: t.cs.String(), Title: t.title, Worktree: t.worktree, Branch: t.branch,
		PermissionMode: r.PermissionMode, Agent: r.Agent, Autonomous: r.Autonomous, Milestone: reqMilestone,
		MilestonePlan: r.Milestone,
	}
	if ctx.DryRun {
		if err := refuseLiveRun(rs, started.Change, ""); err != nil {
			return nil, err
		}
		return RunResponse{Run: RunInfo{
			ID: rid, Workflow: def.Name, Change: started.Change, Title: t.title, Worktree: t.worktree,
			Autonomous: r.Autonomous, Current: def.Start, Steps: []RunStep{}, StartedAt: rfc3339(ctx.Now()), Amendments: []RunAmendment{},
			Messages: []RunMessage{}, Sessions: []RunSession{},
		}, DryRun: true, Queued: !r.Foreground}, nil
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
	if !r.Foreground {
		return queueRunStart(rs, t, def, r, started)
	}
	if _, err := t.open(ctx); err != nil {
		return nil, err
	}
	claimStartLog(ctx, t.common, rid)
	if err := startRun(rs, rid, started, ctx.Now(), def); err != nil {
		return nil, err
	}
	return runEngine(ctx, rs, t.common, rid, started.Change, t.worktree, r.PermissionMode, r.Agent, started.Autonomous, r.NoWait)
}

// milestoneNeedsOwner is the usage error for an autonomous plan run that
// plans a milestone: milestone planning always needs the owner.
func milestoneNeedsOwner(next string) error {
	return Usagef("--autonomous is refused on a plan run that plans a milestone: milestone planning needs the owner, even in autonomous mode").
		WithNext(next)
}

// planIsMilestone reports whether the change set is a record change set,
// which is where a plan run plans a milestone.
func planIsMilestone(ctx *Ctx, change string) bool {
	f := newInboxFold(ctx)
	defer f.close()
	c, ok := f.changeSet(change)
	return ok && c.Open.Kind == changeset.KindRecord
}

// queueRunStart queues the run for the daemon, which starts it headless
// within its maxRuns as rhiza run start --foreground. The start is the
// owner's: whatever change set run start accepts is queued, with no
// milestone, state or leaf filter of the daemon's picks. A change set
// with a live run or a queued start is refused.
func queueRunStart(rs *run.Store, t changeWorktree, def *workflow.Definition, r RunStartRequest, started run.Started) (any, error) {
	change := started.Change
	if err := refuseLiveRun(rs, change, ""); err != nil {
		return nil, err
	}
	err := daemon.UpdateStarts(t.common, func(starts []daemon.Start) ([]daemon.Start, error) {
		if slices.ContainsFunc(starts, func(s daemon.Start) bool { return s.Change == change }) {
			return nil, Usagef("change set %s already has a run queued for the daemon", change).WithNext("rhiza daemon status")
		}
		return append(starts, daemon.Start{
			Change: change, Workflow: r.Workflow, PermissionMode: r.PermissionMode, Agent: r.Agent,
			NoWait: r.NoWait, Autonomous: r.Autonomous,
		}), nil
	})
	if err != nil {
		var e *Error
		if errors.As(err, &e) {
			return nil, err
		}
		return nil, Errorf("queue the run: %v", err)
	}
	h, up := daemon.Running(t.common)
	resp := RunResponse{Run: RunInfo{
		Workflow: def.Name, Change: change, Title: t.title, Worktree: t.worktree, Autonomous: r.Autonomous,
		Steps: []RunStep{}, Amendments: []RunAmendment{}, Messages: []RunMessage{}, Sessions: []RunSession{},
	}, Queued: true}
	if up {
		resp.Daemon = h.String()
	}
	return resp, nil
}

// newPlanChange opens the draft change set run start --new plans in (a
// record with --milestone, else an intention) and derives its worktree.
// It refuses a change set argument and any workflow but plan first, and in
// a real run also the agent-session and can-run refusals, so a refused
// start leaves no draft behind. A dry run writes nothing: the change set
// is the ID the real run would draw, with the worktree its title gives.
func newPlanChange(ctx *Ctx, r RunStartRequest, def *workflow.Definition) (changeWorktree, error) {
	var t changeWorktree
	if strings.TrimSpace(r.Change) != "" || strings.TrimSpace(ctx.ChangeRef) != "" || !ctx.ChangeSet.IsZero() {
		return t, Usagef("--new opens the change set to plan: it takes no change set").WithNext(runStartNext)
	}
	if def.Name != planWorkflow {
		return t, Usagef("--new starts a plan run: workflow %q is not %s", def.Name, planWorkflow).
			WithNext("rhiza run start --workflow plan --new \"<intent>\"")
	}
	kind := changeset.KindIntention
	if r.Milestone {
		kind = changeset.KindRecord
	}
	open := changeset.Open{Kind: kind, Intent: strings.TrimSpace(r.New)}
	if err := open.Validate(); err != nil {
		return t, Usagef("%s", err.Error()).WithNext(runStartNext)
	}
	if !ctx.DryRun {
		if ctx.Getenv != nil {
			if sid := strings.TrimSpace(ctx.Getenv(envSession)); sid != "" {
				return t, Usagef("run start refused: this is agent session %s, and no agent session starts a run", sid).
					WithNext("ask the owner to run it outside an agent session")
			}
		}
		if err := canRun(ctx, def); err != nil {
			return t, err
		}
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return t, err
	}
	current, err := w.current()
	if err != nil {
		return t, err
	}
	actor, forID, err := w.resolveActor(current, ctx)
	if err != nil {
		return t, err
	}
	cs, _, err := w.openChangeSet(open, actor, forID)
	if err != nil {
		return t, err
	}
	if ctx.DryRun {
		common, err := gitCommonDir(ctx)
		if err != nil {
			return t, err
		}
		prefix, err := ctx.Git("rev-parse", "--show-prefix")
		if err != nil {
			return t, Errorf("%v", err)
		}
		t = changeWorktree{w: w, cs: cs, title: open.DisplayTitle(), common: common, prefix: strings.TrimSpace(prefix)}
		mainRoot := filepath.Dir(common)
		t.repo = filepath.Base(mainRoot)
		slug := session.ChangeSetSlug(w.changeSets(), cs)
		t.slug = slug
		t.branch = session.Branch(slug)
		t.worktree = session.Worktree(mainRoot, t.repo, slug)
		return t, nil
	}
	if err := w.setCurrent(cs); err != nil {
		return t, err
	}
	nctx := *ctx
	nctx.ChangeSet, nctx.ChangeRef = cs, ""
	return resolveChangeWorktree(&nctx, "", "no run can work on it")
}

// planWorkflow is the workflow run start --new opens a change set for.
const planWorkflow = "plan"

// startRun journals a new run under the change set's lock, after checking
// that the change set has no other live run, so two starts never both begin.
func startRun(rs *run.Store, rid string, started run.Started, now time.Time, def *workflow.Definition) error {
	unlock, err := rs.LockChange(started.Change)
	if err != nil {
		return Errorf("lock change set %s: %v", started.Change, err)
	}
	defer unlock()
	if err := refuseLiveRun(rs, started.Change, ""); err != nil {
		return err
	}
	if err := run.Start(rs, rid, started, now, def); err != nil {
		return Errorf("start run %s: %v", rid, err)
	}
	return nil
}

// refuseLiveRun is a usage error (exit 2) naming a run of the change set, other than
// except, whose journal has not finished, an interrupted one included.
func refuseLiveRun(rs *run.Store, change, except string) error {
	j, ok, err := rs.LiveOf(change, except)
	if err != nil {
		return Errorf("read runs: %v", err)
	}
	if !ok {
		return nil
	}
	return Usagef("change set %s already has a live run, %s: its journal has not finished", change, j.ID).
		WithNext("rhiza run resume " + j.ID)
}

// shipChangeWorkflow is the workflow that ships a change set through
// build and review, both of which work on its leaves.
const shipChangeWorkflow = "ship-change"

// checkShipLeaves refuses a change set with no leaves: ship-change would
// build nothing, so it is planned first.
func checkShipLeaves(t changeWorktree) error {
	c, ok := t.w.changeSet(t.cs)
	if ok && len(changeLeaves(t.w.model(t.cs), c.Ops)) > 0 {
		return nil
	}
	return Conflictf("change set %s has no leaves: ship-change builds and reviews its leaves, so plan it first", t.cs).
		WithNext("rhiza run start --workflow plan")
}

// refuseAgentSession refuses a run-starting path (verb is start or resume)
// inside an agent session, naming the session. No agent session starts a run.
func refuseAgentSession(ctx *Ctx, verb string) error {
	if ctx.Getenv == nil {
		return nil
	}
	if sid := strings.TrimSpace(ctx.Getenv(envSession)); sid != "" {
		return Usagef("run %s refused: this is agent session %s, and no agent session starts a run", verb, sid).
			WithNext("ask the owner to run it outside an agent session")
	}
	return nil
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
	if st.Abandoned() {
		return nil, Usagef("run %s was abandoned: %s", j.ID, st.Detail).
			WithNext("rhiza run start " + st.Change)
	}
	if st.Autonomous && st.Workflow == planWorkflow && planIsMilestone(ctx, st.Change) {
		return nil, milestoneNeedsOwner("rhiza run status " + j.ID)
	}
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
	// An accepted change set's run is resumed by the daemon; one whose
	// change set has since left the in-flight states is refused here too.
	// A draft or proposed change set (a plan run's) and one the log does
	// not have are let through.
	if !st.Ended {
		// A plan run lands its record during the interview, so a landed
		// change set does not stop a plan run; rejected, withdrawn and
		// superseded ones still do.
		state, ok := changeStateOf(ctx)(st.Change)
		if ok && st.Workflow == planWorkflow && changeset.State(state) == changeset.StateLanded {
			ok = false
		}
		// A non-plan run that reached its land step goes on past it: what
		// is left (wait-main, close-out) happens after the landing.
		if ok && st.Workflow != planWorkflow && changeset.State(state) == changeset.StateLanded && st.Visits["land"] > 0 {
			ok = false
		}
		if ok && changeset.State(state).Terminal() {
			return nil, Conflictf("run %s is not resumed: its change set %s is %s, no longer accepted", j.ID, st.Change, state).
				WithNext("rhiza run abandon " + j.ID)
		}
	}
	def, err := st.ParseDefinition()
	if err != nil {
		return nil, Errorf("run %s: its workflow no longer parses: %v", j.ID, err)
	}
	// Start checked the agent steps; check again, since the journal is a
	// file anyone can edit.
	if err := checkAgentSteps(ctx, def, st.PermissionMode, st.Agent, st.Autonomous, "rhiza run status "+j.ID); err != nil {
		return nil, err
	}
	if retryStep != "" {
		if _, ok := def.Steps[retryStep]; !ok {
			return nil, Errorf("run %s: step %q is not in workflow %q", j.ID, retryStep, def.Name).WithNext("rhiza run status " + j.ID)
		}
	}
	if err := refuseLiveRun(rs, st.Change, j.ID); err != nil {
		return nil, err
	}
	if ctx.DryRun {
		if retryStep != "" {
			// Fold the retried entry the engine would append, without
			// writing it, to show the run as the retry would leave it.
			j.Entries = append(j.Entries, run.Entry{Type: run.TypeRetried, At: ctx.Now(), Step: retryStep, Detail: st.Detail})
			st = j.State()
		}
		return RunResponse{Run: runInfo(ctx, st, false, session.NewStore(session.Dir(common))), DryRun: true, Retry: retryStep, Queued: !r.Foreground}, nil
	}
	if err := refuseAgentSession(ctx, "resume"); err != nil {
		return nil, err
	}
	if err := canRun(ctx, def); err != nil {
		return nil, err
	}
	if !r.Foreground {
		resp, err := queueContinue(ctx, common, j.ID, r.Retry, r.NoWait)
		if err != nil {
			return nil, err
		}
		resp.Run = runInfo(ctx, st, false, session.NewStore(session.Dir(common)))
		resp.Retry = retryStep
		return resp, nil
	}
	resp, err := runEngineRetry(ctx, rs, common, j.ID, st.Change, st.Worktree, st.PermissionMode, st.Agent, st.Autonomous, r.NoWait, r.Retry, nil)
	if err != nil {
		return nil, err
	}
	resp.Retry = retryStep
	return resp, nil
}

// queueContinue queues the continuation of run id for the daemon, which
// drives it headless within its maxRuns as rhiza run resume --foreground.
// It is the owner's action, not a pick: no state or milestone filter
// applies. A run with a continuation already queued is refused. The response
// holds only Queued and Daemon; the caller adds the run.
func queueContinue(ctx *Ctx, common, id string, retry, noWait bool) (RunResponse, error) {
	err := daemon.UpdateContinues(common, func(cs []daemon.Continue) ([]daemon.Continue, error) {
		if slices.ContainsFunc(cs, func(c daemon.Continue) bool { return c.Run == id }) {
			return nil, Usagef("run %s already has a continuation queued for the daemon", id).WithNext("rhiza daemon status")
		}
		return append(cs, daemon.Continue{Run: id, Retry: retry, NoWait: noWait, At: ctx.Now().UTC()}), nil
	})
	if err != nil {
		var e *Error
		if errors.As(err, &e) {
			return RunResponse{}, err
		}
		return RunResponse{}, Errorf("queue the continuation: %v", err)
	}
	resp := RunResponse{Queued: true}
	if h, up := daemon.Running(common); up {
		resp.Daemon = h.String()
	}
	return resp, nil
}

// runRunAbandon journals an unfinished run as failed with the detail
// "abandoned by <actor>", under the run's lock so a run going in another
// process is refused. It removes the run's clean worktree and local branch
// unless another run still uses them; the change set then has no live run,
// so a new one can start. A failed run can be
// abandoned too, which clears it from the inbox; a succeeded or already
// abandoned run is refused.
func runRunAbandon(ctx *Ctx, req any) (any, error) {
	r := req.(RunAbandonRequest)
	ref := strings.TrimSpace(r.Run)
	if ref == "" {
		return nil, Usagef("run abandon needs a run ID").WithNext(runStatusNext)
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
	sessions := session.NewStore(session.Dir(common))
	refuseEnded := func(st run.State) error {
		if st.Ended && (st.Result != run.Failed || st.Abandoned()) {
			return Usagef("run %s is finished: it %s", j.ID, st.Result).WithNext("rhiza run status " + j.ID)
		}
		return nil
	}
	if err := refuseEnded(j.State()); err != nil {
		return nil, err
	}
	from, err := messageSender(ctx)
	if err != nil {
		return nil, err
	}
	detail := run.AbandonedPrefix + from.Handle
	if ctx.DryRun {
		j.Entries = append(j.Entries, run.Entry{Type: run.TypeFinished, At: ctx.Now(), End: run.Failed, Detail: detail})
		return RunAbandonResponse{Run: runInfo(ctx, j.State(), false, sessions), DryRun: true}, nil
	}
	unlock, ok, err := rs.TryLock(j.ID)
	if err != nil {
		return nil, Errorf("lock run %s: %v", j.ID, err)
	}
	if !ok {
		return nil, Conflictf("run %s is running in another process", j.ID).WithNext("rhiza run status " + j.ID)
	}
	defer unlock()
	// Read again under the lock: the holder may have finished it.
	if j, err = rs.Read(j.ID); err != nil {
		return nil, err
	}
	if err := refuseEnded(j.State()); err != nil {
		return nil, err
	}
	e := run.Entry{Type: run.TypeFinished, At: ctx.Now(), End: run.Failed, Detail: detail}
	if err := rs.Append(j.ID, e); err != nil {
		return nil, Errorf("journal the abandoning of run %s: %v", j.ID, err).WithNext("rhiza run status " + j.ID)
	}
	j.Entries = append(j.Entries, e)
	st := j.State()
	return RunAbandonResponse{Run: runInfo(ctx, st, false, sessions), Cleanup: abandonCleanup(ctx, common, rs, j.ID, st)}, nil
}

// worktreeDirty reports whether git status --porcelain -uall output holds
// a change worth keeping the worktree for. Untracked log entries that the
// main checkout holds byte for byte are copies Rhiza made, not work.
func worktreeDirty(porcelain, worktree, main string) bool {
	for l := range strings.SplitSeq(porcelain, "\n") {
		if strings.TrimSpace(l) == "" {
			continue
		}
		rel, ok := strings.CutPrefix(l, "?? ")
		if !ok || !strings.Contains(rel, ".rhiza/log/") {
			return true
		}
		rel = strings.Trim(rel, `"`)
		a, err1 := os.ReadFile(filepath.Join(worktree, rel))
		b, err2 := os.ReadFile(filepath.Join(main, rel))
		if err1 != nil || err2 != nil || !bytes.Equal(a, b) {
			return true
		}
	}
	return false
}

// runChangeSlug is the unique slug of the change set the run st works on.
func runChangeSlug(ctx *Ctx, st run.State) (string, bool) {
	cs, err := id.Parse(st.Change)
	if err != nil {
		return "", false
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return "", false
	}
	return session.ChangeSetSlug(w.changeSets(), cs), true
}

// abandonCleanup removes the abandoned run's worktree without force and
// deletes its local cs/<slug> branch, keeping both when the worktree is
// dirty or another run still uses them, and the branch when it holds
// commits not merged into its upstream or HEAD. The remote branch and any PR are
// left alone. Failures are reported as kept, never as errors: the run is
// already abandoned.
func abandonCleanup(ctx *Ctx, common string, rs *run.Store, id string, st run.State) []RunAbandonStep {
	if st.Worktree == "" || st.Title == "" || ctx.Exec == nil {
		return nil
	}
	main := filepath.Dir(common)
	if samePath(st.Worktree, main) {
		return nil
	}
	slug, ok := runChangeSlug(ctx, st)
	if !ok {
		return nil
	}
	branch := session.Branch(slug)
	git := func(args ...string) ([]byte, error) {
		return ctx.Exec("git", append([]string{"-C", main}, args...), main, nil)
	}
	var out []RunAbandonStep
	kept := func(step, why string) {
		out = append(out, RunAbandonStep{Step: step, Status: "kept", Detail: why})
	}
	if other, err := worktreeUser(rs, id, st); err != nil {
		kept("remove-worktree", "cannot list the runs: "+err.Error())
	} else if other != "" {
		kept("remove-worktree", "run "+other+" still uses it")
	} else if fi, err := os.Stat(st.Worktree); err != nil || !fi.IsDir() {
		git("worktree", "prune")
		out = append(out, RunAbandonStep{Step: "remove-worktree", Status: "already-gone"})
	} else {
		dirty, err := ctx.Exec("git", []string{"-C", st.Worktree, "status", "--porcelain", "-uall"}, st.Worktree, nil)
		switch {
		case err != nil:
			kept("remove-worktree", "cannot read the worktree's status: "+err.Error())
		case worktreeDirty(string(dirty), st.Worktree, main):
			kept("remove-worktree", "the worktree is dirty")
		default:
			// Only verified copies of log entries can remain, which git
			// counts as untracked files and so needs --force to drop.
			args := []string{"worktree", "remove"}
			if strings.TrimSpace(string(dirty)) != "" {
				args = append(args, "--force")
			}
			if _, err := git(append(args, st.Worktree)...); err != nil {
				kept("remove-worktree", err.Error())
			} else {
				out = append(out, RunAbandonStep{Step: "remove-worktree", Status: "done"})
			}
		}
	}
	if out[0].Status == "kept" {
		kept("delete-branch", "the worktree is kept")
		return out
	}
	if _, err := git("show-ref", "--verify", "--quiet", "refs/heads/"+branch); err != nil {
		return append(out, RunAbandonStep{Step: "delete-branch", Status: "already-gone"})
	}
	// -d, not -D: git refuses a branch with commits on neither its
	// upstream nor HEAD, such as an earlier build run's unpushed work.
	if _, err := git("branch", "-d", branch); err != nil {
		if strings.Contains(err.Error(), "not fully merged") {
			kept("delete-branch", "it has commits merged into neither its upstream nor HEAD")
		} else {
			kept("delete-branch", err.Error())
		}
		return out
	}
	return append(out, RunAbandonStep{Step: "delete-branch", Status: "done"})
}

// worktreeUser returns another run that still needs the abandoned run's
// worktree and branch: a live run of the same change set, or any run
// started later in the same worktree. It returns "" when there is none.
func worktreeUser(rs *run.Store, id string, st run.State) (string, error) {
	if live, ok, err := rs.LiveOf(st.Change, id); err != nil {
		return "", err
	} else if ok {
		return live.ID, nil
	}
	runs, err := rs.List()
	if err != nil {
		return "", err
	}
	for _, j := range runs {
		if j.ID == id {
			continue
		}
		o := j.State()
		if o.Worktree != "" && samePath(o.Worktree, st.Worktree) && o.StartedAt.After(st.StartedAt) {
			return j.ID, nil
		}
	}
	return "", nil
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
	st := j.State()
	info := runInfo(ctx, st, rs.Live(j.ID), session.NewStore(session.Dir(common)))
	if info.SignoffWaiting {
		if ss, err := readShipState(filepath.Join(common, filepath.FromSlash(shipDir), st.Change+".json")); err == nil {
			info.PR = ss.URL
		}
	}
	if len(st.Gated) > 0 {
		if w, err := openCtxWorkspace(ctx); err == nil {
			info.Waiting = runWaiting(st.Gated, w.changeSets())
		}
	}
	return info, nil
}

// runWaiting is the change sets among gated, the latest gated entry's,
// that are proposed now in sets, in gated's order.
func runWaiting(gated []string, sets []changeset.ChangeSet) []RunWaiting {
	var out []RunWaiting
	for _, g := range gated {
		for _, c := range sets {
			if c.ID.String() == g && c.State == changeset.StateProposed {
				out = append(out, RunWaiting{Change: g, Title: c.Open.DisplayTitle(), Next: ownerAcceptNext(g)})
			}
		}
	}
	return out
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

// personalWorkflowsDir is the user's own workflows directory:
// $XDG_CONFIG_HOME/rhiza/workflows, else ~/.config/rhiza/workflows. It is
// empty when neither can be known.
func personalWorkflowsDir() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "rhiza", "workflows")
	}
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return filepath.Join(h, ".config", "rhiza", "workflows")
	}
	return ""
}

// lookupWorkflow is the one resolver of a bare workflow name, shared by run
// start, workflow steps naming another workflow and the workflow commands:
// the personal file, else the repo's .rhiza/workflows/<name>.toml, else the
// built-in. It returns the path to record (the file, or builtin:<name>) and
// the source; found is false when no layer has the name.
func lookupWorkflow(root, name string) (path string, data []byte, found bool, err error) {
	var files []string
	if d := personalWorkflowsDir(); d != "" {
		files = append(files, filepath.Join(d, name+".toml"))
	}
	files = append(files, filepath.Join(root, filepath.FromSlash(workflowsDir), name+".toml"))
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err == nil {
			return f, data, true, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return f, nil, false, err
		}
	}
	if data, ok := workflow.Builtin(name); ok {
		return builtinPrefix + name, data, true, nil
	}
	return "", nil, false, nil
}

// resolveWorkflow finds the workflow ref names and returns the path to
// record in the journal (the file, or builtin:<name> for an embedded
// one), its definition and its source. A bare name goes through
// lookupWorkflow: personal file, repo file, built-in.
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
		found, data, ok, lerr := lookupWorkflow(root, ref)
		if lerr != nil {
			return "", nil, nil, Usagef("workflow %q: %v", ref, lerr).WithNext(runStartNext)
		}
		if !ok {
			return "", nil, nil, Usagef("no workflow %q: %s does not exist and there is no built-in of that name (built-in: %s)",
				ref, path, strings.Join(workflow.BuiltinNames(), ", ")).WithNext(runStartNext)
		}
		def, perr := workflow.ParseWith(ref, data, workflowResolver(root))
		if perr != nil {
			if strings.HasPrefix(found, builtinPrefix) {
				return "", nil, nil, Errorf("built-in workflow %s: %v", ref, perr)
			}
			return "", nil, nil, Usagef("workflow %s: %v", found, perr).WithNext(runStartNext)
		}
		return found, def, data, nil
	}
	def, data, err := loadWorkflow(root, path, ref)
	if err != nil {
		return "", nil, nil, err
	}
	return path, def, data, nil
}

// workflowResolver resolves the workflow a workflow step names as run
// start resolves a bare name, through lookupWorkflow.
func workflowResolver(root string) workflow.Resolver {
	return func(name string) ([]byte, error) {
		_, data, ok, err := lookupWorkflow(root, name)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("there is no such workflow in %s, your personal workflows or built in", workflowsDir)
		}
		return data, nil
	}
}

// loadWorkflow reads and validates the workflow at path. A missing file
// and an invalid workflow are usage errors carrying the reason.
func loadWorkflow(root, path, ref string) (*workflow.Definition, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, Usagef("no workflow %q: %s does not exist", ref, path).WithNext(runStartNext)
		}
		return nil, nil, Usagef("workflow %q: %v", ref, err).WithNext(runStartNext)
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	def, err := workflow.ParseWith(name, data, workflowResolver(root))
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
// when the workflow has an agent or leaves step: only an autonomous run
// launches it, and only on an agent session's question. A leaves step's
// builder and verifier, and a checks step's checks, critic, fixer and
// verifier, are checked as agent steps are.
func checkAgentSteps(ctx *Ctx, def *workflow.Definition, mode, agent string, autonomous bool, next string) error {
	agents := false
	for _, name := range def.Order {
		s := def.Steps[name]
		switch s.Kind {
		case workflow.KindAgent:
			agents = true
			if err := checkAgentStep(ctx, s, "an agent step", mode, agent, next); err != nil {
				return err
			}
		case workflow.KindLeaves:
			agents = true
			for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
				r := s.Role(role)
				rs := workflow.Step{Name: s.Name, Agent: r.Agent, PermissionMode: r.PermissionMode}
				if err := checkAgentStep(ctx, rs, "a leaves step whose "+role+" runs", mode, agent, next); err != nil {
					return err
				}
			}
		case workflow.KindChecks:
			agents = true
			for _, c := range s.Checks {
				rs := workflow.Step{Name: s.Name, Agent: c.Agent, PermissionMode: c.PermissionMode}
				if err := checkAgentStep(ctx, rs, "a checks step whose check "+c.Name+" runs", mode, agent, next); err != nil {
					return err
				}
			}
			for _, role := range []string{workflow.RoleCritic, workflow.RoleFixer, workflow.RoleVerifier} {
				r := s.Role(role)
				rs := workflow.Step{Name: s.Name, Agent: r.Agent, PermissionMode: r.PermissionMode}
				if err := checkAgentStep(ctx, rs, "a checks step whose "+role+" runs", mode, agent, next); err != nil {
					return err
				}
			}
		}
	}
	if autonomous && agents {
		return checkAgentStep(ctx, def.Arbiter(), "the arbiter step", mode, agent, next)
	}
	return nil
}

// checkAgentStep checks one step's effective adapter and permission mode
// for checkAgentSteps; what names the step's role in the message.
func checkAgentStep(ctx *Ctx, s workflow.Step, what, mode, agent, next string) error {
	req := agentStepRequest(s, mode, agent)
	if req.PermissionMode == "" {
		return Usagef("step %q is %s with no permission mode: give --permission-mode or set permission_mode on the step", s.Name, what).
			WithNext(next)
	}
	if _, err := startAdapter(ctx, req); err != nil {
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
		if s.Kind != workflow.KindAgent && s.Kind != workflow.KindLeaves && s.Kind != workflow.KindChecks {
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
	resp, err := runEngineRetry(ctx, rs, common, rid, change, worktree, mode, agent, autonomous, noWait, false, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// runEngineRetry is runEngine; with retry it runs the engine's Retry, which
// restarts a failed run at its failed step. With before it runs the
// engine's Continue: before checks the run under its lock and gives the
// entries to journal there, and an *Error from it is returned as it is.
func runEngineRetry(ctx *Ctx, rs *run.Store, common, rid, change, worktree, mode, agent string, autonomous, noWait, retry bool, before func(run.State) ([]run.Entry, error)) (RunResponse, error) {
	sessions := session.NewStore(session.Dir(common))
	e := &run.Engine{
		Store: rs,
		Runner: &stepRunner{
			ctx: ctx, change: change, run: rid, worktree: worktree, main: filepath.Dir(common), mode: mode, agent: agent,
			autonomous: autonomous, sessions: sessions,
		},
		Clock:  runClock,
		Now:    ctx.Now,
		NoWait: noWait,
	}
	if !ctx.DryRun {
		p := runProcess(ctx)
		e.Process = &p
	}
	if ctx.Stderr != nil {
		e.Log = func(line string) { fmt.Fprintf(ctx.Stderr, "rhiza: %s\n", line) }
	}
	if ctx.Rex != nil {
		x := reportToRex(e, ctx.Rex)
		e.Runner.(*stepRunner).launched = func(sid string) { x.launchBlocked(change, sid) }
	}
	// Under the run's lock and the change set's, check that no other run of
	// it is live before the run is reopened or continued.
	e.Admit = func(st run.State) (func(), error) {
		if st.Ended && !retry {
			return func() {}, nil
		}
		unlock, err := rs.LockChange(st.Change)
		if err != nil {
			return nil, Errorf("lock change set %s: %v", st.Change, err)
		}
		if err := refuseLiveRun(rs, st.Change, rid); err != nil {
			unlock()
			return nil, err
		}
		return unlock, nil
	}
	defer e.Runner.(*stepRunner).closeModel()
	run1 := e.Run
	if retry {
		run1 = e.Retry
	}
	if before != nil {
		run1 = func(id string) (run.State, error) { return e.Continue(id, before) }
	}
	st, err := run1(rid)
	var apiErr *Error
	switch {
	case errors.As(err, &apiErr):
		return RunResponse{}, apiErr
	case errors.Is(err, run.ErrLive), errors.Is(err, run.ErrFinished), errors.Is(err, run.ErrNotFailed):
		return RunResponse{}, Conflictf("%v", err).WithNext("rhiza run status " + rid)
	case err != nil:
		return RunResponse{}, Errorf("run %s: %v", rid, err).WithNext("rhiza run status " + rid)
	}
	return RunResponse{Run: runInfo(ctx, st, false, sessions)}, nil
}

// runProcess is the calling process as a run's process: its host and pid
// and, when it runs in Rex, the session, window and block of its tab. Rex names the
// window of the block (`rex block inspect`); when it cannot, the record
// has the block alone.
func runProcess(ctx *Ctx) run.Process {
	window := ""
	if ctx.RexBlock != "" && ctx.Exec != nil {
		args := []string{"block", "inspect"}
		if ctx.RexSession != "" {
			args = append(args, "-s", ctx.RexSession)
		}
		args = append(args, "--json", ctx.RexBlock)
		if out, err := ctx.Exec("rex", args, "", nil); err == nil {
			var in struct {
				Block struct {
					WindowID string `json:"window_id"`
				} `json:"block"`
			}
			if json.Unmarshal(out, &in) == nil {
				window = in.Block.WindowID
			}
		}
	}
	session := ""
	if ctx.RexBlock != "" {
		session = ctx.RexSession
	}
	return run.ThisProcess(session, window, ctx.RexBlock)
}

// stepRunner is the engine's run.Runner: command and wait steps run
// through `sh -c` in the run's worktree, agent steps through session
// start on the run's change set.
type stepRunner struct {
	ctx      *Ctx
	change   string
	run      string
	worktree string
	// main is the main checkout: the folder holding the git common dir.
	main  string
	mode  string
	agent string
	// autonomous marks every session the run launches, so the agent's
	// environment tells it the owner is not there.
	autonomous bool
	sessions   *session.Store
	// launched, when set, is told the session of a milestone request's
	// interview as it launches (see milestoneInterview).
	launched func(sid string)
	// blockPending holds the rejoined sessions whose Rex block check Rex
	// has not answered yet (see run_lost.go).
	lostMu       sync.Mutex
	blockPending map[string]bool
	// wtRoot, modelMu and rm are the run worktree's Rhiza root and the
	// open index its polls read (see run_model_cache.go).
	wtRoot  string
	modelMu sync.Mutex
	rm      *runModel
	// debounce overrides pollDebounce; writeUntil is when the settle after
	// the runner's own write ends.
	debounce   time.Duration
	writeUntil time.Time
}

// Command runs the step's command, its standard output and error
// interleaved and copied to ctx.Stderr. When it exits non-zero the detail is
// the tail of that output (see outputTail).
func (r *stepRunner) Command(step workflow.Step, env []string) (int, string, error) {
	tail := &tailWriter{}
	var buf bytes.Buffer
	code, err := r.sh("exec 2>&1\n"+step.Run, env, io.MultiWriter(r.ctx.Stderr, tail, &buf))
	if err != nil || code == 0 {
		return code, "", err
	}
	if code == run.StopExit {
		return code, lastStopLine(buf.String()), nil
	}
	return code, tail.String(), nil
}

// lastStopLine is the text after the last "stop: " line of out.
func lastStopLine(out string) string {
	detail := ""
	for _, l := range strings.Split(out, "\n") {
		if d, ok := strings.CutPrefix(l, "stop: "); ok {
			detail = d
		}
	}
	return detail
}

// Condition runs the wait step's condition once and returns its exit code;
// on run.StopExit, detail is the last stop: line of its output.
func (r *stepRunner) Condition(step workflow.Step, env []string) (int, string, error) {
	var buf bytes.Buffer
	code, err := r.sh("exec 2>&1\n"+step.Until, env, &buf)
	if err != nil || code != run.StopExit {
		return code, "", err
	}
	io.WriteString(r.ctx.Stderr, buf.String())
	return code, lastStopLine(buf.String()), nil
}

// Caps of a failed command step's kept output.
const (
	stepTailLines = 20
	stepTailBytes = 2048
)

// tailWriter keeps the last stepTailBytes of what is written to it.
type tailWriter struct {
	buf []byte
	cut bool
}

func (t *tailWriter) Write(p []byte) (int, error) {
	t.buf = append(t.buf, p...)
	if over := len(t.buf) - stepTailBytes; over > 0 {
		t.buf = append(t.buf[:0], t.buf[over:]...)
		t.cut = true
	}
	return len(p), nil
}

// String is the kept output: at most stepTailLines whole lines and
// stepTailBytes bytes, cut on a UTF-8 boundary, without trailing newlines.
func (t *tailWriter) String() string {
	b := t.buf
	for len(b) > 0 && !utf8.RuneStart(b[0]) {
		b = b[1:]
	}
	s := strings.TrimRight(string(b), "\r\n")
	lines := strings.Split(s, "\n")
	// A buffer cut mid-line starts with a partial line; drop it unless
	// it is all there is.
	if t.cut && len(lines) > 1 {
		lines = lines[1:]
	}
	if len(lines) > stepTailLines {
		lines = lines[len(lines)-stepTailLines:]
	}
	return strings.Join(lines, "\n")
}

// Holds runs the wait step's condition once; exit 0 means it holds.
func (r *stepRunner) Holds(step workflow.Step, env []string) (bool, error) {
	code, err := r.sh(step.Until, env, nil)
	return code == 0, err
}

// dir is where command and wait steps run: the worktree, or the main
// checkout once the worktree is gone. The close-out step removes the
// worktree, so the steps after it run in the main checkout, at the same
// project folder. The main checkout comes from the git common dir the run
// was started with, not from git run in the vanished folder.
func (r *stepRunner) dir() string {
	if _, err := os.Stat(r.worktree); err == nil {
		return r.worktree
	}
	if r.main == "" {
		return r.ctx.Root
	}
	return filepath.Join(r.main, r.prefix())
}

// prefix is the project folder's path within its checkout: ctx.Root's
// position under the run's worktree or the main checkout, "" at the top.
func (r *stepRunner) prefix() string {
	for _, base := range []string{r.worktree, r.main} {
		if rel, err := filepath.Rel(base, r.ctx.Root); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
			return rel
		}
	}
	return ""
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
// run's change set, on model, and returns its ID. An autonomous run's
// session is recorded autonomous. Its tab is its own.
func (r *stepRunner) StartAgent(step workflow.Step, model string, _ []string) (string, error) {
	req := agentStepRequest(step, r.mode, r.agent)
	req.Change, req.Run, req.autonomous, req.model = r.change, r.run, r.autonomous, model
	req.milestone = r.milestoneOf()
	t, err := resolveChangeWorktree(r.ctx, r.change, "no agent can work on it")
	if err != nil {
		return "", err
	}
	req.Prompt = step.RenderPrompt(workflow.PromptVars{
		Change: t.cs.String(), Slug: t.slug, Title: t.title, Run: r.run, Step: step.Name,
	})
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	sid := resp.(SessionStartResponse).Session.ID
	if r.launched != nil && milestoneInterviewLaunch(r.runState(), step.Name, r.autonomous) {
		r.launched(sid)
	}
	return sid, nil
}

// milestoneInterviewStep is the plan workflow's step whose session
// interviews the owner.
const milestoneInterviewStep = "interview"

// milestoneInterviewLaunch reports whether the agent step's session that
// launches is a milestone request's interview, which waits on the owner
// from launch: not a feature's plan run, not an autonomous run.
func milestoneInterviewLaunch(st run.State, step string, autonomous bool) bool {
	return step == milestoneInterviewStep && st.MilestonePlan && !autonomous && !st.Autonomous
}

// milestoneOf is the milestone the run's request names, from its journal;
// empty for none or when the journal cannot be read.
func (r *stepRunner) milestoneOf() string {
	return r.runState().Milestone
}

// runState is the run's state from its journal; zero when it cannot be read.
func (r *stepRunner) runState() run.State {
	common, err := gitCommonDir(r.ctx)
	if err != nil {
		return run.State{}
	}
	j, err := run.NewStore(run.Dir(common)).Read(r.run)
	if err != nil {
		return run.State{}
	}
	return j.State()
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
// question and the arbiter guide topic. The arbiter writes where the
// question lives (see arbiterWrites).
func (r *stepRunner) StartArbiter(step workflow.Step, question string, _ []string) (string, error) {
	req := agentStepRequest(step, r.mode, r.agent)
	req.Change, req.Run, req.autonomous = r.change, r.run, r.autonomous
	req.Prompt = fmt.Sprintf(arbiterPrompt, question)
	var err error
	if req.writes, req.writesRecord, err = r.arbiterWrites(question); err != nil {
		return "", err
	}
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		if req.writesRecord {
			if why := r.settleArbiterRecord(req.writes, false); why != "" {
				err = withNote(err, why)
			}
		}
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// advisorPrompt is the advisor's whole prompt: the question the session
// asked the owner (a normal run adds no question node) and the guide topic.
const advisorPrompt = "Recommend an answer to this question the owner was asked: %s -- follow rhiza guide advisor."

// advisorStep is the step name of an advisor's session.
const advisorStep = "advisor"

// StartAdvisor starts the arbiter step's adapter and mode as a session
// that recommends an answer to the question session asker asked, text, on
// the run's change set, and returns at once. The session answers nothing:
// it reports its recommendation with run recommend, recorded against the
// session it advises.
func (r *stepRunner) StartAdvisor(step workflow.Step, asker, text string, _ []string) (string, error) {
	req := agentStepRequest(step, r.mode, r.agent)
	req.Step, req.Change, req.Run, req.advises = advisorStep, r.change, r.run, asker
	req.Prompt = fmt.Sprintf(advisorPrompt, text)
	resp, err := runSessionStart(r.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.(SessionStartResponse).Session.ID, nil
}

// worktreeWorkspace opens the model as the run's worktree holds it, where
// its sessions write; read only unless write.
func (r *stepRunner) worktreeWorkspace(write bool) (*workspace, error) {
	if write {
		r.wrote()
	}
	prefix, err := r.ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return nil, err
	}
	return openWorkspace(filepath.Join(r.worktree, strings.TrimSpace(prefix)), r.ctx.Now, r.ctx.IDs, !write || r.ctx.DryRun)
}

// questionHome is the change set whose ops create node q, in the run's
// worktree: the change set the question lives in.
func questionHome(ww *workspace, q id.ID) (changeset.ChangeSet, bool) {
	for _, c := range ww.changeSets() {
		if slices.ContainsFunc(c.Ops, func(o op.Op) bool { return o.Type == op.TypeNodeCreate && o.Node == q }) {
			return c, true
		}
	}
	return changeset.ChangeSet{}, false
}

// arbiterIntent is the intent of the record change set an arbiter answers
// a question in when the question's own change set has landed.
const arbiterIntent = "The arbiter's answer to question %s"

// arbiterWrites is the change set the arbiter of question writes into,
// when that is not the run's: empty when the question lives in the run's
// change set (a session's question), or nowhere that can be told. When
// it lives in another open change set (a review question left in the
// plan's draft), that change set. When its change set has landed (a review
// question written as a record, or left in the run's landed change set), a
// record change set opened here for the answer, with record true; Answer
// lands it once the answer passes.
func (r *stepRunner) arbiterWrites(question string) (string, bool, error) {
	q, err := id.Parse(question)
	if err != nil {
		return "", false, nil
	}
	ww, err := r.worktreeWorkspace(true)
	if err != nil {
		return "", false, err
	}
	home, ok := questionHome(ww, q)
	switch {
	case !ok || home.ID.String() == r.change:
		return "", false, nil
	case !home.State.Terminal():
		return home.ID.String(), false, nil
	case home.State != changeset.StateLanded:
		return "", false, nil
	}
	handle := question
	if n, ok := ww.model(id.ID{}).Node(q); ok && n.Handle != "" {
		handle = n.Handle
	}
	actor, forID, err := ww.resolveActor(id.ID{}, r.ctx)
	if err != nil {
		return "", false, err
	}
	cs, _, err := ww.openChangeSet(changeset.Open{Kind: changeset.KindRecord, Intent: fmt.Sprintf(arbiterIntent, handle)}, actor, forID)
	if err != nil {
		return "", false, err
	}
	return cs.String(), true, nil
}

// settleArbiterRecord lands the arbiter's record change set ref when ok,
// else withdraws it; one that is already terminal is left alone. It
// returns what went wrong, or "".
func (r *stepRunner) settleArbiterRecord(ref string, ok bool) string {
	cs, err := id.Parse(ref)
	if err != nil {
		return err.Error()
	}
	ww, err := r.worktreeWorkspace(true)
	if err != nil {
		return err.Error()
	}
	if c, found := ww.changeSet(cs); !found || c.State.Terminal() {
		return ""
	}
	actor, forID, err := ww.resolveActor(id.ID{}, r.ctx)
	if err != nil {
		return AsError(err).Message
	}
	var body changeset.Body = changeset.Land{}
	what := "land"
	if !ok {
		body, what = changeset.Withdraw{}, "withdraw"
	}
	if _, err := ww.appendEntry(cs, body, actor, forID); err != nil {
		return fmt.Sprintf("%s the arbiter's record change set %s: %s", what, cs, AsError(err).Message)
	}
	return ""
}

// Answer reads the model in the run's worktree, where its sessions write,
// as the change set arbiter session sid writes into sees it (the run's,
// unless the question lives elsewhere: see arbiterWrites), and looks for a
// live Decision with a live answers link to question node question,
// wording arbiterWording and a reasoning. The first such Decision, by link
// ID, is the answer; with none, Why names what the nearest one lacks. A
// record change set opened for the arbiter lands with a valid answer and
// is withdrawn without one.
func (r *stepRunner) Answer(question, sid string) (run.Answer, error) {
	q, err := id.Parse(question)
	if err != nil {
		return run.Answer{}, fmt.Errorf("question %q is not a node ID: %v", question, err)
	}
	ref, record := r.change, false
	if sid != "" {
		rec, err := readSession(r.sessions, sid)
		if err != nil {
			return run.Answer{}, err
		}
		if rec.Writes != "" {
			ref, record = rec.Writes, rec.WritesRecord
		}
	}
	cs, err := id.Parse(ref)
	if err != nil {
		return run.Answer{}, fmt.Errorf("run %s names change set %q: %v", r.run, ref, err)
	}
	v, err := r.readModel("answer:"+cs.String()+":"+q.String(), cs, func(ww *workspace) (any, error) {
		return arbiterAnswer(ww.model(cs), q), nil
	})
	if err != nil {
		return run.Answer{}, err
	}
	ans := v.(run.Answer)
	if record {
		if why := r.settleArbiterRecord(ref, ans.OK); why != "" {
			if ans.OK {
				return run.Answer{Why: why}, nil
			}
			ans.Why += "; " + why
		}
	}
	return ans, nil
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
			Fields: []string{"statement=" + text, "session=" + rec.ID}, lenient: true})
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
		if rec, err = r.markLostIfEnded(rec, hookEndGrace); err != nil {
			return run.AgentStatus{}, err
		}
		rec = r.retryBlock(rec)
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
	return run.AgentStatus{
		Done: true, Outcome: workflow.OutcomeFail, Detail: detail,
		Lost: x.State == session.StateLost, SaidDone: rec.Status.SaidDone(),
	}, nil
}

// AgentEnded records the session's token use as its step ends (see
// recordTokenUse); a warning comes back as the error the engine logs.
func (r *stepRunner) AgentEnded(_ workflow.Step, sid string) error {
	rec, err := readSession(r.sessions, sid)
	if err != nil {
		return err
	}
	if warn := recordTokenUse(r.ctx, r.sessions, rec); warn != "" {
		return errors.New(warn)
	}
	return nil
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
	return rexReport("blocked", n.Block.Kind, n.Block.TabTitle, rexNoticeMsg(n))
}

// rexNoticeMsg is the message of a blocked notice.
func rexNoticeMsg(n run.Notice) string {
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
	return msg
}

// rexReport encodes one OSC 7501 root-record report from app rhiza (see
// rexRecord).
func rexReport(state, kind, title, msg string) []byte {
	return rexRecord("", state, kind, title, msg, nil)
}

// rexRecord encodes one OSC 7501 report from app rhiza on record id (""
// for the root record). The kind is kept only on a blocked state and only
// when Rex knows it, progress only on working or blocked; title and msg
// have their control characters replaced with spaces, are capped on a
// UTF-8 boundary and are base64-encoded, and are left out when empty.
func rexRecord(id, state, kind, title, msg string, progress *int) []byte {
	var b strings.Builder
	b.WriteString("\x1b]7501;state=")
	b.WriteString(state)
	if id != "" {
		b.WriteString(":id=" + id)
	}
	if state == session.StateBlocked {
		switch kind {
		case session.KindPermission, session.KindQuestion, session.KindAuth:
			b.WriteString(":kind=" + kind)
		}
	}
	if progress != nil && (state == session.StateWorking || state == session.StateBlocked) {
		b.WriteString(":progress=" + strconv.Itoa(min(100, max(0, *progress))))
	}
	if state != "clear" {
		b.WriteString(":app=rhiza")
	}
	if t := rexField(title, rexTitleMax); t != "" {
		b.WriteString(":title=" + base64.StdEncoding.EncodeToString([]byte(t)))
	}
	if m := rexField(msg, rexMsgMax); m != "" {
		b.WriteString(":msg=" + base64.StdEncoding.EncodeToString([]byte(m)))
	}
	b.WriteString("\x1b\\")
	return []byte(b.String())
}

// runRex writes a run's program status to Rex: the root record says
// the run's step (and a leaves step's progress), then done or error; a
// child record leaf/<n> per live leaf of a leaves step says its role and
// whether it works or is blocked, and is cleared when the leaf is no
// longer live. While a session's block is noticed (rexNotice), the root
// record is that blocked notice, as before; once no noticed block is
// left, the root record goes back to the run's status.
type runRex struct {
	w io.Writer
	// root is the last root status record, "" before the first.
	root []byte
	// wrote is the last root record written.
	wrote []byte
	// blocks are the open noticed blocks' reports, by session, and order
	// the sessions in the order they blocked.
	blocks map[string][]byte
	order  []string
	// leaves are the child records written, by ID, until cleared.
	leaves map[string][]byte
	// launch is the session of an open launch block (launchBlocked).
	launch string
}

// reportToRex has e's notices and status written to Rex, w, by one runRex.
func reportToRex(e *run.Engine, w io.Writer) *runRex {
	x := &runRex{w: w, blocks: map[string][]byte{}}
	e.Notify = x.notice
	e.Status = x.status
	return x
}

// launchBlocked reports the interview session of a milestone request's
// plan run as blocked on the owner from its launch, as a notice for the
// session: the agent's own first block replaces it, and the run's next
// status that is not working clears it.
func (x *runRex) launchBlocked(change, sid string) {
	x.launch = sid
	x.notice(run.Notice{
		Change: change, Step: milestoneInterviewStep, Session: sid, Blocked: true,
		Block: &run.Block{Kind: session.KindQuestion, Question: "the milestone interview waits on you"},
	})
}

func (x *runRex) write(b []byte) {
	x.w.Write(b)
}

func (x *runRex) writeRoot(b []byte) {
	if !bytes.Equal(b, x.wrote) {
		x.write(b)
		x.wrote = b
	}
}

// notice writes a blocked notice as the root record, or once a block ends
// the next open one's, else the run's status (or a bare working report
// before any status).
func (x *runRex) notice(n run.Notice) {
	if n.Blocked && n.Block != nil {
		if _, ok := x.blocks[n.Session]; !ok {
			x.order = append(x.order, n.Session)
		}
		x.blocks[n.Session] = rexNotice(n)
		x.writeRoot(x.blocks[n.Session])
		return
	}
	delete(x.blocks, n.Session)
	x.order = slices.DeleteFunc(x.order, func(s string) bool { return s == n.Session })
	switch {
	case len(x.order) > 0:
		x.writeRoot(x.blocks[x.order[len(x.order)-1]])
	case x.root != nil:
		x.writeRoot(x.root)
	default:
		x.writeRoot(rexNotice(n))
	}
}

// status writes rep: a clear for each ended leaf, a record for each live
// one, and the root record unless a noticed block holds it. A finished or
// stopped run's root record is written whatever blocks were noticed.
func (x *runRex) status(rep run.StatusReport) {
	if x.leaves == nil {
		x.leaves = map[string][]byte{}
	}
	if x.launch != "" && rep.Step != milestoneInterviewStep {
		delete(x.blocks, x.launch)
		x.order = slices.DeleteFunc(x.order, func(s string) bool { return s == x.launch })
		x.launch = ""
	}
	for _, l := range rep.Ended {
		x.write(rexRecord(rexLeafID(l), "clear", "", "", "", nil))
		delete(x.leaves, rexLeafID(l))
	}
	for _, l := range rep.Live {
		state := session.StateWorking
		if l.Blocked {
			state = session.StateBlocked
		}
		id := rexLeafID(l)
		if b := rexRecord(id, state, l.Kind, l.Handle, rexLeafMsg(l), nil); !bytes.Equal(b, x.leaves[id]) {
			x.write(b)
			x.leaves[id] = b
		}
	}
	x.root = rexRoot(rep)
	if rep.State != run.StatusWorking {
		clear(x.blocks)
		x.order = nil
	}
	if len(x.order) == 0 {
		x.writeRoot(x.root)
	}
}

// rexLeafID is a live leaf's record ID: leaf/<n>, n its position in leaf
// order, since a handle may be longer than an ID segment; a live check's
// is check/<n>, n its position in its step; a review's change-set fixer
// is fix.
func rexLeafID(l run.LeafStatus) string {
	if rexCheckRole(l.Role) {
		return "check/" + strconv.Itoa(l.N)
	}
	if l.Leaf == "" && l.Role == workflow.RoleFixer {
		return "fix"
	}
	return "leaf/" + strconv.Itoa(l.N)
}

// rexCheckRole reports whether role is a checks step session's.
func rexCheckRole(role string) bool {
	return role == workflow.RoleCheck || role == workflow.RoleCritic
}

// rexLeafMsg says a live leaf's role and whether it works or waits.
func rexLeafMsg(l run.LeafStatus) string {
	msg := l.Handle + ": " + l.Role
	switch {
	case l.Arbiter:
		msg += " waits on the arbiter's answer"
	case l.Blocked:
		msg += " blocked on " + rexKindText(l.Kind)
	default:
		msg += " working"
	}
	switch {
	case l.Attempt > 1 && rexCheckRole(l.Role):
		msg += fmt.Sprintf(" (cycle %d)", l.Attempt)
	case l.Attempt > 1:
		msg += fmt.Sprintf(" (attempt %d)", l.Attempt)
	}
	return msg
}

// rexKindText names what a blocked kind waits on.
func rexKindText(kind string) string {
	switch kind {
	case session.KindPermission:
		return "a permission"
	case session.KindQuestion:
		return "a question"
	case session.KindAuth:
		return "a login"
	}
	return "the owner"
}

// rexRoot is the root record of rep: titled with the run's title (else
// its change set), its msg naming the step and a leaves step's leaves
// passed, or how the run finished or why it stopped, on one line.
func rexRoot(rep run.StatusReport) []byte {
	title := rep.Title
	if title == "" {
		title = rep.Change
	}
	var msg string
	switch rep.State {
	case run.StatusDone:
		msg = "run succeeded"
	case run.StatusError:
		msg = "run failed"
	case run.StatusBlocked:
		msg = fmt.Sprintf("%s: step %s waits on you: %s", title, rep.Step, strings.TrimSpace(rep.Question))
	case run.StatusIdle:
		msg = "run stopped at step " + rep.Step
	default:
		msg = "step " + rep.Step
		switch {
		case rep.Progress != nil && rep.Checks > 0:
			msg += fmt.Sprintf(" · %d of %d checks passed", rep.Passed, rep.Checks)
		case rep.Progress != nil:
			msg += fmt.Sprintf(" · %d of %d leaves passed", rep.Passed, rep.Leaves)
		}
	}
	if d, _, _ := strings.Cut(strings.TrimSpace(rep.Detail), "\n"); d != "" && rep.State != run.StatusWorking && rep.State != run.StatusBlocked {
		msg += ": " + d
	}
	return rexRecord("", rep.State, "", title, msg, rep.Progress)
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
func runInfo(ctx *Ctx, st run.State, live bool, sessions *session.Store) RunInfo {
	info := RunInfo{
		ID: st.ID, Workflow: st.Workflow, Change: st.Change, Title: st.Title, Worktree: st.Worktree,
		Autonomous: st.Autonomous, Current: st.Current, Steps: []RunStep{}, StartedAt: rfc3339(st.StartedAt),
		Amendments: []RunAmendment{}, Messages: []RunMessage{}, Sessions: []RunSession{},
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
	if ctx != nil && ctx.Now != nil {
		info.Stalled = runStalledNow(st, live, sessions, ctx.Now())
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
	for _, m := range st.Messages {
		info.Messages = append(info.Messages, RunMessage{
			At: rfc3339(m.At), Step: m.Step, Session: m.Session, Actor: m.Actor,
			To: m.To, Text: m.Text, Summary: m.Summary, Structured: m.Structured,
		})
	}
	for _, s := range st.Sessions {
		rs := RunSession{Step: s.Step, Session: s.Session}
		if u, ok := st.Tokens[s.Session]; ok {
			rs.Tokens = runTokens(u)
		} else {
			rs.Tokens, rs.NoTokens = transcriptTokens(ctx, sessions, s.Session)
		}
		if rs.Tokens != nil {
			addTokens(&info.Totals.Session, rs.Tokens.Session)
			addTokens(&info.Totals.Helpers, rs.Tokens.Helpers)
		}
		info.Sessions = append(info.Sessions, rs)
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
	for _, rd := range st.Rounds {
		rr := RunRound{Round: rd.Number, At: rfc3339(rd.At), Findings: []RunFinding{}, Confirmed: len(rd.Confirmed())}
		for _, f := range rd.Findings {
			rr.Findings = append(rr.Findings, runFinding(f, rd))
		}
		info.Rounds = append(info.Rounds, rr)
	}
	for _, l := range st.Leaves {
		rl := RunLeaf{
			Step: l.Step, Leaf: l.Leaf, Handle: l.Handle, Sessions: []RunLeafSession{}, Attempts: []RunLeafAttempt{},
			Outcome: l.Outcome, Question: l.Question, QuestionHandle: l.QuestionHandle, Round: l.Round(), Commit: l.Commit, Detail: l.Detail,
		}
		for _, rd := range l.Rounds {
			rl.Rounds = append(rl.Rounds, RunLeafRound{
				Round: rd.Round, After: rd.After, Question: rd.Question, QuestionHandle: rd.QuestionHandle, Decision: rd.Decision, At: rfc3339(rd.At),
			})
		}
		for _, ls := range l.Sessions {
			rl.Sessions = append(rl.Sessions, RunLeafSession{
				Role: ls.Role, Attempt: ls.Attempt, Cycle: ls.Cycle, Session: ls.Session, Model: ls.Model, Commit: ls.Commit, At: rfc3339(ls.At),
			})
		}
		for _, a := range l.Attempts {
			s, _ := l.Session(a.Role, a.Attempt)
			rl.Attempts = append(rl.Attempts, RunLeafAttempt{
				Role: a.Role, Attempt: a.Attempt, Cycle: s.Cycle, Outcome: a.Outcome, Commit: a.Commit, Detail: a.Detail, At: rfc3339(a.At),
			})
		}
		for _, f := range l.Findings {
			rl.Findings = append(rl.Findings, RunLeafFinding{Attempt: f.Attempt, ID: f.Finding.ID, Reason: f.Finding.Reason, At: rfc3339(f.At)})
		}
		info.Leaves = append(info.Leaves, rl)
	}
	for _, c := range st.Checks {
		info.Checks = append(info.Checks, runCheckInfo(st, c))
	}
	for _, f := range st.Fixes {
		info.Fixes = append(info.Fixes, RunFix{
			Step: f.Step, Cycle: f.Cycle, Session: f.Session, Model: f.Model, Since: f.Since,
			Outcome: f.Outcome, Commit: f.Commit, Detail: f.Detail, At: rfc3339(f.At),
		})
	}
	for _, c := range st.CodeChecks {
		info.CodeChecks = append(info.CodeChecks, RunCodeCheck{Step: c.Step, Outcome: c.Outcome, Detail: c.Detail, At: rfc3339(c.At)})
	}
	for _, s := range st.Signoffs {
		info.Signoffs = append(info.Signoffs, runSignoff(s))
	}
	info.SignoffWaiting = signoffWaiting(st) && info.State == runStopped
	return info
}

// runCheckInfo is the response form of check c: a cycle for each of its
// step's cycles that ran it, with its sessions and outcome there.
func runCheckInfo(st run.State, c run.CheckState) RunCheck {
	rc := RunCheck{Step: c.Step, Check: c.Check, Cycles: []RunCheckCycle{}}
	if o, ok := st.CheckOutcome(c); ok {
		rc.Outcome = o.Outcome
	}
	if cy, prior, ok := st.CheckOpen(c); cy > 0 && ok {
		rc.OpenCycle, rc.PriorCycle, rc.PriorOutcome = cy, prior.Cycle, prior.Outcome
	}
	for _, cy := range st.Cycles {
		if cy.Step != c.Step || !slices.Contains(cy.Checks, c.Check) {
			continue
		}
		rcy := RunCheckCycle{Cycle: cy.Cycle, Sessions: []RunCheckSession{}, Kind: cy.Kind}
		if cy.Answered() {
			rcy.Decisions = cy.Decisions
		}
		for _, s := range c.Sessions {
			if s.Cycle != cy.Cycle {
				continue
			}
			rs := RunCheckSession{Role: s.Role, Session: s.Session, Model: s.Model, At: rfc3339(s.At)}
			if j, ok := c.Judgement(s.Role, s.Cycle); ok {
				rs.Outcome, rs.Detail = j.Outcome, j.Detail
			}
			rcy.Sessions = append(rcy.Sessions, rs)
		}
		if o, ok := c.Outcome(cy.Cycle); ok {
			rcy.Outcome, rcy.Detail, rcy.At = o.Outcome, o.Detail, rfc3339(o.At)
		}
		rc.Cycles = append(rc.Cycles, rcy)
	}
	for _, rd := range c.Rounds {
		rr := RunCheckRound{Round: rd.Number, Cycle: rd.Cycle, At: rfc3339(rd.At), Findings: []RunFinding{}, Confirmed: len(rd.Confirmed())}
		for _, f := range rd.Findings {
			rr.Findings = append(rr.Findings, runFinding(f, rd))
		}
		for _, v := range rd.Verdicts {
			rr.Verdicts = append(rr.Verdicts, RunVerdict{Finding: v.Finding, Verdict: v.Verdict, Reason: v.Reason})
		}
		rc.Rounds = append(rc.Rounds, rr)
	}
	for _, q := range st.ReviewQuestions {
		if q.Step == c.Step && q.Check == c.Check {
			rc.Questions = append(rc.Questions, RunCheckQuestion{
				Cycle: q.Cycle, Finding: q.Finding, Leaf: q.Leaf, LeafHandle: q.Handle, Question: q.Question, QuestionHandle: q.QuestionHandle,
			})
		}
	}
	return rc
}

// transcriptTokens reads the token use of a session with no journaled
// reading from its transcript, at print time. Nothing is journaled. When
// the transcript cannot be read it returns nil and the reason: the session
// is still running, or its transcript is unreadable, with the warning.
func transcriptTokens(ctx *Ctx, sessions *session.Store, sid string) (*SessionTokenUse, string) {
	unreadable := func(why string) (*SessionTokenUse, string) {
		return nil, "transcript unreadable: " + why
	}
	if ctx == nil || sessions == nil || (ctx.Home == "" && ctx.Getenv == nil) {
		return unreadable("the home folder and environment are unknown")
	}
	rec, err := sessions.Read(sid)
	if err != nil {
		return unreadable(AsError(err).Message)
	}
	a, ok := adapter.Lookup(rec.Agent)
	if !ok {
		return unreadable(fmt.Sprintf("session %s names agent %q, which this rhiza does not know", sid, rec.Agent))
	}
	data, err := readTranscript(ctx, a, sid)
	if err == nil {
		var use adapter.TokenUse
		if use, err = a.TokenUse(data); err == nil {
			return runTokens(runTokenUse(sessionTokenUse(use))), ""
		}
		err = Errorf("read the token use of session %s: %v", sid, err)
	}
	if !rec.Status.Exited() {
		return nil, "still running"
	}
	return unreadable(AsError(err).Message)
}

// runTokens is the response form of a journaled token use reading.
func runTokens(u run.TokenUse) *SessionTokenUse {
	c := func(t run.TokenCounts) SessionTokens {
		return SessionTokens{Input: t.Input, Output: t.Output, CacheCreation: t.CacheCreation, CacheRead: t.CacheRead}
	}
	return &SessionTokenUse{Session: c(u.Session), Helpers: c(u.Helpers)}
}

// addTokens adds t's counts into sum.
func addTokens(sum *SessionTokens, t SessionTokens) {
	sum.Input += t.Input
	sum.Output += t.Output
	sum.CacheCreation += t.CacheCreation
	sum.CacheRead += t.CacheRead
}

// queuedDaemon ends the line of a queued action: which daemon takes it, or
// that none runs.
func queuedDaemon(w io.Writer, holder string) {
	if holder == "" {
		fmt.Fprintf(w, ": no daemon runs, so it waits for one (rhiza daemon)\n")
	} else {
		fmt.Fprintf(w, ": daemon %s continues it headless\n", holder)
	}
}

// Human reports the run, or the plan of a dry run.
func (r RunResponse) Human(w io.Writer) {
	if r.Queued && !r.DryRun && r.Run.ID != "" {
		fmt.Fprintf(w, "Queued the continuation of run %s for the daemon", r.Run.ID)
		queuedDaemon(w, r.Daemon)
		return
	}
	if r.Queued && !r.DryRun {
		fmt.Fprintf(w, "Queued a %s run of %s for the daemon", r.Run.Workflow, r.Run.Change)
		if r.Daemon == "" {
			fmt.Fprintf(w, ": no daemon runs, so it waits for one (rhiza daemon)\n")
		} else {
			fmt.Fprintf(w, ": daemon %s starts it headless\n", r.Daemon)
		}
		return
	}
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

// Human reports the abandoned run, or the plan of a dry run.
func (r RunAbandonResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Run %s would be abandoned\n\n", r.Run.ID)
		runLines(w, r.Run)
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written.\n")
		return
	}
	r.Run.Human(w)
	if len(r.Cleanup) > 0 {
		fmt.Fprintf(w, "\n## Clean-up\n\n")
		for _, c := range r.Cleanup {
			fmt.Fprintf(w, "- %s %s", c.Step, c.Status)
			if c.Detail != "" {
				fmt.Fprintf(w, ": %s", c.Detail)
			}
			fmt.Fprintln(w)
		}
	}
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
			first, rest, _ := strings.Cut(s.Detail, "\n")
			if first != "" {
				fmt.Fprintf(w, ": %s", first)
			}
			fmt.Fprintf(w, "\n")
			for _, l := range strings.Split(rest, "\n") {
				if l != "" {
					fmt.Fprintf(w, "    %s\n", l)
				}
			}
		}
	}
	if len(r.Leaves) > 0 {
		fmt.Fprintf(w, "\n## Leaves\n\n")
		for _, l := range r.Leaves {
			outcome := l.Outcome
			if outcome == "" {
				outcome = "open"
			}
			fmt.Fprintf(w, "- %s (step %s): %s", l.Handle, l.Step, outcome)
			if l.Round > 1 {
				fmt.Fprintf(w, " in round %d", l.Round)
			}
			if l.Detail != "" {
				fmt.Fprintf(w, ": %s", l.Detail)
			}
			fmt.Fprintf(w, "\n")
			if l.Question != "" {
				fmt.Fprintf(w, "  - question %s (%s) blocks it\n", l.QuestionHandle, l.Question)
			}
			for _, rd := range l.Rounds {
				fmt.Fprintf(w, "  - round %d from attempt %d: decision %s answers question %s (%s)\n", rd.Round, rd.After+1, rd.Decision, rd.QuestionHandle, rd.Question)
			}
			for _, ls := range l.Sessions {
				fmt.Fprintf(w, "  - attempt %d %s", ls.Attempt, ls.Role)
				if ls.Cycle > 0 {
					fmt.Fprintf(w, " in cycle %d", ls.Cycle)
				}
				fmt.Fprintf(w, ", session %s", ls.Session)
				for _, a := range l.Attempts {
					if a.Role == ls.Role && a.Attempt == ls.Attempt {
						fmt.Fprintf(w, ": %s", a.Outcome)
						if a.Detail != "" {
							fmt.Fprintf(w, " (%s)", a.Detail)
						}
					}
				}
				fmt.Fprintf(w, "\n")
			}
			for _, f := range l.Findings {
				fmt.Fprintf(w, "  - attempt %d finding %s: %s\n", f.Attempt, f.ID, oneLine(f.Reason))
			}
		}
	}
	if len(r.Checks) > 0 {
		fmt.Fprintf(w, "\n## Checks\n\n")
		for _, c := range r.Checks {
			outcome := c.Outcome
			if outcome == "" {
				outcome = "open"
			}
			if c.OpenCycle > 0 {
				outcome = fmt.Sprintf("open in cycle %d (cycle %d %s)", c.OpenCycle, c.PriorCycle, c.PriorOutcome)
			}
			fmt.Fprintf(w, "- %s (step %s): %s", c.Check, c.Step, outcome)
			if n := len(c.Cycles); n > 0 && c.OpenCycle == 0 {
				fmt.Fprintf(w, " after %d %s", n, shipPlural(n, "cycle", "cycles"))
			}
			fmt.Fprintf(w, "\n")
			for _, cy := range c.Cycles {
				outcome := cy.Outcome
				if outcome == "" {
					outcome = "open"
				}
				fmt.Fprintf(w, "  - cycle %d: %s", cy.Cycle, outcome)
				if cy.Kind == run.CycleOwner {
					fmt.Fprintf(w, ", the owner's changes")
				}
				if len(cy.Decisions) > 0 {
					fmt.Fprintf(w, ", answered by %s", strings.Join(cy.Decisions, ", "))
				}
				for _, s := range cy.Sessions {
					fmt.Fprintf(w, ", %s session %s", s.Role, s.Session)
					if s.Outcome != "" {
						fmt.Fprintf(w, " %s", s.Outcome)
						if s.Detail != "" {
							fmt.Fprintf(w, " (%s)", oneLine(s.Detail))
						}
					}
				}
				fmt.Fprintf(w, "\n")
			}
			for _, rd := range c.Rounds {
				if len(rd.Findings) == 0 {
					fmt.Fprintf(w, "  - round %d (cycle %d): no findings\n", rd.Round, rd.Cycle)
					continue
				}
				fmt.Fprintf(w, "  - round %d (cycle %d): %d %s, %d confirmed\n", rd.Round, rd.Cycle, len(rd.Findings), shipPlural(len(rd.Findings), "finding", "findings"), rd.Confirmed)
				for _, f := range rd.Findings {
					fmt.Fprintf(w, "    - %s\n", findingLine(f))
				}
			}
			for _, q := range c.Questions {
				fmt.Fprintf(w, "  - question %s after cycle %d", q.QuestionHandle, q.Cycle)
				if q.Finding != "" {
					fmt.Fprintf(w, " on %s", q.Finding)
				}
				if q.LeafHandle != "" {
					fmt.Fprintf(w, ", blocks leaf %s", q.LeafHandle)
				}
				fmt.Fprintf(w, "\n")
			}
		}
	}
	if len(r.Fixes) > 0 {
		fmt.Fprintf(w, "\n## Change-set fixes\n\n")
		for _, f := range r.Fixes {
			outcome := f.Outcome
			if outcome == "" {
				outcome = "open"
			}
			fmt.Fprintf(w, "- cycle %d (step %s): fixer session %s: %s", f.Cycle, f.Step, f.Session, outcome)
			if f.Commit != "" {
				fmt.Fprintf(w, ", commit %s", shortRev(f.Commit))
			}
			if f.Detail != "" {
				fmt.Fprintf(w, " (%s)", oneLine(f.Detail))
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.CodeChecks) > 0 {
		fmt.Fprintf(w, "\n## Code checks\n\n")
		for _, c := range r.CodeChecks {
			fmt.Fprintf(w, "- step %s: %s", c.Step, c.Outcome)
			if c.Detail != "" {
				fmt.Fprintf(w, " (%s)", oneLine(c.Detail))
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.Signoffs) > 0 {
		fmt.Fprintf(w, "\n## Sign-offs\n\n")
		for _, s := range r.Signoffs {
			fmt.Fprintf(w, "- %s\n", s.line())
			for _, f := range s.Findings {
				fmt.Fprintf(w, "  - %s %s\n", f.Check, findingLine(f))
			}
		}
	}
	if len(r.Accepts) > 0 {
		fmt.Fprintf(w, "\n## Accepted in the owner's place\n\n")
		for _, a := range r.Accepts {
			fmt.Fprintf(w, "- %s %s change set %s (log entry %s)\n", a.At, a.Verb, a.Change, a.Entry)
		}
	}
	if len(r.Questions) > 0 {
		fmt.Fprintf(w, "\n## Questions put to the arbiter\n\n")
		for _, q := range r.Questions {
			fmt.Fprintf(w, "- %s (%s): %s", q.Handle, q.Question, q.Text)
			switch {
			case q.Decision != "":
				fmt.Fprintf(w, " — decision %s", q.Decision)
				if !q.Answered {
					fmt.Fprintf(w, ", not yet delivered")
				}
			case q.Outcome != "":
				fmt.Fprintf(w, " — not decided: the arbiter's verdict was %s", q.Outcome)
			default:
				fmt.Fprintf(w, " — open")
			}
			fmt.Fprintf(w, "\n")
		}
	}
	if len(r.Rounds) > 0 {
		fmt.Fprintf(w, "\n## Review rounds\n\n")
		for _, rd := range r.Rounds {
			if len(rd.Findings) == 0 {
				fmt.Fprintf(w, "- Round %d: no findings\n", rd.Round)
				continue
			}
			fmt.Fprintf(w, "- Round %d: %d %s, %d confirmed\n", rd.Round, len(rd.Findings), shipPlural(len(rd.Findings), "finding", "findings"), rd.Confirmed)
			for _, f := range rd.Findings {
				fmt.Fprintf(w, "  - %s\n", findingLine(f))
			}
		}
	}
	if len(r.Waiting) > 0 {
		fmt.Fprintf(w, "\n## Waiting on the owner\n\n")
		for _, c := range r.Waiting {
			fmt.Fprintf(w, "- %s %s: proposed: %s\n", c.Change, c.Title, c.Next)
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
	if len(r.Messages) > 0 {
		fmt.Fprintf(w, "\n## Messages\n\n")
		for _, m := range r.Messages {
			note := ""
			if m.Summary != "" {
				note += ", summary: " + m.Summary
			}
			if m.Structured {
				note += ", structured"
			}
			fmt.Fprintf(w, "- %s %s to %s (step %s%s): %s\n", m.At, m.Actor, m.To, m.Step, note, indentRest(m.Text))
		}
	}
	if len(r.Sessions) > 0 {
		fmt.Fprintf(w, "\n## Sessions\n\n")
		for _, s := range r.Sessions {
			if s.Tokens == nil {
				fmt.Fprintf(w, "- %s, session %s: no token use shown, %s\n", s.Step, s.Session, s.NoTokens)
				continue
			}
			fmt.Fprintf(w, "- %s, session %s: %s\n", s.Step, s.Session, tokenUseLine(*s.Tokens))
		}
		fmt.Fprintf(w, "\nTotals: %s\n", tokenUseLine(r.Totals))
	}
	switch r.State {
	case runStopped:
		if r.SignoffWaiting {
			fmt.Fprintf(w, "\nnext: rhiza run summary %s\nnext: rhiza run approve %s\nnext: rhiza run request-changes <file> %s\n", r.ID, r.ID, r.ID)
			break
		}
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

// indentRest indents every non-empty line of s after the first by two
// spaces, so a multi-line text stays inside its markdown bullet.
func indentRest(s string) string {
	lines := strings.Split(s, "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			lines[i] = "  " + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

// tokenUseLine is a token use reading as "tokens <counts>; helper tokens
// <counts>".
func tokenUseLine(u SessionTokenUse) string {
	return "tokens " + tokensLine(u.Session) + "; helper tokens " + tokensLine(u.Helpers)
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
	if r.PR != "" {
		fmt.Fprintf(w, "- Pull request: %s\n", r.PR)
	}
	if r.Detail != "" {
		fmt.Fprintf(w, "- Detail: %s\n", r.Detail)
	}
	if s := r.Stalled; s != nil {
		fmt.Fprintf(w, "- Stalled: %s session %s has not reported since %s\n", s.Step, s.Session, s.Since)
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

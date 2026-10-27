// Package run keeps the journal of a workflow run: an append-only JSONL file
// per run, folded to a state. It is work-layer state and lives outside the
// model's log. It runs no program and names no agent.
package run

import "time"

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
	TypeAnswered   = "answered"
	TypeMessage    = "message"
	TypeTokens     = "tokens"
	TypeFindings   = "findings"
	TypeVerdict    = "verdict"
	TypeGated      = "gated"
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
	Workflow       string `json:"workflow"`
	WorkflowPath   string `json:"workflow_path,omitempty"`
	Definition     string `json:"definition"`
	Change         string `json:"change"`
	Title          string `json:"title,omitempty"`
	Worktree       string `json:"worktree,omitempty"`
	Branch         string `json:"branch,omitempty"`
	PermissionMode string `json:"permission_mode,omitempty"`
	Agent          string `json:"agent,omitempty"`
	// Autonomous is true for a run started without the owner: it skips
	// the owner's gates and questions, recording each. It never changes
	// an agent's permissions.
	Autonomous bool `json:"autonomous,omitempty"`
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
//	message:   Step, Session, Actor (the sender), To, Text, Summary, Structured (a session sent another a message)
//	tokens:    Step, Session, Tokens (a reading of the session's token use; the latest counts)
//	findings:  Findings (a plan reviewer's findings, maybe none; each findings entry
//	           opens the next review round)
//	verdict:   Verdicts (a critic's verdicts on findings of the latest round)
//	gated:     Changes (run gate ended the plan: its change sets, by ID; the
//	           latest gated entry wins)
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
	// the arbiter passed.
	Decision string `json:"decision,omitempty"`
	// Handle is an asked entry's question node handle.
	Handle     string    `json:"handle,omitempty"`
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
}

// Finding is one finding of a plan review round: node Node (its handle
// Handle when it was recorded) has a problem of kind Kind, for Reason. ID
// is f1, f2, ... in the round.
type Finding struct {
	ID     string `json:"id"`
	Node   string `json:"node"`
	Handle string `json:"handle,omitempty"`
	Kind   string `json:"kind"`
	Reason string `json:"reason"`
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
type Round struct {
	Number   int
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
	StartedAt      time.Time
	// Autonomous is the started entry's Autonomous: the run goes on
	// without the owner, and every agent it launches is told so.
	Autonomous bool
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
	// Blocked is what the current agent step waits on, from its last
	// blocked entry; nil once an unblocked, its step's ended or a
	// finished entry follows. BlockedStep and BlockedSession name the
	// step and session it belongs to.
	Blocked        *Block
	BlockedStep    string
	BlockedSession string
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
	// RetryBase holds each step's Visits as they stood at the run's last
	// retried entry. Only visits past it count against max_visits, so a
	// retry gives every step its full cap again; visit numbers still
	// count from the run's start. Nil until a retry.
	RetryBase map[string]int
	// retryNext is the step the last retried entry reopened the run at,
	// until its next entered entry marks that visit as the retry.
	retryNext  string
	retryVisit map[string]int
	// blockAsked is 1 + the index in Asked of the asked entry for the
	// current block, 0 when the block has none (or there is no block).
	blockAsked int
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
}

// State folds the journal.
func (j Journal) State() State {
	st := State{ID: j.ID, Visits: map[string]int{}, Session: map[string]string{}, Tokens: map[string]TokenUse{}}
	for _, e := range j.Entries {
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
				st.Autonomous = s.Autonomous
			}
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
			st.retryNext = ""
		case TypeSession:
			st.Session[e.Step] = e.Session
			st.addSession(e.Step, e.Session)
		case TypeBlocked:
			b := Block{}
			if e.Block != nil {
				b = *e.Block
			}
			st.Blocked, st.BlockedStep, st.BlockedSession = &b, e.Step, e.Session
			st.blockAsked = 0
		case TypeUnblocked:
			st.clearBlock()
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
			if st.BlockedStep == e.Step {
				st.clearBlock()
			}
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
		case TypeAsked:
			st.Asked = append(st.Asked, Asked{
				Step: e.Step, Session: e.Session, Question: e.Question, Handle: e.Handle, Text: e.Detail, At: e.At,
			})
			if st.Blocked != nil && st.BlockedSession == e.Session {
				st.blockAsked = len(st.Asked)
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
			st.Gated = e.Changes
		case TypeFindings:
			st.Rounds = append(st.Rounds, Round{Number: len(st.Rounds) + 1, At: e.At, Findings: e.Findings})
		case TypeVerdict:
			if n := len(st.Rounds); n > 0 {
				st.Rounds[n-1].Verdicts = append(st.Rounds[n-1].Verdicts, e.Verdicts...)
			}
		case TypeStopped:
			st.Stopped = e.Detail
		case TypeFinished:
			st.clearBlock()
			st.Ended = true
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

// BlockAsked is the asked entry of the current block: the question the
// engine raised for it, so a session still blocked is never asked twice
// for one block. It has none when the run is not blocked, or its block
// has not been asked yet.
func (st State) BlockAsked() (Asked, bool) {
	if st.Blocked == nil || st.blockAsked == 0 {
		return Asked{}, false
	}
	return st.Asked[st.blockAsked-1], true
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
	st.Blocked, st.BlockedStep, st.BlockedSession = nil, "", ""
	st.blockAsked = 0
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

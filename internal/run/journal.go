// Package run keeps the journal of a workflow run: an append-only JSONL file
// per run, folded to a state. It is work-layer state and lives outside the
// model's log. It runs no program and names no agent.
package run

import "time"

// Entry types.
const (
	TypeStarted   = "started"
	TypeEntered   = "entered"
	TypeSession   = "session"
	TypeEnded     = "ended"
	TypeStopped   = "stopped"
	TypeFinished  = "finished"
	TypeBlocked   = "blocked"
	TypeUnblocked = "unblocked"
	TypeAmended   = "amended"
	TypeRetried   = "retried"
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
	// RetryBase holds each step's Visits as they stood at the run's last
	// retried entry. Only visits past it count against max_visits, so a
	// retry gives every step its full cap again; visit numbers still
	// count from the run's start. Nil until a retry.
	RetryBase map[string]int
	// retryNext is the step the last retried entry reopened the run at,
	// until its next entered entry marks that visit as the retry.
	retryNext  string
	retryVisit map[string]int
}

// State folds the journal.
func (j Journal) State() State {
	st := State{ID: j.ID, Visits: map[string]int{}, Session: map[string]string{}}
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
		case TypeBlocked:
			b := Block{}
			if e.Block != nil {
				b = *e.Block
			}
			st.Blocked, st.BlockedStep, st.BlockedSession = &b, e.Step, e.Session
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
}

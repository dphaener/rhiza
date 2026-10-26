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
			}
		case TypeEntered:
			st.Visits[e.Step]++
			st.Current = e.Step
			st.CurrentVisit = e.Visit
			if st.CurrentVisit == 0 {
				st.CurrentVisit = st.Visits[e.Step]
			}
			st.Stopped = ""
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
			})
			delete(st.Session, e.Step)
			if st.BlockedStep == e.Step {
				st.clearBlock()
			}
			st.Current = e.Next
			st.CurrentVisit = 0
			st.Stopped = ""
		case TypeStopped:
			st.Stopped = e.Detail
		case TypeFinished:
			st.clearBlock()
			st.Ended = true
			st.Result = e.End
			st.Detail = e.Detail
			st.Current = ""
			st.CurrentVisit = 0
		}
	}
	return st
}

func (st *State) clearBlock() {
	st.Blocked, st.BlockedStep, st.BlockedSession = nil, "", ""
}

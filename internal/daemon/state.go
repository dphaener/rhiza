package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// StateFile is the state file's name, under Dir.
const StateFile = "daemon.state.json"

// LiveRun is a run the daemon sees as live: running or interrupted.
type LiveRun struct {
	ID     string `json:"id"`
	Change string `json:"change"`
	// State is running or interrupted (see run.State.Liveness).
	State string `json:"state"`
}

// Event is something the daemon last started or resumed, and when.
type Event struct {
	// Change is the change set a started run runs; Run is the resumed run.
	Change string    `json:"change,omitempty"`
	Run    string    `json:"run,omitempty"`
	At     time.Time `json:"at"`
	// Runs is, for an unjournaled start, how many runs its change set had
	// at the start, so a later daemon tells a new journal from an old one.
	Runs int `json:"runs,omitempty"`
}

// State is what a daemon keeps for `rhiza daemon status`: written at start
// and after every poll, and left in place when the daemon stops, with
// Stopped set. Whether the daemon runs now is the lock's to say, not this
// file's.
type State struct {
	Process Holder `json:"process"`
	// Version is the version of the binary the daemon runs, recorded when
	// it starts.
	Version string `json:"version,omitempty"`
	// Commit is the git commit that binary was built from, recorded when
	// it starts; health compares it with origin/main.
	Commit   string    `json:"commit,omitempty"`
	LastPoll time.Time `json:"last_poll"`
	Stopped  time.Time `json:"stopped,omitempty"`
	MaxRuns  int       `json:"max_runs"`
	// Settings are the effective settings of the last poll, each with the
	// source it came from (default, git config or flag).
	Settings    []SettingValue `json:"settings,omitempty"`
	Autonomous  bool           `json:"autonomous"`
	Runs        []LiveRun      `json:"runs"`
	LastStarted *Event         `json:"last_started,omitempty"`
	LastResumed *Event         `json:"last_resumed,omitempty"`
	// Unjournaled are the starts whose run never wrote a journal, by
	// change set: the daemon does not start them again until the hold
	// ends, and a daemon that takes over (a re-exec) keeps the holds.
	Unjournaled []Event `json:"unjournaled,omitempty"`
	// Skips are the accepted change sets the last poll did not start, each
	// with its reason, in the order the poll worked them out.
	Skips []Skip `json:"skips,omitempty"`
	// ProjectSkips are the registered projects the last poll did not work,
	// each with why, in the root's order.
	ProjectSkips []ProjectSkip `json:"project_skips,omitempty"`
	// Notified are the IDs of the inbox items the daemon sent a Rex
	// notification for and that still wait, so a later poll, or a later
	// daemon, sends none again while they stay.
	Notified []string `json:"notified,omitempty"`
	// Hooked are the IDs of the inbox items the notify command ran for and
	// that still wait. It is apart from Notified so the hook runs with or
	// without Rex.
	Hooked []string `json:"hooked,omitempty"`
	// Nudged are the stalls the daemon sent a message for and that still
	// last, each as session@report-time, so a session still stalled at the
	// same report is nudged once; a new report is a new stall.
	Nudged []string `json:"nudged,omitempty"`
	// Watchdog is what the watchdog keeps per live run: the progress it
	// last acted at, whether it nudged, and the retry budget it spent.
	Watchdog []WatchdogRun `json:"watchdog,omitempty"`
	// WatchdogHolds are the stalled steps the last poll left to the owner
	// (an attended run, or a spent budget), each an inbox item while it lasts.
	WatchdogHolds []WatchdogHold `json:"watchdog_holds,omitempty"`
	// PollBehind is set while the last poll took longer than the watchdog's
	// multiple of the poll interval, and absent once a poll is within it.
	PollBehind *PollBehind `json:"poll_behind,omitempty"`
	// Events are the IDs of the lifecycle moves the event command has been
	// handled for, run or seeded; EventsSeeded says the first run recorded
	// every move then in the log, so history is not replayed.
	Events       []string `json:"events,omitempty"`
	EventsSeeded bool     `json:"events_seeded,omitempty"`
	// DiskLow is set while the last poll found free space under the
	// daemon's minimum, and absent once a poll finds enough.
	DiskLow *DiskLow `json:"disk_low,omitempty"`
	// UsageHold is set while the last poll held the daemon's launches for
	// usage, and absent once a poll no longer does.
	UsageHold *UsageHold `json:"usage_hold,omitempty"`
	// UsageWindows are the usage windows the last poll found, one per kind:
	// the next poll, and health, build each window from its saved start (or
	// from its reset once it has passed), not from wherever a lookback
	// happens to fall, so a window does not re-anchor as time goes on.
	UsageWindows []UsageWindow `json:"usage_windows,omitempty"`
	// Learned are the window limits the daemon learned from usage-limit
	// hits, one per window kind, the latest hit's.
	Learned []LearnedLimit `json:"learned,omitempty"`
	// Hooks are the notify and event commands the daemon took at start,
	// each with where it came from (flag, git config or unset).
	Hooks []HookValue `json:"hooks,omitempty"`
}

// HookValue is one hook command as the daemon resolved it at start.
type HookValue struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Source  string `json:"source"`
}

// DiskLow is a low-disk finding of a poll.
type DiskLow struct {
	// Since is the first of the polls in a row that found it low.
	Since time.Time `json:"since"`
	// Free and Min are the free bytes and the minimum, as the last poll saw.
	Free uint64 `json:"free"`
	Min  uint64 `json:"min"`
}

// UsageHold is a usage-reserve hold found by a poll.
type UsageHold struct {
	// Since is the first of the polls in a row that held.
	Since time.Time `json:"since"`
	// Window is the held window's kind (5h or 7d).
	Window    string    `json:"window"`
	Used      int64     `json:"used"`
	Projected int64     `json:"projected"`
	Limit     int64     `json:"limit"`
	Reserve   int       `json:"reserve"`
	Reset     time.Time `json:"reset"`
}

// UsageWindow is the span of one usage window a poll found.
type UsageWindow struct {
	Kind  string    `json:"kind"`
	Start time.Time `json:"start"`
	Reset time.Time `json:"reset"`
	// Used is the window's tokens as that poll counted them; health reports
	// it rather than reading every transcript again.
	Used int64 `json:"used,omitempty"`
}

// WriteState replaces the state file, whole.
func WriteState(common string, s State) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(Dir(common), StateFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ReadState is the state file's contents, and false if there is none or it
// does not parse.
func ReadState(common string) (State, bool) {
	var s State
	b, err := os.ReadFile(filepath.Join(Dir(common), StateFile))
	if err != nil || json.Unmarshal(b, &s) != nil {
		return State{}, false
	}
	return s, true
}

// SettingValue is one setting as a poll resolved it.
type SettingValue struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

// Project skip reasons, fixed: the state file and its readers match on them.
const (
	ProjectNotInitialised = "not initialised"
	ProjectNoGit          = "not a git checkout"
	ProjectHeld           = "held by another daemon"
)

// ProjectSkip is a registered project the daemon does not work, and why.
type ProjectSkip struct {
	Project string `json:"project"`
	Reason  string `json:"reason"`
	// Detail is the error behind a reason that has one.
	Detail string `json:"detail,omitempty"`
}

// WatchdogRun is the watchdog's memory of one run.
type WatchdogRun struct {
	Run string `json:"run"`
	// Progress is the last journal or session progress the watchdog saw;
	// a later one resets Nudged.
	Progress time.Time `json:"progress"`
	Nudged   bool      `json:"nudged,omitempty"`
	// Spent is the retries charged to the run's watchdog budget.
	Spent int `json:"spent,omitempty"`
	// Attempts journals each recovery charged to the budget, oldest first,
	// so a restarted daemon still knows what was tried.
	Attempts []WatchdogAttempt `json:"attempts,omitempty"`
	// CIHead is the PR head whose failed CI jobs the watchdog reran, and
	// CIResumed that the run was resumed once they passed: a head is
	// rerun once.
	CIHead    string `json:"ci_head,omitempty"`
	CIResumed bool   `json:"ci_resumed,omitempty"`
	// PRGone is the time of the run's last journal entry when the watchdog
	// found its change set with no open PR: while the journal ends there,
	// the green-but-unmerged pass reads that PR no more.
	PRGone time.Time `json:"pr_gone,omitzero"`
}

// WatchdogAttempt is one recovery the watchdog tried on a run.
type WatchdogAttempt struct {
	At     time.Time `json:"at"`
	Action string    `json:"action"`
}

// WatchdogHold is a stalled step the watchdog only reports.
type WatchdogHold struct {
	Run    string    `json:"run"`
	Change string    `json:"change"`
	Step   string    `json:"step"`
	Since  time.Time `json:"since"`
	Detail string    `json:"detail"`
}

// PollBehind is a poll that ran past its limit.
type PollBehind struct {
	At    time.Time     `json:"at"`
	Took  time.Duration `json:"took"`
	Limit time.Duration `json:"limit"`
}

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
	// Notified are the IDs of the inbox items the daemon sent a Rex
	// notification for and that still wait, so a later poll, or a later
	// daemon, sends none again while they stay.
	Notified []string `json:"notified,omitempty"`
	// Hooked are the IDs of the inbox items the notify command ran for and
	// that still wait. It is apart from Notified so the hook runs with or
	// without Rex.
	Hooked []string `json:"hooked,omitempty"`
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

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
}

// State is what a daemon keeps for `rhiza daemon status`: written at start
// and after every poll, and left in place when the daemon stops, with
// Stopped set. Whether the daemon runs now is the lock's to say, not this
// file's.
type State struct {
	Process     Holder    `json:"process"`
	LastPoll    time.Time `json:"last_poll"`
	Stopped     time.Time `json:"stopped,omitempty"`
	MaxRuns     int       `json:"max_runs"`
	Autonomous  bool      `json:"autonomous"`
	Runs        []LiveRun `json:"runs"`
	LastStarted *Event    `json:"last_started,omitempty"`
	LastResumed *Event    `json:"last_resumed,omitempty"`
	// Unjournaled are the starts whose run never wrote a journal, by
	// change set: the daemon does not start them again while it runs.
	Unjournaled []Event `json:"unjournaled,omitempty"`
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

package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// StartsFile is the queue of runs the owner asked for, under Dir: one per
// clone, shared by its worktrees, never committed. rhiza run start writes
// it; the daemon starts what it holds, as detached headless processes.
const StartsFile = "starts.json"

// StartsLockFile is the lock every writer of the queue holds around its
// read and write, under Dir.
const StartsLockFile = "starts.lock"

// Start is one queued run start: the arguments of the rhiza run start that
// queued it, for the daemon's detached rhiza run start --foreground.
type Start struct {
	Change         string    `json:"change"`
	Workflow       string    `json:"workflow,omitempty"`
	PermissionMode string    `json:"permission_mode,omitempty"`
	Agent          string    `json:"agent,omitempty"`
	NoWait         bool      `json:"no_wait,omitempty"`
	Autonomous     bool      `json:"autonomous,omitempty"`
	At             time.Time `json:"at"`
}

type startsDoc struct {
	Starts []Start `json:"starts"`
}

// ReadStarts is the queue, oldest first; none when there is no file.
func ReadStarts(common string) ([]Start, error) {
	b, err := os.ReadFile(filepath.Join(Dir(common), StartsFile))
	if errors.Is(err, os.ErrNotExist) {
		return []Start{}, nil
	}
	if err != nil {
		return nil, err
	}
	var d startsDoc
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	if d.Starts == nil {
		d.Starts = []Start{}
	}
	return d.Starts, nil
}

// UpdateStarts reads the queue, passes it to update and writes back what
// update returns, holding the queue's lock throughout. update returning an
// error writes nothing.
func UpdateStarts(common string, update func([]Start) ([]Start, error)) error {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return err
	}
	unlock, err := openWait(filepath.Join(Dir(common), StartsLockFile))
	if err != nil {
		return err
	}
	defer unlock()
	starts, err := ReadStarts(common)
	if err != nil {
		return err
	}
	starts, err = update(starts)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(startsDoc{Starts: starts}, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(Dir(common), StartsFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

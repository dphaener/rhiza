package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ContinuesFile is the queue of run continuations the owner asked for, under
// Dir: one per clone, shared by its worktrees, never committed. rhiza run
// resume, approve and request-changes write it; the daemon continues what it
// holds, as detached headless processes.
const ContinuesFile = "continues.json"

// ContinuesLockFile is the lock every writer of the queue holds around its
// read and write, under Dir.
const ContinuesLockFile = "continues.lock"

// Continue is one queued continuation: the run to drive on, for the daemon's
// detached rhiza run resume --foreground.
type Continue struct {
	Run    string    `json:"run"`
	Retry  bool      `json:"retry,omitempty"`
	NoWait bool      `json:"no_wait,omitempty"`
	At     time.Time `json:"at"`
}

type continuesDoc struct {
	Continues []Continue `json:"continues"`
}

// ReadContinues is the queue, oldest first; none when there is no file.
func ReadContinues(common string) ([]Continue, error) {
	b, err := os.ReadFile(filepath.Join(Dir(common), ContinuesFile))
	if errors.Is(err, os.ErrNotExist) {
		return []Continue{}, nil
	}
	if err != nil {
		return nil, err
	}
	var d continuesDoc
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	if d.Continues == nil {
		d.Continues = []Continue{}
	}
	return d.Continues, nil
}

// UpdateContinues reads the queue, passes it to update and writes back what
// update returns, holding the queue's lock throughout. update returning an
// error writes nothing.
func UpdateContinues(common string, update func([]Continue) ([]Continue, error)) error {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return err
	}
	unlock, err := openWait(filepath.Join(Dir(common), ContinuesLockFile))
	if err != nil {
		return err
	}
	defer unlock()
	cs, err := ReadContinues(common)
	if err != nil {
		return err
	}
	cs, err = update(cs)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(continuesDoc{Continues: cs}, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(Dir(common), ContinuesFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

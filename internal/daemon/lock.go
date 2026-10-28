// Package daemon holds what the rhiza daemon keeps in the git common dir:
// its lock, which makes it one per repository however many worktrees the
// repository has.
package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// File is the lock file's name, under Dir.
const File = "daemon.lock"

// Dir is the folder the daemon keeps its files in, under the git common
// dir, beside the run journals and the session records.
func Dir(common string) string { return filepath.Join(common, "rhiza") }

// Holder is the process a lock file names: the daemon that took it.
type Holder struct {
	Host    string    `json:"host"`
	Pid     int       `json:"pid"`
	Started time.Time `json:"started"`
}

// String is the holder as "host:pid".
func (h Holder) String() string { return h.Host + ":" + strconv.Itoa(h.Pid) }

// Lock is the daemon lock of one repository, held until Unlock.
type Lock struct {
	path   string
	unlock func()
}

// Held is the error TryLock returns when another daemon holds the lock:
// Holder names it, as far as its lock file does.
type Held struct{ Holder Holder }

func (e *Held) Error() string {
	if e.Holder.Pid <= 0 {
		return "another daemon holds the lock"
	}
	return "another daemon holds the lock: process " + e.Holder.String()
}

// TryLock takes the daemon lock under the git common dir without waiting,
// and writes this process into the lock file. If a live daemon holds it,
// the error is a *Held. A lock file left by a process that is gone is no
// obstacle: the lock is the operating system's, which it releases when its
// holder dies, so the file's contents only say who held it last and are
// replaced. Where the system has no such lock (see lock_other.go), a file
// naming a pid that runs on this host counts as held.
func TryLock(common string, now time.Time) (*Lock, error) {
	dir := Dir(common)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, File)
	f, unlock, err := open(path)
	if err != nil {
		var h *Held
		if errors.As(err, &h) {
			h.Holder = ReadHolder(common)
		}
		return nil, err
	}
	host, _ := os.Hostname()
	b, err := json.Marshal(Holder{Host: host, Pid: os.Getpid(), Started: now.UTC()})
	if err == nil {
		err = f.Truncate(0)
	}
	if err == nil {
		_, err = f.WriteAt(append(b, '\n'), 0)
	}
	if err != nil {
		unlock()
		return nil, err
	}
	return &Lock{path: path, unlock: unlock}, nil
}

// Unlock releases the lock and empties the lock file, so it names no one.
func (l *Lock) Unlock() {
	if l == nil || l.unlock == nil {
		return
	}
	os.Truncate(l.path, 0)
	l.unlock()
	l.unlock = nil
}

// ReadHolder is the process the lock file names, the zero Holder if it
// names none. It does not say whether that process still runs.
func ReadHolder(common string) Holder {
	var h Holder
	b, err := os.ReadFile(filepath.Join(Dir(common), File))
	if err != nil {
		return Holder{}
	}
	if json.Unmarshal(b, &h) != nil {
		return Holder{}
	}
	return h
}

// Running reports whether a daemon holds the lock now, and who its lock
// file names. It takes the lock and drops it again to find out.
func Running(common string) (Holder, bool) {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return Holder{}, false
	}
	_, unlock, err := open(filepath.Join(Dir(common), File))
	if err == nil {
		unlock()
		return Holder{}, false
	}
	return ReadHolder(common), true
}

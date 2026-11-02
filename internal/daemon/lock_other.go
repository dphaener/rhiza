//go:build !unix

package daemon

import (
	"encoding/json"
	"os"

	"github.com/dphaener/rhiza/internal/run"
)

// open has no flock to use here: the lock file names the holder, and a
// holder whose pid runs on this host (or on another host) holds the lock.
func open(path string) (*os.File, func(), error) {
	if b, err := os.ReadFile(path); err == nil {
		var h Holder
		host, _ := os.Hostname()
		if json.Unmarshal(b, &h) == nil && h.Pid > 0 && (h.Host != host || run.PidAlive(h.Pid)) {
			return nil, nil, &Held{}
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

// openWait has no flock to take here: writers of the file it guards are
// not kept apart.
func openWait(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	return func() { f.Close() }, nil
}

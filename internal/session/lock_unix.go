//go:build unix

package session

import (
	"os"
	"path/filepath"
	"syscall"
)

// lock takes an exclusive flock on the record's sibling lock file,
// .<id>.lock, and returns the function that releases it. The lock file is
// never removed: removing it while another writer waits on it would let two
// writers in at once.
func (s *Store) lock(id string) (func(), error) {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(s.dir, "."+id+".lock"), os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	for {
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		f.Close()
		return nil, err
	}
	return func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, nil
}

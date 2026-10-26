//go:build unix

package run

import (
	"os"
	"path/filepath"
	"syscall"
)

// TryLock takes a non-blocking exclusive flock on the run's sibling lock
// file, .<id>.lock. ok is false when another holder has it. The lock file is
// never removed.
func (s *Store) TryLock(id string) (unlock func(), ok bool, err error) {
	if err := validID(id); err != nil {
		return nil, false, err
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return nil, false, err
	}
	f, err := os.OpenFile(filepath.Join(s.dir, "."+id+".lock"), os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, false, err
	}
	for {
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != syscall.EINTR {
			break
		}
	}
	if err != nil {
		f.Close()
		if err == syscall.EWOULDBLOCK {
			return nil, false, nil
		}
		return nil, false, err
	}
	return func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, true, nil
}

// Live reports whether another holder has the run's lock.
func (s *Store) Live(id string) bool {
	unlock, ok, err := s.TryLock(id)
	if err != nil {
		return false
	}
	if ok {
		unlock()
		return false
	}
	return true
}

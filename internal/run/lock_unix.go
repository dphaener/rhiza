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

// LockChange takes a blocking exclusive flock on the change set's lock file,
// .change-<change>.lock, and returns its unlock. A start or a reopening of a
// run holds it while it checks the change set's other runs and writes its
// journal, so two never both begin. The lock file is never removed.
func (s *Store) LockChange(change string) (unlock func(), err error) {
	if err := validID(change); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(s.dir, ".change-"+change+".lock"), os.O_RDWR|os.O_CREATE, 0o644)
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

// LockJournal takes a blocking exclusive flock on the run's journal lock file,
// .<id>.journal.lock, and returns its unlock. A caller that numbers entries
// from the journal's state holds it from reading that state to appending, so
// two concurrent callers never number from the same state. It is not the
// run's liveness lock, which the engine holds for the run's life. The lock
// file is never removed.
func (s *Store) LockJournal(id string) (unlock func(), err error) {
	if err := validID(id); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(s.dir, "."+id+".journal.lock"), os.O_RDWR|os.O_CREATE, 0o644)
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

// PidAlive reports whether a process with the pid exists on this host.
func PidAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}

// Terminate asks the process with the pid to stop; a process that is gone
// is not an error. A process that leads its process group (a detached run
// leads its own session) is stopped with its whole group, so the step
// children it started do not outlive it.
func Terminate(pid int) error {
	target := pid
	if pgid, err := syscall.Getpgid(pid); err == nil && pgid == pid {
		target = -pid
	}
	if err := syscall.Kill(target, syscall.SIGTERM); err != nil && err != syscall.ESRCH {
		return err
	}
	return nil
}

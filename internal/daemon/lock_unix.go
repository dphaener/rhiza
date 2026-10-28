//go:build unix

package daemon

import (
	"os"
	"syscall"
)

// open takes a non-blocking exclusive flock on the lock file, creating it.
// The system drops the flock when its holder dies, however it dies.
func open(path string) (*os.File, func(), error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, &Held{}
		}
		return nil, nil, err
	}
	return f, func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, nil
}

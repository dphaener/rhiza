package daemon

import (
	"os"
	"time"
)

// Image is what tells two builds of one executable apart: its size and
// modification time.
type Image struct {
	Size int64
	Mod  time.Time
}

// StatImage reads the Image of the file at path.
func StatImage(path string) (Image, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Image{}, err
	}
	return Image{Size: fi.Size(), Mod: fi.ModTime()}, nil
}

// Reexec lets a daemon take over from itself when its executable was
// rebuilt. Stat and Exec are injectable; NewReexec fills in the real ones.
type Reexec struct {
	Path string   // the executable
	Args []string // the daemon's argv, program name first
	Env  []string
	Stat func(path string) (Image, error)
	Exec func(path string, argv, env []string) error

	start Image
}

// NewReexec records the executable's Image now, as the one the daemon
// started with. It returns nil when the executable cannot be found or
// statted, which leaves the daemon without re-exec.
func NewReexec() *Reexec {
	path, err := os.Executable()
	if err != nil {
		return nil
	}
	r := &Reexec{Path: path, Args: os.Args, Env: os.Environ(), Stat: StatImage, Exec: execSelf}
	if r.start, err = r.Stat(path); err != nil {
		return nil
	}
	return r
}

// Start records img as the Image the daemon started with.
func (r *Reexec) Start(img Image) { r.start = img }

// Due reports whether the executable changed since the daemon started and
// no start is pending: a pending start delays the exec to a later poll. An
// executable that cannot be statted (mid-rebuild) is not due.
func (r *Reexec) Due(pending int) bool {
	if r == nil || pending > 0 {
		return false
	}
	img, err := r.Stat(r.Path)
	return err == nil && img != r.start
}

// Do re-executes the daemon with its own arguments. It releases lock first,
// so that the new image takes it; if Exec returns (it fails, or is a test's
// fake), Do takes the lock again and returns it with the error.
func (r *Reexec) Do(common string, lock *Lock, now time.Time) (*Lock, error) {
	lock.Unlock()
	err := r.Exec(r.Path, r.Args, r.Env)
	l, lerr := TryLock(common, now)
	if lerr != nil {
		return nil, lerr
	}
	return l, err
}

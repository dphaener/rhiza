//go:build !unix

package run

// TryLock has no flock to use here: it always succeeds, so a run is never
// seen as live.
func (s *Store) TryLock(id string) (unlock func(), ok bool, err error) {
	if err := validID(id); err != nil {
		return nil, false, err
	}
	return func() {}, true, nil
}

// Live reports whether another holder has the run's lock; here, never.
func (s *Store) Live(string) bool { return false }

// LockChange has no flock to use here: it holds nothing.
func (s *Store) LockChange(change string) (unlock func(), err error) {
	if err := validID(change); err != nil {
		return nil, err
	}
	return func() {}, nil
}

// LockJournal has no flock to use here: it holds nothing.
func (s *Store) LockJournal(id string) (unlock func(), err error) {
	if err := validID(id); err != nil {
		return nil, err
	}
	return func() {}, nil
}

// PidAlive cannot look a process up here: it reports every pid as alive.
func PidAlive(pid int) bool { return pid > 0 }

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

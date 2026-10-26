//go:build !unix

package session

import "sync"

// lockMu stands in for flock where there is none: it serialises writers in
// one process only.
var lockMu sync.Mutex

func (s *Store) lock(string) (func(), error) {
	lockMu.Lock()
	return lockMu.Unlock, nil
}

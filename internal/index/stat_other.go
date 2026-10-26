//go:build !darwin && !linux

package index

import "io/fs"

// statKey has no os fields to read on this platform: the stat fingerprint
// uses names, modification times and modes alone.
func statKey(fs.FileInfo) (osKey, bool) { return osKey{}, false }

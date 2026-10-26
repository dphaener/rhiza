//go:build darwin

package index

import (
	"io/fs"
	"syscall"
)

// statKey returns the change time, inode and device of an os file's info,
// ok false when fi does not come from the os.
func statKey(fi fs.FileInfo) (k osKey, ok bool) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return osKey{}, false
	}
	return osKey{ctime: st.Ctimespec.Nano(), ino: st.Ino, dev: uint64(st.Dev)}, true
}

//go:build unix

package daemon

import "syscall"

// FreeBytes is the space free to an unprivileged writer on the filesystem
// holding dir.
func FreeBytes(dir string) (uint64, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(dir, &st); err != nil {
		return 0, err
	}
	return uint64(st.Bavail) * uint64(st.Bsize), nil
}

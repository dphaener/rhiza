//go:build unix

package daemon

import "syscall"

func execSelf(path string, argv, env []string) error { return syscall.Exec(path, argv, env) }

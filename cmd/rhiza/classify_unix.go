//go:build unix

package main

import (
	"os/exec"
	"syscall"
)

// ownGroup starts cmd in a process group of its own and makes cancelling it
// kill that whole group, so a classifier's children (a hung `claude -p`
// under sh -c) die with it.
func ownGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}

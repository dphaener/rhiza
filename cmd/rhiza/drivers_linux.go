//go:build linux

package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

// isChildOf reports from /proc whether pid's parent is ppid.
func isChildOf(pid, ppid int) bool {
	b, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return false
	}
	s := string(b)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return false
	}
	f := strings.Fields(s[i+1:]) // state, ppid, ...
	return len(f) > 1 && f[1] == strconv.Itoa(ppid)
}

// reapList waits without blocking, by pid, on each listed pid that is still
// a child of parent, and drops from the list each pid that is reaped or no
// longer a child. It never waits on a pid that is not listed.
func reapList(list string, parent int) {
	driverMu.Lock()
	defer driverMu.Unlock()
	var drop []int
	for _, pid := range readDrivers(list) {
		if !isChildOf(pid, parent) {
			drop = append(drop, pid)
			continue
		}
		var ws syscall.WaitStatus
		p, err := syscall.Wait4(pid, &ws, syscall.WNOHANG, nil)
		if p == pid || err == syscall.ECHILD {
			drop = append(drop, pid)
		}
	}
	if len(drop) > 0 {
		dropLocked(list, drop...)
	}
}

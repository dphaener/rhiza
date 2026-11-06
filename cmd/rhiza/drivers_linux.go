//go:build linux

package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

// statFields returns /proc/<pid>/stat's fields after the command name,
// starting with field 3 (state).
func statFields(pid int) []string {
	b, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return nil
	}
	s := string(b)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return nil
	}
	return strings.Fields(s[i+1:])
}

// procStart is pid's start time, field 22 of /proc/<pid>/stat, or "".
func procStart(pid int) string {
	f := statFields(pid)
	if len(f) < 20 {
		return ""
	}
	return f[19]
}

// isChildOf reports from /proc whether pid's parent is ppid.
func isChildOf(pid, ppid int) bool {
	f := statFields(pid) // state, ppid, ...
	return len(f) > 1 && f[1] == strconv.Itoa(ppid)
}

// reapList waits without blocking, by pid, on each listed pid that the
// running daemon self recorded and that is still its child. It drops each
// entry that is reaped or no longer a child, and every entry another daemon
// recorded, without waiting on it: init reaps those drivers, and the pid may
// now belong to one of self's own children that something else waits on.
// It never waits on a pid that is not listed.
func reapList(list string, self daemonID) {
	driverMu.Lock()
	defer driverMu.Unlock()
	var drop []driverEntry
	for _, e := range readDrivers(list) {
		if e.By != self || self.Start == "" || !isChildOf(e.Pid, self.Pid) {
			drop = append(drop, e)
			continue
		}
		var ws syscall.WaitStatus
		p, err := syscall.Wait4(e.Pid, &ws, syscall.WNOHANG, nil)
		if p == e.Pid || err == syscall.ECHILD {
			drop = append(drop, e)
		}
	}
	if len(drop) > 0 {
		dropLocked(list, drop...)
	}
}

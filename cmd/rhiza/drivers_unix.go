//go:build unix

package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// driverMu guards the driver list and the reaper's scan of it.
var driverMu sync.Mutex

// driverListPath is the driver list beside the run logs: <common>/rhiza/drivers.
func driverListPath(logPath string) string {
	return filepath.Join(filepath.Dir(filepath.Dir(logPath)), "drivers")
}

func readDrivers(list string) []int {
	b, _ := os.ReadFile(list)
	var pids []int
	for _, f := range strings.Fields(string(b)) {
		if n, err := strconv.Atoi(f); err == nil && n > 0 {
			pids = append(pids, n)
		}
	}
	return pids
}

func writeDrivers(list string, pids []int) {
	if len(pids) == 0 {
		os.Remove(list)
		return
	}
	var sb strings.Builder
	for _, p := range pids {
		sb.WriteString(strconv.Itoa(p) + "\n")
	}
	tmp := list + ".tmp"
	if os.WriteFile(tmp, []byte(sb.String()), 0o644) == nil {
		os.Rename(tmp, list)
	}
}

func addDriver(list string, pid int) {
	driverMu.Lock()
	defer driverMu.Unlock()
	os.MkdirAll(filepath.Dir(list), 0o755)
	writeDrivers(list, append(readDrivers(list), pid))
}

func dropDrivers(list string, drop ...int) {
	driverMu.Lock()
	defer driverMu.Unlock()
	dropLocked(list, drop...)
}

func dropLocked(list string, drop ...int) {
	var keep []int
	for _, p := range readDrivers(list) {
		gone := false
		for _, d := range drop {
			gone = gone || d == p
		}
		if !gone {
			keep = append(keep, p)
		}
	}
	writeDrivers(list, keep)
}

// reapDrivers is api.Ctx.Reap.
func reapDrivers(common string) {
	reapList(filepath.Join(common, "rhiza", "drivers"), os.Getpid())
}

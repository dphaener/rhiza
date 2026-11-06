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

// daemonID names the daemon process that recorded an entry: its pid and its
// start time (field 22 of /proc/<pid>/stat). A re-exec in place keeps both;
// a fresh start does not.
type daemonID struct {
	Pid   int
	Start string
}

// selfID is the running daemon's daemonID.
func selfID() daemonID {
	pid := os.Getpid()
	return daemonID{Pid: pid, Start: procStart(pid)}
}

// driverEntry is one line of the driver list: "<driver pid> <daemon pid> <daemon start>".
type driverEntry struct {
	Pid int
	By  daemonID
}

func readDrivers(list string) []driverEntry {
	b, _ := os.ReadFile(list)
	var es []driverEntry
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		pid, err := strconv.Atoi(f[0])
		if err != nil || pid <= 0 {
			continue
		}
		e := driverEntry{Pid: pid}
		if len(f) > 1 {
			e.By.Pid, _ = strconv.Atoi(f[1])
		}
		if len(f) > 2 && f[2] != "-" {
			e.By.Start = f[2]
		}
		es = append(es, e)
	}
	return es
}

func writeDrivers(list string, es []driverEntry) {
	if len(es) == 0 {
		os.Remove(list)
		return
	}
	var sb strings.Builder
	for _, e := range es {
		start := e.By.Start
		if start == "" {
			start = "-"
		}
		sb.WriteString(strconv.Itoa(e.Pid) + " " + strconv.Itoa(e.By.Pid) + " " + start + "\n")
	}
	tmp := list + ".tmp"
	if os.WriteFile(tmp, []byte(sb.String()), 0o644) == nil {
		os.Rename(tmp, list)
	}
}

// addDriver records pid on the list as a driver of the running daemon.
func addDriver(list string, pid int) {
	addEntry(list, driverEntry{Pid: pid, By: selfID()})
}

func addEntry(list string, e driverEntry) {
	driverMu.Lock()
	defer driverMu.Unlock()
	os.MkdirAll(filepath.Dir(list), 0o755)
	writeDrivers(list, append(readDrivers(list), e))
}

// dropDrivers drops the running daemon's entries for the given pids.
func dropDrivers(list string, drop ...int) {
	driverMu.Lock()
	defer driverMu.Unlock()
	self := selfID()
	var es []driverEntry
	for _, p := range drop {
		es = append(es, driverEntry{Pid: p, By: self})
	}
	dropLocked(list, es...)
}

func dropLocked(list string, drop ...driverEntry) {
	var keep []driverEntry
	for _, e := range readDrivers(list) {
		gone := false
		for _, d := range drop {
			gone = gone || d == e
		}
		if !gone {
			keep = append(keep, e)
		}
	}
	writeDrivers(list, keep)
}

// reapDrivers is api.Ctx.Reap.
func reapDrivers(common string) {
	reapList(filepath.Join(common, "rhiza", "drivers"), selfID())
}

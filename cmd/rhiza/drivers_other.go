//go:build unix && !linux

package main

// procStart is empty off Linux, which has no /proc to read it from.
func procStart(pid int) string { return "" }

// reapList does nothing off Linux.
func reapList(list string, self daemonID) {}

//go:build unix && !linux

package main

// reapList does nothing off Linux.
func reapList(list string, parent int) {}

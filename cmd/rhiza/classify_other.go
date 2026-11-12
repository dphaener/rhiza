//go:build !unix

package main

import "os/exec"

// ownGroup leaves cmd as it is: only the shell is killed on cancel.
func ownGroup(cmd *exec.Cmd) {}

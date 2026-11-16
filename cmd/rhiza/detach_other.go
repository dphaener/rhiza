//go:build !unix

package main

import "errors"

// detachProgram returns api.Ctx.Detach; detached children need a unix
// system.
func detachProgram(env []string) func(p apiProgram, logPath string) error {
	return func(p apiProgram, logPath string) error {
		return errors.New("detached runs need a unix system")
	}
}

// reapDrivers is api.Ctx.Reap; there is nothing to reap.
func reapDrivers(common string) {}

// adoptChildren is api.Ctx.Adopt; there is nothing to adopt.
func adoptChildren(common string) {}

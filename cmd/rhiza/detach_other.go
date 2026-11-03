//go:build !unix

package main

import "errors"

// detachProgram is api.Ctx.Detach; detached children need a unix system.
func detachProgram(p apiProgram, logPath string) error {
	return errors.New("detached runs need a unix system")
}

package api

import "os"

// childEnv is the environment of a child an operation starts: ctx.Environ
// (the process's when nil) with extra appended, a later entry for a key
// winning as os/exec reads it. It is never nil, so the child never
// inherits the process's environment by default.
func childEnv(ctx *Ctx, extra ...string) []string {
	base := ctx.Environ
	if base == nil {
		base = os.Environ()
	}
	out := make([]string, 0, len(base)+len(extra))
	return append(append(out, base...), extra...)
}

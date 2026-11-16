package api

import "github.com/dphaener/rhiza/internal/userdir"

// childEnv is the environment of a child an operation starts: ctx.Environ
// (none when nil) with extra appended, a later entry for a key winning as
// os/exec reads it. It is never nil, so the child never inherits the
// process's environment.
func childEnv(ctx *Ctx, extra ...string) []string {
	out := make([]string, 0, len(ctx.Environ)+len(extra))
	return append(append(out, ctx.Environ...), extra...)
}

// tempDir is the temporary folder of the Ctx's environment ($TMPDIR, else
// the platform's default), never the process's: os.MkdirTemp and
// os.CreateTemp take it in place of "".
func tempDir(ctx *Ctx) string {
	return userdir.Temp(ctx.Getenv)
}

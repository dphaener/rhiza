package api

// childEnv is the environment of a child an operation starts: ctx.Environ
// (none when nil) with extra appended, a later entry for a key winning as
// os/exec reads it. It is never nil, so the child never inherits the
// process's environment.
func childEnv(ctx *Ctx, extra ...string) []string {
	out := make([]string, 0, len(ctx.Environ)+len(extra))
	return append(append(out, ctx.Environ...), extra...)
}

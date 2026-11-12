package api

// This file holds the head build ship can land and render with. A repo
// whose own rendering the change set's branch may change (Rhiza's) names a
// build of the PR head in its ship-change workflow: a step builds it and
// the workflow's [env] sets SHIP_HEAD_RHIZA to its path. Ship's land and
// render then run with that binary, since the rhiza that runs ship (often
// an installed one, older than the branch) would leave stale rendered
// files. A repo that names none lands and renders in-process.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// envShipRhiza names the binary ship's land and render run with, built
// from the PR head by a step of the repo's ship-change workflow.
const envShipRhiza = "SHIP_HEAD_RHIZA"

// shipHeadBuild returns the head build the run's workflow names, or ""
// when it names none. A named build that is missing is an error: the step
// that builds it did not run.
func shipHeadBuild(ctx *Ctx, _ shipTarget) (string, error) {
	bin := strings.TrimSpace(ctxGetenv(ctx)(envShipRhiza))
	if bin == "" {
		return "", nil
	}
	if ctx.Exec == nil {
		return "", Errorf("ship: the PR head's build %s cannot be run here: no programs can be run", bin)
	}
	if fi, err := os.Stat(bin); err != nil || fi.IsDir() {
		return "", Errorf("ship: %s names %s, which is not a built rhiza: build it in a step before this one, so nothing was landed, rendered or pushed", envShipRhiza, bin)
	}
	return bin, nil
}

// shipActorFlags are the global flags that make the head build act as
// the in-process path would: the caller's --actor and --for, as given.
// With neither, the build resolves them as this process does (the same
// environment and git config).
func shipActorFlags(ctx *Ctx) []string {
	var f []string
	switch {
	case !ctx.Actor.IsZero():
		f = append(f, "--actor", ctx.Actor.String())
	case ctx.ActorRef != "":
		f = append(f, "--actor", ctx.ActorRef)
	}
	switch {
	case !ctx.For.IsZero():
		f = append(f, "--for", ctx.For.String())
	case ctx.ForRef != "":
		f = append(f, "--for", ctx.ForRef)
	}
	return f
}

// shipRunHead runs the head build bin with args (and the actor flags) in
// ctx.Root, and returns its standard output.
func shipRunHead(ctx *Ctx, bin string, args ...string) ([]byte, error) {
	out, err := ctx.Exec(bin, append(args, shipActorFlags(ctx)...), ctx.Root, nil)
	if err != nil {
		return nil, Errorf("ship: %s with the PR head's build: %s", args[0], headFailure(out, err))
	}
	return out, nil
}

// headFailure describes a failed head-build run: err (which carries the
// exit status, and the child's stderr when it printed any), preceded by the
// code and message of the {"error":{...}} body the child's --json printed
// on standard output, when there is one.
func headFailure(out []byte, err error) string {
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(bytes.TrimSpace(out), &body) == nil && body.Error.Message != "" {
		if body.Error.Code != "" {
			return fmt.Sprintf("%s: %s: %v", body.Error.Code, body.Error.Message, err)
		}
		return fmt.Sprintf("%s: %v", body.Error.Message, err)
	}
	return err.Error()
}

// shipLandWith lands change set t.cs: with the head build bin when given,
// else in-process.
func shipLandWith(ctx *Ctx, t shipTarget, bin string) error {
	if bin != "" {
		_, err := shipRunHead(ctx, bin, "land", "--change", t.cs.String(), "--json")
		return err
	}
	land, _ := Lookup("land")
	lctx := *ctx
	lctx.ChangeSet, lctx.ChangeRef, lctx.Record = t.cs, "", false
	_, err := land.Run(&lctx, LifecycleRequest{})
	return err
}

// shipRenderWith renders in ctx.Root, with the head build bin when given,
// else in-process, and returns the rendered files written, then removed.
func shipRenderWith(ctx *Ctx, bin string) ([]string, error) {
	var rr RenderResponse
	if bin != "" {
		out, err := shipRunHead(ctx, bin, "render", "--json")
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(out, &rr); err != nil {
			return nil, Errorf("ship: read render's output from the PR head's build: %v", err)
		}
	} else {
		r, err := runRender(ctx, RenderRequest{})
		if err != nil {
			return nil, err
		}
		rr = r.(RenderResponse)
	}
	return append(append([]string{}, rr.Written...), rr.Removed...), nil
}

// shipRender is how a ship step renders the change set's worktree: with
// the head build the workflow names, in-process when it names none.
func shipRender(ctx *Ctx, t shipTarget) ([]string, error) {
	bin, err := shipHeadBuild(ctx, t)
	if err != nil {
		return nil, err
	}
	return shipRenderWith(ctx, bin)
}

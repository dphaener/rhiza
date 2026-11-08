package api

// This file holds the head build ship uses in Rhiza's own repo. There the
// change set's branch may change how the model renders, and the rhiza that
// runs ship (often an installed one, older than the branch) would land and
// render the branch with its own rendering, leaving stale rendered files.
// So when the worktree's go.mod names Rhiza's module, ship builds rhiza
// from the PR head into <git common dir>/rhiza/ship/<change set>/rhiza
// (work-layer state, never committed, beside the state file) and runs the
// land and render with that build. Every other repo lands and renders
// in-process, as before.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// rhizaModule is Rhiza's Go module path: a go.mod naming it marks Rhiza's
// own repo.
const rhizaModule = "github.com/dphaener/rhiza"

// goModule is the module path that the go.mod file data declares, or ""
// if it declares none. It reads the first module directive, dropping a
// trailing // comment and quotes.
func goModule(data []byte) string {
	for line := range strings.SplitSeq(string(data), "\n") {
		if i := strings.Index(line, "//"); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		rest, ok := strings.CutPrefix(line, "module")
		if !ok || (rest != "" && rest[0] != ' ' && rest[0] != '\t') {
			continue
		}
		rest = strings.TrimSpace(rest)
		if u, err := strconv.Unquote(rest); err == nil {
			rest = u
		}
		return rest
	}
	return ""
}

// isRhizaRepo reports whether root's go.mod names Rhiza's module. No
// go.mod is not Rhiza's repo; an unreadable one is an error.
func isRhizaRepo(root string) (bool, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, Errorf("ship: read go.mod: %v", err)
	}
	return goModule(data) == rhizaModule, nil
}

// headBuildPath is where ship builds rhiza from the change set's PR head.
func (t shipTarget) headBuildPath() string {
	return filepath.Join(t.common, filepath.FromSlash(shipDir), t.cs.String(), "rhiza")
}

// shipHeadBuild builds rhiza from the head checked out in ctx.Root when
// ctx.Root is Rhiza's own repo, and returns the build's path; elsewhere it
// returns "" and builds nothing. A failed build is an error carrying go's
// output.
func shipHeadBuild(ctx *Ctx, t shipTarget) (string, error) {
	ok, err := isRhizaRepo(ctx.Root)
	if err != nil || !ok {
		return "", err
	}
	if ctx.Exec == nil {
		return "", Errorf("ship: building rhiza from the PR head needs go, and no programs can be run here")
	}
	bin := t.headBuildPath()
	if err := os.MkdirAll(filepath.Dir(bin), 0o755); err != nil {
		return "", Errorf("ship: %v", err)
	}
	if _, err := ctx.Exec("go", []string{"build", "-o", bin, "./cmd/rhiza"}, ctx.Root, nil); err != nil {
		return "", Errorf("ship: could not build rhiza from the PR head in %s, so nothing was landed, rendered or pushed: %v", ctx.Root, err)
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
// rhiza built from the PR head in Rhiza's own repo, in-process elsewhere.
func shipRender(ctx *Ctx, t shipTarget) ([]string, error) {
	bin, err := shipHeadBuild(ctx, t)
	if err != nil {
		return nil, err
	}
	return shipRenderWith(ctx, bin)
}

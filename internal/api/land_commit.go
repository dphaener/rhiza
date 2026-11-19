package api

// This file holds land --commit: land, render, and commit only .rhiza with
// the Change: and Actor: trailers a run's commits carry, so a conductor's
// land commit is never taken for a hand commit (Decision
// conductor-land-commit-trailers). Land and render go through the head
// build when SHIP_HEAD_RHIZA names one (runViaHead, with its arguments as
// ever: --commit is never forwarded), and this process makes the commit
// after both return, so a head built without --commit lands the same way.
// A failed land or render commits nothing.

import "strings"

// LandRequest is the request of land.
type LandRequest struct {
	Commit bool `json:"commit,omitempty" flag:"commit" desc:"Then render, stage .rhiza and commit only .rhiza, with the trailers Change: <the change set> and Actor: <the actor that landed it>; anything else staged stays staged and out of the commit"`
}

// landCommitMessage is the subject of land --commit's commit.
const landCommitMessage = "Land and render the change set"

// landCommit lands the change set (v is land), renders, and commits .rhiza.
func landCommit(ctx *Ctx, v lifecycleVerb) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("land --commit makes a commit, so it takes no --dry-run: nothing was landed or committed").
			WithNext("rhiza land --dry-run")
	}
	if ctx.Git == nil {
		return nil, Errorf("land --commit: git cannot be run here, so nothing was landed")
	}
	out, err := v.run(ctx, LifecycleRequest{})
	if err != nil {
		return nil, err
	}
	resp := out.(LifecycleResponse)
	if _, err := runRender(ctx, RenderRequest{}); err != nil {
		return nil, err
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	actor, err := shipActorHandle(w, resp.ChangeSet, ctx)
	w.closeProjection()
	if err != nil {
		return nil, err
	}
	if _, err := ctx.Git("add", "--", ".rhiza"); err != nil {
		return nil, Errorf("land --commit: %v", err)
	}
	staged, err := ctx.Git("diff", "--cached", "--name-only", "--", ".rhiza")
	if err != nil {
		return nil, Errorf("land --commit: %v", err)
	}
	if strings.TrimSpace(staged) == "" {
		return resp, nil
	}
	if _, err := ctx.Git("commit", "-m", landCommitMessage, "-m",
		"Change: "+resp.ChangeSet.String()+"\nActor: "+actor, "--", ".rhiza"); err != nil {
		return nil, Errorf("land --commit: %v", err)
	}
	if resp.Commit, err = gitHead(ctx); err != nil {
		return nil, err
	}
	return resp, nil
}

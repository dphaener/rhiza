package api

// This file holds ship.rebase, which the built-in ship-change workflow
// runs before landing and again before merging: when origin's main has
// moved past the change set's branch, it merges main into the branch.
// It never rebases and never forces a push, so concurrent runs only ever
// fast-forward origin's branch. Conflicts only in rendered files are
// settled by rendering again, since rendered files follow from the model.
// A conflict in code leaves the merge in progress in the worktree and
// records the conflicted files and main's commits in the ship state, for
// a fixer session to resolve and commit; the next call finds that merge
// and pushes it. While CI is unstable, the project's local gate runs on the
// merged head before the push (the gate step before this one saw only the
// head before main was brought in); a failing gate pushes nothing and
// leaves main merged into the branch, unpushed, for the gate's fixer. A
// pushed head is recorded so the next ship check waits for CI on it.

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
)

func init() {
	register(Operation{
		Name:     "ship.rebase",
		Summary:  "Merge origin's main into a change set's branch when main has moved past it, settling rendered-only conflicts by rendering again, and push it plainly; a conflict in code is left as a merge in progress",
		Write:    true,
		Request:  reflect.TypeFor[ShipRebaseRequest](),
		Response: reflect.TypeFor[ShipRebaseResponse](),
		Run:      runShipRebase,
	})
}

// shipMain is the branch the change set's pull request merges into, as
// ship.merge and ship.close-out assume.
const shipMain = "main"

// shipRebaseMessage is the subject of the commit that renders again after
// a merge, when that changed anything.
const shipRebaseMessage = "Render again after merging main"

// ShipRebaseRequest takes nothing: the branch is the change set's.
type ShipRebaseRequest struct{}

// ShipRebaseResponse is what ship.rebase found and did.
type ShipRebaseResponse struct {
	Change        string             `json:"change" desc:"The change set"`
	Branch        string             `json:"branch" desc:"The change set's branch"`
	Onto          string             `json:"onto" desc:"The commit of origin's main the branch was compared with"`
	Behind        bool               `json:"behind" desc:"True if origin's main had commits the branch did not"`
	Rebased       bool               `json:"rebased" desc:"True if origin's main was merged into the branch and the result is ready to push (the name is kept from when this rebased)"`
	Settled       []string           `json:"settled,omitempty" desc:"The rendered files whose conflicts were settled by rendering again"`
	Conflicts     []string           `json:"conflicts,omitempty" desc:"The files that conflicted when the merge stopped on a conflict outside the rendered files; the merge is left in progress in the worktree"`
	MergeConflict *ShipMergeConflict `json:"merge_conflict,omitempty" desc:"The conflict's context for a fixer: the conflicted files, origin/main's commits with their change sets, and the branch's change set"`
	Rendered      []string           `json:"rendered,omitempty" desc:"The rendered files written or removed by rendering after the merge"`
	Commit        string             `json:"commit,omitempty" desc:"The commit of that rendering; empty if it changed nothing"`
	Pushed        string             `json:"pushed,omitempty" desc:"The merged head pushed to the PR branch"`
	Gate          *ShipGateResponse  `json:"gate,omitempty" desc:"The local gate run on the merged head before the push, while CI is unstable; a failing gate pushed nothing"`
	Kept          string             `json:"kept,omitempty" desc:"After a failing gate, the unpushed head left on the branch: origin's main merged into it, for the gate's fixer"`
	Detail        string             `json:"detail,omitempty" desc:"Why the merge stopped; the command exits 1 then"`
}

// ShipMergeCommit is a commit on origin's main that the branch lacks.
type ShipMergeCommit struct {
	Commit  string   `json:"commit" desc:"The commit"`
	Subject string   `json:"subject" desc:"Its subject line"`
	Changes []string `json:"changes,omitempty" desc:"The change sets its Change: trailers name"`
}

// ShipMergeConflict is what a fixer needs to resolve a merge of main that
// conflicts in code: both sides' commits and change sets.
type ShipMergeConflict struct {
	Branch      string            `json:"branch" desc:"The change set's branch"`
	Change      string            `json:"change" desc:"The branch's change set"`
	Onto        string            `json:"onto" desc:"The origin/main commit being merged"`
	Files       []string          `json:"files" desc:"The files still conflicted, including any rendered ones"`
	MainCommits []ShipMergeCommit `json:"main_commits,omitempty" desc:"origin/main's commits since the merge-base, with their change sets"`
}

// ExitCode is 1 when a conflict or a failing local gate stopped the
// rebase, else 0.
func (r ShipRebaseResponse) ExitCode() int {
	if len(r.Conflicts) > 0 || (r.Gate != nil && r.Gate.ExitCode() != 0) {
		return 1
	}
	return 0
}

// splitConflicts sorts conflicted paths, relative to the top of the git
// checkout, into those under the rendered folder (rendered, which ends in
// a slash) and the others, each in the order given.
func splitConflicts(rendered string, paths []string) (inRendered, other []string) {
	for _, p := range paths {
		if strings.HasPrefix(p, rendered) {
			inRendered = append(inRendered, p)
		} else {
			other = append(other, p)
		}
	}
	return inRendered, other
}

// gitLines is git's output split into its non-empty lines.
func gitLines(out string) []string {
	var lines []string
	for l := range strings.SplitSeq(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

// topPath is a pathspec naming p from the top of the git checkout,
// literally.
func topPath(p string) string {
	return ":(top,literal)" + p
}

// runShipRebase merges origin's main into the change set's branch, checked
// out in ctx.Root, when main has moved past it, then pushes it plainly
// and records the pushed head. It never rebases, never forces a push, never
// pushes to main and never settles a conflict outside the rendered files:
// that merge is left in progress for a fixer.
func runShipRebase(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.rebase"); err != nil {
		return nil, err
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	if t.branch == shipMain {
		return nil, Conflictf("ship rebase merges main into a change set's branch, never into %s", shipMain)
	}
	head, err := ctx.Git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return nil, Errorf("ship rebase: which branch: %v", err)
	}
	if b := strings.TrimSpace(head); b != t.branch {
		return nil, Conflictf("ship rebase merges main into the change set's branch %s, but %s is on %s", t.branch, ctx.Root, b)
	}
	if _, err := ctx.Git("fetch", "-q", "origin"); err != nil {
		return nil, Errorf("%v", err)
	}
	onto, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+shipMain)
	if err != nil {
		return nil, Errorf("ship rebase: origin has no %s branch: %v", shipMain, err)
	}
	resp := ShipRebaseResponse{Change: t.cs.String(), Branch: t.branch, Onto: strings.TrimSpace(onto)}
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	rendered := strings.TrimSpace(prefix) + renderedDir + "/"

	// A merge left in progress by an earlier call (a conflict in code) is
	// the fixer's to finish: while files are still conflicted, say so again.
	if _, err := ctx.Git("rev-parse", "--verify", "--quiet", "MERGE_HEAD"); err == nil {
		out, _ := ctx.Git("diff", "--name-only", "--diff-filter=U")
		conflicts := gitLines(out)
		if len(conflicts) == 0 {
			return nil, Conflictf("a merge is in progress in %s with every conflict resolved: commit it (git commit --no-edit), then run ship rebase again", ctx.Root)
		}
		resp.Behind, resp.Conflicts = true, conflicts
		st, err := readShipState(t.statePath())
		if err != nil {
			return nil, err
		}
		resp.MergeConflict = st.MergeConflict
		resp.Detail = fmt.Sprintf("merging origin/%s into %s is still in progress with conflicts in %s", shipMain, t.branch, strings.Join(conflicts, ", "))
		return resp, nil
	}

	if err := shipFollowPR(ctx, t); err != nil {
		return nil, err
	}
	// origin's branch as fetched, which HEAD must hold, so the plain push
	// only ever fast-forwards.
	lease := ""
	if out, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+t.branch); err == nil {
		lease = strings.TrimSpace(out)
		if _, err := ctx.Git("merge-base", "--is-ancestor", lease, "HEAD"); err != nil {
			return nil, Conflictf("origin's %s has commits that %s does not; pull them before merging main", t.branch, ctx.Root)
		}
	}
	orig, err := gitHead(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := ctx.Git("merge-base", "--is-ancestor", resp.Onto, "HEAD"); err == nil {
		// HEAD holds main. A merge committed on top of origin's branch (a
		// fixer's) still has to be rendered, pushed and recorded.
		if lease == "" || lease == orig {
			return resp, nil
		}
		merges, err := ctx.Git("rev-list", "--merges", "--count", lease+"..HEAD")
		if err != nil || strings.TrimSpace(merges) == "0" {
			return resp, nil
		}
		resp.Behind = true
	} else {
		resp.Behind = true
		conflicts, err := shipMergeOnto(ctx, t, resp.Onto, rendered, &resp)
		if err != nil {
			return nil, err
		}
		if len(conflicts) > 0 {
			return shipLeaveConflict(ctx, t, rendered, conflicts, resp)
		}
	}
	resp.Rebased = true

	// Render once more on the merged head, commit what changed, run the
	// local gate, and push. If any of that fails, the branch goes back to
	// where it was, so a retried step merges again rather than finding an
	// unpushed head; only a failing gate keeps the merged head, unpushed,
	// for the gate's fixer.
	if err := shipFinishRebase(ctx, t, rendered, &resp); err != nil {
		if errors.Is(err, errShipGateFailed) {
			return shipKeepGateFailure(ctx, t, resp)
		}
		if _, rerr := ctx.Git("reset", "-q", "--keep", orig); rerr != nil {
			return nil, Errorf("%v; and could not put %s back at %s: %v", err, t.branch, shortSHA(orig), rerr)
		}
		return nil, err
	}
	if resp.Pushed, err = gitHead(ctx); err != nil {
		return nil, err
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	st.Pushed, st.Head, st.Checks, st.MergeConflict = resp.Pushed, "", nil, nil
	if id := shipRunID(ctx); id != "" {
		st.Run = id
	}
	if err := writeShipState(t.statePath(), st); err != nil {
		return nil, err
	}
	return resp, nil
}

// shipKeepGateFailure leaves the merged head the local gate failed on on
// the branch, unpushed, and records it in the ship state: the gate's fixer
// then sees the failure on the code it failed on (the head before the
// merge already passed the gate), and the fix pushes without forcing,
// since the merge descends from the pushed head. ship gate --failed finds
// it there.
func shipKeepGateFailure(ctx *Ctx, t shipTarget, resp ShipRebaseResponse) (any, error) {
	var err error
	if resp.Kept, err = gitHead(ctx); err != nil {
		return nil, err
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	st.GateFailed = resp.Kept
	if err := writeShipState(t.statePath(), st); err != nil {
		return nil, err
	}
	resp.Detail = fmt.Sprintf("the local gate (%s) failed on %s with origin/%s merged in; nothing was pushed, and the merge is left on the branch, unpushed, for the gate's fixer", resp.Gate.Command, t.branch, shipMain)
	return resp, nil
}

// errShipGateFailed is shipFinishRebase's error when the local gate failed
// on the merged head.
var errShipGateFailed = errors.New("the local gate failed")

// shipLeaveConflict records the context of a merge that conflicts in code
// in the ship state and returns the response. The merge stays in progress.
func shipLeaveConflict(ctx *Ctx, t shipTarget, rendered string, conflicts []string, resp ShipRebaseResponse) (any, error) {
	mc := &ShipMergeConflict{Branch: t.branch, Change: t.cs.String(), Onto: resp.Onto, Files: conflicts}
	if base, err := ctx.Git("merge-base", "HEAD", resp.Onto); err == nil {
		out, err := ctx.Git("log", "--format=%H%x09%s%x09%(trailers:key=Change,valueonly,separator=%x2C)", strings.TrimSpace(base)+".."+resp.Onto)
		if err == nil {
			for _, l := range gitLines(out) {
				f := strings.SplitN(l, "\t", 3)
				c := ShipMergeCommit{Commit: f[0]}
				if len(f) > 1 {
					c.Subject = f[1]
				}
				if len(f) > 2 {
					for ch := range strings.SplitSeq(f[2], ",") {
						if ch = strings.TrimSpace(ch); ch != "" {
							c.Changes = append(c.Changes, ch)
						}
					}
				}
				mc.MainCommits = append(mc.MainCommits, c)
			}
		}
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	st.MergeConflict = mc
	if err := writeShipState(t.statePath(), st); err != nil {
		return nil, err
	}
	resp.Conflicts, resp.MergeConflict = conflicts, mc
	resp.Detail = fmt.Sprintf("merging origin/%s into %s conflicts in %s; the merge is left in progress", shipMain, t.branch, strings.Join(conflicts, ", "))
	return resp, nil
}

// shipFinishRebase renders the merged head, commits what that changed,
// runs the local gate on it while CI is unstable (errShipGateFailed when it
// fails), and pushes it to the change set's branch on origin with a plain
// push: never forced, so it fails rather than drop a head it has not seen,
// and never to main.
func shipFinishRebase(ctx *Ctx, t shipTarget, rendered string, resp *ShipRebaseResponse) error {
	var err error
	if resp.Rendered, err = shipRender(ctx, t); err != nil {
		return err
	}
	if resp.Commit, err = shipCommitRendered(ctx, t, rendered, shipRebaseMessage); err != nil {
		return err
	}
	gate, err := shipLocalGate(ctx, &t)
	if err != nil {
		return err
	}
	if !gate.Skipped {
		resp.Gate = &gate
		if !gate.Passed {
			return errShipGateFailed
		}
	}
	if _, err := ctx.Git("push", "origin", "HEAD:refs/heads/"+t.branch); err != nil {
		return Errorf("%v", err)
	}
	return nil
}

// shipMergeOnto merges onto into HEAD, keeping the branch's history. A
// conflict only under rendered is settled by removing those files and
// rendering again. Any other conflict returns every conflicted path and
// leaves the merge in progress. A merge that fails without a conflict is
// aborted and is an error.
func shipMergeOnto(ctx *Ctx, t shipTarget, onto, rendered string, resp *ShipRebaseResponse) ([]string, error) {
	abort := func() { ctx.Git("merge", "--abort") }
	_, err := ctx.Git("merge", "--no-edit", "-m", "Merge origin/"+shipMain+" into "+t.branch, onto)
	if err == nil {
		return nil, nil
	}
	out, derr := ctx.Git("diff", "--name-only", "--diff-filter=U")
	conflicts := gitLines(out)
	if derr != nil || len(conflicts) == 0 {
		abort()
		return nil, Errorf("ship rebase: %v", err)
	}
	inRendered, other := splitConflicts(rendered, conflicts)
	if len(other) > 0 {
		return conflicts, nil
	}
	for _, p := range inRendered {
		if _, err := ctx.Git("rm", "-q", "-f", "--ignore-unmatch", "--", topPath(p)); err != nil {
			abort()
			return nil, Errorf("ship rebase: %v", err)
		}
		if !slices.Contains(resp.Settled, p) {
			resp.Settled = append(resp.Settled, p)
		}
	}
	if _, err := shipRender(ctx, t); err != nil {
		abort()
		return nil, err
	}
	if _, err := ctx.Git("add", "-A", "--", topPath(rendered)); err != nil {
		abort()
		return nil, Errorf("ship rebase: %v", err)
	}
	if _, err := ctx.Git("commit", "-q", "--no-edit"); err != nil {
		abort()
		return nil, Errorf("ship rebase: %v", err)
	}
	return nil, nil
}

// shipCommitRendered stages the rendered folder (rendered, from the top
// of the checkout) and commits it with message and the Change and Actor
// trailers when anything is staged. It returns the commit, or "".
func shipCommitRendered(ctx *Ctx, t shipTarget, rendered, message string) (string, error) {
	if _, err := ctx.Git("add", "-A", "--", topPath(rendered)); err != nil {
		return "", Errorf("%v", err)
	}
	return shipCommitStaged(ctx, t, topPath(rendered), message)
}

// shipCommitStaged commits what is staged under pathspec, with message and
// the Change and Actor trailers, and returns the commit; nothing staged
// commits nothing and returns "".
func shipCommitStaged(ctx *Ctx, t shipTarget, pathspec, message string) (string, error) {
	staged, err := ctx.Git("diff", "--cached", "--name-only", "--", pathspec)
	if err != nil {
		return "", Errorf("%v", err)
	}
	if strings.TrimSpace(staged) == "" {
		return "", nil
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return "", err
	}
	actor, err := shipActorHandle(w, t.cs, ctx)
	if err != nil {
		return "", err
	}
	if _, err := ctx.Git("commit", "-q", "-m", message, "-m", "Change: "+t.cs.String()+"\nActor: "+actor, "--", pathspec); err != nil {
		return "", Errorf("%v", err)
	}
	return gitHead(ctx)
}

// Human says whether the branch was behind, and what the merge did.
func (r ShipRebaseResponse) Human(w io.Writer) {
	switch {
	case !r.Behind:
		fmt.Fprintf(w, "%s already holds origin/%s at %s; nothing to merge.\n", r.Branch, shipMain, shortSHA(r.Onto))
		return
	case r.Gate != nil && !r.Gate.Passed:
		r.Gate.Human(w)
		fmt.Fprintf(w, "Merged origin/%s at %s into %s, but the local gate failed on it: nothing was pushed.\n", shipMain, shortSHA(r.Onto), r.Branch)
		if r.Kept != "" {
			fmt.Fprintf(w, "Left origin/%s merged into %s at %s, unpushed, for the gate's fixer.\n", shipMain, r.Branch, shortSHA(r.Kept))
		} else {
			fmt.Fprintln(w, "The branch is as it was.")
		}
		return
	case len(r.Conflicts) > 0:
		fmt.Fprintf(w, "Merging origin/%s at %s into %s conflicts in:\n", shipMain, shortSHA(r.Onto), r.Branch)
		for _, c := range r.Conflicts {
			fmt.Fprintf(w, "- %s\n", c)
		}
		fmt.Fprintln(w, "The merge is left in progress for a fixer to resolve and commit.")
		return
	}
	fmt.Fprintf(w, "Merged origin/%s at %s into %s.\n", shipMain, shortSHA(r.Onto), r.Branch)
	if len(r.Settled) > 0 {
		fmt.Fprintf(w, "Settled %d rendered %s by rendering again.\n", len(r.Settled), shipPlural(len(r.Settled), "file", "files"))
	}
	if r.Commit != "" {
		fmt.Fprintf(w, "Committed %s.\n", shortSHA(r.Commit))
	}
	if r.Gate != nil {
		fmt.Fprintf(w, "Local gate passed on it (%s).\n", r.Gate.Command)
	}
	fmt.Fprintf(w, "Pushed %s.\n", shortSHA(r.Pushed))
}

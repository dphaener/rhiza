package api

// This file holds ship.rebase, which the built-in ship-change workflow
// runs before landing and again before merging: when origin's main has
// moved past the change set's branch, it rebases the branch onto it.
// Conflicts only in rendered files are settled by rendering again, since
// rendered files follow from the model; any other conflict stops the
// rebase, leaves the branch as it was, and raises a question naming the
// files. A rebased branch is pushed with a lease, and the pushed head is
// recorded so the next ship check waits for CI on it.

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/store"
)

func init() {
	register(Operation{
		Name:     "ship.rebase",
		Summary:  "Rebase a change set's branch onto origin's main when main has moved past it, settling rendered-only conflicts by rendering again, and push it",
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
// a rebase, when that changed anything.
const shipRebaseMessage = "Render again after rebasing onto main"

// shipRebaseLimit bounds the rebase's stops: one per commit replayed, so
// a branch with more commits than this is refused rather than looped on.
const shipRebaseLimit = 1000

// ShipRebaseRequest takes nothing: the branch is the change set's.
type ShipRebaseRequest struct{}

// ShipRebaseResponse is what ship.rebase found and did.
type ShipRebaseResponse struct {
	Change         string   `json:"change" desc:"The change set"`
	Branch         string   `json:"branch" desc:"The change set's branch"`
	Onto           string   `json:"onto" desc:"The commit of origin's main the branch was compared with"`
	Behind         bool     `json:"behind" desc:"True if origin's main had commits the branch did not"`
	Rebased        bool     `json:"rebased" desc:"True if the branch was rebased onto origin's main"`
	Settled        []string `json:"settled,omitempty" desc:"The rendered files whose conflicts were settled by rendering again"`
	Conflicts      []string `json:"conflicts,omitempty" desc:"The files that conflicted when the rebase stopped on a conflict outside the rendered files; the rebase was aborted"`
	Question       string   `json:"question,omitempty" desc:"The handle of the question raised about those conflicts"`
	QuestionChange string   `json:"question_change,omitempty" desc:"The change set the question was written into: the change set itself, or a record change set when it had already landed"`
	Rendered       []string `json:"rendered,omitempty" desc:"The rendered files written or removed by rendering after the rebase"`
	Commit         string   `json:"commit,omitempty" desc:"The commit of that rendering; empty if it changed nothing"`
	Pushed         string   `json:"pushed,omitempty" desc:"The rebased head pushed to the PR branch"`
	Detail         string   `json:"detail,omitempty" desc:"Why the rebase stopped; the command exits 1 then"`
}

// ExitCode is 1 when a conflict stopped the rebase, else 0.
func (r ShipRebaseResponse) ExitCode() int {
	if len(r.Conflicts) > 0 {
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

// runShipRebase rebases the change set's branch, checked out in ctx.Root,
// onto origin's main (or merges main into it, when the branch already
// holds a merge) when main has moved past it, then pushes it with a
// lease and records the pushed head. It never pushes to main and never
// settles a conflict outside the rendered files.
func runShipRebase(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.rebase"); err != nil {
		return nil, err
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	if t.branch == shipMain {
		return nil, Conflictf("ship rebase rebases a change set's branch, never %s", shipMain)
	}
	head, err := ctx.Git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return nil, Errorf("ship rebase: which branch: %v", err)
	}
	if b := strings.TrimSpace(head); b != t.branch {
		return nil, Conflictf("ship rebase rebases the change set's branch %s, but %s is on %s", t.branch, ctx.Root, b)
	}
	if _, err := ctx.Git("fetch", "-q", "origin"); err != nil {
		return nil, Errorf("%v", err)
	}
	if err := shipFollowPR(ctx, t); err != nil {
		return nil, err
	}
	onto, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+shipMain)
	if err != nil {
		return nil, Errorf("ship rebase: origin has no %s branch: %v", shipMain, err)
	}
	resp := ShipRebaseResponse{Change: t.cs.String(), Branch: t.branch, Onto: strings.TrimSpace(onto)}
	if _, err := ctx.Git("merge-base", "--is-ancestor", resp.Onto, "HEAD"); err == nil {
		return resp, nil
	}
	resp.Behind = true

	// The lease: origin's branch as fetched, which the branch must hold,
	// so the forced push never drops a commit pushed from elsewhere.
	lease := ""
	if out, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+t.branch); err == nil {
		lease = strings.TrimSpace(out)
		if _, err := ctx.Git("merge-base", "--is-ancestor", lease, "HEAD"); err != nil {
			return nil, Conflictf("origin's %s has commits that %s does not; pull them before rebasing", t.branch, ctx.Root)
		}
	}
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	rendered := strings.TrimSpace(prefix) + renderedDir + "/"

	orig, err := gitHead(ctx)
	if err != nil {
		return nil, err
	}
	// A branch that already holds a merge (an owner or the orchestrator
	// resolved a conflict by merging main) is merged with main again, not
	// rebased: a rebase would replay the branch's commits flat and bring
	// back the conflicts that merge resolved.
	merges, err := ctx.Git("rev-list", "--merges", "--count", resp.Onto+"..HEAD")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	var conflicts []string
	if strings.TrimSpace(merges) != "0" {
		conflicts, err = shipMergeOnto(ctx, t, resp.Onto, rendered, &resp)
	} else {
		conflicts, err = shipRebaseOnto(ctx, t, resp.Onto, rendered, &resp)
	}
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		resp.Conflicts = conflicts
		resp.Detail = fmt.Sprintf("rebasing %s onto origin/%s conflicts in %s; the rebase was aborted", t.branch, shipMain, strings.Join(conflicts, ", "))
		if err := shipRaiseQuestion(ctx, t, &resp); err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp.Rebased = true

	// Render once more on the rebased head, commit what changed, and push
	// with the lease. If any of that fails, the branch goes back to where
	// it was, so it is still behind and a retried step rebases again
	// rather than finding nothing to do with an unpushed head.
	if err := shipFinishRebase(ctx, t, rendered, lease, &resp); err != nil {
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
	st.Pushed, st.Head, st.Checks = resp.Pushed, "", nil
	if err := writeShipState(t.statePath(), st); err != nil {
		return nil, err
	}
	return resp, nil
}

// shipFinishRebase renders the rebased head, commits what that changed,
// and pushes it to the change set's branch on origin with lease as the
// expected head ("" for none): never forced past a head it has not seen,
// and never to main.
func shipFinishRebase(ctx *Ctx, t shipTarget, rendered, lease string, resp *ShipRebaseResponse) error {
	var err error
	if resp.Rendered, err = shipRender(ctx, t); err != nil {
		return err
	}
	if resp.Commit, err = shipCommitRendered(ctx, t, rendered, shipRebaseMessage); err != nil {
		return err
	}
	if _, err := ctx.Git("push", "--force-with-lease=refs/heads/"+t.branch+":"+lease, "origin", "HEAD:refs/heads/"+t.branch); err != nil {
		return Errorf("%v", err)
	}
	return nil
}

// shipRebaseOnto rebases HEAD onto onto. At each stop whose conflicts are
// all under rendered, it removes those files, renders again, commits the
// stop as the commit being replayed and continues. At a stop with any
// other conflict it aborts the rebase and returns every conflicted path
// of that stop. A rebase that fails without a conflict is aborted and is
// an error.
func shipRebaseOnto(ctx *Ctx, t shipTarget, onto, rendered string, resp *ShipRebaseResponse) ([]string, error) {
	_, err := ctx.Git("rebase", onto)
	for range shipRebaseLimit {
		if err == nil {
			return nil, nil
		}
		out, derr := ctx.Git("diff", "--name-only", "--diff-filter=U")
		conflicts := gitLines(out)
		if derr != nil || len(conflicts) == 0 {
			shipRebaseAbort(ctx)
			return nil, Errorf("ship rebase: %v", err)
		}
		inRendered, other := splitConflicts(rendered, conflicts)
		if len(other) > 0 {
			if aerr := shipRebaseAbort(ctx); aerr != nil {
				return nil, aerr
			}
			return conflicts, nil
		}
		for _, p := range inRendered {
			if _, err := ctx.Git("rm", "-q", "-f", "--ignore-unmatch", "--", topPath(p)); err != nil {
				shipRebaseAbort(ctx)
				return nil, Errorf("ship rebase: %v", err)
			}
			if !slices.Contains(resp.Settled, p) {
				resp.Settled = append(resp.Settled, p)
			}
		}
		if _, err := shipRender(ctx, t); err != nil {
			shipRebaseAbort(ctx)
			return nil, err
		}
		if _, err := ctx.Git("add", "-A", "--", topPath(rendered)); err != nil {
			shipRebaseAbort(ctx)
			return nil, Errorf("ship rebase: %v", err)
		}
		staged, serr := ctx.Git("diff", "--cached", "--name-only")
		if serr != nil {
			shipRebaseAbort(ctx)
			return nil, Errorf("ship rebase: %v", serr)
		}
		// Commit the stop with the replayed commit's message and author,
		// so --continue opens no editor. A stop that settles to no change
		// commits nothing, and --continue drops it.
		if strings.TrimSpace(staged) != "" {
			if _, err := ctx.Git("commit", "-q", "-C", "REBASE_HEAD"); err != nil {
				shipRebaseAbort(ctx)
				return nil, Errorf("ship rebase: %v", err)
			}
		}
		_, err = ctx.Git("rebase", "--continue")
	}
	shipRebaseAbort(ctx)
	return nil, Errorf("ship rebase: the rebase stopped more than %d times", shipRebaseLimit)
}

// shipMergeOnto merges onto into HEAD, keeping the branch's history. A
// conflict only under rendered is settled by removing those files and
// rendering again; any other conflict aborts the merge and returns every
// conflicted path. A merge that fails without a conflict is aborted and
// is an error.
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
		abort()
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

// shipRebaseAbort aborts a rebase in progress, leaving the branch as it
// was before it.
func shipRebaseAbort(ctx *Ctx) error {
	if _, err := ctx.Git("rebase", "--abort"); err != nil {
		return Errorf("ship rebase: could not abort the rebase: %v", err)
	}
	return nil
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

// shipQuestionStatement is the statement of the question a rebase
// conflict raises.
func shipQuestionStatement(t shipTarget, conflicts []string) string {
	return fmt.Sprintf("Rebasing branch %s onto origin/%s conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: %s?",
		t.branch, shipMain, strings.Join(conflicts, ", "))
}

// shipRaiseQuestion writes a question naming resp.Conflicts: into the
// change set while it is in flight, else (it has landed) into a record
// change set, and commits it on the branch (not pushed). A live question
// already holding the handle with the same statement means an earlier
// call raised it, and nothing is written.
func shipRaiseQuestion(ctx *Ctx, t shipTarget, resp *ShipRebaseResponse) error {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return err
	}
	c, ok := w.changeSet(t.cs)
	if !ok {
		return NotFoundf("no change set %s", t.cs).WithNext(changeLsNext)
	}
	record := false
	switch c.State {
	case changeset.StateAccepted:
	case changeset.StateLanded:
		record = true
	default:
		return Conflictf("cannot ship change set %s: it is %s, not accepted", t.cs, c.State).WithNext(legalNext(c))
	}
	statement := shipQuestionStatement(t, resp.Conflicts)
	m := w.model(t.cs)
	// The stored statement holds resolved references, so compare it with the
	// input resolved the same way.
	stored := statement
	if r, err := resolveFieldTexts(writeBase{m: m}, op.KindQuestion, []string{"statement=" + statement}, nil, true); err == nil {
		stored = strings.TrimPrefix(r[0], "statement=")
	}
	resp.Question = handleSlug("rebase-conflict-" + t.slug)
	link := LinkRequest{From: resp.Question, Type: string(op.LinkBlocks), To: t.cs.String()}
	var items []applyItem
	addLink := applyItem{op: "link", build: func(b writeBase) (edit, error) { return buildLink(b, link) }}
	if n, taken := m.HandleTaken(resp.Question); taken {
		if node, ok := m.Node(n); ok && node.Kind == op.KindQuestion && fieldText(node.Fields["statement"]) == stored {
			// An earlier call raised it; it gains the blocks link when it
			// lacks one, else nothing is written.
			if _, err := m.LinkBetween(n, t.cs, op.LinkBlocks); err == nil {
				return nil
			} else if AsError(err).Code != CodeNotFound {
				return err
			}
			items = []applyItem{addLink}
		} else {
			resp.Question = uniqueHandle(resp.Question, func(h string) bool { _, ok := m.HandleTaken(h); return ok })
			link.From = resp.Question
		}
	}
	if items == nil {
		items = []applyItem{{op: "add", build: func(b writeBase) (edit, error) {
			return buildAdd(b, AddRequest{Kind: string(op.KindQuestion), Handle: resp.Question, Fields: []string{"statement=" + statement}, lenient: true})
		}}, addLink}
	}
	wctx := *ctx
	wctx.ChangeSet, wctx.ChangeRef, wctx.Record = t.cs, "", record
	intent := fmt.Sprintf("Raise the rebase conflicts of change set %s", t.cs)
	r, err := runWriteOn(&wctx, "apply", intent, func(w *workspace, cs id.ID) (edit, error) {
		return buildApply(w, cs, ctx.IDs, items)
	})
	if err != nil {
		return err
	}
	resp.QuestionChange = r.(WriteResponse).ChangeSet.String()
	if _, err := ctx.Git("add", "--", store.Dir); err != nil {
		return Errorf("%v", err)
	}
	_, err = shipCommitStaged(ctx, t, store.Dir, "Raise a question on the rebase conflicts")
	return err
}

// Human says whether the branch was behind, and what the rebase did.
func (r ShipRebaseResponse) Human(w io.Writer) {
	switch {
	case !r.Behind:
		fmt.Fprintf(w, "%s already holds origin/%s at %s; nothing to rebase.\n", r.Branch, shipMain, shortSHA(r.Onto))
		return
	case len(r.Conflicts) > 0:
		fmt.Fprintf(w, "Rebasing %s onto origin/%s at %s conflicts in:\n", r.Branch, shipMain, shortSHA(r.Onto))
		for _, c := range r.Conflicts {
			fmt.Fprintf(w, "- %s\n", c)
		}
		fmt.Fprintln(w, "The rebase was aborted; the branch is as it was.")
		if r.Question != "" {
			fmt.Fprintf(w, "Raised question %s.\n", r.Question)
		}
		return
	}
	fmt.Fprintf(w, "Rebased %s onto origin/%s at %s.\n", r.Branch, shipMain, shortSHA(r.Onto))
	if len(r.Settled) > 0 {
		fmt.Fprintf(w, "Settled %d rendered %s by rendering again.\n", len(r.Settled), shipPlural(len(r.Settled), "file", "files"))
	}
	if r.Commit != "" {
		fmt.Fprintf(w, "Committed %s.\n", shortSHA(r.Commit))
	}
	fmt.Fprintf(w, "Pushed %s.\n", shortSHA(r.Pushed))
}

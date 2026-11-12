package api

// The daemon merges the open model-only pull requests (branches model/*): a
// green one is squash-merged and its branch deleted, then the main checkout
// pulls, with no run close-out involved. A red or conflicting one is not
// merged; it is recorded where the inbox reads it.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// daemonMergeModelPRs merges every green open model PR and pulls main. It
// never fails the poll: a failure is in the daemon log and the next tick
// retries.
func daemonMergeModelPRs(ctx *Ctx) {
	if !modelBranchEnabled(ctx) {
		return
	}
	if err := mergeModelPRs(ctx); err != nil && ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon: model PRs: %v\n", err)
	}
}

func mergeModelPRs(ctx *Ctx) error {
	top, common, ok := gitDirsOf(ctx.Root)
	if !ok || !samePath(top+"/.git", common) {
		return nil // not a git repo, or a linked worktree
	}
	git := func(args ...string) (string, error) {
		out, err := ctx.Exec("git", args, top, nil)
		return strings.TrimSpace(string(out)), err
	}
	if _, err := git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/main"); err != nil {
		return nil
	}
	out, err := git("ls-remote", "--heads", "origin", "refs/heads/"+modelBranchPrefix+"*")
	if err != nil {
		return err
	}
	var stuck []daemon.ModelPR
	merged := ""
	var firstErr error
	for _, line := range strings.Split(out, "\n") {
		_, ref, ok := strings.Cut(line, "\t")
		branch := strings.TrimPrefix(ref, "refs/heads/")
		if !ok || branch == ref {
			continue
		}
		cs, _, _ := strings.Cut(strings.TrimPrefix(branch, modelBranchPrefix), "-")
		v, err := gh(ctx, top, "pr", "view", branch, "--json", "number,url,state,headRefOid,body,mergeable,isDraft,statusCheckRollup")
		if err != nil {
			continue // no pull request yet
		}
		pv, err := parsePRView(v)
		if err != nil || !strings.EqualFold(pv.State, "OPEN") || pv.Draft {
			continue
		}
		settled, green, failed := checksVerdict(pv.Checks)
		switch {
		case strings.EqualFold(pv.Mergeable, "CONFLICTING"):
			stuck = append(stuck, daemon.ModelPR{Branch: branch, Change: cs, Since: ctx.Now(), Detail: fmt.Sprintf("model PR %s conflicts with main", pv.URL)})
		case settled && !green:
			stuck = append(stuck, daemon.ModelPR{Branch: branch, Change: cs, Since: ctx.Now(), Detail: fmt.Sprintf("model PR %s has failing checks: %s", pv.URL, strings.Join(failed, ", "))})
		case settled:
			if _, err := gh(ctx, top, "pr", "merge", branch, "--squash", "--match-head-commit", pv.Head); err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("merge %s: %v", branch, err)
				}
				continue
			}
			_, _ = git("push", "origin", "--delete", branch)
			merged = cs
		}
	}
	if err := daemon.WriteModelPRs(common, stuck); err != nil && firstErr == nil {
		firstErr = err
	}
	if merged != "" {
		if err := pullMainAfterModelMerge(ctx, top, common, merged); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// pullMainAfterModelMerge fast-forwards the main checkout after a model PR
// merged, dropping its untracked copies of the merged entries first. A pull
// that cannot be done is recorded as main-stale, as a ship close-out's is.
func pullMainAfterModelMerge(ctx *Ctx, top, common, cs string) error {
	detail := shipMainBlocked(ctx, top)
	if detail == "" {
		dropIncomingLogCopies(ctx, top)
		if _, err := gitAt(ctx, top, "pull", "--ff-only", "origin", "main"); err != nil {
			detail = err.Error()
		}
	}
	if detail != "" {
		detail = "main was not updated: " + detail
		if err := daemon.WriteMainStale(common, daemon.MainStale{Since: ctx.Now(), Change: cs, Detail: detail}); err != nil {
			return err
		}
		return fmt.Errorf("%s", detail)
	}
	return daemon.ClearMainStale(common)
}

// inboxModelPRs is one item per stuck model PR the daemon recorded.
func inboxModelPRs(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	var items []InboxItem
	for _, p := range daemon.ReadModelPRs(common) {
		at := p.Since.UTC().Format(time.RFC3339)
		sum := sha256.Sum256([]byte(strings.Join([]string{InboxModelPR, p.Branch}, "\x00")))
		items = append(items, InboxItem{
			ID: hex.EncodeToString(sum[:])[:7], Kind: InboxModelPR, At: at, Change: p.Change, Detail: p.Detail,
		})
	}
	return items
}

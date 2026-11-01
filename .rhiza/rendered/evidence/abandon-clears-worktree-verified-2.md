# abandon-clears-worktree-verified-2

evidence · `01M6B9CVRNQ3SCN20A9NDB09BX` · actor [claude-code-e4a90f2d](../actor/claude-code-e4a90f2d.md)

## result

PASS

## statement

At de1efb4f ran go test ./cmd/rhiza -run TestScripts/run-abandon -count=1: run-abandon, run-abandon-dirty, run-abandon-shared and run-abandon-failed pass. run-abandon.txtar asserts abandon prints remove-worktree done and delete-branch done, git worktree list drops to 1 and git branch --list cs/* is empty; run-abandon-dirty.txtar asserts an uncommitted file keeps worktree and branch with 'remove-worktree kept: the worktree is dirty'; run-abandon-shared.txtar asserts abandoning a failed run keeps the worktree and branch a newer live run uses, then abandoning that run removes both (! exists on the folder). Read the diff: abandonCleanup now consults worktreeUser (LiveOf plus later runs in the same worktree) before git worktree remove / branch -D, fixing the review blocker; remote branch and PR untouched. go test ./..., go vet and gofmt clean; make check's only failure is self-check's stale .rhiza/rendered for actor records written during this run.

## Links

### Out

- verifies [abandon-clears-worktree](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-clears-worktree.md)

## History

- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch

# abandon-clears-worktree-verified-3

evidence · `01M6B9ZNJ9DZVVPGXC4N0K906G` · actor [claude-code-16e15d72](../actor/claude-code-16e15d72.md)

## result

PASS

## statement

At 70f0876b, go test ./cmd/rhiza -count=1 -run TestScripts/run-abandon passes run-abandon, run-abandon-dirty, run-abandon-unmerged, run-abandon-shared and run-abandon-failed. run-abandon.txtar asserts 'remove-worktree done' and 'delete-branch done', that git worktree list drops to the main checkout only (so the run's worktree folder is removed) and that git branch --list cs/* is empty. run-abandon-dirty.txtar, with an uncommitted scratch.txt in the worktree, asserts 'remove-worktree kept: the worktree is dirty', 'delete-branch kept', that the file still exists and that the cs/ branch remains. Reading the diff: branch deletion now uses git branch -d, so a branch with unmerged commits is kept (run-abandon-unmerged asserts that, and that the folder is gone via ! exists). go fmt, vet and the full go test pass in make check. Its self-check step fails only on rendered actor files and index.md made stale by this run's own model writes, my own actor among them, not by the commit.

## Links

### Out

- verifies [abandon-clears-worktree](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-clears-worktree.md)

## History

- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch

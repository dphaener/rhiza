# abandon-clears-worktree

capability · `01M6ATC08K455FW14QCXPZ9QHD` · parent [run-abandon](../run-abandon.md)

## acceptance

run-abandon.txtar: after rhiza run abandon the worktree folder is gone, git branch --list cs/<slug> is empty and the output reports remove-worktree done; with an uncommitted file in the worktree both are kept and the output says the worktree is dirty

## statement

rhiza run abandon removes the run's clean worktree and deletes its local branch, reporting each as done, already gone, or kept because the worktree is dirty; the remote branch and PR are untouched

## Links

### Out

- decided-by [abandon-removes-clean-worktree](../../../../../decision/abandon-removes-clean-worktree.md)

### In

- verifies [abandon-clears-worktree-verified](../../../../../evidence/abandon-clears-worktree-verified.md)
- verifies [abandon-clears-worktree-qa-1](../../../../../evidence/abandon-clears-worktree-qa-1.md)
- verifies [abandon-clears-worktree-verified-2](../../../../../evidence/abandon-clears-worktree-verified-2.md)
- verifies [abandon-clears-worktree-verified-3](../../../../../evidence/abandon-clears-worktree-verified-3.md)
- verifies [pr-164-abandon-clears-worktree](../../../../../evidence/pr-164-abandon-clears-worktree.md)

## History

- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch

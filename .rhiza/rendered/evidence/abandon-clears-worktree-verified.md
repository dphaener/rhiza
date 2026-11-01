# abandon-clears-worktree-verified

evidence · `01M6B8A2WD9FP22DQYF017MYEB` · actor [claude-code-e296c7cb](../actor/claude-code-e296c7cb.md)

## result

PASS

## statement

At 70424b00 ran go test ./cmd/rhiza -run 'TestScripts/(run-abandon|run-abandon-dirty|describe-golden)$' -count=1: all PASS. run-abandon.txtar asserts 'remove-worktree done' and 'delete-branch done' in the abandon output, git worktree list drops to 1 entry (the run worktree removed by git worktree remove) and git branch --list cs/* is empty. run-abandon-dirty.txtar writes an uncommitted scratch.txt into the run worktree, then asserts 'remove-worktree kept: the worktree is dirty', 'delete-branch kept', the file still exists and the cs/* branch remains. Read the diff of internal/api/run.go: abandonCleanup removes the worktree without force (force only when the sole leftovers are byte-identical copies of .rhiza/log entries), deletes the local cs/<slug> branch only after the worktree is gone, reports done/already-gone/kept, and runs no push or gh call, so the remote branch and PR are untouched. go test ./... passes; make check's fmt/vet/test pass, self-check flags only stale .rhiza/rendered files from this run's actor writes.

## Links

### Out

- verifies [abandon-clears-worktree](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-clears-worktree.md)

## History

- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch

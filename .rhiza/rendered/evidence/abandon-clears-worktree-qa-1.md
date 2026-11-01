# abandon-clears-worktree-qa-1

evidence · `01M6B8EXZ7N2T6WVJ3DY01YBAP` · actor [claude-code-de288672](../actor/claude-code-de288672.md)

## result

PASS

## statement

QA cycle 1: built go build -o $TMPDIR/rhiza-qa ./cmd/rhiza at 5a0d3973 and drove it by hand in scratch git repos (RHIZA_ROOT/CHANGE/RUN unset) with a hold workflow (wait step) started via run start --no-wait. Clean worktree: run abandon --dry-run left the worktree; run abandon printed '## Clean-up' with 'remove-worktree done' and 'delete-branch done', the <repo>.<slug> folder was gone, git worktree list showed only main, git branch --list cs/* was empty, and the cs/<slug> branch pushed beforehand to a bare origin was still listed by git ls-remote (remote untouched); the leftover untracked .rhiza/log copy did not count as dirty. Untracked scratch.txt in the worktree: 'remove-worktree kept: the worktree is dirty', 'delete-branch kept: the worktree is kept', file and branch remained, and a new run started on the change set. Modified tracked file: both kept, dirty. Worktree rm -rf'd first: remove-worktree already-gone (pruned), delete-branch done. Worktree and branch both removed first: both already-gone. A second abandon exits 2. TestScripts/run-abandon and run-abandon-dirty pass. Seen: a clean worktree whose branch holds an unpushed commit still gets branch -D, so that commit is dropped (reported as a finding).

## Links

### Out

- verifies [abandon-clears-worktree](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-clears-worktree.md)

## History

- `01M6ATATQKHW0YV8W1XDP8YECH` Abandoning a run removes its clean worktree and local branch

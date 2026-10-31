# ship-worktree-follows-pr

capability · `01M68582MQEBTQNMGMAGHV53ZG` · parent [shipping](../shipping.md)

## acceptance

A testscript with a bare remote: after open-pr a second clone pushes a commit to the branch; the rebase step fast-forwards instead of exiting 3, land's push succeeds and the worktree's HEAD equals the pushed head; Go tests with a fake git show a dirty worktree or a local-only commit leaves HEAD unchanged and returns a conflict naming them

## statement

Before committing or pushing, rhiza ship rebase, land and open-pr fast-forward a clean worktree to origin's branch when that branch contains the worktree's HEAD, and refuse, naming the commits, when the worktree has uncommitted changes or local commits origin's branch lacks

## Links

### Out

- decided-by [follow-head-only-when-it-contains-pushed](../../../../decision/follow-head-only-when-it-contains-pushed.md)

### In

- verifies [ship-worktree-follows-pr-verified](../../../../evidence/ship-worktree-follows-pr-verified.md)
- verifies [ship-worktree-follows-pr-qa-1](../../../../evidence/ship-worktree-follows-pr-qa-1.md)
- verifies [pr-126-runs-follow-a-moved-pr-head](../../../../evidence/pr-126-runs-follow-a-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

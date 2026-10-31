# ship-worktree-follows-pr-qa-1

evidence · `01M68E9BKNEKA0HXCJY0B35H0P` · actor [claude-code-8d3acb4b](../actor/claude-code-8d3acb4b.md)

## result

PASS

## statement

At 14d9f17, as QA with a scratch build in a scratch repo with a bare origin: a second clone pushed a commit to cs/ship-it. With feature.txt modified, ship rebase exited 3 with 'origin's cs/ship-it has moved to 30b4e3a, but <repo> has uncommitted changes that stop it following: M feature.txt' and HEAD stayed put. Once the worktree was clean, ship rebase exited 0 and HEAD equalled the moved head. With a local commit origin lacked ('Raise a question on the rebase conflicts') and origin ahead, ship rebase exited 3 naming '429a4e6 Raise a question on the rebase conflicts', and HEAD stayed put. The ship-follows-pr-head testscript passes.

## Links

### Out

- verifies [ship-worktree-follows-pr](../capability/rhiza/orchestration/shipping/ship-worktree-follows-pr.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

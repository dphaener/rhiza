# pr-126-runs-follow-a-moved-pr-head

evidence · `01M68FEFTZ2T4VZAEJ7MP46MF1` · actor [darin-haener](../actor/darin-haener.md)

## head

c7b941e2a45f5a323bb19a005c5be8d6cac181b3

## pr

https://github.com/dphaener/rhiza/pull/126

## result

1 check passed: check SUCCESS

## statement

Pull request #126 (https://github.com/dphaener/rhiza/pull/126) passed CI at c7b941e

## status

verified

## Links

### Out

- verifies [follow-moved-pr-head](../capability/rhiza/orchestration/shipping/wait-green-pr/follow-moved-pr-head.md)
- verifies [rebase-keeps-branch-merges](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/rebase-keeps-branch-merges.md)
- verifies [ship-worktree-follows-pr](../capability/rhiza/orchestration/shipping/ship-worktree-follows-pr.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

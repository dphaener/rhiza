# rebase-keeps-branch-merges-qa-1

evidence · `01M68E9CVVK36EPB70J39GHHTE` · actor [claude-code-8d3acb4b](../actor/claude-code-8d3acb4b.md)

## result

PASS

## statement

At 14d9f17, as QA with a scratch build in a scratch repo: when main conflicted on shared.txt, ship rebase raised rebase-conflict-ship-it and exited 1. A second clone then merged main into cs/ship-it with a resolution (14a1dab) and pushed it, and main moved again with unrelated commits. ship rebase exited 0: it merged origin/main into the branch ('Merge origin/main into cs/ship-it' plus a render commit) instead of rebasing, raised no new question, and pushed eb82547. 14a1dab is an ancestor of the pushed head, and origin equals local. The run-ship-change-rebase testscript (cases 3 and 4) passes.

## Links

### Out

- verifies [rebase-keeps-branch-merges](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/rebase-keeps-branch-merges.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

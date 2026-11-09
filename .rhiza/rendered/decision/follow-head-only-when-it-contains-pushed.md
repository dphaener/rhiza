# follow-head-only-when-it-contains-pushed

decision · `01M68557HN470AVYK4NQC26BRH`

## statement

A run follows a PR head it did not push only when that head contains the head the run pushed; a head that drops the run's commits leaves ship check unsettled, with a detail naming the dropped head, and raises no question; ship check only adopts the head, and the write steps fast-forward a clean worktree to it, refusing one with uncommitted changes or local commits the head lacks; a branch that holds a merge of main is merged with main again rather than rebased

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m6-overnight-interview](../source/m6-overnight-interview.md)

### In

- decided-by [follow-moved-pr-head](../capability/rhiza/orchestration/shipping/wait-green-pr/follow-moved-pr-head.md)
- decided-by rebase-keeps-branch-merges (retired)
- decided-by [ship-worktree-follows-pr](../capability/rhiza/orchestration/shipping/ship-worktree-follows-pr.md)

## History

- `01M684XWDT3EKJD4JPXT0JVK2M` Plan milestone 6: Rhiza runs unattended overnight, with the run machinery fixes that needed the orchestrator's hand in milestones 4 and 5
- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

# follow-head-only-when-it-contains-pushed

decision · `01M68557HN470AVYK4NQC26BRH`

## statement

A run follows a PR head it did not push only when that head contains the head the run pushed; a head that drops the run's commits stops with a question; ship check only adopts the head, and the write steps fast-forward a clean worktree to it, refusing one with uncommitted changes or local commits the head lacks; a branch that holds a merge of main is merged with main again rather than rebased

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m6-overnight-interview](../source/m6-overnight-interview.md)

## History

- `01M684XWDT3EKJD4JPXT0JVK2M` Plan milestone 6: Rhiza runs unattended overnight, with the run machinery fixes that needed the orchestrator's hand in milestones 4 and 5

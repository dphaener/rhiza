# rebase-keeps-branch-merges

capability · `01M68581Q6V7XS1VHP9Q48KQJS` · parent [rebase-when-behind](../rebase-when-behind.md)

## acceptance

run-ship-change-rebase.txtar gains a case: the branch conflicts with main, a merge commit resolving it is pushed to origin, main then moves with an unrelated commit; rebase and rebase-merge pass, no rebase-conflict question is written, the pushed head has the earlier merge as an ancestor and the run merges; when main's new commit conflicts again, the question is still raised and the step exits 1

## statement

When a change set's branch already holds a merge of main, rhiza ship rebase merges origin's main into it instead of rebasing, settling rendered-only conflicts by rendering again, so conflicts that merge resolved do not come back

## Links

### Out

- decided-by [follow-head-only-when-it-contains-pushed](../../../../../decision/follow-head-only-when-it-contains-pushed.md)

### In

- verifies [rebase-keeps-branch-merges-verified](../../../../../evidence/rebase-keeps-branch-merges-verified.md)
- verifies [rebase-keeps-branch-merges-qa-1](../../../../../evidence/rebase-keeps-branch-merges-qa-1.md)
- verifies [pr-126-runs-follow-a-moved-pr-head](../../../../../evidence/pr-126-runs-follow-a-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

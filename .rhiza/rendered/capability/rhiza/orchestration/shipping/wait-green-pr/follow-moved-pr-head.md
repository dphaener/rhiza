# follow-moved-pr-head

capability · `01M68580Q11GDHVQY0FF7A653E` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

Go tests beside TestShipCheckWaitsForPushedHead with a fake git and gh: when merge-base --is-ancestor pushed head holds, --settled passes, a green check records the new head as Head and Pushed, and ship merge pins to it; when it does not hold the check stays unsettled naming the dropped head; a testscript pushes from a second clone to a local bare remote and the run's wait-pr passes without timing out

## statement

When a PR's head is not the head the run pushed but contains it, rhiza ship check records it as the pushed head and judges its CI, so wait-pr and the later wait steps pass on a head someone else pushed on top of the run's; a head that drops the run's commits stays unsettled and its detail names the dropped head

## Links

### Out

- decided-by [follow-head-only-when-it-contains-pushed](../../../../../decision/follow-head-only-when-it-contains-pushed.md)

### In

- verifies [follow-moved-pr-head-verified](../../../../../evidence/follow-moved-pr-head-verified.md)
- verifies [follow-moved-pr-head-qa-1](../../../../../evidence/follow-moved-pr-head-qa-1.md)
- verifies [follow-moved-pr-head-verified-2](../../../../../evidence/follow-moved-pr-head-verified-2.md)
- verifies [pr-126-runs-follow-a-moved-pr-head](../../../../../evidence/pr-126-runs-follow-a-moved-pr-head.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing

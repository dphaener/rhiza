# retry-reopens-at-the-wait

decision · `01M6XXFK3NNFP29EST9F4G2CBN`

## reasoning

The ship-change.toml workflow pairs every CI check as a wait step (`rhiza ship check --settled`, 2h timeout, 30s interval) whose on.pass is a command step (`rhiza ship check`), and Engine.Retry reopens at State.FailedStep, so a retry after a fresh push skips the poll and sees 'not settled'. The fix belongs in how retry picks its step, not in each workflow: a per-step retry_at field would need every built-in and custom workflow to opt in, and folding polling into the green command would duplicate the wait kind's timeout and interval and hide the wait from the daemon's step view. Re-entering the preceding wait costs nothing when CI already settled, and keeps 'pending is not failure' true on every path, including the nebari PR #83 case in the change set's intent. When a failed step has no wait predecessor, retry keeps today's behaviour.

## statement

A run retry whose failed step is entered from a wait step's on.pass reopens the run at that wait step, not at the failed step. The wait polls until the checks settle (it passes at its first poll when they already have), then the green step runs as usual and fails only on a real red. This is a generic engine rule, so ship-change's pr-green, rebased-green, landed-green and rebased-merge-green, plan's and plan-milestone's pr-green, and custom workflows with the same wait-then-check shape all get it without new workflow fields. Green commands stay plain `rhiza ship check` with no polling of their own.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-47f5cdc8-1](../question/question-47f5cdc8-1.md)
- derived-from [landed-green-retry-plan](../source/landed-green-retry-plan.md)

### In

- decided-by [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

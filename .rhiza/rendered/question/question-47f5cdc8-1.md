# question-47f5cdc8-1

question · `01M6XXAH5Y23S06ZT8GJX7QTBY`

## session

47f5cdc8-5adc-491d-a6cf-01a1bbd2cfd6

## statement

Ship-change's landed-green failed on pending CI because `rhiza run resume --retry` restarts a failed run at the step that failed (landed-green, a plain `rhiza ship check` command), skipping wait-landed, which is the step that polls until checks settle. The same happens to pr-green, rebased-green and rebased-merge-green in ship-change and plan's green steps. How should a retry wait for CI?

## Links

### In

- answers [retry-reopens-at-the-wait](../decision/retry-reopens-at-the-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

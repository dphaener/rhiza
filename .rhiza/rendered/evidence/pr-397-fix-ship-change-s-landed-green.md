# pr-397-fix-ship-change-s-landed-green

evidence · `01M6YGT5KXH897Y3RBS84T1ASC` · actor [dphaener](../actor/dphaener.md)

## head

99b58b2ce30979da1fbbea59d55e457590916b19

## pr

https://github.com/dphaener/rhiza/pull/397

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #397 (https://github.com/dphaener/rhiza/pull/397) passed CI at 99b58b2

## status

verified

## Links

### Out

- verifies [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

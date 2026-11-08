# retry-reopens-at-wait

capability · `01M6XXHNR6FZCKNSA3SC47TZZK` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/api (or internal/run) shows a retry of a run failed at a command step entered from a wait step's on.pass restarting at the wait, which polls a not-settled until and then passes on to the command; a run failed at a step with no wait predecessor still retries at that step; a testscript (run-resume-retry or a new run-ship-retry-pending with a fake gh) retries a ship-change run failed at landed-green while the new head's checks are pending, and it waits at wait-landed, then passes landed-green once they pass; run with go test -trimpath ./internal/api/... ./internal/run/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-resume-retry|run-ship)' only

## statement

rhiza run resume --retry on a failed run whose failed step is the on.pass target of exactly one wait step in its workflow reopens the run at that wait step instead of the failed step (the retry is journalled, and the response's retry field, at that wait step), so ship-change's pr-green, rebased-green, landed-green and rebased-merge-green and plan's and plan-milestone's green steps wait for pending CI after a retry and fail only on a real red; a failed step with no such wait, or more than one, is retried at itself as today

## Links

### Out

- decided-by [retry-reopens-at-the-wait](../../../../../decision/retry-reopens-at-the-wait.md)
- derived-from [landed-green-retry-plan](../../../../../source/landed-green-retry-plan.md)

### In

- verifies [retry-reopens-at-wait-verified](../../../../../evidence/retry-reopens-at-wait-verified.md)
- verifies [retry-reopens-at-wait-qa-1](../../../../../evidence/retry-reopens-at-wait-qa-1.md)
- verifies [pr-397-fix-ship-change-s-landed-green](../../../../../evidence/pr-397-fix-ship-change-s-landed-green.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

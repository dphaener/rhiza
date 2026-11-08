# retry-reopens-at-wait-verified

evidence · `01M6YFH0KGVY24M32PYYGY8RAX` · actor [claude-code-cb9803b7](../actor/claude-code-cb9803b7.md)

## result

PASS

## statement

Verifier at fe0e16f4. Diff: run.RetryStep(def, failed) returns the single wait step whose on.pass is the failed step, otherwise (none, or more than one) the failed step itself. Both the engine's retry (engine.go, which journals TypeRetried at that step) and api runRunResume (response Retry plus the dry-run fold) call it. Built-in ship-change (pr-green, rebased-green, landed-green, rebased-merge-green, main-green), plan and plan-milestone (pr-green, rebased-green) each have exactly one wait leading to each green step. Ran in a clean env (env -i, because the run session's env makes scripts' rhiza init hit the enclosing repo): go test -trimpath ./internal/run/... ok. TestRetryReopensAtWait: the retry reopens at wait, the wait fails one poll then holds (3 polls), outcomes are wait:pass green:fail wait:pass green:pass, the run succeeds. TestRetryStepChoice: a step with two wait predecessors and a step with none are retried at themselves. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-resume-retry|run-ship)': all 18 scripts PASS. run-ship-retry-pending: a fake gh fails landed-green on pending checks, the retry dry-run says wait-landed, and the retry holds after several checks, then landed-green passes. run-resume-retry-wait: a failure with no wait before it is retried at solo. go test -trimpath ./internal/api/...: Resume/Retry tests PASS. The only failures were TestReadBudget (read medians just over 50ms with the VM load at 7-10; the diff does not touch the read path) and the known .rhiza-root/usage.live.json leak guard, both unrelated to this change.

## Links

### Out

- verifies [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

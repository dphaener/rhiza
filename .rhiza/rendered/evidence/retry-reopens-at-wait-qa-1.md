# retry-reopens-at-wait-qa-1

evidence · `01M6YFT39PH4BY68AGXCYTX32H` · actor [claude-code-b368ab1b](../actor/claude-code-b368ab1b.md)

## result

PASS

## statement

QA at 9b783f07, hands-on with the built binary (go build -trimpath) in a scratch Rhiza root and repo, RHIZA_* unset. Workflow qa: wait-ci (until test -f ci-done, 200ms interval) then green (test -f ci-green). With ci-done present and ci-green absent, run start --foreground failed at green. With ci-done removed (CI pending again), run resume --retry --dry-run printed 'would be retried at step wait-ci' and --json gave retry: wait-ci; the real retry, with a background job creating both files after 1.5s, showed 'wait-ci (visit 2, retried): pass, then green: held after 8 checks', then green passed and the run succeeded, with the retry journalled 'retried at wait-ci'. Real red: a fresh run with ci-done present and ci-green absent, retried, went back to wait-ci, which passed on its first check, and green failed again (the run failed). Workflow qa2: solo (a command no wait leads to) failed, and the dry-run and the retry both chose solo; twin (the on.pass target of two waits, w1 and w2) failed, and the retry chose twin itself, not a wait. The built-in ship-change, plan and plan-milestone each have exactly one wait before every green step (pr, rebased, landed, rebased-merge, main). Scoped tests: go test -trimpath ./internal/run/... -run Retry ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-resume-retry|run-ship-retry-pending)' ok in a clean env (env -i). Under the review session's env, rhiza init in the scripts hits the enclosing repo, which the change did not cause.

## Links

### Out

- verifies [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.

# rex-tab-launch-retry

capability · `01M76XZG799KFE9QRBJNYYZ7KG` · parent [rex-tab-launch](../rex-tab-launch.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run SessionRunRexRetry) with a fake rex and an injected clock assert: two failed tab opens then a success start the session with one tab; four failures fail naming Rex's last error; no test starts a real Rex or agent.

## statement

Opening a run session's Rex tab retries a failed rex call (Rex not answering, or no window ID returned) up to three times with growing waits before failing, so a Rex hiccup no longer holds a run.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [rex-tab-launch-retry-verified](../../../../../evidence/rex-tab-launch-retry-verified.md)
- verifies [rex-tab-launch-retry-qa-1](../../../../../evidence/rex-tab-launch-retry-qa-1.md)
- verifies [rex-tab-launch-retry-qa-2](../../../../../evidence/rex-tab-launch-retry-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

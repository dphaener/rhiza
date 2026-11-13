# ship-gh-transient-retry

capability · `01M76XZG3ZRYA9NTD3SF6QWY0K` · parent [shipping](../shipping.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run ShipGHRetry) with a fake Exec and an injected clock (no real sleep) assert: open-pr whose gh pr create fails twice with 'HTTP 502' then succeeds passes after three calls; a 'GraphQL: Something went wrong' error is retried the same way; four 5xx failures fail the step naming the last error; an HTTP 422 or 'no pull requests found' is not retried.

## statement

Ship steps retry a gh call that fails with a GitHub server error (HTTP 5xx), a timeout or a GraphQL 'Something went wrong' error up to three times with growing waits before failing the step, and fail at once on any other error, so a GitHub blip no longer holds a run.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [ship-gh-transient-retry-verified](../../../../evidence/ship-gh-transient-retry-verified.md)
- verifies [ship-gh-transient-retry-qa-1](../../../../evidence/ship-gh-transient-retry-qa-1.md)
- verifies [ship-gh-transient-retry-qa-2](../../../../evidence/ship-gh-transient-retry-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

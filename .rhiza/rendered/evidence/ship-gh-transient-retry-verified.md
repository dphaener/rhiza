# ship-gh-transient-retry-verified

evidence · `01M791EATH5Q10BWJYWSYJ9SHJ` · actor [claude-code-efa8ac6f](../actor/claude-code-efa8ac6f.md)

## result

PASS

## statement

At 41f90bd44, ran go test -trimpath -count=1 ./internal/api -run ShipGHRetry -v: TestShipGHRetryRecovers (gh pr create failing twice with 'HTTP 502', and twice with 'GraphQL: Something went wrong', passes after 3 calls with two growing injected waits via ctx.Sleep, no real sleep), TestShipGHRetryGivesUp (repeated HTTP 503 fails open-pr after 4 calls/3 waits with the message naming HTTP 503), TestShipGHRetrySkipsOtherErrors (HTTP 422 and 'no pull requests found' fail after 1 call, no wait) all PASS. go test -run Ship and go vet ./internal/api pass. Diff read: every gh call goes through gh() in internal/api/ship.go, which wraps ctx.Exec in retryGH (waits 2s/6s/18s; transient = HTTP 5xx, timeout/timed out, 'something went wrong').

## Links

### Out

- verifies [ship-gh-transient-retry](../capability/rhiza/orchestration/shipping/ship-gh-transient-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

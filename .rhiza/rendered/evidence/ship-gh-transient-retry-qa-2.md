# ship-gh-transient-retry-qa-2

evidence · `01M7963H8HBK4B64P1R3XS4K54` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: gh pr create failing twice with 'HTTP 502: Bad Gateway' then succeeding: open-pr passed after 3 pr create calls; failing 4 times: open-pr failed after 4 calls naming the 502 '(after 3 retries)'; 'HTTP 422: Validation Failed': open-pr failed after 1 call (not retried). go test -trimpath ./internal/api -run ShipGHRetry passes.

## Links

### Out

- verifies [ship-gh-transient-retry](../capability/rhiza/orchestration/shipping/ship-gh-transient-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

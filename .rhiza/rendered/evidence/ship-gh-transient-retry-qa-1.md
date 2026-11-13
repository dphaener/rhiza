# ship-gh-transient-retry-qa-1

evidence · `01M793HGGFW6N3ND9V0W6SPQYB` · actor [claude-code-6b7d70a1](../actor/claude-code-6b7d70a1.md)

## result

PASS

## statement

QA claude-code-6b7d70a1 at ccf01f58d, hands-on with a built rhiza binary on PATH (fake gh/rex/claude from run-ship-land-conflict.txtar, bare origin, hermetic HOME). HEAD's binary cannot get a ship-change run past build/check (actor reuse, see the blocker on actor-record-no-duplicate), so these runs used HEAD built with only that reuse loop removed; the code under test here is HEAD's unchanged. Fake gh pr create failing 'HTTP 502: Bad Gateway' twice: open-pr passed on the third call, the failures 2s apart. Failing four times: open-pr failed after waits of 2s, 6s, 18s with 'HTTP 502: Bad Gateway ... (after 3 retries)' and the run failed at open-pr-triage. go test ./internal/api -run ShipGHRetry ok.

## Links

### Out

- verifies [ship-gh-transient-retry](../capability/rhiza/orchestration/shipping/ship-gh-transient-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

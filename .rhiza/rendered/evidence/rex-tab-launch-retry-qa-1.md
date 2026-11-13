# rex-tab-launch-retry-qa-1

evidence · `01M793HGMD6CRXQE5BCKR032F9` · actor [claude-code-6b7d70a1](../actor/claude-code-6b7d70a1.md)

## result

PASS

## statement

QA claude-code-6b7d70a1 at ccf01f58d, hands-on with a built rhiza binary on PATH (fake gh/rex/claude from run-ship-land-conflict.txtar, bare origin, hermetic HOME). HEAD's binary cannot get a ship-change run past build/check (actor reuse, see the blocker on actor-record-no-duplicate), so these runs used HEAD built with only that reuse loop removed; the code under test here is HEAD's unchanged. Fake rex failing 'rex: server not answering' on its first two calls: the builder's tab opened on the third call (waits 1s, 3s) and the run reached sign-off. Rex always failing: four calls at 0s, +1s, +3s, +9s, then the run failed 'rex could not open the session's tab: ... server not answering (after 3 retries)'. go test ./internal/api -run SessionRun ok.

## Links

### Out

- verifies [rex-tab-launch-retry](../capability/rhiza/runtime/rex-sessions/rex-tab-launch/rex-tab-launch-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

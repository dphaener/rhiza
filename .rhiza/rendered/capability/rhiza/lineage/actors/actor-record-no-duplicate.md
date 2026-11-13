# actor-record-no-duplicate

capability · `01M76XZG6D3HR7G5Y6MHENHMB5` · parent [actors](../actors.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run 'RecordAgentActor') assert: recording the actor of a session whose actor an in-flight change set (another run's branch merged in) already recorded with the same session value writes no new node and returns that actor; a landed or in-flight actor holding claude-code-0123abcd under another session value is not returned, and the new actor gets a unique handle, so no two actors hold one handle; rhiza check on the result reports no duplicate handle.

## statement

Recording an agent actor for a session reuses the actor that a landed or in-flight change set already recorded for that session (its session field equals the session ID) instead of writing a second node, so review's code check no longer fails on a duplicate actor handle. Session IDs are UUIDv7, so sessions launched close together share the base handle <agent>-<first 8 hex>: an actor of another session holding that handle, landed or in flight, is never reused, and the new actor gets a unique handle of its own.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [actor-record-no-duplicate-verified](../../../../evidence/actor-record-no-duplicate-verified.md)
- verifies [actor-record-no-duplicate-verified-2](../../../../evidence/actor-record-no-duplicate-verified-2.md)
- verifies [actor-record-no-duplicate-qa-2](../../../../evidence/actor-record-no-duplicate-qa-2.md)
- verifies [pr-515-held-run-error-sources](../../../../evidence/pr-515-held-run-error-sources.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

# actor-record-no-duplicate-verified-2

evidence · `01M794H0TMZVSP6T72RPJY4A01` · actor [claude-code-a07aa764](../actor/claude-code-a07aa764.md)

## result

PASS

## statement

Verifier claude-code-a07aa764 at 8005f070d (attempt 2). Ran go test -trimpath -count=1 ./internal/api -run 'RecordAgentActor' -v: PASS for TestRecordAgentActorReusesSession, TestRecordAgentActorReusesInFlightSession (an in-flight change set's actor with the same session value is returned and no change set is written), TestRecordAgentActorLandedHandleOtherSession and TestRecordAgentActorInFlightHandleOtherSession (a landed or in-flight actor holding claude-code-0123abcd under another session value, sharing the UUIDv7 prefix, is not returned, the new actor gets another handle, and exactly one actor holds claude-code-0123abcd). The full ./internal/api package passes. Read the diff: recordAgentActor reuses an actor only when its session field equals the session ID (live, or in a non-terminal change set), and uniqueHandle now also avoids handles held in flight. The attempt-1 blocker is fixed: go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(run-ship-land-conflict|run-ship-change|run-ship-conflict-fix)$' passes. rhiza check, run from a clean detached worktree with this commit's build, reports no duplicate handle. Its only failures are render files for the actors of sessions live now, which is expected in a run.

## Links

### Out

- verifies [actor-record-no-duplicate](../capability/rhiza/lineage/actors/actor-record-no-duplicate.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

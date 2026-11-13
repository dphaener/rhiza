# actor-record-no-duplicate-verified

evidence · `01M7916C5NBJYB1RG4CBZ8VREJ` · actor [claude-code-5db19edc](../actor/claude-code-5db19edc.md)

## result

PASS

## statement

Verifier claude-code-5db19edc at d38f497d2: ran go test -trimpath -count=1 ./internal/api -run 'ActorRecord|SessionRun|RecordAgentActor': all PASS, incl. new TestRecordAgentActorReusesLandedHandle (a landed actor already holding claude-code-0123abcd under another session value is returned with its change set, no new change set written, exactly one actor holds the handle) and TestRecordAgentActorReusesInFlightHandle (same for an open, unlanded change set holding it). Diff read: recordAgentActor in internal/api/session.go now scans every change set's actor node.create ops for the base handle <agent>-<sid[:8]> and reuses it before falling back to uniqueHandle. Built the binary and ran rhiza check on the tree: no duplicate handle reported; the only fails are stale render files for uncommitted run actor records (run-env, unrelated).

## Links

### Out

- verifies [actor-record-no-duplicate](../capability/rhiza/lineage/actors/actor-record-no-duplicate.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

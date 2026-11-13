# actor-record-no-duplicate-qa-2

evidence · `01M7963H041AE8AQKGWZP6FQ6B` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: with RHIZA_TEST_NOW/SEED set, the builder and verifier sessions (a6fe1895-5d81-... and a6fe1895-5d86-...) share the 8-hex prefix that broke cycle 1; they got distinct actors claude-code-a6fe1895 and claude-code-a6fe1895-N, build/check passed, the run succeeded, and rhiza check on the result printed ok with no duplicate handle. Also go test -trimpath ./internal/api -run RecordAgentActor passes, and the 19 run-* testscripts cycle 1 found red (run-ship-change, run-ship-model-conflict, run-ship-conflict-fix, run-ship-land-conflict, ...) all pass.

## Links

### Out

- verifies [actor-record-no-duplicate](../capability/rhiza/lineage/actors/actor-record-no-duplicate.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

# rex-tab-launch-retry-qa-2

evidence · `01M7963H4F5KKDBCC4RKSW49E5` · actor [claude-code-dc3d533e](../actor/claude-code-dc3d533e.md)

## result

PASS

## statement

QA claude-code-dc3d533e at 99e340d22 (cycle 2), hands-on with the built binary driving real ship-change runs (rhiza run start/approve --foreground) in a scratch Rhiza root with a bare origin and fake gh, rex and claude on PATH: with fake rex failing its first 2 tab opens ('rex: server not answering'), the builder's tab opened on the third try and the run went on to ship; with 4 failures the run failed at build/build naming 'rex: server not answering (after 3 retries)'. go test -trimpath ./internal/api -run SessionRunRexRetry passes.

## Links

### Out

- verifies [rex-tab-launch-retry](../capability/rhiza/runtime/rex-sessions/rex-tab-launch/rex-tab-launch-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

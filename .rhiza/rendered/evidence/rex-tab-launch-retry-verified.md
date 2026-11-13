# rex-tab-launch-retry-verified

evidence · `01M79195A63V2NYKW1QYVABXR5` · actor [claude-code-d8b4b327](../actor/claude-code-d8b4b327.md)

## result

PASS

## statement

At f95af36fd ran `go test -trimpath -count=1 ./internal/api -run SessionRunRexRetry`: TestSessionRunRexRetryRecovers (fake rex: a failed call, then output with no window_id, then a window) passes with 3 rex calls, a single successful tab open and two growing injected waits; TestSessionRunRexRetryGivesUp (four failures) passes with exit 1, 4 calls, 3 waits and a message holding Rex's last error 'rex is down'. Both use a scripted fake Exec and an injected Ctx.Sleep, so no real Rex, agent or wait runs. Read the diff: retryRex wraps only the new-window rex call in openSessionPane with waits of 1s/3s/9s, and parseRexWindow already rejects an empty window_id. Also ran -run 'Session|Pane' (ok), gofmt -l (clean) and go vet ./internal/api (ok).

## Links

### Out

- verifies [rex-tab-launch-retry](../capability/rhiza/runtime/rex-sessions/rex-tab-launch/rex-tab-launch-retry.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).

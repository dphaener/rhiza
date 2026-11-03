# daemon-self-update-verified

evidence · `01M6FMFSAWEZPYDMB4TV9WG2FT` · actor [claude-code-014090f9](../actor/claude-code-014090f9.md)

## result

PASS

## statement

At c0f1981c I ran go test -count=1 ./internal/daemon ./internal/api -run 'SelfUpdate|DaemonLoop' -v. Every test passed. TestSelfUpdateMovedMainRunsOnce: a moved origin/main is fetched and fast-forwarded, and make update runs once, again only after main moves. TestSelfUpdateUnmovedRunsNothing: equal and ahead run nothing. TestSelfUpdateSkipsAndWaitsForMove: dirty, off main, diverged and a failing make update are logged and not retried until origin/main moves. TestDaemonLoopSelfUpdateWaitsForPending: pending starts at polls 1 and 2 delay the self-update to polls 3 and 4. TestDaemonSelfUpdateSetting: runDaemon with injected git and make, with rhiza.daemon.selfUpdate unset or false, runs no git or make in the main checkout; set to true, it fetches, runs merge --ff-only and runs make update exactly once over 3 polls. A scratch build's rhiza daemon settings in a fresh repo printed 'rhiza.daemon.selfUpdate = false (default)'. make fmt-check test: all 34 packages ok, no FAIL. I read the diff for internal/daemon/selfupdate.go, the daemonLoop gate selfUpdate != nil && pending() == 0, and the settings def; they match the statement. make check's only failure was a stale render for this verifier's own actor record.

## Links

### Out

- verifies [daemon-self-update](../capability/rhiza/orchestration/daemon/daemon-self-update.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs

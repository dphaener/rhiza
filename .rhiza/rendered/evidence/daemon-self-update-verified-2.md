# daemon-self-update-verified-2

evidence · `01M6FPAHEGYZGGZ55MG1HA1488` · actor [claude-code-e81297a5](../actor/claude-code-e81297a5.md)

## result

PASS

## statement

Verified at 08293d6a (attempt 2). Read the diff: SelfUpdate now fetches at most once per Every (5m via SelfUpdateEvery and ctx.Now), and the daemon's fetch runs through ctx.Start with GIT_TERMINAL_PROMPT=0 and is killed after 30s (fetchMainWithin), with Kill closing the output pipe after the grace period so a lingering child cannot block. Ran go test -count=1 -run 'SelfUpdate|FetchMainWithin' on ./internal/daemon and ./internal/api: TestSelfUpdateMovedMainRunsOnce, UnmovedRunsNothing, SkipsAndWaitsForMove (dirty and failing make update not retried until origin/main moves), FetchesAtMostEvery, UsesFetch, TestDaemonLoopSelfUpdateWaitsForPending (pending start delays to a later poll), TestDaemonSelfUpdateSetting (off runs nothing, on runs make update once over 3 polls), FetchMainWithinTimesOut and ExitCode all PASS. make fmt-check vet test pass with no FAIL. A scratch build's rhiza daemon settings prints rhiza.daemon.selfUpdate = false (default). make self-check fails only on stale .rhiza/rendered actor files from this run's records, not on the code.

## Links

### Out

- verifies [daemon-self-update](../capability/rhiza/orchestration/daemon/daemon-self-update.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs

# api-waits-on-run-clock-verified

evidence · `01M7EVPBJTP2EY8SN5NJM9Z8VW` · actor [claude-code-3e50ffb4](../actor/claude-code-3e50ffb4.md)

## result

PASS

## statement

At ad1e4ffd2, under the machine-tests lock and systemd-run cap: go test -count=1 -p 2 -trimpath ./internal/api -run 'TestRunClock|TestTypePointer|TestShipDeploy|TestSpawnSleep|TestRetry' -v exit 0 (TestRunClockFromEnv: RealClock unset/0/abc, ScaledClock at 1000; TestTypePointerScaledClockWaitsForRealWork sends after >=1s real on a 10 min timeout; TestTypePointerScaledClockTimeoutKeepsRealTime fails 'did not reach its prompt within 1s' after 1-3 s real; TestShipDeployScaledClockKeepsPollCount: scaled clock polls 10 times like the RealClock rule, in under 3 s real; TestSpawnSleepClock and TestRetryDefaultSleepClock: spawnSleep, retryGH and retryRex default sleeps go on the scaled clock when the variable is set and real time when it is unset). go test -count=1 -p 2 -trimpath ./internal/api -run 'TestSessionMessage|TestRun' exit 0 in 1.9 s with no 'still polls' Fatal, and existing tests unmodified (only run_clock_test.go added). go vet ./internal/api, rhiza ship generate and rhiza render exit 0. The diff touches only the listed files. It changes no production interval; the remaining time.Sleep/AfterFunc uses are in daemon, profile and watchdog code the leaf excludes; and there are no TODOs.

## Links

### Out

- verifies [api-waits-on-run-clock](../capability/rhiza/integrity/self-check/api-waits-on-run-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

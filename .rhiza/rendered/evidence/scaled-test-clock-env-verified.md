# scaled-test-clock-env-verified

evidence · `01M7EV5R3QC3QA9SKP48P61QD4` · actor [claude-code-ebd5cb4f](../actor/claude-code-ebd5cb4f.md)

## result

PASS

## statement

At 5afc8275d, under flock and the systemd-run cap: go test -count=1 -race -p 2 -trimpath ./internal/run -run 'TestScaledClock|TestTimedOut|TestWait' exit 0 (9.4s; all 19 tests pass, incl. TestWaitTimesOut, TestWaitSleepsOnlyToDeadline, TestWaitBackoffStopsAtTimeout unchanged on fakeClock); go vet ./internal/run exit 0. Diff read: clock.go's Sleep waits max(d/n, min(d,50ms)) and adds only skipped time to Now under a mutex (scaled-sleep-real-poll-floor, scaled-clock-now-skips-sleeps); Wait.TimedOut is exact on non-scaled clocks and adds the min(timeout,3s) real floor on ScaledClock; engine wait() uses it and sleeps its interval past the deadline. Tests cover Sleep(30s) in [50ms,100ms) with Now +30s, Sleep(50ms) ~50ms, 50ms of real time ~50ms of Now, concurrency under -race, ''/unset/0/-5/1.5/abc giving RealClock, the 10m/1s timeout floors, fake/RealClock exact, engine wait passing after 1s real and timing out in [3s,4s). RealClock, the nil-Clock default, workflow intervals and engine_test.go untouched; no TODOs. The full-package run, rhiza ship generate and rhiza render from Checks were not re-run: the brief limits the run to the named tests, and a verifier does not write files.

## Links

### Out

- verifies [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

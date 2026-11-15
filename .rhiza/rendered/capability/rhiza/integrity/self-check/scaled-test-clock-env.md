# scaled-test-clock-env

capability · `01M7EQC3EBD206T6TR63SVYR6Z` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -race ./internal/run -run 'TestScaledClock|TestTimedOut|TestWait' shows: at scale 1000 Sleep(30*time.Second) returns within 100 ms of real time and Now has advanced at least 30 s; real time spent outside Sleep (a 50 ms time.Sleep between two Now calls) advances Now by about 50 ms, not 50 s; Sleep(time.Microsecond) still waits at least 1 ms; Sleep and Now called from several goroutines pass under -race; with RHIZA_TEST_CLOCK_SCALE unset or empty the chosen clock is RealClock; the timeout helper on the scaled clock reports a 10 min timeout not yet passed after 30 s of virtual time but under 3 s of real time, and passed once both hold, reports a 1 s timeout passed only after 1 s of real time, and on RealClock and on an injected fake clock (engine_test.go's fakeClock) reports exactly elapsed >= timeout at once, with no real-time wait; the existing TestWaitTimesOut (6 polls, 5 sleeps), TestWaitSleepsOnlyToDeadline (sleeps [3s 2s]) and TestWaitBackoffStopsAtTimeout pass unchanged on fakeClock; an engine wait step with a 10 min timeout, on the scaled clock, whose condition holds only after 1 s of real time passes with outcome pass, and one whose condition never holds times out within 4 s of real time; failure cases: values 0, -5, 1.5 and abc each give RealClock and no panic. Touches internal/run/engine.go or a new internal/run/clock.go and their _test.go. Must not change RealClock, the Engine's nil-Clock default, any workflow's poll or backoff value, or the existing engine_test.go wait tests' expectations.

## statement

internal/run has a scaled Clock, chosen by RHIZA_TEST_CLOCK_SCALE=<n>, whose Sleep(d) waits d/n of real time (never less than 1 ms) and whose Now is the wall clock plus the time its Sleeps skipped, safe for concurrent use, and a timeout helper that a wait uses to decide its timeout passed: on the scaled clock, the timeout has passed on the clock's Now and at least the smaller of the timeout and 3 s of real time has passed since the wait began; on any other Clock (RealClock, or a fake clock a unit test injects) the helper reports exactly elapsed >= timeout on that clock's Now, with no real-time floor; the engine's wait() decides its step timeout with it and, on the scaled clock past its deadline inside the floor, sleeps its interval rather than zero; with the variable unset, empty or not a positive integer the engine keeps RealClock.

## Links

### Out

- decided-by [scaled-test-clock](../../../../decision/scaled-test-clock.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [scaled-clock-now-skips-sleeps](../../../../decision/scaled-clock-now-skips-sleeps.md)

### In

- depends-on [api-waits-on-run-clock](api-waits-on-run-clock.md)
- blocks [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f1](../../../../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f1.md)
- verifies [scaled-test-clock-env-verified](../../../../evidence/scaled-test-clock-env-verified.md)
- verifies [scaled-test-clock-env-qa-1](../../../../evidence/scaled-test-clock-env-qa-1.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

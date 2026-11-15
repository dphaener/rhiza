# scaled-test-clock-env-qa-1

evidence · `01M7EXVMNTGBHBRXN154N9QHR7` · actor [claude-code-e5658f1e](../actor/claude-code-e5658f1e.md)

## result

PASS

## statement

QA at 27088b322 with a scratch build (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza), rhiza run start --foreground on a wait step (timeout 2s, interval 1s, condition logs each check, never holds): with RHIZA_TEST_CLOCK_SCALE unset, empty, abc, 0, -5 and 1.5 each timed out 'did not hold within 2s' after exactly 3 checks in 2.06-2.09 s wall (RealClock, unchanged); at 1000, 40 checks in 2.11 s. At 1000: a 10m/30s wait whose condition holds after 1 s of real time passed after 18 checks in 1.10 s; a never-holding 10m/30s and 2m/30s wait timed out in 3.09/3.16 s (3 s real floor); a 100ms/30ms wait took 5 checks and 0.17 s on both clocks (sleeps under 50 ms unscaled). go test -count=1 -race ./internal/run -run 'TestScaledClock|TestTimedOut|TestWait|TestClock' ok.

## Links

### Out

- verifies [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

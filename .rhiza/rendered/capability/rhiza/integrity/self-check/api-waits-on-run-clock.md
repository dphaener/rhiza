# api-waits-on-run-clock

capability · `01M7ERPV4AJ3CJ76RBPT5DMSV0` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./internal/api -run 'TestRunClock|TestTypePointer' shows runClock is RealClock with the variable unset and the scaled clock with it set, and on the scaled clock a typePointer wait under the default 10 min timeout for a session that reaches its prompt only after 1 s of real time sends its message, while a 1 s timeout on a session that never reaches it fails only after at least 1 s of real time; go test -trimpath ./internal/api -run 'TestSessionMessage|TestRun' passes with the existing tests on useClock's scriptClock unchanged: session_message_test.go's 1 s timeout cases still fail with 'did not reach its prompt within 1s' without real-time waiting, and no scripted clock Fatals 'the script has ended but the run still polls'. Touches internal/api/run.go, run_spawn.go, api.go, ship.go, session.go, session_message.go, ship_deploy.go and their tests. Must not change a production interval, a workflow .toml, the daemon, feed and session-limit tickers, or the expectations of existing internal/api tests on injected clocks.

## statement

When RHIZA_TEST_CLOCK_SCALE is set, every wait a run makes in internal/api goes by the scaled clock: runClock is chosen from the variable (run.start and run.resume polls, wait conditions and backoff, session_run.go and session_message.go waits), the spawn poll in run_spawn.go, the default ctx.Sleep in api.go, ship.go and session.go, and ship_deploy.go's interval and its deadline; each of these that has a timeout decides it with internal/run's timeout helper, so on the scaled clock a wait on another process's real work keeps at least the smaller of its timeout and 3 s of real time, while on RealClock or a fake clock a test injects (useClock's scriptClock) it keeps exactly the rule of today.

## Links

### Out

- depends-on [scaled-test-clock-env](scaled-test-clock-env.md)
- decided-by [scaled-test-clock](../../../../decision/scaled-test-clock.md)
- decided-by [scaled-clock-now-skips-sleeps](../../../../decision/scaled-clock-now-skips-sleeps.md)

### In

- depends-on [scripts-run-on-scaled-clock](scripts-run-on-scaled-clock.md)
- blocks [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f3](../../../../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f3.md)
- verifies [api-waits-on-run-clock-verified](../../../../evidence/api-waits-on-run-clock-verified.md)
- verifies [api-waits-on-run-clock-qa-1](../../../../evidence/api-waits-on-run-clock-qa-1.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

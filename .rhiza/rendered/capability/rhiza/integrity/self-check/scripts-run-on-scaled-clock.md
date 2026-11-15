# scripts-run-on-scaled-clock

capability · `01M7ERPV9B2R3ZE8V900RBXASQ` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -v ./cmd/rhiza -run 'TestScripts/(run-ship-conflict-fix|run-ship-ci-fix|run-ship-main-dispatch-fail|arbiter-agent-tabs|plan-workflow-ci-infra-wait|session-message-busy)$' passes with each of the first five scripts' elapsed under 30 s (2026-11-15: 195, 173, 130, 108 and 96 s); failure case: a script that clears it (env RHIZA_TEST_CLOCK_SCALE=) runs on RealClock, shown by a short wait-timeout script in cmd/rhiza/testdata/script, run with go test -trimpath ./cmd/rhiza -run 'TestScripts/wait-timeout'. Speedup measure (orchestrator, quiet VM): go clean -testcache && /usr/bin/time -f '%e s' go test -trimpath ./cmd/rhiza, against 639 s. Touches cmd/rhiza/main_test.go and the new wait-timeout script. Must not change a production interval, a workflow .toml, or the daemon, feed and session-limit tickers.

## statement

cmd/rhiza/main_test.go's Setup sets RHIZA_TEST_CLOCK_SCALE to 1000 for every script, so the builtin workflows' 30 s polls and 2 min backoff cost milliseconds in scripts; a script that clears it runs on RealClock.

## Links

### Out

- depends-on [api-waits-on-run-clock](api-waits-on-run-clock.md)
- decided-by [scaled-test-clock](../../../../decision/scaled-test-clock.md)
- decided-by [scaled-clock-now-skips-sleeps](../../../../decision/scaled-clock-now-skips-sleeps.md)

### In

- blocks [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f2](../../../../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f2.md)
- verifies [scripts-run-on-scaled-clock-verified-2](../../../../evidence/scripts-run-on-scaled-clock-verified-2.md)
- verifies [scripts-run-on-scaled-clock-qa-1](../../../../evidence/scripts-run-on-scaled-clock-qa-1.md)
- discusses [note-01m7ewms48pkfg](../../../../note/note-01m7ewms48pkfg.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

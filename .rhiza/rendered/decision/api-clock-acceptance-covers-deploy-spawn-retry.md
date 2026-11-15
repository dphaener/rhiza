# api-clock-acceptance-covers-deploy-spawn-retry

decision · `01M7ES811H4YQVDX7FX04TDDA9`

## reasoning

The finding is right: the leaf's statement moves ship_deploy.go's interval and deadline (today time.Now plus time.Sleep at ship_deploy.go:85-99), run_spawn.go's time.Sleep(spawnPoll) fallback and the retryGH/retryRex default sleeps onto the scaled clock, but the acceptance only exercises runClock selection and typePointer. Dropping those points from the statement would defeat the change set's intent (the 2 s spawn polls are named in it, and ship deploy's wall-time deadline is the exact hazard the scaled-clock-now-skips-sleeps Decision calls out: a scaled sleep against a time.Now deadline would hammer a real HTTP probe about 1000 times as often). So the statement stays and the acceptance gains a runnable check per moved wait. A poll count equal to the RealClock count is the observable that proves deadline and sleep share one clock, and it is deterministic without real network (deployProbe against a local test server or the existing fake). The switch-with-the-variable check for spawn and retries mirrors the runClock selection check already in the acceptance, so it adds no new mechanism. This keeps the leaf's Touches and Must-not lines unchanged.

## statement

Extend api-waits-on-run-clock's acceptance rather than narrow its statement: go test -trimpath ./internal/api -run 'TestShipDeploy|TestSpawnSleep|TestRetry' (new or extended tests) shows that, with RHIZA_TEST_CLOCK_SCALE set, ship deploy against a probe that never turns healthy measures its deadline on the scaled clock it sleeps on, so its poll count for a given timeout and interval equals the RealClock count (it does not poll about scale times as often) and the call returns in far less real time than the timeout; and that the spawn poll's fallback and the default ctx.Sleep of retryGH and retryRex sleep on the scaled clock with the variable set and on time.Sleep (RealClock) with it unset. The existing TestRunClock/TestTypePointer and TestSessionMessage/TestRun points stay as written.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f3](../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f3.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

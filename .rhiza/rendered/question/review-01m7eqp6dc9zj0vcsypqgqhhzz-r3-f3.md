# review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f3

question · `01M7ES1Y2M1XZDF8BW8JPG8AW0`

## statement

Plan review round 3 left finding f3 (missing-acceptance on api-waits-on-run-clock): The statement moves ship_deploy.go's interval and deadline (today time.Now/time.Sleep, ship_deploy.go:85-99), run_spawn.go's spawnSleep (time.Sleep(spawnPoll) when ctx.Sleep is nil) and the default ctx.Sleep of retryGH (ship.go:610) and retryRex (session.go:892) onto the scaled clock. The acceptance checks only runClock's selection and typePointer: nothing runs to show that ship deploy measures its deadline on the same clock it sleeps on (the Decision scaled-clock-now-skips-sleeps names this case) or that the spawn poll and the retries use the scaled clock with the variable set and time.Sleep without it. The critic confirmed it: Verified: ship_deploy.go:85-99 measures its deadline with time.Now and sleeps time.Sleep(interval); run_spawn.go:253-258 falls back to time.Sleep(spawnPoll); the leaf's statement moves these and the retryGH/retryRex default sleeps onto the scaled clock. The acceptance runs only TestRunClock/TestTypePointer and the existing TestSessionMessage/TestRun: nothing shows ship deploy's deadline is measured on the clock it sleeps on (the case the scaled-clock-now-skips-sleeps Decision exists for: a scaled sleep against a wall deadline would hammer a real HTTP probe), nor that the spawn poll and the retries switch with the variable. Statement points with no acceptance; major stands. How should it be settled?

## Links

### Out

- blocks [api-waits-on-run-clock](../capability/rhiza/integrity/self-check/api-waits-on-run-clock.md)

### In

- answers [api-clock-acceptance-covers-deploy-spawn-retry](../decision/api-clock-acceptance-covers-deploy-spawn-retry.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

# review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f2

question · `01M7ES1Y0T5H4SAXCD5P9SVFQ0`

## statement

Plan review round 3 left finding f2 (missing-acceptance on scripts-run-on-scaled-clock): The failure case 'a script that clears it (env RHIZA_TEST_CLOCK_SCALE=) runs on RealClock, shown by a short wait-timeout script' does not say what the script asserts. Under the plan's own real-time floor (min(timeout, 3 s)), a short wait that times out takes the same wall time on the scaled clock as on RealClock, so a timing-based script cannot tell them apart. The acceptance must name an observable that differs, such as the condition's run count or the 'did not hold' detail after N checks. The critic confirmed it: Verified: no wait-timeout.txtar exists yet, and the leaf's failure case says only that the new script 'shows' RealClock. With the real-time floor min(timeout, 3 s), a short wait that times out takes the same wall time on both clocks, and wait()'s timeout detail ('did not hold within <timeout>[: last]') carries no poll count, so neither timing nor the timeout message tells the clocks apart. Without a named observable (condition run count, a 'held after N checks' detail) the point can be met by a vacuous script; an unverifiable acceptance point is missing-acceptance, major. How should it be settled?

## Links

### Out

- blocks [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

### In

- answers [wait-timeout-counts-checks](../decision/wait-timeout-counts-checks.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

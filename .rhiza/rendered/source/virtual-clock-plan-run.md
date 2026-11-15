# virtual-clock-plan-run

source · `01M7EQSVM2ZT6E0ZS24X83W28M`

## quotes

```json
[]
```

## raw_sha256

39e02f900d8393ac36032e85cdbfcf2b706b517704f7a2e68c14d8a95edf1128

## session

134c04c3-1d2d-41cf-a0ba-01a1dd7bc74e

## summary

Autonomous plan run with no owner; the orchestrator told the planner to review and complete the leaves drafted from milestone 38 without widening scope, run no test suite, and record its own choices as assumed Decisions. Research found that the drafted leaf said the scaled clock's Now runs n times faster than the wall clock, while the scaled-test-clock Decision says Now advances by each Sleep's whole d; a multiplied Now would turn real work (process starts, file reads) into 1000x longer virtual time and fire timeouts in internal/run's wait and internal/api's typePointer early. Clock.Now is used only for elapsed time within one process, so the planner chose: Now is the wall clock plus the time Sleeps skipped. It also found ship_deploy.go measures its deadline with time.Now, so its deadline must go by the same clock as its sleep, and that runClock is shared by goroutines, so the clock must be safe for concurrent use. Nothing was left open.

## title

Planning the virtual-clock change set, autonomous

## Links

### In

- derived-from [scaled-clock-now-skips-sleeps](../decision/scaled-clock-now-skips-sleeps.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

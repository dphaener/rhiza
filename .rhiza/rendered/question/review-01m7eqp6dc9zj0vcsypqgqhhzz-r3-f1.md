# review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f1

question · `01M7ES1XZDTXH5WSVV3CEJHG6B`

## statement

Plan review round 3 left finding f1 (contradiction on scaled-test-clock-env): The 1 ms Sleep floor plus the real-time timeout floor turns the scripts' own short waits into near-busy loops, which contradicts the scaled-test-clock Decision (no busy loop on the shared 4-core VM; polls keep their count). 42 wait steps across about 35 scripts in cmd/rhiza/testdata/script use timeout = "30s" with interval = "50ms" (e.g. until = "test -f ../go", a shell spawned per poll). At scale 1000 each 50 ms sleep is cut to the 1 ms floor, so the 30 s virtual deadline passes after about 600 condition spawns. Then, inside the 3 s real floor, engine.go wait() sleeps 'its interval' (1 ms again): up to about 3000 more spawns, against about 60 per 3 s today. The same steps' real budget for background work also drops from 30 s to 3 s. The plan neither bounds the poll rate inside the floor (e.g. a real-time minimum per poll) nor accounts for these scripts. The critic confirmed it: Verified: 43 kind="wait" steps in cmd/rhiza/testdata/script use timeout = "30s" with interval = "50ms" (29 until = "test -f ../go", the rest rhiza ship check), each condition a spawned process. engine.go wait() (902-957) sleeps c.Sleep(min(sleep, timeout-elapsed)) each poll; under the leaf's clock Sleep(50ms) waits the 1 ms floor, so polls run ~50x as often as today (about 1000 spawns/s instead of 20), and the leaf's 'sleeps its interval' inside the real-time floor is again 1 ms. Nothing in the statement or acceptance bounds the real poll rate, and these steps' real budget for background work drops from 30 s to 3 s, a flake risk on the loaded shared VM. This contradicts the scaled-test-clock Decision's promise that polls keep their count at a thousandth of the cost and AGENTS.md's no-busy-loop rule; major stands. How should it be settled?

## Links

### Out

- blocks [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)

### In

- answers [scaled-sleep-real-poll-floor](../decision/scaled-sleep-real-poll-floor.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

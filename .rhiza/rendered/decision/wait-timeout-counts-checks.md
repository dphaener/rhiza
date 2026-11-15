# wait-timeout-counts-checks

decision · `01M7ES64Q438NKMCBE5Q496344`

## reasoning

Read the finding, the leaf's acceptance, the scaled-test-clock, scaled-clock-now-skips-sleeps and scaled-sleep-real-poll-floor Decisions, and engine.go wait() (902-957). The finding is right that wall time cannot separate the clocks (the real-time floor min(timeout, 3 s) makes a short timed-out wait take its full timeout on both) and that the timeout detail carries no count. The check count does differ, and deterministically enough: on RealClock wait() sleeps min(interval, timeout-elapsed), so a 2 s timeout at a 1 s interval checks exactly 3 times; on the scaled clock with the 50 ms real poll floor, Now passes 2 s after 3 checks in about 100 ms of real time, then the floor keeps the wait for 2 s of real time while each poll sleeps its interval, about 50 ms real, giving roughly 40 checks. A threshold of at least 10 leaves wide margin for a loaded VM, and exactly 3 on RealClock is the existing unit-test rule. A 1 s interval is used rather than the scripts' usual 50 ms because sleeps of 50 ms or less are not scaled at all, so a 50 ms interval would give the same count on both clocks. Running both halves in one script also shows the Setup really sets the scale, which a RealClock-only run would not. A counter file is an ordinary script observable and needs no change to wait()'s detail, so no production output changes.

## statement

The wait-timeout script tells the clocks apart by the condition's check count, not by wall time or the timeout message: it runs one workflow whose wait step has timeout = "2s" and interval = "1s" and whose until condition appends a line to a file each time it runs, twice in the same script. With RHIZA_TEST_CLOCK_SCALE cleared (env RHIZA_TEST_CLOCK_SCALE=) the step times out with 'did not hold within 2s' after exactly 3 checks (at 0, 1 and 2 s, as TestWaitTimesOut's checks fall on the deadline), asserted as the file holding exactly 3 lines; with the Setup's scale of 1000 the same step times out after at least 10 checks, since past the virtual deadline each 1 s poll inside the 2 s real-time floor sleeps about 50 ms. scripts-run-on-scaled-clock's failure case names these two counts as what the script asserts; the script costs about 4 s of real time.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f2](../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f2.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

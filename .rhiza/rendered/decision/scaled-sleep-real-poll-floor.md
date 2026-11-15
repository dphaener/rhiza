# scaled-sleep-real-poll-floor

decision · `01M7ES40MJZEMD7J56897K5Z4X`

## reasoning

Read the finding, the scaled-test-clock and scaled-clock-now-skips-sleeps Decisions, scaled-test-clock-env's statement and acceptance, and engine.go wait() (902-957). The scripts' own wait steps use interval = "50ms" 117 times in 91 scripts (one 100ms, one 1s, eleven 30s), each poll spawning a shell or a rhiza process; with a 1 ms floor each would poll about 50 times as often and see its 30 s budget cut to the 3 s real floor, a busy loop on the shared 4-core VM that the scaled-test-clock Decision and AGENTS.md rule out and a flake risk for background work. A per-sleep real-time minimum fixes both at the source: with min(d, 50 ms) the 50 ms polls run exactly as today (same rate, the full 30 s of real time, since they skip nothing), every other poll is bounded to about 20 a second, and the speedup the change set is for is kept: the builtin workflows' 30 s polls, 2 s spawn polls and 2 min backoff, the waits that cost the measured wall time, still cost about 50-120 ms each. 50 ms matches the scripts' interval, so no script needs editing; a lower floor would still multiply their poll rate and shrink their real budget, a higher one buys nothing and slows the long polls. Bounding the rate in wait() alone would leave the other scaled loops (run polls, spawn poll, typePointer, ship deploy) to the 1 ms floor, so the floor belongs in the clock's Sleep.

## statement

The scaled clock's Sleep(d) waits max(d/n, min(d, 50 ms)) of real time instead of max(d/n, 1 ms): a sleep of 50 ms or less is not scaled at all and skips nothing, and a longer one waits at least 50 ms of real time, so no poll on the scaled clock runs more than about 20 times a second of real time; the real-time timeout floor (the smaller of the timeout and 3 s) and the rule that each poll inside it sleeps its interval stay as decided; scaled-test-clock-env's acceptance replaces 'Sleep(time.Microsecond) still waits at least 1 ms' with checks that Sleep(50*time.Millisecond) waits about 50 ms and advances Now by about 50 ms, and that Sleep(30*time.Second) waits at least 50 ms and under 100 ms

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f1](../question/review-01m7eqp6dc9zj0vcsypqgqhhzz-r3-f1.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

# scaled-clock-now-skips-sleeps

decision · `01M7EQT8V4D1RPDFHF8NJ02NSK`

## reasoning

Clock.Now is used only for elapsed time within one process (internal/run engine.go wait, internal/api session_message.go typePointer), so a per-process offset is enough. A Now multiplied by the scale would count real work, such as a 50 ms process start, as 50 s and fire step and message timeouts early. Adding only skipped sleep time keeps each loop's poll count as in production, but it does not by itself keep the real time another process gets: a poll that sleeps the 1 ms floor while Now skips the rest of its 250 ms interval lets a 10 min typePointer wait for a fake agent end after a few seconds of real time, and a 1 s --timeout (session-message-busy.txtar) after about 4 polls in a few ms; engine.go wait() step timeouts likewise wait on other processes. So a timeout also needs a real-time floor: the smaller of the timeout itself and 3 s, enough for a fake agent or a background rhiza command on a loaded 4-core VM while a deliberately expired step timeout in a script costs at most 3 s. On RealClock the real time since the wait began equals Now's elapsed time, so the floor never delays a production timeout. Past the clock's deadline and inside the floor, each poll sleeps its normal interval on the clock, never zero, so the loop does not spin. A loop that sleeps scaled but measures its deadline with time.Now, as ship_deploy.go does, would poll 1000 times as often. runClock is a package variable read by concurrent goroutines, so the offset needs a lock or atomic.

## statement

The scaled clock's Now is the wall clock plus the time its Sleeps skipped (each Sleep(d) adds d minus the real time it waited), not the wall clock multiplied by the scale, and it is safe for concurrent use; every loop that sleeps on it measures its deadline with the same clock's Now, and a timeout on it fires only once the timeout has passed on Now and at least the smaller of the timeout and 3 s of real time has passed since the wait began, so a wait on another process's real work keeps that much real time; on RealClock the rule is the same as today

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [virtual-clock-plan-run](../source/virtual-clock-plan-run.md)

### In

- decided-by [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)
- decided-by [api-waits-on-run-clock](../capability/rhiza/integrity/self-check/api-waits-on-run-clock.md)
- decided-by [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

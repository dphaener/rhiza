# scaled-test-clock

decision · `01M7EQC3D2KVVGGH3Y139X7EAD`

## reasoning

The brief allows a virtual or a scaled clock. Scripts poll for work other processes do (fake agents, background rhiza commands), which needs some real time to happen; a clock whose Sleep returns at once turns each such poll into a busy loop on a shared 4-core VM that has already been OOM-crashed once, while a scaled clock with a floor keeps every poll a real, cheap wait.

## statement

Testscripts get a scaled clock, not one where Sleep returns at once: each Sleep(d) waits d divided by a scale in real time (at least 1 ms) while Now advances by the whole d, so waits, polls, backoffs and timeouts keep their order and count but cost about a thousandth of their wall time; production keeps RealClock

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [scaled-test-clock-env](../capability/rhiza/integrity/self-check/scaled-test-clock-env.md)
- decided-by [api-waits-on-run-clock](../capability/rhiza/integrity/self-check/api-waits-on-run-clock.md)
- decided-by [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).

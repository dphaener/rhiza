# runaway-session-fails-step-qa-1

evidence · `01M7QPPEW5WMT6SX3AY8YSV4ZP` · actor [claude-code-9dc16082](../actor/claude-code-9dc16082.md)

## result

PASS

## statement

QA at b9129e145 (no real daemon or runaway, per brief): go test -trimpath -p 2 -count=1 ./internal/api -run 'TestWatchdogRunaway|TestWatchdogStalled|TestWatchdogPolicy' exit 0 (TestWatchdogRunaway, SpawnedChild, NoHostTable, LeftAlone under_both/no_tree/attended/back_under/threshold_off, all Stalled tests); go test ./internal/daemon -run Resumable exit 0 (a runaway-stopped run is not resumed, a waiting stop is). Code read: watchdogProcTable is nil under RHIZA_TEST_NOW, the stall pass carries the runaway nudge memory from prev, a stopped run is dropped from the stall pass that poll, the stopped detail begins 'watchdog: runaway ' which Resumable skips, and an owner resume re-entering the step clears Stopped.

## Links

### Out

- verifies [runaway-session-fails-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-session-fails-step.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason

# runaway-session-fails-step-verified

evidence · `01M7QNYV5Y96PSJMQF65J2E13Z` · actor [claude-code-18156fbc](../actor/claude-code-18156fbc.md)

## result

PASS

## statement

At f300dfbe1: go test -trimpath -p 2 ./internal/api -run 'TestWatchdogRunaway|TestWatchdogStalled' exit 0 (TestWatchdogRunaway, SpawnedChild, LeftAlone with under-both/back-under/threshold-off/attended/no-tree subtests, NoHostTable, all Stalled tests PASS); -run 'Watchdog|Daemon' exit 0; ./internal/daemon/... exit 0; go vet and gofmt clean; rhiza ship generate and rhiza render in a clean detached worktree exit 0 with no changes. Diff read: watchdogRunaway runs before watchdogStalled on an injected procs reader (nil under RHIZA_TEST_NOW), nudges with count and threshold, after runawayGrace closes panes plus waiting parents and calls terminateRun, journals TypeStopped on the current step naming session/role, processes, RSS and threshold, files a 'runaway session' bug, never resumes (a stopped run is not Ended so the transient pass skips it); nudge memory is WatchdogRun.Runaway in internal/daemon/state.go. Stall limits and nudge text are unchanged (stallReason only extracted sessionRole), there is no /proc read or ps in tests, and the holds are untouched.

## Links

### Out

- verifies [runaway-session-fails-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-session-fails-step.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason

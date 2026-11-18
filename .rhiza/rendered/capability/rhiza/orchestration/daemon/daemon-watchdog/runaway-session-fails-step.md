# runaway-session-fails-step

capability · `01M7QN2RG4WGQBDGHX67B53556` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath -p 2 ./internal/api -run 'TestWatchdogRunaway|TestWatchdogStalled' passes, with a new test driving the watchdog through an injected fake process table (never the host's ps) for an autonomous run whose step session's tree holds 450 processes: the first poll sends that session a nudge naming the count and the threshold; a poll runawayGrace later with the tree still over closes the session's pane, terminates the run, journals a stopped entry on the current step whose detail names the session, 450 processes, the resident memory and the runawayProcs threshold, queues no continue (no resume, no --retry) and files a runaway bug; a spawned child's tree over runawayMemory is caught the same way and named with its role. (failure case) a tree under both thresholds is neither nudged nor stopped; a tree back under the thresholds after the nudge is not stopped and its nudge clears; a threshold set to 0 is never passed; an owner-attended run under autonomousOnly is only logged; a session with no tree in the table is left alone; the daemon reads no host process table under RHIZA_TEST_NOW

## files

internal/api/watchdog_runaway.go, internal/api/watchdog_runaway_test.go, internal/api/daemon_watchdog.go (watchdogDeps gains the process-table reader and a stopped-entry journal dep, watchdogPasses runs the pass, daemonWatchdog wires ps through internal/proctree), internal/api/watchdog_bugs.go (a runaway bug kind), internal/daemon state file holding WatchdogRun (the runaway nudge time)

## must_not

Change the stall pass's limits, nudge text or retry behaviour, or session.Stallable; resume or retry a run stopped as a runaway; read /proc; run ps in a test; touch the pre-start or spawn holds

## statement

Each watchdog poll reads the process tree of every live run's stall sessions and treats one over runawayProcs or runawayMemory like a stall: it nudges the session with the reading, then, still over a threshold runawayGrace later, stops the run and journals its step stopped with the reading (session, process count, resident memory and the threshold passed) as the reason, without resuming it

## Links

### Out

- depends-on [session-process-tree-reading](session-process-tree-reading.md)
- depends-on [runaway-threshold-keys](runaway-threshold-keys.md)
- decided-by [runaway-fails-not-retries](../../../../../decision/runaway-fails-not-retries.md)

### In

- verifies [runaway-session-fails-step-verified](../../../../../evidence/runaway-session-fails-step-verified.md)
- verifies [runaway-session-fails-step-qa-1](../../../../../evidence/runaway-session-fails-step-qa-1.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason

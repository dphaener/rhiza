# stall-covers-spawned-children

capability · `01M7QBVD8Y319MBF5A5DA03TQH` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run 'TestWatchdogStalled|TestDaemonNudgeStalled|TestStepSessions' passes, including a new test that replays the 2026-11-16 case: a conduct run whose conductor session spawned a fixer (a spawn entry, no spawn-ended), the fixer's record autonomous, idle by its hook, with no done mark and its last report past the agent nudge limit; watchdogStalled nudges the fixer's session (not the conductor's) at the first limit, and past the second limit stops the run with both sessions' panes closed and resumes it with --retry, the run's budget charged once; daemonNudgeStalled sends the nudge to the fixer only. (failure case) a fixer that said done, or whose spawn has a spawn-ended entry, is neither nudged nor counted; a conductor waiting on a fixer whose last report is recent is not nudged and its run is not stopped, however old the conductor's own report

## files

internal/api/run.go (stepSessions and its callers' doc comments), internal/api/daemon_watchdog.go (watchdogProgress, watchdogStalled, the stop dep), internal/api/daemon_watchdog_test.go, internal/api/daemon_nudge_test.go, internal/api/run_test.go

## must_not

Change session.Stalled or session.Stallable, the watchdog policy keys or their defaults, rhiza run spawn's own behaviour, or any --json shape; read process trees (the runaway change set's); nudge or stop sessions of owner-attended runs beyond what the existing autonomousOnly policy allows

## statement

The watchdog's stalled-step pass and the daemon's session nudge cover the children a session spawns with rhiza run spawn: an idle child with no done mark is nudged, then its step stopped and retried, and the parent waiting on it is never nudged for its own silence

## Links

### Out

- decided-by [stall-spawned-children-count](../../../../../decision/stall-spawned-children-count.md)

### In

- depends-on [stall-reason-journaled](stall-reason-journaled.md)
- verifies [stall-covers-spawned-children-verified](../../../../../evidence/stall-covers-spawned-children-verified.md)
- verifies [stall-covers-spawned-children-qa-1](../../../../../evidence/stall-covers-spawned-children-qa-1.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch

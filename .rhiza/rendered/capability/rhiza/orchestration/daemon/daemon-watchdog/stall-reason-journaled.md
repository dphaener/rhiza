# stall-reason-journaled

capability · `01M7QBVDARKZGFXHC5P2ARXQHF` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run 'TestWatchdogStalled' passes, including a new test where a stalled spawned fixer's retry appends one warning entry to the run's journal whose detail names the fixer's session ID, the role fixer, its parent session and its last report time, and the charged attempt's action names the same session; and where, with the retry budget spent, the hold's detail and its inbox item name that session and role. (failure case) a step that is only nudged, or that moved since, journals no warning; a stalled step's own session (no spawn) is named as the step's session with no parent

## files

internal/api/daemon_watchdog.go (watchdogStalled, watchdogDeps, daemonWatchdog wiring, watchdogHoldItems), internal/api/daemon_watchdog_test.go

## must_not

Fail the run or add a run state or journal entry type (the existing warning entry carries the reason); change the retry budget or limits; change the bug request filing beyond naming the session in its evidence

## statement

Every stall retry and give-up of the watchdog names its reason on the run: the stalled session, its role or spawn parent, and when it last reported, as a warning in the run's journal, in the charged attempt and in the hold's detail

## Links

### Out

- depends-on [stall-covers-spawned-children](stall-covers-spawned-children.md)
- decided-by [stall-reason-on-run](../../../../../decision/stall-reason-on-run.md)

### In

- verifies [stall-reason-journaled-verified](../../../../../evidence/stall-reason-journaled-verified.md)
- verifies [stall-reason-journaled-qa-1](../../../../../evidence/stall-reason-journaled-qa-1.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch

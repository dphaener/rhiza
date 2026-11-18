# stall-covers-spawned-children-verified

evidence · `01M7QCBDD0HW8BA1079TKJ1G3B` · actor [claude-code-6cb605cb](../actor/claude-code-6cb605cb.md)

## result

PASS

## statement

Verified commit 3c125fa0b. Ran go test -trimpath -p 2 ./internal/api -run 'TestWatchdogStalled|TestDaemonNudgeStalled|TestStepSessions|TestRunStalled' (exit 0; TestWatchdogStalledSpawnedChild, TestDaemonNudgeStalledSpawnedChild, TestStepSessions, TestWatchdogStalled, TestDaemonNudgeStalled, TestRunStalledNow pass), rhiza ship generate (exit 0, 0 changed) and rhiza render (exit 0, only this verifier's own actor files). The full ./internal/api run in Checks was not re-run per the brief (builder ran it; scoped run raised no doubt). Diff read: stepSessions adds open spawn children of the current step; stallSessions drops parents waiting on an open spawn and feeds runStalledAll (daemon nudge) and watchdogProgress; watchdogStalled's stop closes sids plus spawnWaiting. TestWatchdogStalledSpawnedChild replays 2026-11-16: conduct run, conductor spawn of fixer with no spawn-ended, fixer autonomous idle by hook with no done; nothing before 20m, fixer nudged (not conductor) at 21m, at 41m stop R:fixer,conductor and one resume (queueContinue retry=true) with Spent 1; done child and ended spawn not nudged/counted; conductor waiting on a reporting fixer never nudged or stopped despite a 5h-old report. TestDaemonNudgeStalledSpawnedChild sends to the fixer only. must_not respected: session.Stalled/Stallable, policy keys/defaults, run spawn, --json shapes untouched; autonomousOnly hold path unchanged. gofmt clean, no TODOs, only listed files changed besides the actor record.

## Links

### Out

- verifies [stall-covers-spawned-children](../capability/rhiza/orchestration/daemon/daemon-watchdog/stall-covers-spawned-children.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch

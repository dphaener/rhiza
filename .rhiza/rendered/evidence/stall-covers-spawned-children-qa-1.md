# stall-covers-spawned-children-qa-1

evidence · `01M7QDDCNFSEE6JTM03DT8KCZB` · actor [claude-code-9ff9493e](../actor/claude-code-9ff9493e.md)

## result

PASS

## statement

QA with a scratch build of ./cmd/rhiza (fake rex/claude, scratch root): a run whose session spawned a fixer that went idle with no done mark; rhiza daemon --once at 00:19 did nothing, at 00:21 nudged the fixer only (not the conductor), at 00:41 closed both sessions' panes, stopped the run and queued one --retry continue (spent 1 of 2), which the next poll resumed. A fixer re-reporting at 00:40 kept the 50-minute-silent conductor from being nudged or stopped; a fixer that said done (spawn-ended) was never nudged or counted. Scoped go test (TestWatchdogStalled|TestDaemonNudgeStalled|TestStepSessions|TestRunStalled) passes.

## Links

### Out

- verifies [stall-covers-spawned-children](../capability/rhiza/orchestration/daemon/daemon-watchdog/stall-covers-spawned-children.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch

# m40-stalls-conduct-plan

source · `01M7QBTMF7Y263XNM9EEB3SSG9`

## quotes

```json
[
  {
    "by": "orchestrator",
    "text": "On 2026-11-16 a run sat 9 hours because a fixer never ran `rhiza session done` and its conductor's spawn waited on it. ... The daemon should detect and recover (nudge, then retry or fail the step with a reason), for every workflow's agent steps.",
    "where": "hands-off-milestone-interview.md, point 2"
  },
  {
    "by": "orchestrator",
    "text": "The daemon's stall detection needs the same distinction (idle with a live background job vs idle with nothing running), read portably (ps on Linux and macOS), with a longer bound for the background case, and the runaway check reads the same process tree.",
    "where": "planner-notes/m40-stalls.md"
  }
]
```

## raw_sha256

4412a3c62aac6adb0b90db25647276b8cb51e7e3253d5ad9cd9e185f36f282c7

## summary

The conductor of run 01M7QBNRXR planned change set 01M7NVMG (m40-stalls) with no owner present, from the milestone-40 brief (hands-off-milestone-interview.md, point 2 and point 9) and the orchestrator's planner notes for m40-stalls. Research found the watchdog's stalled-step pass (watchdogStalled) and the session nudge (daemonNudgeStalled) read a step's sessions through stepSessions in internal/api/run.go, which covers the step's own session, leaves, checks and change-set fixers but never the children of rhiza run spawn (run.State.Spawns), so a conductor's idle child is never nudged and the conductor's own long-running spawn command is what goes silent. The planner notes' idle-with-a-background-job distinction needs the agent's process tree, which the runaway change set (m40-runaway) reads, so it is drafted apart. No owner words: autonomous run.

## title

Stall change set plan, autonomous conduct run 01M7QBNRXR

## Links

### In

- derived-from [stall-spawned-children-count](../decision/stall-spawned-children-count.md)
- derived-from [stall-reason-on-run](../decision/stall-reason-on-run.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch

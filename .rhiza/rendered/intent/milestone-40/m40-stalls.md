# m40-stalls

intent · `01M7NVQ7YR969RM0WAPQZN97HG` · parent [milestone-40](../milestone-40.md)

## statement

Stall detection covers every agent session of a run, the children a conductor spawns included: a child idle with no done mark is nudged, then its step retried or failed with a reason, whatever the workflow. Measure: a test replays the 2026-11-16 case (a conduct run's fixer idle with no session done while the conductor waits on spawn) and sees the nudge, then the retry.

## Links

### Out

- decided-by [m40-push-wakeup-alarms-folded](../../decision/m40-push-wakeup-alarms-folded.md)

### In

- serves change set `01M7NVMG0EJFRJV03G025C6366` (The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

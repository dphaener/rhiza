# m40-push-wakeup

intent · `01M7NVQ89GPX52MA4CE6R0AQF7` · parent [milestone-40](../milestone-40.md)

## statement

A conductor no longer blocks on spawn: spawn returns once the child starts, and the child's session done or blocked question wakes the conductor with a session message. Measure: change set 01M7M2FB5BNRSDSN7HFWQSC2WN is landed.

## Links

### Out

- decided-by [m40-push-wakeup-alarms-folded](../../decision/m40-push-wakeup-alarms-folded.md)

### In

- serves change set `01M7M2FB5BNRSDSN7HFWQSC2WN` (rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17))

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

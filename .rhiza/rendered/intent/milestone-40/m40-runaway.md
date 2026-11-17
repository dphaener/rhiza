# m40-runaway

intent · `01M7NVQ7XNNV3W6T3PPSZY6SVH` · parent [milestone-40](../milestone-40.md)

## statement

The watchdog treats a runaway agent session of any workflow (its process tree's process count or resident memory past a host threshold, read portably through ps on Linux and macOS) like a stall: it nudges, then fails the step with the reading as its reason. Measure: a test with a fake process table over the threshold sees the nudge, then the step failed with that reason.

## Links

### Out

- decided-by [m40-host-limits-as-holds](../../decision/m40-host-limits-as-holds.md)

### In

- serves change set `01M7NVMDMBVBKRW5B8SB0KMW6T` (The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

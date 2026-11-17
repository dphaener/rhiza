# m40-leaves-per-change-set

decision · `01M7NVQ8XC6XHD41K0E8XCV67R`

## reasoning

Each change set touches a different part of the daemon (the hold poll, the watchdog, spawn, the run engine, request intake), each needs its own code research to write leaf contracts, and the record's done points already fix every acceptance target. Planning nine change sets' leaves here would make one plan review too large to read.

## statement

No change set of milestone 40 gets its leaves in this plan run: each is planned in its own plan run when it is proposed

## wording

(assumed, would have escalated) the planner's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

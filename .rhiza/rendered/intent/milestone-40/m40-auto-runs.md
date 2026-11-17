# m40-auto-runs

intent · `01M7NVQ871787K9GVXZ4DVWK8Z` · parent [milestone-40](../milestone-40.md)

## statement

rhiza.daemon.maxRuns=auto starts another run only while the pre-start checks (load, memory, free tmp, disk) have room, and a number still pins the limit. Measure: a test with readings that have room starts a second run, and with a held reading starts none.

## Links

### Out

- decided-by [m40-maxruns-auto](../../decision/m40-maxruns-auto.md)

### In

- serves change set `01M7NVMXZ1YEGZB56ZEAMPX44B` (rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

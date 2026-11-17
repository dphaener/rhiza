# m40-maxruns-auto

decision · `01M7NVQ8MRHGXM379MRCXFR2EF`

## question

Concurrency (brief point 6): should the daemon pick its own maxRuns from the host's cores and memory in this milestone?

## statement

rhiza.daemon.maxRuns takes auto: the daemon starts another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit

## wording

Yes, auto setting (Recommended), picked by the owner, relayed by the orchestrator

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)
- decided-by [m40-auto-runs](../intent/milestone-40/m40-auto-runs.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

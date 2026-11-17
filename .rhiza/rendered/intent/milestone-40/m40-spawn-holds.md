# m40-spawn-holds

intent · `01M7NVQ7WKPYHQDNDYF5Q6ZT0K` · parent [milestone-40](../milestone-40.md)

## statement

The daemon's load, memory and disk holds, plus a new free-tmp hold, gate every child session spawn of any workflow (rhiza run spawn included), not only run starts: a held spawn waits, and the run journal and the inbox name the reading and the limit. Measure: a test drives a spawn under a low-memory and a low-tmp reading and sees it held, then started once the reading recovers.

## Links

### Out

- decided-by [m40-host-limits-as-holds](../../decision/m40-host-limits-as-holds.md)

### In

- serves change set `01M7NVMAZ0PV2BP11NSNF04ANZ` (The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

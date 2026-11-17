# m40-host-limits-as-holds

decision · `01M7NVQ8GQMKYJYJ8Z0M3QH4J7`

## question

The orchestrator's proposal: no scopes or test slots; extend the load, memory and disk holds to every spawn and to free tmp; treat a runaway session as a stall; keep the leak and the open-pr RSS as bugs. Does that make sense?

## quote

Yes that makes sense to me.

## quote_where

owner, 2026-11-17 orchestrator chat, relayed by the orchestrator

## statement

Host limits are the daemon's holds and watchdog, generic for any workflow: the load, memory and disk holds also gate every child session spawn and gain a free-tmp hold; the watchdog treats a runaway session (process count or memory of its process tree past a host threshold, read portably) like a stall, nudging and then failing the step with a reason. The /tmp projection leak (01M7M3JZFW) and the 2.2 GB RSS of rhiza ship open-pr stay bug fixes, not this milestone's limits

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)
- decided-by [m40-spawn-holds](../intent/milestone-40/m40-spawn-holds.md)
- decided-by [m40-runaway](../intent/milestone-40/m40-runaway.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

# milestone-40

intent · `01M7NVQ7VJTTKA2PX270W95ENA`

## milestone

40

## statement

The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it holds child session spawns while the box lacks room, detects stalled and runaway sessions and recovers them, cleans each run's tmp, retries failed runs from the right step with their findings, hands planners the notes queued with a request, wakes conductors when a child finishes, and records every outside intervention, so intervention stays the exception (2026-11-17: five fixes went from plan to merge under the daemon, #612-#621, but only with an orchestrator session sending machine rules, watching stalls, sweeping tmp, setting maxRuns and retrying failed runs by hand).

## Links

### Out

- decided-by [m40-no-systemd-scopes](../decision/m40-no-systemd-scopes.md)
- decided-by [m40-host-limits-as-holds](../decision/m40-host-limits-as-holds.md)
- decided-by [m40-retry-right-step](../decision/m40-retry-right-step.md)
- decided-by [m40-maxruns-auto](../decision/m40-maxruns-auto.md)
- decided-by [m40-push-wakeup-alarms-folded](../decision/m40-push-wakeup-alarms-folded.md)
- decided-by [m40-before-web-split](../decision/m40-before-web-split.md)
- decided-by [m40-hands-off-measure](../decision/m40-hands-off-measure.md)
- decided-by [m40-leaves-per-change-set](../decision/m40-leaves-per-change-set.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

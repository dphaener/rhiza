# m40-retry-right-step

decision · `01M7NVQ8JNZR7GA4RMQQ36Y7VH`

## question

Failed runs (brief point 7): today the watchdog retries only failures its classifier calls transient, and a real failure goes to the inbox. How should the daemon retry a real failure?

## statement

A failed run restarts at the step its workflow names for the failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failure's output and findings in the next attempt's brief; one retry by default, then the inbox

## wording

Right step + findings (Recommended), picked by the owner, relayed by the orchestrator

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)
- decided-by [m40-retry](../intent/milestone-40/m40-retry.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

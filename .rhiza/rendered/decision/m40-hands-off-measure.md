# m40-hands-off-measure

decision · `01M7NVQ8TSNDCRTCB1B6PXS2Q1`

## reasoning

The owner's bar is 'intervention should be the exception' for any workflow. Ten runs across all three kinds is enough to show that without tying the milestone to weeks of traffic, and the intervention record makes it a count rather than a judgement. The shallow-check and tmp-leak bugs are named because they are the remaining interventions the orchestrator makes; they ship through plan-bug, not here.

## statement

Milestone 40's done measure is 10 consecutive daemon runs on this repo, at least 2 bug, 2 feature and 1 bring-your-own, with zero journalled interventions, the two queued bugs (01M7KW701G, 01M7M3JZFW) fixed and the orchestrator's stopgaps retired

## wording

(assumed, would have escalated) the planner's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)
- decided-by [m40-hands-off](../intent/milestone-40/m40-hands-off.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

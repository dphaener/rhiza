# m40-before-web-split

decision · `01M7NVQ8RT6WWDENG7DCYPQ0NC`

## reasoning

This milestone removes the orchestrator session that the web split's runs would otherwise need beside the daemon, so building it first lets the web split run hands-off.

## statement

Milestone 40 goes before the standalone rhiza web and internal/api split milestone

## wording

(assumed, would have escalated) the orchestrator's assumption, not asked of the owner

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

# m40-byo

intent · `01M7NVQ8CJZTDAPMT4NAC2QBP3` · parent [milestone-40](../milestone-40.md)

## statement

A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, and the spawn holds, stall and runaway detection, failure retry and intervention record apply to its steps. Measure: the change set's testscript passes in make test.

## Links

### In

- serves change set `01M7NVN06JNCCW1NVGPD4WZB16` (A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

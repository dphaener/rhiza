# m40-hands-off

intent · `01M7NVQ8DPH7H40EB1HR95HTFA` · parent [milestone-40](../milestone-40.md)

## statement

Intervention is the exception: 10 consecutive daemon runs on this repo (at least 2 bug, 2 feature and 1 bring-your-own) finish with zero interventions in their journals, with queued bugs 01M7KW701G (shallow checks) and 01M7M3JZFW (projection leak) fixed and the orchestrator's stopgaps retired (conductor-rules.md no longer sent, the rhiza-runwatch and rhiza-tmpsweep units stopped, planner-notes unused). Measure: the intervention count of @m40-interventions over the last 10 finished runs.

## Links

### Out

- decided-by [m40-hands-off-measure](../../decision/m40-hands-off-measure.md)

### In

- serves change set `01M7Q2Q68JE50NCRKGCT5Z3APZ` (Fix watchdog transient failure at step plan/close-out: run 01M7PSG0M7APCTB4V1YKPTM394, request 01M7KHWSWG3ENZ81KH148SK1GD; log: main was not updated: git pull --ff-only origin main: From https://github.com/dphaener/rhiza | * branch main -> FETCH_HEAD | fatal: Cannot rebase onto multiple branches.)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

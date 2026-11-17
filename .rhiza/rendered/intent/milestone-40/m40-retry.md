# m40-retry

intent · `01M7NVQ80NN5H2ZPD51KSEEZHH` · parent [milestone-40](../milestone-40.md)

## statement

A failed run is retried from the step its workflow names for the failure, with the failure's output and findings in the next attempt's brief, once by default, then goes to the inbox. Measure: a test fails a plan run at review-check and sees the retry start at plan with the refusal in the brief, and a second failure reach the inbox.

## Links

### Out

- decided-by [m40-retry-right-step](../../decision/m40-retry-right-step.md)

### In

- serves change set `01M7NVMMPK3A9SRVA9Q63N5SR9` (A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md

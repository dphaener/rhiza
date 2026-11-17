# m40-interventions

intent · `01M7NVQ83YD8JE3TEPTPHSM25G` · parent [milestone-40](../milestone-40.md)

## statement

Every outside intervention in a run is journalled and counted: a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, and a hand commit on the run's branch. Approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted. Measure: a test makes each kind of intervention and sees each in rhiza run status --json and in the daemon's count; a run with none counts zero, and a run whose only owner action is approving its workflow's signoff step also counts zero.

## Links

### In

- serves change set `01M7NWFV54PDK3XDSGH07HK0K7` (The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md
- `01M7NW2XQ9TRY5N3A7MGRAH022` apply /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza-the-daemon-runs-any-workflow/d2542435-eaf2-4a09-ab5d-01a1ebc04a37/scratchpad/fix.json

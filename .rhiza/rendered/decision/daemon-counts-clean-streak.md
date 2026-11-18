# daemon-counts-clean-streak

decision · `01M7R2FXT3P77EXXQMC4RY23ME`

## statement

rhiza daemon status shows each live run's intervention count, and over the project's finished runs the number finished, the number finished with zero interventions and the current streak of consecutive finished runs with zero, newest first

## wording

(assumed, would have escalated) the streak serves the milestone's measure of consecutive runs with zero interventions

## Links

### Out

- derived-from [m40-interventions-plan-research](../source/m40-interventions-plan-research.md)

### In

- decided-by [daemon-status-interventions](../capability/rhiza/workflows/workflow-runs/daemon-status-interventions.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

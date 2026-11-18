# intervention-is-a-person-command

decision · `01M7R2FXPR9X3V95KZABFZM5RA`

## statement

An intervention is an action a person (or an orchestrator acting for one) takes from outside the run: the daemon's own resume, the watchdog's retry, the engine's hand-back and wake messages, messages between the run's own sessions and approving a signoff step the run waits on are never interventions

## wording

(assumed, would have escalated) the line drawn at who made the action: the run and its daemon act as written; anyone else intervenes

## Links

### Out

- derived-from [m40-interventions-plan-research](../source/m40-interventions-plan-research.md)

### In

- decided-by [run-intervention-message](../capability/rhiza/workflows/workflow-runs/run-intervention-message.md)
- decided-by [run-intervention-resume](../capability/rhiza/workflows/workflow-runs/run-intervention-resume.md)
- decided-by [run-intervention-approve](../capability/rhiza/workflows/workflow-runs/run-intervention-approve.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

# intervention-one-entry-type

decision · `01M7R2FXN4BPJRD4ZYWG8FMD04`

## statement

Every intervention is one journal entry type, intervention, with a kind (message, resume, retry, abandon, approve, hand-commit), the actor who made it and a detail; the fold keeps them in order and rhiza run status --json shows the list and its count

## wording

(assumed, would have escalated) one entry type with a kind field, chosen by the conductor so status and the daemon count one list

## Links

### Out

- derived-from [m40-interventions-plan-research](../source/m40-interventions-plan-research.md)

### In

- decided-by [run-intervention-journal](../capability/rhiza/workflows/workflow-runs/run-intervention-journal.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

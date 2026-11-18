# hand-commit-by-trailer

decision · `01M7R2FXRAMCY0W76S1QG78Z7Y`

## statement

A commit on a run's branch is the run's when it is a merge of the base branch or already on the base branch, carries a Leaf: trailer naming a leaf of the run's change set, a Review: trailer naming the run, or an Actor: trailer naming the actor of one of the run's sessions or the actor the run's own process commits as (its command steps' render and land commits); any other commit added to the branch after the run started is a hand commit, journalled once per commit

## wording

(assumed, would have escalated) attribution by the trailers the run's roles already write

## Links

### Out

- derived-from [m40-interventions-plan-research](../source/m40-interventions-plan-research.md)

### In

- decided-by [run-intervention-hand-commit](../capability/rhiza/workflows/workflow-runs/run-intervention-hand-commit.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

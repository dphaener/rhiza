# run-intervention-hand-commit-qa-1

evidence · `01M7R8VX0ZMN8X6KN49TV154JM` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

FAIL

## statement

QA with a scratch build, two-wait workflow: Review: this run, Leaf: of the change set and a merge of main are not counted; plain, Review: other run and unknown Leaf: commits are journalled once each on entering the next step, a later one once at finish, no duplicates, run result unchanged. Major finding: the run's own 'Merge origin/<branch> into <branch>' (ship divergence merge) is counted as a hand commit; minors: a rebase onto main counts main's commits, an Actor: trailer naming the process's owner actor is excused.

## Links

### Out

- verifies [run-intervention-hand-commit](../capability/rhiza/workflows/workflow-runs/run-intervention-hand-commit.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

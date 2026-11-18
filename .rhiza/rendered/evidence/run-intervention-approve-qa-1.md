# run-intervention-approve-qa-1

evidence · `01M7R8VWXAB2FTS9G7RM7JJMPN` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

PASS

## statement

QA with a scratch build: owner rhiza accept of a change set whose live run has accept auto journals one approve ('accepted change set <id> at step hold'); an accept the run waits on (rhiza.autonomy.accept owner), an accept of a change set whose run is finished, a refused second accept and a refused --rule autonomous journal none.

## Links

### Out

- verifies [run-intervention-approve](../capability/rhiza/workflows/workflow-runs/run-intervention-approve.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

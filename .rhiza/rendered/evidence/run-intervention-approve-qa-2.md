# run-intervention-approve-qa-2

evidence · `01M7RAJ17HXC9KYZKMFJTH9ARJ` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2: owner accept of a change set with a live auto-accept run journals one approve at the run's step; dry run, refused repeat and a run whose workflow declares accept=owner journal none; run-intervention-approve script passes.

## Links

### Out

- verifies [run-intervention-approve](../capability/rhiza/workflows/workflow-runs/run-intervention-approve.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

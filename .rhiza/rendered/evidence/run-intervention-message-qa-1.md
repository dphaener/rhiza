# run-intervention-message-qa-1

evidence · `01M7R8VWMZNPN2N643WW9T4WBJ` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

PASS

## statement

QA with a scratch build: session message from outside any run to a run's session journals one message intervention (actor, session, first line); from a session of the same run (RHIZA_RUN+RHIZA_SESSION or RHIZA_SESSION alone) none; to a session of no run sent, none; dry run and blank text none. Rough edge (minor finding): a message to a finished run's session still journals on that run.

## Links

### Out

- verifies [run-intervention-message](../capability/rhiza/workflows/workflow-runs/run-intervention-message.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

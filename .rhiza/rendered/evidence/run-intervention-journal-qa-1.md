# run-intervention-journal-qa-1

evidence · `01M7R8VWHXYZYH7NY4GDXHEN44` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

PASS

## statement

QA with a scratch build in a $TMPDIR root: run status --json shows interventions [] and intervention_count 0 for a clean run, the list in journal order and the count for one with interventions, and text '- Interventions: 3 (resume 2, abandon 1)' plus an Interventions list; a hand-appended entry of kind bogus is not counted and run status warns 'unknown kind "bogus": not counted'.

## Links

### Out

- verifies [run-intervention-journal](../capability/rhiza/workflows/workflow-runs/run-intervention-journal.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

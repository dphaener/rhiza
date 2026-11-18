# daemon-status-interventions

capability · `01M7R2FYDRQGC4KBP7G5SXJ8XB` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscript cmd/rhiza/testdata/script/daemon-status-interventions.txtar: with three finished runs, oldest to newest intervened, clean, clean, finished 3, finished_clean 2 and clean_streak 2 show in --json and text, and a live run row carries its count; (failure case) a project with no finished runs shows 0 for all three and no error; describe-golden.txtar regenerated

## files

internal/api/daemon.go, cmd/rhiza/testdata/script/daemon-status-interventions.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change the daemon's scheduling or the existing daemon status fields

## statement

rhiza daemon status shows each live run's intervention count and, over the project's finished runs, finished, finished_clean (zero interventions) and clean_streak (consecutive finished runs with zero, newest first), in --json and as a text line

## Links

### Out

- depends-on [run-intervention-journal](run-intervention-journal.md)
- decided-by [daemon-counts-clean-streak](../../../../decision/daemon-counts-clean-streak.md)

### In

- verifies [daemon-status-interventions-verified](../../../../evidence/daemon-status-interventions-verified.md)
- verifies [daemon-status-interventions-qa-1](../../../../evidence/daemon-status-interventions-qa-1.md)
- verifies [daemon-status-interventions-qa-2](../../../../evidence/daemon-status-interventions-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

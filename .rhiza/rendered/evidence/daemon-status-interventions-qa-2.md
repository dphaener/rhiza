# daemon-status-interventions-qa-2

evidence · `01M7RAJ1E5JK6X357QZ2902ZNZ` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2: daemon status --json and text show finished 7, finished_clean 2, clean_streak 1 matching each run's status; after daemon --once with fake rex live rows carry intervention_count and ', N intervention'; a fresh project shows 0/0/0 without error.

## Links

### Out

- verifies [daemon-status-interventions](../capability/rhiza/workflows/workflow-runs/daemon-status-interventions.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

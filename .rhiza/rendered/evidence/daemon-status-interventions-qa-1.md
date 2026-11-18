# daemon-status-interventions-qa-1

evidence · `01M7R8VX3R6WQ2Z2HW1NW9JXJV` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

PASS

## statement

QA with a scratch build: daemon status --json and text show finished/finished_clean/clean_streak matching each run's status (9/3/0, then 10/4/1 after a clean newest finish); after rhiza daemon --once with detach faked, live rows carry intervention_count and text ', N interventions'; a fresh project shows 0/0/0 without error.

## Links

### Out

- verifies [daemon-status-interventions](../capability/rhiza/workflows/workflow-runs/daemon-status-interventions.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

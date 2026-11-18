# run-intervention-journal-qa-2

evidence · `01M7RAJ0RZ1GKZ8K3WRBQ9VZVT` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2 with a scratch build: run status --json shows intervention_count 0 and interventions [] on a fresh run, then the list in journal order; text shows '- Interventions: 3 (resume 2, abandon 1)' and the Interventions section; go test ./internal/run -run Intervention ok.

## Links

### Out

- verifies [run-intervention-journal](../capability/rhiza/workflows/workflow-runs/run-intervention-journal.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

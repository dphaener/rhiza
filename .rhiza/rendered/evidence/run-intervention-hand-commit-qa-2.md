# run-intervention-hand-commit-qa-2

evidence · `01M7RAJ1AZ9EV3121BTQFZXH6X` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2: a hand commit on a stopped run's branch is journalled once on entering the next step with its hash, the run succeeding unchanged; a merge of main adds none. A rebase of the branch re-journals the rewritten commit (minor finding f1).

## Links

### Out

- verifies [run-intervention-hand-commit](../capability/rhiza/workflows/workflow-runs/run-intervention-hand-commit.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

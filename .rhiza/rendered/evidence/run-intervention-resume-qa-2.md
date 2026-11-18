# run-intervention-resume-qa-2

evidence · `01M7RAJ132F0G98HEST13CR2G9` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2: person's queued and foreground resume, resume --retry and abandon each journal one matching intervention; dry runs and refused repeats none; RHIZA_DAEMON_RESUME resume none; run start on a failed change set journals a retry (foreground, and queued via retry_by/RHIZA_START_RETRY_BY), the daemon's own start none.

## Links

### Out

- verifies [run-intervention-resume](../capability/rhiza/workflows/workflow-runs/run-intervention-resume.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

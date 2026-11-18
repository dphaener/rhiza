# run-intervention-resume-qa-1

evidence · `01M7R8VWQMWT9NQAV5FBXW2227` · actor [claude-code-50f5e6e7](../actor/claude-code-50f5e6e7.md)

## result

FAIL

## statement

QA with a scratch build: person's queued and foreground resume, resume --retry and abandon each journal the matching intervention with actor and step; dry runs, refused resumes/abandons and RHIZA_DAEMON_RESUME (daemon/watchdog) journal none. Major finding: a person's rhiza run start of a change set with a failed run starts a new run and journals no retry intervention, contrary to the statement.

## Links

### Out

- verifies [run-intervention-resume](../capability/rhiza/workflows/workflow-runs/run-intervention-resume.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

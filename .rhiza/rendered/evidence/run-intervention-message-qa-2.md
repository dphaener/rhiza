# run-intervention-message-qa-2

evidence · `01M7RAJ0WPJ0PAT9B3F173XMWD` · actor [claude-code-884760eb](../actor/claude-code-884760eb.md)

## result

PASS

## statement

QA cycle 2: session message from outside to a run's session journals one message intervention (first line only, actor, step); from a sibling session (RHIZA_RUN+RHIZA_SESSION or RHIZA_SESSION alone), to a no-run session, --dry-run, blank text and unknown session journal none.

## Links

### Out

- verifies [run-intervention-message](../capability/rhiza/workflows/workflow-runs/run-intervention-message.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

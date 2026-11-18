# run-intervention-message

capability · `01M7R2FXXVADBMRDCBMH1JESQ1` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscript cmd/rhiza/testdata/script/run-intervention-message.txtar: a message from outside any run to a run's session shows one message intervention in rhiza run status --json; (failure case) a message sent with RHIZA_RUN and RHIZA_SESSION naming a session of the same run adds none, and a message to a session of no run journals nothing and still sends

## files

internal/api/session_message.go, internal/api/session_message_test.go, cmd/rhiza/testdata/script/run-intervention-message.txtar

## must_not

Change how the message is delivered; journal the engine's own hand-back, amend, lost-session or API-error pointer lines

## statement

rhiza session message to a session of a run journals a message intervention on that run (actor, target session, the text's first line) when the sender is not a session of the same run

## Links

### Out

- depends-on [run-intervention-journal](run-intervention-journal.md)
- decided-by [intervention-is-a-person-command](../../../../decision/intervention-is-a-person-command.md)

### In

- verifies [run-intervention-message-verified](../../../../evidence/run-intervention-message-verified.md)
- verifies [run-intervention-message-qa-1](../../../../evidence/run-intervention-message-qa-1.md)
- verifies [run-intervention-message-qa-2](../../../../evidence/run-intervention-message-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

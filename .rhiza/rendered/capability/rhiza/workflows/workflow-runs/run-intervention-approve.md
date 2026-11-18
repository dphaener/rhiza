# run-intervention-approve

capability · `01M7R2FY5HR9X0ZJRR595VFVPE` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscript cmd/rhiza/testdata/script/run-intervention-approve.txtar: an owner accept of a change set whose run is live under auto accept shows one approve intervention in rhiza run status --json; (failure case) rhiza run approve at a waiting signoff step leaves intervention_count 0, and an accept of a change set with no live run journals nothing

## files

internal/api/lifecycle.go, internal/api/run_signoff.go, cmd/rhiza/testdata/script/run-intervention-approve.txtar

## must_not

Change accept, reject or approve behaviour; journal the workflow's own autonomous gate accepts

## statement

An owner-rule accept or reject of a change set while a live run of it is not waiting for that owner accept (its accept stage is not owner) journals an approve intervention on that run; rhiza run approve at a signoff step the run waits on journals none

## Links

### Out

- depends-on [run-intervention-journal](run-intervention-journal.md)
- decided-by [intervention-is-a-person-command](../../../../decision/intervention-is-a-person-command.md)

### In

- verifies [run-intervention-approve-verified](../../../../evidence/run-intervention-approve-verified.md)
- verifies [run-intervention-approve-qa-1](../../../../evidence/run-intervention-approve-qa-1.md)
- verifies [run-intervention-approve-qa-2](../../../../evidence/run-intervention-approve-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

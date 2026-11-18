# run-intervention-journal

capability · `01M7R2FXVTKFFG30W56PYYSGHW` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test ./internal/run -run Intervention: a journal with two intervention entries folds to both in order, count 2; a journal with none folds to count 0; (failure case) an intervention entry with an unknown kind is refused by append or reported by the fold, never counted silently. Testscript cmd/rhiza/testdata/script/run-status-interventions.txtar shows the count and list in rhiza run status --json and the text line; describe-golden.txtar regenerated

## files

internal/run/journal.go, internal/run/journal_test.go, internal/api/run.go, cmd/rhiza/testdata/script/run-status-interventions.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Record any intervention from a command yet (later leaves do); change any other entry type's fold or the run's state or result

## statement

The run journal has an intervention entry (kind, actor, detail, optional commit) folded into the run's state in order, and rhiza run status shows them: --json an interventions list and an intervention_count (0 when none), text a line counting them by kind

## Links

### Out

- decided-by [intervention-one-entry-type](../../../../decision/intervention-one-entry-type.md)

### In

- depends-on [run-intervention-message](run-intervention-message.md)
- depends-on [run-intervention-resume](run-intervention-resume.md)
- depends-on [run-intervention-approve](run-intervention-approve.md)
- depends-on [run-intervention-hand-commit](run-intervention-hand-commit.md)
- depends-on [daemon-status-interventions](daemon-status-interventions.md)
- verifies [run-intervention-journal-verified](../../../../evidence/run-intervention-journal-verified.md)
- verifies [run-intervention-journal-qa-1](../../../../evidence/run-intervention-journal-qa-1.md)
- verifies [run-intervention-journal-qa-2](../../../../evidence/run-intervention-journal-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted

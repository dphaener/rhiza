# change-new-workflow-flag

capability · `01M71G7NY758150RA00YZ080MV` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript (go test ./cmd/rhiza -run TestScripts/<its name>) opens a change set with change new --workflow, shows the workflow in change show --json, shows run start --dry-run picking it, shows an unknown name refused, and folds with fold-ok

## statement

rhiza change new --workflow <name> opens a change set pinned to that workflow, recorded as a request's workflow is, so its ready pick is that workflow; a name that does not resolve is refused as rhiza request --workflow refuses it

## Links

### Out

- decided-by [m31-polish-request](../../../../decision/m31-polish-request.md)

### In

- verifies [change-new-workflow-flag-verified](../../../../evidence/change-new-workflow-flag-verified.md)
- verifies [change-new-workflow-flag-qa-1](../../../../evidence/change-new-workflow-flag-qa-1.md)
- verifies [pr-425-docs-sign-off-workflow](../../../../evidence/pr-425-docs-sign-off-workflow.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow

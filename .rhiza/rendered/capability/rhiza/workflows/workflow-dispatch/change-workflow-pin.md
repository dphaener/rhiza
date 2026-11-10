# change-workflow-pin

capability · `01M71G7P0PSQXY5RZ89TE5NV9W` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

A testscript pins a draft change set to a workflow and shows change show --json and run start --dry-run picking it, clears it back to the dispatch rules, refuses an unknown name and a landed change set, and folds with fold-ok; internal/api tests cover the operation

## statement

A command pins or clears the workflow of an existing change set that has not landed and has no live run (for example rhiza change workflow <name> | --clear), so a change set opened by a milestone plan can be built with a chosen workflow

## Links

### Out

- decided-by [m31-polish-request](../../../../decision/m31-polish-request.md)
- decided-by [pin-format-five](../../../../decision/pin-format-five.md)

### In

- depends-on [docs-signoff-workflow](../../documentation/docs-signoff-workflow.md)
- verifies [change-workflow-pin-verified](../../../../evidence/change-workflow-pin-verified.md)
- verifies [change-workflow-pin-qa-1](../../../../evidence/change-workflow-pin-qa-1.md)
- verifies [change-workflow-pin-verified-2](../../../../evidence/change-workflow-pin-verified-2.md)
- verifies [pr-425-docs-sign-off-workflow](../../../../evidence/pr-425-docs-sign-off-workflow.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow

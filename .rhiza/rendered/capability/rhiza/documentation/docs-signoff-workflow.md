# docs-signoff-workflow

capability · `01M71G7P1XK0DKX0N1JXMP8T4B` · parent [documentation](../documentation.md)

## acceptance

rhiza workflow check docs passes, rhiza workflow describe docs shows signoff owner, change show --json of each of the three page change sets shows workflow docs, and make self-check passes

## statement

This repo has a workflow named docs (.rhiza/workflows/docs.toml, forked from ship-change) whose autonomy sets signoff = owner, whose builder and fixer prompts write doc prose with the owner's humanize skill under constraints doc-prose-humanized and doc-examples-plain-commands, and the milestone's three page change sets (README and quickstart polish, Concepts polish, Agents and runs docs polish) are pinned to it

## Links

### Out

- depends-on [change-workflow-pin](../workflows/workflow-dispatch/change-workflow-pin.md)
- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [docs-signoff-workflow-verified](../../../evidence/docs-signoff-workflow-verified.md)
- verifies [docs-signoff-workflow-qa-1](../../../evidence/docs-signoff-workflow-qa-1.md)
- verifies [pr-425-docs-sign-off-workflow](../../../evidence/pr-425-docs-sign-off-workflow.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow

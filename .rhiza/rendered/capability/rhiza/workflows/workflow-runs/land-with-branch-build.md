# land-with-branch-build

capability · `01M5X3GC9S0ENG89Y22HF3KPFH` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript where the branch's rendering differs from the installed binary's: after ship land, rhiza check with the head's build passes (no stale rendered files).

## statement

In Rhiza's own repo (its go.mod module is Rhiza's), ship-change's land and render use a rhiza built from the PR head, so a branch that changes how the model renders lands with its own rendering.

## Links

### Out

- decided-by [ship-change-holds-up](../../../../decision/ship-change-holds-up.md)

### In

- verifies [land-with-branch-build-tests](../../../../evidence/land-with-branch-build-tests.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.

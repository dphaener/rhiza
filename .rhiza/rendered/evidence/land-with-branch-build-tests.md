# land-with-branch-build-tests

evidence · `01M5X6XT16R7JVD1GDGR2SC9QD` · actor [claude-code-eac4d56d-verifier](../actor/claude-code-eac4d56d-verifier.md)

## result

go test ./cmd/rhiza -run TestScripts/ship-land-branch-build -count=1 passed at 35733c1

## statement

Ship land builds and renders with the PR head's build in Rhiza's repo.

## status

verified

## tests

TestScripts/ship-land-branch-build

## Links

### Out

- verifies [land-with-branch-build](../capability/rhiza/workflows/workflow-runs/land-with-branch-build.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.

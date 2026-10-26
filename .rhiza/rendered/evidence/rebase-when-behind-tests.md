# rebase-when-behind-tests

evidence · `01M5X6XTQH9PA705SC2JRNBPWR` · actor [claude-code-eac4d56d-verifier](../actor/claude-code-eac4d56d-verifier.md)

## result

go test ./cmd/rhiza -run TestScripts/run-ship-change-rebase -count=1 passed at 35733c1

## statement

Ship rebases a PR that main has moved past, and asks when code conflicts.

## status

verified

## tests

TestScripts/run-ship-change-rebase

## Links

### Out

- verifies [rebase-when-behind](../capability/rhiza/workflows/workflow-runs/rebase-when-behind.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.

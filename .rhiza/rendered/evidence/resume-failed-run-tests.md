# resume-failed-run-tests

evidence · `01M5X6XRM1WM0DBX63V4TWK84R` · actor [claude-code-eac4d56d-verifier](../actor/claude-code-eac4d56d-verifier.md)

## result

go test ./cmd/rhiza -run TestScripts/run-resume-retry -count=1 passed at 35733c1

## statement

resume --retry restarts a failed run at its failed step.

## status

verified

## tests

TestScripts/run-resume-retry

## Links

### Out

- verifies [resume-failed-run](../capability/rhiza/workflows/workflow-runs/resume-failed-run.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.

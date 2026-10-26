# agent-step-done-signal-tests

evidence · `01M5X6XSACCQ2CDD8APT4ANWJY` · actor [claude-code-eac4d56d-verifier](../actor/claude-code-eac4d56d-verifier.md)

## result

go test ./cmd/rhiza -run TestScripts/run-agent-done -count=1 passed at 35733c1

## statement

An agent step with until_when done passes only once the agent says it is done.

## status

verified

## tests

TestScripts/run-agent-done

## Links

### Out

- verifies [agent-step-done-signal](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.

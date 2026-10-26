# agent-step

capability · `01M5TYGHAHETP91SS0S0Y8BPYB` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A testscript with a fake adapter runs an agent step to its outcome.

## statement

An agent step launches its agent through the adapter on the run's change set and ends with the session's outcome.

## Links

### Out

- decided-by [engine-orchestrates](../../../../decision/engine-orchestrates.md)
- decided-by [five-step-kinds](../../../../decision/five-step-kinds.md)

### In

- relates-to [agent-step-prompt](../../../../question/agent-step-prompt.md) — the agent step launches its agent with no prompt today
- verifies [agent-step-tested](../../../../evidence/agent-step-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

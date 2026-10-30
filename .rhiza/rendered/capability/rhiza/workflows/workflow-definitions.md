# workflow-definitions

capability · `01M5SWT9A05VXSMXQME4CPJN2D` · parent [workflows](../workflows.md)

## statement

A workflow is a declarative state machine of steps (agent, command, wait, workflow, foreach) with outcomes and visit limits.

## Links

### Out

- decided-by [five-step-kinds](../../../decision/five-step-kinds.md)
- decided-by [declarative-toml](../../../decision/declarative-toml.md)

### In

- constrains [workflow-runs-another-once](../../../constraint/workflow-runs-another-once.md)
- constrains [workflow-step-no-session-keys](../../../constraint/workflow-step-no-session-keys.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

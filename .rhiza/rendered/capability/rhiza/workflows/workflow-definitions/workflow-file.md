# workflow-file

capability · `01M5TYGH047RKA27R0RBS38HKM` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A test loads a valid file and rejects both bad ones with a message.

## statement

A workflow is a TOML file of named steps, each of kind agent, wait or command, with outcomes naming the next step; an unknown kind or a missing target is refused.

## Links

### Out

- decided-by [declarative-toml](../../../../decision/declarative-toml.md)
- decided-by [five-step-kinds](../../../../decision/five-step-kinds.md)

### In

- verifies [workflow-file-tested](../../../../evidence/workflow-file-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

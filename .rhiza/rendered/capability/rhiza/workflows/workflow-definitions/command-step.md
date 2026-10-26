# command-step

capability · `01M5TYGHN9QSF654FF0V1Q309N` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A testscript runs a command step that exits 0 and one that exits 1 and follows each outcome.

## statement

A command step runs a shell command in the worktree and maps its exit code to an outcome.

## Links

### Out

- decided-by [declarative-toml](../../../../decision/declarative-toml.md)

### In

- verifies [command-step-tested](../../../../evidence/command-step-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

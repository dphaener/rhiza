# command-step

capability · `01M5TYGHN9QSF654FF0V1Q309N` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A testscript runs a command step that exits 0 and one that exits 1 and follows each outcome; a Go test in internal/run shows a command exiting 75 stops the run with its stop: detail and a resume reruns the step

## statement

A command step runs a shell command in the worktree and maps its exit code to an outcome; exit 75 is not an outcome: it stops the run for the owner with the command's last stop: line as the detail, the step is not ended, and resuming reruns it

## Links

### Out

- decided-by [declarative-toml](../../../../decision/declarative-toml.md)

### In

- verifies [command-step-tested](../../../../evidence/command-step-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
- `01M6AG8XH8SVYVXSHFFFMGPQE2` Plan milestone 11: owner gates hold, picked by the orchestrator under the owner's standing instruction

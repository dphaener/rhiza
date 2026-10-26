# wait-step

capability · `01M5TYGHWFYX6PYBNED6R4X1M4` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A test with a fake clock reaches both the held and the timed-out outcome.

## statement

A wait step polls its condition inside the engine until it holds or its timeout passes.

## Links

### Out

- decided-by [five-step-kinds](../../../../decision/five-step-kinds.md)

### In

- verifies [wait-step-tested](../../../../evidence/wait-step-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

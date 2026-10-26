# max-visits

capability · `01M5TYGJ3TN69C01G7J9C3NJWY` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

A testscript loops a step past its cap and sees the run fail.

## statement

A step's `max_visits` caps how often a run enters it, and going over fails the run with the step named.

## Links

### Out

- decided-by [five-step-kinds](../../../../decision/five-step-kinds.md)

### In

- verifies [max-visits-tested](../../../../evidence/max-visits-tested.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.

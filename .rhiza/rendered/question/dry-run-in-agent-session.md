# dry-run-in-agent-session

question · `01M6505NR5JYRG281YJB4AX2JE`

## statement

Should rhiza run start --dry-run also refuse inside an agent session, or only a real start, since a dry run writes and starts nothing?

## Links

### Out

- blocks [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

### In

- answers [run-start-dry-run-allowed](../decision/run-start-dry-run-allowed.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

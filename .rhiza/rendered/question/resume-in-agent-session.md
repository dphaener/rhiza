# resume-in-agent-session

question · `01M6505PXJA79ABGX1DR950R8N`

## statement

Should rhiza run resume also refuse inside an agent session, as run start does, since resuming starts a run's next step?

## Links

### Out

- blocks [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

### In

- answers [run-resume-refuses-in-agent](../decision/run-resume-refuses-in-agent.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

# no-agent-session-starts-runs

constraint · `01M6AXQJNPM78KBAPJHASHJRPY`

## acceptance

Tests show rhiza run start and rhiza run resume both refuse when called from inside an agent session

## statement

No agent session starts a run.

## Links

### Out

- constrains [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)
- constrains [run-resume-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-resume/run-resume-refuses-agent-session.md)

## History

- `01M65G0WKX7Y74QCRB8EAN573P` A repo constraint states that no agent session starts a run, so every run-starting path (run start, run resume, a future daemon) is checked against it, not only rhiza run start

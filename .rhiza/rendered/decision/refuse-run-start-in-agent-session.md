# refuse-run-start-in-agent-session

decision · `01M6505EZA0RTXD77K309NDRT6`

## statement

rhiza run start refuses inside an agent session, when RHIZA_SESSION is set, naming that session

## wording

Refuse run start inside an agent session (RHIZA_SESSION set), naming the session: the change set's intent, the agent's wording (assumed, would have escalated)

## Links

### Out

- derived-from [agents-never-start-runs-plan](../source/agents-never-start-runs-plan.md)

### In

- decided-by [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

# run-start-dry-run-allowed

decision · `01M65F4F8716WGMAEHPGPVKT3Z`

## statement

rhiza run start --dry-run stays allowed inside an agent session; only a real start refuses, since a dry run writes and starts nothing

## wording

Arbiter answer under the standing autonomous instruction: a dry run starts nothing, so it is not refused

## Links

### Out

- answers [dry-run-in-agent-session](../question/dry-run-in-agent-session.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

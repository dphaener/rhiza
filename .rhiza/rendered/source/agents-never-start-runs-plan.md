# agents-never-start-runs-plan

source · `01M65056W4Z6BXNB22RBBYAMQN`

## quotes

```json
[]
```

## raw_sha256

9b10aa8ee30c0f772ad7b0cd689a3bb6f026d53e261382508b5df67b65003c23

## session

2aec8012-a240-4729-976b-01a18a012a74

## summary

Planned change set 01M61GW73 (Agents never start runs) in an autonomous plan run, so no owner was in the loop and no owner words were typed. The intent: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh. Settled from the intent: the refusal keys on RHIZA_SESSION and names the session; the fakes-only rule lives in the build and review guides. Left open, parked as questions: whether --dry-run refuses too, and whether run resume refuses too.

## title

Agents never start runs: plan, autonomous

## Links

### In

- derived-from [refuse-run-start-in-agent-session](../decision/refuse-run-start-in-agent-session.md)
- derived-from [fakes-only-in-build-and-review](../decision/fakes-only-in-build-and-review.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

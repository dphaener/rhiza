# pr-104-agents-never-start-runs

evidence · `01M65HHJTXPKSTDVD0GQ090JMY` · actor [darin-haener](../actor/darin-haener.md)

## head

223505e7bae97b2e2a7d1b27b355dc09948009d9

## pr

https://github.com/dphaener/rhiza/pull/104

## result

1 check passed: check SUCCESS

## statement

Pull request #104 (https://github.com/dphaener/rhiza/pull/104) passed CI at 223505e

## status

verified

## Links

### Out

- verifies [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)
- verifies [build-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-guide-fakes-only.md)
- verifies [review-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

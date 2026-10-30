# build-guide-fakes-only

capability · `01M6505JJY04DWHQZ3EQRFCHAM` · parent [build-guide](../build-guide.md)

## acceptance

A Go guide test finds the fakes-only rule in the build topic, naming Rex, agents and gh, and the topic stays within TopicCap

## statement

rhiza guide build tells verifiers to test runs with fakes only, never real Rex, agents or gh

## Links

### Out

- decided-by [fakes-only-in-build-and-review](../../../../../../decision/fakes-only-in-build-and-review.md)

### In

- verifies [build-guide-fakes-only-verified](../../../../../../evidence/build-guide-fakes-only-verified.md)
- verifies [build-guide-fakes-only-qa-1](../../../../../../evidence/build-guide-fakes-only-qa-1.md)
- verifies [pr-104-agents-never-start-runs](../../../../../../evidence/pr-104-agents-never-start-runs.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.

# review-guide-fakes-only

capability · `01M6505MN5E02V3Y30RMC16FN5` · parent [review-guide](../review-guide.md)

## acceptance

A Go guide test finds the fakes-only rule in the review topic for the compliance, quality, scope and model checks, naming Rex, agents and gh; it finds that qa is hands-on with the built binary and that fakes alone do not count; and the topic stays within TopicCap

## statement

rhiza guide review tells each check except qa to test runs with fakes only, never real Rex, agents or gh; qa runs the built binary hands-on, since fakes alone do not count for it

## Links

### Out

- decided-by [fakes-only-in-build-and-review](../../../../../../decision/fakes-only-in-build-and-review.md)

### In

- verifies [review-guide-fakes-only-verified](../../../../../../evidence/review-guide-fakes-only-verified.md)
- verifies [review-guide-fakes-only-qa-1](../../../../../../evidence/review-guide-fakes-only-qa-1.md)
- verifies [pr-104-agents-never-start-runs](../../../../../../evidence/pr-104-agents-never-start-runs.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
